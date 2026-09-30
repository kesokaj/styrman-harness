---
owner: main
date: 2026-09-30
---
# Verification — real home and an Antigravity-only sandbox before v1.16.1, both CLIs   (2026-09-30)

Range reviewed: v1.16.0..7bf0978. Installed into the real home as 7bf0978 (`sh-agent update` from this checkout at 17:27: 7 written, 259 unchanged, 266 managed, binary 7bf097828165, vendor all, drift none; `skills/init-project/SKILL.md` and `skills/dream/SKILL.md` byte-identical between the tree and both vendors' installed copies). Real-home runs happen in fresh clones of `~/sandbox/gate17/fixture` (the gate16 `greetbox` fixture: an older scaffold with no `TODO.md`, no `.gitignore` and no `.claude/`), outside every checkout. Every `sh-agent` call from a run is pinned to `~/sandbox/gate17/home` (BUG-107). Kit: `~/sandbox/gate17/run.sh` (real home) and `sbrun.sh` (sandbox).

## Status: PASS — the context check on both CLIs, both changed skills (`init-project`, `dream`) on each CLI from the real home (the `all` path), and both skills on Antigravity in a home installed `--vendor agy` (the new path)

The rule: project `AGENTS.md` § Release. The range changes two skills (`init-project`, `dream`), the context (one line) and no agent.

## Claude Code 2.1.285 (`claude -p`, real home)

| Check | Result |
|---|---|
| Context questions, no tools | PASS: "# AGENTS.md — global agent context", 28 skills and 8 subagents, GANNET, `bug`; 7 s |
| `/init-project` update on the fixture (vendor all) | PASS: `TODO.md` created, `.gitignore` has `.claude/worktrees/`, `.claude/settings.local.json` and `.agents/memory`, commit `373d0e7 init: update scaffold…`. The `.claude/` writes were refused, as headless Claude Code does, and reported for the user with the exact link and settings commands (Step 5 guard). So the `.claude` steps ran rather than being skipped; 110 s |
| `/dream` with a named settings file | PASS: `store: <clone>/.agents/memory (1 memory files; default path …: 0)`. The file was named in the turn (`dream-settings.json` in the clone), because this session's own sandbox refuses a shell write of any `.claude/settings.local.json`; 66 s |

## Antigravity CLI 1.2.14 (`agy -p`, real home, `--mode accept-edits`, no skip flag)

No run's stderr has an `auto-denied` line; the day's `cli-*.log` files have 0 `invalid allow entry` lines.

| Check | Result |
|---|---|
| Context questions, no tools | PASS: the same four answers. Store `conversations/0f8a97b9…db` holds the `<RULE[~/.gemini/config/rules/AGENTS.md]>` block with the canary line and no `<truncated` line; 9 s |
| `/init-project` (vendor all) | PASS: `.claude/{skills,agents,commands,rules}` → `../.agents/*` and `settings.local.json` made, the two `.claude/` ignore lines, `TODO.md`, commit `7e0fda6`. On a machine installed for both CLIs, the old behaviour stands; 362 s |
| `/dream` with a named settings file | PASS: `store: <clone>/.agents/memory (1 memory files; …: 0)`; 242 s |

## Antigravity-only sandbox (skill-matrix rows `init-agy-only`, `dream-agy-only`)

The home is `.agents/scratch/gate17-agy`, installed with `sh-agent install --home … --vendor agy` from 7bf0978 (133 files, manifest `"vendor": "antigravity"`). The settings file was written for the run: the matrix's named grants plus `command(for)`, `read_file(*)`, `write_file(*)`; no `toolPermission` key (`grep -c` → 0); `trustedWorkspaces` names only the clone `~/sandbox/gate17/runs/agyonly-1`. The token and settings were deleted by name after each batch.

- Named grants alone: both runs auto-denied within 20 s (`soft-denying tool confirmation "RunCommand"`). The skill's scripts start with `for d in …` and `root=$(git rev-parse …)`, which no `command(<name>)` grant covers; adding `command(for)` did not help.
- The session's permission classifier refused `--dangerously-skip-permissions` until the operator approved it in their own turn ("Ok, do all of these", on the recommendation to allow the skip-flag run in the throwaway clone; `testing/guide.md` § 5 allows it there). The rows now name the flag (BUG-115).

| Row | Result |
|---|---|
| `init-agy-only` | PASS, CHECK exit 0: no `.claude/`, `.gitignore` without a `.claude/` line, `.agents/memory` present and ignored, commit `8caf544 init: scaffold…`; 264 s; 0 auto-denied |
| `dream-agy-only` | PASS, CHECK exit 0: `store: …/agyonly-1/.agents/memory (1 memory files; default path …: 0)`, reached through Step 1 item 0 with no settings file on the machine; 131 s; 0 auto-denied |

Must-fail control: the same `init-agy-only` CHECK on the pre-fix fixture `~/sandbox/agyonly/proj` exits 1 (`.claude/` exists), as the fix's verifier log records (`testing/2026-09-30-1721-verify-agy-only-init.md`). The earlier pre-fix and post-fix sandbox runs in that log do not count as gate evidence (BUG-115); this log's runs replace them.

## What was not exercised

- Claude Code on an Antigravity-only machine: that machine has no Claude Code by definition.
- `/dream`'s hint "run /init-project to move them" on an Antigravity-only machine that still holds Claude Code memory at the default path. The docs lens found it (low); it is deferred to the handoff.
- The Mac and devbox places: they receive the release by `sh-agent update` (operator steps, `deploy/remote.md`).
