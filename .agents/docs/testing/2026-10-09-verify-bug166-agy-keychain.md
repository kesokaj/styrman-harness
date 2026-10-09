---
owner: verifier
date: 2026-10-09
---
# Verification — bug166-agy-keychain   (2026-10-09)

Range reviewed: 9964037 + uncommitted (7 modified, 3 untracked; run.go, run_test.go, docs)

## Status: PASS

Scope: guide § 1 full suite plus the four named new tests ran on Linux; the dispatch's second command (memory, dream, search, cmd/sh-agent) was not run separately: those packages ran inside the full `./...` run, all ok. Darwin keychain branch NOT exercised for real (see Not examined), so not `complete`. guide.md exists, unchanged by me.

## Commands
- 1. pwd /home/ostling-simon/projects/styrman-harness; git status as listed in dispatch (10 paths); HEAD 9964037
- 7. suppressions added in diff (nolint/Skip/.only): none found by grep
- 4. `gofmt -l cmd internal` -> exit 0, printed nothing
- 4. `go vet ./...` -> exit 0
- 4. `go test -count=1 -v ./internal/dream -run 'TestAgyHomeMacKeychain|TestAgyAuthHint|TestRunAgyAuthFailureNamed|TestRunRemovesLeftoverCopies'` -> exit 0; all four `--- PASS` (TestRunRemovesLeftoverCopies, TestAgyHomeMacKeychain, TestAgyAuthHint, TestRunAgyAuthFailureNamed); ok internal/dream 0.481s
- 4. `go test -race -v -count=1 -timeout 30m ./...` (output .agents/scratch/verify-bug166/go-test-1.txt) -> exit 0; counts: pass 649, fail 0, skip 3 at 9964037+dirty for go test -race -v -count=1 -timeout 30m ./...; vs previous: none (earlier logs ran without -v counts). Skips: TestAgyCorpusLeaves, TestRecall, TestHelperIndex (pre-existing env-gated skips, not in the diff). Packages ok: root, cmd/sh-agent (706s), internal/dream, harness, memory, metrics, search, trail.
- 7. suppressions added in diff: none (grep for nolint/Skip/.only in added lines).

## Coverage of changed code
- internal/dream/run.go -- covered by internal/dream/run_test.go (4 new tests). Gap: darwin branch exercised only with goos="darwin" passed on Linux; no real macOS keychain run (owed by operator before release).
- Docs/README/bug records: no code, no test.

## Not examined
- Real macOS run of `sh-agent dream` with ~/Library/Keychains link (owed by operator).
- Separate rerun of the 4-package command (covered by the full run).
- Tests are not mutation-checked (no must-fail control requested).

## Environment
Linux 6.18, go1.27.2, 16 CPUs, no Docker needed (no Dockerfile change).
Tree check (step 8): after the run, status = step-1 set plus this log only; no source change.
