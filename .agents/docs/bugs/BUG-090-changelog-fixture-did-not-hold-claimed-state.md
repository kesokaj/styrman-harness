---
owner: bugs
date: 2026-09-28
id: BUG-090
status: fixed
area: process / behavioural runs of agents/*.md
---
# BUG-090: changelog-agent behavioural fixture did not hold the state its dispatch claimed

## Symptom
The first changelog-agent batch in Phase 2 of `plans/2026-09-28-time-sortable-trail-names.md` dispatched "added a GET /healthz endpoint to cmd/server/main.go", while the fixture's `main.go` was `package main; func main() {}` — no endpoint. The agent correctly refused to write a changelog entry for behavior that was not present, so run 1 produced no artifact.

## Repro
`grep -c healthz .agents/scratch/<fixture>/cmd/server/main.go` on the fixture as first built returns 0, not the 1 the dispatch text implied. Logged in `.agents/docs/testing/2026-09-28-verify-time-sortable-names-phase-2.md` § What went wrong.

## Root cause
The fixture was built from the dispatch prose, not from a command that proves the claimed state, and no assertion checked the fixture before the CLI ran. Same mechanism as BUG-087: the arrangement did not hold the fact the dispatch stated.
`Diagnosis missed:` the global hard rule "a test that sets up an unusual condition proves the condition HOLDS before it is trusted" and BUG-087's own prevention rule ("prove a fixture's claimed state with the same command the agent will use") were both on record and neither was applied when this fixture was built.

## Fix
The batch was stopped. The fixture now commits the endpoint itself, with a guard script (`.agents/scratch/tsn-probe/run.sh`) that greps for the line and aborts before the CLI starts if it is missing.

## Prevention
A behavioural-run fixture is built to hold every fact its dispatch states, and the run script asserts each such fact (a grep for the line, a git log for the commit) before the CLI starts; a fixture fact with no assertion is not trusted.

## Related
BUG-087 (same mechanism: "fixture didn't hold the condition it claimed"), now a second recurrence of the same rule going unapplied — BUG-071 and BUG-081 are earlier instances in the same family (see BUG-087 § Related). Three records now share this prevention rule; it is actionable (a concrete grep/git-log assertion) and generalizable (applies to any behavioural-run fixture, not just this plan), and is a graduation candidate for `context/AGENTS.md` § Testing / behavioural-run guidance — not graduated here per dispatch instruction; the row stays in the main table.
