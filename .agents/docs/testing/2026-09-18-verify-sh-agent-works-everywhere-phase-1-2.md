---
owner: verifier
date: 2026-09-18
---
# Verification — sh-agent-works-everywhere-phase-1 (second pass)   (2026-09-18)

Range reviewed: b393525..d737282

## Status: PASS

Scope: short confirmation pass — gofmt/vet/test/govulncheck, the targeted new-case test, and one live install leg (fresh install + directory-clobber). Complete for what was asked.

## Commands
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -count=1 -v ./...` → exit 0; counts: pass 227, fail 0, skip 0 at d737282 for `go test -race -count=1 -v ./...`; vs previous: `.agents/docs/testing/2026-09-18-verify-sh-agent-works-everywhere-phase-1.md` (pass 227, fail 0, skip 0) → pass 227→227, skip 0→0. Matches the expected 227 (no test added since the first log).
- `go test -race -count=1 -v -run 'TestInstallNamesEveryChange' ./internal/harness/` → exit 0:
  ```
  === RUN   TestInstallNamesEveryChange
  --- PASS: TestInstallNamesEveryChange (0.01s)
  PASS
  ok  	github.com/kesokaj/styrman-harness/internal/harness	1.019s
  ```
- `govulncheck ./...` → exit 0, `No vulnerabilities found.` (go version go1.27.1 linux/amd64, ahead of the machine's default go per AGENTS.md).

### Live install leg
1. `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase-1-2` → exit 0. `grep -c '^add '` on the output = 260; summary line = `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`. Counts match.
2. Directory-clobber: the dispatch's example path `.claude/rules/go-style.md` does not exist in this repo's deploy tree (no `rules/*.md` source file — `rules/` holds only `.gitkeep`, and the destination `.claude` has no `rules/` subdirectory at all). Substituted a real deployed file, `.claude/agents/bugs.md`: removed it and made a directory in its place (`mkdir .claude/agents/bugs.md`), then reran install.
   - Result: exit 1.
   - Output (verbatim):
     ```
     building sh-agent from /home/kesokaj/workspace/styrman-harness
     sh-agent: writing .agents/scratch/verify-phase-1-2/.claude/agents/bugs.md: rename .agents/scratch/verify-phase-1-2/.claude/agents/.sh-agent-3004113953 .agents/scratch/verify-phase-1-2/.claude/agents/bugs.md: file exists
     sh-agent: sh-agent install: exit status 1
     exit status 1
     ```
   - The error names the exact clobbered path (`.claude/agents/bugs.md`). No `update` line (or any per-file line) was printed for it — confirmed by `grep 'bugs.md'` on the output, which returns only the error line, no `update` line.
   - Matches the described Phase 1 behaviour: a directory at a destination fails the install and prints no per-file line for that path.

## Root cause
N/A — PASS.

## Coverage of changed code
No source changed in this range beyond docs/records (per dispatch, "Go sources frozen"); this is a confirmation pass over the prior Phase 1 commit (028a43b) plus a docs/records commit (d737282). No new suppressions or skips observed in the range (not separately diffed this pass — see Not examined).

## Not examined
- Did not diff the range for added `//nolint`, `t.Skip`, `.only(` etc. (dispatch scoped this to a short confirmation pass over specific commands; the first log at 028a43b already covers the suppression sweep for the code-bearing commit).
- Did not run the `--dry-run` install leg (dispatch specified one live leg only).
- Did not exercise CLI smoke checks (claude/gemini sandbox starts) — out of scope for this pass per dispatch.
- Did not test the go-style.md path literally since it does not exist in this repo; substituted an equivalent real path and noted the discrepancy above.

## Environment
Linux 6.12.107+deb13-amd64, go1.27.1 (user-level SDK ahead of /usr/bin/go per AGENTS.md), govulncheck from /home/kesokaj/.go/bin/govulncheck.
