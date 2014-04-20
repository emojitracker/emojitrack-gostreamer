package main

import (
	"net/http"
	"log"
/*	."github.com/azer/debug"*/
)

type connection struct {
	// The HTTP connection.
	w http.ResponseWriter

	// Buffered channel of outbound messages.
	send chan []byte

	// The conceptual "channel" the SSE client is requesting
	/* Yeah, this is a namespace collision with the Go language,
			But the ship has already sailed on that one since the API
			for this has been long defined. */
	channel string
}

func (c *connection) writer() {
	// read as long as channel is open
	for message := range c.send {
		_, err := c.w.Write(message)
		if err != nil {
			break
		}
    if f, ok := c.w.(http.Flusher); ok {
      f.Flush()
    }
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	reqchan := r.URL.Path[10:] //strip out the prepending /subscribe
	//TODO: we should do the above in a clever way so we work on any path

	log.Println("Connection to: ", reqchan)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	c := &connection{ send: make(chan []byte, 256), w: w, channel: reqchan }

	h.register <- c
	defer func() {
		log.Println("Disconnection from: ", reqchan)
		h.unregister <- c
	}()

	c.writer()
}
