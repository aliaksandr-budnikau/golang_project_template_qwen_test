//go:build grpc

// Package grpcserver wraps a gRPC server. The real implementation is compiled
// only with -tags grpc; a no-op stub (server_default.go) is used otherwise so
// that `go build ./...` stays green without the gRPC dependencies.
package grpcserver

import (
	"net"

	"google.golang.org/grpc"
)

// Server wraps a *grpc.Server and the listener it serves on.
type Server struct {
	*grpc.Server
	lis net.Listener
}

// NewServer builds a gRPC server listening on addr.
func NewServer(addr string, opts ...grpc.ServerOption) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{Server: grpc.NewServer(opts...), lis: lis}, nil
}

// Serve starts serving on the listener. It returns when the server stops.
func (s *Server) Serve() error {
	return s.Server.Serve(s.lis)
}

// Stop initiates a graceful stop.
func (s *Server) Stop() {
	s.Server.GracefulStop()
}

// StopNow stops the server immediately.
func (s *Server) StopNow() {
	s.Server.Stop()
}
