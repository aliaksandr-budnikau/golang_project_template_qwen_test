# Makefile for empty_project. Run `make` (or `make help`) to list targets.

.DEFAULT_GOAL := help

BINARY  := empty_project
PKG     := ./cmd/empty_project
BIN_DIR := bin

.PHONY: help build build-grpc run test test-race cover fmt vet lint proto tidy generate docker-build docker-run clean all

help: ## Show this help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

build: ## Build all packages (HTTP-only, default tags)
	@go build ./...

build-grpc: proto ## Build the binary with gRPC enabled
	@go build -tags grpc -o $(BIN_DIR)/$(BINARY) $(PKG)

run: ## Run the HTTP-only service locally
	@go run $(PKG)

test: ## Run unit tests
	@go test ./...

test-race: ## Run tests with the race detector
	@go test -race ./...

cover: ## Generate an HTML coverage report
	@go test -covermode=atomic -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo ">> coverage written to coverage.html"

fmt: ## Format the code
	@gofmt -w .

vet: ## Run go vet
	@go vet ./...

lint: ## Run golangci-lint v2 (installs it if missing)
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@golangci-lint run

proto: ## Generate gRPC/protobuf code and tidy modules
	@./scripts/gen.sh
	@go mod tidy

tidy: ## Run go mod tidy
	@go mod tidy

generate: ## Run go generate ./...
	@go generate ./...

docker-build: ## Build the Docker image
	@docker build -f build/docker/Dockerfile -t $(BINARY):latest .

docker-run: ## Build and run via docker compose
	@docker compose up --build

clean: ## Remove build artifacts
	@rm -rf $(BIN_DIR) coverage.out coverage.html

all: fmt vet test-race lint build ## Format, vet, test, lint, and build
