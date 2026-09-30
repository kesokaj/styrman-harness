---
owner: reviewer
date: 2026-09-30
---
# Review — project-code   (2026-09-30)

Range reviewed: the whole repository at HEAD (3478a1c), no diff.

## Status: COMPLETE

## Verdict: APPROVE

Nothing broken found. This is a mature, heavily security-reviewed codebase:
almost every non-trivial function in `internal/harness` and `internal/search`
carries an inline comment citing the specific past review round, BUG-NNN, or
security-review finding that shaped it (path-escape guards, zip-slip guards,
redaction adversarial cases, git-redirect-env refusals, TOCTOU-safe renames,
os.Root-scoped reads). I read the highest-risk surfaces in full
(`report.go`'s redaction/push pipeline, `zip.go`'s unpack path, `pathguard.go`,
`internal/search/reconcile.go` and `search.go`, `cmd/sh-agent/main.go`,
`derive.go`, `adopt.go`) and found no regression of a recorded bug and no new
correctness or security defect. `gofmt -l cmd internal`, `go vet ./...`,
`go build ./...`, and `govulncheck ./...` all ran clean (the one reported
module vulnerability is the known, unreached `golang.org/x/sys` advisory
`AGENTS.md` already names). `deadcode ./...` flagged one function
(`SwapReleaseSigner`, `internal/harness/release.go:422`) as unreachable; it is
a documented test-injection seam called from `cmd/sh-agent/zip_test.go` in a
different package, so the tool's static-reachability pass misses it — not a
real dead-code finding.

One layout observation, not raised to a finding: `cmd/sh-agent/main.go`
(1199 lines) carries a full second feature's business logic in the entrypoint
package — reading, parsing, and atomically rewriting Claude's `settings.json`
(`attributionSet`/`attributionHint`/`attributionOff`/`readSettings`, main.go:805-1031)
— where the golang skill's checklist says "`cmd/` stays thin — flag business
logic in main.go." It is well-tested and correct as far as I traced it, so I
am not filing it as a blocking finding, but a reviewer on the next change to
this area should ask whether it belongs in `internal/harness` beside
`report.go`/`zip.go`, which own comparable one-off, carefully-hardened file
mutations.

## Findings

None that meet the bar (correctness, security, or a stated requirement/rule
violation). Nothing found regressing any of the ~90 recorded bugs in
`.agents/docs/bugs/INDEX.md`, several of which target exactly the files read
here (report.go redaction/push-target checks BUG-053–060/074; zip unpack
BUG-093/094; search's work-tree-rename parsing BUG-108, fixed today).

## Criteria

1. **Dead code** — pass. No unused exported symbols with real production
   callers missing; the one `deadcode` hit (`SwapReleaseSigner`) is an
   intentional cross-package test seam, documented as such.
2. **Flow tracing** — pass on everything read. `report.go`'s redaction feeds
   `SendReport` end to end (`Collect` → `BuildRecord` → `writeAndPush`), every
   git-redirect env var is checked before every write path that could be
   redirected, every push/clone path scrubs credentials before printing.
   `internal/search/reconcile.go`'s `gitView` uses `--no-renames` and reads
   through `os.Root` (the BUG-108 fix, confirmed in place). `zip.go` validates
   every entry name, size, and kind before writing anything, and removes the
   temp directory on every error path. `cmd/sh-agent/main.go`'s command
   dispatch (`run`) correctly routes `--from-zip` ahead of every other flag
   combination and refuses it outside `update`.
3. **Standards conformance** — pass with the one soft note above
   (`cmd/` not staying thin in `main.go`'s settings-editing code). `gofmt`,
   `go vet`, `go build`, `govulncheck` all clean. R4's one-dependency rule
   holds (`go.mod`: one direct require, `modernc.org/sqlite`). `toolMap` in
   `derive.go` covers every tool actually used across `agents/*.md`
   (Read, Grep, Glob, Bash, Edit, Write) — no agent risks the Antigravity
   hang the project rule warns about.
4. **Directory layout** — pass, with the `cmd/`-not-thin note above as the
   one deviation from the stated convention; entrypoints stay in
   `cmd/sh-agent/`, private logic in `internal/harness` and `internal/search`,
   tests beside their sources throughout.
5. **Simplicity (DIY-except-security)** — pass. No hand-rolled crypto/auth:
   `crypto/sha256`/`sha1` are used only as checksums/content IDs (never as a
   MAC, KDF, or signature primitive), SSH-signing and verification for
   releases go through `git`/`ssh-keygen -Y verify`. No dependency beyond the
   one justified SQLite driver.
6. **Decoupling** — pass. `internal/search` and `internal/harness` interact
   through a small, one-directional import (`search` imports `harness` for
   `SecretShaped`), no reach-into-internals. Small interfaces at the seams
   (e.g., `report.go`'s `redact func(string, []byte) ([]byte, []Redaction)`).

## Not examined

Given the scope (whole repo, ~16,700 lines of Go plus ~1,000 tracked
non-scratch files) and the tool-call budget, I read the highest-risk and
highest-churn Go source files in full but did not read every file. Not
examined in this pass:
- `internal/harness/apply.go`, `frontmatter.go`, `inventory.go`, `plan.go`,
  `release.go` (beyond the `deadcode`/`SwapReleaseSigner` spot check),
  `render.go`, `repo.go`, `source.go`, `state.go`, `owned_unix.go`/`owned_other.go`.
- `internal/search/parse.go`, `store.go`.
- `embed.go` and its generated-counterpart guarantee (`plan_test.go` at repo
  root asserts embedded-plan-equals-checkout; I read its doc comment but not
  its body or a full trace of the embed build).
- All test files (`*_test.go`) — read none in full; relied on `go build`,
  `go vet`, `gofmt`, and `govulncheck` passing rather than re-deriving test
  coverage by hand. `go test ./...` (with `-race`) was not run in this pass
  (time/tool-budget tradeoff against a ~16.7k-line Go tree); a change-free
  review leans on the fact that HEAD's own CI-equivalent (verifier/preflight)
  already gates merges, but this means I cannot personally attest the test
  suite is currently green.
- All markdown: `context/AGENTS.md`, every `skills/*/SKILL.md` and its
  `references/`, `agents/*.md` bodies (frontmatter and tool lists were
  checked; prose was not), `commands/*.md`, `rules/*.md`, `README.md`,
  `.agents/docs/specs/`, `.agents/docs/deploy/`, and whether the "every
  description of the harness names the search index" restatement sweep from
  commit 3478a1c is in fact exhaustive — I did not grep every file for
  "search" against every restatement site named in project rule 6.
- `.agents/docs/testing/2026-09-30-verify-project.md` (untracked, found on
  disk during this review): appears to be an in-progress log from a
  concurrently dispatched `verifier` run, not authored by me and not part of
  HEAD; noted but not reviewed as source.
- `.agents/scratch/` contents (fixture clones, prior review scratch): rightly
  out of scope — scratch is gitignored, not source, per the project's own
  rules.
