---
owner: changelog
date: 2026-09-15
---
# Harness practices Phase 2: plans, phases and memory (G01, G02, G03, G04, G05, G09, G12, G13, G14, G19)     (2026-09-15)

## Summary

Phase 2 of `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`: edits I-7a–e, I-8a–e, I-9a–d, I-10a–b and I-11a–g, applied verbatim from the plan's old/new blocks. A plan now carries an expect on every task, a Checkpoint box per phase, a Cold-reader check box and a Decisions section, and a Checkpoint box takes its evidence from a log or a call. The global context keeps state in the trail and restates status after a compaction. The handoff reads the active plan's state and checks for more drift. Noticed and deferred work goes to the plan's Follow-ups or the handoff's Next Steps, not `TODO.md`. The two stale plans are closed.

## Type

task

## Files changed

- `skills/feature/SKILL.md` — I-7a–e.
  - Step 1 reads `<docs-root>/bugs/INDEX.md` and says so when `ROADMAP.md` "Now" already names a plan with unticked boxes. A single request with gaps gets the numbered form with a recommended answer per gap; a gap left unsettled becomes an `Open:` row. An unread area or an estimated 10+ files gets an `explorer` recon brief at `plans/<date>-recon-<slug>.md`.
  - Step 2's template: each E2E criterion is a yes/no outcome; a Cold-reader check box; `— expect:` on every task; a Checkpoint box per phase; an Assumptions line; a `## Decisions` section with D, `Open:` and `Rejected:` rows; a `— docs:` clause on a task that changes a command, flag, endpoint, config key or deploy step; `[-]` for a box closed without being done.
  - Step 3: the cold-reader brief also asks for criteria that name an implementation choice, unchecked assumptions and missing docs clauses; rebuttals and non-blocking questions become Decisions rows; the Cold-reader box is ticked only on both tokens or a quoted same-turn override.
  - Step 4: a task whose verify names a test writes that test first and sees it fail; a box is ticked when its result meets its expect; a stopping point inside a phase leaves an `— in progress:` note; work found mid-phase is filed where it belongs; in a plan with two or more phases, a phase that changed source gets a `verifier` dispatch with slug `<plan-slug>-phase-<N>`, and its Checkpoint box copies evidence from the log or the call.
  - Step 6 adds `/review --docs` when an interface, command or deploy step changed, closes the plan before the commit, and ends with commit, push and `/handoff`.
- `context/AGENTS.md` — I-8a–e: a phase checkpoint ticks its Checkpoint box; progress that lives only in chat, a compaction summary or the CLI's own memory is lost, so it goes to the trail file that owns it; after a compaction, re-read the plan and `git status` and restate status and next step before the next edit; a Task that meets a design question or outgrows one sitting becomes a Feature for the remainder; the "Workflow skills, also typed as" line is removed. The rendered context stays at 93 lines.
- `skills/handoff/SKILL.md` — I-9a–d: "What's still open" starts at the active plan's first unticked box, then its Follow-ups and `Open:` rows; "Key decisions" starts at the plan's Decisions rows; "Ask for four things" becomes "Ask for six things"; next to contradictions between handoff, ROADMAP, TODO and plan, it now names a plan with an unticked box outside its fenced blocks that `ROADMAP.md` "Now" does not name, and a commit since the previous handoff that changed paths outside the docs root with no `changes/` entry.
- `skills/bug/SKILL.md` — I-10a: adjacent work noticed while fixing is noted under the active plan's Follow-ups, or the handoff's Next Steps when no plan is active.
- `agents/changelog.md` — I-10b: the reply names deferred work to file under the active plan's Follow-ups or the handoff's Next Steps, not `TODO.md`.
- `.agents/docs/plans/2026-09-03-rebuild.md` — I-11a–f: the Mac Antigravity box is ticked as verified on 2026-09-04; five operator boxes are closed as `[-]`, superseded by this plan's Done list.
- `.agents/docs/plans/2026-09-13-learning-loop-project-level.md` — I-11g: its one open box is closed as `[-]`, superseded by this plan's Done list.
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — boxes 2.1–2.5 and the Phase 2 checkpoint ticked with evidence.
- `.agents/docs/testing/2026-09-15-verify-harness-practices-phase2.md` — the Phase 2 verifier log.

## Agents dispatched

verifier. Its first run stopped at its 30-turn limit without writing a log; resumed, it wrote the log.

Review: no `reviewer` ran over the Phase 2 diff. The plan's Done list runs `reviewer` and a `--docs` review once over `c192622..HEAD`.

## Verification

The verifier log reads Status: PASS. The main agent also ran, each bare with its exit status read: `go test -race ./...` exit 0; `gofmt -l .` printed nothing; `go vet ./...` exit 0; a scratch install into `.agents/scratch/phase2-checkpoint-home` exit 0, then `status` with drift none; `render` 93 lines; the changelog agent's `## Records` section once in each vendor's installed copy. The checks of boxes 2.1–2.5 met their expects, with positive controls over `c192622`; box 2.4's control ran before I-11 was applied.

Not run: the sandbox start of each CLI (testing guide §4), deferred to the plan's Phase 5 checkpoint (D36).

## Docs

The changed living docs are the skill, agent and context texts above. `README.md`, the testing guide and the deploy guides are unchanged in this phase.

## Requirements

n/a
