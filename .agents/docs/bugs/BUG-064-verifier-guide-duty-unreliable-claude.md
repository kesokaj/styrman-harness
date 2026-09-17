---
owner: bugs
date: 2026-09-17
id: BUG-064
status: fixed
area: agents/verifier
---
# BUG-064: verifier's guide duty still not reliably performed on Claude Code after BUG-063's fix

## Symptom
Real-home smoke check at harness e67eace: Claude Code's `verifier`, dispatched in a fresh demo repo with no testing guide and an allowlist including Write and mkdir, never attempted `testing/guide.md` — no call, no denial in the transcript — yet reported it honestly (`guide not written: … Scope is not complete`). Antigravity, same check, same repo shape, wrote the guide. Earlier the same day, `claude -p --agents` repro-11 wrote it.

## Repro
`.agents/docs/testing/2026-09-17-verify-real-home-smoke-2.md`; compare against `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md` repro-11 (guide written, same definition, same day).

## Root cause
As observed, not proven: `agents/verifier.md` step 4 (BUG-063's fix) makes the guide check a numbered step run before anything, and the Scope line correctly reports its absence when skipped — the reporting half holds. But the duty half (actually attempting the step) was inconsistent on Claude Code across otherwise-similar runs on the same day with the same agent definition: no tool denial or error explains a skipped run. `.agents/docs/testing/2026-09-17-bug-064-repro.md` shows the same instruction held in every run once placed in the dispatch prompt instead of relying on the subagent's own definition alone — rules inside a subagent's own file are followed intermittently on Claude Code, while instructions in the dispatch prompt were followed in every run today. This is the second time (BUG-043, then BUG-063) a written procedural rule inside `agents/verifier.md` alone had not made Claude Code's verifier perform the guide duty reliably.

## Fix
The operator's chosen design (asked as a short list): the duty moves into the dispatch text the main agent writes, not only the subagent's own definition. `context/AGENTS.md` Definition of Done step 1 now says that when `<docs-root>/testing/guide.md` does not exist, the verifier dispatch must say "testing/guide.md is absent: write it first"; `skills/review/SKILL.md`'s `--tests` row and `skills/preflight/SKILL.md`'s Verification row carry the same sentence. Verified in `.agents/docs/testing/2026-09-17-bug-064-repro.md`: without the sentence the guide was written in 2 of 3 runs (repro-13/14/15); with it, the main agent carried the sentence into the dispatch in 3 of 3 runs and the guide was written 3 of 3 (repro-16/17/18).

## Prevention
A duty a subagent must not skip goes in the dispatch text the main agent writes, not only in the subagent's own definition: rules inside `agents/verifier.md` itself were followed intermittently on Claude Code, while the same instruction placed in the dispatch prompt held in every run today.

## Related
BUG-063 (fixed the reporting half, same mechanism, note appended there); BUG-043 (the original guide duty, same recurring failure to execute on Claude Code); `.agents/docs/testing/2026-09-17-verify-real-home-smoke-2.md`; `.agents/docs/testing/2026-09-17-bug-063-adr-032-repro.md`; `.agents/docs/testing/2026-09-17-bug-064-repro.md`.
