---
owner: main
date: 2026-09-29
---
# Search release readiness: adoption, optimizations, gate runs, security fixes   (2026-09-29-2326)

Plan: `plans/2026-09-29-search-release-readiness.md` (closed). Asked for by the operator: everything done before a version is cut, nothing waiting.

## Changed

- **Context:** `sh-agent search --open` joins the resume reading order, paid for by cutting the `--open` clause from the "Search first" entry (21,996 of 22,000 bytes). `skills/handoff/SKILL.md` restates the order with the step.
- **Optimizations, kept where measured:** insert statements prepared once per reconcile; reference parsing skips lines that cannot match (the same 40,842 references on 1,077 files); the `refs(oid)` index dropped. Build 2.2 s to 1.7 s, index 24.1 MB to 21.7 MB. Rejected: merging the FTS index after a build (the file grew).
- **Security** (`reviews/2026-09-29-security-aae16af-f41e5c0.md`): a database file or companion that is not a regular file is refused; index files are created 0600 and an older 0644 file is tightened on open.
- **Docs:** the spec's failure modes, file modes and measurements; the `search-first` matrix row.

## Measured

- **Adoption** (`testing/2026-09-29-search-adoption.md`, a clean A/B of the whole harness per arm): with the new skills the main agent searched in 3 of 6 large-trail sessions and 9 of 12 sample-app sessions; with the resume step, in 6 of 6 large-trail sessions. Quality was identical in every arm: 42 of 42 sessions delivered.
- **Cost:** no saving overall. Sample apps: median lane cost 18.02 (baseline) against 18.17 USD. Large trail: the search arms cost 15 to 19% more per lane; the trail-question session alone read under half as much and cost 12% less.
- **Antigravity** (`testing/2026-09-29-search-agy-gate.md`): 10 of 10 changed skills and agents completed a run. agy 1.2.13 auto-denies any command containing `$(…)` in headless mode, so `init-project` and `dream` ran with the per-run skip flag in the sandbox.
- **Real home** (`testing/2026-09-29-verify-real-home-search.md`): updated by `sh-agent update` on the operator's word, now at 8e1b7f5; the smoke checks pass on both CLIs.

## Agents dispatched

`reviewer` ×5 (three cold reads, code, docs lens), `security`, `verifier` (`testing/2026-09-29-verify-search-release-readiness.md`, PASS over aae16af..1634787).

## Verification

`gofmt -l cmd internal` empty; `go vet ./...` exit 0; `go test -race ./...` exit 0, 362 passes; mutations O1–O3 and S1–S3 each compile and fail a named test.
