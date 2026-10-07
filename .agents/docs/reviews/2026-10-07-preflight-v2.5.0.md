---
owner: main
date: 2026-10-07
---
# Preflight — v2.5.0

Range reviewed: v2.4.1..8d00cc5

The operator asked "Solve and prove the last issues then cut a new version" and "go and cut the new version when done" (2026-10-07).
- **Version: MINOR.** The release adds behaviour (process rules, the Antigravity print-mode fix, test kits) and breaks nothing an existing place must react to.
- **The repo tracks no version file;** the tag is the version.
- **Base:** the last tag, v2.4.1.

## Failures and blocking findings

None. Every Step 2 item is met, the security report is COMPLETE with low findings only, every reviewer report is approved, and the Antigravity and real-home tiers pass at the release tree.

## Step 1 — Release diff
- 17 commits, clean tree (only dated logs and this report uncommitted).
- **Commit `8d00cc5`:** one line in the agy-rel kit (the `review` row gets the `command(*)` exception). Under the reviewer threshold; the main agent read the diff, with no change to credentials or isolation.
- **Pre-gate commit `8ffa525`:** the BUG-163 fix and the review Step 4 exemption, brought through the Definition of Done. Reviewer `reviews/2026-10-07-review-worktree-4-code-docs.md`, approved with nits, fixed or left in its Outcome. The changes entry `changes/2026-10-07-1536-bug-163-kit-leftovers.md` is written by the main agent.
- **Docs:** no dispatch; the diff was under 30 lines of docs.

## Step 2 — Review
**Reviewer:** no new dispatch; every commit is covered by a COMPLETE report.

| Commits | Covered by |
|---|---|
| `0eb8320` | `reviews/2026-10-05-review-kit-credential-failclosed.md` (218e301..0eb8320) |
| `29cd652` | `…-failclosed-r2.md` (0eb8320..29cd652) |
| `12f75af` | `reviews/2026-10-07-review-worktree-code-docs.md` (approved, round 2) |
| `95b7519` | `…-worktree-2-code-docs.md` (approved, round 2) |
| `3293e25` | `…-worktree-3-code-docs.md` (approved with nits, round 2) |
| `8ffa525` | `…-worktree-4-code-docs.md` (approved with nits) |

Uncovered and read by the main agent, no finding:
- `218e301` (docs wording: guide and `remote.md`)
- `76ef900` (2 lines in preflight, the bridge check)

The handoff commits are records.

**Security:** `reviews/2026-10-07-security-v2.5.0-release.md`, range v2.4.1..3293e25 plus round 2 over `8ffa525`'s fixes. Status COMPLETE, low only, all 68 release-diff paths read, `## Not examined` names none. The earlier range reviews (`…-security-worktree-ab-kit.md`, `…-security-worktree-agy-seed-kits.md`) closed their mediums.

**Verification:**
1. **The verifier log reaches HEAD and passes:** `testing/2026-10-07-verify-v2.5.0-release-r2.md`, `owner: verifier`, range v2.4.1..8d00cc5, PASS (`Previous run:` the 8ffa525 log, also PASS). It covers gofmt, vet, `go test -race` (559/0/3), the size pins, `TestEmbedHoldsOnlyTrackedFiles`, `TestHarnessMap`, the scratch install and `bash -n` of the three kits.
2. **What the release changed was run:**
   - **The commands:** AGENTS.md § Setup commands is unchanged in the range. The guide changed (§ 4 sentences) and added no command. The kits are run by the main agent's logs.
   - **The tiers named by the range's plans** (lean-process 1–3):

| Tier | Log | Result |
|---|---|---|
| Behavioural A/B measurement | `changes/2026-10-07-0945-lean-process.md`, `…-1129-lean-process-2.md` | measurement recorded |
| Seeded `security` check | `testing/2026-10-07-verify-seed-security.md` | high 6/6 in 6/6 runs, max 6/6 in 1/1 |
| Antigravity release rows (AGENTS.md § Release) | `testing/2026-10-07-verify-agy-release-rows.md` | every changed skill and agent has a PASS; `review` passes at the release tree (`review-rel`, 8d00cc5) |
| Real-home smoke, guide § 5 | `testing/2026-10-07-verify-real-home-v2.5.0.md` | 6/6 at 3293e25, and 6/6 at 8ffa525 (the agy agent row clean after the quota reset). 8d00cc5 changes only `.agents/docs/testing/`, which is not installed, so the installed harness is the same. |

3. **No failure is left standing.** Failing rows in the agy log are each cleared by a later PASS of the same row on the final text: `documentation` by `documentation-rel`, `bug` by `bug-fix4`, `scaffold`, `handoff`, `preflight` and `retro` by their `command(*)` re-runs. `review` is cleared too: the `review-rel-1` attempt failed the preamble (the real config changed during a concurrent real-home smoke) and `review-rel-2` was auto-denied, then `review-rel` PASSed at 8d00cc5 with both clean lines.

## Step 3 — Trail freshness
- **The newest `changes/` entry** describes the last commit's work. Pass.
- **No future `date:`.** Pass.
- **No verification log deleted** in the range (`git log --diff-filter=D … testing/` is empty). Pass.
- **`changes/` and `DECISIONS.md` append-only.** One finding: `changes/2026-10-07-1129-lean-process-2.md` was amended before it was first committed (never modified after a commit; `git log --diff-filter=M` shows no modified past entry). Pass.
- **`ROADMAP.md`:** "Now" was cleared when each plan closed; `TODO.md` is empty. Report-only.

## Step 4 — Docs freshness
The docs lens ran inside each worktree review (`--code --docs`). Advisory gaps left:
- `remote.md:175` does not name the reviews that describe this box's credential broker (security, low; owner: main; handoff).
- `kit_leftovers` is pasted in three kits instead of living in `kit-preamble.sh` (low; owner: main; handoff).

## What was not checked
- The Mac.
- `govulncheck` (not installed). `go.mod` and `go.sum` are unchanged in the range.
- Antigravity's quota ran out once during the gate (16:50 UTC). Every affected run was repeated after the reset, and passed.
- An interactive session on either CLI.

## Decisions on non-blocking findings (Step 6)
- **Security, low (all left):** `kit_leftovers` matches by cwd only (the limit is stated in BUG-163). `remote.md:175` lacks the broker note. The planted injection line is the fixture. The gh-token handoff lines are records. These are named in the release entry and the handoff.
- **Code, low (left):** `kit_leftovers` is pasted in three kits instead of living in `kit-preamble.sh` (handoff).
- **`security` at effort `high`:** accepted (ADR-054a, the seeded check); a harder fixture is in the handoff.

VERDICT: SHIP
