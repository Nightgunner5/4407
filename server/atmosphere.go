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
				if t.Gas[matter.Oxygen] > 1 {
					t.Gas[matter.Oxygen] -= 1
					t.Gas[matter.CarbonDioxide] += 1
				} else {
					// TODO: damage
				}
			}
		}
		for z := range State.M {
			State.M[z].Atmos.Tick()
		}
		State.Unlock()

		Players.RLock()
		for p := range Players.C {
			coord := p.Coord
			z := p.Z
			Players.RUnlock()
			State.RLock()
			var a []transferAtmos
			for _, t := range State.M[z].Atmos {
				if t.X < coord.X-25 ||
					t.X > coord.X+25 ||
					t.Y < coord.Y-25 ||
					t.Y > coord.Y+25 {
					continue
				}
				a = append(a, transferAtmos{
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
			State.RUnlock()

			select {
			case p.Send <- struct{ Atmos []transferAtmos }{a}:
			default:
			}

			Players.RLock()
		}
		Players.RUnlock()

		time.Sleep(time.Second)
	}
}
