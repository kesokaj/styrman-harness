---
owner: main
date: 2026-10-05
---
# Releases without a key (ADR-053, ADR-053a)

## Type

feature

## Summary

The operator chose to drop release signing ("I just dont want to use 1 fucking ssh key for everything"), keeping versioning and `sh-agent update`. `4ee867a`: `internal/harness/release.go` loses the compiled-in signer, `ssh-keygen` and `SHA256SUMS.sig`; `update` and `--from-zip` check `SHA256SUMS` (version marker, checksum) only; `ErrZipUnsigned` is `ErrZipNoSums`; R9's Update, Zip place and Release bullets as approved; the release procedure tags with `git tag -a` from any machine where `gh auth status` passes; v2.4.0 carries one last `SHA256SUMS.sig` by the old key for v2.3.2 binaries. The rotation's new key (ADR-052) and its agent were deleted unused. Review fixes: the gh host is pinned (`github.com/kesokaj/styrman-harness`, `TestUpdateFromReleasePinsTheHost`), the zip line says `CHECKED … (integrity only)`, the bridge block derives the old key's fingerprint from its private half and refuses an empty path, and stale signing claims in README, AGENTS.md, delivery.md and the sh-agent matrix are corrected.

## Reviews

Code `reviews/2026-10-05-review-release-unsigned-adr053.md` (APPROVE WITH NITS, all three fixed), docs `reviews/2026-10-05-review-18b9018-4ee867a-docs.md` (nine findings, all fixed), security `reviews/2026-10-05-security-v2.4.0-unsigned-adr053.md` (host pin and wording fixed; Actions and the token left to the operator, ADR-053a).

## Not covered

Turning off GitHub Actions and rotating the `gh` token: the operator's decisions. The `verified vX.Y.Z` manifest token keeps its old word.
