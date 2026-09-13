---
owner: reviewer
date: 2026-09-13
---
# Review — learning-loop-phase-3   (2026-09-13)

Range reviewed: `git diff bd5ed87 -- skills/retro/SKILL.md skills/distill/SKILL.md skills/feedback/SKILL.md context/AGENTS.md skills/init-project/SKILL.md skills/handoff/SKILL.md`, plus the new `skills/init-project/assets/agents-tree/docs/retros/.gitkeep` (untracked). Judged against Decisions D5–D10, the "Retro skill structure" section, and the Phase 3 task list in `.agents/docs/plans/2026-09-13-learning-loop-project-level.md`.

## Verdict: APPROVE WITH NITS

Nothing here would fail a step outright on Linux, and the D5–D10 machinery (scope table, `rule_ref` vocabulary, promotion grep, record template, queue/index schema, apply flow) is implemented consistently across `retro`, `distill`, `feedback` and `context/AGENTS.md` — I did not find a place where one file's schema disagrees with another's. The two real problems are: a portability bug in the new timestamp instruction that fails on the operator's own stated macOS environment, and a stale line in `handoff/SKILL.md` (a file this same diff touched) that still describes `/retro` by its pre-Phase-3 behavior. A third item is a task the plan itself lists as done in this pass but that the diff does not show.

## Findings

- **[medium / confidence: high]** `skills/retro/SKILL.md:21`, `skills/retro/SKILL.md:122`, `skills/distill/SKILL.md:146`, `skills/distill/SKILL.md:152` — both skills now instruct "timestamps from `date -Iseconds`". `AGENTS.md` (global context, line 5) and `context/dev-setup.md` both name macOS (Apple Silicon) as a supported operator machine, and neither this diff nor the existing dev-setup doc installs GNU coreutils. BSD `date` (macOS default) has no `-I`/`-Iseconds` option and exits with `date: illegal option -- I`; a model executing this instruction literally on a Mac fails the step. Fix: use a form BSD and GNU `date` both accept, e.g. `date +%Y-%m-%dT%H:%M:%S%z` (matches the `ts` format already used in the `feedback.jsonl` example), or explicitly require `gdate` and say so.
- **[medium / confidence: high]** `skills/handoff/SKILL.md:192` — "The `/retro` workflow captures corrections and proposes harness improvements" is left over from before this diff. The same file's own artifact table, two lines above (50–51), was updated in this diff to reflect the new dual-scope model, and `retro/SKILL.md`'s rewritten description and D6 make project-scope the more common outcome (a finding about "THIS project's code, layout, commands, or conventions" targets the project, not the harness). The Cross-References line now contradicts the skill it references and the table directly above it in the same file.
- **[low / confidence: medium]** `skills/init-project/assets/AGENTS.md` — the plan's "Surfaces touched" list names `assets/AGENTS.md (Key Files row)` and the Phase 3 task line for `init-project` names three artifacts ("`assets/agents-tree/docs/retros/.gitkeep`, tree listing in SKILL.md, Key Files row"). The `.gitkeep` exists and the tree listing in `skills/init-project/SKILL.md` was updated, but `skills/init-project/assets/AGENTS.md`'s "Key Files" table (the scaffold shipped into new projects) has no diff and no row for `docs/retros/`. Two of three sub-items are done; the third named artifact is untouched.
- **[low / confidence: low]** Phase 3's own task-list line says "Steps 4–7" for `skills/retro/SKILL.md`, but the "Retro skill structure" section it points to (and the shipped file) only defines Steps 2 through 6 — there is no Step 7 anywhere in the plan or the skill. Not a functional bug in the skill, but the plan's own two statements of scope disagree, which will confuse whoever ticks that checkbox.

## Not examined

`internal/harness/*`, `cmd/sh-agent/*` (no Go changes in this diff — Phase 3 is prose-only, consistent with the plan). The sandbox e2e item in Phase 3 ("headless Claude Code `/retro` … produces the record and the index line") was not run; this review is a static read of the prose only. `context/dev-setup.md` was read only for the macOS claim, not diffed (unchanged in this range).

## Criteria

1. Dead code — pass. No orphaned sections; old Step 2–5 table in `retro/SKILL.md` was fully replaced, not left alongside the new one.
2. Flow tracing — mostly pass: `rule_ref` vocabulary, scope table, promotion-grep patterns, record/index schemas trace consistently retro → feedback → distill → context/AGENTS.md → init-project → handoff. Two breaks found: the `date -Iseconds` step is unreachable on macOS as written, and `handoff`'s Cross-References line was not updated to match the caller it documents.
3. Standards conformance — n/a (no Go/TS in this diff); frontmatter lints clean (`go run ./cmd/sh-agent list` shows `retro` parsing fine, `auto`-invocable).
4. Directory layout — pass for what was touched (`retros/` added to both the harness and the scaffold tree); the scaffold's Key Files row is the one gap (see finding above).
5. Simplicity — pass. Grep-based promotion check and plain markdown/JSONL records are the simplest mechanism for the stated requirement; no new dependency.
6. Decoupling — pass. Project-scope vs. harness-scope proposals are cleanly separated by target and by file (project record vs. queue vs. index), matching D5/D6's intent that these are independently editable artifacts.
