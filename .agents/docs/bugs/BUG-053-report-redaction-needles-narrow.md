---
owner: bugs
date: 2026-09-17
id: BUG-053
status: fixed
area: internal/harness (report.go redaction)
---
# BUG-053: `sh-agent report`'s redaction needles missed most real credential shapes, and `--dry-run` could not reveal the miss

## Symptom
Security review (6af54f7..015cf2d) found 14 of 15 probed credential shapes sent unredacted: `GITHUB_TOKEN=…`, `DB_PASSWORD=…`, `Authorization: Bearer …`, `ghp_…`, `sk-ant-…`, `xoxb-…`, `GOCSPX-…`, a `postgres://user:pass@` URL, etc. Only lowercase `token=` was caught. The docs (`deploy/local.md:55`, `deploy/remote.md:100`) told the operator to check `--dry-run` first as a safeguard, but `--dry-run` only prints redactions that *fired* — it cannot show a secret the matcher missed.

## Repro
Probe in `.agents/scratch/secprobe-6af54f7/secprobe_test.go` (gitignored, fake values only): each shape run through `redactLines(…, true)` as a `user_quote`; `probe-run1.log` records 0 redactions for 14/15 shapes.

## Root cause
The shipped needle list (`internal/harness/report.go:29-38` at 015cf2d) was the seven shapes named in the plan's first draft (`-----BEGIN`, `PRIVATE KEY`, `AKIA`, `AIza`, `ya29.`, `token=`, `password=`), matched with case-sensitive `strings.Contains`. It was never checked against the common uppercase/env-var forms of the same secrets, and the `--dry-run` safeguard the docs recommend was never tested against an actual miss — only against a shape it already caught.

## Fix
`prefixNeedles` (vendor prefixes never shorter than the vendor's own) and `keywordNeedles` (case-insensitive secret-assignment regex, `bearer`, URL credentials) in `internal/harness/report.go`, commit 012b4ca. `TestReportRedactShapes` covers all 15 probed shapes plus `api_key=`/`API-KEY:` in JSONL and text, with benign controls; must-fail seen against 015cf2d's `report.go`. `deploy/local.md` § Report back and `deploy/remote.md` § Operator steps now say the dry run shows only what was caught and to read the originals.

## Prevention
Before shipping a redaction or secret-detection feature, build adversarial fixtures of the realistic shapes it must catch (case variants, env-var prefixes, common vendor formats) — not only the shapes named in the design doc — and run them as must-fail tests against the shipped matcher. A safeguard the docs recommend to catch a failure mode must itself be tested against that failure mode (does `--dry-run` show a miss, not just a hit) before the docs claim it.

## Related
BUG-054, BUG-055, BUG-056 (same review, same range); three related findings from one dispatch — the shared rule ("adversarial fixtures before shipping a detection feature; a documented safeguard is tested against the failure it claims to catch") is a candidate for `context/AGENTS.md` Hard rules or a `security` agent checklist item, since it is actionable (a specific fixture-coverage check) and generalizable beyond this file.
