package sseserver

import (
	. "github.com/azer/debug"
	"log"
	"net/http"
)

// Interface to a SSE server.
//
// Exposes a send-only chan `broadcast`, any SSEMessage sent to this channel
// will be broadcast out to any connected clients subscribed to a namespace
// that matches the message.
type sseServer struct {
	Broadcast chan<- SSEMessage
	hub       *hub
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
	s.hub = newHub()
	go s.hub.run()

	// receive msgs to broadcast out to hub
	go func() {
		for {
			s.hub.broadcast <- <-inputStream
		}
	}()

	// return channel for incoming msgs
	return &s
}

// Begin serving SSE connections on specified addr.
// This method blocks forever, as it's basically a setup wrapper around
// `http.ListenAndServe()`
func (s *sseServer) Serve(addr string) {
	// use anonymous function for closure in order to pass value to handler
	// https://groups.google.com/forum/#!topic/golang-nuts/SGn1gd290zI
	http.HandleFunc("/subscribe/", func(w http.ResponseWriter, r *http.Request) {
		sseHandler(w, r, s.hub)
	})

	Debug("Starting server on addr " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
