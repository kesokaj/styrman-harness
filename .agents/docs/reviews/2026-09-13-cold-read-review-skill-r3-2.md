---
owner: reviewer
date: 2026-09-13
---
# Cold read — plan `2026-09-13-review-skill.md`   (2026-09-13)

Range reviewed: N/A — cold-reader check of a plan file, not a diff.
Plan reviewed: `.agents/docs/plans/2026-09-13-review-skill.md`
Surfaces cross-checked: `skills/audit/SKILL.md`, `skills/review-standards/SKILL.md`, `skills/preflight/SKILL.md`, `agents/explorer.md`, `agents/reviewer.md`, `agents/security.md`, `agents/verifier.md`, `agents/documentation.md`, `context/AGENTS.md`, `README.md`, `AGENTS.md`, `.agents/docs/testing/skill-matrix.md`.

## Verdict: REQUEST CHANGES (plan needs two decisions before implementation)

The plan's factual claims about the existing surfaces all check out — I did not find a single misquoted line, wrong file, or misdescribed behavior anywhere in D1–D9 or the Surfaces list. The weakness is not accuracy, it is that two real scope questions are left implicit: what "review this path" actually inspects when the path has no uncommitted diff (the flagship `/review internal/api --simplify` example in the Interfaces block is ambiguous on its own terms), and whether `--fix` may run under `--project` — a combination the old `/audit` never permitted (audit was report-only; Step 5 was docs/changelog/commit, never source edits) and that the new plan neither offers as an example nor explicitly forbids. Both are the kind of fork where two careful implementers ship materially different tools.

## Findings (gaps requiring a decision before implementation)

- **Gap 1 — scope of a bare path target is underspecified** `plan:21` (D2) and `plan:22` (D3) — confidence: high. D2 defines the "working tree" scope precisely (`git diff HEAD --name-only` plus untracked files) and D3 gives verbatim dispatch text for working-tree, range, PR, and `--project` scopes. For "paths" it only says "for paths, the list" — no verbatim brief. The Interfaces block's own example, `/review internal/api --simplify`, is the natural way a user asks "review this directory's current code," but if `internal/api` has no uncommitted changes, is the scope (a) `git diff HEAD -- internal/api` → empty → D9's "empty target: say nothing to review and stop" fires and the example does nothing useful, or (b) the reviewer reads the named files whole, like a scoped `--project`? The reviewer agent's own fallback ("if no diff is discernible, review what the dispatch names") could paper over this, but the plan's D9 empty-target rule is written as if empty-diff-and-no-paths and empty-diff-with-paths-given are the same case, when they are not. Two implementers will genuinely diverge: one wires paths through the same diff-or-nothing logic as the working tree, the other treats named paths as always-whole-content regardless of git state.

- **Gap 2 — `--fix` under `--project` is neither offered nor forbidden** `plan:24` (D5), `plan:26` (D7), `plan:34-44` (Interfaces block) — confidence: high. The Interfaces block's `--project` example (`/review --project --all`) has no `--fix`, and D7 is explicit that audit's old Step 5 (docs/changelog/commit) is not folded in — but says nothing about whether the new `--fix` machinery (D5) is allowed to run repo-wide. `/audit` never edited source; `/review --project --all --fix` would let the main agent apply edits across the *entire repository* based on four agents' combined findings, a materially larger blast radius than `--fix` on a diff. This needs an explicit decision: forbid the combination ("one scope, `--project` never fixes" — analogous to D9's "`--project` together with paths/range/PR is an error"), or state that it is allowed and describe the round/rollback discipline at repo scale. As written, an implementer could go either way and the two results are not interchangeable.

Neither gap is a factual error — the plan's description of every named file (skills/audit, review-standards, preflight, explorer.md, context/AGENTS.md line 69, the skill count, the grep exclusions for `anti-ui-slop`/`gke-basics`) is accurate; I verified each by reading the cited files directly.

## Questions for the implementer

1. **BLOCKING** — When a path target is given and has no uncommitted diff, does `/review` say "nothing to review" (per D9's literal reading) or fall back to whole-file review of the named paths? The Interfaces block's own path example only makes sense under the second reading, but D9 is written for the first. (Gap 1)
2. **BLOCKING** — Is `--fix` permitted together with `--project`? If yes, does the fix loop apply edits repo-wide across all requested lenses in one pass, or is there a smaller per-file/per-round cap? If no, where does the skill say so (D9's scope-conflict error list does not currently include this pair)? (Gap 2)
3. NON-BLOCKING — Precedence when a single token could parse as both an existing path and a valid git ref (e.g., a directory sharing a name with a tag). D9 lists checks in an order (flag → `#n` → range/ref → path) a reasonable implementer follows top to bottom; the collision is rare enough that either resolution is defensible and the plan implicitly picks one via ordering.
4. NON-BLOCKING — Exact mechanism for the "second merged record the same day gets `-2`, `-3`" rule (D4). This is an existing, already-used repo convention (visible in the current `reviews/` directory itself), so a reasonable implementer just follows precedent (glob existing files, increment).
5. NON-BLOCKING — Fitting all of D1–D9's quoted dispatch text into the ≤140-line budget (Phase 1 checklist) is a compression exercise, not a design fork — the content is fully specified even if the final wording is condensed.

## Not examined

`internal/harness/plan.go` and `derive.go` (whether the installer's Antigravity-derivation logic needs any change for a skill named `review` — the plan states no code change is needed since a skill is a file drop, and I did not verify the installer's frontmatter/description validation against the exact new SKILL.md the plan proposes, since it does not exist yet). I did not attempt to estimate whether the proposed skill text fits in 140 lines. I did not review `agents/changelog.md` or `agents/bugs.md` (not touched by this plan). I did not examine the sandbox e2e harness scripts referenced in the Verification section beyond confirming the `skill-matrix.md` rows the plan says to edit currently exist as described.

## Criteria (six-criteria lens, applied to the plan as a design document)

1. Dead code — n/a, no code yet; the plan correctly identifies `skills/audit/` for full deletion rather than a deprecation shim.
2. Flow tracing — mostly pass; the two gaps above are exactly the two places the "hand agents explicit text, don't edit the agents" design (D3) doesn't fully close the loop between what the skill computes and what the unmodified agents will do with it.
3. Standards conformance — pass; every cross-reference to existing skill/agent text (D1, D8) matches the files verbatim.
4. Directory layout — pass; `skills/review/SKILL.md` new, `skills/audit/` deleted, matches existing skill directory conventions.
5. Simplicity — pass; reuses existing agents and existing dispatch mechanics (Agent tool / `invoke_subagent`) rather than inventing new machinery; no new dependency.
6. Decoupling — pass; the "skill computes scope, agents stay unedited" design is a clean seam; the two gaps are at that same seam's edges, not evidence the seam is wrong.

design needs revision — gaps 1 and 2 above.
implementation not ready — questions 1 and 2 are BLOCKING.
