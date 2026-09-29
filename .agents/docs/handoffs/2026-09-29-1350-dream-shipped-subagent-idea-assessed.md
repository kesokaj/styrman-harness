---
owner: main
date: 2026-09-29
---
# Handoff — 2026-09-29 — dream skill shipped (not released); memory/dreaming-as-subagents assessed, awaiting the operator

Session ID: 8a02b8f1-6422-4650-9dfd-3ea55bba9b1b
Project: /home/kesokaj/workspace/styrman-harness
Branch: `main` at 22ab8e7, pushed; this handoff's commit follows it.

## Context

Continuation of `handoffs/2026-09-29-1254-prompts-lens-shipped-dream-plan-proposed.md`. The operator said "Go all on the dream plan", then "Continue with round 3 and build the skill", answered the cold-read cap with "Proceed, add the rows (Recommended)", the apply run with "Go all (Recommended)", the REQUIREMENTS clause with "Approve the edit (Recommended)", raised the subagent idea ("build dreaming and memory as subagents … the same for claude and agy"), and said "dont cut new version yet. Think about the subagent idea."

## Completed

- **`dream` skill** (`skills/dream/SKILL.md`; plan `plans/2026-09-29-dream-skill.md`, closed): consolidates Claude Code's auto memory, one verdict per memory, `promote` as a flag, index rewritten, dry-run gate, record in `retros/`, promotions handed to `/retro`. Phase 1 commit a3cb89c (`changes/2026-09-29-1329-dream-skill-phase-1.md`), Phase 2 commit 22ab8e7 (`changes/2026-09-29-1348-dream-skill-phase-2.md`).
- **Loop wiring:** `/retro` Step 1 source 7, `## Memory`, Step 6; the context's Learning loop names `/dream`; dream records are a third kind in `retros/` and `retros.jsonl` — `REQUIREMENTS.md` line 31 edited with approval (ADR-042), README, handoff and feedback-log skills aligned.
- **First real run** on this checkout's store: `asking-format.md` and `fallback-paths-stay-simple.md` deleted as graduated on the operator's go, `retros/2026-09-29-dream.md`.
- **Evidence:** cold read r1–r3 (`reviews/2026-09-29-review-dream-skill-coldread-r*.md`), behavioural log `testing/2026-09-29-dream-phase-1-behavioural.md` (twelve `/dream` runs, three `/retro` probes, all 3 of 3 vs 0 of 3), verifier logs phase 1, phase 2, r2 (FAIL on a dispatch needle), r3 (PASS), code and docs lenses APPROVE WITH NITS with Outcomes, BUG-096 and BUG-097.
- **Subagent idea assessed:** `plans/2026-09-29-research-memory-subagents.md` (light mode). Recommendation B: one `dreamer` subagent for the consolidation judgement, capture kept in-thread, a vendor-neutral store by pointing Claude Code's `autoMemoryDirectory` at `~/.config/sh-agent/memory/<slug>/` from `/init-project`, no hooks, nothing auto-applied.
- Earlier this session: `/review --prompts` (11b23b6) and the two research records.

## Open / In Progress

- Nothing in progress. `ROADMAP.md` "Now" and `TODO.md` are empty; no plan has an unticked box.
- Awaiting the operator: the subagent recommendation (B) — a Feature if approved.
- Not released on the operator's instruction. Everything since v1.11.0 is unshipped: `/review --prompts`, `dream`, the retro and context changes, ADR-042.
- Requirements edits awaiting approval: none (the retros clause was approved and applied).

## Key Decisions

- D1–D13 and the Rejected rows in the dream plan; D12 records the cold-read cap decision.
- ADR-042: dream records join `retros/`.
- Research B over A/C/D for the subagent idea (in the record; not decided by the operator yet).

## Corrections

- None. The feedback log has no entry for 2026-09-29; every operator turn was a decision on a choice put to them.

## Agent notes (unverified)

- Claude Code 2.1.283+ ships a built-in `/doctor prompt-audit` (memory docs, read 2026-09-29). The harness `--prompts` lens still gives agy the audit and adds the two-model and bug-index rules; whether Claude Code sessions should prefer the built-in is a question for the operator, not acted on.
- Headless `claude -p` expands a slash command into the first user message; no `Skill` tool call appears (D13). A control must not carry `--add-dir` on the harness skill directory (D13).
- BUG-096's prevention rule (grep the target's `## Step` headings before citing them) is a candidate for `review-standards` § 3 at the next retro; the bugs agent flagged it as a third occurrence of the restatement-drift family (BUG-027, BUG-032).
- The retro probes and the apply run wrote nothing on their own; the verifier's leftover `.tmp` was the only stray path this session.

## Environment State

- Branch: `main`, clean, pushed at 22ab8e7 (plus this handoff).
- Real home on this box: still installed from 0d7dc32 (v1.11.0); `/dream`, `/review --prompts` and the new retro step are not there until `sh-agent update` or `sh-agent install` from this checkout.
- Scratch (gitignored, removable): `.agents/scratch/dream-probe/` (fixture, kits, runs 1–6, c1–c6, r1–r3), `prompts-probe/`, `dream-home/`, `prompts-home/`, `verify-*` leftovers.
- The project's auto memory store now holds two files: the read-neighbours rule and the subagent idea.

## Artifacts to Read

- `.agents/docs/plans/2026-09-29-research-memory-subagents.md` — the assessment the operator asked for; the decision is theirs.
- `.agents/docs/changes/2026-09-29-1348-dream-skill-phase-2.md` — what Phase 2 changed, agents dispatched, what was not exercised.
- `.agents/docs/plans/2026-09-29-dream-skill.md` — the closed plan with D1–D13.
- `.agents/docs/testing/2026-09-29-dream-phase-1-behavioural.md` — every probe run, executed and refused calls separate.
- `.agents/docs/bugs/BUG-096-retro-dream-step-citations.md`, `BUG-097-verifier-needle-phrasing-mismatch.md`.

## Next Steps

1. **Operator:** decide on the subagent research record (B recommended); on a go, `/feature` from its "What approving it would do" paragraph.
2. **Operator:** install this checkout into the real home (`sh-agent install` from here, or `sh-agent update` once a tag exists) and run the three smoke checks of `testing/2026-09-03-verify-real-home.md` on each CLI, plus `/dream --dry-run` here (expect the store consistent, two files, nothing graduated). Antigravity has not run `/dream` or the new `/retro` step yet; that smoke is the first time.
3. **Operator:** a release when wanted — `/preflight vX.Y.Z` over `v1.11.0..HEAD` (a MINOR: new behaviour, nothing breaking); not now, per the operator.
4. Nothing else is queued.
