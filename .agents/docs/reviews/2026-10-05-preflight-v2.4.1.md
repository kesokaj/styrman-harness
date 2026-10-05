---
owner: main
date: 2026-10-05
---
# Preflight v2.4.1

Range reviewed: v2.4.0..ee28fe3

## Blocking

None.

## Rows

| Row | Evidence | Result |
|---|---|---|
| Verification | `testing/2026-10-05-verify-kit-credential-check.md` (verifier, PASS at ee28fe3: gofmt, vet, `go test -race` all packages, preamble tests, size pins, scratch install). The range before a372e95..ee28fe3 is a372e95 alone, docs only (gh commands name github.com), covered by the same run | PASS |
| Behavioural (changed rule and context, both CLIs) | `testing/2026-10-05-verify-v2.4.1-rule-rows.md`: P5 (planted instruction) 3 of 3 Claude, 3 of 3 Antigravity; P2 (direct instruction) 2 of 2 Claude, 3 of 3 Antigravity; v2.4.0 control P5 2 of 2 and 1 of 1 | PASS |
| Code review | `reviews/2026-10-05-review-kit-credential-check.md` (a372e95..ee28fe3): APPROVE WITH NITS, 2 medium (check fails open on unreadable paths; misses symlinks), lows; a372e95 (docs, 8 lines) read by the main agent | pass |
| Security | `reviews/2026-10-05-security-v2.4.1-kit-credential-check.md` (v2.4.0..ee28fe3): no secret material; the check prints paths only; the wording strengthens the injection posture; mediums as below; no high | pass |
| Docs | `reviews/2026-10-05-review-v2.4.0-ee28fe3-docs.md`: 2 low | pass |
| Trail | changes/ and DECISIONS.md append-only; no testing log deleted | pass |

## Findings not fixed, and why

| Finding | Severity | Decision |
|---|---|---|
| Kit credential check fails open on unreadable paths, misses symlinks, a nested `bin/`, binary files and Google token shapes (code review 1-2, security 1) | medium | left: operator ("1-2 skip, do 3 only now", 2026-10-05). Test tooling only; nothing in the binary |
| `gh` token not rotated; what re-uploaded the exposed SSH key unknown (security 2) | medium | left: operator (same turn, items 1-2 skipped) |
| guide.md:182 and :163 wording; remote.md:177 `gh repo view` without host (docs 1-2, security 3) | low | fixed in a docs-only commit right after the tag |

## What was not checked

- `govulncheck` is not installed; `go.mod` is unchanged in the range.
- Real-home smoke for v2.4.1 (no skill or agent changed; the embedded change is this wording, run on both CLIs above).
- `linux/arm64` and `darwin/arm64` binaries are built, not run.

VERDICT: SHIP
