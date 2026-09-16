---
owner: bugs
date: 2026-09-16
id: BUG-047
status: fixed
area: docs/deploy (release procedure)
---
# BUG-047: signed-tag release check compared a placeholder to itself

## Symptom
The first draft of the signed-tag release procedure (`remote.md` §
Releases, ADR-026) checked the tag object (type, armor, signature) but
never what the tag pointed at (security review round 1, medium finding).
The first fix — `test "$(git rev-parse vX.Y.Z^{commit})" = "<release
commit>"` with the same `<release commit>` placeholder passed to
`git tag -s` — was tautological (round 2, medium finding: a tag on a
stale commit passed every check), and the verifier's round-1 must-fail
case for it varied a state the flow cannot produce.

## Repro
Not re-run as code; reproduced by inspection of the two review rounds
(`.agents/docs/reviews/2026-09-16-security-release-procedure-adr-026.md`).

## Root cause
The check's expected value (`<release commit>`) was hand-substituted at
the same step that produced the value being checked, so the comparison
could never fail — it verified that a variable equalled itself, not that
the tag pointed where the operator intended.

## Fix
Fixed by deletion of the placeholder: the tag takes no target argument
and sits on HEAD; the checks compare the tag with `git rev-parse HEAD`,
HEAD's parent with `<gated head>` (the head SHA from the gate report's
`Range reviewed:` line — the one external input), and assert
`git diff --quiet <gated head> HEAD -- ':(top)' ':(top,exclude).agents/docs'` (records
only). Rehearsed must-pass and must-fail in
`.agents/scratch/tagcheck/`.

## Prevention
A check's expected value never comes from the same hand-substituted
input as the thing it checks — bind it to something produced elsewhere
(HEAD, a prior gate report), never a placeholder threaded through both
sides. A must-fail case is validated only in a state the real flow can
actually produce, or it validates nothing.

## Related
Second round of findings on one mechanism (the tag-target check),
resolved by deletion rather than another patch, per
`context/AGENTS.md` § Definition of Done step 2 ("stop fixing and ask
what the mechanism is for ... often deleted rather than repaired").
Shares a pattern with BUG-026 (repeated patching of one mechanism across
review rounds) — recommend the "two rounds on one mechanism stops the
patching" rule (already in `context/AGENTS.md`) be checked explicitly by
`reviewer` on round 2 of any dispatch, since BUG-026 and BUG-047 both hit
it before a human/process caught it.
