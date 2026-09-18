---
owner: reviewer
date: 2026-09-18
---
# Review — harness-simplify   (2026-09-18)

Range reviewed: internal/harness and cmd/sh-agent (whole-file read, not a diff): all non-test `.go` files in `internal/harness` (adopt.go, apply.go, derive.go, frontmatter.go, inventory.go, pathguard.go, plan.go, release.go, render.go, repo.go, report.go, source.go, state.go) and `cmd/sh-agent/main.go`.

## Verdict: APPROVE WITH NITS

Nothing here is broken. The codebase is unusually well-defended (every non-obvious choice carries a comment naming the bug or review round that produced it) and I found no dead exported code, no case of two functions silently disagreeing on the same question (the brief's `binOnPath`/`PathShadows` and lexical-vs-`os.SameFile` pairs are already gone — `PathShadows` and `RunningElsewhere` both use `os.SameFile` identity consistently, one call each). The one real issue is `cmd/sh-agent/main.go` carrying a self-contained ~220-line settings.json editor that has nothing to do with CLI orchestration, plus a small cluster of literally-duplicated message strings inside that same code. Everything else is either already minimal or a single-purpose helper with one caller.

## Findings

- **[medium / confidence: medium]** `cmd/sh-agent/main.go:676-895` (`attributionSet`, `attributionHint`, `attributionOff`, `readSettings`, `abandonTemp`) — a self-contained JSON-file editor (read `~/.claude/settings.json`, insert one key, atomic rename) with no dependency on flag parsing or command dispatch. It sits in the 945-line entrypoint file next to install/update/status orchestration, so it and the CLI wiring share one file and one diff for unrelated changes. Every other piece of logic this size in the repo (install, plan, state, release, report) lives under `internal/harness` and is called once from `main.go`; this cluster is the outlier. Shorter version: move these five functions verbatim into a new `internal/harness/attribution.go`, export the two entry points (`AttributionSet`, `AttributionHint`) and call them from `installFiles`/`status` as already written — deletes ~220 lines from `main.go` with no behavior change.

- **[low / confidence: high]** `cmd/sh-agent/main.go:703,711,716,721,725,784,816` — the exact format string `"settings: %s could not be written (%v); add the key by hand\n"` is repeated verbatim at seven call sites inside `attributionSet`/`abandonTemp`, and `"settings: %s could not be edited safely (the result would not parse); add the key by hand\n"` is repeated at two more (lines 756, 775). A wording change (there have already been several review rounds touching this file) has to be made in up to nine places to stay consistent. Shorter version: one `settingsWriteFailed(out io.Writer, path string, err error)` helper (3 lines) replacing the 7 identical `fmt.Fprintf` calls, and a second one-line helper for the "could not be edited safely" pair — net deletion of roughly a dozen duplicated lines.

- **[low / confidence: medium]** `cmd/sh-agent/main.go:691,732-733` vs `internal/harness` pattern in `readSettings` (main.go:822-843) — the 1 MiB bound on `settings.json` is enforced twice, two different ways: `attributionSet` does its own `os.Lstat` + `info.Size() > 1<<20` check and refuses with a message, then falls through to call `readSettings`, which does its *own* `os.Stat` of the same path and silently truncates via `io.LimitReader(f, 1<<20)` instead of rejecting. `attributionOff` (used by `status`) reaches `readSettings` directly and never gets the oversized-file message `attributionSet` gives — the same file is judged two different ways depending which command reads it. Shorter version: make `readSettings` itself reject a file over the bound (return the "is larger than 1 MiB" problem string) and delete `attributionSet`'s separate Lstat-and-size branch (lines 730-735) along with the duplicate `1<<20` constant; one bound, one message, checked once per read instead of twice.

- **[low / confidence: low]** `cmd/sh-agent/main.go:557-562` — the comment on `knownCheckout` explaining why a manifest's `EmbeddedPrefix` value is not handed to `samePath` is printed twice, back to back, identically. Harmless (not user-facing text) but is a straightforward accidental duplication; delete one copy.

## Criteria
1. Dead code — pass. Every exported function in `internal/harness` and every unexported helper in both packages has at least one real (non-test) caller; checked by grepping each identifier's call sites. `redactDocument` initially looked uncalled because it is passed by reference (`add("state.json", redactDocument)`), not called by name — confirmed live.
2. Flow tracing — not in this lens
3. Standards conformance — not in this lens
4. Directory layout — not in this lens
5. Simplicity (DIY-except-security) — one finding: the settings.json editor is a correctly-scoped, dependency-free implementation (no library needed for a one-key JSON patch), but its size-bound check is implemented twice with different outcomes (reject vs. silent truncate). No hand-rolled crypto/auth; `crypto/sha256` and `ssh-keygen -Y verify` in `release.go` are used as designed, not DIY.
6. Decoupling — one finding: the attribution/settings-file cluster in `cmd/sh-agent/main.go` is an independent module living in the wrong package, making the entrypoint file a grab-bag of CLI dispatch plus an unrelated file-editing concern. Everything else split cleanly: each `internal/harness` file owns one concern (state, plan, release, report, pathguard) with small interfaces at the seams (`Source`, `File`, `Report`).

## Not examined
- Every `*_test.go` file (out of scope per the dispatch: non-test `.go` files only).
- `go vet ./...` and `gofmt -l` were run and are clean, but no deeper static analysis (staticcheck, golangci-lint) was run — not installed in this sandbox.
- Runtime behavior was not exercised (no `sh-agent install`/`update`/`report` invocation, per the hard limits: no real-`$HOME` install, no network, no `sh-agent report`).
- Criteria 2 (flow tracing), 3 (standards conformance) and 4 (directory layout) were not evaluated at all, per the brief's scope restriction to criteria 1, 5, 6 plus reuse/efficiency.
