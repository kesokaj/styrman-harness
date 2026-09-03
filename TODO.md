# TODO

Deferred work items. Not blocking current work, but should be addressed later.

<!-- Format: - [ ] Description (added YYYY-MM-DD) -->

## High Priority

- [ ] On the Mac, verify the Antigravity CLI picks up `~/.gemini/config/rules/AGENTS.md`, `agents/<name>/agent.md`, and the skills; if rules are not read from there, change the path in `internal/harness/plan.go` (added 2026-09-03)

## Medium Priority

- [ ] Decide whether `~/.agents/feedback/log.jsonl` on the Linux box (48 entries from the old harness) should be kept as seed data for `/retro` or deleted; the old `~/.agents/bin` scripts were deleted 2026-09-03 (added 2026-09-03)
- [ ] Record the fetch origin for the four Apache-2.0 skills in `skills/VENDORED.md` next time they are re-vendored (added 2026-09-03)

## Low Priority

## Ideas / Nice-to-Have

- [ ] `sh-agent update --check` from a shell prompt hook is deliberately not provided; revisit only if forgetting to update becomes a real problem (added 2026-09-03)
