package handlers

import (
	"net/http"

	"github.com/yourname/empty_project/internal/application"
	"github.com/yourname/empty_project/internal/domain"
)

// PingHandler demonstrates the full request path: HTTP -> application -> domain
// -> infrastructure and back.
type PingHandler struct {
	useCase *application.UseCase
}

// NewPingHandler constructs a PingHandler.
func NewPingHandler(uc *application.UseCase) *PingHandler {
	return &PingHandler{useCase: uc}
}

// pingResponse is the JSON body returned by Ping.
type pingResponse struct {
	Status string      `json:"status"`
	Item   domain.Item `json:"item"`
}

// Ping stores a sample item through the use case, reads it back, and echoes it.
// It handles GET /ping.
func (h *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	item := domain.Item{ID: "sample", Name: "pong"}
	if err := h.useCase.Create(r.Context(), item); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	got, err := h.useCase.Retrieve(r.Context(), item.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, pingResponse{Status: "ok", Item: got})
}
