---
owner: main
date: 2026-09-30
---
# Verification — real home before v1.17.0, both CLIs   (2026-09-30)

Range reviewed: v1.16.2..a275b96, installed into the real home with `sh-agent install` at 21:44. `update` refused the tree, because the agents' untracked review artifacts under `.agents/docs` made it dirty; no source differed from HEAD. The install wrote 17 files, left 251 unchanged, `drift none`. `skills/{project-memory,feedback-log,init-project,dream,retro,distill}/SKILL.md` are byte-identical to both vendors' installed copies. Kit: `~/sandbox/gate19/run.sh`, fresh clones outside every checkout, every `sh-agent` call pinned to `~/sandbox/gate19/home` (BUG-107).

## Status: PASS — the context on both CLIs; `project-memory` on both; `init-project`, `dream`, `retro`, `distill` (with the `dreamer` agent it dispatches) on Antigravity from the real home; `feedback-log` on Antigravity in the agy-only sandbox. One kit defect, found and cleaned: see below.

The rule is project `AGENTS.md` § Release. The range changes six skills, one agent and the context.

| CLI | Check | Result |
|---|---|---|
| Claude Code 2.1.285 | context | PASS: "# AGENTS.md — global agent context", 29 skills, 8 subagents, GANNET, `bug`; 6 s |
| Claude Code | `project-memory` ("Remember that we deploy on Fridays…") | PASS: `.agents/memory/MEMORY.md` `- [deploy-on-fridays](deploy-on-fridays.md) — releases go out on Fridays`, the verb's slug line, no `originSessionId`; 23 s |
| Antigravity 1.2.14 | context | PASS: the same four answers; 8 s |
| Antigravity | `project-memory` | PASS: `- [deploy-on-fridays](deploy-on-fridays.md) — deployments happen on Fridays`; 33 s |
| Antigravity | `/init-project` (older scaffold, vendor all) | PASS: `.claude/{skills,agents,commands,rules}` links; `.claude/settings.local.json` = `{"autoMemoryDirectory": "<clone>/.agents/memory", "autoMemoryEnabled": false}`; `TODO.md`; commit `239d1da init: update scaffold…`; 315 s |
| Antigravity | `/dream --dry-run` | PASS: `store: <clone>/.agents/memory (1 memory files; default path …: 0)`; 181 s |
| Antigravity | `/retro --dry-run` | PASS: findings table empty, "Project scope: 0, Harness scope: 0"; 138 s |
| Antigravity | `/distill` on the dreamer fixture | PASS: `retros/2026-09-30-dreamer.md` (the `dreamer` dispatched) and `retros/2026-09-30-distill-release-by-hand.md`; the skill left proposed as told; 262 s |
| Antigravity (agy-only sandbox, `.agents/scratch/p2-agy` reinstalled from this tree) | `feedback-log` ("No, that is wrong: … tabs …") | PASS on the second attempt: the sandbox `feedback.jsonl` 0 → 1 line; `.agents/memory/tabs-for-indentation.md` and its slug line written by `sh-agent memory add`; no `originSessionId`. The first attempt, with named grants only, was auto-denied on a compound command (BUG-115's kind); the second used `--dangerously-skip-permissions` in the throwaway clone. Token and settings deleted after |

No run's stderr has an `auto-denied` line, except the first sandbox feedback attempt named above.

## Kit defect, found and cleaned
- **What:** the real-home `/distill` run appended one line to the REAL `~/.config/sh-agent/feedback.jsonl` and one to `retros.jsonl`, both naming `~/sandbox/gate19/runs/agy-distill-1`. Distill writes those files with file tools, not through `sh-agent`, so the kit's `--home` wrapper cannot pin them.
- **Found:** the before/after listing of `~/.config/sh-agent`.
- **Cleaned:** exactly those two lines removed (feedback 16 → 15, retros 17 → 16; `grep -c sandbox/gate19` 0 in both). Recorded as a bug.

## What was not exercised
- The real nightly timer (`sh-agent dream --install-timer` on this box): the plan's `who: user` step, in the handoff.
- The Mac and devbox places.
