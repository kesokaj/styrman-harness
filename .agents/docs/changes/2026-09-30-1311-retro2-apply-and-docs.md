---
owner: main
date: 2026-09-30
---
# Retro 2 applied: Next Steps never carry a standing procedure or a second-time item; harness proposals apply in the harness; docs brought current     (2026-09-30)

## Summary
The operator: "Go all. Then cut a new version again and update docs." The three proposals of `retros/2026-09-30-retro-2.md` applied here (the harness checkout): `skills/handoff/SKILL.md` § Next Steps never lists a step the deploy guide records as a standing procedure and never carries an item from the previous handoff a second time — it is done or put to the operator as one line; its Step 3 and the scaffold template `skills/init-project/assets/TODO.md` (and this project's `TODO.md` header) carry the same qualifier, found by the docs lens; `skills/retro/SKILL.md` § Step 6 applies an approved harness-scope proposal in the same run when the retro runs in the harness checkout. Docs pass over the day (`documentation`): `deploy/local.md:66` tightened into one paragraph, no claim added or dropped; README, the specs, the testing guide, the deploy guides and both matrices checked and left as they were, each with the reason in the agent's report.

## Type
task (a retro apply and a docs pass)

## Files changed
- skills/handoff/SKILL.md — two Next Steps paragraphs; Step 3 qualified
- skills/retro/SKILL.md — the Harness-scope bullet
- skills/init-project/assets/TODO.md, TODO.md — the header sentence qualified
- .agents/docs/deploy/local.md — the redaction paragraph, one edit by `documentation`

## Agents dispatched
reviewer docs lens (`reviews/2026-09-30-review-worktree-retro2-docs.md`, REQUEST CHANGES: two restatements, both fixed; a third left, recorded in the retro's Not proposed; Outcome), documentation (the day's docs pass; one file changed), verifier (`testing/2026-09-30-verify-retro2-apply.md`). No code reviewer: 17 lines of prose over five files, read by the main agent.

## Verification
`testing/2026-09-30-verify-retro2-apply.md`. Owed at the release gate: the real-home smoke on both CLIs and one Antigravity run of `handoff`, `retro` and `init-project` (its scaffold template changed).

## Docs
`deploy/local.md:66`; the rest confirmed current.

## Requirements
n/a
