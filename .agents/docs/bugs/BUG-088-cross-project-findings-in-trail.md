---
owner: bugs
date: 2026-09-28
id: BUG-088
status: fixed
area: process (research / handoff / issue lists)
---
# BUG-088: cross-project survey findings landed in this repo's trail and issue lists

## Symptom
The OKF research's consumer-survey agent read three other operator projects
(unifi-home-network, ranta-rankare, valkompass) only to answer whether any of
them uses OKF, but also reported problems it noticed inside those projects.
The drafter put three of those problems into
`.agents/docs/plans/2026-09-28-research-okf.md` §7 "found along the way"
rows #3-#5. The main agent then carried them, unchecked, into the committed
handoff `.agents/docs/handoffs/2026-09-28-research-okf.md` Next Steps 1, into
the change entry's mention of the table, and into a numbered "outstanding
issues" list given to the operator (items 5-7). Operator: "5-7 is not in this
project. Why are you bringing this up?"

## Repro
Not reproduced as a test; reproduced by reading the diff of
`.agents/docs/plans/2026-09-28-research-okf.md` §7 before and after the fix,
and the operator's own correction quoted above.

## Root cause
A cross-project survey was scoped as evidence for one question (does project
X use OKF). Its incidental findings about those projects were treated as
findings belonging to this repo, and nothing in the drafting or handoff step
asked "whose issue is this?" before the findings entered this repo's trail,
handoff Next Steps, and issue lists. The harness already separates
project-specific material by repository (`context/AGENTS.md` § Global vs
project), but that rule covers where customizations live, not where
FINDINGS ABOUT another project are recorded.

## Fix
The three rows were removed from `.agents/docs/plans/2026-09-28-research-okf.md`
§7, with a note explaining why. The already-committed handoff
`.agents/docs/handoffs/2026-09-28-research-okf.md` is left as-is (handoffs are
immutable); the next handoff states the removal. The correction is logged in
`~/.config/sh-agent/feedback.jsonl`.

## Prevention
A finding about another project is not recorded in this project's trail,
handoff Next Steps, or issue lists; a cross-project survey records only the
answer to the question it was run for. Before listing an issue, check which
repository it lives in.

## Related
`.agents/docs/plans/2026-09-28-research-okf.md` §7; handoff
`.agents/docs/handoffs/2026-09-28-research-okf.md`; no prior BUG record covers
this exact mechanism (nearest neighbor is BUG-012, which scopes "change
nothing" to source vs. trail — a different question).
