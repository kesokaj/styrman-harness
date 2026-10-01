---
owner: bugs
date: 2026-10-01
id: BUG-158
status: fixed
area: .agents/docs/testing/kit-preamble.sh (the real-config fingerprint)
---
# BUG-158: the kit preamble ended with exit 2 and no message when `~/.config/sh-agent` did not exist

## Symptom
Sourced on a home with no `~/.config/sh-agent` directory, the preamble ended with exit status 2 and no message, right after printing "kit: real ... before the run:". Second half: the fingerprint of that directory was `ls -lR --time-style=full-iso ... | md5sum`, both GNU-only. Where they are missing (macOS) both commands fail, the fingerprint is the same empty value before and after the run, and the check "the real config changed during the run" could never fire.

## Repro
`TestKitPreambleAssertsCommit` (`internal/harness/catalog_test.go`) runs the preamble under a temporary home: 4 accepting cases came back exit 2. Never seen in a real run: every kit so far ran on this Linux box, where the directory exists and the tools are GNU's. The macOS half was reasoned from the commands, not run.

## Root cause
`ls -la "$dir" 2>/dev/null | sed ...` and `ls -lR ... | md5sum | cut ...` under `set -euo pipefail`: a failing `ls` makes the pipeline fail, and `set -e` ends the script inside a function call or an assignment. Same family as BUG-155 and BUG-157 (pipefail), third shape: the PRODUCER fails, not the consumer. The fingerprint also depended on GNU-only flags, so on another system the check was vacuous.

## Fix
Both pipelines tolerate a missing directory (`{ ... || true; } | ...`). The fingerprint is every path (`find -print`) plus `cksum` of every file, sorted, through `cksum`: tools macOS has too, covering names and contents. Test: the same Go test, with the directory absent (accepts), present and untouched (accepts), a line appended to a file in it (exit 1, "changed during the run"), a new file in it (exit 1).

## Prevention
A script that every kit sources gets a behaviour test under a temporary home before it is trusted, with the unusual states as cases (the directory absent, the tool missing). A check whose two readings can be the same empty value is proven once in its must-fail state (`review-standards` § 2).

## Related
BUG-155, BUG-157 (pipefail family: silent exit, wrong branch, and here a failing producer), BUG-107, BUG-123 (the check this fingerprint serves).
