---
owner: main
date: 2026-09-29
---
# Preflight — v1.12.0   (2026-09-29)

Range reviewed: v1.11.0..74ac6a9

Base v1.11.0 (the last tag); head 74ac6a9, the commit that resolves the gate's findings on top of d0e62f3 (the head the gate agents read, plus the working-tree fixes they re-read in their later rounds). Twenty source paths changed since v1.11.0 (`git diff v1.11.0..74ac6a9 --stat -- . ':!.agents/docs'`); no Go source.

## VERDICT: SHIP

Nothing blocks. Every blocking condition of the gate is clear: no verifier FAIL, BLOCKED or IN PROGRESS; no high-severity security finding; the security report has its `## Not examined` section and names no release-diff path there; the release review ended APPROVE WITH NITS after two rounds of fixes; no `REQUIREMENTS.md` violation (both clauses it gained were approved by the operator, ADR-042 and ADR-043); no project agent shadows a harness review agent.

## Findings, worst first — each fixed before the tag or accepted by the operator

| Severity / confidence | Finding | Disposition |
|---|---|---|
| medium / high (review r1) | a standalone `/dream` appended to `retros.jsonl`, resetting `sh-agent status`'s 14-day retro reminder | fixed now: no index line (ADR-044); re-verified in `-gate-r2.md` |
| medium / high (review r1) | the dreamer's `Since` from the newest `retros/` record, which its own artifact would be | fixed now: `retros/*-retro*.md` at all four sites; re-verified |
| medium / high (security) | the `.gitignore` append in Step 5b ran before the trust check and through a committed symlink | fixed now: symlink test before any write, append after the trust check, Step 6 alike; closed in `-security-…-r2.md` |
| medium / high (security) | the Antigravity resume read and the Antigravity memory write did not check `.agents` | fixed now: "neither it nor `.agents` is a symlink" (the first wording checked only the store, caught in the security round 2); `feedback-log` rule 7 skips a symlinked `.agents`; re-verified in `-gate-r4.md` |
| medium / medium (security) | Claude Code loads whatever `autoMemoryDirectory` names at session start, so a `MEMORY.md` force-added upstream and pulled would be loaded | left: operator — accepted 2026-09-29 ("Accept all three (Recommended)"); inherent to the store inside the checkout (ADR-043), recorded in the spec |
| low (review r2, r3; security r1, r2) | eleven lows: leftover `cat >>` row, a restatement in `feedback-log` and one in `distill`, `~/` expansion in the dreamer's rule, the empty-slug stop's place, ADR-042's heading, the insert branch's JSON string, a matrix row that could not fail, `--prompts --fix` applying proposed-only hunks, the dreamer's draft name unbounded | all fixed now; the two carried-over lows (memory text quoted in committed records; instruction-shaped text in memory files read as data) left: operator — accepted 2026-09-29 |

Every fix is a prose change and was re-verified (`testing/2026-09-29-verify-v1.12.0-gate-r2.md`, `-r3.md`, `-r4.md`, each PASS) and re-read by the reviewer (`reviews/2026-09-29-review-v1.11.0-d0e62f3-code-r2.md` REQUEST CHANGES on lows, `-r3.md` APPROVE WITH NITS) and by security (`reviews/2026-09-29-security-v1.11.0-d0e62f3-r2.md`, its finding 1 fixed after it). The Outcome sections on the three gate reports record each disposition.

## Rows

- **Verification:** `testing/2026-09-29-verify-v1.12.0-gate.md` PASS over d0e62f3 (suite, size test, govulncheck clean, list 28 skills / 8 agents, install 266 files drift none, GANNET on both vendors, the embed test, the vendored digest), then `-gate-r2.md`, `-gate-r3.md`, `-gate-r4.md` PASS over the fixes (21,952 bytes of 22,000 at the last). Task-level logs of the range, all PASS: `verify-prompts-lens-r2.md`, `verify-dream-skill-phase-2-r3.md`, `verify-dreamer-memory-loop-phase-1.md`, `-phase-2.md`, `-final.md`, `-r4.md`; behavioural: `prompts-lens-behavioural.md`, `dream-phase-1-behavioural.md`, `memory-loop-phase-1-behavioural.md`, `memory-loop-phase-2-behavioural.md`.
- **Review (opus):** `reviews/2026-09-29-review-v1.11.0-d0e62f3-code.md` REQUEST CHANGES (two mediums, fixed), `-code-r2.md`, `-code-r3.md` APPROVE WITH NITS; the task-level reviews of the range (`review-worktree-code.md`, `review-dream-code.md`, `review-memory-loop-code.md`) and their Outcomes.
- **Security:** `reviews/2026-09-29-security-v1.11.0-d0e62f3.md` (every one of the 84 paths read, no secret in the range, `## Not examined` present and naming no release path), `-r2.md`; task-level rounds `security-memory-loop.md`, `-r2.md`, `-r3.md` with Outcomes.
- **Trail freshness:** pass — the newest `changes/` entry (`2026-09-29-1517-memory-loop-phase-2.md`) describes the release's work; no future `date:`; `changes/` and `DECISIONS.md` gained only additions since v1.11.0 (ADR-042's heading carries a superseded-in-part mark, the one status edit a past ADR allows, as ADR-040's does); `ROADMAP.md` and `TODO.md` touched since the base, both empty now.
- **Docs freshness:** `reviews/2026-09-29-review-v1.11.0-d0e62f3-docs.md` APPROVE, no findings; the README map unchanged by design (no delivery path moved; the store is project-internal). The gate's own fixes changed three documentation sentences (README's harness-state row, the context's learning loop, `distill`'s cross-reference), each re-read by the review rounds.

## Not from this range
- `TestReinstallBuildsAndInstalls` failed once for the reviewer while the gate verifier's `go test` ran at the same time, and passed alone; a shared scratch path under concurrent suites, not in this diff — handoff Next Steps.

## What was not checked
- Antigravity runs of `/dream`, the resume clause, the capture rule and the `dreamer` (the sandbox home this repo's procedure names does not exist on this box); the operator's real-home smoke after `sh-agent update` covers the first start on each CLI.
- The interactive `/init-project` settings write (refused headless in every probe; this checkout's file was written by hand from the step's text).
- The `init-memory-linked-store` matrix row has no run yet.
- The dreamer's skill-candidate column, 6 of 9 in the probes, reported as a soft column in the Phase 2 log.
- Other `/init-project` writes that predate this range and follow a committed symlink (security round 2, low) — handoff Next Steps.

## Release
- Version: v1.12.0, a MINOR — new behaviour an installed place gains (`/review --prompts`, `/dream`, the memory store, the `dreamer`, two REQUIREMENTS clauses), nothing an existing deployment must react to; a project gains its store the next time `/init-project` runs there.
- Notes: `.agents/scratch/release-notes-v1.12.0.md` (Step 6.1).
- Procedure: `deploy/remote.md` § Releases, each command its own call.
