---
owner: main
date: 2026-09-29
---
# Behavioural run — `dream` skill, Phase 1   (2026-09-29)

Range reviewed: working tree over d8e69bf (`skills/dream/SKILL.md`, new).

## Status: PASS — 3 of 3 skill runs on every seed; 0 of 3 clean controls on the skill's vocabulary

## What was under test

`skills/dream/SKILL.md` as written in the working tree, loaded as a project skill through a symlink `.claude/skills/dream` in a fixture repository (`.agents/scratch/dream-probe/fixture`, commit a6cd45c), run headless: `claude -p` with the prompt `/dream --dry-run` plus one line naming the `--settings` file, `--settings <run>/settings.json` holding `{"autoMemoryDirectory": "<absolute run dir>/mem"}`, `--permission-mode acceptEdits`, `--add-dir` on the harness `skills/dream` (so the symlink resolves) and on the run's store, `--allowedTools` = Skill, Read, Grep, Glob and read-only Bash prefixes, no Write. Real home for credentials (Claude Code 2.1.284); real home untouched. Kit: `.agents/scratch/dream-probe/run.sh`, `mem-template/`.

## Fixture store (six memory files plus `MEMORY.md`, copied fresh per run; every seed asserted by `grep`, `test`, `stat` and a session-id comparison before the runs)

| Seed | File(s) | Expected |
|---|---|---|
| Rule sits verbatim in the clone's `AGENTS.md` | `run-tests-before-claim.md` (feedback) | `graduated` |
| Names `cmd/nothing/main.go` and `deploy.sh --region`, neither exists | `entry-point.md` (project) | `broken` |
| One fact in two files, `originSessionId` 2222… and 3333…, no rule holds it ("listings sorted by name", grepped absent from the global context and the clone) | `listings-by-name.md`, `sort-listings.md` (feedback) | `merge`, flagged `promote` |
| Work still open in the newest handoff (v2 migration, issue #12) | `v2-migration.md` (project) | `keep` |
| Cites `.agents/rules/commit-style.md`, which is mode 000 in the clone (the check cannot run) | `commit-subjects.md` (feedback) | `keep`, why names the unreadable file |
| Index line for a file that does not exist | `MEMORY.md` → `old-deploy-notes.md` | `index` row |

## Round 1 (runs 1–3, controls c1–c3): two setup defects, both mine

- The merge pair restated the asking rule that the harness's global context already holds; every run correctly marked both files `graduated` citing `~/.claude/AGENTS.md:6`. Correct behaviour, wrong seed. Reseeded with a fact no loaded rule holds.
- The controls deleted the skill symlink but kept `--add-dir` on the harness skill directory and left the deletion visible in `git status`; all three controls read `skills/dream/SKILL.md` from the harness and followed it (transcripts: a `Read` of that path in each). Control rebuilt: the removal committed in the clone, no harness `--add-dir`.
- Everything else in round 1 was already 3 of 3: store line, `graduated`, `broken`, `keep`, guard `keep`, `index`, nothing written.

## Round 2 (runs 4–6 with the skill, controls c4–c6 without it)

| Run | Store line names the run's store | `graduated` | `broken` | `merge` on the pair | `promote` on the pair | `keep` | guard `keep` | `index` row | Store unchanged | Clone tree clean | Load marker in first user message |
|---|---|---|---|---|---|---|---|---|---|---|---|
| 4 | yes (6 memory files) | yes | yes | yes | yes | yes | yes | yes | yes | yes | yes |
| 5 | yes | yes | yes | yes | yes | yes | yes | yes | yes | yes | yes |
| 6 | yes | yes | yes | yes | yes | yes | yes | yes | yes | yes | yes |
| c4 | no | no | no | (by-hand "merge" noted) | no | (noted) | no | (noted) | yes | yes | no |
| c5 | no | no | no | (by-hand "merge" noted) | no | no | (noted) | (noted) | yes | yes | no |
| c6 | no | no | no | (by-hand "merge" noted) | no | no | no | no | yes | yes | no |

Each control began "`/dream` isn't installed in this session" and then read the store by hand; none printed the store line, none used the verdicts `graduated`, `broken` or `promote`, none applied the guard rule. So the definition, not the prompt, produces the audit.

The load marker: in a `claude -p` run the slash command is expanded into the first user message (`<command-name>/dream</command-name>` followed by the skill body); no run used the `Skill` tool. Present in 3 of 3 skill runs, absent in 3 of 3 controls (plan D13).

Tool calls, executed and refused counted separately (from `~/.claude/projects/-home-kesokaj-workspace-styrman-harness--agents-scratch-dream-probe-runs-<n>-clone/`):

| Run | Turns | Refused by the permission gate | OS or exit errors |
|---|---|---|---|
| 4 | 18 | 4 (compound Bash commands needing approval; `ls ~/.config/sh-agent` blocked by Claude Code) | 2 (EACCES on the mode-000 guard file — the guard seed working) |
| 5 | 21 | 3 | 1 (EACCES) |
| 6 | 28 | 5 | 1 (EACCES) |
| c4–c6 | 14–17 | 1–2 | 0–1 |

Cost: about $0.55 per skill run, $0.36 to $0.40 per control; twelve runs in all, about $5.70.

## Not exercised

- The apply path (`/dream` without `--dry-run`, a user's go): every run was a dry run, and headless runs never delete; Phase 2's interactive task on this checkout's real store covers it (plan D12).
- Antigravity: no run (plan D10); the installer's frontmatter lint passed on `sh-agent install` into `.agents/scratch/dream-home`, drift none, the skill present under both vendors' skill directories.
- The `feedback.jsonl` route to `promote` (an entry from another session naming the memory's rule): the fixture proved `promote` through the merged pair's differing session ids only; `~/.config/sh-agent` is blocked for headless reads anyway.
- `/retro` reading the verdicts (Phase 2).
- The real home's own store (`asking-format.md` graduation): Phase 2.

## Phase 2 — `/retro --dry-run` reads the dream verdicts (runs r1–r3)

Same fixture, with the working-tree `retro` and `feedback-log` skills linked beside `dream` (fixture commit after a6cd45c), the same reseeded store and `--settings`, prompt `/retro --dry-run` plus "no corrections this session, the evidence is the trail and the memory store", `--add-dir` on the three harness skill directories and the store, read-only tools.

| Run | Turns | `## Memory` section with the store line | `graduated`, `broken` verdicts | merged pair flagged `promote` | pair as a Findings row with `rule_ref null` | says nothing was written | Store unchanged | Clone tree clean |
|---|---|---|---|---|---|---|---|---|
| r1 | 24 | yes | yes | yes | yes | yes | yes | yes |
| r2 | 21 | yes | yes | yes | yes | yes | yes | yes |
| r3 | 28 | yes | yes | yes | yes | yes | yes | yes |

Cost about $0.62 per run. Each run also reported, correctly, that the fixture trail lacks `changes/`, `DECISIONS.md` and `retros/`, and that `feedback.jsonl` holds nothing for the project. The matrix row `dream` check, run against the Phase 1 clones 4–6 by the main agent: PASS in each.

## Phase 2 — the apply run on this checkout's real store

Run by hand from the working-tree skill text (not yet installed in the real home), table put to the operator, answer "Go all (Recommended)"; two `graduated` memories deleted, two kept, index rewritten and consistent, record `retros/2026-09-29-dream.md`, one `kind: dream` line in `~/.config/sh-agent/retros.jsonl`. E2E criterion 1 met.
