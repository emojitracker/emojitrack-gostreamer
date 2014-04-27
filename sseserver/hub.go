package sseserver

import (
	."github.com/azer/debug"
)

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
				if m.Namespace == c.namespace {
					select {
					case c.send <- m.sseFormat():
					default:
						Debug("cant pass to a connection send chan, buffer is full -- kill it with fire")
						delete(h.connections, c)
						close(c.send)
						// go c.ws.Close()
						/* TODO: figure out what to do here...
							 we are already closing the send channel, in *theory* shouldn't the
							 connection clean up? I guess possible it doesnt if its deadlocked or
							 something... is it?

							 we want to make sure to always close the HTTP connection though,
							 so server can never fill up max num of open sockets.
						*/
					}
				}
			}
		}
	}
}
