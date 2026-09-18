---
owner: verifier
date: 2026-09-18
---
# Verification — sh-agent-works-everywhere-phase-1

Range reviewed: b393525..028a43b

## Status: PASS

Scope: complete — Go checks, named tests, govulncheck, and all phase-1 checkpoint legs (first install, second install, touched-source-file update, --dry-run mode probe) ran. `testing/guide.md` exists (dated 2026-09-17), no steps changed, so it was not rewritten.

## Commands

- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0.
- `go test -race -count=1 ./...` → exit 0, all packages `ok`.
- `go test -race -count=1 -v ./...` → exit 0; counts: pass 227, fail 0, skip 0 at 028a43b for `go test -race -count=1 -v ./...`; vs previous: `.agents/docs/testing/2026-09-18-verify-bug079-stale-running-binary-3.md` (pass 226, fail 0, skip 0) → pass 226→227, skip 0→0. The +1 is `TestInstallNamesEveryChange`, the test this range adds; matches the expected 227.
- `go test -race -count=1 -v -run 'TestInstallNamesEveryChange|TestInstallCorrectsTheModeOfAnUnchangedFile|TestInstallStatusResyncUninstall' ./internal/harness/` → all three `--- PASS`.
- `govulncheck ./...` → "No vulnerabilities found." (go1.27.1 on PATH; matches guide's go1.25.8+ expectation).

## Phase-1 checkpoint (plan `.agents/docs/plans/2026-09-18-sh-agent-works-everywhere.md` Phase 1)

All runs used `.agents/scratch/` homes; no git command touched this working tree; the clone was made with `git clone -q "$(pwd)" .agents/scratch/verify-phase-1-clone` and left untouched except one appended line to a copy of `skills/bash/SKILL.md` (chosen because `rules/` and `commands/` in this checkout hold only `.gitkeep` — no deployed file exists there to touch; noted as a deviation from the dispatch's suggested `rules/go-style.md`, which does not exist).

1. **First install** (`.agents/scratch/verify-phase-1`): summary line `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`. `grep -c '^add '` on the captured output = 260, exactly matching the summary's `260 written`.

2. **Second install**, same home: summary `0 written, 260 unchanged, 0 removed, 0 kept — 260 files managed`; `grep -E '^(add|update|remove|mode) '` on the output = 0 lines (no `update ~/.local/bin/sh-agent` this run — `go run` rebuilt the binary but produced byte-identical output, so no line was printed, which is the documented "may appear" case not firing here).

3. **Touched source file** (via `--repo` into a second scratch home, `.agents/scratch/verify-phase-1-home2`): after appending one line to the clone's `skills/bash/SKILL.md` and re-installing with `--repo`, the only change lines were:
   ```
   update  ~/.claude/skills/bash/SKILL.md
   update  ~/.gemini/config/skills/bash/SKILL.md
   update  ~/.local/bin/sh-agent
   ```
   Exactly the two vendor-deployed copies of the touched file as `update`, plus `update ~/.local/bin/sh-agent` — expected here because the binary embeds `skills/` via `go:embed`, so a changed skill file changes the rebuilt binary's bytes. Nothing else changed (257 unchanged, 3 written, 260 files managed).

   Test-setup note (not a product defect): `internal/harness/repo.go` `Reinstall` runs the rebuilt binary's install with `cmd.Dir = repo` (the `--repo` path), so a *relative* `--home` argument passed alongside `--repo` resolves against the repo directory, not the invoking shell's cwd, on the real-install path — but resolves against the invoking shell's cwd on the `--dry-run` path (which never calls `Reinstall`, since it only prints that the binary "would be" rebuilt). This asymmetry meant my first attempt at the dry-run leg landed in a different (empty) home than the one populated by the prior real installs, silently reporting a fresh `add`-everything run against an unpopulated directory. Worked around by using an absolute `--home` path for the dry-run leg, which resolves identically on both paths. This is a real inconsistency in relative-`--home` resolution when combined with `--repo`, outside the diff under review (`apply.go`'s change-line printing); flagging it rather than filing it, since it is not part of this range's coverage obligation. Recommend a follow-up bug record if the operator wants `--home` resolution audited.

4. **`chmod 666` + `--dry-run`** (absolute home, same touched-clone repo): before: `666`. Dry-run output:
   ```
   [dry-run] the binary would be rebuilt from <repo>
   mode    ~/.claude/skills/bash/SKILL.md (0666 -> 0644)
   update  ~/.local/bin/sh-agent
   [dry-run] 1 written, 259 unchanged, 0 removed, 0 kept — 260 files managed
   ```
   After: still `666` (chmod not applied under `--dry-run`, as required). Exactly one `mode … (0666 -> 0644)` line for the probed file; the `update ~/.local/bin/sh-agent` line is the documented dry-run rebuild-would-happen case, consistent with the guide's expectation.

## Root cause (FAIL only)

N/A — Status is PASS.

## Coverage of changed code

- `internal/harness/apply.go` `Install` — the changed print-per-branch logic is exercised by `TestInstallNamesEveryChange` (add/update/update-unreadable-previous/mode/remove branches; per dispatch, no subtests) plus `TestInstallCorrectsTheModeOfAnUnchangedFile` and `TestInstallStatusResyncUninstall`, all `--- PASS`, and confirmed end-to-end against real scratch homes above (add, update, mode, and the dry-run-does-not-chmod branch). No suppressions (`t.Skip`, `//nolint`, `//go:build ignore`, `.skip(`, `.only(`) found in this range — `git diff b393525..028a43b` touches only `internal/harness/apply.go`, its test file, `README.md`, and a spec doc; no such directives added.

## Not examined

- Full skill-matrix (`.agents/docs/testing/skill-matrix.md`) and sandbox CLI starts (guide § 4) — not run; this range touches only `internal/harness/apply.go`, its tests, `README.md`, and the delivery spec, none of `context/`, `agents/`, or `skills/`, so guide §§ 2 and 4 do not apply here.
- `report`/feedback-loop tier (guide § 3b) — not applicable; `internal/harness/report.go` is untouched by this range.
- Did not audit the relative-`--home` + `--repo` resolution asymmetry noted above beyond reproducing it once; no test was written or run against it since it falls outside this range's diff.

## Environment

Linux 6.12.107+deb13-amd64, go1.27.1 (user-level SDK ahead of system `/usr/bin/go`), git clean at 028a43b throughout. No Docker/Colima involved (no Dockerfile in this range).
