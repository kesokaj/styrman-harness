---
owner: changelog
date: 2026-09-03
scope: changeset:09e0aa231bad
---

# Complete harness Go binary implementation with dispatch and gate logic

## Summary
Implemented `cmd/harness` Go binary with five internal packages to replace twelve scripts
(setup.sh ×3, merge-settings.py, derive-agy-agent.py, trail-check.sh, subagent-check.sh,
resume-context.sh, destructive-guard.py, format-hook.sh, feedback-capture.sh, statusline.sh).

Includes: `harness install`, `doctor`, `hook <event>`, `dispatch`, `resume`, `statusline`
subcommands; proportional Definition of Done (resumability rows always fire, accountability
rows fire on thresholds); artifact matching by frontmatter (`owner:`, `date:`) not filename;
secret scanning gate; cost ledger for dispatch runs; verifier as script with model escalation.

## Type
task

## Files changed
29 files: 11 script deletions (tracked), 19 modified (agents, skills, configs, test fixtures),
45 untracked new Go source files. Primary changes:
- `cmd/harness/main.go` — entry point, subcommand router
- `internal/trail/` — diff analysis, scope computation, owed-rows logic, artifact matching
- `internal/gate/` — secret scanning, formatter resolution, resume context
- `internal/dispatch/` — row router, script & model execution, cost ledger
- `internal/install/` — binary/skills/agents deployment, drift detection
- `internal/hookio/` — payload parsing for both vendors
- `claude/agents/*.md`, `shared/skills/*/` — updated to call `harness` subcommands
- `.agents/harness.json` — config with thresholds and verify commands

## Verification
`go test -race ./...` and `go vet ./...` pass; `gofmt` reports no changes.
Binary cross-compiles arm64+amd64. Five scripts tested as stubs; twelve deleted scripts removed
from git index. Review in `.agents/docs/reviews/2026-09-03-review-harness-go-binary.md`.

## Requirements
R2 (hooks ≤2.6 ms), R3 (drift detection), R4 (gates inform), R6 (no model spend on
zero-judgment rows). Review flagged high-severity issues requiring fixes before merge.
