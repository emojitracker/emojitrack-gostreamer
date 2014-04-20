package main

import (
	"net/http"
	"log"
	."github.com/azer/debug"
)

type connection struct {
	// The websocket connection.
	/*    ws *websocket.Conn*/
	w http.ResponseWriter

	// Buffered channel of outbound messages.
	send chan []byte
}

// no reading only writing!!!!
/*func (c *connection) reader() {
    for {
        _, message, err := c.ws.ReadMessage()
        if err != nil {
            break
        }
        h.broadcast <- message
    }
    c.ws.Close()
}*/

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
	// if chan closes (or out of send loop somehow?), try to flush before finish
/*	if f, ok := c.w.(http.Flusher); ok {
		f.Flush()
	}*/
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	reqchan := r.URL.Path[10:] //strip out the prepending /subscribe
	log.Println("Connection to: ", reqchan)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	c := &connection{send: make(chan []byte, 256), w: w}
	h.register <- c
	defer func() {
		Debug("Disconnection")
		h.unregister <- c
	}()
	/*    go c.writer()*/
	/*    c.reader()*/
	c.writer()
}
