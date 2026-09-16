---
owner: reviewer
date: 2026-09-16
---
# Review — worktree-docs-retro (2026-09-16)

Range reviewed: working tree (uncommitted diff: `agents/changelog.md`, `agents/reviewer.md`, `context/AGENTS.md`, `skills/golang/SKILL.md`, `skills/init-project/SKILL.md`, `skills/review-standards/SKILL.md`)

## Verdict: REQUEST CHANGES

One restatement was left standing. `context/AGENTS.md` step 2 adds a new clause — when two dispatched agents contradict each other on one mechanism, fix the rule, not the code — but `skills/review/SKILL.md`'s own restatement of the same convergence rule (the operational logic `/review --fix` runs) says only "stop and ask what the mechanism is for" on a second round hitting the same mechanism or a fix introducing the next round's defect; it has no branch for the two-agents-contradict case at all. An operator following `skills/review/SKILL.md` during a `/review` round would not learn that a contradiction between two dispatched agents means the rule is the defect. Everything else queried for this pass — the checkpoint/EMPTY rule, the filter-that-fires addition, the changelog SHA rule, the golang reference-file sentence, the init-project headless remark, and the review-standards §3 / reviewer criterion-3 pairing — has no other restatement in the living docs that was left stale.

## Findings
- **[medium / confidence: high]** `skills/review/SKILL.md:62` — restates the review convergence/mechanism rule from `context/AGENTS.md` step 2 (same "stop and ask what the mechanism is for" language) but omits the new clause added there in this diff: "when two dispatched agents contradict each other on one mechanism, the rule they both read is the defect — fix the rule, then re-dispatch, never the code on one agent's word." `skills/review/SKILL.md` is what actually governs a `/review --fix` round's stop condition, so the gap is operative, not cosmetic.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens (see Findings: one restatement of a changed rule left unmatched)
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
Checked by targeted grep against each of the seven named subjects (checkpoint/EMPTY, step-2 contradiction, hard-rules filter sentence, changelog SHA rule, review-standards §3/reviewer §3 pairing, golang reference-file sentence, init-project headless remark) across `README.md`, `AGENTS.md` (root), `context/`, `agents/`, non-vendored `skills/*/SKILL.md`, and `.agents/docs/deploy/`. Did not open every skill's reference sub-files (e.g. `skills/*/references/*.md`) line by line beyond the targeted greps shown above — a restatement buried in a reference file rather than a top-level `SKILL.md` could have been missed. Did not check `skills/*/references/` content against the golang reference-file sentence beyond confirming the seven listed files exist and the table is internally consistent.

## Outcome

- Finding 1 (medium) — fixed: `skills/review/SKILL.md` Step 5's round loop now carries the contradiction clause ("When two lens agents contradict each other on one mechanism — one clears what the other flags — the rule they both read is the defect: fix the rule …, then re-dispatch; never the code on one agent's word"), matching `context/AGENTS.md` step 2.
- Not examined (reference sub-files) — closed by reading: `grep -rn 'mechanism' skills/*/references/` returns six lines, all in the golang references about goroutine shutdown mechanisms and Go's export mechanism; none restates the convergence rule.
- Round 1 closed 2026-09-16.
