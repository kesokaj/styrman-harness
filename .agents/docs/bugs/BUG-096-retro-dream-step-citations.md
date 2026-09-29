---
owner: bugs
date: 2026-09-29
id: BUG-096
status: fixed
area: skills/retro (cross-skill step citations)
---
# BUG-096: `skills/retro/SKILL.md` cited `dream`'s step numbers wrong, twice

## Symptom
Two sentences in `skills/retro/SKILL.md` cite `dream`'s numbered steps and both cite the wrong numbers: line 44 said "Steps 1–3 and 5" for the verdict table, omitting dream's Step 4 (the index rows, verdict `index`) and implying application happens at retro's Step 6; line 138 said "Steps 5 and 6: the memory files, the index rewrite", but in `dream` the memory-file edits are Step 3 and the index rewrite is Step 4 — Step 5 is the gate, Step 6 the record.

## Repro
`grep -n "^## Step" skills/dream/SKILL.md` against the two citations in `skills/retro/SKILL.md:44` and `:138`; the numbers do not match. Found by `reviewer`, `.agents/docs/reviews/2026-09-29-review-dream-code.md` findings 1 and 2 (medium, confidence high).

## Root cause
The two sentences were written from memory of the `dream` skill while it was still being drafted in the same task; the step numbers were never checked against `skills/dream/SKILL.md`'s actual `## Step` headings before being committed to `retro`'s file.

## Fix
Both sentences reworded to the correct steps (`skills/retro/SKILL.md:44,138`), per the review's Outcome section: line 44 now cites "Steps 1 to 5" with application "at Step 6 below"; line 138 now cites dream's Step 3 (memory-file edits), Step 4 (index rewrite), and Step 5 (the gate).

## Prevention
Before writing a sentence that cites another file's numbered steps or sections, run `grep -n "^## Step"` (or the equivalent heading pattern) on the target file and keep it on screen while writing the citation; name what the step does, not only its number. This is a specific case of `review-standards` § 3 ("every other restatement of a changed rule is changed with it") — extend it, or the project `AGENTS.md`, to cover cross-file numeric citations specifically, since three records now share this shape (BUG-027, BUG-032, BUG-096 all involve stale/wrong restatements caught only by grepping the subject, not the old wording).

## Related
`.agents/docs/reviews/2026-09-29-review-dream-code.md`; BUG-027, BUG-032 (same prevention-rule family, restatement drift); `skills/review-standards/SKILL.md` § 3.
