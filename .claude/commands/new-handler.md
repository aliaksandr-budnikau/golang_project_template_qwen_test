---
description: Scaffold a new HTTP handler and wire it into the router.
argument-hint: <handler-name>
---

Create a new HTTP handler named `$ARGUMENTS` (kebab-case) following the
project's conventions, mirroring `internal/transport/http/handlers/ping.go`:

1. Add `internal/transport/http/handlers/<name>.go` with a handler struct, a
   `New<Name>Handler` constructor, and a handler method.
2. Register the route in `internal/server/http/router.go` with the correct
   method and path, add it to the `Routes` struct, and construct it in the
   composition root `cmd/empty_project/main.go`.
3. If the handler needs business logic, add a method to
   `internal/application/usecase.go` (and a port in `internal/domain/` if it
   needs persistence or an external call).
4. Add a table-driven test `internal/transport/http/handlers/<name>_test.go`
   using `net/http/httptest`.
5. Run `make build` and `go test ./...` to confirm it compiles and passes.
