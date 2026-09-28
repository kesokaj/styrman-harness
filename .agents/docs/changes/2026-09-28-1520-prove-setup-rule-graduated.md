---
owner: main
date: 2026-09-28
---
# Prove-the-setup rule graduated into AGENTS.md § Testing     (date: 2026-09-28)

## Summary
BUG-087, BUG-090 and BUG-091 are one mistake made three times: trusting a behavioural or sandbox run before proving its setup holds what the run claims. The operator approved graduating it ("go all" on the 2026-09-28 closing list). This repo's `AGENTS.md` § Testing gains a paragraph with three requirements:
- the run script asserts every fixture fact the dispatch states before the CLI starts;
- a headless-agy sandbox home carries the real home's agy `settings.json` (no credential), and tool-using runs are checked for the "auto-denied" line;
- a must-fail control for a `context/` change runs outside a clone of this repository.

The three bug rows move to the Graduated table.

## Type
task

## Files changed
- `AGENTS.md` — § Testing: the new paragraph, 5 lines.
- `.agents/docs/bugs/INDEX.md` — BUG-087, BUG-090 and BUG-091 moved to Graduated.

## Agents dispatched
- `bugs`: moved the rows.
- `verifier`: over the tree before the v1.7.0 gate.

## Verification
- `grep -n 'proves its setup before its output counts' AGENTS.md` shows the paragraph.
- The rows moved, with no other row touched (read in the diff).
- The verifier's log reaches this tree.
- The main agent read the diff: 5 lines, below the reviewer threshold.

## Docs
`AGENTS.md` (project instructions) only.

## Requirements
n/a
