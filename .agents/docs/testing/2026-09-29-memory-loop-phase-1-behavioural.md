---
owner: main
date: 2026-09-29
---
# Behavioural runs — memory loop, Phase 1 (the store)   (2026-09-29)

Range reviewed: working tree over f9b29e0 — `skills/init-project/SKILL.md` Step 5b, `context/AGENTS.md` resume clause, `skills/feedback-log/SKILL.md` rule 7, `skills/dream/SKILL.md` Step 1 counts and split stop. Plan `plans/2026-09-29-dreamer-memory-loop.md`.

## Status: PASS — every Phase 1 expectation met, with one vendor limit recorded (the headless `.claude/` write guard) and six void runs explained

## Kits

`.agents/scratch/memory-probe/`: `fixture/` (an older scaffold: no `retros/`, no `.claude/settings.local.json`, `.claude/skills/init-project` linked to the working tree), `fixture-absent/` (the same plus a committed `{}` settings file), `mem-template/` (two memories and an index, copied into the default per-project path under the real home before a seeded run), `run-init.sh`, `run-remember.sh`, `run-split.sh`; the split fixture is `dream-probe/fixture` with the working-tree `dream` linked. All runs `claude -p`, real home for credentials, prompt on stdin, `--permission-mode acceptEdits`, read-only Bash prefixes plus `mv`, `mkdir`, `cp`, `ln`, `python3` for the init runs. Every seed asserted by the kit before the CLI started (slug computed by the same rule the skill states; the default path emptied first).

## `/init-project` update path with a seeded default store (i1–i3)

| Run | Turns | Store created at `~/.config/sh-agent/memory/<slug>/` | Both files and the index moved, default path empty | `.gitignore` gains `.claude/settings.local.json` | `retros/` created | `.claude/settings.local.json` written | Exact content reported instead |
|---|---|---|---|---|---|---|---|
| i1 | 44 | yes | yes | yes | yes | no (headless guard) | yes |
| i2 | 48 | yes | yes | yes | yes | no | yes |
| i3 | 49 | yes | yes | yes | yes | no | yes |

Claude Code refuses writes under the project's `.claude/` in a headless session (`testing/2026-09-03-verify-sandbox.md` § Vendor behaviour); the skill's Step 5 guard, extended to Step 5b, fired in 3 of 3: the run reported the one-line file for the user. The store, the migration and the `.gitignore` line do not depend on it. Interactive runs are one approval; the E2E in this checkout (Phase 2) is where the write itself is seen. The transcripts read the linked working-tree skill (14 mentions of Step 5b in i1; the installed copy has none), so these runs measured the edit.

## `/init-project` with no default store and a `{}` settings file (a1, a2)

| Run | Turns | `no memory files at <old>` printed | Store created, empty | Settings file | `.before` copy | Note |
|---|---|---|---|---|---|---|
| a1 | 46 | yes | yes | stayed `{}` (write refused), one-key content reported | — | — |
| a2 | 53 | yes | yes | stayed `{}` (write refused), one-key content reported | written | said the file is tracked by git and named `git rm --cached` — now a clause of Step 5b |

## The setting honoured headless (m1, m2)

The kit wrote `.claude/settings.local.json` with the key (the operator's part in a real project) and asked Claude Code to remember a fact.

| Run | Turns | Memory file written in the new store | Anything in the default path |
|---|---|---|---|
| m1 | 4 | yes (`demo-deploy-cadence.md` plus the `MEMORY.md` line) | nothing |
| m2 | 4 | yes (`demo-deploy-schedule.md` plus the line) | nothing |

So Claude Code 2.1.284 honours `autoMemoryDirectory` from a project's local settings file in a headless run of a scratch clone; the trust refusal did not occur here.

## `/dream` split-store stop (s1–s12) and the regression run

- **s1–s3, s4–s6, s7–s9: void.** The real home was installed from commit 4426e58 at 14:03 (the operator's action), which put a `dream` skill of the Phase 1 shape under `~/.claude/skills/`; the `/dream` command resolved to that user-level copy, shadowing the project link to the working tree, so nine runs measured text without the split check. Their first lines kept the installed format; the transcripts show the runs listing the default project directory rather than counting its memory files. Recorded as plan D15 and a bug record; the kits now pass `--setting-sources project,local`.
- **s10–s12, valid:**

| Run | Turns | First line carries both counts | `split store — merge before running` | Verdict table absent | Settings store unchanged | Default store unchanged | Clone tree clean |
|---|---|---|---|---|---|---|---|
| s10 | 7 | yes (6; default path: 2) | yes | yes | yes | yes | yes |
| s11 | 6 | yes | yes | yes | yes | yes | yes |
| s12 | 5 | yes | yes | yes | yes | yes | yes |

- **Regression, `dream-probe` run 8** (no default store, `--setting-sources project,local`): the first line carries the store and `default path …: could not be counted` (the directory was outside the run's allowed directories), the verdict table and the `promote` flag as in Phase 1 of the dream plan, no split line, store unchanged. The skill now says that phrase takes the count's place when the count cannot run.

## Context clause and capture rule

- `context/AGENTS.md` § Resume gained the Antigravity clause; `TestRenderedContextFitsAntigravityBlock` PASS at 21,787 of 22,000 bytes (bare run).
- `skills/feedback-log/SKILL.md` rule 7 names every frontmatter field the spec lists (`name`, `description`, `metadata.node_type`, `type`, `originSessionId`, `modified`) and the index line; read against `specs/sh-agent/memory-loop.md` § The store.

## Not exercised

- The settings write itself (refused headless in every run): Phase 2's E2E in this checkout, interactive.
- Antigravity: the resume clause and the capture rule have not run on agy (plan D7); the operator's real-home smoke.
- A settings file with other keys (the insert-one-line branch): the fixtures had none or `{}`; the branch is stated, not run.
- The hint list when the slug does not match (D9): the absent runs printed the message; no mismatched-slug fixture was built.
- The `~/.config/sh-agent/memory/` directory's mode: created by the runs under the operator's 0700 `~/.config/sh-agent/`; not measured.

## Cost

Ten init and remember runs about $8.20 (the init runs are long, 44–53 turns at about $1.10 each); twelve split runs about $5.10, nine of them void; one regression run.

## Cleanup

The stores the runs created under the real home were removed by name pattern after the log was written: `~/.config/sh-agent/memory/*probe-runs*` (seven directories) and `~/.claude/projects/*probe-runs*/memory` — the second pattern also matched the empty auto-created `memory/` directories of earlier sessions' scratch probes (tsn, retro, prompts); the survey at the start of this session listed every one of them with zero files, so nothing was lost, but the pattern was wider than the seeds it was meant for.
