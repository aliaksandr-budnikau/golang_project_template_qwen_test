#!/usr/bin/env bash
#
# Optional signal-forwarding entrypoint. The binary already handles SIGINT and
# SIGTERM, so the Dockerfile uses the binary directly as ENTRYPOINT. Use this
# script only if you need to run a pre-start hook before exec-ing the binary.
set -euo pipefail

# Example pre-start hook (kept as a no-op):
# echo "starting $(date -u)"

exec "$@"
