---
owner: changelog
date: 2026-09-03
scope: changeset:06cc436fce2d
---

# Harness security fixes, proportionality audit, and ADR implementation

## Summary

Addressed high-severity security and logic findings from code and security reviews of the initial Go consolidation. Security review identified 15 bypasses to the git credential gate and two arbitrary-code-execution paths (repo-controlled `verify[]` commands and hostile repo formatters). Code review flagged a proportionality bug where untracked files (3965 lines of new code in the harness itself) scored zero lines for threshold purposes.

Fixes implemented:
- **ADR-005** (credential gate): Changed from argument parsing to asking git directly (`git diff --cached --name-only` and worktree scan), closing all 15 identified bypasses
- **ADR-006** (formatters): Formatters resolve through PATH only; refuse any binary inside the working tree or `node_modules`, eliminating hostile-repo code execution on every edit
- **ADR-007** (scope-based satisfaction): Switched row satisfaction from date-based (coarse; second commit on same day owed no rows) to scope-based digests in frontmatter, fixing a bug that hid owed rows five times
- **ADR-008** (dispatcher testing): Added stub CLI for dispatcher contract testing in CI, catching a `WaitDelay` timeout bug in the actual implementation

## Type
bugfix

## Files changed
32 files, 896 lines total. Primary changes:
- `internal/gate/secret.go` — credential gate now asks git, not bash; scans worktree and index
- `internal/gate/marker.go`, `internal/gate/marker_unix.go` — formatter resolution via PATH
- `internal/trail/scope.go`, `internal/trail/scope_test.go` — scope digests and worktree computation (ADR-007)
- `internal/dispatch/cli.go`, `internal/dispatch/integration_test.go` — stub CLI for contract testing
- `internal/trail/owed.go` — proportionality logic now counts untracked files
- Root `.agents/docs/architecture/DECISIONS.md` — ADRs 005–008 appended

## Verification
`gofmt -l . | tee /dev/stderr | (! read)` ✓
`go vet ./...` ✓
`go test -race ./...` ✓ (all packages pass, including integration tests)

Per `.agents/docs/testing/2026-09-03-verify-local.md`.

## Requirements
Addressed R4 (gates inform but must not be wrong), R5 (proportional rows), R7 (artifacts identified by scope, not date). Security findings documented in `.agents/docs/reviews/2026-09-03-security-harness-go-rewrite.md` and code findings in `.agents/docs/reviews/2026-09-03-review-harness-go-binary.md`.
