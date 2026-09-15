---
owner: bugs
date: 2026-09-16
id: BUG-036
status: open
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
Not fixed. Remediation named by the finding: one helper using `os.Lstat` that refuses anything but a regular file, used in `RenderContext`, `LoadDocs`, `LoadAgents`, and `LoadSkills`' SKILL.md read, with a test beside `TestSkillFilesAndAdoptCopiesRefuseSymlinks`; devbox's `verify_mirror` also refuses when `find "$HARNESS_SRC" ! -type f ! -type d` prints anything.

## Prevention
A security fix scoped to "the reader that was flagged" is incomplete by default; when a review finds one unsafe reader of a file class (here: symlinks under harness-managed directories), grep every other reader of that same class in the same change, not just the one cited.

## Related
BUG-037 (same review, adjacent adopt.go defect); 2026-09-14 security review (prior symlink finding this one extends); argolis-devbox `2026-09-15-security-worktree.md` finding 1, merged in `2026-09-16-review-worktree-merged.md`.
