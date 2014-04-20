package main

import (
  "log"
  "time"
  "net/http"
)

func main() {
/*  scoreUpdates := fakeRedisStream()*/
  scoreUpdates, _ := RedisGo()

  summarizedScores := ScorePacker(scoreUpdates, time.Duration(17*time.Millisecond))
  epsClients := ConnectionManager()

  go func() {
    for {
      val := <- summarizedScores
      epsClients <- SSEMessage{"",val,"/eps"}
    /*    epsClients <- sseEncoder(summarizedScores)*/
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
