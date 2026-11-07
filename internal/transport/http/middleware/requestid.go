package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// headerRequestID is the header carrying the request id in and out.
const headerRequestID = "X-Request-ID"

type requestIDKey struct{}

// requestIDCtxKey is the context key under which the request id is stored.
var requestIDCtxKey = requestIDKey{}

// WithRequestID returns middleware that assigns a request id (reusing the inbound
// X-Request-ID header when present), stores it in the request context, and
// echoes it in the response.
func WithRequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(headerRequestID)
			if id == "" {
				id = newRequestID()
			}
			w.Header().Set(headerRequestID, id)

			ctx := context.WithValue(r.Context(), requestIDCtxKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestIDFromContext returns the request id set by WithRequestID, or "".
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDCtxKey).(string); ok {
		return v
	}
	return ""
}

// newRequestID returns a random 16-byte hex id.
func newRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
