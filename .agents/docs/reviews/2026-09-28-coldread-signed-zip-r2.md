---
owner: reviewer
date: 2026-09-28
---
# Cold-reader check — signed zip plan, round 2   (2026-09-28)

Range reviewed: `.agents/docs/plans/2026-09-28-signed-zip.md` (plan file, not a diff), against
`internal/harness/zip.go`, `internal/harness/release.go` (lines 190-320), `cmd/sh-agent/main.go`
(lines 1040-1137), `internal/harness/zip_test.go`/`release_test.go`/`cmd/sh-agent/zip_test.go`
(test names only), `REQUIREMENTS.md` R8/R9, `.agents/docs/deploy/remote.md` § Releases,
`.agents/docs/specs/sh-agent/delivery.md` line 32, `README.md:168`, `.agents/docs/architecture/DECISIONS.md`
(ADR list + ADR-038 body), and round 1's report.

## Status: COMPLETE

## Verdict: design needs revision (non-blocking gaps remain) / implementation ready (zero blocking questions)

Round 1's one BLOCKING finding — the `ErrZipUnsigned`-vs-"nothing written" contradiction — is
resolved by the plan's new D3: `UnpackZip`'s `check` hook still obeys one rule (any *error*
refuses, nothing written), and `updateFromZip`'s closure never lets `ErrZipUnsigned` reach
`UnpackZip` as an error — it catches it, records "unverified", and returns `nil`, so `UnpackZip`
proceeds exactly as it would with `check == nil`. I traced this against the current `UnpackZip`
body (`internal/harness/zip.go:51-114`): the hash is computed before `zip.NewReader` and before
any directory is created (`os.MkdirTemp` at line 82), so a `check` inserted "right after the
sha256 … before `zip.NewReader`" is trivially "nothing written" for every error path, including
today's unverified one. I found no remaining contradiction. The docs surface list is also
complete: I grepped `skills/`, `agents/`, `commands/`, `context/`, `REQUIREMENTS.md`, and
`deploy/remote.md` for other restatements of "UNVERIFIED"/"unsigned" and found only the ones the
plan already names (`README.md:168`, `delivery.md:32`, R9's zip-place sub-bullet) plus ADR-038's
own historical text (`DECISIONS.md:332,350`), which is append-only and correctly left alone — a
new ADR-039 supersedes it, the same pattern as ADR-033 superseding ADR-030. What's left is a
handful of non-blocking mechanics the plan leaves to the implementer's judgment; none of them
would produce two materially different builds.

## (1) Gaps

- **Gap (non-blocking) — round 1's blocking contradiction: resolved.** Confirmed against the
  current code (see above). No action needed.
- **Gap (non-blocking) — how the `check` closure reports its verdict back to `updateFromZip` is
  unspecified.** The Interfaces section (plan lines 44-47) says the closure "records the tag and
  returns nil" / "records 'unverified' and returns nil," but `check`'s signature is
  `func(sum string) error` — it has no other return channel. The implementer must invent a
  side-channel (a pointer to a struct the closure fills in, most naturally) so `updateFromZip` can
  build the `VERIFIED`/`UNVERIFIED` first line after `UnpackZip` returns. This is ordinary
  Go and a competent implementer reaches the same shape, but the plan states the effect ("records")
  without naming the mechanism, and a second implementer could as easily reach for a package-level
  variable (wrong, not concurrency-safe if ever reused) absent any hint. Worth one line, not a
  blocker.
- **Gap (non-blocking) — where the existing tag-equality check moves.** `verifyRelease`
  today (`internal/harness/release.go:278-284`) reads the `SHA256SUMS` marker and refuses when it
  doesn't match the caller's own `tag` argument. The plan's `verifySignedSums(keygen, sumsDir,
  signersDir) (sums []byte, tag string, err error)` (plan line 25) has no expected-tag input, so
  that comparison can no longer live inside it — it must move to `verifyRelease` itself, done
  *after* calling `verifySignedSums`, comparing its returned `tag` against the caller-known one.
  The plan doesn't say this in as many words, but the given signature leaves no other place for it,
  so it's decidable from the interfaces text alone, not a real fork.
- **Gap (non-blocking) — carried from round 1, still open: no cross-check between the marker's
  tag and the zip's own file name/version.** `VerifyZipSums` trusts whatever tag the matching
  `SHA256SUMS` claims; the per-file checksum still binds the exact zip bytes to the exact listed
  file name, so a substitution is still caught, just not reported with an "expected vX.Y.Z, got
  vX.Y.Y" message. Still no D-row recording this as deliberate — plan text and round 1's read of
  it agree, so a second implementer reaches the same conclusion, but it remains a should-record,
  not a should-block.
- **Gap (non-blocking) — a folder that holds a stale or foreign `SHA256SUMS`/`SHA256SUMS.sig` next
  to a differently-sourced zip now hard-refuses instead of falling back to today's always-succeeds
  `UNVERIFIED` install.** Concretely: an operator who previously downloaded a signed release into a
  folder, then drops a fresh `main`-branch zip (or a renamed/re-downloaded release zip) into the
  *same* folder, gets "no line names `<file>`" or "checksum mismatch" rather than an UNVERIFIED
  install — a real UX regression from today's behaviour, where an unsigned zip always installs.
  The Failure Modes table already documents this literal outcome (rows 3 and 4), so it is not an
  unstated code path, but the plan never calls out that it is a considered trade-off versus a
  side-effect of the "sums beside the zip, no flag" design (D2) — a one-line D-row ("stale sums
  beside a zip refuse rather than degrade to unverified, because a present-but-wrong signature
  reads as tampering, not absence") would remove the ambiguity about whether this was seen and
  accepted.

## (2) Scope

Unchanged from round 1: nothing built beyond `## Why`/`## Scope › In`. Re-confirmed this round —
`REQUIREMENTS.md` R8/R9's current text, `README.md:168`, `delivery.md:32`, and `deploy/remote.md`
§ Releases' asset list (`SHA256SUMS SHA256SUMS.sig sh-agent-linux-amd64 sh-agent-darwin-arm64
sh-agent-linux-arm64`, "all five") match exactly what the plan's before/after quotes describe
changing, and ADR-039 is the next free ADR number (ADR-038 is the newest in `DECISIONS.md`). No
extra endpoint, no second shape for the same operation (the plan explicitly rejects a `--sums`
flag), no configuration for a value that doesn't vary. The `verifySignedSums(keygen, sumsDir,
signersDir)` split takes two directory parameters instead of one existing `dir`, but this is
because `VerifyZipSums` genuinely needs a *different* place to write the temp signer file (a
scratch tempdir) than where `SHA256SUMS` lives (the operator's own download folder) — a real
second caller with a real second need, not speculative parameterization.

## (3) Readiness

Questions a first-pass implementer would have to ask — **zero BLOCKING**:

1. **NON-BLOCKING** — What Go shape carries the closure's verdict ("unverified" / the tag) back to
   `updateFromZip`? A pointer to a small struct captured by the closure is the natural answer; the
   plan doesn't name it but nothing else is defensible for a single-goroutine call.
2. **NON-BLOCKING** — Does the tag-equality check that `verifyRelease` performs today move into
   `verifyRelease` itself (after calling `verifySignedSums`), since the new split signature has no
   room for it elsewhere? Yes, by elimination from the given signature.
3. **NON-BLOCKING** — Is "no cross-check between the SHA256SUMS marker's tag and the zip's file
   name" intentional? Reads as yes from the interfaces text; carried from round 1, still
   non-blocking.
4. **NON-BLOCKING** — Should a stale/foreign `SHA256SUMS` beside a fresh zip hard-refuse (as the
   Failure Modes table already says) rather than degrade to `UNVERIFIED`? The table already
   answers this; only the "was this seen and chosen, or just fallen out of the mechanics" framing
   is undocumented.

## Not examined

- The exact Go implementation once written (this is still a plan, not a diff) — nothing here was
  checked against actual new code, since none exists yet.
- `cmd/sh-agent/zip_test.go` and `internal/harness/zip_test.go`/`release_test.go` full bodies
  (only `func Test…` names were listed, to check for naming collisions with the plan's four new
  test names — none found).
- `.agents/docs/architecture/DECISIONS.md` ADR-038's full body beyond the two quoted lines (332,
  350) and the "Consequences" section of the earlier ADR-030 entry — not needed for this check but
  not read line-by-line.
- Whether `ssh-keygen`-missing behaviour (hard refuse when `SHA256SUMS` is present but the tool
  isn't) is promised anywhere else as graceful — round 1 flagged this as open and I did not
  re-check it this round; it is still explicit in the plan's Failure Modes table, so it reads as a
  decided choice, not a silent gap.

design needs revision (non-blocking gaps only) / implementation ready
