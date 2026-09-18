---
owner: verifier
date: 2026-09-18
---
# Verification — BUG-079 stale running binary note  (2026-09-18)

Range reviewed: c868ed4..bc745dc

## Status: PASS

Scope: complete — tier 1 (Go checks), tier 3 (scratch-home install/status/drift), `list`,
`govulncheck`, the named test subset, and the live dual-binary PATH check all ran;
`.agents/docs/testing/guide.md` already existed and was not modified (no step change).

## Commands
- `gofmt -l cmd internal` → exit 0, no output
- `go vet ./...` → exit 0
- `go test -race -count=1 ./...` → exit 0, `ok` for all 3 packages (root, cmd/sh-agent, internal/harness)
- `go test -race -count=1 -v ./...` (for counts) → exit 0; counts: pass 229, fail 0, skip 0
  at bc745dc for `go test -race -count=1 -v ./...`
  vs previous: `.agents/docs/testing/2026-09-18-verify-post-v1.2.2-docs-2.md`, pass 215→229,
  skip 0→0. Delta of 14 = TestRunningElsewhere (1 parent + 8 subtests = 9) +
  TestUpdateNamesTheBinaryItRanFrom (1 parent + 3 subtests = 4) +
  TestCheckoutUpdateNamesTheBinaryItRanFrom (1) = 14. Matches exactly.
- `go test -race -count=1 -v -run 'TestRunningElsewhere|TestUpdateNamesTheBinaryItRanFrom|TestCheckoutUpdateNamesTheBinaryItRanFrom' ./...`
  → exit 0 in cmd/sh-agent and internal/harness (root package: "no tests to run", expected —
  none of these live there). All named tests and subtests PASS:
  - TestUpdateNamesTheBinaryItRanFrom: reachable_by_name:_named, not_on_PATH_at_all:_silent,
    --check_says_nothing — all PASS
  - TestCheckoutUpdateNamesTheBinaryItRanFrom — PASS
  - TestRunningElsewhere: 8/8 subtests PASS (running the installed file itself; via a link;
    a system copy reached by name; reported at the PATH entry whichever order; a file no PATH
    entry reaches — not reported; reached only through a link — reported at the PATH entry not
    the exe path; installed file missing but a reachable other copy still reported; exe missing)
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug079/verify-home` → exit 0,
  "260 written, 0 unchanged, 0 removed, 0 kept"
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug079/verify-home` → exit 0,
  "drift none — 259 files match the repo" (repo shown "uncommitted changes" because this
  verifier's own new log/review files are untracked in the working tree — expected, not a defect)
- `go run ./cmd/sh-agent list` → exit 0, 27 skills / 7 agents listed, parses cleanly
- `govulncheck ./...` (go1.27.1) → "No vulnerabilities found." exit 0

## Live check (tests cannot fake this)
Built two distinct-path binaries: `.agents/scratch/verify-bug079/bin-a/sh-agent` and
`.agents/scratch/verify-bug079/bin-b/sh-agent` (identical content, same source rev, but two
separate files/paths — `RunningElsewhere` keys on identity/path, not content, so this is valid).

Positive case (`update`, bin-b on PATH by name) — REFUSED, but for a reason unrelated to the
feature under test: this working tree carries untracked files this verifier itself wrote
(`.agents/docs/testing/2026-09-18-verify-bug079-stale-running-binary.md` and a pre-existing
untracked `.agents/docs/reviews/2026-09-18-review-c868ed4-bc745dc-docs.md` not written by this
run), so `git status --porcelain` is non-clean and `update` refuses per its documented dirty-tree
guard:
```
$ PATH=".agents/scratch/verify-bug079/bin-b:$PATH" .agents/scratch/verify-bug079/bin-b/sh-agent update --home .agents/scratch/verify-bug079/home --repo "$(pwd)"
up to date with upstream (bc745dc)
sh-agent: /home/kesokaj/workspace/styrman-harness has uncommitted changes and `update` installs the working tree, not the fetched commit; commit or stash them, or run `sh-agent install` to deploy them deliberately
```
Per dispatch instructions, fell back to `install` (which does not print the update-only
`staleRunningNote`, so this leg does not directly exercise the note's live text — that mechanism
is instead confirmed via the passing `TestCheckoutUpdateNamesTheBinaryItRanFrom` unit test, which
drives `Reinstall`+`staleRunningNote` end to end against a scratch git repo fixture):
```
$ PATH=".agents/scratch/verify-bug079/bin-b:$PATH" .agents/scratch/verify-bug079/bin-b/sh-agent install --home .agents/scratch/verify-bug079/home --repo "$(pwd)"
building sh-agent from /home/kesokaj/workspace/styrman-harness
260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed
note: .agents/scratch/verify-bug079/home/.local/bin is not on PATH, so `sh-agent` runs .agents/scratch/verify-bug079/bin-b/sh-agent instead; put .agents/scratch/verify-bug079/home/.local/bin ahead of it
settings: .agents/scratch/verify-bug079/home/.claude/settings.json created with "attribution": {"commit": "", "pr": ""} — Claude Code's commit and PR trailers are off
```
(That "is not on PATH" line is the pre-existing `pathHint`, not the BUG-079 `staleRunningNote` —
correct, since `staleRunningNote` is only wired into `update`, never `install`, confirming the
dispatch's claim it is update-only.)

Negative case (bin-b NOT on PATH, `install`) — no stale-binary note printed, as expected:
```
$ .agents/scratch/verify-bug079/bin-b/sh-agent install --home .agents/scratch/verify-bug079/home2 --repo "$(pwd)"
building sh-agent from /home/kesokaj/workspace/styrman-harness
260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed
note: .agents/scratch/verify-bug079/home2/.local/bin is not on PATH, so `sh-agent` runs /home/kesokaj/.local/bin/sh-agent instead; put .agents/scratch/verify-bug079/home2/.local/bin ahead of it
settings: ...
```
`--repo` was accepted on this development build without complaint (not a release version, so no
refusal was expected or seen).

## Root cause
N/A — no failures found.

## Coverage of changed code
- `internal/harness/pathguard.go` `RunningElsewhere` — covered by `internal/harness/pathguard_test.go` `TestRunningElsewhere` (8 subtests, all identity/link/PATH-order/missing-file cases)
- `cmd/sh-agent/main.go` `staleRunningNote` — covered by `cmd/sh-agent/main_test.go` `TestUpdateNamesTheBinaryItRanFrom` (3 subtests: named, not on PATH silent, `--check` silent) and `TestCheckoutUpdateNamesTheBinaryItRanFrom` (checkout/Reinstall leg)
- Both call sites in `main.go` (release-update leg gated on `!*check`, checkout-update leg via `update()`) are each named by a distinct test above — no gap.
- Doc changes (`README.md`, `.agents/docs/deploy/{local,remote}.md`) match the shipped behavior (verified by reading the diff against the code: `staleRunningNote` text, `hash -r`/`rehash`/`type sh-agent` ordering).
- No `//nolint`, `@ts-ignore`, `t.Skip`, `.only(`, or other suppression added in this range (checked via the full diff above — none present).

## Not examined
- The release-update leg (`UpdateFromRelease` → `staleRunningNote`) was not exercised live end
  to end (that path needs a real signed release download, which is out of scope / no network per
  the hard limits); its call site was read directly in the diff and is covered by unit test
  `TestUpdateNamesTheBinaryItRanFrom` at the `staleRunningNote` function level, not through the
  full `UpdateFromRelease` flow.
- `sh-agent report` was not run (hard limit).
- Tier 2 (harness self-check `render | wc -l`), tier 4 (sandbox CLI starts) and tier 5
  (real-home smoke checks) were not run — this range touches no `context/`, `agents/`, or
  `skills/` files, so those tiers do not apply per the guide's own scoping.
- The live positive `update` (not `install`) leg with the stale-binary note's exact printed text
  was not captured live, for the environmental reason above (untracked verifier artifacts made
  the tree non-clean); this is compensated by the unit test but not independently re-observed
  live in this run.

## Environment
- OS: Linux 6.12.107+deb13-amd64
- Go: go1.27.1 linux/amd64 (user-level SDK ahead of system go, per project AGENTS.md)
- govulncheck: clean, no advisories
- No Docker/Colima involved (no Dockerfile in this range)
