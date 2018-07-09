package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/mroth/sseserver"
	"github.com/yvasiyarov/gorelic"
)

// Handles reporting status of this node to our stats block in Redis
// Used so we can monitor rollup status of multiple servers from one place.
func adminReporter(s *sseserver.Server) {
	go func() {
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
			// release redis conn back to pool
			rc.Close()
		}
	}()
}

// Runs the vendor package for reporting to New Relic
func gorelicMonitor() {
	if key := os.Getenv("NEW_RELIC_LICENSE_KEY"); key != "" {
		agent := gorelic.NewAgent()
		agent.NewrelicName = "emojitrack-gostreamer"
		agent.NewrelicLicense = key
		agent.Verbose = false
		agent.Run()
	}
}
