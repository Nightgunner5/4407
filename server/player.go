package main

import (
	"sync"
)

type Player struct {
	Send chan<- interface{}
	Z    int
}

var Players = struct {
	C map[*Player]bool
	sync.RWMutex
}{C: make(map[*Player]bool)}
