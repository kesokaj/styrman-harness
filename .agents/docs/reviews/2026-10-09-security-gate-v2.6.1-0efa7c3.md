---
owner: security
date: 2026-10-09
---
# Security review — gate-v2.6.1-0efa7c3   (2026-10-09)

Range reviewed: v2.6.1..0efa7c3

Commits: 9964037 (v2.6.1 handoff), 11506d1 (BUG-166 fix), 0efa7c3 (BUG-166 handoff). 14 paths, +479/-10.

## Status: COMPLETE

## Verdict: FINDINGS (low only; nothing blocks the gate on security grounds)

No secret leaks, nothing is injectable, and nothing is hand-rolled. The worst issue is a trade-off the trail already records and accepts. On macOS the dream job's headless agy gets a symlink to the operator's real `~/Library/Keychains` (`internal/dream/run.go:301-312`), not a copy. Any Keychain write the child makes therefore reaches the interactive login: a token refresh, a rewrite, or a delete after an auth failure. On Linux it only ever touches a copy. Confidentiality is unchanged. The child has the same uid, so it could already open the real path, and the tool-refusing `settings.json` is still written on both OSes and asserted by the test. Both removal paths, `cleanup` and `removeCopies`, use `os.RemoveAll`, which unlinks a symlink and does not follow it. Tests pin this for both paths, and I re-ran them green. Every finding from the in-range BUG-166 security review (`reviews/2026-10-09-security-bug166-agy-keychain.md`) is fixed or documented as its `## Outcome` says, and I checked each against the code. The macOS behaviour itself has never run on a Mac, except for the operator's PONG probe. That gap belongs to the gate (guide § 3e), not to this review, and is listed under Not examined.

## Findings
- **[low / confidence: medium]** `internal/dream/run.go:301-312`: carried forward from the in-range review as accepted residual risk, not new. The Keychains link removes write isolation of the login on macOS. A dream-run agy that refreshes, rotates or deletes its Keychain item changes the operator's interactive login, and can race an interactive agy. This affects integrity and availability of the login, not disclosure. It is documented in `specs/sh-agent/dream.md:34` and `deploy/local.md:93`. The handoff's Next Step 1 makes "the interactive agy stays logged in afterwards" a pass condition of the owed macOS run. Remediation: none beyond that run. A narrower link would not restore isolation.
- **[low / confidence: low]** `internal/dream/run.go:127-151` with `agyHome`: when the Keychain is locked or the login has expired, agy falls back to an interactive OAuth flow (BUG-166 Repro 2: "printed a browser-login URL and waited"). Nobody can complete it headless. The dream loop then waits the full `cliTimeout` (10 min) for each project, and each wait holds an OAuth prompt open (possibly a loopback callback listener; not verified). The authorization URL goes to `dream-failed/` through `o.clean` with mode 0600. It holds a client id, state and PKCE challenge, which are not credentials. Exploitation would need a same-host actor to race a PKCE-bound callback without seeing the URL, so this is mostly an availability cost: N projects × 10 minutes, nightly. Remediation: after the first project whose failure carries `agyAuthHint`, stop the run (exit 2, log line) instead of retrying every project. The operator deferred the locked-Keychain case on 2026-10-09, so this is recorded and does not block.

## Checked
- **Secrets.** No credential in the code, tests or trail. The test keychain and token are fixture strings under `t.TempDir()`, with `t.Setenv("HOME")`. `AKIAIOSFODNN7EXAMPLE` in the v2.6.1 handoff (`handoffs/2026-10-09-1419-v2.6.1-released.md`) is AWS's published example key. The new log text is constant (`agyAuthHint`) or names paths only (`agyHome` errors, `HOME is %q`). `report` still sends only `proposals/dream/harness/*.md` from the dream tree, listed by name, with a symlink on the way refused (`report.go:404-418`). It never sends `dream.log`, `dream-failed/` or `dream-agy-home-*`. The global search index reads `proposals/` and named files only (`search/reconcile.go:353`), so it cannot reach a leftover link. Private data in the trail: `/Users/spgo` (the operator's Mac username) in BUG-166:18. It has precedent in 9 files at v2.6.1, the repository is PRIVATE (`gh repo view`), and it is no finding.
- **Instruction-shaped text in the trail.** None aimed at a reviewer. The handoff's "Next Steps: cut v2.6.2 … Recommendation: yes" is addressed to the operator as a decision and authorizes nothing. The prior review's Outcome is claims, which I checked against the code (above).
- **Injection.** No new external input reaches a command. The argv is unchanged and the prompt goes on stdin. `agyAuthHint` reads only lines that are not `{` and appends a fixed string, so model text cannot inject into `dream.log`; at worst it gets a misleading hint, which is pinned by `TestAgyAuthHint`.
- **Authz/authn.** The private HOME still gets `settings.json` with `allow:[]`, `request-review` and no trusted workspaces, written after the link on both OSes and asserted in `TestAgyHomeMacKeychain`. The relative-HOME refusal is in place (`run.go:267`). A token file that is a symlink or not a regular file is refused (`copyPrivate`, Lstat). Any Stat error other than ENOENT fails closed.
- **Filesystem walkers that could follow the link.** `cleanup` and `removeCopies` (`RemoveAll`, `e.IsDir()` on the entry's own type) are pinned by `TestAgyHomeMacKeychain` and `TestRunRemovesLeftoverCopies`, which I re-ran: `go test -count=1 ./internal/dream -run …` ok. The kit preamble's credential scan (`testing/kit-preamble.sh:171`, `find -H … -type l` then `grep` on the path) would hit a directory link and report it as unscannable, which fails closed.
- **Crypto.** None touched, nothing hand-rolled. `runtime` is stdlib.
- **Dependencies.** `go.mod` and `go.sum` are unchanged in the range. `govulncheck ./...` on go1.27.2: "Your code is affected by 0 vulnerabilities", with 1 in required modules that is not called (the known GO-2026-5024). `npm audit` does not apply: no JS ecosystem.
- **GCP surface.** Not touched.

## Not examined
- Every path in the release diff was read end to end in the diff: `internal/dream/run.go`, `internal/dream/run_test.go`, `README.md`, `.agents/docs/{bugs/BUG-166-…,bugs/INDEX.md,changes/2026-10-09-1503-…,deploy/local.md,handoffs/2026-10-09-1419-…,handoffs/2026-10-09-1504-…,reviews/2026-10-09-review-bug166-…,reviews/2026-10-09-security-bug166-…,specs/sh-agent/dream.md,testing/2026-10-09-verify-bug166-…,testing/guide.md}`. No release-diff path is left unread.
- Real macOS behaviour was not examined: no Mac here. Unverified are whether agy writes to or deletes its Keychain item through the link, whether tool refusal holds for the real dreamer prompt (Repro 3 used a PONG prompt), and the locked-Keychain and launchd path. The release gate owes these as one real `sh-agent dream --all` on a macOS place (`testing/guide.md` § 3e), which this review does not replace.
- Not mutation-checked by me: the new tests (I may not edit source). The verifier log also says it did not mutation-check. The in-range reviews state mutation results for `TestAgyAuthHint` and the relative-HOME case.
- `internal/dream/collect.go` and the `fakeCLI` and `agyLogin` test helpers were read only where the diff uses them. They are unchanged in the range.

## Outcome

- 1 (dream's agy shares the real Keychain login on macOS): left, accepted and documented (`specs/sh-agent/dream.md`, `deploy/local.md`). The owed macOS run checks that the interactive agy still answers afterwards.
- 2 (a locked Keychain waits out the timeout per project): left. agy's own 60-second login wait (operator probe 2026-10-09) bounds it at about a minute per project, not the 10-minute `cliTimeout`. Stopping the whole run at the first login failure is deferred with the operator's timer item (handoff Next Steps).
