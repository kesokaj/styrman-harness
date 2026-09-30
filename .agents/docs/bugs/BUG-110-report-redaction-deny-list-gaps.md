---
owner: bugs
date: 2026-09-30
id: BUG-110
status: fixed
area: internal/harness/report.go (redaction)
---
# BUG-110: report redaction deny-list missed common credential shapes, then missed escaped/underscored tokens

## Symptom
`internal/harness/report.go:39-58` redacted only a fixed vendor-prefix and
keyword deny-list; JWTs, Stripe (`sk_live_`/`rk_live_`/`whsec_`), `npm_`,
`hf_`, `pypi-`, legacy `sk-` keys, Slack webhook URLs, and a bare
`https://<token>@host` userinfo token with no colon all passed through
unredacted. Once prefixes were added, the keyword needles' leading `\b`
missed a token immediately after a literal `\n` or an underscore, because a
line that failed to decode keeps its escapes literally and `\` and `_` are
both word characters. Found across three security review rounds
(2026-09-30, findings 4 and round-2/round-3 follow-ups).

## Repro
A report containing e.g. an `hf_...` Hugging Face token, or a token preceded
by a literal `\n` in undecoded JSON text, was not flagged by the redaction
pass and would ship unredacted to the feedback repository.

## Root cause
The deny-list was built from the shapes known at the time each prior bug
(BUG-053, BUG-056, BUG-059, BUG-060) was fixed, not from a systematic
enumeration of vendor token formats; a list of instances never converges.
Separately, the keyword regexps' `\b` boundary assumption ("a token starts
after a non-word character") is false for text that still carries escape
sequences or underscore-joined identifiers.

## Fix
`internal/harness/report.go` `prefixNeedles` and the keyword `regexp`
needles: added JWT, Stripe, npm, Hugging Face, PyPI (two forms), Slack
webhook (services and workflows), OpenAI service/admin/legacy key, and a
percent-escaped/padded userinfo-token class; removed the leading `\b` from
needles matched against raw (possibly escaped) text. Rows added across three
rounds in `TestReportRedactShapes`.

## Prevention
A deny-list of credential-shape instances never converges — when a review
names more shapes for the third time on the same mechanism, put a per-kind
generalized rule to the operator instead of writing a fourth batch of
literal prefixes (done here: Outcome item 4, `delivery.md` § report notes
redaction is a safety net, not a guarantee). A needle matched against raw
text never carries a leading `\b`, because escapes and underscores are word
characters. Every added needle gets a benign near-miss case in the same
test.

## Related
BUG-053, BUG-054, BUG-056, BUG-059, BUG-060 (same mechanism, report
redaction, three prior rounds); this is the fourth-plus round on the same
mechanism — recommend the "put a per-kind rule to the operator on a third
occurrence" line move into `context/AGENTS.md` § Hard rules (the "A claim
that something works" bullet family) or `skills/review-standards/SKILL.md`
§ 2, replacing no existing sentence (new class of recurrence, actionable and
generalizable beyond this file).
