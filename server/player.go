package main

import (
	"github.com/Nightgunner5/4407/server/matter"
	"sync"
)

type Player struct {
	Send chan<- interface{}
	matter.Coord
	Z int
}

var Players = struct {
	C map[*Player]bool
	sync.RWMutex
}{C: make(map[*Player]bool)}
