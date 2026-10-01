---
owner: bugs
date: 2026-10-01
id: BUG-148
status: fixed
area: skills/google-catalog/scripts/refresh-catalog.sh
---
# BUG-148: refresh-catalog.sh's swap loop exited 1 silently half way, leaving the destination with no entries

## Symptom
A refresh over a catalog where one entry directory named in `CHECKSUMS` was already gone exited 1 with no output, after deleting the other entries and before copying anything in: the destination was left with no entries. A second trigger: a `CHECKSUMS` with no line containing `/` made `grep /` exit 1 and the script stopped the same way.

## Repro
The reviewer's fixture: entries alpha, beta, gamma, zeta; `rm -rf dest/zeta`; run the script again: exit 1, empty output, alpha, beta and gamma gone (`reviews/2026-10-01-review-catalog-plan-code.md`, finding 1).

## Root cause
Under `set -euo pipefail` the loop body ended with `[ -d "$dest/$old_entry" ] && rm -rf -- ...`. When the last name visited was not a directory the `&&` list returned 1, the `while` loop returned 1, `pipefail` passed it on and `set -e` ended the script with no message. The same shape as BUG-138 (silent exit under `set -e`), plus a partial write because the exit came between the delete and the copy.
`Diagnosis missed:` the first fix covered trigger A with a test and trigger B only in code; the rule "a finding is fixed only when every trigger it names has a test that fails without the fix" was applied to one trigger of two. Review round 2 found it.

## Fix
`if [ -d ... ]; then rm -rf -- ...; fi`, and the list is collected through `{ grep / || true; }`. Tests in `TestRefreshCatalogScript`: a refresh after removing an entry directory must exit 0 and restore the catalog (red before the fix: "a refresh over a catalog with a missing entry directory did not restore it (exit 1)"); a refresh over a destination holding only an empty `CHECKSUMS` must succeed (fails with `grep /` put back). Fixed in the working tree, before any release carried the script.

## Prevention
In a bash script under `set -euo pipefail`, a loop body or pipeline stage never ends in `test && command` or a bare `grep`: write `if ...; then ...; fi` and `{ grep ... || true; }`. A destructive loop is tested with its input in the state a half-done earlier run leaves behind.

## Related
BUG-138 (earlier instance of the silent `set -e` exit, second occurrence with this one), BUG-147, BUG-146; review `reviews/2026-10-01-review-catalog-plan-code.md`. Graduation candidate: the `bash` skill's script checklist; actionable (a missing-entry fixture and a no-match fixture catch it) and generalizable (any `set -e` script).
