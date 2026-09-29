---
owner: bugs
date: 2026-09-29
id: BUG-095
status: fixed
area: process (retro apply — verifier, feature skill, verify log)
---
# BUG-095: retro-applied rule text contradicted its own neighbors twice; verify log counted refused calls as executed

## Symptom
Commit d3cc2b7 applied eight retro proposals to `context/AGENTS.md`, `agents/verifier.md`,
`skills/feature/SKILL.md`, `skills/research/SKILL.md`. The gate code review
(`.agents/docs/reviews/2026-09-28-review-v1.10.0-d3cc2b7-code.md`) returned REQUEST CHANGES, medium: (a)
a verifier rule told a long suite to `tee` its output to a file, so the status read is tee's, not the
suite's, contradicting the same file's "never a pipe whose status hides the command's"; (b) a verifier
rule checked a background job's end with `jobs`, which can never fail since each Bash call is a new
shell; (c) a `skills/feature/SKILL.md` cold-read early exit at the cap contradicted the skill's own gate
list (line 15), box text, and tick rule. The verify log
(`testing/2026-09-28-verify-retro-apply.md`) also counted attempted shared calls as executed, though the
permission gate had refused many. Round-1 fixes (cb6cf73) were reviewed again
(`.agents/docs/reviews/2026-09-28-review-v1.11.0-gate-r2.md`): REQUEST CHANGES again — the cold-read
early exit now skipped real design gaps (reopening BUG-042's shape), and the shared-call rule's
prescribed form was the one the permission gate refused.

## Repro
Not reproduced as a running test; read from the two review artifacts above and
`testing/2026-09-28-verify-retro-apply.md` against commits d3cc2b7 and cb6cf73.

## Root cause
Each retro proposal's rule text was drafted from the proposal's diff without reading the rules already
surrounding the insertion point in the target file — every contradiction found was with a sentence
within 40 lines of the new text (the `tee` rule vs. the pipe-status hard rule; the cold-read early exit
vs. the gate list two lines away). Separately, the verify log's call count read only whether a tool call
was issued, not each result's error/refusal flag, so a permission-gate refusal was recorded as a
successful execution.

## Fix
The operator struck both mechanisms rather than patch a second time: the feature-skill cold-read early
exit and the verifier's shared-call rule were removed; the verifier's `maxTurns` was raised 30→45 to
absorb the shared-call cost instead. Commit to follow this record.

## Prevention
Before inserting a rule from a retro proposal, grep the target file for the mechanism's key words and
read every hit, not just the insertion point. A rule that changes a gate updates every sentence in that
file that lists the gate, in the same edit. A behavioural or verify log's call count reads each tool
result's error/refusal flag and states executed and refused as separate numbers, never one combined
count.

## Related
`.agents/docs/reviews/2026-09-28-review-v1.10.0-d3cc2b7-code.md`,
`.agents/docs/reviews/2026-09-28-review-v1.11.0-gate-r2.md`,
`testing/2026-09-28-verify-retro-apply.md`. Echoes BUG-042 (feature cold-read cap already fixed once for
a contradiction with its own gate list — this is the same shape recurring after a second edit) and
BUG-027 (restatement grep gaps). The "grep the insertion point's neighbors before adding a rule" pattern
now spans three records (BUG-027, BUG-042, BUG-095); recommend it graduate into
`skills/review-standards/SKILL.md` § 3 alongside the existing restatement-grep item, since that section
already owns "every other restatement of a changed rule is changed with it" (BUG-023's graduation) and
this rule is its natural extension — actionable (a grep before insertion) and generalizable beyond this
repo.
