---
owner: main
date: 2026-09-17
---
# Verification — real-home smoke checks, both CLIs   (2026-09-17)

Range reviewed: harness `f89d36f` installed in the real `$HOME` (`sh-agent update`: 0 written, 260 unchanged; `sh-agent status`: binary `f89d36f7a17f`, drift none)

## Status: PASS — three checks, both CLIs

Procedure: `2026-09-03-verify-real-home.md`, with Antigravity's tool checks run as on 2026-09-16 (`2026-09-16-verify-antigravity-real-home-complete.md`): the per-run `--dangerously-skip-permissions` flag, approved by the operator for these two runs today; `~/.gemini/antigravity-cli/settings.json` not modified. Throwaway demo service rebuilt under `.agents/scratch/realhome-smoke/demo` (the earlier one was gone): `cmd/server`, `internal/greet/greet.go` with a test, a minimal `AGENTS.md`, its own git repository (`git init` exit 0, own top level). One uncommitted doc-comment edit to `greet.go` as the dispatch target.

## Runs (real `$HOME`, cwd the demo)

| Check | Claude Code 2.1.274 | Antigravity 1.2.5 |
|---|---|---|
| Context questions, no tools: first heading, skill and subagent counts, DoD step 4, the escape rule, the dispatch-budget rule | PASS, exit 0, no tool use — `# AGENTS.md — global agent context`, 27 skills, 7 subagents; quoted today's BUG-061 clause in step 4, the BUG-057 escape rule and the tool-call budget rule | PASS, exit 0 — the same heading and counts, the same three rules; 0 `Failed to parse` lines in its logs |
| Knowledge skill auto-load (`golang` on `internal/greet/greet.go`, change nothing) | PASS, exit 0 — `Skill:golang` in the transcript; findings: the doubled doc comment, no package comment; no edits | PASS, exit 0 — named `golang` and `review-standards` from `~/.gemini/config/skills/`; the same two findings; no edits |
| Agent dispatch (`verifier` on the uncommitted change) | PASS, exit 0 — `Agent:verifier`, waited for it; `Status: PASS`, log `.agents/docs/testing/2026-09-17-verify-greet-comment.md`, `Scope: complete` | PASS, exit 0 — waited for its subagent; `Status: PASS`, log `.agents/docs/testing/2026-09-17-verify-greet.md` |

Demo afterwards: only the dispatch-target edit, the two verifier logs and two testing guides; no commits beyond the demo's own, no source edited by either CLI.

## Defect observed

Claude Code's `verifier` wrote its absent testing guide as `testing/demo/guide.md`. `agents/verifier.md:37` resolves `<app>` to the `cmd/<app>` directory when there is exactly one — here `server` — before the repository name. Antigravity's `verifier`, in the same repository a few minutes later, wrote `testing/server/guide.md`, as the rule says. Recorded as a bug; the check itself passes, since the dispatch, the wait and the verdict were right.

## Not examined

The full skill matrix (`skill-matrix.md`) from a real home; the Mac; the other places; Antigravity's interactive permission prompts. An empty `crashes/crash_<pid>.log` appeared while the Antigravity dispatch run was live and was gone after it exited — a per-process placeholder, not a crash.

## Environment

Linux (Debian 13, amd64), real `$HOME`, harness `f89d36f`, go1.27.1.
