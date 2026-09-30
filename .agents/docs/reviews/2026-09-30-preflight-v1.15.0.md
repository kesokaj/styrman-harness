---
owner: main
date: 2026-09-30
---
# Preflight — v1.15.0   (2026-09-30)

Range reviewed: v1.14.0..e69c56e

Release request: the operator, 2026-09-30, in their own turn: "I want everything to be completed. All of it." then "go all", after a status list that ended with the release being the last step of "everything"; the previous two releases of the day were cut on the same operator's word. A MINOR: the `verifier` agent's long-suite rule changes (an installed place gains behaviour), plus two fixes to released behaviour (`status`/`search --global` listing proposals by name; `report` redaction rules for nine families). The version follows from `deploy/remote.md` § Releases.

## Findings, worst first

None blocking. Every finding on the range was fixed or rebutted in the task's own records before e69c56e.

| Severity / confidence | Finding | Where | Outcome |
|---|---|---|---|
| high / high (rule, not code) | the nine redaction families were recorded as accepted residue and no decision reversed it | task code review (`reviews/2026-09-30-review-worktree-completion-code.md`) | rebutted: the operator's words of 2026-09-30, "I want everything to be completed. All of it.", after the status list named the residue as the one optional item; recorded in the review's Outcome and in `changes/2026-09-30-1220-everything-completed.md` |
| low / high | two rules without a benign near-miss row | same review, 2 | fixed before the commit |
| low | near shapes of the covered families still pass; the guide reads as full coverage | task security check (`reviews/2026-09-30-security-worktree-completion.md`) | fixed before the commit: six rules widened (rule count unchanged), guide reworded to what the rules reach; the list is closed for instance work (the code comment and the guide say a further shape is a design question) |

## Rows

- **Reviewer**: `reviews/2026-09-30-review-worktree-completion-code.md` (sonnet, REQUEST CHANGES turned by the Outcome: one rebutted with the operator's words, one fixed) over the range's only source commit e69c56e; e83c2d0 is one AGENTS.md line, read by the main agent. Not re-dispatched on opus: the covering report reads the same diff line for line and its one open verdict is the operator's, not the code's.
- **Security**: `reviews/2026-09-30-security-worktree-completion.md` over the same diff: every changed path read (the diff is eight files; `## Not examined` names nothing from the range), 89,887 lines of repository prose probed against the new rules with no false positive, the ReadDir changes read nothing the glob did not; two lows, both fixed, Outcome appended; `govulncheck` the known GO-2026-5024 only.
- **Verification**: `testing/2026-09-30-verify-v1.15.0-gate.md`, PASS at e69c56e (374 pass, 0 fail, 2 skips with reasons; embed, context-size, lint and map tests; scratch install with no drift; govulncheck; three cross-builds with the linux/amd64 binary printing v1.15.0; `report --dry-run`; the installed verifier byte-identical); the task's own log `testing/2026-09-30-verify-completion.md` (PASS, mutations a, b, c, the search tier: rebuild 1.85 s, warm search 36 ms).
- **The project's own blocking condition** (`AGENTS.md` § Release): `testing/2026-09-30-verify-real-home-v1.15.0.md`, PASS — the real home at e69c56e, the three smoke checks on each CLI, and one run of the changed `verifier` on each CLI with a three-minute suite (Claude Code: once, foreground, ten-minute timeout, nothing left behind; Antigravity: the suite recorded, nothing left behind); plus the Claude Code behavioural run `testing/2026-09-30-verifier-long-suites-behavioural.md` (new 0 of 3 background calls, control 2 of 3).
- **Docs freshness** (Step 4): no interface, command or deploy step changed; the guide sentence the range touches (`deploy/local.md:66`) was checked by the security agent and reworded to what the rules reach; no docs-lens dispatch beyond that, recorded here as the main agent's judgement.
- **Trail freshness** (Step 3): pass — the newest `changes/` entry before the release entry (`2026-09-30-1220-everything-completed.md`) describes the range's work; no future `date:`; no modification to `changes/` or `DECISIONS.md` in the range; `ROADMAP.md` and `TODO.md` untouched — both empty by design, no plan active.
- **Requirements**: R4, R8, R9 untouched.
- **Shadowing agents**: none.

## Process notes

- The gate verifier of this release handed back with no live background child, the first gate verifier today to do so; the range's verifier rule is the reason it names.

## What was not checked

- The Mac and the devbox: installed after the release by the operator's `sh-agent update` there; two binaries built, not run.
- Interactive sessions on either CLI.
- A suite longer than ten minutes (the verifier rule's background branch).
- A real `sh-agent report` push and a real `update` from this release: the live paths run at the places after the upload.
- A docs-lens reviewer dispatch over the range (see the docs row: no interface, command or deploy step changed).

Every item above was read against the blocking conditions: a verifier FAIL, BLOCKED or IN PROGRESS log (none: gate, task and real-home logs are PASS); a high security finding (none); a missing owed security report or one naming a release-diff path under Not examined (none); a reviewer REQUEST CHANGES standing (turned in the Outcome by the operator's own words) or a high finding on the code (none); a `REQUIREMENTS.md` violation (none); a project agent shadowing a review agent (none); the project's own condition (met). No item is one a condition names.

VERDICT: SHIP
