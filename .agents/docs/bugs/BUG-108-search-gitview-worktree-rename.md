---
owner: bugs
date: 2026-09-30
id: BUG-108
status: fixed
area: internal/search (git status parsing)
---
# BUG-108: `gitView` read a work-tree rename's source path outside the checkout

## Symptom
Security review of v1.12.1..c872598 (`.agents/docs/reviews/2026-09-29-security-v1.12.1-c872598.md`, finding 1): with a tracked path `abc../x` renamed to `moved.md` and `git add -N moved.md` (a work-tree rename, ` R` in the second status column, not staged), `sh-agent search wapiti` printed `../x:1  doc` and the text of a canary file placed beside the checkout. `--refs ../x` also showed it. The moved file's own content (`alpha`) returned 0 hits.

## Repro
Probe W in the review: commit `abc../x`, move it to `moved.md`, `git add -N moved.md`, run `git status --porcelain=v1 -z`, which prints ` R moved.md\0abc../x\0`. `gitView` (internal/search/reconcile.go) only skipped the rename-source record when the INDEX column (first char) was `R`/`C`; here rename is only in the work-tree column, so `abc../x` was parsed as its own record with the first three bytes cut, yielding `../x`. `filepath.Join(root, "../x")` then read a file one directory above the checkout.

## Root cause
The parser's `-z` record walk advanced past a rename's second (from-path) record only for a staged rename (`X == 'R'||'C'`), not a work-tree rename (`Y == 'R'||'C'`), because `git status` was run without `--no-renames`, so its output shape varies by git's own rename-detection heuristic rather than being pinned by a flag. The plan (`.agents/docs/plans/2026-09-29-search-index.md`, D22) had judged the equivalent mutation M6 (deleting this exact branch) safe on the theory that the unskipped path parses to something nonexistent and is dropped; the review's probe shows a real tracked path can land at that computed location once a `..`-shaped prefix is cut into it.
`Diagnosis missed:` none named — root cause was established directly by the security review's probe.

## Fix
`internal/search/reconcile.go`: `git status` now runs with `--no-renames`, so a rename always arrives as a plain delete record plus a plain add record and the special-case branch (`if x == 'R' || x == 'C' { i++ }`) is deleted entirely — no from-path record to skip. Test: `TestWorktreeRenameStaysInsideCheckout` (internal/search/security_test.go) reproduces probe W's arrangement plus a canary file outside the checkout; asserts the canary's word is never found, the moved file is found at its new path, and `../x` never appears. Failed before the fix, passes after; a mutation restoring the old rename branch compiles and fails it. Two further guards the reviewer suggested (`filepath.IsLocal` on every path; never marking a missing file `skipped`) were tried and dropped — no test could reach them once renames are off.

## Prevention
Parse a tool's output only in the shape a flag pins (`--no-renames`, `-z`, `--porcelain=v1` here), never a shape the tool may vary by its own heuristics (git's rename detection). Before recording a mutation as "equivalent" in a plan, re-examine it for the input shape that would make it live — a mutation judged equivalent deserves the same adversarial-input scrutiny as a live finding, not less.

## Related
Plan `.agents/docs/plans/2026-09-29-search-index.md` D22 (the "equivalent" mutation M6 judgement this bug corrects). Security review `.agents/docs/reviews/2026-09-29-security-v1.12.1-c872598.md` finding 1. See also BUG-082 (mutation-testing checklist, `skills/review-standards/SKILL.md` § 2) for the standing mutation protocol; this bug is evidence that an "equivalent" verdict itself needs the same protocol's scrutiny — worth a line there if a third such case appears.
