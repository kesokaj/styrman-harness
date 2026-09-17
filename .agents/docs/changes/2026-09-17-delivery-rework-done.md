---
owner: changelog
date: 2026-09-17
---
# Definition of Done: sh-agent report fully verified on three security rounds    (2026-09-17)

## Summary

The delivery-rework plan completed Definition of Done on 2026-09-17. Three security review rounds and three code-review rounds examined the `sh-agent report` implementation and its docs. Round 2 (D11) struck the per-value JSON walk in favour of whole-line redaction with a JSON marker in JSONL files, whole-document redaction in `state.json`, push to the feedback URL rather than `origin`, and widened credential needle patterns. Round 3 (D14) narrowed the written claims in `deploy/local.md` about what redaction guarantees, the operator chose the model of "redaction is a safety net" with the guarantee being the operator reading before each report, and eight low findings were within the accepted limits. All findings fixed or rebutted with Outcomes written. Two verification passes confirmed the changes; no regressions. The `context/AGENTS.md` file gained the numbered-table asking rule installed into both vendors' copies.

## Type

task

## Files changed

- `internal/harness/report.go` — D11 and D14: per-value JSON walk removed, whole-line JSONL redaction with JSON marker object, whole-document `state.json` redaction, push to feedback URL with `git push -u <url> HEAD:main` instead of origin, refusal checks origin is the feedback URL and refuses `GIT_*` redirects and linked paths, error message no longer echoes the checkout's origin (which can carry a token), credential needle patterns widened (case-insensitive `secret` assignment, bearer headers with optional quote before colon, vendor prefixes as written, URL credentials with empty user), a key block found before JSON walk closes only at a footer naming the header's label on the header's own line with substring match, decoded-text check reads every value on a line with numbers kept as written
- `internal/harness/report_test.go` — D14: five must-fail test mutants added to `.agents/scratch/r3-mutants/`, fixtures and test cases expanded to cover the five D14 fixes (quoted Authorization key, decoded-text numbers, label-scoped block close, unterminated block reporting, document-level redaction of decoded text), mid-run edit: needle-behind-JSON-escape fixture changed from literal character to real unicode escape (backslash-u0041)
- `.agents/docs/deploy/local.md` § Report back — D14: narrowed written claims to state what the code actually does (block close on header's own line with label substring match, not any line; limits for decoded-text check named: BOM, NaN, Infinity, control characters, truncated records, nesting limits, shadowed keys); added fixture discovery that over-redaction fails safe in the case of a flattened key with dashes stripped
- `.agents/docs/specs/sh-agent/delivery.md` — D11: spec restated that dry-run shows only what was caught; D14: no code change, no doc change (prior findings on the spec remain valid; the defect is the written claim in `local.md`, not the spec)
- `context/AGENTS.md` — added Operator paragraph with the numbered-table asking rule: "Operator: …… When one has more than one option, give me a short list — one line of why per option — and name the recommendation. Never a long preamble, never options without a recommendation." — installed into both vendors' copies (verified by grep of the phrase in both `.claude/AGENTS.md` and `.gemini/config/rules/AGENTS.md` in `.agents/scratch/r2-home`)
- `.agents/docs/bugs/BUG-053` through `BUG-060` — recorded during D11 and D14 review rounds

## Agents dispatched

reviewer (three reports: 6af54f7-015cf2d code, 015cf2d-012b4ca code, 012b4ca-worktree code-docs, 012b4ca-worktree-final), security (three rounds: 6af54f7-015cf2d D11, 015cf2d-012b4ca D13, 012b4ca-worktree D14 with Outcome appended, 012b4ca-worktree-final scoped probe), verifier (two logs: 2026-09-17-report-round2-fixes.md, 2026-09-17-report-round3-fixes.md; first round-2 run stopped by a Claude Code API safeguard and left no log; later run hit its turn limit and was resumed), documentation (four dispatches), bugs (BUG-053 through BUG-060; BUG-053 through BUG-056 recorded in D11, BUG-057 through BUG-060 in D14), changelog

## Verification

Verification at Phase 4 Definition of Done per the plan:

- Verifier round 2 (`.agents/docs/testing/2026-09-17-report-round2-fixes.md`): PASS — gofmt, vet, race test suite clean; govulncheck "No vulnerabilities"; install/status against `.agents/scratch/r2-home` 260 files managed, drift none, both vendor copies contain the new Operator asking-format rule; `report --dry-run` writes nothing; must-fail log `.agents/scratch/r2-mustfail.log` shows 15 FAIL lines, confirming the new test cases exercise the round-2 fixes

- Verifier round 3 (`.agents/docs/testing/2026-09-17-report-round3-fixes.md`): PASS — gofmt, vet, race test suite (69 pass, 0 fail, 1 pre-existing skip) clean; govulncheck "No vulnerabilities"; install/status against `.agents/scratch/r3-home` 260 files managed, drift none; dry-run no-write check (262 files before and after); must-fail mutants in `.agents/scratch/r3-mutants/` (2b-label.log, 3-authz.log, 4-decoded.log, 6-document.log, 7-dryrun.log) each fail their corresponding test, confirming each fix area is exercised

- Reviewer: three code-review reports (6af54f7-015cf2d, 015cf2d-012b4ca, 012b4ca-worktree-code-docs, final 012b4ca-worktree-final) — APPROVE over each range with findings fixed in Outcome sections; one docs-review report (9ad2f66-4a610e9 Phase-3 checkpoint) — APPROVE WITH NITS, two low findings both fixed

- Security: three rounds with Outcomes appended (6af54f7-015cf2d D11, 015cf2d-012b4ca D13, 012b4ca-worktree D14), one scoped final probe (012b4ca-worktree-final with blobs pinned, probes in `.agents/scratch/secprobe-final/`, verdict: FINDINGS (low only), all five inside D14's accepted limits)

## Docs

- `.agents/docs/deploy/local.md` § Report back — narrowed claims on block close and decoded-text limits
- `.agents/docs/specs/sh-agent/delivery.md` — dry-run restated to show only what was caught
- `context/AGENTS.md` — Operator paragraph with numbered-table asking rule added, installed to both vendors
- `.agents/docs/testing/` — two verification logs (round 2, round 3)
- `.agents/docs/reviews/` — six review reports (three code, one docs, three security rounds, one final scoped probe)
- `.agents/docs/bugs/` — eight bug records (BUG-053 through BUG-060)

## Requirements

- R8 (delivery is `sh-agent` and nothing else) — delivered; `sh-agent report` implementation complete, all three security rounds green within D14's accepted limits
- R3 and R4 affirmed — no service added, no CI required, stdlib-only preserved

## Notes

Not part of this entry and deferred to the next commit: `REQUIREMENTS.md` R9, `architecture/DECISIONS.md` ADR-031, `README.md` § How it works with the Mermaid diagram, and `TestHarnessMap` in `cmd/sh-agent/main_test.go`.

Not exercised: a real push to the GitHub feedback repository with these changes, macOS verification, the three real-home smoke checks (`.agents/docs/testing/2026-09-03-verify-real-home.md`).
