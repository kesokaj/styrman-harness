---
owner: verifier
date: 2026-09-14
---
# Verification — upstream items (checkoutShaped fixture, name validation, resync path guard)

Range reviewed: 8046e4d..worktree

## Status: PASS

Scope: gofmt, vet, race suite, targeted new/changed tests, and the Testing
section's context-install check all ran and passed. The sandbox CLI-start
matrix (`.agents/docs/testing/2026-09-03-verify-sandbox.md`) was NOT run —
named gap below, not attempted.

## Commands
- `gofmt -l .` → rc=0, no output (nothing unformatted)
- `go vet ./...` → rc=0, no output
- `go test -race ./... -count=1` → rc=0
  ```
  ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	1.197s
  ok  	github.com/kesokaj/styrman-harness/internal/harness	1.355s
  ```
- `go test ./cmd/sh-agent/ -run 'TestDryRunNeverBuilds|TestPlainCopyNeverRebuilds' -count=1 -v` → rc=0
  ```
  --- PASS: TestDryRunNeverBuilds (0.01s)
  --- PASS: TestPlainCopyNeverRebuilds (0.02s)
  PASS
  ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	0.035s
  ```
- `go test ./internal/harness/ -run 'TestSkillAndAgentNames|TestResyncNever' -count=1 -v` → rc=0
  ```
  --- PASS: TestSkillAndAgentNamesMustBeOneLocalSegment (0.00s)
  --- PASS: TestResyncNeverFollowsANonLocalManifestEntry (0.00s)
  PASS
  ok  	github.com/kesokaj/styrman-harness/internal/harness	0.002s
  ```
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-2026-09-14` → rc=0
  ```
  building sh-agent from /home/kesokaj/workspace/styrman-harness
  260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed
  note: .agents/scratch/verify-2026-09-14/.local/bin is not on PATH; ...
  ```
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-2026-09-14` → rc=0
  ```
  repo       .../styrman-harness (main @ 8046e4d, uncommitted changes)
  ...
  drift      none — 259 files match the repo
  ```
  Reports "drift none" as required (the "259 vs 260" figure is status's own
  file-count accounting, not a drift finding — status explicitly says "none").
- `grep -c 'standing permission, run it' .agents/scratch/verify-2026-09-14/.claude/AGENTS.md` → `1`
- `file internal/harness/pathguard_test.go` → `internal/harness/pathguard_test.go: ASCII text` (confirms BUG-007-family concern does not apply — the new test file is real text, not binary/control-byte content)
- `git diff --numstat internal/harness/pathguard_test.go` → new untracked file, not yet diffed against a committed blob (file is `??` in status); read directly and confirmed plain ASCII via `file`.

All verification commands above were run bare, as their own statement, with
exit status read before any further action (BUG-033/041 / `verify-bare.md`).
No commit was made in this session.

## checkoutShaped fixture — read, not re-run via git stash (git stash forbidden by dispatch)
Read `cmd/sh-agent/main_test.go` lines 195-260. `TestDryRunNeverBuilds` calls
`repo := checkoutShaped(t)`, NOT `filepath.Abs("../..")` directly (that call
only appears in two OTHER tests, lines 122 and 169, and inside `plainCopy`
purely as the SOURCE tree to copy content from). `checkoutShaped` =
`plainCopy(t)` (copies `context/`, `skills/`, `agents/`, `commands/`,
`rules/` from `../..` into a fresh `t.TempDir()` with no `.git`) plus an
`os.Mkdir(.../.git)` on that fresh copy. So the directory passed to
`run([]string{"install", "--dry-run", "--repo", repo, ...})` is a brand-new
temp dir shaped like a checkout (has `.git`) but is NOT `../..` and is NOT a
real git repository. This is exactly what makes the fixture immune to
running inside a downstream mirror (styrman-crew's plain-copy image/harness,
BUG-058 in that repo's own bug log): `../..`'s only role left is as a content
source for the copy, never as the path evaluated for "is this a checkout."

## Root cause
N/A — PASS, no defects found in the reviewed diff.

## Coverage of changed code
- `cmd/sh-agent/main_test.go` (checkoutShaped fixture) — covered by
  `TestDryRunNeverBuilds` and `TestPlainCopyNeverRebuilds`, both run above.
- `internal/harness/inventory.go` (`validName`) — covered by
  `TestSkillAndAgentNamesMustBeOneLocalSegment` in the new
  `internal/harness/pathguard_test.go`.
- `internal/harness/apply.go` (resync loop's `filepath.IsLocal` skip) —
  covered by `TestResyncNeverFollowsANonLocalManifestEntry` in the same new
  file.
- `context/dev-setup.md` (line 8, standing-permission wording) — not code;
  verified indirectly via the installed-home grep
  (`grep -c 'standing permission, run it' .../verify-2026-09-14/.claude/AGENTS.md` = 1),
  confirming the edited sentence is emitted correctly into the assembled
  context and appears exactly once.

## Not examined
- The sandbox CLI-start matrix (`.agents/docs/testing/2026-09-03-verify-sandbox.md`)
  — one start of each CLI (Claude Code, Antigravity/agy) in the sandbox with
  logs grepped for parse errors — was NOT run this session. Named per
  dispatch instruction, not attempted.
- The full skill matrix (`.agents/docs/testing/skill-matrix.md`) was not run;
  no claim of "complete"/"all skills" coverage is made.
- The real-home smoke checks (`.agents/docs/testing/2026-09-03-verify-real-home.md`)
  were not run — this verification used only a scratch home
  (`.agents/scratch/verify-2026-09-14`), never the operator's real `$HOME`,
  per this repo's own rule.
- No `git stash` was used at any point (forbidden by dispatch); the
  `checkoutShaped` claim above was established purely by reading the test
  source, not by re-running it under a stashed/plain-copy tree.

## Environment
- OS: Linux 6.12.107+deb13-amd64
- `go version`: not separately captured this session; `go test`/`go vet`/`go run` all executed without toolchain errors, implying a compatible Go toolchain is on PATH.
- No Docker/Colima involved (this repo has no container build step in scope).
- Working tree: uncommitted changes present as described in the dispatch (main @ 8046e4d, dirty); no commit made by this verification.

## Round 2

Range reviewed: 8046e4d..worktree (post security-review fixes: symlink
refusal in `plan.go`/`adopt.go`, `LoadState` dropping non-local manifest
entries in `state.go`, `identLike` charset added to `validName` in
`inventory.go`, expanded `pathguard_test.go`, reworded `context/dev-setup.md`
line 8 tail).

## Status: PASS

## Commands
- `gofmt -l .` → rc=0, no output
- `go vet ./...` → rc=0, no output
- `go test -race ./... -count=1` → rc=0
  ```
  ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	1.199s
  ok  	github.com/kesokaj/styrman-harness/internal/harness	1.354s
  ```
- `go test ./internal/harness/ -run 'TestSkillAndAgentNames|TestResyncNever|TestLoadStateDrops|TestSkillFilesAndAdopt' -count=1 -v` → rc=0
  ```
  --- PASS: TestSkillAndAgentNamesMustBeOneLocalSegment (0.00s)
  --- PASS: TestResyncNeverFollowsANonLocalManifestEntry (0.00s)
  --- PASS: TestLoadStateDropsNonLocalManifestEntries (0.00s)
  --- PASS: TestSkillFilesAndAdoptCopiesRefuseSymlinks (0.00s)
  PASS
  ok  	github.com/kesokaj/styrman-harness/internal/harness	0.004s
  ```
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-2026-09-14-r2` (fresh scratch home) → rc=0
  ```
  building sh-agent from /home/kesokaj/workspace/styrman-harness
  260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed
  note: .../.local/bin is not on PATH; ...
  ```
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-2026-09-14-r2` → rc=0, `drift      none — 259 files match the repo`
- `grep -c 'Without that record, hand the step back' .agents/scratch/verify-2026-09-14-r2/.claude/AGENTS.md` → `1`
- `go run ./cmd/sh-agent list --repo .` → rc=0, output contains `### Skills (27)` and `### Subagents (7)`

All commands run bare, each its own statement, exit status read before any
further action; no commit made this round.

## Coverage of changed code (round 2 delta)
- `internal/harness/plan.go`, `internal/harness/adopt.go` (symlink refusal
  in the skill walk / adopt copy) — covered by
  `TestSkillFilesAndAdoptCopiesRefuseSymlinks`.
- `internal/harness/state.go` (`LoadState` drops non-local manifest entries;
  `apply.go`'s per-caller guard removed) — covered by
  `TestLoadStateDropsNonLocalManifestEntries`.
- `internal/harness/inventory.go` (`identLike` charset added to
  `validName`) — covered by the expanded cases in
  `TestSkillAndAgentNamesMustBeOneLocalSegment`.
- `context/dev-setup.md` line 8 reword ("Without that record, hand the step
  back with the exact command instead.") — verified via the installed-home
  grep, count = 1.
- `.agents/docs/bugs/INDEX.md`, `BUG-016-*.md`, the two review files — trail
  artifacts, not code; not subject to test coverage.

## Not examined (unchanged from round 1)
- Sandbox CLI-start matrix, full skill matrix, and real-home smoke checks
  still not run this round — same named gaps as round 1.
