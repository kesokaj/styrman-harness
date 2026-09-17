---
owner: main
date: 2026-09-17
---
# `update` from a signed release — Phase 4 of the release-binary plan     (2026-09-17)

## Summary

`sh-agent update` on a place with no harness checkout now moves that place forward from a GitHub release instead of refusing. `gh release view` gives the newest tag; the tag must parse as `vX.Y.Z` and be numerically newer than the running binary, or nothing is downloaded and the reason is printed. `gh release download` fetches this platform's `sh-agent-<goos>-<goarch>`, `SHA256SUMS` and `SHA256SUMS.sig` into a random `~/.config/sh-agent/release-*` directory that is removed on every path out. Verification is two checks and a third fact: `ssh-keygen -Y verify` checks the signature over `SHA256SUMS` against an allowed-signers line compiled into the binary, in the namespace `sh-agent-release`; the downloaded binary's SHA-256 must equal its line in that file; and the file must name the release it belongs to, `# styrman-harness vX.Y.Z`, matching the tag — so the version is part of what was signed (D26). Only then is `~/.local/bin/sh-agent` replaced through a temp file and a rename, and the NEW binary runs its own install, which deploys the harness built into it. `--check` stops after the tag comparison: `up to date (vX.Y.Z)`, or `newer release vA.B.C available` and exit 2. A development build is never overwritten from a release; it is told to download one by hand.

## Type

feature

## Files changed

- `internal/harness/release.go` (new) — `ReleaseRepo`, `ReleaseNamespace`, `releaseSigner`, `sumsVersionMarker`, `ErrNewerRelease`, `AssetName`, `LatestRelease`, `UpdateFromRelease`, `downloadRelease`, `verifyRelease`, `checksumFor`, `compareRelease`, `releaseParts`, `sshKeygen`, `sshKeygenSystemPaths`, `readCapped`, `maxAssetBytes`, `maxSumsBytes`
- `internal/harness/release_test.go` (new) — a fake `gh` and real `ssh-keygen` signatures from a test key, on a PATH holding nothing else: the good path, both up-to-date modes, `--check`, twelve refusals, the downgrade and bad-tag cases, `TestChecksumFor`, `TestCompareRelease`, `TestSSHKeygenPrefersASystemPath`, `TestReleaseSignerMatchesTheDeployGuide`
- `cmd/sh-agent/main.go` — `update` with an embedded source routes to `UpdateFromRelease`, mapping `ErrNewerRelease` to exit 2; `installPassArgs` extracted and shared with the checkout path; usage text and the command comment describe both kinds of update
- `cmd/sh-agent/main_test.go` — `TestReleaseBinaryUpdateCheckExitsTwo`; `TestReleaseBinaryWithoutCheckout`'s update leg rewritten for the release flow with PATH emptied
- `.agents/docs/plans/2026-09-17-release-binary.md` — Phase 4's two boxes ticked with evidence; D26, D27, D28; D22 marked superseded; Phase 5's release-procedure box now also requires the `SHA256SUMS` version marker and a rehearsal

## Agents dispatched

- `verifier` — PASS at 4d45155, `.agents/docs/testing/2026-09-17-release-binary-phase-4.md` (96 pass, 0 fail, 0 skip). That run predates the security fixes below; the counts at 40bf2db are 99 top-level tests and 85 subtests, none failing or skipped, with `gofmt`, `go vet`, `go test -race` and `govulncheck` clean, re-run by the main agent.
- `security` — 1 medium, 6 low, `.agents/docs/reviews/2026-09-17-security-dab2565-4d45155.md`; six fixed, two rebutted in its `## Outcome`. The medium was real and load-bearing: `tag == version` was an equality test and the signed bytes did not name the version, so a genuinely signed OLDER release would install, and one release's signed checksums verified under another release's tag. It found the signature chain otherwise sound — signer, namespace, identity and asset name compiled in or from `runtime`, nothing in a release able to redirect them, nothing outside the download directory written before both checks pass.

## Verification

Measured before any code: `ssh-keygen -Y verify` exits 0 on a good signature and 255 for a wrong key, a wrong namespace and altered content (OpenSSH 10.0p2). Twelve refusals are covered, each asserting the installed binary is byte-identical afterwards, no install ran, and no download directory is left: `gh` missing, `gh` unable to reach the repository, no verifier anywhere, no asset for this platform, a signature from another key, `SHA256SUMS` altered after signing, a checksum mismatch, a release shipping its own allowed-signers line under the compiled-in identity, signed checksums naming another version, signed checksums naming no version, no `.sig`, no `SHA256SUMS`. Must-fail by mutation (`.agents/scratch/p4-mustfail/mutate.py`, which asserts each target is unique, restores from a saved copy and reports a stale mutation as a failure, per BUG-068): all twelve mutations fail the tests that cover them — signature unchecked, checksum uncompared, verification skipped, no verifier required, verifier taken from PATH first, `gh` unchecked, namespace unpinned, signer taken from the release, downgrade allowed, tag unvalidated, the signed version marker unchecked, and a development build overwritten. Two mistakes of the session's own making, both caught by running rather than assuming: a test reached the REAL `gh` and the real repository because `/usr/bin` was on its PATH for `ssh-keygen` (PATH is now a single controlled directory, with `ssh-keygen` linked in), and a generated Go literal had its `\n` escapes turned into newlines, breaking the build (BUG-057's hazard; repaired and then verified by grepping for the escape itself). Not exercised: no real release exists yet, so every download is against the fake `gh`; the three-platform cross-compilation and the real `gh release create` are Phase 5; `darwin/arm64` behaviour is untested on this Linux box.

## Docs

`cmd/sh-agent/main.go`'s usage text and package comment describe both kinds of `update`. `.agents/docs/deploy/local.md` § Update still describes only the checkout kind, and `README.md`, the map, `AGENTS.md` rules 1 and 2, `REQUIREMENTS.md` R8/R9 and `.agents/docs/specs/sh-agent/delivery.md` still describe the unpacked bundle; all of them land in Phase 5's single commit to `main` (D11, R9). Phase 5 must also write the `SHA256SUMS` version marker into `/preflight` Step 6 and `remote.md` § Releases — until it does, no release this code will accept can be cut (D26), which is stated in that box.

## Requirements

R4 holds: stdlib only, `gh` and `ssh-keygen` shelled out, no dependency added. R8 and R9 gain the release-binary way in, in Phase 5's commit, with the operator's approved text. The signer line is the same operator key `.agents/docs/deploy/remote.md` § Releases gives `git verify-tag`, pinned by a test so the two copies cannot drift.
