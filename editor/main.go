package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"github.com/Nightgunner5/4407/server/matter"
	"io"
	"os"
	"strings"
)

func parseInstance(b []byte, current matter.LayoutTile) (ret matter.LayoutTile) {
	defer func() {
		if current > ret {
			ret = current
		}
	}()

	if i := bytes.IndexByte(b, '{'); i != -1 {
		b = b[:i]
	}
	path := string(b)
	switch {
	case strings.HasPrefix(path, "/obj/window"):
		return matter.Window
	case strings.HasPrefix(path, "/turf/simulated/floor"),
		strings.HasPrefix(path, "/turf/unsimulated/floor"):
		return matter.Floor
	case strings.HasPrefix(path, "/turf/simulated/wall"),
		strings.HasPrefix(path, "/turf/unsimulated/wall"):
		return matter.Wall
	case path == "/turf/space":
		return matter.Space
	case path == "/turf/simulated/shuttle/wall":
		return matter.Wall
	case path == "/turf/simulated/shuttle/floor", path == "/turf/simulated":
		return matter.Floor
	case strings.HasPrefix(path, "/area"),
		strings.HasPrefix(path, "/mob"),
		strings.HasPrefix(path, "/obj"):
		fmt.Println(path)
		return current
	default:
		panic(path)
	}
	panic("unreachable")
}

func Parse(r io.Reader) matter.Map {
	stage := 0
	in := bufio.NewReader(r)
	var m matter.Map
	types := make(map[string]matter.LayoutTile)

	var x, y int64
	var cz matter.Layout
	var keyLen int

	for {
		line, err := in.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
		line = line[:len(line)-1]
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		switch {
		case stage == 0 && len(line) > 8 && line[0] == '"' && line[len(line)-1] == ')':
			for i := 1; i < len(line) && line[i] != '"'; i++ {
				keyLen++
			}
			stage++
			fallthrough
		case stage == 1 && len(line) > 8 && line[0] == '"' && line[len(line)-1] == ')':
			if !(line[keyLen+1] == '"' && line[keyLen+2] == ' ' && line[keyLen+3] == '=' && line[keyLen+4] == ' ' && line[keyLen+5] == '(') {
				return nil
			}
			var tt matter.LayoutTile
			key := string(line[1 : keyLen+1])
			line = line[keyLen+6 : len(line)-1]

			inExtra := false
			for i := 0; i < len(line); i++ {
				if !inExtra && line[i] == ',' {
					tt = parseInstance(line[0:i], tt)
					line = line[i+1:]
					i = 0
				}
				if line[i] == '{' {
					inExtra = true
				}
				if line[i] == '}' {
					inExtra = false
				}
			}
			tt = parseInstance(line, tt)
			types[key] = tt

		case stage == 1 && len(line) == 0:
			stage++

		case stage == 4 && len(line) == 0:
			cz = nil

		case (stage == 2 || stage == 4) && len(line) > 11 && line[0] == '(' &&
			line[len(line)-1] == '"' && line[len(line)-2] == '{' &&
			line[len(line)-3] == ' ' && line[len(line)-4] == '=' &&
			line[len(line)-5] == ' ' && line[len(line)-6] == ')':
			stage = 3

			cz = m.NewLevel()
			x, y = 0, 0

		case stage == 3 && len(line) == 2 && line[0] == '"' && line[1] == '}':
			stage++

		case stage == 3 && len(line)%keyLen == 0:
			x = 0
			for i := 0; len(line) != 0; i, line = i+1, line[keyLen:] {
				cz[matter.Coord{x, y}] = types[string(line[:keyLen])]
				x++
			}
			y++

		default:
			return nil
		}
	}

	if stage != 4 {
		return nil
	}

	return m
}

func main() {
	f, _ := os.Open("../server/maps/trunkmap.dmm")
	defer f.Close()
	m := Parse(f)

	m.Compile(0)

	f, _ = os.Create("../server/map.gz")
	defer f.Close()
	g := gzip.NewWriter(f)
	defer g.Close()
	w := gob.NewEncoder(g)
	w.Encode(m)
}
