#!/usr/bin/env bash
#
# Generate gRPC/protobuf Go code. Prefers buf; falls back to protoc + plugins.
# Generated files are written next to the .proto in api/proto/v1/ and are
# git-ignored. Run this (or `make proto`) before building with -tags grpc.
set -euo pipefail

cd "$(dirname "$0")/.."
PROTO_DIR=api/proto/v1

if command -v buf >/dev/null 2>&1; then
  echo ">> generating with buf"
  ( cd api/proto && buf generate )
  exit 0
fi

echo ">> buf not found; falling back to protoc"
if ! command -v protoc >/dev/null 2>&1; then
  echo "error: neither 'buf' nor 'protoc' is available." >&2
  echo "  - Install buf:      https://buf.build/docs/installation"
  echo "  - Or install protoc plus the Go plugins, e.g.:"
  echo "      go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
  echo "      go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
  exit 1
fi

protoc \
  --proto_path="$PROTO_DIR" \
  --go_out="$PROTO_DIR" --go_opt=paths=source_relative \
  --go-grpc_out="$PROTO_DIR" --go-grpc_opt=paths=source_relative \
  service.proto

echo ">> generated code written to $PROTO_DIR"
