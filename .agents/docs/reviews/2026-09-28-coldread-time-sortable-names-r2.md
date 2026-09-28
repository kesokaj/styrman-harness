---
owner: reviewer
date: 2026-09-28
---
# Cold-read review (round 2) — time-sortable trail names   (2026-09-28)

Range reviewed: `.agents/docs/plans/2026-09-28-time-sortable-trail-names.md` (plan document, cold-reader check — no code diff exists yet). Round 1: `.agents/docs/reviews/2026-09-28-coldread-time-sortable-names-r1.md`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS (round 1's blocking gap is closed; two smaller gaps remain, one worth a decision before Phase 2)

Round 1's blocking finding — the changelog agent has no Bash, so its rule couldn't fire — is resolved: D3 keeps the agent's tools unchanged and moves the clock read to the dispatcher, which now passes the stamp in the prompt (Interfaces, `agents/changelog.md` edits, `skills/preflight/SKILL.md:86`). I independently checked every file:line citation in Surfaces touched against the live files (`context/AGENTS.md:65`, `skills/handoff/SKILL.md:43,62,145`, `agents/changelog.md:3,30,39,89`, `skills/init-project/SKILL.md:22,106`, `skills/preflight/SKILL.md:86`) and reran the describing-text sweep (`git grep -nE '(changes|handoffs)/[^ \`]*(date|YYYY|<)' -- context agents skills commands rules README.md AGENTS.md`): all match exactly what the plan claims, including the "exactly four hits" count and the size-test numbers (ran `TestRenderedContextFitsAntigravityBlock`: 20,329 of 22,000 B, matches). I also checked the same-day sort concern the plan's own rule text raises ("so name order is time order"): `sort`/`LC_ALL=C sort` on a mixed set confirms digits sort before letters, so a date-only slug starting with a letter (e.g. `research-okf`) sorts *after* a same-day timed name — the plan's E2E criterion 1 already states this precisely ("today's three date-only handoffs sort after it by name") and tests that the agent applies the explicit override rule rather than raw sort; this is correctly handled, not a gap.

Two remaining items, one worth a decision before Phase 2's behavioural setup is built.

## Findings (cold-reader gaps, not a code review)

- **[medium / confidence: medium]** Plan `## Interfaces`, `/init-project` sentence (lines 89–93) vs. `## Failure modes` ("`/init-project` cannot list a directory. It skips the count and says so in the summary.") — the literal command given, `ls <dir> | grep -E '^[0-9]{4}-...' | grep -vE '...'`, cannot distinguish "directory missing/unreadable" (an `ls: cannot access` on stderr, empty stdout, non-zero `ls` exit swallowed by the pipe) from "directory exists and has zero date-only files" — both produce a count of 0 with no signal to trigger the "skips the count and says so" message the Failure modes section promises. An implementer following the Interfaces command literally never emits that message; one who reads Failure modes first adds an explicit `[ -d "$dir" ]` guard the plan doesn't specify the wording for. Decide which directory-existence check and summary wording to use before Phase 1's `/init-project` edit.
- **[low / confidence: high]** `agents/changelog.md` — Interfaces says "`:23`'s rule 1 becomes: 'Use the stamp given in the dispatch prompt…'" but rule 1 ("1. Use the date given in the dispatch prompt.") is at line 22; line 23 is rule 2 ("2. If absent, run `date +%F` yourself."), which the plan says stays unchanged. Verified by direct read with line numbers. Not blocking — the quoted replacement text unambiguously identifies rule 1 by content regardless of the line number, and Phase 1's verify clause (`grep -nE 'HHMM|%H%M|<stamp>'`) is line-agnostic — but it's a citation error in a plan whose main strength is exact line citations, worth a fix in the same pass.
- **[low / confidence: medium]** Phase 2 "Handoff skill" test — "The new `SKILL.md` is copied into the scratch clone as project skill `handoff-probe` (a non-harness name) and invoked as `/handoff-probe`" doesn't say whether the copy also edits the copied file's own frontmatter `name: handoff` to `name: handoff-probe`. This repo's own installer keys skill identity and install destination off the frontmatter `name:` field, not the directory name (`internal/harness/inventory.go` `validName`/`LoadSkills` comments, `internal/harness/plan.go:58`), and `context/AGENTS.md` § Global vs project requires a project skill never reuse a harness name — a copy that leaves `name: handoff` inside risks the sandbox's already-installed harness `handoff` skill shadowing or colliding with it rather than a clean `/handoff-probe` registration. I could not confirm from this codebase whether Claude Code's native project-skill loader keys the slash command off the directory name or the frontmatter `name:`; two implementers could diverge here (one edits the copied frontmatter, one doesn't), and if it's frontmatter-keyed the n-of-3 result would silently measure the wrong skill.

## Answers to the three questions

**(1) Gaps.** Two, both listed above: the `/init-project` "cannot list a directory" behaviour is promised in Failure modes but not achievable from the literal command given in Interfaces (needs an explicit existence check and message before Phase 1); and the Phase 2 handoff-skill probe's copy step doesn't say whether the copied file's frontmatter `name:` is renamed to avoid a same-name collision with the harness's own installed `handoff` skill. The `agents/changelog.md:23` citation is a smaller, non-blocking misreading of the file (rule 1 is at `:22`). Everything else checked (all Surfaces-touched line citations, the describing-text sweep count, the size-test numbers, the digit-vs-letter same-day sort behaviour) matches the live files and the plan's own claims exactly.

**(2) Scope.** Nothing found beyond what `## Scope` and the operator's "go all" on items 1–8 cover. D3 (keep the agent's tools, move the clock read to the dispatcher) is a necessary fix responding to round 1, not new scope. The Follow-ups (`reviews/`, `testing/` timed names) are explicitly noted and not built. No extra endpoint, shape, or unused configuration.

**(3) Readiness.** One BLOCKING question, one NON-BLOCKING:
- BLOCKING: does the Phase 2 handoff-skill probe rename the copied `SKILL.md`'s frontmatter `name:` field to `handoff-probe`, or leave it as `handoff`? Two implementers would build different test setups, and if the loader is frontmatter-keyed the weaker one measures nothing.
- NON-BLOCKING: the exact wording/mechanism for `/init-project`'s "cannot list a directory" message — a reasonable implementer adds an `[ -d ]` guard and a one-line message consistent with the rest of Step 7's summary style; the Failure modes section already states the intent even though the Interfaces command doesn't implement it.
- NON-BLOCKING: `agents/changelog.md:23` vs `:22` — content unambiguously identifies the target rule.

## Criteria

1. Dead code — n/a, no code in this plan.
2. Flow tracing — passes at every hop I could check: changelog-agent naming now traces from dispatcher (preflight `:86`, description `:3`) through to the agent's revised rule 1, closing round 1's gap; handoff-skill and context-sentence paths trace cleanly with exact line matches.
3. Standards conformance — the describing-text sweep is complete (independently reran with a wider path set covering `.agents/docs/`, `README.md`, `REQUIREMENTS.md`, root `AGENTS.md`; no additional live restatement found beyond the four already-known hits and this repo's own append-only past records, correctly out of scope).
4. Directory layout — n/a, no new files/directories proposed by the plan itself; the templates in `skills/init-project/assets/agents-tree/docs/{handoffs,changes}/.gitkeep` confirm the plan's claim that no template file encodes the old pattern.
5. Simplicity — pass; no new dependency, no speculative abstraction; the `/init-project` regex pair is minimal (with the directory-existence gap noted above).
6. Decoupling — pass; the five touched surfaces (context, handoff skill, changelog agent, init-project, preflight) are independently editable, and the plan's own grep-sweep step is the cross-file-drift check.

## Not examined

- No git diff exists yet — this is a plan, not a change; nothing beyond the plan text and the current-state files it names.
- `.agents/scratch/agy-probe/` and `.agents/scratch/trail-search/` (the measured probe data backing `## Why` and the rejected `INDEX.md` alternative) were not opened; their headline numbers are taken from the plan's own text.
- Behavioural-run mechanics for Claude Code's native project-skill slash-command resolution (directory name vs. frontmatter `name:`) were not verified against product documentation or a live run — flagged as a question rather than a confirmed defect.
- `skills/preflight/SKILL.md` and `skills/init-project/SKILL.md` were read only around the cited line ranges, not in full.

design ready
implementation not ready
