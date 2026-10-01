---
owner: reviewer
date: 2026-10-01
---
# Review — 329a99b-worktree-prompts   (2026-10-01)

Range reviewed: working tree against 329a99b (`git diff 329a99b -- context agents skills README.md ROADMAP.md`), prompts lens (`/review --prompts`)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Two rules left the always-loaded context and landed nowhere (HIGH). Four mechanical contradictions remain between the new context and the files it points at (ADR pointer, ROADMAP "Next", the resume order restated in `handoff`, README map vs AGENTS.md rule 1). Every cited section (`review` Step 4, `review-standards` § 2 and § 8, `bash` § Installed commands, `bug` Step 0, `retro` § When to Suggest) exists, and most other moves landed where the plan says.

## Findings
- **[high / confidence: high]** `context/AGENTS.md:59` — lost rule. Old line 58 said: "before trusting a monitor, watcher or filter that reports nothing, confirm it can emit at all. A green run is evidence only when it demonstrably hit the thing under test; a filter that fires is read before it counts, since a needle matches its own negation". The new line keeps only the must-pass/must-fail clause. The plan (§ Interfaces, "line 58 → first clause + pointer; the rest is already review-standards § 2") assumes the rest is in § 2. It is not: grep for `needle|monitor|watcher|emit|negation` over review-standards, review, bug and feature returns nothing; "demonstrably" was not separately grepped. Fix: add those three clauses to review-standards § 2 (next to the mutation bullet) or restore one sentence in the context.
- **[high / confidence: high]** `context/AGENTS.md:18` — lost rule. Old Records step: "you read any agent-written entry or ADR against the diff before committing it". It is gone, in no skill or agent (grep "against the diff|agent-written" over context, skills and agents is empty), and the plan's Context budget list does not name it as deleted. Fix: restore it in step 4, or move it to `skills/feature` Step 4 / `agents/changelog.md` and cite it.
- **[medium / confidence: high]** `context/AGENTS.md:18` — contradicts and dangles. The context says "you write the ADR (shape: the `feature` skill § Decisions)". But `skills/feature/SKILL.md:67` `## Decisions` is the plan's D-row format, not an ADR shape, and the ADR template is still in `agents/changelog.md:72` (description: "writes ADRs ... after any design decision"). Two authors, one dangling pointer. The plan moves the template to feature § Decisions and strikes ADR authoring from changelog, but neither landed in this tree. Either land both or revert the context sentence.
- **[medium / confidence: high]** `README.md:57,65` vs `AGENTS.md` rule 1 (and rule 4) — README and ADR-049 now draw a dashed olive line "redacted prompt" from the dreamer to the model provider, i.e. the dreamer sends something off the machine. Rule 1 (AGENTS.md, unchanged) still says the timer "sends nothing off it". Rule 4 and the legend sentence still say "A dashed line is a gap: done by hand, or limited to one machine", but ADR-049 makes this dashed line mean "leaves the machine", not a gap. Rule 3 requires the map, docs and ADR in one commit; AGENTS.md rule 1 is part of that set. Fix: edit rule 1's clause and the legend's dashed definition in the same commit. ADR-049 records no operator approval (rule 3 says the operator first); I could not tell from the diff whether it was given (low confidence).
- **[medium / confidence: high]** `context/AGENTS.md:18` vs `skills/init-project/assets/ROADMAP.md`, `ROADMAP.md:3`, `agents/documentation.md:3,35,50` — the context says "ROADMAP.md holds only 'Now' while a plan is open"; the template and this repo's ROADMAP now carry "Next: at most three scheduled milestones", and the documentation agent still maintains "Next / Later" (the template dropped the Later section). Pick one: change the context to "ROADMAP.md holds Now, and Next for scheduled milestones", and drop Later from `documentation.md`.
- **[medium / confidence: medium]** `skills/handoff/SKILL.md:168` restates the resume chain in the old order ("AGENTS.md → TODO.md → this handoff → ROADMAP Now → plan → changes → bugs INDEX → search --open") against the new context order (handoff first, TODO.md after bugs INDEX, `memory show` added, "recon-*/research-* are briefs"). The cold reader would run a different chain from the one the context mandates.
- **[low / confidence: high]** Dated or project-specific text in global files. `context/AGENTS.md:64`: "R2 of the harness lists them" names this repo's requirement id; no other project has an R2. `skills/handoff/SKILL.md` (new paragraph): BUG-073, BUG-136. `skills/review-standards/SKILL.md` § 7/§ 8: BUG-011, BUG-086, BUG-122 and "plan audit-fix, C1". The plan allows only the "(moved from the global context 2026-10-01)" provenance note. Fix: drop the BUG ids and plan ids; point the context at "the installer's list" without R2.
- **[low / confidence: medium]** Smaller losses and stale plan text. The plan said the stop sentence keeps `TaskStop`; the context dropped it (`context/AGENTS.md:35`). Also gone and unnamed by the plan: "the `changes/` entry is its record" for a ticked TODO line, "A handoff is not only for a session that ends mid-task", and the old "if the trail cannot say where work stopped, ask". The plan (`2026-09-30-audit-fix.md:32,39,105`) still says the code-shape rules go to review-standards § 4 and shapes to § 3; they landed in § 8 and § 7. Context and skills agree with each other, so fix the plan text.

## Criteria
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

## Not examined
- `agents/*.md` other than `bugs.md` (not in the diff): only grepped, not read in full (verifier "FIRST verification command" and documentation lane confirmed by grep).
- Rendered size and the canary (`TestRenderedContextUnderTarget`, GANNET): not run, no test run in a prompts review.
- Antigravity-specific behaviour of the new context: not run.
- `skills/init-project/assets/REQUIREMENTS.md` reorder: read, no restatement elsewhere found beyond the grep above.
- The remaining skills' diffs (feature, golang, typescript, preflight, bug, retro, review, init-project) were read in the diff only; their full files were checked just for the section names the context cites.
- Whether `TestHarnessMap` pins the new legend: not checked.

## Outcome

Every finding fixed in the phase; see the next round (round 3, `2026-10-01-review-329a99b-worktree-prompts-r3.md`, APPROVE WITH NITS). The two lost rules landed in review-standards § 2 and context step 4; the escape-in-file-tool rule stays in the bash skill by design.
