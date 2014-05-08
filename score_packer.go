package main

import (
	"encoding/json"
	"time"
)

// ScorePacker receives a channel of string values, and counts their occurences.
//
// Every period duration, it flushes a *summary* of those values to returned
// output channel.  For now, this summary is returned as JSON, since that's how
// I always end up using it.
func ScorePacker(input <-chan string, period time.Duration) <-chan []byte {
	output := make(chan []byte)
	go func() {
		scores := make(map[string]uint)
		ticker := time.Tick(period)
		for {
			select {
			case id := <-input:
				// since int nil-value is zero dont have to worry about init case!
				scores[id] = scores[id] + 1 // increment hash
			case <-ticker:
				if len(scores) > 0 {
					packedscores, _ := json.Marshal(scores)
					output <- packedscores         // output to channel
					scores = make(map[string]uint) // reset hash
				}
			}
		}
	}()
	return output
}
