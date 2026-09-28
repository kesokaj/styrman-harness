---
owner: main
date: 2026-09-28
---
# Deployed skills and agents no longer carry this repository's history     (date: 2026-09-28)

## Summary
The operator: "the preflight skill has my private repo, it should not be there. The preflight skill is universal". Approved as items 4 and 5.

Preflight named no repository or signing key, but it did cite this repository's history as reasons: BUG-033/034/083, a 2026-09-18 anecdote, and a GitHub HTTP 500 story. Eleven other deployed files did the same, and one of them cited styrman-crew's bugs. A project running these files would look for those IDs in its own trail.

Every citation is removed and every rule kept. Where a parenthetical carried a reason, the reason stays in generic words.

## Type
task

## Files changed
- **Agents:** `agents/{changelog,explorer,reviewer,security,verifier}.md`.
- **Skills:** `skills/{distill,feature,golang,init-project,preflight,review,review-standards}/SKILL.md`.
- 12 files in total, 23 lines.
- **Left as they are:**
  - `agents/bugs.md:67`, a generic example row showing the INDEX format;
  - `skills/frontend/SKILL.md:60`, which dates a version fact;
  - the `init-project` `DECISIONS.md` template's `ADR-001` placeholder.
- **Follow-up, not changed:** `skills/find-skills/SKILL.md:97` names `styrman-harness` in a vendored skill (an in-place edit since the first commit).

## Agents dispatched
- `verifier`: PASS (`testing/2026-09-28-verify-universal-cleanup.md`).
- `reviewer`: APPROVE WITH NITS (`reviews/2026-09-28-review-universal-cleanup.md`). Both nits are fixed, and a widened sweep found one more (`init-project`); see its `## Outcome`.

## Verification
- **Verifier:** race suite green; install drift none; 27 skills and 7 agents parse. The positive control grep fires at HEAD (19 hits), and the working tree has 0 outside the named exceptions.
- **Not verified by the verifier:** the three prose edits after its run (verifier `:37`, review `:62`, init-project `:50`). The v1.8.0 gate's verifier covers them.
- **No behavioural run** for the five agent edits, which are citation removals only; the reviewer accepted this deviation.
- **Docs lens:** runs in the v1.8.0 gate over this range.

## Docs
The skill and agent texts are the docs.

## Requirements
n/a
