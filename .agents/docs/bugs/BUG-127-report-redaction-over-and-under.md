---
owner: bugs
date: 2026-10-01
id: BUG-127
status: fixed
area: internal/harness/report.go (redaction)
---
# BUG-127: report redaction hit ordinary prose, then after the first fix kept short real secrets (IN1)

## Symptom
Prose redacted: `git checkout -b feature/x`, `## Secrets:`, `token: field`, `-b is the branch flag`. After the first fix, short real secrets kept: `password=correct`, `DB_PASSWORD="correct horse battery staple"`, `TOKEN_PASSWORD=letmein`, `MYSQL_PWD=letmein`, `--password letmein`.

## Repro
Cases pinned in `TestReportRedactShapes` (both sides).

## Root cause
The "secret assignment" and "cookie header" needles in `internal/harness/report.go` matched shape only. The first fix anchored on the value's length but read only the key's first word and no quoted value, trading over-redaction for under-redaction.

## Fix
One named regex `secretRe` (`report.go:60`) with an `accept` function (`report.go:84`): a quoted value always redacts; unquoted needs >= 4 chars when the whole key holds pass/pwd/secret, >= 8 (or >= 6 with a digit) for token/api-key; the pass family is one kind; `--flag value` counts; `-b`/`--cookie` need `name=value`. `token=abcdefg` and `secret = my passphrase` are pinned as accepted under-redaction. Three review rounds, code and security.

## Prevention
A needle that judges a value states its rule in one comment and pins BOTH a benign near-miss and a must-redact shape per clause in the same commit. A fix that narrows a needle gets security review in the same round, since the safe direction is over-redaction.

## Related
Audit IN1 in `plans/2026-09-30-recon-audit-fix.md`; reviews `reviews/2026-09-30-review-4764c5f-worktree-code{,-r2}.md`, `reviews/2026-10-01-review-4764c5f-worktree-code-r3.md`, `reviews/2026-09-30-security-4764c5f-worktree{,-r2}.md`, `reviews/2026-10-01-security-4764c5f-worktree-r3.md`. BUG-110 is the earlier pass on the same list; this is the third pass on `report.go` needles in one feature, so the rule in project `AGENTS.md` § Rules for working here (report.go paragraph) is a graduation check: it exists, and was extended by this record's prevention (rule per clause, security in the same round). Actionable and generalizable; no new home needed.
