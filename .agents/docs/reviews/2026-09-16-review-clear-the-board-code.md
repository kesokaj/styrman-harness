---
owner: reviewer
date: 2026-09-16
---
# Review — clear-the-board-code (2026-09-16)

Range reviewed: `e08fa65..HEAD` (`cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`,
`internal/harness/adopt.go`, `internal/harness/apply.go`, `internal/harness/harness_test.go`
only, per dispatch — `skills/`, `context/`, `agents/` out of scope for this pass).

## Verdict: APPROVE WITH NITS

Nothing broken and nothing that fails a test. The three hardening changes (update's dirty
refusal, adopt's stage-then-rename, writeAtomic's Sync) are each correct on the path they were
built for, and the new tests are real — every one asserts a behaviour that would fail if the
guard it covers were reverted (verified by reading the pre-image logic each test exercises, and
by `go test -race ./...`, `go vet ./...`, `gofmt -l .` all clean). The one real gap is a
reachability hole in the staged-adopt design: a previous run's staging directory that survives a
crash is not cleaned or even detected until the *next* `adopt` of the same skill, and in that
window `LoadSkills` has no duplicate-name guard (unlike `LoadDocs`), so it is silently absorbed
as a second skill entry rather than refused.

## Findings

- **[medium / confidence: medium]** `internal/harness/adopt.go:67` + `internal/harness/inventory.go:162-213` — the staging directory `dest+".adopt-staging"` is created and renamed *inside* `skills/`. It is only cleaned (`os.RemoveAll(staging)`) at the *start* of the next `Adopt` call for that same skill. If a process is killed after `copyTree` finishes but before `os.RemoveAll(dest)`/`os.Rename` completes, the staging directory survives with a complete `SKILL.md` whose frontmatter `name:` equals the real skill's name. `LoadSkills` (`internal/harness/inventory.go:162`) reads every entry under `skills/` and has no duplicate-name check — unlike `LoadDocs` (`inventory.go:253-274`, which errors with `"%s and %s both define %q"`). Every `sh-agent install`/`status`/`diff` run in that window silently gets two `Skill` entries with the same `Name` but different `Dir`; `Plan` (`internal/harness/plan.go:98-133`) writes both to the identical destination `Rel` path (`~/.claude/skills/<name>/...`), so one copy silently overwrites the other and `status`'s reported skill count is briefly wrong. No test covers a staging directory surviving across two separate CLI invocations. Fix: either have `LoadSkills` reject/report duplicate names the way `LoadDocs` does, or have `install`/`status` refuse (or at least warn on) a `*.adopt-staging` directory found under `skills/`.
- **[low / confidence: medium]** `cmd/sh-agent/main.go:353` / `internal/harness/repo.go:119-121` — `update()`'s new dirty check runs before the only call to `harness.Pull`, so `Pull`'s own dirty refusal (`"the harness repo has uncommitted changes; commit or stash them before \`sh-agent update\`"`) is now unreachable from `update()` and only exercised directly by `TestBehindAndPull`. Not a bug today (the two messages don't contradict), but the two refusals can drift out of sync silently since only one of them is ever seen by an operator running `update`.
- **[low / confidence: low]** `internal/harness/adopt.go:67-68` — the fixed staging name `dest + ".adopt-staging"` passes `identLike`/`validName`, so a legitimately-named skill directory `skills/<name>.adopt-staging` (unlikely, not forbidden) would be silently `os.RemoveAll`'d by the very first line of the staging block before any copy starts, even without `--force` applying to it. Extremely narrow collision, flagged for completeness rather than as a practical risk.

## Criteria

1. Dead code — pass. No orphaned code; the old direct-stat/remove/copy path in `adopt.go` was fully replaced, not left beside the new one.
2. Flow tracing — pass with the one gap above. `update()`'s guard sits correctly before `Pull` and after both `--check` early returns (confirmed by `TestUpdateNeverInstallsDirtyTree`'s five subtests, all of which can fail); `writeAtomic`'s `Sync` follows the same close/remove-on-error pattern as the rest of the function, no fd leak on any path; `Adopt`'s stage/copy/remove/rename sequence correctly leaves the original untouched on every failure branch tested (`TestAdoptFailedCopyKeepsTheExistingSkill`). The one path not traced by any test is the staging-directory-surviving-across-runs case above.
3. Standards conformance — pass. Errors wrapped or returned as-is consistently with the file's existing idiom; `gofmt`/`go vet`/`go test -race` all clean.
4. Directory layout — pass. No new files misplaced; tests stay beside the code they cover.
5. Simplicity — pass. Stage-then-rename and `Sync`-before-rename are the standard minimal patterns for this problem, no new dependency.
6. Decoupling — pass. Changes are localized to `adopt.go`/`apply.go`/`main.go`'s `update`; no cross-package reach-ins introduced.

## Not examined

`skills/`, `context/`, `agents/` content (explicitly out of scope — sibling dispatch). Did not
attempt to actually crash a process mid-`Adopt` to reproduce finding 1 empirically; the gap is
traced through the source (`copyTree` → `RemoveAll(dest)` → `Rename`) and through `LoadSkills`'
absence of a duplicate check, not observed via a live repro.

## Outcome

1. [medium] Adopt's staging directory could survive a crash inside `skills/` and be loaded as a second skill of the same name — **fixed, twice over**: staging now lives at `.agents/scratch/adopt-<name>-<pid>` with a `defer os.RemoveAll`, so the skill walk can never see it (`internal/harness/adopt.go`), and `LoadSkills` now refuses two skills of one name outright, naming both directories, so the hazard is closed however a duplicate arises — a hand-copied directory, a mirror, an interrupted anything (`internal/harness/inventory.go`, `TestLoadSkillsRefusesDuplicateNames`). The adopt test's staging assertion was rewritten around the real invariant (`assertSkillsClean`: nothing in `skills/` that is not a skill, and no two loaded skills sharing a name), since the old one became vacuous with the move.
2. [low] `Pull`'s dirty refusal is unreachable from `update` and the two wordings could drift — **fixed** by cross-referencing: `internal/harness/repo.go` says `update` refuses earlier with a message covering both paths and that this guard is for direct callers; `cmd/sh-agent/main.go` says why the guard sits before both branches. Kept rather than deleted: `Pull` is exported and a direct caller still needs it.
3. [low] a real skill named `<name>.adopt-staging` would be removed — **fixed** by the same move: no staging name exists inside `skills/` any more.

Verified after the fixes: `gofmt -l cmd internal` clean, `go vet ./...` clean, `go test -race -count=1 ./...` passing. Not re-dispatched: the fixes were applied by the main agent and re-read, and the operator was waiting on the session; a third round is named in the handoff as not run.
