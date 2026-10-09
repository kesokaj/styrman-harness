---
owner: main
date: 2026-10-09
---
# BUG-165: search and memory run git through the hardened builder; operator items from the metrics handoff

- **BUG-165 fixed** (`bugs/BUG-165-search-git-status-unhardened.md`). `internal/trail.Command` is now exported and is the one builder for git in a project repository. Search changes:
  - `internal/search/reconcile.go` no longer runs `git status`. A tracked file keeps the index's blob when the index's stat data still matches (`ls-files -s --debug`, read line by line; quoted names skipped; racily clean entries hashed). Otherwise its bytes decide.
  - Untracked files come from `ls-files --others --exclude-standard`.
  - Every noindex marker is checked whatever the prefix.
  - `search.go`'s `check-ignore` and the three git calls in `internal/memory` use the builder. Memory's `check-ignore` had been starting the checkout's fsmonitor program; the security review found it.
  - Search takes about 116 ms per query, the same as v2.5.0.
- Tests, each red before its fix: `TestGitViewRunsNoCleanFilter`, `TestGitViewRunsNoFsmonitor`, `TestGitViewRacilyCleanEditIsHashed` (and its mutation), `TestGitViewDeletedMarkerOutsidePrefix`, `TestIndexStatsIgnoresForgedNames`, `TestAddRunsNoProjectHook`.
- Reviews, two rounds each, both ending in `## Outcome`:
  - `reviews/2026-10-09-review-worktree-code-docs.md` (`reviewer`, code and docs);
  - `reviews/2026-10-09-security-worktree.md` (`security`).
- Verification: `verifier` wrote `testing/2026-10-09-verify-bug-165.md`, PASS. Results: 616 pass, 0 fail, 3 env-gated skips; §§ 3d and 3e exit 0; `govulncheck` reports only three stdlib advisories fixed in go1.27.2 (GO-2026-6604, -6600, -6599), an operator toolchain step (`deploy/local.md`).
- Docs:
  - `specs/sh-agent/search.md` § Freshness, with the partial-clone limit;
  - README's search row;
  - the allowed-command lists in BUG-164, the INDEX and `metrics.md`.
- Operator items from `handoffs/2026-10-09-0848-metrics-store-done.md`, all approved 2026-10-09 ("go all approved. You run this."):
  1. The first real metrics ingest ran on the real home: 103 sessions, 4,968 steps, 104 events, 0 failed. Report at `~/.config/sh-agent/reports/2026-10-09-0849-metrics.html`.
  2. R3 wording edit, ADR-055b, commit f7c607a.
  3. BUG-165 fixed, this entry.
  4. Claude `retro` control run at 2c0f215. The row also fails there, so the 0 of 2 predates the metrics feature; recorded in `testing/2026-10-08-verify-metrics-sandbox.md`. `govulncheck` was run: see above.
  5. The "Harness Delivery Map" artifact page was rebuilt from the README map: https://claude.ai/code/artifact/ff9e4585-c34f-4e2b-9eb9-42ee3b39c979 (version 3).
- Agents dispatched: `reviewer` (2 rounds), `security` (2 rounds), `verifier` (1).
