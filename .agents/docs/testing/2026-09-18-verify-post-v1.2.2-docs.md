---
owner: verifier
date: 2026-09-18
---
# Verification — post-v1.2.2 docs (README + remote.md install chains)

Range reviewed: f3ec5dc..worktree

## Status: PASS

Note: `git status --short --untracked-files=all` at run end also shows
`.agents/docs/testing/2026-09-17-release-binary-e2e.md` modified (an "Added
2026-09-18" section already recording this same v1.2.2 consumer verification,
including the same `test` line exit-0 result). This file was NOT touched by this
verifier run (no Write/Edit issued against it) and was not named in the dispatch's
scope ("README.md and .agents/docs/deploy/remote.md only") - it was already
uncommitted in the tree when this run started. Flagged per the hard rule on
unexpected working-tree changes; nothing here was caused by this run.


Scope: complete for the scope given. Cites .agents/docs/testing/guide.md (already
existed, no changes made to it; guide's §1 Go checks and the general "exercise every
edited chain" principle applied). No `.go` file is in this diff (`git diff HEAD --name-only`
= `.agents/docs/deploy/remote.md`, `README.md` only), so guide §2-§5 (harness self-check,
scratch-home install, report tier, sandbox/real-home CLI starts) do not apply — nothing
under `context/`, `agents/`, or `skills/` changed. Ran guide §1 Go checks anyway per
dispatch instructions, plus a manual exercise (not a guide tier) of the two edited shell
chains against the real published v1.2.2 release.

## Commands

- `gofmt -l cmd internal` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -v -count=1 ./...` → exit 0; output saved to
  `.agents/scratch/verify-post-v122-go/go-test.txt`; all three packages `ok`
  (`github.com/kesokaj/styrman-harness`, `.../cmd/sh-agent`, `.../internal/harness`);
  counts: pass 215 (108 top-level + 107 subtest `--- PASS:` lines), fail 0, skip 0 at
  f3ec5dc for `go test -race -v -count=1 ./...`; vs previous:
  `.agents/docs/testing/2026-09-18-verify-preflight-v1.2.2.md` (pass 215, fail 0, skip 0
  at f91257c+dirty), pass 215→215, skip 0→0 — unchanged, as expected (no Go source in
  this diff).

- `bash -n` on each of the three edited install chains (extracted verbatim to
  `.agents/scratch/verify-post-v122/{readme-gh,readme-browser,remote-md}.sh`):
  all three → exit 0, no syntax errors, no broken `\` continuations.

- Network step (permitted only for this download, per dispatch): in
  `.agents/scratch/verify-post-v122/`,
  `gh release download v1.2.2 --repo kesokaj/styrman-harness --clobber --pattern sh-agent-linux-amd64 --pattern SHA256SUMS --pattern SHA256SUMS.sig --dir .`
  → exit 0. Downloaded `SHA256SUMS` first line: `# styrman-harness v1.2.2`.
  `chmod +x sh-agent-linux-amd64`; `./sh-agent-linux-amd64 version` → `sh-agent v1.2.2`.

- Must-pass direction: exactly the README's `test` line against the real, unmodified
  release files:
  ```
  test "$(./sh-agent-linux-amd64 version)" = "sh-agent $(sed -n 's/^# styrman-harness //p' SHA256SUMS)"
  ```
  → **exit 0** (pass), as expected.

- Must-fail direction: same `test` line against a copy `SHA256SUMS.badversion` whose
  first line was rewritten to `# styrman-harness v1.2.1` (real, unaltered version string
  `sh-agent v1.2.2` on the left, `sed`-extracted `v1.2.1` on the right):
  ```
  test "$(./sh-agent-linux-amd64 version)" = "sh-agent $(sed -n 's/^# styrman-harness //p' SHA256SUMS.badversion)"
  ```
  → **exit 1** (fail), as expected. The `test` line correctly rejects a checksum file
  claiming a different version than the binary reports, confirming the check is neither
  a tautology nor silently vacuous.

- The three `.sh` chain-extraction files: `readme-gh.sh` (README gh route),
  `readme-browser.sh` (README browser route, including the new `mktemp -d` +
  three-file `mv` step), `remote-md.sh` (remote.md § First install on a machine). Read
  side by side with their source fences — each extraction is a verbatim, line-for-line
  copy of the fenced `&&` chain (confirmed by re-reading the `git diff` hunks above
  before extraction).

## Root cause

N/A — no failure found.

## Coverage of changed code

Both changed files are documentation (`README.md`, `.agents/docs/deploy/remote.md`);
no `.go` file in this diff, so "coverage of changed code" (Go test coverage) does not
apply. The behavioural claim in the diff — the new `test` line correctly distinguishes
a matching version string from a mismatched one — is exercised directly above (both
directions), which is the applicable form of coverage for a shell chain in a doc. No
added suppressions, skips, `.only(`, `@ts-ignore`, `//nolint`, etc. found in this diff
(it contains none — grep of the diff hunks shows only comment-text and `&&` chain
edits).

## Not examined

- The `readme-darwin-arm64` and `remote-md` chains were NOT run end-to-end against a
  real `sh-agent-darwin-arm64` or `sh-agent-linux-arm64` binary — only `bash -n`
  syntax-checked, plus the shared `test` line logic exercised via the linux-amd64
  asset (the `test` line text is byte-identical in structure across all three chains
  modulo variable/literal asset name, so the syntax and the comparison logic are the
  same check; the darwin/arm64 binaries themselves were not downloaded or executed).
- The `mktemp -d` + `mv ~/Downloads/...` step at the top of the README browser chain
  was not run for real (would require populating a real `~/Downloads` and was out of
  the network scope the dispatch permitted); `bash -n` confirms it parses, and its
  logic (three-file `mv` into a fresh scratch dir before `cd`) was read, not executed.
- Guide tiers §2 (harness self-check), §3/§3a (scratch-home / plain-copy install),
  §3b (report), §4 (sandbox CLI starts), §5 (real-home smoke checks) were not run:
  none apply, since no `.go`, `context/`, `agents/`, or `skills/` file is in this diff.
- `govulncheck ./...` (an AGENTS.md setup command, not in the guide's numbered tiers)
  was not run — no Go source changed and the dispatch scope named only the guide plus
  the specific shell-chain exercises.
- `/review --docs` was not run — out of verifier scope; not requested by this dispatch.

## Environment

Linux (Debian trixie, 6.12.107+deb13-amd64), repo at f3ec5dc + uncommitted edits to
`README.md` and `.agents/docs/deploy/remote.md` only, tree frozen for the run (no other
files touched; `.agents/scratch/verify-post-v122/` and
`.agents/scratch/verify-post-v122-go/` created, both gitignored). `gh` CLI authenticated
and used only for the one permitted download. Go toolchain per this repo's AGENTS.md
(user-level SDK ahead of `/usr/bin/go`, no `toolchain` directive in `go.mod`).
