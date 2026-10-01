---
owner: bugs
date: 2026-10-01
id: BUG-157
status: fixed
area: .agents/docs/testing/kit-preamble.sh (the KIT_ASSERT_COMMIT check)
---
# BUG-157: a kit asserting a commit that IS in a long-history clone was refused

## Symptom
`kit: commit <sha> not in the clone's log` (exit 2) for a `KIT_ASSERT_COMMIT` that the clone holds, when the clone's history was long. Found by the code review of 7b3e969..a224329 (`.agents/docs/reviews/2026-10-01-review-7b3e969-a224329-code.md`, finding 1, medium): 30 of 30 runs against this repository's history (423 commits) asserting the head commit. Never seen in a real run: every kit so far used the demo fixture, whose log is a few lines.

## Repro
Source the preamble with `KIT_CLONE` a clone of this repository and `KIT_ASSERT_COMMIT` its head's short SHA.

## Root cause
The check was `git log --format=%h | grep -q -- "^$sha" || { ...; exit 2; }` under `set -euo pipefail`. `grep -q` exits at the first match; `git log`, still writing, is killed by SIGPIPE; `pipefail` makes the pipeline's status 141, and the `||` branch refuses. This is the FAIL-OPEN form of BUG-155's class (there: a silent exit; here: no silence, a wrong branch taken). Small logs fit in a pipe, so the demo fixture never showed it.

## Fix
The check asks git directly: `git -C "$KIT_CLONE" merge-base --is-ancestor "$KIT_ASSERT_COMMIT" HEAD`. Checked by hand on a clone of this repository, both states: the head commit and the root commit accepted (20 of 20 for the head), an absent SHA refused with exit 2. The same change widened BUG-155's test (now `TestScriptsHaveNoEarlyExitPipe`) to this file and to `grep -q`, `grep -m` and `sed ...q`; it was red on this line before the fix.

## Prevention
Under `pipefail` a pipeline into a consumer that exits early is wrong in a condition too (`if`, `||`, `&&`): status 141 reads as "no match". Ask the tool directly, or feed the consumer a here-string. Graduated: pinned by `TestScriptsHaveNoEarlyExitPipe` in `internal/harness/catalog_test.go`.

## Related
BUG-155 (same class, the silent form; its test was widened here), BUG-148, BUG-138 (same family).

## Addendum 2026-10-01
The fix (`git merge-base --is-ancestor <value> HEAD`) accepted any revision: `HEAD`, a branch or tag name, `@`, `:/text` passed on every clone, so the assertion proved nothing for such a value (code review and security read of a224329..01255bc, low). Now the value must match `^[0-9a-f]{7,40}$` before git sees it, and git gets `--end-of-options "<sha>^{commit}"`. Checked by hand by the main agent on a clone of this repository: the 7-digit and the 12-digit SHA of the head accepted; `HEAD`, `main`, `abcd`, `--all` refused by the shape check; `deadbeef1` refused as not in the log; an empty value skips the assertion as before. No Go test runs the preamble; the static pipe check is its only test, and that is said in the gate's verifier log.
