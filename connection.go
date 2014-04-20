package main

import (
	"net/http"
	"log"
	."github.com/azer/debug"
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
	cn := c.w.(http.CloseNotifier)
	closer := cn.CloseNotify()

	for {
		select {
			case message := <-c.send:
				_, err := c.w.Write(message)
				if err != nil {
					break
				}
		    if f, ok := c.w.(http.Flusher); ok {
		      f.Flush()
		    }
			case <-closer:
				Debug("closer fired for conn")
				return
		}
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	reqchan := r.URL.Path[10:] //strip out the prepending /subscribe
	//TODO: we should do the above in a clever way so we work on any path

	log.Println("CONNECT\t", reqchan, "\t", r.RemoteAddr)

	headers := w.Header()
	headers.Set("Access-Control-Allow-Origin", "*")
	headers.Set("Content-Type", "text/event-stream; charset=utf-8")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Connection", "keep-alive")
	headers.Set("Server", "emojitrack-gostreamer")

	c := &connection{ send: make(chan []byte, 256), w: w, channel: reqchan }
	h.register <- c

	defer func() {
		log.Println("DISCONNECT\t", reqchan, "\t", r.RemoteAddr)
		h.unregister <- c
	}()

	c.writer()
}
