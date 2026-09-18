---
owner: security
date: 2026-09-18
---
# Security review — post-v1.2.2 worktree fix for findings 1 and 2   (2026-09-18)

Range reviewed: f3ec5dc..worktree (`git diff HEAD -- README.md .agents/docs/deploy/remote.md`; 13 insertions, 8 deletions, 2 files)

## Verdict: FINDINGS

**The `test` line does not close finding 1, and `remote.md:37-38` now says it does.** The new check binds the
downloaded binary to the `SHA256SUMS` *shipped beside it*; finding 1's attacker ships both, genuinely signed,
as a matched pair. An older release's binary + its own `SHA256SUMS` + its own `.sig`, served under a newer
tag, passes `ssh-keygen -Y verify`, passes `sha256sum -c`, and passes the new `test` — because the old binary
prints the old version and the old marker names the old version, so the two sides agree. `sh-agent update`
rejects exactly this (`internal/harness/release.go:281`, "the signed SHA256SUMS says %s but this is release
%s") because it compares the marker to **the tag it asked GitHub for**, an input the doc chain never brings
into the comparison. The comment asserting parity ("fails here as it does in `sh-agent update`") is therefore
a false guarantee in a security procedure, which is worse than round 1's silence. The `~/Downloads` fix is
sound: nothing outside the three named files can be touched, and nothing can be truncated or clobbered.

## Findings

- **[medium / confidence: high]** `.agents/docs/deploy/remote.md:37-38` — the comment claims the `test` line
  makes "an older, genuinely signed release served under a newer name fail here as it does in `sh-agent
  update`". It does not. `verifyRelease` compares the signed marker against the *tag* (`release.go:277-283`)
  and `update` additionally refuses a tag that is not newer than the running binary (`release.go:113-119`);
  the documented chain compares the marker against the *binary*, and the downgrade bundle is internally
  consistent, so both sides match and the chain proceeds to `mv`. Exploitable in principle by anyone who
  controls what the release endpoint serves without holding the signing key (a compromised GitHub account or
  the endpoint itself): the operator is handed a genuine older binary with known defects and every check says
  yes. The operator is worse off than before the change, because the comment tells them the case is covered.
  The only evidence on disk for the new line is a must-PASS run — `.agents/docs/testing/2026-09-17-release-binary-e2e.md`
  (worktree addition, § Added 2026-09-18) records `test … → exit 0` on the genuine, internally consistent
  v1.2.2 bundle, which is green whether or not the line has the property claimed for it; the must-FAIL state
  (an older genuine bundle under a newer tag) was never run, and the repo's own rule is that a check is cited
  only after both. Remediation (no dependency, one `&&` link, replaces nothing): make the version the reader
  *intends* the third input — `remote.md`, where the tag is already a variable: `&& grep -qx "# styrman-harness
  $TAG" SHA256SUMS` with `TAG=vX.Y.Z` also passed to `gh release download`; `README.md:90`, which names no tag
  at all (`gh release download` takes the newest release), needs the reader to type the version they mean.
  Keep the new `test` line as well — it catches a mis-stamped genuine build — but correct the comment to say
  what it checks: the FILE agrees with the checksum file beside it, not with the release you asked for. The
  must-fail rehearsal is two genuine bundles in a scratch directory, no forgery required.

- **[low / confidence: medium]** `README.md:96`, `.agents/docs/deploy/remote.md:35` — the check degenerates to
  a tautology when both sides are empty. `sed -n 's/^# styrman-harness //p'` on a `SHA256SUMS` with no marker
  prints nothing, making the right side the literal `"sh-agent "` (measured, trailing space); `sh-agent
  version` is `fmt.Fprintln(out, "sh-agent", buildVersion())` (`cmd/sh-agent/main.go:106`), which prints
  `"sh-agent "` if `buildVersion()` ever returns the empty string. On the Go versions in use it does not: a
  release build carries `-X main.version=vX.Y.Z`, and with no stamp `buildVersion` falls through to
  `versionFromInfo`, which returns the short VCS revision or `info.Main.Version` — `"(devel)"` for a source
  build, `"unknown"` when there is no build info at all (`main.go:894,916-935`). So the empty-left case needs a
  toolchain that leaves `Main.Version` empty, which I did not observe; confidence dropped to medium for that
  reason. The right-hand side going empty is real and measured, and it is the half that matters: a marker-less
  `SHA256SUMS` makes this check compare against a constant, and `sh-agent update` would reject that same file
  outright (`release.go:279`, "names no release version"). The `grep -qx` remediation above fails closed here;
  `test` does not.

- **[low / confidence: high]** `README.md:96`, `.agents/docs/deploy/remote.md:35` — the `sed` is stricter than
  the parser it stands in for, so the two disagree about what a valid marker is. Go's `sumsVersionMarker`
  (`release.go:48`) is `^#\s*styrman-harness\s+(v[0-9]+\.[0-9]+\.[0-9]+)\s*$` — any spacing, tabs, trailing
  whitespace; the `sed` needs exactly `"# styrman-harness "`. A future `SHA256SUMS` written with `#styrman-`
  or two spaces is accepted by `sh-agent update` and makes the documented chain fail with no message beyond a
  bare non-zero `test`. Fails closed, so this is operability, not exposure. Same class, opposite direction:
  the Go regex takes the FIRST marker (`FindSubmatch`), while `sed` prints ALL of them, so a two-marker file
  yields a two-line right side and the `test` fails (measured) where `update` would accept it. A CRLF
  `SHA256SUMS` yields `v1.2.2\r` and also fails closed (measured), though `sha256sum -c` would have failed on
  it first. Remediation: `grep -qx` on the exact intended line makes doc and code agree on one spelling.

- **[low / confidence: high]** `README.md:121` — the browser route's `mv` is partially atomic and leaks on
  failure. Measured with GNU `mv`: when one of the three sources is absent, the survivors are still moved and
  `mv` exits 1, so the chain stops before `cd` — correct as far as safety goes, but the user's downloads are
  now in a `mktemp -d` they were never shown, the trailing `rm -r "$D"` never runs, and on a box where `/tmp`
  is a tmpfs or is swept by `systemd-tmpfiles` those files are gone after a reboot. Nothing in `~/Downloads`
  can be truncated, clobbered or deleted: the three names are explicit, the destination is a fresh empty
  directory so no collision is possible, and an unrelated file beside them is untouched (verified in a probe).
  `D=$(mktemp -d)` failing leaves `D` empty, which makes the `mv` fail rather than doing anything destructive;
  `rm -r ""` at the tail is an error, not a wildcard. Remediation, if it is worth one line: moving the `mv`
  after `cd "$D"` is no better; naming the failure in the prose ("if it stops here, your three files are in
  `$D`") costs nothing and is what a reader needs at the moment it happens. Note also that a Linux reader of
  this Mac-shaped fence must now swap the asset name in three places (the `mv`, the `chmod`, the `test`) plus
  `shasum`→`sha256sum`; every miss fails closed, but the fence is one substitution harder than before.

- **[low / confidence: low]** `README.md:126` with `README.md:142-146` — on the platform this fence is written
  for, the binary is executed while still quarantined. A browser download on macOS carries
  `com.apple.quarantine`, and README documents `xattr -d` only *after* the move (line 145). The chain now
  executes `./sh-agent-darwin-arm64 version` inside `test "$(…)"`, so if Gatekeeper refuses the exec the chain
  stops with a bare `test` failure and the files stranded in `$D`. It fails closed and it is not a regression
  (the pre-diff line ran `version` too), but the diff moved that execution into a construct that hides the
  reason. Not measurable here — no Mac, and no round has run this fence on one; the worktree e2e log's own
  "Not exercised" list still names "The README browser chain on a Mac" (round 1 finding 2, still open).
  Remediation: run the fence once on the Mac and record what happens, before deciding whether an
  `xattr -d com.apple.quarantine "$ASSET"` link belongs in the chain.

## Checked

- **Secrets** — the only credential-shaped string in the diff is the `allowed_signers` ssh-ed25519 line and the
  operator's email, both unchanged byte-for-byte from f3ec5dc and both public by design (round 1 measured all
  seven occurrences identical; `TestReleaseSignerMatchesTheDeployGuide` pins the `remote.md` copies). No token,
  no key material, no `.env`. `remote.md:163` already records that publication is a separate decision.
- **Injection** — the diff adds shell that the operator pastes. `$ASSET` and `$D` are quoted throughout
  (`"./$ASSET"`, `"$D"`); `sed`'s script is single-quoted; the one unquoted-expansion risk, `rm -r "$D"` with
  an empty `D`, is quoted and therefore an error rather than a wildcard. The `SHA256SUMS` content reaching
  `sed` and `test` is signature-verified *before* it is parsed, and the binary is executed only after it
  matches a hash in that signed file — the ordering is right. No SQL, no template, no path-traversal surface.
- **Authz/authn** — unchanged: the trust root is the ssh-ed25519 signer compiled into the binary and written
  into the doc from a checkout, never from the release page. The gap is not authn but freshness (finding 1).
- **Crypto** — nothing hand-rolled and nothing composed here: `ssh-keygen -Y verify` and `sha256sum`/`shasum`
  as designed, in the `sh-agent-release` namespace the binary uses. No finding.
- **Dependencies** — no dependency added; the diff is documentation. My remediation adds none (`grep -qx` is
  in every POSIX toolchain these chains already invoke).
- **Tooling** — `govulncheck ./...` NOT run in this dispatch: no network by the dispatch's own hard limit, and
  the diff touches no Go code. Round 1 ran it on this tree over v1.2.1..f91257c on go1.27.1: "No vulnerabilities
  found". `npm audit` n/a — no JS in this repository.
- **GCP surface** — none: no IAM, no service account, no Cloud Run, no Secret Manager in scope.
- **Probes** — `.agents/scratch/security-post-v122/probe/` (gitignored, confirmed by `git check-ignore`): GNU
  `mv` partial-failure semantics, `sed` on a two-marker file, on a CRLF file, and on a marker-less file. No git
  command that changes the working tree was run; no install; no network; no `sh-agent report`.

## Not examined

- Both files in scope were read end to end (`README.md` 239 lines, `.agents/docs/deploy/remote.md` 211 lines).
- `internal/harness/release.go` read only at the lines the docs mirror (42-48, 95-125, 180-290); the download
  path's argv construction, `compareRelease`, and the self-replacement/reinstall steps were not read. The claim
  that `update` rejects the downgrade bundle rests on `release.go:277-283` plus the `tag == version` /
  `compareRelease` switch, both of which I read.
- The working tree also carries a modification to `.agents/docs/testing/2026-09-17-release-binary-e2e.md` and
  two untracked files (`reviews/2026-09-18-review-worktree-docs.md`, `testing/2026-09-18-verify-post-v1.2.2-docs.md`)
  that are outside my named range and are other agents' records. I read only the e2e log's diff, for what it
  claims about the new `test` line; the two untracked files I did not open.
- No test run: `go test ./...`, the installer against a scratch home, and `TestReleaseSignerMatchesTheDeployGuide`
  were not executed in this dispatch.
- macOS behaviour (Gatekeeper/quarantine on the browser route, `shasum -a 256 -c - --ignore-missing`) is
  unmeasured here and remains unmeasured anywhere — round 1 finding 2 is still open, and nothing in this diff
  addresses it.
- `gh release download`'s exit status when one `--pattern` matches nothing was not measured (no network); the
  chain appears to fail closed regardless, because a missing asset makes `sha256sum -c --ignore-missing` verify
  nothing and exit 1, and `chmod +x` on an absent file fails.

## Outcome

Written by the main agent after the operator chose option 1 of `.agents/scratch/human-route-downgrade-options.md` ("Go with 1"), in the commit that carries it. The finding was put to the operator rather than fixed on my own judgement because the previous fix on this mechanism was found not to fix it.

- [medium] the `test` line binds the binary to the bundle, not to anything outside it, and the comment claimed otherwise — **fixed**: every chain begins with the reader's own `V=` (README.md:91; `remote.md` .agents/docs/deploy/remote.md:28), the `gh` routes download `"$V"`, `&& grep -qx "# styrman-harness $V" SHA256SUMS` follows the signature check (README.md:97), and the binary test is against `$V`; the `remote.md` comment now says what the two lines do and what the signature alone does not (.agents/docs/deploy/remote.md:39). Verified in the must-fail direction by the verifier against the real bundle with `V=v1.2.3` and with a marker rewritten to `v1.2.1`. What remains is stated in the README comment: `V` is only as good as the reader's source for it.
- [low] tautology when the marker is absent — **fixed** by the same `grep`, which requires the marker to exist and equal `$V`; the `sed` extraction is gone, so the doc and the code no longer differ on marker spellings.
- [low] an absent download strands the other two in an unnamed temp dir — **fixed**: an `ls` of the three names precedes the `mv` (README.md:127), so the chain stops before anything moves.
- [low] the Mac chain runs the binary before `xattr -d` — **fixed**: the quarantine attribute is removed in the chain, after `chmod` and before `version` (README.md:133); `|| true` because `xattr -d` exits 1 when the attribute is already absent (a `gh` download), and a binary that is still quarantined fails closed at `version`.
- `~/Downloads` neighbours — confirmed safe by the finding; no change.
