---
owner: main
date: 2026-09-17
---
# Verification — real-home smoke checks after BUG-064, both CLIs   (2026-09-17)

Range reviewed: harness `f949bdc` installed in the real `$HOME` (`sh-agent update`: 7 written, 253 unchanged; `status`: binary `f949bdc74f37`, drift none)

## Status: PASS — three checks, both CLIs

Procedure as in `2026-09-17-verify-real-home-smoke-2.md`: fresh demo clones `.agents/scratch/realhome-smoke/smoke3-c` (Claude Code) and `smoke3-a` (Antigravity), one uncommitted doc-comment edit, no testing guide. The context question added "(6) what Definition of Done step 1 says a verifier dispatch must say when the testing guide is missing"; the dispatch prompt added "following your Definition of Done". Antigravity's two tool checks used the per-run `--dangerously-skip-permissions` flag, approved by the operator for these runs; `settings.json` not modified.

| Check | Claude Code | Antigravity |
|---|---|---|
| Context, no tools | PASS, exit 0 — heading, 27 skills and 7 subagents, `docs/testing/guide.md` one per project, the escape and budget rules, and (6) "testing/guide.md is absent: write it first" (BUG-064) | PASS, exit 0 — the same six answers; 0 `Failed to parse` lines |
| Skill load (`golang`, change nothing) | PASS, exit 0 — `Skill:golang`; the doubled doc comment; no edits | PASS, exit 0 — `golang` named; the same finding; no edits |
| Agent dispatch (`verifier`) | PASS, exit 0 — the dispatch prompt carried "testing/guide.md is absent" (checked in the transcript); `Status: PASS`; `Scope: complete — … testing/guide.md was absent and is now written`; `.agents/docs/testing/guide.md` on disk | PASS, exit 0 — waited for its subagent; `Status: PASS`; `Scope: complete (… .agents/docs/testing/guide.md exists)`; guide on disk |

The BUG-064 defect of `2026-09-17-verify-real-home-smoke-2.md` did not recur: the main agent put the sentence in the dispatch and the guide was written. One run; the repro log holds the 3-of-3 evidence.

## Not examined

The full skill matrix; the Mac; the other places.

## Environment

Linux (Debian 13, amd64), real `$HOME`, harness `f949bdc`, Claude Code 2.1.274, Antigravity 1.2.5.
