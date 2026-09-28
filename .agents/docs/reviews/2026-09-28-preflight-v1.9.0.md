---
owner: main
date: 2026-09-28
---
# Preflight — v1.9.0   (2026-09-28)

Range reviewed: v1.8.0..b5d4cb3

VERDICT: SHIP

Nothing blocks the release:
- no verifier FAIL, BLOCKED or IN PROGRESS log;
- no high security finding, and the security report read all 36 paths in the release diff (`## Not examined`: none);
- no REQUEST CHANGES from a code review;
- no REQUIREMENTS violation: R9 was changed with the operator's "All approved", 2026-09-28 (ADR-039);
- no project agent shadows a review agent (`.claude/agents/` and `.agents/agents/` are absent).

The gate ran over `v1.8.0..52cebd8`. Every low it found was fixed before the tag, in `8da10d2` and `b5d4cb3`. Those two commits got their own review (`reviews/2026-09-28-review-v1.9.0-gate-fixes.md`, APPROVE WITH NITS, the nit fixed) and a fresh verification over the whole range to `b5d4cb3`. The first gate verification (at `52cebd8`) therefore no longer describes what ships. The r2 log does.

## Findings and decisions, worst first

| # | Source | Severity / confidence | Finding | Decision |
|---|---|---|---|---|
| 1 | security (signed-zip r1–r2, gate) | medium / high | A zip with no `SHA256SUMS`/`.sig` beside it installs `UNVERIFIED`, with exit 0 | **Left open, by decision.** R9 allows it and ADR-039 records it as Open. The release handoff asks the operator whether it should need an explicit `--unverified` (a REQUIREMENTS R9 change). This narrows the v1.8.0 medium (gate row 1): a release zip is now verified. |
| 2 | docs (gate) | low / medium | R9, ADR-039 and `delivery.md` say a zip "without them" installs UNVERIFIED. The code does so only when BOTH files are absent, and refuses when just one is present, which is stricter than the text. | **Left for the operator:** the fix is a REQUIREMENTS edit. The handoff proposes the exact clause. The code's behaviour is the intended one and is tested both ways. |
| 3 | code (gate) | low / medium | A refusal never named the zip the release publishes, so GitHub's "Source code" archive looped the operator | Fixed in `8da10d2`: `zipNamesIn`, and the subtest `github_source_archive`. Must-fail seen. |
| 4 | security (gate) | low / high | The zip install's verdict was not recorded | Fixed in `8da10d2`: the manifest records `;verified vX.Y.Z` or `;unverified`, and status shows it. Legacy form tested in `b5d4cb3`. Must-fail seen for both. |
| 5 | security (gate, carried twice) | low / high | The `releaseSigner` comment omitted the `-ldflags -X` override | Fixed in `8da10d2` (comment). |
| 6 | code (gate) | low / medium | Pid 999999 in the sweep test can be live | Fixed in `8da10d2`: `report-0`. |
| 7 | code (gate) | low / medium | The map's preflight node did not name the zip | Fixed in `8da10d2`. |
| 8 | code (gate) | low / low | The context's `npx skills add` line lacked `-y` | Fixed in `8da10d2`. |
| 9 | code (gate-fix review) | low / high | No test covered status reading a pre-verdict zip manifest | Fixed in `b5d4cb3`: `TestStatusLegacyZipManifest`. Must-fail seen. |

## Reviews

- **Code:**
  - `reviews/2026-09-28-review-v1.8.0-52cebd8-code.md`: opus, APPROVE WITH NITS, all findings fixed (`## Outcome`);
  - `reviews/2026-09-28-review-v1.9.0-gate-fixes.md`: over `52cebd8..b5d4cb3`, APPROVE WITH NITS, fixed.
  - Earlier rounds on the signed zip: `review-signed-zip-code{,-r2}.md`, and the cold reads r1–r2.
- **Security:** `reviews/2026-09-28-security-v1.8.0-52cebd8.md` (FINDINGS, no high; `## Outcome`). Earlier: `security-signed-zip{,-r2}.md`. The fix commits change message text, a comment, a manifest value and tests. They add no new input parsing beyond a `strings.Cut` on a value this binary wrote itself, and the gate-fix code review covers them.
- **Docs:** `reviews/2026-09-28-review-v1.8.0-52cebd8-docs.md` (APPROVE WITH NITS; row 2 above).

## Verification

`testing/2026-09-28-verify-v1.9.0-gate-r2.md`, PASS over `v1.8.0..b5d4cb3`:
- gofmt and vet are clean;
- `go test -race`: 170 pass, 0 fail, 0 skip;
- govulncheck is clean;
- scratch install: 258 files, drift none, GANNET at the end of both vendors' contexts;
- `sh-agent list`: 27 skills and 7 agents.

One correction: the r2 log says the range does not touch `context/`, but `8da10d2` edited `context/AGENTS.md` (`-y`). The size test is part of the passing full suite, and the GANNET check ran on the rendered contexts.

The first gate log (`testing/2026-09-28-verify-v1.9.0-gate.md`, PASS at `52cebd8`, context 20,822 of 22,000 bytes) is kept as the record of that head.

Feature evidence: `testing/2026-09-28-verify-signed-zip-e2e.md` (a scratch signer: VERIFIED, then tampered and refused), and `testing/2026-09-28-verify-find-skills-routing.md` (new 2 of 2, control 0 of 2).

## Trail freshness

- **Newest change entry:** `changes/2026-09-28-1940-v1.9.0-gate-fixes.md` describes the fix commits. `2026-09-28-1925-signed-zip.md` and `2026-09-28-1850-v180-followups.md` cover the rest of the range.
- **Dates:** none is in the future.
- **Append-only:** `changes/` and `DECISIONS.md` gained only additions (`--numstat`: 0 deletions). DECISIONS.md adds ADR-039.
- **ROADMAP.md** "Now" was emptied in `52cebd8`. **TODO.md** is empty.

## Version

MINOR, v1.9.0. It adds a release asset and new `--from-zip` and `status` behaviour, and breaks nothing. The manifest's zip value gains a suffix, but the JSON shape is unchanged: a pre-verdict manifest still reads (tested), and an older binary reads only the part before the suffix. No existing deployment has to react.

## What was not checked

- **macOS:** `--from-zip` verification ran on this Linux box only.
- **The real release zip:** verifying it with the real key is release step 3b (this host) and the operator smoke on each machine.
- **The skill matrix and CLI-start sandbox runs:** not run in this release. The context change is one flag in one sentence.
- **Antigravity:** not asked the find-skills routing question.
