---
owner: main
date: 2026-09-30
---
# Redaction: one rule per kind for the shapes left after the review     (2026-09-30)

## Summary
The operator chose the recommended answers to the two choices the project review left open: the remaining redaction shapes are covered by one rule per kind, and the two forked skills keep their upstream names (no change). `internal/harness/report.go`: the `slack webhook` needle `(?i)hooks\.slack` replaces the two Slack path prefixes and catches any case, a port, a trailing dot, an encoded path and the gov domain; the `url token` needle `://[^/\s:@]{16,}:?@` catches a sixteen-character or longer userinfo token with an empty or no password, escapes and padding included; the `sk-None-` prefix is added. The rule's over-redaction — a long username, or a pathless URL followed by an unrelated `@` on the same line — is accepted in the safe direction and pinned in `TestReportRedactShapes` beside `ASIAN`. Ten rows added to that test, two over-redaction pins, two benign near-misses. `.agents/docs/deploy/local.md:66` describes the two rules (`documentation` wrote the paragraph; the main agent aligned two phrases with the final rules, four words each).

## Type
bugfix (BUG-110's remainder)

## Files changed
- internal/harness/report.go — two rules generalised, one prefix added, comments
- internal/harness/report_test.go — rows, over-redaction pins, near-misses
- .agents/docs/deploy/local.md — the redaction paragraph

## Agents dispatched
security (`reviews/2026-09-30-security-worktree-redaction.md`: three lows, all taken — the over-redaction pinned, the two shorter rules adopted, the stale comment fixed), verifier (`testing/2026-09-30-verify-redaction-rules.md`: PASS, mutation protocol on the three rules), documentation. No reviewer: 36 lines over 3 files, the diff read by the main agent.

## Verification
- `testing/2026-09-30-verify-redaction-rules.md` — PASS at the first form of the rules; the final form re-ran `gofmt`, `go vet`, `go test -race ./...` (4 packages ok) bare, and the mutation of both shorter rules fails `TestReportRedactShapes` on ten lines.
- A rule slip, recorded here rather than hidden: `documentation` was dispatched in parallel with the verifier, and its edit to `deploy/local.md` landed while the verifier ran. The verifier named it. Only a deploy guide moved, no source, so the log describes the tree under test; the next dispatch keeps the tree frozen for every agent.

## Docs
`deploy/local.md:66`. The delivery spec's parenthetical states no number and needed no change.

## Requirements
n/a
