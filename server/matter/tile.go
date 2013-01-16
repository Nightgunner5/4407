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
		a[c] = t
	}

	share := func(c1, c2 Coord) {
		t1, t2 := orig[c1], orig[c2]
		if !t1.Open || !t2.Open {
			return
		}
		n1, n2 := a[c1], a[c2]

		deltaTemp := t1.Temp - t2.Temp
		heatTransfer := 0.0
		h1, h2 := 0.0, 0.0
		for g := range t1.Gas {
			delta := (t1.Gas[g] - t2.Gas[g]) / 5
			h := Gas(g).SpecificHeat()
			h1 += t1.Gas[g] * h
			h2 += t2.Gas[g] * h
			if delta > 0 {
				heatTransfer += h * delta * t1.Temp
			} else {
				heatTransfer += h * delta * t2.Temp
			}
			n1.Gas[g] -= delta
			n2.Gas[g] += delta
		}
		n1.Temp -= heatTransfer
		n2.Temp += heatTransfer
		if deltaTemp/10/heatTransfer > 1 {
			n1.Temp += 0.4 * deltaTemp * h2 / (h1 + h2)
			n2.Temp += 0.4 * deltaTemp * h1 / (h1 + h2)
		}

		a[c1], a[c2] = n1, n2
	}

	for c := range orig {
		share(c, c.Add(-1, 0))
		share(c, c.Add(1, 0))
		share(c, c.Add(0, -1))
		share(c, c.Add(0, 1))
	}
	return a
}
