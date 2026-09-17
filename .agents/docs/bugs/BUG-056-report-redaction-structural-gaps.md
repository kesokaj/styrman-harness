---
owner: bugs
date: 2026-09-17
id: BUG-056
status: fixed
area: internal/harness (report.go redaction — structural gaps)
---
# BUG-056: `sh-agent report`'s redaction had four structural gaps where a needle was present but content was still sent and the run reported success

## Symptom
Security review found four cases where redaction silently failed despite a matching needle being present:
(a) duplicate JSON keys — `json.Valid` passes, the decoded map keeps only the last value, nothing fires on the shadowed (secret-bearing) value, and the raw line is sent unmodified;
(b) `-----END`/`-----BEGIN` concatenated on one line (two PEMs joined without a newline) — the END search starts on the BEGIN line itself, so only that line is removed;
(c) a PEM whose dashes were rewritten (e.g. to en dashes by a chat UI) — the header is removed but no block ever opens, so the body is sent;
(d) in JSONL, a `-----BEGIN` header that is valid JSON on its own (a bare string) — the JSON walk handles it and no block opens, so the body is sent.

## Repro
Each probed against 015cf2d's `report.go`: (a) `{"user_quote":"AKIA…","user_quote":"clean"}` gave 0 redactions, key sent; (b) 2 redactions reported but both key-body lines sent; (c) header line removed but body sent; (d) body sent, JSON walk found nothing.

## Root cause
The redaction pipeline (`internal/harness/report.go:67-80,82-104,195-197` at 015cf2d) treated the JSON walk and the block (BEGIN/END) rule as mutually exclusive and order-dependent, and assumed well-formed, single-key, single-line-boundary input: it never handled duplicate keys (Go's `json.Decode` silently drops all but the last), a block whose footer shares a line with another header, or block markers altered just enough (dash substitution, JSON-string-wrapping) to not match the literal `-----BEGIN`/`-----END` strings. These are exactly the "not the shape I designed for" cases an adversarial or accidental input produces.

## Fix
`internal/harness/report.go`, commit 012b4ca: the key-block rule now runs before the JSON branch; `blockEnd` searches for the footer after the header's offset and keeps a block open when a footer is followed on the same line by another header; `blockOpen` also matches a `BEGIN … PRIVATE KEY` header with rewritten dashes; a JSON line (or `state.json`) the walk leaves still carrying a needle falls through to the whole-line text rule. Deviation from the reviewer's suggested remediation: bare `PRIVATE KEY` deliberately does not open a block (would redact to end of file on ordinary prose naming it), but still removes its own line. All four cases plus a `state.json` duplicate-key case are in `TestReportRedactStructuralGaps`; must-fail seen.

## Prevention
A text-processing redaction pipeline needs fixtures for malformed/adversarial structure, not just malformed *content*: duplicate keys, multi-marker lines, altered delimiters, and a structure that satisfies one parser (JSON) while defeating another (block scanning) in the same pipeline. When two matching strategies (JSON walk, line/block scan) can each individually "succeed" on the same input, test what happens when one finds nothing but the raw text still contains the needle — it must fall through, not be treated as clean.

## Related
BUG-053, BUG-054, BUG-055 (same review, same dispatch). Four findings from one review share the root pattern "shipped without adversarial fixtures for the realistic shapes the feature exists to catch" — recommend this graduate as a standing rule in `context/AGENTS.md` Hard rules or `agents/security.md`'s checklist for any redaction/detection feature; it is actionable (require adversarial fixtures matching named threat shapes before merge) and generalizable well beyond `report.go`.
