package main

import (
	"fmt"
	"net/http"
)

func level(w http.ResponseWriter, i int) {
	State.RLock()
	defer State.RUnlock()

	l := State.Map[i]
	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(w, "[")
	handle(err)
	first := true
	for c, t := range l.Layout {
		if t != 0 {
			if first {
				first = false
			} else {
				_, err = fmt.Fprint(w, ",")
				handle(err)
			}
			_, err = fmt.Fprintf(w, "[%d,%d,%d]", c.X, c.Y, t)
			handle(err)
		}
	}
	_, err = fmt.Fprint(w, "]")
	handle(err)
}
