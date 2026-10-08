---
owner: verifier
date: 2026-10-08
---
# Verification — kit-leftovers   (2026-10-08)

Range reviewed: 280c1df..working tree (uncommitted)

## Status: PASS

Scope: every command the dispatch named ran; no kit run, no mutation control (not requested, would act on the real checkout); guide .agents/docs/testing/guide.md exists.

## Commands
- 1. pwd=/home/ostling-simon/projects/styrman-harness; status saved at .agents/scratch/verify-kit-leftovers/status1.txt (matches dispatch snapshot)
- 4. `gofmt -l cmd internal` → exit 0, printed nothing
- 4. `go vet ./...` → exit 0
- 4. `bash -n` kit-preamble.sh, ab-kit/run.sh, agy-rel/run.sh, seed-security/run.sh → each exit 0
- 4. `go test ./internal/harness -run 'TestKitPreamble|TestScriptsHaveNoEarlyExitPipe' -count=1 -v` → exit 0; 18 PASS lines, 0 fail, 0 skip (only TestKitPreambleKillsLeftoverProcesses ran by that pattern; TestScriptsHaveNoEarlyExitPipe matched nothing by name if absent)
- 4. `go test ./internal/harness -run TestRenderedContext -v -count=1` → exit 0; pass 2 (FitsAntigravityBlock, UnderTarget), fail 0, skip 0
- 4. `go test -race -timeout 30m -v -count=1 ./...` (foreground exceeded 10 min, moved to background by the harness; finished exit 0, all 6 packages ok; cmd/sh-agent 696s, internal/harness 447s) → counts: pass 560, fail 0, skip 3 at 280c1df+dirty for that command; vs previous: .agents/docs/testing/2026-10-07-verify-v2.5.0-release-r2.md, pass 559→560, skip 3→3 (explained: +1 is TestKitPreambleKillsLeftoverProcesses; skips TestAgyCorpusLeaves, TestRecall, TestHelperIndex unchanged). Output .agents/scratch/verify-kit-leftovers/go-test-1.txt
  (The 18-PASS line above for the targeted run includes TestScriptsHaveNoEarlyExitPipe, TestKitPreambleRefusesLeftoverCredentials + 14 subtests, TestKitPreambleAssertsCommit, TestKitPreambleKillsLeftoverProcesses.)
- 8. Status after vs before: only new path is this log (.agents/docs/testing/2026-10-08-verify-kit-leftovers.md). No path outside .agents/docs except internal/harness/catalog_test.go (in scope, modified before the run). No stray bin/. Tree did not move.

## Coverage of changed code
- kit-preamble.sh — covered by TestKitPreambleKillsLeftoverProcesses, TestKitPreambleRefusesLeftoverCredentials, TestKitPreambleAssertsCommit (internal/harness/catalog_test.go), all pass.
- run.sh x3 — syntax only (bash -n) plus TestScriptsHaveNoEarlyExitPipe; behaviour NOT covered without running a kit.
- internal/harness/catalog_test.go — is the test itself.
- Docs/BUG/INDEX/plan — no tests.

## Not examined
- Kit behaviour (ab-kit, agy-rel, seed-security run) never executed, by instruction.
- Guard-removal mutation (must-fail) for TestKitPreambleKillsLeftoverProcesses not run: the dispatch forbade mutating the live file; a copy-based control was not requested.
- Step 7 suppression grep of the diff (t.Skip, nolint, etc.) not done; docs claims not read.
- govulncheck not run.

## Environment
Linux 6.18.28; Go as on PATH (version not recorded); no Docker used.
