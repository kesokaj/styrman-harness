---
owner: main
date: 2026-10-05
---
# Handoff — 2026-10-05 — real home on the checkout; critical-thinking work closed

Session ID: 805233b0-500a-424c-95dc-a250795fcb76
Project: /home/ostling-simon/projects/styrman-harness
Branch: main

## Context

Closes the two items of `handoffs/2026-10-05-1039-direct-instruction-clause.md`, both approved in the operator's own turns ("You fix both 1 and 2. Approved.", then "go option a").

## Completed

- Item 2: `agents/dreamer.md` and `skills/dream/SKILL.md` read the global rule files (`0be3b97`; `changes/2026-10-05-1120-dreamer-reads-global-rules.md`).
- Item 1: this Linux box's real home is now a CHECKOUT place, installed from `0be3b97` with `--repo` (it was the release place v2.3.2). Smoke PASS on both CLIs: `testing/2026-10-05-verify-real-home-critical-thinking.md`.
- Commit `6b6699c` is titled "Handoff: …" but holds only that smoke log: two git calls ran in parallel and the first lost the index lock. Pushed, so not rewritten; this file is the handoff.

## Open / In Progress

- Nothing in progress. Requirements edits awaiting approval: none.

## Key Decisions

- Option A over a v2.3.3 release: the box updates from this checkout from now on (`sh-agent update` pulls and rebuilds); the other places keep their kind.

## Corrections

- None logged (`~/.config/sh-agent/feedback.jsonl` does not exist on this machine).

## Agent notes (unverified)

- The smoke sessions may have registered `~/sandbox/smoke/*` clones in the real global search index; not checked.

## Environment State

- Branch: main, clean after this commit, pushed. Scratch: `.agents/scratch/{ct,dr}*`, `~/sandbox/{ct,dr,smoke}` (no credentials left).

## Artifacts to Read

- `.agents/docs/testing/2026-10-05-verify-real-home-critical-thinking.md` — the install and the smoke rows.

## Next Steps

1. The other places get the rule only when moved forward: a checkout place by `sh-agent update`; a release place only from the next release, which owes Antigravity runs of the changed `dreamer` and `dream` first (project release rule).
