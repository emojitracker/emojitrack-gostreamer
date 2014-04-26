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
