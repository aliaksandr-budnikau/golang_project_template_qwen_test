// Package config loads and validates runtime configuration from the environment.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all runtime configuration for the service.
type Config struct {
	ServiceName string
	HTTP        HTTPConfig
	GRPC        GRPCConfig
	Swagger     SwaggerConfig
	LogLevel    string
}

// HTTPConfig configures the HTTP/REST server.
type HTTPConfig struct {
	// Addr is the listen address, e.g. ":8080".
	Addr string
}

// GRPCConfig configures the gRPC server.
type GRPCConfig struct {
	// Enabled reports whether the gRPC server should start (only honoured in
	// builds compiled with -tags grpc).
	Enabled bool
	// Addr is the listen address, e.g. ":9090".
	Addr string
}

// SwaggerConfig configures the in-process Swagger UI.
//
// The UI is disabled by default so production builds expose no documentation
// surface; enable it for local development (e.g. docker compose) by setting
// SERVICE_SWAGGER_ENABLED=true.
type SwaggerConfig struct {
	// Enabled reports whether the Swagger UI and OpenAPI spec are served at
	// /swagger and /swagger/openapi.yaml.
	Enabled bool
}

// Load reads configuration from the environment, applying defaults for unset
// values.
//
// Recognised environment variables:
//
//	SERVICE_NAME          (default "empty_project")
//	SERVICE_HTTP_ADDR     (default ":8080")
//	SERVICE_GRPC_ADDR     (default ":9090")
//	SERVICE_GRPC_ENABLED  (default "false")
//	SERVICE_SWAGGER_ENABLED (default "false")
//	SERVICE_LOG_LEVEL     (default "info")
func Load() (Config, error) {
	cfg := Config{
		ServiceName: getenv("SERVICE_NAME", "empty_project"),
		HTTP:        HTTPConfig{Addr: getenv("SERVICE_HTTP_ADDR", ":8080")},
		GRPC:        GRPCConfig{Addr: getenv("SERVICE_GRPC_ADDR", ":9090")},
		LogLevel:    getenv("SERVICE_LOG_LEVEL", "info"),
	}

	if enabled := getenv("SERVICE_GRPC_ENABLED", "false"); enabled != "" {
		var err error
		cfg.GRPC.Enabled, err = strconv.ParseBool(enabled)
		if err != nil {
			return Config{}, fmt.Errorf("invalid SERVICE_GRPC_ENABLED %q: %w", enabled, err)
		}
	}

	if enabled := getenv("SERVICE_SWAGGER_ENABLED", "false"); enabled != "" {
		var err error
		cfg.Swagger.Enabled, err = strconv.ParseBool(enabled)
		if err != nil {
			return Config{}, fmt.Errorf("invalid SERVICE_SWAGGER_ENABLED %q: %w", enabled, err)
		}
	}

	return cfg, nil
}

// getenv returns the value of the environment variable key, or fallback if unset
// or empty.
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
