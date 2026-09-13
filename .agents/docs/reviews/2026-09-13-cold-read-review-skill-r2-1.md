---
owner: reviewer
date: 2026-09-13
---
# Review — cold-read-review-skill-r2-1   (2026-09-13)

Range reviewed: plan file `.agents/docs/plans/2026-09-13-review-skill.md` (cold-reader check, not a diff) plus the surfaces it names: `context/AGENTS.md`, `skills/audit/SKILL.md`, `skills/review-standards/SKILL.md`, `skills/preflight/SKILL.md`, `agents/explorer.md`, `agents/reviewer.md`, `agents/security.md`, `agents/verifier.md`, `agents/documentation.md`, `README.md`, `AGENTS.md`, `.agents/docs/testing/skill-matrix.md`, `internal/harness/derive.go`.

## Verdict: REQUEST CHANGES (as a cold-reader check: design needs revision)

The plan is unusually thorough for a one-command fold-in, but it has two concrete defects that a first-pass implementer would hit immediately: its own stated verify command for the `/audit` removal is falsified by two pre-existing, unrelated files containing the substring `/audit`, and it leaves one of the two known "vendor's built-in review command" mentions in `skills/preflight/SKILL.md` (the frontmatter description) out of its surfaces list while covering the other (Step 2 body) explicitly. A third issue is a real inconsistency between D2 and D3 over what "untracked" scope means for the working-tree case — two implementers would build genuinely different scope-computation logic from it.

## Findings (gaps to decide before implementation)

- **[medium / confidence: high]** Plan line 58 (Phase 1, task 2) — the stated verify command `git grep -n '/audit'` shows only history under `.agents/docs/` is false today, before any change: it also matches `skills/gke-basics/references/gke-security.md:152` (`pod-security.kubernetes.io/audit: restricted`, a Kubernetes label unrelated to the workflow skill) and `skills/anti-ui-slop/SKILL.md:47` (`` - Explicit audit: `references/audit.md` ``, that skill's own reference file). Neither is under `.agents/docs/`, so the check as literally written will report two hits and fail even after a fully correct implementation. The plan does not say whether these are known-acceptable exceptions or whether the grep pattern needs to be more specific (e.g. anchored to `` `/audit` `` or `/audit ` word-boundary). Two implementers diverge: one edits the unrelated files to make the raw grep pass, the other invents its own filter. Decide the exact verify pattern (or name the two exceptions) before Phase 1.
- **[low / confidence: high]** Surfaces list, plan line 13, vs `skills/preflight/SKILL.md:3` — the plan names "`skills/preflight/SKILL.md` Step 2" as the surface to repoint, and D8 (line 27) gives the exact new Step-2 wording ("read the diff yourself or run `/review`"). But the frontmatter `description:` on line 3 of that file *also* says "...or reviewing a diff (use the reviewer agent or the vendor's built-in review command)" and is not in the surfaces list or covered by any decision. Left as-is, the skill would describe itself with a now-dead phrase right next to a body that no longer uses it. Decide whether the description line is in scope (and what it should say) before Phase 1.
- **[medium / confidence: medium]** D2 (plan line 21) vs D3 (line 22) — inconsistent definition of the working-tree scope's untracked files. D2: "untracked **source** files from `git status --porcelain`" (implies a filetype/source filter). D3's actual dispatch text: "untracked, not ignored: <from `git status --porcelain`, minus `<docs-root>/` and `.agents/scratch/`>" (no source-type filter — any untracked, non-ignored, non-docs file qualifies, including a stray data file, image, or a new doc dropped outside `<docs-root>`). This changes what a `--docs` or `--code` lens actually sees on a first run, and two implementers reading only one of the two decisions would build different scope-collection code. Pick one rule (recommend D3's, since it's the one written as literal dispatch text) and remove the conflicting wording from D2.

## Readiness — questions a first-pass implementer would have to ask

1. **[NON-BLOCKING]** During `--fix`, who actually edits the source for the `--code`/`--simplify`/`--security` lenses (D5, plan line 24)? `reviewer` and `security` (per `agents/reviewer.md`, `agents/security.md`) have no `Edit` tool and are explicitly forbidden from modifying source, so it must be the main/session agent applying the findings itself, not a subagent. The plan never says this in words, but it's the only consistent reading given the agents' own definitions, and it matches how `--tests` is phrased ("the main agent writes the missing tests"). A competent implementer reaches the same conclusion; still worth stating for symmetry.
2. **[NON-BLOCKING]** D9's disambiguation order for a bare token that is simultaneously an existing path and a valid ref/branch name (e.g. a directory named `main` when branch `main` also exists) is not stated as an explicit precedence rule, but D9's own bullet order (flag → PR → range/ref → path) is a reasonable default to infer, and the collision is rare.
3. **[NON-BLOCKING]** D3/D8's exact new wording for `context/AGENTS.md`'s workflow-skills line and `skills/review-standards/SKILL.md`'s description/body is not quoted verbatim (unlike the Step-2 quote in D8) — a reasonable implementer swaps `/audit` for `/review` in place and updates the "built-in" reference to name `/review`, consistent with D8's stated intent.

No BLOCKING questions found beyond the D2/D3 scope inconsistency above, which I am counting as a **gap** (a decision the plan must make, not a question an implementer could reasonably default past) rather than a question, since the plan itself asserts two different things.

## Not examined

- The full text of `.agents/docs/testing/skill-matrix.md` beyond the `audit`/`preflight` rows and its harness (lines ~1–30); did not check every existing row for other `/audit` mentions.
- `internal/harness/plan.go` / `apply.go` (how `skills/` directories are actually installed) — not needed for this plan since it is a pure file-drop, but not read in depth.
- Whether `.agents/docs/changes/`, ADR numbering (`ADR-017`), or `.agents/docs/architecture/DECISIONS.md` already contain a conflicting ADR-017.
- Full text of `agents/security.md` and `agents/verifier.md` bodies beyond frontmatter (frontmatter was sufficient to confirm no-Edit-tool).
- Antigravity-side behavior of parallel `invoke_subagent` calls (assumed working from the existing `/audit` skill's identical pattern, not independently verified against `internal/harness/derive.go` beyond the tool map).

## Criteria

1. Dead code — n/a (plan review, not a diff); `skills/audit/` deletion is the point of the plan and is correctly scoped as removal, not a comment.
2. Flow tracing — one real break found: the D9→preflight `/review` wording is traced correctly by D8, but the parallel occurrence in `skills/preflight/SKILL.md:3` is not traced to any decision (see finding above).
3. Standards conformance — plan format matches the `feature`-skill plan shape (Why/Scope/Decisions/Surfaces/Interfaces/Failure modes/Verification/Phases/Risks); frontmatter has `owner`/`date` as required.
4. Directory layout — new skill lands at `skills/review/SKILL.md`, matching every other workflow skill; correct.
5. Simplicity — no new dependency; reuses existing agents/skills rather than inventing new machinery; appropriate.
6. Decoupling — D3's design (skill computes scope, hands agents plain text, agents unedited) is a clean seam; good.
