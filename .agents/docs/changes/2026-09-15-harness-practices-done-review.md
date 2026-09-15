---
owner: main
date: 2026-09-15
---
# Harness practices: Done review of c192622..22c106c, two rounds     (2026-09-15)

## Summary

The Done list's `reviewer` and `--docs` review ran over the source range `c192622..22c106c` of `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`, each on Opus with the repository's changed texts pasted (plan D9). Round 1 raised seven findings, all fixed in the working tree. Round 2 confirmed five fixes, showed one finding contradicted ADR-022 (rebutted, its clause reverted), and found that the fixes to preflight Step 4 and to `/init-project`'s `CLAUDE.md` branch introduced three defects. Under Definition of Done step 2 the loop stopped there, and two questions go to the operator instead of a third patch. Five bug records were written, including the Phase 6 matrix kit's pipe-escape defect.

## Type

task

## Files changed

- `skills/init-project/SKILL.md` — :44 no longer keeps an existing `CLAUDE.md`'s content with `@AGENTS.md` added; when the user declines to fold it, the file is left untouched (R2, plan D42). :8 and :66 no longer claim that Antigravity loads project skills and agents; :66 names the 2026-09-15 precedence test and says project rules were not tested.
- `skills/preflight/SKILL.md` — :65 every Step 4 docs finding blocks the release, is fixed as ordinary work after NO-SHIP, and `/preflight` runs again; :77 lists "a Step 4 docs finding" as blocking; :85 and :90 no longer carry docs fixes into the release commit or the tag (plan D43). :50 is unchanged from `22c106c`: a round 1 clause was added and then reverted (plan D44).
- `skills/handoff/SKILL.md` — :49 deferred work goes under the active plan's Follow-ups, or the handoff's Next Steps when no plan is active.
- `agents/bugs.md` — :76 names the five graduation homes of `context/AGENTS.md:31`.
- `skills/feature/SKILL.md` :97 and `skills/bug/SKILL.md` :53 — the `documentation` trigger names behavior, interfaces, the plan, or a deploy path.
- `.agents/docs/testing/skill-matrix.md` — task 6.1's rows and waiver (I-24a–d), and Procedure step 3 says a kit unescapes `\|` in check cells (BUG-024).
- `.agents/docs/reviews/2026-09-15-review-c192622-22c106c-{code,docs,merged}.md` — the two lens reports (round 2 replaced round 1) and the merged record with both rounds and the Outcome.
- `.agents/docs/bugs/BUG-024` to `BUG-028` and `INDEX.md` — the kit's pipe escape (fixed); `/init-project` against R2 (fixed); preflight Step 4's routing (open); stale describing text, a fourth occurrence of that shape (fixed); a finding applied against ADR-022 (fixed).
- `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` — D42, D43, D44; boxes 6.1–6.3 ticked and 6.4 in progress.
- `ROADMAP.md` — "Now" names the Phase 6 matrix and the Done list as the next steps.

## Corrections to earlier entries

`changes/2026-09-15-harness-practices-phase5.md` describes two behaviours reversed here: preflight Step 4 fixing docs itself with the fixes riding in the release commit, and `/init-project` keeping an existing `CLAUDE.md`'s content with `@AGENTS.md` added. It also says preflight Step 1 asks for the reviewer "when the release diff is a bug fix"; Step 1's clause is about uncommitted work that "reaches 50+ lines or 5+ files or is a bug fix", and Step 2's release-diff threshold is unchanged (ADR-022).

## Agents dispatched

reviewer (code lens, rounds 1 and 2), reviewer (docs lens, rounds 1 and 2), all on Opus; bugs (BUG-024 to BUG-028). This entry was written by the main agent.

## Verification

The main agent ran, each bare with its exit status read: `go test -race ./...`, `gofmt -l .`, `go vet ./...`, a scratch install and `status`, and `render`; the results are in the commit that carries this entry. Not run: a sandbox start of each CLI for these fixes, because the Phase 6 matrix rows are running in the sandbox home; the Done list's final verification covers it.

## Open

- Requirement R2: may a project's `CLAUDE.md` keep other content with `@AGENTS.md` added, or must it stay the single line? (plan D42)
- The release gate: how should preflight Step 4 handle docs findings, including a finding a prior gate rebutted? (plan D43, BUG-026)

## Requirements

n/a
