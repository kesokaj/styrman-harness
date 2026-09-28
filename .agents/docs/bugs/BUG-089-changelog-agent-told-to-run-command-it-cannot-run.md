---
owner: bugs
date: 2026-09-28
id: BUG-089
status: fixed
area: agents/changelog
---
# BUG-089: changelog agent's date rule named a command it has no Bash to run

## Symptom
`agents/changelog.md` (pre-`e5a49d4`) stated date rule 2 as "If absent, run `date +%F` yourself." The agent's `tools:` frontmatter is `Read, Grep, Glob, Edit, Write` — no Bash — so the instruction could never execute; the agent had no way to satisfy its own rule.

## Repro
Read `agents/changelog.md` at any commit before `e5a49d4` and compare its date rule 2 against its `tools:` list: the rule names a shell command, the list has no Bash. Found in cold read round 3 of `plans/2026-09-28-time-sortable-trail-names.md` (`.agents/docs/reviews/2026-09-28-coldread-time-sortable-names-r3.md`); round 1 flagged the same shape for the new rule the plan first proposed (`.agents/docs/reviews/2026-09-28-coldread-time-sortable-names-r1.md`, finding 1).

## Root cause
The rule was written against what a human or a Bash-capable agent could do, not against the actual `tools:` grant of this agent. Nothing cross-checked an instruction's verbs (run, execute, `date`, `git`) against the tool list of the agent it was written into.

## Fix
`e5a49d4` (plan D3/D4): the dispatcher now passes the stamp computed by `date +%F-%H%M` in the dispatch text (documented in the agent's `description:` and in `skills/preflight/SKILL.md`); the Bash-less "run it yourself" rule is deleted from `agents/changelog.md`.

## Prevention
An instruction in an agent or skill that names a command is checked against that agent's `tools:` list before it ships; an agent without Bash receives the command's output in its dispatch text instead of an instruction to run it.

## Related
Round 1 finding in `.agents/docs/reviews/2026-09-28-coldread-time-sortable-names-r1.md` (same shape, caught earlier in the same plan). BUG-090 (this plan's Phase 2 fixture defect) is a separate mechanism from the same plan.
