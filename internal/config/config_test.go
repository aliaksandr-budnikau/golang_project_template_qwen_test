package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	// An empty value makes getenv() fall back to its default, so loading with
	// every variable empty yields the documented defaults.
	for _, k := range []string{
		"SERVICE_NAME", "SERVICE_HTTP_ADDR", "SERVICE_GRPC_ADDR",
		"SERVICE_GRPC_ENABLED", "SERVICE_LOG_LEVEL",
	} {
		t.Setenv(k, "")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg.ServiceName != "empty_project" {
		t.Errorf("ServiceName = %q, want %q", cfg.ServiceName, "empty_project")
	}
	if cfg.HTTP.Addr != ":8080" {
		t.Errorf("HTTP.Addr = %q, want %q", cfg.HTTP.Addr, ":8080")
	}
	if cfg.GRPC.Addr != ":9090" {
		t.Errorf("GRPC.Addr = %q, want %q", cfg.GRPC.Addr, ":9090")
	}
	if cfg.GRPC.Enabled {
		t.Errorf("GRPC.Enabled = true, want false")
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "info")
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("SERVICE_NAME", "prod")
	t.Setenv("SERVICE_HTTP_ADDR", ":9000")
	t.Setenv("SERVICE_GRPC_ADDR", ":9091")
	t.Setenv("SERVICE_GRPC_ENABLED", "true")
	t.Setenv("SERVICE_LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}

	if cfg.ServiceName != "prod" {
		t.Errorf("ServiceName = %q, want %q", cfg.ServiceName, "prod")
	}
	if cfg.HTTP.Addr != ":9000" {
		t.Errorf("HTTP.Addr = %q, want %q", cfg.HTTP.Addr, ":9000")
	}
	if cfg.GRPC.Addr != ":9091" {
		t.Errorf("GRPC.Addr = %q, want %q", cfg.GRPC.Addr, ":9091")
	}
	if !cfg.GRPC.Enabled {
		t.Errorf("GRPC.Enabled = false, want true")
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
}

func TestLoadInvalidGrpcEnabled(t *testing.T) {
	t.Setenv("SERVICE_GRPC_ENABLED", "notabool")
	if _, err := Load(); err == nil {
		t.Fatalf("Load() expected an error for invalid SERVICE_GRPC_ENABLED, got nil")
	}
}
