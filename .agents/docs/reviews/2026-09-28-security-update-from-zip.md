---
owner: security
date: 2026-09-28
---
# Security review — update-from-zip   (2026-09-28)

Range reviewed: 1cb2608..worktree (`git diff 1cb2608` plus the 12 untracked files in `git status --short --untracked-files=all`)

## Status: COMPLETE

## Verdict: FINDINGS

The unpacker's path handling is sound. I found no zip-slip route on Linux or macOS, and duplicates, a file and a directory with the same name, and case-folding collisions on APFS all fail closed. The worst problem is the trust model, not the parser. `update --from-zip` puts an unsigned archive into the global instruction set of two agents that run shell commands with the operator's credentials. A release binary, whose value is signature verification, now installs that content without checking it. Nothing binds the archive to this repository either: any zip with `context/AGENTS.md` and `skills/` is accepted. ADR-038 accepts this risk explicitly and README states it honestly. I record it at medium so the gate sees it, but it is not a new defect in the code. The remaining findings are low. The byte limits leave directory creation unbounded. The recorded sha256 is computed through a second open of the file, not the one that is unpacked. The manifest value switches off `status`'s content comparison, reopening the stale-content class of incident that ADR-034 was written for. Several texts now state security claims that the code no longer keeps.

## Findings

- **[medium / confidence: high]** `cmd/sh-agent/main.go` `updateFromZip`, `README.md` § Install from a zip, ADR-038 "Rejected" bullet 4: **the unsigned zip becomes agent instructions.**
  - What is exposed: the installed `context/`, `skills/`, `agents/`, `commands/` and `rules/` are what every Claude Code and Antigravity session on the machine reads and acts on. A substituted or wrong zip (a fork, a look-alike repository, an old `styrman-harness-main (3).zip` in Downloads) therefore amounts to arbitrary command execution through the agent.
  - What checks it: nothing. The only content test is `zipRoot`/`IsRepo`, which checks shape (`context/AGENTS.md` plus `skills/`), not origin.
  - The printed and recorded `sha256` is an identifier, not a verification: GitHub publishes no hash for "Download ZIP" or "Source code (zip)" archives, and `SHA256SUMS.sig` covers binaries only.
  - This is also the only way in where a signed release binary installs content its compiled-in signer never vouched for (ADR-034 superseded in part).
  - The ADR's premise, "a zip is unsigned", is a release-process choice, not a constraint.
  - Remediation, no new dependency:
    - `/preflight` publishes a `git archive` zip of the tag as a release asset listed in `SHA256SUMS`, which is already signed.
    - `update --from-zip` verifies the zip's `crypto/sha256` against that list when the operator passes the sums and signature files beside it (e.g. `--sums SHA256SUMS`), through the existing `ssh-keygen -Y verify` path in `release.go` and its compiled-in signer line.
    - For an unverified zip, such as "Download ZIP" on `main`, the command prints an explicit `UNVERIFIED` line in its own output, not only in README.
  - If the operator reaffirms ADR-038's acceptance, this does not block the gate.
- **[low / confidence: medium]** `internal/harness/zip.go:85-89` (directory entries) and `:116` (`MkdirAll` of each file's parent): **the byte limits do not bound directory creation.**
  - `maxZipFileBytes` and `maxZipTotalBytes` count file content only, and nothing limits a name's length or depth.
  - Each of the 5,000 permitted entries can create a chain of new directories up to PATH_MAX. On Linux (4096 bytes) that is about 2,000 levels of one-character names; on macOS (1024) about 500.
  - In principle one archive of a few tens of MB can create millions of directories (inodes plus one block each) under `~/.config/sh-agent/zip-*` before any limit fires, and the cleanup `RemoveAll` then has to walk all of them.
  - Reasoned from `os.MkdirAll` semantics, not executed.
  - Remediation: refuse names longer than a small bound (for example 512 bytes) or deeper than a fixed element count (the real trees are fewer than 10 deep), checked in `zipNameSafe` so the refusal happens before anything is created.
- **[low / confidence: high]** `cmd/sh-agent/main.go` `updateFromZip` → `fileSHA256(zipPath)`, then `harness.UnpackZip(zipPath, …)` → `zip.OpenReader`: **hash and content are not bound (TOCTOU).**
  - The file is opened twice. If it changes in between (a browser still finishing the download, a sync client replacing it), the `zip:sha256:<hex>` in the manifest and the output describes bytes that were never installed.
  - Remediation: open once. Inside the 0700 unpack directory, copy the archive through `io.TeeReader` into a `sha256` hash, then unpack from that private copy with `zip.OpenReader`, or hash it with `io.NewSectionReader` over the same `*os.File` that `zip.NewReader(f, size)` reads. The recorded hash then names what was unpacked.
- **[low / confidence: high]** `cmd/sh-agent/main.go` `status` (`fromZip` block) and `ResolveSource`/`knownCheckout`: **the manifest value is a switch that turns off content comparison.**
  - A `zip:sha256:` value in `state.json` makes `status` skip `Plan` and the drift comparison for every binary, including a release binary that carries the signed tree. It then exits 0 on presence alone (`Lstat`).
  - This is the class of incident behind ADR-034: stale content in use while `status` reports nothing wrong. The line reads "not compared" rather than "drift none", so the text is honest, but an old zip installed today goes unnoticed for as long as the home lives.
  - Not a trust boundary: `state.json` and `~/.local/bin/sh-agent` are both writable by the same user, so a tamperer can silence `status` either way.
  - Prefix hygiene: `ResolveSource` keys on `zip:`, while `status` and `knownCheckout` key on `zip:sha256:`. A `zip:` value without `sha256:` is excluded as a candidate by one and passed to `samePath` as a path by the other. Harmless today.
  - Remediation:
    - One helper, e.g. `harness.ZipSum(repo) (sum string, ok bool)`, used at all three sites.
    - In a release binary, also report whether the installed files match the embedded release ("differs from release vX.Y.Z"), so staleness stays visible.
    - Per-file `crypto/sha256` in the manifest would restore real drift detection; ADR-038 rejected it on versioning grounds, which is a policy choice, not a technical limit.
- **[low / confidence: high]** Texts that describe security properties the code no longer has:
  - `internal/harness/repo.go:23-33`: the `ResolveSource` comment still says "A RELEASE binary installs the harness built into it and nothing else … not one named on the command line", and `--from-zip` now does exactly that.
  - `.agents/docs/architecture/DECISIONS.md` ADR-034's header is still "Accepted", with no "Superseded in part by ADR-038" mark, although that file's preamble says to mark a superseded ADR.
  - The output line "the binary itself is not replaced" (main.go `updateFromZip`, asserted by `TestUpdateFromZip`), the usage text "the binary is not replaced", README ("never replaced this way") and R9 ("the binary itself is not replaced") are all wrong in a case that matters. `installFiles` → `Plan(src, self, …)` writes the RUNNING binary to `~/.local/bin/sh-agent` (`plan.go:188-189`). Running `go run ./cmd/sh-agent update --from-zip …`, or an unverified binary from Downloads, therefore overwrites the verified release binary on PATH.
  - Remediation: reword to "no newer binary is fetched; the running binary is installed as is". R9's sentence needs the operator's approval (REQUIREMENTS is immutable without it). Update the `ResolveSource` comment and annotate ADR-034's header.
- **[low / confidence: high]** Cleanup gaps (`cmd/sh-agent/main.go` `updateFromZip`, `internal/harness/zip.go:70-74`):
  - SIGINT, SIGTERM or a kill during the unpack or install leaves `~/.config/sh-agent/zip-*` behind. So does a panic inside `UnpackZip`: the deferred check sees `err == nil`.
  - Nothing ever sweeps leftovers, unlike `report-*`, which `sweepInterruptedClones` (`report.go:409-420`) globs and removes when the owning process is dead.
  - The error from `os.RemoveAll` is discarded.
  - Exposure is small: each directory is 0700 and `report`'s `Collect` reads only named files (`report.go:231-290`), so a leftover tree is never sent. It does survive `uninstall`, whose `os.Remove` of the config directory fails on a non-empty directory.
  - The dry-run path that unpacks under the system temp directory (plan D10) has no test asserting it cleans up. `TestUpdateFromZipDryRunWritesNothing` checks only the home.
  - Remediation:
    - Name the directory `zip-<pid>-*` and, at the start of `updateFromZip`, sweep only the directories whose process is dead, as `sweepInterruptedClones` does. A blind sweep would delete a concurrent run's tree.
    - Report a failed `RemoveAll`.
    - Assert the system-temp cleanup, e.g. with `t.Setenv("TMPDIR", t.TempDir())`.
- **[low / confidence: high]** Refusals claimed in comments or the plan but untested (`internal/harness/zip_test.go`):
  - a header that understates the uncompressed size (the `unpackFile` comment claims it cannot carry more through; it holds by construction, through `io.LimitReader` and `archive/zip`'s own size check, but no test builds one);
  - a directory entry and a file with the same path, in either order (refused by `O_EXCL` / `ENOTDIR`);
  - aliased names (`top//x`, `top/./x`, `./context/…`);
  - the `filepath.IsLocal` fallback in `zipNameSafe`, which is unreachable on Unix after the three explicit checks (defence in depth, fine).
  - Remediation: add fixture cases so a future refactor cannot silently lose these refusals.

## Checked

- **Zip-slip / path handling:**
  - Every entry is name-checked archive-wide before anything is created (`zip.go:47-51`). Refused: absolute names, backslashes, `..` elements, and anything `filepath.IsLocal` refuses. NUL bytes make `os.OpenFile` fail.
  - The checked name is the central-directory name, which is the one `archive/zip` uses, so a local header's name cannot differ from it.
  - Stripping the prefix cannot produce an escaping `rel`. A `top//x` form yields a leading `/`, whose empty first element fails the `zipTrees` filter, and `filepath.Join` never treats it as absolute anyway.
  - Symlinks are never created: a Unix symlink entry inside the five trees is refused, one outside is skipped, and a name ending in `/` is forced to a directory by `archive/zip`. Only `MkdirAll` and `O_EXCL` regular files are created, inside a fresh `MkdirTemp` directory, so no existing link can be followed.
  - The zip format has no hard links. Mode bits from the archive are discarded (0644 files, 0755 directories), so no setuid survives.
  - APFS case and normalisation collisions, and aliased duplicates, hit `O_EXCL` or `ENOTDIR` and refuse the whole zip. The case-sensitive `zipTrees` filter only ever narrows what is written.
- **Decompression:**
  - Written bytes are counted: `io.LimitReader(rc, max+1)` per file, on top of `archive/zip`'s own check against the header size. The total overshoot is bounded at 128 MiB + 16 MiB + 1.
  - Entries outside the five trees are never opened. Unsupported methods and encrypted entries fail `f.Open`.
  - The entry cap is applied after `zip.OpenReader` has parsed the central directory, whose memory is proportional to the file size (Go's own guard). This is acceptable for an operator-chosen file. Directories are the gap above.
- **Temporary directory permissions:**
  - `os.MkdirTemp` creates the directory 0700, and the parent is created 0700 only when absent.
  - `Install` → `TightenConfigDir` (`apply.go:30`) later corrects `~/.config/sh-agent` to 0700.
  - The dry-run fallback in the system temp directory is also a 0700 directory with a random name, so it is safe on a shared `/tmp`.
- **Cleanup:** every error return in `UnpackZip` sets the named `err`, so the deferred `RemoveAll(tmp)` runs; returns that shadow `err` still assign the result. `updateFromZip` defers `RemoveAll(dir)` before any later return. The gaps are listed above.
- **Secrets:**
  - A regex scan over the full diff and all untracked files, for GitHub, AWS and Slack tokens, private-key armour, `refresh_token`, `client_secret`, `password`, API keys and email addresses, matched only a test fixture's symlink target text `/etc/passwd`.
  - The pattern was confirmed to fire on a synthetic token.
  - The manifest records only the hash, never the local zip path, and `report` sends `state.json` through `redactDocument`.
- **Injection:**
  - No shell, SQL or template use. Entry names reach errors quoted with `%q`.
  - The ≤12 bytes of manifest text echoed into errors, and the manifest paths given to `Lstat` in `status`, come from a file writable only by the same user.
  - No git command runs in the unpack directory: `Install` skips `Inspect` when `Origin` is set, and the path never reaches `Reinstall` or `go build`.
- **Authz/authn:**
  - The `--from-zip` branch comes before `ResolveSource` and refuses `--check`, `--repo` and non-`update` commands. `$SH_AGENT_REPO` is silently ignored on this path, where ADR-034's style would refuse it; this is harmless because the file is named explicitly.
  - A release binary still refuses a recorded or discovered path: `ResolveSource` is unchanged for release versions, and the `zip:` value is never used as a path.
  - An older development build would read the value as a relative candidate path. That gives no new capability: planting a harness-shaped tree in the working directory already wins by design.
  - The trust model is the medium finding.
- **Crypto:** only `crypto/sha256` from the standard library, as a checksum. Nothing is hand-rolled.
- **Dependencies:**
  - No new module; the code uses the standard library only (`archive/zip`, `crypto/sha256`, `encoding/hex`).
  - `govulncheck ./...` under go1.27.1: "No vulnerabilities found." (exit 0).
  - `npm audit` does not apply: there is no `package.json`, `web/` or `ui/`.
- **GCP surface:** none touched. No IAM, service accounts, deploy flags or Secret Manager use in the diff.
- **Trail markdown** (read for secrets, credentials, private data and instruction-shaped text):
  - Files read: BUG-093, both plans, the three cold-read reviews, the three verify logs, and every hunk of the nine modified docs.
  - None contain credentials. The private data present is the operator's GitHub handle and repository name, which appear throughout the repository already.
  - No text is addressed to an agent in a way that would redirect it; the scan also looked for "ignore previous instructions" phrasing.

## Not examined

- **Blast radius outside the diff, not read end to end:**
  - `cmd/sh-agent/main.go` outside the diff hunks, `run()` lines 95-160 and `installFiles` (`update()`, `reorder()`'s body, `adopt`, `diff`, `list`/`render`, the unchanged middle of `status`);
  - `internal/harness/plan.go` beyond lines 180-195;
  - `inventory.go` `checkManagedTree`;
  - `apply.go` beyond lines 95-160 and the `TightenConfigDir` call site;
  - `release.go`, cited only for its `ssh-keygen` verification path in a remediation, not re-read;
  - `report.go` beyond `Collect`.
  - These are unchanged code the new path reaches. Every release-diff path itself was read: every hunk of each of the 13 modified files and every line of the 12 untracked files.
- **macOS:** APFS case and normalisation behaviour and the PATH_MAX of 1024 are reasoned, not run on a Mac. The e2e log also names macOS as not exercised.
- **Not executed (defensive review only):** no hostile archive was built. The directory-amplification bound and the TOCTOU window are reasoned from the code and the standard library's semantics, not demonstrated.
- Two reports appeared during this run and were not read, because they are other agents' records outside the dispatched range: `.agents/docs/reviews/2026-09-28-review-update-from-zip-code.md` and `.agents/docs/reviews/2026-09-28-review-update-from-zip-docs.md`.
- `go test` was not rerun by this review. The test results cited are the verifier's logs, and Phase 1's log predates the BUG-093 change to `zip.go` (its line numbers differ).
