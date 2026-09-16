---
owner: explorer
date: 2026-09-16
---
# Recon — diagnose-adopt-force-deletion   (2026-09-16)

## Question
Why does `sh-agent adopt <src> --force` compute a deletion target that can be
`<repo>/skills` (src=`.`), `<repo>` itself (src=`..`), or the same directory
as the source (src=`skills/<name>`)? Rank root causes, no fix, no destructive
repro.

## Summary
Single root cause, three manifestations. `internal/harness/adopt.go:24`
calls `filepath.EvalSymlinks(c)` on the raw candidate string, and Go's
`EvalSymlinks` does **not** absolutize a relative path that contains no
symlink components — it returns the input essentially unchanged (`.` stays
`.`, `..` stays `..`, `skills/foo` stays `skills/foo`). `adopt.go:37` then
does `dest := filepath.Join(repo, "skills", filepath.Base(from))`. When
`from` is `.`, `filepath.Base(".")` is `.`, so `dest` cleans to
`<repo>/skills` — the whole skills tree. When `from` is `..`,
`filepath.Base("..")` is `..`, so `dest` cleans to `<repo>` itself. When
`from` is a relative multi-segment path like `skills/foo` run from the repo
root, `from` stays relative while `dest` is built as an absolute path under
`repo`; the two strings differ but name the identical directory on disk, so
`os.RemoveAll(dest)` at `adopt.go:42` deletes the source before
`copyTree(from, dest)` at `adopt.go:48` ever reads it. I verified the
`EvalSymlinks` behavior directly with a throwaway Go probe (deleted after
use, nothing under source or `.agents/scratch` left behind) — see Evidence.

## Relevant files
- `internal/harness/adopt.go` — the entire bug surface: candidate resolution (15-32), dest computation (37), remove-then-copy ordering (38-50).
- `internal/harness/harness_test.go:485` `TestAdopt` — existing test; check whether it exercises `src="."`, `src=".."`, or same-directory adoption (it likely only covers a by-name/by-path adopt into a fresh dest, not these edge cases — see Open questions).
- `internal/harness/pathguard_test.go:112` `TestSkillFilesAndAdoptCopiesRefuseSymlinks` — shows the project already treats symlink-handling in this area as a known-risky seam (guards against copying symlink targets), but that guard is orthogonal to this bug (it fires inside `copyTree`, after the destructive `RemoveAll` has already run).

## Flows
1. `cmd/sh-agent/main.go` case `"adopt"` (~line 202) parses `--force` and calls `harness.Adopt(repo, home, fs.Arg(0), *force)`. `repo` comes from `harness.FindRepo`, independent of `src`.
2. `Adopt` (`adopt.go:15`) builds `candidates` = `[src, home/.agents/skills/src, home/<ClaudeRoot>/skills/src, home/<AgyRoot>/skills/src]` and iterates, taking the **first** candidate where `filepath.EvalSymlinks(c)` succeeds and `<resolved>/SKILL.md` exists. `src` itself is always tried first, so if the operator's cwd (or `src` literally) already contains a `SKILL.md`, `candidates[0]` (`src` as given) wins before any of the home-rooted candidates are tried — meaning a bare `.` or `..` never gets a chance to be reinterpreted as an absolute path via the other candidates.
3. `from` is set to whatever `EvalSymlinks` returned for that winning candidate — for `.` and `..` this is the literal string `.`/`..`, unresolved to an absolute path, confirmed by direct probe (below).
4. `dest = filepath.Join(repo, "skills", filepath.Base(from))` (`adopt.go:37`) — `filepath.Base` on `.` returns `.`, on `..` returns `..`; `filepath.Join` cleans these away, producing `<repo>/skills` and `<repo>` respectively.
5. `adopt.go:38-47`: if `dest` exists and `--force` is set, `os.RemoveAll(dest)` runs unconditionally, with **no check that `dest` and `from` (or its absolute form) are distinct, and no check that `dest` is not `repo` or `repo/skills` itself.**
6. `copyTree(from, dest)` (`adopt.go:48`, defined 54) then walks `from` — if `from` was inside what just got deleted (cases 1 and 3), the walk fails or copies nothing meaningful; if `from` is disjoint from the wiped `dest` (case 2, `..`), the walk still proceeds but the deletion of `<repo>` has already destroyed the checkout.

## Evidence (path computation only, nothing deleted)
Ran a throwaway `go run` probe (written under `.agents/scratch/probe/`, deleted immediately after; nothing under source was touched) printing `filepath.EvalSymlinks` and `filepath.Base` for `.`, `..`, and `skills/foo`, from inside a `skills/foo/` directory containing `SKILL.md`, and separately from the parent (`repo`-root-equivalent) directory:

```
(run from inside skills/foo, which has SKILL.md)
EvalSymlinks(".")          = "."           base = "."
EvalSymlinks("..")         = ".."          base = ".."
EvalSymlinks("skills/foo") = "" (lstat error: no such "skills" here)

(run from the parent, i.e. repo-root-equivalent)
EvalSymlinks(".")          = "."           base = "."
EvalSymlinks("..")         = ".."          base = ".."
EvalSymlinks("skills/foo") = "skills/foo"  base = "foo"
```

This exactly reproduces all three reported symptoms:
- src=`.` (from inside a skill dir) → `from="."` → `dest = Join(repo,"skills",".") = <repo>/skills`.
- src=`..` (from inside a skill dir) → `from=".."` → `dest = Join(repo,"skills","..") = <repo>` itself.
- src=`skills/foo` (run from repo root) → `from="skills/foo"` (relative, unresolved) while `dest="<repo>/skills/foo"` (absolute) — different strings, same directory on disk → `RemoveAll(dest)` deletes the very tree `copyTree` is about to read as `from`.

## Conventions
- The codebase already has one precedent for treating `adopt`'s copy path as security-sensitive: `pathguard_test.go` and the symlink refusal in `copyTree` (adopt.go:65-67), added per "security review 2026-09-14" per its comment. This bug sits upstream of that guard.
- Per this repo's `AGENTS.md`, path-affecting bugs graduate to `BUG-NNN` records with a prevention rule; no such record currently references `adopt.go` (not confirmed by index read — see Open questions).

## Risks
- All three reported symptoms trace to the same two lines (`adopt.go:24`, `adopt.go:37`); a fix must not treat them as three separate cases.
- `os.RemoveAll(dest)` currently runs before any read of `from`'s absolute identity — order-of-operations risk independent of the `EvalSymlinks` quirk: even with correct absolutization, nothing stops `dest` from being computed as `repo` or `repo/skills` if `filepath.Base(from)` is legitimately empty/`.`/`..` for some other resolved path shape.

## Falsification
- Candidate (EvalSymlinks non-absolutization + Base(dot-path)): falsified if `filepath.EvalSymlinks` on this Go toolchain/OS returns an absolute path for a plain `.`/`..` input — probe above (Go 1.x on this Linux box, verified directly) shows it does not; would need to be re-checked only if the Go version changes this stdlib behavior.
- Candidate (same-directory delete-before-copy for `skills/<name>` case): falsified if `from` were absolutized before comparison to `dest`, or if `Adopt` checked `from == dest` (via `filepath.Abs`) before calling `RemoveAll`; currently no such check exists (adopt.go:38-47), so not falsified.

## Not examined
- Did not read `TestAdopt` (`harness_test.go:485`) in full to confirm which of these three cases it already covers vs. leaves untested.
- Did not check `.agents/docs/bugs/INDEX.md` for an existing BUG record on this area.
- Did not examine `FindRepo`/`ClaudeRoot`/`AgyRoot` definitions in detail; assumed from AGENTS.md description that `.agents/skills` is a symlink into `skills/`, which is consistent with case 3 reproducing through that symlink too, but did not independently trace `harness.ClaudeRoot`/`AgyRoot` constants.
