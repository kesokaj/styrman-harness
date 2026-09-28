---
owner: security
date: 2026-09-28
---
# Security review — signed-zip   (2026-09-28)

Range reviewed: 0b762db..worktree (`git diff 0b762db`: 13 modified files, plus 4 untracked files: `plans/2026-09-28-signed-zip.md`, `reviews/2026-09-28-coldread-signed-zip-r1.md`, `reviews/2026-09-28-coldread-signed-zip-r2.md`, `testing/2026-09-28-verify-signed-zip-e2e.md`)

## Status: COMPLETE

## Verdict: FINDINGS

The verification works. When it runs, it checks the right things in the right order, and I found no bypass in the verify-then-unpack path. The problem is that it only runs when two files sit beside the zip, and whoever supplies the zip can leave them out. A zip with nothing beside it installs with exit 0. The output line says `UNVERIFIED` partway through, and that is all. So against an attacker who hands over the zip, the v1.8.0 gate's medium ("an unsigned zip becomes the agents' instructions") is unchanged. It is narrowed to one case: the operator fetched all three release files and something altered them in transit. It is not closed. The plan calls this change "the fix", and ADR-039 does not record what remains open. Second, when verification does run and fails because the files were altered, the refusal message explains how to skip the check. That hint is attached to the failures that point to tampering: a bad signature and a missing `SHA256SUMS`. The harmless stale-pair case gets no such hint. Everything else is low.

## Findings

- **[medium / confidence: high]** `internal/harness/release.go:425-428` (`VerifyZipSums`) with `cmd/sh-agent/main.go` `updateFromZip` (the `check` closure and the trust line). **What is wrong:** verification is triggered by `SHA256SUMS`/`SHA256SUMS.sig` being present in the zip's own folder. **Why it matters:** the same party who supplies the zip controls those files. With both absent, `ErrZipUnsigned` becomes nil and the zip unpacks and installs with exit 0. The only disclosure is `UNVERIFIED — no signed SHA256SUMS beside it; trust rests on your own download`, one clause in a long stdout line, and an agent running the command may never pass it on. This includes a zip named `styrman-harness-vX.Y.Z.zip`, which claims to be a signed release. **The record:** the v1.8.0 medium still holds for any zip that reaches the operator without its sums. The plan (`## Why`) presents the change as the fix for that medium. ADR-039's Consequences do not say the risk remains for an attacker-supplied zip. **Remediation:**
  - (a) State the residual in ADR-039's Consequences, and carry it as a narrowed, still-open row in the v1.9.0 gate report.
  - (b) Operator design decision: require an explicit opt-in to install without verification, for example an `--unverified` flag. Without the flag, refuse when no signed sums are present. R9's `main`-zip path keeps working with one deliberate word, and neither a person nor an agent can reach an unverified install by accident. ADR-039 rejected a `--sums` flag and refusing every unsigned zip; this is neither of those. It needs an R9 wording change and an ADR, per project rule 3.

- **[medium / confidence: medium]** `internal/harness/release.go:442` (`VerifyZipSums`). **What is wrong:** every `verifySignedSums` failure is wrapped with `(move SHA256SUMS and SHA256SUMS.sig out of <dir> to install the zip unverified on purpose)`. Measured with the E2E's scratch-signer binary, dry-run, in `.agents/scratch/security-signed-zip/`:
  - `SHA256SUMS` edited after signing → `SHA256SUMS signature does not verify (ssh-keygen: Signature verification failed: incorrect signature) (move SHA256SUMS and SHA256SUMS.sig out of rel to install the zip unverified on purpose)`, exit 1.
  - `.sig` present, `SHA256SUMS` absent → the same hint, exit 1.
  - The harmless case plan D7 wrote the hint for — a genuine pair beside a differently named zip — gets `no line names styrman-harness-main.zip — use the release's own file name`, with no move-away hint.

  **Why it matters:** the hint sits on exactly the failures that mean someone altered the files. In a harness whose users are agents that act on tool output, that makes the refusal a self-documented, one-step downgrade to installing whatever the zip says. **Remediation:**
  - On signature, marker and missing-file failures, say the files are not genuine and to download the release's three files again. Drop the move-away clause there.
  - If a hint stays anywhere, put it only on the `no line names` stale-pair path, where the `use the release's own file name` hint for a `main` zip is also misleading.
  - Add a test that the signature-failure message does not contain the bypass text.

- **[low / confidence: medium]** `cmd/sh-agent/main.go` `updateFromZip`. **What is wrong:** a genuine older release's zip, with its own genuine `SHA256SUMS` and `.sig`, verifies and installs with `VERIFIED against … release v<older>`. Nothing compares it with the running binary's version, and the manifest records only `zip:sha256:<hex>`, so nothing can be compared later either. `UpdateFromRelease` refuses "sideways or backwards" for exactly this reason (`release.go`, the `compareRelease(tag, version) <= 0` case, security review 2026-09-17); the zip path does not. **Remediation:** when `IsReleaseVersion(buildVersion())`, warn or refuse when `compareRelease(verifiedTag, buildVersion()) < 0`. Both are existing helpers, so no new dependency.

- **[low / confidence: medium]** `internal/harness/release.go:40` (`var releaseSigner`) and `deploy/remote.md` § Releases steps 1–4. **What is wrong:** the trust anchor is a package-level string, so `-ldflags -X …/internal/harness.releaseSigner=…` replaces it at link time. The E2E did exactly that.
  - Measured: the E2E binary trusts `e2e@signer` and still contains the release signer's literal (`grep -a -c` = 1 for both strings). Grepping the bytes of a built binary therefore cannot show which signer is in effect.
  - The release procedure checks the checksums (`sha256sum -c`) and the signature (`ssh-keygen -Y verify`), but never runs a built binary against the freshly signed `SHA256SUMS`. A release built with a stray `-X` (the E2E build line reused, say) would publish binaries that trust another key, and every documented check would still pass.

  **Remediation:**
  - After step 3, and before `gh release create`, run `./sh-agent-linux-amd64 update --from-zip styrman-harness-vX.Y.Z.zip --home <scratch> --dry-run` in the release directory. Require `VERIFIED against the signed SHA256SUMS of release vX.Y.Z`. It uses existing commands and proves which signer the host-platform binary actually uses.
  - Alternative: make the signer a `const` and give tests a non-string (for example func-typed) hook, which `-X` cannot set. This removes the E2E's method, so it is the operator's call.
  - Separately, the comment at `release.go:34` ("a variable only so a test can substitute a test key's line") is now stale.

- **[low / confidence: high]** `.agents/docs/deploy/remote.md` § Releases, the paragraph "After the upload: `gh release view …` lists the five assets". **What is wrong:** the diff changed the loop comment to "all six", but this sentence still says five. If the zip upload fails and the reader checks against the paragraph, the release passes as complete, and zip places are left with only unsigned `main` zips. **Fix:** "six".

- **[low / confidence: high]** Test coverage of the control. **What is wrong:**
  - No CLI test drives the `VERIFIED` branch. The `verifiedTag` wiring and the `VERIFIED` line are covered only by the one-off E2E; plan Phase 1 names `TestUpdateFromZipVerified|TestUpdateFromZipTampered`, which do not exist.
  - No test covers `.sig` present with `SHA256SUMS` absent. My probe shows it refuses correctly.
  - No test pins the refusal texts (see the second finding).

  **Remediation:** a `cmd/sh-agent` test that builds the binary with the test signer into `t.TempDir()` (`go build -ldflags -X`) and runs the verified and tampered cases, plus a `TestVerifyZipSums` subtest for "sums missing, sig present".

- **[low / confidence: high]** `internal/harness/release.go:425-428`. **What is wrong:** "unsigned" is decided when both `os.Stat` calls fail for any reason, not only `fs.ErrNotExist`. An ELOOP or EIO on both therefore falls through to an unverified install instead of a refusal. **Why it is low:** anyone who can cause those errors can also delete the files, so it adds nothing to the first finding. But a verification guard whose precondition check errors should refuse. **Fix:** `errors.Is(err, fs.ErrNotExist)` for both, via `os.Lstat`, so that a dangling symlink counts as present and refuses.

- **[low / confidence: medium]** `deploy/remote.md` § Releases step 1b (`git archive … vX.Y.Z`). **What it does:** publishes nothing secret (see Checked), but ships the whole tracked tree, including `.agents/docs/` (729 entries). That tree names internal hostnames, machine paths, the operator's email and a report record's contents; the same section's own paragraph already lists them. **Where it ends up:** on a zip place, typically a corporate workstation, the downloaded zip stays in the operator's Downloads, while `UnpackZip` writes only the five trees. **Why it is low:** the content is identical to GitHub's own "Source code (zip)" for the tag, and the repository is private. **Optional minimization:** `git archive --format=zip --prefix=styrman-harness-vX.Y.Z/ -o … vX.Y.Z context skills agents commands rules`. Check that R9's "a zip of its tree" wording still fits.

## Checked

- **Verification soundness (asked):**
  - The signature is verified over the in-memory `SHA256SUMS` bytes (stdin), against the compiled-in signer: `-I <identity>`, `-n sh-agent-release`, and an `allowed_signers` holding only that line.
  - The `# styrman-harness vX.Y.Z` marker is read only from verified bytes.
  - The zip's line is matched by exact `filepath.Base` name. Two lines naming it with different hashes are refused; two with the same hash are harmless.
  - The sum is taken from the one buffer that is unpacked: `UnpackZip` reads once, hashes, calls `check`, then `zip.NewReader(bytes.NewReader(data))`.
  - All tests I ran passed (`-race -count=1`): `TestVerifyZipSums` (6 subtests), `TestUnpackZipCheckStopsBeforeWriting`, `TestUpdateFromZipForeignSignatureRefused` (which uses the real compiled-in signer) and `TestUpdateFromZipUnsignedSaysUnverified`.
- **Bypass attempts (asked):**
  - Name listed twice: refused when the hashes differ, and the file is signed anyway.
  - Path tricks: only `filepath.Base` reaches `checksumFor`. A zip name with `(1)` or another folder's name gets "no line names" and is refused. Entry-level path safety is unchanged and was re-read: `zipNameSafe` with `filepath.IsLocal`, symlink and irregular entries refused, `O_EXCL`, size and depth caps.
  - Symlinked `SHA256SUMS`: it is followed and its content verified. Measured: a symlink to the genuine file with the genuine zip gives `VERIFIED`; with the byte-flipped zip it gives `checksum mismatch`.
  - TOCTOU: none between verify and unpack. The zip and `SHA256SUMS` are verified from memory; `ssh-keygen` reads `.sig` once, and swapping it can only make verification fail. The unpack directory is `os.MkdirTemp` (0700).
- **Stale or foreign pair (asked):** refuses and does not fall back to unverified. Measured (probe 3); tested ("foreign key", "renamed zip"). One file present refuses: tested for a missing `.sig`, measured for a missing `SHA256SUMS`.
- **Temp `allowed_signers` (asked):** `os.MkdirTemp("", "sh-agent-verify-*")` (random name, 0700), file written 0600, `defer os.RemoveAll`. It is never written into the zip's possibly shared folder. No `sh-agent-verify-*` was left in `/tmp` after the E2E or my probes. A SIGKILL would leave only the public signer line.
- **Unsigned path (asked):** it still prints `UNVERIFIED`. `TestUpdateFromZipUnsignedSaysUnverified` passed, and it also asserts that `VERIFIED against` is absent.
- **E2E `-ldflags -X` (asked):**
  - The scratch key (`.agents/scratch/e2e-signed/keys/k`, 0600, gitignored by `.gitignore:40`) signed the E2E `SHA256SUMS`. Its public key (`…AAAAILR1…`) is not the release key (`…AAAAIBkx0…`).
  - The operator's release key signed nothing in this range.
  - The scratch-signer binary exists only under `.agents/scratch/e2e-signed/` and its `h1/.local/bin/` copy. The real `~/.local/bin/sh-agent` contains no `e2e@signer` (`grep -a -c` = 0; release-signer literal present). See the fourth finding for why a byte grep proves presence, not which signer is in effect.
  - The E2E left `rel/styrman-harness-v0.0.1.zip` byte-flipped (`orig.zip` is the genuine one). That is harmless, but anyone re-running from `rel/` gets `checksum mismatch`.
- **Secrets:** none in the diff or the 4 untracked files, all read end to end. The signer line and email are public-key material already committed. BUG-094's record names host "dave" and metrics only. Tracked-tree scan:
  - Secret-shaped paths: 3 hits, all docs whose names contain "secret" or "token".
  - Key and token shapes: only fake fixtures in `internal/harness/report_test.go` and docs describing redaction.
- **Instruction-shaped text:** none in the trail files beyond normal plan and review prose. The one instruction that matters is the refusal hint (second finding).
- **Injection:** `exec.Command` with no shell. Paths are option arguments to `-f` and `-s`, so a `-`-leading folder cannot become a flag. `SHA256SUMS` is parsed only after the signature holds.
- **Authz/authn:** no new auth surface. The trust anchor is unchanged, and `ssh-keygen` is still located system path first.
- **Crypto:** nothing hand-rolled. It is OpenSSH `ssh-keygen -Y verify` (SSHSIG) and stdlib `crypto/sha256` used as designed.
- **Dependencies:** none added (`go.mod` untouched). `govulncheck ./...` on go1.27.1: "No vulnerabilities found." (exit 0). `npm audit`: not applicable, since no JS package is in the diff.
- **GCP surface:** not touched.
- **`git archive` step (asked):**
  - It archives the tag's tracked tree only, so ignored `.agents/scratch/` (the release folder, keys, `allowed_signers`) cannot enter it.
  - There is no `.gitattributes`, so no `export-subst` or `export-ignore`.
  - The content equals GitHub's source zip for the tag. The repository is private and assets need repository read access (the same section).
  - The one tracked symlink, `.agents/skills`, lies outside the five trees and is skipped.
  - Nothing secret is published; see the last finding for the optional minimization.
- **Tools and actions to disclose:** probes ran in `.agents/scratch/security-signed-zip/` with the E2E's scratch-signer binary, all `--dry-run` into a scratch home, which stayed empty. I also invoked `~/.local/bin/sh-agent update --from-zip … --home <scratch> --dry-run` once. It is an older build without `--from-zip` and exited at flag parsing ("flag provided but not defined"), so nothing ran.

## Not examined

No release-diff path went unread: every hunk of the 13 modified files and all 4 untracked files were read end to end. Not examined, outside those hunks:

- **Unchanged code in touched files:**
  - `cmd/sh-agent/main.go` outside `updateFromZip`: `installFiles`, `sweepStaleZipDirs`, argument parsing including BUG-094's `zipSet` refusal.
  - The test helpers `newRelease`, `signRelease`, `markedZip` and `isolate`, seen by name only. `isolate`'s PATH isolation was not checked.
  - `release_test.go`. I did not run the `TestUpdateFromRelease*` tests, so the claim that the `verifyRelease` refactor behaves as before rests on reading the diff.
- **Docs outside the diff hunks:** README, REQUIREMENTS, DECISIONS and `delivery.md` outside the changed lines; `remote.md` outside § Releases.
- **Not run:** the full suite, `go vet` and `gofmt` (the verifier's job).
- **Not exercised:** macOS (ssh-keygen location, `$TMPDIR`), and the real v1.9.0 zip under the real key, which does not exist yet.
- **Not verified:** BUG-094's claim that feedback-repo commit `1703cea` removed the pushed report (not fetched).
