---
owner: main
date: 2026-09-03
---
# Next steps done: preflight base, verifier orientation, two skills pruned, ADR-010, retro applied, first real install     (date: 2026-09-03)

## Summary

The operator asked for every open item to be done. Nine tasks; two stay with the operator because they need the Mac.

1. **Real home on this box** — first `sh-agent install` into the operator's `$HOME` on the operator's explicit instruction (the manifest records the commit). The Mac (`sh-agent update`, one interactive start per CLI) is still the operator's.
2. **`/preflight` base ref** — only `*-preflight-*` reports count, a `..worktree` line names its base before the dots, and an unattended first release takes the root commit and says so. Closes the SHIP/NO-SHIP divergence seen in the matrix.
3. **`verifier` orientation** — `pwd` and `git status --short` first; a described change set that is not in the tree is `Status: BLOCKED` with what was seen, never FAIL; the log's Status line gains a `Scope:` line.
4. **Harness state prompts** — decided (ADR-010): no vendor-settings writes (a `REQUIREMENTS.md` non-goal), one approval per session accepted and documented in README and `deploy/local.md` with the user-side setting named.
5. **Skill matrix** — decided (ADR-010): recorded as `testing/skill-matrix.md` (every prompt, check, flag); scripts stay in scratch.
6. **Pruning** — `grill-me` (launcher for `grilling`; `/grilling` works on both CLIs) and `skill-repair` (repairs a `skills`-CLI manifest the harness does not use) removed; 28 skills. `VENDORED.md`, the context command list, and the repo `AGENTS.md` updated.
7. **Sandbox cleanup** — 62 test clones removed, seeded feedback entries and proposals removed from the sandbox home.
8. **Retro** — one real correction this session (sandbox testing reported complete at 3 of 30 skills) logged to `~/.config/sh-agent/feedback.jsonl`; proposal file written and both proposals applied: the verifier's `Scope:` line and the repo `AGENTS.md` Testing rule that "complete" is only said after the full matrix.

## Type

task

## Files changed

- `skills/preflight/SKILL.md`, `agents/verifier.md` — rules above
- `skills/grill-me/`, `skills/skill-repair/` — deleted; `skills/VENDORED.md`, `context/AGENTS.md`, `AGENTS.md` — updated
- `README.md`, `.agents/docs/deploy/local.md` — the one-time approval and the user-side setting
- `.agents/docs/architecture/DECISIONS.md` — ADR-010
- `.agents/docs/testing/skill-matrix.md` — new; `.agents/docs/plans/2026-09-03-rebuild.md` — section 12

## Agents dispatched

None; diff read by the main thread (markdown only, no Go change).

## Verification

`gofmt`, `go vet`, `go test -race` clean; scratch install 256 files, no drift, inventory 28 skills; sandbox reinstalled, Antigravity lists `grilling`, `bug`, `feature` and its log has no parse errors; sandbox sanity runs on Claude Code: an unattended `/preflight` now names the root commit as base and ignores non-preflight reports; a `verifier` dispatched on a change set that does not exist orients, reports `Status: BLOCKED`, writes no log, and stays in its directory.

## Requirements

R4 (install idempotent, removes what the repo no longer has), R7 (state stays in `~/.config/sh-agent/`), non-goal "editing settings.json" upheld.
