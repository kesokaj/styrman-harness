---
owner: main
date: 2026-09-21
---
# Verification — real-home smoke checks on v1.6.0, both CLIs   (2026-09-21)

Range reviewed: harness `825bf86` = `v1.6.0` installed in the real `$HOME` (`sh-agent update` 17:28, 258 files, drift none). The operator's "Update locally and run tests. Make sure it works." Local checks first: `gofmt -l cmd internal` clean, `go vet ./...` clean, `go test -race -count=1 ./...` exit 0, scratch install `.agents/scratch/x` drift none.

## Status: PASS

Procedure as in `2026-09-21-verify-real-home-smoke-5.md`: fresh demo clones `smoke6-c` (Claude Code) and `smoke6-a` (Antigravity), each CLI in its own subshell; one uncommitted doc-comment edit for check 3; no commit in either clone. Antigravity checks 2 and 3 with the per-invocation `--dangerously-skip-permissions` in the throwaway clone (guide § 5).

| Check | Claude Code | Antigravity |
|---|---|---|
| Context questions, no tools: first heading, skill and subagent counts, DoD step 3, the grilling command, the canary word, the Project layout line | PASS — `# AGENTS.md — global agent context`, 27 / 7, step 3 quoted, `/grilling`, GANNET, "`AGENTS.md` alone, with no `CLAUDE.md` or `GEMINI.md`, except a one-line stub where…" | PASS — same; "no truncation is reported"; named `/grill-me` from habit beside `/grilling` (as on 2026-09-17); store: RULE block present, 0 `<truncated` markers, canary line inside |
| Knowledge skill auto-load (`golang` on `greet.go`, no edits) | PASS — `golang` loaded, two style nits, no edits | PASS — `golang` and `review-standards` loaded, two findings, no edits |
| Agent dispatch (`verifier` on the uncommitted edit) | PASS — log `2026-09-21-verify-greet-comment.md`, Status PASS; the verifier also wrote the missing `guide.md` | PASS — log `2026-09-21-verify-greet.md`, `## Status: PASS`; `Failed to parse skill` count 0 |

Mac (Antigravity-only, binary place): `sh-agent update` on 2026-09-21 17:39 by the operator — "verified and installed sh-agent-darwin-arm64 v1.6.0", 2 written, 127 unchanged, 1 removed (the init-project `CLAUDE.md` asset). Its canary question remains the operator's.

## Not exercised

The devbox; the full skill matrix; a first-session-after-upgrade run elsewhere.
