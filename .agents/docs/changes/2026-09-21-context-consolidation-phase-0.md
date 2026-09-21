---
owner: main
date: 2026-09-21
---
# Context consolidation, Phase 0 — cap recorded, canary last, smoke check can now fail     (2026-09-21)

## Summary

First phase of `.agents/docs/plans/2026-09-21-context-consolidation.md` (BUG-085). `REQUIREMENTS.md` § Known limits records that Antigravity's global rules are one 24,000-byte block (ADR-035). `sh-agent install` now ends the rendered context with a fixed canary line, `Canary: this context ends with the word GANNET.`, and the real-home smoke check asks each CLI for that word, so a truncated file fails a check for the first time. The context is still 27,163 bytes; Phase 1 prunes it under the 22,000-byte budget and pins it with a test.

## Type

feature phase (source + docs)

## What changed

- `internal/harness/render.go` — `Canary` constant; `RenderContext` appends it after the inventory; doc comment says why it is last.
- `internal/harness/harness_test.go` — `TestInventoryAndRender` asserts the rendered context ends with the canary (red first on the inventory line, green after).
- `REQUIREMENTS.md` § Known limits — the block-cap bullet (operator-approved text); `.agents/docs/architecture/DECISIONS.md` — ADR-035 (committed together, 5e08629).
- `README.md` (What gets installed where), `AGENTS.md` § Testing, `.agents/docs/testing/2026-09-03-verify-real-home.md` (context row), `.agents/docs/testing/guide.md` § 5 — the canary question and what a FAIL is.
- Plan: cold-reader check ticked (three rounds, cap reached, operator's "approved" for the one residue, D15); Phase 0 boxes and checkpoint ticked.

## Verification

`testing/2026-09-21-context-consolidation-phase-0.md` — PASS: gofmt, vet, `go test -race -count=1 ./...` 250 pass 0 fail; scratch-home install 260 files, status drift none; rendered file ends with the canary, identical on both vendors, 27,163 bytes (over the cap by design until Phase 1). Not exercised: guide § 4 sandbox start (Phase 1 checkpoint), § 5 real-home smoke (Phase 4), behavioural runs (Phase 2).

## Review

Source diff under 50 lines and two files (render.go +9, harness_test.go +6): read by the main agent, no reviewer dispatched. `/review --docs` over Phase 0 and Phase 1 together runs at the Phase 1 checkpoint, where the context itself changes.

## Agents dispatched

`reviewer` ×3 (cold-reader rounds, `reviews/2026-09-21-coldread-context-consolidation-r1..r3.md`), `changelog` (ADR-035), `verifier` (the log above).
