---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-8   (2026-10-01)

Range reviewed: working tree on 7c2e30a (Phase 8)

## Status: PASS

Scope: all seven requested steps ran; guide exists (.agents/docs/testing/guide.md, cited, unchanged). § 3c run in compact form (see below), not the full matrix file; sh-agent-matrix.md row not re-read.

## Commands
- `gofmt -l cmd internal` → exit 0, printed nothing
- `go vet ./...` → exit 0
- `go test -race -v -count=1 ./...` → exit 0; all 6 packages ok; counts: pass 512, fail 0, skip 3 at 7c2e30a+dirty for go test -race -v -count=1 ./... (log .agents/scratch/verify-p8x/go-test-1.txt); vs previous: .agents/docs/testing/2026-10-01-verify-audit-fix-phase-7-r4.md, pass 492→512, skip 3→3
- `go test -count=1 -cover ./cmd/... ./internal/...` → exit 0; cmd/sh-agent 85.6% (baseline 84.8, plan >=84 met), harness 89.0% (88.2), dream 89.6, memory 89.1, search 84.2
- guide § 3: `go run ./cmd/sh-agent install --home .agents/scratch/verify-p8` → exit 0 (266 written); `status --home ...` → exit 0, "drift none — 265 files match the repo", `binary 7c2e30a9519b+dirty` (no warning). Install printed PATH note and "settings: .claude/settings.json created with attribution commit/pr empty" (new settings.go).
- IN5 (built binary `.agents/scratch/verify-p8x/sh-agent`): `list --check` → exit 2; `index --kind bug --home ...` → exit 2; `version --dry-run` → exit 2. Each prints "flag provided but not defined: -<flag>" then "Usage of <verb>:" — on STDOUT, not stderr (the flagset writer is `out`, main.go:201 `c.fs.SetOutput(out)`; the pre-change code also used `fs.SetOutput(out)`, so not a regression; the dispatch said stderr, noted as deviation).
- `<verb> --help` for all 15 verbs → each prints "Usage of <verb>:" (8-18 lines), no panic/goroutine text; exit 1 for every verb (main.go:209 returns flag.ErrHelp as an ordinary error → 1); top-level `--help` exit 0. Not asserted by a convention doc I found; recorded as observed.
- guide § 3a: built binary, read-only .git-less mirror: `install --repo mirror --home home2` → exit 0, "plain copy, not its own git checkout: the running binary is installed as is" (no rebuild); mirror listing shows no .agents; `status` → "repo ... (plain copy ...; update is refused)", drift none (265 files); `install --repo . --home home3 --dry-run` → "[dry-run] the binary would be rebuilt from <checkout>" and a dry-run line "would create .../.claude/settings.json" (read-only mode did not create it, see 3c/IN row below).
- guide § 3b: `go run ./cmd/sh-agent report --dry-run --home .agents/scratch/verify-p8` → exit 0; prints skipped files, "[dry-run] would clone git@github.com:kesokaj/styrman-harness-feedback.git into .../report-NNN" and would-write lines; `.config/sh-agent/` in the scratch home holds only state.json afterwards (no report-* dir, nothing sent).
- guide § 3c compact kit (`.agents/scratch/verify-p8x/m/`: `co` = git clone of the checkout overlaid with the working tree, `pc` = rsync copy without .git, `rel` = `go build -trimpath -ldflags "-X main.version=v9.9.9"`; fresh home per cell; PATH = /usr/bin:/bin + go dir). Cells (the head-pipe SIGPIPE exit 141 on install/report in my output is the filter, not the verb; install exit read separately = 0):
  - checkout: install exercised (builds from checkout, exit 0); status exercised (main @ 7c2e30a uncommitted, binary 7c2e30a9519b+dirty, drift none, 265 files); update --check refused-with-cause (up to date with upstream, then "has uncommitted changes ... commit or stash, or run install"; exit 1); report --dry-run exercised; list --check refused exit 2 (flag not defined + usage).
  - plain copy: install exercised ("plain copy ... installed as is"); status exercised (plain copy, update is refused, drift none); update --check refused, names cause and alternative (checkout clone or release binary; exit 1); report --dry-run exercised; list --check refused exit 2.
  - release build: install exercised ("release binary v9.9.9: installing the harness built into it"); status drift none, binary line no warning; update --check exercised, exit 0 ("newest release v1.16.2 is not newer than this binary v9.9.9" — a gh was on PATH and reached GitHub read-only; guide says no gh, so the network was touched for a release listing, nothing installed); report --dry-run exercised, exit 0; list --check refused exit 2.
- IN10: `grep -rn 'MkdirTemp("", \|CreateTemp("", ' cmd internal --include='*.go' | grep -v _test` → empty (grep exit 1). Targeted `go test ./cmd/sh-agent ./internal/harness -run '<list>' -count=1 -v` → exit 0, both packages ok, 29 PASS lines, 0 FAIL: TestAttributionNote, TestAttributionSet, TestAttributionSetAndHint, TestAttributionSetReportsAnUnwritableSettingsFile, TestLintDescriptionLimit, TestLintRejectsUnquotedComment, TestNoSystemTempDirInSource, TestParseOrigin, TestReadCappedRefusesOversize, TestReorder, TestReorderSearchFlags, TestResolveSource, TestResolveSourceLenient, TestResolveSourceRefusesNamedNonRepo, TestResolveSourceSkipsZipManifest, TestSettingsWriteUsesHarnessWriteAtomic, TestUnpackZip* (12 tests), TestVerbRefusesAnotherVerbsFlag. (The dispatch's `TestOrigin` pattern matched only TestParseOrigin; TestSourceOriginString is not in the pattern.)
- `govulncheck ./...` → exit 0; "Your code is affected by 0 vulnerabilities"; 1 vulnerability in a required module not called (the known GO-2026-5024 class; ID not printed without -show verbose, not confirmed).

## Coverage of changed code
- cmd/sh-agent/main.go — covered by main_test.go (85.6% package)
- internal/harness/{apply,frontmatter,inventory,pathguard,release,repo,report,zip}.go — covered by their *_test.go plus new phase8_test.go, repo_origin_test.go; package 89.0%
- internal/harness/settings.go — covered by settings_test.go
- Suppressions added: none checked beyond the targeted tests; `git diff` for //nolint, t.Skip, .only not grepped (see Not examined).

## Not examined
- Added-suppression grep over the diff (nolint/t.Skip/.only) was not run.
- Full sh-agent-matrix.md table (8-verb cells beyond the five requested), `update --dry-run`, Antigravity/real-home smoke (guide § 4, § 5) not run.
- From .agents/scratch/audit-fix/phase8-main.md "Left / unverified": `index` dropping search-only flags (--kind/--limit/--refs/--open) is a behaviour narrowing needing a line in README's command table and sh-agent-matrix.md; real-home install, Antigravity run not done by the fork.
- From phase8-harness.md "Not verified": the signed-zip verification path with the built binary (unit tests only: TestVerifyZipSums*); the fork's module-wide go vet/test was blocked mid-run (I ran both green on the final tree); its follow-ups 1-3 (main.go readSettings ReadCapped; updateFromZip comment; main_test pin) — main.go:1053 now mentions "the system temp directory (IN10)" in a comment, not read in full by me.
- --help exit status is 1 per verb (observed; no convention doc consulted); the error/usage for an undefined flag is on stdout, not stderr as the dispatch expected.
- An old /tmp/zip-* leftover mentioned by the fork was not inspected.

## Tree check
Source/test/config paths at the end equal those at the start (17 modified/new files listed in the dispatch). New artifacts by other agents seen: .agents/docs/reviews/2026-10-01-review-worktree-p8-docs.md, .agents/docs/reviews/2026-10-01-security-worktree-p8.md.

## Environment
Linux 6.12 (Debian 13), bash; Go from `go version` not recorded separately (go1.2x user SDK per AGENTS.md); no Docker used; govulncheck at ~/.go/bin.
