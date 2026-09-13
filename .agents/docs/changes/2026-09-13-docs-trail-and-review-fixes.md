---
owner: changelog
date: 2026-09-13
---
# Phase 4 — docs trail and review fixes     (2026-09-13)

## Summary
Closed the trail for 830f1fb..b64ffe8 (catch-up entry, bug records, ROADMAP/TODO); updated all repo docs and deployment guides to stop naming `/tmp` and instead use `.agents/scratch/` per the hard rule (BUG-007); fixed code issues from round 1 security review: `--dry-run` now never builds or executes (prints the intended action and inspects with the running binary); `install`/`sync`/`update` refuse a checkout found only by walking the working directory when the manifest records a different one, with a message to pass `--repo`; responded to two govulncheck advisories (stdlib fixes in Go ≥1.24.6, see TODO); added `TestDryRunNeverBuilds` and `TestKnownCheckout` to the installer tests. Second review round over the fixes is pending (see `reviews/2026-09-13-*-learning-loop-round-2.md`).

## Type
task

## Files changed
- `cmd/sh-agent/main.go` — dry-run check, known-checkout validation
- `cmd/sh-agent/main_test.go` — TestDryRunNeverBuilds, TestKnownCheckout (added)
- `internal/harness/harness_test.go` — TestReinstallBuildsAndInstalls, TestInstallMarksDirtyTree (added)
- `.agents/docs/architecture/DECISIONS.md` — ADR-013 (rebuild not copy), ADR-014 (all skills model-invocable); marked ADR-002, ADR-009 superseded
- `.agents/docs/bugs/INDEX.md` — rows for BUG-006, BUG-007, BUG-008, BUG-009
- `.agents/docs/bugs/BUG-007-repo-docs-tmp-scratch-contradiction.md` — new
- `.agents/docs/bugs/BUG-008-retro-commits-skipped-repo-trail.md` — new
- `AGENTS.md`, `README.md`, `.agents/docs/deploy/local.md` — removed `--home /tmp/...` examples, documented `binary` line and `retros.jsonl`
- `ROADMAP.md`, `TODO.md` — updated Now/Next/Later, added `sh-agent update` and Phase 3 gate items
- `testing/2026-09-13-verify-research-pass.md`, `testing/2026-09-13-verify-learning-loop-phases-1-2.md` — new
- `reviews/2026-09-13-review-learning-loop-phases-1-2.md`, `reviews/2026-09-13-security-learning-loop-phases-1-2.md` — new
- `plans/2026-09-13-learning-loop-project-level.md` — Phase 4 marked complete

## Agents dispatched
verifier (phases 1–2 smoke tests), reviewer (round 1, code), security (round 1, code), bugs ×2 (BUG-007, BUG-008), changelog ×3 (catch-up 830f1fb..b64ffe8, ADRs+TODO, this entry), documentation ×1 (ROADMAP by docs agent), reviewer + security (round 2 pending over fixes).

## Verification
Unit: `gofmt -l .` (nothing), `go vet ./...` (clean), `go test -race -count=1 ./...` (`ok cmd/sh-agent 1.128s`, `ok internal/harness 1.461s`). Second review round running as this entry is written; results in `reviews/2026-09-13-*-learning-loop-round-2.md` files.

## Requirements
R4 (Records stay in the trail) — all changes, bugs, and ADRs appended to `DECISIONS.md`, `bugs/INDEX.md`, `changes/`. R6, R7 (harness state, retro queue in `~/.config/sh-agent/` only) — no project-level retro mechanism until Phase 3; addressed in ADR-013, ADR-014.
