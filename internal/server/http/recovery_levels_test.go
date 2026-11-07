// Package httpserver_test proves that a panic anywhere in the request path is
// recovered by the middleware and that its full stack trace is logged.
//
// A panic that originates deep in the infrastructure layer propagates up
// through the application layer and the transport handler before the outermost
// Recovery middleware catches it, returns a clean 500, and logs the trace. This
// file exercises that at four distinct depths, including the real /ping API
// wired through the real production router.
package httpserver_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yourname/empty_project/internal/application"
	"github.com/yourname/empty_project/internal/domain"
	"github.com/yourname/empty_project/internal/transport/http/handlers"
	"github.com/yourname/empty_project/internal/transport/http/middleware"

	httpserver "github.com/yourname/empty_project/internal/server/http"
)

// --- Panicking stand-ins, one per layer of the request path. ---------------

// panickingRepo is an infrastructure adapter that panics inside its own
// methods, simulating a panic that originates at the persistence layer.
type panickingRepo struct{}

func (panickingRepo) Save(_ context.Context, _ domain.Item) error {
	panic("boom: infrastructure layer")
}

func (panickingRepo) Get(_ context.Context, _ string) (domain.Item, error) {
	panic("boom: infrastructure layer")
}

// panickingUseCase is an application-service stand-in whose orchestration logic
// panics, simulating a panic that originates at the application layer.
type panickingUseCase struct{}

func (panickingUseCase) Create(_ context.Context) { panic("boom: application layer") }

// --- Harness ---------------------------------------------------------------

// captureLogger returns a JSON slog.Logger and the buffer it writes to so the
// emitted log lines can be inspected.
func captureLogger(t *testing.T) (*slog.Logger, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return logger, &buf
}

// runRequest invokes h (already wrapped with Recovery) for GET path and returns
// the status code the client received.
func runRequest(h http.Handler, path string) int {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	h.ServeHTTP(rec, req)
	return rec.Code
}

// traces returns every "trace" field logged by a "recovered from panic" record
// found in buf, in order.
func traces(t *testing.T, buf *bytes.Buffer) []string {
	t.Helper()
	dec := json.NewDecoder(buf)
	var out []string
	for dec.More() {
		var rec map[string]any
		if err := dec.Decode(&rec); err != nil {
			break
		}
		if rec["msg"] == "recovered from panic" {
			if tr, ok := rec["trace"].(string); ok {
				out = append(out, tr)
			}
		}
	}
	return out
}

// check asserts that a panic through h was recovered to a 500 and that a
// stack trace mentioning wantFrame was logged, then prints that trace.
func check(t *testing.T, title string, h http.Handler, wantFrame string) {
	t.Helper()
	logger, buf := captureLogger(t)
	h = middleware.Recovery(logger)(h)

	code := runRequest(h, "/")
	if code != http.StatusInternalServerError {
		t.Fatalf("[%s] status = %d, want %d", title, code, http.StatusInternalServerError)
	}

	got := traces(t, buf)
	if len(got) == 0 {
		t.Fatalf("[%s] expected a logged stack trace, none found. Full log:\n%s", title, buf.String())
	}
	trace := got[len(got)-1]
	if !strings.Contains(trace, wantFrame) {
		t.Errorf("[%s] stack trace does not mention originating frame %q:\n%s", title, wantFrame, trace)
	}

	t.Logf("\n===== %s =====\nclient received: HTTP %d\nstack trace logged by recovery middleware:\n%s",
		title, code, trace)
}

// TestRecoveryCatchesPanicAtEveryLayer demonstrates that a panic originating at
// any layer of the request path is recovered and its full stack trace logged.
func TestRecoveryCatchesPanicAtEveryLayer(t *testing.T) {
	// 1. Panic in the transport handler itself.
	check(t, "transport handler panics",
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			panic("boom: transport handler")
		}),
		"recovery_levels_test.go") // the panic site lives in this test

	// 2. Panic originating in the application layer.
	check(t, "application layer panics",
		http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			panickingUseCase{}.Create(r.Context())
		}),
		"panickingUseCase.Create")

	// 3. Panic originating in the infrastructure layer, travelling up through the
	//    real application.UseCase — so the trace shows infra -> application ->
	//    transport -> recovery.
	infra := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		uc := application.NewUseCase(panickingRepo{})
		_ = uc.Create(r.Context(), domain.Item{ID: "1", Name: "x"})
	})
	check(t, "infrastructure layer panics", infra, "panickingRepo.Save")

	// 4. The real /ping API, wired through the real production router, panics at
	//    the infrastructure layer.
	logger, buf := captureLogger(t)
	uc := application.NewUseCase(panickingRepo{})
	router := httpserver.NewRouter(logger, &httpserver.Routes{
		Health: handlers.NewHealthHandler(),
		Ping:   handlers.NewPingHandler(uc),
	})

	code := runRequest(router, "/ping")
	if code != http.StatusInternalServerError {
		t.Fatalf("[real /ping] status = %d, want %d", code, http.StatusInternalServerError)
	}
	got := traces(t, buf)
	if len(got) == 0 {
		t.Fatalf("[real /ping] expected a logged stack trace, none found. Full log:\n%s", buf.String())
	}
	trace := got[len(got)-1]
	for _, want := range []string{
		"handlers.(*PingHandler).Ping",
		"application.(*UseCase).Create",
		"panickingRepo.Save",
	} {
		if !strings.Contains(trace, want) {
			t.Errorf("[real /ping] stack trace does not show frame %q:\n%s", want, trace)
		}
	}
	t.Logf("\n===== real /ping panics at infrastructure layer =====\nclient received: HTTP %d\nstack trace:\n%s",
		code, trace)
}
