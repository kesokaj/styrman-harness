---
owner: main
date: 2026-09-17
---
# Handoff — 2026-09-17 — release-binary: Phase 1 done, Phase 2 one task in

Session ID: 4e712c10-b55b-4126-ae29-579c7d94ac09
Project: /home/kesokaj/workspace/styrman-harness
Branch: release-binary (pushed; `main` untouched at 22b92ff; this checkout is left on the branch)

## Context

The operator wants places to install the harness from a downloaded release binary with the harness built in, instead of a clone, and `sh-agent report` to stop keeping a full clone of the feedback repository on every machine. Every release also gets release notes. The operator chose to continue the implementation in a fresh session.

## Completed

- f7ee817 — Phase 1: `.agents/docs/plans/2026-09-17-release-binary.md` (D1–D19), recon `.agents/docs/plans/2026-09-17-recon-release-binary.md`, cold read `.agents/docs/reviews/2026-09-17-coldread-release-binary.md` (three rounds, proceeded at the cap), ADR-033. The operator approved the R8 sentence and the R9 bullets as quoted in the plan's Interfaces, and D5 (signed checksums instead of the tag).
- 7f3c0f9 — Phase 2, task 1: `embed.go` (package `styrman`, `//go:embed all:context all:skills all:agents all:commands all:rules`, 133 files) and `embed_test.go` (must-fail seen, then pass). D19 release notes added to the plan.

## Open / In Progress

- The plan's first unticked box: Phase 2 "`Source` type; every reader takes it; link refusal kept for `Dir != \"\"`" — not started. The readers to change are listed in the recon brief § 1 (plan.go, inventory.go, render.go, frontmatter.go, derive.go, apply.go, state.go); `ResolveSource` and the per-command table are in the plan's Interfaces (D17).
- Nothing in `REQUIREMENTS.md`, README, AGENTS.md or the map changes before Phase 5, which lands on `main` as ONE commit (D11, R9).
- Requirements edits awaiting approval: none — R8 and R9 texts approved, applied in Phase 5.
- Branch `close-map-gaps`: still waiting for another place's queued proposal; its planned ADR now needs number 034 (032 and 033 are taken), and its reading clone moves to `~/.config/sh-agent/feedback-read` (plan D12).
- `/feedback` on this box: still to re-try after `sh-agent update` from `main` (BUG-065's interactive check) — switch this checkout to `main` first.

## Key Decisions

- D1–D19 in the plan. Most load-bearing: D5 signed checksums; D8 throwaway partial clone with one retry; D11 one squash commit to `main`; D14 embedded fallback only for a release binary; D17 `ResolveSource` and git helpers only for a tree on disk; D19 release notes in four groups.

## Corrections

- none logged this session.

## Agent notes (unverified)

- Recon brief § 6 wrongly says `TestHarnessMap` does not exist; it is at `cmd/sh-agent/main_test.go:987`.
- Dispatch rules that held all day: tool-call budget plus "write the artifact early"; the missing-guide sentence for the verifier (the guide exists here at `.agents/docs/testing/guide.md`).
- Agent-written ADR and changes drafts needed correction five times today; write the entry yourself from the diff when you hold the evidence (context § Definition of Done step 4).

## Environment State

- Checkout on `release-binary`, clean after this handoff's commit. This box's installed harness is ce8ed6e — do NOT run `sh-agent update` while on this branch.

## Artifacts to Read

- `.agents/docs/plans/2026-09-17-release-binary.md` — Interfaces, Failure modes, Verification, the unticked boxes.
- `.agents/docs/plans/2026-09-17-recon-release-binary.md` — every reader with file:line.

## Next Steps

1. Phase 2, task 2: the `Source` type and every reader moved to it, tests first (`fstest.MapFS`, the existing symlink tests unchanged in meaning).
2. Phase 2, tasks 3–4: `Plan` from embedded equals from checkout; `ResolveSource`, the per-command table, `state.json` `embedded:<version>`, the version stamp. Then the Phase 2 checkpoint (verifier dispatch, commit, push).
3. Phases 3–5 as the plan orders them.
