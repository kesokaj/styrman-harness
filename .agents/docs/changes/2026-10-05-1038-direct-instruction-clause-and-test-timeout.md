---
owner: main
date: 2026-10-05
---
# Direct instruction carried out under the critical-thinking rule; test timeout for one-CPU boxes

## Summary

Closed the Antigravity row P2 that `2026-10-05-verify-critical-thinking.md` left FAIL (ADR-050's open consequence): the rule's stress-test bullet now ends "and a direct instruction is carried out: its risks go in the reply's table, after the work". Approved by the operator in their own turn: "go approved 1-4, skip 5" (1 the clause, 2 the skip flag for test runs in throwaway clones, 3 the timeout, 4 the `rule:<name>` vocabulary).

## Type

task

## Files changed

- `rules/critical-thinking.md` — the clause (Antigravity block 16,907 of 22,000 bytes; context unchanged at 14,700).
- `AGENTS.md` setup line, `README.md:286`, `.agents/docs/testing/guide.md` § 1, § 3d, § 3e, `.agents/docs/deploy/local.md:18`, `skills/harness-audit/references/areas.md:36` — `go test -race -timeout 30m`: on this one-CPU box `cmd/sh-agent` takes 631-790 s under `-race`, past Go's 10-minute default.
- `skills/feedback-log/SKILL.md:58`, `skills/retro/SKILL.md:55,60` — `rule:<name>` also names the harness's global `rules/<name>.md`; a harness rule gap may target it.

## Agents dispatched

documentation (twice: the setup and guide § 1 lines with README and local.md; guide § 3d/§ 3e), reviewer (docs lens, `reviews/2026-10-05-review-worktree-docs-2.md`, APPROVE WITH NITS; Outcome: 3 fixed, 1 left: operator), verifier (PASS). Changelog written by main from the diff.

## Verification

- `testing/2026-10-05-verify-critical-thinking-r2.md` (main, PASS): Antigravity P2 with the clause 3 of 3 carried out; the old rule with the same flag 1 of 2 (old rule across both rounds: refused 2 of 3); Antigravity P4 1 of 1; Claude Code P2 2 of 2, P1 1 of 1.
- `testing/2026-10-05-verify-critical-thinking-r2-go.md` (verifier, PASS): `gofmt`, `go vet`, `go test -race -timeout 30m ./...` 553 pass, 0 fail, 3 env-gated skips; `list`, `render` 87 lines, size pins, scratch install drift none, rule byte-identical in both targets.

## Notes

- Left for the operator: `agents/dreamer.md:32` and `skills/dream/SKILL.md:30,39` list only `.agents/rules/` as where a rule already exists, not the global rule files; changing the agent needs its own behavioural run.
- Item 5 (a probe that shows the anti-sycophancy rule's effect) skipped by the operator.
