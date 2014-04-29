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

	// set up the public interface
	var s = sseServer{}

	// start up our actual internal connection hub
	// which we keep in the server struct as private
	var h = newHub()
	s.hub = h
	go h.run()

	// expose just the broadcast chanel to public
	// will be typecast to send-only
	s.Broadcast = h.broadcast

	// return handle
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

	http.HandleFunc("/admin/status.json", func(w http.ResponseWriter, r *http.Request) {
		adminHandler(w, r, s.hub)
	})

	Debug("Starting server on addr " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
