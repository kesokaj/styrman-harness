---
owner: security
date: 2026-10-09
---
# Security review — bug166-agy-keychain   (2026-10-09)

Range reviewed: HEAD (9964037)..working tree (uncommitted; `.agents/scratch/bug166.diff`): `internal/dream/run.go`, `internal/dream/run_test.go`, `.agents/docs/{deploy/local.md,specs/sh-agent/dream.md,testing/guide.md,bugs/INDEX.md}`

## Status: COMPLETE

## Verdict: FINDINGS (low only, nothing blocks)

Nothing here is exploitable, and no secret leaks. The worst problem is a design trade-off that the docs do not state. On Linux the dream job's private HOME gets a *copy* of the agy login, so whatever the headless child does to its credential (refreshes it, rotates it, deletes it after an auth failure) stays in the throwaway directory. On macOS the diff *links* the real `~/Library/Keychains` into that HOME, so every Keychain write the child makes lands in the operator's real login keychain, at the same time as any interactive agy session. Confidentiality does not change. The child runs as the same uid and could already open the real path. Keychain item ACLs are per application, not per HOME. The tool-refusing `settings.json` is still the only real barrier, and it is unchanged. The cleanup paths do not follow the link. `cleanup()` and `removeCopies` both use `os.RemoveAll`, which unlinks a symlink instead of descending into it. `removeCopies` also selects only real directories at the top level. No other remover or walker in the module matches `dream-agy-home-*`.

## Findings
- **[low / confidence: medium]** `internal/dream/run.go:299-309`: the Keychains link ends the private HOME's write isolation for the login on macOS. A dream-run agy that refreshes, overwrites or deletes its Keychain item (for example on the "Authentication required" path) changes the operator's interactive login, and it can race an interactive agy that refreshes at the same moment. On Linux the same actions touch only the copy. This is availability and integrity of the operator's login, not disclosure. Remediation: say it in `specs/sh-agent/dream.md` and `deploy/local.md`. Today they say "linked (never copied)" but not "writes reach the real login". On the owed macOS `dream --all` run (`testing/guide.md` § 3e), also confirm that a failed or locked-Keychain run leaves the interactive agy logged in afterwards. A narrower link, such as `login.keychain-db` alone, would not restore isolation, so no code change is proposed.
- **[low / confidence: high]** `internal/dream/run.go:370-386` (`removeCopies`), test gap: the claim that "a sweep of a killed run's HOME never reaches the real Keychain" is true today but not pinned. `TestAgyHomeMacKeychain` checks only `cleanup()`. `TestRemoveCopies*` (run_test.go:266) builds a stale `dream-agy-home-123` with no link in it. If removal later changed to a walking delete (`filepath.Walk` plus `os.Remove`, or a scrub before removal), it could follow the link into `~/Library/Keychains` and no test would fail. Remediation: in the removeCopies test, put a `Library/Keychains` symlink to a sentinel directory inside the stale `dream-agy-home-*` and assert that the sentinel's file survives. Mutation-check it once with a following delete.
- **[low / confidence: medium]** `internal/dream/run.go:271,276-278,307`: `real := os.Getenv("HOME")` is not checked for being absolute. If HOME is empty or relative, `Stat("Library/Keychains")` is checked against the process cwd, while the symlink target, also relative, resolves against `<private home>/Library/`. The check and the use then point at different places. This fails closed (a dangling link, so agy fails auth), which is why it is low. The token path has the same issue from before this diff. Remediation: refuse with exit 2 unless `filepath.IsAbs(real)` before any login path is built.

## Checked
- Secrets: no credential in code, tests or docs. The test keychain and token are fixture strings under `t.TempDir()`, and `t.Setenv("HOME")` keeps the tests off the real home. `agyAuthHint` appends a constant string and echoes no output. The new `fail` messages name paths only. The `err.Error()` at :132 is `agy: exit status N` or a timeout, as before. `keepFailure` still scrubs through `harness.Redact`. `report` does not send `dream.log` or `state.json`: it sends only `proposals/dream/harness/*` (report.go:404-418). One note, not a finding: `agyAuthHint` matches model-influenced stdout, so a failed run whose answer contains "Authentication required" gets a misleading hint in the log. It misdirects the operator and discloses nothing.
- Removal and walk paths over `~/.config/sh-agent` (every `RemoveAll`, `WalkDir`, `ReadDir` in non-test code was grepped):
  - `cleanup` and `removeCopies` in dream: `RemoveAll` does not follow the link, and `e.IsDir()` reads the entry's own type.
  - `search/dream.go:107` (`dream-agy-` MkdirTemp, removes only its own tmp)
  - `metrics/agy.go:293-302` (`metrics-agy-*` only)
  - `report.go:559-577` (`report-*` dirs, skips symlinks)
  - `SweepStaleDirs` (`zip-`, `release-`, skips symlinks)
  - `Uninstall` (manifest paths, `os.Remove`, not recursive)
  - None of them can reach through `dream-agy-home-*/Library/Keychains`.
- TOCTOU: there is a gap between `Stat` of the Keychains directory and the child resolving the link. Exploiting it needs the same uid, and the private HOME comes from `MkdirTemp` with mode 0700 under conf. There is no privilege boundary to cross. The pre-existing one-hour-age race (a dry run's `removeCopies` against a live run longer than an hour) is unchanged, and the link does not make it worse.
- Injection: no new external input reaches a command. The argv is unchanged, and the prompt is still sent on stdin.
- Authz/authn: the tool-refusal settings file is still written into the private HOME on both OSes, and the new test asserts `"allow":[]`.
- Crypto: none touched, nothing hand-rolled.
- Dependencies: `go.mod`/`go.sum` unchanged, and `runtime` is stdlib. `govulncheck ./...` (go1.27.2, `$(go env GOPATH)/bin/govulncheck`) reports no vulnerabilities reached, plus 1 in required modules that is not called (the known GO-2026-5024). `npm audit` does not apply because there is no JS ecosystem.
- GCP surface: not touched.

## Not examined
- Real macOS behaviour: whether the Security framework or agy finds the keychain only through `$HOME`, and whether agy deletes or rewrites its Keychain item on an auth failure. There is no macOS box here, so finding 1's confidence is medium.
- The `agyLogin`/`fakeCLI` test helpers beyond their use in the diff. `internal/dream/collect.go` beyond `scrub`.

## Outcome

- 1 (the dream job's agy writes to the real Keychain on macOS): fixed in the docs. `specs/sh-agent/dream.md` and `deploy/local.md` say it shares the interactive login and a refresh reaches both. The owed macOS run (`testing/guide.md` § 3e, BUG-166) checks that the interactive agy stays logged in afterwards.
- 2 (no test pins that removeCopies leaves the link's target alone): fixed. `TestRunRemovesLeftoverCopies` now leaves a private HOME behind with a `Library/Keychains` link to a sentinel folder, and the sentinel survives the sweep.
- 3 (HOME not checked for being absolute): fixed. `agyHome` refuses a `HOME` that is not absolute. `TestAgyHomeMacKeychain` has a relative-HOME case, which fails with the check disabled.
