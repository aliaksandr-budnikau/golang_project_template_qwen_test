package swagger

import (
	"bytes"
	"os"
	"testing"
)

// TestEmbeddedSpecIsInSync guards the //go:generate copy: the Swagger UI serves
// the embedded openapi.yaml, so if api/openapi/openapi.yaml (the source of
// truth) was changed without running `go generate ./...`, the UI would show a
// stale contract. This test fails in that case.
func TestEmbeddedSpecIsInSync(t *testing.T) {
	// go test runs with the package directory as the working directory.
	src, err := os.ReadFile("../../../../api/openapi/openapi.yaml")
	if err != nil {
		t.Fatalf("reading source spec: %v", err)
	}

	if !bytes.Equal(src, spec) {
		t.Fatal("embedded openapi.yaml is stale: run `go generate ./...` to refresh it from api/openapi/openapi.yaml")
	}
}
