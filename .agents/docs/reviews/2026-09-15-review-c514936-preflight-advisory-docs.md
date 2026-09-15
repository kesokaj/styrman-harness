---
owner: reviewer
date: 2026-09-15
---
# Review — c514936-preflight-advisory-docs   (2026-09-15)

Range reviewed: working tree against c514936 — `git diff c514936 -- skills/preflight/SKILL.md` (numstat 3/3: lines :65, :77, :81; blob read b72b629, matching the diff's index line), docs lens per `skills/review/SKILL.md:48` plus the paths added by `AGENTS.md:46`

## Verdict: APPROVE WITH NITS

Nothing is broken. No living doc still describes the old behaviour: no file says Step 4 docs findings block, that a rebuttal carries forward from a prior gate report, that `/preflight` runs again after a docs NO-SHIP, or that the SHIP path has no findings left. The four findings are all low severity. The most useful one: :65 gives every advisory gap an owner "who fixes it as ordinary work". Before, only low-severity gaps took that route. Now high-severity gaps do too, and :65 still names no place a later session would find that work. Second: :65 says a docs finding blocks at the Definition of Done's `/review --docs`, but that row does not run for every change (F2). This is the fifth round of findings on the Step 4 mechanism (BUG-026). None of the four needs new machinery. Each is fixed by naming an existing queue or a limit in one clause, or by deleting the promise that has no queue.

## Findings

- **[low / confidence: medium]** `skills/preflight/SKILL.md:65` — narrowed: `c514936:skills/preflight/SKILL.md:65` → `skills/preflight/SKILL.md:65`. The old text had two routes. A blocking gap was "fixed after NO-SHIP as ordinary work … and `/preflight` runs again", which forced the fix before release. Every other gap was "recorded as `left` with an owner", an Outcome status defined at `skills/review/SKILL.md:56`. The new text drops both and keeps only "record each in the gate report with its severity, confidence and an owner, who fixes it as ordinary work". So every gap now takes the no-trigger route, including high-severity ones.
  - **Why it fails:** on SHIP, the gate report is committed (`:85`) and never read again. The resume order (`context/AGENTS.md:9`) does not include `reviews/`. `TODO.md` and `ROADMAP.md` are empty at Done (`context/AGENTS.md:18`). A home for deferred work exists by meaning in `skills/handoff/SKILL.md:49` ("the active plan's Follow-ups, or the handoff's Next Steps"), but :65 does not point there. The owner's fix therefore depends on the main agent recognising a gate-report line as deferred work.
  - **Fix:** add one clause to :65, such as "and lists it in the active plan's Follow-ups, or the handoff's Next Steps when no plan is active". The alternative is to delete "who fixes it as ordinary work" and leave the gap as a record only.
- **[low / confidence: medium]** `skills/preflight/SKILL.md:65` — narrowed: `c514936:skills/preflight/SKILL.md:77` → `skills/preflight/SKILL.md:65`. The old :77 blocked the release on a high-severity docs gap, however the change had entered. The new sentence "a docs finding blocks where the Definition of Done's `/review --docs` runs, in the task that made the change" reads as a coverage guarantee.
  - **Why it fails:** that row runs only "when the change touches an interface, a command, or a deploy step" (`context/AGENTS.md:16`). It does not run for Trivial work (`context/AGENTS.md:25`, steps 4 and 5 only) or for docs-only changes of ten lines or fewer (`context/AGENTS.md:21`). In those cases a docs gap now blocks nowhere, and the sentence does not say so. The global hard rule says a coverage claim names what it does not cover. This repository's `AGENTS.md:46` widens the row to every change under `context/`, `agents/` and `skills/`, but `/preflight` is deployed to projects without that rule.
  - **Fix:** name the limit, for example "…runs (a change to an interface, command or deploy step); other gaps are only recorded here". Or drop the clause.
- **[low / confidence: medium]** `skills/preflight/SKILL.md:65` — dropped: `c514936:skills/preflight/SKILL.md:65`. The removed rule read the newest prior gate report's `## Outcome` and cited a rebutted finding instead of raising it again. The removal is deliberate (ADR-024), and no living doc restates the rule.
  - **Why it fails:** the docs lens greps the living docs at head (`skills/review/SKILL.md:48`), so a gap already rebutted in an earlier gate report comes back in every later release. :65 then assigns it again to an owner "who fixes it", and :65 offers no rebut outcome. The owner either re-fixes something already judged not a defect or re-rebuts it each release. Nothing blocks, so the cost is repeated noise, not a wrong release.
  - **Fix:** either say the owner may rebut it per `/review` Step 4, or accept the noise and say nothing more.
- **[low / confidence: high]** `.agents/docs/changes/2026-09-15-harness-practices-phase7-gate-docs-advisory.md:29` — "The option C docs-lens review reads DOCS-REVIEW-VERDICT-MARKER" is still a placeholder, so this statement is not yet true. **Fix:** replace it with this report's verdict line before the commit. The rest of the entry matches the scope. Its descriptions of :65, :77 and :81 are accurate. Its five `grep -cF` counts (1, 1, 0, 0, 0) and the `^## ADR-024: The release gate` count (1) were re-run and reproduce.

### Checked, no finding
- Old-wording `git grep` at the working tree over `README.md`, `AGENTS.md`, `REQUIREMENTS.md`, `context/`, `agents/`, `commands/`, `rules/` (both empty), `skills/` less the vendored six in `skills/VENDORED.md`, `.agents/docs/specs/` (empty), `.agents/docs/deploy/` and `.agents/docs/testing/sh-agent/guide.md`. `.agents/rules`, `.agents/agents` and `.agents/commands` do not exist, and `.agents/skills` is a symlink searched at `skills/`. Patterns searched: `rebutt`, `Step 4`, `docs gap`, `docs finding`, `runs again`, `no findings left`, `needs none`, `prior gate`, `Steps 2`, `preflight-*`, `high-severity docs`, `blocks the release`, `NO-SHIP`, `docs lens`, `--docs`, `preflight`/`release gate`/`gate report`. Subject-noun pass per BUG-027: `block` near docs/release/ship/gate. The only stale-looking hits are the scope's own lines. `context/AGENTS.md:16`, `skills/feature/SKILL.md:79,97`, `skills/bug/SKILL.md:53` and `skills/review/SKILL.md:56,60` describe task-time review, which the change does not alter.
- `REQUIREMENTS.md`: nothing requires docs findings to block a release. R3 ("no gates") is not engaged, because the gate is written process, not a hook. No violation.
- `.agents/docs/testing/skill-matrix.md:28`, the `preflight` row, checks only for `SHIP` and a report file. It is unaffected. The matrix is not a living doc under this brief, so this is noted, not searched in depth.
- Observation outside the living docs, not a finding: the BUG-026 record's "Phase 7 resolution" paragraph describes the never-committed coverage version in the present tense ("now runs the `--docs` lens only over the commits … not already covered"). Its later "Option C resolution" paragraph is correct. The file name still ends `-open` while its status is `fixed`.

### Unverified
- Whether a real `/preflight` run on either CLI records Step 4 gaps with an owner and still reaches SHIP with gaps present. Nothing was executed.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- No `/preflight` run in the sandbox or any CLI; the behaviour is judged from text only.
- `README.md`, the deploy guides, `.agents/docs/testing/sh-agent/guide.md` and `.agents/docs/testing/skill-matrix.md` were searched with `git grep`, not read whole.
- The other working-tree files (`DECISIONS.md` ADR-023/024, BUG-026, `INDEX.md`, the plan) were read only as intent context. They are not living docs under the brief and were not reviewed as scope. The plan file `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (D49, task 7.7) was not opened.
- The prior round's report `.agents/docs/reviews/2026-09-15-review-c514936-preflight-docs.md` was not opened. Its findings were taken from the BUG-026 record's summary.

## Outcome

- Finding 1, the owner's work has no named place on SHIP — `fixed at skills/preflight/SKILL.md:65` — each gap is filed under the active plan's Follow-ups, or the handoff's Next Steps when no plan is active.
- Finding 2, where a docs finding still blocks — `fixed at skills/preflight/SKILL.md:65` — only in the task that changed an interface, a command or a deploy step, where the Definition of Done runs `/review --docs`; any other docs gap blocks nowhere.
- Finding 3, a gap rebutted at one gate is raised again at the next — `rebutted:` a trade-off recorded in ADR-024's Consequences and plan D49: the gate keeps no state across runs for docs findings, and an advisory gap raised again costs a line in the report, never a blocked release.
- Finding 4, the changes entry's placeholder — `fixed at .agents/docs/changes/2026-09-15-harness-practices-phase7-gate-docs-advisory.md` — the verdict is filled in.

