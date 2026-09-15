---
owner: bugs
date: 2026-09-15
id: BUG-022
status: open
area: skills/golang (references)
---
# BUG-022: golang-ecosystem reference tells agents to register private code with third-party indexers

## Symptom
`skills/golang/references/golang-ecosystem.md:1030`: "**Register for discoverability** — add to Context7, DeepWiki, OpenDeep, zRead. Even for private libraries." This sends private source to third-party indexing services. Found by the 2026-09-15 docs-lens sweep, F51.

## Repro
Not reproduced as an executed action; found by review — the line was grepped and read in context.

## Root cause
Imported third-party guidance (the golang-ecosystem reference) kept a publishing/discoverability step without checking whether it applies to private, non-public code.

## Fix
Not fixed. This task only records the finding; no source changed.

## Prevention
Imported skill text is checked for steps that send code or data to an external service before it is vendored or kept; a step that "registers" or "publishes" project content to a third party is flagged and scoped to public/OSS work only, never applied blanket to "even private libraries."

## Related
2026-09-15 docs-lens sweep, `.agents/docs/reviews/2026-09-15-review-sweep-merged.md` F51.
