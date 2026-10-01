---
owner: bugs
date: 2026-10-01
id: BUG-155
status: fixed
area: skills/google-catalog/scripts/refresh-catalog.sh (the per-entry checks)
---
# BUG-155: the refresh ended with exit 141 and no message when an entry held very many files of a refused kind

## Symptom
The refresh exited 141 with no output when an upstream entry held very many files of a refused kind. Found by the security read of 4165f81..7b3e969 (`.agents/docs/reviews/2026-10-01-security-4165f81-7b3e969.md`, fourth finding, by reading), then reproduced by the main agent: an entry with 1,500 dotfiles (about 200 KB of paths) gave exit 141 where the refusal "entry ... holds a dotfile ..." was owed. In the script since the catalog was built (v2.1.0); never hit on the real upstream.

## Repro
`TestRefreshCatalogScript`, fixture "many dotfiles" (`internal/harness/catalog_test.go`), red before the fix: `want exit 1 naming "holds a dotfile", got exit 141`.

## Root cause
Three checks read `bad=$(find ... -print | head -1)` under `set -euo pipefail`. `head` exits after one line; `find`, still writing more than a pipe holds (64 KB on Linux), is killed by SIGPIPE; `pipefail` makes the pipeline's status 141 and `set -e` ends the script inside the assignment, before the `fail` on the next line. Small inputs fit in the pipe, so every earlier test passed.

## Fix
The three finds end `-print -quit` (no pipe). The class was closed in the same change: the two other command substitutions that piped into `head -1` (`fetched()`'s awk, the manifest date's sed) are now one awk each that exits after its first match; the one remaining `| head` (`grep ... | head -20 >&2 || true`, the writer log excerpt) is guarded by `|| true`.

## Prevention
Under `pipefail`, never `producer | head -N` in a command substitution or as a statement: the producer's SIGPIPE ends the script with 141 and no message once the output outgrows a pipe. Make the producer stop itself (`find -quit`, `awk '{...; exit}'`), and test each refusal with enough input to fill a pipe.

## Related
BUG-148 (same script, a silent exit under pipefail in the swap loop), BUG-138 (same family). BUG-156 (same script, same review). Other records naming the script: BUG-146, BUG-147, BUG-150.
Graduated 2026-10-01: the prevention rule is now `TestScriptsHaveNoBareHeadPipe` in `internal/harness/catalog_test.go` (fails on a `| head` line not ending `|| true` in any `skills/*/scripts/*.sh` that sets `pipefail`).
