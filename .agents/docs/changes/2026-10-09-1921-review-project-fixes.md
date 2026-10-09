---
owner: main
date: 2026-10-09
---
# Project review fixes: 21 rows fixed, 3 struck, 5 deferred, two review rounds     (2026-10-09)

## Summary
The operator approved the merged project review's recommendations by row (`reviews/2026-10-09-review-project-merged.md`). Fixed in `3bffdcf`: the agent lint refuses `mcpServers` (row 1); the dream CLI runs in its own process group with a `WaitDelay` (row 2, BUG-169); the delivery claim's one home names v2.6.2 and `metrics`, pinned by `TestDeliverySpecCurrent` and run in the release procedure (row 3, BUG-167); seven usage-error sites exit 2 with the reason on the writer, `search.IsUsage` across the package boundary (row 4, BUG-168); `version` and `main`'s exit status have tests (row 13); a benign near-miss per redaction needle and three credential kinds (rows 8, 9); the review and init-project skill lines (rows 5, 6); `PRIORITY_TODO.md` deleted (row 10); the layout blocks name `internal/trail` (row 11); three stale one-liners (row 12); `trail.DocsRoot`, the dreamer reading skills through the installer's frontmatter parser, `installPassArgs` reused, `metricsSummary`, `readSmall` sized by the file, `IsBusy` and `Lint` deleted, three exports unexported (rows 14–20); the corpus figure per box (row 22); govulncheck by its full path (row 29). Review rounds 1 and 2 in `6cec23b`: an explicit `--vendor all` forwarded to the child install again (a regression the simplify row introduced), the dream group killed on SIGINT, SIGTERM, SIGHUP and SIGQUIT and after a clean exit with a lingering helper no longer a failure, the needles by kind with `-u` anchored on whitespace and the auth key in any quoting, and the delivery paragraph naming only what v2.6.2 ships. Struck with reasons: rows 7, 26, 27. Deferred to the operator: rows 21, 23, 24, 25, 28, the agent-key allow-list, and a `PRIVATE_KEY=` kind. The review cap was reached (round 2 landed on the mechanism round 1 reshaped): the round-2 fixes are tested and verified, and a third round is the operator's call (the record's § Round 2).

## Type
task (the fixes of a review; three bug records)

## Files changed
- `internal/harness/inventory.go`, `harness_test.go` — `mcpServers` refused for agents
- `internal/dream/run.go`, `run_test.go` — `callWithin`: process group, `Cancel`, `WaitDelay`, signals forwarded, group killed after `Wait`, `ErrWaitDelay` is success; four tests
- `cmd/sh-agent/main.go`, `usage_test.go`, `delivery_test.go`, `zip_test.go`, `search_test.go`; `internal/search/search.go` — exit 2 at seven usage sites, `usageError`/`IsUsage`; `TestUsageErrorsExit2`, `TestMainExitStatus`, `TestVersionVerb`, `TestDeliverySpecCurrent`, `TestInstallPassArgsForwardsAnExplicitVendor`; `installPassArgs`, `metricsSummary`; five tests re-pinned to exit 2
- `internal/harness/report.go`, `report_test.go` — three needle kinds, near-misses, pinned over-redactions
- `internal/harness/frontmatter.go`, `release.go`, `phase8_test.go`; `internal/search/export.go`, `reconcile.go`; `internal/memory/memory.go` — `Lint` and `IsBusy` gone; `latestRelease`, `loadDocs`, `refusal`
- `internal/trail/trail.go`, `trail_test.go`; `internal/metrics/claude.go`, `events.go`; `internal/dream/dream.go`, `collect.go`, `collect_test.go` — `DocsRoot`, `docsOf`, `harness.Split`/`Parse` in the dreamer, `readSmall`
- `AGENTS.md`, `README.md`, `.agents/docs/specs/sh-agent/{delivery,dream,search}.md`, `.agents/docs/deploy/{local,remote}.md`, `.agents/docs/testing/guide.md`, `skills/{review,init-project,harness-audit}/SKILL.md`, `skills/harness-audit/references/areas.md` — the docs of every behaviour above
- `PRIORITY_TODO.md` — deleted
- `.agents/docs/bugs/BUG-167-*.md`, `BUG-168-*.md`, `BUG-169-*.md`, `INDEX.md` — the three records
- `.agents/docs/reviews/2026-10-09-review-project-merged.md` — § Round 1, § Outcome, § Round 2

## Agents dispatched
reviewer (opus, code+docs, two rounds), security (two rounds), verifier (once, last). Their artifacts: `reviews/2026-10-09-review-dbd776d-3bffdcf-code-docs.md`, `reviews/2026-10-09-security-dbd776d-3bffdcf.md`, `testing/2026-10-09-verify-review-project-fixes.md`.

## Verification
`testing/2026-10-09-verify-review-project-fixes.md`: PASS at `6cec23b` on a frozen tree, 682 pass, 0 fail, 5 skip, govulncheck no new advisory, two must-fail controls red then green. Every fix had its test red before it; the main agent's own mutation of the three needles failed 46 assertion lines.

## Notes
- The review cap: a third round has not run; the operator decides (handoff Next Steps 1).
- Row 17's simplify claim "the child's behaviour is identical" was false when the home's state records a narrower vendor; the round-1 reviewer caught it and `installPassArgs` now forwards an explicit vendor, which also corrects `update`.
