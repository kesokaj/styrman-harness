---
owner: bugs
date: 2026-09-15
id: BUG-031
status: fixed
area: context (headless runs)
---
# BUG-031: a print-mode run ends its turn while a dispatched agent is still running

## Symptom
`documentation` matrix row on Claude Code (`aL1`, long-limit diagnostic, 3600 s cap): `claude -p`
exited 0 at 1072 s and the row's check FAILED. The main turn ended while a background `security`
agent was still running; its last output was "Nothing more can go ahead until security reports
back: code fixes wait for security; docs wait for the code fixes." It had dispatched verifier,
reviewer, security and a docs-lens reviewer, but README.md and `deploy/local.md` were never updated,
no changes entry was written, nothing was committed.

## Repro
Evidence in `.agents/scratch/phase-notes.txt` and `.agents/scratch/phase6-kit/results/`. On
Antigravity, the `handoff` row a1 printed exactly "And I will wait for the subagent to report back."
and `agy -p` does wait ("root agent idle; waiting for 1 background task(s) (bounded by
--print-timeout)"), ending at its 20-minute print timeout instead — a different failure mode, not
reproduced as a false-pass. Claude Code's `documentation` row passed on its matrix re-run (a2, 1253
s), so this did not fail twice.

## Root cause
`claude -p` ends its turn once its own reasoning is done, without waiting for a still-running
background agent to return. A search of `context/`, `agents/` and `skills/` (references excluded) for
background/in-flight/print-mode/`claude -p`/`agy -p`/wait-for-an-agent wording finds only user gates
(`skills/bug/SKILL.md:15`, `skills/feature/SKILL.md:15`, `skills/review/SKILL.md:60`) and handoff's
"before ending" line, present even in the pre-plan c192622 copy. No harness text tells a print-mode
run to collect dispatched agents' results before ending its turn; the gap predates this plan.

## Fix
Fixed 2026-09-16 (Phases 2-3, clear the board): `context/AGENTS.md` now
states a headless run waits for every dispatched agent to return before
ending its turn, and says so if one has not.

## Prevention
A headless run waits for every agent it dispatched to return before it ends
its turn, and says so in its output when one has not. Graduated:
`context/AGENTS.md`, next to "Unattended runs on Claude Code set a `/goal`
naming the check that must pass."

## Related
BUG-013 (headless run's own status must be read before trusting a green run). This is the third
record whose prevention is "validate a check/run in the state it must fail as well as the state it
must pass" family only tangentially — its own rule (wait for dispatched agents before ending a
print-mode turn) is distinct and not yet written anywhere in the harness; recommend it land in
`context/AGENTS.md`.
