---
owner: reviewer
date: 2026-10-01
---
# Review — 4764c5f-worktree-code-r3   (2026-10-01)

Range reviewed: working tree against commit 4764c5f (`git diff 4764c5f`; 16 files under cmd/sh-agent, internal/dream, internal/harness, internal/search). Round 3; round 2 is `2026-09-30-review-4764c5f-worktree-code-r2.md`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

All four round-2 answers hold. Nothing is broken in the code. One restatement is still stale: `.agents/docs/deploy/remote.md:75` says `update --check` "works on a dirty checkout", and the diff now makes it fail on one (IN4). It was carried from round 1 as a docs-pass item. I also ran `gofmt -l cmd internal` (clean), `go vet ./...` (clean) and `go test -race ./...` (all six packages ok) on the working tree.

## Findings
- **[low / confidence: high]** `.agents/docs/deploy/remote.md:75` — "sh-agent update --check ... works on a dirty checkout" contradicts `cmd/sh-agent/main.go` `update()`, where the dirty check now sits above the `--check` return and `--check` returns an error. A description that outlives its behaviour reads as a guarantee. Fix in the documentation pass (step 3 of Definition of Done), not in this diff; the docs pass must also check any timer advice near remote.md:92 against it. This is round 1's finding 1 and round 2's "answer 1", still open.

## Round 2 answers
1. Resync-removal error now returns through `saveOnError` (`apply.go:147`); `TestInstallResyncFailureKeepsManifest` exists. Fixed.
2. `Status` (`apply.go:222`) and `Diff` (`:256`) use the same `symlink || underLink` predicate as Install. Fixed.
3. `underLink` (`:407`) now sits above `pruneEmptyDirs`' doc-comment (`:421`), each on its own func. Fixed.
4. Linked `~/.claude` keeping every file is intended: Install prints `keep` lines, and the summary line prints the `kept` count (`apply.go:186`). Accepted; it is not silent and is pinned by `TestUninstallAndLinkedDirs`.

## Criteria
1. Dead code — pass: `copyDB` is still used by `openCopy`; no orphaned helper.
2. Flow tracing — pass: traced `saveOnError` (every error return after a write), `underLink` (Install, resync, Status, Diff, Uninstall), `update()` ordering, `answer()` (sole caller `realRun`), `removeCopies(conf, o.Now())` (`Now` defaulted at `dream.go:56`), `openCopy` callers (three), and `find()` chunk-id dedupe (`h.id` set at the single place a hit is built, before append). Remaining untraced edge (not a finding): `trailDates` uses repo-root-relative paths from git against root-relative paths, so a project that is a repo subdirectory falls back to mtime, the old behaviour.
3. Standards conformance — minor fail: remote.md:75 stale (above); gofmt, vet and tests are green.
4. Directory layout — pass: no new source files; tests beside code.
5. Simplicity (DIY-except-security) — pass: stdlib only; the redaction `accept` func is a filter, not auth or crypto.
6. Decoupling — pass: `copyFile` and `copyRetryWait` are in-package seams; no cross-package reach-ins.

## Not examined
- Test bodies in run_test.go, dream_test.go, collect_test.go, search_test.go, main_test.go, memory_test.go, and most of harness_test.go beyond the named tests (existence and green status only; no must-fail mutations run by me).
- `trailDates` with a staged rename (`R  old -> new` in porcelain), quoted or non-ASCII paths, and git-call timeouts.
- Behavioural, sandbox and real-home runs; govulncheck; the Antigravity run (AGENTS.md § Testing).

## Outcome

1. `remote.md:75` and the timer advice near `:92` — fixed by the documentation pass of this phase (deploy guide, README update row, delivery spec), same commit.
Round 1 (`2026-09-30-review-4764c5f-worktree-code.md`): 1 fixed (docs pass), 2 fixed (per-key value rule, `TestReportRedactShapes`), 3 rebutted (corpus measurement: 971 leaves, 0 argument keys, before and after), 4 fixed (`saveOnError` skips links, `underLink`). Round 2 (`-r2.md`): 1 fixed (resync error path, `TestInstallResyncFailureKeepsManifest`), 2 fixed (`Status`/`Diff` use `underLink`), 3 fixed (comment placement), 4 intended behaviour, pinned by `TestUninstallAndLinkedDirs` (keep lines and a `kept` count, not silent). No fourth round: code converged on APPROVE WITH NITS with the one docs item.
