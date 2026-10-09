---
owner: reviewer
date: 2026-10-09
---
# Review — bug166-agy-keychain   (2026-10-09)

Range reviewed: HEAD (9964037)..working tree (uncommitted; `.agents/scratch/bug166.diff`), lenses code + docs

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing is broken in the code. `go test -race ./internal/dream` passes. The one defect is a stale restatement: README still says agy runs from a private HOME "holding a copy of the login", which is no longer true on macOS, where the Keychains folder is linked and not copied. The macOS behaviour itself rests on the operator's probe (BUG-166 Repro 3) and no real macOS `dream --all` has run yet. The record and the testing guide both say so, so I did not count it as a finding.

## Findings
- **[low / confidence: high]** `README.md:189` — the `sh-agent dream` row says "agy runs from a private HOME holding a copy of the login". The diff changed that claim in `specs/sh-agent/dream.md:34` and `deploy/local.md:93`, but not here. On macOS the login is a link to `~/Library/Keychains`, and the token file is copied only when it exists. Fix: reword to "holding the login (a copy of the token file, or on macOS a link to the Keychains folder)". `dream.md:79` ("no Antigravity login") is still true.
- **[low / confidence: medium]** `internal/dream/run.go:348` — `agyAuthHint` matches the substring "Authentication required" in the combined stdout and stderr. agy's stdout holds model output, so a project transcript that quotes this phrase can add the "not logged in" hint to an unrelated failure line. Only the log line is affected: the appended text is fixed, so nothing is injected, and the failure itself is real. Fix: match stderr only, or the stream-json event shape, if the operator's probe shows where agy prints it. No test for the stdout case.
- **[low / confidence: low]** `internal/dream/run.go:268-270` — any Stat error on the token file (EACCES, not only ENOENT) counts as "no token". On Linux this still returns the missing-login error. On darwin with Keychains present it skips the token and goes on. Harmless, because the Keychain is then the login, but it is silent. No fix needed unless you want the distinction.

## Criteria
1. Dead code — pass. `runtime` is used, nothing is left behind, and the new function has callers.
2. Flow tracing — pass. The only `agyHome` caller (`run.go:99`) passes `runtime.GOOS`. Both failure paths (non-zero exit, no answer) carry the hint, and `TestRunAgyAuthFailureNamed` runs both exit codes through `Run`, so removing the hint from either path turns a case red. `TestAgyHomeMacKeychain` covers linux-requires-token, darwin Keychain-only, darwin with both, and neither. It checks the link target and that cleanup does not reach through the link into the real folder. Those assertions fail without the change; the signature change alone breaks compilation. Cleanup (`os.RemoveAll`) and the stale sweep (`removeCopies`, `run.go:384`) remove the link without following it. Not covered by a test: the `MkdirAll` and `Symlink` error paths, a Keychains path that is a file and not a directory, and a darwin run through `Run`.
3. Standards conformance — fail (low). README restatement above. BUG-166 record and INDEX row exist and match the code. `deploy/local.md`, `dream.md` and `testing/guide.md` agree with the code. The guide's coverage claim names its gap (a real macOS run is owed).
4. Directory layout — pass. The change stays in `internal/dream`, with tests beside the code.
5. Simplicity (DIY-except-security) — pass. Stdlib only, no new dependency, no crypto written here.
6. Decoupling — pass. `goos` is injected so the OS branch is testable, with no new seams.

## Not examined
- Behaviour on a real macOS machine: whether agy honours the linked Keychains folder, and a locked Keychain under launchd. Both are the operator's step per BUG-166.
- Security of linking the real Keychains folder into the HOME of a child that reads untrusted transcripts. A separate security review (`reviews/2026-10-09-security-bug166-agy-keychain.md`) is in progress; I did not duplicate it.
- The `changes/` entry and the handoff, which are not in the diff. The rest of the range was read.

## Outcome

- 1 (README restatement): fixed. `README.md` § `sh-agent dream` now says the private HOME holds a copy of the token file and, on macOS, a link to `~/Library/Keychains`.
- 2 (hint matched model text): fixed. `agyAuthHint` skips stream-json lines (`{`…), so only agy's own lines count. `TestAgyAuthHint` gained a model-text case and an agy-prompt-on-stdout case; removing the `{` check turns the first one red.
- 3 (any Stat error counts as no token): rebutted. On darwin the Keychain is then the login, and on Linux the run still fails closed with the missing-login error. A wrong error message is the worst case; no data is lost and nothing is granted.
- Not examined (MkdirAll/Symlink error paths, Keychains as a file, darwin through `Run`): left. Error paths return after cleanup like their neighbours. Keychains as a file is refused by the `IsDir` check. A darwin run through `Run` needs a real macOS run, owed before release (`testing/guide.md` § 3e).
