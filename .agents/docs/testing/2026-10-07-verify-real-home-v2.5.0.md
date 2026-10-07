---
owner: main
date: 2026-10-07
---
# Verification — real-home install at 3293e25 and the guide § 5 smoke on both CLIs (v2.5.0 gate)

Range reviewed: v2.4.1..3293e25.

**The install:** one install into the real home, at 3293e25, on the Linux box, the operator's own (`go run ./cmd/sh-agent install --repo <this checkout>`, per `deploy/local.md` § Build and install).
- **The authority.** The operator said "Go fix all" and "Solve and prove the last issues then cut a new version" (2026-10-07). When asked, they chose "Checkout place (Recommended)": this home had been a release-binary place on v2.3.2, behind the v2.4.0 bridge, and the installer refused the switch without `--repo`.
- **The checks** are the three of `testing/2026-09-03-verify-real-home.md` on each CLI.

## Status: PASS

6 of 6 rows passed, each against its own criterion, at 3293e25.

| CLI | check | result | seconds |
|---|---|---|---|
| Claude Code | context: GANNET, the Antigravity poll rule present, DoD step 2 quoted | 1 of 1, denied 0 | 13 |
| Claude Code | skill: Go review names `golang`, changes nothing | 1 of 1, denied 0 | 31 |
| Claude Code | agent: `verifier` over a comment edit, log written | 1 of 1, denied 0 | 49 |
| Antigravity | context, plus the store check | 1 of 1, denied 0 | 111 |
| Antigravity | skill | 1 of 1, denied 0 | 169 |
| Antigravity | agent: `verifier` log `2026-10-07-verify-greet-comment.md` | 1 of 1 | — |

**The skill rows.** The smoke script wrote its `out.txt` and `err.txt` into the clone, so its own "changes nothing" test first printed FAIL. Judged again with those two files excluded, both CLIs named `golang` and left 0 other changes.

**The context row's skill count.** Claude Code's answer said "35 of 36 … `review-standards` missing". A direct question then got "Yes", with its description quoted. The installed `review-standards/SKILL.md` is byte-identical to the repo. The first answer was a miscount, not a missing skill.

**The Antigravity store check** (conversation `476be3ad…`):
- one global `<RULE[~/.gemini/config/rules/AGENTS.md]>` block (the other `AGENTS.md` block is the demo project's own)
- one `<RULE[…critical-thinking.md]>` block
- GANNET present
- no `<truncated` line

**After the install:**
- `sh-agent status`: installed from commit 3293e25, binary 3293e25a553b, drift none (1789 files match), no shadowed line.
- `hash -r; sh-agent version`: 3293e25a553b.
- `type -a`: one `sh-agent`.

## Not examined
- The Mac.
- The nightly dream timer (not installed).
- An interactive session on either CLI (the operator's).
