---
owner: reviewer
date: 2026-09-30
---
# Review — project-simplify   (2026-09-30)

Range reviewed: the whole repository at HEAD (3478a1c), no diff.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing broken. The codebase is unusually well-documented for its complexity (nearly every non-obvious branch cites a BUG-/D-/ADR- number for the edge case it guards), `deadcode ./...` finds no genuinely unreachable production code, `go mod tidy -diff` shows no unused module, and the one cross-package dependency between the two internal packages (`internal/search` calling `harness.SecretShaped`) is exactly the kind of small, single-function reuse this lens wants to see rather than a duplicated reimplementation. The one real finding is a duplicated atomic-write sequence that has already needed a manual two-place fix once (evidenced by the code's own comments); the rest are lower-confidence architectural observations, not defects.

## Findings

- **[medium / confidence: high]** `internal/harness/apply.go:309-347` (`writeAtomic`) and `cmd/sh-agent/main.go:908-935` (the settings-file writer inside `attributionSet`) independently implement the same temp-file-then-rename sequence: `CreateTemp` in the destination directory, write, `Sync`, `Chmod`, `Close`, `Rename`, with cleanup on each failing step. This is not a style duplication — the apply.go comment at line 328-332 says the fsync step "was fixed there first" in the settings writer and had to be ported across, i.e. the exact failure mode (a fix landing in one copy and needing manual porting to the other) has already happened once. Shorter fix: extract the shared steps (temp-create, write, sync, chmod, close, rename, cleanup-on-error) into one helper in `internal/harness` that returns a plain `error`; have `attributionSet` call it and translate the returned error into its own `"settings: %s could not be written (%v)"` messages, instead of re-running the five-step sequence inline. This removes ~25 duplicated lines and the two-place-fix hazard.

- **[low / confidence: low]** `cmd/sh-agent/main.go` (1199 lines) mixes CLI dispatch (`run`) with several self-contained pieces of logic that have no dependency on `flag`, `os.Args`, or CLI output framing and could live in `internal/harness` as independently testable units: the settings-JSON attribution editor (`attributionSet`, `attributionHint`, `attributionOff`, `readSettings`, `settingsUnwritable`, `abandonTemp`, roughly lines 815-1035) and the PATH-shadow reporter (`pathHint`, lines 788-814, which is a thin wrapper over `harness.PathShadows`/`harness.RunningElsewhere`, both already in `internal/harness`). Moving them would shrink the one file two engineers currently share for "install output," "attribution notice," and "update/zip" work, letting those be edited in parallel without touching the same 1199-line file. Not a correctness issue — flagged for the "files independently editable" test in this lens, at low confidence since the current grouping (everything the `install`/`update` commands print) is a defensible single concern too.

## Criteria

1. Dead code — pass. `deadcode ./...` reports one candidate (`internal/harness/release.go:422 SwapReleaseSigner`) but it is called from `cmd/sh-agent/zip_test.go` (a different package's test), so it is a deliberate test seam, not dead code; the tool's non-test build graph simply can't see it. No commented-out code blocks, no unreachable branches, no unused module found (`go mod tidy -diff` empty).
2. Flow tracing — not in this lens
3. Standards conformance — not in this lens
4. Directory layout — not in this lens
5. Simplicity (DIY-except-security) — pass, with the one reuse finding above counted here too (the duplicated write sequence is the main simplicity gap found); the security-sensitive path-safety code (`zipNameSafe` in `internal/harness/zip.go`) layers a manual absolute/`..`/backslash check in front of `filepath.IsLocal`, which looks redundant at first read but is deliberate defense against Windows-backslash paths that `IsLocal` does not special-case on Unix — not flagged, since this is exactly the kind of belt-and-suspenders the DIY-except-security exemption asks reviewers not to second-guess.
6. Decoupling — pass, with the one low-confidence note above; the `internal/search` → `internal/harness` dependency is a single well-scoped function (`SecretShaped`) called one way only, which is the interface size this criterion wants, not a reach-in.

## Not examined

Markdown surfaces (`context/`, `agents/*.md`, `commands/`, `rules/`, individual `skills/*/SKILL.md` bodies) were surveyed only at the directory/manifest level (skill count against `AGENTS.md`'s "Skills (28)" and `skills/VENDORED.md`'s provenance table, which reconcile); their prose was not read line by line for dead/redundant instructions. `internal/search/store.go`'s SQLite schema-versioning path and `internal/harness/report.go`'s clone/push flow (both flagged in the initial survey as candidates) were read only at the function-signature level, not traced for efficiency; nothing suspicious surfaced in that pass, but a full read could still find something this budget didn't reach. Test files were not reviewed (out of lens: flow tracing/standards, and reuse across test helpers specifically was not checked).
