---
owner: reviewer
date: 2026-09-28
---
# Cold-reader check — signed zip plan   (2026-09-28)

Range reviewed: `.agents/docs/plans/2026-09-28-signed-zip.md` (plan file, not a diff), against
`internal/harness/release.go`, `internal/harness/zip.go`, `cmd/sh-agent/main.go`,
`REQUIREMENTS.md` R9, `.agents/docs/deploy/remote.md` § Releases, `README.md`.

## Status: COMPLETE

## Verdict: design needs revision / implementation not ready

The plan has one real contradiction at its center: the `UnpackZip` interface says the `check`
hook's error is "returned unchanged, with nothing written" for *any* error, but `updateFromZip`'s
spec requires that a `check` error of `ErrZipUnsigned` specifically fall through to today's
successful UNVERIFIED install rather than refuse. As written these two clauses cannot both be
true, and the two implementers I can picture would resolve it two different ways (special-case
`ErrZipUnsigned` inside `UnpackZip` itself, vs. call `UnpackZip` a second time with `check=nil`
after catching `ErrZipUnsigned` from the first call) — with materially different consequences
(the second breaks the "read once, hash and unpack from the same buffer" invariant the existing
doc comment on `UnpackZip` states as a security property). That single gap should be closed
before Phase 1 starts; the rest of the plan is otherwise tight and its docs/requirements/ADR
surface list checks out against the current `README.md`/`REQUIREMENTS.md`/`deploy/remote.md`
text.

## (1) Gaps

- **BLOCKING — the `ErrZipUnsigned` fallback contradicts the `check`-hook contract.**
  `UnpackZip` (plan line 32): "[`check`'s] error is returned unchanged, with nothing written."
  `updateFromZip` (plan line 45): "`ErrZipUnsigned` → `UNVERIFIED`, as today" — i.e. the zip must
  still be unpacked and installed. Today's unverified path (`cmd/sh-agent/main.go:1079-1096`)
  unconditionally calls `UnpackZip` then `installFiles` and writes files. If `check` returning
  `ErrZipUnsigned` makes `UnpackZip` abort "with nothing written" (as the interface says for
  every error), there is no path left by which the zip still gets installed UNVERIFIED — unless
  `UnpackZip` special-cases `ErrZipUnsigned` internally (not described anywhere: which lines, does
  `zip.NewReader` still run, does the rest of the unpack proceed as if `check` were nil?), or
  `updateFromZip` calls `UnpackZip` a second time with `check=nil` (not described either, and it
  would read+hash the file twice, breaking the "read once… sum and unpack both come from that one
  buffer" property `zip.go`'s doc comment currently guarantees — `internal/harness/zip.go:44-45`).
  The Verification section's own unit-test bullet "`check`'s error stops the unpack before any
  write" (plan line 75) reads as universal, which sharpens the contradiction rather than resolving
  it: is that bullet about a *non*-`ErrZipUnsigned` error only? The plan does not say.
- **Gap — where does `VerifyZipSums` (and `ErrZipUnsigned`) live?** Surfaces touched names
  `release.go` for the `verifyRelease` split and `zip.go` for the `UnpackZip` check hook, but
  `VerifyZipSums` composes `verifySignedSums` (release.go) with a zip-specific checksum lookup and
  is called as `check`'s callee from `updateFromZip`. Non-blocking (release.go is the reasonable
  default, since `verifySignedSums`/`checksumFor` already live there and are unexported), but the
  plan should say so since it changes which file the new exported symbols and the sentinel error
  land in.
- **Gap — the symlink risk note names a check that no phase task performs.** Risks (plan line 97)
  asserts `git archive` keeps `.agents/skills` as a symlink entry and that this is "Checked by
  `unzip -l` on the E2E zip," but neither the E2E criterion (line 78) nor any Phase 2 task
  mentions running `unzip -l` or asserting on its output. In the current tree this risk is
  actually moot — `.agents/skills` sits under `.agents/`, which is not one of `zipTrees`
  (`context`, `skills`, `agents`, `commands`, `rules`; `internal/harness/zip.go:38`), so
  `checkZipTrees` (`internal/harness/zip.go:127-129`, the `continue` on an unmatched top segment)
  never inspects it, and even if it did, `checkZipTrees` unconditionally refuses any symlink entry
  found inside the five trees (`internal/harness/zip.go:141-142`) — so the "risk" cannot actually
  manifest as a silent inclusion; worst case is a hard refusal, which is arguably the point. The
  gap is that the plan states an assumption "checked by X" that X never executes in the verify
  clauses as written — low severity, but the phrasing overclaims verification that isn't wired up.
- **Not a checked assumption: no cross-check between the SHA256SUMS marker's tag and the zip's own
  file-name version.** `VerifyZipSums` (unlike `verifyRelease`, which takes a caller-known `tag`
  and refuses a mismatch, `internal/harness/release.go:282-284`) has no expected-tag input — it
  trusts whatever `# styrman-harness vX.Y.Z` marker is inside the *matching* signed SHA256SUMS.
  This is probably intentional (there is no "expected version" for a zip the operator points at by
  hand), and the per-file checksum still binds the specific `zip` bytes to the specific listed
  file name, so a straight substitution attack is caught. But the plan never states this as a
  decision, so a reviewer can't tell "no cross-check" from "forgot to specify one." Non-blocking —
  a reasonable implementer reads it the same way I did — but worth a one-line D-row.

## (2) Scope

Nothing built beyond `## Why`/`## Scope › In`. The zip-verification behavior, the R9/ADR/docs
surface list, and the Out bullets (no new flag, no `main`-zip verification, no manifest change)
all match what the operator approved in Next Steps 1b. No extra endpoint, no second shape for the
same operation, no configuration for a value that doesn't vary. Confirmed the current `README.md`
(`README.md:168`, the "zip is unsigned… UNVERIFIED" sentence) and `REQUIREMENTS.md` R9's zip-place
sub-bullet (`REQUIREMENTS.md:44`) are exactly the strings the plan's before/after quotes name, so
the docs-touch list is neither over- nor under-scoped relative to what currently needs to change.

## (3) Readiness

Questions a first-pass implementer would have to ask:

1. **BLOCKING** — How does `ErrZipUnsigned` from `check` actually result in a written, UNVERIFIED
   install, given `UnpackZip`'s stated "any error → nothing written" contract? (Gap 1 above.) Two
   careful implementers would build different things here (special-case inside `UnpackZip` vs. a
   second call from `updateFromZip`), with different correctness properties.
2. **NON-BLOCKING** — Which file does `VerifyZipSums`/`ErrZipUnsigned` live in? A reasonable
   implementer puts it in `release.go` beside `verifySignedSums`/`checksumFor`, which it directly
   reuses.
3. **NON-BLOCKING** — Is the `unzip -l` symlink check meant to be an added E2E assertion, or is
   the Risks note just documenting why the existing symlink refusal in `checkZipTrees` already
   covers it? Either answer is defensible and doesn't change the code, only whether one assertion
   line gets added to the E2E task.
4. **NON-BLOCKING** — No cross-check between the SHA256SUMS marker tag and the zip's own name is
   intentional read from the interfaces text; a competent implementer reaches the same reading.

## Not examined

- `.agents/docs/architecture/DECISIONS.md` ADR entries beyond confirming ADR-038 is the latest
  (039 is free) — did not read ADR-038's full text for consistency with D1/D2's "by: implementer"
  framing.
- `.agents/docs/specs/sh-agent/delivery.md` § The model — named in Surfaces touched but not read
  against the plan's wording; I did not check whether its current text also needs the same
  "unsigned only" restatement fix beyond what grep found (it appeared in the grep for
  UNVERIFIED/unsigned across the repo, consistent with being on the touch list, but I did not read
  its prose).
- `cmd/sh-agent/zip_test.go` and `internal/harness/zip_test.go`/`release_test.go` current test
  bodies beyond the `newRelease`/`writeZip` helper signatures — did not verify the named future
  test names (`TestVerifyZipSums`, `TestUnpackZipCheck`, `TestUpdateFromZipVerified`,
  `TestUpdateFromZipTampered`) don't collide with existing tests.
- The full `ssh-keygen`-missing behavior change in context (a machine with SHA256SUMS present but
  no `ssh-keygen` now hard-refuses instead of falling back to UNVERIFIED as it would today with no
  SHA256SUMS present) — this is explicit in the plan's failure-modes table so it reads as a
  decided design choice, not an unstated gap, but I did not check whether any other doc surface
  promises graceful degradation when `ssh-keygen` is absent.

design needs revision / implementation not ready
