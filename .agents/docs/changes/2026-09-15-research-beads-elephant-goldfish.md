---
owner: changelog
date: 2026-09-15
---
# Research: beads, elephant-goldfish, goldfish against harness     (date: 2026-09-15)

## Summary
Researched three GitHub repositories (gastownhall/beads, vshvedov/elephant-goldfish, lukacf/goldfish) for skills, agents, commands, and rules that could be added to the harness. Generated a 40-group recommendation table covering 73 candidates with priority guidance (11 add now, 15 add next, 6 later, 8 skip). Identified harness drift in the removed missing-repro gate sentence and seven open plan boxes. Analyzed bundling cost, tool conflicts, and fits.

## Type
task

## Files changed
- `.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md` — new research record: 40-group recommendation table, suggested first bundle, drift found in the harness itself, cost, conflicts, corrections from verification, and what was not verified
- `.agents/docs/bugs/BUG-017-workflow-resume-reran-proposers.md` — new bug record from workflow run wf_921298bb-241 hitting model usage limit and resume re-running proposers
- `.agents/docs/bugs/INDEX.md` — new row for BUG-017

## Agents dispatched
Workflow run wf_921298bb-241: 6 readers and 3 proposers completed; 37 verification and critic agents failed on a model usage limit; a resumed run re-started the 3 proposers and was stopped (BUG-017). Workflow run wf_1e892864-c31, a new run over the cached candidates: 49 agents, 0 errors (36 verification lenses, 1 completeness critic, 2 second-round proposers, 9 second-round verification lenses, 1 set judge). bugs (BUG-017); changelog (this entry).

## Verification
Main agent spot-checked harness drift claims with `git show 87a862d -- context/AGENTS.md` (removed gate) and `grep` of open plan boxes and priority-on-vendors statement. Checked repo footprints (license headers, build flags, hook paths, MCP tool counts) against shallow clones in `.agents/scratch/research/`. Research file checked for control bytes: none found. No verifier run (docs-only task). No reviewer run (no source diff); the main agent read the research record itself.

## Docs
no docs changed — research record and recommendations only, no interface, command, or requirement amended

## Requirements
n/a
