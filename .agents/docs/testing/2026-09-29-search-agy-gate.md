---
owner: main
date: 2026-09-29
---
# Antigravity gate runs — every skill and agent the search feature changed   (2026-09-29)

Plan: `plans/2026-09-29-search-release-readiness.md`, Phase 3. Rule: project `AGENTS.md` § Release (no tag without one Antigravity run of each changed skill and agent).

## Status: PASS — 10 of 10 scenarios completed; 2 of them needed the per-run skip flag (a vendor limit, below)

## Setup

agy 1.2.13, headless (`agy -p … --mode accept-edits`). `HOME` was a sandbox home under `~/sandbox/search/gate/`, installed from this tree with `sh-agent install --vendor antigravity` (f41e5c0 for the re-runs). It held a copied login and a `settings.json` written per run: the run's clone as the only trusted workspace, and command grants by name. `sh-agent` was a wrapper pinning `--home` to the sandbox home (BUG-107). Each run used a fresh clone outside every checkout (BUG-106) of the `tasklog` sample project: a Go CLI with passing tests and a trail of five bug records, plans, reviews and handoffs. The login and the settings were deleted by name after the batch. Kit: `.agents/scratch/gate-kit/` (`run.sh`, `check.sh`, `drive.sh`; gitignored, ADR-010).

## Results

| Changed file | Scenario | Attempt that passed | Check | Search output in the conversations |
|---|---|---|---|---|
| `skills/review-standards` | apply the review criteria to `cmd/tasklog/main.go`, chat only | 1 | the answer names `review-standards` | no |
| `skills/research` | `/research`, light mode, a storage question | 1 | ends with the continue question | yes |
| `skills/handoff` | `/handoff` (re-run after the resume step changed its text) | 4 (first of the re-run) | a fourth handoff file | yes |
| `skills/feature` | `/feature`, stop after the plan and the cold read | 1 | plan with boxes, `expect:`, `## Decisions` | yes |
| `skills/bug` | `/bug` on a seeded regression (`TrimLeft` for `TrimSpace`) | 3 | tests pass, BUG-006 written, the regression gone | yes |
| `skills/retro`, `agents/dreamer` | `/retro --dry-run`, which dispatches the `dreamer` | 1 | `retros/*-dreamer.md`, `Status: COMPLETE` | yes |
| `agents/dreamer` | direct dispatch | 1 | the same artifact check | yes |
| `agents/security` | direct dispatch on a one-commit diff, final wording | 1 | `reviews/*-security-gate.md`, `Status: COMPLETE` | yes |
| `skills/init-project` | the update path on an older scaffold (no `retros/`) | 4, with `--dangerously-skip-permissions` | `retros/.gitkeep` exists | no |
| `skills/dream` | `/dream --dry-run` on a seeded store | 4, with `--dangerously-skip-permissions` | verdict table naming the memories | yes |

`agents/reviewer`, `explorer`, `documentation`, `bugs` and `changelog` completed on Antigravity in `testing/2026-09-29-search-index-behavioural.md`; their definitions have not changed since. The context is covered by the real-home smoke, `testing/2026-09-29-verify-real-home-search.md`.

## Vendor limit found: command substitution is auto-denied

agy 1.2.13 auto-denies, in headless mode, any command containing `$(…)`, whatever `permissions.allow` holds. A probe with grants for `ls`, `wc`, `echo` and `cat` gave:

| Command | Result |
|---|---|
| `ls \| wc -l` | ran |
| `n=5; echo $n` | ran |
| `echo count=$(ls \| wc -l)` | auto-denied |

Adding `command(n)`, `command(n=*)` or `command(bash)` changed nothing. The scripts of `/init-project` and `/dream` use `$(…)`, so their first three attempts each ended with "no output produced … auto-denied". Those two skills passed with the per-run flag `--dangerously-skip-permissions`, inside the sandbox home and a throwaway clone. `testing/2026-09-21-verify-real-home-smoke-5.md` used the same flag for the same kind of limit. From the real home the limit does not bite: its settings say `toolPermission: always-proceed`.

Two commands were added to the grants by name after a denial, as the rule asks: `[` (a shell test in the skills' scripts) and, in the efficiency kit, `which`, `rg`, `tee`, `true`, `chmod`, `ps`, `uname` and `false`.

## Print mode ending early

The first two `bug` attempts ended with "dispatched the `verifier` subagent … the tree remains frozen": `agy -p` returned while the subagent was still running, so no bug record existed yet. The third attempt completed. This is the vendor behaviour the context's headless rule describes; the direct dispatches here told the agent not to end its turn before the subagent returned, and all of them completed on the first attempt.

## What was not exercised

- Skills the feature did not change.
- The Mac and the devbox.
- agy's behaviour when a grant list meets `$(…)` in a future version; the probe is three commands and can be re-run.
