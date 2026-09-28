---
owner: main
date: 2026-09-28
---
# OKF record: other projects' findings removed (BUG-088)     (date: 2026-09-28)

## Summary
The operator corrected a scope leak. Problems found inside unifi-home-network, ranta-rankare and valkompass during the OKF consumer survey had been listed as this repo's issues:
- in the research record's found-along-the-way table;
- in the 2026-09-28 handoff's Next Steps;
- in an outstanding-issues list given to the operator.

The three rows are removed from the record, with a note explaining why. The survey's conclusion that no project uses OKF stays, because the recommendation rests on it. The committed handoff is immutable and is superseded by today's second handoff. This entry supersedes the mention of rows #3–#5 in `changes/2026-09-28-research-okf.md`.

## Type
bugfix

## Files changed
- `.agents/docs/plans/2026-09-28-research-okf.md` — §7 table: rows #3–#5 removed; its heading says "in this repo" and explains the removal.
- `.agents/docs/bugs/BUG-088-cross-project-findings-in-trail.md` — new record. Prevention rule: a finding about another project is not recorded in this project's trail, handoff Next Steps or issue lists.
- `.agents/docs/bugs/INDEX.md` — a row for BUG-088.

## Agents dispatched
- `bugs`: BUG-088.

## Verification
Docs only. The main agent read the record's diff (1 line added, 4 removed) and the bug record. No verifier or reviewer ran, because no source changed. The correction is logged in `~/.config/sh-agent/feedback.jsonl`.

## Docs
No docs changed.

## Requirements
n/a
