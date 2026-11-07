//go:build grpc

package grpcserver

import (
	"context"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithRecovery returns a grpc.ServerOption that installs interceptors which turn
// a panic in a downstream RPC handler into a codes.Internal error and log the
// full stack trace via the structured logger.
//
// This mirrors the HTTP Recovery middleware (see
// internal/transport/http/middleware): a panic is recovered, its stack trace
// logged, and the caller handed a clean error instead of crashing the server.
//
// A single option installs both a unary and a stream interceptor, so panics in
// streaming RPCs are covered too.
func WithRecovery(logger *slog.Logger) grpc.ServerOption {
	return func(s *grpc.Server) {
		s.UnaryInterceptor(recoveryUnary(logger))
		s.StreamInterceptor(recoveryStream(logger))
	}
}

// recoveryUnary logs the stack of any panic in a unary RPC and returns
// codes.Internal so the client receives a clean error rather than a bogus
// success.
func recoveryUnary(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("recovered from panic",
					"rpc", info.FullMethod,
					"error", rec,
					"trace", string(debug.Stack()),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		resp, err = handler(ctx, req)
		return
	}
}

// recoveryStream logs the stack of any panic in a streaming RPC and returns
// codes.Internal.
func recoveryStream(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("recovered from panic",
					"rpc", info.FullMethod,
					"error", rec,
					"trace", string(debug.Stack()),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		err = handler(srv, stream)
		return
	}
}
