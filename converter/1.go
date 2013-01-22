package main

import (
	"compress/gzip"
	"encoding/gob"
	"github.com/Nightgunner5/4407/matter"
	"os"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var Import []struct {
		Layout map[matter.Coord]matter.LayoutTileTurf
	}
	{
		f, err := os.Open("../server/map.gz")
		handle(err)
		g, err := gzip.NewReader(f)
		handle(err)
		handle(gob.NewDecoder(g).Decode(&Import))
		g.Close()
		f.Close()
	}

	var Export matter.Map
	for _, l := range Import {
		ll := Export.NewLevel()
		for c, t := range l.Layout {
			ll[c] = matter.LayoutTile{Turf: t}
		}
	}
	Export.Compile(8)

	{
		f, err := os.Create("../server/map.gz")
		handle(err)
		g := gzip.NewWriter(f)
		handle(gob.NewEncoder(g).Encode(Export))
		g.Close()
		f.Close()
	}
}
