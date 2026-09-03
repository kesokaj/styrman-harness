---
owner: bugs
date: 2026-09-03
id: BUG-003
status: fixed
area: context
---
# BUG-003: two scratch rules in the deployed context contradicted each other

## Symptom
In the first live sandbox test, Claude Code answered: "the context contradicts itself here: the Environment & Stack section says disposable work goes in the CLI scratchpad (/tmp/claude-<uid>/) … I'd follow the Hard-rules version (.agents/scratch/)."

## Repro
`sh-agent render | grep -n scratch` showed both the old bullet in `context/dev-setup.md` and the new hard rule in `context/AGENTS.md`.

## Root cause
The scratch rule was rewritten in `context/AGENTS.md` (ADR-008) while `context/dev-setup.md` kept its earlier bullet; the two files are concatenated at install and nobody grepped the rendered output.

## Fix
Removed the bullet from `context/dev-setup.md`; the hard rule is the only statement.

## Prevention
When a rule changes, run `sh-agent render | grep -n <keyword>` before deploying and read every hit; the deployed file is the concatenation, not one source file.

## Related
`.agents/docs/testing/2026-09-03-verify-sandbox.md`
