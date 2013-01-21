package main

import (
	"fmt"
	"github.com/Nightgunner5/4407/matter"
	"io/ioutil"
	"net/http"
)

func init() {
	for i := 0; i < int(matter.TileCount); i++ {
		http.HandleFunc(fmt.Sprintf("/tile/%d.png", i), staticHandler(fmt.Sprintf("tile-%d.png", i)))
	}

	http.HandleFunc("/icon/status-cond.png", staticHandler("status-cond.png"))
}

func staticHandler(fn string) func(http.ResponseWriter, *http.Request) {
	file, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(file)
	}
}
