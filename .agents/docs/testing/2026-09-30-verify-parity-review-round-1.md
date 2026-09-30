---
owner: verifier
date: 2026-09-30
---
# Verification — parity-review-round-1   (2026-09-30)

Range reviewed: f1ca525..e2cffc1

## Status: PASS

Scope: complete for the tiers named in the dispatch (§ 1, § 3, § 3b, cover, two mutations); guide exists; § 2, § 4, § 5 not requested.

## Commands
- `gofmt -l cmd internal` → exit 0, printed nothing
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0 (.agents/scratch/verify-pr1/go-test-1.txt); counts: pass 453, fail 0, skip 2 at e2cffc1 for go test -race -v -count=1 ./...; vs previous: .agents/docs/testing/2026-09-30-verify-parity-memory-dream-phase-4-5.md, pass 443→453, skip 2→2 (TestRecall, TestHelperIndex, env-gated)
  - Tests added by range, all PASS: TestDreamTimerDryRun, TestReadState, TestCollectDoesNotFollowLinkedTrail, TestRunAgyPrivateHomeAndStdin, TestRunAgyDeniedToolIsNoAnswer, TestRunAgyWithoutLoginFailsBeforeProjects, TestRunPreProjectFailureIsExit2AndRecorded, TestRunCorruptStateIsExit2AndKept, TestRunRemovesLeftoverCopies, TestRenderFoldsLineBreaks, TestInstallTimerRefusesUnsafePaths, TestLeftoverLockFileIsNoLock, TestHeldLockRefusesAfterRetry, TestIndexMatchIsTheLinkTarget (14 new funcs; TestAgyConversationText and others modified)
- `go test -cover -count=1 ./internal/dream ./internal/memory ./internal/search` → exit 0; dream 87.2%, memory 89.1%, search 84.2%
- § 3 `sh-agent install --home .agents/scratch/verify-pr1/home` → exit 0 (268 written); `status` → "drift none — 267 files match the repo"
- § 3b `sh-agent report --dry-run --home <scratch>` (scratch home holding proposals/dream/harness/x.md, proposals/dream/proj/p.md, proposals/dream/projdream.md) → exit 0; would write only state.json, proposals/dream/harness/x.md, report.json; neither project dream file listed (output .agents/scratch/verify-pr1/report.txt)
- Mutation (a): removed ` WHERE step_type IN (14, 15)` in internal/search/dream.go (go vet ok, compiled) → `go test ./internal/search -run TestAgyConversationText -v` exit 1; `--- FAIL: TestAgyConversationText`, "carries a tool's output (security review 2026-09-30)". Restored from copy; sha256 ae22506c… identical to before.
- Mutation (b): internal/dream/run.go agyHome settings replaced with `{"toolPermission":"always-proceed"}` (go vet ok, compiled) → `go test ./internal/dream -run TestRunAgyPrivateHomeAndStdin -v` exit 1; `--- FAIL`, "child settings grant something: {"toolPermission":"always-proceed"}". Restored; sha256 9bd54d1c… identical to before.

## Coverage of changed code
- internal/dream/{collect,dream,run,state,timer}.go — covered by collect_test.go, run_test.go, timer_test.go (new tests above); package 87.2%
- internal/memory/memory.go — memory_test.go (+ cmd/sh-agent/memory_test.go); 89.1%
- internal/search/dream.go — dream_test.go; 84.2%
- cmd/sh-agent/main.go — cmd/sh-agent/memory_test.go
- skills/*, docs: no Go code
- Suppressions added (nolint, Skip, ts-ignore, .only): none in the diff.

## Not examined
§ 2 harness self-check, § 4 sandbox CLI starts, § 5 real-home (not requested; no claude/agy started). Line-by-line coverage of the new branches beyond package percentages. Scratch home under .agents/scratch/verify-pr1 left in place (gitignored).

## Environment
Linux 6.12 amd64, Go from PATH; no Docker used. Tree clean at start; at end only this log is new (source/test/config unchanged, sha-verified for the two mutated files).
