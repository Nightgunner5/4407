package main

import (
	"code.google.com/p/go.net/websocket"
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
	}
}

func sockWrite(conn *websocket.Conn, send <-chan interface{}) {
	for packet := range send {
		websocket.JSON.Send(conn, packet)
	}
}
