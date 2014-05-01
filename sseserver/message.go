package sseserver

import (
	"fmt"
)

// SSEMessage is a message suitable for sending over a Server-Sent Event stream.
//
// The following fields are provided:
//   Event (string) - an event scope for the message [optional].
//   Data  ([]byte) - the message payload.
//   Namespace (string) - namespace to match a message to a client subscription.
//
// For more information on the SSE format itself, check out this article:
// http://www.html5rocks.com/en/tutorials/eventsource/basics/
//
// Note `Namespace` is not part of the SSE spec, it is merely used internally to
// map a message to the appropriate HTTP virtual endpoint.
//
type SSEMessage struct {
	Event     string
	Data      []byte
	Namespace string
}

// sseFormat is the formatted bytestring for a SSE message, ready to be sent.
func (msg SSEMessage) sseFormat() []byte {
	if msg.Event != "" {
		return []byte(fmt.Sprintf("event:%s\ndata:%s\n\n", msg.Event, msg.Data))
	} else {
		return []byte(fmt.Sprintf("data:%s\n\n", msg.Data))
	}
}
