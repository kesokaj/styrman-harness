---
owner: bugs
date: 2026-09-17
id: BUG-061
status: fixed
area: agents/changelog (ADR authoring)
---
# BUG-061: changelog agent's ADR-031 draft asserted facts not in its dispatch or the diff

## Symptom
Dispatched to append ADR-031 (R9, the four procedures set in stone), the `changelog` agent wrote a draft that: (1) asserted "two other machines consume this harness" and named two consumer projects, against R8's rule that this repository names no consumer project; (2) said the project AGENTS.md "§ Rules for working here establish three strict rules" when the dispatch and diff put five rules in a new § How the harness works; (3) called cross-machine retro "deferred", a decision nobody made; (4) cited R9 by line numbers.

## Repro
Not reproduced as a test; found by the reviewer (finding 1, `.agents/docs/reviews/2026-09-17-review-77f47f7-map-r9.md`) catching (2), and by the main agent, while applying that fix, re-reading the ADR against the diff and finding (1), (3), and (4).

## Root cause
The agent drew on older ADRs' history and its own inference rather than confining itself to the dispatch prompt and the diff under review; nothing forced every sentence in the ADR to trace to one of those two sources, so plausible-sounding but ungrounded claims (consumer count and names, a rule count, a "deferred" status, line-number citations that drift) passed as fact.

## Fix
Main agent rewrote the ADR-031 text from the diff before commit, removing the consumer-project claim, correcting the rule count/section name, removing the fabricated "deferred" characterization, and citing R9 by name/section instead of line numbers.

## Prevention
An ADR or changes entry states only what its dispatch and the diff say — every factual claim traceable to one of the two, no fill-in from memory of older ADRs or other repos. The main agent reads an agent-written ADR against the diff itself before committing it, not only via the reviewer. A claim about another project's use of the harness never enters a decision record (R8); this is the second time R8's no-consumer-naming rule was violated by a subagent (see BUG-015 for the installer-side case) and reinforces graduating a "no consumer names in this repo" check into the `changelog` agent definition or a `.agents/rules/` file, plus a citation rule (name sections, never line numbers, since lines drift).

## Related
BUG-015 (an installer change traced to downstream copies; it predates R8, 2026-09-17), BUG-018/BUG-027 (trim/restatement drift caught only by grepping the diff), review `.agents/docs/reviews/2026-09-17-review-77f47f7-map-r9.md` finding 1.
