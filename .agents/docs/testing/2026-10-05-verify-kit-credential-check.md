---
owner: verifier
date: 2026-10-05
---
# Verification — kit-credential-check   (2026-10-05)

Range reviewed: 218e301..0eb8320

## Status: FAIL

Scope: steps 1-4 as dispatched ran; full -race of cmd/sh-agent not owed (no non-test Go source changed). Guide exists (.agents/docs/testing/guide.md, changed in range, wording only). Not complete: one failure.

## Commands
- 1. `gofmt -l cmd internal` → exit 0, printed nothing
- 2. `go vet ./...` → exit 0
- 4. `go test -count=1 -timeout 30m -v ./internal/...` (output .agents/scratch/verify-kit/go-test-1.txt) → exit 1; counts: pass 449, fail 1, skip 3 at 0eb8320 for go test -count=1 -timeout 30m -v ./internal/...; vs previous: none with this command (v2.4.1 log used other commands)
  - ok internal/dream, internal/memory, internal/search; FAIL internal/harness
  - FAIL: TestScriptsHaveNoEarlyExitPipe: `catalog_test.go:557: kit-preamble.sh pipes into a consumer that exits early, without a guard (BUG-155): unreadable=$(find -H "$KIT_HOME" -path "$KIT_HOME/bin" -prune -o -print 2>&1 >/dev/null | head -5)`
  - skips (pre-existing, no t.Skip added in range): TestAgyCorpusLeaves, TestRecall, TestHelperIndex
- 4. `go test -race -count=1 -run 'TestKitPreamble' -v ./internal/harness/` → exit 0; counts: pass 14, fail 0, skip 0 at 0eb8320
- 4. `bash -n .agents/docs/testing/kit-preamble.sh` → exit 0

## Root cause
The rewritten credential check in kit-preamble.sh (new line, range diff) runs `find ... 2>&1 >/dev/null | head -5` in a script that uses pipefail. That is exactly the shape BUG-155's guard test (TestScriptsHaveNoEarlyExitPipe, internal/harness/catalog_test.go:557) refuses: head can exit early, find gets SIGPIPE, and under pipefail the pipeline status becomes non-zero. The new TestKitPreambleRefusesLeftoverCredentials (12 subcases) passes, but the full package fails on the pre-existing guard. The author ran only TestKitPreamble* and so missed it. Not pre-existing: the line is added in this range. Fix is the author's (guard the pipe or avoid head); not patched here.

## Coverage of changed code
- .agents/docs/testing/kit-preamble.sh — covered by internal/harness/catalog_test.go TestKitPreambleRefusesLeftoverCredentials (pass) and TestKitPreamble* (14 pass); also trips TestScriptsHaveNoEarlyExitPipe (FAIL)
- internal/harness/catalog_test.go — is the test itself
- .agents/docs/testing/guide.md, changes entry — docs, no test
- Suppressions added in range: none seen in diff stat; full diff not grep-audited for nolint/skip beyond the 3 existing skips being unchanged.

## Not examined
cmd/sh-agent tests (not owed); full -race of ./internal/...; govulncheck; Docker; the changes entry content; whether other tests besides the one failing would fail under shell edge cases.

Step 8: tree check at end: untracked `.agents/docs/reviews/2026-10-05-review-kit-credential-failclosed.md` appeared (another agent's artifact, listed as seen); no source, test or config path changed.

## Environment
Linux 6.18.28, go1.27.1 linux/amd64, no Docker used.
