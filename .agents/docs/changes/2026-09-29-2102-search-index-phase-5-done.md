---
owner: main
date: 2026-09-29
---
# Search index, Phase 5 and Done: efficiency runs, final verification   (2026-09-29-2102)

Plan: `plans/2026-09-29-search-index.md` (closed).

- **Efficiency runs** (`testing/2026-09-29-search-efficiency.md`): harness v1.12.1 against this branch, on two sample apps (three sessions each) and a larger-project task.
  - Claude Code: 28 of 28 sessions delivered in both arms.
  - Cost: no measurable difference. The median lane cost is 17.78 against 18.78 USD, with overlapping ranges at n = 2 per app per arm.
  - Who searched: the main agent ran `sh-agent search` in 0 of 28 sessions; subagents made all 33 calls (`security` 22, `reviewer` 11). Main-agent adoption is a follow-up for the operator (D36).
  - Antigravity: the app lanes are not comparable, because print mode ended before dispatched subagents returned (D37). The larger-project pair answered 5 of 5 in both arms.
- **Kit hygiene:** agy grants were added by name when a scenario was auto-denied (`which`, `rg`, `tee`, `true`, `chmod`, `ps`, `uname`, `false`, the app binaries). The copied agy tokens and run settings were deleted by name after the batches. The real `~/.config/sh-agent/` gained no file today.
- **Final verification:** `testing/2026-09-29-verify-search-index-final-2.md` PASS at 4c9a327, covering the `LC_ALL=C` fix made after the first final log. The first final log's "215 pass" counts top-level tests; its raw output holds 358 passes with subtests, and HEAD 359.
- **Plan closed**, `ROADMAP.md` "Now" emptied. The real-home smoke checks are the operator's, in the handoff.

Agents dispatched: `verifier`.
