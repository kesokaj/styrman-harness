---
owner: reviewer
date: 2026-09-21
---
# Review — worktree-docs   (2026-09-21)

Range reviewed: working tree — `git diff HEAD -- .agents/docs/bugs/INDEX.md .agents/docs/plans/2026-09-21-context-consolidation.md context/AGENTS.md context/dev-setup.md skills/bug/SKILL.md skills/feature/SKILL.md skills/golang/SKILL.md skills/review-standards/SKILL.md`, plus untracked `internal/harness/size_test.go`

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The consolidation is mostly a clean, faithful trim — the three named moves (BUG-066 and BUG-068/082's mutation protocol to `skills/review-standards/SKILL.md` § 2, BUG-069's guard rule to `skills/golang/SKILL.md` § Error Handling) land correctly and `.agents/docs/bugs/INDEX.md` is updated to match. But two concrete rules were deleted from `context/AGENTS.md` Hard rules with no new home anywhere in the living docs, not folded into either of the two skills that absorbed the rest of that paragraph. Both are real procedural content, not citations or examples.

## Findings

- **[medium / confidence: high]** `context/AGENTS.md` (review step, was line 16 at HEAD) — dropped: the clause "a saved diff file only when the tree will stand still, and the reviewer pins the blobs it read" is gone from the rewritten step 2 (`... given the diff only (a git range or a path list, never the intent); below that read...`). This is a real procedural rule about how a reviewer must handle a saved-diff dispatch (only when the tree is frozen; pin the blobs read), and it is the exact mechanism this review's own dispatch instructions rely on ("a saved diff file only when the tree will stand still, and the reviewer pins the blobs it read"). `git grep` for "pins the blobs" and "tree will stand still" across `agents/`, `skills/`, `README.md`, `commands/` finds nothing — `agents/reviewer.md`'s own description only says "Give it the diff or the plan only, never the intent," which does not carry the blob-pinning/tree-freeze condition. No new home; fix by restoring the clause in `context/AGENTS.md` or moving it into `agents/reviewer.md`.
- **[medium / confidence: high]** `context/AGENTS.md` (instrumentation bullet, was line 57 at HEAD) — dropped: "a self-hosting suite refuses to start when its port is already bound, takes the port from an overridable variable, and probes the process it started rather than 'something answers'" is removed from the "Silence from your own instrumentation is not evidence" bullet, leaving only the general principle ("A green run is evidence only when it demonstrably hit the thing under test"). This is a concrete, actionable test-design rule (not just an example citation), and `git grep` for "self-hosting", "port is already bound", "overridable variable", "probes the process" across the repo (excluding changes/bugs/reviews/retros/handoffs/DECISIONS.md) finds no other occurrence. No new home.
- **[low / confidence: medium]** `context/dev-setup.md` Backend bullet dropped "stdlib first, every dependency justified" — not a finding on its own since `context/AGENTS.md` Hard rules already states this globally in the same file set ("Otherwise stdlib first, every dependency justified"), so it is redundant removal, not a drop. Listed here only because it was flagged during the check and rejected after verification — no action needed.

Everything else checked out:
- `context/AGENTS.md` Project layout: the removed ASCII tree (`.agents/{rules,skills,agents,workflows,commands}`, `scratch/`, the full `docs/` subtree) is preserved verbatim in `skills/init-project/SKILL.md` (lines ~69-72), which is the canonical scaffolding source; the new prose points at `/init-project`. Not a drop.
- The removed sentence "The harness is the `styrman-harness` repo, managed by `sh-agent status|update|list`; a new skill, agent, command, or rule is a file dropped into it, never code" has its substance in this project's own `AGENTS.md` ("Adding a skill, agent, command, or rule is a file drop — never a code change"), which overrides the global context here per context/AGENTS.md's own rule 1. Not a drop.
- The attribution-settings sentence dropped from `context/AGENTS.md` Git section is fully preserved in `README.md` (Install, Daily commands table) and `.agents/docs/deploy/{local,remote}.md`. Not a drop.
- `context/dev-setup.md` Auth bullet dropped the literal `gcloud auth configure-docker …` command, but `skills/gcloud-setup/SKILL.md:105` still has it, and the bullet already points at `(gcloud-setup skill)`. Not a drop.
- `skills/retro/SKILL.md` still defines the `Scope: project` / `Scope: harness` tags dropped from the context's parenthetical; the context text still describes the same behaviour in prose. Not a drop.
- `internal/harness/size_test.go` (new, untracked): its `antigravityRulesBlockBytes = 24000` and `renderedContextBudgetBytes = 22000` match `REQUIREMENTS.md` § Known limits and `.agents/docs/bugs/INDEX.md` BUG-085 and `.agents/docs/architecture/DECISIONS.md` ADR-035 exactly. No mismatch.
- `.agents/docs/bugs/INDEX.md` BUG-066/068/069/082 rows point at the correct new sections, and those sections (`skills/review-standards/SKILL.md` § 2, `skills/golang/SKILL.md` § Error Handling) do contain the corresponding content, verified by reading the diff.
- `.agents/docs/plans/2026-09-21-context-consolidation.md`: only a checkbox tick ("go all", 2026-09-21) — no claim to verify against code.
- `skills/bug/SKILL.md` and `skills/feature/SKILL.md`: both add a one-line pointer to the mutation protocol now centralized in `review-standards` § 2; consistent with the move.

## Criteria

1. Dead code — n/a, docs lens.
2. Flow tracing — n/a, docs lens.
3. Standards conformance — pass except the two dropped rules above (no restatement anywhere in the living docs the brief scopes: `README.md`, `AGENTS.md`, `.agents/{rules,skills,agents,commands}/`, `.agents/docs/specs/`, `.agents/docs/testing/guide.md`, `.agents/docs/deploy/`, `context/`, `agents/`, non-vendored `skills/`).
4. Directory layout — n/a, docs lens; the report was written to `.agents/docs/reviews/` as required.
5. Simplicity (DIY-except-security) — n/a, docs lens.
6. Decoupling — n/a, docs lens.

## Not examined

Nothing in scope was skipped: all eight changed/untracked files were read in full (diff for the eight tracked files, whole-file read of the untracked `size_test.go`), and every removed/changed rule's subject nouns were grepped against the living-docs set the brief names. `.agents/docs/testing/guide.md` was checked only for the self-hosting/port terms (not read end-to-end) since the brief's grep target was that one rule; no other claim in this diff touches testing procedure.
