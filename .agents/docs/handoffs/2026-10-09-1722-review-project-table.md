---
owner: main
date: 2026-10-09
---
# Handoff — 2026-10-09 — project review done, fixes await the operator's pick

Session ID: 728a4838-8b3c-4819-a62d-2b01e563a77f
Project: /home/ostling-simon/projects/styrman-harness
Branch: main

## Context

`/review --project --all --high` over c267990 (v2.6.2, the live tag). Six lens agents ran; their artifacts and the merged record are committed. No source, test, doc or config changed.

## Completed

- The record: `.agents/docs/reviews/2026-10-09-review-project-merged.md` — verdict per lens, 29 numbered rows highest priority first, observations, Not checked. Change entry: `.agents/docs/changes/2026-10-09-1722-review-project-all.md`.
- Commands: all PASS (gofmt, vet, race suite 649/0/3, scratch install, list, govulncheck clean from the full path).

## Open / In Progress

- Nothing in progress. The fixes have not started: they wait for the operator to approve rows by number from the merged record.
- Requirements edits awaiting approval: none.

## Key Decisions

- Review only, as asked: the fixes are planned from the table, not applied here.
- Row 7 (the agents' "at most ten lines" cap) is recommended as strike; rows 21, 23, 24, 25, 28 as defer; rows 26 and 27 as strike (accepted risk on file; no trust boundary).

## Corrections

- None.

## Agent notes (unverified)

- The code lens did not reproduce row 2 (the dream grandchild case); its mechanism is the documented `os/exec` wait on pipe copies.
- The security lens did not test whether Claude Code honours an inline `mcpServers` in a user-level agent file (row 1).

## Environment State

- Branch: main, clean after this commit, pushed. Tag v2.6.2 unchanged. `.agents/scratch/verify-project/` holds the verifier's test.json and coverage files (gitignored); `.agents/scratch/home` holds the scratch install.

## Artifacts to Read

- `.agents/docs/reviews/2026-10-09-review-project-merged.md`: start here; each row names its lens artifact.

## Next Steps

1. Operator: pick rows from the merged record (`go approved …`). Recommended first batch: rows 1–4 (the two code mediums, the delivery spec, the exit-2 vocabulary), then 5–6 and 8–13 as one docs-and-small-fixes task, then 14–20 as one simplify task. Each batch is a Task under the Definition of Done: fix, test that fails without the fix, verifier, round 2 for the security medium (row 1), Outcome appended to the record.
2. Operator decision carried from the v2.6.2 handoff: whether `sh-agent dream` stops the run at the first agy login failure, and whether the Mac's checkouts should be dreamt over (run `sh-agent search <word>` once in each).
