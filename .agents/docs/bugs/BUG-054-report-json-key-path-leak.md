---
owner: bugs
date: 2026-09-17
id: BUG-054
status: fixed
area: internal/harness (report.go redactJSON)
---
# BUG-054: a secret in a JSON key was written verbatim into `report.json` as the redaction's recorded path

## Symptom
Security review found that when a needle sits inside a JSON object KEY (not its value), the sent copy of the source file redacts correctly, but `report.json` — the one file every run sends — records the full key, secret included, as the `Path` of the `Redaction`.

## Repro
Probe fixture `{"tool_output":{"https://api.example.test/v1?token=FAKE…":"404"}}` in a `retros.jsonl` copy: sent `retros.jsonl` is redacted, but the generated `report.json` contains `"path": "tool_output.https://api.example.test/v1?token=FAKE…"`.

## Root cause
`redactJSON` (`internal/harness/report.go:175-181` at 015cf2d) calls `fire(p, n, true)` on a key hit using `p`, the full key path including the offending key text, and stores that path unmodified in the `Redaction` struct, which is marshalled into `report.json` and printed to stdout. `TestReportRedactDocument` asserted only `len(recs) == 1` for the key case, so a test existed but checked count, not content — it could not catch the leaked value in the recorded path.

## Fix
`redactJSON` collapses a key hit to one record with `Path` cleared (`internal/harness/report.go`, commit 012b4ca), since the whole value is already replaced and a path adds nothing. `TestReportRedactStructuralGaps` asserts the marshalled record carries no `FAKEVALUE`; must-fail seen.

## Prevention
A test that verifies redaction correctness must assert on content (the marshalled/sent bytes contain no secret), not just record count — `len(recs) == 1` proves a record was made, not that it is safe. When a field is meant to carry metadata only (a path, a location), assert directly that it never carries the value it points at.

## Related
BUG-053, BUG-055, BUG-056 (same review). See BUG-053's Related note on the shared adversarial-fixture / content-assertion rule.
