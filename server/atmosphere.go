package main

import (
	"github.com/Nightgunner5/4407/matter"
	"time"
)

type transferAtmos struct {
	Temp          float64
	Oxygen        float64
	Nitrogen      float64
	CarbonDioxide float64
	NitrousOxide  float64
	Plasma        float64
	X             int64
	Y             int64
}

func atmosphere() {
	for {
		var playerLocations []struct {
			XY matter.Coord
			Z  int
		}
		Players.RLock()
		for p := range Players.C {
			playerLocations = append(playerLocations, struct {
				XY matter.Coord
				Z  int
			}{p.Coord, p.Z})
		}
		Players.RUnlock()

		State.Lock()
		for _, p := range playerLocations {
			t := State.M[p.Z].Atmos.Get(p.XY)
			if t != nil {
				t.Temp += 2.5
			}
		}
		for z := range State.M {
			State.M[z].Atmos.Tick()
		}
		State.Unlock()

		State.RLock()
		a := make([]struct{ Atmos []transferAtmos }, len(State.M))
		for i, l := range State.M {
			for _, t := range l.Atmos {
				a[i].Atmos = append(a[i].Atmos, transferAtmos{
					t.Temp,
					t.Gas[matter.Oxygen],
					t.Gas[matter.Nitrogen],
					t.Gas[matter.CarbonDioxide],
					t.Gas[matter.NitrousOxide],
					t.Gas[matter.Plasma],
					t.X,
					t.Y,
				})
			}
		}
		State.RUnlock()

		Players.RLock()
		for p := range Players.C {
			select {
			case p.Send <- a[p.Z]:
			default:
			}
		}
		Players.RUnlock()

		time.Sleep(time.Second)
	}
}
