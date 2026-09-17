---
owner: bugs
date: 2026-09-17
id: BUG-062
status: fixed
area: agents/verifier
---
# BUG-062: Claude Code's verifier resolved `<app>` to the repo directory name, Antigravity to `cmd/<app>`

## Symptom
In the demo repository under `.agents/scratch/realhome-smoke/demo` (single binary `cmd/server`, no existing testing guide), Claude Code's `verifier` wrote its guide to `.agents/docs/testing/demo/guide.md`. Antigravity's `verifier`, dispatched minutes later in the same repository, wrote `testing/server/guide.md` — matching G20's first branch. The 2026-09-16 log `.agents/docs/testing/2026-09-16-verify-antigravity-real-home-complete.md` shows Antigravity also resolving `server` there.

## Repro
Dispatch `verifier` (Claude Code) in a fresh repo with exactly one `cmd/<app>` directory and no `AGENTS.md`-named binary override, and observe which name it writes its guide under vs. `ls -d cmd/*/`. Not independently re-run here; observed once in the 2026-09-17 real-home smoke check (`.agents/docs/testing/2026-09-17-verify-real-home-smoke.md`, § Defect observed).

## Root cause
`agents/verifier.md:37` stated the resolution order (`cmd/<app>` when exactly one → AGENTS.md-named binary → repo directory name) but named no check for its first branch. In 3 of 3 runs before the fix, Claude Code's verifier never ran `ls -d cmd/*/` (confirmed by its log and tool trace) and fell straight to the third branch, naming the guide after the repo/`AGENTS.md` heading (`demo`) instead of `cmd/server`. The explorer's independent recon (`.agents/docs/plans/2026-09-17-recon-bug-062.md`) ranked the same two candidates highest — the `AGENTS.md` heading (`# AGENTS.md — demo`) anchoring the resolution, and the rule's missing explicit pre-flight check — agreeing with the main agent's hypothesis; install/derive corruption was ruled out (installed text byte-identical to source on both vendors).

Observation (repro-1, not part of the root cause): the verifier attempted `mkdir -p` on an absolute path inside the harness repository instead of the demo repo; the permission system denied it and the verifier retried inside the demo repo, so nothing was created there (checked).

## Fix
`agents/verifier.md` step 2 now runs `ls -d cmd/*/` first and states an `AGENTS.md` heading does not override it; the log template gains an `App: <name> (<branch>)` line; `agents/documentation.md` restates the resolution. Uncommitted, lands in the fix commit. Verified in `.agents/docs/testing/2026-09-17-bug-062-repro.md`: 2 of 2 runs with the working-tree definition ran `ls -d cmd/*/` and logged `App: server (cmd/server)`, writing the guide to `testing/server/guide.md`.

Operator's choice at the third-round question (recorded by the explorer as a design question, since this is the third failure on the `<app>`-resolution mechanism after BUG-043): "Patch: make it a command" over deleting the `<app>` level.

## Prevention
An ordered resolution rule names the check for its first branch, not only the outcome, so an agent cannot silently fall through; the artifact records which branch fired, making a wrong resolution visible without a repro.

## Related
BUG-043 (agents/verifier § What to run graduated); BUG-063 (same repro run, separate missing-guide defect); `.agents/docs/testing/2026-09-17-verify-real-home-smoke.md`; `.agents/docs/testing/2026-09-16-verify-antigravity-real-home-complete.md`; `.agents/docs/testing/2026-09-17-bug-062-repro.md`; `.agents/docs/plans/2026-09-17-recon-bug-062.md`.
