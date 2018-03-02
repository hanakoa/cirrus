package cirrus

import (
	"log"
	"sync"
	"time"
)

const (
	//NumNodes = 1024

	// NumNodes is the number of nodes we offer (max 1024)
	NumNodes = 10
	// GrpcPort is the port our gRPC server runs on
	GrpcPort = 50051
)

var (
	// HeartbeatPeriodicity is how often apps should heartbeat before we requisition their node IDs.
	HeartbeatPeriodicity = getHeartbeatPeriodicity()
	// StaleCheckPeriodicity is how often we check for stale node IDs
	StaleCheckPeriodicity = time.Second * 5
)

// AvailableNodeIds is a channel of available node IDs.
// As nodes die, their node IDs will be requisitioned and sent into the channel.
var AvailableNodeIds = make(chan int, NumNodes)

// maps node IDs to heartbeats
// TODO use int32, since protobuf has no vanilla int?

// Heartbeats stores the last time nodes sent a heartbeat.
var Heartbeats map[int]time.Time

// Server is the main struct used to run Cirrus.
// It runs a gRPC server for heartbeats, as well as handles requisitioning of stale node IDs.
type Server struct {
}

// how long an apps can abstain from heartbeat-ing its node ID
// before we consider it stale
func getHeartbeatPeriodicity() time.Duration {
	return time.Second * 10
}

func (f *Server) run() {
	Heartbeats = make(map[int]time.Time, NumNodes)

	log.Println("Seeding available node IDs...")
	for i := 0; i < NumNodes; i++ {
		AvailableNodeIds <- i
	}
	log.Println("Finished seeding...")

	var wg sync.WaitGroup
	wg.Add(1)
	go PruneStaleEntries(StaleCheckPeriodicity)

	wg.Add(1)
	s := &GrpcServer{port: GrpcPort}
	go s.run()

	wg.Wait()
}

// PruneStaleEntries checks for stale node IDs.
func PruneStaleEntries(sleepDuration time.Duration) {
	for {
		log.Println("Checking for stale node IDs...")
		for nodeID, heartbeatTime := range Heartbeats {
			if time.Now().After(heartbeatTime.Add(HeartbeatPeriodicity)) {
				log.Printf("Node %d has been requisitioned", nodeID)
				AvailableNodeIds <- nodeID
				delete(Heartbeats, nodeID)
			}
		}
		time.Sleep(sleepDuration)
	}
}
