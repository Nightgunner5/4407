package main

import (
	"time"
)

func atmosphere() {
	for {
		State.Lock()
		for i := range State.M {
			State.M[i].Atmos.Tick()
		}
		State.Unlock()

		time.Sleep(time.Second)
	}
}
