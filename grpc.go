package frost

import (
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/kevinmichaelchen/frost/pb"
)

// GrpcServer is a gRPC server that provides an endpoint
// which allows apps to send "heartbeats" to request node IDs.
type GrpcServer struct {
	port int
}

func (service *GrpcServer) Heartbeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	nodeId := <-AvailableNodeIds
	return &pb.HeartbeatResponse{AppID: in.AppID, NodeID: int32(nodeId)}, nil
}

func (service *GrpcServer) Run() {
	address := fmt.Sprintf(":%d", GrpcPort)
	log.Printf("Listening on %s\n", address)
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