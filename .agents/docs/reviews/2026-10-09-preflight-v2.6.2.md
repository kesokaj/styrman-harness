---
owner: main
date: 2026-10-09
---
# Preflight — v2.6.2   (2026-10-09)

Range reviewed: v2.6.1..f69d63e

A PATCH release: it fixes released behaviour only. On macOS, `sh-agent dream` refused every run because it required agy's Linux token file, but agy keeps its login in the macOS Keychain (BUG-166). No skill, agent, rule or context file changed. The Go source change is `internal/dream/run.go` and its tests.

## Failures found during the gate

| # | Failure | Severity / confidence | Resolution |
|---|---|---|---|
| 1 | No real macOS dream run recorded, as `testing/guide.md` § 3e now requires (reviewer finding 1) | high / high | cleared. The operator's Mac run PASSED (`testing/2026-10-09-verify-v2.6.2-macos-dream.md`): `ok 1 proposals`, and the interactive agy still answered afterwards |
| 2 | No test pinned the `IsDir` check (reviewer 3b) | low / high | fixed in f69d63e. `TestAgyHomeMacKeychain` puts a file at Library/Keychains and goes red with the check dropped |
| 3 | The guide did not require the interactive-agy check (reviewer 3c) | low / high | fixed in f69d63e (`testing/guide.md` § 3e) |

## Reviews

- **Reviewer** (opus, code + docs): `reviews/2026-10-09-review-v2.6.1-0efa7c3-code-docs.md`, Status COMPLETE, verdict REQUEST CHANGES. Every finding is resolved in its `## Outcome`:
  - 1: cleared by the Mac run.
  - 2: rebutted. The operator's probe printed "Authentication required" verbatim; that output is now quoted in BUG-166.
  - 3a: rebutted on the number. agy gives up after its own 60 s login wait, not `cliTimeout`. The docs now say so.
  - 3b and 3c: fixed.
  - Severity: medium/low except 1, and 1 is cleared.
  - The fix commit f69d63e (3 lines of test, 3 doc lines and the reports) was read by me, not by a second reviewer round.
- **Earlier BUG-166 rounds** (on the uncommitted diff): `reviews/2026-10-09-review-bug166-agy-keychain.md` and `reviews/2026-10-09-security-bug166-agy-keychain.md`, both COMPLETE, every finding fixed or rebutted.
- **Security** (range touches secrets: the agy login): `reviews/2026-10-09-security-gate-v2.6.1-0efa7c3.md`, Status COMPLETE.
  - `## Not examined` names no release-diff path. Highest severity low.
  - 1, the shared Keychain login: accepted and documented.
  - 2, a locked Keychain costs about 60 s per project: left. Stopping the run at the first login failure is deferred with the timer item (handoff).
  - `govulncheck` (go1.27.2) reports nothing the code reaches. `go.mod` is unchanged.

## Verification

1. **Verifier log reaching HEAD:** `testing/2026-10-09-verify-v2.6.2-gate.md`, owner verifier, `Range reviewed: v2.6.1..f69d63e`, Status PASS. It ran `gofmt`, `go vet` and `go test -race -timeout 30m ./...`: 649 passed, 0 failed, 3 env-gated skips, unchanged since v2.6.1. It also ran the § 3e subset (280 passed) and the four BUG-166 tests with -v.
2. **What the release changed was run.**
   - AGENTS.md is unchanged in the range.
   - `testing/guide.md` § 3e adds the real macOS dream run with the interactive-agy check. It ran and passed: `testing/2026-10-09-verify-v2.6.2-macos-dream.md` (owner main, range reaching HEAD, rows 2 of 2 and 1 of 1).
   - No plan was added or modified; ROADMAP Now names none.
3. **No failure is left standing:** no FAIL or BLOCKED log in the range.

## Trail freshness

- The newest changes entry, `2026-10-09-1503-bug166-dream-agy-macos-keychain.md`, describes the range: pass.
- No future dates: pass.
- No testing log was deleted: pass.
- `changes/` and DECISIONS.md are append-only (no M in the range): pass.
- ROADMAP and TODO were not touched. Both are empty of open work, which is correct.

## Docs freshness

The docs lens rode in the reviewer dispatch. README, `specs/sh-agent/dream.md`, `deploy/local.md` and `testing/guide.md` are updated; context/, agents/ and skills/ do not mention the login. No open docs gap.

## Non-blocking findings — decisions

- Security 1 (dream's agy shares the real Keychain login on macOS): left: accepted, documented, and checked by the Mac run's interactive-agy row.
- Security 2 (about 60 s per project with a locked Keychain): left: operator. The fix would add behaviour (stop the run at the first login failure); it goes to the handoff with the deferred timer item.
- The Mac's global search index lists no projects, so `dream --all` there covers nothing: not a defect of this range. It goes to the handoff.

## What was not checked

- The nightly launchd timer against a locked Keychain (deferred by the operator, 2026-10-09). It is not a changed skill or agent and no blocking condition names it.
- `dream --all` over a non-empty global index on the Mac. `--project` ran the same code path past target selection.
- A second reviewer round over f69d63e (test and docs only, read by me).
- No Antigravity run of a changed skill or agent is owed: the range changes none.

VERDICT: SHIP
