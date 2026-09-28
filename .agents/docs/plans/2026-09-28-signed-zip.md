---
owner: main
date: 2026-09-28
---
# Plan — a signed release zip that `--from-zip` verifies   (2026-09-28)

## Why

The v1.8.0 gate's open medium: an unsigned zip becomes the agents' instructions, and `UNVERIFIED` only discloses that (`reviews/2026-09-28-preflight-v1.8.0.md` row 1). The handoff put the fix to the operator as Next Steps 1b: "/preflight could publish a `git archive` zip under the signed SHA256SUMS, and --from-zip would verify it … changes R9's Release procedure, so it needs your approval". The operator answered, 2026-09-28: "All approved. Cut new version if needed. All needs to be fixed."

## Scope

- **In:**
  - Each release publishes `styrman-harness-vX.Y.Z.zip`, a `git archive` of the tag with the top folder `styrman-harness-vX.Y.Z/`, listed in the signed `SHA256SUMS`.
  - `update --from-zip` verifies a zip when `SHA256SUMS` and `SHA256SUMS.sig` sit in the same directory as the zip. Otherwise it installs the zip marked `UNVERIFIED`, as today.
  - REQUIREMENTS R9 (the Release bullet and the zip place), ADR-039, the release procedure (`deploy/remote.md` § Releases), README, and `delivery.md`.
- **Out:**
  - a new flag;
  - verifying a `main` zip, which has no signature to check;
  - a manifest change.
- **Follow-ups, not built:** a zip handed over WITHOUT its sums still installs, marked UNVERIFIED, as R9 allows. This change narrows the v1.8.0 medium rather than closing it. Requiring an explicit `--unverified` would be a REQUIREMENTS R9 change and is put to the operator (ADR-039 Open; security reviews r1–r2).

## Surfaces touched

- `internal/harness/release.go`: `verifyRelease` (lines 219-294) splits into `verifySignedSums(keygen, sumsDir, signersDir) (sums []byte, tag string, err error)` and the per-asset check. `verifyRelease` keeps its behaviour and error texts; the existing `TestUpdateFromRelease*` tests pin them.
- `internal/harness/zip.go`: `UnpackZip` gains a check that runs after hashing and before anything is written.
- `cmd/sh-agent/main.go` (`updateFromZip`); tests in `internal/harness/zip_test.go` and `cmd/sh-agent/zip_test.go`.
- `REQUIREMENTS.md` R9; `.agents/docs/architecture/DECISIONS.md` ADR-039; `.agents/docs/deploy/remote.md` § Releases (the build, checksum and upload lines, and "all five" becoming six); `README.md` § Install from a zip; `.agents/docs/specs/sh-agent/delivery.md` § The model.

## Interfaces

- **`func UnpackZip(zipPath, parent string, check func(sum string) error) (dir, sum string, err error)`.** `check` may be nil. It runs once, right after the sha256 of the one read buffer is known, before `zip.NewReader`, and its error is returned unchanged, with nothing written.
- **`func VerifyZipSums(zipPath, sum string) (tag string, err error)`**, with `ErrZipUnsigned` beside it, in `internal/harness/release.go` next to `verifySignedSums` and `checksumFor`.
  - It looks for `SHA256SUMS` and `SHA256SUMS.sig` in `filepath.Dir(zipPath)`.
  - Neither present: it returns `"", ErrZipUnsigned`, a new sentinel.
  - Only one present: an error naming the missing one.
  - Both present:
    1. It writes the compiled-in signer to a temp directory, `os.MkdirTemp("", "sh-agent-verify-*")`, removed in a `defer`.
    2. `verifySignedSums` checks the signature and reads the `# styrman-harness vX.Y.Z` marker.
    3. `checksumFor(sums, filepath.Base(zipPath))` must equal `sum`.

    Any failure is an error, and the zip is refused.
  - It returns the marker's tag.
- **`updateFromZip`** passes a `check` closure that calls `VerifyZipSums` and handles `ErrZipUnsigned` ITSELF (D3):
  - `ErrZipUnsigned`: the closure records "unverified" and returns nil, so `UnpackZip` goes on to unpack, and the output says `UNVERIFIED` as today;
  - success: the closure records the tag and returns nil, and the first line reads `zip <file> (sha256 <hex>, VERIFIED against the signed SHA256SUMS of release vX.Y.Z)`, and the rest is unchanged;
  - any other error: the closure returns it, `UnpackZip` returns it unchanged with nothing written, and the command exits non-zero.
- **The release asset:** `git archive --format=zip --prefix=styrman-harness-vX.Y.Z/ -o <release dir>/styrman-harness-vX.Y.Z.zip vX.Y.Z`, added to the `sha256sum` line and the upload loop in `deploy/remote.md` § Releases.
- **REQUIREMENTS R9.**
  - The Release bullet: "…and publishes that tag's binaries with checksums signed by the same key." becomes "…and publishes that tag's binaries and a zip of its tree (`git archive`), with checksums signed by the same key."
  - The zip place sub-bullet gains: "When the release's `SHA256SUMS` and `SHA256SUMS.sig` sit beside the zip, the zip is verified against them and refused if either check fails; a zip without them is installed marked `UNVERIFIED`."
  - The approval note records "All approved", 2026-09-28, and ADR-039.

## Failure modes

| Situation | Result |
|---|---|
| The `.sig` is missing, or only one of the two files is there | Refused, naming the missing file |
| The signature is bad, or a different key signed it | Refused: "SHA256SUMS signature does not verify (…)" |
| The zip is renamed (a browser's `(1)` suffix), so `SHA256SUMS` has no line for it | Refused: "no line names <file>", with a hint to use the release's file name |
| The checksum differs | Refused: "checksum mismatch" |
| `ssh-keygen` is missing | Refused: "ssh-keygen not found; cannot verify" |

In every case the verification runs before any write, so the home is untouched.

## Verification

- **Unit tests** (`internal/harness`, with `newRelease`'s key and a zip built by `writeZip`):
  - verified;
  - `.sig` missing;
  - a foreign key;
  - checksum mismatch;
  - renamed zip;
  - no sums (`ErrZipUnsigned`);
  - `check`'s error stops the unpack before any write.
- **The existing `TestUpdateFromRelease*`** tests stay green unchanged.
- **CLI:** `cmd/sh-agent`, isolated with an empty `PATH` except `ssh-keygen`, linked into a temp bin the way the release tests do. A verified zip prints `VERIFIED` and installs; a tampered checksum refuses and leaves the home empty.
- **E2E criterion:** a release-versioned build verifies a `git archive` zip of `HEAD`, signed with a scratch key in the release layout, and prints `VERIFIED`. Flipping one byte of the zip refuses it with nothing written. The real v1.9.0 zip is verified the same way after the release, with the real key (operator smoke).

## Phases

- [x] Cold-reader check — r1 needed revision (fixed as D3); r2 was implementation ready with zero BLOCKING questions, and its non-blocking points are D4–D7 (`reviews/2026-09-28-coldread-signed-zip-r{1,2}.md`)
### Phase 1 — Verification and CLI
- [x] Refactor `verifyRelease` → `verifySignedSums` — verify: `go test -race -run 'TestUpdateFromRelease|TestVerifier' ./internal/harness`, AI — expect: pass unchanged
- [x] `UnpackZip` check hook and `VerifyZipSums` — verify: `go test -race -run 'TestVerifyZipSums|TestUnpackZipCheck' ./internal/harness`, AI — expect: the unit cases above, each seen red first
- [x] `updateFromZip` wiring and output — verify: `go test -race -run 'TestUpdateFromZipVerified|TestUpdateFromZipTampered' ./cmd/sh-agent`, AI — expect: `VERIFIED` printed; a tampered zip refused with an empty home
- [x] Checkpoint — expect: full suite green — evidence: `testing/2026-09-28-verify-signed-zip.md` PASS, 2026-09-28
### Phase 2 — Requirements, ADR, release procedure, docs, E2E, one commit
- [x] R9 text, ADR-039, `deploy/remote.md` § Releases, README, `delivery.md` — docs: those files — verify: `/review --docs` — expect: no stale "unsigned only" statement
- [x] E2E criterion — verify: the scratch-key run, AI
- [x] `reviewer` and `security` over the diff; `verifier` — expect: findings fixed or rebutted
- [x] Checkpoint — expect: one commit, pushed — evidence: the signed-zip commit on `main` (`git log --grep 'signed release zip'`), 2026-09-28

## Risks

Assumptions: none open. `git archive` keeps the tracked symlink `.agents/skills`, which sits outside the five trees, so `checkZipTrees` never opens it (`internal/harness/zip.go`); a symlink inside the five trees is refused. Both are covered by the zip plan's tests.

## Decisions

- D1 — 2026-09-28 — The check hook sits inside `UnpackZip`, so the sum is taken from the same buffer that is unpacked (security review r1, low 3) and verification precedes every write — by: implementer.
- D2 — 2026-09-28 — The sums are found beside the zip, with no flag. The operator downloads three files from the Releases page into one folder — by: implementer.
- D3 — 2026-09-28 — `UnpackZip`'s hook stays one rule (any error refuses, nothing written). The unsigned case is not an error at that level: `updateFromZip`'s closure turns `ErrZipUnsigned` into "unverified" and returns nil. One read, one buffer — by: cold read round 1 (`reviews/2026-09-28-coldread-signed-zip-r1.md`, finding 1).
- D4 — 2026-09-28 — The closure reports back through variables it captures in `updateFromZip` (`verifiedTag string`) — by: implementer, cold read r2.
- D5 — 2026-09-28 — `verifySignedSums` returns the marker's tag, and `verifyRelease` compares it with the tag it expects, keeping its existing error text — by: implementer, cold read r2.
- D6 — 2026-09-28 — No separate check of the file name against the tag: `checksumFor(sums, base(zip))` already requires that exact name in the signed file, so a renamed or foreign zip is refused — by: implementer, cold read r1–r2.
- D7 — 2026-09-28 — `SHA256SUMS` and its `.sig` found beside a zip that does not match them refuse it, rather than degrading to UNVERIFIED. Signed sums in the zip's folder say it is meant to be a release zip, so a mismatch is the case to stop at. The error names the release the sums belong to, and never suggests removing them (security review r1, finding 2: that hint was a one-step bypass) — by: implementer, cold read r2; amended after security r1.
- Rejected:
  - a `--sums` flag, one more thing to type wrong;
  - refusing every unsigned zip, which would break the approved `main` zip path (R9).
