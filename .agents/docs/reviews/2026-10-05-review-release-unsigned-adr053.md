---
owner: reviewer
date: 2026-10-05
---
# Review — release-unsigned-adr053   (2026-10-05)

Range reviewed: 18b9018..4ee867a

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing is broken. Every refusal that should remain is still there: no SHA256SUMS, a missing or wrong version marker, a checksum mismatch, never backwards, and an unreadable SHA256SUMS is refused rather than treated as absent. For each one, a mutation made the test that names it fail. The signing code is gone completely: no dead branch, no half-removed helper in the source, and `go vet` and `gofmt` are clean. The remote.md § Releases commands, the v2.4.0 bridge included, ran as written from `.agents/scratch/release-v2.4.0` in a scratch clone, and each one failed when it should. What remains are three low findings. The first is an orphaned test helper. The second is two describing-text restatements that still say releases are signed. The third is a test comment that claims more than the test proves.

## Findings
- **[low / confidence: high]** `internal/harness/release_test.go:88` — `mustRun` has no caller left (grep across `internal/` and `cmd/`: only its definition), because the ssh-keygen setup helpers that used it were removed. It is dead code: delete it.
- **[low / confidence: high]** `AGENTS.md:84` ("the signed-tag release procedure") and `.agents/docs/specs/sh-agent/delivery.md:83` ("`/preflight`, which signs the checksums with the same key as the tag") still describe releases as signed. ADR-053 reverses that claim. delivery.md is rule 6's one home for the delivery claim, so a present-tense signing statement there contradicts its own § The model at line 32. Fix: reword both lines. The rest of the restatement grep was clean: `context/`, `skills/preflight`, `agents/`, README, REQUIREMENTS, local.md and remote.md were already updated. `testing/sh-agent-matrix.md:66` is a dated run record and was left alone.
- **[low / confidence: medium]** `internal/harness/release_test.go:99` (and the matching comment on `TestUpdateFromZipVerifiedAndTampered` in `cmd/sh-agent/zip_test.go`) — the comments say that because PATH holds only the fake gh, "an update that needed [ssh-keygen] would fail here". That holds for a lookup on PATH. It does not hold for the verifier this range removed, which looked in `/usr/bin/ssh-keygen` first, regardless of PATH. A reintroduced `sshKeygen()` would pass these tests on any machine that has OpenSSH installed. Fix: narrow the comment to what PATH isolation actually proves.

## Evidence
- `gofmt -l cmd internal`: empty. `go vet ./internal/harness ./cmd/sh-agent`: clean. `go test -count=1 -run 'Release|Zip|VerifyZip|Checksum|NonCanonical|Embed|HarnessMap'` passed on both packages at 4ee867a. This is a targeted run, not the full `-race` suite.
- Mutations (`.agents/scratch/rev-adr053/mutate.sh`, run in a scratch clone at 4ee867a): each target occurred exactly once, the file was restored from a saved copy, `git status` was clean after each, and every mutated build compiled.
  - M1, the `named != tag` check removed: caught by `TestUpdateFromReleaseRefusals/checksums_of_another_release`.
  - M2, the missing-marker refusal removed: caught by `.../checksums_naming_no_version`.
  - M3, the binary hash compare removed: caught by `.../binary_does_not_match_its_checksum` and `.../SHA256SUMS_altered`.
  - M4, the zip hash compare removed: caught by `TestVerifyZipSums/checksum_mismatch`, `.../github_source_archive`, `TestUpdateFromZipWrongChecksumRefused` and `TestUpdateFromZipVerifiedAndTampered`.
  - M5, any stat error read as "no sums": caught by `TestVerifyZipSums/sums_unreadable`.
  - M6, the backwards guard removed: caught by `TestUpdateFromReleaseNeverGoesBackwards`.
  - M7, the missing-SHA256SUMS stat removed: caught by `.../no_checksums_in_the_release`.
  - M8, `main.go` swallowing a `VerifyZipSums` refusal: caught by both cmd zip tests.
  - M9, a `--pattern SHA256SUMS.sig` re-added to the download: not caught. This is harmless. The fake gh does not fail when a pattern matches nothing, and neither does the real gh while another pattern matches, so it changes no behaviour.
- Release rehearsal (`.agents/scratch/rev-adr053/rehearse.sh`, scratch clone, v2.4.0 tagged in the clone only, host binary only):
  - Steps 1, 1b, 2 and 3 each exited 0, and the `VERIFIED against the SHA256SUMS of release v2.4.0` grep matched.
  - Bridge fingerprint test: the documented `SHA256:5HpYQ…` is the fingerprint of the documented retired public key (`ssh-keygen -lf`). The test exits 1 for a foreign key and for an unset RETIRED.
  - `ssh-keygen -Y sign -f "$RETIRED"` works with no agent.
  - The verify against the documented retired line exits 255 for a foreign key and 0 with the matching line.
  - `rm … && test ! -e` exits 0.
  - `git -C "$(git rev-parse --show-toplevel)" archive v2.3.2 | tar -x -C ../v232src` works from the subdirectory. An absent tag makes the pipeline exit 2.
  - The v2.3.2 build exits 0. The v2.3.2 binary built as written refuses a foreign-key `.sig` ("SHA256SUMS signature does not verify"). Built with the test signer through `-X`, it prints the exact grep'd `VERIFIED against the signed SHA256SUMS of release v2.4.0` line.
  - With no `.sig`, a v2.3.2 binary installs the zip `UNVERIFIED`, as delivery.md:32 says.
  - The new binary ignores a `.sig` beside the zip and refuses a zeroed sums line ("checksum mismatch").
  - One detail is not a finding: a re-run of the sign step over an existing `SHA256SUMS.sig` prompts to overwrite, and it exits 0 on no answer. If the old signature is kept, the verify line after it catches the stale signature.

## Criteria
1. Dead code — fail (low): `mustRun` is orphaned in the test file. The source is clean: `ReleaseNamespace`, `releaseSigner`, `sshKeygen`, `SwapReleaseSigner`, `verifySignedSums` and the `bytes` import are all removed, and no references remain.
2. Flow tracing — pass. `UpdateFromRelease` → `downloadRelease` (pattern list reduced to two) → `verifyRelease` → `readSums` → `checksumFor`. `updateFromZip` → `UnpackZip(check)` → `VerifyZipSums` → `ErrZipNoSums`, which `main.go` maps to UNVERIFIED, while every other error refuses. A refusal leaves nothing behind (`TestVerifyZipSumsWritesNothingUnderTempDir`, and the refusal tests' binary/marker/glob assertions).
3. Standards conformance — pass, apart from the two stale restatements above. Errors are wrapped, and `fs.ErrNotExist` is checked with `errors.Is`.
4. Directory layout — pass. No new files in source.
5. Simplicity (DIY-except-security) — pass. Removing code leaves SHA-256 integrity checking with stdlib `crypto/sha256` as designed, and nothing is hand-rolled. That authenticity now rests on the GitHub credential and HTTPS is the operator's decision, recorded in ADR-053 and in R9 as amended. It is not reviewed here as a defect.
6. Decoupling — pass. `harness` exposes `VerifyZipSums` and `ErrZipNoSums` to `cmd`, and the test-only `SwapReleaseSigner` seam is gone.

## Not examined
- The full `go test -race -timeout 30m ./...` suite. Only the targeted subset above was run.
- The release rehearsal built only the host's linux/amd64 binary. The darwin/arm64 and linux/arm64 builds, step 4 (`gh release create/upload/view`), the tag commands and `gh auth status` were not run. They would publish, or they need the real tag.
- The bridge was rehearsed with a test key. The real retired key's signature, and the stock v2.3.2 binary's acceptance of it, can only be checked with that key.
- README.md's install chains beyond a glance at the `gh release download` and `sha256sum -c` lines.
- The records-only files in the range: DECISIONS.md, BUG-160, changes, handoffs, other reviews and testing logs.
