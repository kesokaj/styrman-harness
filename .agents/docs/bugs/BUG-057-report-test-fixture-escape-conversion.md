---
owner: bugs
date: 2026-09-17
id: BUG-057
status: fixed
area: internal/harness (report_test.go fixture authoring)
---
# BUG-057: a file-editing tool converted a JSON escape into its literal character, so a test named for that escape never exercised it

## Symptom
`internal/harness/report_test.go`'s case named `"needle behind a JSON unicode escape"` held the literal string `AKIA…`, not the six-byte escape `\u0041KIA…` (this record was itself written with the escape converted, and corrected by script — the defect it describes) its name and surrounding comment claimed to add coverage for. The same silent escape-to-character conversion also altered a Go comment in `report.go` (an escaped example became the literal character) and one line of a security agent's own report file.

## Repro
`git hash-object` on the pinned test blob, then `od -c internal/harness/report_test.go` around the fixture line: the bytes read as the literal character, not backslash-u-zero-zero-four-one. Confirmed by stubbing `decoded()` to return `""` in a scratch copy — `TestReportRedactStructuralGaps` still passed in full (`.agents/scratch/mutant-test/`), so the case did not depend on the decode path it claimed to cover.

## Root cause
The tool used to write the fixture silently converted the JSON unicode escape sequence into the actual Unicode character while writing the file, rather than preserving the literal backslash-escape bytes. Because the resulting character (from the `AKIA` prefix) still matched the raw-text needle pass, the test kept passing and the mismatch between the test's name/intent and what it exercised was invisible until a mutant (`decoded()` disabled) was run against it.

## Fix
`internal/harness/report_test.go`, uncommitted at the time of writing on top of 012b4ca (landed in the commit carrying D13 and D14): fixtures rewritten through a script that writes explicit escape bytes, checked with `grep -n` and `od -c` before commit. Must-fail confirmed by stubbing `decoded()` (`.agents/scratch/r2-decoded-mustfail.log`) and by per-fix mutants (`.agents/scratch/r3-mutants/*.log`).

## Prevention
After writing any fixture or text containing an escape sequence (`\u`, `\n` inside a string literal meant to stay literal, etc.) via an editing tool, check the actual bytes (`grep -n` / `od -c`) before trusting the file. A test's name is not evidence that it exercises the path it names — run a mutant that removes the path under test and confirm the test fails.

## Related
BUG-053, BUG-054, BUG-055, BUG-056 (same review, same dispatch, `.agents/docs/reviews/2026-09-17-review-012b4ca-worktree-code-docs.md` finding 1, `.agents/docs/reviews/2026-09-17-security-012b4ca-worktree.md` finding 6). BUG-056 already recommends graduating "adversarial fixtures + must-fail mutants" into `context/AGENTS.md` Hard rules or `agents/security.md`; this record's "check fixture bytes after any escape-bearing write" is a narrower, generalizable companion rule (it applies to any tool-authored fixture, not just redaction) and belongs in the same place.
