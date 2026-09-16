---
owner: main
date: 2026-09-16
---
# BUG-036: symlinks in the repo's managed directories were deployed as their targets     (2026-09-16)

## Summary

`sh-agent` refused a symbolic link only inside a skill directory and in `adopt`'s copy. `context/`, `agents/`, `rules/`, `commands/` and `workflows/` were read with `filepath.Glob` and `os.ReadFile`, which follow links, so a link committed there deployed its target's bytes into every consumer home — and for `context/` into the `AGENTS.md` both CLIs send to the model.

The fix went through three review rounds. The first version guarded each reader; round 1 found the guard swallowed in `LoadSkills`, and round 2 found that a managed directory that is itself a link still slipped through. Two rounds on one mechanism is this harness's signal to ask what the mechanism is for, and the operator chose one boundary check. `checkManagedTree` now runs once, from `LoadInventory` and `RenderContext`, and refuses a link at any of the six managed directories or at any entry directly inside the five flat ones. The readers use plain `os.ReadFile` again. `readRegular` survives only for `skills/<name>/SKILL.md`, because `skills/<name>` may itself be a link that `LoadSkills` resolves by design. Round 3 approved.

Found by argolis-devbox's security review of the a81ef4b mirror refresh (finding 1), recorded as BUG-036.

## Type

bug

## Files changed

- `internal/harness/inventory.go` — `checkManagedTree`, called first in `LoadInventory`; `readRegular` kept for the `SKILL.md` read, with `LoadSkills` now returning every error but a missing file.
- `internal/harness/render.go` — `RenderContext` runs `checkManagedTree` before reading.
- `internal/harness/pathguard_test.go` — `TestContextAgentsAndDocsRefuseSymlinks` (six file cases), `TestSymlinkedSkillMdIsRefused`, `TestSkillDirWithoutSkillMdIsSkipped`, `TestManagedDirsRefuseSymlinks` (six directory cases).
- `.agents/docs/bugs/BUG-036-symlink-refusal-incomplete.md` and `INDEX.md` — status fixed.
- `.agents/docs/bugs/BUG-038-verifier-git-checkout-discarded-fix.md` and `INDEX.md` — a new open record: the verifier ran `git checkout --` against tracked files to reproduce the pre-fix state, discarding this fix, then recovered it. Its rule is fixed next, as its own task.
- `.agents/docs/testing/2026-09-16-verify-bug-036-symlink-reads.md` — the verifier log, its third run.
- `.agents/docs/reviews/2026-09-16-review-bug-036-code.md` (round 3, APPROVE, with the three-round ledger in its Outcome) and `…-bug-036-docs.md` (re-run, APPROVE).
- `.agents/docs/plans/2026-09-16-recon-diagnose-symlink-reads.md` — the independent diagnosis.

## Agents dispatched

- `explorer`: the goldfish diagnosis, given the symptom and repro only.
- `verifier`: three runs, one per round.
- `reviewer`: the code lens three times, the docs lens twice.
- `bugs`: BUG-038, and the closing of BUG-036.

## Verification

- **Diagnosis** agreed with the hypothesis: one guard added in cceee54 to two of five directories.
- **Every test fails without its fix.** The six file cases fail against d167cf8 with `Plan err = <nil>`; the `SKILL.md` case fails with the swallowing `continue` restored; the six directory cases fail against the leaf-only guard. Each run was checked for its `=== RUN` lines, because one earlier check had passed vacuously with "no tests to run", and two others had failed to build rather than to assert (BUG-029).
- **Verifier, Status PASS:** counts pass 54, fail 0, skip 1, up from 47 by the last test's parent and six subtests; `gofmt`, `go vet` and `go test -race ./...` clean; `list` 27 skills and 7 agents; `render` 93 lines; a scratch install with drift none. This repository's own `skills/` is a real directory of real directories, so the check leaves it untouched.
- **Reviews:** code REQUEST CHANGES, REQUEST CHANGES, APPROVE; docs APPROVE twice. The Outcome ledger is in the code report.
- **What went wrong on the way:** the verifier's first run discarded the uncommitted fix with `git checkout --` and recovered it (BUG-038); the main agent verified the tree independently afterwards.

## Docs

No living doc changed. The docs lens found none claiming a link may sit at any managed directory or inside the flat ones, and `AGENTS.md`'s description of `.agents/skills` as a symlink stays accurate: the check looks at `repo/skills`, never at `.agents/skills`.

## Requirements

n/a
