package main

import (
  "time"
  "encoding/json"
)

/*type ScoreUpdate struct {
  id string
  score int
}*/

/****
 * receives a channel of string values, and counts their occurences.
 * every PERIOD duration, flush a *summary* of those values to output channel.
 ****/
 // TODO: change to uint?

func ScorePacker(inputStream <-chan string, period time.Duration) <-chan []byte {
  outputStream := make(chan []byte)
  go func() {
    scores := make(map[string]int)
    ticker := time.Tick(period)
    for {
      select {
        case id := <-inputStream:
          //increment hash, since int nil-value is zero dont have to worry about init case!
          //(could distribute workers with shared memory and atomic uint)
          scores[id] = scores[id] + 1
        case <-ticker:
          packedscores, _ := json.Marshal(scores)
          outputStream <- packedscores   // output to channel
          scores = make(map[string]int)  // reset
      }
    }
  }()
  return outputStream
}
