package main

import (
	"fmt"
	"net/http"
)

func level(i int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		State.RLock()
		defer State.RUnlock()

		l := State.M[i]
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, "[")
		first := true
		for c, t := range l.Layout {
			if t != 0 {
				if first {
					first = false
				} else {
					fmt.Fprint(w, ",")
				}
				fmt.Fprintf(w, "[%d,%d,%d]", c.X, c.Y, t)
			}
		}
		fmt.Fprint(w, "]")
	}
}
