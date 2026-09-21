---
owner: main
date: 2026-09-21
---
# Context consolidation, Phase 4 — real home proven, BUG-085 fixed, plan closed     (2026-09-21)

## Summary

Last phase of `.agents/docs/plans/2026-09-21-context-consolidation.md`. On the operator's word the real home on this box was updated to 94a6596 (`sh-agent update`, drift none, binary unshadowed), and the three smoke checks passed on both CLIs from a fresh demo clone; Antigravity's own conversation store held the rules block with the canary line and zero truncation markers, the first real-home Antigravity session without a cut since v1.2.0. BUG-085 is fixed. The plan is closed: every box ticked or moved; the two boxes only the operator can take — the Mac smoke and `/preflight v1.5.0` — are in the handoff's Next Steps.

## Type

feature phase (verification, records, docs)

## What was done

- `sh-agent update` into the real home (17 written, 243 unchanged); `hash -r`, `sh-agent version` = `94a65966d102`, `status` drift none, rendered context 20,205 bytes ending with the canary on both vendors.
- `testing/2026-09-21-verify-real-home-smoke-5.md` — PASS ×3 on both CLIs; Antigravity's skill-load and verifier checks reran with the per-invocation `--dangerously-skip-permissions` in the throwaway clone because the real settings hold no allow list today (settings file untouched, as on 2026-09-16).
- Must-fail mutation of the size test (94a6596): 2,010 bytes appended in place made the test fail naming 215 bytes over; restored by copy, tree clean.
- `bugs`: BUG-085 status fixed in the record and the index. `documentation` (94a6596): README's `rules/` row, this repo's `AGENTS.md` setup commands and guide § 2 name the size test.
- Plan closed; `ROADMAP.md` "Now" cleared; `TODO.md` empty. The sandbox home's copied credential files were deleted after Phase 2.

## Verification

`go test -race -count=1 ./...` exit 0 (bare, main agent; no source changed in this phase, plan D17); the smoke log above.

## Review

No source in this phase; the documentation edits (11 lines) were read by the main agent against the diff.

## Agents dispatched

`bugs` (BUG-085 status), `documentation`.

## Feature totals (ae9512f..this commit)

Context 27,163 → 20,205 bytes, pinned at 22,000; canary last; retro and bugs regrowth stop; three procedures into skills; goldfish gaps closed; ADR-035; BUG-085 fixed, BUG-086 and BUG-087 recorded. Cold read three rounds; reviewers ×9 (cold read 3, code 1, docs 3, goldfish 3); verifier ×2; classification and refute workflows ×2; behavioural runs 3 of 3. Release pending: v1.5.0 MINOR from the operator's turn.
