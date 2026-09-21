---
owner: main
date: 2026-09-21
---
# Context consolidation, Phase 2 — the loop that grew the context now has a delete step     (2026-09-21)

## Summary

Third phase of `.agents/docs/plans/2026-09-21-context-consolidation.md` (BUG-085). The retro skill's target table gains a row for rules that only retell a fixed bug or duplicate a skill (delete or merge), and its apply step runs the size test after any addition to `context/*.md` or `rules/*.md`, so a proposal that overflows the 22,000-byte budget must name what pays for it or go to a skill. The `bugs` agent's recommendation of a permanent home in the context now names the existing sentence the rule would replace or merge into. Three headless runs of the changed agent in clones of this repo did exactly that, 3 of 3.

## Type

feature phase (skill, agent, docs)

## What changed

- `skills/retro/SKILL.md` — Step 2 target table: one row (`Rule that only retells a fixed bug, or duplicates a skill`); Step 6 harness-scope bullet: the size-test sentence (plan D14).
- `agents/bugs.md` § Rules — one bullet after the "three records share a prevention rule" bullet (plan D10; no tool added, `derive.go` untouched).
- Bug records BUG-086 (the installer piped into `head` died on SIGPIPE) and BUG-087 (a behavioural fixture did not hold the state its dispatch claimed, twice), with index rows.
- Plan: Phase 2 boxes and checkpoint ticked.

## Verification

- `go test -race -count=1 ./...` exit 0, gofmt clean (main agent, bare).
- Scratch home and the sandbox home reinstalled: drift none; both CLIs listed the retro skill and quoted its description; zero `Failed to parse skill`.
- `testing/2026-09-21-context-consolidation-phase-2-behavioural.md`: three `claude -p --agents` runs from the working-tree `agents/bugs.md` in clones of this repo, subject BUG-086 — 3 of 3 records recommend widening the existing pipe-into-filter sentence and cite the 22,000-byte pin; subagent tool order read from the transcripts. A first arrangement in demo-service clones failed on its own fixture and could not exercise the rule (BUG-087). Not exercised: Antigravity's derived agent beyond install, status and a start; a control run with the old definition.

## Review

Diff of three lines in two files, read by the main agent. Docs lens (`reviews/2026-09-21-review-worktree-docs-2.md`): APPROVE, no findings; the test and the REQUIREMENTS section the new sentence names exist; no stale restatement.

## Agents dispatched

`reviewer` (docs lens), `bugs` (BUG-086, BUG-087).
