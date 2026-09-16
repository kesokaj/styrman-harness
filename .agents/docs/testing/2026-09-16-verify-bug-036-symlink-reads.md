---
owner: verifier
date: 2026-09-16
---
# Verification — bug-036-symlink-reads (redesigned fix, re-run)   (2026-09-16)

Range reviewed: uncommitted working-tree diff against HEAD (d167cf8), scoped to
`internal/harness/inventory.go`, `internal/harness/render.go`,
`internal/harness/pathguard_test.go` (`git diff HEAD -- <those three files>`).

Cites `.agents/docs/testing/sh-agent/guide.md` §1-§3.

## Status: PASS

Scope: §1 Go checks, §2 harness self-check, §3 scratch-home install/drift, §3a
applicability check, plus an independent reproduction of the reviewer's
leaf-guard-vs-tree claim. §4 sandbox CLI starts and §5 real-home smoke checks
were not run (out of scope per dispatch); prior logs cited below.

## Commands

- `gofmt -l .` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -v ./...` → exit 0; output saved to
  `.agents/scratch/verify-bug-036-3/go-test-1.txt`.
  counts: pass 54, fail 0, skip 1 at d167cf8+dirty for `go test -race -v ./...`
  vs previous: `.agents/docs/testing/2026-09-16-verify-bug-036-symlink-reads.md`
  (prior run, same slug, same sha+dirty), pass 47→54, skip 1→1. The rise of 7
  is exactly `TestManagedDirsRefuseSymlinks`'s parent line plus its six
  subtests (`context`, `agents`, `rules`, `commands`, `workflows`, `skills`) —
  confirmed by grepping the new test's block in the saved log; no other test
  name changed count. Fail stayed 0, skip stayed 1 (unrelated pre-existing
  skip, unchanged by this diff).
- `go run ./cmd/sh-agent list` → exit 0; `### Skills (27)` and
  `### Subagents (7)` headers both print, matching the repo's own layout.
- `go run ./cmd/sh-agent render | wc -l` → `93` (expected).
- Scratch-home install:
  `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug-036-home-3`
  → `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`.
  `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug-036-home-3`
  → `drift none — 259 files match the repo` (the 260th managed path is the
  binary itself, not diffed as a content file); `repo has 27 skills, 7 agents,
  0 commands, 0 rules`. `list` inside this run and independently
  (`find skills -maxdepth 1 -type d | wc -l` → 28, i.e. 27 skill dirs plus the
  `skills/` root; `find skills -maxdepth 1 -type l | wc -l` → 0) confirm this
  repo's own `skills/` is a real directory of real directories, unaffected by
  `checkManagedTree`'s special-casing of `skills` as a possible link target.

## §3a — plain-copy / rebuild-decision tier

Does not apply. The diff touches only `internal/harness/inventory.go`,
`render.go`, and their test file — none of `IsCheckout`, `cmd/sh-agent`'s
rebuild decision, or `internal/harness/repo.go` (checked with
`grep -n "IsCheckout" internal/harness/repo.go cmd/sh-agent/*.go`, none of
which appear in the reviewed diff). §3a is for changes to that rebuild logic
specifically, which this fix does not touch.

## Coverage of changed code

- `checkManagedTree(repo)` (new, in `inventory.go`) — covered by
  `TestManagedDirsRefuseSymlinks` (6 subtests: one per managed directory,
  including the `skills` directory-itself case) and exercised transitively by
  `TestContextAgentsAndDocsRefuseSymlinks` (leaf-level symlinks under each of
  the five flat directories) and by every other `Plan`/`LoadInventory` test
  that runs against a repo without any symlink present (regression coverage
  for the non-error path). Called from both `LoadInventory` and
  `RenderContext`; `RenderContext`'s call site has no dedicated new test but
  is reached by every existing render test plus `TestManagedDirsRefuseSymlinks`
  running through `Plan`, which calls both `LoadInventory` and (via the
  install path used elsewhere in the suite) rendering — not observed to be
  independently exercised for `RenderContext` alone; see Not examined.
- `readRegular` (only remaining call site: `LoadSkills`'s
  `skills/<name>/SKILL.md` read) — covered by
  `TestSymlinkedSkillMdIsRefused` (symlinked SKILL.md refused) and
  `TestSkillDirWithoutSkillMdIsSkipped` (missing SKILL.md still skipped, not
  an error) and `TestSkillFilesAndAdoptCopiesRefuseSymlinks` (pre-existing).
- `LoadSkills`'s error-handling change (return non-`ErrNotExist` errors
  instead of silently `continue`ing) — covered by the same two tests above.
- No added `//nolint`, `//lint:ignore`, `t.Skip` (beyond the pre-existing
  "symlinks unavailable here" skip guards, which carry their own reason text
  and are unchanged behaviorally — they already existed in the prior round),
  or other suppression found in the reviewed diff.

## Independent reproduction of the reviewer's leaf-guard-vs-tree claim

Built a minimal scratch copy of the module
(`.agents/scratch/verify-bug-036-leafguard-repro/`: `go.mod`, `cmd/sh-agent`,
`internal/harness`), then overwrote only `inventory.go` and `render.go` in
that copy with the pre-review leaf-only-guard versions saved at
`.agents/scratch/leafguard/{inventory.go,render.go}` (current
`pathguard_test.go`, with `TestManagedDirsRefuseSymlinks`, left in place).
`go build ./...` in that copy succeeded (no compile drift between the two
guard designs). Then:

```
go test ./internal/harness/... -run TestManagedDirsRefuseSymlinks -v
```

Output saved to `.agents/scratch/verify-bug-036-3/leafguard-repro.txt`:
exit 1, six `=== RUN TestManagedDirsRefuseSymlinks/<dir>` lines (one per
`context`, `agents`, `rules`, `commands`, `workflows`, `skills` — ruling out a
vacuous "no tests to run" pass), all six failing with exactly
`Plan err = <nil>, want a symlink refusal naming .../<dir>`. This matches the
main agent's claim verbatim and independently confirms it: the leaf-only
guard does not catch a symlink at the managed-directory level itself.

The real tree's own run (`go test -race -v ./...` above, and a targeted
`-run TestManagedDirsRefuseSymlinks -v` implicit in that same log) shows all
six subtests passing, also with six `=== RUN` lines present in
`.agents/scratch/verify-bug-036-3/go-test-1.txt`.

No source, test, or config file was modified to do this — only files under
`.agents/scratch/` were created or copied.

## Not examined

- §4 sandbox CLI starts (Claude Code and Antigravity) — not run per dispatch;
  no context/agents/skills wording changed by this diff, only the file-read
  guard in Go code, so no new sandbox scenario is implicated by this change.
  Cited as not re-run: `.agents/docs/testing/2026-09-03-verify-sandbox.md`.
- §5 real-home smoke checks — not run per dispatch. Cited as not re-run:
  `.agents/docs/testing/2026-09-03-verify-real-home.md` and
  `.agents/docs/testing/2026-09-04-verify-mac-real-home.md`.
- Full skill-matrix (`.agents/docs/testing/skill-matrix.md`) — not run; this
  change is a Go-code-only fix with no skill/agent content change.
- `RenderContext`'s specific call to `checkManagedTree` was not isolated in a
  dedicated unit test distinct from `LoadInventory`'s call; both call sites
  are reached transitively through `Plan` in the existing suite, but no test
  was found that calls `RenderContext` directly with a symlinked managed
  directory to prove that call site independently (as opposed to via `Plan`,
  which calls `LoadInventory` first and would already fail there).
- Did not review `.agents/docs/bugs/INDEX.md`, the plan file, or the two new
  review files (`2026-09-16-review-bug-036-code.md`,
  `2026-09-16-review-bug-036-docs.md`) for their own conclusions — out of
  scope for this verification pass, which was scoped to the three named
  source/test files.

## Environment

Linux 6.12.107+deb13-amd64, `go` version as resolved by `go run`/`go build`
(not pinned separately in this run — no version drift indicated by any
command). No Docker/Colima involved (no Dockerfile changed). Git: repo dirty
at `d167cf8` with the three files under review plus untracked docs files
listed in the dispatch's git-status snapshot; no git command in this run
altered the working tree (only `git status`, `git diff`, `git rev-parse` were
used).
