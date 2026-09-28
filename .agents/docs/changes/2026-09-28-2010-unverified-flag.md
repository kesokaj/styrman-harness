---
owner: main
date: 2026-09-28
---
# An unsigned zip needs `--unverified`     (date: 2026-09-28)

## Summary
The operator approved the v1.9.0 handoff's Next Steps ("Go approved all").

**Refused without the flag.** `sh-agent update --from-zip <file>` now refuses a zip with neither `SHA256SUMS` nor `SHA256SUMS.sig` beside it, before anything is written, unless `--unverified` is passed. It then installs as before, marked `UNVERIFIED`.

**The flag allows a missing signature, never a failed one.** A zip with only one of the two files, a bad signature, a stale pair or a mismatch is refused with or without the flag.

**The flag belongs to `update --from-zip` only.** Anywhere else it is refused.

This closes the security medium carried since v1.8.0, and the R9 precision note.

## Type
feature

## Files changed
- **Code:** `cmd/sh-agent/main.go` (the flag, its refusal, `updateFromZip`'s new parameter, and the usage lines); the `ErrZipUnsigned` comment in `internal/harness/release.go`.
- **Tests:**
  - new: `TestUpdateFromZipUnsignedRefusedWithoutFlag` (the real run and the dry run) and `TestUnverifiedOnlyWithFromZip`;
  - extended: `TestUpdateFromZipForeignSignatureRefused` also runs with `--unverified`;
  - 13 unsigned-zip invocations now pass `--unverified`, so the refusal tests (zip-slip, not-a-zip, install refusal) still exercise their own refusal and not the new one. They were red without it.
- **Requirements and decisions:** `REQUIREMENTS.md` R9 (approved text) and `DECISIONS.md` ADR-040.
- **Docs:** `README.md` (§ Install from a zip, and the map label of the `main` zip edge) and `delivery.md` § The model.

## Verification
- `go test -race ./...` is green.
- **Must-fail mutations**, each compiled and failed a named test:
  - the refusal dropped failed `TestUpdateFromZipUnsignedRefusedWithoutFlag`;
  - the out-of-place flag check dropped failed `TestUnverifiedOnlyWithFromZip`;
  - the flag bypassing a failed check failed `TestUpdateFromZipForeignSignatureRefused`.

## Requirements
R9's zip-place clause is replaced with the text proposed in `handoffs/2026-09-28-1948-v1.9.0-released.md` Next Steps 1, approved "Go approved all" (ADR-040).

## Version
MAJOR (v2.0.0), per the version rule's "a new required flag".
