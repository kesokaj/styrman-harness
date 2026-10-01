---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-9-r3   (2026-10-01)

Range reviewed: 2ba4a50..working tree (uncommitted Phase 9)

## Status: PASS

Scope: complete (guide § 1, size pins, scratch install, coverage map); no behavioural/Antigravity runs requested.

## Commands
- 1. status count at start: 48
- 1. `gofmt -l cmd internal` -> exit 0, no output
- 2. go vet ./... -> exit 0
- 3. `go test -race ./...` -> exit 0; all 6 packages ok (root, cmd/sh-agent, internal/dream, internal/harness, internal/memory, internal/search); dream and memory were (cached) in this run, re-run uncached in step 4 and ok
- 3. named harness tests -v -> exit 0; PASS: TestVendoredRowsCarryAPin, TestAgentDescriptionRules, TestLintDescriptionLimit, TestRenderedContextFitsAntigravityBlock, TestRenderedContextUnderTarget, TestSkillLinksResolve. TestHarnessMap lives in cmd/sh-agent/main_test.go (not internal/harness), so it did not match there; it ran and passed under go test ./... in cmd/sh-agent
- 3. size pins (logged): Antigravity block 14992 bytes (context 14992, rules 0), budget 22000; hand-written context 14107 bytes (render 14992 with inventory), target 15000
- 4. `go test ./... -count=1` -> exit 0, 6 packages ok; `go test ./... -count=1 -v | grep -c '^--- PASS'` = 341 top-level passes (subtests are indented and not counted); no FAIL lines in tail
- 5. `go run ./cmd/sh-agent install --home .agents/scratch/verify-r3-home` -> exit 0: `266 written, 0 unchanged, 0 removed, 0 kept — 266 files managed`
- 5. `go run ./cmd/sh-agent status --home .agents/scratch/verify-r3-home` -> `binary     2ba4a50351ab+dirty` (no warning); `repo has   34 skills, 8 agents, 0 commands, 0 rules`; `drift      none — 265 files match the repo`
- 6. TestMemoryUsageErrorExits2 (cmd/sh-agent/memory_test.go:161) -> `--- PASS`
- 6. dream `%d project(s) failed` runErr passed to finish: covered by TestRunInvalidOutputKeepsRawAndWatermark (internal/dream/run_test.go ~303-311), which asserts st.LastRun.Error contains "failed"; passes
- 7. status count end: 49 (48 + this log); equal to expectation

## Coverage of changed code
- internal/dream/run.go -- TestRunInvalidOutputKeepsRawAndWatermark (asserts only substring "failed", not the full count text)
- cmd/sh-agent/main.go -- TestMemoryUsageErrorExits2
- internal/harness/size_test.go -- is a test
- cmd/sh-agent/memory_test.go, internal/dream/run_test.go -- tests

Suppressions: not scanned for added nolint/skip patterns beyond the Go files above (no diff review of skills/*.md for them).

## Not examined
Suppression/skip grep across the full diff; Docker (no Dockerfile in scope); Antigravity/behavioural runs; govulncheck; pre-existing-failure analysis (none failed).

## Environment
Linux 6.12 amd64, bash; Go per `/usr/bin/go` or user SDK (version not recorded); no Docker. Raw outputs in .agents/scratch/verify-audit-fix-phase-9-r3/.
