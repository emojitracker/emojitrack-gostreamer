package sseserver

import (
	. "github.com/azer/debug"
	"log"
	"net/http"
)

type connection struct {
	w         http.ResponseWriter // The HTTP connection
	send      chan []byte         // Buffered channel of outbound messages.
	namespace string              // Conceptual "channel" SSE client is requesting
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
	namespace := r.URL.Path[10:] // strip out the prepending "/subscribe"
	// TODO: we should do the above in a clever way so we work on any path

	log.Println("CONNECT\t", namespace, "\t", r.RemoteAddr)

	headers := w.Header()
	headers.Set("Access-Control-Allow-Origin", "*")
	headers.Set("Content-Type", "text/event-stream; charset=utf-8")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Connection", "keep-alive")
	headers.Set("Server", "emojitrack-gostreamer")

	c := &connection{send: make(chan []byte, 256), w: w, namespace: namespace}
	h.register <- c

	defer func() {
		log.Println("DISCONNECT\t", namespace, "\t", r.RemoteAddr)
		h.unregister <- c
	}()

	c.writer()
}
