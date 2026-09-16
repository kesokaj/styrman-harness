---
owner: changelog
date: 2026-09-16
---
# Retro proposals applied — retro 2 (2026-09-16)

## Summary
Six harness-scope findings from `/retro --harness`, all approved by the operator in one turn and applied. Proposals refine procedural prose across skills and agents: `/review` Step 3 grepping project test suites (crew BUG-061); `verifier` "Long suites, once" bullet (run minutes-long suite once in background, poll output file); `/preflight` Step 2 Security naming dispatch coverage ("read end to end or named under `## Not examined`"); `context/AGENTS.md` § Hard rules extending BUG-044's git-status rule with `git log --oneline @{u}..HEAD` and `git tag` (untracked commits or tags a stop sign); `skills/retro/SKILL.md` Step 6 rule for proposal Status lines ("applied <date> (<commit>)" written after review and verification exist); and `agents/security.md` probe repository gating (created with `git init`, exit status read, every command uses `git -C`, BUG-045 closed). Bugs index: BUG-023 and BUG-044 moved to Graduated; BUG-045 (security probe fell through, undone by main agent, nothing pushed); BUG-046 (sh-agent-update timer installed dirty tree 17:04:48, open); BUG-047 (release tag-binding check tautological, fixed by deletion) recorded by bugs agent.

## Type
process/rules

## Files changed
- `skills/review/SKILL.md` — Step 3 post-dispatch brief greps project test suites (changed: "the same" → "its")
- `agents/verifier.md` — § What to run, new bullet for long suites (run once in background, poll output file for result line)
- `skills/preflight/SKILL.md` — Step 2 Security bullet (every release-diff path read end to end or named `## Not examined`)
- `context/AGENTS.md` — § Hard rules (extended after any dispatch that can write: `git status --short --untracked-files=all`, `git log --oneline @{u}..HEAD`, `git tag`; untracked commits or tags are stop signs)
- `skills/retro/SKILL.md` — Step 6 rule (proposal Status line "applied <date> (<commit>)" written after its review and verification exist)
- `agents/security.md` — probe repository created with `git init`, exit status read, every later command named with `git -C`
- `.agents/docs/bugs/INDEX.md` — BUG-023 and BUG-044 moved to Graduated table

## Agents dispatched
One Definition of Done pass over the whole change, shared with `2026-09-16-signed-tag-release-procedure.md`: reviewer ×2 (docs lens, round 1 REQUEST CHANGES, round 2 APPROVE), security ×2 (the release procedure; its round 2 also produced BUG-045), verifier ×2, bugs ×2 (BUG-023 and BUG-044 graduated; BUG-045, 046, 047 recorded), changelog ×1 (both entries).

## Verification
`go test -race -v -count=1 ./...` PASS (pass 79, fail 0, skip 1). Install to `.agents/scratch/verify-retro-2` succeeds with drift none and `sh-agent list` 27/7. Sandbox reinstalled three times as text changed; each start of Claude Code and Antigravity headless quoted the new Hard-rules sentence from deployed context verbatim; Antigravity's newest `cli-*.log` with zero `Failed to parse` lines. Logs: `.agents/docs/testing/2026-09-16-verify-retro-2-and-release-procedure.md` (round 2; round 1 compared by mistake against older log and was corrected).

## Docs
`skills/review/SKILL.md`, `agents/verifier.md`, `skills/preflight/SKILL.md`, `context/AGENTS.md`, `skills/retro/SKILL.md`, `agents/security.md` all changed (procedural prose, no user interface, no deploy step, no README change).

## Requirements
n/a — no interface, config, or deploy path changed in user-facing direction.
