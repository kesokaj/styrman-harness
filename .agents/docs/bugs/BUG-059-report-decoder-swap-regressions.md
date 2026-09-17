---
owner: bugs
date: 2026-09-17
id: BUG-059
status: fixed
area: internal/harness (report.go redaction — decoder/needle swap regressions)
---
# BUG-059: replacing a matcher and a decoder each introduced a regression the old implementation had handled

## Symptom
The D13 rewrite (whole-line/whole-document redaction, replacing the per-value JSON walk) introduced two new gaps while fixing round 2's findings: (1) the new `authorization` needle regex missed a quoted header name (`"Authorization": "Basic …"`) that the pre-existing secret-assignment needle's quote-tolerant pattern would have caught had it been reused; (2) replacing the JSON walk's `json.Decoder` with `UseNumber()` by a plain `json.Unmarshal` into `any` made a line holding a number beyond float64 (`1e400`) fail to decode, so an escaped needle on that line was sent unredacted — a case the old decoder handled.

## Repro
`.agents/docs/reviews/2026-09-17-security-012b4ca-worktree.md` findings 3 and 4: `headers: {"Authorization": "Basic FAKEVALUE"}` sent with 0 redactions under the new `authorization["']?` gap (before D14's widened pattern); a JSONL line containing `1e400` alongside an escaped needle sent unredacted, confirmed by running the same input against 012b4ca's `UseNumber()`-based walk (`probe-base.log`), which caught it.

## Root cause
The rewrite replaced two working pieces (a matcher and a decoder) with new ones scoped to the round-2 findings, without re-running the old implementation's already-accepted inputs (quoted header names, numbers beyond float64) against the replacements before deleting the old code. Both regressions are "we handled that before, and stopped handling it" — invisible in a diff review that checks new behavior against new tests, not old behavior against old ones.

## Fix
`internal/harness/report.go`, commit 012b4ca (D14): `authorization` needle widened to `(?i)authorization["']?\s*:\s*\S`; `decoded()` now reads a line with a `json.Decoder` using `UseNumber()`, looping over every value, closing both the `1e400` regression and multi-record lines. Covered by `TestReportRedactStructuralGaps` ("number beyond float64 …", "two records on one line") and `TestReportRedactShapes`, mutants `3-authz` and `4-decoded` in `.agents/scratch/r3-mutants/`.

## Prevention
When replacing a decoder or a matcher, run the old implementation's already-accepted inputs (edge-case numbers, quoting styles, encodings it previously handled) as must-pass cases against the new one before deleting the old code — a swap is a regression risk even when it targets a real, named defect.

## Related
BUG-056 through BUG-058, BUG-060 (same review chain, `.agents/docs/reviews/2026-09-17-security-012b4ca-worktree.md`). Three of this session's records (BUG-057 "check fixture bytes", BUG-058 "check every trigger", this one "re-run old accepted inputs") share the pattern "a fix verified against the finding it targets, not against everything the replaced code used to guarantee" — recommend a standing reviewer/security checklist item ("list what the old code accepted/covered; confirm the replacement still does") in `agents/reviewer.md` or `agents/security.md`; actionable (a specific before/after input diff) and generalizable beyond `report.go`.
