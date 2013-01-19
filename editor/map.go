package main

import (
	"fmt"
	"github.com/Nightgunner5/4407/server/matter"
	"net/http"
)

func levelCount(w http.ResponseWriter) {
	State.RLock()
	defer State.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	_, err := fmt.Fprint(w, len(State.Map))
	handle(err)
}

func setTile(w http.ResponseWriter, z, x, y int64, t uint32) {
	State.Lock()
	defer State.Unlock()

	State.Map[z].Layout[matter.Coord{x, y}] = matter.LayoutTile(t)
}

func newLevel(w http.ResponseWriter) {
	State.Lock()
	defer State.Unlock()

	w.Header().Set("Content-Type", "application/json")
	State.Map.NewLevel()
	_, err := fmt.Fprint(w, len(State.Map)-1)
	handle(err)
}

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
