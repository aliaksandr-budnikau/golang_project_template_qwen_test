# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.26-alpine AS build

WORKDIR /src

# Cache dependencies first so the module cache is only rebuilt when go.mod changes.
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Build a static binary.
COPY . .
ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED}
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/empty_project ./cmd/empty_project

# ---- Runtime stage ----
FROM alpine:3.20 AS runtime

RUN apk add --no-cache curl \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

COPY --from=build --chown=appuser:appgroup /out/empty_project /usr/local/bin/empty_project

USER appuser
EXPOSE 8080 9090

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD curl -fsS http://localhost:8080/healthz || exit 1

# The binary handles SIGINT/SIGTERM itself. Use build/scripts/entrypoint.sh as
# ENTRYPOINT instead if you need a pre-start hook.
ENTRYPOINT ["/usr/local/bin/empty_project"]
