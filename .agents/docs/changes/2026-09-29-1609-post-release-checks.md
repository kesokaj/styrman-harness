---
owner: main
date: 2026-09-29
---
# After v1.12.0: every unchecked tier run, two fixes, one gate rule     (date: 2026-09-29)

## Summary
The operator asked, after the release, whether everything in the gate had been fixed, then "I want everything checked before tag" and approved the recommendations. This entry closes the gate's "What was not checked" list and the two leftovers it named:

- **Real-home smoke on both CLIs** (`testing/2026-09-29-verify-real-home-v1.12.0.md`): the real home brought to v1.12.0 with `sh-agent update`; the three smoke checks PASS on Claude Code and Antigravity; the first Antigravity runs of `/dream`, the resume clause and the `dreamer` through `/retro` PASS; the interactive settings write of `/init-project` seen working; the `init-memory-linked-store` matrix row run once, PASS.
- **`TestReinstallBuildsAndInstalls`** stats its own `sh-agent-build-<pid>` instead of globbing the shared scratch (BUG-104; mutation protocol run).
- **`/init-project`** refuses to write through any symlinked path, at every step (the older writes the v1.12.0 security round had named as out of range).
- **`AGENTS.md` § Release:** a gate whose "What was not checked" names an Antigravity run of a changed skill or agent is NO-SHIP until it exists.
- The `dream` and `init-memory-linked-store` matrix rows corrected (the store inside the clone; Step 6's ignore lines allowed).

## Type
fix (v1.12.1: the `/init-project` guard is a fix to released behaviour; the test and the rule are repository-only)

## Files changed
`internal/harness/harness_test.go`, `skills/init-project/SKILL.md`, `AGENTS.md`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/testing/2026-09-29-verify-real-home-v1.12.0.md` (new), `.agents/docs/bugs/BUG-104-reinstall-build-glob-flake.md` and the index.

## Agents dispatched
- `bugs` — BUG-104.
- `verifier` — over the working tree (its log is named in the v1.12.1 gate report).
- No `reviewer`: the source diff is ten lines across three files, read by the main agent (the test's own-path check, one guard sentence, one Release rule); no `security`: no auth, crypto, secrets, input parsing or dependency changed; no `documentation`: the skill file is its own documentation and no README statement changed.

## Verification
The real-home log above; the mutation check on the test (the mutant compiled, the named test failed, restored, green); the verifier's log; the v1.12.1 gate.
