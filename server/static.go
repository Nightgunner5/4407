package main

import (
	"fmt"
	"github.com/Nightgunner5/4407/server/matter"
	"net/http"
)

func init() {
	for i := 0; i < int(matter.TileCount); i++ {
		http.HandleFunc(fmt.Sprintf("/tile/%d.png", i), tileHandler(i))
	}
}

func tileHandler(i int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(tileicon[i])
	}
}
