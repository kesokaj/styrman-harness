---
owner: main
date: 2026-09-30
---
# Verification — parity plan Phase 2: memory on both CLIs, three setups   (2026-09-30)

Range reviewed: ae8f965..f554da8 (Phase 2's text), with the binary built from f554da8. Plan `plans/2026-09-30-parity-memory-dream.md` Phase 2 checkpoint: E2E 1 on the three setups, E2E 2 on the claude setup in both states, and the `dreamer` behavioural run.

## Status: PASS — E2E 1 on all three setups and across CLIs, E2E 2's control and change, `dreamer` 4 of 5 on the matrix check (5 of 5 on the store)

## Kit
All paths are outside every checkout. `git -C ~/sandbox/p2/runs rev-parse` fails; the run script asserts it.
- **Fixture**: `~/sandbox/p2/fixture`, the gate16 `greetbox` module plus a commit adding `/init-project`'s ignore lines (`9e4b88e`). Each run clone asserts `git check-ignore -q .agents/memory/x.md` before the CLI starts.
- **Homes**: `.agents/scratch/p2-{claude,agy,all}`, installed from f554da8 with `--vendor claude|agy|all` (135, 134 and 268 files as `install`'s "written" line counts them, the binary included; the gate verifier counted 134, 133 and 267 managed files beside it). Each run asserts its manifest's vendor. The control home `.agents/scratch/p2-claudectl` was installed from a worktree at v1.16.2 (the pre-change harness: no `project-memory` skill, old feedback-log rule 7).
- **Environment**: `HOME=<home>` for every run, so each `sh-agent` call the agent makes is pinned there (BUG-107). The real `~/.config/sh-agent` listing was hashed before and after: unchanged.
- **Claude Code 2.1.285**: `CLAUDE_CONFIG_DIR=<home>/.claude` with a copied `.credentials.json`, chmod 600. The probe "Reply with the single word ok." answered `ok`, so the plan's assumption holds. Memory settings went through `--settings` (this session cannot write `.claude/settings.local.json`).
- **Antigravity 1.2.14**: token copied chmod 600; the settings file was written for the run with named grants, no `toolPermission` (`grep -c` 0), and `trustedWorkspaces` naming only the run clones. No skip flag was needed.
- Every run's stderr was captured (`2> <run>.err`) and holds nothing but the kit's own `rc=… secs=…` line, so 0 `auto-denied` lines. The same capture recorded agy's `auto-denied` line in the gate17 runs, so it can emit one. The copied credentials were deleted by name after the batch.

## E2E 1 — memory parity

The prompt is "Remember that we deploy on Fridays. Save it the way this project saves memories, then reply with one line." Recall is a fresh session: "When do we deploy in this project? Answer with the weekday only."

| Setup | Remember (CLI → file, index line) | Recall (CLI → answer) |
|---|---|---|
| claude | cc → `deploy-on-fridays.md`, `- [deploy-on-fridays](deploy-on-fridays.md) — deploys happen on Fridays` (14 s) | cc → `Friday` (12 s) |
| agy | agy → `friday-deploys.md`, its line (24 s) | agy → `Friday` (23 s) |
| all, clone m1 | cc → `deploy-fridays.md` (9 s) | **agy** → `Friday` (11 s) |
| all, clone m2 | agy → `deploy-on-fridays.md` (29 s) | **cc** → `Friday` (9 s) |

Every index line is in the verb's slug shape. The fixture holds the fact nowhere but the store. Transcripts show the calls:
- Claude Code (sandbox `projects/*.jsonl`): `sh-agent memory add --type project …` and `sh-agent memory show` (at session start, in both setups).
- Antigravity conversation stores: `sh-agent memory add` and `sh-agent memory show` in both the agy and the all home.

## E2E 2 — Claude Code's native writes are off

The prompt is "Remember that I prefer tabs over spaces in every file."

| State | Home | `--settings` | Result |
|---|---|---|---|
| must-fail control | `p2-claudectl` (v1.16.2 harness) | `autoMemoryDirectory` only | `indent-with-tabs.md` **with `originSessionId`**, index line titled "Indent with tabs": Claude Code's native writer. The control produced its condition. |
| the change | `p2-claude` (f554da8) | both keys, `autoMemoryEnabled: false` | `tabs-over-spaces.md`, no `originSessionId` in any file; the verb's slug line |

## `dreamer` behavioural run (AGENTS.md:65)

`claude -p --setting-sources project,local --agents <json built from the working-tree agents/dreamer.md>` in the claude home, with the matrix row's contract on stdin (`Store: <clone>/.agents/memory … Promoted: vet-before-tag-1.md, vet-before-tag-2.md`). Each run used a fresh clone seeded by `~/sandbox/p2/dreamer/fixture.sh`, which asserts:
- two `bugs/INDEX.md` rows for one defect with the prevention rule written nowhere ("second occurrence");
- a five-step release in a handoff ("same five steps") and a `changes/` entry;
- two `feedback` memories from two session ids, and one open `project` memory;
- no "before tag" rule in the fixture or the global context.

| Run | Artifact COMPLETE | Store read at `<clone>/.agents/memory` (Scope line) | Both promoted files | Bug ids cited | skill/procedure candidate | `open-migration` left alone | Tree clean | Matrix check |
|---|---|---|---|---|---|---|---|---|
| r1 | yes | yes | yes | yes | yes | yes | yes | pass |
| r2 | yes | yes | yes | yes | yes | yes | yes | pass |
| r3 | yes | yes | yes | no: cites `bugs/INDEX.md:<line>` and "two bug rows" instead of the ids | yes | yes | yes | fail on the id needle |
| r4 | yes | yes | yes | yes | yes | yes | yes | pass |
| r5 | yes | yes | yes | yes | yes | yes | yes | pass |

One more variance, not in the check: the Scope line counts 4 memory files in r1 and r4 but 3 in r2, r3 and r5, over identical stores of 3 memories plus `MEMORY.md`. r1 and r4 counted the index as a memory. The dreamer's counting is unchanged by this phase.

The result is 4 of 5 on the matrix check as written, and 5 of 5 on what this change touched: the store path, which every run read at `<root>/.agents/memory`. r3 used the bug rows as evidence but cited them by line. The needle measures the citation form, which the 2026-09-29 definition did not change, and this change did not touch.

## Not exercised
- `/init-project` writing `autoMemoryEnabled: false` into a real settings file. Headless Claude Code refuses writes under `.claude/`; the key's effect is what E2E 2 shows.
- The rollout hint in a live session: its unit test covers it (`TestShowRolloutHint`).
- Antigravity's handling of `/project-memory` typed as a slash command. The skill loaded by description in the agy runs (the verb calls follow its text), but the command form was not typed.
