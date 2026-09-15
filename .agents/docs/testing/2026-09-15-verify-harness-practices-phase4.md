---
owner: verifier
date: 2026-09-15
---
# Verification — harness practices G01-G26, Phase 4 (2026-09-15)

Range reviewed: f3013a6..working tree (nothing from Phase 4 committed; `git diff f3013a6 --stat` plus one untracked file, `.agents/docs/changes/2026-09-15-harness-practices-phase4.md`)

## Status: PASS

Scope: complete — gofmt/vet/test chain, `sh-agent list`/`render`/`install`/`status`, all five plan boxes (4.1-4.5) re-run against the live tree and a fresh scratch home, `## Records` check on all five changed agents in both vendor copies, suppression scan over the full diff. Sandbox CLI starts deferred to Phase 5 per instructions (not this run's scope).

## Commands

- `gofmt -l .` → exit 0, no output (clean).
- `go vet ./...` → exit 0, no output.
- `go test -race -v ./...` → exit 0, saved to `.agents/scratch/verify-phase4/go-test-1.txt`.
  counts: pass 31, fail 0, skip 1 at f3013a6+dirty for `go test -race -v ./...`
  vs previous: none (no earlier `testing/*.md` log carries a `counts:` line; grep confirmed).
- `go run ./cmd/sh-agent list` → exit 0, prints 27 skills / 7 agents, no parse errors.
- `go run ./cmd/sh-agent render | wc -l` → `93` (matches expect).
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase4-home` → directory did not exist beforehand (confirmed); exit 0, "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed".
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-phase4-home` → exit 0, "drift  none — 259 files match the repo" (matches expect "drift none").
- `git diff f3013a6 --stat -- cmd internal` → empty. §3a of `.agents/docs/testing/sh-agent/guide.md` does not apply — no `cmd`/`internal` changes in range.

### Plan box 4.1 (I-16a–d) — `agents/verifier.md`
- "Every command on its own" → 1
- "Test counts:" → 1
- "Inherited failures:" → 1
- "Suppressions added:" → 1
- "What this run did not run or read" → 1
All match expect (each `1`).

### Plan box 4.2 (I-17a–f)
- "counts as not reached until you have read the rest" → reviewer.md 1, verifier.md 1, security.md 1, explorer.md 1
- "What you did not reach — files" in `agents/reviewer.md` → 1
- "What you did not reach — paths" in `agents/security.md` → 1
All match expect.

### Plan box 4.3 (I-18a–n)
- "Diagnosis missed:" → `agents/bugs.md` 1, `skills/bug/SKILL.md` 1, `skills/retro/SKILL.md` 1 (expect: at least 1 each)
- "Graduated" → `agents/bugs.md` 3, `skills/retro/SKILL.md` 1, `context/AGENTS.md` 1 (expect: at least 1 each)
- "in the Graduated table once its prevention rule lives" → `.agents/docs/bugs/INDEX.md` 1, `skills/init-project/assets/agents-tree/docs/bugs/INDEX.md` 1 (expect 1 each)
- `render | wc -l` → 93 (expect 93)
All match expect.

### Plan box 4.4 (ADR-021)
- `grep -c '^## ADR-021: Bug index rows move' .agents/docs/architecture/DECISIONS.md` → 1 (expect 1).

### Plan box 4.5 (derived Antigravity agents), run against fresh scratch home `.agents/scratch/verify-phase4-home`
- "counts as not reached" in derived `agent.md`: reviewer 1, verifier 1, security 1, explorer 1 (expect 1 each)
- "Every command on its own" in derived verifier `agent.md` → 1 (expect 1)
All match expect.

### `## Records` section count — changed agents, source vs derived (Antigravity)
- bugs: src=1, derived=1
- explorer: src=1, derived=1
- reviewer: src=1, derived=1
- security: src=1, derived=1
- verifier: src=1, derived=1
All five changed agents carry exactly one `## Records` section in both the Claude source and the Antigravity-derived copy.

### Suppressions added (from `git diff f3013a6`)
Scanned the full diff for added `//nolint`, `//lint:ignore`, `@ts-ignore`, `@ts-expect-error`, `eslint-disable`, `t.Skip`, `t.SkipNow`, `//go:build ignore`, `.skip(`, `.only(`. The only hit is a line of prose in `agents/verifier.md` (line 197 of the diff) that describes the suppression rule itself, not an actual added suppression. No suppressions or skips were added in this range.

## Coverage of changed code
No Go source under `cmd/` or `internal/` changed in this range (`git diff f3013a6 --stat -- cmd internal` is empty), so the Go test suite's pass/fail counts are the whole coverage story for code and unaffected by Phase 4. All Phase 4 changes are to `agents/*.md`, `context/AGENTS.md`, `skills/*/SKILL.md`, and `.agents/docs/*` — text checked above with `grep -c` against each plan box's exact expectation, all passing.

## Root cause (FAIL only)
N/A — Status PASS.

## Not examined
- The sandbox start of each CLI (`claude -p` / `agy -p` smoke check) — deferred to the plan's Phase 5 checkpoint per instructions, not run here.
- No frontend/container tiers apply (no `web/`/`ui/` package.json, no Dockerfile changed).
- Did not read the untracked `.agents/docs/changes/2026-09-15-harness-practices-phase4.md` file's contents — noted only as an untracked file in range, not part of the verify-command scope given.

## Environment
- OS: Linux 6.12.107+deb13-amd64
- Go: `go version` not separately captured this run, but `go test`/`go vet`/`gofmt` all ran successfully via the toolchain on PATH.
- Docker/Colima: not applicable (no Dockerfile changed).
