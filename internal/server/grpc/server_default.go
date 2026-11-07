//go:build !grpc

package grpcserver

import "fmt"

// errNotBuilt is returned by the no-op stubs used when built without -tags grpc.
var errNotBuilt = fmt.Errorf("gRPC server not built: rebuild with `-tags grpc` (run `make proto` first)")

// Server is a no-op stub used when the service is built without -tags grpc.
type Server struct{}

// NewServer returns errNotBuilt when built without gRPC support.
func NewServer(_ string) (*Server, error) {
	return nil, errNotBuilt
}

// Serve returns errNotBuilt when built without gRPC support.
func (s *Server) Serve() error {
	return errNotBuilt
}

// Stop is a no-op when built without gRPC support.
func (s *Server) Stop() {}

// StopNow is a no-op when built without gRPC support.
func (s *Server) StopNow() {}
