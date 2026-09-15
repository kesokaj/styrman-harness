---
owner: verifier
date: 2026-09-15
---
# Verification — harness-practices-phase5   (2026-09-15)

Range reviewed: b43dde1..working tree (live tree, nothing from Phase 5 committed yet) — `git diff b43dde1 --stat` plus one untracked file, `.agents/docs/testing/2026-09-15-verify-g26-precedence.md` (5.1's own artifact, not re-verified here).

## Status: PASS

Scope: ran the repo's full verification chain (gofmt, vet, race tests, sh-agent list/render/install/status) and the verify command of every checkable Phase 5 box (5.2, 5.3 minus its ADR-022 count, 5.4, 5.5, 5.6, 5.7); did not run the sandbox CLI starts the Phase 5 checkpoint also names (main agent's job) or 5.1's own e2e test (already logged and ticked).

## Commands

- `gofmt -l .` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -v ./...` → exit 0. Output saved to `.agents/scratch/verify-phase5/go-test-1.txt`.
  counts: pass 31, fail 0, skip 1 at b43dde1+dirty for `go test -race -v ./...`
  vs previous: `.agents/docs/testing/2026-09-15-verify-harness-practices-phase4.md`, pass 31→31, skip 1→1 (no fall in pass, no rise in skip)
- `go run ./cmd/sh-agent list` → exit 0; 27 skills, 7 agents printed, no parse errors.
- `go run ./cmd/sh-agent render` → exit 0; 93 lines (`.agents/scratch/verify-phase5/render.txt`) — matches expect.
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase5-home` → exit 0 (home did not exist beforehand); "260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed".
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-phase5-home` → exit 0; `drift      none — 259 files match the repo` — matches "drift none" expect (the 260 vs 259 count difference is the binary itself, excluded from the file-drift compare, not a discrepancy in the managed set).
- `git diff b43dde1 --stat -- cmd internal` → no output: no `cmd/` or `internal/` changes in this range, so §3a of `.agents/docs/testing/sh-agent/guide.md` (rebuild-before-install) does not apply here.

### Box 5.2 verify
- `grep -cF "never reuses a harness name" .agents/scratch/verify-phase5/render.txt` → `1`
- `grep -cF "show what a vendor loaded" .agents/scratch/verify-phase5/render.txt` → `1`
- `grep -cF "never a harness" skills/retro/SKILL.md skills/distill/SKILL.md` → `skills/retro/SKILL.md:2`, `skills/distill/SKILL.md:1`
- Expect `1` and `1`; `2` and `1` — all match.

### Box 5.3 verify (ADR-022 count excluded — another agent is appending it now; see Not examined)
- `git grep -nF "50+ lines or 5+ files" -- context skills agents` → exactly two lines, both `skills/preflight/SKILL.md` (line 40, the Step 1 "or is a bug fix" line; line 50, the Step 2 release-diff line) — matches expect.
- `grep -cF 'over the fix when it changed an interface' skills/bug/SKILL.md` → `1` — matches expect.

### Box 5.4 verify
- `skills/preflight/SKILL.md`: "Blocking means any of" 1, "exactly as `/review` Step 3 gives it" 1, "fix those in skills, agents, rules" 1, "Step 4's docs fixes, the changelog entry" 1, "plus its records and Step 4's docs fixes" 1.
- `agents/documentation.md`: "or a deploy path changed" 1, "same deploy path" 1, "deploy guides, the project AGENTS.md setup and deploy sections" 1, "Setup commands and Deploy sections only" 1.
- All nine patterns print `1` — matches expect.

### Box 5.5 verify
- `grep -c "^## Gates that need the user" skills/feature/SKILL.md skills/bug/SKILL.md` → `1`, `1`.
- `grep -cF "The three-round cap is a gate" skills/review/SKILL.md` → `1`.
- `grep -cF "whatever autonomous instruction is in effect" skills/bug/SKILL.md` → `0`.
- `grep -cF "Never invent a repro on your own" skills/bug/SKILL.md` → `1`.
- `grep -cF "Two points always wait for the user" skills/bug/SKILL.md` → `1`.
- Expect `1` each; `1`; `0`; `1`; `1` — all match.

### Box 5.6 verify
- `grep -cF "derive the build, run and test commands" skills/init-project/SKILL.md` → `1`.
- `grep -cF 'add \`@AGENTS.md\` as its first line' skills/init-project/SKILL.md` → `1`.
- Expect `1` each — match.

### Box 5.7 verify
- `grep -cF "or fixes a bug" .agents/scratch/verify-phase5-home/.claude/agents/reviewer.md .agents/scratch/verify-phase5-home/.gemini/config/agents/reviewer/agent.md` → `1`, `1`.
- `grep -cF "or a deploy path changed" .agents/scratch/verify-phase5-home/.claude/agents/documentation.md .agents/scratch/verify-phase5-home/.gemini/config/agents/documentation/agent.md` → `1`, `1`.
- Expect `1` for each file — match.

### `## Records` count (extra check requested by the dispatch)
- `grep -c "^## Records" agents/reviewer.md agents/documentation.md` → `1`, `1`.
- Same grep against the scratch-home derived files `.gemini/config/agents/reviewer/agent.md` and `.gemini/config/agents/documentation/agent.md` → `1`, `1` — the `## Records` header survives the Claude → Antigravity derivation in both agents.

## Root cause (FAIL only)

Not applicable — Status is PASS.

## Coverage of changed code

This range touches only markdown (skills, agents, context, this repo's own docs) — no `cmd/` or `internal/` Go source changed (`git diff b43dde1 --stat -- cmd internal` was empty). Coverage is judged by the plan's own grep/render checks (boxes 5.2–5.7 above), all of which match their expect. No added `//nolint`, `t.Skip`, `.only(`, `.skip(`, `@ts-ignore`, `@ts-expect-error` or `eslint-disable` in `git diff b43dde1 -- '*.md'` — checked directly, no hits.

## Not examined

- Box 5.3's `grep -c '^## ADR-022: Every bug fix gets the reviewer' .agents/docs/architecture/DECISIONS.md` — explicitly excluded by the dispatch because another agent is appending that ADR concurrently; re-check once that write lands.
- Box 5.1 (the I-19 precedence test) and its own log `.agents/docs/testing/2026-09-15-verify-g26-precedence.md` — already ticked with evidence in the plan; not re-verified here.
- The Phase 5 Checkpoint's sandbox-install-and-CLI-start requirement ("this tree installed into the sandbox home and one start of each CLI quoting Definition of Done step 2's reviewer trigger... `grep -c "Failed to parse"` on the newest Antigravity `cli-*.log` printing `0`") — run by the main agent per the dispatch, not this log.
- Phase 6 and the Done section of the plan — out of scope for this dispatch.
- Content review of the untracked `.agents/docs/testing/2026-09-15-verify-g26-precedence.md` beyond confirming it exists as untracked — not read for correctness.
- No frontend, container, or non-Go tiers apply to this range; none run.

## Environment

Linux (6.12.107+deb13-amd64), repo at `b43dde1+dirty`. `go` and `git` present on PATH; scratch install used `.agents/scratch/verify-phase5-home` (did not exist beforehand, per the plan's constraint), never `/tmp` or the real `$HOME`.
