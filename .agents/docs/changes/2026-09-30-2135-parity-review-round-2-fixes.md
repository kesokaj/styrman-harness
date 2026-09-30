---
owner: main
date: 2026-09-30
---
# Parity feature: review round 2 fixed (code, security)     (2026-09-30)

## Summary
Round 2 over f1ca525..e2cffc1 confirmed round 1's fixes and found:
- **Antigravity tool-call arguments** (security, medium): a reply step's tool-call arguments (JSON objects: commands run, files written) still reached the prompt; they are now dropped. On a real conversation copy, 32 argument runs dropped and 21 prose runs kept.
- **Tests owed** (code, medium): two round-1 fixes had no test that fails without them. Both have one now, each shown by a must-fail mutation: an Antigravity read error keeps the watermark, and a failed state write is exit 2.
- **Lows**:
  - start-of-run cleanup spares copies under an hour old (a dry run takes no lock);
  - the symlink check covers `.agents` itself;
  - `targets()` loses unreachable branches, and `Run` refuses usage errors (`TestRunUsage`);
  - the timer keeps only absolute PATH entries;
  - `answer()` needs a SUCCESS result;
  - agy's private settings set `toolPermission` explicitly (a real sandbox run answered);
  - the step types' source is recorded (plan D8b);
  - the spec states the cross-vendor flow and `--setting-sources local`'s auth limit.

Also here: BUG-116 to BUG-122, the records of what went wrong in the feature.

## Type
bug fixes within the feature (review round 2)

## Files changed
- `internal/dream/{collect,dream,run,timer}.go` and tests
- `internal/search/dream.go` and test
- `.agents/docs/specs/sh-agent/dream.md`
- the plan (D8b)
- the two round-2 review reports with Outcomes
- `.agents/docs/testing/2026-09-30-verify-parity-review-round-1.md`
- BUG-116..122 and INDEX

## Agents dispatched
reviewer (`reviews/2026-09-30-review-f1ca525-e2cffc1-code.md`), security (`reviews/2026-09-30-security-f1ca525-e2cffc1.md`), verifier (round-1 log PASS, 453 tests; two mutations), bugs (BUG-116..122).

## Verification
`go test -race -count=1 ./...` exit 0. Must-fail mutations for the two added tests, each compiled and failing, restored byte-identical. A real agy sandbox run: exit 0, 2 proposals.

## Requirements
R3, R4, R10 (ADR-047, ADR-047a).
