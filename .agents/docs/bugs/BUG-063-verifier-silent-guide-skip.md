---
owner: bugs
date: 2026-09-17
id: BUG-063
status: fixed
area: agents/verifier
---
# BUG-063: verifier's log said `Scope: complete` while a required testing guide was never written

## Symptom
In repro-4 of `.agents/docs/testing/2026-09-17-bug-062-repro.md`, the verifier's combined call `mkdir -p .agents/docs/testing/server && … gofmt …; echo …` was denied by the run's tool allowlist (no `echo`). It never retried the `mkdir`, wrote no `testing/server/guide.md` although `agents/verifier.md` requires one whenever the guide is absent, and its log said `Scope: complete` without naming the missing guide.

## Repro
`.agents/docs/testing/2026-09-17-bug-062-repro.md` § "Separate defect seen in repro-4"; the denial is specific to that run's allowlist (no `echo`), not independently re-run with a different allowlist.

## Root cause
The guide duty was a paragraph after § What to run, outside the verifier's numbered steps, so a run could skip it without ever trying (repro-5: no attempt at all) or drop it after one denied command (repro-6). The Scope line's definition of "complete" was tied to tiers and items run, not to the guide's existence, so the log said complete either way. Before: 2 of 2 runs with the guide path denied logged `Scope: complete` without naming the guide (`.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md`, repro-5, repro-6).
`Diagnosis missed:` the explorer was not dispatched (location judged obvious); the main agent's hypothesis covered only the Scope wording, not the skipped duty itself (repro-5 never attempted the guide) — a brief asking "why was the duty not attempted" would have found it.

## Fix
`agents/verifier.md` step 4 makes the guide check a numbered step run before anything else, and a guide not written is named on the Scope line as `guide not written: <why>`, which is never `complete` (`agents/verifier.md:39`). Verified in `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md`: repro-7, 8, 10 (guide denied) named `guide not written`; repro-9, 11 (open) wrote the guide. Residual: repro-10's line opened "complete for the tiers …" before the `guide not written` clause — a wording slip, not a Scope-complete claim.

## Prevention
A duty an agent must perform is a numbered step it runs, not a paragraph beside the steps; "complete" is defined by the artifacts it owes, not only the commands it ran.

## Related
BUG-043 (the guide duty); BUG-062 (same repro run, same log format, different mechanism; ADR-032 removed the `<app>` resolution this fix's Scope check now covers); `.agents/docs/testing/2026-09-17-bug-062-repro.md`; `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md`.
BUG-064 (this fix's reporting half held in the real-home smoke check, but the guide duty itself was still skipped without any tool denial on Claude Code — third occurrence of this mechanism). In BUG-064's repro, repro-13 skipped the guide and still logged `Scope: complete` (naming the missing guide inside the line) — this bug's reporting rule is not reliable on Claude Code either; BUG-064's fix (the dispatch-text sentence) is the mitigation.
