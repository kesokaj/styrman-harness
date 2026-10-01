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

## Addendum 2026-10-01
The security read of 7b3e969..a224329 (finding F3, low) found channels no redirect covers: a git trace target named by the caller's environment or config (`GIT_TRACE`, `GIT_TRACE2_EVENT`: reproduced, the override's value was written into both files), git's terminal credential prompt (shows the URL; not probed), the process list while git runs, and `bash -x`. Fix: the script exports `GIT_TERMINAL_PROMPT=0` and sets every `GIT_TRACE*` variable and the three `GIT_TRACE2*` targets to 0 for its own calls; its header now says the value is in the process list while git runs and that `bash -x` prints it, so a token belongs in a credential helper. Test: `TestRefreshCatalogScript` runs with `GIT_TRACE` and `GIT_TRACE2_EVENT` pointing at files and an override holding a marker; both files held the marker before the fix. Not closed by code: the process list and `bash -x`, stated in the header. Lesson for Prevention: "where it flows" includes git's own side channels (trace, prompt), not only its stdout and stderr.

## Addendum 2026-10-01 (second)
The first addendum's fix was itself wrong in one line and incomplete in two, found by the code review and the security read of a224329..01255bc before any release held it.
- (a) `GIT_CURL_VERBOSE=0` turned curl's trace ON: git checks only that the variable is set (security read probe: 31 trace lines with the request path). Now `unset GIT_CURL_VERBOSE`, pinned by a text assertion in `TestRefreshCatalogScript` because no run shows the difference (the calls' stderr is discarded).
- (b) The loop zeroed an exported `GIT_TRACE_REDACT`, where 0 switches redaction off. It is now left out of the loop and exported as 1.
- (c) `GIT_TERMINAL_PROMPT=0` does not stop an askpass program (`GIT_ASKPASS`, `core.askPass`, `SSH_ASKPASS`), which is handed a prompt naming the URL. The script now exports an empty `GIT_ASKPASS`, which also stops the fallback to the other two (probed by the security read). Test: an https fixture server answering 401 and an askpass stub that logs its arguments; the stub ran twice before the fix ("Username for 'https://127.0.0.1:<port>'", "Password for ...").
- (d) The trace test set only environment targets; it now also runs under a `trace2.eventTarget` in git's config (the main agent proved the precondition by hand: the marker was written to the target without the script's settings, and not with `GIT_TRACE2_EVENT=0`).

Prevention, added: an environment variable is "off" only by the semantics of the program that reads it (unset, 0, empty and 1 each mean something different per variable): read the consumer's rule for each one, and test the channel it guards, not the assignment.
