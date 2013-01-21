package main

import (
	"compress/gzip"
	"encoding/gob"
	"github.com/Nightgunner5/4407/server/matter"
	"io"
	"net/http"
	"os"
	"runtime/pprof"
	"sync"
)

func ReadMap(r io.Reader) (matter.Map, error) {
	g, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer g.Close()

	var m matter.Map
	err = gob.NewDecoder(g).Decode(&m)
	return m, err
}

var State struct {
	M matter.Map
	sync.RWMutex
}

func main() {
	f, err := os.Open("map.gz")
	if err != nil {
		panic(err)
	}
	m, err := ReadMap(f)
	f.Close()
	if err != nil {
		panic(err)
	}
	State.M = m

	f, _ = os.Create("cpu.prof")
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	go atmosphere()

	err = http.ListenAndServe(":4407", nil)
	if err != nil {
		panic(err)
	}
}
