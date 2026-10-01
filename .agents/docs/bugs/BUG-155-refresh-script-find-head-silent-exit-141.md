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
Under `pipefail`, never `producer | head -N` in a command substitution or as a statement: the producer's SIGPIPE ends the script with 141 and no message once the output outgrows a pipe. Make the producer stop itself (`find -quit`, `awk '{...; exit}'`), and test each refusal with enough input to fill a pipe. The fail-open form: the same pipe in a condition (`producer | grep -q x || refuse`) does not end the script, it takes the wrong branch, because 141 reads as "no match" (BUG-157); ask the tool directly or feed the consumer a here-string.

## Related
BUG-148 (same script, a silent exit under pipefail in the swap loop), BUG-138 (same family). BUG-156 (same script, same review), BUG-157 (the fail-open form). Other records naming the script: BUG-146, BUG-147, BUG-150.
Graduated 2026-10-01: the prevention rule is now `TestScriptsHaveNoEarlyExitPipe` in `internal/harness/catalog_test.go` (fails on an unguarded pipe into `head`, `grep -q`, `grep -m` or `sed ...q` in any tracked `skills/*/scripts/*.sh` or `.agents/docs/testing/*.sh` that sets `pipefail`).
2026-10-01 (next commit): the test was first `TestScriptsHaveNoBareHeadPipe` (renamed) and was widened via helper `earlyExitPipeProblem`; red on two real lines: `skills/review/scripts/restatement-grep.sh:55` (`printf ... | grep -q ... || { keep=0; break; }`, harmless until a line outgrows a pipe, now a here-string) and the kit preamble's line (BUG-157).
2026-10-01 (second): the scanner as committed in 4487a7a treated everything inside double quotes as data, so `v="$(find ... | head -1)"` passed (code review of 01255bc..4487a7a, medium); a command substitution is now code again inside double quotes, and a `sed` script with `q` inside braces is matched. Pointer, not part of this bug: the same review and the security read found the resolution-order defect in the refresh script's pin (`checkout --detach <short pin>` takes a tag of that name first; older than this release); the script now refuses a pin whose checked-out commit does not start with it. Same cause as BUG-157's second addendum. (Superseded, see below: the prefix check still passed.)
2026-10-01 (third, next commit): the security read of 4487a7a..5bbd789 (finding 2, low) showed a tag on a commit made to share the pin's prefix still passed the prefix check, so a pin is now a full 40-digit commit id, refused by shape before any clone, and the post-checkout check is gone (superseded, see below). Also: the test's `sed` pattern as committed in 5bbd789 let through `/re/!q`, `/re/Iq`, `/re/Mq` and `\%re%q` (code review of 4487a7a..5bbd789, medium); the character before the command is now anything but a lowercase letter or underscore, with `Q` and an exit code taken too; 29 must-refuse lines.
2026-10-01 (fourth): removing the post-checkout check was a mistake, found before any release held it (security read of 5bbd789..42670ff, medium, probed with git 2.47.3): `git checkout --detach <full id>` takes a local branch of that exact name first, and a clone makes the source's default branch local, so a source whose default branch is named like the pin chose the content, with exit 0, even for a pin it did not contain. Now the script checks out `<pin>^{commit}` and compares the checked-out commit's id with the pin. Test in `TestRefreshCatalogScript`: a fixture source whose default branch is named as the pin and points at another commit; the other commit's file was vendored before the fix (seen red). Lesson: a check removed because "the new rule makes it unreachable" needs the same proof as a check added; here the rule (40 digits) closed tags and short ids, not branches.
