// Package httpserver assembles the HTTP/REST server: the middleware chain, the
// router, and the server lifecycle. It is named httpserver to avoid colliding
// with the net/http import.
package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/yourname/empty_project/internal/transport/http/handlers"
	"github.com/yourname/empty_project/internal/transport/http/middleware"
)

// Routes holds the handlers wired into the router.
//
// Swagger is optional: when non-nil it mounts the Swagger UI and OpenAPI spec at
// /swagger and /swagger/openapi.yaml (set by the composition root only when
// SERVICE_SWAGGER_ENABLED is on).
type Routes struct {
	Health  *handlers.HealthHandler
	Ping    *handlers.PingHandler
	Swagger http.Handler
}

// NewRouter builds the HTTP handler: the middleware chain wrapping a
// method-aware http.ServeMux.
//
// Middleware order (outermost first): recovery -> request id -> logging -> mux.
func NewRouter(logger *slog.Logger, r *Routes) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", r.Health.Liveness)
	mux.HandleFunc("GET /readyz", r.Health.Readiness)
	mux.HandleFunc("GET /ping", r.Ping.Ping)

	if r.Swagger != nil {
		mux.Handle("/swagger", r.Swagger)
		mux.Handle("/swagger/", r.Swagger)
	}

	var h http.Handler = mux
	h = middleware.Logging(logger)(h)
	h = middleware.WithRequestID()(h)
	h = middleware.Recovery(logger)(h)
	return h
}
