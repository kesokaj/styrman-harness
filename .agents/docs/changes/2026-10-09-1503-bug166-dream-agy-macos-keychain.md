---
owner: main
date: 2026-10-09
---
# BUG-166: the dream job logs in to Antigravity on macOS

## Type

bug fix

## Summary

On macOS, `sh-agent dream` refused every run with "no Antigravity login at …/antigravity-oauth-token". agy keeps its login in the Keychain there and finds it only through `$HOME/Library/Keychains`. The dream job's private HOME had neither. Three operator probes on the Mac established this (BUG-166 § Repro).

## Changes

- **`internal/dream/run.go`.**
  - `agyHome(conf, goos)`: on darwin it links `~/Library/Keychains` into the private HOME (never copied); the token file is copied when present. The run is refused only when neither exists, or when `HOME` is not absolute.
  - `agyAuthHint`: names a browser-login prompt in agy's own (non-stream-json) output lines on both failure paths, so a locked Keychain no longer reads as `invalid output`.
- **`internal/dream/run_test.go`.** `TestAgyHomeMacKeychain`, `TestAgyAuthHint` and `TestRunAgyAuthFailureNamed` are new. `TestRunRemovesLeftoverCopies` now checks that the sweep leaves the link's target alone.
- **Docs.** `README.md`, `specs/sh-agent/dream.md` and `deploy/local.md` describe the macOS login. `testing/guide.md` § 3e: a change to the login path owes a real macOS `dream --all` before release.
- **Records.** `bugs/BUG-166-…` and its INDEX row.

## Agents dispatched

- **reviewer:** `reviews/2026-10-09-review-bug166-agy-keychain.md`, APPROVE WITH NITS. Two findings fixed, one rebutted.
- **security:** `reviews/2026-10-09-security-bug166-agy-keychain.md`, three lows, all fixed.
- **verifier:** `testing/2026-10-09-verify-bug166-agy-keychain.md`, PASS: 649 tests, 0 failures.
- The goldfish diagnosis was skipped: the operator's probes located the cause in one function.

## Not covered

- The darwin branch is tested only with `goos="darwin"` on Linux. A real macOS `dream --all` with the fixed binary is owed before the release that carries it.
- The nightly launchd run against a locked Keychain is unchecked (deferred by the operator, 2026-10-09).
