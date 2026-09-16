---
owner: changelog
date: 2026-09-16
---
# BUG-041: security agent graded stdlib crypto/rand UUID generation as hand-rolled cryptography     (2026-09-16)

## Summary

Sandbox `/feature` run on Antigravity had the harness `security` agent flag a UUID v4 built from `crypto/rand` bytes as "hand-rolled cryptography, strictly prohibited" while the harness `reviewer` judged the same lines as stdlib use exempt from the DIY ban. Both agents' reports closed with `Not examined: nothing`, suppressing visibility of their disagreement. Root cause: all four restatements of the "never hand-roll crypto" rule (context § Hard rules, review-standards § 5, security § 4, reviewer criterion 5) never defined what "hand-rolled" means, so consuming `crypto/rand` bytes read to one agent as rolling crypto rather than using stdlib as designed. Fixed: all four restatements now define hand-rolled as a primitive or protocol implemented or composed in the target code and named in the finding; stdlib `crypto/*` used as designed (bytes or identifier from `crypto/rand`, a checksum, TLS) is explicitly no finding. When a dependency is requested, the finding must name the stdlib alternative rejected (e.g. `crypto/rand.Text()` in Go 1.24). The three `Not examined` lines now require naming the gap if full coverage is not claimed.

## Type

bug fix

## Files changed

- `context/AGENTS.md` — § Hard rules: "Never hand-roll auth, crypto, hashing, sessions, JWT/OIDC" now defines it as "a primitive or protocol implemented or composed here" and names the stdlib use that is exempt.
- `agents/security.md` — § 4 Crypto and § 5 Dependencies updated to define hand-rolled and name the stdlib alternative for dependency findings; § Not examined unconditionally requires naming the gap if full coverage is not claimed.
- `agents/reviewer.md` — criterion 5 updated to state the stdlib-safe forms explicitly, with "as designed" applying to all items in the list.
- `agents/verifier.md` — § Not examined unconditionally requires naming gaps.
- `skills/review-standards/SKILL.md` — § 5 exemption updated to define hand-rolled precisely and list stdlib crypto uses that are not DIY.

## Agents dispatched

reviewer (`reviews/2026-09-16-review-bug-g6-security-exemption.md` APPROVE WITH NITS; two findings fixed: unconditional `Not examined` rule in security.md, dangling modifier in reviewer.md criterion 5), verifier (shared log: rendered context 93 lines, list 27 skills 7 agents, no drift), bugs, changelog.

## Verification

Verifier log `.agents/docs/testing/2026-09-16-verify-linkshelf-sim-fixes.md`: no rendered-text changes in context output (93 lines, consistent with prior). Grep of the four restatements (context, security, reviewer, review-standards) confirms all four now define hand-rolled and name stdlib exemption consistently.

Owed but not done in this commit: sandbox row `security-rand-id` after-the-fix (a2) on both CLIs (the before-the-fix control at a1 showed Antigravity's security agent filed crypto/rand UUID as hand-rolled; after the fix it should not), recorded later in testing logs.

## Docs

`context/AGENTS.md` (§ Hard rules), `agents/security.md`, `agents/reviewer.md`, `agents/verifier.md`, `skills/review-standards/SKILL.md` all updated to restate the exemption consistently; no stale old wording left in any living doc (checked by git grep for full-service scope).

## Requirements

n/a
