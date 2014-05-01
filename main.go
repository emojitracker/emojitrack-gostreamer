package main

import (
	"log"
	"strings"
	"time"

	"github.com/mroth/emojitrack-gostreamer/sseserver"
)

func main() {
	// get us some data
	log.Println("Connecting to redis...")
	scoreUpdates, detailUpdates := RedisGo()

	// set up SSE server interface
	s := sseserver.NewServer()
	clients := s.Broadcast

	// fanout the scoreUpdates to two destinations
	rawScoreUpdates := make(chan RedisMsg)
	epsfeeder := make(chan RedisMsg)
	go func() {
		for scoreUpdate := range scoreUpdates {
			rawScoreUpdates <- scoreUpdate
			epsfeeder <- scoreUpdate
		}
	}()

	// Handle packing for epschan
	/*
	   This first goroutine basically grabs out just the data field of a RedisMsg,
	   and converts it to a string, because that's what my generic scorepacker
	   function expects to receive (for now).

	   Then, we just pipe that chan into a ScorePacker.
	*/
	scoreVals := make(chan string)
	epsScoreUpdates := ScorePacker(scoreVals, time.Duration(17*time.Millisecond))
	go func() {
		for {
			scoreVals <- string((<-epsfeeder).data)
		}
	}()

	// goroutines to handle passing messages to the proper connection pool
	// TODO: ask someone smart about whether each of these should be their own
	// goroutine, since the select here was kinda pointless since we dont need branching
	go func() {
		for msg := range rawScoreUpdates {
			clients <- sseserver.SSEMessage{
				Event:     "",
				Data:      msg.data,
				Namespace: "/raw",
			}
		}
	}()
	go func() {
		for val := range epsScoreUpdates {
			clients <- sseserver.SSEMessage{
				Event:     "",
				Data:      val,
				Namespace: "/eps",
			}
		}
	}()
	go func() {
		for msg := range detailUpdates {
			dchan := "/details/" + strings.Split(msg.channel, ".")[2]

			clients <- sseserver.SSEMessage{
				Event:     msg.channel,
				Data:      msg.data,
				Namespace: dchan,
			}
		}
	}()

	/*  go func() {
	    for {
	      select {
	        case msg := <- rawScoreUpdates:
	          clients <- SSEMessage{"",msg.data,"/raw"}
	        case val := <- epsScoreUpdates:
	          clients <- SSEMessage{"",val,"/eps"}
	        case msg := <- detailUpdates:
	          dchan := "/details/" + strings.Split(msg.channel, ".")[2]
	          clients <- SSEMessage{msg.channel,msg.data,dchan}
	      }
	    }
	  }()*/

	// start the monitor reporter to periodically send our status to redis
	go reporter(s)

	// share and enjoy
	port := ":8001"
	log.Println("Starting server on port " + port)
	log.Println("HOLD ON TO YOUR BUTTS...")

	// this method blocks by design
	s.Serve(port)
}

/*
  general patterns.

  redis -> chan scoreupdates
        -> chan detailstream

  *scoreupdates -> scorepacker -> chan epsstream
                -> chan rawstream

  *rawstream -> raw_pool => N clients
  *epsstream -> eps_pool => N clients
  *detailstream -> detail_pool => 4 clients for foo
                               -> 1 client  for bar
                               => 7 clients for xxx

  ^^^^ somehow buffered??

  status messages emitted from each pool on timer
  chan  <- raw_pool

  accumulator gofunc for reading status msgs from each chan
  emit on ticker to redis write...

  OR, a crazy DRY way to handle redis reporting...
    ...just HTTP hit localhost node for status, haha!
*/
