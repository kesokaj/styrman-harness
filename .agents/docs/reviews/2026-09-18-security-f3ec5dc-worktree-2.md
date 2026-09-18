---
owner: security
date: 2026-09-18
---
# Security review — post-v1.2.2 worktree fix, round 2 (option 1)   (2026-09-18)

Range reviewed: f3ec5dc..worktree (round 2) — `git diff HEAD -- README.md .agents/docs/deploy/remote.md`; 42 insertions, 30 deletions, 2 files.

## Verdict: FINDINGS

**The mechanism is right and round 1's findings 1, 2 and 3 are closed; what is left is that the README tells the
reader to take `V` from the one source the adversary controls.** Making `V` the reader's out-of-band input is the
correct fix: `grep -qx "# styrman-harness $V"` binds the SIGNED checksum file to a version the download did not
choose, and `test "$(… version)" = "sh-agent $V"` binds the file's own claim to it too, so the round-1 downgrade
bundle (an older release's binary + its own `SHA256SUMS` + its own `.sig`, internally consistent, served under a
newer name) now stops at the `grep`. The residual is entirely in where `V` comes from, and `README.md:82-83`
names three sources of which the third — "the Releases page read on a machine you already trust" — is the
attacker-controlled endpoint itself: under exactly the adversary finding 1 named (someone who controls what the
release endpoint serves without holding the key), deleting or hiding the newest release makes the page say the
newest is the old version, the reader types that, and every check agrees. Trusting the *machine* is not trusting
the *source*. Everything else in the diff holds up under measurement: the `xattr -d … || true` cannot mask a
failure that matters, the `ls` closes the partial-`mv` leak for the ordinary case, no `~/Downloads` neighbour can
be touched, and the corrected `remote.md` comment no longer claims a parity the chain does not have.

## Answers to the five checkable questions

1. **Finding 1, human route: closed against the served-bundle substitution; open against the source of `V`.**
   `gh release download "$V"` now names a tag, and both new links compare the *signed* marker and the binary's
   own output to `$V`, which came from outside the download — so the whole-bundle swap fails at the `grep`
   (the old bundle's marker is the old version). Two things remain. (a) An attacker who controls the reader's
   source of `V` defeats it entirely, and README's own list offers such a source (finding A below). (b) A reader
   who copies `V=v1.2.2` literally from an old README is **pinning, not being downgraded**: the value came from
   the reader's side, the chain delivers exactly what was asked for, and the security property on offer — "the
   version you named, or nothing" — holds. It is a freshness/maintenance problem, not an integrity one
   (finding B): nothing bumps that literal and no test pins it, so after the next release every new install
   silently lands on v1.2.2 until `sh-agent update` runs. Freshness for a FIRST install cannot be established by
   any chain: `update`'s protection is `compareRelease(tag, version) < 0` against the *installed* binary
   (`internal/harness/release.go:113-119`), and a first install has no installed binary to compare against.
2. **Finding 2 (the tautology): closed.** The right-hand side is now a literal the reader supplied, not a value
   `sed`-extracted from the same file, so the two sides can no longer go empty together. Measured: with `V` empty
   the pattern is `# styrman-harness ` (trailing space), which matches no genuine marker — exit 1, fails closed
   (probe 7); the degenerate `grep -qx ""` case cannot arise because the pattern always carries the prefix
   (probe 8); and `gh release download ""` would fail first. Round 1's finding 3 (the `sed` disagreeing with
   `sumsVersionMarker`) is closed in the direction that matters: the `grep` pattern is now byte-identical to what
   the release procedure writes (`remote.md:145`, `printf '# styrman-harness %s\n'`), so doc and producer agree
   on one spelling. It stays stricter than the *consumer* regex (`\s*`, tabs, trailing space: probe 5 = exit 1)
   and, on this box, more permissive for a two-marker file (probe 4) — both fail closed or are unreachable
   (finding C).
3. **The `|| true` around `xattr -d`: cannot mask a failure that matters — it fails closed.** Every failure it
   swallows is re-tested by the very next link: if the attribute was not cleared and Gatekeeper refuses the exec,
   `$( … version)` is empty or garbage and `test` exits 1 (measured: empty vs expected → exit 1, probe 15), so the
   chain stops before `mv`. The failures it swallows are (i) attribute absent — harmless, (ii) `xattr` not
   installed, which is every Linux reader who copies the Mac fence (measured: absent here, probe 16) — harmless,
   (iii) removal genuinely refused — caught by the `version` link. The ordering is also right in the way that
   matters: the attribute is stripped only AFTER `ssh-keygen -Y verify` and the checksum, so no safety attribute
   is ever removed from an unverified file. One honest residual (finding D): `2>/dev/null` also hides *why*, so a
   Mac reader whose removal failed sees a bare `test` failure; and no round has run this fence on a Mac.
4. **The `ls` check: no path reaches a `~/Downloads` neighbour.** Measured in the probe — the three names are
   literal, there is no glob, the destination is a fresh empty `mktemp -d` so no collision or clobber is possible,
   and a neighbouring file beside them is untouched in every failure mode tried (probe 13). `ls` exits 2 when a
   member is missing, stopping the chain before `mktemp` and before `mv` (probe 11), which closes round 1's
   partial-`mv` leak for the ordinary case. Two residuals, neither reaching a neighbour: a TOCTOU window between
   `ls` and `mv` still strands the survivors in an unnamed `$D` (measured: `mv` exit 1, one file left in the temp
   dir, probe 12), so README's "stops the chain before anything moves" is true of the check, not of the whole
   chain; and `ls` passes on a *directory* named `SHA256SUMS`, which `mv` then relocates (probe 9-10) before
   `ssh-keygen … < SHA256SUMS` fails on it — fails closed, and `test -f` would be the precise check (finding E).
5. **The corrected `remote.md` comment: yes, with one omission.** Its three claims check out —
   the `grep` and `test` do tie the signed file and the binary to the version the reader named; `sh-agent update`
   does compare the signed marker to the tag it asked for (`release.go:277-282`); and "without them a genuinely
   signed OLDER release, served whole, passes the signature and checksum steps" is exactly right. Round 1's false
   guarantee ("fails here as it does in `sh-agent update`") is gone. What the sentence still leaves unsaid is the
   half of `update` the chain cannot have: `update` also refuses a tag that is not newer than the running binary
   (`release.go:113-119`), and it never runs the downloaded asset's `version` to compare it against the tag — it
   execs the REPLACEMENT with `install`, after verification and after the atomic replace (`release.go:155-177`),
   so the `test` link is an addition rather than a mirror, and the doc chain's ordering — run it before the move
   — is the stricter of the two. A reader can take "what `sh-agent update` does" as parity and conclude the chain protects freshness;
   it protects identity (finding B's remediation is one clause here).

## Findings

- **A. [medium / confidence: high]** `README.md:82-83` (and, more weakly, `.agents/docs/deploy/remote.md:24-25`) —
  the doc names the compromised channel as an acceptable source for `V`. "V is the version you expect, from a
  source outside the download: the operator, a checkout's changelog, or the Releases page read on a machine you
  already trust." The third is not outside the download — it is the same endpoint, rendered as HTML. Under
  finding 1's adversary (a compromised GitHub account, or control of what that endpoint serves, without the
  signing key) the page is made to show an older release as the newest; the reader takes `V` from it, and the
  genuine older bundle satisfies the signature, the checksum, the `grep` and the `test`. The operator is handed
  a known-defective binary with every check green — round 1's outcome reached through the display channel
  instead of the asset channel. `sh-agent update` is immune only because it compares against the version
  *already installed*, which a first install does not have. Remediation, no dependency, doc-only: strike the
  third option, leaving "the operator, or a checkout's changelog" (both genuinely out-of-band), and say in one
  clause what the chain guarantees — *the version you named, or nothing; not that it is the newest* — with the
  pointer that `sh-agent update` is what refuses to go backwards afterwards. `remote.md:24-25` says "you cut it,
  or read it on a machine you already trust"; "you cut it" is sound for that file's audience (the operator), so
  the same trim applies to the second clause only.

- **B. [low / confidence: high]** `README.md:91`, `README.md:126` — the version literal `V=v1.2.2` is hardcoded in
  both README chains, nothing bumps it and no test pins it (`git grep v1.2.2` finds these two lines and nothing
  else in tracked source, docs or skills; `skills/preflight/SKILL.md` has no step naming README). After the next
  release every reader who copies the block verbatim installs v1.2.2 — genuinely, verifiably, and stale — and no
  check anywhere in the chain or in `sh-agent status` says so. This is a pin, not a downgrade: the value came
  from the reader's side and the chain does exactly what it promises, which is why it is low and not medium. It
  is still the most likely way a real reader ends up on an old binary. `remote.md:28` avoids it with the
  `vX.Y.Z` placeholder, which fails closed (`gh release download "vX.Y.Z"` errors, and the `grep` could not
  match). Remediation, pick one: (a) use the same `vX.Y.Z` placeholder in README with "the newest version, from
  the Releases page or the operator" — fails closed, costs the reader one substitution; (b) keep the literal and
  add a line to the release procedure (`remote.md` § Releases, after `gh release view`) that bumps both README
  occurrences, plus a test asserting the README literal equals the newest tag. (a) is smaller and needs no code.

- **C. [low / confidence: high]** `README.md:97,131`, `.agents/docs/deploy/remote.md:34` — `grep -qx "# styrman-harness $V"`
  makes `$V` a BASIC REGULAR EXPRESSION, so each `.` in the version is a wildcard: measured, the pattern for
  `v1.2.2` also matches a line reading `# styrman-harness v1X2Y2` (probe 2), while `grep -Fqx` does not (probe 3).
  I could not construct a reachable exploit and say so plainly: `-x` fixes the length, and the release procedure
  plus `sumsVersionMarker` (`release.go:48`) only ever produce `v<digits>.<digits>.<digits>`, which for a
  same-length collision forces the dots back into the dot positions — so no *genuine signed* marker can collide,
  and only signature-verified bytes ever reach this line. Remediation is free and removes the class: `grep -Fqx`.
  Two related deviations, both measured and both harmless for the same ordering reason: this box's `grep` is
  **ugrep 7.8.4, not GNU grep**, and it strips a trailing CR, so a CRLF `SHA256SUMS` MATCHES here (exit 0, probe
  A/B) where round 1 measured the `sed` form failing closed under GNU grep — and its `grep -v '^#'` normalises
  the body, so `sha256sum -c` then passes on a CRLF file too (probe E). Unreachable: CRLF bytes would not carry
  the operator's signature. And `grep -qx` matches ANY marker line while the Go parser takes the FIRST (probe 4),
  a doc/code divergence that needs a two-marker file the operator signed. None of these is exposure; all of them
  are the reason to prefer the exact-match form.

- **D. [low / confidence: medium]** `README.md:133`, with `README.md:150` — the `xattr` step is written for a
  platform no round has run it on. It fails closed (question 3 above) and the prose at line 150 is now accurate
  about what the chain does, so this is not a blocker; but `2>/dev/null || true` means a Mac reader whose removal
  genuinely failed gets a bare non-zero `test` and no reason, and the worktree e2e log's own "Not exercised" list
  still names the README browser chain on a Mac. Confidence is medium because macOS Gatekeeper's exact behaviour
  for a direct command-line exec of a quarantined ad-hoc-signed arm64 binary is not measurable here. Remediation:
  run the fence once on the Mac and record the result before the next release; if it is worth a word, the prose
  can say "if the version check fails on a Mac, run `xattr -l` on the file in `$D`".

- **E. [low / confidence: high]** `README.md:124-127` — the prose claims more than the check does. "The `ls` stops
  the chain before anything moves when a file is missing" is true of a file missing *at the check*, and not of
  one that vanishes between the `ls` and the `mv`: measured, the survivors are still moved, `mv` exits 1, and the
  reader's files sit in a `mktemp -d` they were never shown (probe 12) — round 1's finding 4, narrowed, not
  closed. Nothing in `~/Downloads` outside the three named files is read, moved, clobbered or deleted in any
  failure mode tried (probe 13). Also `ls` is a presence check that passes on a directory (probe 9). Remediation,
  one line of prose rather than more shell: "if the chain stops before the move, your three files are in `$D`" —
  and `test -f … -a -f … -a -f …` (or three `test -f` links) instead of `ls` if the precision is wanted.

## Checked

- **Secrets** — the `allowed_signers` ssh-ed25519 line and the operator's email appear in the diff unchanged
  byte-for-byte from f3ec5dc (both public by design; `TestReleaseSignerMatchesTheDeployGuide` pins the
  `remote.md` copies). No token, key material, `.env` or credential path added. Nothing is echoed into a log.
- **Injection** — the diff is shell the operator pastes. `$V`, `$ASSET` and `$D` are quoted at every expansion
  (`"$V"`, `"./$ASSET"`, `"$D"`); the only place `$V` is unquoted-in-effect is as a `grep` pattern, which is
  finding C. `~/Downloads/...` is tilde expansion, which is not word-split, so a `$HOME` with a space is safe.
  `rm -r "$D"` with an empty `D` is an error, not a wildcard. Ordering is right throughout: the `SHA256SUMS`
  bytes are signature-verified BEFORE anything parses them, and the binary is executed only after its hash
  matches a line in that signed file. The one pipeline (`grep -v '^#' | sha256sum -c -`) has the filter UPSTREAM
  of the verifier, so the pipeline's status is the verifier's — the safe orientation of the harness's own rule.
  `sha256sum -c --ignore-missing` with no listed file present exits 1, not 0 (probe F): no vacuous pass.
- **Authz/authn** — unchanged and correct: the trust root is the ssh-ed25519 signer taken from a checkout or a
  trusted binary, never from the release page, and the doc still says so. The new links add *identity* (which
  release this is) to an already-sound *integrity* check. Freshness remains outside the chain's reach on a first
  install; that is finding A/B, not an authn defect.
- **Crypto** — nothing hand-rolled and nothing composed here: `ssh-keygen -Y verify` in the `sh-agent-release`
  namespace and `sha256sum`/`shasum -a 256`, both as designed, both matching what `internal/harness/release.go`
  does. No finding.
- **Dependencies** — none added; the diff is documentation. My remediations add none: `grep -F`, `test -f` and a
  prose trim are in every POSIX toolchain these chains already invoke.
- **Tooling** — `govulncheck ./...` NOT run in this dispatch: the dispatch forbids network and the diff touches
  no Go code. Round 1 ran it on this tree over v1.2.1..f91257c on go1.27.1 with "No vulnerabilities found".
  `npm audit` n/a — no JavaScript in this repository.
- **GCP surface** — none in scope: no IAM binding, service account, Cloud Run flag, Secret Manager reference or
  key file anywhere in the diff or its blast radius.
- **Probes** — `.agents/scratch/security-post-v122-2/probe/` (confirmed gitignored with `git check-ignore`):
  `grep -qx` BRE-dot wildcarding, `-F` behaviour, empty-`V`, loose/CRLF/two-marker markers, `ls` on a missing
  member and on a directory, the `ls`→`mv` TOCTOU, neighbour survival, `test` against empty output, and
  `sha256sum -c --ignore-missing` with nothing present. No git command that changes the working tree was run; no
  install, no network, no `sh-agent report`. Note for whoever repeats these: `grep` on this box is **ugrep
  7.8.4**, whose CR handling differs from GNU grep — the CRLF results are machine-specific and are called out as
  such in finding C.

## Not examined

- Both files in scope were read end to end this round (`README.md` 243 lines, `.agents/docs/deploy/remote.md`
  214 lines); no tool result in this review came back truncated.
- `internal/harness/release.go` was read only where the docs mirror it: the marker regex (line 48), the
  not-newer/not-a-release switch (100-130), the temp-dir/verify/replace/reinstall body (130-199) and the verify
  block (270-295). `checksumFor`, `compareRelease`, `writeAtomic` and `downloadRelease`'s error handling past
  line 199 were not read; the claim that `update` never runs the asset's `version` against the tag rests on
  those regions plus a grep for every `exec.Command` in the file, not on a full read of it.
- macOS is unmeasured, here and everywhere: the browser fence, `xattr -d com.apple.quarantine`, Gatekeeper's
  treatment of a quarantined ad-hoc-signed arm64 binary, and `shasum -a 256 -c - --ignore-missing`. Round 1
  finding 5 stays open; finding D is its successor.
- `gh release download "$V"` behaviour was not exercised (no network): its exit status when the tag does not
  exist, when a `--pattern` matches nothing, and whether `--clobber` interacts with `--dir` as assumed are all
  taken from the doc and from round 1's reasoning, not measured.
- No test run in this dispatch: `go test ./...`, `TestReleaseSignerMatchesTheDeployGuide`, `TestHarnessMap` and
  the installer against a scratch home were not executed.
- The working tree also carries a modification to `.agents/docs/testing/2026-09-17-release-binary-e2e.md` and
  four untracked records (`reviews/2026-09-18-review-worktree-docs.md`, the round-1 security report,
  `testing/2026-09-18-verify-post-v1.2.2-docs.md` and `-docs-2.md`). Only the round-1 security report was read;
  the others are outside my named range and are other agents' records. In particular I did not check whether any
  verification log records a must-FAIL run of the new `grep` link — round 1's complaint that only a must-PASS
  run existed is unverified for this round.

## Outcome

Written by the main agent, in the commit that carries the fix. The mechanism was confirmed right; the three residuals are all taken.

- A [medium] the README named the Releases page as a source for `V` — **fixed**: struck; the comment names the operator who cut it or a checkout's changelog, and says why not the page (README.md:83); the guarantee is stated as "the version you named, or nothing — never the newest" (README.md:87).
- B [low] `V=v1.2.2` hardcoded twice in README — **fixed**: both are `V=vX.Y.Z` (README.md:93), the same placeholder `remote.md` uses, so a stale copy-paste fails at the download instead of pinning silently; measured: the literal placeholder makes the `grep` exit 1.
- C [low] BRE dots in the `grep` pattern — **fixed**: `grep -qxF` in all three chains (README.md:99, .agents/docs/deploy/remote.md:34); measured against the real v1.2.2 bundle: `V=v1.2.2` exit 0; `V=v1.2.3`, `V=v1X2X2`, `V=vX.Y.Z` and empty `V` exit 1.
- (5) the `remote.md` comment's omission — **fixed**: it now says `update` also refuses a not-newer tag, which no first install can mirror (.agents/docs/deploy/remote.md:41).
- (3) `xattr … || true` fails closed — confirmed; no change. (4) the `ls`→`mv` window and `ls` passing on a directory of that name — **left** (main agent): the reader's own machine and three named files; nothing outside them is reachable, as the finding measured; README's "before anything moves" describes the check, which is what it says.
