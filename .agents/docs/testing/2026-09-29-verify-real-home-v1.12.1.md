---
owner: main
date: 2026-09-29
---
# Verification — `/init-project` update path on Antigravity, real home at v1.12.1   (2026-09-29)

Range reviewed: v1.12.1 (781021e) installed in the real home (`sh-agent update`, binary 781021e34aa7); the run from the real `$HOME` in a fixture under `.agents/scratch/agy-linked/` (a trusted path for Antigravity).

## Status: PASS — 1 of 1 run; the one row the v1.12.1 gate had listed as not checked on Antigravity

The v1.12.1 range changed one skill, `init-project` (Step 1's sentence and the Step 3 and 4 guards). `AGENTS.md` § Release asks for one Antigravity run of each changed skill before the tag; the gate report (`reviews/2026-09-29-preflight-v1.12.1.md`) had only the Claude Code run of the linked-store row and said so under "What was not checked", while its verdict paragraph read as if the rule were met. This run closes the gap after the tag; the addendum on that report records the overstatement.

## Fixture (asserted before the CLI started)

A one-file Go service, an older scaffold (`AGENTS.md`, `REQUIREMENTS.md`, empty `TODO.md`/`ROADMAP.md`, a bug index, no deploy guides, no `.claude/` links), `.agents/memory` a committed symlink to an empty directory outside the clone, and the default store at `~/.claude/projects/<slug>/memory/` seeded with one memory plus the index; one commit. The script asserted the link, its git tracking, the empty target and the two default-store files before `agy` ran.

## Antigravity CLI (`agy -p`, `--mode accept-edits`, real home, installed skills)

| Check | Result |
|---|---|
| Update path | PASS — every project file kept; the customization directories, trail directories, deploy guides and `.claude/*` links added; one commit `init: update scaffold to current templates`, push skipped as asked |
| Step 3/4 guards (the range's change) | PASS — the report names `.agents/memory` as "symlink tracked by git; kept and not written through"; the link is intact afterwards |
| Step 5b trust check | PASS — stopped ("a committed symlink would carry the operator's memory wherever it points"); the target directory still empty; the two default-store files untouched; `.claude/settings.local.json` absent; "moved: nothing" |
| Step 6 | PASS — `.claude/settings.local.json` and `.agents/memory` appended to `.gitignore` (the matrix row allows this) |
| Logs | stderr empty (no `auto-denied`); no `invalid allow entry` in `cli.log`; exit 0 |

## Not verified here

The `init-memory` and `init-memory-absent` rows on Antigravity (the happy path that writes the settings file and moves a store); `/distill` dispatching the `dreamer` on either CLI; the Mac and the devbox. Headless Claude Code refuses the settings write, so the happy path's write is verified only interactively (`testing/2026-09-29-verify-real-home-v1.12.0.md`).

## Commands

`.agents/scratch/agy-linked/` held the fixture, `out.txt`, `err.txt` and `status.txt`; removed after this log, with the seeded default store for the fixture's slug.
