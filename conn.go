package main

import (
	"golang.org/x/net/websocket"
	"fmt"
	"time"
)

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string
}

func (c *connection) writer() {
	for {
		select {
		case message := <-c.send:
			err := websocket.Message.Send(c.ws, message)
			if err != nil {
				break
			}
		}
	}
	c.ws.Close()
}

func (c *connection) timeTeller() {
	t := time.Tick(2 * time.Second)
	for now := range t {
		fmt.Println(now)
		h.broadcast <- now.String()
	}
}

func wsHandler(ws *websocket.Conn) {
	c := &connection{send: make(chan string, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	c.writer()
}
