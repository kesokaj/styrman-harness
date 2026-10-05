---
owner: main
date: 2026-10-05
---
# Preflight — v2.4.0   (2026-10-05)

Range reviewed: v2.3.2..b549d71

Version: MINOR. The range adds behaviour and removes nothing a place used: the first always-on rule file (`rules/critical-thinking.md`), the operator's reply-and-approval protocol in the context (ADR-050), the dreamer and `dream` reading the global rule files, `rule:<name>` naming them in `feedback-log`/`retro`, and `-timeout 30m` on the test commands. The operator first said "ship v2.3.3" (a number the agent had proposed), then, told the version rule, "cut v2.4.0 … go approve all" in their own turn of 2026-10-05.

## Findings, worst first

No blocking finding.

| # | Finding | Severity / confidence | Source | Decision |
|---|---|---|---|---|
| 1 | `rules/critical-thinking.md:14` "Disagree before complying" unconditional, contradicting line 8 (the wording behind the earlier Antigravity refusals) | medium / medium | code review | fixed before the tag (`b549d71`); P2 re-run at b549d71: Antigravity 3 of 3, Claude Code 2 of 2 |
| 2 | Nightly dream job could load the new reply format and fail its one-JSON-object output (`internal/dream/dream.go:177`) | medium / low | code review | rebutted: the job's exact Claude command (`--setting-sources local …`) answered `NONE` to a canary probe while the default-sources control named GANNET and the rule file; Antigravity runs from a private HOME holding only the login and a grant-nothing settings file (`internal/dream/run.go:249-283`), so no rules there |
| 3 | `testing/guide.md:198` expected one `<RULE[...]>` block | medium / high | docs review | fixed before the tag (`b549d71`, by `documentation`) |
| 4 | `REQUIREMENTS.md:55` says a second Antigravity rules file "is dropped without a marker"; today's stores show both blocks loaded | medium / low | docs review | left: operator — `REQUIREMENTS.md` changes only with their approval; the proposed edit is in the handoff |
| 5 | `skills/harness-audit/SKILL.md:36` bare `go test -race ./...` | low / high | docs review | fixed before the tag (`b549d71`); Antigravity run 1 of 1 at b549d71 |
| 6 | `skills/feature/SKILL.md:21` still says "numbered list" for a backlog | low | code review | accepted: a backlog restatement is a scope list the skill owns, not a reply of choices; named in the release entry |
| 7 | `skills/feature/SKILL.md:15` "Everything else runs without asking" vs the context's approval gate | low | code review | accepted: it runs inside a scope the operator already approved at the skill's gate, which the context counts as approval |
| 8 | `skills/retro/SKILL.md:140` "ask which to apply" keeps the older approval wording | low | docs review | accepted: the meaning matches; named in the release entry |

## Rows of the Definition of Done

| Row | Evidence |
|---|---|
| Reviewer | `reviews/2026-10-05-review-v2.3.2-0c0846c-code.md` (opus, APPROVE WITH NITS) over v2.3.2..0c0846c; the 2-line source diff of `b549d71` read by main (below threshold) |
| Security | not owed: no path in the range touches auth, crypto, secrets, IAM, input parsing or dependencies; no Go source changed |
| Verification 1 | `testing/2026-10-05-verify-v2.3.3-gate.md` (owner verifier, PASS): range v2.3.2..worktree = b549d71's tree; after it only records changed |
| Verification 2 | edited commands (`-timeout 30m` in `AGENTS.md`, guide § 1/3d/3e) run as quoted by that log; main-agent tiers in `testing/2026-10-05-verify-v2.3.3-gate-rows.md` (PASS, every row n of N) |
| Verification 3 | `testing/2026-10-05-verify-critical-thinking.md` (FAIL, Antigravity P2 0 of 1) cleared by `-gate-rows.md` (P2 3 of 3 at b549d71; fixes `07564fb`, `b549d71`) and earlier by `-critical-thinking-r2.md` |
| Docs (Step 4) | `reviews/2026-10-05-review-v2.3.2-0c0846c-docs.md` (APPROVE WITH NITS); findings above |
| Pre-gate commit | `b549d71` (Step 1: fixes 1, 3, 5) |

## Trail freshness (Step 3)

- Newest `changes/` entries describe this range: pass.
- No future `date:` (only template placeholders): pass.
- No verification log deleted: pass.
- `DECISIONS.md` modified once, by `da620c3`, 11 lines added and 0 removed (ADR-050 appended): pass.
- `ROADMAP.md`, `TODO.md` untouched since v2.3.2 (both empty): report only.

## The project's own release rule

Real-home smoke on both CLIs: context at b549d71 (gate-rows log), skill and agent rows at 0be3b97 (`testing/2026-10-05-verify-real-home-critical-thinking.md`); one Antigravity run of each changed skill and agent (`dreamer`, `dream`, `feature`, `requirements`, `feedback-log`, `retro`, `harness-audit`): gate-rows log, each 1 of 1.

## What was not checked

- The anti-sycophancy rule's effect against a control that flatters or apologises (operator: skip): not a blocking condition.
- Real-home skill and agent rows at b549d71 itself (taken at 0be3b97; b549d71 changed no skill or agent those rows exercised).
- `govulncheck` (no Go or dependency change in the range).
- The other three places.

None of these is named by a blocking condition.

VERDICT: SHIP
