---
owner: main
date: 2026-09-13
---
# Handoff — 2026-09-13 — improvement review applied in full; the harness is at e5833a3

Session ID: session_01CtMYM9AGp33edxBQGURNaB        (same machine: `claude --resume`; matches the feedback log's session field)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main, clean, pushed at e5833a3 (plus this handoff's commit)

## Context

Second half of the 2026-09-13 session. The operator approved every item of `.agents/docs/reviews/2026-09-13-improvement-review.md` ("go all"), then restated the Definition of Done: "Nothing should be in todo and roadmap when we are done." This handoff supersedes `2026-09-13-learning-loop-and-improvement-review.md` (newest wins).

## Completed

- Items 1–17 and 19 of the improvement review applied; 18 in part; 20 deferred by decision. Record: `.agents/docs/changes/2026-09-13-improvement-review-applied.md` (one line per item). Commits: 8a8b276 (only the deleted skill — BUG-011) and e5833a3 (everything else).
- Definition of Done step 4 now says `ROADMAP.md` and `TODO.md` are EMPTY at Done and the main agent owns both; operator actions and ideas live in the newest handoff. `.agents/docs/changes/2026-09-13-done-means-empty-todo-and-roadmap.md`; feedback entry logged (restatement, `AGENTS.md § Definition of Done`). Both files are empty now.
- Go: `ownCheckout` is a `.git` lookup that never asks git; a nested plain copy (the crew mirror) is refused by `update` and reports empty; `status` reads the installed binary's build info. Three security rounds on that mechanism, converged (`.agents/docs/reviews/2026-09-13-security-improvement-items*.md`); verifier PASS (`.agents/docs/testing/2026-09-13-verify-improvement-items.md`); reviewer APPROVE WITH NITS after two out-of-turns runs (`.agents/docs/reviews/2026-09-13-review-improvement-items.md`), its one real finding (frontend SPA vs `createServerFn`) fixed.
- Skills: 27 (`gcloud-auth-verification` merged into `gcloud-setup`); docker rewritten for a Go multi-stage image; frontend rewritten to what `styrman-crew/web` really does (SPA, client embedded via go:embed, no Nitro); golang references cut and split with zero dangling links; handoff template trimmed (frontmatter, repo-relative paths, commit step, no `--global`); bug/feature/research/feedback skills per the review.
- Agents: no `memory: project`; explorer writes its brief early; wording is vendor-neutral.
- Both CLIs smoke-checked from the sandbox after the change (context, `golang` auto-load, `explorer` dispatch; Antigravity 1.2.2, log clean); the Done rule quoted back correctly by both.
- BUG-010 (batch edit half-applied) and BUG-011 (chained commit) recorded; the global hard rule on chains into a commit now names staging too (third occurrence across projects). `REQUIREMENTS.md` R7 clarified with approval (ADR-016).

## Open / In Progress

Nothing is queued in this repo. Actions only the operator can take:

- On this box and on the Mac: `sh-agent update` (or `go run ./cmd/sh-agent install`) once — the installed 2026-09-03 binary cannot rebuild itself (BUG-006); then `sh-agent status` shows a `binary` line without a warning. Then the three smoke checks in `.agents/docs/testing/2026-09-03-verify-real-home.md` from each real home, on each CLI (Antigravity is 1.2.2 now).
- Optional `REQUIREMENTS.md` Known-limits edit (needs approval): replace "The Mac build is confirmed at the operator's next `sh-agent update` there (see the rebuild plan)." with "Confirmed on the Mac for Antigravity on 2026-09-04 (`.agents/docs/testing/2026-09-04-verify-mac-real-home.md`); Claude Code on the Mac is unverified."
- Upgrade Go to ≥ 1.24.6 on both machines (govulncheck: GO-2025-3956, GO-2026-4602; stdlib, toolchain-fixed, rebutted as a code change).
- Refresh the downstream mirrors in styrman-crew (`image/harness/`) and argolis-devbox (`deploy/harness/`) — a task in each of those projects; both are on the pre-BUG-006 harness (`.agents/docs/deploy/remote.md` "Downstream mirrors").
- Item 20 (per-machine `~/.config/sh-agent/` queue and index): decide after the first `/retro` on the Mac.
- Re-run the full skill matrix (`.agents/docs/testing/skill-matrix.md`) when convenient; only spot checks ran today.

## Key Decisions

- Done means `TODO.md` and `ROADMAP.md` are empty; the main agent owns both; the changelog agent no longer touches TODO.
- The plain-copy check never asks git (`.git` lookup); git's own environment, config and SSH command are the operator's, accepted in a comment next to `Pull`; no env allow-list.
- Rejected: an env allow-list for git; editing the Known-limits sentence without approval; amending or force-pushing 8a8b276.

## Corrections

- "To clarify. Nothing should be in todo and roadmap when we are done." — feedback log 2026-09-13T20:16:49Z (`signal: restatement`, `rule_ref: AGENTS.md § Definition of Done`).

## Agent notes (unverified)

- Both `explorer` smoke runs answered in chat only and wrote no brief; the prompts said "change nothing", which the agents may have read as "write nothing". Not proven a defect.
- The reviewer agent ran out of turns twice on a diff that included the golang reference cut; the scoped third run worked. A large-diff review should be scoped by path from the start.
- `ownCheckout` accepts a `.git` file (worktrees, submodules) — deliberate; a copy carrying a `gitdir:` file would report that other repo.
- The 260-vs-259 file count in install/status output is the binary, excluded from the drift comparison by design (BUG-001).

## Environment State

- Branch: main, clean; origin at HEAD. Real home on this box: still the 2026-09-03 binary and content behind HEAD until the operator's `sh-agent update` (R6).
- Sandbox home `/home/kesokaj/sandbox/home` at the working tree of e5833a3; clones under `/home/kesokaj/sandbox/runs/` (`*-retro-p3`, `*-smoke`); evidence under `.agents/scratch/e2e*/` (gitignored).
- Running processes: none.

## Artifacts to Read

- `.agents/docs/changes/2026-09-13-improvement-review-applied.md` — what changed per item
- `.agents/docs/reviews/2026-09-13-improvement-review.md` — the 20 items and their outcome
- `.agents/docs/bugs/INDEX.md` — BUG-006..011 from today
- `.agents/docs/architecture/DECISIONS.md` — ADR-013..016
- `.agents/docs/handoffs/2026-09-13-learning-loop-and-improvement-review.md` — the first half of the day (superseded by this file)

## Next Steps

1. Operator: `sh-agent update` here and on the Mac; smoke checks from each real home; then the Known-limits sentence if wanted.
2. Refresh the two downstream mirrors from their own projects.
3. First real `/retro` in a project under the new skill; then `/retro --harness` here to consume anything queued.
