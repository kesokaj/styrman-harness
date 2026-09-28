---
owner: bugs
date: 2026-09-28
id: BUG-092
status: fixed
area: process / behavioural runs of agents/*.md
---
# BUG-092: behavioural probe's truncated fixture left `claude -p` running inside the harness checkout, one run wrote a fake change entry

## Symptom
During v1.7.0 gate fixes, a probe script extracted a shell function `demo()` from another script with `sed -n '/^demo()/,/^}/p'`. The sed range ended at the first `}` inside embedded Go code in the function's heredoc, truncating the sourced function. `d=$(demo …)` came back empty; `( cd "$d" && … )` with an empty `$d` is a no-op `cd`, so all four `claude -p` runs executed inside `styrman-harness` instead of the intended sandbox repo. One run (the control definition) wrote `.agents/docs/changes/2026-09-28-add-healthz-endpoint.md`, a fake change entry, into this repository.

## Repro
Run the probe's `demo()` extraction with `sed -n '/^demo()/,/^}/p' <source script>` where the function body contains a heredoc embedding Go code with its own closing brace; observe the extracted text stops at that brace, not the function's real end.

## Root cause
The probe script extracted shared fixture code with a line-range `sed` pattern instead of sourcing the file whole; a `}` inside an embedded heredoc satisfied the range's end pattern early. The truncated fixture was never asserted before use, so the resulting empty `$d` fed a bare `cd "$d"`, which silently no-ops rather than failing, leaving the CLI runs in the wrong directory. This happened minutes after the prove-the-setup rule (BUG-087/090/091) graduated into `AGENTS.md` § Testing; the new script did not carry the assertion the graduated rule requires, though the earlier `run.sh` had one.

## Fix
Fixture function moved to its own file, `.agents/scratch/tsn-probe/lib.sh`, sourced whole (no line-range extraction). Added `check_demo`, asserting own-repo root, the endpoint line, and the HEAD commit before every CLI start. Every `cd "$d"` changed to `cd "$d" || exit 1`. The assertion was validated: fires on an empty path and on a non-root directory, passes on a good fixture. The fake change entry was untracked, read, and deleted by name before any commit. Log: `.agents/docs/testing/2026-09-28-verify-v1.7.0-gate-fixes.md` § What went wrong.

## Prevention
A probe script never extracts shell code from another file with a line-range pattern (a `}` inside a heredoc ends the range early); shared fixture code lives in its own file and is sourced whole. Every `cd` into a fixture directory is `cd <dir> || exit 1`, preceded by the fixture assertion that the graduated `AGENTS.md` § Testing rule ("a behavioural or sandbox run proves its setup before its output counts") already requires.

## Related
BUG-087, BUG-090, BUG-091 — the same prove-the-setup rule graduated into `AGENTS.md` § Testing in a9717cc; this is a recurrence minutes after graduation, caused by a new script (not the earlier `run.sh`) that lacked the assertion. The rule's existing home already covers this; the gap was a new script not carrying it, not a missing rule — no new graduation target recommended beyond re-confirming the existing paragraph covers script rewrites, not only the scripts audited at graduation time.
