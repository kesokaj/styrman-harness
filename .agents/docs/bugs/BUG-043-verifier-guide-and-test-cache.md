---
owner: bugs
date: 2026-09-16
id: BUG-043
status: fixed
area: agents/verifier
---
# BUG-043: verifier skipped the first-use guide and let a cached Go test result hide an e2e regression

## Symptom
Two symptoms, both from the Claude Code lane of the linkshelf simulation (`.agents/docs/testing/2026-09-16-linkshelf-simulation.md`, gaps G16 and G17):

(a) Across twelve verifier dispatches on Claude Code, no `testing/linkshelf/guide.md` was ever written, although `context/AGENTS.md` step 3 says the verifier creates it on first use. Antigravity's third dispatch, on the same project, did write one.

(b) The project's own run had `go test -race ./...` print `ok linkshelf/cmd/linkshelf (cached)` while a `-count=1` rerun showed five of six e2e tests failing. The e2e tests build the binary at run time, which Go's test cache cannot see. The project recorded this itself as its own BUG-002.

## Repro
Sandbox `/home/kesokaj/sandbox/runs/sim-claude-linkshelf`: `.agents/docs/testing/` holds four verifier logs and no guide file; `.agents/docs/bugs/BUG-002-go-test-cache-hid-e2e-regression.md` records symptom (b) from inside that project.

## Root cause
The old guide sentence read "a procedure that no guide describes," which an `AGENTS.md` that lists the test commands satisfies on a plain reading — so the verifier judged the guide already described and never wrote one. Separately, the counts-run instruction named `-race` and `-v` but said nothing about the test cache, so a plain `go test -race ./...` was accepted as a real run even when Go served a cached pass for a binary-building e2e suite.

## Fix
`agents/verifier.md` (uncommitted at dispatch time, committed with this record): the guide is now written whenever `testing/<app>/guide.md` does not exist as a file — "an `AGENTS.md` that lists the commands does not replace it." The counts run always carries `-count=1` — "a verifier never takes a cached result … an `ok … (cached)` line is not a run."

## Prevention
An agent's "first time" duty is keyed to whether a specific file exists, never to a judgement call about whether another file "describes" the same thing. A verifier's counted test run is never a cached one — always force `-count=1` on the counts command, in this repo and in any project's own verification guide.

## Related
`.agents/docs/reviews/2026-09-16-review-bug-g16-verifier-cache-guide.md` (APPROVE WITH NITS; its one medium finding — a conditional `-count=1` trigger being undecidable — was fixed by deleting the condition). `.agents/docs/testing/2026-09-16-verify-bug-043-verifier.md` (verification, in progress) plus a live verifier dispatch in the sandbox at the fixed definition, owed and tracked in the simulation log's Re-runs section. Same simulation as BUG-039..042.
