---
owner: changelog
date: 2026-09-15
---
# Harness practices Phase 1: docs lens and lost-rule tracking (G06, G08)     (2026-09-15)

## Summary

Completed Phase 1 of the harness-practices plan (`.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`). Applied edits I-1a, I-2a, I-3a, I-4a and I-5a–e, recorded BUG-018 as fixed (I-6), and ran the one-time docs sweep (task 1.6). The bug skill's missing-repro gate holds whatever autonomous instruction is in effect; a retro diff that deletes or rewrites a rule lists each removed rule's new home or `dropped: <why>`; the review skill's `--docs` lens hunts dropped and narrowed rules; this repo's `AGENTS.md` Testing section adds `/review --docs` with its extra living docs; and the hardcoded skill and subagent counts are gone from `AGENTS.md`, `README.md` and the testing guide. The sweep's seven reports over six areas hold 98 findings; 9 became rehearsed plan edits (under D39, D40), 13 became open bugs (BUG-019 to BUG-023), 76 filed as Follow-ups.

## Type

task

## Files changed

- `skills/review/SKILL.md` — I-1a: `--docs` lens hunts old wording at the scope's head, greps living docs for dropped and narrowed rules, reports findings with `--project` classification
- `skills/bug/SKILL.md` — I-2a: restored "whatever autonomous instruction is in effect" clause to missing-repro gate
- `skills/retro/SKILL.md` — I-3a: diff that deletes or rewrites a rule lists each removed rule's new home or `dropped: <why>`
- `AGENTS.md` — I-4a: the Testing section adds `/review --docs` over the range for changes to `context/`, `agents/` or `skills/`, with this repo's extra living docs; I-5d: the skill and subagent counts are removed from :7
- `README.md` — I-5a to I-5c: "27 skills" and "seven subagents" removed from :3, :86 and :87
- `.agents/docs/testing/sh-agent/guide.md` — I-5e: the `sh-agent list` comment no longer states counts; it says the section headers print them
- `.agents/docs/bugs/BUG-018-bug-skill-missing-autonomy-clause.md` — status: fixed
- `.agents/docs/bugs/INDEX.md` — BUG-018 row updated to status: fixed
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — boxes 1.1–1.6 and Phase 1 checkpoint ticked; sweep edits I-7c–e, I-12f, I-15e–f, I-20e, I-22f, I-25a–c and decisions D39, D40 added; D15 records the re-run rehearsal (101 edits, 0 failures, 29 files, 111 expectations hold)
- `.agents/docs/reviews/2026-09-15-review-sweep-merged.md` — merged findings from seven area reviews (task 1.6); its Outcome names one owner per finding: 9 as rehearsed plan edits, 13 as BUG-019 to BUG-023, 76 under Follow-ups
- `.agents/docs/reviews/2026-09-15-review-sweep-{context-agents,repo-docs,repo-docs-2,skills-1,skills-2,skills-3,skills-4}.md` — the seven area reports of the sweep
- `.agents/docs/bugs/BUG-019-*.md` to `BUG-023-*.md` and their `INDEX.md` rows — open bugs from the sweep
- `.agents/docs/testing/2026-09-15-verify-harness-practices-phase1.md` — the Phase 1 verifier log

## Agents dispatched

documentation (applied I-5a–c and I-5e; declined I-4a and I-5d as outside its lane, so the main agent applied them per D7), reviewer (task 1.2 control, seven 1.6 area reviews), bugs (BUG-019–023), verifier

Review: no `reviewer` ran over the Phase 1 source diff. The plan's Done list runs `reviewer` and a `--docs` review once over `c192622..HEAD`.

## Verification

Verifier log `.agents/docs/testing/2026-09-15-verify-harness-practices-phase1.md` reports Status PASS. Main agent verified: `go test -race ./...` exit 0; `gofmt -l .` prints nothing; `go vet ./...` exit 0; scratch install to `.agents/scratch/phase1-checkpoint-home` exit 0, then status with drift none; `render` output 93 lines (expected). The verifier re-ran the verify commands of boxes 1.1, 1.3, 1.4 and 1.5, which met their expects; boxes 1.2 and 1.6 carry their own evidence in the plan. Sandbox CLI start deferred to Phase 5 checkpoint (plan D36).

## Docs

- `skills/review/SKILL.md` — behavior change: `--docs` lens now hunts dropped/narrowed rules
- `skills/bug/SKILL.md` — behavior restored: autonomy clause on missing-repro gate
- `skills/retro/SKILL.md` — behavior change: rule-tracking on diff deletion/rewrite
- `AGENTS.md` (this repo) — requirement added: `/review --docs` for changes to context/, agents/, skills/
- `README.md` — hardcoded counts removed
- `.agents/docs/testing/sh-agent/guide.md` — the `list` sanity comment names no counts

## Requirements

n/a
