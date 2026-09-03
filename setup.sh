#!/usr/bin/env bash
# Bootstrap only. The installer is `harness install` — this exists because the
# binary cannot install itself before it has been built once.
#
#   ./setup.sh                        deploy both targets
#   ./setup.sh --dry-run              show what would happen
#   ./setup.sh --target claude        deploy one target
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")"
if [ -d .git ] && [ -d .githooks ]; then
  git config core.hooksPath .githooks 2>/dev/null || true
fi

exec go run ./cmd/harness install "$@"
