package main

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"github.com/Nightgunner5/4407/server/matter"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"runtime/pprof"
)

func ReadMap(r io.Reader) (matter.Map, error) {
	g, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer g.Close()

	var m matter.Map
	err = gob.NewDecoder(g).Decode(&m)
	return m, err
}

func clamp(f float64) uint8 {
	if f >= 255 {
		return 255
	}
	if f <= 0 {
		return 0
	}
	return uint8(f)
}

func main() {
	f, err := os.Open("map.gz")
	if err != nil {
		panic(err)
	}
	m, err := ReadMap(f)
	f.Close()
	if err != nil {
		panic(err)
	}

	const (
		tileMax = 4
	)
	var tileicon [tileMax]image.Image
	for i := range tileicon {
		f, err := os.Open(fmt.Sprintf("tile-%d.png", i))
		if err != nil {
			panic(err)
		}

		tileicon[i], err = png.Decode(f)
		f.Close()
		if err != nil {
			panic(err)
		}
	}

	f, _ = os.Create("cpu.prof")
	defer f.Close()
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	img := image.NewRGBA(image.Rect(int(m[0].Min.X)<<2, int(m[0].Min.Y)<<2, int(m[0].Max.X)<<2, int(m[0].Max.Y)<<2))
	const (
		tempMax = 300
		airMax  = 110
	)
	for i := 0; i < 1500; i++ {
		for _, t := range m[0].Atmos {
			r := image.Rect(int(t.X<<2), int(t.Y<<2), int(t.X<<2)+4, int(t.Y<<2)+4)
			draw.Draw(img, r, tileicon[m[0].Layout[t.Coord]], image.ZP, draw.Src)
			overlay := image.NewUniform(color.NRGBA{
				clamp(t.Temp * 255 / tempMax),
				clamp(t.Total() * 255 / airMax),
				0,
				clamp((t.Total() + 5) * 200 / airMax),
			})
			draw.Draw(img, r, overlay, image.ZP, draw.Over)
		}

		f, err := os.Create(fmt.Sprintf("atmos-%04d.png", i))
		if err != nil {
			panic(err)
		}

		err = png.Encode(f, img)
		f.Close()
		if err != nil {
			panic(err)
		}

		fmt.Println(i)

		m[0].Atmos.Tick()
	}
}
