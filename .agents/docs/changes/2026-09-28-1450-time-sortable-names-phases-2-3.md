---
owner: main
date: 2026-09-28
---
# Time-sortable trail names, Phases 2–3: behavioural proof and review     (date: 2026-09-28)

## Summary
The Phase 1 rule text (`e5a49d4`) was proven behaviourally. Every probe with the new text wrote timed names, and every control with the old text wrote date-only names. The reviews found two sentences that still described the old date step in words, which the path-pattern sweep could not match; both are fixed (D5). Two bugs are recorded. E2E 1, a cold agy resume on a fresh clone of the pushed HEAD, runs after this commit, because it needs this commit's timed handoff to exist.

## Type
feature

## Files changed
- `skills/handoff/SKILL.md` — Step 3.1 names `date +%F-%H%M` for the file name (docs-lens finding 1, code-review finding 1).
- `skills/preflight/SKILL.md` — the dirty-tree changelog dispatch passes the stamp (docs-lens finding 2).
- `.agents/docs/plans/2026-09-28-time-sortable-trail-names.md` — Phase 2 and the review box ticked; D5.
- `.agents/docs/testing/2026-09-28-verify-time-sortable-names-phase-2.md` — the behavioural runs.
- `.agents/docs/testing/2026-09-28-verify-time-sortable-names-fixes.md` — verifier, PASS.
- `.agents/docs/reviews/2026-09-28-review-time-sortable-names-{code,docs}.md` — the reviews, with `## Outcome`.
- `.agents/docs/bugs/BUG-089-…`, `BUG-090-…` and `INDEX.md` — two records.
- `ROADMAP.md` — "Now" at Phase 3.

## Agents dispatched
- `reviewer` ×2: code (APPROVE WITH NITS, 1 finding, fixed) and the docs lens (REQUEST CHANGES, 2 findings, both fixed).
- `bugs`:
  - BUG-089: the changelog agent was told to run `date`, which it cannot;
  - BUG-090: a behavioural fixture did not hold the state its dispatch claimed, a recurrence of BUG-087.
- `verifier`: PASS on the fixed tree.

## Verification
- **Changelog agent** (`--agents`, no Bash, stamp in the dispatch): **3 of 3** timed; the control from `e5a49d4~1`, 0 of 3.
- **Handoff skill** (as `handoff-probe`): **3 of 3** timed; the control, 0 of 3.
- **Context rule** (`--append-system-prompt`): **3 of 3** timed; the control, 0 of 3.
- **Sandbox start:** Claude Code and agy each 1 of 1 (the timed form, "the timed name is newer", GANNET); 0 parse errors in the agy log.
- **E2E 2:** `/init-project` warned on a clone of this repo with the exact count (160) and renamed nothing; it stayed silent on a timed-only project.
- **Verifier after the fixes:** 251 tests pass, drift none, and both installed vendor copies carry the fixed lines.
- **Not exercised:** a new behavioural run for the two one-line fixes; the real home, which is still on v1.6.0.

## Docs
The skill and agent texts are the docs.

## Requirements
n/a
