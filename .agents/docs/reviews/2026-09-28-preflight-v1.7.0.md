---
owner: main
date: 2026-09-28
---
# Preflight — v1.7.0   (2026-09-28)

Range reviewed: v1.6.0..56621b6

VERDICT: SHIP

No finding blocks: no verifier FAIL, no high security finding, no REQUEST CHANGES from a code review, no REQUIREMENTS violation, and no shadowing project agent (none exists under `.claude/agents` or `.agents/agents`). Every non-blocking finding below carries its decision.

## Pre-gate commits (Step 1)

- `a9717cc`: prove-the-setup rule graduated into `AGENTS.md` § Testing. Verified before the gate: `testing/2026-09-28-verify-v1.7.0-gate.md`, PASS.
- `f916e35`: round-1 fixes.
- `56621b6`: round-2 fixes.

## Findings and decisions, worst first

| # | Round (range) | Source | Severity / confidence | Finding | Decision |
|---|---|---|---|---|---|
| 1 | 1–3 | security | medium / high | `testing/skill-matrix.md:13-14` grants broad sandbox permissions for the full matrix (`command(*)`, `write_file(*)`, `Bash(bash:*)`, `Bash(python3:*)`, `Bash(curl:*)`), contradicting the new `AGENTS.md` § Testing rule | **Accepted, deferred.** The problem predates this release, and the matrix is not run in it. Filed in the handoff's Next Steps, to be fixed before the next full-matrix run, because that run reads this file. |
| 2 | 1 (v1.6.0..a9717cc) | security | medium / high | the new rule copied the real home's agy `settings.json` (`toolPermission: always-proceed`, real checkouts trusted) into sandbox homes | **Fixed** in `f916e35`, then corrected in `56621b6` (row 5). The copy in the sandbox home was deleted when named. |
| 3 | 1 | security | low / medium | the fallback secret scan (gitleaks absent on this host) missed `sk-ant-`, `refresh_token`, `.credentials.json` and `antigravity-oauth-token` | **Fixed** in `f916e35` and `56621b6` |
| 4 | 1 | code | low / medium | `agents/changelog.md` Output omitted the date-only flag that rule 1 promises | **Fixed** in `f916e35`. Behavioural: new definition 2 of 2, control 2 of 2 (it already flagged), so the file now agrees with itself. |
| 5 | 2 (a9717cc..f916e35) | code, security | low / high | the example allow list included `list_dir(*)` (agy 1.2.8: "invalid allow entry") and `find`/`git`/`sed` labelled read-only | **Fixed** in `56621b6`: the bullet states the principle (fewest grants; a `command(<name>)` grant allows every argument; clone-only trust; a `cli.log` "invalid allow entry" check), not a list. It was the second round on one mechanism, so the operator was asked; answer: "go with your recommendation and ship v1.7.0". |
| 6 | 2 | code, security | low / high | `refresh_token` was missing from the fallback scan | **Fixed** in `56621b6` |
| 7 | 2 | docs | low / medium | `context/AGENTS.md`'s secret-shaped-path list does not name CLI token files as preflight now does | **Accepted, deferred.** A context change needs its own full cycle (size test, behavioural runs). Filed in the handoff. |
| 8 | 3 (f916e35..56621b6) | security | low / medium | the rule's two checks (auto-denied, invalid allow entry) show only that a run was not blocked; a settings file with `always-proceed` would pass both | **Accepted.** The rule itself forbids copying the real home's file. A check for a `toolPermission` key is a follow-up for the same task as row 1. |
| 9 | 3 | security | low / high that it fires | the fallback scan's pattern names match prose in the trail files that ride this commit | **Judged, file by file.** Every hit is a pattern name in review or testing prose: `reviews/2026-09-28-review-a9717cc-f916e35-{code,docs}.md`, `reviews/2026-09-28-security-{v1.6.0-a9717cc,a9717cc-f916e35,f916e35-56621b6}.md`, `testing/2026-09-28-verify-v1.7.0-gate-2.md`. A strict token-shape scan found 0 credential shapes in those files and 0 in every line added over `v1.6.0..56621b6`; a planted `sk-ant-api03-…` fixture fired, which validates it. |
| 10 | 3 | security | low / medium | BUG-091 § Prevention still says to copy the real home's settings | **Accepted.** Bug records are immutable apart from their status. The rule graduated into `AGENTS.md` § Testing, which now says the opposite and is the operative text; BUG-091's index row points there. |

Round 1 docs lens (`reviews/2026-09-28-review-v1.6.0-a9717cc-docs.md`): APPROVE, no findings. Round 3 has no code review: the diff is 2 lines, and the main agent read it.

## Reviews (Step 2, Step 4)

- **Code:**
  - `reviews/2026-09-28-review-v1.6.0-a9717cc-code.md` (opus): APPROVE WITH NITS;
  - `reviews/2026-09-28-review-a9717cc-f916e35-code.md` (opus): APPROVE WITH NITS;
  - earlier coverage of `17fa89b..e5a49d4`: `reviews/2026-09-28-review-time-sortable-names-code.md`.
- **Security:**
  - `reviews/2026-09-28-security-v1.6.0-a9717cc.md`;
  - `reviews/2026-09-28-security-a9717cc-f916e35.md`;
  - `reviews/2026-09-28-security-f916e35-56621b6.md`.

  All three: FINDINGS, no high; every path read, nothing under Not examined.
- **Docs:**
  - `reviews/2026-09-28-review-v1.6.0-a9717cc-docs.md`: APPROVE;
  - `reviews/2026-09-28-review-a9717cc-f916e35-docs.md`: APPROVE WITH NITS (row 7).

## Verification (Step 2)

- `testing/2026-09-28-verify-v1.7.0-gate-3.md`: PASS over `v1.6.0..56621b6`.
  - gofmt and vet clean;
  - `go test -race`: 132 top-level and 251 including subtests pass, 0 fail;
  - govulncheck clean;
  - scratch install: drift none, GANNET on both vendors;
  - 27 skills and 7 agents;
  - embed guard clean.
  - Its hand-back summary called `56621b6` a `context/AGENTS.md` change. It changed `AGENTS.md` (project) and `skills/preflight/SKILL.md`; the log's commands are unaffected.
- Behavioural and E2E evidence for the release's features:
  - `testing/2026-09-28-verify-time-sortable-names-phase-2.md`;
  - `-e2e1.md`;
  - `testing/2026-09-28-verify-v1.7.0-gate-fixes.md`.

## Trail freshness (Step 3)

- **Newest change entry:** by the new rule, `changes/2026-09-28-1549-v1.7.0-gate-fixes-2.md`. It describes the last work in the range. By plain name order, `2026-09-28-research-okf.md` still sorts last: today is the switch-over day, which the rule covers.
- **No future `date:`:** checked.
- **Append-only:** `git log --diff-filter=M v1.6.0..HEAD -- changes DECISIONS.md` lists only `c65a8fd`, which appended ADR-037 (8 added, 0 deleted).
- **ROADMAP and TODO:** touched in range; both are empty.

## What was not checked

- No behavioural run of the `56621b6` bullet itself. It is procedure text, and its evidence is the round-1 probe, whose allow list the new bullet would reject.
- The full skill matrix was not run (row 1).
- Behavioural runs of the changelog agent, the handoff skill and the context rule ran on Claude Code only. Antigravity got the sandbox start and E2E 1.
- The real home, the Mac and the devbox after this tag: the operator's smoke checks.
