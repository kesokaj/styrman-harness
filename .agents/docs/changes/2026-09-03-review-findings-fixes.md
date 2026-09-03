---
owner: changelog
date: 2026-09-03
scope: changeset:af2dcaa6c4d3
---

# Gate, scope, and verifier prompt fixes from review findings

## Summary

Fixed eight high- and medium-confidence findings from reviewer reports of prior changesets.

**Scope computation & gates:** Trail gates (`trailgate.go`, `subagent.go`) now load `.agents/harness.json` before computing scope, fixing mismatched digests between `harness check` and `harness hook stop`. Added `HeadScope()` for commit-boundary logic, isolating it from worktree-aware `Scope()` used by owed-rows logic. Recompute scope after agents finish to detect tree mutations between dispatch start and artifact write.

**Head-based trail measurement:** Fixed merge-commit handling (`git.go`: first-parent only) and proportionality logic (`owed.go`: count untracked files correctly). Verifier tests now commit before asserting to catch date-based false positives.

**Verifier prompt:** `rowPrompt` now includes a `RowVerifier` case with captured command output and failure logs, matching the agent instruction to "read the dispatch prompt for the exact command(s) that failed".

**Glob compilation:** Fixed UTF-8 handling (`source.go`: iterate runes, not bytes) so non-ASCII source globs like `src/år/*.go` match correctly. Addresses silent disabling of Definition-of-Done rows.

**Ancillary fixes:** Dead code removal (`check.go:70` unreachable case). Agy plugin install predicate now consistent with mirror skip. Scope added to `harness dispatch --dry-run` output and stop gate notices.

## Type
bugfix

## Files changed
32 files, 1002 lines total. Primary changes in:
- `internal/trail/{scope.go, head.go, git.go, owed.go, check.go}` — scope/HEAD/merge handling, owed logic
- `internal/gate/{trailgate.go, subagent.go}` — config loading, scope consistency
- `internal/dispatch/dispatch.go` — verifier prompt with command output, scope recomputation
- `internal/trail/source.go` — UTF-8 glob compilation
- `internal/install/install.go` — agy predicate consistency
- `claude/agents/verifier.md`, `antigravity/plugins/harness/agents/verifier.md` — instruction updates
- Test files: `trail_test.go`, `gate_test.go` — regression tests for merge commits, scope isolation, glob UTF-8

## Verification
`gofmt -l . | tee /dev/stderr | (! read)` ✓
`go vet ./...` ✓
`go test -race ./...` ✓ — all packages pass

Per `.agents/docs/testing/2026-09-03-verify-local.md`.

## Requirements
Addressed findings in `.agents/docs/reviews/2026-09-03-review-scope-coverage-gates.md` (high & medium findings 1–3, 7–8), `review-agy-plugin-and-head-trail-fixes.md` (finding 1), `review-glob-matcher-and-doc-sync.md` (findings 1–2). Maintains R4 (gates inform, block only at boundary), R5 (proportional rows), R7 (scope-based identification).
