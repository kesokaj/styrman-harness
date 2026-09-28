---
owner: reviewer
date: 2026-09-28
---
# Review — signed-zip-code, round 2   (2026-09-28)

Range reviewed: 0b762db..worktree (`git diff 0b762db`: README.md, AGENTS.md, REQUIREMENTS.md, ROADMAP.md, DECISIONS.md ADR-039, deploy/remote.md, specs/sh-agent/delivery.md, cmd/sh-agent/{main.go,zip_test.go}, internal/harness/{release.go,zip.go,zip_test.go}; plus the untracked plan `.agents/docs/plans/2026-09-28-signed-zip.md`, read by grep)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Round 1's three blocking items are fixed. The map, its `linkStyle` indices and the six-asset count are all right, and the VERIFIED path now has a test that fails under the `verifiedTag` mutation (reproduced below). Two problems remain.

1. **Release step 5 runs too late, and only on linux/amd64.** It runs after `gh release upload`, so the binary it checks has already shipped. That contradicts ADR-039's "so a binary built with a stray test signer never ships". It also executes `./sh-agent-linux-amd64` unconditionally, which cannot run on a Mac or an arm64 box that holds the signing key.
2. **Security finding 2 is fixed in code but not pinned.** Re-inserting the exact old move-away hint compiles and leaves every test green. The only negative assertion sits on the checksum-mismatch path, which never carried the hint.

The stale-pair remedy from round-1 finding 4 still gives the wrong advice: it tells the operator to fetch the OLD release's zip.

## Findings

- **[medium / confidence: high]** `.agents/docs/deploy/remote.md:171-173` (new step 5) runs after step 4's `gh release create` and `gh release upload` loop (lines 165-169).
  - **Timing:** by the time the check runs, the assets are public. A failure there means a release with a bad signer is already published. ADR-039 Consequences (`DECISIONS.md`, "so a binary built with a stray test signer never ships") and the step's own comment ("never ships") claim the opposite.
  - **Host:** the step executes `./sh-agent-linux-amd64`. This guide lets a release be cut on any machine holding the signing key (§ Releases, line 96; line 215 "before that machine cuts a release"), and the places include two Macs. On darwin/arm64 or linux/arm64 the command fails with an exec format error, and the gate stops only after the publish.
  - **Fix:** move step 5 between step 3 (sign) and step 4 (release create). Run the binary that matches the host, for example `./sh-agent-$(go env GOOS)-$(go env GOARCH)`, since all three share one `-ldflags`. Keep the output going into a file.
  - The rest of step 5 is sound. It writes to `../zipcheck.out`, not a pipe. `--dry-run` runs the same `UnpackZip` check, so it prints the trust line. `SHA256SUMS`/`.sig` sit in the same directory as the zip, which is the cwd. The dry run on the fresh `../zipcheck` home unpacks under the system temp directory (main.go:1057-1060).

- **[medium / confidence: high]** `cmd/sh-agent/zip_test.go:694` and `internal/harness/release.go:447-449`: removing security finding 2's hint is not pinned by any test.
  - **Mutation:** in the scratch copy `.agents/scratch/rev2-signed-zip/tree`, the target text occurred exactly once. I made the `verifySignedSums` error return in `VerifyZipSums` `fmt.Errorf("%w (move SHA256SUMS and SHA256SUMS.sig out of %s to install the zip unverified on purpose)", err, dir)`, which is security's quoted round-1 text.
  - **Result:** `go vet` was clean, so the mutation compiled. Every `TestUpdateFromZip*` test and `TestVerifyZipSums` PASSED under `-v`. I restored from the saved copy, and `cmp` confirmed the restore.
  - **Why the test cannot catch it:** the `move SHA256SUMS` needle is checked only after the checksum-mismatch run, and the mismatch error never carried the hint, even in round 1. The paths that did carry it are the bad signature (`TestUpdateFromZipForeignSignatureRefused`, the `foreign key` subtest) and a missing `SHA256SUMS`. Those assert only a positive substring.
  - Per the global DoD, a finding is fixed only when each trigger it names has a test that fails without the fix. Security named two triggers: signature failure and `.sig` without `SHA256SUMS`.
  - **Fix:** assert that no `move` hint appears on the foreign-signature refusal and on a `.sig`-only folder. Neither case exists as a unit test yet; see round-1 low 6.

- **[medium / confidence: high]** `internal/harness/release.go:453-455`: the stale-pair remedy is still wrong. This is round-1 finding 4, only half fixed.
  - **The case:** genuine v1.9.0 `SHA256SUMS`/`.sig` left in `~/Downloads`, next to a newly downloaded `styrman-harness-v1.9.1.zip` (the browser saved the new pair as `SHA256SUMS (1)`). The error reads: `the signed SHA256SUMS beside it are release v1.9.0's and do not list styrman-harness-v1.9.1.zip: no line names styrman-harness-v1.9.1.zip — download that release's own zip, under its own name`.
  - **Why it is wrong:** "that release" is v1.9.0, so the advice is to fetch and install the OLDER release. That is a silent downgrade for an operator, or an agent, that follows tool output. The advice fits only the renamed-zip case, which is the one the unit test covers (`renamed zip`, which asserts `no line names`).
  - **Missing test:** the unit case round 1 asked for, sums signed for tag A beside a zip named for tag B, was not added.
  - **Fix:** name both possibilities, for example "download the SHA256SUMS and SHA256SUMS.sig of the release this zip belongs to, or that release's zip under its own name". Or parse `styrman-harness-(vX.Y.Z).zip` from the file name, and when it names a different version than the sums, say the sums are stale. Add the tag-A/tag-B case.

- **[low / confidence: high]** `.agents/docs/plans/2026-09-28-signed-zip.md:106` (D7) still says "The error names the files, so the operator can move them away to install unverified on purpose". That is the hint security finding 2 removed, and ADR-039 says the messages "never suggest removing the checksums". This is the active plan, which the next session reads on resume, and no later D-entry supersedes it. **Fix:** append a D-entry that withdraws D7's last sentence and cites the security review. Do not edit D7 in place if the log is append-only.

- **[low / confidence: high]** ADR-039 says "The messages name what to download", but the refusals that mean tampering do not. `SHA256SUMS signature does not verify (…)`, `<zip> has no SHA256SUMS`, `<zip> has no SHA256SUMS.sig` and the no-marker error carry no remedy at all. Security finding 2's first remediation bullet ("say the files are not genuine and to download the release's three files again") was not applied. **Fix:** either wrap those errors in `VerifyZipSums` with that sentence, or narrow the ADR's claim.

- **[low / confidence: high]** `internal/harness/release.go:414-425`: the `ErrZipUnsigned` doc comment was split from its declaration. `SwapReleaseSigner` was inserted between them, so `SwapReleaseSigner`'s godoc now opens with "ErrZipUnsigned: no SHA256SUMS … The caller installs it marked UNVERIFIED.", and the exported `ErrZipUnsigned` (line 425) has no doc. This is the describing-text rule. **Fix:** move the two comment lines back above `var ErrZipUnsigned`.

## Round-1 findings — status

1. **Map, high: FIXED.**
   - `README.md:20` relabels the Releases node "platform binaries + release zip + signed SHA256SUMS".
   - `:37` relabels the preflight edge "signed binaries and zip".
   - `:41` adds the dashed `releases -. "release zip + signed SHA256SUMS, by hand; verified" .-> zipPlace`.
   - The unverified `main`-zip line at `:40` is kept.
   - **`linkStyle` checked against the edge order.** There are 21 edges, indices 0-20:
     - 0 edit→ship
     - 1 ship→harnessRepo
     - 2 ship→preflight
     - 3 preflight→harnessRepo
     - 4 preflight→releases
     - 5 harnessRepo→checkoutPlace
     - 6 releases→binaryPlace
     - 7 harnessRepo→zipPlace
     - 8 releases→zipPlace
     - 9-11 the three →installer edges
     - 12 installer→homeFiles
     - 13 homeFiles→sessions
     - 14-20 the report loop
   - Teal `0,1,5,6,7,8`, plum `2,3,4`, indigo `9-13` and amber `14-20` cover every index exactly once. The new edge is teal like its sibling releases→binaryPlace. Every pre-change index shifted by exactly +1 from index 8 on.
   - The rest of finding 1 is fixed too:
     - `AGENTS.md:14` rule 2 now says "binaries and a zip of its tree, with signed checksums";
     - README Release (`:65`) is fixed;
     - remote.md now says "all six" (`:170`) and "the six assets" (`:182`).
   - One gap: `TestHarnessMap` still does not pin labels, so a future relabel drift stays uncaught. That is the existing design, not a finding.
2. **VERIFIED path untested, medium: FIXED AND TESTED.** `TestUpdateFromZipVerifiedAndTampered` (`cmd/sh-agent/zip_test.go:640`) goes through `harness.SwapReleaseSigner`. I reproduced the mutation in a scratch copy: `verifiedTag = tag` occurred once, was replaced with `_ = tag`, compiled, and made the test fail at `zip_test.go:682` ("a verified zip does not say so"). The file was restored from a saved copy and `cmp` matched. The tampered half flips one byte, gets `checksum mismatch`, and asserts an empty home.
3. **Remedy messages, medium: PARTLY FIXED.**
   - No "move the sums away" hint remains in the code (`grep` over `cmd`, `internal`, README, AGENTS, REQUIREMENTS, ADRs, deploy, specs, context, skills, agents).
   - Three gaps remain, each a finding above:
     - the plan's D7 still states the hint;
     - the removal has no test that fails without it;
     - the stale-pair message now gives downgrade advice.

Round-1 lows:
- `AGENTS.md`/README "binaries with signed checksums": fixed.
- R9's one-file wording and the missing `.sig`-only test: not addressed. They are not re-raised, but the `.sig`-only test is now also needed by the second finding.

## Other checks asked for
- **ADR-039's Open bullet** (`DECISIONS.md`, Consequences) records honestly that a zip with nothing beside it still installs with exit 0. It marks this as an operator question requiring an R9 change, and it cites the security medium. This matches the security report's first finding, and nothing in the diff claims otherwise. README's zip section says `main` zips are unsigned. delivery.md now states the verified/unverified split once, and README/AGENTS point to it: rule 6 holds.
- **REQUIREMENTS R9** was edited with a recorded approval ("All approved", ADR-039).
- **Gates:** `gofmt -l cmd internal` is clean, `go vet ./...` is clean, and `go test -race -count=1 ./...` passes for all three packages. No `t.Parallel` is in `cmd/sh-agent`, so the global signer swap does not race.

## Criteria
1. **Dead code: pass.** `SwapReleaseSigner` has a caller (the CLI test), and `ErrZipUnsigned` is used by main.go. Nothing is orphaned.
2. **Flow tracing: pass with gaps.** `updateFromZip` → `UnpackZip(check)` → `VerifyZipSums` → `verifySignedSums`/`checksumFor` is consistent, `ErrZipUnsigned` maps to UNVERIFIED, and every other error refuses before any write. The gaps: one test cannot fail (the second finding), and one remedy message is wrong (the third finding).
3. **Standards conformance: fail (low/medium).** Restatements left standing: the plan's D7, and ADR-039's "messages name what to download". The godoc is split from its identifier. The formatting and vet gates are verified as run.
4. **Directory layout: pass.** Nothing new lands outside the existing files.
5. **Simplicity (DIY-except-security): pass.** Signatures use OpenSSH `ssh-keygen -Y`, checksums use stdlib `sha256`, and there is no new dependency. `SwapReleaseSigner` is an exported test hook in production code. It is reachable only inside this module (`internal/`), and security round 1 already discussed the hook trade-off, so it is not a finding here.
6. **Decoupling: pass.** `UnpackZip` takes a plain `func(sum string) error`, and verification lives in `release.go`.

## Not examined
- The plan file, the cold-read reviews and the E2E/verify logs were read only through grep. I did not check the plan's tick state against the work.
- The BUG-094 record and `INDEX.md` diffs were not read in this round.
- The release procedure was not executed. The step 5 finding comes from reading the guide's step order and host assumptions, not from a run on a Mac.
- macOS was not examined.
