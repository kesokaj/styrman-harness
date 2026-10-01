---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-8-r2   (2026-10-01)

Range reviewed: working tree on 7c2e30a (Phase 8, round 2)

## Status: FAIL

Scope: every requested step ran; one expectation (step 6 takesValue grep empty) not met; guide exists (.agents/docs/testing/guide.md).

## Commands
- `gofmt -l cmd internal` → exit 0, empty output
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0, all 6 packages ok; counts: pass 515, fail 0, skip 3 at 7c2e30a+dirty for go test -race -v -count=1 ./... (output .agents/scratch/verify-p8r2/go-test-1.txt); vs previous: not compared
- `go test -cover -count=1 ./cmd/... ./internal/harness` → exit 0; cmd/sh-agent 85.7% (>=84), internal/harness 89.0% (>=88)
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-p8r2` → exit 0, 266 written; `status` → exit 0, `binary 7c2e30a9519b+dirty` no warning, `drift none — 265 files match the repo`
- built binary B=.agents/scratch/verify-p8r2-bin: `list --check` exit 2 (flag not defined, usage); `version extra` exit 2 "version takes no arguments: [\"extra\"]"; `list - --check --home H` exit 2 "list takes no arguments: [...]"; `update -- --from-zip nowhere.zip --home H` exit 2 "update takes no arguments: [...]" and nothing fetched; `search -- --from-zip --home H` exit 0 (ran as a query, 10 hits); `adopt --help` prints usage, exit 1 (flag.ContinueOnError path; usage printed as expected, exit not in the dispatch's expectations)
- linked settings: `install --home verify-p8r2-link` with .claude -> target dir: exit 0; output "settings: a directory above .../.claude/settings.json is a symlink; add the key to its target by hand"; `ls target/settings.json` -> No such file (as required). Link dir removed afterwards.
- targeted `go test -run ... -v -count=1` → exit 0; PASS: TestReorder, TestVerbRefusesAnotherVerbsFlag, TestDoubleDashKeepsFollowingArgumentsPositional, TestVerbRefusesStrayArguments, TestReorderSearchFlags, TestUpdateFromZipDryRunCleansTempDir, TestNoSystemTempDirInSource, TestVerifyZipSumsWritesNothingUnderTempDir, TestReadCappedRefusesOversize, TestAttributionSetRefusesALinkedClaudeDir; 0 FAIL
- greps: `grep -rn takesValue cmd internal guide.md` → NOT EMPTY: cmd/sh-agent/main.go:936,941,954 (a live local variable inside reorder(), not a stale list; guide.md has none); `exit status is 2` README → 1; `2 on a usage error` search.md → 1; `beside the zip` delivery.md 1, README.md 1
- `govulncheck ./...` → exit 0; 0 affecting vulnerabilities (1 in a required module, not called)

## Root cause  (FAIL only)
Only the takesValue grep expectation failed: main.go reorder() defines `takesValue` as a local map built from the FlagSet (comment documents it). Either the expectation is stale (it targeted a removed hand-kept list) or the name should be absent; caller to decide. No functional failure found.

## Coverage of changed code
Not re-mapped this run beyond coverage figures (cmd 85.7%, harness 89.0%) and the targeted tests above; settings.go covered by settings_test.go, phase8_test.go, repo_origin_test.go exist. No suppression scan done.

## Not examined
Per-file coverage map and suppressions diff scan; vs-previous test count comparison (previous log not read); behavioural/§3 tiers beyond those listed.

## Environment
Linux 6.12 (Debian 13), Go from PATH, govulncheck at ~/.go/bin. Docker not relevant. Tree status at end equals start plus this log only.
