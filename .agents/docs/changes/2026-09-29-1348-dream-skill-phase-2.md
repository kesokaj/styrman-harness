---
owner: main
date: 2026-09-29
---
# `dream` skill, Phase 2: the loop, the record kind, the first real run     (date: 2026-09-29)

## Summary
The `dream` skill (Phase 1, `changes/2026-09-29-1329-dream-skill-phase-1.md`) is wired into the learning loop and has run for real:

- **`/retro` reads it first:** Step 1 gains source 7, memory evidence (the skill's dry run, Steps 1 to 5); the record template gains `## Memory`; Step 6 applies the approved verdict edits with the proposals. Three `/retro --dry-run` probes on the seeded fixture carried the store line, the verdicts and the merged pair as a finding with `rule_ref null`, 3 of 3.
- **The global context** (§ Learning loop) names `/dream` and says every retro, distill and dream leaves a record in `retros/`; size test 21,643 of 22,000 bytes.
- **Dream records are a third kind in `retros/` and `retros.jsonl`** — `REQUIREMENTS.md` line 31 says so with the operator's approval ("Approve the edit (Recommended)", ADR-042), and README, the handoff skill's artifact table and the feedback-log skill's pointer say the same in the same words.
- **First real run,** by hand from the working-tree skill on this checkout's own store: two memories `graduated` (their rules live in the installed global context at lines 6 and 17), two kept; on the operator's "Go all (Recommended)" both were deleted, the index rewritten, `retros/2026-09-29-dream.md` quotes them in full, one `kind: dream` line in `retros.jsonl`. E2E criterion 1 met.
- `skill-matrix.md` gains a `dream` row whose check passed on the Phase 1 clones.

The plan is closed: every box ticked, one moved (the real-home install and smoke, the operator's, in the closing handoff). `ROADMAP.md` "Now" is empty again. Not released: the operator asked for no new version yet.

## Type
feature (phase 2 of 2)

## Files changed
`skills/retro/SKILL.md` (Steps 1, 4, 6); `context/AGENTS.md` (two clauses in § Learning loop); `REQUIREMENTS.md` (one clause, approved); `README.md` (one row, by `documentation`); `skills/handoff/SKILL.md` (one row); `skills/feedback-log/SKILL.md` (one sentence); `.agents/docs/architecture/DECISIONS.md` (ADR-042); `.agents/docs/testing/skill-matrix.md` (row `dream`); the plan; `ROADMAP.md`; the retros record.

## Agents dispatched
- `verifier` ×3 — `testing/2026-09-29-verify-dream-skill-phase-2.md` PASS; `-r2.md` FAIL on one grep whose needle in the dispatch did not match the file's own wording (the change was right; BUG record below); `-r3.md` PASS after the wording was aligned across all five restatements.
- `reviewer` (code lens) — `reviews/2026-09-29-review-dream-code.md`, APPROVE WITH NITS: two wrong step citations from retro into dream, fixed (Outcome on the report).
- `reviewer` (docs lens, the `/review --docs` of this change) — `reviews/2026-09-29-review-dream-docs.md`, APPROVE WITH NITS: four enumerations fixed, one with the operator's approval (Outcome on the report).
- `documentation` — README row 210.
- `bugs` — two records for this round: the cross-skill step citations, and the dispatch needle that failed a correct change.

## Verification
- `testing/2026-09-29-dream-phase-1-behavioural.md` § Phase 2: `/retro --dry-run` probes r1–r3, 3 of 3; the apply run on the real store; the matrix row check on clones 4–6.
- Not exercised: Antigravity runs of `/dream` and the new `/retro` step (D10; the real-home smoke after `sh-agent update` is the operator's), the `feedback.jsonl` route to `promote`, `--fix`-style application of dream verdicts from inside `/retro` on a real store (the retro probes were dry runs).

## Notes
- The verifier left an empty `verify-dream-skill-phase-2.md.tmp` beside its log; deleted by name.
- The memory docs read for D7 also document Claude Code's built-in `/doctor prompt-audit` (2.1.283+); it bears on the `--prompts` lens shipped earlier today, noted in the handoff.
