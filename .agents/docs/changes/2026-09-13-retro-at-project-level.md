---
owner: changelog
date: 2026-09-13
---
# Phase 3 — learning loop at project level     (2026-09-13)

## Summary
Completed Phase 3 of the learning-loop plan: `skills/retro/SKILL.md` rewritten with target-first scope decision (two-column per-category target table, `rule_ref` vocabulary, two spacing-agnostic promotion greps); `skills/distill/SKILL.md` Step 6 records project-level `.agents/docs/retros/<date>-distill-<name>.md` and index line; `skills/feedback/SKILL.md` adopts `rule_ref` vocabulary forward; `context/AGENTS.md` gains Learning loop paragraph and Project layout `retros/` row; `skills/init-project/assets/agents-tree/docs/retros/.gitkeep` added; `skills/handoff/SKILL.md` artifact table cross-references retro records. Project-scope proposals written in full to `<docs-root>/retros/<date>-retro.md`; harness-scope proposals recorded as pointers (full diff stays in `~/.config/sh-agent/proposals/`, per R7); one index line appended to `~/.config/sh-agent/retros.jsonl`; `--harness` flag and headless behaviour documented; timestamps use `date -u +%FT%TZ` (macOS omits nanoseconds and UTC label in `-Iseconds`). Phase 3 residue from round-3 review resolved: the `rule_ref` vocabulary and the two spacing-agnostic promotion greps both approved "use your recommendation" (operator's turn 2026-09-13, after round-3 findings).

## Type
feature (phase 3 of 5)

## Files changed
- `skills/retro/SKILL.md` — Steps 2–6 rewritten; target table (project/harness per category); `rule_ref` vocabulary and promotion rules; project record template; `--harness` consumption; permission-denied path
- `skills/distill/SKILL.md` — Step 6 records project record and index line
- `skills/feedback/SKILL.md` — `rule_ref` field adopts controlled vocabulary; Log Format example; cross-reference to retro outputs
- `context/AGENTS.md` — Learning loop paragraph (scope rule, index, cross-machine hand-over); Project layout gains `retros/` row
- `skills/init-project/assets/agents-tree/docs/retros/.gitkeep` — new
- `skills/init-project/SKILL.md` — tree listing updated
- `skills/handoff/SKILL.md` — artifact table row for retro records

## Agents dispatched
reviewer ×1 (Phase 3 prose), changelog (you)

## Verification
Frontmatter lint: `go run ./cmd/sh-agent list` (no parse errors); rendered context 92 lines; sandbox home with working tree installed: headless `/retro` with two inline corrections each in fresh clones of the demo on Claude Code 2.1.270 and Antigravity 1.1.25 → each wrote `.agents/docs/retros/2026-09-13-retro.md` (78 and 66 lines) with two `Scope: project` proposals targeting the project `AGENTS.md`, each appended one line to `~/.config/sh-agent/retros.jsonl` (verified JSON), no harness-scope proposal so no queue file, nothing applied (headless), Antigravity log `cli-20260913_214533.log` zero `Failed to parse` errors; the Claude run detected that the index grep required a slash and could not match a root-level `AGENTS.md` target — fixed in the skill and the plan (the directory part is now optional). Reviewer APPROVE WITH NITS: fixed macOS `date` command, stale handoff cross-reference, plan wording.

## Requirements
R7 holds as written (harness-scope content only under `~/.config/sh-agent/`; the project record is the project's own state — ADR-015 below).
