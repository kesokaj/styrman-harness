---
owner: verifier
date: 2026-09-15
---
# Verification — harness practices G01-G26, Phase 2   (2026-09-15)

Range reviewed: cec4d4e..working tree (live tree; nothing from Phase 2 committed yet)

`git diff cec4d4e --stat`: `.agents/docs/plans/2026-09-03-rebuild.md` (12), `.agents/docs/plans/2026-09-13-learning-loop-project-level.md` (2), `agents/changelog.md` (2), `context/AGENTS.md` (6), `skills/bug/SKILL.md` (2), `skills/feature/SKILL.md` (33), `skills/handoff/SKILL.md` (6). `git status --short` shows no untracked files.

## Status: PASS

Scope: guide §1-3 run in full (§3a not applicable); plan boxes 2.1-2.5 and their positive controls run exactly as given. §4 sandbox CLI start deferred (see Not examined).

## Commands

### Guide §1 — Go checks
No `.go` file is in the range (`git diff cec4d4e --stat -- '*.go'` printed nothing). Ran anyway:
- `gofmt -l .` → no output, exit 0
- `go vet ./...` → no output, exit 0
- `go test -race ./...` → `ok github.com/kesokaj/styrman-harness/cmd/sh-agent 1.234s`, `ok github.com/kesokaj/styrman-harness/internal/harness 1.374s`, exit 0

### Guide §2 — Harness self-check
- `go run ./cmd/sh-agent list` → printed skills (27) and subagents (7) tables, no parse errors, exit 0
- `go run ./cmd/sh-agent render | wc -l` → `93`

### Guide §3 — Scratch-home install and drift
- `.agents/scratch/verify-phase2-home` did not exist before this run
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase2-home` → `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`, exit 0
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-phase2-home` → `drift none — 259 files match the repo`, exit 0
- `grep -c "^## Records" agents/changelog.md` → `1`
- `grep -c "^## Records" .agents/scratch/verify-phase2-home/.gemini/config/agents/changelog/agent.md` → `1`
- Both the source agent definition and its derived Antigravity agent.md carry exactly one `## Records` section.

### Guide §3a — Plain-copy install
`git diff cec4d4e --stat -- cmd internal` printed nothing → §3a does not apply, not run.

### Plan box 2.1 (I-7a-e) — expect: every pattern ≥1 in `skills/feature/SKILL.md`, last pattern 0
`grep -cF` for each pattern in `skills/feature/SKILL.md`:
"— expect:"=2, "- [ ] Checkpoint —"=1, "## Decisions"=3, "- [ ] Cold-reader check —"=1, "— in progress:"=1, "recon-<slug>"=1, "phase-<N>"=1, "close the plan"=1, "<docs-root>/bugs/INDEX.md"=1, "then commit, push and `/handoff`"=1, "over the range when an interface"=1; last: "then tick the boxes in the plan"=0. Matches expect exactly.
Positive control over c192622 (`git show c192622:skills/feature/SKILL.md` written to scratch): every pattern printed `0`, and the last pattern printed `1`. Matches the stated positive control exactly.

### Plan box 2.2 (I-8a-e) — expect: `1`, `0`, `1`, `93`
`go run ./cmd/sh-agent render` → `.agents/scratch/render.txt`; `grep -cF "restate status and next step before the next edit"`=1, `grep -cF "Workflow skills, also typed"`=0, `grep -cF "can also be typed as"`=1, `wc -l`=93. Matches expect exactly.

### Plan box 2.3 (I-9a-d, I-10a-b) — expect: `1`, and nothing printed
`grep -cF "Ask for six things" skills/handoff/SKILL.md` → `1`. `git grep -nF 'should put in \`TODO.md\`' -- agents skills context` → no output (exit 1). Matches expect exactly.
Positive control over c192622: `agents/changelog.md:94:Reply in **at most ten lines**...should put in \`TODO.md\`...` — matches the stated control line exactly.

### Plan box 2.4 (I-11a-g) — expect: only this plan's path
Fence-aware awk (I-9d command) over `.agents/docs/plans/*.md` → printed only `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`. Matches expect.
Positive control: `git show cec4d4e:<path>` for the three plan files (`2026-09-15-harness-practices-g01-g26.md`, `2026-09-03-rebuild.md`, `2026-09-13-learning-loop-project-level.md`) written into `.agents/scratch/verify-phase2-control-plans/`, then the same awk command run over that directory's `*.md` → printed all three files (the two old plans plus this plan's own path, which always has unticked boxes). Matches "before I-11 it also prints the two old plans."

### Plan box 2.5 (Describing text) — expect: nothing printed
`git grep -nE "tick its tasks|then tick the boxes|Workflow skills, also typed|should put in .TODO" -- context skills agents README.md AGENTS.md` → no output (exit 1). Matches expect.
Positive control over c192622: printed 4 lines (`agents/changelog.md:94`, `context/AGENTS.md:11`, `context/AGENTS.md:71`, `skills/feature/SKILL.md:72`) — at least three lines as stated. Matches.

## Root cause (FAIL only)
N/A — all checks passed.

## Coverage of changed code
No `.go` files changed in this range; the change set is entirely context/skills/agents/plan markdown, exercised by the grep/awk verify commands above rather than `go test`.

## Not examined
- Guide §4 (sandbox start of each CLI, log grepped for parse errors) — deferred to the Phase 5 checkpoint per the plan, not run in this phase.
- Guide §5 (real-home smoke checks) — not applicable; no install to a real home occurred.
- Guide §3a (plain-copy install) — not applicable, no `cmd`/`internal` changes in range.
- Phase 2 Checkpoint line's own further items (e.g. re-running `go test -race ./...` as part of the checkpoint evidence) are covered by the §1 run above; no separate second run performed.

## Environment
Linux (Debian 6.12.107+deb13-amd64), Go per `go.mod` in this checkout, no Docker/Colima involved in this phase.
