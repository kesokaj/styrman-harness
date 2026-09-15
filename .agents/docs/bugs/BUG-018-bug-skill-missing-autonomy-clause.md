---
owner: bugs
date: 2026-09-15
id: BUG-018
status: open
area: skills/bug (gates)
---
# BUG-018: the bug skill's missing-repro gate no longer holds against an "autonomous" instruction

## Symptom
`skills/bug/SKILL.md:29` reads: "If there is no repro and the bug is not obvious from one file, **ask the user for one** ... Never invent a repro" — with no clause stopping this gate from being skipped on the strength of a "work autonomously" instruction. The feature skill keeps that clause for its own four gates (`skills/feature/SKILL.md:15`). A bug session run under an autonomous instruction therefore has no written rule requiring it to stop for a repro.

## Repro
`git show 87a862d -- context/AGENTS.md` shows commit 87a862d (2026-09-03, "Trim the context, move procedures into feature and bug skills, tier the reviewer") removing this line from the context: "Four gates need the user and are never skipped on the strength of an 'autonomous' instruction: scope confirmation for a feature, the cold-reader check before code, a missing repro, and a review or design loop that hits its round cap." Then `grep -n repro skills/bug/SKILL.md` shows the ask-for-a-repro sentence at line 29 with no autonomy clause. Both commands run by the main agent 2026-09-15.

## Root cause
The 2026-09-03 trim (87a862d) moved the four gates out of the global context into `skills/feature/SKILL.md` and `skills/bug/SKILL.md`. The feature skill's copy kept the autonomy clause; the bug skill's copy of the missing-repro gate did not. ADR-009 records the trim as "cut to under 70 lines with every rule kept," but no step checked each removed line for a landing spot in its new home.

## Fix
Not applied. Proposed in `.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md` (group G06): restore the clause at `skills/bug/SKILL.md:29` (e.g. "whatever autonomous instruction is in effect"), and give the `/review --docs` lens a walk over removed lines that flags a rule leaving with no home. G21 would later replace the clause with a full gates paragraph in the bug skill.

## Prevention
A change that trims or moves instructions lists each removed rule with its new home (file:line) or with why it was dropped, and a reviewer checks that list against the diff — do not rely on a line-count target ("under 70 lines") as the sole acceptance criterion for a trim.

## Related
BUG-003, BUG-007 (wording left behind after a rule moved elsewhere) — three records now share this pattern; the prevention rule should graduate to the project's `AGENTS.md` or a `.agents/rules/` file (e.g. "moved-rule checklist for trims"). `.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md` (group G06).
