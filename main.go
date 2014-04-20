package main

import (
  "log"
  "time"
  "net/http"
  "strings"
)

func main() {
  // set up a connection pool to receive clients on
  clients := ConnectionManager()

  // get us some data
/*  scoreUpdates := fakeRedisStream()*/
  scoreUpdates, detailUpdates := RedisGo()

  // create a channel of just the values from the RedisMsg for scoreupdates
  // suitable to sending to my generic scorepacker
  scoreVals := make(chan string)
  go func() {
    for {
      msg := <- scoreUpdates
      scoreVals <- string(msg.data)
    }
  }()
  // then send it to that scorepacker
  summarizedScores := ScorePacker(scoreVals, time.Duration(17*time.Millisecond))

  // goroutine to handle passing messages to the proper connection pool
  // TODO: ask someone smart about whether each of these should be their own
  // goroutine, since the select here is kinda pointless since we dont need branching
  go func() {
    for {
      select {
        case val := <- summarizedScores:
          clients <- SSEMessage{"",val,"/eps"}
        case msg := <- detailUpdates:
          dchan := "/details/" + strings.Split(msg.channel, ".")[2]
          clients <- SSEMessage{msg.channel,msg.data,dchan}
      }
    }
  }()

  http.HandleFunc("/subscribe/", sseHandler)
  if err := http.ListenAndServe(":8001", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }

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

*/
