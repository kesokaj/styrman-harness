---
owner: bugs
date: 2026-09-18
id: BUG-080
status: fixed
area: cmd/sh-agent/main.go (install/update, checkout path)
---
# BUG-080: relative `--home` resolved against the repo directory on a real install, against cwd on `--dry-run`

## Symptom
`sh-agent install --repo <checkout> --home rel-home` from some directory put the installed files under `<checkout>/rel-home` on a real install, while `install --dry-run` with the same flags resolved `rel-home` against the invoking directory. Found by the Phase 1 verifier of `.agents/docs/plans/2026-09-18-sh-agent-works-everywhere.md`, `testing/2026-09-18-verify-sh-agent-works-everywhere-phase-1.md`.

## Repro
`TestRelativeHomeResolvesAgainstTheInvokingDirectory` (`cmd/sh-agent/main_test.go`): the child install created `rel-home` inside a scratch clone of the repository instead of under the test's working directory.

## Root cause
The checkout path of `install`/`update` rebuilds the binary and runs the build's own install as a child process with `cmd.Dir = repo` (`internal/harness` `Reinstall`); `installPassArgs` handed the child `--home` as typed, so a relative path was resolved by the child against the repository directory, not the caller's cwd. `--dry-run` skips the rebuild/re-exec and resolves the same flag in the parent, against the real cwd — the two modes diverged in *where* they acted, not just in whether they wrote.

## Fix
`cmd/sh-agent/main.go` `run`: `home` is made absolute with `filepath.Abs` right after it is resolved, before any verb runs, so the re-exec'd child and the parent agree on the same directory. Test asserts the manifest lands under `<cwd>/rel-home` and that no `rel-home` appears inside the repository; a must-fail mutation disabling the `Abs` call fails the test.

## Prevention
A path handed to a child process started with a different working directory (`cmd.Dir`) is made absolute before the handoff — the child cannot know the parent's cwd. This is also an instance of the `context/AGENTS.md` read-only-mode rule: `--dry-run` must resolve paths exactly as the real run does, or it is not a dry run of the real run.

## Related
BUG-070 (`--check`/read-only-mode divergence causing damage), BUG-046 (`update` installing the working tree). Three records now share the read-only/real-run divergence mechanism (BUG-070, BUG-080, and BUG-046 adjacently) — candidate for graduating a standing check into `context/AGENTS.md` § Hard rules or a reviewer checklist item that diffs dry-run vs real-run path resolution.
