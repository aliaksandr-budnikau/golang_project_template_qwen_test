//go:build !grpc

// Command empty_project is the HTTP/REST entrypoint. It loads configuration,
// wires the dependency graph, and runs the HTTP server with graceful shutdown.
//
// For the HTTP+gRPC variant, build with the gRPC tag after generating the
// stubs (`make proto`):
//
//	go build -tags grpc ./cmd/empty_project
//
// See main_grpc.go.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourname/empty_project/internal/application"
	"github.com/yourname/empty_project/internal/config"
	"github.com/yourname/empty_project/internal/infrastructure/persistence"
	httpserver "github.com/yourname/empty_project/internal/server/http"
	"github.com/yourname/empty_project/internal/server/http/swagger"
	"github.com/yourname/empty_project/internal/telemetry"
	"github.com/yourname/empty_project/internal/transport/http/handlers"
)

func main() {
	if err := run(); err != nil {
		slog.Default().Error("service exited with error", "err", err)
		os.Exit(1)
	}
}

// run loads config, wires dependencies, and runs the HTTP server until a
// termination signal is received, then shuts it down gracefully.
func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := telemetry.NewLogger(cfg.LogLevel)
	slog.SetDefault(logger)
	logger.Info("starting service", "name", cfg.ServiceName)

	// --- Composition root: construct concrete types and wire them together. ---
	repo := persistence.NewMemoryRepo()
	useCase := application.NewUseCase(repo)
	randomUseCase := application.NewRandomUseCase()

	health := handlers.NewHealthHandler()
	ping := handlers.NewPingHandler(useCase)
	random := handlers.NewRandomHandler(randomUseCase)

	routes := &httpserver.Routes{Health: health, Ping: ping, Random: random}
	if cfg.Swagger.Enabled {
		routes.Swagger = swagger.Handler()
	}
	handler := httpserver.NewRouter(logger, routes)
	srv := httpserver.NewServer(cfg.HTTP.Addr, handler)

	// --- Run with graceful shutdown on SIGINT / SIGTERM. ---
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("http server listening", "addr", cfg.HTTP.Addr)
	return srv.Start(ctx)
}
