package handlers

import (
	"net/http"

	"github.com/yourname/empty_project/internal/application"
)

// RandomHandler serves GET /random: it imitates a short (~1s) unit of work and
// answers with a random number between 1 and 154.
type RandomHandler struct {
	useCase *application.RandomUseCase
}

// NewRandomHandler constructs a RandomHandler.
func NewRandomHandler(uc *application.RandomUseCase) *RandomHandler {
	return &RandomHandler{useCase: uc}
}

// randomResponse is the JSON body returned by Random.
type randomResponse struct {
	Status string `json:"status"`
	Number int    `json:"number"`
}

// Random runs the use case and echoes its result. It handles GET /random.
// The response is marked Cache-Control: no-store — every call yields a fresh
// number, so it must not be cached by intermediaries.
func (h *RandomHandler) Random(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")

	number, err := h.useCase.Generate(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, randomResponse{Status: "ok", Number: number})
}
