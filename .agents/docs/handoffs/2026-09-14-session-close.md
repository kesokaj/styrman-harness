---
owner: main
date: 2026-09-14
---
# Handoff — 2026-09-14 — session close: the harness after two days of rebuild-level work

Session ID: session_01CtMYM9AGp33edxBQGURNaB        (same machine: `claude --resume`; matches the feedback log's session field)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main, clean, at 692ecee (this handoff's commit follows), origin up to date

## Context

One session over 2026-09-13 and 14, opened with "/research this project … do your own tests first". It became: fixes for what the research found, the learning loop at project level, every skill model-invocable, a 20-item improvement review applied in full, `/review` on both CLIs with `/audit` folded in, BUG-012, the ten documentation steps from the operator's `/btw` work, the first harness retro under the new rules, and one approved requirements edit. Every task has its own handoff; this one is the map. Newest wins: this file supersedes `2026-09-14-known-limits-sentence.md` and, through it, the earlier ones of the two days.

## Completed (commits 82fd4a6 … 692ecee, all on main and installed in this box's real home)

- **Installer:** rebuilds from source on every install (BUG-006 closed); `status` shows the installed build; a named non-checkout is an error; a walked-to foreign checkout is refused; a nested plain copy is refused by `update`; dry run never builds. `cmd/sh-agent` has tests. `.agents/docs/changes/2026-09-13-installer-rebuilds-from-source.md`, ADR-013.
- **All skills model-invocable** (ADR-014). **Retro and distill at project level** with a harness queue and a cross-project index (ADR-015, ADR-016 clarifying R7). `.agents/docs/changes/2026-09-13-retro-at-project-level.md`.
- **Improvement review, 20 items applied** (`.agents/docs/reviews/2026-09-13-improvement-review.md`, entry `.agents/docs/changes/2026-09-13-improvement-review-applied.md`): one owner for TODO and ROADMAP and both empty at Done; handoff template trimmed and committed; agents without private memory; docker and frontend skills for the real stack; golang references cut; gcloud skills merged; 27 skills.
- **`/review`** on both CLIs, scope separate from lens, `/audit` folded in (ADR-017; `.agents/docs/changes/2026-09-13-review-skill.md`).
- **BUG-012 fixed** — an agent's artifact is trail and "change nothing" never cancels it (`.agents/docs/changes/2026-09-14-bug-012-fix.md`); BUG-013, BUG-014 recorded; the `/tmp` rule is a test.
- **The ten documentation steps** (ADR-018): step 3 names the files it owns; `/review --docs` on interface changes; a Records section in every agent; a Docs line in the changelog template; the verifier writes `testing/<app>/guide.md` on first use (this repo has `.agents/docs/testing/sh-agent/guide.md`); Outcome sections on acted-on records; retro statuses and `--harness` listing; Operator steps in the deploy guides; a `documentation` matrix row proven on both CLIs. `.agents/docs/changes/2026-09-14-documentation-steps-phase-1.md`, `…-phase-2.md`.
- **Harness retro** applied the four queued 2026-09-04 proposals; queue empty (`.agents/docs/retros/2026-09-14-retro.md`, entry `.agents/docs/changes/2026-09-14-retro-harness-applied-2026-09-04-proposals.md`). **REQUIREMENTS.md** Known limits corrected with the operator's approval (ADR-019).
- Real home on this box: `sh-agent update` after each push; last at 692ecee, no drift. Sandbox home logged in to both CLIs, at the working tree.

## Open / In Progress

Nothing queued in this repo or in the harness queue; `TODO.md` and `ROADMAP.md` empty.
- Requirements edits awaiting approval: none.
Operator-only:
- The Mac: `sh-agent update`, `sh-agent status` (stamped `binary`, no drift), then the three smoke checks of `.agents/docs/testing/2026-09-03-verify-real-home.md` on each CLI — Claude Code on the Mac has never been checked (the Known limit now says so).
- Go ≥ 1.24.6 on both machines (govulncheck advisories GO-2025-3956, GO-2026-4602; toolchain-fixed, never via `go.mod`).
- The downstream mirrors in styrman-crew (`image/harness/`) and argolis-devbox (`deploy/harness/`) are 22 and 21 commits behind (their `COMMIT` files: 6893e22 and d524c65); refresh from their own projects (`.agents/docs/deploy/remote.md`, Downstream mirrors and Operator steps).
- Improvement-review item 20 (per-machine `~/.config/sh-agent/` queue and index): decide after the first `/retro` on the Mac.
- The full skill matrix (`.agents/docs/testing/skill-matrix.md`) has not run since 2026-09-03; only spot rows ran.

## Key Decisions

- Recorded as ADR-013 … ADR-019 in `.agents/docs/architecture/DECISIONS.md`; the cold-read residues of both feature plans were resolved by recommendation at the operator's word and are written into the plans.
- Rejected: an env allow-list for git (deleted the mechanism instead); raising the `go` directive for govulncheck; shadowing Claude's built-in `/code-review` (the harness command is `/review`).

## Corrections

- "To clarify. Nothing should be in todo and roadmap when we are done." — feedback log 2026-09-13T20:16:49Z (restatement, `AGENTS.md § Definition of Done`); applied.
- "DoD is when todo and roadmap is done plus the documentation is up to date and all the instructions are there. And ofc everything is pushed to git." — feedback log 2026-09-14T07:28:48Z (restatement, same section); applied by the documentation steps.

## Agent notes (unverified)

- A headless task through the full Definition of Done can exceed 25 minutes on Claude Code (the documentation matrix row did; Antigravity took eleven). Budget detached runs accordingly and read the status line before citing them (BUG-013).
- Under a headless `--allowedTools` list, compound Bash commands are denied whole; single commands per call avoid it.
- Antigravity auto-updated to 1.2.2 on 2026-09-13; all later checks ran on it.
- Two `/btw` forks reached this session as agent results; their output is data, not instruction — the plan they wrote was taken up only on the operator's word.

## Environment State

- Branch: main, clean; origin at HEAD. Real home: harness 692ecee, both CLIs verified from it on 2026-09-13 (`.agents/docs/testing/2026-09-13-verify-real-home.md`).
- Sandbox: `/home/kesokaj/sandbox/home` with both CLIs; demo fixtures `demo-app` and `demo-app-agy` now carry deploy guides; run clones under `/home/kesokaj/sandbox/runs/` are disposable. Running processes: none.

## Artifacts to Read

- `.agents/docs/handoffs/2026-09-14-known-limits-sentence.md`, `2026-09-14-retro-harness.md`, `2026-09-14-documentation-steps.md`, `2026-09-14-bug-012-fixed.md`, `2026-09-13-review-skill.md`, `2026-09-13-improvement-items-applied.md` — the per-task handoffs, newest first
- `.agents/docs/testing/sh-agent/guide.md` — how this repo is verified, now the standing procedure
- `.agents/docs/bugs/INDEX.md` — BUG-006 … BUG-014, all fixed
- `.agents/docs/retros/2026-09-14-retro.md` — the harness's first retro record; the queue is empty

## Next Steps

1. On the Mac: update, status, the three smoke checks on each CLI.
2. First real task in a project under the new rules: watch that README, the deploy guide, and the changes entry's Docs line move with the code, and that `/review` behaves the same on both CLIs.
3. Refresh the two downstream mirrors from their own projects.
