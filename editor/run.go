package main

import (
	"compress/gzip"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/Nightgunner5/4407/matter"
	"net"
	"net/http"
	"os"
	"sync"
)

var (
	startnew = flag.Bool("new", false, "start with an empty map")
	filename = flag.String("file", "../server/map.gz", "the file to edit")
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func read() (m matter.Map) {
	f, err := os.Open(*filename)
	handle(err)
	defer f.Close()
	g, err := gzip.NewReader(f)
	handle(err)
	defer g.Close()
	r := gob.NewDecoder(g)
	err = r.Decode(&m)
	handle(err)
	return
}

func save() {
	State.Lock()
	defer State.Unlock()

	State.Compile(16)

	f, err := os.Create(*filename)
	handle(err)
	defer f.Close()
	g := gzip.NewWriter(f)
	defer g.Close()
	r := gob.NewEncoder(g)
	err = r.Encode(State.Map)
	handle(err)
	fmt.Println("Map saved to ", *filename)
}

var State struct {
	matter.Map
	sync.RWMutex
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Fatal error:", r.(error))
			os.Exit(1)
		}
	}()

	flag.Parse()

	if !*startnew {
		fmt.Println("Reading map from", *filename)
		State.Map = read()
		fmt.Println("Map read successfully")
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	handle(err)
	defer ln.Close()

	fmt.Println("Starting editor web server at", "http://"+ln.Addr().String())

	handle(http.Serve(ln, nil))
}
