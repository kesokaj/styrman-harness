---
owner: bugs
date: 2026-09-16
id: BUG-037
status: open
area: internal/harness/adopt.go
---
# BUG-037: `sh-agent adopt --force` can delete the harness checkout

## Symptom
Security review of the argolis-devbox mirror refresh (`argolis-devbox/.agents/docs/reviews/2026-09-15-security-worktree.md`, finding 2) found `internal/harness/adopt.go:37-44` computes `dest := filepath.Join(repo, "skills", filepath.Base(from))` where `from` is `filepath.EvalSymlinks` of an operator- or agent-typed source. For source `.` or `..`, `EvalSymlinks` returns them unchanged, so `dest` becomes `<repo>/skills` or `<repo>` itself; `os.Stat(dest)` succeeds, the refusal message suggests `--force`, and `--force` runs `os.RemoveAll(dest)` — deleting every skill, or the whole checkout including `.git` and unpushed work. A second case: a source already equal to `<repo>/skills/<name>` (e.g. `adopt skills/foo --force` from the repo root, or via the `.agents/skills` link) gives `dest == from`, so `RemoveAll` deletes the source before `copyTree` reads it.

## Repro
Not reproduced destructively. `argolis-devbox/.agents/scratch/security-worktree-probe/main.go` computed the paths only (no deletion) and confirmed `src="." → dest="<repo>/skills"` and `src=".." → dest="<repo>"`.

## Root cause
`Adopt` builds `dest` from `filepath.Base(from)` without validating that `from` names a genuine skill subdirectory distinct from `repo` or `repo/skills` — `.` and `..` are valid path syntax that survive `EvalSymlinks` unchanged and collide with the destination root. The bug predates this plan's diff; this diff only added `copyTree`'s symlink refusal without touching the destination computation it wraps.

## Fix
Not fixed. Remediation named by the finding: make `from` absolute; run `validName` on `filepath.Base(from)`; refuse when `from` and `dest` are equal or one contains the other; copy into a temp directory under `skills/` and rename it over `dest` only after the copy succeeds; a test for `.`, `..`, and a self-adopt. Until it lands, never pass `--force` to `adopt`.

## Prevention
A command that composes a deletion path from user-typed input validates the resolved path is inside its intended subtree and distinct from any ancestor of the deletion root before calling `RemoveAll` — `.` and `..` are not edge cases to special-case after the fact, they are the default inputs `EvalSymlinks` passes through unchanged.

## Related
BUG-036 (same review, the symlink-refusal gap in the sibling code this diff touches); `context/AGENTS.md`'s instruction to promote skills with `sh-agent adopt` (the path that makes this reachable in normal use); argolis-devbox `2026-09-15-security-worktree.md` finding 2, merged in `2026-09-16-review-worktree-merged.md`.
