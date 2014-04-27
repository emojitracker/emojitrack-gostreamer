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
	connections map[*connection]bool // Registered connections.
	broadcast chan SSEMessage 			 // Inbound messages to propogate out.
	register chan *connection 			 // Register requests from the connections.
	unregister chan *connection 		 // Unregister requests from connections.
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
			Debug("new connection being registered for " + c.namespace)
			h.connections[c] = true
		case c := <-h.unregister:
			Debug("connection told us to unregister for " + c.namespace)
			delete(h.connections, c)
			close(c.send)
		case m := <-h.broadcast:
			for c := range h.connections {
				if m.namespace == c.namespace {
					select {
					case c.send <- m.sseFormat():
					default:
						Debug("cant pass to a connection send chan, buffer is full -- kill it with fire")
						delete(h.connections, c)
						close(c.send)
						// go c.ws.Close()
					}
				}
			}
		}
	}
}
