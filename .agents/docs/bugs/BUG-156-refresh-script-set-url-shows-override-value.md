---
owner: bugs
date: 2026-10-01
id: BUG-156
status: fixed
area: skills/google-catalog/scripts/refresh-catalog.sh (fetch, a source override's value)
---
# BUG-156: `git remote set-url` showed a source override's value on failure

## Symptom
On a re-run with a source override (`CATALOG_AGENTS_CLI_REPO` or `CATALOG_GOOGLE_SKILLS_REPO`), when the existing work clone's git config could not be written (a stale `.git/config.lock` from an interrupted run, a read-only or full disk), git printed `fatal: could not set 'remote.origin.url' to '<the override's value>'` and the script exited 128 with no message of its own. The script's header, `skills/VENDORED.md` and the testing guide all say no message ever shows such a value (it may hold a token in a URL's path). Found by the security read of 4165f81..7b3e969 (first finding, reproduced with git 2.47.3); older than that range (in v2.2.0).

## Repro
`TestRefreshCatalogScript`: plant `work/agents-cli/.git/config.lock`, run with an override whose path holds the marker `SECRETLOCK`. Red before the fix: the marker and git's `error:`/`fatal:` lines were in the output.

## Root cause
v2.2.0's fix for git's own messages (gate row L3) discarded the output of `git clone` and `git fetch` but not of `git remote set-url`, the third git call that carries the value. The fix was made per call that had been seen failing, not per call that carries the value.

## Fix
`git -C "$dir" remote set-url -- origin "$2" >/dev/null 2>&1 || fail "cannot point the <role> clone at its source (git's own message is not shown ...; remove <dir> and run again)"`. All five git calls in the script were then read: set-url, fetch and clone carry the value and discard git's output; `checkout` and `rev-parse` do not carry it.

## Prevention
A value that must never be shown is tracked by where it FLOWS, not by which call was seen to print it: list every command that receives it and silence or wrap each. The test plants a marker in the value and greps the whole output for it on each failure path. This is `AGENTS.md`'s "needles grow by KIND" rule applied to call sites.

## Related
BUG-155 (same script, same review), BUG-110 (redaction by kind). Other records naming the script: BUG-146, BUG-147, BUG-148, BUG-150.
