package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/Nightgunner5/4407/server/matter"
	"log"
	"net/http"
)

func init() {
	http.Handle("/ws", websocket.Handler(socket))
}

func socket(conn *websocket.Conn) {
	defer conn.Close()

	addr := conn.Request().RemoteAddr
	log.Printf("Connect: %s", addr)
	defer log.Printf("Disconnect: %s", addr)

	send := make(chan interface{}, 16)

	p := &Player{
		Send: send,
	}
	Players.Lock()
	Players.C[p] = true
	Players.Unlock()
	defer func() {
		Players.Lock()
		delete(Players.C, p)
		Players.Unlock()

		close(send)
	}()

	go sockWrite(conn, send)

	send <- level(0)

	for {
		var packet *Packet
		if err := websocket.JSON.Receive(conn, &packet); err != nil {
			log.Printf("Read error: %s: %s", addr, err)
			return
		}
		switch {
		case packet.Position != nil:
			Players.RLock()
			c := p.Coord
			z := p.Z
			Players.RUnlock()

			if (c.X-packet.Position.X)*(c.X-packet.Position.X)+(c.Y-packet.Position.Y)*(c.Y-packet.Position.Y) != 1 {
				break
			}

			State.RLock()
			t := State.M[z].Layout[*packet.Position]
			State.RUnlock()

			if t == matter.Wall || t == matter.Window {
				break
			}

			Players.Lock()
			p.Coord = *packet.Position
			Players.Unlock()
		}
	}
}

func sockWrite(conn *websocket.Conn, send <-chan interface{}) {
	for packet := range send {
		websocket.JSON.Send(conn, packet)
	}
}
