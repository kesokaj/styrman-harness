---
owner: reviewer
date: 2026-10-09
---
# Review — project-simplify   (2026-10-09)

Range reviewed: the whole repository at HEAD (c267990), no diff.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing in this lens is broken. The one piece of plainly dead code is `search.IsBusy` (`internal/search/export.go:21`), which nothing calls, not even a test. The largest simplification is the docs-root resolver, which is written out three times in two packages (twice inside `internal/metrics`). Next is a hand-rolled frontmatter reader in `internal/dream/collect.go` that duplicates `harness.Split` and `harness.Parse`, even though dream already imports harness. Everything else is small duplication inside `cmd/sh-agent/main.go`, plus a 1 MiB buffer allocated for each trail file on every metrics ingest. The rule-pinned shapes the dispatch named (redaction needles by kind, the size pins, the lint's line grammar) are outside this lens and are not argued against. The split git builders (`trail.Command` versus `harness.gitOut`, BUG-165) and the `search` to `metrics` export seam (ADR-055) were checked and are deliberate. They are not findings.

Savings if every finding is taken: about 90 production lines, plus about 40 test lines.

## Findings

- **[low / confidence: high]** `internal/search/export.go:20-21` — `IsBusy` has no caller in `cmd/`, `internal/` or any `_test.go`. ADR-055's list of exported helpers (`.agents/docs/architecture/DECISIONS.md:670`) does not include it. Shorter: delete both lines. `metrics` retries nothing on BUSY, so no future caller is waiting for it. Saves 2 lines.

- **[low / confidence: high]** `internal/metrics/claude.go:446-451` duplicates `docsRootOf` at `internal/metrics/events.go:227-234` inside the same package. A third copy, which differs only in its fallback, is at `internal/search/reconcile.go:73-80`. Shorter: in `writer.docsRoot`, replace the loop with `d := docsRootOf(root)` (saves 6 lines). Then move the resolver to `internal/trail` as `trail.DocsRoot(root string) string`, which returns `""` when neither directory exists. `search`, `metrics` and `dream` already import trail. search keeps its `if d == "" { d = ".agents/docs" }`. That saves another ~8 lines and leaves one place to change the rule. Do NOT fold in the `.agents/docs`-only checks at `internal/search/search.go:202,321,366`: those deliberately register only scaffolded checkouts (security review, finding 2).

- **[low / confidence: high]** `internal/dream/dream.go:112-115` and `internal/dream/collect.go:71-74` are the same four-line default (`docs := p.DocsRoot; if docs == "" { docs = ".agents/docs" }`). Shorter: one `func docsOf(p search.Project) string` in dream, called from both. Saves 3 lines.

- **[low / confidence: medium]** `internal/dream/collect.go:289-316` — `frontField` and its `unquote` re-implement a frontmatter reader (28 lines). dream already imports `internal/harness`, which exports `Split` and `Parse` (`internal/harness/frontmatter.go:22,76`). Shorter, at `collect.go:274`:
  ```go
  front, _, _ := harness.Split(body)
  fm := harness.Parse(front)
  name, desc := fm["name"], fm["description"]
  ```
  This saves ~24 lines. Side effect (confidence: low): `frontField` reads only the key line, so a block-scalar `description: >` comes back as the literal `>` and the dreamer sees that as the skill's description. `Parse` folds the block. I did not check whether any installed skill uses a block scalar.

- **[low / confidence: medium]** `internal/harness/frontmatter.go:221-222` — the exported `Lint` has no production caller (the installer calls `ParseStrict` at `inventory.go:190,241,299`). It exists for tests (`harness_test.go:88,101,2378-2398`), and `phase8_test.go:165` exists only to prove that `Lint` and `ParseStrict` agree. Shorter: delete `Lint`. Tests call `_, err := ParseStrict(x)`, and the agreement test goes. Saves 2 production lines and ~8 test lines. Keep it only if it is meant as an API for a planned caller. If so, its doc comment should say that.

- **[low / confidence: high]** `cmd/sh-agent/main.go:648-654` (`runInstall`) builds the `--home`/`--vendor` pass-through by hand, duplicating `installPassArgs` at `main.go:1097-1110`. The one difference is that runInstall forwards `--vendor all` when it is given explicitly, and the child's behaviour is identical either way. Shorter: `pass, err := installPassArgs(c.homeFlag, c.home, c.vendorFlag)`. Saves 5 lines and leaves one place that decides what a child install inherits.

- **[low / confidence: medium]** `cmd/sh-agent/main.go:453-504` — the `show` and `report` arms of `runMetrics` repeat the same sequence: no-args check, `metricsShowOptions`, `existingMetricsStore`, `defer Close`, `Show` (~14 lines each). Shorter: `func (c *ctx) metricsSummary(sub string) (*metrics.Summary, error)`, which opens, shows, closes and returns the summary, with each arm keeping only its output step. Saves ~12 lines. It also brings `runMetrics` (113 lines) under the 80-line mark.

- **[low / confidence: medium]** `internal/metrics/events.go:251-260` — `readSmall` allocates a fresh 1 MiB buffer for every handoff, verify log, review and bug file of every project root on every `metrics ingest` (`events.go:109,121,136`), and reads it with a single `f.Read`. The cost scales as files × 1 MiB of garbage on each ingest. Shorter, with the same line count: `b, _ := io.ReadAll(io.LimitReader(f, 1<<20)); return string(b)`. The allocation then follows the file's size, and a short read is no longer possible.

- **[low / confidence: low]** `cmd/sh-agent/main.go:875-1033` — `status` (158 lines) splits along a natural seam at line ~971. Above it is the "installed / binary / PATH / retro / dream" header, which needs only `home`. Below it is the "repo has / drift" half, which needs `src` and `vendor`. Shorter: `statusInstalled(home, src, out) (harness.State, bool, error)` and `statusDrift(src, home, vendor, fromZip, out) error`. No lines saved, but each half can then be tested and edited on its own. `harness.Install` (`internal/harness/apply.go:19`, 168 lines) and `search.gitView` (`internal/search/reconcile.go:174`, 161 lines) are the other two functions over 150 lines. I did not read their bodies (see Not checked).

- **[low / confidence: low]** `testmain_test.go:14` — the same 8-line `TestMain` (md5-identical body) is copied into 6 packages: root, `cmd/sh-agent`, `internal/{harness,trail,search,memory}`. Near-copies sit in `internal/dream/main_test.go` and `internal/metrics/testmain_test.go`. Shorter: an `internal/testenv` package with `func ScrubGit()`, and each `TestMain` becomes `testenv.ScrubGit(); os.Exit(m.Run())`. Saves ~25 lines and leaves one place to add the next variable. The trade-off is one extra package. Defer is a fair answer.

## Checked and not findings
- `internal/harness` imports neither `modernc.org/sqlite` nor `search`/`metrics` (`go list` import graph). The rule holds.
- `go build` runs once per install, only in `harness.Reinstall` (`internal/harness/repo.go:308`). An update also builds once.
- The git exec builders are split on purpose: `trail.Command` for project repositories, `harness.gitOut`/`gitOutC` for the harness's own checkout (comment at `internal/trail/trail.go:120-124`, BUG-165).
- `internal/search/export.go` is the seam ADR-055 chose ("Shared helpers move", `DECISIONS.md:670`). Every export except `IsBusy` has a caller in `internal/metrics`.
- `embed.go` has to sit at the repository root because `go:embed` cannot reach a parent directory. `plan_test.go`, `init_project_test.go`, `embed_test.go` and `testmain_test.go` test that package and the tree beside it, so they belong at the root too.
- `cmd/sh-agent/main.go`'s dispatch is a table (`verbs`, line 145) with per-verb flag registration. The flat `ctx` field bag is the cost of that design, and I see no concrete failure in it.
- An unexported-function scan per package found no function without a production caller.

## Criteria
1. Dead code — fail (minor): `search.IsBusy` is dead and `harness.Lint` is used only by tests. No unexported function lacks a caller, and I found no commented-out blocks in the files I read.
2. Flow tracing — not in this lens
3. Standards conformance — not in this lens
4. Directory layout — not in this lens (asked for anyway: the root-level Go files belong where they are, see above)
5. Simplicity (DIY-except-security) — pass with nits: there is duplication (docs root, frontmatter, pass-args, metrics arms) but no new dependency and no hand-rolled security primitive. In the code I read, `crypto/sha256`/`crypto/rand` are used as designed.
6. Decoupling — pass: harness stays free of the driver, the git builders split on purpose, and `metrics`→`search` is the seam ADR-055 recorded. `internal/trail` is the natural home for the shared docs-root resolver.

## Not checked
- Bodies not read: `internal/harness/report.go` (1061 lines, apart from the function sizes), `apply.go` `Install`, `zip.go`, `release.go`, `settings.go`, `adopt.go`, `plan.go`, `inventory.go` beyond its function list; `internal/search/reconcile.go` `gitView`, `search.go` `refs`/`reconcileGlobal`, `dream.go` (search), `parse.go` beyond lines 95-140, `store.go`; `internal/metrics/show.go`, `agy.go`, `claude.go` apart from 436-522, `verb.go` past line 60, `report.go`; `internal/dream/run.go`, `timer.go`, `collect.go` outside 44-76 and 255-320; `internal/memory/memory.go` bodies; `internal/trail/trail.go` outside 120-148.
- Test files: only helper names, the `TestMain` copies, and `plan_test.go`'s head were checked. I did not look for fixture duplication inside `harness_test.go` (2689 lines), `catalog_test.go` (2033) or `cmd/sh-agent/main_test.go` (1909).
- Dead-code scans: top-level exported `func`/`type`/`const`/`var` declarations and unexported funcs only. Not covered: identifiers inside `const (...)`/`var (...)` blocks, struct fields, `ctx` fields that no flag sets, and exported names used only inside their own package. Each scan counts word matches, so a name shared with another identifier could hide a dead one.
- Pre-flight: no `sh-agent search --refs` (there is no diff), and I did not read `REQUIREMENTS.md` or `.agents/docs/bugs/INDEX.md` in full. Rule-pinned shapes were not examined, as the dispatch directed.
- Efficiency: matchCommit's map-order prefix match (`internal/metrics/events.go:218`) and the per-query reconcile in `search` were not measured.
