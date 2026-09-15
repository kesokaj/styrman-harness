---
owner: main
date: 2026-09-15
---
# Handoff — 2026-09-15 — harness practices G01–G26: the plan is at the cold-read cap, waiting on the operator

Session ID: b12af1f5-5a11-41f7-8a72-867bf23337a4        (same machine: `claude --resume b12af1f5-5a11-41f7-8a72-867bf23337a4`)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

After the research of 2026-09-15 (3803b6c), the operator chose all 26 add-now and add-next groups and started the feature: "go with the first bundle g01-g26, start /feature", then "go all, use your default thresholds for G20". The plan went through three cold-read rounds of three readers each without both tokens reading ready, so the feature skill's cap gate stopped the work for the operator's decision. No harness file has changed.

## Completed

- Plan `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`: 24 interfaces as 88 fenced edits in application order, revised after round 1 (c71aacf) and round 2 (55a860f). A rehearsal applies all 88 to copies of the files at `c192622` with 0 failures, and all 81 task expectations hold against the rehearsed end state (D15).
- Cold-read reports: `.agents/docs/reviews/2026-09-15-review-coldread-harness-practices-r{1,2,3}-{a,b,c}.md`.
- `ROADMAP.md` "Now" names the plan.

## Open / In Progress

- The plan's first unticked box is the Cold-reader check, noted in progress. Its two cap `Open:` rows hold the question and the recommendation: cut the release gate's coverage rule (I-21a–c), which drew findings in all three rounds, and resolve the listed narrow items by recommendation.
- Carried unchanged: `/retro --harness` for the two queued proposals in `~/.config/sh-agent/proposals/2026-09-14-retro.md`; the argolis-devbox mirror refresh; the sandbox CLI-start matrix for cceee54.
- Requirements edits awaiting approval: none.

## Key Decisions

- D1–D34 in the plan's `## Decisions`; the round-3 residue is its two cap `Open:` rows.

## Corrections

- none (no feedback-log entry for this project today).

## Agent notes (unverified)

- The rehearsal script is scratch (`.agents/scratch/rehearsal/rehearse.py`), not committed (ADR-010); the plan's "How edits are written, applied and dispatched" section specifies what it does, so it can be rebuilt.
- In-session agents on Claude Code load the installed definitions, not this checkout's (no `.claude/` here); D9 says how the plan's dispatches carry changed text.

## Environment State

- Branch: main; 55a860f pushed; this handoff, the plan's cap note and the round-3 reports ride in the next commit.
- Running processes: none.
- Scratch (gitignored): `.agents/scratch/research/`, `.agents/scratch/rehearsal/`, `.agents/scratch/coldread-r{1,2,3}.{json,txt}`.

## Artifacts to Read

- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — the plan; start at its Phases section and the cap `Open:` rows.
- `.agents/docs/reviews/2026-09-15-review-coldread-harness-practices-r3-c.md` — the round-3 reader with the two blocking questions on the release gate and the naming test.
- `.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md` — the research behind the groups.

## Next Steps

1. Operator: decide at the cap — cut or keep the release gate's coverage rule, and whether Phase 1 starts after the narrow fixes or after a fourth cold-read round.
2. Then apply the decision to the plan, re-run the rehearsal and the expectation checks, tick the Cold-reader check box with the operator's words quoted, and start Phase 1.
