---
owner: verifier
date: 2026-10-05
---
# Verification — kit-credential-check-r2   (2026-10-05)

Range reviewed: 218e301..29cd652

## Status: PASS

Previous run: FAIL at 0eb8320 (TestScriptsHaveNoEarlyExitPipe), log .agents/docs/testing/2026-10-05-verify-kit-credential-check.md (left unchanged: a FAIL log is never updated in place, so this run has its own file).

Scope: steps 1-5 as dispatched all ran at 29cd652; no non-test Go source changed so full cmd/sh-agent -race not owed. Guide exists (testing/guide.md, wording change only). Not `complete`: cmd/sh-agent and govulncheck not run.

## Commands
- 1. `gofmt -l cmd internal` → exit 0, printed nothing
- 2. `go vet ./...` → exit 0
- 3. `go test -count=1 -timeout 30m -v ./internal/...` (-v added for counts; output .agents/scratch/verify-kit/go-test-3.txt) -> exit 0; counts: pass 452, fail 0, skip 3 at 29cd652 for go test -count=1 -timeout 30m -v ./internal/...; vs previous: .agents/docs/testing/2026-10-05-verify-kit-credential-check.md, pass 449->452, skip 3->3 (explained: pass +3 from tests added/split in range, and the earlier run counted the failure). ok dream, harness (118s), memory, search. Skips unchanged: TestAgyCorpusLeaves, TestRecall, TestHelperIndex.
- 4. `go test -race -count=1 -run 'TestKitPreamble|TestScriptsHaveNoEarlyExitPipe' -v ./internal/harness/` -> exit 0; counts: pass 17, fail 0, skip 0 at 29cd652 for that command; vs previous: none with this exact command (0eb8320 ran TestKitPreamble only: pass 14). `--- PASS: TestScriptsHaveNoEarlyExitPipe`.
- 5. `bash -n .agents/docs/testing/kit-preamble.sh` -> exit 0

What failed at 0eb8320 (TestScriptsHaveNoEarlyExitPipe): PASSES now, in both the isolated run (step 4) and the whole-package run (step 3).

## Coverage of changed code
- .agents/docs/testing/kit-preamble.sh -- covered by internal/harness/catalog_test.go TestKitPreamble* and TestScriptsHaveNoEarlyExitPipe (pass)
- internal/harness/catalog_test.go -- is the test
- guide.md, changes entry, reviews, prior log -- docs, no test
- Suppressions added in range (diff grep): catalog_test.go added `t.Skip("root reads through permissions")` and `t.Skip("no symlinks here")` (diff lines ~243, ~281), each with a reason on its line; neither triggered (skip count unchanged). No nolint, .only or ignore.

## Not examined
cmd/sh-agent tests; full -race over ./internal/...; govulncheck; Docker; content of changes/review entries; the mutation (must-fail) check of the guard, not requested.

Step 8: tree vs start (clean): new untracked `.agents/docs/reviews/2026-10-05-review-kit-credential-failclosed-r2.md` (another agent's artifact, seen) and this log. No source, test or config path changed.

## Environment
Linux 6.18.28, go1.27.1 linux/amd64, no Docker used.
