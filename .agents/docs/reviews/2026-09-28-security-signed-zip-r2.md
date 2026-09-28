---
owner: security
date: 2026-09-28
---
# Security review — signed-zip, round 2   (2026-09-28)

Range reviewed: 0b762db..worktree

The tree changed while I was reviewing it. Between 19:18:34 and 19:19:03 the main agent edited `internal/harness/zip_test.go`, `internal/harness/release.go`, the plan, `DECISIONS.md` and `deploy/remote.md`, and `reviews/2026-09-28-review-signed-zip-code-r2.md` appeared. This verdict covers the tree as hashed at 19:19:51 (sha256, first 16 hex characters):

| File | sha256 |
|---|---|
| `release.go` | `582f817363552016` |
| `zip.go` | `df7e65f94b9a2fcc` |
| `internal/harness/zip_test.go` | `ac8555746309fc49` |
| `main.go` | `b2986d0ffcc19269` |
| `cmd/sh-agent/zip_test.go` | `8aff9491a5fadee5` |
| `DECISIONS.md` | `88f94bd7cc3233f9` |
| `remote.md` | `ad200daae09b6210` |
| `README.md` | `d7ab842e418379da` |
| plan | `c5d3045d38fc3c8d` |

I re-hashed these at 19:22:45 and none had changed. The full list is in `.agents/scratch/sec-r2/hash-final.txt`. Anything edited after that time is not covered.

## Status: COMPLETE

## Verdict: FINDINGS (none blocking this change)

**The worst problem is round 1's medium, and it is still open.** A zip handed over without its `SHA256SUMS`/`SHA256SUMS.sig` still installs with exit 0, marked `UNVERIFIED`. The code for that path has not changed. ADR-039 now records it honestly as **Open**, and the operator will be asked about an `--unverified` flag, which is an R9 change.

**Why it does not block this change.** It is acceptable as non-blocking:
- The change only narrows the exposure v1.8.0 shipped with, where every zip installed unverified.
- It adds no new way into the unverified path.
- Closing the path needs a REQUIREMENTS edit that only the operator can approve.

**Conditions.** It stays an open medium in the v1.9.0 gate report and becomes an operator question in the handoff. The plan does not yet say this. It reads `Follow-ups, not built: none` and `Risks: Assumptions: none open`, and its `## Why` still calls this change "the fix".

**The other two asked items are done.**
- **Bypass hint:** it is gone from every `VerifyZipSums` message and from the plan's D7. A compiled mutation shows the tests now catch its return on the bad-signature and missing-`.sig` paths. Two gaps remain. The other trigger round 1 named (`.sig` present, `SHA256SUMS` absent) is still untested. A reworded hint also passes every test.
- **"Six assets":** fixed. Release step 5 is now step 3b. It runs before `gh release create`, on the host's own binary, and I judge it sound.

**`SwapReleaseSigner`** cannot be called at runtime. It is not in the production binary (measured with `go tool nm`).

Everything else is low.

## Findings

- **[medium / confidence: high] Carried from round 1. OPEN by decision; not blocking this change.**
  - **Where:** `internal/harness/release.go` `VerifyZipSums` (`sumsErr != nil && sigErr != nil` → `ErrZipUnsigned`), together with `cmd/sh-agent/main.go:1085-1087`, which turns that sentinel into an install.
  - **What is wrong:** whoever supplies the zip decides whether it is verified, because they can leave the two files out. The weakness and the remediation are unchanged from round 1: an explicit `--unverified`, as an R9 change.
  - **Judgement: acceptable as non-blocking:**
    - (i) Every path this change adds refuses or verifies. The only path to an unverified install is the one v1.8.0 already had, so exposure strictly shrinks.
    - (ii) The fix changes R9, which is the operator's to approve.
    - (iii) ADR-039's Consequences now state the residual and cite the medium.
  - **Conditions:**
    - (a) The v1.9.0 gate report carries it as an open medium, narrowed but not closed. The v1.8.0 row must not be closed.
    - (b) The handoff's Next Steps asks the operator the `--unverified` question.
    - (c) Add one line to the plan, `.agents/docs/plans/2026-09-28-signed-zip.md`, pointing at ADR-039's Open bullet. It is needed because three places still contradict ADR-039:
      - `## Why` still quotes the handoff calling this "the fix";
      - `Follow-ups, not built: none`;
      - `Risks: Assumptions: none open`.

      A gate writer reading the plan would close the row.

- **[low / confidence: high] Round-1 finding 2's removal is only partly pinned.** `internal/harness/zip_test.go:460-525` (`TestVerifyZipSums`).
  - **Setup:** three mutations of the `verifySignedSums` error return in `VerifyZipSums`, in the scratch copy `.agents/scratch/sec-r2/tree`. Each compiled (`go vet` exit 0) and ran against `TestVerifyZipSums` and `TestUpdateFromZip*`.
  - **Blanket:** round 1's exact hint on every `verifySignedSums` error. `TestVerifyZipSums/sig_missing` and `TestVerifyZipSums/foreign_key` FAIL, so this one is caught.
  - **Narrow:** the same hint only when `SHA256SUMS` is absent and `.sig` present. This is the second trigger round 1 named. Every test PASSES, so it survives. No subtest covers a `.sig` with no `SHA256SUMS`.
  - **Reworded:** `(delete the checksum files to install the zip unverified)` on every error. Every test PASSES, so it survives. The negative needle `move SHA256SUMS` matches one phrasing only.
  - **Remediation:**
    - Add a subtest, "sums missing, sig present": it refuses, is not `ErrZipUnsigned`, and names `SHA256SUMS`.
    - Make the negative assertion case-insensitive `unverified` on every refusal subtest. A refusal has no reason to mention an unverified install.
  - The copy was restored from a saved file (`cmp` clean). The live `release.go` was never touched: its hash was `582f817363552016` before and after.

- **[low / confidence: high] README steers readers to the unverified path.** `README.md` § Install from a zip, the example block.
  - **What is wrong:** the only command shown is `sh-agent update --from-zip ~/Downloads/styrman-harness-main.zip`. It comes straight after "A zip of `main` … works too, but it is unsigned:". The recommended path, the release zip plus two files, has no example.
  - **Why it matters:** anyone who copies the command, a person or an agent, takes the `UNVERIFIED` path the first finding leaves open.
  - **Fix:** make the example `~/Downloads/styrman-harness-vX.Y.Z.zip`, with `SHA256SUMS` and `SHA256SUMS.sig` in the same folder, and show the `main` zip second with its `UNVERIFIED` note.

- **[low / confidence: medium] Downgrade. Carried from round 1; neither addressed nor recorded.** `cmd/sh-agent/main.go` `updateFromZip`.
  - **What is wrong:** a genuine older release zip, with its genuine pair, verifies and installs.
  - **What changed:** the older version is now at least visible in `VERIFIED … of release v<older>`.
  - **What is still missing:** nothing compares it with `buildVersion()`. ADR-039 does not mention it, and no report records it as fixed or rebutted.
  - **Remediation (unchanged):** when `IsReleaseVersion(buildVersion())`, warn or refuse on `compareRelease(verifiedTag, buildVersion()) < 0`. Both are existing helpers. Alternatively, rebut it in an Outcome.

- **[low / confidence: high] Stat errors fall through to unverified. Carried from round 1; neither addressed nor recorded.** `internal/harness/release.go` `VerifyZipSums`.
  - **What is wrong:** any `os.Stat` error on both files, not only a missing file, falls through to `ErrZipUnsigned` and an unverified install. A verification guard whose check errors should refuse instead.
  - **Fix (unchanged):** `errors.Is(err, fs.ErrNotExist)` on both, via `os.Lstat`.

- **[low / confidence: high] `SwapReleaseSigner`: hardening, not a vulnerability.** `internal/harness/release.go:414-421`.
  - **Answer to the question:** nothing can call it at runtime; the evidence is under Checked.
  - **Residual 1:** nothing stops a future non-test caller from compiling.
  - **Residual 2:** the comment at `release.go:34-35`, "It is a variable only so a test can substitute a test key's line", is stale. Package tests assign the variable directly, the CLI test goes through `SwapReleaseSigner`, and the variable is also what `-ldflags -X` replaces. That is the case the E2E relied on and step 3b now guards.
  - **Fix:**
    - A test that parses the module's non-`_test.go` files (stdlib `go/parser`, no dependency) and fails on any reference to `SwapReleaseSigner` other than its declaration.
    - A comment that names `-X` and step 3b.
  - It is not safe for concurrent use: an unsynchronised write to a package variable. That is fine while `cmd/sh-agent` has no `t.Parallel`, which reviewer r2 checked.

- **[low / confidence: medium] Release zip carries the whole tree. Carried from round 1; optional.** `.agents/docs/deploy/remote.md` step 1b.
  - **What it does:** `git archive … vX.Y.Z` puts the whole tracked tree, `.agents/docs/` included, into a release asset.
  - **Why it is low:** the content equals GitHub's own source zip, and it is readable only by accounts that can already read the repository.
  - **Optional minimisation (unchanged):** archive only `context skills agents commands rules`.

## Round-1 findings — status

| Round 1 | Status at this snapshot |
|---|---|
| medium: sums can be omitted | **Open by decision** (ADR-039 Consequences, "Open"). Code unchanged. Acceptable as non-blocking under the conditions in the first finding. |
| medium: bypass hint in refusals | **Fixed** in code and trail. `VerifyZipSums` returns `verifySignedSums` errors unwrapped. The plan's D7 is amended, and ADR-039 says messages never suggest it. Grep over the tree, excluding scratch and dated reviews: `move SHA256SUMS` survives only as three negative test needles. **Partly pinned** (second finding). |
| low: downgrade | **Open**, unrecorded (fourth finding). |
| low: `-X` trust anchor, no built-binary check | **Fixed** by step 3b; see below. The comment nit is in the sixth finding. |
| low: "five assets" | **Fixed.** `remote.md:170` "all six", `:182` "the six assets", ADR-039 "six assets". No living doc still says five; the remaining hits are dated handoffs, plans and retros. |
| low: test coverage | **VERIFIED branch fixed** by `TestUpdateFromZipVerifiedAndTampered`. Reviewer r2 reproduced the `verifiedTag` mutation, and 13 `TestUpdateFromZip*` PASS here. **`.sig`-only case still open** (second finding). |
| low: stat errors read as "unsigned" | **Open**, unrecorded (fifth finding). |
| low: `git archive` scope | **Open**, optional (last finding). |

### The three asked items

**1. Sums can be omitted.** Non-blocking, with conditions (a) to (c) in the first finding.

**2. Bypass hint removed from every `VerifyZipSums` message.** Confirmed. The message for each case:

| Case | Message |
|---|---|
| Signature fails; `.sig` or `SHA256SUMS` missing; no marker | `verifySignedSums`'s own text, with no hint |
| `ssh-keygen` missing | says the zip "cannot be verified without it" |
| Stale pair or unlisted name | "put the zip in a folder with its own release's SHA256SUMS and SHA256SUMS.sig, under the release's own file name" |
| Mismatch | "download it again" |

The stale-pair text no longer sends anyone to an older release; that was reviewer r2's downgrade point. It is pinned by the `stale sums of another release` subtest.

That text could lead someone to move a zip into an empty folder and install it `UNVERIFIED`. This gives an attacker nothing beyond the first finding:
- The message appears only after the signature has verified, when the file name is not listed.
- An attacker who controls the folder can simply leave the sums out.

**3. "Six assets" and release step 5.** "Six" is fixed. Step 5 became **step 3b** at 19:18:52 (`remote.md:160-163`), and I judge it sound:
- **Placement:** it runs after step 3's signature check and before `gh release create`. A failure stops the release before anything is published. The tag is already pushed at that point, but it is correct; only the binaries need rebuilding.
- **Which binary:** it runs `./sh-agent-$(go env GOOS)-$(go env GOARCH)`, so it no longer requires a linux/amd64 host. A host with no matching build, or a mismatched `GOENV`, fails closed before publishing.
- **Dry run:** `--dry-run` into scratch `../zipcheck` runs the same `UnpackZip` check. Round 1 measured that the verification branches run into an empty scratch home under `--dry-run`.
- **Output capture:** stdout goes to a file (`>` truncates, so a binary that fails to start leaves an empty file, and `grep` fails). The `grep` needle includes the tag and cannot match the `UNVERIFIED` line.
- **What it proves, together with step 3:** step 3's `ssh-keygen -Y verify` shows the signature is by the documented line. Step 3b shows the host binary's compiled-in line accepts that signature. A stray `-X` with another key therefore stops the release.
- **Coverage:** it checks one binary. The other two come from the same loop with the same `-ldflags`, and ADR-039 now says "the host's own built binary".
- **Limit:** I judged step 3b by reading it, not by running it.

**`SwapReleaseSigner` at runtime.** Nothing can call it:
- `internal/` makes it importable only inside this module.
- Its only caller is `cmd/sh-agent/zip_test.go:675`, by grep over the tree with scratch excluded.
- Measured: a production build (`go build -trimpath ./cmd/sh-agent`, snapshot copy) has **0** `SwapReleaseSigner` symbols in `go tool nm`. As controls, `harness.VerifyZipSums` has 2 and `harness.releaseSigner` has 1. The linker drops the function, so no path in the shipped binary reaches it.
- `-ldflags -X` sets only string variables and cannot reach a function.
- No plugin or reflection path exists to a package-level function.

It therefore adds no attack surface beyond the existing `var releaseSigner`, which step 3b now guards. Its godoc was re-separated from `ErrZipUnsigned` at 19:19:03.

## Checked

- **Secrets:** none. I read the diff and all untracked files. A pattern scan over the diff and untracked files found no secret shapes (private-key armour, `ghp_`, `github_pat_`, `AKIA`, `xox*`, `AIza`, `sk-`, `"private_key"`). The signer line and email are public-key material, already committed. Test keys are generated in `t.TempDir()`. The trail files carry only machine paths already present across the trail.
- **Instruction-shaped text:** the plan's D7 bypass sentence is gone (amended 19:18:52). Reviewer r1's finding 4, which recommended adding the move-away hint to more errors, is superseded on the record by reviewer r2 and ADR-039. I found no other agent-directed instruction in the trail files.
- **Injection:** no new exec surface. The `ssh-keygen` arguments are unchanged: option-argument paths, no shell. The tag in the `VERIFIED` line comes only from signed bytes, and `sumsVersionMarker` limits it to `v[0-9]+\.[0-9]+\.[0-9]+`.
- **Authz/authn:** the trust anchor is unchanged (the compiled-in line), and `SwapReleaseSigner` is absent from the binary (see above).
- **Crypto:** nothing hand-rolled. It uses OpenSSH `ssh-keygen -Y verify` (SSHSIG) and stdlib `crypto/sha256` as designed.
- **Dependencies:** `go.mod` untouched. `govulncheck ./...` in the snapshot copy printed "No vulnerabilities found." (exit 0). `npm audit` does not apply: there is no JS package.
- **GCP surface:** not touched.
- **Tests I ran (snapshot copy, `-race -count=1`):**
  - `internal/harness`, `-run 'TestVerifyZipSums|TestUnpackZip|TestUpdateFromRelease|TestVerif'`: ok. This also covers the `verifyRelease` → `verifySignedSums` split, which round 1 did not run.
  - `cmd/sh-agent`, `-run TestUpdateFromZip` with `-v`: 13 PASS, 0 FAIL, 0 SKIP.
- **Coverage of the release diff:**
  - **Read end to end by me:**
    - untracked files, whole: the plan (whole, then `## Why`, `## Failure modes`, `## Decisions` and lines 11-53 and 65-97 again after its 19:18:52 edit), `reviews/…-review-signed-zip-code.md`, `reviews/…-review-signed-zip-code-r2.md`, `reviews/…-security-signed-zip.md`, `testing/…-verify-signed-zip.md`, `testing/…-verify-signed-zip-e2e.md`;
    - every hunk of the 14 modified files, re-read where they moved: ADR-039 in full, `remote.md` § Releases lines 94-185, current `TestVerifyZipSums`, `release.go` 405-470;
    - `release.go` 25-50;
    - the whole of `updateFromZip`.
  - **Relied on round 1:** `reviews/2026-09-28-coldread-signed-zip-r1.md` and `-r2.md`. Both are unchanged since round 1 read them end to end. Their mtimes are 18:55 and 18:59, before round 1 took its range: its range lists the 19:03:40 E2E log but not the 19:05:59 verify log.
- **Actions to disclose:**
  - I created `.agents/scratch/sec-r2/`: a snapshot tree, a production binary, the `nm` dump, test, mutation and govulncheck logs, and hash lists.
  - I edited no live file, ran no `sh-agent` command against any home, and ran no git command that changes the tree.
  - After the probes the tree still showed 14 modified files, with `release.go` unchanged.

## Not examined

- **Parts of modified files outside their diff hunks:** `DECISIONS.md` apart from ADR-039; README, REQUIREMENTS and `delivery.md` outside the changed lines; `remote.md` outside lines 94-185; `main.go` outside `updateFromZip`.
- **Test helpers read only where they appear in the diff:** `isolate`, `markedZip`, `newRelease`, `signRelease`.
- **Not run by me:** the full suite, `gofmt` and `go vet ./...`. Those belong to the verifier; the verify log predates the 19:18-19:19 edits, so it needs a re-run.
- **Release step 3b** was not executed, and the real v1.9.0 zip under the real key does not exist yet.
- **macOS:** the `go env GOOS/GOARCH` form on a Mac host, `ssh-keygen` location, `$TMPDIR`.
- **Not verified:** BUG-094's claim that feedback-repo commit `1703cea` removed the pushed report. I did not fetch it.
- **Later edits:** anything changed in the tree after 19:22:45.
