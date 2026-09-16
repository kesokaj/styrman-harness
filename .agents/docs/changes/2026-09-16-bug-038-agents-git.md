---
owner: main
date: 2026-09-16
---
# BUG-038: a read-only agent ran a destructive git command against the tracked tree     (2026-09-16)

## Summary

While verifying the BUG-036 fix, the `verifier` agent ran `git checkout --` on tracked files to reproduce the pre-fix state, which threw the uncommitted fix away. It recovered the content and reported the near-miss, and the main agent confirmed the tree independently, but the harness had no rule against it: the agent definitions forbade editing source, which reads as "no edit tool", while `git checkout`, `restore`, `stash`, `clean` and `reset` rewrite tracked files with no editor involved.

The rule now exists in both places a subagent reads. Each read-only agent with Bash — `verifier`, `reviewer`, `security`, `explorer` — gains one identical line in its `## Records` section: never a git command that changes the working tree, even to reproduce an earlier state; copy into `.agents/scratch/` or read the old content with `git show <rev>:<path>`. The same rule is one clause on the Git paragraph of `context/AGENTS.md`, added after the docs lens found the context silent on it. The three agents without Bash are left out because they cannot run git.

## Type

bug

## Files changed

- `agents/verifier.md`, `agents/reviewer.md`, `agents/security.md`, `agents/explorer.md` — one line each, after "Never: source, tests, config, or another agent's records".
- `context/AGENTS.md:45` — one clause appended to the Git paragraph; the rendered context stays at 93 lines.
- `.agents/docs/bugs/BUG-038-verifier-git-checkout-discarded-fix.md` and `INDEX.md` — status fixed.
- `.agents/docs/reviews/2026-09-16-review-bug-038-code.md` (APPROVE) and `…-docs.md` (APPROVE WITH NITS), each with its `## Outcome`.

## Agents dispatched

- `reviewer`: the code lens (ADR-022) and the docs lens (a change under `agents/`).
- `bugs`: closed BUG-038 in its own fix commit.

No `verifier` ran: the diff holds no Go code, and the checks below are the ones this repository's `AGENTS.md` names for a change to `agents/` and `context/`; the main agent ran each bare and read its exit status.

## Verification

- `go run ./cmd/sh-agent list` exit 0: 27 skills, 7 agents. `render | wc -l` prints 93 before and after the context clause.
- Scratch installs (`.agents/scratch/bug038-home`, `bug038-home-2`) exit 0 with `drift none — 259 files match the repo`; `## Records` appears exactly once in each of the seven `agents/*.md`, each `.claude/agents/*.md` and each derived `.gemini/config/agents/*/agent.md`; the new line is in both vendors' copies of all four agents, and the context sentence in both vendors' `AGENTS.md`.
- Sandbox (`/home/kesokaj/sandbox/home`), this tree installed, drift none, from fresh clones of the demo repositories. **The agent line:** Claude Code quoted it verbatim from `~/.claude/agents/verifier.md:100` once its `agents/` directory was passed with `--add-dir` (a headless session's reads are confined to its working directory; the first attempt could not read the file and said so); Antigravity quoted it from its derived copy. **The context sentence:** both CLIs quoted it verbatim with no tools. The two Antigravity logs those starts created, `cli-20260916_063710.log` and `cli-20260916_063906.log`, have `Failed to parse` 0 times. Both clones stayed clean.
- **Reviews:** code APPROVE, no findings; docs APPROVE WITH NITS, its medium finding (the context silent on the rule) fixed in the same commit and its low one closed by the record update.

## Docs

`context/AGENTS.md` is the living doc that changed. The docs lens found no other living doc describing what a subagent may do with git, and no skill that dispatches these agents asks them to reproduce an earlier state in a way the rule now refuses.

## Requirements

n/a
