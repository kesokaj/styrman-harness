---
owner: changelog
date: 2026-09-17
---
# `sh-agent report` — Phase 2 of delivery rework     (2026-09-17)

## Summary

Implemented `sh-agent report [--feedback-repo DIR] [--dry-run]`, which collects redacted feedback (feedback.jsonl, retros.jsonl, state.json, proposals/*.md) from ~/.config/sh-agent, emits a metrics record (report.json), and sends all to a feedback repository. Runs before FindRepo so machines installed from unpacked bundles can report. Redaction is structural (whole JSON value, whole line, or whole BEGIN–END block) per D9, never needle-span. Feedback repository is cloned on first use through ~/.config/sh-agent/feedback; --feedback-repo names an explicit checkout (never cloned into, rejected if it is a harness checkout). Stages by name only, reads git diff --cached --stat, commits and pushes as separate calls. First push into empty repository sets upstream; later runs push plainly. Dry-run prints the plan and writes nothing.

## Type

task

## Files changed

- cmd/sh-agent/main.go — `report` case before FindRepo; `reorder` takesValue map includes `--feedback-repo`; usage text updated
- cmd/sh-agent/main_test.go — `TestReportWithoutCheckout` (no harness checkout, dry-run validates behavior)
- internal/harness/report.go (new) — `Collect`, `BuildRecord`, `ResolveFeedback`, `SendReport`, and redaction helpers
- internal/harness/report_test.go (new) — 8 top-level tests covering redaction, collection, record shape, git flow, feedback resolution

## Agents dispatched

verifier (PASS, .agents/docs/testing/2026-09-17-delivery-rework-phase-2.md), bugs (BUG-051 git env isolation, BUG-052 identical-commit testing). Reviewer and security deferred to Phase 4's Done row per plan.

## Verification

Verifier: 118 pass, 0 fail, 1 pre-existing skip; govulncheck clean (go1.27.1, note: AGENTS.md's two-stdlib-advisory claim is stale); go test -cover shows 69.5% (cmd) and 83.6% (internal); dry-run smoke test against a scratch home cloned nothing, wrote nothing, as required. Every function and case in report.go and the report verb mapped to covering test; two bugs found in test isolation (BUG-051, BUG-052) and fixed.

## Docs

No docs changed — Phase 3 writes README, deploy/local.md, deploy/remote.md, and skils/retro/SKILL.md to describe report and the four places. AGENTS.md note about two stdlib advisories flagged for Phase 3 (go1.27.1 renders them absent).

## Requirements

n/a — R8 satisfied in Phase 1; no interface, command flag, or behavior visible to a user changed beyond the new `sh-agent report` verb itself.
