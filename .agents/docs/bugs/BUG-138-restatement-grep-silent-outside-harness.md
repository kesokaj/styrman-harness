---
owner: bugs
date: 2026-10-01
id: BUG-138
status: fixed
area: skills/review (restatement-grep.sh)
---
# BUG-138: restatement-grep.sh exits 1 with no output in any project but the harness checkout

## Symptom
`skills/review/scripts/restatement-grep.sh` (run by the review skill's `--docs` lens before reading by hand) exited 1 with no output on every call outside the harness checkout. The reviewer read "no restatements" and listed it as unverified (its `## Unverified`: "exited 1 with no output every time"). Seen 2026-10-01 in the skill-matrix row `review-docs` on Claude Code: the fixture's `README.md:7` held the word `name`, the script printed nothing.
Second defect: with exactly one matching living file, `grep -n` omitted the file name, so output was `1: text` instead of `README.md:1: text`, breaking the documented `path:line: text` contract.

## Repro
In any git repository without `skills/VENDORED.md`: `bash <path>/restatement-grep.sh -i -- name` gives exit 1, no output, even with a tracked file containing `name`.
Test: `go test ./internal/harness -run TestRestatementGrepRunsOutsideHarness` (new, `internal/harness/restatement_script_test.go`). Before the fix: "want exit 0 with a README hit, got exit 1"; after the first fix: missing `README.md:1:` prefix; passes after both.

## Root cause
`set -euo pipefail` plus a command substitution whose last pipeline stage was `grep -v -E '^(Skill)$'`. With no VENDORED.md the input is empty, that grep exits 1, pipefail fails the substitution, and `set -e` exits the script with status 1 at the `vendored=$(...)` assignment (line 39 before the fix). In the harness checkout VENDORED.md exists, so the pipeline succeeded and every test ran there.
The `-H` defect: `grep -n` prints the file name only when given two or more files.

## Fix
`skills/review/scripts/restatement-grep.sh:39-43` wraps the whole vendored pipeline in `{ ...; } || true`; lines 57 and 62 use `grep -nH`. In the working tree, same change as this record.

## Prevention
1. A script a skill ships for every project is tested in a fixture that is NOT the harness checkout (a temp git repo with none of the harness's files), in both the hit and no-hit state, asserting exit code and output shape; the harness tree makes every optional input present, so a run there proves nothing about other projects.
2. Under `set -e -o pipefail`, a command substitution whose pipeline may legitimately produce nothing ends in `|| true` on the whole group, never only on its first stage. Parent rule: "silence from your own instrumentation is not evidence" (review-standards § 2) was not applied to this script.

## Related
BUG-137 (same script's mode bit; the embedded test caught that one because it ran in the checkout, the shape that missed this one). Introduced by changes `2026-10-01-0610-audit-fix-phase-2.md`. Found by reading the reviewer's artifact in the matrix run, not by a test; no explorer dispatch (a `bash -x` trace stopped at the `vendored=` line).
