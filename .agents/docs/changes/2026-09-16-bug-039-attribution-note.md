---
owner: changelog
date: 2026-09-16
---
# BUG-039: Claude Code headless commits carried Co-Authored-By trailer despite no-attribution rule     (2026-09-16)

## Summary

The context rule forbade `Co-Authored-By` trailers but had no mechanism — it competed with Claude Code's own `attribution.commit` setting. Three of three headless commits in the sandbox ran with the trailer. Fixed: context § Git now names the CLI setting (`"attribution": {"commit": ""}`) and what suppresses the reminder. `cmd/sh-agent/main.go` gained `attributionHint()` and `attributionOff()` to check whether that setting is present in `<home>/.claude/settings.json` and print an advisory note if not, called at both `install` and `status` command paths.

## Type

bug fix

## Files changed

- `context/AGENTS.md` — § Git sentence names the CLI setting and its reminder, and the command to check it.
- `cmd/sh-agent/main.go` — `attributionHint()` (line 483) and `attributionOff()` (line 499) helper functions; call sites at lines 165 (install) and 319 (status).
- `cmd/sh-agent/main_test.go` — `TestAttributionNote` with nine subtests covering both call sites and error branches (no settings file, setting present/false/missing/custom, invalid JSON, directory, FIFO).
- `README.md` — Install section and `status` output row updated to name the setting checked.
- `.agents/docs/deploy/local.md` and `.agents/docs/deploy/remote.md` — Install and status output sections updated.

## Agents dispatched

reviewer (`reviews/2026-09-16-review-bug-g1-attribution-note.md` APPROVE WITH NITS; three findings fixed: test coverage for vendor guard and JSON error branch, error-handling decision for unreadable advisory file), security (`reviews/2026-09-16-security-bug-g1-attribution-note.md` FINDINGS; four findings fixed: read bounded to regular files with size limit, error message no echo of file bytes, path printed not hardcoded `~/.claude`, parse failure prints both the advisory and the trailer warning), documentation (README and deploy guides), verifier, bugs, changelog.

## Verification

Verifier log `.agents/docs/testing/2026-09-16-verify-linkshelf-sim-fixes.md`: `TestAttributionNote` all nine subtests PASS (including `antigravity_only_never_notes`, `malformed_json_is_advisory_no_echo`, `a_directory_is_advisory_and_never_opened`, `a_fifo_is_advisory_and_never_opened`). Install and status commands both printed the attribution note line against the scratch home lacking the setting. Counts pass 64, fail 0, skip 1 (the skip is pre-existing `TestFindRepo`, environment-dependent).

Owed but not done in this commit: sandbox row `attribution-off` (headless commit with the setting present in the home), recorded later in testing logs.

## Docs

`README.md` (Install section and status row), `context/AGENTS.md` (§ Git), `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md` all updated to name the setting and its CLI reminder; the settings key the code checks and prints is `"attribution": {"commit": ""}` (the `"pr"` key is a separate optional recommendation for PR footers, not checked by the tool).

## Requirements

n/a
