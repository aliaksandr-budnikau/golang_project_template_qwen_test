package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourname/empty_project/internal/application"
)

func TestRandomHandler(t *testing.T) {
	tests := []struct {
		name string
		req  func() *http.Request
		want func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "returns a number in range",
			req:  func() *http.Request { return httptest.NewRequest(http.MethodGet, "/random", nil) },
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				if got := resp.Code; got != http.StatusOK {
					t.Fatalf("status = %d, want %d", got, http.StatusOK)
				}

				var body randomResponse
				if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
					t.Fatalf("decoding response: %v", err)
				}
				if body.Status != "ok" {
					t.Errorf("status field = %q, want %q", body.Status, "ok")
				}
				if body.Number < 1 || body.Number > 154 {
					t.Errorf("number = %d, want in [1, 154]", body.Number)
				}
				if ct := resp.Header().Get("Content-Type"); ct != "application/json" {
					t.Errorf("Content-Type = %q, want %q", ct, "application/json")
				}
				if cc := resp.Header().Get("Cache-Control"); cc != "no-store" {
					t.Errorf("Cache-Control = %q, want %q", cc, "no-store")
				}
			},
		},
		{
			name: "cancelled request context yields 500",
			req: func() *http.Request {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return httptest.NewRequest(http.MethodGet, "/random", nil).WithContext(ctx)
			},
			want: func(t *testing.T, resp *httptest.ResponseRecorder) {
				t.Helper()
				if got := resp.Code; got != http.StatusInternalServerError {
					t.Fatalf("status = %d, want %d", got, http.StatusInternalServerError)
				}
			},
		},
	}

	// A tiny work delay keeps the tests fast; the delay itself is covered by
	// the application-layer test.
	h := NewRandomHandler(application.NewRandomUseCase(application.WithWorkDelay(time.Millisecond)))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.Random(rec, tt.req())
			tt.want(t, rec)
		})
	}
}
