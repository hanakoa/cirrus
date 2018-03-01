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
	nodeId := <-AvailableNodeIds
	log.Printf("[gRPC] -- Granting node ID %d to app %s", nodeId, in.AppID)
	Heartbeats[nodeId] = time.Now()
	// TODO you should pass in your current node ID so we can do a quick lookup
	// TODO if there are no available nodes left, we shouldn't hang
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