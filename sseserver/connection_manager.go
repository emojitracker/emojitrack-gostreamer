package sseserver

import (
	"net/http"
	"log"
	."github.com/azer/debug"
)

/* PUBLIC INTERFACE */
// TODO: MOVE INTO OWN FILE

// Interface to a SSE server.
//
// Exposes a send-only chan `broadcast`, any SSEMessage sent to this channel
// will be broadcast out to any connected clients subscribed to a namespace
// that matches the message.
type sseServer struct {
	Broadcast chan<- SSEMessage
}

// Creates a new sseServer and returns a reference to it.
func SSEServer() *sseServer {
	// channel to receive msgs to broadcast.
	// we make here as bidirectional so we can read from it,
	// but cast to write-only in public interface.
	inputStream := make(chan SSEMessage)

	// set up the public interface
	var s = sseServer{
		Broadcast: inputStream,
	}

	// start up our actual internal connection hub
	go h.run()

	// receive msgs to broadcast out to hub
	go func() {
		for {
			msg := <-inputStream
			h.broadcast <- msg
		}
	}()

	// return channel for incoming msgs
	return &s
}

// Begin serving SSE connections on specified addr.
// This method blocks forever, as it's basically a setup wrapper around
// `http.ListenAndServe()`
func (s *sseServer) Serve(addr string) {
	http.HandleFunc("/subscribe/", sseHandler)

	Debug("Starting server on addr " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}


/* PRIVATE */

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
					}
				}
			}
		}
	}
}
