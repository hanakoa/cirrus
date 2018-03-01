package frost

import (
	"fmt"
	pb "github.com/kevinmichaelchen/frost/pb"
	"github.com/kevinmichaelchen/my-go-utils"
	"golang.org/x/net/context"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestFrost(t *testing.T) {
	log.Println("Seeding available node IDs...")
	for i := 0; i < NumNodes; i++ {
		AvailableNodeIds <- i
	}
	log.Println("Finished seeding...")
	f := &Server{}

	var wg sync.WaitGroup

	wg.Add(1)
	go f.run()

	time.Sleep(time.Second * 2)

	conn := utils.InitGrpcConn(fmt.Sprintf("%s:%d", "localhost", GrpcPort), 3, time.Second*5)
	client := pb.NewHeartbeatServiceClient(conn)

	for i := 0; i < NumNodes; i++ {
		wg.Add(1)
		// uuid.New().String()
		go spawnTestApp(strconv.Itoa(i), client)
	}

	wg.Wait()
}

func spawnTestApp(appID string, client pb.HeartbeatServiceClient) {
	// Send initial heartbeat to obtain initial node ID
	request := &pb.HeartbeatRequest{AppID: appID}
	if response, err := client.Heartbeat(context.Background(), request); err != nil {
		panic(err)
	} else {
		log.Printf("[App %s] -- Acquired node ID: %d", appID, response.NodeID)
	}

	var nodeID int32
	for {
		// Send a heartbeat periodically, with the current Node ID
		request := &pb.HeartbeatRequest{AppID: appID, NodeID: nodeID}
		if response, err := client.Heartbeat(context.Background(), request); err != nil {
			panic(err)
		} else {
			log.Printf("[App %s] -- Acquired node ID: %d", appID, response.NodeID)
			nodeID = response.NodeID
		}

		// Sleep for random time
		time.Sleep(time.Duration(rand.Intn(15)) * time.Second)
	}
}
