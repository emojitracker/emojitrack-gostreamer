package main

import (
  "fmt"
)

type SSEMessage struct {
	event   string
	data    []byte
	channel string
}

func (msg SSEMessage) sseFormat() string {
	if msg.event != "" {
		return fmt.Sprintf("event:%v\ndata:%v\n\n", msg.event, string(msg.data))
	} else {
		return fmt.Sprintf("data:%v\n\n", string(msg.data))
	}
}

/* does it actually make sense to have this all be its own routine?
   after all a connection hub is only going to manage sse stuff, so maybe it shoudl understand natively?
   or maybe this actually gets spawned inside the CH logic rather than in main... */

/* Transforms a chan of SSEMessages into a serialized stream of bytes */
/* so this whole idea is maybe uncessary and dumb no more writing code on 1am on school night */
func sseEncoder( in <-chan SSEMessage ) chan<- []byte {
  out := make(chan<- []byte)
  go func() {
    for {
      msg := <-in
      out <- []byte( msg.sseFormat() )
      /* TODO: can above be expressed as `out <- <-in.sseFormat()` ? */
    }
  }()
  return out
}
