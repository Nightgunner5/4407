package matter

import (
	"fmt"
	"math"
)

type Coord struct {
	X, Y int64
}

func (c Coord) Add(x, y int64) Coord {
	c.X += x
	c.Y += y
	return c
}

type Tile struct {
	Gas          [gasCount]float64
	Temp         float64
	Open         bool
	Space        bool
	HeatTransfer float64
	HeatCapacity float64
}

func TileIndoor() Tile {
	return Tile{
		Gas: [gasCount]float64{
			0:        100,
			Oxygen:   TileContentsOxygen,
			Nitrogen: TileContentsNitrogen,
		},
		Temp:         RoomTemperature,
		Open:         true,
		HeatTransfer: 0.04,
		HeatCapacity: 225000,
	}
}

func TileWall() Tile {
	return Tile{
		Gas: [gasCount]float64{
			0: 100,
		},
		Temp:         RoomTemperature,
		HeatTransfer: 0.0005,
		HeatCapacity: 312500, //a little over 5 cm thick , 312500 for 1 m by 2.5 m by 0.25 m steel wall
	}
}

func TileWindow() Tile {
	return Tile{
		Gas: [gasCount]float64{
			0: 100,
		},
		Temp:         RoomTemperature,
		HeatTransfer: 0.03,
		HeatCapacity: 250000,
	}
}

func TileSpace() Tile {
	return Tile{
		Gas: [gasCount]float64{
			0: 100,
		},
		Temp:         TempSpace,
		Open:         true,
		Space:        true,
		HeatTransfer: 0.4,
		HeatCapacity: 700000,
	}
}

func (t Tile) Pressure() float64 {
	moles := t.Total()

	return moles * R / ATM / TileVolume * t.Temp
}

func (t Tile) Total() float64 {
	moles := 0.0

	for i, g := range t.Gas {
		if i != 0 {
			moles += g
		}
	}

	return moles
}

func (t Tile) check() bool {
	if math.IsNaN(t.Temp) || t.Temp <= 0 || t.Temp > 10000 {
		return false
	}

	for g := range t.Gas {
		if math.IsNaN(t.Gas[g]) || t.Gas[g] < 0 || (g == 0 && t.Gas[g] != 100) {
			return false
		}
	}

	return true
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

func (a Atmosphere) Tick() {
	share := func(t1, t2 *coordTile) {
		t1.Temp -= 2.7
		t2.Temp -= 2.7

		lc, rc := t1.HeatCapacity, t2.HeatCapacity // ending capacity
		ls, rs := lc*t1.Temp, rc*t2.Temp           // starting heat
		lt, rt := 0.0, 0.0                         // heat transferred

		for g := Gas(1); g < gasCount; g++ {
			l, r := t1.Gas[g]*GasMoveFraction, t2.Gas[g]*GasMoveFraction

			if !t1.Open || !t2.Open {
				l, r = 0, 0
			}

			ls += g.SpecificHeat() * t1.Gas[g] * t1.Temp
			rs += g.SpecificHeat() * t2.Gas[g] * t2.Temp

			t1.Gas[g] += r - l
			t2.Gas[g] += l - r

			lt += g.SpecificHeat() * l * (t1.Temp - t2.Temp)
			rt += g.SpecificHeat() * r * (t2.Temp - t1.Temp)

			lc += g.SpecificHeat() * t1.Gas[g]
			rc += g.SpecificHeat() * t2.Gas[g]
		}

		lt *= t1.HeatTransfer
		rt *= t2.HeatTransfer

		if lc > 0.001 && rc > 0.001 {
			ls += rt - lt
			t1.Temp = ls / lc
			rs += lt - rt
			t2.Temp = rs / rc
		}

		if dt, dh := ls/lc-rs/rc, ls-rs; math.Abs(dt) > 5 {
			if dh > 0 {
				dh *= t1.HeatTransfer
			} else {
				dh *= t2.HeatTransfer
			}

			t1.Temp -= dh / lc
			t2.Temp += dh / rc
		}

		t1.Temp += 2.7
		t2.Temp += 2.7
	}

	maybeShare := func(left int, c Coord) {
		right := a.coordIndex(c)
		if right < len(a) && a[right].Coord == c {
			share(&a[left], &a[right])
		}
	}

	for i := range a {
		c := a[i].Coord
		maybeShare(i, c.Add(-1, 0))
		maybeShare(i, c.Add(1, 0))
		maybeShare(i, c.Add(0, -1))
		maybeShare(i, c.Add(0, 1))
		if !a[i].Tile.check() {
			fmt.Printf("%#v\n", a[i])
			panic("NaN")
		}
	}
}
