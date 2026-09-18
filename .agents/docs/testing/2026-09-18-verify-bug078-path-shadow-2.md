---
owner: verifier
date: 2026-09-18
---
# Verification — BUG-078 review round 1 (2026-09-18)

Range reviewed: 4cdc0b5..8ce58a0

## Status: PASS

Scope: complete — testing/guide.md exists (no changes needed to it); all AGENTS.md §
Setup commands ran, the four named tests ran individually with `-count=1 -v`, the
whole suite ran with `-count=1` and again with `-count=1 -v` for counts, the live
(un-fakeable) install/status shadow check ran against a real scratch home not on
PATH, `sh-agent list` and `govulncheck ./...` ran.

## Commands

- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -count=1 ./...` → exit 0, `ok` in all three packages
  (`github.com/kesokaj/styrman-harness`, `.../cmd/sh-agent`, `.../internal/harness`).
- `go test -race -count=1 -v ./...` (for counts) → exit 0, `PASS` in all three packages.
  counts: pass 108 top-level + 105 subtests (213 `--- PASS:` lines total), fail 0,
  skip 0, at 8ce58a0+dirty (uncommitted testing-log file only, no source diff) for
  `go test -race -count=1 -v ./...`.
  vs previous: `.agents/docs/testing/2026-09-18-verify-bug078-path-shadow.md`
  (52574a0), pass 108+103=211 → 108+105=213, skip 0→0.
  explained: `TestResolveExecutable` (6 subtests: executable file, symlink-resolves,
  no-execute-bit, dangling symlink, directory, missing path) was replaced by
  `TestExecutableAt` (4 subtests: no-execute-bit, dangling symlink, directory,
  missing path) — the "executable file" and "symlink resolves to its target" cases
  moved into `TestPathShadows` as identity subtests, since `PathShadows` now uses
  `os.SameFile` directly rather than a separate resolver (review round 1). Net
  `TestExecutableAt`: -2. `TestPathShadows` grew from 10 to 14 subtests (+4): added
  "a hard link to the installed binary is not a shadow", "a double trailing slash",
  "a dot segment in the installed dir", "an installed binary that does not exist yet
  still reports shadows" (the old symlink/no-repeat cases were renamed/rescoped, not
  duplicated — table read directly, see below). Net +4-2 = +2, matching 211→213.
- `go test -race -count=1 -v -run 'TestPathShadows|TestExecutableAt' ./internal/harness/`
  → exit 0. `TestPathShadows`: 14/14 subtests PASS. `TestExecutableAt`: 4/4 subtests
  PASS. `PASS`, `ok ... 1.018s`.
- `go test -race -count=1 -v -run 'TestInstallAndStatusNameAPathShadow|TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice' ./cmd/sh-agent/`
  → exit 0. Both PASS (0.09s, 0.06s). `PASS`, `ok ... 1.172s`.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug078-2` → exit 0.
  Operator's real PATH used (`.local/bin` first, holding a real `sh-agent` binary,
  none of the scratch home's bin dirs on it):
  ```
  building sh-agent from /home/kesokaj/workspace/styrman-harness
  260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed
  note: .agents/scratch/verify-bug078-2/.local/bin is not on PATH, so `sh-agent` runs /home/kesokaj/.local/bin/sh-agent instead; put .agents/scratch/verify-bug078-2/.local/bin ahead of it
  settings: .agents/scratch/verify-bug078-2/.claude/settings.json created with "attribution": {"commit": "", "pr": ""} — Claude Code's commit and PR trailers are off
  ```
  Confirmed `/home/kesokaj/.local/bin/sh-agent` is a real file on PATH
  (`ls -la` shows a 7MB executable; `which -a sh-agent` resolves to it first).
  No `warning: ... comes before` line printed (correct — the scratch home's bin dir
  is off PATH, so the "note" branch is the right one, not "warning").
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug078-2` → exit 0.
  ```
  repo       /home/kesokaj/workspace/styrman-harness (main @ 8ce58a0, uncommitted changes)
  installed  2026-09-18T12:53:18+02:00 from commit 8ce58a0+dirty (vendor: all, 260 files, manifest ~/.config/sh-agent/state.json)
  binary     8ce58a0ef558+dirty
  retro      none yet — run /retro at the end of a session with corrections
  repo has   27 skills, 7 agents, 0 commands, 0 rules
  claude     ~/.claude  (AGENTS.md, CLAUDE.md, skills/, agents/, commands/, rules/)
  agy        ~/.gemini/config  (rules/AGENTS.md, skills/, agents/, workflows/, rules/)
  drift      none — 259 files match the repo
  ```
  No `shadowed` line printed — correct, since the scratch home's bin dir is not on
  PATH so `PathShadows` returns `onPath=false` and `status`'s shadow block only
  fires when `onPath` is true.
- `go run ./cmd/sh-agent list` → prints skills/agents tables; header counts
  "Skills (27)" match repo-has line; no parse errors.
- `govulncheck ./...` → "No vulnerabilities found." exit 0.

## Root cause

N/A — all PASS.

## Coverage of changed code

Diff `4cdc0b5..8ce58a0` touches only `internal/harness/pathguard.go` (new),
`internal/harness/pathguard_test.go` (new), `cmd/sh-agent/main.go`,
`cmd/sh-agent/main_test.go`, plus docs/README/reviews (no Go).

- `internal/harness/pathguard.go`: `PathShadows` — covered by
  `TestPathShadows` (14 table subtests, internal/harness/pathguard_test.go).
  `executableAt` — covered directly by `TestExecutableAt` (4 subtests) and
  indirectly by every `TestPathShadows` subtest that exercises identity/no-execute
  cases.
- `cmd/sh-agent/main.go`: `pathHint` (rewritten to call `PathShadows` once) —
  covered by `TestOffPathInstallNamesWhatAnswersAndDoesNotWarnTwice` (off-PATH
  "note" branch, live-checked above) and by the on-PATH "warning" branch inside
  `TestInstallAndStatusNameAPathShadow`. `status`'s new shadow block (reading
  `harness.PathShadows` and printing `shadowed`) — covered by
  `TestInstallAndStatusNameAPathShadow` (cmd/sh-agent/main_test.go); live-checked
  above showing the `onPath=false` case correctly suppresses the line.
  No gap: every changed non-test Go function is named by a test that exercises it.
- No suppressions added: `git diff 4cdc0b5..8ce58a0 -- '*.go'` grepped for
  `nolint|lint:ignore|t\.Skip|go:build ignore|\.only\(|\.skip\(` — no matches.

## Not examined

- `.agents/docs/testing/guide.md` §2 (harness self-check `render | wc -l`), §3a
  (plain-copy install), §3b (report dry-run), §4 (sandbox CLI starts), §5
  (real-home smoke checks) were not run: this range touches only
  `internal/harness/pathguard.go` and `cmd/sh-agent/main.go` (PATH-shadow logic),
  not `context/`, `agents/`, `skills/`, the checkout-vs-mirror rebuild path, or
  `report.go`, so those tiers are out of scope per the guide's own trigger
  conditions.
- The three non-Go review/docs artifacts and README/local.md/remote.md prose
  changes in the diff (`.agents/docs/reviews/*`, `README.md`,
  `.agents/docs/deploy/{local,remote}.md`) were read for suppression-grepping
  only, not proofread for content accuracy — that is `/review --docs` territory,
  not this verification pass.
- `sh-agent report` was not run (hard limit — never run it from a verification
  session).

## Environment

Linux 6.12.107+deb13-amd64, go1.27.1 (per repo's documented user-level SDK ahead
of /usr/bin/go), govulncheck clean. Working tree: only the new/updated testing-log
files under `.agents/docs/testing/` are untracked/dirty; no source file was
modified during this run. Scratch home used:
`.agents/scratch/verify-bug078-2` (not on operator's real PATH).
