---
owner: main
date: 2026-09-16
---
# Handoff — 2026-09-16 — harness security session done

Session ID: b12af1f5-5a11-41f7-8a72-867bf23337a4        (same machine: `claude --resume b12af1f5-5a11-41f7-8a72-867bf23337a4`; matches the feedback log's session field)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

The two downstream mirrors' security reviews of the a81ef4b refresh left five findings to this repository, recorded as BUG-033 to BUG-037 when the harness practices plan closed. The operator asked for the harness security session, BUG-037 first. All five are fixed, plus BUG-038, found on the way. Each ran under the bug skill: a failing test or reproduction first, the smallest fix, the reviewer on every fix (ADR-022), a verifier log, the record closed in its own fix commit.

## Completed

- **BUG-037** (d167cf8) — `sh-agent adopt . --force`, `adopt ..` or a self-adopt could delete `skills/` or the whole checkout. `Adopt` absolutizes the source, validates its name and refuses an overlap before any `RemoveAll`. Tests `TestAdoptRefusesDestructiveSources` and `TestWithin`.
- **BUG-036** (714ad23) — a symlink under the managed directories was deployed as its target. One boundary check, `checkManagedTree`, from `LoadInventory` and `RenderContext`; three review rounds, the second redesign chosen by the operator ("One boundary check"). Four tests.
- **BUG-038** (703c6cd) — the verifier ran `git checkout --` on tracked files and discarded the BUG-036 fix, then recovered it. A read-only agent never runs a git command that changes the tree: one line in the four read-only definitions and one clause in the context's Git paragraph.
- **BUG-035** (61344e9) — the `CLAUDE.md` guidance sent personal notes into `rules/`, which ships everywhere. Three sites now say to move only rules meant for every machine, never credentials or personal or host-specific notes.
- **BUG-033 and BUG-034** (dbde0aa) — `/review` and `/preflight` refuse a project agent named like a harness review agent, by file name or frontmatter `name:`, and `/preflight` blocks on an owed security report that is missing, has no `## Not examined` section, or names a release-diff path under it (D52, a narrow reversal of D35 for that row; both by the operator). Reproduced and re-tested in the sandbox on Claude Code; matrix rows `review-shadow`, `review-shadow-name`, `preflight-shadow`.
- **Real home:** this tree installed at dbde0aa (`binary dbde0aa7f769`, drift none) and the three smoke checks PASS on both CLIs, `testing/2026-09-16-verify-real-home-2.md`.

## Open / In Progress

- Nothing in flight. Open records: BUG-019 to BUG-023 (sweep findings in skills and references) and BUG-031 (a print-mode run, or an agent, can end its turn before what it started returns).
- The mirrors are at a81ef4b and do not carry this session's fixes (Next Steps 2 and 3).
- Requirements edits awaiting approval: none.

## Key Decisions

- **BUG-036:** one boundary check over per-reader guards, after two review rounds landed on the same mechanism — "One boundary check (Recommended)".
- **BUG-033:** refuse to dispatch on a clash, over renaming the harness agents or keeping the convention — "Refuse to dispatch on a clash (Recommended)".
- **BUG-034:** the security row only, a narrow reversal of D35 — "Security row only (Recommended)"; D52.
- **Deferred, agreed by both reviewers:** `adopt`'s copy-into-a-temp-directory-then-rename (BUG-037's record); a direct `RenderContext` test with a linked `context/` (BUG-036's record).
- **Rejected:** rewording `/review`'s description — the new check makes its "harness's own agents" claim true.

## Corrections

- None logged: `~/.config/sh-agent/feedback.jsonl` has no entry for this repository today.

## Agent notes (unverified)

- The shadow reproduction ran the project copy, and the model then declined the planted instructions and did a real review. That is luck; the check exists because the next shadow may not be so honest.
- Whether Claude Code resolves a subagent dispatch by file name or by frontmatter `name:` is unmeasured; the check matches both.
- BUG-029's failure shape recurred three times this session — a marker check that fired in the wrong state, two must-fail runs that failed to build rather than to assert, and one that passed with "no tests to run" — each caught before it counted. A retro candidate: a must-fail run requires its test's `=== RUN` line.
- The verifier could not see the sandbox transcripts until told they live under the sandbox home's `.claude/projects/`; its last run then confirmed the dispatch counts with `jq`.

## Environment State

- Branch: main, clean, pushed at dbde0aa plus this handoff.
- Running processes: none.
- Pending: nothing here. styrman-crew and argolis-devbox still hold their a81ef4b refresh commits, not pushed.

## Artifacts to Read

- `.agents/docs/bugs/INDEX.md` — BUG-033 to BUG-038, all fixed; the record of each fix in its Fix section.
- `.agents/docs/changes/2026-09-16-bug-037-adopt.md`, `…-bug-036-symlink-reads.md`, `…-bug-038-agents-git.md`, `…-bug-035-claude-md-guidance.md`, `…-bug-033-034-review-gate.md` — one entry per fix, each naming its agents, verification and review rounds.
- `.agents/docs/reviews/2026-09-16-review-bug-03*` — every review report, each with its `## Outcome`.
- `.agents/docs/testing/2026-09-16-verify-bug-03*` and `…-verify-real-home-2.md` — the verifier logs and the real-home check.
- `.agents/docs/testing/skill-matrix.md` — the three new rows.
- `.agents/docs/handoffs/2026-09-16-harness-practices-done.md` — the previous handoff; its Next Steps 1, 4, 6 and 7 still stand.

## Next Steps

1. **On the Mac:** `sh-agent update --check`, `sh-agent update`, `sh-agent status`, then the three smoke checks of `.agents/docs/testing/2026-09-03-verify-real-home.md` on each CLI. This box is installed at dbde0aa.
2. **Refresh both mirrors again** from this repository's HEAD, by each `VENDORED.md` recipe, so they carry the fixes for the findings they raised; commit locally as before. Then argolis-devbox: push (it deploys to the VM). styrman-crew: ship through a release per `versioning.md`; its handoff recommends v1.34.0.
3. **In this checkout:** `/retro --harness` for the two proposals from styrman-crew queued in `~/.config/sh-agent/proposals/2026-09-14-retro.md`, plus this session's candidates: BUG-029's recurrence, and whether BUG-032's and BUG-038's rules graduate.
4. **Still open from the previous handoff:** option C's filing clause in `/preflight` Step 4 (Antigravity filed its docs gaps nowhere a resume reads), and BUG-031's wait rule in the context and the agent definitions.
5. **The two deferred remediations,** as tasks when convenient: `adopt`'s temp-directory copy, and the direct `RenderContext` symlink test.
