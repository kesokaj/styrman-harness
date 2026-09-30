---
owner: verifier
date: 2026-09-30
---
# Verification — retro-apply   (2026-09-30)

Range reviewed: working tree — changed files: AGENTS.md agents/verifier.md skills/bug/SKILL.md skills/feature/SKILL.md skills/golang/SKILL.md skills/review-standards/SKILL.md; untracked, not ignored: none outside .agents/docs; read the diff with git diff HEAD -- <files>.

## Status: PASS

Scope: complete for this change's tiers — prose-only edit (one agent definition, four skills, project AGENTS.md); no Go source, test, or config file changed. Ran guide § 1 (Go checks), the harness self-check subset the dispatch named, and guide § 3 (scratch-home install + drift). `testing/guide.md` already exists (dated 2026-09-17); no update needed since the steps did not change. Sandbox (§4), real-home smoke (§5), and Antigravity runs are out of scope for this dispatch and are listed under Not examined, backed by the two cited prior artifacts.

Fact recorded per dispatch: `.agents/docs` files appeared during this run from parallel agents (a `changelog` entry under `changes/`, two `reviews/` files, and an edit to `.agents/docs/bugs/INDEX.md`) — none of it is source, tests, or config, and `git status --short` at the end of this run confirms the tracked-file change set is unchanged from the dispatch's list plus that one docs edit.

## Commands
- `git diff HEAD -- AGENTS.md agents/verifier.md skills/bug/SKILL.md skills/feature/SKILL.md skills/golang/SKILL.md skills/review-standards/SKILL.md | wc -l` → 79 lines changed; confirms prose-only, no `.go` file in the diff.
- `gofmt -l cmd internal` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race ./...` (run once in background per the long-suite rule, output to `.agents/scratch/verify-retro-apply/go-test-race.txt`) → exit 0.
  ```
  ok  	github.com/kesokaj/styrman-harness	1.114s
  ok  	github.com/kesokaj/styrman-harness/cmd/sh-agent	14.735s
  ok  	github.com/kesokaj/styrman-harness/internal/harness	3.328s
  ok  	github.com/kesokaj/styrman-harness/internal/search	(cached)
  ```
  counts: not printed per-test (package-level `ok` lines only, no `-v`); counts: not applicable for this command's output shape. `internal/search` reported `(cached)` — untouched by this change, consistent with no source edits.
  vs previous: no earlier counts line for this exact command found in `testing/`; vs previous: none.
- `go test -count=1 -run 'TestRenderedContextFitsAntigravityBlock|TestLint|TestHarnessMap' -v ./internal/harness` → exit 0. Only `TestLintRejectsUnquotedComment` and `TestRenderedContextFitsAntigravityBlock` matched and ran (both PASS); `TestHarnessMap` does not live in `internal/harness` (it is defined in `cmd/sh-agent/main_test.go`), so the dispatch's exact command did not exercise it there — see gap closed below.
  - `TestRenderedContextFitsAntigravityBlock`: "Antigravity rules block is 21996 bytes (context 21996, rules 0), budget 22000" — PASS, under the ADR-035 budget.
  - counts: pass 2, fail 0, skip 0 at 5dea7c4+dirty for `go test -count=1 -run 'TestRenderedContextFitsAntigravityBlock|TestLint|TestHarnessMap' -v ./internal/harness`; vs previous: none found.
- Gap-closing command (not in the dispatch's literal list, run because the named regex missed its target package): `go test -count=1 -run 'TestHarnessMap' -v ./cmd/sh-agent` → exit 0, `--- PASS: TestHarnessMap`. counts: pass 1, fail 0, skip 0 at 5dea7c4+dirty for `go test -count=1 -run 'TestHarnessMap' -v ./cmd/sh-agent`; vs previous: none found.
- `go run ./cmd/sh-agent install --home .agents/scratch/home` (scratch home rebuilt fresh with `rm -rf .agents/scratch/home` first) → exit 0. "266 written, 0 unchanged, 0 removed, 0 kept — 266 files managed". Settings note printed (attribution trailers off) — informational, not a failure.
- `go run ./cmd/sh-agent status --home .agents/scratch/home` → exit 0. Key lines:
  ```
  repo       /home/kesokaj/workspace/styrman-harness (main @ 5dea7c4, uncommitted changes)
  installed  2026-09-30T10:24:25+02:00 from commit 5dea7c4+dirty (vendor: all, 266 files, manifest ~/.config/sh-agent/state.json)
  binary     5dea7c455bf4+dirty
  repo has   28 skills, 8 agents, 0 commands, 0 rules
  drift      none — 265 files match the repo
  ```
  Drift reports "none". Note: install reported 266 files managed while status compares 265 (one file, e.g. the manifest itself, is expected to be excluded from the drift comparison) — not a discrepancy in outcome, both report zero drift/zero unchanged-mismatch.
- `go run ./cmd/sh-agent list` → exit 0. Printed 28 skills and 8 subagents with descriptions; the four changed skills (bug, feature, golang, review-standards) and the changed agent (verifier) all listed and parsed without error.

## Installed-copy vs working-tree diff check
- `diff .agents/scratch/home/.claude/agents/verifier.md agents/verifier.md` → no output (match).
- `diff .agents/scratch/home/.claude/skills/bug/SKILL.md skills/bug/SKILL.md` → no output (match).
- `diff .agents/scratch/home/.claude/skills/feature/SKILL.md skills/feature/SKILL.md` → no output (match).
- `diff .agents/scratch/home/.claude/skills/golang/SKILL.md skills/golang/SKILL.md` → no output (match).
- `diff .agents/scratch/home/.claude/skills/review-standards/SKILL.md skills/review-standards/SKILL.md` → no output (match).

All five installed copies are byte-identical to the working tree.

Note: the changed project-root `AGENTS.md` is this repo's own rules file (per this repo's AGENTS.md § Project layout, "Root: AGENTS.md alone"), not `context/AGENTS.md` (the deployed global context) — it is not installed by `sh-agent install` into a target home, so no scratch-home diff applies to it. Confirmed only `AGENTS.md` changed in this commit's diff, not `context/AGENTS.md`, via `git diff HEAD --stat -- AGENTS.md context/AGENTS.md` (1 file changed).

## Coverage of changed code
Not applicable — no `.go` file is in the diff; the change is prose across one agent definition, four `SKILL.md` files, and the project `AGENTS.md`. No suppressions (`//nolint`, `t.Skip`, `.only(`, etc.) apply since no source or test file changed.

## Cited prior artifacts (not rerun, per dispatch instruction)
- `.agents/docs/testing/2026-09-30-verifier-foreign-artifact-behavioural.md` — behavioural run of the `agents/verifier.md` change.
- `.agents/scratch/vkit/runs/skillcheck.json` — headless CLI start exercising the working-tree skills; recorded subtype `success`, empty stderr.

## Not examined
- Guide § 4 (sandbox start of each CLI) — not run this session; not required by the dispatch, and the change is prose-only with the behavioural/headless runs above already covering the agent and skills text.
- Guide § 5 (real-home smoke checks, including the GANNET canary and the Antigravity conversation-store check) — not run; no install into a real `$HOME` occurred or is planned in this dispatch (installer was only run against `.agents/scratch/home`).
- Antigravity agent.md derivation for `verifier` was written by the install step but not separately diffed/inspected beyond the installer's own "0 unchanged, 0 removed" success signal; not opened by hand.
- Guide § 2's `go run ./cmd/sh-agent render | wc -l` line-count check (under 100 lines) — not run; the dispatch's explicit command list did not include it and `TestRenderedContextFitsAntigravityBlock` (the byte-budget pin, run above) is the guide's stated authoritative check for `context/*.md`/`rules/*.md` changes, which this change set does not touch anyway (only `AGENTS.md`, `agents/verifier.md`, and `skills/*/SKILL.md` changed — no `context/` or `rules/` files).
- The full skill-matrix (`testing/skill-matrix.md`) — not run; this dispatch's scope is the retro-apply prose change, not a "complete"/"all skills" claim.
- Sections 3a/3b/3c/3d of the guide (plain-copy install, `report`, verb×place matrix, search index) — not applicable, no code touching those areas changed.

## Environment
Linux 6.12.107+deb13-amd64, repo at /home/kesokaj/workspace/styrman-harness, HEAD 5dea7c4+dirty on `main`. `go` toolchain as configured in this environment (no explicit version pin per `go.mod`, per project convention). Scratch artifacts: `.agents/scratch/verify-retro-apply/go-test-race.txt`, `.agents/scratch/home/` (fresh scratch home, gitignored).
</content>
