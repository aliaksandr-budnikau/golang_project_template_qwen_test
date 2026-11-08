//go:build grpc

// Command empty_project is the HTTP+gRPC entrypoint, compiled with -tags grpc.
// It is mutually exclusive with main.go (built without the tag).
//
// Generate the gRPC stubs first with `make proto`.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/yourname/empty_project/internal/application"
	"github.com/yourname/empty_project/internal/config"
	"github.com/yourname/empty_project/internal/infrastructure/persistence"
	grpcserver "github.com/yourname/empty_project/internal/server/grpc"
	httpserver "github.com/yourname/empty_project/internal/server/http"
	"github.com/yourname/empty_project/internal/server/http/swagger"
	"github.com/yourname/empty_project/internal/telemetry"
	grpcsvc "github.com/yourname/empty_project/internal/transport/grpc"
	"github.com/yourname/empty_project/internal/transport/http/handlers"
)

func main() {
	if err := runGRPC(); err != nil {
		slog.Default().Error("service exited with error", "err", err)
		os.Exit(1)
	}
}

// runGRPC runs both the HTTP and gRPC servers and shuts both down gracefully.
func runGRPC() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := telemetry.NewLogger(cfg.LogLevel)
	slog.SetDefault(logger)
	logger.Info("starting service (http+grpc)", "name", cfg.ServiceName)

	// --- Composition root ---
	repo := persistence.NewMemoryRepo()
	useCase := application.NewUseCase(repo)
	randomUseCase := application.NewRandomUseCase()

	// HTTP.
	health := handlers.NewHealthHandler()
	ping := handlers.NewPingHandler(useCase)
	random := handlers.NewRandomHandler(randomUseCase)
	routes := &httpserver.Routes{Health: health, Ping: ping, Random: random}
	if cfg.Swagger.Enabled {
		routes.Swagger = swagger.Handler()
	}
	httpHandler := httpserver.NewRouter(logger, routes)
	httpSrv := httpserver.NewServer(cfg.HTTP.Addr, httpHandler)

	// gRPC.
	grpcSrv, err := grpcserver.NewServer(cfg.GRPC.Addr, grpcserver.WithRecovery(logger))
	if err != nil {
		return err
	}
	grpcsvc.Register(grpcSrv)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := httpSrv.Start(ctx); err != nil {
			logger.Error("http server error", "err", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := grpcSrv.Serve(); err != nil {
			logger.Error("grpc server error", "err", err)
		}
	}()

	logger.Info("servers listening", "http", cfg.HTTP.Addr, "grpc", cfg.GRPC.Addr)

	<-ctx.Done()
	logger.Info("shutdown signal received")
	grpcSrv.Stop() // initiates a graceful stop; unblocks grpcSrv.Serve()
	wg.Wait()
	return nil
}
