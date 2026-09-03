---
owner: bugs
date: 2026-09-03
id: BUG-005
status: fixed
area: agents/verifier
---
# BUG-005: `verifier` ran in a fresh worktree and could not see uncommitted work

## Symptom
Sandbox task run on Claude Code: `verifier` returned FAIL, "the described change set does not exist in this working tree" — clean `git status`, no `internal/version/`. The main thread retried, was refused when it pointed the agent at the primary tree, and verified by hand (demo-app `BUG-001`).

## Repro
Any Claude Code session with uncommitted changes: dispatch `verifier`; it starts in `.claude/worktrees/agent-<id>/` at `HEAD`.

## Root cause
`agents/verifier.md` declared `isolation: worktree`. Claude Code honours it by checking `HEAD` out into a throwaway git worktree, so the agent sees the last commit, never the change under test — and the Definition of Done runs verification *before* the commit. Its log also lands in the worktree and is lost. Antigravity ignores the key, which is why the same task passed there.

## Fix
Removed `isolation: worktree` from `agents/verifier.md`; the description no longer promises an "isolated context". The agent still has its own context window and never edits source, which was the point of the isolation. Verified in the sandbox by dispatching `verifier` on an uncommitted change after reinstall.

## Prevention
A Claude-only agent key is not free: `isolation: worktree` is only right for an agent that must work on a committed base, never for one that judges uncommitted work. Any change to an agent's frontmatter is followed by one dispatch in the sandbox on an uncommitted change.

## Related
`.agents/docs/testing/2026-09-03-verify-sandbox.md`; the demo repo's `.agents/docs/bugs/BUG-001-verifier-worktree-isolation.md`, written by the run that hit it.
