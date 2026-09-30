---
owner: main
date: 2026-09-30
---
# Preflight — v1.16.0   (2026-09-30)

Range reviewed: v1.15.0..bf36817

Release request: the operator, 2026-09-30, in their own turn: "Go all. Then cut a new version again and update docs." A MINOR: installed places gain rule text in three skills (`handoff`, `retro`, `init-project`'s scaffold template); no code change since v1.15.0; the docs pass changed one deploy-guide paragraph.

## Findings, worst first

None blocking. The docs lens over the task found two restatements, both fixed before b452c59, and one left with its reason recorded (`reviews/2026-09-30-review-worktree-retro2-docs.md` § Outcome).

## Rows

- **Reviewer**: the range is 17 lines of prose over five files; below the reviewer threshold in lines and read by the main agent; the docs lens (`reviews/2026-09-30-review-worktree-retro2-docs.md`, REQUEST CHANGES turned by two fixes and one recorded rebuttal) covers the same diff line for line.
- **Security**: no row. No changed path touches auth, crypto, secrets, IAM, input parsing or dependencies.
- **Verification**: `testing/2026-09-30-verify-v1.16.0-gate.md`, PASS at bf36817 (374 pass, 0 fail, 2 skips with reasons; embed, context-size, lint and map tests; scratch install with no drift; govulncheck; three cross-builds with the linux/amd64 binary printing v1.16.0; installed copies byte-identical; no background job at hand-back); the task's own log `testing/2026-09-30-verify-retro2-apply.md`, PASS.
- **The project's own blocking condition** (`AGENTS.md` § Release): `testing/2026-09-30-verify-real-home-v1.16.0.md`, PASS — the real home at bf36817, the context check and one run of each of the three changed skills on each CLI, eight runs, every scenario first attempt, no auto-denied line, `~/.config/sh-agent/` unchanged. The handoff runs show both new rules in effect on both CLIs.
- **Docs freshness** (Step 4): the docs lens above; and the operator's requested docs pass (`documentation`, one paragraph changed, every other living doc checked and left with a reason, in `changes/2026-09-30-1311-retro2-apply-and-docs.md`).
- **Trail freshness** (Step 3): pass — the newest `changes/` entry before the release entry describes the range's work; no future `date:`; no modification to `changes/` or `DECISIONS.md` in the range; `ROADMAP.md` untouched and empty, `TODO.md` touched (its header sentence qualified).
- **Requirements**: untouched.
- **Shadowing agents**: none.

## What was not checked

- The Mac and the devbox: they pick the release up with `sh-agent update`; two binaries built, not run.
- Interactive sessions on either CLI.
- The retro skill's harness-apply sentence in a fixture (see the real-home log: it ran for real in this session instead).
- A real `sh-agent report` push and a real `update` from this release.

Every item above was read against the blocking conditions: a verifier FAIL, BLOCKED or IN PROGRESS log (none); a high security finding (no row owed); a missing owed security report (none owed); a reviewer REQUEST CHANGES standing or a high finding (the docs lens's verdict is turned in its Outcome; no code finding); a `REQUIREMENTS.md` violation (none); a shadowing agent (none); the project's own condition (met for all three changed skills on Antigravity and the smoke on both CLIs). No item is one a condition names.

VERDICT: SHIP
