package main

import (
	"time"
	"log"
)

// PruneStaleEntries checks for stale node IDs.
func PruneStaleEntries(heartbeatPeriodicity, sleepDuration time.Duration) {
	for {
		log.Println("Checking for stale node IDs...")
		now := time.Now()
		for nodeID, heartbeatTime := range Heartbeats {
			if now.After(heartbeatTime.Add(heartbeatPeriodicity)) {
				delete(NodeIds, nodeID)
				delete(Heartbeats, nodeID)
			}
		}
		printAvailableNodeIds()
		time.Sleep(sleepDuration)
	}
}