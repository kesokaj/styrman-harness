---
owner: changelog
date: 2026-09-14
---
# BUG-015: plain-copy install no longer rebuilds     (2026-09-14)

## Summary

`sh-agent install` from a plain copy (a harness tree with no `.git`, i.e. the downstream mirrors styrman-crew `image/harness/` and argolis-devbox `deploy/harness/`) no longer rebuilds the binary with `go build`; it installs the running binary as is, prints `plain copy, not its own git checkout: the running binary is installed as is` (`[dry-run] … would be installed as is` under --dry-run), and prints no toolchain warning. From its own git checkout nothing changed. The rebuild gating is narrowed from "do we have `go`?" to "is this a git checkout?" preventing the rebuild-on-plain-copy failure observed in read-only downstream-mirror scenarios.

## Type

bugfix

## Files changed

- cmd/sh-agent/main.go — `shouldRebuild(goOnPath, child, ownCheckout)` gated by `harness.IsCheckout`
- cmd/sh-agent/main_test.go — `TestPlainCopyNeverRebuilds` (failed with "go.mod file not found" before fix), `TestShouldRebuild` extended with ownCheckout dimension
- internal/harness/repo.go — `ownCheckout` renamed to exported `IsCheckout`
- AGENTS.md — Rules for working here section updated with plain-copy behavior
- README.md — command table row 24 updated
- .agents/docs/deploy/local.md:27 — describing sentence moved with behavior (one medium finding from reviewer, fixed)
- .agents/docs/deploy/remote.md — Downstream mirrors section already accurate, cites BUG-015
- .agents/docs/testing/sh-agent/guide.md — new section 3a documenting plain-copy install messages

## Agents dispatched

verifier, reviewer (code), reviewer (docs), bugs, changelog

## Verification

Verifier log `.agents/docs/testing/2026-09-14-verify-bug-015.md`:
- `gofmt -l .` → exit 0, no output. PASS
- `go vet ./...` → exit 0, no output. PASS
- `go test -race ./...` → exit 0, all packages pass. PASS
- Scratch-home install and status with no drift, rebuild path intact for checkout. PASS
- Real binary built from checkout installed into read-only git-less mirror: printed "plain copy" line, no toolchain warning, created no `.agents/` in mirror, status correctly reported "plain copy" with exit 0. PASS
- Checkout path under --dry-run still announces rebuild. PASS
- Not exercised: sandbox CLI starts and real-home smoke checks (no context/skills/agents content changed); go-absent plain-copy variant.

## Docs

Four files with describing-text changes moved alongside the behavior:
- `AGENTS.md` — Rules for working here section (line 35)
- `README.md` — install command table row (line 24)
- `.agents/docs/deploy/local.md` — local setup section (line 27, fixed per reviewer docs lens)
- `.agents/docs/deploy/remote.md` — Downstream mirrors section, already updated

## Downstream

Both downstream mirrors (styrman-crew `image/harness/` and argolis-devbox `deploy/harness/`) must refresh to this commit or later. The prompt for that refresh is `.agents/docs/plans/2026-09-14-mirror-refresh.md` (committed in the following handoff commit).

## Requirements

BUG-015 is the manifest of a prevention rule in the bug record: "A behavior that depends on the environment the binary runs in must be designed against every deployment shape the trail names." Compliance with R1 (deployment shapes named in deploy guides) enforces the rule before the next similar defect lands. n/a for feature requirements.
