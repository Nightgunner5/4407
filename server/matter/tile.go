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

func TileWindow() Tile {
	return Tile{
		Temp:         RoomTemperature,
		HeatTransfer: 0.03,
	}
}

func TileSpace() Tile {
	return Tile{
		Temp:         TempSpace,
		Open:         true,
		Space:        true,
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

func (t Tile) check() bool {
	if math.IsNaN(t.Temp) {
		return false
	}

	for g := range t.Gas {
		if math.IsNaN(t.Gas[g]) || t.Gas[g] < 0 || (g == 0 && t.Gas[g] != 0) {
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
		if n1.Open && n2.Open && (math.Abs(deltaTemp) > 0.001 || math.Abs(n1.Pressure()-n2.Pressure()) > 0.01) {
			ht1, ht2 := 0.0, 0.0
			hct1, hct2 := 0.0, 0.0
			hcn1, hcn2 := 0.0, 0.0
			hcd1, hcd2 := 0.0, 0.0
			for g := range t1.Gas {
				h := Gas(g).SpecificHeat()

				hct1 += n1.Gas[g] * h
				hct2 += n2.Gas[g] * h

				delta := (n1.Gas[g] - n2.Gas[g]) / 5

				if (-0.0001 < delta && delta < 0.0001) || (t1.Gas[g] < 0.001 && t2.Gas[g] < 0.001) {
					hcn1 += n1.Gas[g] * h
					hcn2 += n2.Gas[g] * h
					continue
				}

				n1.Gas[g] -= delta
				n2.Gas[g] += delta

				hcn1 += n1.Gas[g] * h
				hcn2 += n2.Gas[g] * h

				if delta > 0 {
					ht1 += h * delta * deltaTemp
					hcd1 += h * delta
				} else {
					ht2 -= h * delta * deltaTemp
					hcd2 -= h * delta
				}
			}

			if hcn1 > 0.05 && hcn2 > 0.05 && (hcd1 > 0.0001 || hcd2 > 0.0001) {
				n1.Temp = (hct1*n1.Temp - hcd1*t1.Temp + hcd2*t2.Temp - ht1 + ht2) / hcn1
				n2.Temp = (hct2*n2.Temp + hcd1*t1.Temp - hcd2*t2.Temp + ht1 - ht2) / hcn2
			}

			if hct1 + hct2 > 0.001 && hct2 > hcn2*0.9 && hct2 < hcn2*1.1 {
				heat := t1.HeatTransfer * t2.HeatTransfer * deltaTemp / (hct1 + hct2)

				n1.Temp -= heat * hct2
				n2.Temp += heat * hct1
			}
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
		if !other[i].Tile.check() {
			fmt.Printf("%#v\n", other[i])
			panic("NaN")
		}
	}
	//for i := range other {
	//	if other[i].Space {
	//		other[i].Gas = [gasCount]float64{}
	//		other[i].Temp = TempSpace
	//	}
	//}
	return other, orig
}
