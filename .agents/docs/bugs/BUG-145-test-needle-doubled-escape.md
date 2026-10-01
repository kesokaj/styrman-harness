---
owner: bugs
date: 2026-10-01
id: BUG-145
status: fixed
area: internal/dream (collect_test.go, feedback redaction)
---
# BUG-145: the escaped-secret redaction test passed with its fix reverted, then failed with the fix present

## Symptom
`TestFeedbackRedactsEscapedSecret` (`internal/dream/collect_test.go`), the test for the dream's JSON-lines feedback redaction, PASSED with the fix reverted on its first must-fail control, so it proved nothing. After a needle change, a second control failed the other way: FAIL with the fix present.

## Repro
Write a Go raw-string literal from a Python heredoc with `\\\\u0026`. It lands in the Go source as `\\u0026` (two backslashes), so the JSON line decodes to the literal text `\u0026` instead of `&`. The escaped-secret shape under test never existed in the fixture. Separately, the first needle (an `AKIA...` key with one letter written as a `\uXXXX` escape) was one the raw redaction already caught, so it could not show the gap even with a correct escape.

## Root cause
A fixture meant to carry an escape went through two layers of string quoting (heredoc, then Go literal) and the written file was never checked for the escape itself. The global hard rule for exactly this ("after writing text meant to keep an escape, grep -n for the escape itself") was not followed. The needle was chosen from assumption, not from the gate's measured shape (`\u0026sig=`), so the control had nothing to separate fixed from unfixed.

## Fix
The fixture line now reads `\u0026sig=<32 hex>` with one backslash, checked with `grep -n 'u0026sig'`. The must-fail control now fails without the fix and the test passes with it. Files: `internal/dream/collect_test.go` (and `collect.go` for the redaction it covers), uncommitted at record time.

## Prevention
- A test that sets up an unusual condition proves the condition holds first: for an escape, `grep -n` the written file for the escape itself, as the hard rule says.
- A must-fail control that PASSES means the fixture, not the code, is wrong: fix the fixture before touching the implementation. Choose the needle from the measured shape, not from a guess.

## Related
Second occurrence of the escape rule being skipped in this project; the rule already existed (global hard rules, "a test that sets up an unusual condition proves the condition holds first"), so it was not followed. Graduation candidate: a `review-standards` § 2 line "a must-fail control that passes is a fixture defect first"; actionable (the mutation protocol already runs the control) and generalizable. BUG-071 graduated the "condition HOLDS" clause. BUG-127 (redaction shapes), BUG-142 (process slips in this plan family).
