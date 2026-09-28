---
owner: security
date: 2026-09-28
---
# Security review, round 2 — update-from-zip   (2026-09-28)

Range reviewed: 1cb2608..worktree (`git diff 1cb2608`: 17 modified files, 240+/33-, unchanged during this run; plus the 16 untracked files in the dispatch snapshot)

## Status: COMPLETE

## Verdict: FINDINGS

The worst problem is still the round-1 medium, which stays open by decision. An unsigned archive becomes the global instruction set of two agents that run commands with the operator's credentials, and that holds under a signed release binary too. The `UNVERIFIED` line is adequate as an interim disclosure, but it is not a control, and it counts as an interim only if two things are written down (finding 1). The code fixes for round-1 lows 2 and 3 are sound.

Everything new is low:
- BUG-093's defect class recurs at a later refusal point. Reproduced: a zip that `Plan` refuses leaves `~/.config/sh-agent` on a fresh home.
- The stale-directory sweep builds a glob from an unescaped home. Reproduced: it deleted another home's directory.
- `mkdirAllTracked` counts `EEXIST` as "created".
- The skill-matrix grants widen while its text overstates a protection.
- Several texts still state the removed "binary is not replaced" property.
- The 256 MiB archive cap is about 90 times the real archive.

The untested refusals from round 1 are still untested.

## Findings

- **[medium / confidence: high]** Round-1 medium, **open by decision** (plan § Scope "Out"; ADR-038 "Rejected" bullet 4): `update --from-zip` installs an unverified archive as agent instructions, including from a signed release binary.
  - **Is the `UNVERIFIED` label adequate as an interim?** Yes as disclosure, no as a control.
  - What holds:
    - It is on the first output line of both real and dry runs, and `TestUpdateFromZip` pins it.
    - README § Install from a zip and `delivery.md` say the same thing.
    - I saw it printed first in my own probe run.
  - What it does not do:
    - It is printed while the run that installs is already under way. Nothing asks for confirmation, and nothing stops the install.
    - The printed sha256 has no reference value to compare against. GitHub publishes none for "Download ZIP" or "Source code (zip)", so the operator cannot act on it.
    - `status` shows `from zip sha256:…` but never "unverified".
    - Cosmetic: a zip that is refused later still prints "installing its harness …" before the error (probe A below).
  - It is adequate as an interim only if:
    1. The operator's acceptance is re-recorded. ADR-038 still gives "a zip is unsigned" as the reason, and that is a choice about the release process, not a constraint.
    2. The signed-zip option exists as an open operator decision in a trail file. I found it in none of the files I read (the plan, ADR-038, ROADMAP). Plan D12 routes only the R9 wording to the release report, so today the option lives only in chat.
  - Remediation is unchanged from round 1, with no new dependency:
    - `/preflight` publishes a `git archive` zip of the tag and lists it in the already-signed `SHA256SUMS`.
    - `--from-zip` checks the zip's `crypto/sha256` against that list through `release.go`'s existing `ssh-keygen -Y verify` path.
    - Optionally, the `status` line says `unverified` for a zip install.

- **[low / confidence: high]** `cmd/sh-agent/main.go` `updateFromZip` → `installFiles`: **BUG-093's class recurs at a later refusal point.**
  - `UnpackZip` creates the parent chain once its own checks pass and removes it only on its own errors.
  - `Plan` can still refuse the unpacked tree afterwards: a skill with no description, a secret-shaped file, an unmapped tool, a `#` frontmatter value.
  - When it does, the deferred `os.RemoveAll(dir)` removes only the unpack directory. `~/.config/sh-agent`, and `~/.config` if it was created, stay behind.
  - **Reproduced** (probe A): a fresh scratch home and a zip adding `skills/zz-nodesc/SKILL.md` with no description gave exit 1, `skill zz-nodesc: SKILL.md has no description`. The home then held `.config/sh-agent`.
  - This contradicts the comment "a refused zip writes nothing" and BUG-093's prevention rule. It is not a security exposure: the directory is empty and mode 0700. On the code reviewer's scale it may rate higher.
  - Remediation:
    - `updateFromZip` records whether the parent existed before the call. When `installFiles` fails, it removes the parent with `os.Remove` (empty only) if it did not exist.
    - Add a fresh-home CLI test with a zip that `Plan` refuses.

- **[low / confidence: high]** `cmd/sh-agent/main.go` `sweepStaleZipDirs`: **the glob pattern is built from an unescaped home, so the sweep reaches outside the named home.**
  - **Reproduced** (probe B), with both `zip-old` directories aged three hours:
    - `update --from-zip notzip.zip --home '<scratch>/h[1]'` deleted `<scratch>/h1/.config/sh-agent/zip-old`, a different home's directory.
    - It left the named home's own stale `zip-old` in place.
  - Reach: only `zip-*` directories (not symlinks, because `Lstat` + `IsDir`) older than an hour, under any `.config/sh-agent` the pattern matches, within the user's own permissions. A home containing `*`, `?` or `[` is unusual.
  - `RemoveAll` errors are still discarded here, in `UnpackZip`'s cleanup, and in the `updateFromZip` defer. Round 1's remediation for this was not taken.
  - The same mechanism predates this diff in `sweepInterruptedClones` (`report.go:411`, `report-*`, gated by pid).
  - Remediation, all stdlib:
    - Use `os.ReadDir(parent)` plus `strings.HasPrefix(e.Name(), "zip-")`, so no pattern is ever matched against a path.
    - Print a `note:` when a removal fails.
    - Give `report.go` the same fix.
  - **The dispatch's symlink question:**
    - A `zip-*` symlink is skipped (`Lstat`).
    - `os.RemoveAll` unlinks symlinks inside the tree rather than following them, and uses the `*at` syscalls on Linux.
    - Swapping a directory for a link between `Lstat` and `RemoveAll` needs a process running as the same user, which is not a boundary.
    - If `~/.config/sh-agent` is itself a symlink (which `state.go` supports), the sweep acts in its target, which is still sh-agent's own directory.
    - Age is measured by the directory's mtime, which is set when the unpack starts. A concurrent run would lose its tree only if its install took more than an hour. Negligible.

- **[low / confidence: high]** `internal/harness/zip.go:203-209` `mkdirAllTracked`: **`EEXIST` is recorded as "created".**
  - If something creates a path between the `Lstat` walk and `os.Mkdir`, the path is still appended to `created`.
  - The deferred cleanup then calls `os.Remove` on an entry this call never made. `os.Remove` deletes an empty directory, and it deletes a file outright.
  - Exploitability is negligible: it needs a race by a process running as the same user. It is still a guard in front of a destructive action that can be wrong.
  - Remediation: append only when `os.Mkdir` returned nil.
  - Mode `0o700` is fine. It creates `~/.config` as 0700 on a home that lacks one, which is stricter than `MkdirAll`'s 0755. When `Lstat` treats a non-ENOENT error as "missing", the following `Mkdir` fails, so that case is benign.

- **[low / confidence: medium]** `.agents/docs/testing/skill-matrix.md` step 2 (Phase 4, item 7, in this diff): **the sandbox grants are widened, and the text overstates a protection.**
  - What changed: the diff adds `Bash(bash:*)`, `Bash(python3:*)` and unrestricted `Bash(curl:*)` for Claude, and `command(bash|python3|curl)` for agy. It drops the v1.7.0 gate's narrowing, while the sandbox home holds a copied login token.
  - What the text gets right: the allowlist is not a boundary, because `go` and `git` already run arbitrary code, so no new class of capability is added.
  - What it overstates: "The protection is that the copied token lives only for the batch and is deleted by name after it."
    - Deleting the local copy bounds the window. It does not revoke a copy exfiltrated during the batch.
    - Unrestricted `curl` makes exfiltration a single command for any run that ingests hostile text.
    - The `bash -n` / `py_compile` rationale explains `bash` and `python3`, not outbound `curl`, which was localhost-only before.
  - The operator approved matrix option 1 (plan D5), hence low.
  - Remediation:
    - Keep `curl` limited to localhost unless a named scenario needs more.
    - State the residual risk: an exfiltrated token outlives the batch.
    - Sign the sandbox session out, or revoke it, after each batch, rather than only deleting the file.

- **[low / confidence: high]** Texts that still state the removed property (round-1 low 5, partly closed):
  - `cmd/sh-agent/main.go:42`, the usage line "(the binary is not replaced)", and `:1045`, the `updateFromZip` doc comment "the binary is not replaced".
    - `installFiles` → `Plan(src, self, …)` writes the RUNNING binary to `~/.local/bin/sh-agent` (`plan.go:188-189`).
    - So an unverified binary run from Downloads replaces the verified one on PATH.
    - These are code texts, fixable without approval. The output line and README are now accurate.
  - `internal/harness/repo.go:23-33`: the `ResolveSource` comment "A RELEASE binary installs the harness built into it and nothing else … not one named on the command line" is unchanged.
  - `DECISIONS.md`: ADR-034's header is still "— Accepted", with no "Superseded in part by ADR-038" mark.
  - ADR-038's "What is unpacked" bullet lists 16 MiB, 128 MiB and 5,000 entries. It omits the 256 MiB archive cap and the 1,024-byte and 32-level name bounds.
  - R9's "the binary itself is not replaced" is operator-approved text. D12 routes it to the release report, which is acceptable if that happens.

- **[low / confidence: medium]** `internal/harness/zip.go:52-62` with `release.go:352` `readCapped`: **the memory bound is loose.**
  - What is held: up to 256 MiB in one buffer, which `io.ReadAll`'s growth can transiently double. `zip.NewReader` then parses the entire central directory before the 5,000-entry cap applies.
  - Worst case, reasoned from `archive/zip`'s structures and not measured: a 256 MiB archive of minimal directory records holds several million entries at roughly 5 times their on-disk size, a peak of around 1.5–2 GB.
  - The real archive is 2,853,062 bytes (`git archive` of HEAD, 852 files; the e2e zipball had 936 entries), so the cap is about 90 times the need.
  - This is a local denial of service from a file the operator chose, and an improvement on round 1, which had no cap at all.
  - Remediation:
    - Lower `maxZipArchiveBytes` to about 32 MiB, still more than 10 times the real archive.
    - When `Stat` gives a size, read with `make([]byte, size)` and `io.ReadFull`.
    - The `readCapped` comment ("a file a release supplied") now also covers the operator's zip.

## Round-1 status

- **Low 2, name length and depth: CLOSED.**
  - `zipNameSafe` refuses names over `maxZipNameBytes` (1,024) and names with 32 or more slashes (`maxZipNameDepth`). The check runs archive-wide, before anything is created.
  - Tested only far over the bounds (`TestUnpackZipRefusesDeepOrLongNames`: 41 levels, 1,104 bytes), never at the boundary. That test also does not assert that nothing was created, although the code's ordering guarantees it.
  - Residual: 5,000 entries × about 30 new directories ≈ 150,000 directories, roughly 600 MiB of directory blocks on ext4. That is bounded, against millions before. A depth closer to the real trees (fewer than 10 levels) would cut it further.
- **Low 3, sum from the same buffer as the unpack: CLOSED.**
  - The file is read once. The sum, `zip.NewReader`, the check pass and the write pass all use the same immutable buffer. `TestUnpackZipSumIsOfTheBytesRead` pins the value.
  - Output wording: CLOSED. The line reads "UNVERIFIED — a zip carries no signature; trust rests on your own download … installing its harness and the running binary …; nothing is downloaded and nothing in the zip replaces the binary". It is accurate and pinned by a test.
- **Cleanup gaps: PARTLY CLOSED.**
  - Parent-chain removal: closed for refusals inside `UnpackZip`, via `mkdirAllTracked` and the deferred reverse removal, tested by `TestUnpackZipLateRefusalCreatesNoParent` (symlink, duplicate). **Open** for refusals after `UnpackZip` returns (finding 2).
  - Stale-directory sweep: added, keyed by age (one hour) rather than pid. It covers what a kill, a SIGPIPE or a panic left behind, on the next real run. `TestUpdateFromZipSweepsStaleDirs` tests it, including that a dry run sweeps nothing and a fresh directory is kept. Defect: finding 3.
  - `RemoveAll` errors: still discarded.
  - Dry-run cleanup under the system temp directory: still untested (no `TMPDIR` assertion).
- **Untested refusals: NOT CLOSED.** None of the three testable items has a fixture:
  - A header that understates the uncompressed size. `checkZipTrees`' comment now makes this claim too.
  - A directory entry and a file with the same path, in either order. It is refused late, by `O_EXCL` or `ENOTDIR`.
  - Aliased names (`top//x`, `top/./x`). They are refused late, by `O_EXCL`, after `filepath.Join` cleans them.
  - The `filepath.IsLocal` fallback is unreachable on Unix, which is fine.
  - The new tests cover other things: long and deep names, the archive cap, late refusals on a fresh parent, refusal before any command, and the empty value.
- **Low 4, the manifest value switches off content comparison:** unchanged by design.
  - ADR-038 records it, and `diff` now prints "not compared" too.
  - The shared-helper remediation was not taken: three sites still parse `zip:sha256:` separately, and `ResolveSource` keys on `zip:`. Harmless.
- **Low 5, texts:** partly closed (finding 6).

## Checked

- **New code, per the dispatch:**
  - `readCapped` at 256 MiB: finding 7. A FIFO or `/dev/zero` passes the `Stat` check with size 0 but is still capped by the `LimitReader`.
  - `checkZipTrees` decompressing into `io.Discard`: sound.
    - It counts the bytes actually produced, per file (`LimitReader`, max+1) and in total (it stops past 128 MiB).
    - `archive/zip`'s checksum reader rejects an understated size with `ErrFormat` and checks the CRC at EOF.
    - Entries outside the five trees are never opened, and the write pass re-reads the same bytes.
    - The cost is that each harness entry is decompressed twice.
  - `mkdirAllTracked` with `0o700`: finding 4.
  - `sweepStaleZipDirs`: finding 3.
  - The `fs.Visit` block: sound.
    - `Visit` sees flags set on the command line, so `--from-zip=` and `--from-zip ""` are refused (tested).
    - Every refusal comes before any command, tested for eight commands with an empty `PATH`, asserting no output and no writes.
    - Only resolving the home runs before it.
    - `--force`, `--feedback-repo` and `$SH_AGENT_REPO` are silently ignored, which is harmless.
    - `update --from-zip --dry-run` takes `--dry-run` as the file name and fails at `open`, which is safe.
- **BUG-094, whether anything left the machine beyond the record:** no.
  - `gh api repos/kesokaj/styrman-harness-feedback/commits`: the newest commit is `c5d47b297ff5`, 2026-09-28T15:21:17Z, "report: dave 2026-09-28-172114". Its parent is from 2026-09-18.
  - It adds one file, `reports/dave/2026-09-28-172114/report.json`, 16 lines:
    - `ts` 2026-09-28T17:21:14+02:00, `host` dave, linux/amd64;
    - `vendor`, `files`, `harness_version` and `installed_at` all null;
    - `harness_commit` 1cb2608, `drift` 257, `skills` 27, `agents` 7, `go` go1.27.1, `redactions` [].
  - The record's "metrics only" matches. It does not list `ts` or the null fields, which are part of every metrics record.
  - The only branch is `main`, at that commit. The repository is private and was last pushed at 15:21:19Z. No later commit exists, which is consistent with "the mutation run pushed nothing".
  - The prevention holds. The test empties `PATH`, and `report` runs git as `exec.Command("git", …)` (`repo.go:250,257`) and gh as `exec.Command("gh", …)`, both resolved through `PATH`.
  - Removing the pushed record is correctly left to the operator (status `open`).
- **Secrets:**
  - I regex-scanned the full diff and every untracked file for:
    - GitHub, AWS, Slack, Google and Anthropic token shapes;
    - private-key and OpenSSH armour and `ssh-ed25519` keys;
    - `refresh_token`, `client_secret` and `access_token`;
    - `password=` and `api_key=`;
    - email addresses;
    - "ignore previous instructions" phrasing.
  - The scan fired on a synthetic `ghp_` token (the positive control).
  - It had two matches, both in round 1's own description of its scan.
  - `.gitignore` and `secretShaped` gain `*oauth-token*`, which is an improvement.
- **Injection:** there is no shell, SQL or template use. Entry names reach errors quoted with `%q`, and no git command runs in the unpack directory.
- **Authz/authn:** a release binary still refuses `--repo` and `$SH_AGENT_REPO`, and the `zip:` value is never a candidate path (decoy test). The trust model is finding 1.
- **Crypto:** only `crypto/sha256` from the standard library, used as a checksum. Nothing is hand-rolled.
- **Dependencies:**
  - No new module.
  - `govulncheck ./...` under go1.27.1 printed "No vulnerabilities found." and exited 0.
  - `npm audit` does not apply: there is no `package.json`.
- **GCP surface:** none touched.
- **Trail markdown, read end to end:**
  - Files: BUG-093, BUG-094, both plans, the three cold reads, the code and docs reviews, the three verify logs, and every hunk of the modified docs (DECISIONS, local, remote, delivery, skill-matrix, AGENTS, README, REQUIREMENTS, ROADMAP, VENDORED).
  - Credentials: none.
  - Private data: the operator's GitHub handle and the host name `dave`. Both already appear in this private repository and in the private feedback repository.
  - Nothing is instruction-shaped: no text aimed at an agent redirects it.
- **Probes** in `.agents/scratch/secr2-4139464/` (gitignored, 9 MB, left in place):
  - A: a benign zip whose extra skill has no description, installed into a fresh home.
  - B: a non-zip file, with `--home '<scratch>/h[1]'` and two aged `zip-old` directories.
  - The binary was built from the working tree. No real home was touched.

## Not examined

- **`.agents/docs/reviews/2026-09-28-review-update-from-zip-code-r2.md` (122 lines).** It appeared during this run (mtime 17:32:14, not in the dispatch snapshot). It is another agent's record and I did not read it. It is a release-diff path, so under the dispatch's rule it blocks the gate until it is read; a short dispatch over that one file clears it.
- `reviews/2026-09-28-review-v1.7.0-381e86d-code.md`, the basis for the skill-matrix change (Phase 4). It is outside the diff.
- Unchanged code the path reaches, not re-read end to end this round: `release.go` beyond `readCapped`, `report.go` beyond its sweep, `apply.go` beyond the `Origin` hunk, `plan.go` outside lines 175-200, `inventory.go` and `state.go`.
- macOS (APFS, PATH_MAX 1024): reasoned, not run.
- The memory multiplier in finding 7 and the directory residual under low 2: reasoned, not measured. No hostile archive was built. The probes were benign.
- `go test` was not run by this review. The verifier logs predate D11; the e2e rerun covers D11 at the CLI level.
