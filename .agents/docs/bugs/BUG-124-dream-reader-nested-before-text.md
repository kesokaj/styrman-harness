---
owner: bugs
date: 2026-10-01
id: BUG-124
status: fixed
area: search/dream reader
---
# BUG-124: a tool call's arguments in a string-only protobuf sub-message reached the dreamer's prompt (SD1)

## Symptom
Tool-call JSON arguments, with the field key and length bytes glued in front, were kept as a prose leaf and sent to the dreamer.

## Repro
`proseLeaves(pb(10, pb(1, []byte(<JSON, 32-126 bytes>))))` kept the JSON with its key and length bytes in front.

## Root cause
`walk()` in `internal/search/dream.go` tried a length-delimited field as text before trying it as a nested message. A sub-message holding only a string is all printable bytes, so `textField` accepted it, and `prose()`'s JSON guard failed on the glued leading byte instead of seeing the JSON.

## Fix
Nested parse first, text leaf only when the nested parse fails (`internal/search/dream.go:352-365`). `TestProseLeaves` gained the two shapes. `TestAgyCorpusLeaves` measured 971 leaves and 0 argument keys on 95 real databases before and after.

## Prevention
A schemaless decoder tries the structured parse before the text reading. A corpus measurement ("0 of 971") is evidence about one corpus, never a structural guarantee: the structural case gets its own unit fixture.

## Related
Audit SD1 in `plans/2026-09-30-recon-audit-fix.md`; plan `plans/2026-09-30-audit-fix.md` Phase 1; verify log `testing/2026-09-30-verify-audit-fix-phase-1.md`. Diagnosis was the audit itself (goldfish dispatch skipped: location mechanically obvious).
