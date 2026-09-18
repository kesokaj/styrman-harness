---
owner: verifier
date: 2026-09-18
---
# Verification — sh-agent-works-everywhere-phase-2-3   (2026-09-18)

Range reviewed: e7fec88..42840fb

## Status: PASS

Scope: complete — every command listed in the dispatch ran, plus the live scratch checks; testing/guide.md exists (unchanged this run).

## Commands
- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -count=1 -v ./...` → exit 0. Output saved to `.agents/scratch/verify-sh-agent-works-everywhere-phase-2-3/go-test-1v.txt`.
  counts: pass 250, fail 0, skip 0 at 42840fb for `go test -race -count=1 -v ./...`
  vs previous: `.agents/docs/testing/2026-09-18-verify-sh-agent-works-everywhere-phase-2-2.md`, pass 248→250, skip 0→0 (expected: this range adds two subtests to TestInstallMakesTheConfigDirPrivate, confirmed — that test now has 4 subtests: absent, 0755, dry_run_names_it_and_leaves_it, a_symlinked_runtime_dir_is_left_alone)
- `GOOS=windows go build ./...` → exit 0.
- `GOOS=darwin GOARCH=arm64 go build ./...` → exit 0.
- `go test -race -count=1 -v -run 'TestInstallMakesTheConfigDirPrivate|TestThrowawayParentIsPrivate|TestDiffUsesTheHomeNotTmp' ./internal/harness/` → exit 0, all PASS (verbatim in `.agents/scratch/verify-sh-agent-works-everywhere-phase-2-3/targeted.txt`):
  ```
  --- PASS: TestThrowawayParentIsPrivate (0.00s)
  --- PASS: TestDiffUsesTheHomeNotTmp (0.00s)
  --- PASS: TestInstallMakesTheConfigDirPrivate (0.01s)
      --- PASS: TestInstallMakesTheConfigDirPrivate/absent (0.00s)
      --- PASS: TestInstallMakesTheConfigDirPrivate/0755 (0.00s)
      --- PASS: TestInstallMakesTheConfigDirPrivate/dry_run_names_it_and_leaves_it (0.00s)
      --- PASS: TestInstallMakesTheConfigDirPrivate/a_symlinked_runtime_dir_is_left_alone (0.00s)
  ```
- `govulncheck ./...` → "No vulnerabilities found." (go1.27.1 on this box, no advisories).

## Live scratch checks (`.agents/scratch/verify-phase-2-3`)
- Fresh `go run ./cmd/sh-agent install --home .agents/scratch/verify-phase-2-3` → exit 0, 260 files managed; `.config/sh-agent` created directly at `stat -c %a` = `700` (ConfigDir creates the runtime dir private only — no prior mode to tighten, so no "mode" line on a fresh install).
- `chmod 755 .../.config/sh-agent` → confirmed `755`.
- `install --dry-run` → exit 0; printed `mode    ~/.config/sh-agent (0755 -> 0700)`; dir left at `755` afterward (dry run names it and touches nothing — confirmed).
- Plain `install` (dir still 755, all files unchanged) → printed `mode    ~/.config/sh-agent (0755 -> 0700)` then `0 written, 260 unchanged, 0 removed, 0 kept`; dir now `700`.
- To confirm ORDER against an add line (a summary-only run has no add/update lines), re-chmod'd to `755` and deleted one installed file (`.claude/skills/bash/SKILL.md`), then ran `install` again. First three lines quoted verbatim:
  ```
  building sh-agent from /home/kesokaj/workspace/styrman-harness
  mode    ~/.config/sh-agent (0755 -> 0700)
  add     ~/.claude/skills/bash/SKILL.md
  ```
  The mode line precedes the first add line, as required; dir back to `700`.
- Symlink case: fresh home `.agents/scratch/verify-phase-2-3-symlink` whose `.config/sh-agent` is a symlink to `.agents/scratch/verify-phase-2-3-symlink-target/sh-agent` (mode 755 on the real dir). `install` → exit 0; no `mode` line in the output (`grep -n mode` on the output file returned nothing); target dir still `755` afterward. Confirmed: symlink left alone, target mode untouched.

## Root cause  (FAIL only)
N/A — Status PASS.

## Coverage of changed code
- `internal/harness/state.go` (`ConfigDir`, `TightenConfigDir`) — covered by `internal/harness/audit_test.go` (unit) and `TestInstallMakesTheConfigDirPrivate` (subtests: absent, 0755, dry_run_names_it_and_leaves_it, a_symlinked_runtime_dir_is_left_alone), plus the live scratch runs above exercising the same four scenarios end-to-end through the built binary.
- `internal/harness/apply.go` (3 added lines — diff/report path, "never change an existing dir's mode") — covered by `TestDiffUsesTheHomeNotTmp` and the audit_test.go changes in this range.
- No production source file in this range lacks a covering test.
- Suppressions added in this range (`git diff e7fec88..42840fb` on `*.go`, grepped for `nolint|skip|ignore|expect-error|eslint-disable|.only(`): none found.

## Not examined
- `TestInstallMakesTheConfigDirPrivate`'s failed-chmod-is-a-note path was not separately exercised live (no real permission-denied environment was constructed in scratch; only the test-suite's coverage of it was relied on — that subtest is not one of the two new ones named in the dispatch, so it was not enumerated above, but it ran as part of the full `./...` suite).
- No frontend/JS surface in this repo — not applicable.
- Containers/Docker not touched in this range — not run.
- `sh-agent status`/`adopt`/`report` behavior around the same mode change was not separately probed live beyond what `audit_test.go` covers in the full suite run.

## Environment
- OS: Linux 6.12.107+deb13-amd64.
- Go: go1.27.1 linux/amd64 (user-level SDK ahead of /usr/bin/go, per AGENTS.md).
- No Docker/Colima involved in this range.
- Working tree: clean except this log file (untracked) at start and throughout; no git command used that changes the tree; all installs targeted `.agents/scratch/verify-phase-2-3*`, never a real home.
