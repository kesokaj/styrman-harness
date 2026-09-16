---
owner: bugs
date: 2026-09-16
id: BUG-036
status: fixed
area: internal/harness (symlink handling)
---
# BUG-036: the symlink refusal is incomplete

## Symptom
Security review of the argolis-devbox mirror refresh (`argolis-devbox/.agents/docs/reviews/2026-09-15-security-worktree.md`, finding 1) found the symlink refusal added by this plan only covers the skill walk (`internal/harness/plan.go:108`) and `adopt`'s copy (`internal/harness/adopt.go:65`). `RenderContext` (`internal/harness/render.go:25,42`), `LoadDocs` (`internal/harness/inventory.go:197`) and `LoadAgents` (`internal/harness/inventory.go:158`) still read `context/*.md`, `rules/`, `commands/`, `workflows/*.md` and `agents/*.md` with `filepath.Glob` plus `os.ReadFile`, which follow symlinks.

## Repro
Not reproduced by exploit; read from the code paths named above. A symlink committed under `context/`, `rules/`, `commands/`, `workflows/` or `agents/`, pointing at any file the installing user can read (e.g. a credential file in their home), is copied at 0644 into `~/.claude/` and `~/.gemini/config/`; for `context/`, its content is also concatenated into the AGENTS.md both CLIs send to the model provider. Devbox's `032_setup-harness.sh` mirror-verify gate (`sha256sum -c`, `find -type f`) does not catch an added link either. Confirmed not reachable at a81ef4b: `git ls-tree` shows the tree holds one symlink, `.agents/skills`, which install does not read as content (verified by the reviewing security agent, not independently re-run here).

## Root cause
The symlink refusal was added at the one place a prior review (2026-09-14) flagged — the skill walk and adopt's copy — without auditing every other reader that walks harness-managed directories and follows links by default (`os.ReadFile`/`filepath.Glob` do not distinguish a regular file from a symlink target).

## Fix
`internal/harness/inventory.go`: a new `checkManagedTree(repo)` runs once from `LoadInventory` and refuses a symbolic link at any of `context/`, `agents/`, `commands/`, `rules/`, `workflows/` or `skills/` itself, and at every entry directly inside the five flat directories. `skills/<name>` may still be a link — `npx skills add -g` leaves one and `LoadSkills` resolves it. `readRegular`, an `Lstat`-first read, covers the one read the boundary check cannot reach, `skills/<name>/SKILL.md`; `LoadSkills` now skips a directory only when that file is missing and returns every other error (closing the swallowed-refusal gap round 1 of review found). `internal/harness/render.go`: `RenderContext` runs `checkManagedTree` before reading `context/`.

`internal/harness/pathguard_test.go` adds `TestContextAgentsAndDocsRefuseSymlinks` (a link in place of a file, six cases including `context/AGENTS.md`), `TestSymlinkedSkillMdIsRefused`, `TestSkillDirWithoutSkillMdIsSkipped`, and `TestManagedDirsRefuseSymlinks` (a link in place of each of the six directories). Every test fails without its fix: the file cases against d167cf8, the `SKILL.md` case with the swallowing `continue` restored, and the directory cases against the leaf-only guard.

Independent diagnosis (`.agents/docs/plans/2026-09-16-recon-diagnose-symlink-reads.md`) agreed with the hypothesis: one guard was added in cceee54 to the skill walk and `adopt`'s copy, never to the three loaders. Three review rounds (`.agents/docs/reviews/2026-09-16-review-bug-036-code.md`, Outcome section): round 1 found the refusal swallowed in `LoadSkills`, round 2 found a managed directory that is itself a link still followed — two rounds on one mechanism, escalated to the operator, who chose one boundary check over more per-reader guards; round 3 approved. The docs lens (`…-bug-036-docs.md`) approved twice. Verifier PASS at `.agents/docs/testing/2026-09-16-verify-bug-036-symlink-reads.md`: counts pass 54, fail 0, skip 1.

Follow-ups, deferred: a direct `RenderContext` test with a linked `context/` (that call site is reached only through `Plan` today); and, outside this repository, argolis-devbox's decision whether its `verify_mirror` boot step refuses anything but regular files and directories, carried in its own handoff.

## Prevention
A security fix scoped to "the reader that was flagged" is incomplete by default; when a review finds one unsafe reader of a file class (here: symlinks under harness-managed directories), grep every other reader of that same class in the same change, not just the one cited.

## Related
BUG-037 (same review, adjacent adopt.go defect); 2026-09-14 security review (prior symlink finding this one extends); argolis-devbox `2026-09-15-security-worktree.md` finding 1, merged in `2026-09-16-review-worktree-merged.md`.
