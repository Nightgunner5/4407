package main

import (
	"github.com/Nightgunner5/4407/matter"
)

type transferMap struct {
	X, Y int64
	Turf matter.LayoutTileTurf
}

func level(i int) struct{ Map []transferMap } {
	State.RLock()
	defer State.RUnlock()

	l := State.M[i]
	var s struct{ Map []transferMap }
	for c, t := range l.Layout {
		if t != (matter.LayoutTile{Turf: matter.Space}) {
			s.Map = append(s.Map, transferMap{c.X, c.Y, t.Turf})
		}
	}
	return s
}
