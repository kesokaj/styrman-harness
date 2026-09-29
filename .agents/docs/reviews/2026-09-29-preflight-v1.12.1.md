---
owner: main
date: 2026-09-29
---
# Preflight — v1.12.1   (2026-09-29)

Range reviewed: v1.12.0..c5853b3

Base v1.12.0 (the last tag); head c5853b3, the commit that carries the docs-lens fix on top of 62ea57b, the head the gate rows read (the fix re-verified in round 2). Three source paths changed (`git diff v1.12.0..c5853b3 --stat -- . ':!.agents/docs'`): `internal/harness/harness_test.go` (a test), `skills/init-project/SKILL.md` (one sentence), `AGENTS.md` (one Release sentence). A PATCH: a fix to released behaviour (`/init-project` refuses symlinked writes at every step), nothing new an installed place gains.

## VERDICT: SHIP

Nothing blocks: the verifier row is PASS, there is no security row to owe (no path touches auth, crypto, secrets, IAM, input parsing or dependencies), the reviewer threshold is not reached, no `REQUIREMENTS.md` violation, no project agent shadows a harness one. The repository's own added condition (`AGENTS.md` § Release: the Antigravity runs of changed skills and agents before a tag) is met by `testing/2026-09-29-verify-real-home-v1.12.0.md`: `/init-project`'s changed text was run on Claude Code (the linked-store row) and the range changes no other skill or agent; every skill and agent v1.12.0 added has now run on both CLIs.

## Findings
| Severity / confidence | Finding | Disposition |
|---|---|---|
| medium / high (docs lens) | `skills/init-project/SKILL.md` Step 1 claimed every write is symlink-guarded while Steps 3 and 4's copy loops carried no guard | fixed before the tag in c5853b3: Step 3 skips a symlinked root file, Step 4 refuses a symlinked `.agents` (if/else) and skips a symlinked file or directory; exercised on a linked `.agents` in scratch; re-verified in `testing/2026-09-29-verify-v1.12.1-gate-r2.md` PASS |
| low / low (docs lens) | a file count in the real-home log read as a copy-paste | clarified in the log (that fixture seeds four memories plus the index) |

## Rows
- **Verification:** `testing/2026-09-29-verify-v1.12.1-gate.md` PASS over 62ea57b (suite, size test 21,952 bytes, govulncheck clean, install 266 files drift none, GANNET on both vendors, the pre-tag checks), then `-gate-r2.md` PASS over the fix; `testing/2026-09-29-verify-post-release-v1.12.0.md` PASS over the same tree before the commit (suite 171 pass, the flaky test alone, install drift none); `testing/2026-09-29-verify-real-home-v1.12.0.md` PASS on both CLIs.
- **Review:** below threshold (ten lines, three files): the main agent read the diff — the test's own-path check with its mutation run (BUG-104), the guard sentence, the Release sentence — recorded in `changes/2026-09-29-1609-post-release-checks.md`. No dispatch.
- **Security:** no row owed; said here rather than skipped: the range's test change, guard sentence and rule text touch none of the named surfaces, and the guard narrows what `/init-project` writes.
- **Trail freshness:** pass — the newest `changes/` entry describes the range; no future `date:`; `changes/` and `DECISIONS.md` gained only additions; `ROADMAP.md` and `TODO.md` untouched since v1.12.0 because both are empty (report-only).
- **Docs freshness:** `reviews/2026-09-29-review-v1.12.0-62ea57b-docs.md` APPROVE WITH NITS — the two findings above, one fixed now (small and testable, so the verification row was re-run), one clarified; Outcome on the report. The README map unchanged; `skills/preflight/SKILL.md` already allows a project's added blocking condition, so the new Release rule needs no companion edit there.

## What was not checked
- The Mac and the devbox: the operator's `sh-agent update` there and their smoke.
- `/distill` dispatching the `dreamer`; the `init-memory` and `init-memory-absent` rows on Antigravity; `/init-project`'s update path on Antigravity (its changed sentence ran on Claude Code only, on the linked-store row).

## Release
- Version v1.12.1; notes `.agents/scratch/release-notes-v1.12.1.md`; procedure `deploy/remote.md` § Releases, each command its own call.
