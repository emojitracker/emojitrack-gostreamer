package sseserver

import (
	"testing"
	"time"
)

func mockHub(conns int) (h *hub) {
	h = newHub()
	go h.run()
	for i := 0; i < conns; i++ {
		c := &connection{
			send: make(chan []byte, 256),
			created:   time.Now(),
			namespace: "/test",
		}
		h.register <- c
	}
	return h
}

func benchmarkBroadcast(conns int, b *testing.B) {
	h := mockHub(conns)

	for n := 0; n < b.N; n++ {
		h.broadcast <- SSEMessage{"", []byte("foo bar woo"), "/test"}
		h.broadcast <- SSEMessage{"event-foo", []byte("foo bar woo"), "/test"}

		// mock reading the connections
		// in theory this happens concurrently on another goroutine but here we will
		// exhaust the buffer quick if we dont force the read
		for c := range h.connections {
			<-c.send
			<-c.send
		}
	}
}

func BenchmarkBroadcast1(b *testing.B)    { benchmarkBroadcast(1, b) }
func BenchmarkBroadcast10(b *testing.B)   { benchmarkBroadcast(10, b) }
func BenchmarkBroadcast100(b *testing.B)  { benchmarkBroadcast(100, b) }
func BenchmarkBroadcast500(b *testing.B)  { benchmarkBroadcast(500, b) }
func BenchmarkBroadcast1000(b *testing.B) { benchmarkBroadcast(1000, b) }
