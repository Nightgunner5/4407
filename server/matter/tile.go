package matter

type Tile struct {
	Gas  [gasCount]float64
	Temp float64
}

func TileIndoor() *Tile {
	return &Tile{
		Gas: [gasCount]float64{
			Oxygen:   TileContentsOxygen,
			Nitrogen: TileContentsNitrogen,
		},
		Temp: RoomTemperature,
	}
}

func TileSpace() *Tile {
	return &Tile{
		Temp: TempSpace,
	}
}

func (t *Tile) Pressure() float64 {
	moles := t.Total()

	return moles * R / ATM / TileVolume * t.Temp
}

func (t *Tile) Total() float64 {
	moles := 0.0

	for _, g := range t.Gas {
		moles += g
	}

	return moles
}
