---
owner: changelog
date: 2026-09-03
scope: changeset:f6f6c28eee97
---

# Implement ADR-009: narrative row exception for changelog     (2026-09-03)

## Summary
Implemented the narrative row exception: the changelog row is now satisfied by any uncommitted entry from today while work is in progress. At a commit, the tree is clean and the exception lapses — the entry must then cover HEAD's scope exactly, restoring the original bug-fix property.

Changelog entries record tasks (what was done and why), not code states. A byte-exact scope scope made the row unsatisfiable on a working tree: the entry was stale as soon as it was written, and the changelog agent's anti-duplication rule declined to write another, leaving the row due despite valid entries on disk.

## Type
task

## Files changed
- `internal/trail/trail.go` — added `NarrativeRow()` method on Row; added `hasInProgressEntry()` method to check for uncommitted entries from today
- `internal/trail/owed.go` — added narrative row exception logic to `Owed()` (lines 109–125); changelog now excused from scope match while tree is dirty
- `internal/trail/scope_test.go` — added `TestNarrativeExceptionRequiresCurrentTaskEntry` (regression: verifies committed entries don't satisfy new work, superseded scope forms don't satisfy anything, uncommitted current-form entries do)
- `internal/trail/scope_test.go` — added `TestFreshnessDoesNotInheritTheException` (verifies exception is dispatcher-only, freshness still requires exact coverage)
- `REQUIREMENTS.md` (R7) — documented narrative vs. code-state row distinction; clarified that changelog is scoped to task, not code state; exception lapses at commit boundary

## Verification
`go test -race ./...` and `go vet ./...` pass (verified in `.agents/docs/testing/2026-09-03-verify-local.md`). New tests enforce the property: an uncommitted entry satisfies the row while work is dirty, but the entry is stale to freshness checking and the exception drops at commit boundary.

## Requirements
ADR-009 (Accepted 2026-09-03). Solves the deferred work from 2026-09-03-stop-gate-scope-unification entry: re-evaluates scope scheme to distinguish narrative rows from code-state rows. Satisfies R7 (artifacts identified by scope, not date) and R5 (proportional DoD).
