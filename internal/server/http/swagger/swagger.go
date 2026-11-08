// Package swagger serves a third-party Swagger UI and the OpenAPI spec, in
// process, so the API can be explored by clicking "Try it out".
//
// Both assets are embedded at build time, so the feature is self-contained and
// adds no runtime or build dependency. It is wired into the router only when
// SERVICE_SWAGGER_ENABLED is set, keeping the production surface documentation-free.
//
// Refresh the embedded spec from its source of truth with `go generate ./...`.
package swagger

import (
	_ "embed"
	"net/http"
)

//go:generate cp ../../../../api/openapi/openapi.yaml ./openapi.yaml

//go:embed ui.html
var uiHTML []byte

//go:embed openapi.yaml
var spec []byte

// Handler returns an http.Handler that serves the Swagger UI and the OpenAPI
// spec. It is meant to be mounted at /swagger by the router:
//
//	/swagger            -> UI (301 to /swagger/)
//	/swagger/           -> UI
//	/swagger/openapi.yaml -> embedded spec (fetched by the UI, same origin)
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/swagger":
			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
		case "/swagger/", "/swagger/index.html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(uiHTML)
		case "/swagger/openapi.yaml":
			w.Header().Set("Content-Type", "application/yaml")
			_, _ = w.Write(spec)
		default:
			http.NotFound(w, r)
		}
	})
}
