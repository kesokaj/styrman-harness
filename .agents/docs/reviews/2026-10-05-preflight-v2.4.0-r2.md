---
owner: main
date: 2026-10-05
---
# Preflight v2.4.0 (r2) — releases without a key

Range reviewed: v2.3.2..6556bbe

Supersedes `reviews/2026-10-05-preflight-v2.4.0.md` (v2.3.2..b549d71), whose release commit `caa37d6` was never tagged: it compiles in a signer whose key was exposed.

## Blocking

None.

## Rows

| Row | Evidence | Result |
|---|---|---|
| Verification | `testing/2026-10-05-1601-verify-v2.4.0-gate-r5.md` (verifier, v2.3.2..6556bbe): gofmt, vet, `go test -race` all six packages ok, 42 targeted tests, scratch install drift none, list | PASS |
| Main-agent tiers | `testing/2026-10-05-verify-v2.4.0-gate-rows-r2.md` (v2.3.2..6556bbe): `context skills agents commands rules` byte-identical b549d71..6556bbe, so the behavioural, Antigravity and real-home rows of `testing/2026-10-05-verify-v2.3.3-gate-rows.md` carry; host-pin mutation | PASS |
| Earlier failing logs | `2026-10-05-1317-verify-v2.4.0-gate.md` (PASS at 41af344), r3 (PASS 48ae45c), r4 (PASS 18b9018); none FAIL or BLOCKED | clear |
| Code review | v2.3.2..b549d71: `reviews/2026-10-05-review-v2.3.2-0c0846c-code.md` and the first gate; b549d71..18b9018: `review-signer-rotation-code.md`, `-r2` (findings fixed in 48ae45c, 18b9018); 18b9018..4ee867a: `review-release-unsigned-adr053.md` APPROVE WITH NITS (3 lows fixed in 6556bbe); 4ee867a..6556bbe (5 Go files, ~60 lines of code and tests) read by the main agent | pass |
| Security | `security-v2.4.0-unsigned-adr053.md` (18b9018..4ee867a): 3 medium; `-r2` (4ee867a..6556bbe): host pin and bridge fixed, 3 low. Earlier rounds over the signer rotation are superseded by ADR-053. No high finding stands | pass |
| Docs | `review-18b9018-4ee867a-docs.md`: nine findings, fixed in 6556bbe | pass |
| Trail | append-only `changes/` and `DECISIONS.md` (0 removed lines v2.3.2..HEAD); no testing log deleted; no future dates; ROADMAP/TODO untouched (report-only) | pass |
| Shadowing agents | none under `.claude/agents` or `.agents/agents` | pass |
| REQUIREMENTS | R9 changed with the operator's approval in their own turn, ADR-053, one commit (4ee867a) | pass |

## Findings not fixed, and why

| Finding | Severity | Decision |
|---|---|---|
| GitHub Actions enabled (security r1 #1) | medium | left: operator — put to them with the release (ADR-053a) |
| `gh` token not rotated (security r1 #3, BUG-160) | medium | left: operator — they skipped it on 2026-10-05; put to them again |
| First-install commands in README and remote.md name no host (security r2 #1) | low | fixed right after the tag in a docs-only commit; no binary is affected |
| Bridge delete line exits 0 with an empty `RETIRED` (security r2 #2) | low | left: the line runs once, for v2.4.0, with an explicit path checked by the line before it |
| Paste-cache copy of the old key (security r2 #3) | low | deleted during the bridge step of this release |
| Security r2 did not read `DECISIONS.md` lines 202-604 | — | rebutted: the file is append-only (0 removed lines v2.3.2..HEAD) and the range's only change to it, the ADR-053a append, was read; earlier rounds read the rest |
| Version | — | MINOR (v2.4.0), the operator's: new behaviour, nothing an existing deployment must react to beyond updating while v2.4.0 is newest |

## What was not checked

- `govulncheck` is not installed on this box.
- No new Antigravity or real-home run: the embedded harness is byte-identical to b549d71, whose rows passed.
- `linux/arm64` and `darwin/arm64` binaries are built, not run; the v2.3.2-accepts-v2.4.0 check runs on linux/amd64.

VERDICT: SHIP
