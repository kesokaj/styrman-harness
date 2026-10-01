---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-3-r2 (2026-10-01)

Range reviewed: 329a99b..working tree (dirty); guide: .agents/docs/testing/guide.md; previous: testing/2026-10-01-verify-audit-fix-phase-3.md

## Status: PASS

Scope: guide §1, §2 (except CLI start), §3, bash -n, modes. Not run: §4 CLI start and behavioural run (recorded by the main agent in testing/2026-10-01-verify-size-tier-behavioural.md).

## Commands
- `gofmt -l cmd internal` -> exit 0, no output
- `go vet ./...` -> exit 0
- `go test -race -v -count=1 ./...` -> exit 0, all 6 packages ok; counts: pass 485, fail 0, skip 3 at 329a99b+dirty for go test -race -v -count=1 ./...; vs previous: testing/2026-10-01-verify-audit-fix-phase-3.md, pass 484->485, skip 3->3 (skips: TestAgyCorpusLeaves, TestRecall, TestHelperIndex). Output .agents/scratch/verify-audit-fix-phase-3-r2/go-test-1.txt
- `go run ./cmd/sh-agent list` -> exit 0, 50 lines
- `go run ./cmd/sh-agent render | wc -c` -> exit 0, 14882 bytes (under 15,000 target and 22,000 limit)
- TestRenderedContextFitsAntigravityBlock, TestRenderedContextUnderTarget -> PASS
- `go test ./cmd/sh-agent -run TestHarnessMap` -> ok
- `go test . -run TestEmbedded -v` -> TestEmbeddedPlanEqualsCheckout PASS (BUG-137 fixed)
- `install --home .agents/scratch/verify-p3r2` -> exit 0; `status` -> "drift none - 269 files match the repo"; scratch home removed (confirmed gone)
- `bash -n` kit-preamble.sh, restatement-grep.sh -> exit 0 both
- `ls -l`: kit-preamble.sh -rw-r--r--, restatement-grep.sh -rw-r--r-- (both 0644)

## Coverage of changed code
- Only Go change: internal/harness/size_test.go TestRenderedContextUnderTarget, read: const target 15000, renders context from checkout and fails above it; passes at 14882. Suppressions added: none seen in that file.
- Other changes are markdown/shell, verified by install/status and bash -n.

## Not examined
§4 CLI start, behavioural run, real-home smoke checks, skill matrix, govulncheck, search §3d (no internal/search change).

## Environment
Linux 6.12, bash; tree status start vs end differs only by this log.
