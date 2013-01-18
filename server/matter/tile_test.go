package matter

import (
	"testing"
)

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
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a.Tick()
	}
}

func BenchmarkAtmosphereTick5x5(b *testing.B) {
	b.StopTimer()
	a := makeBenchmark(5)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a.Tick()
	}
}

func BenchmarkAtmosphereTick15x15(b *testing.B) {
	b.StopTimer()
	a := makeBenchmark(15)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		a.Tick()
	}
}
