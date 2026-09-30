---
owner: main
date: 2026-09-30
---
# Preflight — v1.14.0   (2026-09-30)

Range reviewed: v1.13.1..7d6dca0

Release request: the operator, 2026-09-30, in their own turn: "Cut a new version". A MINOR: installed places gain rule text in four skills (`golang`, `bug`, `review-standards`, `feature`) and one agent (`verifier`), from the retro applied in 7535bc0; the project `AGENTS.md` gained two rules that are not deployed; no Go source changed since v1.13.1. The version was not named; MINOR follows from `deploy/remote.md` § Releases ("anything an installed place gains").

## Findings, worst first

None. No agent reported a finding on this range beyond the two wording nits fixed before 7535bc0 (`reviews/2026-09-30-review-worktree-retro-code.md` § Outcome).

## Rows

- **Reviewer**: the range's only source diff is 7535bc0 (12 lines of prose over six files), covered by `reviews/2026-09-30-review-worktree-retro-code.md` (APPROVE WITH NITS, both fixed before the commit, Outcome appended); the rest of the range is trail. Not re-dispatched: the covering report reads the same diff line for line. Read by the main agent as well.
- **Security**: no row. No changed path touches auth, crypto, secrets, IAM, input parsing or dependencies; the instruction text added was read by the reviewer and the docs lens for contradiction with its neighbours.
- **Verification**: `testing/2026-09-30-verify-v1.14.0-gate.md`, PASS at 7d6dca0 (373 pass, 0 fail, 2 skips with reasons; embed, context-size, lint and map tests; scratch install with no drift; govulncheck with the known GO-2026-5024 only; three cross-builds with the linux/amd64 binary printing v1.14.0; installed copies byte-identical to the tree).
- **The project's own blocking condition** (`AGENTS.md` § Release: the real-home smoke on both CLIs plus one Antigravity run of each changed skill and agent): `testing/2026-09-30-verify-real-home-v1.14.0.md`, PASS — the real home installed at 7d6dca0 from this checkout, then from a fixture clone outside every checkout: the three smoke checks on Claude Code and on Antigravity, and one run of `golang`, `review-standards`, `bug`, `feature` and the `verifier` on each CLI, twelve runs, every scenario on its first attempt, no auto-denied line, the real `~/.config/sh-agent/` unchanged. The verifier change also has `testing/2026-09-30-verifier-foreign-artifact-behavioural.md` (17 runs: new 9 of 9, control 1 of 8 false positive).
- **Docs freshness** (Step 4): `reviews/2026-09-30-review-worktree-retro-docs.md`, APPROVE over the same diff, no findings; every added sentence checked against its twin in `context/AGENTS.md`, `review-standards` § 2 and `preflight` Step 2.
- **Trail freshness** (Step 3): pass — the newest `changes/` entry before the release entry (`2026-09-30-1025-retro-apply-seven-rules.md`) describes the range's work; no future `date:`; `changes/` and `DECISIONS.md` gained no modification in the range (`git log --diff-filter=M` empty); `ROADMAP.md` and `TODO.md` untouched since the base — both empty by design, no plan active; lag reported, not a violation.
- **Requirements**: R4, R8, R9 untouched; no `REQUIREMENTS.md` edit.
- **Shadowing agents**: none under `.claude/agents/` or `.agents/agents/`.

## Process notes

- The gate verifier stopped with background work after its log was complete and was stopped by the main thread once the artifact was on disk (the queued 2026-09-23 proposal's shape, third occurrence today).
- The real home was moved to the gated head before the tag, by `sh-agent update` from this checkout, as the v1.12.0 checks did; the release changes nothing about it afterwards.

## What was not checked

- The Mac and the devbox: installed after the release by the operator's `sh-agent update` there; `darwin/arm64` and `linux/arm64` binaries were built, not run.
- Interactive sessions on either CLI; every run was headless.
- The `bug` skill's class check on a class with more than one site on Antigravity (the fixture has one site; Claude Code produced the class-check line on the same fixture).
- The `review-standards` § 2 additions applied by a reviewer to a real consolidation (none in the fixture); the skill loaded and applied its criteria on both CLIs.
- A real `sh-agent report` push and a real `update` from this release: the live paths run at the places after the upload.

Every item above was read against the blocking conditions: a verifier FAIL, BLOCKED or IN PROGRESS log (none: the gate log and the real-home log are PASS); a high security finding (no row owed, none); a missing owed security report (none owed); a reviewer REQUEST CHANGES or high finding (none); a `REQUIREMENTS.md` violation (none); a project agent shadowing a review agent (none); the project's own condition — the real-home smoke on both CLIs and an Antigravity run per changed skill and agent (all present in `testing/2026-09-30-verify-real-home-v1.14.0.md`). No item is one a condition names.

VERDICT: SHIP
