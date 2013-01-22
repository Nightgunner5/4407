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
			if m[i].Layout[c].Empty() {
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
				switch m[i].Layout[Coord{x, y}].Turf {
				case Space:
					m[i].Atmos.Set(Coord{x, y}, TileSpace())
				case Wall:
					m[i].Atmos.Set(Coord{x, y}, TileWall())
				case Floor:
					m[i].Atmos.Set(Coord{x, y}, TileIndoor())
				case Window:
					m[i].Atmos.Set(Coord{x, y}, TileWindow())
				case Airlock:
					m[i].Atmos.Set(Coord{x, y}, TileWall())
				case HeatVent:
					m[i].Atmos.Set(Coord{x, y}, TileHeater())
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

type LayoutTile struct {
	Turf    LayoutTileTurf
	Area    string
	Objects []LayoutObject
}

func (t LayoutTile) Empty() bool {
	return t.Turf == Space && t.Area == "" && len(t.Objects) == 0
}

type LayoutTileTurf uint32

const (
	Space LayoutTileTurf = iota
	Wall
	Floor
	Window
	Airlock
	HeatVent

	TileCount
)

func (t LayoutTileTurf) String() string {
	return t.GoString()
}

func (t LayoutTileTurf) GoString() string {
	switch t {
	case Space:
		return "Space"
	case Wall:
		return "Wall"
	case Floor:
		return "Floor"
	case Window:
		return "Window"
	case Airlock:
		return "Airlock"
	case HeatVent:
		return "HeatVent"
	}
	return strconv.FormatUint(uint64(t), 10)
}

type LayoutObject struct {
	Icon       string
	IconOffset uint16
}
