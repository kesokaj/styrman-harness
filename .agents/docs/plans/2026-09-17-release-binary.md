---
owner: main
date: 2026-09-17
---
# Plan — release binary with the harness built in; `report` keeps nothing   (2026-09-17)

## Why

The operator wants a place to get the harness by downloading one `sh-agent` binary from a GitHub release — "Just downloading the binary makes a lot of sense. Because it will only work If I logged in to github or have the ssh key" — instead of cloning the whole repository, whose `.agents/docs` trail then sits beside their projects and can confuse a session opened there. They also found that `sh-agent report` keeps a full clone of `styrman-harness-feedback` on every place, so each machine holds every other machine's reports ("But why can I see dave reports?"). Scope confirmed in chat on 2026-09-17 with "go all, throwaway clone for 12". Recon: `.agents/docs/plans/2026-09-17-recon-release-binary.md`.

## Scope

In (the operator's numbered list):
1. The release binary replaces the unpacked bundle as R9's second way in; a checkout stays the first. R9 edit on exact text, ADR.
2. `context/`, `skills/`, `agents/`, `commands/`, `rules/` embedded with Go `embed` from a package at the module root.
3. `install` source order: `--repo`, a checkout found as today, then the embedded files.
4. `update` on a binary-only place: newest release via `gh release download`, verified, binary replaced, reinstalled. (Verification is by signed checksums, D5 — the scope line said "the SSH-signed tag".)
5. Platforms linux/amd64, darwin/arm64, linux/arm64, cross-compiled locally in `/preflight`, attached to the release.
6. A binary-only place changes only when a release is cut.
7. The first download documented.
8. macOS quarantine documented, for a browser download only.
10. The unpacked bundle leaves the docs and the map; its code path stays.
11. README "Install without cloning": `gh`, `go install` over SSH, the browser; PATH.
13. Every release carries release notes grouped as New features, Improvements and optimisation, Fixes, and Other updates (D19).
12. `report` makes a throwaway partial clone holding only this place's folder, deletes it after the push; `update` removes an old kept clone.

Out: 9 — a download page, a timer, Windows. Apple code signing. Changing `adopt`, `uninstall` or `diff` beyond reading from a source.

Follow-ups noted, not built: the `close-map-gaps` branch reads `~/.config/sh-agent/feedback`; when it lands it keeps its own full clone at `~/.config/sh-agent/feedback-read` (D12).

## Surfaces touched

- New `embed.go` at the module root (package `styrman`), `//go:embed all:context all:skills all:agents all:commands all:rules`.
- `internal/harness/plan.go`, `inventory.go`, `render.go`, `frontmatter.go`, `derive.go`, `apply.go`, `state.go`: readers move from `repo string` to a `Source` (Interfaces).
- `internal/harness/repo.go`: `FindRepo`, `IsCheckout`, `Inspect`, `Reinstall`; new release code in `internal/harness/release.go`.
- `internal/harness/report.go` and `report_test.go`: `ResolveFeedback`, `SendReport`, `commitAndPush`, the four clone-pinning tests.
- `cmd/sh-agent/main.go` and `main_test.go`: source selection, `update`, `status`, `version`; `TestHarnessMap` must still pass.
- `REQUIREMENTS.md` R8 (first sentence) and R9; `.agents/docs/architecture/DECISIONS.md` ADR-033; `README.md` § How it works (map) and a new § Install without cloning; `AGENTS.md` § How the harness works rules 1 and 2 (D18); `.agents/docs/deploy/local.md` (§ Update, § From an unpacked bundle, § Report back), `.agents/docs/deploy/remote.md` § Releases and § The four places; `.agents/docs/specs/sh-agent/delivery.md`; `skills/preflight/SKILL.md` Step 6.

## Interfaces

**R9 — proposed exact replacement of the Install, Update, Report back and Release bullets (operator approval required, Phase 1):**

- **Install:** a place either clones `styrman-harness` with the operator's own credential and runs `sh-agent install`, or downloads a release's `sh-agent` binary with that credential — `gh release download`, `go install` over SSH, or the Releases page — and runs `sh-agent install`, which installs the harness built into the binary. Either way installed files go into the vendor directories and no repository is created.
- **Update:** the operator edits, verifies, commits and pushes `main` in their checkout. A checkout place runs `sh-agent update` when the operator decides — fetch, refuse a dirty checkout, fast-forward, rebuild, reinstall. A binary-only place runs `sh-agent update`, which downloads the newest release, verifies its signed checksums, replaces the binary and reinstalls; it changes only when a release is cut. Nothing is pushed to a machine.
- **Report back:** each place runs `sh-agent report`, which sends a redacted copy of `~/.config/sh-agent/` and a metrics record to `styrman-harness-feedback` through a throwaway clone that holds only that place's folder and is deleted after the push; the harness is improved from the reports in the operator's checkout, and that change returns through Update.
- **Release:** `/preflight vX.Y.Z`, from the operator's own turn, gates the tree, pushes an SSH-signed tag on `main`, and publishes that tag's binaries with checksums signed by the same key. Checkout places follow `main`; binary-only places update to the newest release.

**R8 — proposed exact replacement of its first sentence (operator approval required, Phase 1):** "The harness reaches a machine one way: a checkout the operator owns, or a release's `sh-agent` binary with the harness built in, each installed with `sh-agent install` and kept current with `sh-agent update`." The rest of R8 is unchanged.

**Source.** `type Source struct { FS fs.FS; Dir string; Version string }` (`Version` set for the embedded files only, D20). `Dir` is the checkout or plain-copy path (`FS` is `os.DirFS(Dir)`), or `""` for the embedded files (`FS` is `styrman.Files`). Every reader takes a `Source`; symlink refusal (BUG-036) stays a `os.Lstat` on `filepath.Join(Dir, name)` when `Dir != ""` — `embed` cannot hold a link, so the embedded source needs none. `state.json` records `"repo": "<dir>"` or `"repo": "embedded:<version>"`.

**Source selection** (`install`, `status`, `diff`, `list`, `render`, `update`, `report`): an explicit `--repo` or `$SH_AGENT_REPO` is used or it is an error, exactly as today — never a fall-through (D15). Otherwise the working directory or a parent, then the recorded repo when it is still a harness tree, then the embedded files — but the embedded files only in a release binary, one whose version is a `vX.Y.Z` tag (release ldflags or `go install …@vX.Y.Z`); a development build with no tree found keeps today's `harness repo not found …` error (D14). `adopt` requires a checkout and says so.

**`ResolveSource` replaces `FindRepo`** (`internal/harness/repo.go:20`): `func ResolveSource(explicit, home, version string) (Source, error)` — `explicit` is `--repo` (else `$SH_AGENT_REPO`, read inside as today); `version` is the running binary's version; the embedded fallback applies only when `version` matches `^v\d+\.\d+\.\d+$`. Every git helper — `IsCheckout`, `Inspect`, `Behind`, `Pull`, `Reinstall` — is called only with `src.Dir != ""`; with an embedded source none is called (D17).

**Per command with an embedded source (`src.Dir == ""`)** (D17):

| Command | Behaviour |
|---|---|
| `install`, `sync` | install from `styrman.Files`; copy the running binary (D13) |
| `status` | first line `source     embedded vX.Y.Z (release binary; update downloads the newest release)` in place of the `repo` line; `binary` line shows the version; drift compared against the embedded files |
| `diff`, `list`, `render` | read the embedded files |
| `update` | the release flow (Interfaces § `update` on a binary-only place); with `--check`, only step 1: `up to date (vX.Y.Z)` exit 0, or `newer release vA.B.C available` exit 2 |
| `adopt` | `adopt needs a harness checkout; this install is a release binary`, exit 1 |
| `uninstall`, `version` | unchanged — neither reads the source |
| `report` | `BuildRecord` takes the `Source`; `harness_commit` is `null` and a new `harness_version` field carries `vX.Y.Z`; `drift`, `skills`, `agents` computed from the embedded files |

With a checkout or plain-copy source every command behaves exactly as today.

**Binary handling on install** (`cmd/sh-agent/main.go` today decides rebuild vs copy with `IsCheckout(repo)`): with the embedded source the running binary is copied to `~/.local/bin/sh-agent` and never rebuilt, and `IsCheckout` is never called with an empty path; a checkout or plain copy keeps today's rule (BUG-006, BUG-015) (D13).

**Version.** Release builds: `-trimpath -ldflags "-X main.version=vX.Y.Z"`, `CGO_ENABLED=0`. A `go install …@vX.Y.Z` binary: `debug.ReadBuildInfo().Main.Version`. Otherwise today's VCS stamp. `sh-agent version` prints it; `status` prints `source embedded vX.Y.Z` for an embedded install. `status` reads the INSTALLED file's version from its module version when that is a `vX.Y.Z` tag, because `-trimpath` keeps `-ldflags` out of the build info (D21).

**Release assets** (per tag): `sh-agent-linux-amd64`, `sh-agent-darwin-arm64`, `sh-agent-linux-arm64`, `SHA256SUMS` (`sha256sum` format), `SHA256SUMS.sig` = `ssh-keygen -Y sign -f <operator key> -n sh-agent-release SHA256SUMS`.

**`update` on a binary-only place** (no checkout found by the source order): requires `gh` on PATH, else `update: this install has no checkout; install gh and run gh auth login, or download a newer release by hand (README § Install without cloning)`, exit 1.
1. `gh release view --repo kesokaj/styrman-harness --json tagName` → latest tag; equal to this binary's version → `up to date (vX.Y.Z)`, exit 0.
2. `gh release download <tag> --pattern sh-agent-<GOOS>-<GOARCH> --pattern SHA256SUMS --pattern SHA256SUMS.sig --dir <tmp under ~/.config/sh-agent/>`.
3. `ssh-keygen -Y verify -f <allowed-signers file written from the signer line compiled into the binary> -I <signer identity> -n sh-agent-release -s SHA256SUMS.sig < SHA256SUMS`; then the binary's SHA-256 against its `SHA256SUMS` line (`crypto/sha256`). Any failure: nothing replaced, message names the check, exit 1.
4. Replace `~/.local/bin/sh-agent` by rename from the verified file (same filesystem), then run the new binary's `install`; the temp directory is removed on every path.

**Report throwaway clone.** `--feedback-repo DIR` keeps today's behaviour exactly (a persistent, verified checkout the operator names). Without it:
1. `git clone --quiet --filter=blob:none --no-checkout --depth 1 <FeedbackURL> ~/.config/sh-agent/report-<pid>`; an empty repository (no commits) is cloned plainly.
2. `verifyFeedback` on it; `git sparse-checkout set --no-cone reports/<host>/`; `git checkout`.
3. Write the run directory, stage, commit and push exactly as today (`push -u <FeedbackURL> HEAD:main`).
4. `os.RemoveAll` of the clone on every exit path (deferred). A rejected push: remove, re-clone, re-apply the same run once; a second rejection → `push rejected twice; re-run sh-agent report — nothing is kept locally`, exit 1.
`--dry-run` clones nothing, as today.

**Release notes** (D19): `/preflight` Step 6 writes `.agents/scratch/release-notes-vX.Y.Z.md` from the `changes/` entries dated after the previous tag's commit, one line per entry under `## New features` (Type feature), `## Improvements and optimisation` (Type task or refactor that changes behaviour), `## Fixes` (Type bugfix, with its BUG id), `## Other updates` (rules, skills, docs, context), plus `## Install` pointing at README § Install without cloning; an empty group is omitted. The operator reads it before the upload; `gh release create vX.Y.Z --notes-file <it>` attaches it. The same text goes into the release's `changes/` entry.

**Old kept clone cleanup** (`update`, both kinds of place): `~/.config/sh-agent/feedback` is removed when it is a git checkout whose origin is `FeedbackURL`, with a clean `git status --porcelain` and no commits ahead of its upstream; the output says `removed the kept feedback clone ~/.config/sh-agent/feedback (report no longer keeps one)`. Otherwise it is left and named with the reason.

## Failure modes

| When | What the operator sees |
|---|---|
| `install` with no tree found | a release binary installs its embedded files; a development build prints today's `harness repo not found …` (D14) |
| `--repo` or `$SH_AGENT_REPO` naming something that is not a harness tree | today's error, no fall-through to the embedded files (D15) |
| `update`, binary-only, no `gh` | the gh message above, exit 1, nothing changed |
| `update`, `gh` not logged in or no repo access | `update: gh release view failed: <gh's first line>`, exit 1 |
| No release asset for this platform | `update: release vX.Y.Z has no sh-agent-<os>-<arch>`, exit 1 |
| Signature or checksum mismatch | `update: SHA256SUMS signature does not verify` / `checksum mismatch for sh-agent-<os>-<arch>`, exit 1, binary untouched |
| `ssh-keygen` missing | `update: ssh-keygen not found; cannot verify the release`, exit 1 |
| `report`, partial clone unsupported by the server | falls back to a plain `--depth 1 --no-checkout` clone, output names it |
| `report`, push rejected twice | message above, exit 1, no directory left |
| A kept feedback clone with local work | left in place, named with `has uncommitted changes` or `has N unpushed commits` |

## Verification

Tiers: unit, integration (local bare repositories, fake `gh` on PATH), e2e (a real release).

- Unit: readers against `fstest.MapFS`; the embedded source lists every skill directory and all 15 `.gitkeep` files; `Plan` from the embedded source equals `Plan` from the checkout (paths and bytes) at the same commit.
- Integration: install from the embedded source into a scratch home → `status` drift none, same file set as from the checkout; `update` binary-only path against a fake `gh` script serving a binary, `SHA256SUMS` and a signature made with a test key (the test compiles in that key's signer line), for good signature, bad signature, bad checksum, missing asset, missing `gh`; `report` against a local bare repository holding two hosts' reports: the run lands, the clone directory is gone afterwards, the other host's report blobs were never in the clone (checked before removal through a test hook), a rejected push retries once; old-clone cleanup removes a clean clone and leaves a dirty one.
- E2E criteria:
  - Golden path: after `/preflight v1.1.0`, on the work Mac, the README's `gh` one-liner downloads the binary and `sh-agent install` installs the harness with no clone; `sh-agent status` shows `source embedded v1.1.0`, drift none; `sh-agent report` sends a run and leaves no `~/.config/sh-agent/feedback` or `report-*` directory.
  - Edge: on that place, `sh-agent update` with no newer release prints `up to date (v1.1.0)`; with the signature file of a different release swapped in (a scratch copy of the assets, served by the fake `gh`), it refuses and leaves the binary untouched.

## Phases

- [x] Cold-reader check — three rounds, `.agents/docs/reviews/2026-09-17-coldread-release-binary.md`; never reached both tokens, ended at the cap. Blocking items G1–G7 resolved by D13–D17; round 3's G8 resolved as D18. At the cap the operator chose "Proceed to implementation", 2026-09-17

### Phase 1 — decisions (no code)

- [x] Operator approves the R8 sentence in Interfaces (added by cold read round 1) — verify: the operator's words quoted here, who: user — expect: approval in their own turn. Evidence: "Approve as written", the operator's own turn, 2026-09-17, on the exact sentence quoted in Interfaces
- [x] Operator approves the R9 text in Interfaces, verbatim, and D5 (signed checksums instead of a signed tag) — verify: the operator's words quoted here, who: user — expect: approval in their own turn. Evidence: "Plan approved", the operator's own turn, 2026-09-17, given on this plan with the R9 text and D5 as written; the cold-reader check was still running, so no code followed
- [x] ADR-033 drafted from D1–D12 by the main agent, read against this plan — verify: read, who: AI — expect: every decision row represented, no claim beyond them. Evidence: ADR-033 appended to `.agents/docs/architecture/DECISIONS.md` by the main agent from D1–D18 and read against this plan: every D-row represented or named in the plan it points to; the supersession of ADR-031 D3 and ADR-030 D2 stated (D18); no claim beyond the plan, 2026-09-17
- [x] Branch `release-binary` created from `main`; this plan and ADR-033 committed there — verify: `git log --oneline -1` on the branch, who: AI — expect: the commit on the branch, `main` untouched. Evidence: branch `release-binary` created from `main` at 22b92ff; the plan, the recon brief, the three-round cold read, ADR-033 and the ROADMAP pointer committed on it; `main` untouched, 2026-09-17

### Phase 2 — the embedded source

- [x] `embed.go` at the module root and a test that every file under the five directories is in `styrman.Files`, the 15 `.gitkeep` files included — verify: unit, `go test ./...`, who: AI — expect: fails before `embed.go` exists, passes after. Evidence: `embed_test.go` `TestFilesMatchTree` failed with 133 `not embedded` errors against an `embed.go` holding an empty `Files`, passed after `//go:embed all:context all:skills all:agents all:commands all:rules`; `go list` reports 133 embedded files, equal to the files under the five directories, 2026-09-17
- [x] `Source` type; every reader takes it; link refusal kept for `Dir != ""` — verify: unit against `fstest.MapFS` plus the existing symlink tests, who: AI — expect: all existing tests pass unchanged in meaning. Evidence: `internal/harness/source.go` (`Source`, `DirSource`, `readRegular` on a Source); `LoadInventory`, `LoadSkills`, `LoadAgents`, `LoadDocs`, `RenderContext`, `Plan` take a `Source`, `checkManagedTree` runs only for `Dir != ""`; callers in `main.go` and `report.go` pass `DirSource(repo)` until task 4. `source_test.go` failed to build before `Source` existed, then passed: a MapFS plan (paths, bytes, 0444 → 0644), nine refusals from a MapFS, a non-regular SKILL.md, DirSource and MapFS planning identically, a linked `skills/<name>` followed and a dangling one refused. The existing tests changed only `repo` → `DirSource(repo)`; with the link guards in `checkManagedTree` and `readRegular` disabled, four symlink tests failed, restored they pass. `gofmt`, `go vet`, `go test -race -count=1 ./...` clean; `install --home .agents/scratch/source-home` then `status`: drift none, 2026-09-17
- [x] `Plan` from embedded equals `Plan` from checkout — verify: unit, who: AI — expect: identical file lists and bytes. Evidence: `plan_test.go` (package `styrman_test`) `TestEmbeddedPlanEqualsCheckout` compares paths, bytes, modes and KeepForeign for all three vendors, passes; must-fail seen twice: `skills` embedded without `all:` → 52 files planned from the tree only; `skills/golang/SKILL.md` made 0755 → a mode mismatch on both vendors (embed records 0444, so an executable skill file would lose its execute bit from a release binary — the test is the alarm for that); both restored, the file's original 0664 included, 2026-09-17
- [x] `ResolveSource`, the per-command table (install, status, diff, list, render, adopt, report record), `state.json` `embedded:<version>`, `version` from ldflags or build info — verify: integration, install from embedded into `.agents/scratch/embed-home`, `status`, who: AI — expect: drift none, `source embedded` — docs: `README.md`. Evidence: `ResolveSource` replaces `FindRepo` (`repo.go`), `EmbeddedSource`, `IsReleaseVersion`, `EmbeddedPrefix` (`source.go`); `Install` takes a `Source` and records `embedded:<version>` with no commit and no git call; `BuildRecord` takes a `Source`, new `harness_version` field; `main.go`: `var version` (ldflags), install copies the running binary for an embedded source, `status` prints `source     embedded vX.Y.Z (…)` and `from release vX.Y.Z`, `diff`/`list`/`render` read the source, `adopt` refuses with the table's message, `update` refuses (placeholder until Phase 4, D22), `releaseFromInfo` for the installed file (D21). Tests: `TestResolveSource`, `TestResolveSourceRefusesNamedNonRepo`, `TestEmbeddedInstallStateAndRecord`, `TestReleaseBinaryWithoutCheckout`, `TestReleaseFromInfo`; must-fail by mutation, each restored: the release gate removed → `TestResolveSource` and `TestReleaseBinaryWithoutCheckout` fail; the `embedded:` skip removed → `TestResolveSource` fails; an explicit non-tree falling through → `TestResolveSourceRefusesNamedNonRepo` and `TestReleaseBinaryWithoutCheckout` fail. Integration: `CGO_ENABLED=0 go build -trimpath -ldflags "-X main.version=v0.0.0"`, run from `/` with `SH_AGENT_REPO` unset: `install --home .agents/scratch/embed-home` 260 files; `status` `source     embedded v0.0.0`, `from release v0.0.0`, drift none, exit 0 (its `binary` line showed the revision, the test build being from a dirty tree); `state.json` `"repo": "embedded:v0.0.0"`, `"commit": ""`; `update` refused, exit 1; a checkout install into `.agents/scratch/embed-cmp` has the same file set, bytes and modes. Not exercised: a `go install …@vX.Y.Z` binary (no tag with the embed exists yet), the `binary` line of a real clean-tag build. README deferred to Phase 5's docs task, which lands on `main` in one commit (D11), 2026-09-17
- [x] Checkpoint — gofmt, vet, `go test -race -count=1 ./...` clean; verifier log; branch commit and push. Evidence: verifier PASS `.agents/docs/testing/2026-09-17-release-binary-phase-2.md` (81 pass, 0 fail, 0 skip at 21ba8dc; checkout install and a release-shaped embedded install both drift none; identical file sets); `reviewer` APPROVE with 0 findings `.agents/docs/reviews/2026-09-17-review-f7ee817-21ba8dc-code.md`; `security` 1 medium and 6 low `.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md` — findings 1, 2 and the non-hermetic test fixed (`TestEmbedHoldsOnlyTrackedFiles`, `knownCheckout`'s `embedded:` refusal with `TestKnownCheckoutRefusesEmbeddedManifest`, the temp-directory probe in `TestResolveSource`), five rebutted in writing in that report's `## Outcome`, and the sixth (should a release binary refuse a home a checkout installed?) raised to the operator; changes entry `.agents/docs/changes/2026-09-17-release-binary-phase-2.md`, 2026-09-17

### Phase 3 — `report` keeps nothing

- [x] Throwaway partial clone, sparse `reports/<host>/`, deferred removal, retry once on rejection; `--feedback-repo` unchanged — verify: integration against a two-host bare repository, who: AI — expect: run lands, clone gone, other host's blobs absent before removal, retry observed — docs: `.agents/docs/deploy/local.md`. Evidence: `cloneThrowaway` (`--filter=blob:none --no-checkout --depth 1`, `verifyFeedback`, `sparse-checkout set --no-cone reports/<host>/`, `checkout` skipped for an empty repository), `writeAndPush`, `errPushRejected`, and `SendReport`'s two-attempt loop with the removal deferred inside each attempt; `--feedback-repo` keeps one attempt, the dirty-file notice and today's rejection message. `TestSendReportThrowawayClone` (run lands on the remote, clone gone, no kept checkout, no `report-*` left; through `testHooks.beforeRemove`: the other host's directory never checked out and its canary never in a local blob, read with `GIT_NO_LAZY_FETCH=1`) and `TestSendReportRetriesRejectedPush` (one rejection retried and re-applied, two rejections → `push rejected twice … nothing is kept locally`, nothing left). Mechanics settled first in `.agents/scratch/p3-probe/probe.sh`: a `file://` remote ignores `--filter` unless it sets `uploadpack.allowFilter`, and `cat-file` on a promised object lazily FETCHES it — both would have made the isolation check vacuous. 2026-09-17
- [x] The four clone-pinning tests rewritten for the throwaway model, each must-fail against the old code — verify: unit/integration, who: AI — expect: fail on HEAD's `report.go`, pass after. Evidence: `TestSendReport` now drives the default path (dry run touches nothing, first run into an empty repository, second run's `-2` directory, `assertKeepsNothing` after each) and reads what was sent from the ORIGIN, not from a local checkout; the named-checkout half moved to `TestSendReportNamedCheckout` (written into, never removed, stray file left alone, rejected push keeps the commit); `TestResolveFeedbackFailures` drops the two default-path cases the model removes and gains the unreachable-URL case (nothing left behind); `TestSendReportRefusesRedirectsAndLinks`, `TestSendReportPushesOnlyToFeedbackURL` and `TestSendReportDryRunShowsDroppedRecords` unchanged in meaning. Must-fail: all five new/rewritten tests failed to BUILD against the old `ResolveFeedback` signature, so the guards were then mutated one at a time instead (`.agents/scratch/p3-mustfail/mutate.py`, which asserts each target is unique and restores from a saved copy, per BUG-068): removal not deferred, `--filter` dropped, `sparse-checkout` dropped, retry removed, and the dirty check dropped — each failed exactly the test that covers it, and the file was restored. 2026-09-17
- [x] Old kept clone cleanup in `update` — verify: integration, clean clone removed, dirty one kept and named, who: AI — expect: as stated — docs: `.agents/docs/deploy/local.md`. Evidence: `RemoveKeptClone` removes `~/.config/sh-agent/feedback` only when it is a checkout, `verifyFeedback` passes, `status --porcelain --untracked-files=all` is empty and nothing is ahead of upstream; every other case is left in place and named. Called at the top of `update` in `cmd/sh-agent/main.go`, reported and never fatal. `TestRemoveKeptClone`: clean removed and said so, dirty kept ("has uncommitted changes"), one unpushed commit kept ("has 1 unpushed commit(s)"), another repository kept ("not the feedback repository"), nothing there says nothing. `.agents/docs/deploy/local.md` § Report back rewritten for the throwaway model, including the two rejection paths and what `update` does with an old clone. 2026-09-17
- [ ] Checkpoint — as Phase 2; `security` over the phase (credential-bearing directory, pushes to a remote)

### Phase 4 — `update` from a release

- [ ] `internal/harness/release.go`: latest tag, download, `ssh-keygen -Y verify`, checksum, rename, reinstall — verify: integration with a fake `gh` and a test key, the five cases in Verification, who: AI — expect: all five behave as the Failure modes table says
- [ ] Signer line compiled in from `.agents/docs/deploy/remote.md` § Releases' allowed-signers line — verify: unit, who: AI — expect: the constant equals that line
- [ ] Checkpoint — as Phase 2; `security` over the phase (signature verification, binary replacement)

### Phase 5 — release, docs, map (lands on `main` as ONE commit, R9)

- [ ] `skills/preflight/SKILL.md` Step 6 and `.agents/docs/deploy/remote.md` § Releases: cross-compile the three binaries, `SHA256SUMS`, `SHA256SUMS.sig`, `gh release create vX.Y.Z` with the assets — verify: `/review --docs`, who: AI — expect: exact commands, the signer check before upload — docs: `.agents/docs/deploy/remote.md`
- [ ] Release notes in `/preflight` Step 6 and `remote.md` § Releases (Interfaces § Release notes) — verify: `/review --docs`, and a dry build of the notes for v1.1.0 from today's `changes/` entries, read by the operator, who: both — expect: four groups as specified, no empty group — docs: `.agents/docs/deploy/remote.md`
- [ ] README § Install without cloning (the operator-approved draft of 2026-09-17: `gh`, `go install` over SSH, the browser with `xattr`, PATH), § How it works map (the binary as the second way in, the bundle gone, report's throwaway clone), `AGENTS.md` rule 1, `local.md`, `remote.md` § The four places, `delivery.md` — verify: `/review --docs`, who: AI — expect: no living doc describes a kept feedback clone or the bundle as a way in — docs: `README.md`
- [ ] `AGENTS.md` § Setup commands names `LoadSkills` as the `os` advisory's call path; since Phase 2 the readers go through `os.DirFS` — re-state the path from a `govulncheck` run on an older Go, or drop the function name — verify: read, who: AI — expect: the note names a path that exists
- [ ] `/preflight` Step 6 refuses a release when `git status --porcelain --ignored --untracked-files=all -- context skills agents commands rules` names anything, and builds the binaries from the tag's own tree; `.agents/docs/deploy/remote.md` § Releases says the same (D23) — verify: `/review --docs` plus one refused dry run against a planted ignored file, who: AI — expect: the refusal names the file — docs: `.agents/docs/deploy/remote.md`
- [ ] R8 and R9 edits as approved; `TestHarnessMap` passes; ADR-033 final — verify: unit, who: AI — expect: pass
- [ ] Operator approves the redrawn map (R9) — verify: quoted, who: user — expect: approval
- [ ] Merge `release-binary` into `main` as one commit (`git merge --squash`), changes entry, push — verify: `git show --stat HEAD`, who: AI — expect: code, docs, map, R9, ADR in one commit
- [ ] E2E: `/preflight v1.1.0` from the operator's turn; on the work Mac the golden path and edge in Verification — verify: e2e, who: both — expect: as written; log `.agents/docs/testing/<date>-release-binary-e2e.md`
- [ ] Checkpoint and Done — verifier, reviewer over the whole range, security, documentation, bugs, handoff

## Risks

Assumptions:
- GitHub serves `--filter=blob:none` partial clones over SSH and accepts a push from a `--depth 1` clone — SETTLED 2026-09-17 by one real `sh-agent report` from this box, authorised by the operator (D25): the blob-less `--no-checkout --depth 1` clone over SSH succeeded, `sparse-checkout` limited it to `reports/dave/`, and the push landed as `09586ff` (`reports/dave/2026-09-17-200615`, five files). The clone was removed and `~/.config/sh-agent/` holds no `report-*`. The push stalled on the network for about four minutes before completing — slow, not refused. Log: `.agents/docs/testing/2026-09-17-release-binary-phase-3-real-report.md`. Note: a local bare repository ignores `--filter` unless it sets `uploadpack.allowFilter`, so the local tiers seed that.
- `go install github.com/kesokaj/styrman-harness/cmd/sh-agent@latest` resolves the newest semver tag: `v1.0.0` has no embed, so `@latest` gives a binary that cannot install until `v1.1.0` exists — stated in the README section; checked after `v1.1.0` with `go install …@latest` in a scratch `GOPATH`.
- `ssh-keygen -Y verify` exists on macOS and Debian (OpenSSH ≥ 8.1) — checked with `ssh-keygen -Y 2>&1 | head -1` on this box and the work Mac, first task of Phase 4.

## Decisions

D1 — 2026-09-17 — The release binary replaces the unpacked bundle as the second way in; R9's four bullets change to the exact text in Interfaces, subject to approval — by: the operator, "go all" on scope items 1 and 10
D2 — 2026-09-17 — Embedding from a root package `styrman` with `all:` so the 15 `.gitkeep` files are included — by: implementer, from recon risk 3 (`embed` skips dotfiles without `all:`)
D3 — 2026-09-17 — Readers take a `Source{FS, Dir}`; link refusal stays an `os.Lstat` on the disk path when `Dir` is set — by: implementer, from recon: every reader takes a path today and BUG-036 must not regress
D4 — 2026-09-17 — Source order `--repo`, env, working directory, recorded path, embedded — by: the operator, scope item 3, with the recorded path kept from today
D5 — 2026-09-17 — `update` verifies `SHA256SUMS.sig` (`ssh-keygen -Y verify`, namespace `sh-agent-release`) and each binary's checksum, not the git tag: a tag's signature covers the commit, never a release asset — needs the operator's approval since scope item 4 said "the SSH-signed tag" — by: implementer, flagged for the operator
D6 — 2026-09-17 — Version from `-ldflags -X main.version` for releases, `debug.ReadBuildInfo` for `go install`, the VCS stamp otherwise — by: implementer
D7 — 2026-09-17 — Binary-only `update` shells out to `gh` and `ssh-keygen`; no Go dependency added (R4) — by: the operator, scope item 4, and R4
D8 — 2026-09-17 — `report` without `--feedback-repo` uses a throwaway partial clone under `~/.config/sh-agent/report-<pid>`; a rejected push is retried once from a fresh clone — by: the operator, "throwaway clone for 12"; the retry by implementer, since a throwaway clone cannot keep a rejected commit
D9 — 2026-09-17 — `update` removes a kept `~/.config/sh-agent/feedback` only when clean, origin-verified and not ahead — by: the operator, scope item 12
D10 — 2026-09-17 — Three platforms, `CGO_ENABLED=0 -trimpath`, built and uploaded inside `/preflight` Step 6 — by: the operator, scope item 5
D11 — 2026-09-17 — Work happens on branch `release-binary`, one squash commit to `main` with code, docs, map, R9 and ADR-033 (R9's one-commit rule) — by: implementer, from R9
D12 — 2026-09-17 — The `close-map-gaps` branch, when it lands, keeps its reading clone at `~/.config/sh-agent/feedback-read`, since D9 removes `~/.config/sh-agent/feedback` — by: implementer
D13 — 2026-09-17 — An install from the embedded source copies the running binary and never rebuilds; `IsCheckout` is never asked about an empty path — by: cold read round 1, blocking 1
D14 — 2026-09-17 — The embedded source is a fallback only for a release binary (version a `vX.Y.Z` tag); a development build outside a tree keeps today's error, so a stale embedded copy is never installed silently — by: cold read round 1, blocking 2
D15 — 2026-09-17 — An explicit `--repo` or `$SH_AGENT_REPO` never falls through to the embedded files — by: cold read round 1, blocking 4, keeping today's deliberate behaviour
D16 — 2026-09-17 — R8's first sentence changes with R9, on its own exact text and approval — by: cold read round 1, blocking 3
D17 — 2026-09-17 — `ResolveSource` replaces `FindRepo`; git helpers are called only for a tree on disk; every command's behaviour with an embedded source is fixed by the table in Interfaces, `update` branching to the release flow on `src.Dir == ""` — by: cold read round 2, blocking G6 and G7
D18 — 2026-09-17 — The rule "a tag delivers nothing" changes with R9, in three places, in the Phase 5 commit: `AGENTS.md` § How the harness works rule 2 becomes "**A tag is a release.** `/preflight vX.Y.Z` makes a signed tag and publishes that tag's binaries with signed checksums; checkout places follow `main`, binary-only places update to the newest release."; `.agents/docs/deploy/remote.md` § Releases' opening sentence says a tag marks a version and carries the release binaries that binary-only places update to, while checkout places still follow `main`; ADR-031 is not edited (append-only) — ADR-033 states it supersedes ADR-031 D3's "a tag delivers nothing". The `Source selection` sentence also names `update` (round 3, non-blocking G9) — by: cold read round 3, blocking G8, resolved at the cap
D19 — 2026-09-17 — Every release carries release notes in four groups built from the `changes/` entries since the previous tag, read by the operator and attached with `--notes-file` — by: the operator, "We need to write a release note when we do a new version. Like updates, optimization and fixes, new features etc.", added after the cold read; a docs-and-procedure addition that changes no interface the cold read settled
D20 — 2026-09-17 — `Source` carries `Version`, set only for the embedded files: `Install` records `embedded:<version>`, `status` and the report record show it, and no reader has to be passed the running binary's version separately — by: implementer, Phase 2 task 4
D21 — 2026-09-17 — `status` reads an installed release binary's version from its build info's module version when that is a `vX.Y.Z` tag (`releaseFromInfo`); `-ldflags` cannot be read back because `-trimpath` leaves it out of the build info (measured with go1.27.1: `go version -m` shows `-ldflags` without `-trimpath`, none with it). The checkout path keeps `versionFromInfo`, so a development build at a clean tag is still compared by revision — by: implementer, Phase 2 task 4
D22 — 2026-09-17 — Until Phase 4, `update` with an embedded source refuses with `this install is the release binary vX.Y.Z and has no checkout; download a newer release by hand (README § Install without cloning)`, exit 1, calling no git helper; Phase 4 replaces it with the release flow — by: implementer, so the branch never calls git with an empty path between phases
D23 — 2026-09-17 — `go:embed` reads the WORKING TREE, so a release must be built from a tree with nothing untracked or ignored under the five directories: `embed_test.go` pins the embedded set to `git ls-files` at test time, and `/preflight` Step 6 refuses at release time (Phase 5) — by: implementer, from the Phase 2 security review's finding 1, reproduced with an ignored `rules/.env` that was embedded while every other test stayed green
D24 — 2026-09-17 — A release binary's install REPLACES an install this home got from a checkout: it installs the embedded files, prunes what is no longer planned and rewrites the manifest to `embedded:vX.Y.Z`, without the `knownCheckout` guard. Converting a place to binary-only needs no flag; the cost is that running a downloaded binary from outside the checkout on the operator's own box swaps that home's install — by: the operator, "Replace, as it does now", their own turn, on the Phase 2 security review's finding 6
D25 — 2026-09-17 — ONE real `sh-agent report` from this box into `styrman-harness-feedback` is authorised, to settle whether GitHub serves a `--filter=blob:none` partial clone over SSH and accepts a push from a `--depth 1` clone (plan § Risks). It runs only after Phase 3's local-bare-repository tests pass, writes only `reports/<this host>/`, and its result is recorded in the Phase 3 testing log — by: the operator, "Yes, when Phase 3's code is ready", their own turn
Defaults for round 1's non-blocking items: the checksum file is signed in `/preflight` Step 6 right after the tag's own verification, with the same key; the compiled-in signer line is a package variable a test overrides; `go install` is documented for `v1.1.0` or newer — by: implementer
Rejected: embedding in `cmd/sh-agent` — `go:embed` cannot reach `../../skills`
Rejected: verifying the git tag only — it does not cover the downloaded file (D5)
Rejected: the GitHub API for `report` — would make `gh` required for reporting (the operator chose the throwaway clone)
Rejected: a Go SSH-signature library — a dependency against R4; `ssh-keygen` is the vetted verifier already on every place
Open: whether `go install` should be a documented way in before `v1.1.0` exists — default: the README names `v1.1.0` or newer explicitly
