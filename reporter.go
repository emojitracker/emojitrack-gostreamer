package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/mroth/sseserver"
	"github.com/yvasiyarov/gorelic"
)

// Handles reporting status of this application to external services.

func adminReporter(s *sseserver.Server) {
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

func gorelicMonitor() {

	if envIsStaging() || envIsProduction() {
		if key := os.Getenv("NEW_RELIC_LICENSE_KEY"); key != "" {
			agent := gorelic.NewAgent()
			agent.NewrelicName = "emojitrack-gostreamer"
			agent.NewrelicLicense = key
			agent.Verbose = false
			agent.Run()
		}
	}

}
