---
owner: reviewer
date: 2026-09-15
---
# Review — c514936-preflight-docs   (2026-09-15)

Range reviewed: working tree vs c514936 — `git diff c514936 -- skills/preflight/SKILL.md` (post-image blob `d0120ff63a23f93a1057bf7acff7263a7f78a2a6`, hashed after reading). Lens: `--docs`. Living docs searched with `git grep` in the working tree, no revision: `README.md`, `AGENTS.md`, `context/`, `agents/`, `skills/` minus the six skills in `skills/VENDORED.md` (grilling, codebase-design, find-skills, anti-ui-slop, gke-basics, terraform), `commands/`, `rules/`, `.agents/docs/testing/sh-agent/guide.md`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/deploy/`. `.agents/docs/specs/` is empty. `.agents/{rules,agents,commands}` do not exist, and `.agents/skills` is a symlink to `skills/`.

## Verdict: REQUEST CHANGES

The worst problem: under the new Step 4, a prior `reviews/*-docs*` report counts as covering its commits whatever its verdict and whether or not its findings were resolved. Step 5 now blocks only on what the new, smaller docs run reports. A docs review that returned REQUEST CHANGES and was never closed can therefore be skipped, and the release ships past it. This repository already has such a report (F1). The new sentence about rebutted findings also points at a `## Outcome` the named glob does not reach: `/review` writes Outcome into a `-merged.md` record. It also only holds while a finding can come only from the commits it covered (F2). No living doc repeats the removed wording. Process note: these findings land on the same mechanism as BUG-026's three earlier rounds. Under Definition of Done step 2 and BUG-026's prevention rule they go to the operator; they are not patched here.

## Findings

- **[medium / confidence: medium]** `narrowed: skills/preflight/SKILL.md:65 @c514936 (lens over all of <base>..HEAD) and :77 @c514936 → skills/preflight/SKILL.md:65, :77` — Before, Step 4 ran the docs lens over the whole release range, so an unfixed docs gap from an earlier review was found again and blocked. Now the lens runs only "over the commits in `<base>..HEAD` not already covered by a prior `reviews/*-docs*` or `reviews/*-preflight-*` report, deciding coverage as Step 2 does". Step 2 (`:50`) decides coverage only from the `Range reviewed:` line. Step 5 (`:77`) blocks on "a Step 4 REQUEST CHANGES verdict or high-severity docs finding", which means the new run's findings only. Nothing makes an open finding in a covering report block. Concrete case in this repo: `.agents/docs/reviews/2026-09-15-review-641eccd-worktree-docs.md:9` reads `## Verdict: REQUEST CHANGES` with `Range reviewed: 641eccd..fe21671` (:7). Its record `2026-09-15-review-641eccd-worktree-merged.md` has `## Round 1` (:26) and no `## Outcome`. A `/preflight` with a base before 641eccd would treat those commits as covered, skip them, and could print SHIP without that verdict ever being weighed. Fix: a report covers its commits only once it is resolved, meaning its verdict is not REQUEST CHANGES, or its record's `## Outcome` gives every blocking finding `fixed` or `rebutted`. Otherwise its open blocking findings count as Step 4 blocking findings. Unchanged Step 2 (`:50`) has the same property for code reviews; that is outside this scope and not raised.

- **[medium / confidence: medium]** `narrowed: skills/preflight/SKILL.md:65 @c514936 ("a finding it rebutted, the same statement at the same file:line, is cited in the gate report as rebutted there and not raised again") and :77 @c514936 (", not counting a finding a prior gate report rebutted") → skills/preflight/SKILL.md:65 ("A finding rebutted in a covering report's ## Outcome is therefore not raised again")` — Two things were lost.
  (a) The "therefore" holds only if a docs finding can come only from the commits that produced it. The docs lens (`skills/review/SKILL.md:48`) instead greps the living docs at the scope's head for every statement the scope changes. A later, uncovered commit that edits the same rule again brings back the same stale hit at the same file:line. The docs reviewer is never shown the earlier rebuttal. With the Step 5 exception deleted, that re-raised finding blocks, so BUG-026's R2-1 symptom (a correctly rebutted docs finding blocks the release) returns in this case.
  (b) The `## Outcome` the sentence names is usually not in a file the glob matches. `/review` Step 4 (`skills/review/SKILL.md:56`) writes `## Outcome` into `<date>-review-<scope>-merged.md` whenever `--fix` or more than one lens ran, and "the reviewing agents never amend their own past reports". `reviews/*-docs*` does not match `-merged.md`. In this repo, `2026-09-15-review-c192622-22c106c-docs.md` has no `## Outcome`; the Outcome for that review is in `2026-09-15-review-c192622-22c106c-merged.md`.
  Fix: name the covering report's record (`reviews/*-merged*` together with its `*-docs*` artifact) as the place rebuttals are read. Keep the old match test (same statement at the same file:line): a re-raised finding that matches is cited as rebutted and does not block.

- **[low / confidence: medium]** `narrowed: skills/preflight/SKILL.md:65 @c514936 ("is cited in the gate report as rebutted there") → skills/preflight/SKILL.md:65 ("when a prior report covers every commit, cite it in the gate report")` — Before, the gate report always recorded a rebutted docs finding it did not raise again. Now Step 4 requires a citation only when prior reports cover every commit. When coverage is partial, the covering reports and their rebuttals leave no trace in the gate report, unless the reader takes "deciding coverage as Step 2 does" to include Step 2's per-report "cite the covering report" (`:50`). Fix: in Step 4, cite every covering report, as Step 2 does.

- **[low / confidence: high]** `dropped: skills/preflight/SKILL.md:81 @c514936` ("the SHIP path has no findings left and needs none") — The deleted clause was false: Step 4 records non-blocking gaps as `left` on SHIP (BUG-026 item 3), and ADR-024's Consequences records the deletion. It has no replacement, though. Step 4 (`:65`) says non-blocking gaps are "recorded as `left` with an owner" but not where. Step 5's report contents (`:71-75`) have no such item. Step 6 (`:81`) requires `## Outcome` only after NO-SHIP. So on SHIP, where the `left` lines and owners go is undefined, and new `:65` reads rebuttals from "a covering report's `## Outcome`". Fix: one clause in Step 5, for example "non-blocking Step 4 gaps go in the gate report's `## Outcome` as `left` with an owner". Alternatively, rebut by citing ADR-024 if the current wording is judged enough.

Old-wording `git grep` results: none. `newest prior`, `prior gate report`, `not counting`, `needs none`, `SHIP path`, `no findings left`, `not raised again`, `raised again`, `rebut`, `Outcome`, `gate report` and `preflight` have no hit in the living docs that still describes the removed behaviour. The only hits are the new text in `skills/preflight/SKILL.md` and the consistent `context/AGENTS.md:16` ("Findings are fixed or rebutted in writing, in the report's `## Outcome` section"). Plan task 7.7 counts checked: `not already covered` 2, `newest prior` 0, `not counting a finding` 0, `no findings left` 0, `^## ADR-024` 1.

Newest `changes/` entry: `.agents/docs/changes/2026-09-15-harness-practices-phase7-gate-coverage.md` (untracked). Its Files changed section matches the diff at :65, :77 and :81. It repeats F2's claim ("so a finding rebutted in a covering report's `## Outcome` is not raised again"). Its Verification section still holds the placeholder `DOCS-REVIEW-VERDICT-MARKER`, which is work in progress for this review's verdict and not raised as a finding.

Unverified:
- Whether an agent reads "deciding coverage as Step 2 does" as including Step 2's per-report citation (this decides whether F3 applies).
- Whether an agent's "coverage is your judgement" (`:50`) would in practice refuse to count an unresolved REQUEST CHANGES report as covering (this decides F1's likelihood, not its presence in the text).
- No `/preflight` run was exercised; every finding comes from reading the text.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- The working-tree changes to `.agents/docs/architecture/DECISIONS.md`, the plan, `BUG-026-preflight-step4-routing-open.md` and `bugs/INDEX.md`: read for intent only, since the brief excludes them from the living docs, and not reviewed.
- Unchanged Steps 1–3 and the SHIP sub-steps of `skills/preflight/SKILL.md`, except as context for the changed lines.
- The prior review reports `2026-09-15-review-641eccd-worktree-{docs,merged}.md` and `2026-09-15-review-c192622-22c106c-{docs,merged}.md`: only their verdict, `Range reviewed:`, `## Round` and `## Outcome` headings were read, not their findings in full.
- Vendored skills (excluded by the brief). `.agents/docs/specs/` is empty.

## Outcome

- Every finding in this report — `left: operator` — this is the fourth round of review findings on the release gate's docs routing (BUG-026), so under Definition of Done step 2 no patch was attempted. The coverage change it reviewed stays uncommitted while the operator decides; the recommendation put to the operator is to make the gate's docs review advisory, since the Definition of Done's task-time `/review --docs` already blocks on docs findings.

## Outcome update (2026-09-15)

- Every finding in this report — resolved by replacement: the coverage text it reviewed was never committed. The operator approved option C (plan D49, ADR-024), which makes the release gate's docs review advisory and removes the coverage rule, the rebuttal carry-forward and the docs blocking clause these findings concern. Option C has its own docs-lens review, `reviews/2026-09-15-review-c514936-preflight-advisory-docs.md`.
