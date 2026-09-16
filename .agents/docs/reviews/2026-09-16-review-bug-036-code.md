---
owner: reviewer
date: 2026-09-16
---
# Review — bug-036-code   (2026-09-16)

Range reviewed: `git diff HEAD -- internal/harness/inventory.go internal/harness/render.go internal/harness/pathguard_test.go` (working tree, uncommitted, on top of d167cf8)

## Verdict: APPROVE

The round-2 redesign is implemented correctly. `checkManagedTree` runs once at the boundary from both `LoadInventory` and `RenderContext`, and every call site of either (`internal/harness/plan.go:76/80`, `cmd/sh-agent/main.go:183/191/195/283`, `harness_test.go`) goes through one of them — there is no remaining path that reads `context/`, `agents/`, `commands/`, `rules/`, `workflows/` or `skills/` bytes without first passing the check. `LoadAgents` and `LoadDocs` are only reachable from `LoadInventory`, confirmed by grep. `readRegular` now has exactly one call site (`skills/<name>/SKILL.md` in `LoadSkills`), matching its comment. `gofmt -l .`, `go vet ./...`, and `go test -race ./...` all pass (0 findings, pass 47/fail 0/skip 1 including the added tests), and the four new/changed tests (`TestManagedDirsRefuseSymlinks`, `TestContextAgentsAndDocsRefuseSymlinks`, `TestSymlinkedSkillMdIsRefused`, `TestSkillDirWithoutSkillMdIsSkipped`) each fail meaningfully when the guard they cover is removed (per the round-2/3 ledger and my own re-run against HEAD).

The `skills` exception is asymmetric by design and correctly so: a symlink at `skills/` itself is caught by the unconditional `info.Mode()&os.ModeSymlink != 0` check in `checkManagedTree` *before* the `dir == "skills"` skip is reached (verified: that branch runs ahead of the skip, and `TestManagedDirsRefuseSymlinks/skills` exercises exactly this). A symlink at `skills/<name>` is intentionally let through `checkManagedTree` (`npx skills add -g` legitimately leaves one) and is instead resolved by `LoadSkills`'s own `filepath.EvalSymlinks`, with `readRegular` guarding the one leaf read (`SKILL.md`) that resolution reopens to a symlink attack. That is the correct division of labor and matches the operator's "one boundary check, not more per-reader guards" call — with the SKILL.md leaf as the one reader the boundary provably cannot reach.

One non-blocking observation, not raised as a finding because it fails safely rather than incorrectly: in `checkManagedTree`, a managed-directory name that exists as a plain file (not a symlink, not a directory — e.g. someone commits a file named `rules`) hits `!info.IsDir() || dir == "skills"` and is silently skipped with no error, for all five flat directories, not just `skills`. This doesn't let a symlink payload through (that path is caught earlier by the `ModeSymlink` check), and the subsequent `os.ReadFile`/`filepath.Glob`/`os.ReadDir` call in `RenderContext`/`LoadDocs`/`LoadAgents` still fails (ENOTDIR) rather than silently succeeding — so there's no concrete failure mode, just an unused branch in the "sensible for a non-directory" question the dispatch asked me to check.

## Findings

None — correctness, security, and rule findings only, and I found none.

## Criteria

1. Dead code — pass. `readRegular` has exactly one live call site matching its comment; no leftover branches from the round-2 redesign.
2. Flow tracing — pass. Every `LoadInventory`/`RenderContext` caller (traced through `plan.go`, `cmd/sh-agent/main.go`, `harness_test.go`) now passes through `checkManagedTree`; the four new/changed tests each fail without the mechanism they assert (confirmed against the round-2 leaf-only guard saved in `.agents/scratch/leafguard/` per the ledger, and re-run green against the current tree).
3. Standards conformance — pass. Errors wrapped or returned with `%w`/named path context throughout; `errors.Is(err, os.ErrNotExist)` used correctly instead of a bare `continue`; `gofmt`/`vet`/`test -race` all clean.
4. Directory layout — pass. New tests land in the existing `internal/harness/pathguard_test.go`, beside the code they cover.
5. Simplicity (DIY-except-security) — pass. One stdlib-only boundary check replaces several ad hoc leaf guards, per the operator's explicit direction; no new dependency.
6. Decoupling — pass. `checkManagedTree` is a single, package-internal seam called from exactly two entry points; no cross-package reach-in introduced.

## Not examined

- `internal/harness/plan.go`'s skill-copy `WalkDir` symlink refusal, `internal/harness/adopt.go`'s `copyTree`, and the round-1 arbitrary-file-read fix — all pre-existing, untouched by this diff, not re-audited beyond confirming they still compose correctly with the new boundary check.
- `cmd/sh-agent/main.go` beyond the four `LoadInventory`/`RenderContext` call sites and the unrelated `selfExecutable` read (confirmed not a repo-content read).
- `.agents/docs/testing/2026-09-16-verify-bug-036-symlink-reads.md`: out of my scope (not one of the three named files), but its "Coverage of changed code" section describes `readRegular` as covering `LoadAgents`, `LoadDocs`, and `RenderContext` directly — that is the round-2-*before*-redesign shape, not the code currently in the working tree (which uses `checkManagedTree` + plain `os.ReadFile` for those three, and confines `readRegular` to `SKILL.md` alone). Worth a second look by whoever owns that log, since a stale coverage claim in a verification artifact is exactly the kind of thing the project's hard rules ask to catch — but it names no defect in the source under review here.
- Did not re-drive the sandbox/real-home smoke checks (§4/§5 of the testing guide); relied on `go test -race ./...`, `go vet`, `gofmt -l .` for this round.

## Outcome

Three rounds; the report above is round 3's. The ledger:

- **Round 1 [high], `LoadSkills` swallowed the refusal** — `fixed` at `internal/harness/inventory.go` (`LoadSkills`): it continues only on `os.ErrNotExist` and returns every other error. Tests `TestSymlinkedSkillMdIsRefused` and `TestSkillDirWithoutSkillMdIsSkipped` in `pathguard_test.go`; with the swallowing `continue` restored in a scratch copy the first fails on both assertions.
- **Round 1 [low], the verifier log's coverage claim** — `fixed`: the verifier's re-run corrected it (`testing/2026-09-16-verify-bug-036-symlink-reads.md`).
- **Round 2 [high], a managed directory that is itself a link is followed** — `fixed` by redesign. A second round on the same mechanism went to the operator, who chose "One boundary check" over more per-reader guards. `checkManagedTree` runs once, from `LoadInventory` and `RenderContext`; `RenderContext`, `LoadAgents` and `LoadDocs` read with plain `os.ReadFile` again; `readRegular` remains only for `skills/<name>/SKILL.md`. Test `TestManagedDirsRefuseSymlinks`, six directories; against the leaf-only guard saved in `.agents/scratch/leafguard/` all six fail with `Plan err = <nil>`, checked for six `=== RUN` lines in both runs.
- **Round 3, APPROVE.** Its note that a plain file at a managed directory's name is skipped by the check — `left`: it fails safely downstream with `ENOTDIR`, as the report says, and no reader can deploy from it. Its side note on the verifier log — `rebutted`: the verifier rewrote that log while the review ran; the version on disk records pass 54 and lists `checkManagedTree` in its coverage section.
- **Follow-up, in BUG-036's record:** a direct `RenderContext` test with a linked `context/`; today that call site is reached only through `Plan`, which runs the check in `LoadInventory` first.
