---
owner: main
date: 2026-09-16
---
# Handoff — 2026-09-16 — linkshelf simulation done, retro proposals awaiting the operator

Session ID: c50e0fd4-6273-4e35-b1ef-7639dff8bb07        (same machine: `claude --resume c50e0fd4-6273-4e35-b1ef-7639dff8bb07`; matches the feedback log's session field)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

The operator asked for a harness simulation: in the sandbox home, have Claude Code and Antigravity each build the same small Go service (linkshelf) from an empty repository, `/init-project` then `/feature` through the Definition of Done, headless; compare both results against the harness's rules and the golang skill; turn every gap into a fix under the bug skill or a plan Follow-up; then `/retro` at harness scope, install into this box, `/handoff`. "Stop only at a gate." The one gate reached is the retro's: which proposals to apply.

## Completed

- **Both builds ran to Done without a human reply** (plan `plans/2026-09-16-linkshelf-simulation.md`, log `testing/2026-09-16-linkshelf-simulation.md`). Claude Code: 75 min, 6 commits, exactly the brief's routes, stdlib only, 52 tests, four project bug records, independent review APPROVE. Antigravity: 26 min, 5 commits, over-built (three DELETE shapes, GET-by-id, flags plus env), one dependency taken on a false security finding, review APPROVE WITH NITS. Result repositories: `/home/kesokaj/sandbox/runs/sim-{claude,agy}-linkshelf`.
- **19 gaps (G1–G19, plus G20 found at the end), 7 fixed under the bug skill, each with a reviewer round, a verifier log, a record and a sandbox re-run at the fixed tree:** BUG-039 attribution trailer (7b8743f: Git rule names the CLI setting; `sh-agent install`/`status` print an advisory note; `TestAttributionNote`), BUG-040 deploy template placeholders (ff83cf3), BUG-041 the security exemption's boundary and `Not examined` never `nothing` (10c2155), BUG-042 feature-skill cap, scope and `supersedes` rows (140650c), BUG-043 verifier guide-on-absence and `-count=1` (7a6ea90, second wording f51461d after the real-home smoke found the first skippable).
- **Follow-ups filed in the plan, not built:** G4, G5, G10–G13, G20.
- **Retro written** (`retros/2026-09-16-retro.md`, `~/.config/sh-agent/retros.jsonl`): eight harness-scope proposals queued in `~/.config/sh-agent/proposals/2026-09-16-retro.md`, none applied.
- **This box installed at f51461d** (`binary f51461dc1d95`, `drift none`), six real-home smoke checks PASS (`testing/2026-09-16-verify-real-home-3.md`); Antigravity's settings widened for the checks and restored identical.
- Memory: `project-linkshelf-simulation.md` and an updated `reference-sandbox-headless.md` in the Claude memory directory.

## Open / In Progress

- **The retro gate:** which of the eight proposals to apply (Next Steps 1). Nothing else is in flight.
- The operator's `~/.claude/settings.json` on this box lacks `"attribution": {"commit": ""}`; `sh-agent status` says so until it is added (Next Steps 2).
- Still open from earlier handoffs: the Mac install; the two mirrors at a81ef4b; option C's filing clause in `/preflight` Step 4; BUG-031's wait rule; the two deferred remediations (adopt's temp-directory copy, the direct `RenderContext` symlink test). BUG-019 to BUG-023 and BUG-031 remain open records.
- Requirements edits awaiting approval: none.

## Key Decisions

- D1–D4 in the plan: one empty repository per CLI with a bare local origin; the brief passed verbatim so the DELETE shape reached the feature skill; a wide allow-list so a Dockerfile was the agent's choice; Phase 3 started on settled gaps while the Claude lane ran, because a sandbox reinstall mid-run would have changed the agents that lane loads.
- BUG-041: define "hand-rolled" (a primitive or protocol implemented here) rather than soften the ban; a dependency-adding finding names the stdlib alternative it rejected.
- BUG-043's `-count=1`: unconditional, by deleting the reviewer-flagged condition — a verifier never takes a cached result.
- BUG-043's guide duty: after two rounds on one sentence, the `<app>` ambiguity (G20) was filed, not reworded a third time (BUG-026).
- The retro's proposals stay queued: applying is the operator's call; headless, the record says so.

## Corrections

- None: `~/.config/sh-agent/feedback.jsonl` has no entry for this repository today; the operator's two mid-session messages were "keep going, notify me when the Claude lane finishes".

## Agent notes (unverified)

- Antigravity read only `skills/golang/SKILL.md`, never the reference files; Claude Code read `golang-io-web.md` and `golang-testing.md`. Whether that is the model or the CLI's skill loading is unmeasured (retro Proposal 6 names the two files in `SKILL.md` either way).
- Claude Code's verifier FAIL on real-home check 3 (a dangling comment after `func Hello`) was a fair describing-text call, not a harness fault.
- The Go 1.24.4 toolchain on this box carries 23 stdlib vulnerabilities per `govulncheck` (both CLIs' security agents reported it); an environment item, not a harness one.
- The Antigravity result repository's own BUG-001 record carries a prevention rule that contradicts the golang skill ("always import vetted UUID libraries"); it is a sandbox artefact, left as evidence.

## Environment State

- Branch: main, clean, pushed at f51461d plus this handoff.
- Running processes: none (all detached drivers finished; `pgrep -f sim-kit` is empty).
- Sandbox home: installed from the working tree at a634c3d+the BUG-043 wording (`rerun3.sh`), `drift none` at the time; its `~/.claude/settings.json` now holds `{"attribution": {"commit": "", "pr": ""}}` (the `attribution-off` row) — delete it to reproduce BUG-039 there. Its Antigravity settings are unchanged.
- This box: installed from f51461d, `drift none`; Antigravity settings restored identical to their backup.
- Pending: nothing uncommitted.

## Artifacts to Read

- `plans/2026-09-16-linkshelf-simulation.md` — every box ticked; Follow-ups G4, G5, G10–G13, G20; decisions D1–D4.
- `testing/2026-09-16-linkshelf-simulation.md` — both CLIs side by side: runs, checklists, transcript audit, gaps, the `security-rand-id` row before and after, all re-runs.
- `testing/2026-09-16-verify-real-home-3.md` — the install into this box and the six smoke checks.
- `bugs/INDEX.md` rows BUG-039 to BUG-043 and their records; BUG-043 has an Update section with a `Diagnosis missed:` line.
- `reviews/2026-09-16-review-linkshelf-claude-code.md` and `…-agy-code.md` — the independent code reads of the two built services.
- `retros/2026-09-16-retro.md` and `~/.config/sh-agent/proposals/2026-09-16-retro.md` — the eight proposals with their diffs.
- `changes/2026-09-16-linkshelf-simulation.md` and `changes/2026-09-16-bug-039-…` to `…-bug-043-…` — the entries.
- `.agents/scratch/sim-kit/` (not committed): `gaps.md` with evidence per gap, `results/` with every run's output, transcript copies and `summary.tsv`.

## Next Steps

1. **Pick the retro proposals to apply** (`~/.config/sh-agent/proposals/2026-09-16-retro.md`; say numbers or "go all"): (1) changelog `git show --stat` clause — crew P3 remainder; (2) DoD step 4 EMPTY at every checkpoint — crew P4 as drafted; (3) two agents contradicting on one mechanism means the rule is the defect; (4) a filter that fires is evidence only after its matching line is read — BUG-029 graduates; (5) the reviewer greps every restatement of a changed rule — fifth occurrence, and BUG-043's second fix was exactly this; (6) the golang skill names its HTTP and testing reference files; (7) init-project notes Claude Code's headless loop refusal; (8) graduate BUG-040 and BUG-043. Applying is a task in this repository under its Definition of Done (`sh-agent install`, one start of each CLI in the sandbox), then `sh-agent install` here.
2. **On this box:** add `"attribution": {"commit": "", "pr": ""}` to `~/.claude/settings.json`; `sh-agent status` then prints no note.
3. **On the Mac:** `sh-agent update --check`, `sh-agent update`, `sh-agent status`, the three smoke checks of `testing/2026-09-03-verify-real-home.md`.
4. **The mirrors** (styrman-crew, argolis-devbox) are still at a81ef4b: refresh from this HEAD by each `VENDORED.md` recipe; argolis pushes deploy, crew ships by release (its handoff recommends v1.34.0).
5. **Follow-ups** when convenient: G20 (define `<app>` in the verifier), G11 is Proposal 6, G4 is Proposal 7; option C's filing clause in `/preflight` Step 4 and BUG-031's wait rule from earlier handoffs; the two deferred remediations.
