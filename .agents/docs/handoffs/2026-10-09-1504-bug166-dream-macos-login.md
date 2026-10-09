---
owner: main
date: 2026-10-09
---
# Handoff — 2026-10-09 — BUG-166: dream login on macOS

Session ID: 728a4838-8b3c-4819-a62d-2b01e563a77f
Project: /home/ostling-simon/projects/styrman-harness
Branch: main

## Context

On the operator's Mac (v2.6.1), `sh-agent dream --all` refused to run: "no Antigravity login at …/antigravity-oauth-token". agy keeps its login in the macOS Keychain and finds it only through `$HOME/Library/Keychains`. Three operator probes established this (BUG-166 § Repro).

## Completed

- Fix committed and pushed in `11506d1`:
  - `internal/dream/run.go`: `agyHome(conf, goos)` links `~/Library/Keychains` on darwin; `agyAuthHint` names a login prompt.
  - Tests in `internal/dream/run_test.go`.
  - Docs: `README.md`, `specs/sh-agent/dream.md`, `deploy/local.md`, `testing/guide.md` § 3e.
- Reviewer, security and verifier all ran; every finding is fixed or rebutted in the reports' `## Outcome`; the verifier passed.

## Open / In Progress

- Not released: the Mac still runs v2.6.1 and is still broken until a release carries `11506d1`.
- The darwin branch is tested only with `goos="darwin"` on Linux. A real macOS run is owed before the release (`testing/guide.md` § 3e).
- Requirements edits awaiting approval: none.

## Key Decisions

- The Keychain folder is linked, never copied. The token file stays optional on darwin. Linux is unchanged.
- Rejected:
  - Dropping the file check alone: agy would hang on a browser login.
  - Extracting the Keychain item into a token file: unknown format, Keychain prompts, fragile.
  - A dedicated dream login (`dream --login`): needs a design decision, not needed now that the link works.
- Deferred by the operator: the nightly launchd run against a locked Keychain. The failure line now names that case.

## Corrections

- None.

## Agent notes (unverified)

- With the "reply PONG" prompt, the model in all three probes tried `ls -la` first and returned an empty response. If the real dreamer prompt does the same, macOS projects fail as "invalid output". Only the owed real run shows it.
- On macOS the dream job's agy shares the real Keychain login, so a token refresh it makes is the interactive agy's too (security finding 1, documented).

## Environment State

- Branch: main, clean after this handoff's commit. Nothing running.

## Artifacts to Read

- `.agents/docs/bugs/BUG-166-dream-agy-login-macos-keychain.md`: cause, probes, fix, prevention rule.
- `.agents/docs/changes/2026-10-09-1503-bug166-dream-agy-macos-keychain.md`: what changed and what is not covered.
- `.agents/docs/reviews/2026-10-09-{review,security}-bug166-agy-keychain.md`: findings and outcomes.
- `.agents/docs/testing/2026-10-09-verify-bug166-agy-keychain.md`: PASS log.

## Next Steps

1. Operator decision: cut v2.6.2 with `/preflight v2.6.2`. Recommendation: yes, since the dream job is unusable on macOS until a release ships the fix. Before the tag, the gate owes a real macOS `sh-agent dream --all` with a binary built from `11506d1`. Every project's `dream.log` line must read `ok N proposals` or `unchanged`, never `not logged in` or `invalid output`, and the interactive agy must stay logged in afterwards.
