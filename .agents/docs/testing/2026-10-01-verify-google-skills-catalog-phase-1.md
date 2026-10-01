---
owner: verifier
date: 2026-10-01
---
# Verification — google-skills-catalog-phase-1   (2026-10-01)

Range reviewed: 754a7ce..working tree (Phase 1)

## Status: FAIL

Scope: steps 1-9 all ran; guide exists. One FAIL line (refresh script relative --work path), otherwise every check green.

## Commands
- 1. `gofmt -l cmd internal` -> exit 0, prints nothing. Start non-catalog status: 11 lines (.agents/scratch/verify-start-status.txt); unstaged diffstat in .agents/scratch/verify-start-diffstat.txt
- 2. `go vet ./...` -> exit 0
- 3. `go test -race -count=1 ./...` -> exit 0 all 6 pkgs ok; wall 1m1.8s (cmd/sh-agent 61.3s, harness 6.2s, search 7.3s, dream 6.5s, memory 1.2s, root 1.4s)
- 3. `go test ./... -v -count=1` (to .agents/scratch/verify-go-test-1.txt) -> exit 0; top-level `^--- PASS` 355; all `--- PASS` incl subtests 532; fail 0; skip 3 (see below)
- 4. named harness tests: 10 PASS (TestCatalogMatchesChecksums, IndexMatchesEntries, FilesAreCarryable, InstallsClean, RouterDescription, RefreshCatalogScript 2.65s, NoTmpInHarnessProse, VendoredRowsCarryAPin, RestatementGrepSkipsVendoredCatalog, SkillLinksResolve), no SKIP; ok
- 5. TestEmbedHoldsOnlyTrackedFiles PASS; TestRenderedContext: Antigravity block 15155 bytes (budget 22000); hand-written 14252 (target 15000)
- 6. catalog entries 141; INDEX entries 141; CHECKSUMS lines 780 (785 staged files minus INDEX, MANIFEST, CHECKSUMS, 2 LICENSE = 780); executables 0; symlinks 0
- 7. Must-fail controls (each compiled, each failed the named test, restored by cp/cmp; chmod back to 644):
  - (a) +1 byte spanner-basics/SKILL.md -> FAIL TestCatalogMatchesChecksums: "spanner-basics/SKILL.md differs from its CHECKSUMS line: the vendored file was edited"
  - (b) INDEX line removed + ghost added -> FAIL TestCatalogIndexMatchesEntries, both reported: "INDEX.md lacks the line ..." and "INDEX.md has a line no entry produces: - `ghost` — x"
  - (c) chmod 755 -> FAIL TestCatalogFilesAreCarryable: "has an executable bit"
  - (d) `[a-z0-9/-]+` -> `[a-z0-9-]+` -> FAIL TestRestatementGrepSkipsVendoredCatalog: catalog hit skills/x/catalog/e/SKILL.md in output
  - (e) isVendoredPath early return (lines 46-48) removed -> FAIL TestSkillLinksResolve: dangling skill link google-catalog/catalog/... (several)
  - after: non-catalog status equals start (only difference is this log file, untracked, created after the start snapshot); unstaged `git diff --stat` identical to start; 0 unstaged modifications under catalog
- 8. `go run ./cmd/sh-agent install --home .agents/scratch/verify-catalog-p1-home` -> exit 0, 13.5s wall (includes build), 1836 written. status: binary 754a7ce6dd7a+dirty (no warning); repo has 35 skills, 8 agents; drift none (1835 files match). `list`: Skills (35). catalog dir entries: .claude 146, .gemini/config 146 (141 entries + INDEX, MANIFEST, CHECKSUMS, 2 LICENSE).
- 9. `bash skills/google-catalog/scripts/refresh-catalog.sh --dest .agents/scratch/verify-catalog-p1-dest --work .agents/scratch/catalog-refresh 2c39459 d5d9052` (RELATIVE paths, as dispatched) -> exit 1 after 3m32s (network OK, clones done): "the catalog is empty: no entry with a SKILL.md under .agents/scratch/catalog-refresh/build" from TestCatalogMatchesChecksums, TestCatalogIndexMatchesEntries, TestCatalogFilesAreCarryable; script ends "the catalog built from these pins is refused". FAIL line.
- 9. Same with ABSOLUTE --dest/--work -> exit 0, 2.3s (clones reused); 141 entries (131 google-skills-cloud, 2 developers, 1 agents-cli-langchain, rest per script); `diff -r verify-catalog-p1-dest skills/google-catalog/catalog` -> exit 0, no difference. The committed catalog is what the script produces at the pins.
- Counts note (4/3): 3 SKIPs in the full run are TestAgyCorpusLeaves, TestRecall, TestHelperIndex (none in the catalog tests; not in the diff; guards of search/dream corpora; not compared with an earlier log, none checked) -> pre-existing? not verified.
- 8(step). Tree check: non-catalog status identical to start apart from this log; unstaged diffstat identical.

## Root cause  (FAIL)
`refresh-catalog.sh` line ~111 runs `cd "$repo" && CATALOG_DIR="$build" CATALOG_WRITE=1 go test ./internal/harness ...`. `go test` runs with cwd = internal/harness, so a relative `--work` (hence relative `$build`) resolves under internal/harness and finds nothing; the tests then call the (real) built catalog empty. The default `work` is absolute (`$repo/.agents/...`), so the default invocation works; a relative `--work`/`--dest` does not. The script neither absolutizes the paths nor rejects relative ones. Reported, not patched: either `realpath`/`cd && pwd` the --work and --dest at parse time, or refuse relative paths with a message.

## Coverage of changed code
- skills/google-catalog/catalog (data) — internal/harness/catalog_test.go (checksums, index, carryable, installs clean), mutation-proven a/b/c
- skills/google-catalog/SKILL.md — TestRouterDescription
- skills/google-catalog/scripts/refresh-catalog.sh — TestRefreshCatalogScript (passes, 2.65s) and the step-9 run; relative --work path NOT covered by any test
- isVendoredPath in TestNoTmpInHarnessProse / TestSkillLinksResolve — mutation e proved SkillLinksResolve; TestNoTmpInHarnessProse use of the helper not mutation-tested
- skills/review/scripts/restatement-grep.sh — TestRestatementGrepSkipsVendoredCatalog, mutation d
- skills/VENDORED.md Vendored catalogs section — TestVendoredRowsCarryAPin passes

## Not examined
Content of the 141 vendored skills beyond the tests; TestNoTmpInHarnessProse helper mutation; the 3 skips' history; git diff of VENDORED.md/README/plan text; refresh script's behaviour with changed pins or --exclude; Antigravity/Claude sandbox starts (guide § 2/4).

## Environment
Linux 6.12.107 debian 13, go from PATH (not recorded), no Docker needed (no Dockerfile in range). Scratch: .agents/scratch/verify-*.
