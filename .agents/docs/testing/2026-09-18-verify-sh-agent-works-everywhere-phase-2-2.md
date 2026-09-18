---
owner: verifier
date: 2026-09-18
---
# Verification — sh-agent-works-everywhere-phase-2 (round 2)   (2026-09-18)

Range reviewed: 3209cf3..e7fec88

## Status: PASS

Scope: complete — gofmt/vet/test-race, three cross-compiles, targeted -run subset, live install/stat with tighten-back-to-700 check, govulncheck, coverage mapping. testing/guide.md exists (unchanged).

## Commands

- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -count=1 ./...` → exit 0, all 3 packages ok:
  ```
  ok  	github.com/kesokaj/styrman-harness	1.085s
  ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	6.474s
  ok  	github.com/kesokaj/styrman-harness/internal/harness	2.780s
  ```
- `go test -race -count=1 -v ./...` (for counts) → exit 0. counts: pass 248, fail 0, skip 0 at e7fec88 for `go test -race -count=1 -v ./...` (counted via `grep -c '^\s*--- (PASS|FAIL|SKIP):'`, matching indented subtest lines; also cross-checked `=== RUN` = 248).
  vs previous: `.agents/docs/testing/2026-09-18-verify-sh-agent-works-everywhere-phase-2.md`, pass 244→248, skip 0→0. Matches the dispatch's prediction exactly (+2 subtests of `TestInstallMakesTheConfigDirPrivate`, +1 `TestNoteRunningElsewhereThroughASymlinkedHome`, +1 net from `TestOwnedByCaller` moving to a unix-only file with no count change) = +4.
- `GOOS=windows go build ./...` → exit 0 (WIN_OK).
- `GOOS=darwin GOARCH=arm64 go build ./...` → exit 0 (DARWIN_OK).
- `GOOS=linux GOARCH=arm64 go build ./...` → exit 0 (LINUXARM_OK).
- `go test -race -count=1 -v -run 'TestInstallMakesTheConfigDirPrivate|TestThrowawayParentIsPrivate|TestNoteRunningElsewhereThroughASymlinkedHome|TestNoteRunningElsewhereOnlyForTheRealHome|TestOwnedByCaller|TestDiffUsesTheHomeNotTmp' ./internal/harness/` → exit 0. All 6 named tests PASS (8 including the 2 subtests of `TestInstallMakesTheConfigDirPrivate`: `absent`, `0755`), 0 fail, 0 skip.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase-2-2` → succeeded, 260 files managed. `stat -c '%a' .agents/scratch/verify-phase-2-2/.config/sh-agent` → `700`.
- Live tightening check: `chmod 755 .agents/scratch/verify-phase-2-2/.config/sh-agent` → stat confirmed `755`; re-ran `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase-2-2` (exit 0); `stat -c '%a' .agents/scratch/verify-phase-2-2/.config/sh-agent` → `700`. Confirms `ConfigDir` tightens an already-present, loosened directory on every install, not only on first creation.
- `govulncheck ./...` → "No vulnerabilities found." (go1.27.1 on this box, consistent with AGENTS.md's note that a modern toolchain reports clean).

## Root cause
N/A — Status PASS, nothing failed.

## Coverage of changed code
Changed `.go` files in range (`git diff --stat 3209cf3..e7fec88 -- '*.go'`):
- `internal/harness/apply.go` — exercised by `TestInstallMakesTheConfigDirPrivate` and the full install test suite (package ok).
- `internal/harness/audit_test.go` (test file itself, expanded) — ran as part of `go test`, all subtests pass.
- `internal/harness/owned_other.go` / `owned_unix.go` — covered by `owned_unix_test.go` (new) and `TestOwnedByCaller`, both run explicitly above and pass. `owned_other.go` (non-unix build) is exercised only by cross-compile (`GOOS=windows`), which confirms it compiles; no runtime test possible from this Linux box — noted under Not examined.
- `internal/harness/pathguard.go` — has `pathguard_test.go`; part of the full-package pass (248/248).
- `internal/harness/release.go` — has `release_test.go`; part of the full-package pass.
- `internal/harness/report.go` — has `report_test.go`; part of the full-package pass.
- `internal/harness/state.go` (`ConfigDir` addition) — covered by `TestInstallMakesTheConfigDirPrivate`, `TestThrowawayParentIsPrivate`, `TestDiffUsesTheHomeNotTmp`, all run explicitly and pass, plus the live tighten-to-700 reproduction above.

No suppressions (`//nolint`, `t.Skip`, `.only(`, etc.) found added in this range — the diff is source + test only, no skip/ignore directives present in the changed hunks.

## Not examined
- `owned_other.go`'s runtime behavior (non-unix `ownedByCaller` stub) — only compiled via `GOOS=windows go build`, never executed (no Windows host available here).
- Full behavioural CLI smoke run (`.agents/docs/testing/2026-09-03-verify-sandbox.md` steps: start each CLI in the sandbox and grep for parse errors) was not re-run this round — the dispatch did not name it and no `agents/`/`context/`/`skills/` files changed in this range (confirmed: diff was Go-source only, `*.go` files listed above).
- Unrelated to this range: the working tree at verification time also showed a modified `.agents/docs/plans/2026-09-18-sh-agent-works-everywhere.md` and an untracked `.agents/docs/testing/2026-09-18-verify-sh-agent-works-everywhere-phase-3.md`, apparently from concurrent Phase 3 work in the same checkout. Neither was touched by this verifier; flagging since the tree is supposed to stay frozen during a dispatch (per AGENTS.md Definition of Done step 1) — did not affect the Go-only range under test here, but the main agent should confirm no overlap with this dispatch's scope.

## Environment
OS: Linux 6.12.107+deb13-amd64. Go: go1.27.1 linux/amd64 (user-level SDK, per AGENTS.md). No Docker/Colima involved (no Dockerfile changed).
