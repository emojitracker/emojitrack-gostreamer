package main

import (
  "fmt"
)

type SSEMessage struct {
	event     string
	data      []byte
	namespace string
}

func (msg SSEMessage) sseFormat() []byte {
	if msg.event != "" {
		return []byte(fmt.Sprintf("event:%s\ndata:%s\n\n", msg.event, msg.data))
	} else {
		return []byte(fmt.Sprintf("data:%s\n\n", msg.data))
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
