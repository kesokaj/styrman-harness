---
owner: bugs
date: 2026-09-16
id: BUG-037
status: fixed
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
`internal/harness/adopt.go`: after the candidate loop resolves `from`, `Adopt` now calls `filepath.Abs` on `from`, runs `validName("skill", src, filepath.Base(from))` on the resolved name, computes `dest` as an absolute path, and refuses — before any `os.RemoveAll` — when either path contains the other, via a new unexported helper `within(parent, child string) bool` (built on `filepath.Rel`, true when child is parent itself or lies inside it). The refusal reads `%s and %s overlap; adopt copies a skill directory from outside the repo's skills/`.

`internal/harness/harness_test.go` adds `TestAdoptRefusesDestructiveSources` (four subtests: `.`, `..`, a skill adopted onto itself, a relative `skills/alpha` from the repo root) and `TestWithin` (a table test for the helper). Both fail without the fix: run against `git show HEAD:internal/harness/adopt.go` in a scratch copy, all four subtests fail (reporting `skills/` removed, the repo removed, and the source removed), and `TestWithin` fails when `within` is replaced by `strings.HasPrefix(child, parent)` on `/a/skills-old` vs `/a/skillsdelta`.

Verification: `gofmt -l .` empty, `go vet ./...` exit 0, `go test -race ./...` exit 0 — logged in `.agents/docs/testing/2026-09-16-verify-bug-037-adopt.md`. Reviewed at `.agents/docs/reviews/2026-09-16-review-bug-037-code.md` and `…-bug-037-docs.md`, both APPROVE.

Independent diagnosis: `.agents/docs/plans/2026-09-16-recon-diagnose-adopt-force-deletion.md` agreed with the hypothesis and sharpened it — `filepath.EvalSymlinks` does not absolutize a relative path, so `.`, `..`, and `skills/<name>` reached `filepath.Base` unchanged. One root cause, three symptoms.

Follow-up, deferred (not part of this fix, addresses a different failure): copy into a temporary directory under `skills/` and rename it over `dest`, so a `--force` replacement can't lose the old skill when the copy fails halfway.

## Prevention
A command that composes a deletion path from user-typed input validates the resolved path is inside its intended subtree and distinct from any ancestor of the deletion root before calling `RemoveAll` — `.` and `..` are not edge cases to special-case after the fact, they are the default inputs `EvalSymlinks` passes through unchanged.

## Related
BUG-036 (same review, the symlink-refusal gap in the sibling code this diff touches); `context/AGENTS.md`'s instruction to promote skills with `sh-agent adopt` (the path that makes this reachable in normal use); argolis-devbox `2026-09-15-security-worktree.md` finding 2, merged in `2026-09-16-review-worktree-merged.md`.
