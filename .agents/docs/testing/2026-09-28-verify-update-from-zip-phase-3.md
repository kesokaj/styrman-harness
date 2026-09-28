---
owner: verifier
date: 2026-09-28
---
# Verification — update-from-zip phase 3   (2026-09-28)

Range reviewed: 1cb2608..worktree

## Status: PASS

Scope: complete — ran all 8 numbered steps from the dispatch; `.agents/docs/testing/guide.md` exists.

## Commands
1. `gofmt -l cmd internal` → exit 0, no output (clean).
   `go vet ./...` → exit 0, no output.
   `go test -race -count=1 -timeout 180s -v ./...` → exit 0 (bare). Output: `.agents/scratch/verify-update-from-zip-p3/go-test-1.txt`. `ok` for all 3 packages (root, cmd/sh-agent, internal/harness).
   counts (top-level `^--- (PASS|FAIL|SKIP):`): pass 160, fail 0, skip 0 at 1cb2608+dirty for `go test -race -count=1 -timeout 180s -v ./...`; incl. subtests: 291 `--- PASS/FAIL/SKIP:` lines total, all PASS.
   vs previous: `.agents/docs/testing/2026-09-28-verify-update-from-zip-phase-1.md`, pass 139→160 (top-level), skip 0→0. Rise in pass explained: new zip.go/zip_test.go tests added in this range (`git diff 1cb2608 --numstat` shows new files), not a fall/skip-rise — no FAIL condition triggered.
2. `govulncheck ./...` → exit 0, "No vulnerabilities found." (go1.27.1 in use, satisfies AGENTS.md go1.25.8+ requirement).
3. `go test -count=1 -run TestRenderedContextFitsAntigravityBlock -v ./internal/harness` → exit 0. Size line: "Antigravity rules block is 20550 bytes (context 20550, rules 0), budget 22000" — under the 22,000-byte ADR-035 pin.
4. `go run ./cmd/sh-agent install --home .agents/scratch/p3` → 258 written, 0 unchanged, 0 removed, 0 kept (fresh scratch home). `go run ./cmd/sh-agent status --home .agents/scratch/p3` → `drift none — 257 files match the repo` (the 258th managed file is `~/.local/bin/sh-agent`, the binary itself, not part of the drift-comparison file set — consistent with prior runs). Both vendors' rendered context confirmed ending in GANNET:
   - `.agents/scratch/p3/.claude/AGENTS.md` (CLAUDE.md is a one-line `@AGENTS.md` stub) tail: "...Canary: this context ends with the word GANNET."
   - `.agents/scratch/p3/.gemini/config/rules/AGENTS.md` tail: same, ends "...Canary: this context ends with the word GANNET."
   Never touched the real $HOME (used `--home .agents/scratch/p3` throughout).
5. `go run ./cmd/sh-agent list` → "Skills (27)" and "Subagents (7)" confirmed by grep. `go run ./cmd/sh-agent --help` → shows `update --from-zip FILE [--dry-run]  zip place: unpack a repository zip you downloaded, install its harness and the running binary (nothing in the zip replaces the binary), delete the unpacked copy`.
6. `git status --porcelain --ignored --untracked-files=all -- context skills agents commands rules` → only ` M skills/VENDORED.md` (a tracked, intentional modification per the change set — not untracked or ignored). No untracked or ignored files under these paths. `go test -count=1 -run TestEmbedHoldsOnlyTrackedFiles ./...` → PASS (root package); other two packages report "no tests to run" (expected, test lives only in root package).
7. Network-safety check of `cmd/sh-agent/zip_test.go` (15 top-level tests calling `run(...)`):
   - `TestUpdateFromZip`, `TestUpdateFromZipRefusals`, `TestUpdateFromZipDryRunWritesNothing`, `TestStatusAfterZip`, `TestKnownCheckoutZip`, `TestUpdateFromZipRefusedZipWritesNothing`, `TestFromZipEmptyValueRefused`, `TestDiffAfterZip`, `TestUpdateFromZipSweepsStaleDirs`, `TestUpdateFromZipRefusedContentWritesNothing`, `TestUpdateFromZipSweepWithGlobCharsInHome`, `TestUpdateFromZipDryRunCleansTempDir` — all call `isolate(t)` (line 72: `t.Setenv("PATH", t.TempDir())`), which sets `PATH` to an empty temp directory. Network-safe.
   - `TestFromZipRefusedOnOtherCommands` — does not call `isolate(t)` but directly does `t.Setenv("PATH", t.TempDir())` at line 262 (same effect, empty PATH). Network-safe.
   - `TestReleaseUpdateDryRunChangesNothing` — does not call `isolate(t)`; instead sets `PATH` (line 374) to a temp dir (`bin`) containing only a hand-written stub `gh` script (lines 368-374) that logs its args and prints a fake version, never touching the network. No real `git`/`gh` reachable. Network-safe by construction.
   - `TestCheckoutUpdateDryRunChangesNothing` (line 468) — does **NOT** call `isolate(t)` and does **NOT** restrict `PATH`; it uses the real `git` binary from the ambient PATH (via `exec.Command` in the test body, lines 480-493, and inside `run([]string{"update", "--dry-run", "--repo", repo, ...}, ...)`). This is the one test in the file that can reach a real `git`/`gh` binary. However every git operation it performs targets only `t.TempDir()`-local paths: `gitFixture` (main_test.go:704) creates `origin` as a local `--bare` repo under `t.TempDir()` and clones/pushes only between local paths, with `GIT_CONFIG_GLOBAL`/`GIT_CONFIG_NOSYSTEM` set to prevent picking up the operator's real git config; no remote URL, `http(s)://`, or SSH host ever appears. The command under test is invoked with `--repo <local path>`, which is the "checkout" mode, not the release/`gh` path, so `gh` (network-capable) is never invoked either. Conclusion: reachable git binary, but no reachable network — root cause of safety is the all-local fixture, not a PATH restriction. Named here per the dispatch's request to name any test that "can" reach git/gh; this one technically can reach the `git` binary but not any network endpoint.
8. `git diff --numstat 1cb2608` (tracked-file changes only; untracked new files are not part of a plain `git diff`) → all rows numeric (no `-`/`-` binary markers): DECISIONS.md, bugs/INDEX.md, deploy/local.md, deploy/remote.md, delivery.md, skill-matrix.md, .gitignore, AGENTS.md, README.md, REQUIREMENTS.md, ROADMAP.md, cmd/sh-agent/main.go, cmd/sh-agent/main_test.go, internal/harness/apply.go, internal/harness/plan.go, internal/harness/repo.go, internal/harness/source.go, skills/VENDORED.md.

## Suppressions added
`git diff 1cb2608 -- '*.go'` grepped for `nolint`, `lint:ignore`, `t.Skip`, `go:build ignore`, `.only(`, `.skip(` in added (`+`) lines → none found. No suppressions or skips added in this range.

## Coverage of changed code
- `internal/harness/zip.go` (new) — covered by `internal/harness/zip_test.go` (new), which passed (all subtests under `TestUnpackZip*`, `TestSecretShapedCLITokenFiles`, etc.).
- `cmd/sh-agent/main.go` (modified, +170/-8) — covered by `cmd/sh-agent/zip_test.go` (new, 15 tests) plus existing `main_test.go` (all passing).
- `internal/harness/apply.go`, `internal/harness/plan.go`, `internal/harness/repo.go`, `internal/harness/source.go` (modified) — exercised by the full `go test -race -count=1 ./...` run (all PASS) and by `TestHarnessMap`/`TestEmbedHoldsOnlyTrackedFiles` for the plan/embed changes; not individually diffed line-by-line for a 1:1 mapping beyond the aggregate pass.
- Doc-only changes (DECISIONS.md, bugs/INDEX.md, deploy/*.md, delivery.md, skill-matrix.md, AGENTS.md, README.md, REQUIREMENTS.md, ROADMAP.md, VENDORED.md) — no test coverage applicable; reviewed separately by `/review --docs` runs already on disk under `.agents/docs/reviews/` (not re-examined here, out of scope for this verifier phase).

## Not examined
- The content/correctness of the `.agents/docs/reviews/*` and `.agents/docs/bugs/BUG-093*`, `BUG-094*` files themselves (reviewer artifacts; out of scope per dispatch — "reviewers run in parallel... ignore those").
- No frontend/`web`/`ui` directory exists in this repo — not applicable.
- No Dockerfile changed in this range — container build/smoke test not run (not applicable per the dispatch's own scope, which did not ask for it).
- Full behavioural/manual CLI smoke run (starting `claude`/`gemini` interactively) per `testing/guide.md`'s tier-4 sandbox procedure was not run in this phase (not in the dispatch's numbered list); prior phases (phase-1, phase-2) covered targeted `run(...)`-level test execution.
- Did not independently line-by-line diff `internal/harness/apply.go`, `plan.go`, `repo.go`, `source.go` against their existing `_test.go` coverage beyond confirming the aggregate suite passes; a stricter per-hunk coverage audit was not performed.

## Environment
OS: Linux 6.12.107+deb13-amd64 (Debian trixie). Go version: go1.27.1 linux/amd64 (user-level SDK ahead of /usr/bin/go, per AGENTS.md). `govulncheck ./...`: "No vulnerabilities found." Colima/Docker: not applicable (no Dockerfile change in range).
