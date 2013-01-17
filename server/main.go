package main

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"github.com/Nightgunner5/4407/server/matter"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
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

	var atmosBuf matter.Atmosphere
	img := image.NewNRGBA(image.Rect(0, 0, 149, 149))
	const (
		tempMax  = 300
		oxyMax   = 25
		nitroMax = 85
		airMax   = 110
	)
	for i := 0; i < 10000; i++ {
		for _, t := range m[0].Atmos {
			img.SetNRGBA(int(t.X+74), int(t.Y+74), color.NRGBA{
				uint8(t.Temp / tempMax * t.Total() / airMax * 255),
				uint8(t.Gas[matter.Oxygen] / oxyMax * 255),
				uint8(t.Gas[matter.Nitrogen] / nitroMax * 255),
				255,
			})
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

		m[0].Atmos, atmosBuf = m[0].Atmos.Tick(atmosBuf)
	}
}
