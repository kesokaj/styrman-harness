---
owner: explorer
date: 2026-09-28
---
# Recon — sh-agent update --from-zip   (2026-09-28)

## Status: COMPLETE

## Question
Add `sh-agent update --from-zip <file>`: unpack a GitHub-downloaded zip into a temp dir,
refuse zip-slip entries, install with the RUNNING binary (no rebuild, no replace), delete
the temp dir; must work from a release binary, which `ResolveSource` today refuses (ADR-034).
Map status drift, Install/Plan manifest recording, update's flag handling, README map +
TestHarnessMap, every living-doc "two ways in" statement, embedded-tree sizes, report's use
of `State.Repo`.

## Summary
The scope statement is in direct tension with a documented, ADR-backed design decision:
`.agents/docs/specs/sh-agent/delivery.md` § How a machine gets the harness says "Two ways,
and no third," and AGENTS.md rule 1 says "no third [way]" implicitly via "Two ways in, one
way back, nothing else." A zip unpack today is explicitly the demoted **third case** — "a
plain copy… it is no longer one of the ways in… it is the fallback." Making `update --from-zip`
work *from a release binary* turns that fallback into a second way a release-binary place can
move forward (alongside the signed-checksum download), which is exactly the kind of path
change AGENTS.md rule 3 requires going to the operator first, with an ADR, landing with the
map/doc changes in one commit — before writing any code. Mechanically the feature is buildable:
`cmd/sh-agent/main.go`'s `update` case (main.go:198-230) already branches on `repo == ""` for
the release-binary path calling `harness.UpdateFromRelease` — a `--from-zip` branch would sit
beside it, build a `harness.Source{Dir: <tempdir>, FS: os.DirFS(tempdir)}` (see `source.go:18`),
call `harness.Install`/`harness.Plan` directly (never `harness.Reinstall`, which always spawns
`go build`), and `os.RemoveAll` the temp dir afterward. The unresolved design problems are: (a)
the ADR/doc conflict above, (b) what `State.Repo` should record once the temp dir is gone
(status's `EmbeddedPrefix` branch at main.go:383 and `knownCheckout` at main.go:573 both assume
`st.Repo` is either a real path or an `embedded:` release string — a stale temp path breaks
both), and (c) zip-slip / top-folder handling (GitHub's zip has one top folder,
`styrman-harness-<branch>/`, that must be stripped before `IsRepo`/`Plan` see the tree).

## Relevant files
- `internal/harness/repo.go:22-48` — `ResolveSource`; release-binary branch (34-48) refuses
  any `--repo`/`$SH_AGENT_REPO` unconditionally and returns `EmbeddedSource(version)`. This is
  the function the scope names as "superseded… for the one file the operator names" — but
  `--from-zip` most naturally bypasses `ResolveSource` entirely (own resolution in `main.go`,
  producing a `DirSource` pointed at the unpacked temp dir) rather than adding a branch inside
  it, matching how `report()` (main.go:304-319) already builds its own `Source` instead of
  going through `ResolveSource` when it needs release-binary-specific behaviour.
- `internal/harness/repo.go:196-236` — `Reinstall`: rebuilds via `go build ./cmd/sh-agent` then
  runs that build's `install --repo <repo>`. Scope explicitly forbids this for `--from-zip`
  ("never rebuilding, never replacing the binary… installs that tree with the RUNNING binary").
  `update`'s existing checkout path (`cmd/sh-agent/main.go:474-524`, see `update()`) calls
  `harness.Reinstall` at line 517 — `--from-zip` must NOT reuse this function; it needs its own
  small path that calls `harness.Install`/`harness.Plan` in-process, the same shape `install`
  uses for a release binary (main.go:137-146, `installFiles`) or a plain copy (main.go:175-182).
- `internal/harness/source.go:18,36-54` — `Source` struct, `EmbeddedSource`, `EmbeddedPrefix =
  "embedded:"`, `IsReleaseVersion`. A `DirSource` helper (used by `absSource` in repo.go:90-96)
  is the existing constructor for a Dir-backed Source; `--from-zip` would use the same shape.
- `internal/harness/apply.go:19-148` (`Install`) — for a Dir source (`src.Dir != ""`, line
  119-123) it records `repo = src.Dir` (the literal temp path) and `commit = Inspect(repo).Commit`
  in `State{Repo: repo, Commit: commit, ...}` (line 134). `Inspect` (repo.go:135-148) checks
  `IsCheckout` first (`.git` present) — the unpacked zip tree has no `.git`, so `Inspect` returns
  an EMPTY `RepoInfo{}` silently (no error), giving `commit == ""` and a manifest that records
  `Repo: <temp-dir-path>` with an empty commit — a path that will not exist once the temp dir is
  removed. This is the finding behind item 1 of the dispatch: no `checkManagedTree`/link-refusal
  call happens automatically inside `Install` — that runs earlier, inside `RenderContext` and
  `LoadInventory` (`render.go:31`, `inventory.go:90,128`), which `Plan` calls (plan.go:76-83,
  98-99) BEFORE `Install` ever sees the files, so it does run, just not from `Install` itself.
- `internal/harness/plan.go:53-199` (`Plan`) — reads exactly `context/*.md` (via
  `RenderContext`), `skills/<name>/**`, `agents/<name>.md`, `commands/<name>.md` +
  `workflows/<name>.md` (comment says both map the same as commands; only `inv.Commands` is
  actually walked — `workflows/` is not a separate directory read, contrary to the doc comment
  at plan.go:61), and `rules/<name>.md` — five source directories total, matching the scope's
  list (context, skills, agents, commands, rules). No separate `workflows/` dir exists in the
  repo today (confirmed: `find commands rules -type f` → 2 files total, `commands/` is empty).
- `cmd/sh-agent/main.go:198-230` (`update` case) — the exact insertion point. `repo == ""` (a
  release binary) already branches to `harness.UpdateFromRelease` (line 218); a
  `--from-zip` flag would need to be checked even earlier (before or alongside this branch),
  since scope requires it to work "from a RELEASE binary" specifically, and today the checkout
  branch (`repo != ""`, line 227-230) calls `update()` → `harness.Reinstall`, which is also
  wrong for `--from-zip`. A `*zipFlag != ""` check likely wants to short-circuit both branches.
- `cmd/sh-agent/main.go:346-472` (`status`) — `repo := src.Dir` (347); when `repo == ""`
  (release binary, line 350-352) prints `source embedded vX.Y.Z`, no git state. For a DEV BUILD
  today, `repo` is always a real, resolvable checkout path (ResolveSource never returns a
  vanished Dir for a dev build in the current code) — so nothing today exercises "status after
  install from a now-deleted temp dir." Adding `--from-zip` for a dev build is out of scope
  (scope says the feature targets release binaries) but if `ResolveSource`/`status` share the
  Dir-source code path afterward, `Inspect(repo)` at line 354 would run `git -C <deleted-dir>
  rev-parse …`, which fails quietly into `RepoInfo{}` (IsCheckout returns false via `os.Lstat`
  error) — `status` would print `repo <deleted-path> (plain copy, not its own git checkout;
  update is refused)` (line 364), which is misleading (there is no such repo at all) and
  `update` would then be permanently refused for that home (matches `errPlainCopy`,
  repo.go:118-119) even though the operator wants to run `--from-zip` again. This is the
  strongest argument for a distinct `ZipPrefix` (`"zip:sha256:<hex>"` or similar) recorded in
  `State.Repo` instead of the literal temp path, paralleling `EmbeddedPrefix` exactly.
- `cmd/sh-agent/main.go:383` — `from := "commit " + st.Commit; if release, ok :=
  strings.CutPrefix(st.Repo, harness.EmbeddedPrefix); ok { from = "release " + release }`. A
  `ZipPrefix` would need a third branch here: `else if zipInfo, ok :=
  strings.CutPrefix(st.Repo, harness.ZipPrefix); ok { from = "zip " + zipInfo }` (or similar),
  printed instead of `commit ` + an empty/stale string.
- `cmd/sh-agent/main.go:553-580` (`knownCheckout`) — already has the pattern to copy: it
  special-cases `EmbeddedPrefix` (line 573) to avoid calling `samePath` against a release
  string. `knownCheckout` is called from the `install`/`update`-checkout branches
  (main.go:147, 227) — NOT from the release-binary `update` branch (line 198-226) — so if
  `--from-zip` reuses the release-binary branch's shape, `knownCheckout` likely does not need
  to change; if a `ZipPrefix`-recorded home is later re-resolved through `ResolveSource`'s
  dev-build walk (repo.go:75, which already special-cases `EmbeddedPrefix` to avoid treating it
  as a candidate path), that same `strings.HasPrefix(st.Repo, EmbeddedPrefix)` guard
  (repo.go:75) would ALSO need a `ZipPrefix` exclusion, or a stale zip manifest entry gets
  treated as a literal (deleted) candidate directory.
- `internal/harness/report.go:311-319` (`report()`'s own Source construction, mirrors
  main.go:304-319) and `report.go:345-358` — `report` does NOT read `st.Repo` from the loaded
  state at all (confirmed by grep: only `report.go:332` calls `LoadState`, and that result's
  `.Repo`/`.Commit` fields are never referenced — only `.Vendor`, `.Files`, `.InstalledAt`).
  Instead `report` re-resolves its OWN `src` independently (release binary → `EmbeddedSource`,
  else → fresh `ResolveSource` call) and uses `Inspect(src.Dir)` on that fresh source
  (report.go:354-358) to get `harness_commit`. **This means a stale `State.Repo` after a zip
  install does NOT directly break `report`'s output fields** — but if a future `report` call's
  own `ResolveSource` (for a dev build only; release binaries skip it, report.go:312-313) walks
  up from the working directory or falls back to `st.Repo` as a candidate path (repo.go:75-77),
  a `ZipPrefix`-tagged `st.Repo` must be excluded the same way `EmbeddedPrefix` already is,
  or `report` tries to `Inspect`/`LoadInventory` a directory that no longer exists.
- `README.md:182` and `.agents/docs/specs/sh-agent/delivery.md` § How a machine gets the
  harness — "Two ways, and no third" / "A plain copy… it is no longer one of the ways in… it
  is the fallback for a machine that can reach neither of the two, not a way in." Both need an
  ADR-driven rewrite if `--from-zip` becomes a real path for a release-binary place, not merely
  a code change.
- `README.md:9-45` (mermaid map, § How it works) and `cmd/sh-agent/main_test.go:1341-1373`
  (`TestHarnessMap`) — the test only checks that every `sh-agent <verb>` mentioned in the map
  exists in `usage` (main.go:35) and that `install`, `update`, `report` are all present
  (main_test.go:1362-1372); it does NOT parse flags, so `--from-zip` needs no new map NODE
  unless the feature adds a genuinely new LANE (e.g. "operator's own machine, GitHub web UI
  download" → temp dir → installer, bypassing both `harnessRepo` and `releases` nodes in the
  map, README.md:18-30). Per AGENTS.md rule 3, if this is a new path on the map, it is a design
  question requiring an ADR and a map edit in the same commit as the code.
- AGENTS.md rule 1 ("Two ways in, one way back, nothing else") and rule 6 ("The delivery claim
  has one home… stated once, in `.agents/docs/specs/sh-agent/delivery.md` § The model") — rule
  6 means the actual prose fix, if approved, is ONE edit to delivery.md § The model / § How a
  machine gets the harness; README.md:182, `.agents/docs/deploy/remote.md:56-61`, and
  `.agents/docs/deploy/local.md:45-55` all currently just point at / restate that same claim in
  one sentence each and would need the pointer sentence updated, not full rewrites (per rule 6's
  "point at it in one sentence and never restate it").

## Flows
- Existing release-binary `update`: `main.go` `update` case, `repo == ""` branch (198-226) →
  `harness.UpdateFromRelease(home, version, checkOnly, installArgs, out)` (`release.go:90`) →
  downloads+verifies signed release asset → replaces `~/.local/bin/sh-agent` → runs that new
  binary's `install`. This is the pattern `--from-zip` diverges from: no download, no binary
  replacement, no rebuild — just unzip → `Install`/`Plan` with the CURRENT running binary as
  `selfExe` (Plan's `selfExe []byte` param, plan.go:68,188-189) so `~/.local/bin/sh-agent`
  keeps being deployed as it is today (`installFiles`, main.go:277, not read in full this pass —
  worth a quick look before implementing, it is the shared helper both `install` branches use).
- Proposed `--from-zip` flow (not yet in code): validate flag/file → `os.MkdirTemp` under
  `ConfigDir(home)` or repo-analogous scratch (per project convention, never `/tmp` —
  `render.go`'s `diffOne` and `release.go` both already avoid `/tmp` for exactly this reason,
  apply.go:222-229) → unzip with a path-traversal guard (`filepath.Clean`/`strings.HasPrefix`
  check per entry — GitHub zips have one top-level folder to strip, e.g.
  `styrman-harness-main/context/...`) → `harness.IsRepo(tempdir)` sanity check (repo.go:99-105)
  → build `Source{Dir: tempdir, FS: os.DirFS(tempdir), Version: <running version>}` → call
  `harness.Plan` + `harness.Install` in-process exactly as `installFiles` does (main.go:277,
  not fully read) → `os.RemoveAll(tempdir)` → print a status line distinguishing this from both
  existing update paths.

## Conventions
- Every refusal in `repo.go` names WHERE a bad path/value came from (the `--repo`/
  `$SH_AGENT_REPO` error at repo.go:44) and is actionable; a `--from-zip` error (bad zip,
  traversal attempt, not-a-repo-shaped unpack) should match this voice.
- Scratch/temp files never go to `/tmp`; existing code uses `ConfigDir(home)`-relative temp
  dirs (apply.go:222-229 `Diff`, and `release.go`'s download directory at repo.go-adjacent
  `release.go:132-138` pattern) — `--from-zip`'s unpack dir should follow the same rule.
  `Reinstall` (repo.go:205) uses `.agents/scratch/sh-agent-build-<pid>` INSIDE the source repo
  itself — not applicable here since there is no source repo to nest inside for a release
  binary; the `ConfigDir(home)` pattern is the right model.
- `EmbeddedPrefix` (source.go:44) is the existing precedent for "manifest records something
  other than a live directory path" — a `ZipPrefix` constant alongside it is the natural shape
  the dispatch's own suggested `zip:sha256:<hex>` value fits (hashing the zip's bytes, not the
  temp dir, since the temp dir is deleted).
- Every state-machine-ish `if strings.HasPrefix(st.Repo, EmbeddedPrefix)` site found (repo.go:75,
  main.go:383, main.go:573) needs the equivalent `ZipPrefix` branch — three call sites, not one.

## Risks
1. **Design/doc conflict (highest risk).** delivery.md and AGENTS.md rule 1 currently say
   "two ways in… no third," with the unpacked-zip case explicitly demoted to "fallback… not a
   way in." Scope's own framing ("this path supersedes that [ADR-034] for the one file the
   operator names") suggests awareness of tension with ADR-034 specifically, but does not
   address the delivery-model "no third way" claim. Per AGENTS.md rule 3 this needs the
   operator's sign-off and an ADR BEFORE code, landed in one commit with docs + map (rule 3,
   rule 6). This is a recon finding to raise, not something to silently code around.
2. **Reinstall vs Install split.** It would be easy to accidentally route `--from-zip` through
   the existing checkout `update()` function (main.go:474) which unconditionally calls
   `harness.Reinstall` (rebuild) — scope explicitly forbids this. The new code must call
   `harness.Install`/`harness.Plan` directly, mirroring `installFiles` (main.go:277, not read
   this pass — read it before implementing).
3. **Manifest staleness.** Recording the literal temp dir path in `State.Repo` (as `Install`
   does today for any Dir source, apply.go:117-123) leaves `status`/`update`/`knownCheckout`
   pointing at a directory that no longer exists after cleanup — three call sites listed above
   need a `ZipPrefix` (or similar) to handle this, matching the `EmbeddedPrefix` precedent.
4. **Zip-slip.** No existing code in this repo unzips anything (grep found no `archive/zip`
   usage) — this is new territory; the entry-refusal logic (scope: "refuses any entry that
   would leave it") has no existing pattern to copy from within this codebase.
5. **`workflows/` doc-comment mismatch** (plan.go:61) is a pre-existing minor inaccuracy, not
   caused by this feature, but worth noting if editing plan.go nearby.

## Not examined
- `cmd/sh-agent/main.go` `installFiles` function body (main.go:277) — named but not read in
  full; read before implementing, since `--from-zip` should likely reuse it directly.
- `internal/harness/release.go` beyond lines 80-139 (download, checksum/signature verification,
  binary replacement) — not needed for `--from-zip` (no download/replace) but worth a skim for
  the temp-directory-cleanup idiom it already uses.
- `internal/harness/source.go` full file (only lines 18-54 read) — `DirSource`,
  `FS`/`os.DirFS` wiring, `src.path()` helper referenced in plan.go/inventory.go.
- `internal/harness/inventory.go` `checkManagedTree` full body (only referenced by line number,
  not read) — confirm exactly what it refuses (symlinks only, per plan.go:108's comment, or more).
- Existing tests: `internal/harness/*_test.go`, `cmd/sh-agent/main_test.go` other than
  `TestHarnessMap` — no test-naming conventions or existing `--check`/`--dry-run` test patterns
  surveyed; needed before writing the plan's test list.
- `.agents/docs/plans/2026-09-17-delivery-rework.md` (referenced by delivery.md as the PLANNED
  doc) — not read; may already discuss or foreclose a zip-install path.
- `skills/init-project`, `skills/find-skills` install-path mentions (dispatch item 5) — not
  grepped; likely low-risk (probably just point at README § Install without cloning) but
  unconfirmed.
- `context/` directory's own statements of the delivery model (dispatch item 5) — not grepped.
- Full `installPassArgs`/flag-parsing interaction for a new `--from-zip DIR` flag: whether it
  needs its own entry in `reorder()`'s `takesValue` map (main.go:622-627) — almost certainly
  yes (a flag taking a value) but not yet drafted/confirmed against `flag.FlagSet` registration
  style used for `--feedback-repo` etc. (main.go:84-90).

## Sizes (dispatch item 6)
`du -sb`: `context` 19,658 bytes; `skills` 852,476 bytes; `agents` 47,801 bytes; `commands` 0
bytes (empty dir); `rules` 0 bytes at top level but contains 2 files (`find` reports 2 files
under `rules/`, likely one level down or `du` miscounting an empty-looking top — worth
re-verifying with `find rules -type f -exec du -b {} +` before setting an unpack byte limit).
Total tracked files across all five dirs: 132. A zip of the whole repo (GitHub's "Download ZIP")
will be larger than these five dirs alone since it includes `.agents/`, `cmd/`, `internal/`,
`.git`-adjacent metadata is excluded by GitHub zips but everything else in the repo is included —
an unpack size/file-count limit should be set generously above ~1MB / 200 files, or scoped to
validate only that the five expected top-level dirs exist post-unpack rather than capping total
zip size tightly.
