# Project structure

This is the canonical map of the repository: every directory, every file, and
why it is where it is. It complements [`README.md`](README.md) (quick start) and
[`docs/architecture.md`](docs/architecture.md) (the layering rationale).

## Layering

The code is organised as a layered, dependency-inverted (hexagonal / clean)
architecture. Dependencies point **inward**:

```
    transport ──▶ application ──▶ domain ◀── infrastructure
        │             │                  │
   server      cmd/empty_project (composition root)
   config      telemetry
   pkg, api
```

- `domain` imports nothing internal.
- `application` imports only `domain` ports.
- `infrastructure` and `transport` import `application` / `domain`.
- `server` imports `transport` + `application`.
- `cmd/empty_project` imports everything (it is the only place that builds
  concrete types and wires them together).

## Full tree

```
.
├── .claude/                    # Claude Code configuration for this repo
│   ├── agents/
│   │   ├── api-designer.md     #   reviews REST/gRPC contract changes
│   │   └── go-reviewer.md      #   reviews Go code changes
│   ├── commands/
│   │   ├── lint.md             #   /lint  — run the quality gate
│   │   └── new-handler.md      #   /new-handler <name> — scaffold an HTTP handler
│   ├── settings.json           #   team settings (checked in)
│   └── settings.local.json.example  # copy to settings.local.json (git-ignored)
│
├── .editorconfig               # editor settings (tabs for .go/Makefile, 2-space otherwise)
├── .gitignore                  # ignore binaries, coverage, generated *.pb.go, local settings
├── .golangci.yml               # golangci-lint configuration (linters + settings)
│
├── api/                        # API contracts (the public surface)
│   ├── openapi/
│   │   └── openapi.yaml        #   REST contract (OpenAPI 3.0): /healthz, /readyz, /ping
│   └── proto/                  #   gRPC contract + buf config
│       ├── buf.yaml            #     buf module + lint/breaking config
│       ├── buf.gen.yaml        #     buf codegen plugins (protoc-gen-go, -grpc)
│       └── v1/
│           └── service.proto   #     proto3 service + messages (generated *.pb.go is git-ignored)
│
├── build/                      # build + deploy assets (kept out of the Go tree)
│   ├── docker/
│   │   ├── .dockerignore
│   │   └── Dockerfile          #     multi-stage, static binary, non-root, HEALTHCHECK
│   └── scripts/
│       └── entrypoint.sh       #   optional pre-start hook (exec "$@")
│
├── CLAUDE.md                   # guidance for Claude Code working in this repo
├── README.md                   # project overview + quick start
│
├── cmd/
│   └── empty_project/          # composition root + entrypoints
│       ├── main.go             #   HTTP-only entrypoint            //go:build !grpc
│       └── main_grpc.go        #   HTTP+gRPC entrypoint            //go:build grpc
│
├── docker-compose.yml          # run the image; maps 8080 (HTTP) and 9090 (gRPC)
│
├── docs/                       # human-readable design docs
│   ├── adr/
│   │   └── 0001-record-architecture-decisions.md   # how to write ADRs
│   └── architecture.md         # the layering + dependency rule (rationale)
│
├── go.mod                      # module github.com/yourname/empty_project (placeholder)
│
├── internal/                   # private code (not importable outside this module)
│   ├── application/            # use cases; orchestrate the domain, depend on ports
│   │   ├── usecase.go
│   │   └── usecase_test.go
│   ├── config/                 # load + validate configuration from the environment
│   │   ├── config.go
│   │   └── config_test.go
│   ├── domain/                 # entities + ports; NO internal dependencies
│   │   └── entity.go           #   Item + Repository interface
│   ├── infrastructure/         # concrete adapters implementing domain ports
│   │   ├── external/
│   │   │   └── client.go       #     sample external-service client
│   │   └── persistence/
│   │       └── repository.go   #     in-memory Repository implementation
│   ├── server/                 # assembles + runs a transport (router, middleware, lifecycle)
│   │   ├── grpc/               #   gRPC server wrapper (opt-in)
│   │   │   ├── server.go            # //go:build grpc        real impl
│   │   │   └── server_default.go    # //go:build !grpc       no-op stub (keeps default build green)
│   │   └── http/
│   │       ├── router.go       #     middleware chain + method-aware mux
│   │       ├── server.go       #     graceful-shutdown wrapper
│   │       └── swagger/        #     opt-in Swagger UI (go:embed; on when SERVICE_SWAGGER_ENABLED=true)
│   │           ├── swagger.go  #       handler: /swagger UI + /swagger/openapi.yaml
│   │           ├── ui.html     #       thin mount; third-party swagger-ui renders the spec
│   │           └── openapi.yaml#       embedded copy of api/openapi/openapi.yaml (refresh: make generate)
│   ├── telemetry/
│   │   └── logger.go           #   structured logging (slog)
│   └── transport/              # protocol boundaries (thin: translate to/from application)
│       ├── grpc/               #   gRPC service implementations (opt-in)
│       │   ├── server.go            # //go:build grpc        real impl
│       │   └── server_default.go    # //go:build !grpc       no-op stub
│       └── http/
│           ├── handlers/       #     HTTP handlers (thin)
│           │   ├── health.go   #       /healthz, /readyz
│           │   ├── ping.go     #       /ping  (HTTP → application → domain → infra)
│           │   └── helpers.go  #       writeJSON helper
│           └── middleware/
│               ├── logging.go  #        request logging
│               ├── recovery.go #        panic recovery
│               └── requestid.go#        request-id propagation
│
├── LICENSE
├── Makefile                    # build/run/test/lint/proto/docker targets (`make help`)
│
├── pkg/                        # reusable, importable-by-others code
│   └── util/
│       └── util.go
│
├── scripts/
│   └── gen.sh                  # codegen entry point (buf, falling back to protoc); used by `make proto`
│
├── STRUCTURE.md                # this file
│
└── test/
    └── integration/            # end-to-end tests
        └── .gitkeep           #   keeps the (currently empty) dir in version control
```

## Build variants

| Variant | Command | Includes gRPC |
|---------|---------|---------------|
| Default | `go build ./...`, `make build` | No — HTTP-only, pure stdlib, zero external deps |
| gRPC | `make proto` then `go build -tags grpc` | Yes — needs generated stubs + `google.golang.org/grpc` |

The gRPC code is **build-tagged** so the default build never needs the generated
stubs or the gRPC dependency. Each tagged file has a paired stub:

- `//go:build grpc` — the real implementation (`server.go`, `main_grpc.go`).
- `//go:build !grpc` — a no-op stub (`server_default.go`, `main.go`) that keeps
  `go build ./...` green without the tag.

The `cmd/empty_project` entrypoints are mutually exclusive for the same reason:
exactly one of `main.go` / `main_grpc.go` is compiled at a time.

## Conventions (where things go)

- **New port** → `internal/domain/`.
- **New use case** → `internal/application/`.
- **New adapter** (persistence / external) → `internal/infrastructure/`.
- **New HTTP handler** → `internal/transport/http/handlers/` (thin) + register in
  `internal/server/http/router.go` + wire in `cmd/empty_project/main.go`.
  Use the `/new-handler <name>` slash command to scaffold this.
- **REST contract** → `api/openapi/openapi.yaml`.
- **gRPC contract** → `api/proto/v1/*.proto`.
- **Design decisions** → a new record in `docs/adr/`.
- **Tests** → one `<name>_test.go` beside each `<name>.go`, table-driven.

Keep the import direction inward. See `docs/architecture.md` for the rationale
and `docs/adr/` for the recorded decisions.
