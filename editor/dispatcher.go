package main

import (
	"net/http"
	"strconv"
	"strings"
)

func init() {
	http.HandleFunc("/", dispatch)
}

func dispatch(w http.ResponseWriter, r *http.Request) {
	path := strings.FieldsFunc(r.URL.Path, func(r rune) bool {
		return r == '/'
	})

	switch len(path) {
	case 0:
		home(w)
		return

	case 2:
		switch path[0] {
		case "level":
			l, err := strconv.ParseUint(path[1], 10, 64)
			if err != nil {
				break
			}
			level(w, int(l))
			return
		}
	}

	http.NotFound(w, r)
}
