package frost

import (
	"fmt"
	pb "github.com/kevinmichaelchen/frost/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

// GrpcServer is a gRPC server that provides an endpoint
// which allows apps to send "heartbeats" to request node IDs.
type GrpcServer struct {
	port int
}

func (service *GrpcServer) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Printf("[gRPC] -- Processing heartbeat for app %s", in.AppID)
	if nodeID := in.GetNodeID(); nodeID != 0 {
		// TODO do we need to lock Heartbeats, since the prune goroutine is deleting from it?
		if time.Now().Before(Heartbeats[int(nodeID)].Add(HeartbeatPeriodicity)) {
			log.Printf("[gRPC] -- App %s already has a valid node ID... Extending ID lifetime...", in.AppID)
			Heartbeats[int(nodeID)] = time.Now()
			return &pb.HeartbeatResponse{AppID: in.AppID, NodeID: nodeID}, nil
		}
		log.Printf("[gRPC] -- App %s has expired node ID %d", in.AppID, nodeID)
	}

	// TODO if you don't create more than 1024, there should always be an available nodeID,
	// however, we should probably return an error if there are no available nodes,
	// instead of letting the client hang
	// https://stackoverflow.com/questions/3398490/checking-if-a-channel-has-a-ready-to-read-value-using-go
	// TODO check if channel has ready-to-read value, otherwise return error
	nodeID := <-AvailableNodeIds
	log.Printf("[gRPC] -- Granting node ID %d to app %s", nodeID, in.AppID)
	Heartbeats[nodeID] = time.Now()
	return &pb.HeartbeatResponse{AppID: in.AppID, NodeID: int32(nodeID)}, nil
}

func (service *GrpcServer) run() {
	address := fmt.Sprintf(":%d", GrpcPort)
	log.Printf("Listening on %s", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Starting grpc server...")
	server := grpc.NewServer()

	// Register our services
	pb.RegisterHeartbeatServiceServer(server, service)

	// Register reflection service on gRPC server.
	reflection.Register(server)
	log.Println("Registered grpc services...")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
