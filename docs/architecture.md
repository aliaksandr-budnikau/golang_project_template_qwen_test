# Architecture

This service uses a layered, dependency-inverted structure (a light form of
hexagonal / clean architecture). Dependencies point **inward**:

```
    transport ──▶ application ──▶ domain ◀── infrastructure
```

## Packages

- **domain** — business entities and the ports (interfaces). No dependencies on
  other internal packages.
- **application** — use cases that orchestrate the domain. Depends only on domain
  ports, never on concrete implementations.
- **infrastructure** — concrete adapters (persistence, external clients) that
  implement the domain ports.
- **transport** — the protocol boundaries (HTTP, gRPC). Thin: translate to and
  from the application layer.
- **server** — assembles and runs a transport (router, middleware, lifecycle).
- **config** — loads and validates configuration.
- **telemetry** — logging / tracing / metrics.
- **cmd/empty_project** — the composition root: constructs concrete types, wires
  them together, and runs the service with graceful shutdown.

## Dependency rule

Outer layers may import inner layers, never the reverse:

- `domain` imports nothing.
- `application` imports `domain`.
- `infrastructure` and `transport` import `application` / `domain`.
- `server` imports `transport` and `application`.
- the composition root (`cmd/empty_project`) imports everything.

## Build variants

- **Default** (`go build ./...`, `make build`): HTTP-only, pure standard library,
  zero external dependencies.
- **gRPC** (`make proto` then `go build -tags grpc`): adds the gRPC server. The
  generated stubs are produced into `api/proto/v1/` and are git-ignored, so the
  default build never requires them.

See `docs/adr/` for architecture decision records.
