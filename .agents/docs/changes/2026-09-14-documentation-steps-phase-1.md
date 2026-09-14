---
owner: changelog
date: 2026-09-14
---
# Documentation steps (phase 1 of 2)     (date: 2026-09-14)

## Summary

Phase 1 of ten documentation gaps in the harness context and agent definitions. Completed items 1, 3, 8, 9 of `.agents/docs/plans/2026-09-14-documentation-steps.md` — the Definition of Done now names which files `documentation` owns (step 3); mandates `/review --docs` when an interface, command, or deploy step changes (step 2); every agent definition gains a "Records" section stating what it writes, where, and what it never touches; the changelog template gains a "Docs" section. Three cold-reader rounds (six reviews) found and resolved residue on the backfill method for proposal statuses, rejected proposal status lines, and guide file naming; all incorporated in phase 2. The operator confirmed scope and plan before code.

## Type

feature

## Files changed

- context/AGENTS.md — DoD step 2 appends `/review --docs` rule (line 16); step 3 clarified to name README, specs, testing guides, deploy guides, project AGENTS.md, and enforces that new commands/endpoints are not done until documented (lines 17-18); `/review` added to workflow skills list (line 69).
- README.md — Layout section line 87 names `review` as the review command on both CLIs.
- agents/verifier.md — Records section added (lines 85-89) with write patterns, scratch rule, and never-touch list.
- agents/changelog.md — template line 63-64 gains "## Docs" section explaining what to record; Records section added (lines 86-90).
- agents/documentation.md — Records section added before Output.
- agents/explorer.md — Records section added before Output.
- agents/reviewer.md — Records section added before Output.
- agents/security.md — Records section added before Output.
- agents/bugs.md — Records section added before Output.
- skills/preflight/SKILL.md — Step 4 "Docs freshness" (lines 63-65) now runs `/review --docs` over the release diff; description updated (line 3).

## Agents dispatched

reviewer ×6 (cold reads with residue resolution), changelog (this entry)

## Verification

`go run ./cmd/sh-agent list` renders context at 91 lines (budget 100). `gofmt -l . && go vet ./... && go test -race -count=1 ./...` all pass; `/tmp` prose test (BUG-014) green. `grep '^## Records' agents/*.md | wc -l` → 7 (all agents carry the section). Sandbox install of the working tree succeeded; derived Antigravity agents carry the Records section unchanged. Context question on both CLIs (Claude Code and Antigravity) quoting the new DoD step 3 clause "a new command, flag, endpoint, or deploy step is not done until the file that describes it says so" confirmed both CLIs load the context and questions work.

## Docs

README.md line 87, context/AGENTS.md DoD steps 2 and 3 (interfaces to agents and users), agent Records sections (internal rule statements).

## Requirements

R7 (procedures are now documented in the harness: context governs when docs are written, Records sections govern where).
