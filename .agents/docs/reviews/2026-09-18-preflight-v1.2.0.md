---
owner: main
date: 2026-09-18
---
# Release gate — v1.2.0     (2026-09-18)

Range reviewed: v1.1.0..a1d2ee8

## VERDICT: SHIP

No blocking finding stands. The code review returned REQUEST CHANGES over bba384d with two medium findings; both were fixed before the tag in a1d2ee8, the verification row was re-taken at that commit, and the review's `## Outcome` records each fix. No high-severity finding, no `REQUIREMENTS.md` violation, and no project agent shadows a harness review agent (`.agents/agents/` and `.claude/agents/` hold none, checked before dispatching).

## What this release is

A small release on top of v1.1.0, and the reason it exists is content rather than code: an installed place gains the six prevention rules graduated into `context/AGENTS.md` and two additions to `skills/preflight/SKILL.md`. It also fixes what `update` prints about the old kept feedback clone. MINOR by this repo's rule — an installed place GAINS a changed context rule and a changed skill — not a PATCH, and nothing here requires an existing deployment to react.

## Rows

| Row | Verdict | Artifact |
|---|---|---|
| Verification, v1.1.0..bba384d | PASS — 102 top-level + 86 subtests, 0 fail, 0 skip; `govulncheck` clean; the new context and preflight text confirmed in BOTH deployed copies | `testing/2026-09-18-verify-preflight-v1.2.0.md` |
| Verification, after the gate's fixes | PASS — same counts reproduced at a1d2ee8; all 11 `TestRemoveKeptClone` subtests pass with zero skips, the stale-tracking-ref one genuinely reproducing | `testing/2026-09-18-verify-preflight-v1.2.0-fixes.md` |
| Code review | REQUEST CHANGES → 2 medium, 2 low, all fixed; `## Outcome` written | `reviews/2026-09-18-review-preflight-v1.2.0-code.md` |
| Security | FINDINGS — 1 low, fixed by the same restructure; all 19 release-diff paths read end to end, none deferred | `reviews/2026-09-18-security-preflight-v1.2.0.md` |
| Docs lens | Not dispatched for this range; see "What was not checked" | — |

## Findings, worst first

All five fixed in a1d2ee8; none carried into the release.

1. **[medium/high] the release guide's copy-paste block contradicted its own prose.** The block still ended `gh release create … <all five assets>` while the text five lines above said to create with `--notes-file` alone and upload one asset per call. `/preflight` item 5 executes the guide as written, so v1.2.0 would have used the fragile form — under the very fix meant to prevent it, added by yesterday's retro. The rule had been written BESIDE the commands instead of INTO them. The block now creates the release, loops one `gh release upload … --clobber` per asset with a retry note, and lists the assets afterwards.
2. **[medium/high] the kept-clone message's settling command was narrower than the count it explained.** The count is `rev-list --count --all --not --remotes`; the suggested command checked `HEAD` only. Measured in `.agents/scratch/fetchhead-check/`: with one commit on a second local branch the count is `1`, the old command prints nothing, and `local.md` — changed in the same range — said an empty list means the clone holds nothing unique and can be removed by hand. That was an instruction to hand-delete the only copy of that commit. The command now reads `git log --oneline --all --not --remotes FETCH_HEAD`, verified to print exactly that commit, and `local.md` says what it covers and why.
3. **[low/high] the message interpolated the clone's path into a command to paste** (security). A home path with a space or a shell metacharacter would have ridden into copy-paste-ready text; the program never executes it, so this was suggested-command injection rather than live. Fixed by removing both interpolations — `leave()` already names the directory — and the test asserts `git -C` never appears.
4. **[low] the new subtest keyed its skip on `@{upstream}`** rather than on the quantity the guard reads. It now skips only when `rev-list --count --all --not --remotes` is `0` or errors, so the shape is confirmed rather than inferred from a symptom; the re-verification shows it passing, not skipping.
5. **[low] preflight's non-blocking-findings rule sat after the items that tag and push**, while its deadline is "before the tag". Moved above item 1 and reworded.

## What was not checked

- **No docs-lens dispatch for this range.** The docs changed here are `README.md` (one block, itself the fix for BUG-075 and run as written), the two deploy guides and the preflight skill; the code review read all four against the code, which is where finding 1 came from. A `/review --docs` pass would have added little over that, and this is recorded rather than quietly skipped.
- **The Mac.** Nothing in this release was exercised on darwin: `darwin/arm64` is built and published, never run. Every Mac leg of the v1.1.0 plan is still owed and is now owed against v1.2.0 instead.
- **`sh-agent update` from v1.1.0 to v1.2.0 against the real release** — it cannot be tested before v1.2.0 exists. The path is covered by the unit tests and by the v1.1.0 rehearsal; the first real cross-release update will be the Mac's, or this box's next.
- **The two-CLI sandbox start**, because that sandbox home no longer exists here; the verifier substituted a grep of both deployed copies, which proves delivery rather than a clean CLI start.
- **`report` from a release-binary install**, and the swapped-signature edge case against a real release: both owed on the Mac.
- The security dispatch ran no tests and no mutations (the verifier's tier); the main agent ran the mutations for the changed message, two well-formed ones, each failing two subtests.

## Version

v1.2.0 — MINOR. An installed place gains changed context rules and a changed skill; nothing an existing deployment must react to. The tag is the only version tracking.
