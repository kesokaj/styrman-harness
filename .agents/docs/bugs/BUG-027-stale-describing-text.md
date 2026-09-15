---
owner: bugs
date: 2026-09-15
id: BUG-027
status: fixed
area: skills and agents (describing text)
---
# BUG-027: plan's describing-text check missed restatements in other words

## Symptom
Done review found stale restatements the plan's edits left behind:
- `skills/init-project/SKILL.md:8` and `:66` still said Antigravity loads project rules, skills and
  agents, after I-19 recorded outcome C (Antigravity loaded no project skill or agent).
- `skills/handoff/SKILL.md:49` still sent deferred work to `TODO.md` after I-10.
- `agents/bugs.md:76` listed three graduation homes after I-18 made five.
- `skills/feature/SKILL.md:97` and `skills/bug/SKILL.md:53` kept a `documentation` trigger without
  "the plan" and "a deploy path" after I-25.

## Repro
Not reproduced as a running scenario; found by the Done review grepping each line against its rule's
intended current wording.

## Root cause
The plan's describing-text check grepped each replaced sentence's *old wording* (task 2.5's patterns).
A table row or a restated list says the same rule in different words and does not match a wording grep.

## Fix
Fixed at the cited lines (merged review F3, F4, F6, F7).

## Prevention
A rule change greps the rule's subject nouns — "deferred work", "Antigravity", "graduat", "`documentation`
when" — not only its old wording.

## Related
BUG-003, BUG-007, BUG-019 share this shape (a rule moves/changes and stale restatements survive because
the check looked for old wording, not the subject). Four occurrences now recommends this prevention rule
graduate into a permanent home — `context/AGENTS.md`'s describing-text rule or a `reviewer`/`review-standards`
checklist item — rather than being re-derived per bug record.
