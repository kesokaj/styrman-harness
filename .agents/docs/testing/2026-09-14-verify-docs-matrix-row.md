---
owner: main
date: 2026-09-14
---
# Verification — the `documentation` skill-matrix row on both CLIs   (2026-09-14)

Range reviewed: the working tree of `plans/2026-09-14-documentation-steps.md` Phase 2, installed into the sandbox home

## Status: PASS for the row's check on both CLIs; Claude Code's run did not finish its own close-out

Scope: the new `documentation` row of `skill-matrix.md` (prompt: add a `SHUTDOWN_TIMEOUT` environment variable through the Definition of Done; check: `SHUTDOWN_TIMEOUT` in `README.md`, `.agents/docs/deploy/local.md`, and a `changes/` entry), run once per CLI in fresh clones of the demo fixtures. Not run: the rest of the matrix.

## Fixture change first

The first attempt used a `PORT` task. It showed two fixture defects, not harness defects: the demo service already read `PORT`, and it had no `.agents/docs/deploy/` (scaffolded before the guides existed). Both demo repos gained the three deploy guides (`local.md` filled, no mention of the new setting; commits 9ae4eff in `demo-app`, e2cad5c in `demo-app-agy`) and the row's task became `SHUTDOWN_TIMEOUT`.

## Runs (detached, 1500-second budget each, status files under `.agents/scratch/e2e-docs3/`)

| CLI | Exit | Row check | What the task produced |
|---|---|---|---|
| Antigravity 1.2.2 | `agy-docs2 exit 0 at 2026-09-14T08:04:44Z` (about 11 minutes) | PASS | commit `818bd68` with the feature, README, `deploy/local.md` (Configure row and the stop procedure), a changes entry carrying the new `## Docs` section, a verify log, a new `testing/demo-app/guide.md` written by the verifier under the new rule, a handoff |
| Claude Code 2.1.270 | `claude-docs2 exit 124 at 2026-09-14T08:18:17Z` (budget exhausted) | PASS on the working tree | the feature (`internal/shutdown/`), tests, README, `deploy/local.md`, a changes entry, two reviewer rounds, a security review, two verify logs, a `--docs` lens report (`reviews/2026-09-14-review-worktree-docs.md`), two bug records — all staged, not committed; no handoff. The Definition of Done it ran in full consumed the budget before the commit. |

Antigravity logs: 0 `Failed to parse`. Neither run asked a question.

## Reading

The behaviour the row exists to prove holds on both CLIs: a task that adds a setting updates README and the deploy guide and records the docs in its changes entry; the verifier creates the app's testing guide; the docs lens runs. The Claude Code timeout is a budget effect of the full Definition of Done on a 25-minute headless run, recorded here and in the handoff; a longer budget or a lighter row prompt would let it finish the commit. Not a harness defect; not re-run today.

## Environment

Linux (Debian 13, amd64), Claude Code 2.1.270, Antigravity CLI 1.2.2, harness working tree after d87df30.
