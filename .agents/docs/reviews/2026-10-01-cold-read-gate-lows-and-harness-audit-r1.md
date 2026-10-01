---
owner: reviewer
date: 2026-10-01
---
# Cold read — gate lows and harness-audit plan, round 1   (2026-10-01)

Range reviewed: the plan file and the surfaces it names at HEAD (`.agents/docs/plans/2026-10-01-gate-lows-and-harness-audit.md`; read: gate report rows L1-L10, `internal/dream/timer.go` and `dream.go`, `cmd/sh-agent/main.go` (runUninstall, runDream, status), `internal/harness/inventory.go` and `frontmatter.go`, `internal/search/reconcile.go` and `search.go`, `skills/google-catalog/scripts/refresh-catalog.sh`, `internal/harness/catalog_test.go`, `size_test.go`, ADR-014, the recon/audit plans' ID prefixes)

## Status: COMPLETE

## Verdict: design needs revision; implementation not ready

One decision in the plan contradicts a standing operator decision and is recorded as "by: implementer" (D4, below, BLOCKING). Phase 1 is otherwise close to implementable: items 2, 3, 4, 5, 10, 11 have small, checkable shapes. Phase 2 (the skill) is a one-paragraph shape for the plan's largest deliverable, with several choices only the author can make.

## Gaps

G1 (BLOCKING) D4 and interface 6 make `harness-audit` slash-only (`disable-model-invocation: true`). ADR-014 (`DECISIONS.md:109`, operator 2026-09-13: "I want ALL skills to be able to be executed by the agent, not hidden behind user only") removed the flag from every harness skill and reserved it for third-party skills; README.md:243 states "none of the harness's own do". The plan neither cites nor supersedes ADR-014, D4 is "by: implementer", and the Why does not ask for it. The stated reason for D4 (the listing budget) is not forced: the plan's own fallback (Risks) is a description cut to the 264 bytes, which a model-invocable skill needs anyway (name plus description at most 264 bytes; `go test ./internal/harness -run TestSkillListingBudget -v` shows 19,236 of 19,500 over 35 skills). Also unlisted consequence: no skill is slash-only today, so the first one flips `slashNote` (`inventory.go`) for every deployed context ("Skills marked `/name` are typed by the user; the rest load when..."), which the size pins in `size_test.go` and the context's wording both see. Needs the operator's choice (keep ADR-014: model-invocable with a trigger-list description that keeps the user's-words safeguard, or a new ADR superseding it plus the README:243 edit and the inventory sentence in the same change).

G2 (BLOCKING) "Runs in the harness checkout only" has no detection rule, and the E2E contradicts the failure-mode text. The failure mode names `sh-agent status` as the way to find the checkout; E2E (c) runs the skill in a clone of the repository, which is not the checkout `status` reports. Either detection is by tree shape (`context/AGENTS.md`, `skills/`, `cmd/sh-agent` present) and the clone passes, or by comparing with status's source and the E2E cannot pass. Pick one.

G3 (BLOCKING for the skill's deliverable) The skill's procedure is not decided, only its output file names. Not in the plan: which agent does recon per area and which does the audit (the 2026-10-01 slices are `reviewer` reports, the recon brief is `explorer` territory); the tool-call budget per dispatched agent and whether areas run in parallel; the severity scale; what "permanent IDs" means across audits (the recon brief's IDs are the previous audit's; how a new finding gets the next number); the default of `--area` omitted (all nine areas = a full audit, which Scope says is the operator's to start, "Out"); the default of `--reaudit` with no brief (newest `plans/*-recon-*harness-audit*.md`?); what `<n>` in `audit-<n>-<area>` counts; "the tally appended to a new report" (which file, what shape). Also the area table does not match the source audit: the recon (`plans/2026-09-30-recon-audit-fix.md` headings) has C, WS, TS, AG, IN, SD, TR, G; the plan lists `DL`, which exists nowhere in the recon or the audit-fix plan; and the re-audit's six slice reports (context-trail, workflow-skills, technical-skills, agents, installer, search-memory-dream) cut the nine prefixes differently from `--area <ids>`. The plan must say what `DL` is and how an area maps to a slice, or the area file `references/areas.md` will be invented.

G4 Item 1, interfaces (timer):
- `uninstall`'s "found nothing" text is not specified (only `--remove-timer`'s two lines and the legacy notice), though the L1 row names `uninstall`. `runUninstall` (`main.go:292-306`) gates on `dream.TimerInstalled`, which checks both the Linux and the macOS paths whatever the OS, while `RemoveTimer` removes only the running OS's files (`timer.go:226-246`): on a mismatch `removed` is false after `TimerInstalled` was true. Say which one decides the message.
- `LegacyTimerNotice(home)` takes no `goos`, but the check command differs per OS (`systemctl ... FragmentPath` vs `launchctl print`); say whether it prints one by the files found or both.
- `StatusLine` returns early in the `state unreadable` branch (`dream.go:196-198`); say whether the legacy suffix goes on every return. `TestStatusLine` pins exact strings.
- The dry run of `--remove-timer` today prints `[dry-run] would remove the nightly dream timer; installed now: <bool>` (`main.go:339-343`); "prints the same notice" must say whether that line stays. `TestDreamTimerDryRun` asserts the home stays empty, so the notice must read only.
- The signature change `RemoveTimer(...) (bool, error)` touches `main.go:301`, `main.go:347` and six calls in `timer_test.go`; fine, but the plan's surface list omits the test files.

G5 Item 6, search header. `gitView(root, prefix)` is called twice (`search.go:172` per project, `search.go:317` per project in global scope with prefix `.agents/docs/`), and `sc.status()` (`search.go:428`) is shared by the words header (`:643`) and `--open` (`:877`); `index` prints elsewhere. Not decided: which of these carry the suffix; what the header says in global scope (several checkouts: `<dir>/` alone is ambiguous without a project name); and the markers are collected from the whole checkout, not only under the prefix (`reconcile.go:105-112`), so in global scope a marker outside `.agents/docs/` would be reported as "left out" though it hides nothing from that query. Default suggestion: project scope only, `<project>:<dir>/` in global scope, markers filtered by the prefix. Also the scope struct needs a new field; the plan only names `gitView`'s return.

G6 Items 7 and 11, script. The `--work` rule ("`<work>/build` exists and `<work>/origins.txt` does not") guards one of five things the script writes there: `fetch` also does `git -C <work>/agents-cli remote set-url origin` and `checkout --detach` on any directory that has a `.git` (`refresh-catalog.sh:125-131`), so an operator's own clone named `agents-cli` or `google-skills` in `--work` is rewritten with no `build` directory present; `<work>/writer.log` and `origins.txt` are overwritten. A positive marker (the script writes one into a directory it made; refuse a non-empty directory without it) is the same size and closes the class. Also: `work` is never passed through `phys()` today (`:57-58`), so the "resolves to `/`, home or the checkout root" check needs that call plus the same unreadable-ancestor handling the guard has (`phys` returns 1); `--work` inside `--dest` or equal to it is not addressed (the swap's `cp -R "$build/."` into a destination that contains the build). The script header (lines 17-19, 31-33) already claims `--work` is guarded ("an empty --dest or --work, /, the home directory ... are refused") while the code checks only emptiness: the plan should say the header text is corrected, not only extended. `phys()` collapse: bash keeps exactly two leading slashes and folds three or more (`cd -P //tmp && pwd -P` prints `//tmp`, `///tmp` prints `/tmp`; checked here), so the "run of slashes" wording is wider than needed but harmless; the verify fixture `//` plus an absolute path does produce the condition.

G7 Item 4, notes test. The helper the Phase 1 box calls ("the test's helper is called on fixture lines") does not exist: the check is inline in the loop (`catalog_test.go:353-383`); name it (a function taking a line, returning an error). The five words `except`, `unless`, `only`, `once`, `after` need a match rule: whole word and case-insensitive, or `read-only` and `thereafter` trip it. I read every current clause naming `--quiet` in the manifest against the proposed split (`;` and ". "): none holds a banned word, so the manifest stays green, but the claim "green on the manifest" is not stated as checked. Clause check covers only the clause that names the flag; an exception in the neighbouring sentence ("No `--quiet`. Except after approval, ...") passes: say that is accepted.

G8 Item 3, notes content. L3 is "fix" but no wording is given for either note. The scaffold note must become a refusal of `--auto-approve`/`--force` (the router forbids both whatever a note says); the filestore note must deal with an entry script that hard-codes `--quiet` (not "run without `--quiet`", which cannot be done by running the script). Wording is a security-relevant judgment left to the implementer.

G9 Item 10, scope. The fix as written ("the catalog tests give git an environment without `GIT_*`") is satisfied by the fixture helper at `catalog_test.go:443-444` alone, but the verify (`GIT_DIR=/nonexistent go test ... -run TestRefreshCatalogScript`) also reaches every `bash script` invocation (`:502, :590, :630, :646, :685, :704`), whose script runs `git clone`, `git -C ... rev-parse` and a nested `go test`; they take `os.Environ()` too. All of them must drop `GIT_*`, and the keys to drop (all `GIT_*`, or a named list) are not stated. The verify condition does hold: `git -C dir` with `GIT_DIR=/nonexistent` is refused ("not a git repository"), so it fails today as claimed. Other tests that run git (embed and search tests) are outside this item; say so.

G10 Item 5, lint. Message says "does not end or parse", but `Split` can only fail on the opening not matching or the block not closing (`frontmatter.go:20-32`); a parse error from `ParseStrict` is already refused (`inventory.go:299-302`). A BOM followed by a closed block is refused too under the rule as written (Split fails), with a message that says "does not end"; decide whether a BOM-prefixed readable block is refused (the plan's rule) and say the message names the BOM. Skills and agents already refuse a Split failure; the plan should say it checked that (not read here).

G11 Item 12 / D3. "BUG-139 stays a measured rate ... closed as fixed" is self-contradictory; and the claim that the dispatch rule "never names ROADMAP.md or TODO.md" exists was not checked (the `documentation` agent's lane sentence at `agents/documentation.md:60` exists). BUG-142 graduates to the global context's Definition of Done step 1, which does carry the frozen-tree sentence (present in the deployed context read this session); state the bug record's Graduated-table move as the evidence.

G12 Docs clauses missing or incomplete. A new skill: `skills/VENDORED.md` is listed under Surfaces but no phase box names why (a harness-owned skill is not vendored); `context/` inventory is generated, so no edit. The skill dispatches agents: `agents/` dispatch texts and the context's "Skills (N)" line are generated. The refresh script's header text for `--work` and the overrides, `google-catalog/SKILL.md` or `README` mention of the overrides (grep not run here). `LoadDocs` refusal: the project `AGENTS.md` bullet on frontmatter refusals is named; the README sentence is named; `testing/guide.md` § 2 for a lint change is not.

G13 Verification wording. E2E (a) is behaviour, fine, and its condition holds (a scratch home is suffixed because `accountHome` differs, `timer.go:62-77`). E2E (b) holds: `skills/google-catalog/catalog/.sh-agent-noindex` is tracked. E2E (c) "at least one finding with an ID, a severity and a `file:line`" is checkable by grep; "changes no tracked file outside `.agents/docs`" needs `git status --porcelain` in the clone, which the plan does not name. Phase 1's item 6 box says "the header with one, two and the root marker": the root marker is `dir == ""` in `reconcile.go:179-181` (prints `./` only if the new code maps it), so a fixture needs a marker committed at the root of a temp repo.

G14 Risks. The Antigravity assumption (loads and runs a slash-only skill) has prior evidence (`testing/2026-09-03-verify-sandbox.md:37`: agy hides slash-only skills from the model, `/init-project` expanded in print mode) that the plan does not cite; moot if G1 resolves to model-invocable.

## Scope findings (built beyond the request in `## Why`)

S1 Slash-only invocation (G1): not asked, contradicts ADR-014; no D-row says who asked. Gap.
S2 `--area`, `--reaudit [<brief>]` and the two-mode skill: the request says "a guide that does the work we have done with our reports" with "a report and action plan" as output. A re-audit mode is a second shape of the operation; D5 says the operator's words support it, which they do not name (they name a repeatable audit). Treat as a gap unless the operator confirms; `--area` is needed for the behavioural run, which is defensible.
S3 `references/areas.md`: a second file for the skill; defensible as the per-area checklist, but not named in the request. Low.
S4 Item 4's carve-out test and item 11's override restriction are items from the gate report plus the operator's "fix everything left"; in scope. Item 11 is not a gate row (L1-L10), it comes from "the script's environment-override residue" in the Why; fine.
S5 Nothing else extra: no new endpoint, flag or verb; the plan's own Out section says so.

## Questions for the author

1. BLOCKING — ADR-014: model-invocable (keep) or supersede it with the operator's approval? (G1)
2. BLOCKING — How does the skill decide it is in "the harness checkout", given the E2E runs in a clone? (G2)
3. BLOCKING — Which agents, budgets, severity scale and ID-numbering does the skill prescribe; what is `DL`; how do areas map to the six slices; defaults of `--area` omitted and `--reaudit` without a brief? (G3)
4. NON-BLOCKING — `uninstall`'s nothing-found message and which function decides it. Default: `removed`, from `RemoveTimer`, prints "no dream timer found for <home>". (G4)
5. NON-BLOCKING — notice text per OS: print the check command for the OS of the files found. (G4)
6. NON-BLOCKING — legacy suffix on every `StatusLine` return. Default: yes. (G4)
7. NON-BLOCKING, but decide before coding — search header in global scope and on `--open`; markers filtered by prefix. (G5)
8. NON-BLOCKING — `--work` rule: positive marker or the `build`/`origins.txt` pair; guard `agents-cli` and `google-skills` clones; `phys()` on work; work inside dest. (G6)
9. NON-BLOCKING — notes test helper name and whole-word match; wording of the two reworded notes. (G7, G8)
10. NON-BLOCKING — `GIT_*` dropped for every subprocess in the catalog tests, not only fixture git. (G9)
11. NON-BLOCKING — BUG-139's closing state and the check that the dispatch rule exists. (G11)
12. NON-BLOCKING — extra docs clauses: `testing/guide.md` § 2 for the lint, the script header, the overrides' mention. (G12)

## Criteria

design ready: no (G1, G2, G3 and the lesser gaps above)
implementation not ready: three BLOCKING questions (1, 2, 3)

## Not examined

Not run or read: `internal/harness/derive.go` and skill/agent loaders' refusal of a failed Split (asserted only); the global-scope code path of `search` beyond `search.go:300-330`; `index` output code; `skills/google-catalog/SKILL.md` for the overrides; the five audit-3 slice reports beyond the first lines of one; `skill-matrix.md` and `testing/guide.md` § 4 row shapes; the full text of the gate report beyond L1-L10; BUG records' contents; whether the dispatch rule in D3 exists. No code, plan or test was changed.
