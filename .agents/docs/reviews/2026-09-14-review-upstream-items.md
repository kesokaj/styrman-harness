---
owner: reviewer
date: 2026-09-14
---
# Review — upstream-items   (2026-09-14)

Range reviewed: 8046e4d..worktree

HEAD pinned: 8046e4d5f24b009bd2461ebef28931cade7a4acd
Blobs read (`git hash-object`):
- cmd/sh-agent/main_test.go   4676f8bd0f4dc4a04a210ce15dfe87540463c65b
- internal/harness/inventory.go   5a7c641fdd06263f6712110025cf31504b38370b
- internal/harness/apply.go   c1b2a55f2f2efcb8024b8efe2b9c961d0d50c3cc
- context/dev-setup.md   a14ab2f4152e8a823aa5f3950dcbbedbb7bfb765
- internal/harness/pathguard_test.go (untracked)   4397427fa4f4a387bdda65f22ed0e5d29763e453

`go build ./...`, `go vet ./...`, `go test ./cmd/sh-agent/... ./internal/harness/...` all ran bare and green (not piped into a commit chain).

## Verdict: APPROVE WITH NITS

Nothing broken. `validName` correctly closes the path-traversal gap (it needs the explicit backslash check — `filepath.IsLocal` alone does not catch a literal `\` on Linux, since Linux does not treat it as a separator), the `apply.go` resync guard is a no-op for every real manifest entry the installer itself ever writes, and `checkoutShaped` genuinely decouples the dry-run test from being run inside a downstream mirror. The one real nit is a self-contradiction risk in the `dev-setup.md` edit, worth a second look before it propagates to every project via `sh-agent install`.

## Findings

- **[low / confidence: medium]** `context/dev-setup.md:8` — the edit adds "AND that project's `.claude/settings.json` records the standing permission" as a precondition for running a root step, but the sentence still ends with the unchanged clause "Project `.claude/settings.json` records the standing permission." The same fact now appears twice in one bullet — once as a gating condition, once as a trailing restatement — which reads as accidental duplication rather than a deliberate two-part rule. Since this line is the single source deployed to every project's global context, worth confirming the intent (tightening: verbal instruction alone is no longer sufficient, the file must already grant it) was meant, and then removing the now-redundant trailing sentence, rather than shipping both copies.
- **[informational / confidence: high]** `internal/harness/inventory.go:114,155` — `validName` is applied to the fallback names too (`e.Name()` for a skill directory, `filepath.Base(p)` trimmed of `.md` for an agent file). For the skill fallback, `e.Name()` comes from `os.ReadDir` and can never be `""`, `"."`, `".."`, or contain a separator — the guard is unreachable there in practice, which is fine (defense in depth, not a defect) but confirms it is not doing any real work for that branch. For the agent fallback, a file literally named `"...md"` on disk would trim to `".."` and DOES fail the guard — an intentional, if unusual, refusal rather than a gap.

## Not examined

- `internal/harness/plan.go`, `derive.go` in full (only grepped for the `filepath.Join(ClaudeRoot|AgyRoot, …)` call sites feeding `apply.go`'s manifest).
- The Antigravity (`agy`) install path end-to-end; only the Claude-side manifest entries were traced.
- `cmd/sh-agent/main_test.go`'s other tests (`TestKnownCheckout`, `TestPlainCopyNeverRebuilds`) beyond confirming they still compile and pass — not line-by-line re-verified since they are unchanged by this diff.
- No attempt to run the real `sh-agent install --home <scratch>` end-to-end smoke check from `AGENTS.md`'s Setup commands; relied on `go test`/`go vet`/`go build` only, given the diff is unit-test and guard-focused.

## Criteria

1. Dead code — pass. No orphaned code; the new `validName` helper and its test file are both used.
2. Flow tracing — pass. `validName` is wired into both `LoadSkills` and `LoadAgents` for both the parsed and fallback names; the `apply.go` guard sits exactly on the untrusted-input boundary (a manifest file under the user's home) and does not change behaviour for any path the installer itself constructs, since every `File.Rel` is `filepath.Join(ClaudeRoot|AgyRoot, name, …)` with `name` already passed through `validName`. `checkoutShaped`/`IsCheckout`/`Inspect` interaction verified directly: `IsCheckout` only `Lstat`s for a `.git` entry (empty dir satisfies it), and `Inspect` swallows the resulting `git` command failures into empty `RepoInfo` fields rather than erroring, exactly as the comment claims — and `TestDryRunNeverBuilds`'s call path (`knownCheckout` short-circuits on an explicit `--repo`) never touches git validity at all, so the test is genuinely independent of where the suite itself is checked out.
3. Standards conformance — pass. Idiomatic use of `filepath.IsLocal` (Go 1.20+) plus the necessary explicit separator check; error messages name the offending file per the project's own convention (`AGENTS.md`: "refused with the file named").
4. Directory layout — pass. New test file `internal/harness/pathguard_test.go` sits beside the code it covers (`internal/harness/`), matching the project's "tests beside the code" convention; no stray files.
5. Simplicity — pass. Pure stdlib (`path/filepath`, `strings`); no new dependency; the guard is minimal and single-purpose.
6. Decoupling — pass. The guard is local to the package boundary where untrusted names/paths enter (`inventory.go` parse boundary, `apply.go` manifest-read boundary); no reach-ins.

## Outcome

- **1 [low/medium] redundant trailing restatement — FIXED.** The sentence
  now ends "Look at the target before deleting or overwriting it. Without
  that record, hand the step back with the exact command instead." — the
  precondition once, and what to do when it is absent. Re-verified:
  sandbox install rc 0, `status` drift none, the rendered
  `.claude/AGENTS.md` carries the new wording.
- **2 [informational] `validName` on the skill fallback is unreachable —
  ACCEPTED as is.** Kept: one guard at one boundary for both name sources
  is simpler than a second code path that trusts `ReadDir`; it costs
  nothing.
