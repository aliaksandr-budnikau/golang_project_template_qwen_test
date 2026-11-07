//go:build grpc

// Package grpcsvc contains the gRPC service implementations. The real
// implementation is compiled only with -tags grpc; a no-op stub
// (server_default.go) is used otherwise so that `go build ./...` stays green
// without the gRPC dependencies or generated code.
package grpcsvc

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/yourname/empty_project/api/proto/v1"
)

// PingService implements the generated pb.PingServer.
type PingService struct {
	pb.UnimplementedPingServer
}

// NewPingService constructs a PingService.
func NewPingService() *PingService {
	return &PingService{}
}

// Ping handles the Ping RPC.
func (s *PingService) Ping(_ context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong: " + req.Message}, nil
}

// Register registers the gRPC services on the given server.
func Register(s grpc.ServiceRegistrar) {
	pb.RegisterPingServer(s, NewPingService())
}
