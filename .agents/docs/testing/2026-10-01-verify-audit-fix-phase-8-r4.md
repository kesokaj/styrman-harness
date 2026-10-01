---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-8-r4   (2026-10-01)

Range reviewed: working tree on 7c2e30a (Phase 8, round 4: final)

## Status: PASS

Scope: complete for the dispatched commands; guide .agents/docs/testing/guide.md exists. Real-home install and Antigravity runs were not in scope.

## Commands
- `gofmt -l cmd internal` → exit 0, empty
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0; 6 packages ok (root, cmd/sh-agent, dream, harness, memory, search); counts: pass 517, fail 0, skip 3 at 7c2e30a+dirty for go test -race -v -count=1 ./... ; vs previous: .agents/docs/testing/2026-10-01-verify-audit-fix-phase-8-r3.md, pass 515→517, skip 3→3 (output .agents/scratch/verify-p8r4/go-test-1.txt)
- `go test -cover -count=1 ./cmd/... ./internal/harness` → exit 0; cmd/sh-agent 85.8% (>=84), internal/harness 89.0% (>=88)
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-p8r4` → exit 0 (266 written); `status --home ...` → exit 0, `binary 7c2e30a9519b+dirty` (no warning), `drift none — 265 files match the repo`
- Built binary (`go build -o .agents/scratch/verify-p8r4-bin ./cmd/sh-agent` → 0), empty PATH, HOME=scratch dir:
  - `update --home -- --from-zip nowhere.zip` → exit 1; "sh-agent: nowhere.zip: reading the archive: open nowhere.zip: no such file or directory" (zip branch ran; no git fetch, no "up to date")
  - `update --home .agents/scratch/verify-p8r4 -- --from-zip nowhere.zip` → exit 2; `update takes no arguments: ["--from-zip" "nowhere.zip"]`
  - `list --home .agents/scratch/verify-p8r4 --check` → exit 2
- Targeted `-v` run → all PASS, none failing: TestReorder, TestDoubleDashKeepsFollowingArgumentsPositional, TestVerbRefusesStrayArguments, TestUpdateFromZipDecidedByParsedFlags, TestReorderSearchFlags (cmd: ok); TestReadCappedRefusesOversizeStream, TestAttributionSetRefusesALinkedClaudeDir (harness: ok)
- `find skills -type f -perm /111` → empty
- `git status ... | grep -v '^ M\|^D \|^ A\|^A \|^AM\|^??'` → empty (grep exit 1)

## Coverage of changed code
- cmd/sh-agent/main.go — main_test.go (TestUpdateFromZipDecidedByParsedFlags, TestVerbRefusesStrayArguments)
- internal/harness readCapped — phase8_test.go; settings.go — settings_test.go
- Suppressions added: not scanned this run (not in the dispatch).

## Not examined
Suppression scan of the diff; real-home install; Antigravity runs. Tree check: status at end equals start except this log (another artifact, mine); no source/test/config path moved.

## Environment
Linux 6.12 (Debian 13), Go as on PATH (version not recorded), no Docker used.
