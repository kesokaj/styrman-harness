---
owner: bugs
date: 2026-09-16
id: BUG-044
status: fixed
area: agents/security (probe files outside scratch)
---
# BUG-044: security agent's crafted-input probes landed as untracked files in the repo root

## Symptom
After the `security` agent's review of the settings-write change
(`.agents/docs/reviews/2026-09-16-security-attribution-write.md`, "ran 20
crafted inputs against a binary built from the current blob"), `git status`
showed seven untracked files in the repository root whose *names* were the
probe contents themselves: `["x"]`, `null`, `{}`, `{"a":1}\n{"b":2}`, a
BOM-prefixed object, a CRLF object, `{"model":"x"}` — each ~100 bytes. The
verifier's second run (`.agents/docs/testing/2026-09-16-verify-attribution-write.md`)
flagged them; the main agent inspected and removed each by exact name.

## Repro
Not re-run. Inferred from the debris shape: the probe loop wrote each
crafted input's *content* into a path built from that same content (a
variable mix-up in a `printf`/redirect — content used where a filename
variable was meant), executed with the working directory at the repo root
instead of a scratch subdirectory.

## Root cause
The security agent improvised its probe harness ad hoc in the repo root
instead of writing fixed-name probe files under `.agents/scratch/<slug>/`.
`agents/security.md:77` already states "anything temporary you create goes
under the project's `.agents/scratch/` (gitignored), never `/tmp`" — the
rule existed and was not followed. Compounding it: nothing after the
dispatch checked the working tree for stray untracked files before the
review was accepted as done; the debris surfaced only on the verifier's
second pass.

## Fix
No code changed. Debris removed by exact name after inspection
(`git status --porcelain -z` listing, `os.stat` of each path, then
`os.remove`); this record documents the failure.

## Prevention
Before writing ad hoc probe inputs, a review or verification agent puts
them under `.agents/scratch/<slug>/` with names it chooses itself, never
derived from the input content and never in the working directory — this
is already the agent's own stated rule (`agents/security.md:77`,
mirrored in `agents/verifier.md`), so the gap here is that it was not
followed, not that it was missing.
Second rule, new: after any dispatch that ran probes or crafted-input
tests, the main agent runs `git status --short --untracked-files=all` and
treats an untracked path outside `.agents/docs/` as a stop sign — the same
posture the hard rules already require for an unexplained `D` on tree
changes.

## Related
BUG-014 (verifier wrote scratch to `/tmp` instead of `.agents/scratch/` —
same class of rule-not-followed, different agent); BUG-038 (agent's own
"Scratch"/"Never" list in its definition, same section format). Three
records now share "the agent's own scratch rule exists but wasn't
followed" (BUG-014, BUG-038's git-tree rule, BUG-044) — recommend a
standing post-dispatch check (`git status --short --untracked-files=all`,
stop on anything outside `.agents/docs/`) graduate into
`context/AGENTS.md`'s Definition of Done step 1 (verifier) and step 2
(review), since it is actionable (one grep-able command) and generalizable
(applies to any agent that writes temp files, not just security).
Logged in `.agents/docs/reviews/2026-09-16-security-attribution-write.md`
and `.agents/docs/testing/2026-09-16-verify-attribution-write.md`.
