---
owner: reviewer
date: 2026-10-05
---
# Review — worktree-docs-2   (2026-10-05)

Range reviewed: working tree vs HEAD (fe77469): `.agents/docs/deploy/local.md`, `.agents/docs/testing/guide.md`, `AGENTS.md`, `README.md`, `rules/critical-thinking.md`, `skills/feedback-log/SKILL.md`, `skills/retro/SKILL.md` (docs lens; the untracked `testing/2026-10-05-verify-critical-thinking-r2.md` was read only for its range line). Does not overwrite `reviews/2026-10-05-review-worktree-docs.md`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The four copies of the full-suite command (`AGENTS.md:23`, `README.md:286`, `deploy/local.md:18`, `testing/guide.md:17`) agree, and the numbers behind them reproduce from a record (`~790 s` for `cmd/sh-agent` under `-race` on `nproc=1` in `testing/2026-10-05-verify-critical-thinking-go.md:17`; `~62 s` in `testing/2026-10-01-verify-google-skills-catalog-phase-4.md:42`). What is broken: the guide now says `cmd/sh-agent` outlasts Go's 10-minute default on this box, but two other commands in the same guide still run that package under `-race` with no `-timeout`, so as written they fail here. The `rule:<name>` widening in `feedback-log` and `retro` is not carried to the places that decide where a harness rule finding goes.

## Findings
- **[medium / confidence: high]** `.agents/docs/testing/guide.md:123` and `:135` (§ 3d, § 3e) vs `guide.md:17-18` — `go test -race ./internal/search ./cmd/sh-agent` and `go test -race ./internal/memory ./internal/dream ./internal/search ./cmd/sh-agent` carry no `-timeout`; § 1 now states `cmd/sh-agent` takes ~13 min under `-race` on a one-CPU box, past the 10-minute default, so these two commands panic with "test timed out" on the box the guide is written for. Fix: add `-timeout 30m` to both lines (or one sentence in § 3d/3e pointing at § 1's flag).
- **[low / confidence: medium]** `skills/harness-audit/references/areas.md:36` (and `skills/harness-audit/SKILL.md:36`, which tells the auditor to run `AGENTS.md` § Setup commands, so that one follows the change) — the Baseline row lists a literal `go test -race ./...` with no timeout; run as written on this box it times out in `cmd/sh-agent`. Fix: `-timeout 30m` in the row, or "the Setup-commands line, verbatim".
- **[low / confidence: medium]** `skills/retro/SKILL.md:60` vs `:55` — the category table gives "Rule gap or restatement" the harness target `context/AGENTS.md` only. Step 2 item 1 now says `rule:<name>` also names the harness's global `rules/<name>.md`, and `rules/critical-thinking.md` exists, so a finding about that rule has `rule_ref` `rule:critical-thinking` but no harness target row that can hold the fix. Fix: harness cell becomes `context/AGENTS.md`, or `rules/<name>.md` for an always-on/path-scoped rule.
- **[low / confidence: medium]** `agents/dreamer.md:32` and `skills/dream/SKILL.md:30` and `:39` — "what already exists" / "graduated" list the global context, project `AGENTS.md` and `.agents/rules/`, not the installed global rule files (`~/.claude/rules/`, `~/.gemini/config/rules/`; `internal/harness/plan.go:62`). A candidate that `rules/critical-thinking.md` already states is therefore not recognised as existing or graduated. Not caused by this diff alone (the rule file landed in da620c3), but this diff's `rule:<name>` change is the first text to treat global rule files as a rule home. Fix: add the installed global rule files to both lists.

Checked and clean: `rules/critical-thinking.md:8` against `context/AGENTS.md:5` (both say a direct instruction is its own approval; the rule adds where the risks go); `README.md:240` now names `trigger: always_on` (the earlier review's gap is closed); `README.md:278` and `AGENTS.md:43` rules/ lines still match one shipped rule; no `*_test.go` asserts the old command text (grep for `test -race` in `internal/` and `cmd/` found none).

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- The Status and body of the untracked `testing/2026-10-05-verify-critical-thinking-r2.md` (only its range line, grep hits); whether round 2 resolved Antigravity P2 is unverified.
- Whether the verifier agent (`agents/verifier.md:55`, `go test -race -v ./...`) should learn the timeout: it is generic across projects, so left out.
- Antigravity's behaviour on the new "direct instruction is carried out" wording (a behavioural run, not a docs check).
- Records excluded by brief: `changes/`, `DECISIONS.md`, `handoffs/`, `reviews/`, `retros/`, bug records, dated logs, plans.

## Outcome

- [medium] guide § 3d/§ 3e commands lack the timeout — fixed by `documentation`: `.agents/docs/testing/guide.md:123` and `:135` carry `-timeout 30m`.
- [low] `skills/harness-audit/references/areas.md:36` bare `go test -race ./...` — fixed at `skills/harness-audit/references/areas.md:36`.
- [low] `skills/retro/SKILL.md:60` harness rule-gap target lacks `rules/<name>.md` — fixed at `skills/retro/SKILL.md:60`.
- [low] `agents/dreamer.md:32`, `skills/dream/SKILL.md:30,39` list only `.agents/rules/` as where a rule already exists — left: operator. Changing an agent definition needs its own behavioural run (project `AGENTS.md` § Testing), outside the approved items; put to the operator as one line with a recommendation.
