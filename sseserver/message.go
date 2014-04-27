package sseserver

import (
	"fmt"
)

type SSEMessage struct {
	Event     string
	Data      []byte
	Namespace string
}

func (msg SSEMessage) sseFormat() []byte {
	if msg.Event != "" {
		return []byte(fmt.Sprintf("event:%s\ndata:%s\n\n", msg.Event, msg.Data))
	} else {
		return []byte(fmt.Sprintf("data:%s\n\n", msg.Data))
	}
}
