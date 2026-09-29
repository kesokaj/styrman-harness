---
owner: reviewer
date: 2026-09-29
---
# Review — memory-loop-docs   (2026-09-29)

Range reviewed: `1b8e84b^..HEAD` plus the working tree (commits 1b8e84b..HEAD plus uncommitted changes to `agents/dreamer.md` (new, `AM`), `skills/init-project/SKILL.md`, `skills/dream/SKILL.md`, `skills/feedback-log/SKILL.md`, `skills/retro/SKILL.md`, `skills/distill/SKILL.md`, `context/AGENTS.md`, `AGENTS.md`, `.gitignore`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/specs/sh-agent/memory-loop.md` (new), `README.md`, `.agents/docs/deploy/{local,remote}.md`).

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing in the diff is factually wrong against the code (`internal/harness/report.go` sends exactly `feedback.jsonl`, `retros.jsonl` and `proposals/`, never a directory walk of `~/.config/sh-agent/`, so every claim that the memory store is never reported checks out), and the slug rule, store path and settings-file name are stated identically everywhere they appear (`memory-loop.md`, `dream/SKILL.md`, `init-project/SKILL.md`, `skill-matrix.md`). The worst problem is that the new `dreamer` agent and the new memory store were wired into the mechanics (`/retro`, `/distill`, `init-project`, the map's file list) but not into the prose that already enumerates the learning loop's agents and the contents of `~/.config/sh-agent/` — three living docs restate that enumeration and none of them were extended, and `REQUIREMENTS.md`'s stated split between harness-wide state and project state doesn't have room for a project-scoped store parked inside the harness-state directory.

## Findings

- **[medium / confidence: high]** `REQUIREMENTS.md:31` — this line frames `~/.config/sh-agent/` as holding "the harness's own state (manifest, feedback log, the harness proposal queue, the cross-project index)" and draws the contrast that "a project's retro and distill records are the project's own state, in its trail." The new per-project memory store (`~/.config/sh-agent/memory/<slug>/`, one directory per project, described in `specs/sh-agent/memory-loop.md` as "the project's own state" in substance) fits neither side of that split: it is project-scoped like the trail, but lives under the harness-state directory for a stated reason (shared across CLIs, gitignored, machine-local). `REQUIREMENTS.md` was not touched by this change to say so. Per this repo's own hard rule, a trajectory change to what lives where goes through a proposed edit + ADR before `REQUIREMENTS.md` is updated — worth raising to the operator now rather than after the fact, since the architecture already shipped (Phase 1 committed) without that reconciliation.
- **[medium / confidence: high]** `context/AGENTS.md:41` (§ Learning loop), `skills/feedback-log/SKILL.md` § "Where the loop's outputs land", `skills/handoff/SKILL.md` § Cross-References — all three enumerate the learning loop's participants and where their records land ("run `/retro`... run `/distill`... `/dream` tidies..."; "Distill and dream records go there too"; "The `changelog` subagent writes..."), and none names the new `dreamer` agent even though this same diff has `/retro` (source 8) and `/distill` (Step 1) dispatch it and its artifact lands in the very `retros/` directory two of these three passages already describe. This is the exact case the dispatch's item (b) asks to check, and it fails in all three places, not one.
- **[low / confidence: medium]** `README.md` § How it works, `configDir` map node — the label enumerates the directory's contents (`feedback.jsonl, retros.jsonl, proposals/, state.json`) but was not extended to mention `memory/<slug>/`, which now also lives there. The report-back edge itself is correctly unaffected (same files sent), so this is a labeling gap on the node's content claim rather than a path that moved — worth a one-line addition to the node text, not a new edge.

## Criteria

1. Dead code — n/a (docs lens); no orphaned doc sections found, the changed passages are all still referenced.
2. Flow tracing — pass with the note above: `/retro` → `dreamer` → `retros/<date>-dreamer.md` → cited under `/retro`'s Findings is traced consistently across `agents/dreamer.md`, `skills/retro/SKILL.md`, `skills/distill/SKILL.md`, and `skill-matrix.md`'s `dreamer` row; the one break is the prose summaries that don't mention the new hop (finding 2).
3. Standards conformance — pass for the checked items: `sh-agent report`'s file allowlist (`feedback.jsonl`, `retros.jsonl`, `proposals/`) matches every doc claim (`memory-loop.md`, `remote.md`, `README.md`); the slug/store/settings-file wording is identical across every surface it appears in (item e); fails on the "every restatement changed with it" rule for the `dreamer` mentions (finding 2) and arguably for `REQUIREMENTS.md`'s state split (finding 1).
4. Directory layout — n/a (docs lens); the spec's store path and the skill files agree on where things land.
5. Simplicity — n/a (docs lens); no new dependency or hand-rolled mechanism described.
6. Decoupling — n/a (docs lens); `/retro` and `/distill` consume the dreamer's artifact rather than the conversation, matching the project's subagent contract (context § Work).

## Not examined

- The uncommitted plan file `.agents/docs/plans/2026-09-29-dreamer-memory-loop.md` itself (only grepped for `REQUIREMENTS`/ADR references, not read in full) and its decisions D13–D16.
- `internal/harness/inventory.go` and the `## Installed` generation path (excluded per the dispatch's item (d): agent-count changes there are not a finding).
- Whether Claude Code's actual runtime behavior matches the "Trust" bullet in `memory-loop.md` (a workspace-trust claim about Claude Code internals) — unverified from this repo; taken on faith as stated.
- `.agents/docs/plans/2026-09-29-dreamer-memory-loop.md`'s cold-read review records and whether findings 1–3 above were already raised and answered there (I did not open the `2026-09-29-review-dreamer-memory-loop-coldread-r*.md` or `2026-09-29-review-worktree-docs.md` records in `.agents/docs/reviews/`, several of which look like they may cover this exact ground from an earlier round).
- `.agents/docs/deploy/local.md` and `remote.md` beyond the two changed lines quoted in the diff — the rest of both files was not re-read for drift.

## Outcome
- `REQUIREMENTS.md:31` names no memory store — fixed with the operator's approval ("Approve the edit (Recommended)", 2026-09-29): "its auto memory store in `.agents/memory/` (gitignored, one per checkout, never sent by `sh-agent report`)"; ADR-043. The store moved inside the project on the operator's "Option 2" between this review and the fix (plan D17), which is why the clause names `.agents/memory/` rather than the reviewed path.
- `dreamer` absent from the loop's enumerations — fixed: `context/AGENTS.md` § Learning loop ("`/retro` reads it first and dispatches `dreamer` …"), `skills/feedback-log/SKILL.md` § Where the loop's outputs land, `skills/handoff/SKILL.md` § Cross-References; size test 21,855 bytes.
- README map `configDir` label omits `memory/<slug>/` — rebutted after D17: the store is no longer under `~/.config/sh-agent/`, so the label is accurate as it stands; the `documentation` agent removed the two README statements and the two deploy-guide clauses that had placed it there.
