package main

import (
	"fmt"
	"time"

	"github.com/mroth/emojitrack-gostreamer/sseserver"
)

func reporter(s *sseserver.Server) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C
		fmt.Println(s.Status())
		//get redis conn from pool
		//report to redis
		//release redis conn
	}
}
