package main

import (
	"log"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/mroth/sseserver"
)

func main() {
	// set up SSE server interface
	s := sseserver.NewServer()
	clients := s.Broadcast

	// get us some data from redis
	log.Println("Connecting to Redis...")
	initRedisPool()
	scoreUpdates, detailUpdates := myRedisSubscriptions()

	// fan out the scoreUpdates to two destinations
	rawScoreUpdates := make(chan redis.Message)
	epsfeeder := make(chan redis.Message)
	go func() {
		for scoreUpdate := range scoreUpdates {
			rawScoreUpdates <- scoreUpdate
			epsfeeder <- scoreUpdate
		}
	}()

	// Handle packing scores for eps namespace.
	//
	// This first goroutine basically grabs just the data field of a Redis msg,
	// and converts it to a string, because that's what my generic scorepacker
	// function expects to receive.
	//
	// Then, we just pipe that chan into a ScorePacker.
	scoreVals := make(chan string)
	go func() {
		for {
			scoreVals <- string((<-epsfeeder).Data)
		}
	}()
	epsScoreUpdates := ScorePacker(scoreVals, time.Duration(17*time.Millisecond))

	// "Fan in", creating proper namespaced SSEMessages depending on the
	// context, and delivers them to the sseserver broadcast hub for delivery to
	// clients.
	//
	// Also runs in yet another goroutine.
	go func() {
		for {
			select {
			// rawPublisher
			case msg := <-rawScoreUpdates:
				clients <- sseserver.SSEMessage{
					Event:     "",
					Data:      msg.Data,
					Namespace: "/raw",
				}
			// epsPublisher
			case val := <-epsScoreUpdates:
				clients <- sseserver.SSEMessage{
					Event:     "",
					Data:      val,
					Namespace: "/eps",
				}
			// detailPublisher
			case msg := <-detailUpdates:
				detailID := strings.Split(msg.Channel, ".")[2]
				namespace := "/details/" + detailID
				clients <- sseserver.SSEMessage{
					Event:     msg.Channel,
					Data:      msg.Data,
					Namespace: namespace,
				}
			}
		}
	}()

	// monitoring in staging and production
	if envIsStaging() || envIsProduction() {
		adminReporter(s) // periodically reports stats on node status into redis
	}

	// share and enjoy
	port := envPort()
	log.Println("Starting server on port", port)
	log.Println("HOLD ON TO YOUR BUTTS...")

	// this method blocks by design
	s.Serve(port)
}
