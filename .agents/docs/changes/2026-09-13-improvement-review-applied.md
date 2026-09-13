---
owner: changelog
date: 2026-09-13
---
# Improvement review applied — 20 items approved, 18 fully applied, 2 deferred     (2026-09-13)

## Summary
Improvement review panel of eight readers with 55 agents raised 46 proposals, 21 survived skeptical refutation. Operator approved all 20 final items ("go all"). Applied: items 1–17 and 19 (R7 clarification with ADR-016). Item 18 (vendor re-check on Mac + real home) completed in sandbox on Antigravity 1.2.2; real home and Mac are operator's scope. Item 20 (queue/index multi-machine location) deferred. Verification: fmt, vet, and test clean; 27 skills / 7 agents; sandbox e2e all pass on both CLIs; security round 1 found and fixed a git-env pollution issue; review round 2 is in progress.

## Type
task

## Files changed
- `.agents/docs/architecture/DECISIONS.md` — ADR-016 appended (R7 clarification)
- `REQUIREMENTS.md` — R7 sentence added (harness state vs project records scope)
- `ROADMAP.md` — "Now" section cleared after improvement review completion
- `TODO.md` — TODO item for item 20 (deferred: per-machine queue/index decision) added
- `cmd/sh-agent/main.go` — `ownCheckout()` added, `gitEnv()` pollution fix, status command prints binary stamp
- `internal/harness/repo.go` — checkout guard refactored, `gitEnv()` for every git command
- `skills/golang/references/` — 78 dangling links stripped, cut to stack scope (golang-ecosystem, golang-core, golang-testing, golang-io-web, golang-structure.md); golang-samber.md created for DI pattern
- `skills/docker/` — rewrote multi-stage Dockerfile example and amd64 build rule, deleted registry section
- `skills/frontend/` — rewrote TanStack Start config guidance (verified against styrman-crew SPA mode, client embedded, no Nitro)
- `skills/gcloud-setup/SKILL.md` — merged gcloud-auth-verification rules (project-discovery line removed); skill count 27
- `skills/gcloud-auth-verification/` — deleted
- `context/AGENTS.md` — removed redundant duplicate bullets, removed colima unconditional command, removed project memory from three agents, renamed `styrman` to `sh-agent` in Dev-setup section
- `.agents/docs/deploy/remote.md` — "Downstream mirrors" section added (styrman-crew, argolis-devbox), refresh instructions for stale vendored copies
- `context/dev-setup.md` — removed unconditional colima, removed arm64 parenthetical
- Skills templates (`assets/`) — init-project AGENTS.md and DECISIONS.md templates aligned with produced shape
- Feature skill — e2e never deploys; routes through preflight after Step 6
- Research skill — full mode now hands into feature skill at Step 1 (unified plan format)
- Feedback skill — logs corrections of order/procedure (not just substantive findings)
- `.gitignore` — added secret-shaped patterns for `.env*`, `*.pem`, `*.key`, `*.tfvars`, `*.tfstate`

## Agents dispatched
four editing agents (golang references, docker, frontend, Go binary), verifier, reviewer (×2), security (×2), bugs (BUG-010 — git-env pollution), changelog (this entry)

## Verification
- `gofmt -l .` — no output
- `go vet ./...` — clean (after git-env fix)
- `go test -race -count=1 ./...` — all pass, both packages (git-env pollution test added)
- `go run ./cmd/sh-agent install --home /tmp/scratch` — drift none, binary line stamped, no install errors
- `go run ./cmd/sh-agent list` — 27 skills listed, 7 agents, zero parse errors
- Sandbox e2e on both CLIs (Antigravity 1.2.2): context loaded, golang auto-triggered, explorer dispatched, smoke checks pass; evidence `.agents/scratch/e2e-items/`, Antigravity log `cli-20260913_220545.log` 0 parse errors
- Security round 1 review: `.agents/docs/reviews/2026-09-13-security-improvement-items.md` — FINDINGS: medium (git inheritance) fixed, low (gitignore) fixed, low (govulncheck) rebutted TODO
- Verifier log: `.agents/docs/testing/2026-09-13-verify-improvement-items.md` — PASS
- Reviewer round 2: `.agents/docs/reviews/2026-09-13-review-improvement-items.md` (reading; no file yet; reference as "see file")
- Security round 2: `.agents/docs/reviews/2026-09-13-security-improvement-items-round-2.md` (reading; no file yet; reference as "see file")

## Requirements
R7 (R7 clarified by ADR-016), R5 (manifest-only install, no hand-installed skill wiped)
