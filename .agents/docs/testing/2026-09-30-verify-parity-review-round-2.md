---
owner: verifier
date: 2026-09-30
---
# Verification — parity-review-round-2   (2026-09-30)

Range reviewed: e2cffc1..eeab1b7

## Status: PASS

Scope: guide § 1, § 3 (scratch install + drift), cover of internal/dream and internal/search, one must-fail mutation, two cross-builds. Not run: § 2, § 3a-3d, § 4, § 5 (see Not examined). Guide cited: .agents/docs/testing/guide.md (exists).

## Commands
- `gofmt -l cmd internal` → exit 0, printed nothing
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0 (.agents/scratch/verify-r2/go-test-1.txt); counts: pass 461, fail 0, skip 2 at eeab1b7 for go test -race -v -count=1 ./...; vs previous: .agents/docs/testing/2026-09-30-verify-parity-review-round-1.md, pass 453→461, skip 2→2 (TestRecall, TestHelperIndex, env-gated)
  All 6 packages ok (root, cmd/sh-agent, internal/dream, internal/harness, internal/memory, internal/search).
- Tests added by the range, all PASS: TestCollectNamesAnUnreadableTranscript, TestCollectDoesNotFollowALinkedAgentsDir, TestRunUsage, TestRunAgyReadErrorKeepsWatermark, TestRunStateWriteFailureIsExit2, TestRemoveCopiesKeepsFreshOnes, TestAnswerNeedsASuccessfulResult, TestInstallTimerKeepsOnlyAbsolutePath; plus extended TestAgyConversationText (PASS).
- `go test -count=1 -cover ./internal/dream ./internal/search` → ok; dream 89.2%, search 84.3%.
- Scratch install: `go run ./cmd/sh-agent install --home .agents/scratch/verify-r2-home` → exit 0; `status` → "drift none — 267 files match the repo".
- `GOOS=darwin GOARCH=arm64 go build ./...` → exit 0
- `GOOS=linux GOARCH=arm64 go build ./...` → exit 0
- Mutation (must-fail): internal/search/dream.go:240 removed ` || strings.HasPrefix(s, `{"`)` (the source uses a backtick literal, not `"{\""`; same clause). Mutated code built (exit 0); `go test -run TestAgyConversationText ./internal/search` → FAIL: "dream_test.go:170: text ... carries a tool call's arguments (security review round 2)". Restored from a copy; sha256 a7d927f4...e8d5c identical before/after; `git status` clean apart from this log.

## Coverage of changed code
- internal/dream/collect.go — collect_test.go (two new tests)
- internal/dream/dream.go, run.go — run_test.go (new tests)
- internal/dream/timer.go — timer_test.go
- internal/search/dream.go — dream_test.go (mutation-proven)
- Suppressions added: none checked beyond grep of added test names; no t.Skip in the added test list.

## Not examined
§ 2 harness self-check, § 3a-3d, § 4 and § 5 sandbox/real-home CLI starts (no claude/agy started, as instructed); macOS runtime behaviour of syscall.Flock (cross-compile only); added-suppression grep over the diff not run separately; no docs diff read.

## Environment
Linux 6.12 (Debian 13) amd64; Go as installed on PATH; no Docker used.
