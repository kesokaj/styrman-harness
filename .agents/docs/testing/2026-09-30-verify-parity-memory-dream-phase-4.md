---
owner: main
date: 2026-09-30
---
# Verification — parity plan Phase 4: `sh-agent dream` runs, both CLIs   (2026-09-30)

Range reviewed: 1bf5f4f..worktree (the run path, `report`'s harness-target proposals, `/retro` source 9, D7b). The binary is built from the working tree (`.agents/scratch/sh-agent-p4`). Plan `plans/2026-09-30-parity-memory-dream.md` Phase 4: E2E 3, 4 and 6.

## Status: PASS — E3, E4 and E6, after one real defect in the agy invocation was found and fixed (D7b)

## Kit
- **Fixtures**: `~/sandbox/p4/proj-{a,b,c,d}`, each a clone of the p2 fixture seeded by `~/sandbox/p2/dreamer/fixture.sh`: two bug rows for one defect, a five-step release in a handoff and a change, two `feedback` memories from two sessions, an open `project` memory. Each asserts "seeds ok". All sit outside every checkout.
- **Claude Code runs**: `--home .agents/scratch/p4home`, whose `search.db` indexes only the fixtures (`sh-agent index --home`); the child `claude` keeps the real environment, as the plan says.
- **Antigravity runs**: the whole process under `HOME=.agents/scratch/p2-agy` (manifest vendor `antigravity`). The token was copied chmod 600. The settings file has an EMPTY allow list and does not trust `dream-work`, as on the real home. Token and settings were deleted by name after.

## E3 — golden path (Claude Code)

| Run | Result |
|---|---|
| `sh-agent dream --all --home p4home --vendor claude` | exit 0, 53 s; `proj-a ok 2 proposals`, `proj-b ok 2 proposals`, `- run 2 projects, 4 proposals` |
| Files | `proposals/dream/2026-09-30-proj-{a,b}.md` and `proposals/dream/drafts/2026-09-30-proj-a-tag-release/SKILL.md`, `proposals/dream/drafts/2026-09-30-proj-b-project-release/SKILL.md` |
| Content (proj-a) | 1. "go vet clean before every release tag", `**Status:** queued`, `**Target:** project`, `**Kind:** rule`, evidence BUG-001/002 and both memories; 2. a project skill for the hand release, evidence the change entry and the handoff |
| Second run, nothing changed | exit 0; `proj-a unchanged`, `proj-b unchanged`; proposal files 4 → 4 at that moment (the later Claude Code `--project proj-a` run below added `2026-09-30-proj-a-2.md` and a draft) |

## E4 — edges
- A memory file `leak.md` holding `api token = FAKEVALUE-SECRET-123`: `--project proj-a --dry-run` printed the prompt, and `grep -c FAKEVALUE` returned 0. The fake-CLI unit test pins the same on the real call's stdin.
- A deleted root: `proj-c` was indexed, then removed, and the next run logged `proj-c skip gone`.
- `sh-agent status` on an installed home prints `dream      never run; timer not installed`. On a home with nothing installed, `status` stops before its retro and dream lines, which is existing behaviour.

## E6 — Antigravity (agy-only sandbox)

| Attempt | Invocation | Result |
|---|---|---|
| 1 | `agy -p <prompt> --mode plan --disable-slash-commands --output-format text` (D7a) | **fail invalid output**. stderr: `warning: --mode plan has no effect while slash command expansion is disabled.` and `no output produced — a tool required the "command" permission … auto-denied`. The job kept the output in `dream-failed/`, left the watermark, exited 1 |
| probe | `--mode plan --output-format text` | no answer, auto-denied |
| probe | `--mode plan --output-format json --json-schema <schema>` | `"response":""`, auto-denied |
| probe | the prompt prefixed with "IMPORTANT: every tool call in this session is refused and ends it with no answer … Answer directly with the JSON object only." | the full JSON, 0 denied |
| 2–4 | the job with that warning in its instructions (both CLIs) and `--mode plan --output-format text` | **3 of 3 ok**, 2 proposals and a draft each; same-day files `…-proj-d.md`, `-2`, `-3` |

The D7a probe had used a one-line prompt, and a one-line prompt never tempts a tool call; the real 14 KB prompt did (see the plan's D7b). Claude Code with the same warning: `--project proj-a` exit 0 (printed to the session, not saved), and `p4home/.config/sh-agent/dream.log` line `2026-09-30T18:46:40Z …/proj-a ok 1 proposals`. The run's own `run1.txt`/`run1.err` are empty because a successful run prints nothing; the log is its record.

## Not exercised
- The timer (Phase 5).
- A real `sh-agent report` push carrying `proposals/dream/harness/`: `Collect` is unit-tested.
- A night on which agy still calls a tool despite the warning. The failure path is unit-tested; the frequency is unknown beyond 3 of 3.

Audited by `testing/2026-09-30-verify-parity-memory-dream-phase-4-5.md`: no claim contradicted. The three imprecisions it named (draft paths, the 4 → 4 count, the unrecorded exit status) are corrected above.
