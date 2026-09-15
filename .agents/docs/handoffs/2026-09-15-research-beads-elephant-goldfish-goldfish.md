---
owner: main
date: 2026-09-15
---
# Handoff — 2026-09-15 — research: beads, elephant-goldfish and goldfish against the harness

Session ID: b12af1f5-5a11-41f7-8a72-867bf23337a4        (same machine: `claude --resume b12af1f5-5a11-41f7-8a72-867bf23337a4`)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

The operator asked for three repositories (gastownhall/beads, vshvedov/elephant-goldfish, lukacf/goldfish) to be researched against the harness for drift, more tasks and phases, smaller verifications and planning, and for a table of what to add or not. Research only: nothing in `context/`, `agents/`, `skills/` or the Go code changed.

## Completed

- Commit `3803b6c` (pushed): the research record `.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md`. It holds 40 groups over 73 verified candidates (11 add now, 15 add next, 6 later, 8 skip), a suggested first bundle, drift found in the harness itself, cost, conflicts, corrections from verification, and what was not verified.
- Same commit: BUG-017 (`.agents/docs/bugs/BUG-017-workflow-resume-reran-proposers.md`, INDEX row) and the changes entry `.agents/docs/changes/2026-09-15-research-beads-elephant-goldfish.md`.
- This handoff's commit: BUG-018 (open), the regression the research found in the bug skill's repro gate, and the cold-reader brief `.agents/docs/plans/2026-09-15-coldread-handoff-2026-09-15.md` (no missing files, no contradictions).

## Open / In Progress

- No group is implemented and no plan exists; ROADMAP "Now" and TODO.md stay empty until the operator picks groups.
- Harness drift found and not fixed, each tied to a group in the record: the missing-repro gate lost its autonomy clause in 87a862d (BUG-018, open; fixed by G06); seven operator boxes stay open in two finished plans (G05); the vendor-precedence claim at `context/AGENTS.md:37` is untested for Claude Code skills (G26); the repo setup line chains `gofmt -l`, vet and test (G07); skill and agent counts are hand-typed (G08).
- Carried unchanged from the 2026-09-14 handoff: `/retro --harness` for the two queued proposals in `~/.config/sh-agent/proposals/2026-09-14-retro.md`; the argolis-devbox mirror refresh; the sandbox CLI-start matrix for cceee54.
- Requirements edits awaiting approval: none.

## Key Decisions

- Neither beads nor goldfish is adoptable as a tool: beads needs a Dolt database, git hooks and a settings.json hook; goldfish is an MCP server with SQLite and a daemon. Only their practices transfer, as text.
- The tiers follow the set judge, which saw cumulative cost: context edits land on existing lines, and the plan-template additions merge into one Decisions section and one Checkpoint box instead of separate sections.
- G22 (reviewer on every bug fix) amends a recorded threshold and G24 (bug index compaction) reopens an approved refusal; both wait for the operator.

## Corrections

- none (no feedback-log entry for this project today).

## Agent notes (unverified)

- The Claude Code skill precedence behind G26 comes from a documentation summary one verification lens read, not a live test.
- The 60-row styrman-crew bug index behind G24 was counted by a lens and not re-checked.
- Suggested build order from the record: G06, G01, G02, G03 with G04, G05, G07, G08, then G09, then G10 with G11.

## Environment State

- Branch: main, clean at `3803b6c` plus this handoff's commit; pushed.
- Running processes: none.
- Scratch (gitignored): `.agents/scratch/research/` holds the three shallow clones, `briefs/`, `batches/`, `harness-map.json` and the workflow extracts; `.agents/scratch/stage-research.txt` is the staging list.

## Artifacts to Read

- `.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md` — the table and everything behind it.
- `.agents/docs/bugs/BUG-017-workflow-resume-reran-proposers.md` — how to reuse finished workflow agents without a resume.
- `.agents/docs/changes/2026-09-15-research-beads-elephant-goldfish.md` — what this task recorded and which agents ran.
- `.agents/docs/bugs/BUG-018-bug-skill-missing-autonomy-clause.md` — the open regression, third of its kind after BUG-003 and BUG-007.

## Next Steps

1. Operator: choose which groups to build from the research record, and decide G22 and G24.
2. Then `/feature` over the chosen groups as one plan, with G06 first so the docs lens reviews the edits that follow.
3. BUG-018 is the third rule lost when wording moved (BUG-003, BUG-007): when G06 lands, graduate its prevention rule into this repo's AGENTS.md.
4. Carried: `/retro --harness` for the two queued proposals; argolis-devbox mirror refresh; sandbox CLI-start matrix on cceee54.
