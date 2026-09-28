---
owner: main
date: 2026-09-28
---
# Verification — time-sortable names, E2E 1: agy cold resume   (2026-09-28)

Range: `9a17e0a` (pushed HEAD). Plan: `plans/2026-09-28-time-sortable-trail-names.md`, Phase 3.

## Status: PASS (with a weak control)

The criterion held in every run. The must-fail control did not fail, so this E2E shows that the rule works, not that it was necessary in this setup. See below.

## The condition, proven before the runs

Two fresh clones:
- **This repo at `9a17e0a`:**
  - by name, `ls handoffs | tail -1` is `2026-09-28-research-okf.md` (the oldest of today's);
  - every file has one mtime (clone time);
  - the true newest by `git log` is `2026-09-28-1450-time-sortable-names.md`.
- **A clean demo** (no harness source), holding `2026-09-28-zz-older-session.md`, committed first, "Next: write the tests", and `2026-09-28-1500-newer-session.md`, committed second, "Next: deploy feature A to staging":
  - name-sort last is the older file;
  - one mtime;
  - `git log` newest is `1500`.

## Runs (agy 1.2.8, `-p --mode plan`, same cold-resume question)

| Run | Context | Project | Handoff cited | Result | Tool calls |
|---|---|---|---|---|---|
| new ×2 | new (sandbox home `.agents/scratch/tsn`, harness at `9a17e0a`) | this repo | `2026-09-28-1450-time-sortable-names.md` | 2 of 2 right | 9, 11 |
| new ×2 | new | demo | `2026-09-28-1500-newer-session.md`, "Deploy feature A to staging" | 2 of 2 right | 7, 7 |
| control ×2 | old (real home, v1.6.0) | demo | `2026-09-28-1500-newer-session.md` | 2 of 2 right | 9, 9 |

**How each version decided.** The control listed the directory, read both handoffs, and broke the tie with `git log -n 5 --stat`. The new-context runs read both and needed no `git log`. With two same-day files, reading both is cheap, so the rule's gain here is two calls, not a different answer. The expensive case is the one this plan's `## Why` measured before the change: 8 of 17 calls, and the wrong handoff opened first.

## What went wrong on the way (BUG-091)

- **First attempt, new context:** both runs produced no output. Headless agy auto-denied the first shell command, because the sandbox home had no `~/.gemini/antigravity-cli/settings.json`. The real home's file was copied in: permissions and trusted workspaces, no secret.
- **First attempt, control:** it ran on the harness clone and answered right by citing the new rule, which it read from `context/AGENTS.md` in the clone's tree. A control for a context change cannot run inside the harness repo, which is why the demo project was added.

## Not exercised

- Claude Code for E2E 1 (Phase 2 covered its context runs).
- The real home after `sh-agent update`: the operator's smoke check, in the handoff.

The credential copy was `chmod 600` and deleted by name after each batch; `ls` confirmed it was gone. Both clones were unmodified afterwards (`git status` empty).
