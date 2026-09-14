---
owner: bugs
date: 2026-09-14
id: BUG-014
status: fixed
area: agents/verifier
---
# BUG-014: verifier installed its scratch home under /tmp instead of .agents/scratch/

## Symptom
Verifying the BUG-012 fix, the `verifier` agent installed its scratch home at
`/tmp/sh-agent-verify-bug012` (`.agents/docs/testing/2026-09-14-verify-bug-012-fix.md`),
though the global hard rule says scratch is `.agents/scratch/` in the project, never
`/tmp`, and the repo's `AGENTS.md` repeats it for scratch homes.

## Repro
Dispatch the verifier for a change that needs a scratch install; its testing log
names the `/tmp` path it used.

## Root cause
The global rule speaks to "every generated or temporary file"; `agents/verifier.md`
listed what commands to run without naming where its own scratch install goes. An
agent reading its own definition first takes that silence as freedom to use `/tmp`
by habit — same shape as BUG-007, where repo docs said `/tmp` while the context
already said never.

## Fix
`agents/verifier.md`: added one sentence — every scratch home, build directory, or
temporary file goes under the project's `.agents/scratch/`, never `/tmp`; the global
scratch rule applies to the agent too.

## Prevention
An agent definition must restate the global rules that govern its own side effects
(where it writes, what it never touches) — a global rule alone does not reach a
subagent reliably. Before writing or editing an `agents/*.md` file, check it names
its own scratch/output location explicitly rather than relying on the global context
being read.

## Related
BUG-007, BUG-003 — same "/tmp instead of .agents/scratch/" mechanism, now three
occurrences (BUG-003 in context docs, BUG-007 in repo docs, BUG-014 in an agent
definition). This prevention rule should graduate to a standing check: every
`agents/*.md` and `docs/*.md` file that runs commands states its own scratch path,
verified by a repo-wide grep for `/tmp` as part of `/init-project` or a pre-commit
review checklist.
