---
owner: reviewer
date: 2026-09-28
---
# Review — cold-read: time-sortable trail names   (2026-09-28)

Range reviewed: `.agents/docs/plans/2026-09-28-time-sortable-trail-names.md` (plan document, cold-reader check — no code diff exists yet)

## Status: COMPLETE

## Verdict: REQUEST CHANGES (design needs revision on one point before Phase 2 is trustworthy; everything else in the plan is unusually precise)

The plan is well-specified for the handoff-skill and context-sentence surfaces (exact file:line citations checked against the actual files and all matched), but it builds the "changelog agent's naming" scope item on an assumption the agent's own definition contradicts: `agents/changelog.md:4` gives the agent `tools: Read, Grep, Glob, Edit, Write` — no Bash — so the new rule's primary path ("The file name's time comes from `date +%H%M` run when you write the file") cannot execute in a real dispatch, and nothing in the plan updates the tools list, the agent's own dispatch `description:`, or the global Definition-of-Done dispatch text to supply an HHMM instead. Phase 2's behavioural test for the changelog agent masks this by granting `Bash(date:*)` through `--allowedTools`, a capability the shipped agent definition does not have, so a PASS there is not evidence of what a real dispatch will produce.

## Findings

Cold-reader findings, not a code review — severities describe how much implementation work would stall or misfire.

- **[high / confidence: high]** `agents/changelog.md:4` (tools list) vs. Interfaces "Changelog agent" (plan lines 69–73) — the new date rule assumes the changelog agent can run `date +%H%M` itself, but its `tools:` frontmatter has no Bash. Real dispatches never mention this scope item's `description:` line (`agents/changelog.md:3`, "passing today's date from `date +%F` in the prompt" — not HHMM) or the global Definition of Done step 4 (`context/AGENTS.md:18`), neither of which the plan's Surfaces-touched list names for editing. Absent one of those changes, every real changelog dispatch falls to the rule's own fallback ("cannot run commands → write the date-only name and flag it") — the changelog half of this feature ships inert. This must be decided before Phase 1: either (a) add Bash (or a scoped `Bash(date:*)`) to the agent's tools and treat that as the security-relevant capability change it is, or (b) commit the dispatcher (main agent / global context) to always pass `HHMM` in the prompt, and update `description:` and `context/AGENTS.md:18` to say so.
- **[medium / confidence: high]** Phase 2 "Changelog agent" verify clause (plan lines 122–125) — `--allowedTools` adds `Bash(date:*)` for the behavioural run, which is not in the agent's shipped `tools:` list. This test fixture does not produce the condition it claims to measure: it proves the *rule text* can be followed when Bash is available, not that a production dispatch (no Bash) will produce a timed name. n of 3 passing here would be read as validating the shipped agent, and it does not.
- **[low / confidence: medium]** `## Why` quotes "then 'go all' on scope items 1–8," but no artifact in this repo's trail (checked `.agents/docs/changes/2026-09-28-*`, `.agents/docs/handoffs/2026-09-28-*`, and a repo-wide grep) defines what items 1–8 were. A cold reader cannot check the built Scope (context sentence, handoff skill, changelog agent, init-project warning) against the actual authorization — only against the plan author's own restatement of it in `## Scope`.

## Criteria

1. Dead code — n/a, no code in this plan; the plan itself adds/removes no dead artifacts.
2. Flow tracing — fails at one hop: the changelog-agent naming rule does not trace through to a caller capable of executing it (see high finding above). The handoff-skill and context-sentence paths trace cleanly: `context/AGENTS.md:65` citation, `skills/handoff/SKILL.md:43/62/145`, and `skills/init-project/SKILL.md:22/106` line citations were all read against the live files and match exactly what the plan claims is there.
3. Standards conformance — the Surfaces-touched grep sweep (`git grep -nE '(changes|handoffs)/[^ \`]*(date|YYYY|<)' …`) was independently re-run against other files mentioning `YYYY-MM-DD` (`agents/bugs.md`, `skills/retro/SKILL.md`, `skills/feature/SKILL.md`, `skills/init-project/assets/**`, `TODO.md`, `.agents/docs/deploy/local.md`); none of them describe `handoffs/`/`changes/` naming, so the plan's narrower sweep is not missing a restatement — pass.
4. Directory layout — n/a, no new files/directories proposed.
5. Simplicity — pass; no new dependency, no speculative abstraction. The `/init-project` regex pair is the minimum needed, with a known, already-acknowledged blind spot (a date-only slug beginning with 4 digits is excluded from the warning count the same way it defeats sorting) — accepted risk, not a new one.
6. Decoupling — pass; the four surfaces (context, handoff skill, changelog agent, init-project) are edited independently and the plan's own grep-sweep step is the mechanism that catches cross-file drift.

## Not examined

- No git diff exists yet (this is a plan, not a change) — nothing to review beyond the plan text and the named current-state files.
- `internal/harness/size_test.go` was read only for the budget constant and test shape, not run; the plan's claim of "20,329 of 22,000 B" today was not independently recomputed (would require building/running the test).
- `.agents/scratch/agy-probe/` (the measured cold-start probe backing `## Why`) and `.agents/scratch/trail-search/` (backing the rejected `INDEX.md` alternative) were not opened; their headline numbers are taken from the plan's own text, not verified against the raw run logs.
- `agents/bugs.md`, `skills/retro/SKILL.md`, and other files found to reference `YYYY-MM-DD` were checked only for the narrow `handoffs/`/`changes/` naming pattern (none found), not read in full for unrelated issues.

## Answers to the three questions

**(1) Gaps.** One blocking gap: the changelog agent has no Bash tool (`agents/changelog.md:4`) and the plan neither grants it one nor commits the dispatcher to supply `HHMM`, so the "changelog agent's naming" scope item as specified will not fire in production even though its behavioural test (with `Bash(date:*)` added to `--allowedTools`) will likely pass — a verification criterion masking a real implementation choice that was never made. Everything else checked (context sentence location, handoff skill lines, init-project lines, the describing-text sweep) matches the live files exactly.

**(2) Scope.** Nothing found that the request in `## Why` did not name, other than the unverifiable "scope items 1–8" reference (finding 3) — I could not confirm the four built surfaces are the full, correct set against that quoted authorization because the items list is not on record anywhere in this repo's trail. No extra endpoint, shape, or unused configuration was added beyond what `## Scope` states.

**(3) Readiness.** One BLOCKING question for the author:
- Does the changelog agent gain a Bash tool for this feature, or does the dispatcher always pass `HHMM`? (BLOCKING — two implementers would ship different things: one broadens the agent's tool surface, the other changes dispatch-prompt convention and two other files the plan does not currently list as touched.)

Two NON-BLOCKING notes, resolvable by a reasonable implementer without asking:
- Exact insertion point of the new context sentence relative to "…are append-only; comments say why, not what." (plan gives the sentence and the anchor phrase; either reading produces a correct, if stylistically different, paragraph).
- The `/init-project` date-only regex undercounts a slug beginning with 4 digits — already an accepted, named failure mode for sorting, and the plan's Failure modes section covers the same root cause.

design needs revision
