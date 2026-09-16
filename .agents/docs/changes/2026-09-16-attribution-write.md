---
owner: changelog
date: 2026-09-16
---
# `sh-agent install` writes Claude Code's attribution setting when absent     (date: 2026-09-16)

## Summary
Implements ADR-025 to resolve BUG-039: `sh-agent install` (and its alias `sync`, for vendor `all` or `claude` only) now writes `"attribution": {"commit": "", "pr": ""}` to `~/.claude/settings.json` when the `attribution` key is absent, creating the file if missing. A present key of any value is left alone and the advisory note stays. File edits use textual insertion re-parsed before writing so the operator's other keys and formatting remain byte-identical. `--dry-run` reports and writes nothing. No state of the file fails install; the setting is advisory (R4). `uninstall` leaves the key (the trailer must stay off with or without the harness). Supersedes ADR-010's "never touch vendor settings" clause on this one point only; the permission-write refusal it stated still stands.

## Type
feature

## Files changed
- `cmd/sh-agent/main.go` — `attributionSet` (file creation, textual insertion, Lstat guard, size/format checks, Sync before rename), `readSettings` (shared helper for `attributionSet` and `attributionOff`), `abandonTemp` (temp-file cleanup with error reporting), `attributionOff` refactored to use `readSettings`, install/sync case wired to call `attributionSet` before `attributionHint`
- `cmd/sh-agent/main_test.go` — `TestAttributionSet` (14 subtests: missing file, other keys with formatting/mode preservation, CRLF, empty object, off/custom value, malformed, non-object JSON, null, >1 MiB bound, symlink guard, dry-run, vendor-agy no-op), `TestAttributionNote` cases `no settings file` and `other settings only` updated for new order (install leg expects `settings:` line and no note; status leg has no note because key now present); new case `status alone on a home without the key` keeps the note
- `REQUIREMENTS.md` — R4 sentence approved 2026-09-16: "`install` also sets one key in the operator's `~/.claude/settings.json`, `attribution` (commit and PR trailers off), only when the key is absent; it never changes or removes another key, and `uninstall` leaves it."; § Non-goals exception approved 2026-09-16: narrowed "Editing `settings.json` or any vendor setting" to name the exception and ADR-025
- `context/AGENTS.md` § Git — "`sh-agent install` sets it when absent; `status` says when a different value is set"
- `README.md` — § Install: attribution sentence and line reworded to "the installer touches no vendor setting except the attribution key (ADR-025)"; `status` row updated for the absent-file note-firing condition
- `.agents/docs/deploy/local.md` — § Verify and § First session after install (line 37 reworded to "the installer touches no vendor setting except the attribution key (ADR-025)")
- `.agents/docs/deploy/remote.md` — § First install
- `.agents/docs/architecture/DECISIONS.md` — ADR-025 (full decision: context, decision, consequences naming R4/Non-goals changes and the operator's approval); ADR-010's heading marked "Superseded in part by ADR-025"

## Agents dispatched
reviewer ×5 (three cold-read rounds with round-3 findings: one doc-line residue left as Open rows; code lens APPROVE WITH NITS: three findings fixed — null handling, error-path cleanup, symlink/size guards; docs lens REQUEST CHANGES: two findings fixed — status note-firing condition restatement, README antecedent noun), security ×1 (five low-severity findings: TOCTOU symlink/FIFO race fixed with `O_EXCL`, missing `Sync` fixed, non-object JSON output fixed, CRLF line-ending handling confirmed, uninstall promise updated; four earlier defects closed: dangling/live symlinks rejected, >1 MiB guard added, bare `null` caught, `abandonTemp` helper added), verifier ×3 (third and final pass: three runs, last on final tree after security fixes; `go test -race -v -count=1 ./...` PASS: counts pass 79, fail 0, skip 1; install/status/reinstall e2e PASS; coverage: `TestAttributionSet` branch-by-branch vs `attributionSet` live; gaps: `abandonTemp` error paths and `readSettings` open failure untested), documentation ×1 (README, deploy/local.md, deploy/remote.md), changelog ×2 (ADR-025 and this entry), bugs ×1 (BUG-044: security agent's crafted-input probes left seven untracked files in repo root; removed; prevention rule added to Definition of Done step 1).

## Verification
Unit: `go test -race -v -count=1 ./...` PASS (pass 79, fail 0, skip 1). Integration: `go run ./cmd/sh-agent install --home .agents/scratch/verify-attribution-write-3` created settings file with key and printed `settings: …/.claude/settings.json created with "attribution": {"commit": "", "pr": ""} — Claude Code's commit and PR trailers are off`; `status` printed no note (key now present, rendered 93 lines). E2E: second install to same home idempotent (0 written, no settings line). Both CLIs: `sh-agent list` returned expected counts (Skills 27, Subagents 7). Logs: `.agents/docs/testing/2026-09-16-verify-attribution-write.md` (third pass after security fixes; two prior passes overwritten; full matrix of crafted inputs exercised in `.agents/scratch/secrev/`).

## Docs
README.md (Install section, status row), .agents/docs/deploy/local.md (Verify and First session after install), .agents/docs/deploy/remote.md (First install), context/AGENTS.md (§ Git sentence). REQUIREMENTS.md changed — R4 and § Non-goals — both approved by operator 2026-09-16.

## Requirements
R4 (one-sentence addition approved 2026-09-16), § Non-goals (exception clause approved 2026-09-16). Both in REQUIREMENTS.md, committed in the same commit as ADR-025.
