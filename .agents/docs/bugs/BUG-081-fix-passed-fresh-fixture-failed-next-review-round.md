---
owner: bugs
date: 2026-09-18
id: BUG-081
status: fixed
area: process (review — fix verification against a fresh fixture)
---
# BUG-081: three fixes each passed their own test and the next review round found the same mechanism still broken

## Symptom
Three times a fix passed its own test and the next round of review found the same mechanism still broken: `ConfigDir`'s `MkdirAll(…, 0700)` tested only on a fresh home where no 0755 dir pre-existed; the BUG-079 identity-check note gated on file identity but was tested only on the branch that installs a binary; a lexical "on PATH" check was tested only on an unsymlinked dir. Found in `.agents/docs/reviews/2026-09-18-security-160aa16-3209cf3.md` and `-3209cf3-e7fec88.md` § Outcome, and `.agents/docs/reviews/2026-09-18-review-harness-merged.md` § Outcome.

## Repro
Not reproduced fresh here — see the three review rounds cited above for the failing arrangements: a pre-existing loose (0755) config dir, the update branch that installs nothing, and a symlinked PATH entry.

## Root cause
Each fix's own test built a brand-new fixture in which the fix could not lose, rather than reconstructing the specific arrangement the review finding had measured (the pre-existing loose directory, the no-install branch, the symlink). A green test therefore proved the fix worked on an easy case, not that it closed the finding.

## Fix
Each mechanism was simplified after asking the operator: `TightenConfigDir` fixes the mode line unconditionally; the BUG-079 note is written only after a binary was actually written; `PathShadows` compares by `os.SameFile` instead of lexical path text.

## Prevention
`context/AGENTS.md` § Definition of Done, step 2 already states: a fix's test builds the arrangement the finding measured — the pre-existing loose directory, the branch that installs nothing, the symlinked path — not a fresh fixture where the fix cannot lose. This rule was added as part of the session that produced this bug and covers it going forward.

## Related
`BUG-078`, `BUG-079` (same session, `Diagnosis missed:` lines feed this finding); retro `2026-09-18-retro.md` finding 3, Proposal 3.
