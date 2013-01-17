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
				if math.IsNaN(h) || math.IsInf(h, 0) {
					panic("NaN")
				}
				hct1 += n1.Gas[g] * h
				hct2 += n2.Gas[g] * h

				delta := (t1.Gas[g] - t2.Gas[g]) / 8
				if math.IsNaN(delta) || math.IsInf(delta, 0) {
					panic("NaN")
				}
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
					if math.IsNaN(t1.Temp) || math.IsInf(t1.Temp, 0) {
						panic("NaN")
					}
					ht1 += h * delta * deltaTemp
					hcd1 += h * delta
				} else {
					if math.IsNaN(t2.Temp) || math.IsInf(t2.Temp, 0) {
						panic("NaN")
					}
					ht2 -= h * delta * deltaTemp
					hcd2 -= h * delta
				}
			}
			if math.IsNaN(hcn1) {
				panic("NaN")
			}
			if math.IsNaN(hcn2) {
				panic("NaN")
			}
			if math.IsNaN(hct1) {
				panic("NaN")
			}
			if math.IsNaN(hct2) {
				panic("NaN")
			}
			if math.IsNaN(hcd1) {
				panic("NaN")
			}
			if math.IsNaN(hcd2) {
				panic("NaN")
			}
			if math.IsNaN(ht1) {
				panic("NaN")
			}
			if math.IsNaN(ht2) {
				panic("NaN")
			}
			if hcn1 > 0.05 && hcn2 > 0.05 && (hcd1 > 0.0001 || hcd2 > 0.0001) {
				n1.Temp = (hct1*n1.Temp - hcd1*t1.Temp + hcd2*t2.Temp - ht1 + ht2) / hcn1
				if math.IsNaN(n1.Temp) || n1.Temp <= 0 {
					fmt.Println("(", hct1, "*", n1.Temp, "-", hcd1, "*", t1.Temp, "+", hcd2, "*", t2.Temp, "-", ht1, "+", ht2, ") /", hcn1)
					panic("NaN")
				}
				n2.Temp = (hct2*n2.Temp + hcd1*t1.Temp - hcd2*t2.Temp + ht1 - ht2) / hcn2
				if math.IsNaN(n2.Temp) || n2.Temp <= 0 {
					fmt.Println("(", hct2, "*", n2.Temp, "+", hcd1, "*", t1.Temp, "-", hcd2, "*", t2.Temp, "+", ht1, "-", ht2, ") /", hcn2)
					panic("NaN")
				}
			}
			if hct1 > 0.0003 && hct2 > 0.0003 && hct2 > hcn2*0.9 && hct2 < hcn2*1.1 {
				heat := t1.HeatTransfer * t2.HeatTransfer * deltaTemp / (hct1 + hct2)
				if math.IsNaN(heat) {
					panic("NaN")
				}

				n1.Temp -= heat * hct2
				if math.IsNaN(n1.Temp) || n1.Temp <= 0 {
					fmt.Println(heat, "*", hct2)
					panic("NaN")
				}
				n2.Temp += heat * hct1
				if math.IsNaN(n2.Temp) || n2.Temp <= 0 {
					fmt.Println(heat, "*", hct1)
					panic("NaN")
				}
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
	return other, orig
}
