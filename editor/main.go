package main

import (
	"compress/gzip"
	"encoding/gob"
	"github.com/Nightgunner5/4407/server/matter"
	"os"
)

func C(x, y int) matter.Coord {
	return matter.Coord{int64(x), int64(y)}
}

func main() {
	f, _ := os.Create("map.gz")
	defer f.Close()

	g := gzip.NewWriter(f)
	defer g.Close()

	w := gob.NewEncoder(g)

	var m matter.Map
	l := m.NewLevel()
	for x := -10; x <= 10; x++ {
		for y := -10; y <= 10; y++ {
			if x == 10 && -3 < y && y < 3 {
				l[C(x, y)] = matter.Floor
			} else if x == -10 || x == 10 || y == -10 || y == 10 {
				l[C(x, y)] = matter.Wall
			} else {
				l[C(x, y)] = matter.Floor
			}
		}
	}

	m.Compile(64)

	w.Encode(m)
}
