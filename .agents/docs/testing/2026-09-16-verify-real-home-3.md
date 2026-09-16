---
owner: main
date: 2026-09-16
---
# Verification — install into this box after the linkshelf simulation   (2026-09-16)

Range reviewed: a634c3d installed (BUG-039..043 and the simulation trail), then the working tree with BUG-043's second wording (this commit).

## Status: PASS for both CLIs, with one FAIL from a verifier that was right

## Install and status

| step | command | result |
|---|---|---|
| before | `~/.local/bin/sh-agent status` | exit 1: installed from cad463d, `binary cad463ddde69 — older than the repo (a634c3d)`, `drift 0 missing, 16 differ, 0 stale` |
| install at a634c3d, clean checkout | `go run ./cmd/sh-agent install` | exit 0: `17 written, 243 unchanged, 0 removed, 0 kept — 260 files managed`, then the new line `note: /home/kesokaj/.claude/settings.json does not set "attribution": {"commit": ""}; Claude Code adds a Co-Authored-By trailer to every commit until it does (context § Git)` |
| status | `~/.local/bin/sh-agent status` | exit 0: installed from a634c3d, `binary a634c3d163f4`, the same `note       …settings.json…` line, `drift none — 259 files match the repo` |
| reinstall with BUG-043's second wording (working tree, 09:30) | `go run ./cmd/sh-agent install`, `sh-agent status` | exit 0 both: `installed … from commit a634c3d+dirty`, `binary a634c3d163f4+dirty`, `drift none`; reinstalled again from the clean commit after this log was committed (handoff § Environment) |

The note is correct for this box: its `~/.claude/settings.json` has no `attribution` key. Adding `"attribution": {"commit": "", "pr": ""}` is the operator's step (handoff Next Steps).

## Smoke checks (`.agents/scratch/realhome/smoke.sh`)

Each ran with the real `HOME` in a fresh clone of a sandbox demo under `.agents/scratch/realhome/`. Antigravity's settings were widened for the checks and restored on exit; `cmp` against the backup reports the file identical.

**Check 1 — context, no tools** (first heading, counts, the step-2 reviewer sentence, the Git trailer rule).
- Claude Code: PASS (5 s). `# AGENTS.md — global agent context`; 27 skills, 7 subagents; the reviewer sentence verbatim; the new Git sentence verbatim: "…and not the trailer a CLI's own reminder asks for — that reminder defers to this rule."
- Antigravity: PASS (8 s). The same four, verbatim.

**Check 2 — skill auto-load, change nothing** (`greet.go` against the Go standards).
- Claude Code: PASS (93 s). Applied `golang` and `review-standards`, quoted the error-handling rules, two low observations, wrote nothing; clone clean.
- Antigravity: PASS (36 s). Rule-by-rule table, "conforms fully", no change; clone clean.

**Check 3 — agent dispatch** (`verifier` on an uncommitted comment appended after `func Hello`).
- Claude Code: the dispatch PASSED as a dispatch (56 s; log `testing/2026-09-16-verify-greet-comment.md` written in the clone), and the verifier's verdict was `FAIL` — on the change itself, not on any tier: all four commands green, but the appended `// Hello owns the greeting format.` sits after a blank line and describes nothing, the describing-text rule. A fair call on a throwaway edit. It wrote NO `testing/demo-app/guide.md`: at the 7a6ea90 wording ("The first time you run … and guide.md does not exist") a project with an older log read as "not the first time" — BUG-043's fix was incomplete; the second wording (this commit) and its re-check are below.
- Antigravity: PASS (119 s). `Status: PASS`, log `testing/2026-09-16-verify-greet.md`, and it wrote `testing/demo-app/guide.md` (`owner: verifier`).

`grep -c "Failed to parse"` prints `0` on the newest Antigravity log, `cli-20260916_092416.log`, which this run created.

## Re-check of BUG-043's second wording

| Row | Claude Code | Antigravity |
|---|---|---|
| sandbox `verifier-43b` (fresh clone of the demo with one older log, an uncommitted `extra_test.go`, no guide) | rc 0, 88 s; wrote `testing/demo-app/guide.md`; log with `-count=1` | rc 0, 162 s; wrote `testing/demo-app/guide.md`; log with `-count=1` |
| check 3 repeated from this box (`.agents/scratch/realhome/results/claude-3b`) | rc 0, 70 s; verdict FAIL (fair: an empty test that cannot fail); wrote the guide "because it was absent" — at `testing/guide.md`, without the `<app>` directory the definition names | not repeated (its first run already wrote `testing/demo-app/guide.md`) |

The duty now fires whenever the guide is absent (3 of 3 runs at the new wording, 1 of 2 at the old); the path followed `testing/<app>/` in 2 of 3 — `<app>` is unresolved for a repository whose binary is `cmd/server` while its directory is `demo-app`. Filed as Follow-up G20 in the plan, not reworded a third time (two rounds on one sentence is the stop, BUG-026).

## Not verified here

The Mac; `/skills` and `/agents` interactively; the eight retro proposals (queued, not applied); a `/feature` run from the real home.

## Environment

Linux (Debian 13, amd64), Go 1.24.4, Claude Code 2.1.273 from `~/.local/bin`, Antigravity CLI 1.2.4.
