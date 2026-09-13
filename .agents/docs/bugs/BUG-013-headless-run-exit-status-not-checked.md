---
owner: bugs
date: 2026-09-13
id: BUG-013
status: fixed
area: process/verification
---
# BUG-013: main agent cited a headless run's output tail as evidence without reading its own exit-status file first

## Symptom
The main agent cited the sandbox run `/review --project --all` (Claude Code headless,
clone `/home/kesokaj/sandbox/runs/claude-project`) as evidence that "five artifacts and
the merged record" were written, and passed that claim to the changelog agent. The
`verifier` agent (`.agents/docs/testing/2026-09-13-verify-review-skill.md`, Status FAIL)
read `.agents/scratch/e2e-review/status` and found `claude-project exit 124` — the run had
hit the 580-second `timeout` wrapper. Its actual output called the merged record a draft
with the security lens still running and the explorer's brief missing.

## Repro
Launch a multi-subagent headless run under `timeout 580`; read only the tail of its
stdout, never the status line it wrote itself.

## Root cause
The main agent inspected the output file and the artifact listing but not the exit status
it had itself recorded in the status file — the very instrument it set up. The global
hard rule on instrumentation silence was applied to the agents' logs but not to the
wrapper's exit code, and the 580-second budget was too short for six subagents including
the security agent on the strongest model.

## Fix
Re-launched detached from the tool timeout with a 1500-second budget and an explicit
status file; evidence is now cited only after that status line reads exit 0 and the
merged record has no pending row; `verifier` re-dispatched afterwards.

## Prevention
Every headless run cited as evidence is read from its exit status first (`<name> exit N`),
then its artifacts; a non-zero or missing status line is a FAIL of the evidence, whatever
the output tail says. A run that fans out to several subagents gets a budget beyond the
tool's ten-minute limit by running detached (`setsid nohup … &`) with its own status file.

## Related
BUG-011 (a status not read before a commit); the global hard rule on instrumentation
silence — this is the second bug (with BUG-011) where a written status/exit signal
existed but was not read before acting; a project rule requiring "read the status file
before citing a run as evidence" may belong in AGENTS.md.
