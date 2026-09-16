---
owner: changelog
date: 2026-09-16
---
# Go code: clear the board Phase 1 (2026-09-16)

## Summary
The Go half of clearing every open item from `.agents/docs/plans/2026-09-16-clear-the-board.md` Phase 1 and the review rounds that followed. `update` now refuses a dirty checkout on every path before any pull, build or install (BUG-046, ADR-027); `Inspect` counts untracked files explicitly to catch wholly-untracked skill directories that `status.showUntrackedFiles=no` would hide (security review); `Plan` refuses secret-shaped filenames in skill directories because `.gitignore` keeps them out of git, so review never sees them, and `install` would copy them to both vendor directories on every machine (security review); `LoadSkills` refuses duplicate skill names, since they install to one destination and walk order would decide which text every machine reads (code review); `Adopt` stages its copy before renaming to leave the existing skill untouched if copy fails (attribution Follow-up); `writeAtomic` syncs before the rename (attribution Follow-up); test fixtures run under `gitCleanEnv` (security review); `gofmt -l` gate scoped from whole repo to `cmd internal` because gofmt walks `.agents/scratch` where Go tool does not (BUG-048).

## Type
bugfix

## Files changed
- `cmd/sh-agent/main.go` — `update` refuses a dirty checkout on every path, calls `harness.Inspect` before any pull
- `internal/harness/repo.go` — `Inspect` passes `--untracked-files=all` to `git status --porcelain`; `Pull`'s dirty refusal cross-referenced with `update`
- `internal/harness/plan.go` — `Plan` refuses secret-shaped filenames via new `secretShaped` helper (`.env`, `.env.*`, `*.pem`, `*.key`, `*.tfvars`, `*.tfstate`, `*credentials*.json`, `*-key.json`)
- `internal/harness/inventory.go` — `LoadSkills` refuses duplicate names, reports both directories
- `internal/harness/adopt.go` — `Adopt` stages copy under `.agents/scratch/adopt-<name>-<pid>`, renames into place
- `internal/harness/apply.go` — `writeAtomic` syncs before the rename
- `cmd/sh-agent/main_test.go`, `internal/harness/harness_test.go`, `internal/harness/repo_test.go`, `internal/harness/plan_test.go`, `internal/harness/inventory_test.go`, `internal/harness/adopt_test.go`, `internal/harness/apply_test.go` — new tests: `TestUpdateNeverInstallsDirtyTree` (clean-current, dirty-current, dirty-behind, clean-behind, `--check`), `TestInspectCountsUntrackedFiles`, `TestPlanRefusesSecretShapedFiles` (seven shapes plus control), `TestLoadSkillsRefusesDuplicateNames`, `TestAdoptFailedCopyKeepsTheExistingSkill`, `TestWriteAtomicIsAtomicAndDurable`, `TestRenderContextRefusesSymlinks`, `TestAbandonTempReportsRemoveFailure`, `TestReadSettingsProblems`; `TestNoTmpInHarnessProse` now walks whole non-vendored skill directories plus `commands/` and `rules/` via new `vendoredSkills` parser of `skills/VENDORED.md`; test fixtures use `gitCleanEnv` stripping `GIT_DIR`, `GIT_WORK_TREE`, `GIT_CONFIG*` (security review low, BUG-045 class in test code)
- `AGENTS.md` § Setup commands — `gofmt -l cmd internal` (caveat added in `.agents/docs/testing/sh-agent/guide.md` and `agents/verifier.md`)
- `README.md` — setup commands updated to match
- `.agents/docs/deploy/local.md` — gofmt scope updated; reason stated

## Agents dispatched
verifier ×1, reviewer ×2 (code and docs lenses), security ×1, bugs ×1, changelog ×1

## Verification
`.agents/docs/testing/2026-09-16-verify-clear-the-board.md` — PASS, `go test -race -v -count=1 ./...` pass 92 fail 0 skip 1 (up from 79). Install into scratch home reports drift none. `sh-agent list` reports 27 skills, 7 agents. Three guards validated in must-fail state: tested with dirty refusal removed (install proceeded), with adopt staging removed (left directory), with prose test against pre-fix tree (passed, so test now catches what it did not before).

## Docs
No docs changed in the Go changes themselves. Docs changes (skills, context, rules) are Entry 2. Deploy path and setup commands (AGENTS.md, README, deploy/local.md, testing guide, verifier agent) updated for `gofmt` scope and `BUG-048` caveat; all listed above.

## Requirements
n/a — all fixes are internal refactoring and hardening, no new requirement, no interface change.
