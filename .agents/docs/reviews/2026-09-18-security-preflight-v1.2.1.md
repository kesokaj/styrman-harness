---
owner: security
date: 2026-09-18
---
# Security review — preflight v1.2.1   (2026-09-18)

Range reviewed: v1.2.0..288e52d

## Verdict: FINDINGS (three low; nothing that blocks the tag)

The trust-boundary change itself is **closed**. For a release binary `ResolveSource` returns
`EmbeddedSource(version)` before any disk lookup can run: the working-directory walk, the
`LoadState`/`st.Repo` candidate and `--repo`/`$SH_AGENT_REPO` are all on the other side of the
`IsReleaseVersion` branch, and `src.Dir == ""` propagates to every consumer — `repo := src.Dir`
gates `install`, `status`, `update`, `adopt`, and `Plan`/`LoadInventory`/`RenderContext` are
never called with a source that did not come from `ResolveSource` (verified: the only other
callers are `report.go:361/367`, which take `ReportRun.Source`, itself from `ResolveSource`).
No contrived `state.json` reaches a release binary at all, so `embedded:` prefixes, `..` and
symlinks in `repo` are moot on that path. The version cannot be forged at runtime: it comes
from an `-ldflags` var or `debug.ReadBuildInfo`, both fixed at build time, and the regexp is
`^v[0-9]+\.[0-9]+\.[0-9]+$` — anchored, ASCII-only, and Go's `$` does not admit a trailing
newline, so no crafted string gets in. Nothing an attacker controls at runtime (env, manifest,
cwd) flips a build between the two kinds. The three findings below are peripheral: a
defence-in-depth gap on the newly widened delete path, a reporting regression, and internal
identifiers in the trail.

## Findings

- **[low / confidence: medium]** `internal/harness/report.go:504` (`RemoveKeptClone`), newly
  reached from `cmd/sh-agent/main.go:206` — every guard in front of the `os.RemoveAll` except
  `os.Lstat` and `IsCheckout` asks **git in the operator's environment** (`gitOut` → plain
  `exec.Command("git", "-C", dir, …)`). This same file already treats that as a defeater:
  `SendReport` refuses to run at all when any of `gitRedirects` (`GIT_DIR`, `GIT_WORK_TREE`,
  `GIT_INDEX_FILE`, `GIT_OBJECT_DIRECTORY`, `GIT_COMMON_DIR`) is set, because "every git call
  below would commit and push in another repository whatever the checkout checks say". The
  delete path has no such refusal. With `GIT_DIR` alone set, `rev-parse --show-toplevel`
  answers with the cwd (`dir`, via `-C`), so `verifyFeedback` passes, and `remote.origin.url`,
  `status`, `stash list` and `rev-list --count --all --not --remotes` are all answered by the
  redirected repository rather than by the clone about to be deleted — i.e. the "clean, fully
  pushed" verdict can be rendered about something other than the directory that is removed.
  In practice the redirect usually makes `status` noisy and falls into the safe `leave()`
  branch, which is why this is low and not medium; it needs a redirect target whose worktree
  content matches `dir`. What makes it worth recording now is the **blast radius of the
  diff**: before this change a release-binary place never deleted anything on `update`, and
  those places (unpacked bundles, `/opt/devbox/...`, wrapper scripts) are exactly where a
  stray `GIT_*` is plausible. Remediation, one line and no new dependency: call the existing
  `gitRedirects` loop at the top of `RemoveKeptClone` and `leave()` when one is set (or give
  these specific `exec.Command`s a `cmd.Env` with the five variables stripped, as `gitOutC`
  already does for `LC_ALL`).
- **[low / confidence: high]** `cmd/sh-agent/main.go:284` (`report`) — `report` resolves its own
  source, so on a release binary the new refusal fires whenever `$SH_AGENT_REPO` is merely
  *exported in the environment* (a developer's shell profile, a wrapper). `report` then prints
  `report: this is the release binary …` and continues with `src = harness.Source{}`, which
  makes `BuildRecord` take the `src.FS == nil` branch: `harness_commit`, `harness_version`,
  `drift`, `skills` and `agents` all go null in the pushed record. That is the same failure
  *shape* as BUG-077 — a machine that reports nothing wrong while saying nothing useful — and
  it silently removes the drift signal the delivery model relies on. Not a confidentiality or
  integrity break (nothing wrong is written; the refusal text is not put into the record or
  the notes, only onto the local terminal). Remediation: in `report`, a release binary should
  use `EmbeddedSource` regardless of `--repo`/`$SH_AGENT_REPO` (report reads a source, it never
  *installs* one, so ADR-034's reason does not apply), or at minimum the note appended to the
  record should say "refused a named repo" rather than "no harness checkout".
- **[low / confidence: medium]** `.agents/docs/handoffs/2026-09-18-v1.2.0-released.md`,
  `.agents/docs/testing/2026-09-17-release-binary-e2e.md`,
  `.agents/docs/bugs/BUG-077-*.md`, `.agents/docs/changes/2026-09-18-release-binary-only-embedded.md`
  — no credentials, keys or tokens anywhere in the release diff (grepped; `git diff` scan for
  `api_key|secret|token|password|BEGIN … PRIVATE KEY|ghp_|AKIA|Bearer` is empty), but the trail
  now commits internal identifiers: the hostnames `spgo-mac.roam.internal`, `argolis-devbox`
  and `devbox`, the host path `/opt/devbox/deploy/harness`, a feedback-repo commit id
  (`c4b13a8`) and a verbatim `report.json` body (host, os, arch, vendor, file counts). This is
  fine while `styrman-harness` is private (the handoff records that it is), and it is exactly
  the evidence the bug records need, so the remediation is not redaction but a decision: if the
  repository is ever made public — and its own README documents a public-looking
  `go install …@vX.Y.Z` / `gh release download` path — these strings go out with it. Worth an
  explicit line in `remote.md` that publication requires a pass over `.agents/docs/` first.

## Answers to the six questions

1. **Is the rule closed?** Yes — see the verdict. No path was found by which a release binary
   reads a tree on disk. The one on-disk read that remains is `selfExecutable()` in
   `installFiles` (the running binary reads itself to install itself), which is pre-existing,
   out of range, and not a tree.
2. **Does the `--repo` refusal disclose anything?** No. The echoed value is the operator's own
   argument or their own exported `$SH_AGENT_REPO`, printed to the local terminal (stderr via
   `main`, or stdout on the `report` path) and to nothing else; it does not enter the record
   pushed to the feedback repository. One low-confidence caveat: the echo is unscrubbed, so a
   value that happened to be a credential-bearing URL (`https://user:token@host/…`) would be
   reprinted verbatim — `scrubCredentials` exists in `report.go` for exactly that shape and is
   not applied here. The dev-build message (`%s is not a harness checkout`) has the same
   property and predates this diff.
3. **Is the dev/release split forgeable?** No. `buildVersion()` reads the `-ldflags` `version`
   var, else `debug.ReadBuildInfo()` (`Main.Version`, `vcs.revision`, `vcs.modified`) — all
   sealed at build time. No environment variable, manifest field or file influences it; the
   manifest's `embedded:` value is read only by `status` and `knownCheckout`, never by
   `ResolveSource` on the release branch. Someone who builds their own tree at a local `vX.Y.Z`
   tag gets a "release binary" that installs *that tree's own* embedded files — no trust
   boundary is crossed, since the binary's own content was already trusted as the binary.
   Worth noting for the docs only: the code comment and ADR-034 say the embedded files are
   "signed and verified on the way in", which is true of the `UpdateFromRelease` download path
   (`SHA256SUMS` + `ssh-keygen -Y verify`, D26–D28), not of anything checked at install time.
4. **`shedKeptClone` in the binary-only `update`.** It can remove exactly one path,
   `filepath.Abs(home)/.config/sh-agent/feedback` — no interpolation, no argument reaches the
   path. `home` is `--home` or `os.UserHomeDir()`. The refusals hold: `os.Lstat` (not `Stat`)
   leaves a symlink where it is; a non-checkout, a checkout whose toplevel or `remote.origin.url`
   is not the feedback repository, a dirty tree (`--untracked-files=all --ignored`), a stash, a
   non-zero `rev-list --count --all --not --remotes`, and any git failure in any of those all
   `leave()` rather than delete. `--check` is guarded explicitly at `main.go:205` and pinned by
   `TestReleaseBinaryUpdateShedsKeptClone`, which asserts the clone survives `--check` and is
   removed (with the message) on a real run. Two residual notes: the environment gap in
   finding 1, and the ordering — `shedKeptClone` runs *before* `UpdateFromRelease`, so the
   removal happens even when the update then refuses (no `gh`, offline, a non-release or
   backwards tag). The second is by design (`update` is not a read-only mode) and the deletion
   is fully guarded, but it is a side effect on a failure path and is worth one sentence in the
   change entry.
5. **`internal/harness/report.go` message changes.** **Premise mismatch, on the record:**
   `report.go` has **no changes in v1.2.0..288e52d** (`git log v1.2.0..288e52d --
   internal/harness/report.go` is empty; the tree is clean). Those two message changes are
   BUG-074's, already inside the v1.2.0 tag and outside this gate's range. Reviewed the current
   text anyway: no decision path sits in either string; the `leave()` reason names `dir` in
   prose only, and the settling command deliberately does **not** interpolate the path
   (`git log --oneline --all --not --remotes FETCH_HEAD` is literal, with the comment "a path
   with a space or a shell metacharacter would ride into text the operator pastes"). The count
   `%s` is a `rev-list --count` output. Nothing new can carry a metacharacter into pasteable
   text.
6. **Every release-diff path.** All fourteen read end to end; see `## Not examined` for the
   single qualification. The trail markdown carries no secrets; the instruction-shaped text in
   it (ADR-034's decision, the BUG prevention rules, the e2e command blocks) is
   operator/agent-authored in-repo content, not untrusted input, and the only commands an agent
   might replay are `gh release download`, `sha256sum -c`, `ssh-keygen -Y verify` and read-only
   `git log`/`rev-list` — none writes to a remote.

## Checked

- **Secrets** — `git diff v1.2.0..288e52d` scanned for `api[_-]?key|secret|token|password|
  BEGIN … PRIVATE KEY|ghp_|AKIA|Bearer`: no match. No `.env`, `*.pem`, `*.key`, `*.tfvars` or
  credential JSON in the range. Internal hostnames/paths in the trail: finding 3. `verifyFeedback`
  still refuses to print `remote.origin.url` ("it can carry a token") and `scrubCredentials`
  still covers the URLs that are printed.
- **Injection** — no SQL, no shell string, no template. Every external value in the range is
  either a path handed to `filepath.Join`/`os.Lstat`/`exec.Command` as a discrete argv element
  (never a shell), or a message operand. The new refusal strings interpolate `version` (build-
  time) and the operator's own `--repo`/`$SH_AGENT_REPO`. `UpdateFromRelease` still validates
  the GitHub-supplied tag against `IsReleaseVersion` before it reaches argv (a `-`-leading tag
  would be read as a flag) — unchanged and still correct with the new regexp usage.
- **Authz/authn** — no endpoints, no sessions. The authorisation-shaped checks in range are
  `knownCheckout` (unchanged, and now unreachable for a release binary since `repo == ""`
  short-circuits before it) and the `RemoveKeptClone` guard chain (finding 1). The ordering
  change in `ResolveSource` moves the release check *earlier*, which is the fail-closed
  direction.
- **Crypto** — nothing hand-rolled in the range and nothing crypto-adjacent added. The release
  verification path (`ssh-keygen -Y verify`, `sha256sum`, the compiled-in signer line,
  `/usr/bin/ssh-keygen` before `PATH`) is untouched by this diff and is a vetted external
  verifier, not a home-grown construction. `IsReleaseVersion` is a version regexp, not a
  security primitive, and is not being asked to authenticate anything.
- **Dependencies** — `go.mod` unchanged, still zero non-stdlib requires (`module
  github.com/kesokaj/styrman-harness`, `go 1.24.4`, no `toolchain` per ADR-029). `govulncheck
  ./...` on go1.27.1: **"No vulnerabilities found."** `npm audit` not applicable — no
  `package.json` anywhere in the repository (npm is present on the box; there is no JS/TS
  project to audit).
- **GCP surface** — none in range: no IAM binding, service account, Cloud Run flag, Secret
  Manager reference or key file. The only GCP-adjacent strings are the internal hostnames in
  finding 3 (`argolis-devbox` names an Argolis-style environment, which is why that finding
  exists at all).

## Not examined

- **`internal/harness/release.go` beyond `UpdateFromRelease`'s first 40 lines** — the download,
  checksum and signature verification body was not re-read. It is unchanged in this range
  (v1.2.0 content) and was gated by the 2026-09-17 and 2026-09-18 reviews, but it *is* blast
  radius for the reordered `shedKeptClone` call and I did not re-verify it here.
- **`internal/harness/apply.go` / `plan.go` / `inventory.go` / `render.go` bodies** — confirmed
  by grep that they take `Source` and by `apply.go:96` that an embedded install records
  `EmbeddedPrefix+src.Version`, but the file bodies were not read end to end. Unchanged in
  range.
- **`report.go` in full** — read lines 340–395 and 480–625 (record building, `RemoveKeptClone`,
  `verifyFeedback`, `gitRedirects`, `SendReport`'s first lines). `Collect`, the redaction
  pipeline, `cloneThrowaway` and `commitAndPush` were not read in this run. Unchanged in range.
- **No runtime verification was performed.** No `sh-agent` binary was built or executed, no
  install was run against a scratch home, `report` was not run in any mode, and no scratch
  probe repository was created. Every statement above is from reading the source and the diff;
  finding 1's redirect behaviour in particular is reasoned from git's documented precedence,
  not measured — hence medium confidence.
- **The v1.2.0 tag's own contents** were not re-reviewed; this report covers v1.2.0..288e52d
  only. The `report.go` messages named in question 5 sit inside v1.2.0 and were read but not
  gated here.
- **Untracked file `.agents/docs/testing/2026-09-18-verify-preflight-v1.2.1.md`** (present in
  the working tree at review time, not in the range) was not read.

## Outcome

Read against the diff by the main agent, 2026-09-18. All three fixed before the tag.

1. **[low] `RemoveKeptClone` asked git without refusing the redirect variables**, on a path that DELETES and that this release newly extends to binary-only places — while `SendReport` refuses them a few lines away for exactly that reason. Fixed with the same `gitRedirects` loop: any of them set and the clone is kept, with the variable named. `TestRemoveKeptClone` gained a subtest; removing the loop fails it.
2. **[low] an exported `$SH_AGENT_REPO` made `report` send a null record on a release binary.** This was a regression introduced by ADR-034 an hour earlier: `report` went through the same refusal `install` owes, caught the error, and reported with `harness_version`, `drift`, `skills` and `agents` all null — the useless-record shape BUG-077 was about. Fixed: `report` READS a source and never installs one, so a release binary now takes `EmbeddedSource` directly, and a `--repo` passed to `report` is answered with a note rather than a refusal. `TestReportOnReleaseBinaryIgnoresRepoFlags` covers it.
3. **[low] the trail commits internal identifiers** — hostnames, machine paths, a report body — fine inside a private repository, a decision outside one. Fixed by writing the check into `.agents/docs/deploy/remote.md` § Releases, where publication would be considered: no secret is in them, that is verified every release, but publication is a separate decision from releasing.

The verdict that matters is the one this release was cut for, and it was reached independently rather than taken from the dispatch: the trust boundary is closed — `ResolveSource` returns the embedded source before any disk lookup, `src.Dir == ""` propagates to every consumer, no caller takes a source from elsewhere, and the version cannot be forged at runtime.

Premise correction accepted: `internal/harness/report.go` has no changes in `v1.2.0..288e52d`; the two message fixes the dispatch asked about shipped inside v1.2.0. That was an error in my dispatch text, not in the review.
