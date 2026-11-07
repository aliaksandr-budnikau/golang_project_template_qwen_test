// Package telemetry provides observability primitives (logging, tracing, metrics).
package telemetry

import (
	"log/slog"
	"os"
	"strings"
)

// NewLogger builds a structured JSON *slog.Logger writing to stderr.
//
// level is one of "debug", "info", "warn", "error" (case-insensitive); any
// unrecognised value falls back to info.
func NewLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: lvl})
	return slog.New(handler)
}
