---
owner: verifier
date: 2026-09-21
---
# Verification — project-agents-md-native Phase 1   (2026-09-21)

Range reviewed: v1.5.0..b013b52

## Status: PASS

Scope: guide.md § 1 (Go checks), § 2 (harness self-check, size test), § 3 (scratch-home
install and drift), plus the plan-specified spot checks (asset removal, global stub
retained, repo-root CLAUDE.md gone, SKILL.md wording). testing/guide.md exists and was
followed as written, not rewritten (no step change needed). Behavioural/e2e/CLI-start
evidence for this range lives in `.agents/docs/testing/2026-09-21-verify-agents-md-native.md`
and the plan's Phase 1 task-box notes — cited, not rerun, per dispatch instruction. No
agy or claude CLI invocation was run by this dispatch; the real home was never touched;
no source was edited; nothing was deleted outside `.agents/scratch/verify-v160`.

Scope: complete for the tiers this dispatch was asked to run.

## Commands

- `gofmt -l cmd internal` → exit 0, no output (no files need formatting). PASS
- `go vet ./...` → exit 0, no output. PASS
- `go test -race -count=1 ./...` → exit 0; all packages `ok` (root, cmd/sh-agent,
  internal/harness). Output: `.agents/scratch/verify-v160/go-test-race.txt`. PASS
  - Also ran `-v` variant for counts (guide § 1 does not itself specify `-v`, added per
    the general test-counts rule): `.agents/scratch/verify-v160/go-test-race-v.txt` →
    exit 0. counts: pass 251, fail 0, skip 0 at b013b52 for
    `go test -race -count=1 ./...`.
    vs previous: `.agents/docs/testing/2026-09-21-verify-preflight-v1.5.0.md`, which
    ran the same command per guide § 1 as written (no `-v`) and recorded
    `counts: not printed`; no prior `-v` run of this exact command found in `testing/`
    to diff pass/fail counts against. Not a regression signal either way — noted as a
    baseline for future `-v` runs of this command.
- `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1 -v`
  → exit 0. `size_test.go:48: Antigravity rules block is 20328 bytes (context 20328,
  rules 0), budget 22000` — `--- PASS`. Matches the dispatch's expected value (20328 of
  22000) exactly, and matches the plan's Phase 1 task-box note (2026-09-21: "size test
  PASS, 20,328 of 22,000"). PASS
- `go run ./cmd/sh-agent list` → exit 0. Output shows "### Skills (27)" with 27 table
  rows and "### Subagents (7)" with 7 table rows. Matches expected 27 skills / 7 agents.
  PASS
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-v160/home` → exit 0.
  PASS
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-v160/home` → exit 0;
  output line `drift      none — 257 files match the repo`. PASS
- `test ! -e .agents/scratch/verify-v160/home/.claude/skills/init-project/assets/CLAUDE.md`
  → exit 0 (path absent — the asset is gone from the install, as the plan intends).
  PASS
- `tail -1 .agents/scratch/verify-v160/home/.claude/AGENTS.md` → `Canary: this context
  ends with the word GANNET.` — the expected canary line, install rendered correctly.
  PASS
- `test -f .agents/scratch/verify-v160/home/.claude/CLAUDE.md` → exit 0 (file present).
  Contents: a generated-by comment plus `@AGENTS.md` — the global stub is still
  written, per plan scope (global stub explicitly Out of this plan). PASS
- `test ! -e CLAUDE.md` (repo root) → exit 0 (absent). PASS — matches Phase 1 task box
  "`git rm CLAUDE.md` at the repo root" and `git diff --stat` below confirming its
  deletion in-range.
- `grep -c "not written" skills/init-project/SKILL.md` → `1`, exit 0. Matches expected
  count of 1. PASS

## Root cause

N/A — no failures.

## Coverage of changed code

`git diff v1.5.0..HEAD --stat` (14 files changed, 306 insertions, 8 deletions):
docs/plan/review/ADR/handoff/testing-log files, `CLAUDE.md` (deleted),
`REQUIREMENTS.md`, `ROADMAP.md`, `context/AGENTS.md`, `skills/init-project/SKILL.md`,
`skills/init-project/assets/CLAUDE.md` (deleted). No `.go` file is in this range — no
Go source changed, so no test-to-source mapping applies for this tier; the Go suite
above (`go test -race ./...`) still ran in full as guide § 1 requires for every change,
and passed against the unmodified Go code. All changed files are documentation/markdown
or the removed `CLAUDE.md` assets, which this dispatch's spot checks (install output,
grep, `test -e`) directly exercise.

No suppressions (`//nolint`, `t.Skip`, `.only(`, etc.) found in this diff — it contains
no Go or test source.

## Not examined

- The behavioural/CLI-start/integration evidence for this range (sandbox Claude Code
  and Antigravity starts quoting the new Project layout line, the headless
  `/init-project` probe run, the e2e golden-path answer from this checkout's
  `AGENTS.md`) — per the dispatch, this is cited from
  `.agents/docs/testing/2026-09-21-verify-agents-md-native.md` and the plan's Phase 1
  task-box note, not rerun by this dispatch.
- `/review --docs` over the range (Phase 1 task box "Review row" is unchecked in the
  plan as of this run) — out of scope for this verifier dispatch, which only covers
  the Checkpoint's test tier.
- `govulncheck ./...` — not requested by this dispatch's command list; not run.
- guide § 3a (plain-copy install), § 3b (`report`), § 3c (verb × place-kind matrix),
  § 4 (sandbox CLI start), § 5 (real-home smoke checks) — none requested by this
  dispatch; § 4/§ 5 evidence is the cited prior testing log per the dispatch
  instructions.
- `govulncheck`, `agy`/`claude` CLI runs — explicitly excluded by the dispatch
  ("never run agy or claude").

## Environment

OS: Linux 6.12.107+deb13-amd64 (Debian trixie). Go: go1.27.1 linux/amd64 (user-level
SDK, per `AGENTS.md` Setup commands note). No Docker/Colima involved — no Dockerfile in
this range. Scratch home: `.agents/scratch/verify-v160/home` (gitignored, created and
left in place per dispatch instructions — not deleted, since dispatch said "delete
nothing outside `.agents/scratch/verify-v160`").
