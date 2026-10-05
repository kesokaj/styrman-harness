---
owner: changelog
date: 2026-10-05
---
# Critical thinking rule and operator approval protocol

## Summary

Deployed the harness's first always-on rule file (`rules/critical-thinking.md`), an anti-sycophancy protocol (stress-test before agreeing, critique first, no flattery, no apologies). Added the operator's reply-and-approval protocol to the context's Operator paragraph: a short explanation and recommendation, a numbered table highest priority first, and no work until the operator approves by number (`go all approved`, `go approved 1-3, skip 5`). Aligned three skills (feature, dream, requirements) with the new approval protocol and documented the first always-on rule key in README.

## Type

task

## Files changed

From `git diff HEAD --stat` ad2130b:
- `rules/critical-thinking.md` — new; anti-sycophancy protocol (trigger: always_on, 14 lines)
- `context/AGENTS.md` — Operator paragraph: added approval protocol; struck SHORT LIST
- `.agents/docs/architecture/DECISIONS.md` — ADR-050 (11 lines added)
- `AGENTS.md` — project layout line: rules/ no longer empty
- `README.md` — rules row (line 240): documents `trigger: always_on`; layout line 278: `rules/` holds one file
- `skills/feature/SKILL.md`, `skills/dream/SKILL.md`, `skills/requirements/SKILL.md` — updated to cite context's approval protocol; requirements/SKILL.md dropped dead quote
- `.agents/docs/reviews/2026-10-05-review-worktree-docs.md` — reviewer findings (docs lens), `## Outcome` by main: 4 fixed (README:240 by documentation, the rest by main), 1 rebutted, 1 later nit fixed by main
- `.agents/docs/testing/2026-10-05-verify-critical-thinking-go.md` — verifier test log (cmd/sh-agent install, `-timeout 30m` for single-CPU box)
- `.agents/docs/testing/2026-10-05-verify-critical-thinking.md` — behavioural test log (Claude Code PASS, Antigravity FAIL: P2 direct-instruction refusal under anti-sycophancy rule)

## Agents dispatched

reviewer (docs lens, one round, APPROVE WITH NITS), documentation (2 runs: README:278, README:240), verifier (PASS), changelog; this entry was corrected by main against the diff

## Verification

- `verifier`: `testing/2026-10-05-verify-critical-thinking-go.md` — PASS; `gofmt`, `go vet`, `go test -race -count=1 -timeout 30m ./...` (553 pass, 0 fail, 3 env-gated skips; `-timeout 30m` because this box has one CPU and `cmd/sh-agent` takes 635-790 s under `-race`, past the 10-minute default), `list`, `render` (87 lines), scratch install `.agents/scratch/verify-ct-home` drift none, rule file byte-identical in both targets, sizes reproduced (hand-written 14,700 bytes, Antigravity block 16,815 of 22,000)
- `documentation`: `README.md` 278 and 240; main: `skills/{feature,dream,requirements}/SKILL.md` restatements
- Behavioural test `testing/2026-10-05-verify-critical-thinking.md` — Status FAIL (Claude Code: approval protocol visible, anti-sycophancy no measurable difference vs control; Antigravity: P2 direct instruction refused under rule, control ran it 3 of 3); the fix ("A direct instruction is its own approval") is the operator's choice (logged in handoff Next Steps)

## Docs

- `context/AGENTS.md` — Operator paragraph: approval protocol
- `AGENTS.md` — Layout line updated (rules/ not empty)
- `README.md` — rules row (`trigger: always_on`), layout line
- `rules/critical-thinking.md` — new, living doc
- Three skills (`feature`, `dream`, `requirements`) updated to name the approval protocol
- No interface, command, or deploy step changed (rules are always-on, no new CLI verb)

## Requirements

n/a

## Notes

Open issue on Antigravity (from the behaviour log, one valid sample): a direct instruction was refused under the new rule, which the control home carried out 3 of 3, so the context's "A direct instruction is its own approval" does not yet hold there. The fix is the operator's choice, in the handoff's Next Steps. A `--dangerously-skip-permissions` re-run was refused by the session's permission classifier and not attempted another way.
