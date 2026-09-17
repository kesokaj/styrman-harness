---
owner: main
date: 2026-09-17
---
# Handoff — 2026-09-17 — delivery rework Done; R9 sets the four procedures in stone

Session ID: 4e712c10-b55b-4126-ae29-579c7d94ac09
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

The session resumed after the machine went down mid-Definition-of-Done for the delivery-rework plan (`sh-agent report`, the loop back). It finished that plan through two more security rounds, then the operator asked for a diagram of how the project works and for Install, Update, Report back and Release to be "set in stone" — now `REQUIREMENTS.md` R9.

## Completed

- `77f47f7` — delivery rework Done. Security round 2 (plan D13): the per-value JSON walk struck for whole-line redaction; push names the feedback URL, not origin; refusal no longer prints the origin; two shapes widened. Round 3 (D14): after three rounds on the same mechanisms the operator set the guarantees — redaction is a safety net (the guarantee is reading what was added under `~/.config/sh-agent/` before each report), the push trusts the operator's own git config — plus label-scoped key-block close, quoted `Authorization`, `decoded` reading every value, dry run showing unterminated blocks, limits stated in docs. Work-Mac reporting settled (D12). `context/AGENTS.md` gained the numbered-table asking rule. BUG-053–060. Plan fully ticked; ROADMAP "Now" empty.
- `8c3205e` — R9 approved on its exact text; ADR-031; `README.md` § How it works (Mermaid map: solid works, dashed is a gap); project `AGENTS.md` § How the harness works (five strict rules); `TestHarnessMap` in `cmd/sh-agent/main_test.go`. BUG-061.
- Interactive view of the map: private claude.ai artifact https://claude.ai/artifact/YZm6KTqsWwGVv1f8ffHCsz (version 2 aligned with the README map). Outside the repository; the README map is the source (AGENTS.md rule 5).

## Open / In Progress

- Nothing queued in the trail: no unticked plan box, `TODO.md` and `ROADMAP.md` empty, tree clean, pushed.
- Requirements edits awaiting approval: none.

## Key Decisions

- D13, D14 in `.agents/docs/plans/2026-09-17-delivery-rework.md` — every round-2 and round-3 finding and the operator's choice on each ("take the recommendations as written", both times).
- ADR-031 — R9 over AGENTS.md-only rules (an agent may edit AGENTS.md without approval); the README map is the source, the interactive page only a view.
- Rejected in D14: refusing URL rewrites in code (the operator's own config redirects every push on the machine), and chasing multi-line or unlisted credential shapes.

## Corrections

- none logged this session (`~/.config/sh-agent/feedback.jsonl` has no entry from it).

## Agent notes (unverified)

- The file tools converted `\u0041` escapes into `A` at least four times this session (a test fixture, a Go comment, a security report, BUG-057 itself, and this handoff). Write escape-bearing text through a script and check bytes (BUG-057).
- Agent-written ADR and changes drafts invented facts twice (consumer project names, a "deferred" decision, a section that did not change); read them against the diff before committing (BUG-061). BUG-061's author suggests graduating a "no consumer names" check into `agents/changelog.md` — not done.
- Subagents hit turn limits or an API safeguard block three times; a dispatch that says "write the report early, keep to ~N tool calls" finished reliably.
- The operator's "Do you understand what im saying?" was a push for stronger enforcement than AGENTS.md rules, answered by R9 — not logged as a correction.

## Environment State

- Branch: main, clean, up to date with origin at `8c3205e`. Tags: `v1.0.0` only.
- Running processes: none.
- Scratch worth keeping until the next session: `.agents/scratch/r3-mutants/`, `r2-mustfail.log`, `r2-decoded-mustfail.log`, `r9-mustfail.log`, `map/` (Mermaid source, PNG render, drafts), the round-2/3 decision write-ups `2026-09-17-report-round{2,3}-decisions.md`.

## Artifacts to Read

- `REQUIREMENTS.md` R9 — the four procedures, immutable without approval.
- `README.md` § How it works — the map; `AGENTS.md` § How the harness works — the rules around it.
- `.agents/docs/deploy/local.md` § Report back — what `report` sends, its stated limits, its failure table.
- `.agents/docs/changes/2026-09-17-delivery-rework-done.md`, `2026-09-17-harness-map-r9.md` — what landed and what was not exercised.
- `.agents/docs/bugs/INDEX.md` — BUG-053–061 are new.

## Next Steps

1. Operator: at each place, `sh-agent update` (the context and `report` changed), then `sh-agent status`, then the three real-home smoke checks (`.agents/docs/testing/2026-09-03-verify-real-home.md`) — not run this session.
2. Operator: read what was added under `~/.config/sh-agent/` since the last report, `sh-agent report --dry-run`, then `sh-agent report` from this box — the D13/D14 code has not pushed to the real GitHub feedback repository yet.
3. Worth a plan if the operator wants it: close the dashed gaps on the map — something that reads other places' reports from `styrman-harness-feedback` into `/retro --harness`. Under R9 this needs the operator's approval, an ADR and one commit with the map.
4. Optional: `/retro` on this session's bug records (BUG-057, BUG-061 are rules about agent output) — two graduation candidates, the escape byte-check and the changelog agent's "only what the diff says".
