---
owner: reviewer
date: 2026-09-18
---
# Review — worktree-docs   (2026-09-18)

Range reviewed: working tree — changed files: README.md, .agents/docs/deploy/remote.md (git diff HEAD -- README.md .agents/docs/deploy/remote.md); no untracked, non-ignored files in scope (`.agents/docs/testing/2026-09-18-verify-post-v1.2.2-docs.md` is untracked but was excluded per the dispatch's file list and read only for context, not as review scope).

## Verdict: APPROVE

The diff adds a `test` gate comparing the downloaded binary's `sh-agent version` output against the signed `SHA256SUMS` version marker, and switches the README browser route to a `mktemp -d` scratch directory, in all three install chains (README gh route, README browser route, `remote.md` § First install). This is exactly the fix the newest `changes/` entry (`2026-09-18-release-v1.2.2.md`, "Known and accepted, not fixed in this tag") describes as owed in the commit after the v1.2.2 push, and it matches: SHA256SUMS' marker line is written as `# styrman-harness vX.Y.Z` (confirmed in `remote.md` § Releases build recipe), `sh-agent version` prints `sh-agent vX.Y.Z` (confirmed in `cmd/sh-agent/main.go`'s `version` case and `buildVersion`/ldflags), and the `sed` extraction matches that format — the comparison is correct. All three chains keep the `test` line before the final `mv`, so a failed check still stops the move, and the surrounding prose ("the download is verified before it is moved", README's `install`/`status` line, `remote.md`'s BUG-078 paragraph) is unchanged and still accurate. The browser fence's new `mktemp -d` + explicit `mv` of the three named files into it, then `cd ~ && rm -r "$D"`, matches the updated comment ("in a scratch directory so nothing else in ~/Downloads is touched") and no longer risks clobbering fixed names in `~/Downloads`. No stale restatement of the old comparison-free chain, the old `~/Downloads`-in-place browser route, or the old "prints its version before it is moved" wording was found elsewhere in the living docs.

## Findings

None. Checked for stale/dropped restatements via `git grep` of the old wording's subject nouns ("prints its version before", "ASSET\" version", "sh-agent-linux-amd64 version", "cd ~/Downloads", "in ~/Downloads") across README.md, `.agents/docs/deploy/{local,remote}.md`, `.agents/docs/testing/guide.md`, `.agents/docs/specs/sh-agent/delivery.md`; the only remaining hits are inside `.agents/docs/testing/2026-09-17-release-binary-e2e.md` (a dated log, out of scope by the brief) and past `reviews/`/`bugs/`/`changes/` records (also out of scope), none of which assert current install-chain behaviour.

- `local.md` § Verify and `remote.md` § Verify were not touched by this diff and do not need to be: both already point at `sh-agent version` unqualified + `sh-agent status` + `which -a`, which is the post-move check, distinct from the new pre-move `test` gate this diff adds; no inconsistency.
- `.agents/docs/testing/guide.md` has no tier covering the manual install-chain verification at all (pre-existing, not introduced or worsened by this diff) — not a new finding, but noted under Not examined below since the brief asked for a consistency check and there is nothing there to compare against.
- `.agents/docs/specs/sh-agent/delivery.md` does not describe the human verification route (gh/browser chains) at all, only `sh-agent update`'s internal checks — this diff does not change what delivery.md claims, so no conflict.
- README's `xattr` quarantine note (line 142-146) sits after the gh/browser fences and is unaffected by the `mktemp -d` change — `xattr -d` still targets `~/.local/bin/sh-agent`, the final installed path, regardless of which scratch directory the file passed through.

## Criteria

1. Dead code — n/a — docs lens
2. Flow tracing — pass: the `test` gate's inputs (binary's own `version` output, SHA256SUMS' marker line) are produced exactly as the build recipe in `remote.md` § Releases writes them; the comparison holds for a real release.
3. Standards conformance — pass: matches the global context's describing-text rule (comment above each fence still describes what the fence does) and the harness's "single `&&` chain, failed check stops the move" convention.
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — pass: plain `test`/`sed`/`mktemp`, no new dependency, no hand-rolled verification (signature/checksum verification itself is unchanged, still `ssh-keygen -Y verify` + `sha256sum -c`).
6. Decoupling — n/a — docs lens

## Not examined

- `.agents/docs/testing/guide.md`'s absence of any tier for the manual install-chain checks — flagged above as a pre-existing gap, not verified against a "should cover this" requirement since no rule was found mandating it.
- Did not exercise any of the three shell chains (hard limit: no network, no real install).
- Did not read `.agents/docs/testing/2026-09-18-verify-post-v1.2.2-docs.md` past its header (it was IN PROGRESS/unfinished at time of review and is not part of the reviewed scope).
- Did not diff against `.agents/docs/specs/` files other than `delivery.md`, or against every skill/agent/rule file for unrelated restatements of "version" or "SHA256SUMS" wording outside the grep set already run.
