---
owner: bugs
date: 2026-09-15
id: BUG-021
status: fixed
area: skills/distill
---
# BUG-021: distill's feedback.jsonl entry omits the `project` field, breaking retro's cross-project filter

## Symptom
The feedback entry `distill` appends to `~/.config/sh-agent/feedback.jsonl` (`skills/distill/SKILL.md:158-169`) has no `"project"` field — only `ts`, `session`, `category`, `signal`, `original`, `corrected`, `user_quote`, `rule_ref`. `skills/feedback/SKILL.md:13,43` says every entry carries `project`, and is what lets `/retro`/`/distill --from-log` spot cross-project patterns. `skills/retro/SKILL.md:65-66` excludes this project's own entries with `grep -vE '"project": ?"<this project>"'`; an entry with no `project` field is never matched and never excluded, so a same-project correction is counted as a cross-project hit and can promote a rule to the harness falsely. Found by the 2026-09-15 docs-lens sweep, F21.

## Repro
Not reproduced as a failure in a live run; found by review — the JSON shape at `skills/distill/SKILL.md:158-169` was compared against `skills/feedback/SKILL.md:43` and the exclusion filter at `skills/retro/SKILL.md:65-66`.

## Root cause
`distill`'s own feedback-append example was written independently of the `feedback` skill's entry shape and never included the `project` field the retro filter depends on.

## Fix
Fixed 2026-09-16 (Phases 2-3, clear the board): `skills/distill/SKILL.md`
now defers to `skills/feedback/SKILL.md` for the entry's shape instead of
carrying its own drafted example missing `project`.

## Prevention
Every writer of `~/.config/sh-agent/feedback.jsonl` uses the `feedback`
skill's entry shape rather than a locally drafted example. Graduated:
`skills/distill/SKILL.md`.

## Related
2026-09-15 docs-lens sweep, `.agents/docs/reviews/2026-09-15-review-sweep-merged.md` F21.
