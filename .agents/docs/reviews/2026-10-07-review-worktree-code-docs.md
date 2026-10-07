---
owner: reviewer
date: 2026-10-07
---
# Review — worktree-code-docs   (2026-10-07)

Range reviewed: working tree vs HEAD (4793d5a): context/AGENTS.md, skills/{feature,review,bug,handoff,preflight,retro,operate-gcp}/SKILL.md, agents/{bugs,changelog,documentation}.md (ROADMAP.md and plans/ ignored). Code lens and docs lens in one pass.

## Status: COMPLETE

## Verdict: APPROVE (final, after round 2; round 1 was REQUEST CHANGES)

The worst problem: the new combined `--code --docs` dispatch is specified by two instructions that cannot both be followed. The Docs brief's item 1 marks every review criterion `n/a — docs lens`, but review Step 3 now sends it, under a `## Docs lens` heading, in the same dispatch as the code brief that demands all six criteria. Beyond that, the "failure" definition now contradicts itself, and several living restatements of the changed rules were left stale (operate-gcp, three init-project templates). No contradiction between context/AGENTS.md, bug, feature, review, preflight, handoff and retro on the two-round cap, the 30-line docs rule or the "main agent writes records" rule was found. The size and description pins pass (`go test ./internal/harness -run TestRenderedContext|TestSkillListingBudget|TestAgentDescription|TestLintDescriptionLimit`).

## Findings
- **[medium / confidence: high]** `skills/review/SKILL.md:52,60` vs `:50` — Step 3 puts the code brief and the Docs brief in ONE dispatch ("the code brief, then the Docs brief under `## Docs lens`"), but Docs brief item 1 says "Keep your report template with every criterion marked `n/a — docs lens`". Read literally the reviewer must both apply the six criteria and mark all six n/a; a second agent contradiction of this kind is what Step 5 calls "the rule is the defect". Fix: word item 1 conditionally ("when the docs lens is the whole dispatch, mark every criterion n/a; in a combined dispatch the criteria belong to the code brief and items 2-7 populate a separate `## Docs lens` section"). Preflight Step 2 (`skills/preflight/SKILL.md:~58`) and context step 2 inherit the same ambiguity.
- **[medium / confidence: medium]** `skills/bug/SKILL.md:21`, `agents/bugs.md:3` — the failure list still includes "a test failing for a reason other than the change under test" and "a fix that took more than one attempt", and bugs.md still lists "a failing test that took more than one attempt" as a dispatch trigger, while the same sentences now say a defect the same uncommitted change's review/tests/verifier caught "is not a failure" / "never for a defect ... tests caught". Those items are by nature caught by the change's own tests before commit, so a print-mode run can argue either way (the BUG-042 pattern: a gate that counts nothing checkable). Fix: state the boundary as the one checkable fact (the defect reached a commit, a deploy, or the user), and say explicitly that a multi-attempt fix or foreign test failure inside one uncommitted change is recorded in the review `## Outcome`/changes entry, not as a BUG, or keep them as BUG triggers and carve them out of "Not a failure".
- **[medium / confidence: high]** stale restatements the diff left standing:
  - `skills/operate-gcp/SKILL.md:80` — "Every incident gets a number (the global rule: anything that went wrong)"; the global rule now reads "every failure that escaped its change's own gates" (`context/AGENTS.md:33`).
  - `skills/init-project/assets/agents-tree/docs/architecture/DECISIONS.md:5` — "(the changelog agent writes them)"; ADRs are the main agent's (context step 4, `agents/changelog.md:73`). Pre-dates this diff but is exactly the restatement of the changed record-writing claim.
  - `skills/init-project/assets/agents-tree/docs/bugs/INDEX.md:3,10` — "Maintained by the `bugs` agent" / "Rows move here, by the `bugs` agent"; context now says the main agent writes records, `bugs` only when it lacks the evidence (`context/AGENTS.md:33,18`).
  - `skills/init-project/assets/agents-tree/docs/deploy/local.md:8` — "the `documentation` agent keeps it current (Definition of Done step 3)"; step 3 now sends edits up to 30 lines to the main agent.
  These templates are copied into every new project, so the stale claim is propagated.
- **[low / confidence: medium]** `skills/feature/SKILL.md:93` vs `:57` — the user gate on the one-report cold read hangs on "a BLOCKING question only the user can settle", a classification the author makes with no checkable count; an unattended run can label everything NON-BLOCKING and write D-rows, which is the self-answer path BUG-042 warned about (the count there was reports, which is now exactly one, so the count itself holds; the escape is the classification). Also line 57 allows `[-]` for "the user's same-turn override quoted with its date" while Step 3 lists only the one-phase/under-5-files skip. Suggest: make the reviewer's BLOCKING marks the input (every BLOCKING item is either answered by the user's words or recorded as an `Open:` row), and list the override in Step 3.
- **[low / confidence: medium]** `skills/bug/SKILL.md:55`, `skills/operate-gcp/SKILL.md:83` — a record the main agent writes "in the `bugs` agent definition's format" inherits `owner: bugs` from the template (`agents/bugs.md:30`), misattributing authorship. Say `owner: main` when the main agent writes it.

## Criteria
1. Dead code — pass; nothing orphaned (the removed "three-round"/"cold-read round 4" wording has no living remnant: grep over context, skills, agents, rules, README found none).
2. Flow tracing — fail on one hop: review Step 3 combined dispatch vs Docs brief item 1. Round cap (context 16, review 100/112, bug 15, feature 15, preflight 78) consistent; `--docs` rides preflight Step 2's dispatch, falls back to a lone docs run.
3. Standards conformance — fail: the stale restatements above; size/description/listing pins pass.
4. Directory layout — pass; no new files in scope.
5. Simplicity (DIY-except-security) — pass; nothing hand-rolled, the change removes dispatches.
6. Decoupling — pass; one home per claim (30-line rule lives in context step 3, restated by pointer in feature/bug/preflight; bugs/changelog/documentation descriptions carry the agent-side half).

## Not examined
Behavioural runs of the changed skills/agents on either CLI (guide § 4); the interactive README map; skills/init-project/SKILL.md body beyond grep hits; `skills/review/SKILL.md` lines outside Steps 3-5; full diff of `skills/handoff` (read, only the two edits). README.md:259 was read and states nothing the diff changed.

## Round 2

Verified the fixes against the working tree; `go test ./internal/harness -count=1` passes (76s).

- Finding 1 (combined dispatch vs Docs brief item 1) — fixed: `skills/review/SKILL.md:60` is now conditional (alone: all criteria n/a; inside `--code --docs`: the code brief fills the criteria, the list is the `## Docs lens` section). No other text restates the old item.
- Finding 2 (failure boundary) — fixed: one fact ("reached a commit, a deploy or the user") is stated identically in `skills/bug/SKILL.md:21`, `agents/bugs.md:3`, `context/AGENTS.md:11,33`, `skills/operate-gcp/SKILL.md:80`; the "in however many attempts" clause resolves the overlap with the "fix took more than one attempt" list item. Residual (not a finding): the list items and the boundary sentence still sit side by side, and the boundary wins by its own wording.
- Finding 3 (stale restatements) — fixed: operate-gcp:80, DECISIONS.md:5, bugs/INDEX.md:3,10, deploy/local.md:8 now agree with context steps 3-4. A repo-wide grep for the old wording ("anything that went wrong", "changelog agent writes", "by the `bugs` agent", "documentation agent keeps") found no living remnant.
- Low `owner:` finding — fixed: `owner: main` stated at `skills/bug/SKILL.md:21,55`; `agents/bugs.md:30` template keeps `owner: bugs` for a dispatched agent, which is correct.
- Low "author-judged BLOCKING" (`skills/feature/SKILL.md:93`) — rebuttal accepted: the cold reader marks BLOCKING, and each user-only item ends as a D-row, rebuttal or `Open:` row a later reader can check. The Step 3 vs line 57 override mismatch is left as a nit.

No defects introduced by the fixes found.

## Outcome

Written by the main agent, 2026-10-07. Round 2 resumed the same agent with `SendMessage` (review Step 5) and approved.

| # | Finding | Outcome |
|---|---|---|
| 1 | The combined `--code --docs` brief contradicted itself | **fixed** at `skills/review/SKILL.md:60`: Docs brief item 1 is now conditional |
| 2 | The failure boundary could be argued either way | **fixed**: one fact, "reached a commit, a deploy or the user", worded the same at `skills/bug/SKILL.md:21`, `agents/bugs.md:3`, `context/AGENTS.md:11,33` and `skills/operate-gcp/SKILL.md:80` |
| 3 | Stale restatements | **fixed**: `skills/operate-gcp/SKILL.md:80`, and in the init-project templates `docs/architecture/DECISIONS.md:5`, `docs/bugs/INDEX.md:3,10` and `docs/deploy/local.md:8` |
| low | A main-written bug record still said `owner: bugs` | **fixed**: `owner: main` at `skills/bug/SKILL.md:21,55` |
| low | The BLOCKING gate is judged by the author | **rebutted**: the cold reader marks an item BLOCKING; the author decides only whether the user alone can settle it, and records every item as a checkable D-row or rebuttal (`skills/feature/SKILL.md:93`). Round 2 accepted this. |
| nit | Step 3 and the checkbox at line 57 differ on the user override | **left**: implementer; the override wording is unchanged from before this change |
