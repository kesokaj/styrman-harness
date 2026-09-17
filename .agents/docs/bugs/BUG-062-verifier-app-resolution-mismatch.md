---
owner: bugs
date: 2026-09-17
id: BUG-062
status: open
area: agents/verifier
---
# BUG-062: Claude Code's verifier resolved `<app>` to the repo directory name, Antigravity to `cmd/<app>`

## Symptom
In the demo repository under `.agents/scratch/realhome-smoke/demo` (single binary `cmd/server`, no existing testing guide), Claude Code's `verifier` wrote its guide to `.agents/docs/testing/demo/guide.md`. Antigravity's `verifier`, dispatched minutes later in the same repository, wrote `testing/server/guide.md` — matching G20's first branch. The 2026-09-16 log `.agents/docs/testing/2026-09-16-verify-antigravity-real-home-complete.md` shows Antigravity also resolving `server` there.

## Repro
Dispatch `verifier` (Claude Code) in a fresh repo with exactly one `cmd/<app>` directory and no `AGENTS.md`-named binary override, and observe which name it writes its guide under vs. `ls -d cmd/*/`. Not independently re-run here; observed once in the 2026-09-17 real-home smoke check (`.agents/docs/testing/2026-09-17-verify-real-home-smoke.md`, § Defect observed).

## Root cause
Not diagnosed. `agents/verifier.md:37` states an ordered resolution (`cmd/<app>` when exactly one → AGENTS.md-named binary → repo directory name), and only Claude Code's run skipped straight to the third branch despite the first branch's precondition being met. What differs: the demo's `AGENTS.md` heading is `# AGENTS.md — demo`, which may have been read as naming the project/binary rather than as a title — untested.

## Fix
Not fixed.

## Prevention
A rule with an ordered resolution states the check to run for its first branch (e.g. `ls -d cmd/*/`) rather than describing only the outcome, so an agent cannot silently fall through; a verifier log names which branch it used, making a wrong resolution visible without re-running the dispatch.

## Related
BUG-043 (agents/verifier § What to run graduated); `.agents/docs/testing/2026-09-17-verify-real-home-smoke.md`; `.agents/docs/testing/2026-09-16-verify-antigravity-real-home-complete.md`.
