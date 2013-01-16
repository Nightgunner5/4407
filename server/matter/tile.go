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

	HeatTransfer float64
}

func TileIndoor() Tile {
	return Tile{
		Gas: [gasCount]float64{
			Oxygen:   TileContentsOxygen,
			Nitrogen: TileContentsNitrogen,
		},
		Temp:         RoomTemperature,
		Open:         true,
		HeatTransfer: 0.04,
	}
}

func TileWall() Tile {
	return Tile{
		Temp:         RoomTemperature,
		HeatTransfer: 0.0005,
	}
}

func TileSpace() Tile {
	return Tile{
		Temp:         TempSpace,
		Open:         true,
		HeatTransfer: 0.4,
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

func (orig Atmosphere) Tick(other Atmosphere) (new, old Atmosphere) {
	a := other
	if a == nil {
		a = make(Atmosphere, len(orig))
	}
	for c, t := range orig {
		a[c] = t
	}

	share := func(c1, c2 Coord) {
		t1, ok1 := orig[c1]
		t2, ok2 := orig[c2]
		if !ok1 || !ok2 {
			return
		}
		n1, n2 := a[c1], a[c2]

		deltaTemp := t1.Temp - t2.Temp
		heatTransfer := 0.0
		h1, h2 := 0.0, 0.0
		if t1.Open && t2.Open {
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
		}
		n1.Temp -= heatTransfer
		n2.Temp += heatTransfer
		if deltaTemp/10-heatTransfer > 0 {
			ht := t1.HeatTransfer * t2.HeatTransfer
			n1.Temp -= ht * deltaTemp * h2 / (h1 + h2)
			n2.Temp += ht * deltaTemp * h1 / (h1 + h2)
		}

		a[c1], a[c2] = n1, n2
	}

	for c := range orig {
		share(c, c.Add(-1, 0))
		share(c, c.Add(1, 0))
		share(c, c.Add(0, -1))
		share(c, c.Add(0, 1))
	}
	return a, orig
}
