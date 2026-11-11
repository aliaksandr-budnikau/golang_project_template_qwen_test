# CLAUDE.md

Guidance for Claude Code when working in this repository.

## What this is

A Go HTTP/REST + gRPC microservice scaffold. See `README.md` for the full
overview and `docs/architecture.md` for the layering.

## Architecture (dependency rule)

Dependencies point **inward**:

```
transport → application → domain ← infrastructure
```

- `domain` imports nothing internal.
- `application` depends on `domain` ports (interfaces), never on implementations.
- `infrastructure` and `transport` depend on inner layers.
- `cmd/empty_project` is the composition root: it is the only place that
  constructs concrete types and wires them together.

When adding code, place it in the layer that matches its responsibility and keep
the import direction inward. New ports go in `domain`; new implementations go in
`infrastructure`.

## Build variants

- Default build (`make build` / `go build ./...`): **HTTP-only, zero external
  dependencies**.
- gRPC is opt-in and build-tagged: `make proto` then `go build -tags grpc`.
  Generated code lives in `api/proto/v1/` and is git-ignored.

Because of this, do **not** add gRPC (`google.golang.org/grpc`) or other
third-party imports to the default-build path; gate them behind `//go:build grpc`
with a `//go:build !grpc` stub, as `internal/server/grpc` and
`internal/transport/grpc` already do.

## Commands

Run these before committing:

```sh
make all        # fmt, vet, test -race, lint, build
```

Common ones: `make build`, `make test`, `make lint`, `make run`, `make proto`.

## Conventions

- Format with `gofmt` (run `make fmt`); the linter config is `.golangci.yml`.
- Handlers are thin and live in `internal/transport/http/handlers/`; business
  logic lives in `internal/application/`.
- Table-driven tests, one `<name>_test.go` beside each `<name>.go`.
- Use the provided slash commands: `/new-handler <name>` to scaffold a handler,
  `/lint` to run the quality gate. See `.claude/commands/`.
- Use the subagents in `.claude/agents/` (`go-reviewer`, `api-designer`).

## Before finishing a task

- Run  `make build` and `go test ./...` (and `make lint`) to confirm the default
  build is green.
- If you touched a `.proto`, run `make proto` and rebuild with `-tags grpc`.

