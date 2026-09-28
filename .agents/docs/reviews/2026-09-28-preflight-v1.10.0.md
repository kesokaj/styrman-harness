---
owner: main
date: 2026-09-28
---
# Preflight — v1.10.0   (2026-09-28)

Range reviewed: v1.9.0..e33a4e1

VERDICT: SHIP

Nothing blocks the release:
- no verifier FAIL, BLOCKED or IN PROGRESS log at the head;
- no high security finding, and the security report names no release-diff path under `## Not examined`;
- no REQUEST CHANGES from a code review;
- no REQUIREMENTS violation: R9's zip-place clause is back to the v1.9.0 text, with the operator's words quoted in its note (ADR-041);
- no project agent shadows a review agent.

**How the range got here:**
- `--unverified` was added (`b1ecab3`, ADR-040), gated as a planned v2.0.0, and hardened (`d3e493b`).
- The operator then reversed it before any tag: "Don't complicate things. Just accept the zip. This is a workaround for environments that can't use the normal way." (`c30b8da`, ADR-041).
- The v2.0.0 gate artifacts (`reviews/2026-09-28-*-v1.9.0-b1ecab3*`, `testing/2026-09-28-verify-v2.0.0-gate*.md`) describe code that does not ship. They are kept as records only.
- This gate ran over `v1.9.0..c30b8da`, and its two test lows were fixed in `e33a4e1`.

## Findings and decisions, worst first

| # | Source | Severity / confidence | Finding | Decision |
|---|---|---|---|---|
| 1 | security (v1.8.0 → v1.10.0) | medium → accepted | A zip without the full signed pair installs `UNVERIFIED` | **Accepted by the operator (ADR-041):** "Just accept the zip." `UNVERIFIED` in the output and `;unverified` in the manifest disclose it. A zip with the pair is still verified, and refused if the check fails. |
| 2 | security (gate) | low / high | Extra words after `update --from-zip <file>` are ignored silently | **Left, by the operator's decision.** ADR-041 names stray-argument refusals among the hardening the operator stopped. The worst case installs `UNVERIFIED`, the documented behaviour. |
| 3 | security (v1.9.0 gate, carried) | low / high | An older signed release's zip installs `VERIFIED`; there is no downgrade check | **Accepted (ADR-039, ADR-041):** a zip place has no other way to roll back, and the first line names the release. |
| 4 | code (gate) | low / medium | No end-to-end test for a half pair | Fixed in `e33a4e1`: `TestUpdateFromZipHalfPairSaysUnverified`. Must-fail seen. |
| 5 | code (gate) | low / medium | `TestUnpackZipLateRefusalCreatesNoParent` asserted no refusal text | Fixed in `e33a4e1`. Must-fail seen. |

## Reviews

- **Code:** `reviews/2026-09-28-review-v1.9.0-c30b8da-code.md` (opus, APPROVE WITH NITS; `## Outcome`). The `e33a4e1` change adds tests only, and its must-fails are recorded there.
- **Security:** `reviews/2026-09-28-security-v1.9.0-c30b8da.md` (FINDINGS, one low; `## Outcome`).
- **Docs:** `reviews/2026-09-28-review-v1.9.0-c30b8da-docs.md` (APPROVE, no findings).

## Verification

`testing/2026-09-28-verify-v1.10.0-gate-r2.md`, PASS over `v1.9.0..e33a4e1`:
- gofmt and vet are clean;
- `go test -race`: 171 pass, 0 fail, 0 skip;
- govulncheck is clean;
- scratch install: drift none, and GANNET is untruncated on both vendors.

The first gate log (`verify-v1.10.0-gate.md`, PASS at `c30b8da`, 170 pass) says the range added and removed no tests. That is inaccurate: tests were added in `b1ecab3`/`d3e493b` and some were removed in `c30b8da`, netting 170. Its PASS stands.

## Trail freshness

- **Newest change entry:** `changes/2026-09-28-2054-v1.10.0-gate-tests.md`. With `2044-zip-accepts`, `2020-v2.0.0-gate-fixes` and `2010-unverified-flag`, it covers the range.
- **Dates:** none is in the future.
- **Append-only:** `changes/` gained additions only. `DECISIONS.md` gained ADR-040 and ADR-041. Its one changed line is ADR-039's heading, now marked "its one-file refusal superseded by ADR-041"; ADR-040's heading is marked "Superseded by ADR-041".
- **ROADMAP.md** "Now" and **TODO.md** are empty.

## Version

MINOR, v1.10.0. Relative to v1.9.0, a zip with only one of the two files now installs `UNVERIFIED` instead of being refused. Nothing becomes required, and nothing an existing deployment must react to.

## What was not checked

- **macOS:** the zip-place path has not run there.
- **The real release zip:** verifying it with the real key is release step 3b on this host, and the operator smoke check on each machine.
- **The skill matrix and CLI-start sandbox runs:** `context/`, `agents/` and `skills/` did not change in this range.
