package main

import (
	"time"
	"log"
	"sync"
	"github.com/google/uuid"
	"bytes"
	"strconv"
)

const (
	//NumNodes = 1024
	NumNodes = 3
	GrpcPort = 50051
)

// maps node IDs to app IDs
var NodeIds map[int]string

// maps node IDs to heartbeats
var Heartbeats map[int]time.Time

func main() {
	NodeIds = make(map[int]string, NumNodes)
	Heartbeats = make(map[int]time.Time, NumNodes)

	// how long an apps can abstain from heartbeat-ing its node ID
	// before we consider it stale
	heartbeatPeriodicity := time.Second * 10

	// how often we check for stale node IDs
	staleCheckPeriodicity := time.Second * 5

	seedTestData()

	var wg sync.WaitGroup
	wg.Add(1)
	go PruneStaleEntries(heartbeatPeriodicity, staleCheckPeriodicity)

	wg.Add(1)
	s := GrpcServer{port: GrpcPort}
	go s.Run()

	wg.Wait()
}

func seedTestData() {
	now := time.Now()
	for i := 0; i < NumNodes; i++ {
		NodeIds[i] = uuid.New().String()
		Heartbeats[i] = now
	}
}

func getAvailableNodeIds() []int {
	var nodeIds []int
	for i := 0; i < NumNodes; i++ {
		if _, ok := NodeIds[i]; !ok {
			nodeIds = append(nodeIds, i)
		}
	}
	return nodeIds
}

func printAvailableNodeIds() {
	var buffer bytes.Buffer
	nodeIds := getAvailableNodeIds()
	for _, n := range nodeIds {
		buffer.WriteString(" | ")
		buffer.WriteString(strconv.Itoa(n))
	}
	s := buffer.String()
	if len(s) > 0 {
		log.Println("Available nodes:", s)
	} else {
		log.Println("No nodes available... 😢")
	}
}
