---
owner: main
date: 2026-09-18
---
# The not-found error names why the embedded fallback was refused     (2026-09-18)

## Summary

The operator followed README § Install without cloning on their Mac, ran `sh-agent install --vendor antigravity` from `/tmp`, and got `harness repo not found: run from inside the checkout, pass --repo, or set SH_AGENT_REPO`. Their reply was the right question: "Do I need the repo on my machine? it is supposed to get it from git?? I should just download the binary." The code was behaving correctly — the binary on their PATH was `sh-agent 2615454169fb`, a development build predating the embedded harness (5.4 MB against the release's 6.7 MB), and D14 refuses the embedded fallback for anything that is not a `vX.Y.Z` release so a stale copy is never installed silently. The defect was the message: it listed only the ways out and never said that THIS build carries no harness, so "wrong directory" and "wrong binary" — opposite problems with opposite fixes — read identically.

`ResolveSource` now says which case applies and adds the release-binary path: `harness repo not found and this build is 2615454169fb, a development build rather than a vX.Y.Z release, so it has no harness built in: run from inside the checkout, pass --repo, set SH_AGENT_REPO, or install from a release binary (README § Install without cloning)`. A build with no version stamp at all gets the same sentence with "this build carries no version stamp".

## Type

bugfix

## Files changed

- `internal/harness/repo.go` — `ResolveSource`'s not-found error, with the operator's case in a comment beside it
- `internal/harness/harness_test.go` — `TestResolveSource` now asserts, for each of its eight non-release version strings, that the message names the missing harness, mentions the release-binary path, and quotes the version
- `.agents/docs/bugs/BUG-076-resolve-source-not-found-message.md` (new) and `INDEX.md` — the record and its row
- `.agents/docs/bugs/INDEX.md` — also repairs BUG-075's row, which had landed below the `## Graduated` header with main-table columns; moved into the main table and ordered newest-last (the `bugs` agent flagged it, and it was confirmed before moving)

## Agents dispatched

`bugs`, for BUG-076; its record was read against the diff before committing. No reviewer: ten lines of source and fourteen of test, below the threshold (context § Definition of Done step 2), so the main agent read the diff — the only behavioural change is the error string, and the eight assertions fail against the old one.

## Verification

`gofmt`, `go vet`, `go test -race -count=1 ./...` clean. The new message was read from a real development build, not composed and assumed: `go build -o .agents/scratch/devbuild-sh-agent ./cmd/sh-agent` then running it from `/tmp` printed the sentence above with this checkout's own stamp. Must-fail: restoring the old message fails all eight assertion groups in `TestResolveSource`. The operator's own machine then confirmed the whole path end to end with the correct binary: `sh-agent v1.2.0`, `install --vendor antigravity` writing 2 changed of 130 files, and `status` showing `source embedded v1.2.0`, `binary v1.2.0`, drift none.

## Docs

No doc changed: README § Install without cloning already said what to do, and the message now points at it. The two things the Mac proved — that a `darwin/arm64` binary runs at all, and that a one-vendor install is 130 files — are logged in `.agents/docs/testing/2026-09-17-release-binary-e2e.md` with the rest of the Mac legs as they land.

## Requirements

None touched.
