package main

import (
	"time"
	"log"
	"sync"
	"github.com/google/uuid"
)

var nodeIds map[string]int
var heartbeats map[string]time.Time

func main() {
	nodeIds = make(map[string]int, 1024)
	heartbeats = make(map[string]time.Time, 1024)

	// apps have to heartbeat every 30 seconds
	heartbeatPeriodicity := time.Second * 30

	seedTestData()

	var wg sync.WaitGroup
	wg.Add(1)
	go pruneStaleEntries(heartbeatPeriodicity, time.Second * 2)
	wg.Wait()
}

func seedTestData() {
	now := time.Now()
	for i := 0; i < 10; i++ {
		heartbeats[uuid.New().String()] = now
	}
}

// pruneStaleEntries checks for stale node IDs.
func pruneStaleEntries(heartbeatPeriodicity, sleepDuration time.Duration) {
	for {
		log.Println("Checking for stale node IDs...")
		now := time.Now()
		for appID, heartbeatTime := range heartbeats {
			if now.After(heartbeatTime.Add(heartbeatPeriodicity)) {
				log.Printf("App %s is expired. Last heartbeat was %s", appID, heartbeatTime)
			}
		}
		time.Sleep(sleepDuration)
	}
}