---
owner: main
date: 2026-09-29
---
# Verification — real home on the Linux box after v1.12.0, both CLIs   (2026-09-29)

Range reviewed: v1.12.0 (95f8fca) installed as 0beb429 (`sh-agent update` from this checkout, 23 written, 266 managed, binary 0beb429bfa0a); the checks run from the real `$HOME` in fresh clones of a fixture service under `.agents/scratch/smoke/` (a trusted path for Antigravity, whose real settings say `toolPermission: always-proceed`).

## Status: PASS on both CLIs — the three smoke checks, plus the first runs of `/dream`, the resume clause, the `dreamer` and the linked-store row

The operator asked, after v1.12.0, that everything the gate had listed as not checked be checked ("I want everything checked before tag"). These are those runs; the rule that makes them a gate condition is in `AGENTS.md` § Release.

## Fixture

A three-file Go service (`/healthz`, `/hello`, `internal/greet`), an `AGENTS.md` with one rule, a bug index with a recurring row, one handoff, and a seeded `.agents/memory/` store (four memories: one graduated by the `AGENTS.md` rule, one broken path, a two-session pair on one fact, plus the index), `.claude/settings.local.json` pointing at it, both gitignored; the fixture commits clean. One clone per run.

## Claude Code 2.1.284 (`claude -p`, real home, installed skills)

| Check | Result |
|---|---|
| Context questions, no tools | PASS — first heading "AGENTS.md — global agent context", 28 skills and 8 subagents, last word GANNET, `dream` and `dreamer` named for their roles |
| Skill auto-load (`golang` on `greet.go`, no edits) | PASS — named `golang`, no edits (its `go` commands were refused by the run's tool list, as intended) |
| Agent dispatch (`verifier` on an uncommitted comment edit) | PASS as a dispatch — ran in the clone, wrote `testing/2026-09-29-1602-verify-greet.md` with a Status line; its verdict is FAIL because the fixture has no test covering the edit, which is the verifier's correct reading of a test-less fixture |
| `/dream --dry-run` | PASS — `store: <clone>/.agents/memory (4 memory files; default path …: 0)`, `graduated`, `broken`, `promote` present, no split, store unchanged, tree clean |
| `/retro --dry-run` | PASS — the `dreamer` dispatched, `retros/2026-09-29-dreamer.md` written and cited, the pair as findings with `rule_ref null` |
| `init-memory-linked-store` matrix row (`.agents/memory` a committed symlink to an empty directory outside the clone; the default store seeded) | PASS — Step 5b stopped at the trust check ("a committed link would carry your memory wherever it points"), nothing behind the link, the five default-store files untouched, no settings file; Step 6 appended the ignore lines, which the row now allows |
| Interactive settings write | PASS — this session's Write tool created `.claude/settings.local.json` in a scratch clone without a guard refusal; headless runs refuse it, as every probe recorded |

## Antigravity CLI 1.2.13 (`agy -p`, real home, `--mode accept-edits`)

| Check | Result |
|---|---|
| Context questions, no tools | PASS — the same four answers, GANNET from the canary line; the Resume Protocol's memory file named as `.agents/memory/MEMORY.md`, present in the clone |
| Skill auto-load (`golang` on `greet.go`, no edits) | PASS — `golang` and `review-standards` loaded, tree clean |
| `/dream --dry-run` | PASS — the same store line with both counts, verdicts present, no split, store unchanged, tree clean |
| `/retro --dry-run` | PASS — the `dreamer` dispatched (the derived `agent.md`), its artifact written and cited, findings with `rule_ref null` |
| Logs | no `auto-denied` line in any run's stderr; no `invalid allow entry` in `cli.log` |

## Also this pass

- `TestReinstallBuildsAndInstalls` no longer globs the shared scratch: it stats its own `sh-agent-build-<pid>` (BUG-104); mutation protocol run (the mutant compiled, the named test failed on the leaked directory, restored, green).
- `/init-project` refuses to write through any symlinked path (one sentence in Step 1, the older writes the security round had named as out of range).
- `AGENTS.md` § Release: a gate that has not run the changed skills and agents on Antigravity is NO-SHIP.

## Not verified here

The Mac and the devbox (the operator's `sh-agent update` there); `/distill` dispatching the `dreamer`; the `init-memory` and `init-memory-absent` rows on Antigravity; `/init-project`'s update path on Antigravity.

## Commands

`sh-agent update` → `23 written, 243 unchanged, 0 removed, 0 kept — 266 files managed`; `sh-agent status` → installed from 0beb429, binary 0beb429bfa0a; `tail -c 60` of both rendered contexts ends with GANNET. Run kits: `.agents/scratch/smoke/run-claude.sh`, `run-agy.sh` (removed after this log).
