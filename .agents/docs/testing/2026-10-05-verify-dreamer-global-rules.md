---
owner: main
date: 2026-10-05
---
# Verification — dreamer and dream read the global rule files, behavioural runs   (2026-10-05)

Range reviewed: 507244a..worktree. The tree differs from 507244a in `agents/dreamer.md:32` (What you read § 3 adds the global rule files, `~/.claude/rules/*.md` or `~/.gemini/config/rules/*.md`) and `skills/dream/SKILL.md:30,39` (the same, for the read list and the `graduated` verdict). Operator approval: "You fix both 1 and 2. Approved." (2026-10-05).

## Status: PASS

## Setup

- **Fixture**: `~/sandbox/dr/fixture` (outside every checkout), commit `9afd9de`: `AGENTS.md`, an empty bug index, and one handoff whose Corrections quote the operator: "Don't apologise when I correct you, just state the corrected logic." Each run clone gets an untracked store `.agents/memory/` (gitignored) with one `feedback` memory, `no-apologies-on-correction.md`, and its `MEMORY.md` line. The preamble asserted the commit, the handoff line and the index line; all 4 runs printed `kit: real config unchanged`.
- **Home**: `.agents/scratch/dr-home`, `--vendor claude` from the working tree (896 files), `~/.claude/rules/critical-thinking.md` present (asserted per run: "No apologies when challenged"). `HOME` and `CLAUDE_CONFIG_DIR` pinned there; credentials copied chmod 600 and deleted by name after.
- **Definitions**: `claude -p --agents` with `dreamer` built from the working-tree file (new) or `git show HEAD:agents/dreamer.md` (old); the prompt (stdin) only asks to dispatch `dreamer` with the dispatch line `Store: <clone>/.agents/memory. Docs root: .agents/docs. Since: all. Promoted: none.` and does not repeat the change. Stream transcripts show one `Agent` call with `subagent_type: dreamer` and a `SubagentHandback` in every run.

## Pass criterion

The artifact `retros/<date>-dreamer.md` cites `rules/critical-thinking.md` and has no candidate row about apologies (the fact is already a loaded rule).

## Results

| Definition | Pass |
|---|---|
| new | **2 of 2** (Left alone: "already a rule: global `rules/critical-thinking.md`") |
| old (control) | 2 of 2 |

**The probe does not discriminate on Claude Code.** Claude Code loads `~/.claude/rules/*.md` into every session, subagents included, so the old dreamer sees the rule without being told to read it. The change makes the read list name what the agent already has in context on Claude Code; where a rule file is not in the dreamer's loaded context, it now says to read it. That case is not exercised here.

## Not exercised

- Antigravity runs of the changed `dreamer` and `dream` (required before the next tag by the project's release rule).
- The `dream` skill's `graduated` verdict in a live `/dream` run (skill prose; covered by the install and parse check).
