---
owner: main
date: 2026-09-29
---
# Handoff — 2026-09-29 — research: an embedded SQLite index for the agent; verdict CONDITIONAL, build nothing

Project: /home/kesokaj/workspace/styrman-harness
Branch: `main`; this handoff's commit carries the research record.

## Context

The operator ran `/research` (full mode) on whether the harness should give the agent a local index: embedded SQLite with FTS5 and sqlite-vec, inside a Git-tracked project, with no daemon and no binary churn in Git. The previous handoff is `2026-09-29-1631-v1.12.1-released.md`.

## Completed

- **Research record** `plans/2026-09-29-research-sqlite-index.md`. Verdict CONDITIONAL: SQLite through `modernc.org/sqlite` is the right engine if an index is ever needed, but none is needed today. Recommendation A: keep agent-driven grep and build nothing. The record carries the adoption rule with measurable triggers, the decision matrix, the Git boundary (a blob-ID-keyed cache under `.agents/scratch/`, reconciled at query time) and the failure modes.
- **Measured here:** a whole-tree `git grep` in 11 ms; FTS5 queries under 1 ms; a contentless FTS5 index of 3.4 MB (trigram 22.9 MB); the pure-Go driver with sqlite-vec v0.1.9 working without cgo and cross-compiling to darwin/arm64 (+8.4 MB binary); a 20,000-vector brute-force KNN in 8–9 ms; the sqlite-vec `ncruces` Go binding broken as documented.
- Probe files deleted from scratch. No source, skill, context or requirement changed.

## Open / In Progress

- **Awaiting the operator:** accept recommendation A (nothing to build), or choose C or D. Choosing either makes it a Feature whose first step is an R4 edit put to the operator as exact text. Requirements edits awaiting approval: none.

## Key Decisions

- Full mode stopped after the recommendation: the research skill's spec and `feature` handoff apply to the recommended option, and it builds nothing. Stated in the record.

## Corrections

- None.

## Next Steps

1. Operator: decide on the research record's recommendation.
2. Nothing else is queued.
