package shared

import (
	"github.com/Nightgunner5/4407/matter"
)

type transferMap struct {
	X, Y int64
	Turf matter.LayoutTileTurf
	Area string
}

func Level(l matter.Layout) struct{ Map []transferMap } {
	var s struct{ Map []transferMap }
	for c, t := range l {
		if !t.Empty() {
			s.Map = append(s.Map, transferMap{c.X, c.Y, t.Turf, t.Area})
		}
	}
	return s
}
