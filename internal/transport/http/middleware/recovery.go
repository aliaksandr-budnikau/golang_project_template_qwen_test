// Package middleware provides reusable net/http middleware.
package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery returns middleware that turns a panic in a downstream handler into a
// 500 response and logs the stack trace.
func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("recovered from panic",
						"error", rec,
						"trace", string(debug.Stack()),
					)
					// The response may have been partially written already, in
					// which case http.Error is a no-op.
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
