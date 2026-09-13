# TODO

Deferred work items. Not blocking current work, but should be addressed later.

<!-- Format: - [ ] Description (added YYYY-MM-DD) -->

## High Priority

- [x] On the Mac, verify the Antigravity CLI picks up `~/.gemini/config/rules/AGENTS.md`, `agents/<name>/agent.md`, and the skills — verified on Linux (2026-09-03) and on Mac (2026-09-04, `.agents/docs/testing/2026-09-04-verify-mac-real-home.md`)
- [ ] Operator: run `sh-agent update` (or `go run ./cmd/sh-agent install`) once on this Linux box and once on the Mac — the installed 2026-09-03 binary cannot rebuild itself (BUG-006); then `sh-agent status` must show a `binary` line without a warning (added 2026-09-13)

## Medium Priority

- [ ] Upgrade Go to ≥ 1.24.6 on the Linux box and the Mac (govulncheck: GO-2025-3956 in os/exec LookPath, GO-2026-4602 in os; both stdlib, fixed by the toolchain, not by code) (added 2026-09-13)
- [ ] sh-agent: a nested plain copy (a mirror inside another git repo) inherits the enclosing repo's commit in Inspect and status, and update could pull that repo; pin the git root (rev-parse --show-toplevel must equal the checkout) — improvement-review item, security round 3 residual (added 2026-09-13)

## Low Priority

## Ideas / Nice-to-Have

- [ ] `sh-agent update --check` from a shell prompt hook is deliberately not provided; revisit only if forgetting to update becomes a real problem (added 2026-09-03)
