---
owner: verifier
date: 2026-10-01
---
# Verification — google-skills-catalog-final   (2026-10-01)

Range reviewed: b0789f1..working tree

## Status: PASS

Scope: guide §1 commands, §3f catalog coverage, four must-fail controls, scratch-home install, list. Two coverage notes below are narrower-than-asked tests, not failures. Sandbox CLI rows not run (main agent's).
Guide: .agents/docs/testing/guide.md (exists, read via AGENTS.md pointer; not rewritten).

Nothing broken. Notes: (1) refresh script `--dest` refusals of `/`, HOME and the checkout have no case of their own: TestRefreshCatalogScript runs only `--dest ""` and a non-catalog directory with content (the generic "not empty, no CHECKSUMS" guard). The `/`/HOME/checkout `case` lines in refresh-catalog.sh:61-66 are untested by name. (2) The `[a-z0-9-]` refusal is exercised only by a name with a space ("bad name"); no uppercase or other-character case.

## Commands
- 1. pwd /home/kesokaj/workspace/styrman-harness, HEAD b0789f1, date 2026-10-01; change set present.
- 4. `gofmt -l cmd internal` → exit 0, empty output.
- 4. `go vet ./...` → exit 0.
- 4. `go test -race -count=1 ./...` → exit 0: ok root 1.4s, cmd/sh-agent 63.5s, internal/dream 6.7s, internal/harness 8.1s, internal/memory 1.2s, internal/search 7.7s (.agents/scratch/verify-final/go-test-1.txt).
- 4. `go test -count=1 -v ./...` → exit 0; counts: pass 536 (359 top-level), fail 0, skip 3 at b0789f1+dirty for go test -count=1 -v ./...; vs previous: none compared (earlier logs use other commands; not diffed). Skips: TestAgyCorpusLeaves, TestRecall, TestHelperIndex (pre-existing environment-gated skips, not in the diff).
- 7. Added skips in diff, all with reason on the line: catalog_test.go `t.Skip("runs a nested go test")`, `t.Skip("inside the write mode itself")`, `t.Skip("no symlinks here")` x2. No nolint, ts-ignore, .only.
- Must-fail controls (backup in .agents/scratch/verify-final/, restored, sha256 OK, `git diff --stat` identical to before):
  - a. swap changed to `[ -d … ] && rm -rf …` → `go test ./internal/harness -run TestRefreshCatalogScript -count=1` exit 1: "a refresh over a catalog with a missing entry directory did not restore it (exit 1)". FAILED as required.
  - b. root-marker rewrite made a no-op (prefix "./" kept) → `go test ./internal/search -run TestNoIndexMarkerSkipsTheDirectory -count=1` exit 1: "a root marker left files in the index". FAILED as required.
  - c. `.sh-agent-noindex` moved to scratch → `-run TestCatalogFilesAreCarryable` exit 1: "the catalog has no .sh-agent-noindex marker". FAILED as required; moved back (0 bytes, tracked, no status change).
  - d. one byte appended to catalog/agent-platform-alert-configuration/SKILL.md → `-run TestCatalogMatchesChecksums` exit 1: "differs from its CHECKSUMS line: the vendored file was edited". FAILED as required; restored.
- 5. `go run ./cmd/sh-agent install --home .agents/scratch/verify-final-home` → exit 0, 1784 files written. Status (HOME=scratch, installed binary): `binary b0789f1b61f6+dirty` no warning; `drift none — 1783 files match the repo`.
- 5. `diff -rq skills/google-catalog <home>/.claude/skills/google-catalog` → empty (exit 0); same for `.gemini/config/skills/google-catalog` → empty.
- 6. `go run ./cmd/sh-agent list` → exit 0; Skills (35), Subagents (8).

## Coverage of changed code
- skills/google-catalog/scripts/refresh-catalog.sh: TestRefreshCatalogScript (internal/harness/catalog_test.go:381): missing-entry-dir swap (line ~558-565, control a proves it); name rule via "name with a space" case (partial, note 2); `--dest ""` and non-catalog dir (partial for `/`, HOME, checkout, note 1); symlink, symlinked entry, symlinked LICENSE refusals (lines 640-657); manifest `## Excluded` merge (line ~589-598).
- internal/search/reconcile.go: TestNoIndexMarkerSkipsTheDirectory (search_test.go ~647-683), covers marker in a subdirectory, in docs/, and at the root (control b proves root).
- internal/harness/catalog_test.go: marker presence in TestCatalogFilesAreCarryable (control c); TestCatalogWriteModeNeedsCatalogDir (refuses CATALOG_WRITE=1 without CATALOG_DIR; cleanup restores CHECKSUMS, and asserts it unchanged).
- skills/google-catalog/SKILL.md: prose; TestRouterDescription and the install tests cover it.
- Not checked: "marker in any git state" (tracked, untracked, modified) beyond what the single test shows; I did not read each state in it.

## Not examined
Sandbox CLI starts and behavioural rows (main agent); Antigravity runs; `govulncheck`; docs freshness / `/review --docs`; guide §5 real-home steps; each of the other 35 skills' tests individually; whether the marker test covers every git state (tracked/untracked/modified).

## Environment
Linux 6.12 (Debian 13), go from PATH (version not recorded), no Docker used.

## Tree check
Step 8: `git status --short --untracked-files=all` (excluding catalog deletions) compared to the start: the only new path is this log. No source, test or config path changed.
