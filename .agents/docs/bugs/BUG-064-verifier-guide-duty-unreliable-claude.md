---
owner: bugs
date: 2026-09-17
id: BUG-064
status: open
area: agents/verifier
---
# BUG-064: verifier's guide duty still not reliably performed on Claude Code after BUG-063's fix

## Symptom
Real-home smoke check at harness e67eace: Claude Code's `verifier`, dispatched in a fresh demo repo with no testing guide and an allowlist including Write and mkdir, never attempted `testing/guide.md` — no call, no denial in the transcript — yet reported it honestly (`guide not written: … Scope is not complete`). Antigravity, same check, same repo shape, wrote the guide. Earlier the same day, `claude -p --agents` repro-11 wrote it.

## Repro
`.agents/docs/testing/2026-09-17-verify-real-home-smoke-2.md`; compare against `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md` repro-11 (guide written, same definition, same day).

## Root cause
Not established. `agents/verifier.md` step 4 (BUG-063's fix) makes the guide check a numbered step run before anything, and the Scope line correctly reports its absence when skipped — the reporting half holds. But the duty half (actually attempting the step) is inconsistent on Claude Code across otherwise-similar runs on the same day with the same agent definition: no tool denial or error explains this run's skip. This is the second time (BUG-043, then BUG-063) a written procedural rule has not made Claude Code's verifier perform the guide duty reliably.

## Fix
Not fixed. No wording change is proposed here.

## Prevention
None established. Two prose-based fixes (BUG-043, BUG-063) have not produced reliable execution of this one step on Claude Code; a third wording change is not recommended. This is a design question for the operator — e.g. whether the guide should be written by a different owner, at a different moment, or checked post-hoc rather than trusted from the log — not another rewording of step 4.

## Related
BUG-063 (fixed the reporting half, same mechanism); BUG-043 (the original guide duty, same recurring failure to execute on Claude Code); `.agents/docs/testing/2026-09-17-verify-real-home-smoke-2.md`; `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md`.
