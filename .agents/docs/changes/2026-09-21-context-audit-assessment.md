---
owner: main
date: 2026-09-21
---
# Context-audit prompt assessed; Antigravity truncation found (BUG-085)     (2026-09-21)

## Summary

The operator handed over an "Agent context audit" prompt and asked whether it would help and whether `rules/` should be created. The answer is a proposal file, not a change: `.agents/docs/plans/2026-09-21-context-audit-assessment.md`. Measuring the harness against the prompt's own scorecard found that the Antigravity CLI loads only the first 24,000 bytes of the deployed 27,114-byte global context, in every release since v1.2.0 — recorded as BUG-085, open. No source, skill, agent or context file changed.

## Type

docs (assessment, bug record)

## What was done

- Measured every always-on surface, the skill and agent descriptions, and the growth of `context/AGENTS.md` since ADR-009 (7,349 → 23,194 chars at 67 → 74 lines).
- Fetched the Antigravity rules page and the Claude Code memory and skills pages; confirmed the 12,000-char documented cap, the 200-line target, the `paths:` rule trigger (on read), the 1,536-char skill-description cut.
- Asked both CLIs, from the real home with no tools, for the last heading of the global context: Antigravity answered `## Environment & Stack` and reported a 3,114-byte truncation; Claude Code answered `## Installed`. The CLI's own `<truncated N bytes>` marker was found in its conversation store for 2026-09-18 and 2026-09-21, both at byte 24,000.
- Three independent refuters (Workflow `wf_cc624239-e08`: measurements, fit, rules lenses) checked the draft; their corrections are applied and listed in the plan file.
- `bugs` dispatched: BUG-085 record under `.agents/docs/bugs/` and its `INDEX.md` row, status open.

## Recommendation on record

Adopt the prompt's method in part — defects first, pruning classes with counts, one-row-per-move layout, a goldfish on the three always-on files — inside the consolidation Feature queued since 2026-09-18. Strike its enforcement layer (R3, ADR-007), its Cursor/Codex rows, its goldfish over every skill and its adherence experiment. `rules/` later, as an output of the pruning, after a two-file probe shows the cap is per file. Four open questions for the operator are in the plan file.

## Agents dispatched

`bugs` (BUG-085). The refute pass ran as a Workflow of three general subagents, read-only.

## Not done

No fix. `ROADMAP.md` "Now" and `TODO.md` stay empty: nothing is active until the operator decides. The prompt file `agent-context-audit-prompt.md` is untracked at the repo root pending Open question 1.
