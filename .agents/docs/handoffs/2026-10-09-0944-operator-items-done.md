---
owner: main
date: 2026-10-09
---
# Handoff — 2026-10-09 — the metrics handoff's five operator items done; one toolchain step left

Session ID: 5dd0d90e-7c79-4c04-8f88-4cd1feee99e8
Project: /home/ostling-simon/projects/styrman-harness
Branch: main

## Context

After the metrics store closed (`handoffs/2026-10-09-0848-metrics-store-done.md`), the operator approved all five queued items and said "You run this." All five are done. Details are in `changes/2026-10-09-*-bug-165-git-hardening-search-memory.md`.

## Completed

- **First real metrics ingest, on the real home, at the operator's instruction:** 103 sessions, 4,968 steps, 104 events. Report at `~/.config/sh-agent/reports/2026-10-09-0849-metrics.html`.
- **R3 wording edit:** REQUIREMENTS.md and ADR-055b (f7c607a).
- **BUG-165 fixed** (db001db).
  - Search and memory run git only through `internal/trail.Command`; search never runs `git status`. Search time per query is unchanged, about 116 ms.
  - Reviewer and security ran two rounds each.
  - Verifier: `testing/2026-10-09-verify-bug-165.md` PASS, 616 tests.
- **Claude `retro` control at 2c0f215:** the matrix row fails there too, so it predates the metrics feature (`testing/2026-10-08-verify-metrics-sandbox.md`).
- **`govulncheck`:** installed (user-level, `~/go/bin`) and run.
- **"Harness Delivery Map" artifact rebuilt from the README map:** https://claude.ai/code/artifact/ff9e4585-c34f-4e2b-9eb9-42ee3b39c979, version 3.

## Open / In Progress

- **`govulncheck` reports three stdlib advisories in go1.27.1, all fixed in go1.27.2:** GO-2026-6604 (`os`, via `os.Root` in search), GO-2026-6600 and GO-2026-6599 (`html/template`, via `metrics report`). The fix is a newer Go on this box; `go.mod` is never pinned (ADR-029). This is the operator's toolchain step, `.agents/docs/deploy/local.md` § Operator steps.
- **Requirements edits awaiting approval:** none.

## Key Decisions

- **The control run settles the Claude `retro` row.** It fails at 2c0f215 too, so the row's expectation of a harness-scope proposal is the matrix's question (left: operator), not a regression.
- **Search does not refuse a partial clone,** a stated limit (security low). Refusing would stop search in every partial clone.

## Corrections

- None logged.

## Agent notes (unverified)

- This box's `go` resolves to `/usr/local/bin/go` (go1.27.1). `deploy/local.md` describes a user-level SDK symlinked at `~/.local/bin/go`, which does not exist now. The guide's description of this box may be stale.

## Environment State

- Branch: main, clean, pushed at db001db (before this handoff). Newest tag v2.5.0. No processes running.
- Sandbox evidence is under `~/sandbox/metrics-sbx/` (runs `claude`, `claude-1`, `claude-control`, `agy`). The scratch kits are `.agents/scratch/kits/` (gitignored).

## Artifacts to Read

- `.agents/docs/bugs/BUG-165-search-git-status-unhardened.md`: the hardening rule as built.
- `.agents/docs/testing/2026-10-09-verify-bug-165.md`: the verification and the `govulncheck` output.

## Next Steps

1. **Operator:** move this box to go1.27.2 or later (`deploy/local.md` § Operator steps, which also needs its `~/.local/bin/go` description checked against `/usr/local/bin/go`). Then run `govulncheck ./...` and `sh-agent install`. This is the one item queued before the next release.
2. **Then:** cut the next release with `/preflight vX.Y.Z` when you choose. That range carries the metrics store, BUG-164 and BUG-165, and the changed `retro` skill; that skill's Antigravity run exists (`testing/2026-10-08-verify-metrics-sandbox.md`).
