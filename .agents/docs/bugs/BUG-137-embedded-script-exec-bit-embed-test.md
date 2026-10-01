---
owner: bugs
date: 2026-10-01
id: BUG-137
status: fixed
area: embed / release binary
---
# BUG-137: a 0755 script under an embedded tree failed TestEmbeddedPlanEqualsCheckout

## Symptom
`TestEmbeddedPlanEqualsCheckout` (root package) failed in `testing/2026-10-01-verify-audit-fix-phase-3.md`: `skills/review/scripts/restatement-grep.sh` planned as 0644 from `embed.FS` and 0755 from the tree. A release binary would have installed the review script with a different mode than a checkout install.

## Repro
Commit any file with mode 0755 under `context/`, `skills/`, `agents/`, `commands/` or `rules/`, then run `go test .`.

## Root cause
Phase 2 of plan `2026-09-30-audit-fix.md` (commit 329a99b) added the script as 0755 and ran no `go test`, because "no Go source changed": the guide's § 1 header said to skip the Go checks when no `.go` file is in the diff. The root package's embed test compares those five trees with the binary's embedded copy, and an embedded filesystem carries no exec bit (`plan.go:137-138` takes the mode from file info). So a non-Go file in an embedded tree is an input to a Go test, and the skip rule hid it.

## Fix
Both scripts are now 0644 and the script is always invoked as `bash <path>` (the review skill already says so). Guide § 1's header now says every change runs the Go checks, including one that touches only the five embedded trees. The script's header comment explains the mode.

## Prevention
A file under the five embedded trees is a source change for the embed tests whatever its extension: run guide § 1 on every change, never skipped for "no .go file in the diff". Never give a file in those trees an exec bit; invoke scripts as `bash <path>`.

## Related
Plan `2026-09-30-audit-fix.md` Phase 2 (commit 329a99b, the regression); Phase 3 verify log `testing/2026-10-01-verify-audit-fix-phase-3.md`; BUG-048 (the gofmt walk), the neighbouring embed/guide lesson.
