package matter

import (
	"strconv"
)

type Map []Level

func (m *Map) NewLevel() Layout {
	l := make(Layout)
	*m = append(*m, Level{Layout: l})
	return l
}

func (m Map) Compile(padding int64) {
	for i := range m {
		var min, max Coord
		for c := range m[i].Layout {
			if m[i].Layout[c] == Space {
				delete(m[i].Layout, c)
				continue
			}
			if c.X < min.X {
				min.X = c.X
			}
			if c.Y < min.Y {
				min.Y = c.Y
			}
			if c.X > max.X {
				max.X = c.X
			}
			if c.Y > max.Y {
				max.Y = c.Y
			}
		}
		min.X -= padding
		min.Y -= padding
		max.X += padding
		max.Y += padding

		m[i].Min = min
		m[i].Max = max

		m[i].Atmos = make(Atmosphere, 0, (max.X-min.X)*(max.Y-min.Y))
		for y := min.Y; y < max.Y; y++ {
			for x := min.X; x < max.X; x++ {
				switch m[i].Layout[Coord{x, y}] {
				case Space:
					m[i].Atmos.Set(Coord{x, y}, TileSpace())
				case Wall:
					m[i].Atmos.Set(Coord{x, y}, TileWall())
				case Floor:
					m[i].Atmos.Set(Coord{x, y}, TileIndoor())
				case Window:
					m[i].Atmos.Set(Coord{x, y}, TileWindow())
				}
			}
		}
	}
}

type Level struct {
	Layout Layout

	Atmos    Atmosphere
	Min, Max Coord
}

type Layout map[Coord]LayoutTile

type LayoutTile uint32

const (
	Space LayoutTile = iota
	Wall
	Floor
	Window
)

func (t LayoutTile) String() string {
	return t.GoString()
}

func (t LayoutTile) GoString() string {
	switch t {
	case Space:
		return "Space"
	case Wall:
		return "Wall"
	case Floor:
		return "Floor"
	case Window:
		return "Window"
	}
	return strconv.FormatUint(uint64(t), 10)
}
