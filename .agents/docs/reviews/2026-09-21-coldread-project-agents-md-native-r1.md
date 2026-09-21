---
owner: reviewer
date: 2026-09-21
---
# Cold-reader review — project-agents-md-native plan, round 1   (2026-09-21)

Range reviewed: `.agents/docs/plans/2026-09-21-project-agents-md-native.md` and the surfaces it names (`skills/init-project/SKILL.md`, `context/AGENTS.md`, `REQUIREMENTS.md`, `AGENTS.md`, `README.md`, `.agents/docs/architecture/DECISIONS.md` ADR-005, `.agents/docs/testing/skill-matrix.md`, `internal/harness/size_test.go`, `internal/harness/plan.go`, `internal/harness/render.go`).

## Status: COMPLETE

## Verdict: implementation not ready

Two blocking gaps, both about surfaces the plan's own "Surfaces touched" section never names.

## (1) Gaps

- **[high / confidence: high]** `REQUIREMENTS.md:11` (R2 — "Standards, not vendor formats") states: "The context follows the agents.md standard and is named `AGENTS.md` at every level (`~/.claude/CLAUDE.md` **and a project's `CLAUDE.md`** are one-line imports)." This is the exact claim the plan reverses (a project's `CLAUDE.md` no longer exists), and it is a named, numbered requirement — immutable without the operator's explicit approval per the global rule ("propose the exact edit → approval → ADR → update `REQUIREMENTS.md` and `ROADMAP.md` in one commit"). The plan's "Surfaces touched" list (lines 19-24) never mentions `REQUIREMENTS.md` or R2, its Decisions section carries no D-row proposing an R2 edit or citing the operator's approval of that exact text, and Phase 1's grep clause ("grep `AGENTS.md`/README for a claim that this repo has one") does not include `REQUIREMENTS.md`. Shipping this plan as scoped either leaves R2 stale (a description outliving the behaviour it describes, and a violation of the immutable-requirements rule) or requires an out-of-process edit nobody proposed. This is also a live recurrence risk of **BUG-025** (`.agents/docs/bugs/INDEX.md:26`): "I-23b's edit left project `CLAUDE.md` not a one-line import, violating REQUIREMENTS.md R2" — its prevention rule ("a plan's cold-reader check and rehearsal compare every edit against REQUIREMENTS.md before it lands") is exactly what this check is catching now.
- **[high / confidence: high]** `.agents/docs/testing/skill-matrix.md:22`, the `init-deploy` row, is the fixture the plan's own Verification section names for its behavioural tier ("`/init-project` headless in a fresh scratch repo … answers in the prompt as the matrix's `init-deploy` row"). Its prompt text explicitly instructs the agent: "Do everything the skill says, including the `.claude` symlinks, **`CLAUDE.md`**, `.gitignore` lines, and the commit." After this plan's `SKILL.md` edit, the skill no longer writes `CLAUDE.md` — so the fixture's own prompt now conflicts with the skill it is meant to exercise, and directly contradicts the plan's stated success criterion at line 39 ("`ls CLAUDE.md` afterwards must fail"). This is a test fixture in a verify clause that no longer produces (and actively works against) the condition the plan's E2E criterion claims. The plan does not list `skill-matrix.md` as a surface touched, so nothing directs an implementer to update this row before running the behavioural tier.
- **[low / confidence: low]** The `Why` section's evidence ("Verified today on this box … as the memory page documents") is not backed by any file in this repo — `HERON-4412`/`2.1.277` appear nowhere else under `.agents/docs/`. Not blocking (it is background rationale, not an implementation detail), but worth a one-line testing/change record so the claim is checkable later rather than living only in chat memory.

## (2) Scope

Everything built (SKILL.md edits, this repo's own `CLAUDE.md` removal, the `context/AGENTS.md` Project layout line, ADR-036, v1.6.0 MINOR) is covered by a D-row citing the operator's "Your recommendation. Approved." (D1, D2) or the implementer's own measurement (D3). No unapproved scope creep found.

## (3) Readiness

Could a first-pass implementer ship without asking? No — two BLOCKING questions:

- **BLOCKING**: Does R2 (`REQUIREMENTS.md:11`) get edited in the same commit as ADR-036, and what is the exact proposed wording — has the operator approved that specific text (not just the general "Your recommendation. Approved." on the Why paragraph, which never quotes an R2 edit)?
- **BLOCKING**: Is `.agents/docs/testing/skill-matrix.md:22`'s `init-deploy` prompt updated to drop the `CLAUDE.md` instruction before the plan's own behavioural verification tier runs, and is that update in scope for Phase 1 or a prerequisite?
- **NON-BLOCKING**: Should a testing/ log record the 2.1.277 native-read verification cited in `Why`, so it's checkable outside chat memory?

## Criteria (per the six-point lens, applied to the plan document)

1. Dead code — n/a to a plan file; the plan does correctly identify `skills/init-project/assets/CLAUDE.md` for deletion (`git rm`).
2. Flow tracing — traced correctly through `SKILL.md` Step 3's loop/bullet and the update-rule at `:22`; missed the R2 and skill-matrix restatements above.
3. Standards conformance — fails the global "every other restatement of a changed rule is changed with it" rule (`REQUIREMENTS.md` R2) and the project's own immutable-REQUIREMENTS rule.
4. Directory layout — n/a; no new files proposed outside existing conventions (ADR under `architecture/DECISIONS.md`, review under `reviews/`).
5. Simplicity — appropriately minimal; correctly rejects a `--claude-stub` flag as unneeded code.
6. Decoupling — n/a; single-skill, single-doc change with clean boundaries.

## Not examined

`internal/harness/harness_test.go` and `pathguard_test.go` (opened only enough to confirm the size-test budget mechanics, not read end to end — out of scope for this plan since global stub is Out). `.agents/docs/testing/2026-09-21-*` historical verification logs (spot-checked for the `20,205` figure only). `skills/init-project/assets/CLAUDE.md` and other `assets/` templates were not opened (not needed beyond confirming the deletion target exists in the plan's own citation). Did not run any command, per the read-only instruction.
