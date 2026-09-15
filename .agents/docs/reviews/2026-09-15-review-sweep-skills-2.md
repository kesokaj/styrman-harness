---
owner: reviewer
date: 2026-09-15
---
# Review — sweep-skills-2   (2026-09-15)

Range reviewed: working tree, `skills/golang/` read whole — `SKILL.md` (59 lines) and all seven `references/*.md` (5863 lines). This is not a diff (the tree was clean for `skills/golang`). Lens: `--docs --project`. Compared against `context/AGENTS.md`, `context/dev-setup.md`, `AGENTS.md`, `REQUIREMENTS.md`, `skills/review-standards/references/go-review.md`, `skills/VENDORED.md`, `.agents/docs/testing/skill-matrix.md`, `go.mod`, and the other files inside the scope.

## Verdict: REQUEST CHANGES

`SKILL.md` is accurate, but all seven reference files are outdated. They contradict the skill's own coding standards, which the review agents enforce. Worst first:
- `golang-testing.md` teaches a `testing/synctest` call that does not compile. `golang-architecture.md:725` in the same skill says that API was replaced.
- `golang-architecture.md:544` says `Close()` errors need no check. `go-review.md:23` says the opposite, because an unchecked close on a writer can lose the final flush.
- `golang-io-web.md` makes Cobra + Viper the default CLI stack, registered in `init()`, and uses testify in its handler test. This breaks stdlib-first, no-`init()` and no-assertion-library (`SKILL.md:27,40,41`, `context/AGENTS.md:49`).
- `golang-ecosystem.md` sends benchmark "audit trail" files to `/tmp`, against `context/AGENTS.md:53`.
- About 150 pointers name deep-dive guides, `assets/examples/*.go`, `templates/*` and `*.md` debug files that were cut in e5833a3 and no longer exist.

## Classification (--project)

| File | Class | Proof |
|---|---|---|
| `skills/golang/SKILL.md` | accurate | `SKILL.md:49-55` names exactly the 7 files, and each sub-topic it lists exists as a `## Sub-topic:` header. `SKILL.md:40` matches `context/AGENTS.md:49` and `context/dev-setup.md:5`. Its rules are the ones the references contradict (see findings). |
| `references/golang-core.md` | outdated | `:61` `wg.Go()` "Go 1.24+" (it is 1.25). `:230` maps "8-entry buckets and overflow chains" (Swiss tables since 1.24). `:399` sends samber/oops readers to Ecosystem, but oops lives in `golang-samber.md:539`. |
| `references/golang-testing.md` | outdated | `:143,161` synctest shown as "Experimental" with `synctest.Run(func(t *testing.T)`. `:252` `is.NotNil` breaks its own `:28`. `:513-527` lists reference files that do not exist. |
| `references/golang-architecture.md` | outdated | `:544` Close errors unchecked. `:534` and `:612` point at itself for content that is in `golang-structure.md`. `:647,668` `/golang-modernize` does not exist. |
| `references/golang-structure.md` | outdated | `:497` "see the decision table below": the table was replaced by `:561`. `:489` "<10 services" threshold breaks `:561` and `SKILL.md:26`. |
| `references/golang-io-web.md` | outdated | `:413` `assert.New(t)`. `:893` Cobra+Viper default. `:929-1047` twelve `assets/examples/*.go` pointers, but there is no `assets/` directory. |
| `references/golang-ecosystem.md` | outdated | `:792,795,799` use `/tmp`. `:323` says every project MUST have `.golangci.yml`, against `SKILL.md:23`. `:464-496` names debug files that do not exist. `:240` still teaches the tools.go pattern, which `golang-architecture.md:760` supersedes. |
| `references/golang-samber.md` | outdated | `:3` is accurate (stdlib first), but `:192` marks a plain loop "✗ boilerplate, error-prone" and `:722` says `oops.Recover()` "MUST" be used. `:627-638` does not compile. It has 20+ pointers to guides that do not exist. |

No file is wholly redundant.

## Findings

### High

None blocking by `REQUIREMENTS.md`. `REQUIREMENTS.md:17-19` (R4, stdlib only) covers the harness binary, not user projects.

### Medium

- **[medium / confidence: high] divergent** `skills/golang/references/golang-testing.md:143` and `:161`, against `golang-architecture.md:725`.
  - **Wrong:** the file calls `testing/synctest` "Experimental … not yet covered by Go's compatibility guarantee" and shows `synctest.Run(func(t *testing.T) { … })`. `Run` never took a `*testing.T`: in Go 1.24 it was `Run(func())`, and only behind `GOEXPERIMENT=synctest`. The package became stable in Go 1.25 as `synctest.Test(t, func(t *testing.T))`, which `golang-architecture.md:725` itself lists as the replacement. On this stack (`go.mod` 1.24.4, docker skill `golang:1.24`) the sample fails to compile either way.
  - **Should say:** "`testing/synctest` (Go 1.25+; Go 1.24 only with `GOEXPERIMENT=synctest`): `synctest.Test(t, func(t *testing.T) { … })`."
- **[medium / confidence: high] divergent** `skills/golang/references/golang-core.md:61`, against `golang-architecture.md:750`.
  - **Wrong:** core says "Go 1.24+: `wg.Go()`"; architecture says Go 1.25+. `sync.WaitGroup.Go` shipped in Go 1.25. Code written from core fails to build on the stack's pinned 1.24 (`go.mod:3`).
  - **Should say:** "Go 1.25+: `wg.Go()`".
- **[medium / confidence: high] divergent** `skills/golang/references/golang-architecture.md:544`, against `skills/review-standards/references/go-review.md:23` and `skills/golang/SKILL.md:19`.
  - **Wrong:** "`.Close()` errors: acceptable to not check — `defer f.Close()` is fine without error handling." The review checklist requires `Close()` errors to be checked on writers, because a swallowed close loses the final flush. A writer would ship code the reviewer then flags, and the data-loss path is real.
  - **Should say:** "An unchecked `defer Close()` is acceptable on read-only resources. On writers (files written, `gzip`/`bufio` writers) check the `Close()` error."
- **[medium / confidence: high] divergent** `skills/golang/references/golang-io-web.md:413-432`, against `golang-testing.md:28,342` and `SKILL.md:41`.
  - **Wrong:** the "Testing HTTP Handlers" sample uses `is := assert.New(t)` (testify) and `is.NoError`, `is.Equal`, `is.Contains`. The testing reference and SKILL.md say no assertion library by default and hand-written fakes.
  - **Should say:** use stdlib checks (`if err != nil { t.Fatalf(…) }`, `if rec.Code != http.StatusOK { t.Errorf(…) }`) with a hand-written fake.
- **[medium / confidence: high] divergent** `skills/golang/references/golang-io-web.md:893`, against `SKILL.md:27`, `SKILL.md:40`, `context/AGENTS.md:49` and `context/dev-setup.md:5`.
  - **Wrong:**
    - `:893` says "Use Cobra + Viper as the default stack for Go CLI applications"; stdlib `flag` only "for trivial single-purpose tools" (`:897`).
    - `:903-910` adds six third-party packages as the Quick Reference.
    - `:944` registers subcommands "in `init()`", which SKILL.md forbids unless absolutely required.
    - `:963-965` binds everything to Viper's global instance, which is global mutable state (`SKILL.md:26`).
    - The harness's own CLI is stdlib only (`REQUIREMENTS.md:19`).
  - **Should say:** "stdlib `flag` (with a small subcommand switch) is the default. Cobra/Viper is a dependency to justify when a deep subcommand tree or shell completions demand it. Wire subcommands from `main`, not `init()`."
- **[medium / confidence: high] divergent** `skills/golang/references/golang-architecture.md:200` and `golang-samber.md:192`, against `SKILL.md:40`, `SKILL.md:55`, `golang-samber.md:3` and `golang-architecture.md:199`.
  - **Wrong:**
    - `golang-architecture.md:200`: "for filter/group-by/chunk, use `github.com/samber/lo`".
    - `golang-samber.md:139`: "Everything else … requires manual for-loops. `lo` fills this gap".
    - `golang-samber.md:192`: the plain loop is marked "✗ Manual — boilerplate, error-prone".
    - `golang-samber.md:243`: "Compose lo functions … instead of writing nested loops".
    - These make a dependency the default, while SKILL.md says samber only "when stdlib is demonstrably not enough" and `golang-architecture.md:199` quotes "a little copying is better than a little dependency".
  - **Should say:** "write the loop; `lo` only when the stdlib is demonstrably not enough, justified in the change entry". Drop the ✗ label on the loop.
- **[medium / confidence: high] divergent** `skills/golang/references/golang-core.md:373` and `golang-samber.md:722`, against `SKILL.md:18`, `SKILL.md:55` and `golang-samber.md:3`.
  - **Wrong:** `golang-core.md:373` says "Use `samber/oops` for production errors needing stack traces, user/tenant context…" with no stdlib-first qualifier. `golang-samber.md:722` says "`oops.Recover()` MUST be used in goroutine boundaries", an unconditional MUST for a third-party library.
  - **Should say:** "wrap with `fmt.Errorf("…: %w")` and attach variable data as `slog` attributes; `samber/oops` only when that is demonstrably not enough". For `golang-samber.md:722`: "if the project uses oops, recover at goroutine boundaries with `oops.Recover()`".
- **[medium / confidence: high] stale** `skills/golang/references/golang-structure.md:497` (also `:489`, `:493-494`, `:506-507`), against `golang-structure.md:561`, `:600` and `SKILL.md:26`, `:41`.
  - **Wrong:**
    - `:497` "Choose the right DI library for your project size and team — see the decision table below": the table was cut, and `:561` now says "the default answer is no".
    - `:489` "Small projects (< 10 services) SHOULD use manual constructor injection" implies a library above 10 services. `:554` says 15+ and `:600` says under 20, so the file gives three thresholds.
    - `:493-494` (lazy init, singletons/transients) and `:506-507` (centralized container, lazy loading) are leftovers of the cut container material.
  - **Should say:** "Dependencies are wired by hand in `main()`; a DI library is a dependency to justify, default no (see DI Libraries)." Delete the decision-table pointer and the container-lifecycle items.
- **[medium / confidence: high] stale** `skills/golang/references/golang-ecosystem.md:792`, `:795`, `:799`, against `context/AGENTS.md:53`.
  - **Wrong:** benchmark baselines go to `/tmp/report-1.txt`, and "Keep all `/tmp/report-*.txt` files as an audit trail". The global rule is never `/tmp`; scratch lives in `.agents/scratch/`. `/tmp` is also cleared on reboot, so the "audit trail" does not survive.
  - **Should say:** `… | tee .agents/scratch/bench/report-1.txt` and `benchstat .agents/scratch/bench/report-1.txt .agents/scratch/bench/report-2.txt`.
- **[medium / confidence: high] divergent** `skills/golang/references/golang-ecosystem.md:323` (also `:321`, `:376`), against `SKILL.md:23` and `AGENTS.md:15`.
  - **Wrong:** "Every Go project MUST have a `.golangci.yml` — it is the source of truth". SKILL.md makes `go vet` + `staticcheck` the gate. This repo's own gate is `gofmt -l . && go vet ./... && go test -race ./...`, and it has no `.golangci.yml`. `golang-architecture.md:665` also says run golangci-lint. Two sources name different lint gates.
  - **Should say:** one gate, stated in SKILL.md. For example: "`go vet` + `staticcheck`; golangci-lint when the project already configures it".
- **[medium / confidence: medium] divergent** `skills/golang/references/golang-ecosystem.md:105-113` (also `:24`, `:27`, `:32`), against `context/AGENTS.md:13-15` and `context/dev-setup.md:5`.
  - **Wrong:** a second "Definition of Done for Observability" says a feature is not done without Prometheus metrics with PromQL comments, spans on every service method and DB query, Grafana dashboards and alert rules, and PostHog/Segment RUM events. `:24` says "Every HTTP endpoint MUST have latency and error rate metrics"; `:27` says "Set up OpenTelemetry tracing on new projects". This conflicts with the global Definition of Done, and it adds Prometheus, OTel and PostHog as defaults on a stdlib-first, Cloud Run stack.
  - **Should say:** a checklist framed "when the service has a metrics/tracing backend", not a Definition of Done. Or cut it.
- **[medium / confidence: medium] security** `skills/golang/references/golang-ecosystem.md:1030`.
  - **Wrong:** "Register for discoverability — add to Context7, DeepWiki, OpenDeep, zRead. **Even for private libraries.**" This tells an agent to submit private source to third-party indexers, which leaks it (compare `skills/VENDORED.md:24`, "this repository is not published").
  - **Should say:** "Public libraries only; never register a private repository with an external indexer."
- **[medium / confidence: high] stale — files that do not exist, cut in e5833a3.** `find skills/golang -type d` shows only `references/`.
  - **Wrong:**
    - `golang-io-web.md:929, 933, 944, 948, 969, 989, 1010, 1026, 1030, 1039, 1043, 1047` point at `assets/examples/{main,root,serve,flags,args,config,version,exit_codes,output,signal,completion,cli_test}.go`.
    - `golang-ecosystem.md:464-496` points at `compilation.md`, `common-go-bugs.md`, `testing-debug.md`, `diagnostic-tools.md`, `concurrency-debug.md`, `pprof.md`, `performance-debug.md`.
    - `golang-ecosystem.md:995`, `:1019`, `:1060` point at `templates/README.md`, `templates/CHANGELOG.md`, `templates/llms.txt`.
    - An agent told to "copy the template from templates/README.md" finds nothing and may invent its contents.
  - **Should say:** inline the gist or delete the pointer.

### Low

- **[low / confidence: high] stale — named guides that do not exist.** They are written as links but point at nothing, so the reader is sent looking for guides that were never installed. Delete them, or inline one line.
  - `golang-core.md:29, 65, 79, 164-168, 226, 249, 272, 276, 289, 299, 381-385, 486, 508`
  - `golang-testing.md:105, 189, 338, 511, 513-527`
  - `golang-architecture.md:115, 174, 315-323, 562, 582, 616-621, 705, 709`
  - `golang-structure.md:378, 430, 436-439, 443, 447, 514`
  - `golang-io-web.md:496, 618, 747, 828, 839, 861-866`
  - `golang-ecosystem.md:49-65, 257-269, 372, 572-591, 653, 670-688, 745, 784, 826-833, 910, 915, 991, 1013, 1032, 1042, 1054, 1088`
  - `golang-samber.md:42, 110, 112, 180, 273, 275, 304, 344, 379, 405, 471, 528, 784, 892, 909, 945, 947, 1014, 1040, 1079, 1095`
  - `golang-ecosystem.md:271` is an empty `## Cross-References` header.
- **[low / confidence: high] stale — cross-references to the wrong file or to each other.**
  - `golang-architecture.md:534` points at "the Architecture reference" for `var _ Interface = (*Type)(nil)`. Should point at `golang-structure.md` (the pattern is at `:146-154`).
  - `golang-architecture.md:612` points at Architecture for 12-factor. Should point at `golang-structure.md` (`:384`).
  - `golang-core.md:399` points at Ecosystem for the samber/oops API. Should point at `golang-samber.md` (`:539`).
  - `golang-core.md:109` sends graceful shutdown to Architecture, and `golang-architecture.md:634` sends it to Core. It lives only in `golang-io-web.md:355`.
  - `golang-io-web.md:436` and `:461` send "comprehensive HTTP testing" to Testing, and `golang-testing.md:105` sends it to "HTTP Testing" (not in the skill).
  - `golang-io-web.md:872` and `golang-testing.md:390` send database integration tests to each other; neither has them.
  - `golang-io-web.md:654` (gRPC error mapping in Core) and `:655` (bufconn in Testing): neither target has it.
  - `golang-core.md:550` (integer overflow section in Ecosystem): absent.
  - `golang-testing.md:386`, `:392`, `golang-core.md:178`, `:640` and `golang-architecture.md:211`, `:373` send readers to Ecosystem for linter configuration. `golang-ecosystem.md:347-349` has none.
  - `golang-ecosystem.md:901` (Structure, for where doc files belong): absent.
- **[low / confidence: high] stale** `golang-architecture.md:647` and `:668` name an "explicit `/golang-modernize` invocation". There is no such command (`commands/` is empty) and no such skill among the 27. Should say "when the user asks for a full modernization scan".
- **[low / confidence: high] stale** `golang-testing.md:532`, `golang-ecosystem.md:37`, `:596` and `golang-samber.md:119` point to the `samber/cc-skills@promql-cli` skill. It is not in `skills/` and not in `skills/VENDORED.md`. Remove these pointers.
- **[low / confidence: high] divergent** `golang-testing.md:252` `is.NotNil(result)`, against `golang-testing.md:28` and `SKILL.md:41` (no assertion library). The `:251` indentation is also broken. Should say `if result == nil { t.Fatal("Process returned nil") }`.
- **[low / confidence: medium] divergent** `golang-testing.md:27` and `golang-core.md:27`: "Packages with goroutines SHOULD use `goleak.VerifyTestMain`" makes `go.uber.org/goleak` a default. That conflicts with `golang-testing.md:28` and `SKILL.md:40`. Should say "`goleak` is a test dependency to justify; the `-race` run and explicit goroutine joins come first".
- **[low / confidence: high] divergent** `golang-testing.md:201`, `:209`, `:226` use `for i := 0; i < b.N; i++` as the benchmark pattern. The same file at `:440` says `b.N` loops "should migrate to `b.Loop()`", and so does `golang-architecture.md:752`. Use `for b.Loop()`.
- **[low / confidence: high] divergent** `golang-ecosystem.md:240-255` presents the `tools.go` blank-import pattern as current. `golang-architecture.md:760` says "Move tool deps to `go.mod` tool directives (Go 1.24+)", and the stack is on 1.24. Should say `go get -tool <pkg>` and a `tool` directive, with `tools.go` only below Go 1.24.
- **[low / confidence: high] divergent** `golang-ecosystem.md:168` ("Using `go get -u` to upgrade an existing dependency is safe") against `golang-ecosystem.md:223` ("Prefer `go get -u=patch` … minor upgrades may … change behavior"). Should say "`go get -u=patch` is the routine upgrade".
- **[low / confidence: high] divergent** `golang-architecture.md:472` ("init() … across files in **filename alphabetical order**") against `golang-core.md:661` ("`init()` execution order across files is unspecified"). Should say "the spec only encourages lexical filename order; the go tool follows it; never depend on it".
- **[low / confidence: high] divergent** `golang-core.md:628` discards `sql.Open`'s error with `_` in the recommended `sync.Once` example, against `SKILL.md:19` and `golang-core.md:363`. Should use `sync.OnceValues` returning `(*sql.DB, error)`.
- **[low / confidence: high] divergent** `golang-io-web.md:530` uses `log.Printf("method=%s …")` in the recommended gRPC logging interceptor, against `SKILL.md:30` (slog JSON) and `golang-core.md:372`. Should say `slog.InfoContext(ctx, "rpc", "method", info.FullMethod, …)`.
- **[low / confidence: medium-high] outdated** `golang-core.md:230`: "Maps are hash tables with 8-entry buckets and overflow chains". Since Go 1.24 (the stack's version) maps are Swiss tables. Should say "Swiss tables since Go 1.24 (groups of 8 slots, open addressing); buckets with overflow chains before".
- **[low / confidence: high] outdated** `golang-io-web.md:30`: "`ReverseProxy` with `.Rewrite` (Go 1.26+)". `Rewrite` has existed since Go 1.20; 1.26 is only when `Director` was deprecated (`golang-architecture.md:727`). Should say "`.Rewrite` (Go 1.20+)".
- **[low / confidence: high] outdated — code samples that do not compile.**
  - `golang-io-web.md:752-755`: `GetUser(id string)` uses an undefined `ctx`.
  - `golang-io-web.md:770-773`: `err :=` redeclares the named result `err` ("no new variables on left side of :=").
  - `golang-samber.md:627-638`: `return oops…` inside a handler with no result.
  - `golang-samber.md:730`: `With("panic_value", r)` with `r` undefined.
  - `golang-samber.md:609-611`: `r.db.Query(query, id)` without a context and with rows never closed, against `golang-io-web.md:676`, `:678`.
  - `golang-samber.md:667`: interpolates `req.ProductID` into the message of a "good" example, against its own MUST at `:710`.
- **[low / confidence: medium] divergent** CI-oriented guidance against `context/AGENTS.md:45` ("no CI, no hooks, ever"): `golang-ecosystem.md:267` (Dependabot/Renovate auto-merge PRs), `:1010` (CI build-status badge), and `golang-architecture.md:647` (Agent 5 scans the "CI pipeline").
- **[low / confidence: medium] outdated** `golang-ecosystem.md:412` and `:417`: "`run.timeout` (default: 5m)" and "`run.skip-dirs`". These are golangci-lint v1 settings, while the same file targets v2 (`:334`, `:416`). v2 has no default timeout and moved directory exclusion to `linters.exclusions.paths`. Unverified against the current golangci-lint docs.
- **[low / confidence: medium] provenance** `skills/golang/` derives from `samber/cc-skills`: it has that project's Persona/Modes format, `golang-testing.md:532` names `samber/cc-skills@promql-cli`, and `golang-ecosystem.md:33` links samber's alert library. It was imported in 9dcdc4c and cut in e5833a3, but `skills/VENDORED.md` has no source or licence row for it. `VENDORED.md:14-26` records attribution for the other third-party-derived skills. Add a row, or state in `VENDORED.md` that `golang` is a fork. The upstream licence is unverified.

## Unverified

- `golang-architecture.md:653`, `:693-695`: "latest Go is 1.26". Go 1.27 would be due around August 2026. The local toolchain is 1.24.4, so this could not be checked.
- `golang-architecture.md:726-727` (Go 1.26 deprecations of `EncryptPKCS1v15` and `ReverseProxy.Director`) and `:701` (`modernize` linter since golangci-lint v2.6.0).
- `golang-io-web.md:978`: `cobra.ExactValidArgs`, possibly deprecated upstream.
- All `samber/*` API signatures and tables in `golang-samber.md` (hot, lo package paths `lo/it` and `exp/simd`, mo, oops, ro, slog-*) were not checked against upstream.

## Criteria

1. Dead code — n/a — docs lens. The doc equivalent is the stale pointers above.
2. Flow tracing — n/a — docs lens. Cross-references were traced file to file; see the stale findings.
3. Standards conformance — n/a — docs lens. The divergences from `SKILL.md`, `context/AGENTS.md` and `go-review.md` are listed as findings.
4. Directory layout — n/a — docs lens.
5. Simplicity (DIY-except-security) — n/a — docs lens. The dependency-by-default findings (Cobra/Viper, lo, oops, testify, goleak, DI) are listed.
6. Decoupling — n/a — docs lens.

## Not examined

- Nothing in the named scope was skipped: all 8 files were read whole.
- Old wording of rules removed by e5833a3 was not `git grep`ed across the living docs. The scope is not a diff, so the `--project` clause applied instead of dropped/narrowed analysis.
- No external source (Go release notes, golangci-lint, samber repositories) was fetched. Version claims rest on the internal contradictions cited, plus reviewer knowledge where marked.

## Re-check commands

```bash
git grep -n 'synctest.Run\|Experimental:\*\* .testing/synctest' -- skills/golang
git grep -n 'Go 1.24+: .wg.Go' -- skills/golang
git grep -n 'acceptable to not check' -- skills/golang
git grep -n 'assert.New(t)\|is.NotNil\|is.NoError' -- skills/golang
git grep -n 'Cobra + Viper as the default\|registering them in .init()' -- skills/golang
git grep -n 'use .github.com/samber/lo.\|Manual — boilerplate\|MUST be used in goroutine boundaries\|Use .samber/oops. for production' -- skills/golang
git grep -n 'decision table below\|< 10 services' -- skills/golang
git grep -n '/tmp/' -- skills/golang
git grep -n 'MUST have a .\.golangci.yml' -- skills/golang
git grep -n 'Definition of Done for Observability\|Every HTTP endpoint MUST have' -- skills/golang
git grep -n 'Even for private libraries' -- skills/golang
git grep -n 'assets/examples/\|templates/\(README\|CHANGELOG\|llms\)\|common-go-bugs.md\|pprof.md\|compilation.md' -- skills/golang
git grep -n '/golang-modernize\|cc-skills@promql-cli' -- skills/golang
git grep -n 'tools.go Pattern\|go get -u. to upgrade an existing dependency is safe' -- skills/golang
git grep -n 'filename alphabetical order\|db.conn, _ = sql.Open\|log.Printf("method=' -- skills/golang
git grep -n '8-entry buckets\|.Rewrite. (Go 1.26+)\|run.skip-dirs' -- skills/golang
git grep -n 'Dependabot\|Build Status' -- skills/golang
git grep -n 'See the Architecture reference (.golang-architecture.md.) for the .var _\|samber/oops API\|for graceful shutdown patterns\|for goroutine lifecycle and graceful shutdown' -- skills/golang
git grep -n 'golang\|cc-skills' -- skills/VENDORED.md
```
