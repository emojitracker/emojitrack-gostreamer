package main

import (
	."github.com/azer/debug"
)


func ConnectionManager() chan SSEMessage {
	inputStream := make(chan SSEMessage)
	go h.run()
	go func() {
		for {
			msg := <-inputStream
			h.broadcast <- msg
		}
	}()
	return inputStream
}

type hub struct {
	// Registered connections.
	connections map[*connection]bool
	// Inbound messages to propogate out.
	broadcast chan SSEMessage
	// Register requests from the connections.
	register chan *connection
	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan SSEMessage),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			Debug("new connection being registered for " + c.channel)
			h.connections[c] = true
		case c := <-h.unregister:
			Debug("connection told us to unregister for " + c.channel)
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				if m.channel == c.channel {
					select {
					case c.send <- m.sseFormat():
					default:
						Debug("cant write to a connection, assuming it needs to be cleaned up")
						delete(h.connections, c)
						close(c.send)
						// go c.ws.Close()
					}
				}
			}
		}
	}
}
