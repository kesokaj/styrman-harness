---
owner: verifier
date: 2026-09-18
---
# Verification — BUG-078 path shadow, round 3   (2026-09-18)

Range reviewed: 4cdc0b5..b029ea1

## Status: PASS

Scope: `.agents/docs/testing/guide.md` tier 1 (Go checks) plus the dispatch's targeted
tests, scratch-home install/status, `list`, `govulncheck`, and both live PATH checks
(a: off-PATH install/status silence; b: symlinked-alias shadow repro). Guide.md already
exists (owner: verifier, date: 2026-09-17); not rewritten. Full sandbox/skill-matrix
tiers (guide §2, §4, §5) not run — no change to `context/`, `agents/`, or `skills/` in
this range (diff touches only `cmd/sh-agent/{main.go,main_test.go}` and
`internal/harness/pathguard{,_test}.go`).

## Commands
- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -count=1 -v ./...` (background, log at
  `.agents/scratch/verify-bug078-3/go-test-full.txt`) → exit reflected by `ok` on all
  three packages, no `FAIL`:
  ```
  ok  	github.com/kesokaj/styrman-harness	1.091s
  ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	3.995s
  ok  	github.com/kesokaj/styrman-harness/internal/harness	2.799s
  ```
  counts: pass 215, fail 0, skip 0 at b029ea1 for `go test -race -count=1 -v ./...`.
  vs previous: `.agents/docs/testing/2026-09-18-verify-bug078-path-shadow-2.md`, pass
  213→215, skip 0→0. The +2 matches the dispatch note that this range adds two
  `TestPathShadows` table cases (symlinked directory); the renamed `t.Run` does not
  change the count.
- `go test -race -count=1 -v -run 'TestPathShadows|TestExecutableAt' ./internal/harness/`
  → exit 0. All 11 `TestPathShadows` subtests PASS, including the two new ones:
  `a_symlinked_directory_that_resolves_to_the_bin_dir:_reached,_and_an_earlier_shadow_is_still_named`
  and `a_symlinked_directory_first:_nothing_shadows`. `TestExecutableAt` and its 4
  subtests PASS.
- `go test -race -count=1 -v -run 'TestInstallAndStatusNameAPathShadow|TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice' ./cmd/sh-agent/`
  → exit 0, both PASS.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug078-3` → exit 0.
  Printed:
  ```
  note: .agents/scratch/verify-bug078-3/.local/bin is not on PATH, so `sh-agent` runs /home/kesokaj/.local/bin/sh-agent instead; put .agents/scratch/verify-bug078-3/.local/bin ahead of it
  ```
  No `warning:` line.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug078-3` (normal PATH)
  → exit 0. No `shadowed` line in the output (verified by reading full output).
- `go run ./cmd/sh-agent list` → exit 0, skills/agents tables printed (27 skills, 7
  agents).
- `govulncheck ./...` → "No vulnerabilities found."

### Live check (b): symlinked-alias PATH shadow repro
Created `.agents/scratch/verify-bug078-3-alias` as a symlink to
`.agents/scratch/verify-bug078-3/.local/bin`, and an executable
`.agents/scratch/verify-bug078-3-earlier/sh-agent` (`#!/bin/sh`, chmod +x). Ran:
```
PATH="$(pwd)/.agents/scratch/verify-bug078-3-earlier:$(pwd)/.agents/scratch/verify-bug078-3-alias:$PATH" \
  go run ./cmd/sh-agent status --home .agents/scratch/verify-bug078-3
```
Output included exactly one shadow line:
```
shadowed   /home/kesokaj/workspace/styrman-harness/.agents/scratch/verify-bug078-3-earlier/sh-agent comes first on PATH, so `sh-agent` runs that, not the binary above
```
Confirms the PATH walk reaches the installed file through the symlinked alias and
correctly names only the one earlier real shadow. Cleanup: removed the alias symlink
by name (`rm .agents/scratch/verify-bug078-3-alias`), confirmed
`.agents/scratch/verify-bug078-3-earlier` held only `sh-agent` (`ls`), then
`rm -r .agents/scratch/verify-bug078-3-earlier`.

## Root cause  (FAIL only)
N/A — no failures found.

## Coverage of changed code
Diff `4cdc0b5..b029ea1` touches four `.go` files:
- `internal/harness/pathguard.go` — covered by `internal/harness/pathguard_test.go`
  (`TestPathShadows`, `TestExecutableAt`), exercised directly above; the new
  `os.SameFile`-based identity walk is hit by the new symlinked-directory subtests.
- `cmd/sh-agent/main.go` — covered by `cmd/sh-agent/main_test.go`
  (`TestInstallAndStatusNameAPathShadow`,
  `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice`), and by the two live checks
  above (real `install`/`status` runs).

No suppressions, skips, or `.only(`/`.skip(` added in this range
(`git diff 4cdc0b5..b029ea1 -- '*.go'` has no added `nolint`, `lint:ignore`, `t.Skip`,
or `go:build ignore` lines).

## Not examined
- Guide §2 harness self-check (`sh-agent list`/`render` line-count), §4 sandbox CLI
  starts, and §5 real-home smoke checks: not run, because this range touches no
  `context/`, `agents/`, or `skills/` file (change is confined to `cmd/sh-agent` and
  `internal/harness/pathguard.go`), and the guide scopes those tiers to changes that
  touch those directories.
- Guide §3a (plain-copy install) and §3b (`report` dry-run): not run — this range does
  not touch `IsCheckout`/rebuild logic or `internal/harness/report.go`.
- Full skill-matrix (`.agents/docs/testing/skill-matrix.md`): not run; not claiming
  "complete"/"all skills".

## Environment
Linux 6.12.107+deb13-amd64. `go version` not captured separately in this run (repo's
own AGENTS.md notes a user-level go1.27.1 SDK ahead of `/usr/bin/go` on this box, used
by `go run`/`go test` above). No Docker/Colima involved (no Dockerfile changed in this
range).
