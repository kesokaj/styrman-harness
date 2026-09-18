---
owner: main
date: 2026-09-18
---
# Retro 2 (2026-09-18) applied: five proposals, the delivery claim consolidated     (2026-09-18)

## Summary

The five proposals of `retros/2026-09-18-retro-2.md` are applied on the operator's "apply all this time" — three harness-scope (queue `~/.config/sh-agent/proposals/2026-09-18-retro.md`, second section) and two project-scope: showing the release notes is not a gate (`preflight`); a vocabulary's every enumeration and a reversed claim's every restatement change in the same edit (`context/AGENTS.md`); `reviewer`, `security` and `explorer` write an IN PROGRESS artifact after a third of their budget; the project `AGENTS.md` § Testing says how a behavioural run is made and what it proves (BUG-083 graduates); and the delivery claim — what a tag delivers, to which kind of place — now has one home, `specs/sh-agent/delivery.md` § The model, with README, the project `AGENTS.md` (rules 2 and 6, § Release), `deploy/remote.md` (four sentences) and `delivery.md`'s own two older sentences turned into pointers (BUG-084 graduates). The operator's instruction: no further retro after this; the next tag is the last of the day.

## Type

context, skills, agents, project rules and docs

## Files

- `skills/preflight/SKILL.md` Step 6.1 — showing the notes is not a gate; continue to the tag, push and upload.
- `context/AGENTS.md` § Hard rules — "Two shapes hide best": a VOCABULARY's every enumeration; a reversed CLAIM's every restatement, one grep before the first review (BUG-084).
- `agents/reviewer.md`, `agents/security.md`, `agents/explorer.md` § Write the artifact EARLY — an IN PROGRESS file after a third of the budget.
- `AGENTS.md` § Testing — the behavioural-run recipe (dispatch does not repeat the rule; prompt on stdin; `--allowedTools`; transcript read; n of N) and "a scratch `HOME` cannot authenticate" (BUG-083); § How the harness works rule 6 — the delivery claim has one home; rule 2 and § Release point at it.
- `.agents/docs/specs/sh-agent/delivery.md` — new § The model; `:7` and § What is NOT being built point at it.
- `README.md:245`, `.agents/docs/deploy/remote.md` (`:8`, `:91`, `:95`, `:128`, `:195`) — pointers.
- `.agents/docs/bugs/INDEX.md`, `BUG-083`, `BUG-084` — graduated.

## Verification

- `gofmt`, `go test ./...` ok; `go run ./cmd/sh-agent install --home .agents/scratch/verify-v140/home` from empty → drift none; `list` 27 skills, 7 agents. `verifier` dispatched over the working tree — `testing/2026-09-18-verify-retro-2-applied.md`.
- Not exercised: a CLI start (the real home is the operator's; the sandbox procedure was not run today); a behavioural run of the three agents' new sentence (BUG-083 says the definition alone binds in some runs only; the dispatch sentence carries it).

## Definition of Done

- Review: `reviewer` (opus, code lens) and `reviewer` (docs lens) over the working tree — `reviews/2026-09-18-review-worktree-retro-2-{code,docs}.md`. Round 1: code REQUEST CHANGES (the three agents' templates had no Status field — the BUG-084 shape, made in the diff that wrote its rule; README:59/:61, rule 1 and `delivery.md:36-37` still restated the claim; three lows), docs REQUEST CHANGES (README:59/:61 — the sweep grep matched only the singular). Round 2: docs APPROVE; code APPROVE WITH NITS, its open items all on rule 6's scope — decided by the main agent, no round 3: rule 6 covers the claim, a command's description in a guide is not a restatement; `delivery.md:37` trimmed; § The model completed. § Outcome in each record.
- Verifier: PASS twice (`testing/2026-09-18-verify-retro-2-applied.md`, § Round 2 on the frozen tree after the review fixes); the tree moved under round 1 only in `.agents/docs/` (records), outside its scope.
- Agents dispatched: `verifier`, `reviewer` ×2, `bugs`.
- Statuses set to `applied 2026-09-18 (<commit>)` in the queue file and the record after the commit.
