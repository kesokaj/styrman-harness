---
owner: main
date: 2026-09-28
---
# Verification — signed release zip, E2E   (2026-09-28)

Range reviewed: 0b762db..worktree. Plan: `plans/2026-09-28-signed-zip.md`.

## Status: PASS

## Setup, proven before the runs

- **A scratch signer.** A scratch ed25519 key (`.agents/scratch/e2e-signed/keys/k`, comment `e2e@signer`). The binary was built with `-ldflags "-X main.version=v0.0.1 -X '…/internal/harness.releaseSigner=e2e@signer ssh-ed25519 …'"`, so it is a release-versioned build that trusts only the scratch key. The operator's release key signed nothing.
- **The release zip.** `git archive --format=zip --prefix=styrman-harness-v0.0.1/ -o rel/styrman-harness-v0.0.1.zip HEAD`, the same command the release procedure runs. Python `zipfile` showed a single top folder, `styrman-harness-v0.0.1`, and one symlink entry, `.agents/skills`, which lies outside the five trees and is never written.
- **The signed checksums.** `rel/SHA256SUMS` holds `# styrman-harness v0.0.1` and the zip's line. `ssh-keygen -Y sign -n sh-agent-release` produced `SHA256SUMS.sig`.

## Runs

| # | Run | Result |
|---|---|---|
| 1 | `update --from-zip rel/styrman-harness-v0.0.1.zip --home h1` | exit 0. The first line: `… VERIFIED against the signed SHA256SUMS of release v0.0.1`. `status`: `from zip sha256:be00544d3025`, `258 of 258 manifest files present`. |
| 2 | One byte of the zip flipped, then the same command with `--home h2` | exit 1. `checksum mismatch for styrman-harness-v0.0.1.zip against the signed SHA256SUMS of v0.0.1`, **0** entries in `h2`, and no `sh-agent-verify-*` directory left in `/tmp`. |

Unit tests (`TestVerifyZipSums`) cover these cases: verified; no sums, which is unsigned; `.sig` missing; a foreign key; a checksum mismatch; a renamed zip. `TestUnpackZipCheckStopsBeforeWriting` shows that a check that refuses leaves nothing behind. The CLI tests are `TestUpdateFromZipForeignSignatureRefused` (red before the wiring, "err = <nil>") and `TestUpdateFromZipUnsignedSaysUnverified`.

## Not exercised

- The real v1.9.0 release zip verified with the real key. That runs after the release, as an operator smoke check.
- macOS.
