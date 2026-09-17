---
owner: changelog
date: 2026-09-17
---
# Requirement R9: four procedures set in stone; map pins them to the binary     (2026-09-17)

## Summary

R9 was added to `REQUIREMENTS.md` with the operator's explicit approval on 2026-09-17, fixing the four procedures that operate the harness — Install, Update, Report back, Release — as the immutable mechanism for how the harness reaches machines and how experience comes back. The map in `README.md` § How it works draws them as a Mermaid diagram (solid lines work today; two dashed lines are gaps: reports from other places are read by hand, and `/retro --harness` applies only the proposals queued on its own machine); AGENTS.md § How the harness works establishes five strict rules that keep the map honest; ADR-031 records the design and the gate that pins it. `TestHarnessMap` in `cmd/sh-agent/main_test.go` fails when the README has no section, when the map names an `sh-agent` verb the binary does not have, or when install, update or report is missing from it. The mechanism ensures that a procedure renamed or removed fails the gate, never silently.

## Type

task

## Files changed

- `REQUIREMENTS.md` — R9 added (lines 37–46), fixing the four procedures with the operator's explicit approval and ADR-031 citation
- `README.md` — § How it works added (lines 5–58): contract statement, Mermaid diagram (16 edges, colour-coded by procedure), prose for all four procedures and the two gaps, and note that `TestHarnessMap` pins the map to the binary
- `AGENTS.md` — § How the harness works added (lines 9–17): five numbered strict rules (two ways in one way back, tags deliver nothing, changing a path is a design question, dashed means a gap, the map is pinned by a test) with enforcement statement; no other section of AGENTS.md changed
- `.agents/docs/architecture/DECISIONS.md` — ADR-031 added (lines 249–257), recording the design and operator approval with exact text, decisions, alternatives rejected, and consequences naming R9 and the two gaps
- `cmd/sh-agent/main_test.go` — `TestHarnessMap` added (lines 987–1014), a small regex pass over the README mermaid block against the binary's `usage` constant, checking for the three required verbs and rejecting bogus ones

## Agents dispatched

changelog, reviewer (`.agents/docs/reviews/2026-09-17-review-77f47f7-map-r9.md`, APPROVE WITH NITS, both findings fixed: ADR-031 citation corrected from "Rules for working here" to "How the harness works" with rule count fixed to five, AGENTS.md rule 5 rewritten to name the interactive page as a private artifact not in the tree and clarified rule 4 on dashed lines), bugs (record for the ADR draft's inventory mistakes — named consumer projects and misdescribed rule count — struck before commit).

## Verification

- Must-fail: `.agents/scratch/r9-mustfail.log` — test binary given a doctored README with `reportx` (non-existent verb) and missing `update`/`report` lines fails with three errors, each caught by the test's three assertions (bogus verb, missing verbs)
- Must-pass: `.agents/scratch/r9-test.log` — `go test -race -count=1 ./...` exit 0, two packages pass (cmd/sh-agent 3.340s, internal/harness 1.784s)
- Main agent verified: gates passed, README mermaid syntax manually checked (16 edges, colour groups partition all exactly once, dashed edges mark the two acknowledged gaps, node labels valid), `usage` constant read (`install`, `update`, `report` all present), Mermaid diagram renders (rendered to `.agents/scratch/map/map.png` with mermaid-cli and checked)

## Docs

- `README.md` — § How it works added: Mermaid map, contract statement, four procedures with prose, note on `TestHarnessMap` pinning the map
- `AGENTS.md` — § How the harness works added: five numbered strict rules, enforcement statement
- `.agents/docs/architecture/DECISIONS.md` — ADR-031 appended: context, design, alternatives rejected, consequences

## Requirements

- R9 (approved by the operator 2026-09-17 on the exact text in `REQUIREMENTS.md` lines 37–46) — four procedures set in stone, drawn in README.md, pinned by test, changed only with operator approval and ADR
- R8 (delivery is `sh-agent` and nothing else) — affirmed; R9 draws only the Install, Update, Report back, Release procedures
- R3 (no hooks, no gates, no CI) — affirmed; R9 adds no gate, only a test pin

## Not exercised

- The map viewed on GitHub itself (rendered locally with mermaid-cli only).
- The interactive page is a private claude.ai artifact outside the repository; it was not regenerated from the final README map text.
- An agent-written first draft of this entry claimed an AGENTS.md release-section change and misdescribed the map's gaps; both corrected by the main agent against the diff before commit (BUG-061's rule).
