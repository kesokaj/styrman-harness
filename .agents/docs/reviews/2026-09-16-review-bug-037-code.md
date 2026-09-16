---
owner: reviewer
date: 2026-09-16
---
# Review — bug-037-code   (2026-09-16)

Range reviewed: `git diff HEAD -- internal/harness/adopt.go internal/harness/harness_test.go` (working tree, unstaged)

## Verdict: APPROVE

Nothing is broken. The fix closes all three manifestations named in BUG-037 and the diagnosis (`.` from inside a skill dir, `..` from a subdirectory of a skill, a relative or absolute path that already names `<repo>/skills/<name>`), including the path reached through the `.agents/skills` symlink and a relative multi-segment `src` evaluated with `cwd == repo`. `within` is correct on equality, parent/child, prefix-name siblings, and trailing separators. The new test fails against the pre-fix code and passes against the fix, confirmed both by static trace and by execution.

## Findings

None that meet the correctness/security/rule bar. Two verified-clean areas worth recording since the brief asked about them explicitly:

- **Relative multi-segment path (diagnosis case 3), no direct unit test** — `TestAdoptRefusesDestructiveSources` covers `.`, `..`, and a self-adopt with an *absolute* `src`; it does not chdir to the repo root and call `Adopt(repo, home, "skills/foo", true)`, the exact reproduction the diagnosis probed (`from` stays relative while `dest` is absolute). I verified this path is in fact closed by writing a throwaway test under `.agents/scratch/bug037-verify/` (deleted after use, nothing under source touched): `Adopt` refused with `"...skills/foo and .../skills/foo overlap..."` and the source survived. The fix is correct here because `filepath.Abs` is applied uniformly regardless of whether `from` arrived relative or absolute, so the existing tests exercise the same code path by a different route. Not a finding — recording it as a coverage note only, since `review-standards` §2 (tests) asks whether a specific named manifestation has its own assertion, and this one doesn't, though the mechanism that fixes it is shared and is tested.
- **Deferred remediation (temp-dir + rename)** — I agree with leaving it out of this diff. It addresses a different failure mode (a `copyTree` that fails halfway loses the old skill because `dest` was already `RemoveAll`'d) rather than the destructive-target bug this diff closes, and folding it in would mix an atomicity refactor into a targeted security fix. It should stay tracked as a named follow-up (BUG-037's `## Fix` section already lists it) rather than be silently dropped.

## Verification performed

- `git show HEAD:internal/harness/adopt.go` against the diff, confirming the pre-fix `dest := filepath.Join(repo, "skills", filepath.Base(from))` is exactly the code the diagnosis traced.
- Traced all three `TestAdoptRefusesDestructiveSources` subtests against the pre-fix body by hand: "dot" would `RemoveAll(repo/skills)` (destroying `skills/alpha`) and mis-report `dest`; "dotdot" would `RemoveAll(repo)` itself; "self-adopt" would `RemoveAll` the source before `copyTree` reads it, so the test's second assertion (`source was removed`) would fail — all three are real regression guards, not tautologies. `newRepo` and `skillDir` build everything under `t.TempDir()`, so no real destruction occurs even if I had executed against the reverted file (I did not revert `adopt.go` myself — sandboxed reasoning plus `git show` was sufficient and avoids any destructive operation on the working tree).
- Ran `go test ./internal/harness/... -run TestAdopt -v`: both `TestAdopt` and `TestAdoptRefusesDestructiveSources` (all three subtests) PASS against the current fix.
- Ran `gofmt -l .` (empty), `go vet ./...` (clean), `go test -race ./...` (both packages `ok`).
- Wrote and ran a throwaway test (`.agents/scratch/bug037-verify/probe_test.go`, deleted after use) exercising the relative multi-segment case from repo root; refused correctly, source intact. Confirms `within` closes this path too.
- Checked `within`'s four required cases by inspection: equality (`rel == "."` → true), parent-vs-descendant (`rel` has no `..` prefix → true), sibling whose name prefixes another (`/skills/foo` vs `/skills/foobar` → `rel = "../foobar"`, correctly `false`), trailing separator (`filepath.Rel` calls `Clean` on both inputs internally → no effect).
- Checked `cmd/sh-agent/main.go`'s adopt usage text and `Adopt`'s doc comment: both still describe current behavior; nothing narrowed needs a text update.
- Checked `repo` is always `filepath.Abs`'d before reaching `Adopt` (`FindRepo`, `repo.go:28,44`), but also confirmed the fix does not rely on that — `dest` is `filepath.Abs`'d unconditionally inside `Adopt`, so a relative `repo` argument (a future caller, or a unit test) is still handled correctly.
- Checked `.agents/docs/bugs/INDEX.md`: BUG-037 (and the bug file itself) are still marked `open`. That is expected — the dispatch scoped this diff to `adopt.go`/`harness_test.go` only; updating the record's status is the `bugs` agent's job in a later step, not part of this code diff, so not a finding here.

## Criteria

1. Dead code — pass. No unused code added; `within` has two live call sites; nothing commented out or left as scaffolding.
2. Flow tracing — pass. `from`/`dest` are made absolute before every comparison and before `RemoveAll`; the refusal is checked before any deletion; `validName` additionally rejects malformed resolved names. All three named manifestations plus the symlink-reached and relative-path variants trace through the same `Abs` + `within` mechanism, verified by execution above.
3. Standards conformance — pass. Errors wrapped/returned per `golang` skill (`fmt.Errorf`, no swallowed errors); the new test is table-free but uses `t.Run` subtests appropriately; `t.Helper()` used in the shared `skillDir` closure.
4. Directory layout — pass. Both changed files stay in `internal/harness/`, test beside source, no new files.
5. Simplicity — pass. stdlib-only (`path/filepath`, `strings`); no new dependency; `within` is a small, single-purpose helper, not speculative abstraction.
6. Decoupling — pass. `within` is package-private, minimal surface, no reach-ins; does not touch `copyTree`'s existing symlink refusal (BUG-036) or any other reader.

## Not examined

- BUG-036 (symlink handling in `RenderContext`/`LoadDocs`/`LoadAgents`) — out of scope for this diff, not re-reviewed.
- `pathguard_test.go` and the rest of the symlink-refusal test suite — not read; the diff doesn't touch that code and the diagnosis calls it orthogonal.
- Sibling-prefix and trailing-separator cases for `within` were checked by inspection and `filepath.Rel`/`Clean` semantics, not by a dedicated executed test (none exists in the diff for those two cases specifically).
- Whole-repo test suite beyond `internal/harness` and `cmd/sh-agent` (both `ok` from `go test -race ./...`); no other packages exist in this module to miss.

## Outcome

- **Note 1, the relative multi-segment case had no dedicated test** — `fixed` at `internal/harness/harness_test.go`, in the new subtest "a relative path into the repo's own skills/": it runs `Adopt(repo, home, "skills/alpha", true)` from the repo root and requires an error with the source still present. Against `git show HEAD:internal/harness/adopt.go` in a scratch copy it fails with "the source was removed"; it passes with the fix.
- **Note 2, copying into a temporary directory and renaming it over `dest`** — `left: the operator's next harness session`. It addresses a different failure, a `--force` replacement losing the old skill when the copy fails halfway, and is named in BUG-037's Fix section and in the handoff's Next Steps.
- **The verifier's coverage note on `within`** (`testing/2026-09-16-verify-bug-037-adopt.md`) — `fixed` at `internal/harness/harness_test.go`, `TestWithin`: equality, parent to child, child to parent, prefix-sharing siblings, a trailing separator and unrelated paths. With `within` replaced by `strings.HasPrefix(child, parent)` in a scratch copy it fails on `/a/skills-old` and `/a/skillsdelta`; it passes in the tree.
