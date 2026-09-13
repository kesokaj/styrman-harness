---
owner: main
date: 2026-09-13
---
# Handoff — 2026-09-13 — `/review` shipped on both CLIs, `/audit` folded in; two open bugs on headless explorer writes

Session ID: session_01CtMYM9AGp33edxBQGURNaB
Project: /home/kesokaj/workspace/styrman-harness
Branch: main, clean once this handoff's commit is pushed

## Context

The operator: "Im missing the code-review feature in agy … this harness should be the same for agy and claude." Then "go all, fold /audit into /review". Plan: `.agents/docs/plans/2026-09-13-review-skill.md` (one phase, all boxes ticked). Supersedes `2026-09-13-real-home-updated.md`.

## Completed

- `skills/review/SKILL.md`: one command on both vendors — scope (working tree, range, PR, paths, `--project`) separate from lens (`--code`, `--security`, `--simplify`, `--docs`, `--tests`, `--all`), `--fix` as a capped loop, `--high`; the skill computes the scope text and per-dispatch slugs, the agents are unchanged. `skills/audit/` deleted; `/review --project --all` is the audit. References repointed; 27 skills. ADR-017; entry `.agents/docs/changes/2026-09-13-review-skill.md`.
- Cold read: three rounds of three readers (`.agents/docs/reviews/2026-09-13-cold-read-review-skill-*.md`); round-3 residue resolved by recommendation and written into the plan's D3/D9.
- Sandbox e2e: `/review` on Claude Code and Antigravity each wrote `reviews/2026-09-13-review-worktree-code.md` in a fresh clone with the scope text on `Range reviewed:`; `/review --project --all` on Claude Code (detached, 14 minutes) wrote all seven artifacts and the merged record with one verdict per lens; negative prompts loaded nothing; Antigravity logs clean.
- Definition of Done: reviewer REQUEST CHANGES → both findings fixed (`reviews/2026-09-13-review-review-skill.md`); verifier FAIL on my false citation of a timed-out run (`testing/2026-09-13-verify-review-skill.md`) → re-run `testing/2026-09-13-verify-review-skill-2.md`. BUG-013 records the false citation.

## Open / In Progress

Nothing queued in this repo. Operator-only, carried from the previous handoff: the Mac (`sh-agent update`, three smoke checks, Claude Code there never checked); the Known-limits sentence; Go ≥ 1.24.6 on both machines; the downstream mirrors in styrman-crew and argolis-devbox; item 20 after the first Mac retro; the full skill matrix (`testing/skill-matrix.md` now has `review` and `review-project` rows).

Open bug: **BUG-012** — the `explorer` agent, dispatched headless with a narrow Bash allow-list, reports a brief it never wrote; with `bash`/`printf` allowed it writes. Leading cause: its brief goes through a shell heredoc. The fix candidate is one line in `agents/explorer.md` (write with the Write tool) plus a controlled retest — a task for the next session, not started.

## Key Decisions

- Not `code-review`: Claude Code's built-in owns that name. The skill never names a vendor built-in.
- `--fix` is refused with `--project`; paths are read whole; ranges use the agents' three-dot diff; one scope class per run.
- Antigravity's first `/review` launch printed nothing and left no log file; the re-run with `--log-file` passed — treated as a one-off launch failure, not a skill defect (unverified).

## Corrections

- None in this part of the session.

## Agent notes (unverified)

- A headless Claude Code dispatch that fans out to six agents needs more than the tool's ten-minute budget; run it detached with its own status file and read the status line before citing it (BUG-013).
- The security agent on the strongest model is the slow lens; `/review --project --all` should be expected to take a quarter of an hour.

## Environment State

- Real home on this box: harness at the commit before this handoff's; run `sh-agent update` to pick up `/review`. Sandbox home at the working tree. Clones under `/home/kesokaj/sandbox/runs/{claude,agy}-review`, `claude-project`, `claude-project2` hold the e2e artifacts; evidence under `.agents/scratch/e2e-review/` (gitignored). No running processes.

## Artifacts to Read

- `skills/review/SKILL.md` — the command
- `.agents/docs/plans/2026-09-13-review-skill.md` — decisions D1–D9 and the residue section
- `.agents/docs/changes/2026-09-13-review-skill.md` — what shipped and how it was verified
- `.agents/docs/bugs/INDEX.md` — BUG-012 (open), BUG-013

## Next Steps

1. Operator, here: `sh-agent update`, then try `/review` on a real change in a project on each CLI.
2. Next session: BUG-012 — make the explorer write with the Write tool, retest headless twice, close the record.
3. The Mac items and the mirrors, as before.
