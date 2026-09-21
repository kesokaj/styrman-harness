---
owner: main
date: 2026-09-21
---
# Review — working tree (context consolidation Phase 1), code + docs lenses   (2026-09-21)

Range reviewed: working tree — changed files: .agents/docs/bugs/INDEX.md, .agents/docs/plans/2026-09-21-context-consolidation.md, context/AGENTS.md, context/dev-setup.md, skills/bug/SKILL.md, skills/feature/SKILL.md, skills/golang/SKILL.md, skills/review-standards/SKILL.md; untracked, not ignored: internal/harness/size_test.go.

## Verdicts

- code lens — APPROVE, no findings (`2026-09-21-review-worktree-code.md`)
- docs lens — REQUEST CHANGES, two findings (`2026-09-21-review-worktree-docs.md`)

## Findings

| # | Lens | Severity / confidence | Finding |
|---|---|---|---|
| 1 | docs | medium / high | `context/AGENTS.md` review step (was line 16) dropped "a saved diff file only when the tree will stand still, and the reviewer pins the blobs it read"; no home in `agents/`, `skills/`, `README.md`, `commands/`. |
| 2 | docs | medium / high | `context/AGENTS.md` instrumentation bullet (was line 57) dropped the self-hosting-suite test-design rule (refuses to start when its port is bound, port from an overridable variable, probes the process it started); no restatement anywhere. |

## Not checked

The pruning review's per-row byte arithmetic (trusted by the code lens); the sandbox CLI starts (a separate log, `testing/2026-09-21-verify-sandbox-pruned-context.md`).

## Outcome

1. rebutted: dropped on purpose — `agents/reviewer.md` has no blob-pinning duty (the docs lens confirms: no home anywhere), so the allowance was never enforced; removing it makes the rule stricter, not narrower (a git range or a path list only). Recorded as plan D16 and in the pruning review, row 9.
2. fixed: `skills/review-standards/SKILL.md` § 2, the must-fail-mutation checklist item now ends with the self-hosting-suite sentence; loaded by the reviewer, `/review` and `/preflight`, which is where a suite's design is judged.
