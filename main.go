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

	// fanout the scoreUpdates to two destinations
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
	epsScoreUpdates := ScorePacker(scoreVals, time.Duration(17*time.Millisecond))
	go func() {
		for {
			scoreVals <- string((<-epsfeeder).Data)
		}
	}()

	// goroutines to handle passing messages to the proper connection pool.
	//
	// I could use a select here and do as one goroutine, but having each be
	// independent could be slightly better for concurrency as these actually do
	// have a small amount of overhead in creating the SSEMessage so this is
	// theoretically better if we are running in parallel on appropriate hardware.

	// rawPublisher
	go func() {
		for msg := range rawScoreUpdates {
			clients <- sseserver.SSEMessage{
				Event:     "",
				Data:      msg.Data,
				Namespace: "/raw",
			}
		}
	}()

	// epsPublisher
	go func() {
		for val := range epsScoreUpdates {
			clients <- sseserver.SSEMessage{
				Event:     "",
				Data:      val,
				Namespace: "/eps",
			}
		}
	}()

	// detailPublisher
	go func() {
		for msg := range detailUpdates {
			dchan := "/details/" + strings.Split(msg.Channel, ".")[2]

			clients <- sseserver.SSEMessage{
				Event:     msg.Channel,
				Data:      msg.Data,
				Namespace: dchan,
			}
		}
	}()

	// monitoring in staging and production
	if envIsStaging() || envIsProduction() {
		// start the monitor reporter to periodically send our status to redis
		go adminReporter(s)
		// newrelic perf monitoring
		gorelicMonitor()
	}

	// share and enjoy
	port := envPort()
	log.Println("Starting server on port " + port)
	log.Println("HOLD ON TO YOUR BUTTS...")

	// this method blocks by design
	s.Serve(port)
}
