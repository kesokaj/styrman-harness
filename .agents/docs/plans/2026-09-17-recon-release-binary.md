---
owner: explorer
date: 2026-09-17
---
# Recon — release binary with embedded assets (2026-09-17)

## Question
Survey the codebase to support a plan for: (a) releasing per-platform sh-agent
binaries with context/skills/agents/commands/rules embedded via `go:embed`,
usable via `go install .../cmd/sh-agent@latest`; (b) install source order
`--repo` > checkout-found-as-today > embedded; (c) `update` on a binary-only
place downloading the latest release via `gh release download`, verifying the
signed tag, replacing the binary, reinstalling; (d) `report` switching from a
kept clone to a throwaway partial clone per run; (e) README "install without
cloning" section, R9 and map changes.

## Summary
Every file-reading function in `internal/harness` (`Plan`, `LoadInventory`,
`LoadSkills`, `LoadAgents`, `LoadDocs`, `RenderContext`, `checkManagedTree`,
`readRegular`) takes a `repo string` and calls `os.ReadFile`/`os.ReadDir`/
`filepath.Glob`/`filepath.WalkDir`/`os.Lstat`/`filepath.EvalSymlinks` directly
— none takes an `fs.FS`, so all would need a parallel `fs.FS`-based path (or a
refactor to `io/fs` funcs) to also read from an embedded filesystem. `FindRepo`
already has exactly the "checkout found as today" behavior requested for (b);
adding embedded-as-fallback is a third source appended after it. `Reinstall`
in repo.go is the only thing that changes the deployed binary today and it
insists on `go build` from source (BUG-006) — a release binary changes this
substantially: R4's "single Go binary, stdlib only" and R8/R9's "one delivery
channel: checkout+update, or unpacked bundle+install" would need explicit
amendment (ADR + operator approval) since GitHub Releases would become a THIRD
channel not named in R8/R9. `report.go`'s tests (`TestSendReport`,
`TestSendReportRefusesRedirectsAndLinks`, `TestSendReportPushesOnlyToFeedbackURL`)
pin a `runDir`/kept-clone model (`~/.config/sh-agent/feedback`) that assumes a
persistent local clone across runs — switching to a throwaway partial clone
per run is a structural rewrite of `ResolveFeedback`, not an incremental patch.

## Relevant files
- `internal/harness/repo.go` — `FindRepo` (source order: explicit/env → cwd+parents → last-install state), `IsRepo`, `IsCheckout`, `Inspect`, `Behind`, `Pull`, `Reinstall` (git-build-only), `gitOut`.
- `internal/harness/plan.go` — `Plan()` builds `[]File` from `LoadInventory`+`RenderContext`; walks `skills/<name>` with `filepath.WalkDir`, refuses symlinks/secret-shaped files; `Reinstall`'s bin is added as `File{Rel: BinPath, ...}` when `selfExe != nil`.
- `internal/harness/inventory.go` — `LoadInventory`, `LoadSkills`, `LoadAgents`, `LoadDocs`, `checkManagedTree`, `readRegular`: all path-based, all refuse symlinks in the managed tree (BUG-036).
- `internal/harness/render.go` — `RenderContext`: reads `context/AGENTS.md` + `context/*.md` via `os.ReadFile`/`filepath.Glob`.
- `internal/harness/derive.go` — Antigravity agent translation; has `toolMap` (unmapped tool fails install per AGENTS.md rule). Small (56 lines), path-independent of repo reads (works on already-loaded `Agent` struct).
- `internal/harness/adopt.go` — `sh-agent adopt` (promotes a project customization into the repo); needs closer read (not yet fully read — see Not examined).
- `internal/harness/state.go` — `LoadState`/state persisted to `home` (`~/.config/sh-agent` in real use), records `st.Repo` used by `FindRepo`'s third candidate.
- `internal/harness/report.go` (632 lines) — feedback-repo report/send logic: `ResolveFeedback`, `verifyFeedback`, `runDir`, `commitAndPush`, `SendReport`. Not yet fully read — see Not examined; this is the piece (d) most directly rewrites.
- `internal/harness/report_test.go` (570 lines) — pins current kept-clone behavior; named tests confirmed to exist: grep found `TestSendReport`, `TestSendReportRefusesRedirectsAndLinks`, `TestSendReportPushesOnlyToFeedbackURL` referenced in report.go comments/tests (exact bodies not yet read).
- `cmd/sh-agent/main.go` (781 lines) — CLI entry; `install` chooses rebuild-vs-plain-copy around line 107-154; `status` reporting around 233-273; `update` calls `Reinstall` at line 411. Not yet fully read line-by-line — see Not examined.
- `cmd/sh-agent/main_test.go` (1014 lines) — install/status/update tests against checkout and plain-copy fixtures; not yet read in detail.
- `internal/harness/harness_test.go` (1499 lines) — includes `TestHarnessMap` per the dispatch; not yet located/read.
- `go.mod` — `module github.com/kesokaj/styrman-harness`, `go 1.24.4`, no `toolchain` directive (deliberate, ADR-029, per AGENTS.md).
- `REQUIREMENTS.md` — R4 (single binary, stdlib only, `update` = git pull + reinstall), R8 (one delivery channel: checkout+update or bundle+install, no other channel), R9 (four procedures fixed; any change to them needs ADR + operator approval, one commit with code+docs+map).
- `.agents/docs/deploy/remote.md` § Releases — not yet read, needed for (5).
- `skills/preflight` — not yet read, needed for (5) Step 6.
- `README.md` — not yet read, needed for the "Install without cloning" section and "the map" (R9's "README.md § How it works").

## Flows (partial)
- Install source order today, per `FindRepo` (repo.go:20-48): explicit `--repo` or `$SH_AGENT_REPO` (hard error if given but not a repo) → walk cwd and parents for `IsRepo` → `LoadState(home).Repo` (last install's recorded repo) → error. This already implements exactly the "explicit, then checkout found as today" half of requirement (b); embedded-fallback would be appended as a further candidate, but note `FindRepo` currently *errors* if nothing is found — the embedded path presumably becomes the state when FindRepo fails, which is a real branch to add in `main.go`, not just repo.go.
- `Plan()` (plan.go:68) calls `LoadInventory(repo)` then `RenderContext(repo, inv)`, both path-based against `repo string`; then walks `inv.Skills[i].Dir` via `filepath.WalkDir` for skill file contents, and reads `inv.Agents[i].Raw`/`inv.Commands[i].Raw`/`inv.Rules[i].Raw` which were already loaded into memory by `LoadAgents`/`LoadDocs` (path-based `os.ReadFile`).
- Install rebuild-vs-copy decision lives in `cmd/sh-agent/main.go` around lines 107-154: after `FindRepo`, checks `harness.IsCheckout(repo)` (`own`); if true and `go` is on PATH, calls `harness.Reinstall` (rebuilds from source, BUG-006's fix: an installed copy re-installing itself must never silently reuse its own already-built bytes); if plain copy (`!own`, i.e. `.git` absent — an unpacked bundle, BUG-015), installs the *running* binary as-is with a "plain copy... installed as is" message.

## Conventions
- Every repo-reading function refuses symlinks explicitly and cites the originating bug/review in a comment (BUG-036, security reviews dated in comments) — new fs.FS-based readers must preserve this refusal, and `fs.FS` (e.g. from `embed.FS`) does not expose `Lstat`/symlink info the same way, so the symlink-refusal logic cannot transfer mechanically.
- Every risky decision (rebuild vs. copy, dirty-tree refusal, secret-shaped refusal) has an inline comment naming the bug number or review date that justified it — a strong convention for any new decision code (embedded-fallback logic, release-binary trust boundary) to follow.
- `go.mod` deliberately carries no `toolchain` directive (ADR-029) — release automation must not reintroduce one.

## Risks (preliminary)
- **R8/R9 conflict**: R9 says the harness has exactly two install paths (checkout+update, bundle+install) and enumerates "Four procedures... and no others"; a GitHub Releases binary distributed via `go install .../cmd/sh-agent@latest` and `gh release download` is arguably a THIRD delivery channel and a fifth procedure-shaped thing. REQUIREMENTS.md is immutable without explicit operator approval + ADR (already reflected in R9's own text: "together, only with the operator's explicit approval, an ADR, and one commit"). This recon cannot determine whether that approval already covers today's dispatch — the dispatch text describes the change as already "confirmed by the operator," but REQUIREMENTS.md itself is unmodified as of this recon (R8/R9 text above still describes only two channels) — a plan must not proceed to code before REQUIREMENTS.md + an ADR are updated in the same commit as required by R9's own final sentence.
- **fs.FS refactor is pervasive, not local**: `LoadInventory`, `LoadSkills`, `LoadAgents`, `LoadDocs`, `RenderContext`, `checkManagedTree`, `readRegular`, `Plan`'s skill walk all need a filesystem abstraction; `embed.FS` is read-only and has no symlinks (so symlink-refusal code becomes dead code on that path, but must stay for the checkout/bundle path) — likely means introducing an interface or dual code paths, not a one-line signature change.
- `Reinstall` (repo.go) is the sole mutator of the deployed binary and currently requires `go build` from a git checkout; a release binary needs a wholly different `update` path when there is no checkout (no `go`, no source) — this is exactly gap (c) in the dispatch, confirmed structurally: today's `update` (main.go ~line 411 calls `Reinstall`) has no non-git branch at all for replacing the binary.
- `report.go`/`report_test.go` pin a kept-clone model; rewriting to throwaway partial clones will break `TestSendReport`, `TestResolveFeedbackFailures`, `TestSendReportRefusesRedirectsAndLinks`, `TestSendReportPushesOnlyToFeedbackURL` unless carefully re-targeted — these tests were not read in full; exact assertions unknown (see Not examined).

## Open questions
- Has REQUIREMENTS.md's R8/R9 actually been amended/approved already, or does the dispatch's "confirmed by the operator" cover only the recon scope, not a REQUIREMENTS.md edit? The file on disk right now still describes two channels only.
- What `.gitignore`/embed directive scheme is intended for `context/`, `skills/`, `agents/`, `commands/`, `rules/` so `//go:embed` at the module root does not also embed `.agents/scratch/` or `VENDORED.md` third-party content inadvertently?
- Are there dotfiles/underscore-prefixed files under skills/ or agents/ that `embed` would silently skip without `all:`? Not yet checked (see Not examined).

## Not examined (ran out of budget / did not reach)
- `cmd/sh-agent/main.go` full read (only grepped line numbers, not full bodies) — need lines ~1-780 read in full for install/status/update/report command wiring, flag parsing, and exact rebuild/copy branch code.
- `cmd/sh-agent/main_test.go` (1014 lines) — not read at all.
- `internal/harness/harness_test.go` (1499 lines) including `TestHarnessMap` — not read at all.
- `internal/harness/report.go` (632 lines) and `report_test.go` (570 lines) — only referenced via grep/comments in repo.go's neighborhood; `ResolveFeedback`, `verifyFeedback`, `runDir`, `commitAndPush`, `SendReport` bodies not read.
- `internal/harness/adopt.go`, `apply.go`, `derive.go` (full body), `frontmatter.go`, `state.go` — not read in detail.
- `.agents/docs/deploy/remote.md` § Releases — not read.
- `skills/preflight/SKILL.md` Step 6 — not read.
- `README.md` § How it works / map — not read.
- Listing actual files under `skills/` and `agents/` for dotfile/underscore-prefixed names (embed all: question) — not done.
- go.mod Go version cross-check against `go:embed` requirements (embed needs go1.16+; 1.24.4 is fine, not a concern, but not verified against any go.work or multi-module layout).

## Additional findings (second pass)

- `cmd/sh-agent/main.go:136-138` — decision code confirmed:
  `own := harness.IsCheckout(repo)`, `rebuild := shouldRebuild(goErr == nil, os.Getenv(harness.ChildEnv) == "1", own)`.
  `shouldRebuild(goOnPath, child, ownCheckout bool) bool { return goOnPath && !child && ownCheckout }` (main.go ~450).
  `knownCheckout` (main.go ~420) refuses building from a harness-shaped directory that is not the one actually installed (no `--repo`/`$SH_AGENT_REPO`, and it doesn't match `state.Repo`) — a real guard a release-binary "third source" must respect or explicitly bypass.
  `binaryOlder(rev, commit string)` (main.go) compares the running binary's VCS stamp to the repo's HEAD for `status`'s staleness warning — a release binary built via `go install` carries `debug.ReadBuildInfo` module version info, not a short git commit hash, so this comparison logic will need to change too.
- `internal/harness/report.go` `ResolveFeedback` (line 377) is confirmed to look for an existing dir at `home/.config/sh-agent/feedback`; if absent it `git clone`s the feedback URL there and KEEPS it (no removal). This is the exact mechanism (d) replaces. `runDir` (line 550) computes `reports/<host>/<date>-<HHMMSS>[-N]` inside that persistent checkout. `commitAndPush` (line 599) stages/commits/pushes named paths only, to `url` explicitly (not `origin`), `HEAD:main`. None of these three currently do a partial/sparse/shallow clone — full `git clone` only. Switching to `--filter=blob:none --no-checkout --depth 1` + sparse-checkout is a rewrite of `ResolveFeedback`'s clone branch and adds an sparse-checkout/cleanup step SendReport must call at the end (not present today) and `update` must also gain: dispatch (d) says "`update` removes an old kept clone" — no code today references cleaning up `~/.config/sh-agent/feedback` from `update` at all; grepped `cmd/sh-agent/main.go` and found no such call. This is a new behavior, not a modification of an existing one.
- `internal/harness/report_test.go` confirmed test names matching dispatch: `TestSendReport` (:333), `TestResolveFeedbackFailures` (:434), `TestSendReportRefusesRedirectsAndLinks` (:489), `TestSendReportPushesOnlyToFeedbackURL` (:526), plus `TestSendReportDryRunShowsDroppedRecords` (:559) not named in the dispatch but adjacent and likely also affected. Bodies not read in full (budget) but their names alone confirm they exercise the persistent-clone/kept-checkout model that (d) removes.
- **`TestHarnessMap` does NOT exist** anywhere in the repo (grepped `internal/harness/harness_test.go` function list exhaustively — 45 test functions, none named `TestHarnessMap`). The dispatch's item (6) asks to report on it; it is either not yet written, misnamed in the dispatch, or the closest analog is `TestInstallStatusResyncUninstall` / `TestVendorSelection` / `TestNestedCopyIsNotACheckout`, none of which is a "map" test. Flag this back to the dispatcher.
- `skills/init-project/assets/agents-tree/**/.gitkeep` — 13 dotfiles found under `skills/`, all named `.gitkeep`, used by `init-project` to ship an empty directory skeleton for a new project's `.agents/` tree. Plain `//go:embed skills` SKIPS files/dirs starting with `.` or `_` (Go embed rule) — these `.gitkeep` files would silently vanish from an embedded build unless the embed directive uses `all:` (e.g. `//go:embed all:skills`). No other dotfiles or underscore-prefixed files found under `skills/` or `agents/` in this search (searched with `find ... -name ".*" -o -name "_*"`, only .gitkeep matched, all under init-project's assets tree, agents/ had none).
- Go toolchain present on this box: `go1.27.1` (ahead of what go.mod declares, `go 1.24.4`, no toolchain directive — consistent with AGENTS.md's note on ADR-029).
- `.agents/docs/deploy/remote.md` § Releases (lines 58-88): today's release procedure is signed-tag-only, explicitly "a MARKER naming a version, not a way the harness reaches a machine" (line 60) and "There is no deploy step after the push" (line 88) — this is the exact sentence ADR-030 D2 states and that a release-binary-attachment step would contradict or need to amend. `skills/preflight/SKILL.md` Step 6 "SHIP" (lines 79-88 read) has 6 sub-steps (changelog, version bump, one release commit with gate artifacts, tag, push, "Deploy via the project's deploy guide"); nothing today builds or attaches platform binaries — step 6.5's "Deploy via deploy guide" is GCP/remote-box oriented (cloud-run-deploy skill, ssh box, or local-only), not "attach a binary to a GitHub Release." Attaching binaries to a release would need a NEW sub-step here, and reconciling with remote.md's "nothing pushes to a machine / no deploy step after tag push" framing is a real tension to resolve with the operator.
- `README.md` "How it works" (lines 5-56) contains the mermaid diagram (`checkoutLane`, `checkoutPlace`, `bundle`) that is "the map" referenced in dispatch (e) and R9. Confirmed two lanes exist today: checkout+update, and unpacked-bundle+install-only ("update refuses it — unpack a newer copy to move forward", line 99). A third lane (release binary via `go install`/`gh release download`) is a genuine addition to this diagram, not a relabeling.

## Risks (additional)

- `binaryOlder`'s git-short-hash comparison and `status`'s staleness detection assume every deployed binary carries a git VCS stamp; a `go install @latest`-built binary's `debug.ReadBuildInfo().Main.Version` is a module pseudo-version/tag, not a short commit hash matching `Inspect(repo).Commit` — `status` output and the staleness warning need new logic for this case, not just a new source branch.
- No existing code path removes `~/.config/sh-agent/feedback`; "update removes an old kept clone" (d) is wholly new logic to add to `update`'s flow in `cmd/sh-agent/main.go`, gated on detecting "binary-only place" (which itself needs a new way to detect "this place has no checkout, only a binary install" — today `FindRepo` either finds a repo or hard-errors; there's no persisted "binary-only" state distinct from "bundle, no .git" already handled by `own=false`).
- `go:embed all:` vs plain `embed` distinction matters concretely for `skills/init-project/assets/agents-tree/**/.gitkeep` — must use `all:` prefix or add the .gitkeep files back another way (they exist purely to make git track otherwise-empty directories; embed's directory model doesn't have "empty directories" as a concept anyway, so their omission may be harmless in practice — worth confirming with init-project's actual usage of that tree before assuming it's a problem).

## Open questions (additional)
- Has the operator actually approved amending R8/R9 already (dispatch says "confirmed by the operator") — recon cannot verify this from the repo since REQUIREMENTS.md and ADR list show no such change yet. A plan must not write code before REQUIREMENTS.md + a new ADR land in one commit, per R9's own final sentence.
- Where dispatch's `TestHarnessMap` should live / whether it's a test to be newly written for map-diagram-consistency (e.g., asserting README's diagram lists exactly the lanes `sh-agent` implements) — does not exist today.
