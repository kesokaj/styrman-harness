---
owner: main
date: 2026-09-16
---
# Preflight gate — v1.0.0   (2026-09-16)

Range reviewed: 4aa19be..0a67f3f

## VERDICT: SHIP

The first release this repository has cut. Base ref `4aa19be` is the commit both downstream consumers carry, so the release diff is exactly what they have not seen: 17 commits, 75+ files. One blocking finding was raised and cleared; the mechanism behind it was withdrawn rather than patched a third time.

## What was blocking, and what happened to it

**The release review returned REQUEST CHANGES on a high finding, and it was right.** `go.mod` had gained `toolchain go1.25.8` earlier the same day to clear two standard-library advisories. Under the default `GOTOOLCHAIN=auto`, a pin is honoured by DOWNLOADING the toolchain, so on a host with neither that toolchain cached nor a reachable proxy every `go` invocation in the module exits 1 with `toolchain not available`. That is the shape of a consumer host, and `sh-agent install` and `update` both shell out to `go build`.

The first fix set `GOTOOLCHAIN=local` inside `Reinstall`. Measured on the real path, it did not hold: the OUTER `go run` is itself subject to the pin and fails before any of this code executes. The round-2 verifier and the main agent reproduced that independently and reported the same two exit statuses.

**Resolution: the pin is withdrawn (ADR-029).** Three rounds of findings on one mechanism is where this project's own rule says to stop patching and ask what the mechanism is for. The answer was that it bought patched builds in one checkout at the price of a build-time network dependency for everyone else. `GOTOOLCHAIN=local` stays inside `Reinstall` on its own merits — an installer must never need to download a compiler.

## The advisories this release ships with, and why that is the right call

`govulncheck ./...` reports two standard-library advisories as affecting this code. The security review read both advisory records rather than their titles and established that **neither is reachable here**:

| Advisory | Reported in | Why it cannot fire in this program |
|---|---|---|
| GO-2026-4602 | `os`, via `LoadSkills` | It is an escape from an `os.Root`. This module constructs none — no `os.Root`, `os.OpenRoot` or `os.DirFS` anywhere in it. |
| GO-2025-3956 | `os/exec`, via `Reinstall`'s `LookPath` | It triggers on `LookPath` receiving `""`, `"."` or `".."`. The only call passes the literal `"go"`. |

Both are symbol matches. Against them stands a pin that made the tool unbuildable offline, and an unbuildable tool is the worse property, because it is how every later fix reaches a machine. This exception is recorded here and in ADR-029 precisely so a later session does not re-add the pin on the strength of a `govulncheck` line. The standing fix remains an operator-run Go upgrade: they are fixed in go1.24.6 and go1.25.8, Debian trixie ships `golang-1.24-go 1.24.4-1` and nothing newer, and the toolchain is root-owned.

## Step 2 — the three independent rows

| Row | Verdict | Artifact |
|---|---|---|
| Verifier | PASS, two rounds plus an appended correction | `testing/2026-09-16-verify-preflight-v1.0.0.md` |
| Reviewer (code, opus) | REQUEST CHANGES → **lifted** in round 2 | `reviews/2026-09-16-review-preflight-v1.0.0.md` |
| Security | FINDINGS, none blocking; round 2 says ship | `reviews/2026-09-16-security-preflight-v1.0.0.md` |

Every finding across both rounds is answered in its report's `## Outcome` section — fixed, rebutted with a checkable reason, or left with a named owner. Fifteen were fixed in code or docs during the gate, including four the reviewers proved by probe: a `defer` that deleted the prepared copy on every `adopt` failure path, a mode mask that never reached an already-deployed file, a `secretShaped` guard on the exit but not the entrance, and a release-procedure pathspec that missed a source change in a sibling directory.

Two findings are `left` with owners: the consumer's select-then-verify tag selection has no walk-down and no version floor, which is styrman-crew's code and is written into its channel brief; and the operator's Go toolchain upgrade.

## Step 3 — trail freshness

Newest `changes/` entry describes this work. No artifact carries a future date. `changes/` and `DECISIONS.md` are pure appends across the range — zero deleted lines in `DECISIONS.md`. `ROADMAP.md` and `TODO.md` were not touched since the base ref, which is correct rather than lagging: both are empty by this project's rule and have been all day.

`REQUIREMENTS.md` changed in the range. R4 gained a sentence about `update` refusing a dirty checkout, asked for and approved by the operator in their own turn, recorded as ADR-027 and committed with the requirement edit.

## Step 4 — docs freshness (advisory, non-blocking)

`reviews/2026-09-16-review-4aa19be-HEAD-docs.md`, APPROVE WITH NITS. Two low findings: the README and AGENTS.md do not mention two new user-reachable `install` refusals, and `agents/verifier.md` keeps `gofmt -l .` as the generic default with its exception clause beside it. Both filed for ordinary work; neither blocks.

## Final gate measurements at 0a67f3f

```
gofmt -l cmd internal                → exit 0, no output
go vet ./...                         → exit 0
go test -race -v -count=1 ./...      → exit 0; pass 109, fail 0, skip 1
install into a scratch home          → 260 written; drift none; 27 skills, 7 agents
offline build (empty GOMODCACHE,
  GOTOOLCHAIN=auto, GOPROXY=off)     → exit 0
govulncheck ./...                    → 2 affecting, both unreachable (above)
```

The one skip is `TestFindRepo`, pre-existing and unrelated. Test count rose 92 → 109 across the gate as findings were fixed with tests.

## What was not checked

- **No third review round.** The round-2 fixes were applied by the main agent, re-read, and two were proven by reproducing the reviewers' own scenarios end to end. No reviewer or security agent read them.
- **Antigravity's real-home skill-load and agent-dispatch checks.** Deployment and discovery are verified from the real home; execution is verified only in the sandbox. The permission grant needed is the operator's and this session's classifier refused it three times.
- **The consumer side.** Nothing in styrman-crew or argolis-devbox was exercised. The channel that will consume this tag does not exist yet.
- **The release procedure end to end.** This tag is its first exercise; its commands were rehearsed in scratch repositories, must-pass and must-fail, not against a real release until now.
- The Mac, and the full skill matrix.
