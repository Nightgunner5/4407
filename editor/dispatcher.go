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

	case 1:
		switch path[0] {
		case "levels":
			levelCount(w)
			return

		case "save":
			save()
			return
		}

	case 2:
		switch path[0] {
		case "level":
			l, err := strconv.ParseInt(path[1], 10, 64)
			if err != nil {
				break
			}
			level(w, int(l))
			return

		case "levels":
			switch path[1] {
			case "new":
				newLevel(w)
				return
			}
		}

	case 5:
		switch path[0] {
		case "set":
			z, err := strconv.ParseInt(path[1], 10, 64)
			if err != nil {
				break
			}
			x, err := strconv.ParseInt(path[2], 10, 64)
			if err != nil {
				break
			}
			y, err := strconv.ParseInt(path[3], 10, 64)
			if err != nil {
				break
			}
			t, err := strconv.ParseUint(path[4], 10, 32)
			if err != nil {
				break
			}
			setTile(w, z, x, y, uint32(t))
			return
		}
	}

	http.NotFound(w, r)
}
