# empty_project

A best-practice starting point for a Go HTTP/REST + gRPC microservice.

The **default build is HTTP-only and dependency-free** (pure standard library).
The gRPC layer is opt-in and generated, so `go build ./...` works out of the box.

## Requirements

- Go 1.26+
- (Optional, for gRPC) `buf` or `protoc` + the Go plugins. See *gRPC* below.
- (Optional) `make`, `docker`, `docker compose`

## Quick start

```sh
make help          # list all targets
make build         # go build ./...  (HTTP only)
make run           # start the HTTP server on :8080
make test          # go test ./...
make all           # fmt, vet, test -race, lint, build
```

Then:

```sh
curl -s localhost:8080/healthz
curl -s localhost:8080/ping
```

Stop the server with `Ctrl-C` (it shuts down gracefully).

## Layout

```
cmd/empty_project/        composition root + entrypoints (main.go, main_grpc.go)
internal/
  config/                 configuration loading/validation
  telemetry/              structured logging
  domain/                 entities + ports (no deps)
  application/            use cases (depends on domain ports)
  infrastructure/         concrete adapters (persistence, external clients)
  transport/
    http/                 handlers + middleware
    grpc/                 gRPC service implementations (build-tagged)
  server/
    http/                 router + server lifecycle (+ opt-in Swagger UI)
    grpc/                 gRPC server wrapper (build-tagged)
pkg/                      reusable, importable code
api/
  openapi/                REST contract
  proto/                  gRPC contract + buf config
build/docker/             Dockerfile + .dockerignore
scripts/                  code generation
test/integration/         end-to-end tests
docs/                     architecture + ADRs
```

Dependencies point inward: `transport → application → domain ← infrastructure`.
See `docs/architecture.md`.

## Configuration

All configuration comes from the environment (defaults in parentheses):

| Variable               | Default         | Meaning                     |
|------------------------|-----------------|-----------------------------|
| `SERVICE_NAME`         | `empty_project` | service name (logging)      |
| `SERVICE_HTTP_ADDR`    | `:8080`         | HTTP listen address         |
| `SERVICE_GRPC_ADDR`    | `:9090`         | gRPC listen address         |
| `SERVICE_GRPC_ENABLED` | `false`         | start the gRPC server       |
| `SERVICE_SWAGGER_ENABLED` | `false`     | serve the Swagger UI at `/swagger` |
| `SERVICE_LOG_LEVEL`    | `info`          | `debug`/`info`/`warn`/`error` |

## Swagger UI

A Swagger UI is served in-process at `/swagger`, so the REST surface can be
explored and run without a client. It is **off by default** — production exposes
no documentation surface — and enabled with `SERVICE_SWAGGER_ENABLED=true`.

```sh
make docker-run      # docker compose enables it, so the UI is on locally
# open http://localhost:8080/swagger
```

The UI is a thin mount for the third-party
[swagger-ui](https://github.com/swagger-api/swagger-ui); the OpenAPI spec is
embedded from `api/openapi/openapi.yaml` (refresh the copy with `make generate`).

## gRPC (opt-in)

The gRPC code is build-tagged (`//go:build grpc`) and generated, so the default
build needs neither the generated stubs nor the gRPC dependencies. To build with
gRPC:

```sh
make proto        # generate api/proto/v1/*.pb.go and tidy modules
make build-grpc   # go build -tags grpc -o bin/empty_project ./cmd/empty_project
SERVICE_GRPC_ENABLED=true ./bin/empty_project
```

Generation requires `buf` (preferred) or `protoc` + the Go plugins:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Docker

```sh
make docker-build     # build the image
make docker-run       # build + run via docker compose (maps 8080/9090)
```

## Tooling

- `make lint` runs [golangci-lint](https://golangci-lint.run/) (installs it if
  missing); config in `.golangci.yml`.
- `make cover` writes `coverage.html`.
- `make clean` removes build artifacts.

## Module path

The module is `github.com/yourname/empty_project` — a placeholder. Rename it when
you have a real repository:

```sh
go mod edit -module github.com/you/yourmodule
# then update the import paths, e.g.:
grep -rl 'github.com/yourname/empty_project' --include='*.go' . | xargs sed -i '' \
  's#github.com/yourname/empty_project#github.com/you/yourmodule#g'
make proto && make build   # regenerate gRPC code against the new path
```
