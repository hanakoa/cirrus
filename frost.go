package frost

import (
	"time"
	"sync"
	"log"
)

const (
	//NumNodes = 1024
	NumNodes = 3
	GrpcPort = 50051
)

var AvailableNodeIds = make(chan int, NumNodes)

// maps node IDs to heartbeats
var Heartbeats map[int]time.Time

type FrostServer struct {
}

func (f *FrostServer) Run() {
	Heartbeats = make(map[int]time.Time, NumNodes)

	// how long an apps can abstain from heartbeat-ing its node ID
	// before we consider it stale
	heartbeatPeriodicity := time.Second * 10

	// how often we check for stale node IDs
	staleCheckPeriodicity := time.Second * 5

	var wg sync.WaitGroup
	wg.Add(1)
	go PruneStaleEntries(heartbeatPeriodicity, staleCheckPeriodicity)

	wg.Add(1)
	s := &GrpcServer{port: GrpcPort}
	go s.Run()

	wg.Wait()
}

// PruneStaleEntries checks for stale node IDs.
func PruneStaleEntries(heartbeatPeriodicity, sleepDuration time.Duration) {
	for {
		log.Println("Checking for stale node IDs...")
		for nodeID, heartbeatTime := range Heartbeats {
			if time.Now().After(heartbeatTime.Add(heartbeatPeriodicity)) {
				log.Printf("Node %d is newly available\n", nodeID)
				AvailableNodeIds <- nodeID
				delete(Heartbeats, nodeID)
			}
		}
		time.Sleep(sleepDuration)
	}
}