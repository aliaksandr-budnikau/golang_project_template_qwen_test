package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// Server wraps an *http.Server with a bounded graceful-shutdown timeout.
type Server struct {
	srv             *http.Server
	shutdownTimeout time.Duration
}

// NewServer constructs an HTTP Server that will listen on addr serving handler.
func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
		shutdownTimeout: 10 * time.Second,
	}
}

// Start serves HTTP until ctx is cancelled, then shuts the server down
// gracefully within the configured timeout. It returns nil on a clean stop and
// the ListenAndServe error otherwise.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		return s.srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
