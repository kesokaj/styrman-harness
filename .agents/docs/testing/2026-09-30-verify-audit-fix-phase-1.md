---
owner: verifier
date: 2026-09-30
---
# Verification — audit-fix-phase-1   (2026-09-30)

Range reviewed: 4764c5f..working tree (uncommitted, 16 files, 839+/94-). Guide: .agents/docs/testing/guide.md

## Status: PASS

Scope: guide §1, §3, §3b, §3d, §3e and the corpus measurement ran. Hunk-by-hunk coverage was not mapped; coverage is judged per named test (below). Real CLIs and real home untouched.

## Commands
- `gofmt -l cmd internal` → exit 0, no output
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0; counts: pass 478, fail 0, skip 3 at 4764c5f+dirty for go test -race -v -count=1 ./... (log .agents/scratch/verify-audit-fix-phase-1/go-test-1.txt); vs previous: .agents/docs/testing/2026-09-30-verify-v1.17.0-gate-r2.md, pass 462→478, skip 2→3; explained: the new env-gated TestAgyCorpusLeaves (t.Skip with reason), the other two skips unchanged (TestRecall, TestHelperIndex)
- §3 `go run ./cmd/sh-agent install --home .agents/scratch/verify-p1` → exit 0; `status` → "drift none — 267 files match the repo"
- §3b `go test ./internal/harness -run 'TestReportRedact' -count=1 -v` → 6 PASS (JSONL, Shapes, StructuralGaps, CleanIsByteIdentical, Text, Document); `report --dry-run --home .agents/scratch/verify-p1` → exit 0, printed destination and would-write lines, absent files skipped, nothing pushed
- §3d `sh-agent index --rebuild --home .agents/scratch/s` → exit 0, 2.07 s (limit 10 s); warm `search verifier` → exit 0, 0.056 s (limit 100 ms); `go test ./internal/search -count=1` → ok
- §3e `go test -race -count=1 ./internal/memory ./internal/dream ./internal/search ./cmd/sh-agent` → ok all four (also `./internal/search ./cmd/sh-agent` alone ok)
- corpus: `SH_AGENT_CORPUS=$PWD/.agents/scratch/corpus go test ./internal/search -run TestAgyCorpusLeaves -count=1 -v` → PASS; `corpus: 95 databases, 971 prose leaves kept, 0 with a tool-argument key; samples []` (matches D30: 971 and 0)

## Coverage of changed code
Each named test exists and PASSED in the -v run:
- internal/search/dream_test.go — TestProseLeaves, TestCopyDBRetriesOnce, TestCopyDBNeverOpensSource, TestAgyCorpusLeaves (measurement)
- internal/search/search_test.go — TestNoDuplicateHitAcrossRounds
- internal/dream/run_test.go — TestTextFieldCannotForgeASection, TestAnswerToleratesExtraKeys, TestPromptEscapesFence, TestAnswerLineTooLong, TestRemoveCopiesUsesInjectedClock
- internal/dream/collect_test.go — TestSessionTextSkipsCommandOutput, TestTrailDateFromGit
- internal/harness/report_test.go — TestReportRedactShapes; harness_test.go — TestInstallMidFailureKeepsManifestComplete, TestInstallDoesNotFollowLinkedDest, TestResyncSkipsLink
- cmd/sh-agent/memory_test.go — TestDreamVerb; main_test.go — TestUpdateNeverInstallsDirtyTree (8 PASS lines incl. subtests)
- Every changed source file has a changed test file in the range. Not verified: that the named tests cover every changed hunk, nor that they fail without the fix (no mutation run).
- Suppressions added: two t.Skip, both with a reason ("root ignores directory modes"; "measurement, not a gate: set SH_AGENT_CORPUS=..."). No .only, nolint or ignore directives.

## Not examined
Mutation check of the new tests; per-hunk coverage; `go test -cover`; §3c matrix rows, §4 and §5 (no context/agents/skills change, no real CLI by instruction); govulncheck; the TestRecall measurement (SH_AGENT_RECALL not set).

## Environment
Linux 6.12 (Debian 13), go from PATH (version not recorded), no Docker used. Tree at end: only change beyond the 16 modified files is this log (mine); no source/test/config path changed during the run.
