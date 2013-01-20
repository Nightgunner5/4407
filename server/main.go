package main

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"github.com/Nightgunner5/4407/server/matter"
	"io"
	"io/ioutil"
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

var tileicon [matter.TileCount][]byte

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

	for i := range tileicon {
		f, err := ioutil.ReadFile(fmt.Sprintf("tile-%d.png", i))
		if err != nil {
			panic(err)
		}
		tileicon[i] = f
	}

	for i := range State.M {
		http.HandleFunc(fmt.Sprintf("/level/%d", i), level(i))
	}

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
