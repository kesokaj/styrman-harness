---
owner: reviewer
date: 2026-09-18
---
# Review — c868ed4-bc745dc-docs (2026-09-18)

Range reviewed: c868ed4..bc745dc (git diff c868ed4...bc745dc)

## Verdict: APPROVE WITH NITS

The three places that had to agree on BUG-079's remedy — README.md, `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md` — do agree: same order (`sh-agent version` → `hash -r`/`rehash` → `type sh-agent` → `which -a sh-agent`) and the same remedy text (`update` replaces `~/.local/bin/sh-agent` and nothing else; a shell's cached lookup keeps running the old file until `hash -r`/`rehash`). The one real gap is `.agents/docs/specs/sh-agent/delivery.md`, which still describes the release-binary `update` path as "replaces the binary and reinstalls" with no mention of the running-vs-installed distinction this scope just spent three files establishing — a reader who trusts the delivery model alone would not learn that the binary they invoked may not be the one just replaced. Nothing in the diff reroutes a line on the README `How it works` map: the new `note:` line is diagnostic output inside the existing `update` node, not a new edge or node.

## Findings

- **[low / confidence: medium]** `.agents/docs/specs/sh-agent/delivery.md:33` — "`update` there asks GitHub for the newest release, refuses a tag that is not `vX.Y.Z` or is not newer, verifies `SHA256SUMS.sig` against the key compiled into the binary and each binary against its checksum, replaces the binary and reinstalls." "replaces the binary" is unqualified and, read next to the now-detailed Verify sections, invites the wrong reading (that the invoked binary itself is replaced). This is pre-existing wording the scope did not touch, but the scope was specifically asked to check it against the new BUG-079 nuance and it has not been updated to say "replaces `~/.local/bin/sh-agent`" or otherwise flag the running-binary distinction. The same unqualified phrase also persists, untouched by this diff, at `README.md:59`, `README.md:163`, `.agents/docs/deploy/remote.md:59`, and the comment at `.agents/docs/deploy/local.md:40` — none of these were in scope to fix here, but they are the same gap and should get the same one-line qualifier when next touched.
- **[informational — standing question, no code owed]** `.agents/docs/testing/guide.md` § 5 ("Real-home smoke checks") still lists only three checks (context questions, one skill load, one agent dispatch) via `.agents/docs/testing/2026-09-03-verify-real-home.md` / `2026-09-04-verify-mac-real-home.md`; neither referenced file mentions `sh-agent version`, `type sh-agent`, `hash -r`, or the `shadowed` line (confirmed by grep: zero hits). This scope did not add such a check either. Per the repeated standing question from prior docs reviews: I read this as an accepted e2e-only tier (BUG-078/BUG-079 verification lives in the deploy guides' Verify sections and the bug records' own manual e2e logs, not in the automated real-home smoke checklist) rather than a hard finding — but it is now the third scope in a row to leave it unresolved, so it is named again rather than silently dropped.

## Criteria (docs lens)

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens (the three Verify-section edits and `staleRunningNote`'s message text were traced for wording agreement; they agree)
3. Standards conformance — n/a — docs lens; describing-text rule applied above (delivery.md finding)
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- Did not run `gofmt`/`go vet`/`go test -race` (code correctness is out of this docs-lens brief; the new `RunningElsewhere` code and its tests were read only for wording that docs must match, not for correctness).
- Did not check every other file under `.agents/docs/specs/` beyond `delivery.md`, nor `skills/` contents beyond a targeted grep for "replaces the binary" / "replaces … binary" phrasing — no other hits surfaced in that grep across `context/`, `agents/`, `skills/` (excluding vendored), so treated as clean, but a full read of every skill file was not performed.
- Did not check `.agents/docs/bugs/INDEX.md` or the BUG-079 record itself — both were committed in `c868ed4`, outside the reviewed range.
- Did not evaluate whether a `changes/` entry is owed for the `bc745dc` fix commit — no such entry exists in the diff; flagged here only as a note, not scored as a finding, since it falls outside the specific comparison points the brief named.
