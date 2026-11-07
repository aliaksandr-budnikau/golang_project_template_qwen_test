//go:build !grpc

package grpcsvc

// Register is a no-op when built without gRPC support. It exists so that
// `go build ./...` (without the gRPC tag) compiles this package without the
// generated protobuf code.
func Register() {}
