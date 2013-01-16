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

type coordTile struct {
	Coord
	Tile
}

func coordLess(a, b Coord) bool {
	if a.Y != b.Y {
		return a.Y < b.Y
	}
	return a.X < b.X
}

type Atmosphere []coordTile

func (a Atmosphere) coordIndex(c Coord) int {
	// Define f(-1) == false and f(n) == true.
	// Invariant: f(i-1) == false, f(j) == true.
	i, j := 0, len(a)
	for i < j {
		h := i + (j-i)/2 // avoid overflow when computing h
		// i ≤ h < j
		if a[h].Coord != c && !coordLess(a[h].Coord, c) {
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	// i == j, f(i-1) == false, and f(j) (= f(i)) == true  =>  answer is i.
	return i
}

func (a Atmosphere) Get(c Coord) *Tile {
	i := a.coordIndex(c)
	if i < len(a) && a[i].Coord == c {
		return &a[i].Tile
	}
	return nil
}

func (a *Atmosphere) Set(c Coord, t Tile) {
	i := a.coordIndex(c)
	if i == len(*a) {
		*a = append(*a, coordTile{c, t})
		return
	}
	if (*a)[i].Coord == c {
		(*a)[i].Tile = t
		return
	}
	*a = append((*a)[:i+1], (*a)[i:]...)
	(*a)[i] = coordTile{c, t}
}

func (orig Atmosphere) Tick(other Atmosphere) (new, old Atmosphere) {
	if cap(other) < len(orig) {
		other = make(Atmosphere, len(orig))
	}
	other = other[:len(orig)]
	copy(other, orig)

	share := func(i, j int) {
		t1, t2 := orig[i], orig[j]
		n1, n2 := &other[i], &other[j]

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
	}

	maybeShare := func(left int, right Coord) {
		i := orig.coordIndex(right)
		if i < len(orig) && orig[i].Coord == right {
			share(left, i)
		}
	}

	for i := range orig {
		maybeShare(i, orig[i].Coord.Add(-1, 0))
		maybeShare(i, orig[i].Coord.Add(1, 0))
		maybeShare(i, orig[i].Coord.Add(0, -1))
		maybeShare(i, orig[i].Coord.Add(0, 1))
	}
	return other, orig
}
