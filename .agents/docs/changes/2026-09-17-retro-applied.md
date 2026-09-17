---
owner: changelog
date: 2026-09-17
---
# Retro proposals applied: six rules graduated     (2026-09-17)

## Summary
The operator approved all six proposals from `.agents/docs/retros/2026-09-17-retro.md`, drawn from the delivery rework session's three security rounds and one user restatement. Five rules went to `context/AGENTS.md` (BUG-057, BUG-058, BUG-061, the dispatch budget, and the REQUIREMENTS.md home for rules that must not drift), one to `agents/changelog.md` (BUG-061), and the dispatch text in `skills/review/SKILL.md` and `skills/preflight/SKILL.md` were updated to carry the dispatch budget rule.

## Type
task

## Files changed
- `context/AGENTS.md` — five sentences added, in § Hard rules (two), § Definition of Done steps 2 and 4, and § Work: check escapes survived file writes (BUG-057), every trigger in a finding must have a failing test (BUG-058), main agent reads agent-written records against diff before committing (BUG-061), dispatch to reviewer/security/verifier names a tool-call budget below maxTurns and repeats "write artifact early" (Proposal 3), a rule the user wants never to drift is proposed as a REQUIREMENTS.md edit first (Proposal 6)
- `agents/changelog.md` — one sentence added to Rules: every factual claim must trace to dispatch or diff, never invented (BUG-061)
- `skills/review/SKILL.md` — Step 3: every dispatch gets a tool-call budget below the agent's `maxTurns` and "write the artifact early, then update it"
- `skills/preflight/SKILL.md` — Step 2: the same addition to the reviewer, security and verifier dispatches

## Agents dispatched
verifier, reviewer (docs lens), bugs (BUG-057, BUG-058, BUG-061 graduated to rules), changelog

## Verification
- Verifier log `.agents/docs/testing/2026-09-17-retro-context.md` PASS: gofmt, vet, `go test -race -count=1 ./...` all clean (124 pass); scratch install drift none, 259 files match; phrases from the BUG-058, dispatch-budget, BUG-061 and REQUIREMENTS.md sentences confirmed in both vendors' installed context copies; "Every factual claim traces to the dispatch or the diff" confirmed deployed to changelog agent. One check void: escape needle reached verifier as plain text; main agent re-ran verification with script building the needle from character codes — exactly one real occurrence in context/AGENTS.md and both vendor copies.
- Reviewer log `.agents/docs/reviews/2026-09-17-review-24c524e-retro-docs.md` APPROVE WITH NITS: all six sentences traced to their bug records; one gap found (tool-call budget rule not yet in skill templates) and fixed in same commit, re-verified with scratch install (drift none, new rule deployed).
- After fix: main agent re-ran `go run ./cmd/sh-agent install --home .agents/scratch/retro-home` (drift none, skill text confirmed) and `go test -count=1 ./...` (pass).

## Docs
- `context/AGENTS.md` — § Hard rules, § Definition of Done and § Work
- `agents/changelog.md` — § Rules section appended
- `skills/review/SKILL.md` — Step 3 template updated
- `skills/preflight/SKILL.md` — Step 2 template updated
- Trail artifacts: `.agents/docs/retros/2026-09-17-retro.md` (proposals), `.agents/docs/reviews/2026-09-17-review-24c524e-retro-docs.md` (review), `.agents/docs/testing/2026-09-17-retro-context.md` (verification)
- No README, deploy guides, or user-facing docs changed

## Requirements
n/a

## Not exercised
- Headless CLI starts in sandbox home (no logged-in sandbox on this box, per plan row)
- Real-home smoke checks (not in dispatch scope; only scratch home used per project dev rule)
- `/review --docs` of the two skill edits: the review covered the context and changelog diff, and its finding produced the skill edits, which were checked by scratch install only
