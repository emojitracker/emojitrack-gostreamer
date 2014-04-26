package main

import (
  "time"
  "encoding/json"
)

/****
 * receives a channel of string values, and counts their occurences.
 * every PERIOD duration, flush a *summary* of those values to output channel.
 ****/

func ScorePacker(inputStream <-chan string, period time.Duration) <-chan []byte {
  outputStream := make(chan []byte)
  go func() {
    scores := make(map[string]uint)
    ticker := time.Tick(period)
    for {
      select {
        case id := <-inputStream:
          //increment hash, since int nil-value is zero dont have to worry about init case!
          scores[id] = scores[id] + 1
        case <-ticker:
          packedscores, _ := json.Marshal(scores)
          outputStream <- packedscores   // output to channel
          scores = make(map[string]uint)  // reset
      }
    }
  }()
  return outputStream
}
