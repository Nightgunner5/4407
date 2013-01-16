package matter

type Coord struct {
	X, Y int64
}

func (c Coord) Add(x, y int64) Coord {
	c.X += x
	c.Y += y
	return c
}

type Tile struct {
	Gas  [gasCount]float64
	Temp float64
	Open bool
}

func TileIndoor() Tile {
	return Tile{
		Gas: [gasCount]float64{
			Oxygen:   TileContentsOxygen,
			Nitrogen: TileContentsNitrogen,
		},
		Temp: RoomTemperature,
		Open: true,
	}
}

func TileWall() Tile {
	return Tile{
		Temp: RoomTemperature,
	}
}

func TileSpace() Tile {
	return Tile{
		Temp: TempSpace,
		Open: true,
	}
}

func (t *Tile) tick(others ...Tile) {
	c := 1
	for _, o := range others {
		if o.Open {
			t.Temp += o.Temp
			for i := range o.Gas {
				t.Gas[i] += o.Gas[i]
			}
			c++
		}
	}
	t.Temp /= float64(c)
	for i := range t.Gas {
		t.Gas[i] /= float64(c)
	}
}

func (t Tile) Pressure() float64 {
	moles := t.Total()

	return moles * R / ATM / TileVolume * t.Temp
}

func (t Tile) Total() float64 {
	moles := 0.0

	for _, g := range t.Gas {
		moles += g
	}

	return moles
}

type Atmosphere map[Coord]Tile

func (orig Atmosphere) Tick() Atmosphere {
	a := make(Atmosphere, len(orig))
	for c, t := range orig {
		if t.Open {
			t.tick(t, orig[c.Add(-1, 0)], orig[c.Add(1, 0)], orig[c.Add(0, -1)], orig[c.Add(0, 1)])
		}
		a[c] = t
	}
	return a
}
