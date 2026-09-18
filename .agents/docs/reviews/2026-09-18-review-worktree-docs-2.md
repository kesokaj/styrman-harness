---
owner: reviewer
date: 2026-09-18
---
# Review — worktree-docs-2   (2026-09-18)

Range reviewed: working tree — `git diff HEAD -- README.md .agents/docs/deploy/remote.md .agents/docs/testing/2026-09-17-release-binary-e2e.md`

## Verdict: APPROVE WITH NITS

The version-binding hardening (V pinned outside the download, `grep -qx "# styrman-harness $V"`, `test "$(... version)" = "sh-agent $V"`) is applied consistently to both README's `gh` route, README's browser route, and remote.md's release-binary route; the xattr fold reads correctly and no stray standalone xattr paragraph survives elsewhere; and the "newest release" phrasing throughout the repo is correctly scoped to `sh-agent update` only — the manual download instructions never claimed that and still don't. One real defect: the new `grep -qx` version check in both README.md and remote.md uses the version string as an unescaped basic regular expression, so a `.` in `vX.Y.Z` matches any character rather than a literal dot — a narrow but concrete weakening of the exact check this diff was written to add.

## Findings

- **[medium / confidence: medium]** `README.md:127` and `.agents/docs/deploy/remote.md:34` — `grep -qx "# styrman-harness $V" SHA256SUMS` uses `$V` (e.g. `v1.2.2`) as a `grep` basic regular expression, not a literal string: the unescaped `.` characters match any character, so the exact-version binding these two lines were added to enforce (per the diff's own stated rationale — defeating "a genuinely signed OLDER release, served whole") is looser than a literal-string match would be. Use `grep -qxF` (fixed-string) or escape the dots before matching.
- **[low / confidence: low]** `.agents/docs/testing/guide.md` — no tier cites the manual release-binary verify chain (README § Install without cloning, remote.md § Verify), even though it changed in this diff and `.agents/docs/testing/2026-09-17-release-binary-e2e.md`'s new "Added 2026-09-18" section exercises it end-to-end a second time. Tier 3b explicitly notes when a check is "an e2e-tier check the main agent runs and logs, not a bare `verifier` dispatch" (real `report` run); no equivalent sentence exists for this chain, so it's ambiguous by the guide's own convention whether the gap is intentional. Not asserting it must be added — flagging the ambiguity per the dispatch's specific question.

## Checked and clean (worth recording given the brief's specific asks)

- "Newest release" / "the latest" phrasing (README §§ How it works, Install without cloning, Update; `local.md`; `delivery.md`) is attributed only to `sh-agent update`, never to the manual `gh`/browser download chains — no stale claim found.
- `xattr` appears only in README's folded browser chain and in the (exempt, dated) e2e log; no other living doc still instructs a separate post-install `xattr` step.
- The `# styrman-harness $V` header format and the `sh-agent $V` output of `version` (`cmd/sh-agent/main.go:106`, `fmt.Fprintln(out, "sh-agent", buildVersion())`) match what `internal/harness/release.go` actually writes and what `release_test.go` fixtures use — the new `grep`/`test` lines are checking against the real format (modulo the regex-escaping finding above).

## Criteria

1. Dead code — n/a, docs lens; nothing removed without a replacement home except the standalone xattr paragraph, which was correctly folded rather than orphaned.
2. Flow tracing — the version-binding addition traces consistently across all three call sites (README gh route, README browser route, remote.md route) and matches the actual release format; one weakening found (grep regex, above).
3. Standards conformance — n/a, docs lens.
4. Directory layout — n/a, docs lens.
5. Simplicity (DIY-except-security) — n/a, docs lens; note the changed lines are exactly the kind of hand-checked signature/version verification the security exemption covers, so the regex-escaping finding is treated as medium not low.
6. Decoupling — n/a, docs lens.

## Not examined

- `.agents/docs/specs/sh-agent/delivery.md` and `context/AGENTS.md` were grepped for "newest release"/"latest" but not read in full for other drift outside this scope's three files.
- `.agents/docs/deploy/local.md` was grepped, not read end to end.
- No live execution of any of the shell chains (hard limit: no network, no real install) — the `grep -qx` regex finding is a static read, not a reproduced failure.

## Outcome

Written by the main agent, in the commit that carries the fix.

- [medium] `$V` as an unescaped BRE — **fixed**: `grep -qxF` in all three chains (README.md:99, .agents/docs/deploy/remote.md:34); measured: `V=v1X2X2` exits 1 against the real marker.
- [low] `testing/guide.md` has no tier for the manual install chain — **left**, owner `documentation`, filed with the earlier `guide.md` § 5 finding in the handoff's Next Steps; the e2e log is where the chain is exercised today.
