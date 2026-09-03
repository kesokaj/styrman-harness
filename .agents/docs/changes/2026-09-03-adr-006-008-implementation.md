---
owner: changelog
date: 2026-09-03
scope: worktree:9366002f2d81
---

# Implement ADRs 006–008: formatters, scope-based satisfaction, dispatcher testing

## Summary
Implemented three architectural decisions:

**ADR-006** (Formatters): Resolved formatters through PATH only; refuse any binary that
resolves inside the working repository or `node_modules`. Removes the 419 ms per-edit tax
on PostToolUse hooks and prevents hostile repositories from auto-executing local tools.

**ADR-007** (Scope-based satisfaction): Changed row satisfaction from date-based to
scope-based frontmatter. Every artifact now declares `scope:` (`worktree:<digest>` for
uncommitted work, `commit:<sha>` for clean trees). Rows fire for the current scope only,
fixing the bug where a second commit on the same day owed no rows.

**ADR-008** (Dispatcher testing): Added a stub CLI that mimics `claude -p` and `agy -p`
for testing the dispatcher contract in CI without credentials or spend. Caught a timeout
bug (`exec.CommandContext` cancellation did not bound wall time; fixed with `WaitDelay`).

## Type
task

## Files changed
25 source files (789 lines). Primary changes in:
- `internal/gate/` — formatter resolution (ADR-006), resume context injection
- `internal/trail/` — scope computation and worktree digest (ADR-007), owed-rows logic
- `internal/dispatch/` — CLI stub and testing (ADR-008)
- `internal/install/derive.go` — agent derivation with new capabilities
- Root `.agents/docs/architecture/DECISIONS.md` — ADRs 006–008 appended (already recorded separately)

## Verification
`go test -race ./...` and `go vet ./...` pass; `gofmt` reports no changes needed.
Per `.agents/docs/testing/2026-09-03-verify-local.md`.

## Requirements
ADRs 006–008 settle design decisions affecting gates and the dispatch contract.
No explicit REQUIREMENTS.md item; these are refinements to R3 (drift detection) and R4
(gates inform). See `DECISIONS.md` for rationale.
