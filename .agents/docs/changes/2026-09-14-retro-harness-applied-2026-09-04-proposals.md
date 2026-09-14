---
owner: changelog
date: 2026-09-14
---
# Retro applied: harness 2026-09-04 proposals (date: 2026-09-14)

## Summary
Consumed the queued proposal backlog from the 2026-09-04 harness retro (styrman-crew session). Applied four proposals: added the retro-before-handoff sentence to the Learning loop of context/AGENTS.md; added two Hard rules bullets (behaviour change includes text that describes it; git diff/ls-files verification before delete); added checklist item "every new test can fail" to review-standards SKILL.md § 2. Context renders to 93 lines (within 100-line budget). Created first retro record under project trail.

## Type
task

## Files changed
- context/AGENTS.md — one sentence added to the Learning loop (retro before handoff) and two Hard-rules bullets; nothing moved
- skills/review-standards/SKILL.md — "every new test can fail" added as final item of § 2 Flow tracing
- .agents/docs/retros/2026-09-14-retro.md — new retro record (first under project-level retro rule)

## Agents dispatched
none

## Verification
go test -race -count=1 ./...: all green. Rendered context 93 lines (budget 100). Sandbox home installed and both CLIs quoted new order sentence and new git rule correctly (agy exit 0, claude exit 0); Antigravity lint 0 parse errors. Full outcome in .agents/docs/retros/2026-09-14-retro.md "Verification of the applied set" section.

## Docs
No README, specs, testing guides, or deploy guides changed — context rule evolution, documented in retro record and this entry.

## Requirements
R7.
