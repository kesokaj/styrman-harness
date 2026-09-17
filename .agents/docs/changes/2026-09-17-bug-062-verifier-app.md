---
owner: changelog
date: 2026-09-17
---
# BUG-062: verifier `<app>` resolution — run `ls -d cmd/*/` first     (2026-09-17)

## Summary
Claude Code's verifier named `<app>` after the repository directory and its `AGENTS.md` heading instead of checking for a single `cmd/<app>` directory: in the demo repository (only `cmd/server`), it wrote `testing/demo/guide.md` instead of `testing/server/guide.md`. The fix: `agents/verifier.md` step 2 now explicitly runs `ls -d cmd/*/` and states it is not overridden by an `AGENTS.md` heading. The log template adds an `App: <name> (source)` line, making the resolution branch visible on every run. `agents/documentation.md` restates the three-branch order to keep both files aligned. Root cause: the resolution rule named the outcome of each branch but never said how to check them in order.

## Type
bugfix

## Files changed
- `agents/verifier.md` — step 2 now runs `ls -d cmd/*/` first, states `AGENTS.md` heading does not override it, and extracts the last path element as `<app>`; template line 71 gains `App: <name> (cmd/<name> | named in AGENTS.md | repository directory)` to record which branch fired
- `agents/documentation.md` — its restatement of the order now names the `ls -d cmd/*/` check and that an `AGENTS.md` heading does not override it

## Agents dispatched
explorer, reviewer, verifier, bugs, changelog

## Verification
- Verifier log `.agents/docs/testing/2026-09-17-verify-bug-062.md` PASS: `gofmt`, `vet`, `go test -race -v -count=1 ./...` all clean (124 pass, 1 skip, no fail); scratch install drift none, 259 files match; rendered `App:` line confirmed in both vendors' installed copies.
- Review log `.agents/docs/reviews/2026-09-17-review-1500bcf-bug-062.md` APPROVE WITH NITS: fix is sound, two low-severity wording findings found and fixed in same commit (strip `cmd/` and trailing `/` from `ls` output explicitly stated, "branch" renamed to "source" to avoid git-branch confusion), re-installed with drift none.
- Behavioural evidence from `.agents/docs/testing/2026-09-17-bug-062-repro.md`: 3 of 3 runs before the fix wrote `testing/demo/guide.md` (repo-name fallback, the bug); 2 of 2 runs after the fix ran `ls -d cmd/*/` and logged `App: server (cmd/server)`; repro-3 wrote `testing/server/guide.md`, repro-4 wrote no guide at all (BUG-063, below).

## Docs
- `agents/verifier.md` — agent-prompt file documenting verifier behavior; step 2 and template line 71 changed
- `agents/documentation.md` — agent-prompt file; restatement of resolution order added
- Trail artifacts: `.agents/docs/bugs/BUG-062-verifier-app-resolution-mismatch.md` (problem and fix), `.agents/docs/plans/2026-09-17-recon-bug-062.md` (goldfish diagnosis), `.agents/docs/reviews/2026-09-17-review-1500bcf-bug-062.md` (review), `.agents/docs/testing/2026-09-17-bug-062-repro.md` (repro log), `.agents/docs/testing/2026-09-17-verify-bug-062.md` (this run's verification)
- No README, specs, deploy guides, or user-facing docs changed

## Requirements
n/a

## Not exercised
- Antigravity with the fixed definition (operator's `sh-agent update` after push)
- Real-home install into `~/.claude` or `~/.gemini` (operator step, out of dispatch scope)
- Full skill matrix (`.agents/docs/testing/skill-matrix.md`)

## Note
A separate defect (BUG-063) was recorded in the same repro run: verifier's log said `Scope: complete` while a required guide was not written (a chained shell call was denied by the run's tool allowlist and never retried). Open, not fixed.
