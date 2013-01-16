package matter

import (
	"testing"
)

func atmosEquals(tag string, a, b Atmosphere, epsilon float64, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("%q: len(%d) != len(%d)", tag, len(a), len(b))
	}

	for i := range a {
		c := a[i].Coord
		t1, t2 := a.Get(c), b.Get(c)
		if t1.Open != t2.Open {
			t.Errorf("%q: open %+v %b != %b", tag, c, t1.Open, t2.Open)
		}
		diff := t1.Temp - t2.Temp
		if -epsilon > diff || diff > epsilon {
			t.Errorf("%q: temp %+v %v != %v", tag, c, t1.Temp, t2.Temp)
		}
		for g := range t1.Gas {
			diff := t1.Gas[g] - t2.Gas[g]
			if -epsilon > diff || diff > epsilon {
				t.Errorf("%q: %s %+v %v != %v", tag, Gas(g), c, t1.Gas[g], t2.Gas[g])
			}
		}
	}
}

func TestAtmosphereSanity(t *testing.T) {
	const epsilon = 0.001

	var a Atmosphere

	a.Set(Coord{0, 0}, TileSpace())
	a.Set(Coord{0, 100}, TileIndoor())
	b, a := a.Tick(nil)
	atmosEquals("Two unrelated tiles", a, b, epsilon, t)

	a = nil
	a.Set(Coord{0, 0}, TileSpace())
	a.Set(Coord{0, 1}, TileWall())
	a.Set(Coord{0, 2}, TileIndoor())
	b, a = a.Tick(b)
	atmosEquals("Space, wall, floor", a, b, epsilon, t)

	a = nil
	a.Set(Coord{0, 0}, TileIndoor())
	a.Set(Coord{1, 1}, TileIndoor())
	a.Set(Coord{0, 1}, Tile{
		Gas:  a.Get(Coord{0, 0}).Gas,
		Temp: WaterFreezes,
		Open: true,
	})
	for i := 0; i < 100; i++ {
		a, b = a.Tick(b)
	}
	atmosEquals("Floor, cold floor, floor", a, b, epsilon, t)
}

func makeBenchmark(size int64) Atmosphere {
	a := make(Atmosphere, 0, size*size)
	for i := int64(0); i < size-1; i++ {
		a.Set(Coord{i, 0}, TileSpace())
		a.Set(Coord{0, i + 1}, TileSpace())
		a.Set(Coord{i + 1, size - 1}, TileSpace())
		a.Set(Coord{size - 1, i}, TileSpace())
	}
	for x := int64(1); x < size-1; x++ {
		for y := int64(1); y < size-1; y++ {
			a.Set(Coord{x, y}, TileIndoor())
		}
	}
	return a
}

func BenchmarkAtmosphereTick3x3(b *testing.B) {
	b.StopTimer()
	a := makeBenchmark(3)
	var c Atmosphere
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a, c = a.Tick(c)
	}
}

func BenchmarkAtmosphereTick5x5(b *testing.B) {
	b.StopTimer()
	a := makeBenchmark(5)
	var c Atmosphere
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a, c = a.Tick(c)
	}
}

func BenchmarkAtmosphereTick15x15(b *testing.B) {
	b.StopTimer()
	a := makeBenchmark(15)
	var c Atmosphere
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a, c = a.Tick(c)
	}
}
