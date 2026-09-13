---
owner: changelog
date: 2026-09-13
---
# Review rounds 2–3: checkout guard closing     (2026-09-13)

## Summary
Second and third review rounds on the installer's checkout guard (Phase 4 entry). Round 2: `reviewer` APPROVE WITH NITS (prose nits fixed), `security` FINDINGS: the guard could be bypassed with a named path (`--repo`, `$SH_AGENT_REPO`) that is not a checkout because `FindRepo` fell through to the directory walk; it failed open on a corrupt manifest; the test shared a build-dir glob. Per the global rule on a second round against the same mechanism, the fix is a deletion: `FindRepo` now refuses a named path that is not a checkout (`TestFindRepoRefusesNamedNonRepo`); `knownCheckout` fails closed on a manifest error; the test checks only its own state. Round 3: all low severity, mechanism declared converged, no round 4. Residual findings noted but NOT fixed: first install has no manifest to compare; `git -C <dir>` can ascend to an enclosing repo (pre-existing, improvement review); `knownCheckout` trusts non-emptiness, relying on `FindRepo`'s invariant; status/diff still run git in the discovered checkout. Also from round 3: `Reinstall` now prints `building sh-agent from <repo>`; the vacuous build-dir assertion in `TestDryRunNeverBuilds` replaced by checks on what a build would leave.

## Type
task

## Files changed
- `internal/harness/repo.go` — `FindRepo` refuses named non-checkout path; `knownCheckout` fails closed on manifest error
- `cmd/sh-agent/main.go` — print line `building sh-agent from <repo>`
- `internal/harness/harness_test.go` — `TestFindRepoRefusesNamedNonRepo` (added); `TestDryRunNeverBuilds` checks build output not assertions
- `.agents/docs/reviews/2026-09-13-review-learning-loop-round-2.md` — new (reviewer APPROVE WITH NITS)
- `.agents/docs/reviews/2026-09-13-security-learning-loop-round-2.md` — new (security FINDINGS)
- `.agents/docs/reviews/2026-09-13-security-learning-loop-round-3.md` — new (security, all low, converged)

## Agents dispatched
reviewer ×1, security ×2

## Verification
Unit: `gofmt -l .` (nothing), `go vet ./...` (clean), `go test -race -count=1 ./...` (ok cmd/sh-agent, ok internal/harness); live: install from the crew mirror refused with the named non-checkout message; dry run prints the would-rebuild line and writes nothing.

## Requirements
R6 (installer safety) — guard refined to refuse bypass vectors; R7 (harness state) — unchanged.
