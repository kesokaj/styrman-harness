---
owner: main
date: 2026-09-03
---
# Context trimmed, procedures moved to skills, reviewer tiered, sandbox prepared     (date: 2026-09-03)

## Summary

Applied the improvements the operator accepted after the honest take (plan section 11, ADR-009):

- `context/AGENTS.md` cut from 200 to 67 lines, every rule kept; compaction instruction and `/goal` line added; inventory embedded as names only (`RenderInventoryCompact`; `sh-agent list` keeps the full table).
- New skills `feature` and `bug` carry the work protocol and the bug path, model-invocable and callable as `/feature` and `/bug`.
- `reviewer` on Sonnet with a "what counts as a finding" section; `/preflight` dispatches it on Opus (Claude) or runs on `pro` (Antigravity). `changelog` template gains "Agents dispatched"; `bugs` records capped at 40 lines.
- Sandbox for live testing on this box: `/home/kesokaj/sandbox/home` (harness installed, git identity only) and `/home/kesokaj/sandbox/demo-app` (minimal Go HTTP service with a test, committed). The operator logs in to both CLIs there; the agent then drives them headless.

## Type

task

## Files changed

- `context/AGENTS.md` — rewritten short
- `skills/feature/SKILL.md`, `skills/bug/SKILL.md` — new
- `agents/reviewer.md`, `agents/changelog.md`, `agents/bugs.md`, `skills/preflight/SKILL.md` — tiering, finding bar, dispatch line, record cap
- `internal/harness/inventory.go`, `render.go`, `harness_test.go` — compact inventory
- `README.md`, `AGENTS.md`, `.agents/docs/plans/2026-09-03-rebuild.md`, `.agents/docs/architecture/DECISIONS.md` (ADR-009)

## Agents dispatched

none (harness work; verification run by the main agent)

## Verification

`gofmt`, `go vet`, `go test -race` clean; end-to-end install on a scratch home; see `testing/2026-09-03-verify-context-trim.md`.

## Requirements

R2, R5, R7.
