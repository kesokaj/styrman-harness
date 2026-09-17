---
owner: changelog
date: 2026-09-17
---
# Phase 3: Delivery docs rewritten around the four places   (2026-09-17)

## Summary

Phase 3 of the delivery-rework plan concluded with docs rewritten to remove all references to downstream mirrors and per-project copies, refocusing on the four places the operator uses this harness and the two ways to get it — checkout via `sh-agent update`, or unpacked bundle via `sh-agent install`. The feedback repository `styrman-harness-feedback` created as the delivery loop's return channel. Reviewer checkpoint passed with two low findings, both fixed in the range.

## Type

task

## Files changed

- `.agents/docs/deploy/remote.md` — retired § Downstream mirrors and § Releases-as-delivery; replaced with § The four places; removed other-project trust-root check; preserved versioning rule and signing recipe; timer paragraph now directs to `sh-agent update` at each place
- `.agents/docs/deploy/local.md` — added § Update (for checkout), § From an unpacked bundle, § Report back (command, flags, redaction semantics, null-field rules, failure table)
- `README.md` — added `report` row to command table; added paragraph on two ways in and loop back; added feedback checkout to installed-where table
- `AGENTS.md` — Key-files row updated to drop "mirrors" reference; Release section points to ADR-030 and describes tags as markers, not delivery channel; plain-copy rule rewritten for unpacked bundle; added delivery spec and `report.go` to Key-files; govulncheck notes corrected for go1.27.1
- `skills/retro/SKILL.md:132` — bullet now directs every place to pick up changes with `sh-agent update` rather than naming downstream mirrors
- `cmd/sh-agent/main_test.go` — comments in `TestDryRunNeverBuilds`, `TestPlainCopyNeverRebuilds` and `plainCopy` reworded from "downstream mirror" / "mirror" to "unpacked bundle" / "plain copy"; assertion message on line 313 updated
- `.agents/docs/specs/sh-agent/delivery.md` — table rows reworded from "argolis-devbox" and "styrman-crew" to "A Linux project host" and "A second Linux project host" per R8

## Agents dispatched

reviewer (docs lens), no others — Phase 4 (e2e and Definition of Done) remains

## Verification

Verification at Phase 3 checkpoint per the plan:
- `git grep -nE 'image/harness|deploy/harness|Downstream mirrors|downstream mirrors' -- README.md AGENTS.md context skills agents .agents/docs/deploy ':!skills/mermaid-syntax'` returned 2 hits before (remote.md:47, skills/retro/SKILL.md:132) and 0 after — must-fail observed, must-pass re-run clean
- `sh-agent install --home .agents/scratch/phase3-home` wrote 260 files; `status` drift none; new wording present in both `.claude/skills/retro/SKILL.md` and `.gemini/config/skills/retro/SKILL.md`
- NOT run: headless CLI starts — sandbox home no longer exists on this box
- Reviewer checkpoint: APPROVE WITH NITS over commit range 9ad2f66..4a610e9 — one docs cross-check (local.md § Report back matched to code message by message); two low findings (test comment wording; delivery spec table naming), both fixed, Outcome appended, gate re-run clean

## Docs

- `.agents/docs/deploy/remote.md` — § The four places, § signed tags as release markers, removed other-project sections
- `.agents/docs/deploy/local.md` — § Update, § From an unpacked bundle, § Report back
- `README.md` — `report` row, two-ways-in paragraph, feedback checkout entry
- `AGENTS.md` — Key-files update, Release section, plain-copy rule, govulncheck notes
- `skills/retro/SKILL.md` — harness-scope bullet reworded
- `cmd/sh-agent/main_test.go` — test comments and assertion message brought current
- `.agents/docs/specs/sh-agent/delivery.md` — table rows reworded per R8

## Requirements

- R8: no consumer project named in living docs — delivered; review finding on specs table fixed; all mirror language retired from design descriptions
- R3 and R4 affirmed — no service added, no CI required, local-by-default and stdlib-only preserved in the docs

