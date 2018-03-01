package frost

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/kevinmichaelchen/frost/pb"
	"time"
)

// GrpcServer is a gRPC server that provides an endpoint
// which allows apps to send "heartbeats" to request node IDs.
type GrpcServer struct {
	port int
}

func (service *GrpcServer) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	log.Printf("[gRPC] -- Processing heartbeat for app %s", in.AppID)
	if nodeId := in.GetNodeID(); nodeId != 0 {
		// TODO do we need to lock Heartbeats, since the prune goroutine is deleting from it?
		if time.Now().Before(Heartbeats[int(nodeId)].Add(HeartbeatPeriodicity)) {
			log.Printf("[gRPC] -- App %s already has a valid node ID... Extending ID lifetime...", in.AppID)
			Heartbeats[int(nodeId)] = time.Now()
			return &pb.HeartbeatResponse{AppID: in.AppID, NodeID: nodeId}, nil
		} else {
			log.Printf("[gRPC] -- App %s has expired node ID %d", in.AppID, nodeId)
		}
	}

	nodeId := <-AvailableNodeIds
	log.Printf("[gRPC] -- Granting node ID %d to app %s", nodeId, in.AppID)
	Heartbeats[nodeId] = time.Now()
	// TODO if you don't create more than 1024, there should always be an available nodeID,
	// however, we should probably return an error if there are no available nodes,
	// instead of letting the client hang
	return &pb.HeartbeatResponse{AppID: in.AppID, NodeID: int32(nodeId)}, nil
}

func (service *GrpcServer) Run() {
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