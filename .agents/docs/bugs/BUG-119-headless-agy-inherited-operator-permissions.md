---
owner: bugs
date: 2026-09-30
id: BUG-119
status: fixed
area: internal/dream (agy and Claude child CLIs)
---
# BUG-119: the headless agy child inherited the operator's permissions

## Symptom
Security round 1, finding 1 (`.agents/docs/reviews/2026-09-30-security-cb799ff-f1ca525.md`): on this box agy's settings say toolPermission always-proceed, so an instruction injected into a transcript could have run a command at 03:00 unattended; only the prompt's warning stood in the way.

## Repro
Run the nightly dream with the operator's real settings and a transcript containing an instruction to run a command.

## Root cause
The child CLI ran under the operator's own HOME and configuration while reading untrusted text; safety rested on a prompt warning, which is not a control. The D7b fix (BUG-116) left this in place.

## Fix
`internal/dream/run.go` `agyHome`: a private HOME with a copy of the login and settings granting nothing, removed after the run. Verifier mutation (grant always-proceed) failed TestRunAgyPrivateHomeAndStdin; a real sandbox run with always-proceed in the operator's settings produced 2 proposals and never saw the grant. Claude Code gained `--strict-mcp-config --setting-sources local`; the Antigravity reader stopped sending tool output (step_type 14/15 only).

## Prevention
An unattended child CLI fed untrusted text runs with a configuration that grants nothing, built by the caller, never with the operator's own.

## Related
BUG-116, BUG-115 (sandbox settings lesson), BUG-117, BUG-113 (skills/dream).
