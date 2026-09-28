---
owner: verifier
date: 2026-09-28
---
# Verification — v1.8.0 follow-ups   (2026-09-28)

Range reviewed: e4ae34c..worktree

## Status: PASS

Scope: gofmt, go vet, go test -race (all pkgs), context-size test, govulncheck, install/status drift, list count, tracked-files git status, embed test, numstat. Complete — every command below ran and testing/guide.md exists.

## Commands

1. `gofmt -l cmd internal` → exit 0, no output (clean).
2. `go vet ./...` → exit 0, no output.
3. `go test -race -count=1 -timeout 180s -v ./...` (bare, output saved to `.agents/scratch/verify-v180-followups/go-test-race.txt`) → all 3 packages `ok` (`github.com/kesokaj/styrman-harness` 1.086s, `.../cmd/sh-agent` 11.157s, `.../internal/harness` 2.916s), no FAIL lines; counts: pass 164, fail 0, skip 0 at e4ae34c+dirty for `go test -race -count=1 -timeout 180s -v ./...`; vs previous: `.agents/docs/testing/2026-09-28-verify-v1.8.0-gate.md`, pass 162→164, skip 0→0 (rise in pass, no fall, no rise in skip — no FAIL condition).
4. Context size: `go test -count=1 -run TestRenderedContextFitsAntigravityBlock -v ./internal/harness` → `size_test.go:48: Antigravity rules block is 20822 bytes (context 20822, rules 0), budget 22000`; PASS.
5. `govulncheck ./...` (backgrounded, output saved to `.agents/scratch/verify-v180-followups/govulncheck.txt`) → "No vulnerabilities found." (go1.27.1, matches AGENTS.md's "clean on go1.25.8+").
6. `go run ./cmd/sh-agent install --home .agents/scratch/fu` → succeeded, 258 files written (0 unchanged, 0 removed), both `~/.claude` and `~/.gemini/config` trees populated (scratch home, never the real `$HOME`).
7. `go run ./cmd/sh-agent status --home .agents/scratch/fu` → `drift  none — 257 files match the repo`; repo shown as `main @ e4ae34c, uncommitted changes` (expected, working tree is dirty); 27 skills, 7 agents, 0 commands, 0 rules reported.
8. Both vendors' rendered context checked directly in the scratch home:
   - `tail -c 80 .agents/scratch/fu/.claude/AGENTS.md` → ends `Canary: this context ends with the word GANNET.`
   - `tail -c 80 .agents/scratch/fu/.gemini/config/rules/AGENTS.md` → same ending.
   - `grep -n "found elsewhere" <both files>` → both contain the new sentence: "a skill found elsewhere (`find-skills`, `npx skills add`) goes into its checkout's `skills/` — ... then `sh-agent install`, never only into a vendor directory."
9. `go run ./cmd/sh-agent list` → `### Skills (27)` and `### Subagents (7)` — matches expected counts.
10. `git status --porcelain --ignored --untracked-files=all -- context skills agents commands rules` → only tracked modifications:
    ```
     M agents/bugs.md
     M context/AGENTS.md
     M skills/VENDORED.md
     M skills/find-skills/SKILL.md
    ```
    Nothing untracked, nothing ignored, under those paths.
11. `go test -count=1 -run TestEmbedHoldsOnlyTrackedFiles ./...` → PASS in the root package (`github.com/kesokaj/styrman-harness`); "no tests to run" in `cmd/sh-agent` and `internal/harness` (the test lives only in the root package — expected, exit 0 both).
12. `git diff --numstat` → all rows numeric (no `-`/binary markers):
    ```
    2  2  .agents/docs/testing/skill-matrix.md
    2  2  REQUIREMENTS.md
    1  1  agents/bugs.md
    32 1  cmd/sh-agent/zip_test.go
    1  1  context/AGENTS.md
    6  4  internal/harness/repo.go
    11 1  internal/harness/report.go
    3  1  internal/harness/zip.go
    21 0  internal/harness/zip_test.go
    1  1  skills/VENDORED.md
    4  7  skills/find-skills/SKILL.md
    ```

## Root cause  (FAIL only)

N/A — no FAIL.

## Coverage of changed code

- `internal/harness/repo.go` — covered by `internal/harness` package tests (part of the `go test -race ./...` run above; no `_test.go` diff for this file itself, but exercised transitively — not independently verified beyond the full-suite pass; see Not examined).
- `internal/harness/report.go` — covered similarly.
- `internal/harness/zip.go` — covered by `internal/harness/zip_test.go`, which itself gained a new 21-line hunk in this change set (net-new test lines, no removed lines).
- `cmd/sh-agent/zip_test.go` — 32/1 diff, itself a test file, exercised by the `cmd/sh-agent` package pass above.
- `context/AGENTS.md`, `skills/find-skills/SKILL.md`, `skills/VENDORED.md`, `agents/bugs.md`, `REQUIREMENTS.md`, `.agents/docs/testing/skill-matrix.md` — prose/doc changes; verified via the rendered-context checks (item 8), the `list` skill count (item 9), and the context-size test (item 4); no unit-test coverage expected or applicable for prose.
- Suppressions added: none found in the diff (`git diff` above shows no added `//nolint`, `t.Skip`, `@ts-ignore`, `.only(`, or similar directive lines).

## Not examined

- `internal/harness/repo.go` and `internal/harness/report.go` diffs were not individually mapped line-by-line to specific test functions beyond confirming the package's tests all pass; a targeted `-run` against the specific new logic in these two files was not performed (out of scope per the dispatch's fixed command list).
- The `zip.go`/`zip_test.go` new zip-related logic was verified by the full `-race` package pass and the new-lines-only diff (21/0 in `zip_test.go`) but not read line-by-line for semantic correctness of assertions.
- No `/review --docs` run — out of scope for this verification dispatch (that is a separate reviewer/documentation step per this repo's AGENTS.md).
- The `TestHarnessMap` test (pinning the README delivery map) was not separately isolated to confirm it ran as part of `go test -race ./...`; it is presumed covered by the full-suite pass but its specific PASS line was not individually grepped out.
- A parallel reviewer's file possibly added under `.agents/docs/reviews/` during this run was not checked (per dispatch instruction to ignore it).

## Environment

OS: Linux 6.12.107+deb13-amd64. Go: go1.27.1 (user-level SDK ahead of `/usr/bin/go`, per AGENTS.md). Repo state: `main @ e4ae34c` plus the working-tree edits named in the dispatch (11 modified files, 1 untracked testing doc). Scratch home used for install/status: `.agents/scratch/fu` (gitignored, never the real `$HOME`).
