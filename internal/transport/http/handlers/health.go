package handlers

import (
	"net/http"
	"time"
)

// HealthHandler reports service liveness and readiness.
type HealthHandler struct {
	startedAt time.Time
}

// NewHealthHandler constructs a HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{startedAt: time.Now()}
}

// Liveness handles GET /healthz — it reports that the process is alive.
func (h *HealthHandler) Liveness(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Readiness handles GET /readyz — it reports the startup time as a readiness
// proxy.
func (h *HealthHandler) Readiness(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ready",
		"startedAt": h.startedAt.UTC().Format(time.RFC3339),
		"uptime":    time.Since(h.startedAt).String(),
	})
}
