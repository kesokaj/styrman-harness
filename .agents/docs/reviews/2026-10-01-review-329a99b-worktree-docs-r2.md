---
owner: reviewer
date: 2026-10-01
---
# Review — 329a99b worktree, docs lens, round 2   (2026-10-01)

Range reviewed: working tree against 329a99b (`git diff 329a99b` plus untracked files)

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Eight of the nine round-1 findings are fixed in the files as they stand; `handoff:36` is not. Two low residues remain: the missing exception at `:36` and a "holds only items the plan will consume" clause that sits beside the new "empty at every checkpoint" rule.

## Round-1 findings, re-checked
1. `skills/handoff/SKILL.md:168` resume chain now matches `context/AGENTS.md:9` (handoff, Now, plan, changes, bugs INDEX, TODO, search --open, memory show; plan-postdates exception). Fixed.
2. `:66` now says TODO empty at every checkpoint and ROADMAP holds Now and Next. Head fixed; see finding below for the tail.
3. `:36` NOT fixed: item 5 still says the newest handoff wins on conflict with no exception (the exception is only at `:168`). See finding 1.
4. `agents/documentation.md`: ROADMAP is "not yours" (:48), description and :41 name TODO and ROADMAP Now as the main agent's. Fixed.
5. `AGENTS.md` rule 1 names the redacted send to the CLI's model; rule 4 carries the ADR-049 exception. Fixed.
6. `README.md:71` says olive. Fixed (README:65 legend is consistent).
7. `skills/docker/SKILL.md:108` points at `docs/deploy/local.md` § Operator steps; matches `context/dev-setup.md:8`. Fixed.
8. 15,000-byte pin named in `AGENTS.md:25-26`, `README.md:275-277`, `testing/guide.md:25-27` (`-run` regexp runs both tests). Fixed.
9. Template `deploy/local.md` § Operator steps comment names the standing root-step permission. Fixed.

## Findings
- **[low / medium]** `skills/handoff/SKILL.md:36` — still reads "the Resume Protocol gives the newest handoff the win on conflict" with no "except an open plan that postdates it" (`context/AGENTS.md:9` carries it). Add the clause; it is the sentence a corrections writer relies on.
- **[low / medium]** `skills/handoff/SKILL.md:66` — after "`TODO.md` is empty at every checkpoint" the same sentence still says "`TODO.md` holds only items the plan will consume"; `context/AGENTS.md:18` says only "empty at every checkpoint". The same wording is in the templates (`skills/init-project/assets/TODO.md:3`, `ROADMAP.md:3`) where it is qualified by "Empty at every checkpoint", and the repo's own `TODO.md:3` says "Empty when the work is done". Reads as a mid-plan working list, so not a defect by itself; trim the handoff clause to "TODO.md is empty at every checkpoint" to remove the apparent conflict.

## Restatement greps (restatement-grep.sh)
- "TODO.md second", "Later / Ideas": only the recon brief (plans/2026-09-30-recon-audit-fix.md) names them, as past findings. Clean.
- "22,000": historic plans plus the live pins above; REQUIREMENTS/README/AGENTS/guide all consistent. Clean.
- "amber": README:65 (report back, correct), frontend skill examples (colour names, unrelated), no stale "Dream (amber)". Clean.
- "settings.json": `deploy/local.md:104,110` (attribution key and Claude's permission prompt, unrelated to the root-step record), `memory-loop.md:47` ("may also live in .claude/settings.json; the harness records it in docs/deploy"), `docker` fixed. Clean.

## Criteria
1. Dead code — n/a, docs lens
2. Flow tracing — n/a, docs lens
3. Standards conformance — restatement sweep above; two low nits
4. Directory layout — n/a
5. Simplicity — n/a
6. Decoupling — n/a

## Not examined
Code under internal/ and cmd/; rendered context bytes were not re-measured (render is 86 lines, under the 100-line guide claim); untracked files are not seen by the restatement script (read by name only: BUG-137 and the review/verify logs, not read for content); the interactive map page.

## Outcome

Both nits fixed: handoff line 36 carries the open-plan exception; line 66 now says TODO.md is empty at every checkpoint and holds, between checkpoints, only items the active plan will consume. APPROVE WITH NITS stands as the docs-lens verdict for Phase 3.
