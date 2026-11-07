// Package handlers contains the HTTP handlers — the HTTP boundary of the
// application. Handlers are thin: they translate HTTP into calls to the
// application layer and the results back into HTTP.
package handlers

import (
	"encoding/json"
	"net/http"
)

// writeJSON writes v as a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
