---
owner: bugs
date: 2026-09-17
id: BUG-058
status: fixed
area: internal/harness (report.go redaction — review Outcome accuracy)
---
# BUG-058: a review Outcome marked a finding "fixed by deleting the mechanism" while one of the finding's two triggers still leaked

## Symptom
Round 2's security review Outcome recorded finding 2 (a key body sent past its needle line) as fixed by deleting the JSON walk. The finding actually named two triggers; the second — a header line that also mentions a footer, closing the block above the key body it should still cover — lived in `blockEnd`, which the JSON-walk deletion never touched. Round-3 security finding 2(b) found the key body still sent.

## Repro
`internal/harness/report.go`'s `blockEnd` (015cf2d/012b4ca) closes a block at any `-----END`/`END…PRIVATE KEY` found after the header, including one on the header's own line naming an unrelated footer; text like `pasted -----BEGIN OPENSSH PRIVATE KEY----- up to -----END below` followed by body lines sent the body with 2 redactions recorded (`.agents/docs/reviews/2026-09-17-security-012b4ca-worktree.md` finding 2(b)).

## Root cause
The round-2 fix addressed finding 2 by removing the JSON walk that produced its most visible symptom (a per-value redaction bypass), without re-checking the finding's other named trigger against the new code. The Outcome recorded "fixed" against the finding's headline case only, not against every mechanism the finding described.

## Fix
`internal/harness/report.go`, commit 012b4ca (D14): a block now closes on a header's own line only when a footer carries that header's own label, not any footer mention. Covered by `TestReportRedactStructuralGaps` ("header line that mentions a footer"), mutant `2b-label` in `.agents/scratch/r3-mutants/`.

## Prevention
When a review Outcome marks a finding fixed, check every trigger the finding names — not only the headline reproduction — against a must-fail test before recording "fixed". "Fixed by deleting the mechanism" is only true if the deletion removes every path the finding described.

## Related
BUG-053 through BUG-057 and BUG-059, BUG-060 (same review chain, same dispatch). This is the second time in this session an Outcome's "fixed" call needed re-verification per-trigger (see BUG-057's mutant discipline); if a third instance appears the underlying rule — an Outcome entry names each trigger it checked, not just the finding number — should graduate into `agents/reviewer.md` or `agents/security.md`'s Outcome-writing guidance.
