---
owner: main
date: 2026-09-03
---
# Global context revision: work protocol, checkpoints, bugs memory, scratch, cadence     (date: 2026-09-03)

## Summary

Third revision of `context/AGENTS.md` after the operator's review, plus the research they asked for.

- Operator line is "fullstack software engineer". Resume protocol now reads the active plan and the bug index, and has a **Checkpoints** subsection saying what is written and when.
- Definition of Done is an ordered sequence: tests and verification → review → docs → records → commit and push.
- `REQUIREMENTS.md` stays immutable; a trajectory change goes proposal → approval → ADR → update of REQUIREMENTS and ROADMAP.
- New **`bugs` subagent** (seventh): every failure becomes `BUG-NNN` under `.agents/docs/bugs/` with a prevention rule, indexed in `INDEX.md`, read at resume. This repo records BUG-001 and BUG-002 from the rebuild.
- Layout gains `docs/bugs/` and **`.agents/scratch/`** (gitignored) — the only place for generated or temporary files. Hard rule added; init-project assets and gitignore lines updated.
- **Global vs project customizations** section: both layers load, the project layer has priority (verified: Claude Code docs — additive, project rules after user rules; Antigravity — workspace above global). Project `.claude/` symlinks are documented for rules too.
- **Learning loop cadence:** retro after two or more corrections or after 14 days; `sh-agent status` now prints the days since the last retro proposal; distill on the third repeat. Harness state (feedback log, proposals, global handoffs) moved to `~/.config/sh-agent/`; the four skills updated.
- **Work protocol:** trivial / task / feature / bug paths; plan files with phases and task checkboxes; a cold-reader check of the plan by `reviewer` before any code; phases as checkpoints; goldfish diagnosis and a failing test for bugs; four testing tiers; e2e before "done"; four user gates.
- Research read in full: Böckeler's "Harness engineering for coding agent users" and `vshvedov/elephant-goldfish`. Adopted and rejected strategies in `.agents/docs/plans/2026-09-03-research-harness-strategies.md`; ADR-008.

## Type

task

## Files changed

- `context/AGENTS.md` — rewritten (200 lines deployed with dev-setup and inventory)
- `agents/bugs.md` — new
- `skills/init-project/SKILL.md`, `skills/init-project/assets/agents-tree/{docs/bugs/INDEX.md,scratch/.gitkeep}` — bugs and scratch in the scaffold
- `skills/{feedback,retro,distill,handoff}/SKILL.md` — state paths under `~/.config/sh-agent/`
- `internal/harness/state.go` (`LastRetro`), `cmd/sh-agent/main.go` (status line), `harness_test.go`
- `.agents/docs/bugs/{INDEX.md,BUG-001-status-stale-binary.md,BUG-002-adopt-flag-order.md}`, `.agents/scratch/.gitkeep`, `.gitignore`
- `README.md`, `AGENTS.md`, `REQUIREMENTS.md` (R7), `.agents/docs/plans/2026-09-03-rebuild.md` (section 10), `.agents/docs/architecture/DECISIONS.md` (ADR-008)
- `.agents/docs/plans/2026-09-03-research-harness-strategies.md` — new

## Verification

`gofmt`, `go vet`, `go test -race` clean (`TestLastRetro` added). End-to-end on a scratch home: 250 files, 7 agents on both vendors, status shows the retro line and no drift, uninstall clean. See `testing/2026-09-03-verify-context-revision.md`.

## Requirements

R2, R5, R7 (new).
