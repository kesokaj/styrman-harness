---
owner: bugs
date: 2026-09-17
id: BUG-063
status: open
area: agents/verifier
---
# BUG-063: verifier's log said `Scope: complete` while a required testing guide was never written

## Symptom
In repro-4 of `.agents/docs/testing/2026-09-17-bug-062-repro.md`, the verifier's combined call `mkdir -p .agents/docs/testing/server && … gofmt …; echo …` was denied by the run's tool allowlist (no `echo`). It never retried the `mkdir`, wrote no `testing/server/guide.md` although `agents/verifier.md` requires one whenever the guide is absent, and its log said `Scope: complete` without naming the missing guide.

## Repro
`.agents/docs/testing/2026-09-17-bug-062-repro.md` § "Separate defect seen in repro-4"; the denial is specific to that run's allowlist (no `echo`), not independently re-run with a different allowlist.

## Root cause
Not diagnosed beyond what is observed: a chained shell call failed partway, the agent did not retry the required step, and its own log did not check its claimed scope against what the required-artifact rule (`agents/verifier.md` guide duty) actually demands.

## Fix
Not fixed.

## Prevention
A required artifact that was not written is named under the log's Scope as not done, so a skipped duty is never reported complete.

## Related
BUG-043 (the guide duty); BUG-062 (same repro run, same log format, different mechanism); `.agents/docs/testing/2026-09-17-bug-062-repro.md`.
