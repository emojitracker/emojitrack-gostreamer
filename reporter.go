package main

import (
	"encoding/json"
	"time"

	"github.com/mroth/emojitrack-gostreamer/sseserver"
)

func reporter(s *sseserver.Server) {
	ticker := time.NewTicker(5 * time.Second)
	for {
		// block waiting for ticker
		<-ticker.C

		// get redis conn from pool
		rc := redisPool.Get()

		// report to redis
		report, _ := json.Marshal(s.Status())
		serverNode := s.Status().Node
		rc.Do("HSET", "admin_stream_status", serverNode, report)

		// release redis conn
		rc.Close()
	}
}
