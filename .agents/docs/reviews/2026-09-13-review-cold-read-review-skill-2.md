---
owner: reviewer
date: 2026-09-13
---
# Review — cold-read-review-skill-2   (2026-09-13)

Range reviewed: plan file `.agents/docs/plans/2026-09-13-review-skill.md` (no code diff exists yet — this is a pre-implementation cold read), cross-checked against `agents/reviewer.md`, `agents/security.md`, `agents/verifier.md`, `agents/explorer.md`, `agents/documentation.md`, `skills/audit/SKILL.md`, `skills/review-standards/SKILL.md`, `skills/preflight/SKILL.md`, `context/AGENTS.md`, `README.md`, `AGENTS.md`, `.agents/docs/testing/skill-matrix.md`.

## Verdict: REQUEST CHANGES (design needs revision before Phase 1 starts)

The worst problem: three of the plan's five lenses (`--code`, `--simplify`, `--docs`) all route to the same `reviewer` subagent (D3), and that agent's own filename convention (`agents/reviewer.md:59`, `$docs_root/reviews/$(date +%F)-review-<slug>.md`) and the plan's own artifact-slug rule (D3: `review-<slug>`, D3's slug list has no lens component) do not distinguish between them. Any combination of two or more of those lenses — including the flagship `--all` example the plan itself gives twice (interfaces block, lines 39/43) — dispatches the same agent in parallel to write the *same path*, a race/clobber. The plan's own Verification step (3) expects "four artifacts (explorer brief..., reviewer, security, verifier)" for `--project --all` (line 54), i.e. it implicitly assumes only one reviewer artifact even though `--all` fans out to three separate reviewer dispatches under D3 — the plan contradicts itself on this point. This must be decided (lens-qualified filenames, or a single reviewer dispatch carrying all three briefs) before Phase 1, because it changes both the skill's dispatch logic and the merged-record linking logic.

A second, related problem: `--simplify` ("criteria 1, 5 and 6 only... style is not a finding") and `--docs` (compare against README/specs, no six-criteria section at all) ask the `reviewer` agent to produce reports that do not fit its own hardcoded template (`agents/reviewer.md:37`, "every report addresses all six — explicitly, even when a criterion passes"; :80-86, a fixed `## Criteria` section listing all six). `agents/reviewer.md` is not in the plan's "Surfaces touched" list (line 32), so the plan is relying on the per-dispatch brief silently overriding the agent's own persistent instructions — workable in principle (the existing `skills/audit/SKILL.md:30` already customizes the reviewer's brief for its own report), but the plan never says this is the mechanism, and never specifies what report shape a `--docs`-lens or `--simplify`-lens reviewer run should actually produce instead of the six-criteria template. A first-pass implementer has no template to write into the skill file for these two lenses.

Third, a concrete arithmetic error: current skill count is 27 (`ls -d skills/*/ | wc -l` = 27, confirmed live). Deleting `skills/audit/` (-1) and adding `skills/review/` (+1) nets to 27, not 26. The plan states the target as 26 in three places (Scope, line 13; Verification, line 54) — this is wrong and, if followed literally, produces an incorrect README.md/AGENTS.md skill count and a verification tier (`go run ./cmd/sh-agent list` → 26 skills) that will never pass against a correct implementation, or will pass only if the implementer additionally and silently drops a skill nobody asked to drop.

## Findings
- **[high / confidence: high]** plan `.agents/docs/plans/2026-09-13-review-skill.md:22` (D3) vs `agents/reviewer.md:59` — three lenses (`--code`, `--simplify`, `--docs`) route to `reviewer`, all writing to an identical `reviews/<date>-review-<slug>.md` path when combined (as `--all` does, per the plan's own usage examples at lines 39 and 43); a parallel dispatch collision, not merely a style nit. Fix: qualify the artifact slug by lens (e.g. `review-<slug>-code.md`), or route `--simplify`/`--docs` through a single reviewer dispatch that carries all active reviewer-routed briefs at once and writes one file.
- **[high / confidence: medium]** plan line 22 (D3, `--simplify`/`--docs` briefs) vs `agents/reviewer.md:37,80-86` — the agent's own report template requires all six numbered criteria every time; the plan's alternate briefs assume a different report shape but the plan neither edits `agents/reviewer.md` (absent from Surfaces touched, line 32) nor specifies the alternate template inline. An implementer must invent the missing template or guess whether the brief text is meant to override the agent file.
- **[medium / confidence: high]** plan lines 13 and 54 ("skill counts → 26"; "`go run ./cmd/sh-agent list` → 26 skills") — arithmetic error. Current count is 27; deleting `audit` and adding `review` nets to 27, not 26. Verify with `ls -d skills/*/ | wc -l` before and after.

## Gaps — must be decided before implementation (design)

1. **Reviewer-dispatch collision for combined lenses** (cite: plan:22 D3, `agents/reviewer.md:59`) — as above. Not decided.
2. **Report template for `--simplify`/`--docs` reviewer dispatches** (cite: plan:22 D3, `agents/reviewer.md:37,80-86`) — the plan assumes the reviewer agent can be steered off its fixed six-criteria template by brief text alone; this is untested and unstated as a mechanism.
3. **Skill-count arithmetic** (cite: plan:13, plan:54, live `ls -d skills/*/` = 27) — target number is wrong.
4. **`--project --all` "four artifacts" claim** (cite: plan:54) undercounts: under D3, `--all` on `--project` dispatches `explorer` + `reviewer`(×3 briefs: code/simplify/docs) + `security` + `verifier` = up to 6 distinct dispatches, not "four artifacts." The e2e check needs to state the real expected artifact count once (1) is resolved.

## Readiness — questions a first-pass implementer must ask (or answer alone)

- **BLOCKING** — "When two or more of `--code`, `--simplify`, `--docs` are active together, how are the reviewer dispatches and their artifact filenames disambiguated?" Two competent implementers would build genuinely different things here (one one-file-per-lens, one one-shared-reviewer-call, one unknowingly ships a race).
- **BLOCKING** — "What report structure does a `reviewer` dispatch produce for `--simplify` and `--docs`, given `agents/reviewer.md`'s template is fixed to the six-criteria format and is not in scope to edit?" Without an answer, an implementer either edits an out-of-scope file or ships a skill whose brief the agent may ignore.
- NON-BLOCKING — "How is a positional target token distinguished from a lens flag during parsing?" D2 enumerates target and lens forms exhaustively enough that a reasonable implementer can write the parser (flags starting `--` matching one of the five lens names, `--fix`, `--high`, or `--project` are flags/target-selectors; everything else positional is target).
- NON-BLOCKING — "What happens when `--project` is combined with a positional path target (e.g. `/review internal/api --project`)?" Not specified, but a sensible default (error, or `--project` wins) is a one-line implementer decision with no stated consequence either way.
- NON-BLOCKING — "Does `--fix` re-dispatch the reviewing agent after a fix to confirm it, for lenses other than `--tests`?" D5 references "the Definition of Done's review rule," whose "three rounds" language already implies iterative re-review; a reasonable implementer follows that reference rather than needing a fresh answer.
- NON-BLOCKING — should the new skill carry an `argument-hint:` frontmatter field, as `audit` and `preflight` do? Not mentioned; a reasonable implementer adds one matching the usage block.

## Not examined

Did not read `internal/harness/plan.go` or `derive.go` to confirm the installer needs no change for a plain skill add/delete (the plan's Surfaces-touched list implies none, and skill add/delete is described elsewhere as a pure file-drop — this was taken on the strength of `AGENTS.md`'s own "Adding a skill... is a file drop" rule, not independently re-verified against the installer source). Did not run the sandbox e2e scenarios (none exist yet — this is a pre-implementation review). Did not check `commands/` or `.agents/rules/` for any existing `/audit` or `/review` references beyond the files the plan already names.

## Criteria (applied to the plan document itself, not to code)
1. Dead code — N/A, no code yet; the plan does correctly identify `skills/audit/` for full deletion.
2. Flow tracing — fail: the reviewer-dispatch fan-out (finding 1) breaks caller→callee consistency for the `reviewer` agent's own filename contract.
3. Standards conformance — fail: plan's D3 assumes `agents/reviewer.md` behavior it doesn't actually specify or permit editing.
4. Directory layout — pass: new skill at `skills/review/SKILL.md` matches existing skill layout; no stray placement.
5. Simplicity — pass: no new dependency, reuses existing subagents; the design intent (route to existing actors) is sound even though its wiring has gaps.
6. Decoupling — fail (partial): reusing one subagent for three lenses without an explicit multiplexing contract couples the lenses' outputs unsafely (finding 1).
