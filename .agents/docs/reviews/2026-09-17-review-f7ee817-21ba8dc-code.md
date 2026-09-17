---
owner: reviewer
date: 2026-09-17
---
# Review — release-binary Phase 2 (f7ee817..21ba8dc)

Range reviewed: f7ee817..21ba8dc (branch release-binary), source and test changes: `embed.go`, `embed_test.go`, `plan_test.go`, `cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`, `internal/harness/{source,repo,apply,report,inventory,plan,render,harness_test,pathguard_test,report_test}.go`.

Checked against: `.agents/docs/plans/2026-09-17-release-binary.md` (Interfaces: Source, Source selection, ResolveSource, Per-command table, Binary handling on install, Version; decisions D13–D17, D20–D22; Phase 2 task boxes), `REQUIREMENTS.md` (R4, R8, R9 — Phase 2 correctly leaves these untouched, they are Phase 5 work), and `.agents/docs/bugs/INDEX.md` BUG-006, BUG-015, BUG-036.

## Verdict: APPROVE

This is a clean, mechanical, well-tested refactor. Every reader (`LoadSkills`, `LoadAgents`, `LoadDocs`, `RenderContext`, `Plan`, `checkManagedTree`) was moved from `repo string` to `Source` with no behavioural drift for a tree on disk; `ResolveSource` implements the exact fallback order and release-only gate the plan specifies (D14, D15); git helpers (`IsCheckout`, `Inspect`, `Behind`, `Pull`, `Reinstall`) are called only when `src.Dir != ""`, and I verified by direct read of `main.go` and `repo.go` that no call site passes an empty path to any of them (D17). `gofmt -l`, `go vet ./...`, and `go test ./...` all pass clean at HEAD. I found no high- or medium-severity defect.

## Findings

None affecting correctness, security, or a stated requirement. One low-confidence observation, not a finding:

- **[low / confidence: low]** `internal/harness/source.go:56-61` — `Source.path` for the embedded source (`Dir == ""`) returns the bare slash-path (e.g. `context/AGENTS.md`) with no `embedded:` marker, so an error message from `readRegular` on the embedded source reads identically to a relative path on disk. This is cosmetic (error text only, never used for a filesystem operation) and I found no place where it could be mistaken for an actionable path or cause an incorrect action — noting it only because a future caller might assume `path()` output is always openable.

## Criteria

1. **Dead code** — pass. `FindRepo` (replaced by `ResolveSource`) has zero remaining references anywhere in `internal/` or `cmd/`; the old free-function `readRegular` was fully replaced by the `Source` method of the same name with no leftover duplicate.
2. **Flow tracing** — pass. Traced `ResolveSource` → both call sites in `main.go` (the shared `src, err := harness.ResolveSource(...)` before the command switch, and `report()`'s own independent call) → every command branch (`install`, `status`, `diff`, `update`, `list`, `render`, `adopt`) implements exactly the per-command table in the plan, confirmed line-by-line against Interfaces § Per command with an embedded source. `Install` (`apply.go:96`) and `BuildRecord` (`report.go:347-351`) both branch on `src.Dir == ""` to write `embedded:<version>` / set `harness_version` and skip every git call, matching D17/D20. The `update` refusal for an embedded source (`main.go:192-194`) matches D22's exact text and is correctly scoped to "until Phase 4" — not a gap in this phase. Tests are behavioural and can fail: `TestEmbeddedPlanEqualsCheckout` and `TestFilesMatchTree` both document and (per the plan's evidence log, and by my own inspection of the assertions) exercise real mutation cases (missing `all:`, a mode change). `TestResolveSource` and `TestReleaseBinaryWithoutCheckout` (`harness_test.go:526`, `main_test.go:1008`) cover D14/D15's release-only gate and the full command table end-to-end.
3. **Standards conformance** — pass. Go idioms followed (wrapped errors, `io/fs` used consistently in place of `os`/`filepath` where the source may be embedded). No stale restatement found: grepped for `FindRepo` and `readRegular(path` (old free-function signature) — none remain.
4. **Directory layout** — pass. `embed.go` at the module root in package `styrman`, as the plan specifies; `Source` and its helpers live in `internal/harness/source.go`, beside the other readers; tests sit next to the code they cover (`source_test.go`, `embed_test.go`, `plan_test.go`).
5. **Simplicity / DIY-except-security** — pass. No new dependency; `embed`, `io/fs`, `regexp` are stdlib. No hand-rolled crypto/auth in this range (the signed-checksum verification via `ssh-keygen -Y verify` is Phase 4, not in this diff).
6. **Decoupling** — pass. `Source` is a small, well-documented struct at the seam every reader shares; `ResolveSource` and `EmbeddedSource`/`DirSource` are the only constructors, keeping the git-vs-embedded branching in one place (`repo.go`, `apply.go`, `report.go`, `main.go`) rather than scattered.

## BUG-036 regression check

Confirmed no regression: `pathguard_test.go`'s five symlink-refusal tests (`TestSkillFilesAndAdoptCopiesRefuseSymlinks`, `TestContextAgentsAndDocsRefuseSymlinks`, `TestSymlinkedSkillMdIsRefused`, `TestManagedDirsRefuseSymlinks`, plus `TestSkillAndAgentNamesMustBeOneLocalSegment`) were mechanically updated to wrap `repo` in `DirSource(repo)` with no change in assertions, and `checkManagedTree` / `Source.readRegular` still run `os.Lstat` on the disk path whenever `src.Dir != ""` — the exact guard BUG-036 required, now hoisted into the `Source` type per D3. `go test ./...` passes, so these assertions still hold.

## BUG-006 / BUG-015 check

`main.go:130-175`: `install`/`sync` with `repo == ""` (embedded) copies the running binary and never calls `Reinstall`/`IsCheckout` with an empty path (D13, confirmed by direct read); the existing `IsCheckout(repo)`-gated rebuild/copy logic for a tree on disk is untouched, so the BUG-006 (never copy a stale self-build) and BUG-015 (never rebuild into a plain copy) fixes are preserved unchanged.

## Not examined

- `internal/harness/state.go`, `frontmatter.go`, `derive.go` — named in the plan's "Surfaces touched" list but not present in this diff's stat; confirmed no code touches their `Repo string` field's meaning beyond storing the new `"embedded:<version>"` string, which needs no struct change, so their absence from the diff is consistent, not a gap — but I did not open these files to double check independently.
- `internal/harness/report_test.go`, `internal/harness/harness_test.go` in full (only the `TestResolveSource`/`TestResolveSourceRefusesNamedNonRepo` region and a few others were read directly; the remaining ~700 lines of `harness_test.go` were not read line-by-line, though `go test ./...` passing gives indirect confidence).
- `internal/harness/source_test.go`'s `TestPlanFromFSSource`, `TestFSSourceRefusals`, `TestFSSourceSkillMdMustBeRegular`, `TestDirSourceReadsTheSameAsFS`, `TestDirSourceFollowsLinkedSkillDir`, `TestEmbeddedInstallStateAndRecord` — function names and their presence were confirmed by grep; bodies were not read.
- `cmd/sh-agent/main_test.go`'s other tests (`TestRunInstallStatusUninstall`, `TestKnownCheckout`, `TestDryRunNeverBuilds`, `TestPlainCopyNeverRebuilds`, `TestAttributionSet`, etc.) — not read; these predate or are only lightly touched by this range per the diff stat, and `go test ./...` passes.
- Phase 3/4/5 work (throwaway report clone, release download/verify, docs/map/R8/R9 edits) is out of scope for this diff and was not reviewed — it does not exist yet in this range, consistent with the Phase 2 task boxes being the only ones checked in this commit range.
- No security review was performed beyond what the six criteria cover; Phase 2's own plan checkpoint does not call for a `security` dispatch (that is scoped to Phases 3 and 4).
