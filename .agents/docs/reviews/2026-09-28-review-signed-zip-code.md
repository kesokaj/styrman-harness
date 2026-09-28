---
owner: reviewer
date: 2026-09-28
---
# Review — signed-zip-code   (2026-09-28)

Range reviewed: 0b762db..worktree (`git diff 0b762db` plus untracked `.agents/docs/plans/2026-09-28-signed-zip.md`, `.agents/docs/testing/2026-09-28-verify-signed-zip-e2e.md`)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The code is sound. The split keeps the release path's behaviour, verification runs before unpacking, and the signature, sums and zip bytes each come from one read. The blocker is in the docs. The diff changes R9's Release procedure: releases now publish a zip that a zip place verifies. It leaves the README map untouched. REQUIREMENTS.md says "the four procedures and the map change together … one commit carrying the code, the docs and the map". Project AGENTS.md rule 3 calls a release-behaviour change that does not touch the map incomplete. Two more restatements still describe the old procedure: the deploy guide's "five assets" and AGENTS.md/README "binaries with signed checksums". On the code side, no test covers the VERIFIED success path; a compiled mutation proves it. The stale-sums refusal also gives the wrong remedy.

## Findings

- **[high / confidence: high]** `README.md:20,37,40` (map) vs `REQUIREMENTS.md:46,48` and `AGENTS.md` rule 3. R9's Release bullet now says releases publish "a zip of its tree". The map still labels the node `Releases<br/>platform binaries + signed SHA256SUMS` and the edge `preflight -- "signed binaries" --> releases`. The only zip line still runs from `harnessRepo` (`"zip, downloaded by hand; …"`). No line shows a release zip, verified, reaching the zip place. `TestHarnessMap` checks verbs only, so it cannot catch this. **Fix, in this same commit:**
  - relabel `releases` (binaries + zip + signed SHA256SUMS) and the preflight edge;
  - add a dashed `releases -. "release zip + SHA256SUMS, downloaded by hand; --from-zip verifies" .-> zipPlace` edge (dashed because the download is by hand);
  - keep the `main`-zip line as the unverified path.

  Per rule 3, this rerouting goes to the operator. ADR-039's "All approved" covered R9 and the release asset; the finding cannot tell whether that included the map.

- **[medium / confidence: high]** `.agents/docs/deploy/remote.md:179`. The prose after the upload block still says "lists the five assets". The diff changed the block's own comment to "all six" at line 170. `/preflight` executes this guide as written, and two counts that disagree are the kind of thing that made the v1.2.0 gate stop (`reviews/2026-09-18-preflight-v1.2.0.md` finding 1). **Fix:** "six assets".

- **[medium / confidence: high]** `cmd/sh-agent/main.go:1091` (`verifiedTag = tag`) has no test.
  - **Mutation:** in a scratch copy (`.agents/scratch/rev-signed-zip/tree`), replace `verifiedTag = tag` with `_ = tag`. The mutation had exactly one occurrence, compiled, and was restored from a saved copy; `cmp` confirmed the restore.
  - **Result:** all 12 `TestUpdateFromZip*` tests ran (`-v`, none skipped) and PASSED. The command's headline success output, `VERIFIED against the signed SHA256SUMS of release vX.Y.Z`, is therefore untested. A regression that prints UNVERIFIED for a verified zip, or the reverse wiring, would ship green.
  - The plan's Phase 1 row names `TestUpdateFromZipVerified` and `TestUpdateFromZipTampered`. Neither exists. The E2E log (`testing/2026-09-28-verify-signed-zip-e2e.md:24`) names the tests that do exist, but it covers this path only by hand.
  - `TestUpdateFromZipUnsignedSaysUnverified` passes on the pre-change code too. It pins D3 (unsigned is not refused) and is not a red-first test.
  - **Fix:** a CLI test that signs `markedZip` with a key it sets as `harness.releaseSigner`, or via an exported test hook, in the release layout. It asserts `VERIFIED against the signed SHA256SUMS of release v…` and a populated home.

- **[medium / confidence: high]** `internal/harness/release.go:444-450` (`VerifyZipSums`, the `checksumFor` and mismatch errors).
  - **What the docs promise:** ADR-039 says "Any failure refuses the zip, including a stale or foreign pair of files; the message says to move them away". Plan D7 says the same.
  - **What the code does:** it adds that hint only to `verifySignedSums` errors. A genuine but STALE pair gets no hint. The likely case is a Downloads folder holding the previous release's `SHA256SUMS`/`.sig`, where the browser saved the new pair as `SHA256SUMS (1)`, or a `main` zip saved beside an old pair. That case fails `checksumFor` with `no line names styrman-harness-v1.9.1.zip — use the release's own file name`, which is wrong advice: the zip already has its own name. A mismatch gives `checksum mismatch …` with no remedy either.
  - **Fix:** name the release the sums belong to and append the same move-away hint to both errors. Add a unit case: sums signed for tag A beside a zip named for tag B.

- **[low / confidence: high]** `AGENTS.md:14` (rule 2) and `README.md:64` still say `/preflight` "publishes that tag's binaries with signed checksums". This no longer matches R9's amended Release bullet. **Fix:** "binaries and a zip of its tree".

- **[low / confidence: medium]** `REQUIREMENTS.md:43` vs `internal/harness/release.go:427`. R9 says a zip "without them is installed marked UNVERIFIED". "Them" is the pair, so a zip with only ONE of the two files arguably falls under that sentence. The code (plan D7, the safer reading) refuses it. **Fix:** make R9 and ADR-039 say it explicitly ("…when either sits beside the zip…"); the R9 edit needs operator approval. The case where `.sig` is present but `SHA256SUMS` is missing has no unit test; only the reverse, `sig missing`, is tested. With `ssh-keygen` absent, the one-file case reports the missing verifier rather than the missing file, contrary to the plan's failure table. That is harmless because it still refuses.

## Dispatch checks

- **`verifyRelease` → `verifySignedSums` split.** Behaviour is preserved: the same order (stat sums, read capped, stat sig, signer identity, write `allowed_signers` into the release dir, verify, marker), then the tag compare and the per-asset check in `verifyRelease`. Every message is byte-identical through `what = "release "+tag`, except `names no release version`, which now reads "of release vX" instead of "of vX". `release_test.go:372` pins only the substring, and the change is harmless. The targeted `TestUpdateFromRelease*`/`TestVerif*` tests pass with `-race`.
- **`VerifyZipSums` decision table:**
  - neither file → `ErrZipUnsigned`;
  - either file alone → refused (`has no SHA256SUMS.sig` / `has no SHA256SUMS`);
  - both → signature, marker, own-name line, compare.

  The one-file fail-open mutation (`&&`→`||`) fails `TestVerifyZipSums`, and so does the dropped checksum compare. Both compiled and a named test failed.
- **`UnpackZip`'s check hook** runs after hashing the single buffer and before `zip.NewReader` and every write. A mutation that ignores the check's error fails `TestUnpackZipCheckStopsBeforeWriting`. The caller's `MkdirAllTracked(parent)` runs before the check, and `RemoveCreated` removes it on refusal; `TestUpdateFromZipForeignSignatureRefused` asserts an empty home.
- **The `updateFromZip` closure** is correct: `ErrZipUnsigned`→nil, other errors propagate, success records the tag. Only the success branch is untested (finding 3).
- **The temp signers directory** comes from `os.MkdirTemp("", "sh-agent-verify-*")` with `defer os.RemoveAll`, and it is not written into the user's download folder. It is cleaned on every return path. No test asserts the cleanup; the E2E run observed it by hand.
- **R9 and ADR-039 vs the code** match, apart from findings 4 and 6.
- **`deploy/remote.md` step 1b.**
  - The `git archive` prefix and name match what `VerifyZipSums` looks up: it needs `filepath.Base` = `styrman-harness-vX.Y.Z.zip`, which matches the `sha256sum` line.
  - The upload loop lists six assets, and the block comment says "all six".
  - The prose at line 179 still says five (finding 2).
- **BUG-094 isolation:**
  - The `internal/harness` tests use only a temp key and `ssh-keygen`, with no network.
  - The CLI tests call `isolate`, which empties `PATH` and uses a scratch cwd. Neither reachable path calls `report`/`gh`/`git`.
  - `ssh-keygen` comes from `sshKeygenSystemPaths`, not `PATH`.
  - `HOME` is not isolated, but every run passes `--home`, and a refusal regression would install into `t.TempDir()`, not push anything.

## Criteria
1. **Dead code: pass.** No leftovers. The removed dated review citations were comments only.
2. **Flow tracing: pass with one gap.** The paths from `updateFromZip` through `UnpackZip` to the check, then `VerifyZipSums` → `verifySignedSums` → `checksumFor`, are all consistent. The errors propagate unwrapped where the tests expect it. The verified-output branch is untested (finding 3).
3. **Standards conformance: fail.** Restatements of the changed Release procedure were left standing: the map, remote.md:179, AGENTS.md:14 and README.md:64. `gofmt -l` is clean and `go vet` passes.
4. **Directory layout: pass.** `VerifyZipSums` sits beside `verifySignedSums` in `release.go`, as the plan says, and the tests sit beside the code.
5. **Simplicity (DIY-except-security): pass.** The code signs and verifies with OpenSSH `ssh-keygen -Y` and uses stdlib `sha256`; nothing is hand-rolled, and there is no new dependency.
6. **Decoupling: pass.** `UnpackZip` takes a plain `func(sum string) error` hook. It does not import verification.

## Not examined
- The full `go test -race ./...` suite. I ran only the targeted packages and patterns (`TestUpdateFromRelease|TestVerif|TestVerifyZipSums|TestUnpackZip|TestUpdateFromZip|TestHarnessMap`), all green.
- `.agents/docs/bugs/BUG-094-*.md` and `INDEX.md`: I checked the status flip for consistency only, not the claim that commit `1703cea` exists in the feedback repository.
- The cold-read reviews r1/r2: I read them only through grep hits.
- The E2E run itself was not reproduced.
- macOS was not examined.
