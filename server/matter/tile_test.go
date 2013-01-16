package matter

import (
	"testing"
)

func atmosEquals(tag string, a, b Atmosphere, epsilon float64, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("%q: len(%d) != len(%d)", tag, len(a), len(b))
	}

	for c := range a {
		if _, ok := b[c]; !ok {
			t.Errorf("%q: b is missing %+v", tag, c)
		}
	}

	for c := range b {
		if _, ok := a[c]; !ok {
			t.Errorf("%q: a is missing %+v", tag, c)
		}
	}

	for c, t1 := range a {
		if t2, ok := b[c]; ok {
			if t1.Open != t2.Open {
				t.Errorf("%q: open %b != %b", tag, t1.Open, t2.Open)
			}
			diff := t1.Temp - t2.Temp
			if -epsilon > diff || diff > epsilon {
				t.Errorf("%q: temp %v != %v", tag, t1.Temp, t2.Temp)
			}
			for g := range t1.Gas {
				diff := t1.Gas[g] - t2.Gas[g]
				if -epsilon > diff || diff > epsilon {
					t.Errorf("%q: %s %v != %v", tag, g, t1.Gas[g], t2.Gas[g])
				}
			}
		}
	}
}

func TestAtmosphereSanity(t *testing.T) {
	const epsilon = 0.001

	a := Atmosphere{
		Coord{0, 0}:   TileSpace(),
		Coord{0, 100}: TileIndoor(),
	}
	b := a.Tick()
	atmosEquals("Two unrelated tiles", a, b, epsilon, t)

	a = Atmosphere{
		Coord{0, 0}: TileSpace(),
		Coord{0, 1}: TileWall(),
		Coord{0, 2}: TileIndoor(),
	}
	b = a.Tick()
	atmosEquals("Space, wall, floor", a, b, epsilon, t)

	a = Atmosphere{
		Coord{0, 0}: TileIndoor(),
		Coord{1, 1}: TileIndoor(),
	}
	a[Coord{0, 1}] = Tile{
		Gas:  a[Coord{0, 0}].Gas,
		Temp: WaterFreezes,
		Open: true,
	}
	b = a
	for i := 0; i < 10; i++ {
		a = b
		b = a.Tick()
	}
	atmosEquals("Floor, cold floor, floor", a, b, epsilon, t)
}

func makeBenchmark(size int64) Atmosphere {
	a := make(Atmosphere, size*size)
	for i := int64(0); i < size-1; i++ {
		a[Coord{i, 0}] = TileSpace()
		a[Coord{0, i + 1}] = TileSpace()
		a[Coord{i + 1, size - 1}] = TileSpace()
		a[Coord{size - 1, i}] = TileSpace()
	}
	for x := int64(1); x < size-1; x++ {
		for y := int64(1); y < size-1; y++ {
			a[Coord{x, y}] = TileIndoor()
		}
	}
	return a
}

func BenchmarkAtmosphereTick3x3(b *testing.B) {
	b.StopTimer()
	a := makeBenchmark(3)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a = a.Tick()
	}
}

func BenchmarkAtmosphereTick5x5(b *testing.B) {
	b.StopTimer()
	a := makeBenchmark(5)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a = a.Tick()
	}
}

func BenchmarkAtmosphereTick15x15(b *testing.B) {
	b.StopTimer()
	a := makeBenchmark(15)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a = a.Tick()
	}
}
