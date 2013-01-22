package main

import (
	"github.com/Nightgunner5/4407/matter"
)

func level(i int) struct{ Map [][3]int64 } {
	State.RLock()
	defer State.RUnlock()

	l := State.M[i]
	var s struct{ Map [][3]int64 }
	for c, t := range l.Layout {
		if t != (matter.LayoutTile{Turf: matter.Space}) {
			s.Map = append(s.Map, [3]int64{c.X, c.Y, int64(t.Turf)})
		}
	}
	return s
}
