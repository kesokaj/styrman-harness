---
owner: reviewer
date: 2026-09-16
---
# Review — bug-g1-attribution-note   (2026-09-16)

Range reviewed: working tree, three named paths (`cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go`, `context/AGENTS.md`), uncommitted. Other working-tree edits present in the repo (agents/*.md, skills/*.md, etc.) were ignored per the dispatch.

Blobs read: `main.go` 3a7663e894a36bcca0dd568be8df3786254d762a, `main_test.go` 22e56970b0d2e2ff122c82fe3548e023b610d2c4, `context/AGENTS.md` 2fce8698f7defe085483a2a81a32ccace9e366e2 (all match the diff's post-image hashes).

## Verdict: APPROVE WITH NITS

The feature is correctly implemented and the message text matches between `install`, `status`, and the new `context/AGENTS.md` § Git sentence; `gofmt`, `go vet`, and `go test -race ./cmd/sh-agent/` all pass. The gaps are two untested branches in `attributionHint`'s own error handling, both plausible in the field and both invisible to `TestAttributionNote` as written, plus a flow question about whether an unrelated-file read error should be allowed to fail the whole `install`/`status` command.

## Findings

- **[medium / confidence: high]** `cmd/sh-agent/main_test.go:360-397` — `TestAttributionNote` never passes `--vendor agy` (or `antigravity`). The gating `if vendor == harness.VendorAll || vendor == harness.VendorClaude` at `main.go:165` and `main.go:319` is exactly the behaviour the dispatch asked to be verified, and it is real, working code (confirmed manually: `install --vendor agy` prints no attribution line) — but no test exercises it. If that condition were dropped or inverted, every existing test case would still pass. Add one more `run()` call with `--vendor agy` asserting the note is absent.

- **[medium / confidence: medium]** `cmd/sh-agent/main.go:494-497` — the malformed-JSON branch (`json.Unmarshal` error → print "...is not valid JSON..." and return `nil`) has zero test coverage; none of the four `TestAttributionNote` cases write invalid JSON to `settings.json`. This is a real user-facing path (a hand-edited `settings.json` with a syntax error is plausible) and it is also the one branch where a regression has the worst effect: if `return nil` were ever changed to `return err`, `sh-agent install` — the tool's core command — would hard-fail for any operator whose `settings.json` has a syntax error, rather than just skip the advisory note. Add a fifth case with e.g. `{"attribution": {` and assert the "not valid JSON" text (and, separately, that `install`/`status` still exit 0 and still finish their other output).

- **[low / confidence: medium]** `cmd/sh-agent/main.go:484-487` — any `os.ReadFile` error other than `fs.ErrNotExist` (e.g. a permission-denied `settings.json`) is propagated up and aborts `install`/`status` entirely, before the drift report or the rest of `status`'s output is printed. `settings.json` is explicitly "the operator's file" per the function's own doc comment (never written by sh-agent) and the note it feeds is advisory; failing the whole command over an unreadable advisory file is a harsher failure mode than the feature seems to intend, and it is untested either way. Worth a decision (log-and-continue vs. current behaviour) even if the current choice is kept deliberately.

## Criteria
1. Dead code — pass: no unused symbols, imports, or leftover debug output introduced.
2. Flow tracing — mostly pass, with the low-severity finding above: `attributionHint`'s non-ENOENT read error propagates through `status`/`install` and short-circuits their remaining output; every other path (missing file, empty/false, non-empty/true, missing key) is correctly wired to both call sites and to the printed text.
3. Standards conformance — pass: Go error wrapping isn't needed here (`os.ReadFile`'s error is returned bare, which is consistent with sibling code in the same file, e.g. `resolveVendor`/`selfExecutable`); doc comment above `attributionHint` accurately describes intent; the new context sentence and the two code messages agree word-for-word ("attribution": {"commit": ""}).
4. Directory layout — pass: change stays inside existing files, no new files.
5. Simplicity — pass: stdlib only (`encoding/json`, `io/fs`), no new dependency; the `any`-typed switch for string/bool is the minimum needed to be safe against an unknown or evolving schema for `attribution.commit`, though its bool branch was not directly checked against Claude Code's documented schema (only the string form `""` is referenced in the context sentence and observed in test fixtures).
6. Decoupling — pass: `attributionHint` is a single, self-contained helper called identically from both commands; no cross-package reach-ins.

Staleness sweep (README.md, `.agents/docs/deploy/*.md`, `internal/harness/*.go` doc comments, AGENTS.md) — no stale text found: README's one-line summary of `sh-agent status`'s output was already non-exhaustive before this change (it omits the pre-existing foreign-CLAUDE.md note too), so it is not contradicted, just still incomplete by design.

## Not examined
- Whether Claude Code's actual `settings.json` schema documents `attribution.commit` as ever being a JSON boolean (only the string form is evidenced in this repo); could not verify against upstream docs from this sandbox.
- The rest of the working tree's uncommitted changes (agents/*.md, skills/*.md, cmd/sh-agent/main.go's other unrelated hunks if any, etc.) — out of scope per the dispatch.
- Runtime behavior on macOS/Antigravity paths — only Linux, this box, was used for the manual `--vendor agy` check.

## Outcome

- Finding 1 (medium) — fixed: `TestAttributionNote` gained the case `antigravity only never notes` (`--vendor agy` on install and status); shown to fail with the vendor gate removed from status (`--- FAIL: TestAttributionNote/antigravity_only_never_notes`), passes with it.
- Finding 2 (medium) — fixed: the case `malformed json is advisory` asserts the "is not valid JSON" note and a zero exit from both commands.
- Finding 3 (low) — decided and fixed: settings.json is advisory — sh-agent never writes it and no state of it may fail install or status — so a read error other than not-exist now prints "could not be read (<err>)" and returns nil; the case `unreadable file is advisory` (settings.json as a directory) shown to fail with `return err` restored, passes with the note. The decision is the comment above `attributionHint`.
- Round 1 closed 2026-09-16; no second round owed (coverage and one robustness decision, no mechanism finding). `gofmt -l .` empty, `go vet ./...` and `go test -race ./...` exit 0 after the changes.
