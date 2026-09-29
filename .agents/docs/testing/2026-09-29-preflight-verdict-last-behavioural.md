---
owner: main
date: 2026-09-29
---
# Behavioural runs — preflight writes its verdict last (BUG-105)   (2026-09-29)

Change under test: `skills/preflight/SKILL.md` Step 5 gained "Write the verdict line LAST, when every other section is final: read each item under 'What was not checked' against every blocking condition below, the project's own included, and an item a condition names makes the verdict NO-SHIP whatever the rows above say (BUG-105)."

## Status: PASS — with the rule the verdict is written last in 6 of 6 reports; without it, in 0 of 6

## Setup

Claude Code 2.1.284, headless, `--setting-sources project,local`. The arm's skills were copied into the clone's `.claude/skills/`; the agents were passed with `--agents` (the working-tree definitions); the arm's rendered context was appended. Arms: the working tree (with the rule) and the install made before the change (without it). Clones sat under `~/sandbox/search/pf/`, outside every checkout.

Fixture: the `tasklog` sample project, tagged `v0.1.0`, plus one commit adding a project skill `.agents/skills/notes/` and a `## Release` section to its `AGENTS.md`: "a gate whose 'What was not checked' names an Antigravity run of a skill or agent the release range changed is NO-SHIP until that run exists." The run script asserted the tag and the rule's sentence before each start. Prompt: `/preflight v0.2.0`, no remote, never tag, run nothing on Antigravity.

Two variants, three runs per arm each:

- **Plain:** no run of the `notes` skill is recorded anywhere.
- **Wrong CLI (BUG-105's own case):** a testing log records a Claude Code run of the skill, PASS 3 of 3, and no Antigravity run.

## Results

| Variant | Arm | Verdict NO-SHIP, the missing Antigravity run named | Verdict line after "What was not checked" |
|---|---|---|---|
| Plain | with the rule | 3 of 3 | 3 of 3 (lines 45 to 57 of the report) |
| Plain | without | 3 of 3 | 0 of 3 (lines 9 to 11) |
| Wrong CLI | with the rule | 3 of 3 | 3 of 3 (lines 50 to 57) |
| Wrong CLI | without | 3 of 3 | 0 of 3 (lines 11 to 13) |

Both texts reached the right verdict on both fixtures, so the verdict itself does not tell the arms apart here. The ordering does: with the rule the verdict is the last thing written, after the unchecked items are final.

**Antigravity 1.2.13** (sandbox home installed from the working tree, the plain fixture): completed on the first attempt with no skip flag; `VERDICT: NO-SHIP`, the missing run named, the verdict line (41) after the section (16).

## A void batch

The first four runs put the arm's agents into the clone's `.claude/agents/`. Preflight's own check for project agents that shadow the review agents stopped all four before any review, in both arms, so they measured the shadow check and not the rule. The agents were then passed with `--agents`.

## What was not exercised

- A release range in which the old text writes SHIP against its own unchecked items, as v1.12.1's gate did; both fixtures here state the rule plainly enough that either text obeys it.
- Interactive sessions.
