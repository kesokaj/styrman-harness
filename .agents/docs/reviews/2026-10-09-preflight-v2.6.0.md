---
owner: main
date: 2026-10-09
---
# Preflight — v2.6.0   (2026-10-09)

Range reviewed: v2.5.0..5d8ab96

This is a MINOR release. It adds `sh-agent metrics` (ingest, show, report, export) and its feedback into dream and `/retro`, and fixes BUG-164 and BUG-165. Nothing an existing place must react to.

## Failures found during the gate (all fixed before the tag)

| # | Failure | Severity / confidence | Resolution |
|---|---|---|---|
| 1 | The Antigravity `/retro` run on 56b4aaa recorded nothing about the metrics it read: the record template had no place for them | medium / high | fixed in c478d86 (`## Metrics` section). Re-run on both CLIs at c478d86 and 5d8ab96, all checks passing |
| 2 | `testing/2026-10-09-verify-v2.6.0-final.md` is FAIL: `GOOS=windows go vet ./...` | none (dispatch error) / high | Windows is not a release target, and both failing lines predate v2.5.0. The later PASS log `-final-r2.md` says: "The earlier FAIL came from the dispatch asking for a Windows vet, not from the code." `-final-r3.md` (PASS, reaches HEAD) then ran vet on the three release targets only |
| 3 | Security over 6118921..c478d86: replay de-duplication crossed labels; retro Step 1 had no data guard; `report` usage omitted `--until` | low / medium | fixed in 5d8ab96, each code fix under a test red before it (`reviews/2026-10-09-security-6118921-c478d86.md` § Outcome) |
| 4 | Reviewer over 56b4aaa..c478d86: the committed matrix `retro` check did not require `## Metrics` | low / medium | fixed in 5d8ab96; check seen passing on the new records and failing on three old records and an empty clone |
| 5 | Claude Code login expired mid-gate. It cut two verifier runs short and failed two real-home smoke rows | none (environment) / high | all re-run after `/login`; the failed outputs are kept, and the smoke log says so |

## Step 2: review, security, verification

**Shadowing:** no project agent shadows `reviewer`, `security`, `verifier` or `explorer`; the repo has no `.claude/agents/` or `.agents/agents/` override.

**Reviewer**, all `## Status: COMPLETE`:

| Report | Range | Verdict |
|---|---|---|
| `reviews/2026-10-08-review-2c0f215-166aa39-code-docs.md` | 2c0f215..166aa39 | approve with nits, Outcome written |
| `reviews/2026-10-09-review-166aa39-a967512-code-docs.md` | 166aa39..56b4aaa, rounds 1–8, opus, code and docs | round 8 APPROVE |
| `reviews/2026-10-09-review-56b4aaa-c478d86-code-docs.md` | 56b4aaa..c478d86 | approve with nits, fixed |
| `reviews/2026-10-09-review-worktree-code-docs.md` | BUG-165 fix | approve with nits |

- v2.5.0..2c0f215 is covered by the kit-leftovers reviews (`2026-10-08-review-worktree-5-code-docs.md`).
- c478d86..5d8ab96 is round 1's fixes, all low. Round 2 was skipped under `review` Step 4: every code fix has a test that failed without it, the verifier ran after, and no security finding was medium or higher. The skill and matrix edits there were read by me.

**Security**, all `## Status: COMPLETE`:
- `reviews/2026-10-08-security-2c0f215-166aa39.md`
- `reviews/2026-10-09-security-worktree.md`
- `reviews/2026-10-09-security-v2.5.0-a967512-part{1,2,3}.md`, together every path of v2.5.0..a967512, part 1 rounds 1–5 also covering 6118921
- `reviews/2026-10-09-security-6118921-c478d86.md`, 13 of 13 paths read
- `reviews/2026-10-08-security-worktree-kit-leftovers.md`

No high finding anywhere. No release-diff path is named under `## Not examined`. c478d86..5d8ab96 is that last report's fixes; its Outcome records them.

**Verification:**
1. **A verifier log reaches HEAD and passes:** `testing/2026-10-09-verify-v2.6.0-final-r3.md`, `owner: verifier`, range v2.5.0..5d8ab96, PASS.
   - gofmt clean; vet clean on host, linux/amd64, darwin/arm64 and linux/arm64;
   - `go test -race -timeout 30m -v -count=1 ./...`: 638 pass, 0 fail, 3 skip (the three env-gated measurements);
   - govulncheck: no reachable finding;
   - § 3g canary must-fail control on a copy.

   After 5d8ab96 only dated records changed.
2. **What the release changed was run:**
   - § 3g's commands: the race suite and the canary control in r3. The real-data check from copies is in `-final.md` at c478d86: Claude tool calls 2707 against 2707 unique ids, Antigravity 190 against 190 tool steps, a second ingest adds 0. Since then only `show.go`'s key changed, by adding the label (real data carries no label).
   - The matrix's `retro` row check was changed in this range and run in the sandbox kits.
   - The main agent's tiers, from the testing guide § 2/§ 4 and the plan's sandbox runs:
     - `testing/2026-10-09-verify-v2.6.0-sandbox.md`, PASS, range to 5d8ab96: Claude Code 7 of 7 and Antigravity 8 of 8 at c478d86 and at 5d8ab96.
     - `testing/2026-10-09-verify-real-home-v2.6.0.md`, PASS, range to 5d8ab96: the § 5 smoke, 6 of 6 at 5d8ab96.
3. **No failure is left standing:**
   - `-final.md` FAIL is cleared as row 2 above.
   - `2026-10-08-verify-metrics-sandbox.md` FAIL (the old matrix `retro` check, 0 of 2) is cleared. Its check was changed with a control at 2c0f215 showing the failure predates the feature (`changes/2026-10-09-1011-pre-release-fixes.md`), and it was re-run passing on both CLIs in the sandbox log.
   - The Antigravity kit's credential FAILs (both times AWS's example key, deleted by name, re-scanned clean with a planted-key control) are cleared in the sandbox log.

**Project rule (AGENTS.md § Release)**, for `skills/retro`, the one skill changed in the range:
- the real-home smoke on both CLIs at the release head, 6 of 6;
- Antigravity runs of it at c478d86 and 5d8ab96, both passing.

## Step 3: trail freshness

| Check | Result |
|---|---|
| The newest `changes/` entry (`2026-10-09-1250-v2.6.0-tail-review-fixes.md`) describes the release diff's last work | pass |
| No future `date:` | pass |
| Deleted under `testing/` | `ab-kit/metrics.py`: a kit script replaced by `sh-agent metrics show` (plan D24, Phase 5 entry), not a verification log. pass |
| `changes/` and `DECISIONS.md` append-only | three commits touch DECISIONS.md; each diff is insertions only (19, 8, 8 lines). pass |
| `ROADMAP.md` / `TODO.md` touched since v2.5.0 | yes (9187d89); `TODO.md` holds no items. pass |

## Step 4: docs freshness

The docs lens rode in every reviewer dispatch above. The one docs gap it found (the matrix check) is fixed. No open docs finding.

## Accepted, with reasons

- **Security 6118921..c478d86, finding 2 (low):** retro copies other sessions' verbatim quotes into its record. Rebutted in that record's Outcome: the Findings column has always required the user's words, and the entries are this project's, in this project's trail.
- **`internal/metrics/open_other.go`** (`!unix`) is not compiled by any release target and has no test. It is the one-line `os.Open` fallback for a platform nothing ships to.
- **Pre-existing:** `GOOS=windows go vet` fails on `syscall.Flock` (`internal/memory`) and `Setpgid` (`internal/harness` test). Windows is not a target; this is not a defect of this release.

## What was not checked

- **The Mac:** no darwin run, beyond `GOOS=darwin go vet` and the cross build.
- **A plain-copy place, a release-binary place and a zip place:** not updated from this tag before release. That is each place's own `sh-agent update` after the tag, per `deploy/remote.md`.
- **pass^k:** one sandbox run per CLI per commit, so no repeat-run consistency.
- **Opus on Claude Code** for the sandbox `retro` runs (Sonnet).
- **The nightly dream timer's ingest on a real schedule:** the timer is not installed on this box. `dream --dry-run` and the dream tests cover the wiring.
- **Interactive sessions** on either CLI.

None of these is a blocking condition. The project's Antigravity rule is met by the runs above.

## Bridge check

v2.5.0 was not a bridge release; v2.4.0 was, and v2.5.0 already followed it. No bridge closes with this tag.

VERDICT: SHIP
