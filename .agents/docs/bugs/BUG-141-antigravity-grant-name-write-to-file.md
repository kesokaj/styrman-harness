---
owner: bugs
date: 2026-10-01
id: BUG-141
status: fixed
area: testing (skill-matrix Antigravity kit)
---
# BUG-141: headless Antigravity rows died "auto-denied" for a shell command the per-name grant list could not express

## Symptom
Phase 9 Antigravity skill-matrix batch (2026-10-01, `.agents/scratch/audit-fix/kit-agy-matrix.sh`, sandbox home `~/sandbox/audit-fix/sbAgy`): rows `bash`, `mermaid`, `docker`, `init-deploy`, `init-update`, `handoff`, `cloud-run` (earlier `spec`, `requirements`, `review-resume`) ended with `jetski: no output produced — a tool required the "command" permission that headless mode cannot prompt for, so it was auto-denied`. Three rounds widening `command(<name>)` grants (the documented list, then `cd`, `sed`, `tee`) changed nothing. The message names neither the tool nor the command.

## Repro
Sandbox `settings.json` with per-name `command(<name>)` grants and `write_file(*)`; run a row whose skill issues `node`, `docker`, `gcloud`, or a `for d in …; do` loop through `agy -p`: no output, the auto-denied line on stderr.

## Root cause
The denied thing was a shell command every time, never a file write. `write_file(*)` was correct all along (`write_to_file` is only the tool's display name; an earlier draft of this record blamed it wrongly). The grant store accepts four actions only: `command`, `read_file`, `write_file`, `mcp`; `cli.log` showed `ignoring invalid allow entry "write_to_file(*)": unknown action "write_to_file"`. `command(<name>)` matches one command name, so a name outside the list or a construct (loop, heredoc) is denied. The last `"CommandLine"` in each row's conversation database: `node` (mermaid, rendering with mmdc), `docker`, `gcloud`, a `for d in …; do` loop (init-deploy and init-update; it is the init-project skill's own Step 5b text, which `command(<name>)` cannot express), and a multi-line `grep -rn` / heredoc (bash row). Three rounds guessed at names because the denial names no command; reading the database settled it.

## Fix
The kit runs those rows with `command(*)` in their sandbox clone; `testing/skill-matrix.md` step 2 and `testing/guide.md` § 4 now name the four actions and the database diagnosis. Round 5 re-runs the rows.

## Prevention
A headless Antigravity row that dies "auto-denied" is diagnosed from its conversation database (last `CommandLine`) and `cli.log` (`invalid allow entry`) BEFORE any grant is widened. A grant is one of the four actions the store knows; a skill whose text holds a shell construct (loop, heredoc) is run with `command(*)` in a sandbox clone, never with a longer name list.

## Related
BUG-091, BUG-107, BUG-116 (earlier Antigravity sandbox setup bugs); the `agy` version in `testing/2026-09-03-verify-sandbox.md`. No existing rule covered this: the documented per-name list was the rule and could not express the commands.
