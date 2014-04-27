package sseserver

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
	namespace string
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
	namespace := r.URL.Path[10:] //strip out the prepending /subscribe
	//TODO: we should do the above in a clever way so we work on any path

	log.Println("CONNECT\t", namespace, "\t", r.RemoteAddr)

	headers := w.Header()
	headers.Set("Access-Control-Allow-Origin", "*")
	headers.Set("Content-Type", "text/event-stream; charset=utf-8")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Connection", "keep-alive")
	headers.Set("Server", "emojitrack-gostreamer")

	c := &connection{ send: make(chan []byte, 256), w: w, namespace: namespace }
	h.register <- c

	defer func() {
		log.Println("DISCONNECT\t", namespace, "\t", r.RemoteAddr)
		h.unregister <- c
	}()

	c.writer()
}
