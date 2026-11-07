---
description: Run the project quality gate (format, vet, lint) and report findings.
---

Run the project's quality gate and report the results concisely:

1. `gofmt -l .` — list unformatted files; if any are listed, run `gofmt -w .`
   and report which changed.
2. `go vet ./...`
3. `make lint` — installs and runs golangci-lint.

Report findings grouped by severity with `file:line`. If everything is clean,
say so.
