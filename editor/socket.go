package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/Nightgunner5/4407/matter"
	"github.com/Nightgunner5/4407/shared"
	"log"
	"net/http"
	"sync"
)

func init() {
	http.Handle("/ws", websocket.Handler(socket))
}

type Client struct {
	Send chan<- interface{}
}

var Clients = struct {
	M map[*Client]bool
	sync.RWMutex
}{
	M: make(map[*Client]bool),
}

type Packet struct {
	Set *struct {
		Z    int
		X, Y int64
		Turf matter.LayoutTileTurf
	}
	NewLevel *struct{}
	Save     *struct{}
}

func MapZPacket(i int) interface{} {
	State.RLock()
	defer State.RUnlock()

	return struct {
		MapZ struct {
			M interface{}
			Z int
		}
	}{struct {
		M interface{}
		Z int
	}{shared.Level(State.M[i].Layout).Map, i}}
}

func socket(conn *websocket.Conn) {
	defer conn.Close()

	addr := conn.Request().RemoteAddr
	log.Printf("Connect: %s", addr)
	defer log.Printf("Disconnect: %s", addr)

	send := make(chan interface{}, 16)

	c := &Client{
		Send: send,
	}
	Clients.Lock()
	Clients.M[c] = true
	Clients.Unlock()
	defer func() {
		Clients.Lock()
		delete(Clients.M, c)
		Clients.Unlock()

		close(send)
	}()

	go sockWrite(conn, send)

	{
		State.RLock()
		zLevels := len(State.M)
		State.RUnlock()

		for i := 0; i < zLevels; i++ {
			send <- MapZPacket(i)
		}
	}
	for {
		var packet *Packet
		if err := websocket.JSON.Receive(conn, &packet); err != nil {
			log.Printf("Read error: %s: %s", addr, err)
			return
		}
		if packet.Set != nil {
			State.Lock()
			tile := State.M[packet.Set.Z].Layout[matter.Coord{packet.Set.X, packet.Set.Y}]
			tile.Turf = packet.Set.Turf
			State.M[packet.Set.Z].Layout[matter.Coord{packet.Set.X, packet.Set.Y}] = tile
			State.Unlock()

			packet := MapZPacket(packet.Set.Z)

			Clients.RLock()
			for c := range Clients.M {
				select {
				case c.Send <- packet:
				default:
				}
			}
			Clients.RUnlock()
		}

		if packet.NewLevel != nil {
			State.Lock()
			State.M.NewLevel()
			l := len(State.M) - 1
			State.Unlock()

			packet := MapZPacket(l)

			Clients.RLock()
			for c := range Clients.M {
				select {
				case c.Send <- packet:
				default:
				}
			}
			Clients.RUnlock()
		}

		if packet.Save != nil {
			save()
		}
	}
}

func sockWrite(conn *websocket.Conn, send <-chan interface{}) {
	for packet := range send {
		websocket.JSON.Send(conn, packet)
	}
}
