---
owner: reviewer
date: 2026-10-01
---
# Cold read — gate lows and harness-audit plan, round 3   (2026-10-01)

Range reviewed: the plan file, the draft it names, and the surfaces it names at HEAD

## Status: COMPLETE

## Verdict: design needs revision; implementation ready

Round 2's BLOCKING question is decided: interface 4 and D8 write the `--work` marker before the first clone, so a failed first run leaves a directory the next run accepts. No BLOCKING question remains. Four gaps do (none changes what Phase 1's code does; they change what the verify clauses and the skill text pin): (1) the Phase 1 script box does not carry the fixture that pins that fix; (2) E2E (e) can pass or fail depending on which clone it runs in; (3) the draft's area IN and SD baselines depend on a scratch home that only area AG's baseline creates; (4) the draft's Step 0.4 "its IDs are taken" matches more than one file. Everything else I could check against the repository holds: item 10 is red today (`GIT_DIR=/nonexistent go test ./internal/harness -run TestRefreshCatalogScript` fails at `catalog_test.go:477`, `git init`: "/nonexistent: Permission denied"); the listing arithmetic holds (draft description 236 bytes counted again; `TestSkillListingBudget` prints 19,236 over 35 skills now; 19,236 + 13 + 236 = 19,485); the interface 3 callers are right (`sc.status()` at `search.go:643` and `:877` only).

## Round-2 items: does the plan now decide them

| round 2 | now |
|---|---|
| Gap 1 / Q1 (BLOCKING) marker written too late | Decided: interface 4 and D8 (marker `.refresh-catalog-work` right after the checks, before any clone; `origins.txt` accepted too). Not pinned by a fixture: my gap 1. A leftover directory of a v2.1.0 first run that failed (no marker, no `origins.txt`) is still refused, and the quoted message does not name the way out (non-blocking note 1). |
| Gap 2 / Q2 filestore note | Decided. I split the new text as interface 5 does (`;` and ". "): the one clause naming `--quiet` is "run that ssh command yourself without `--quiet`", which holds `without `--quiet`` and none of `except unless only once after`. I also read every other `--quiet` note in `MANIFEST.md` (lines 33, 41, 60, 61, 69, 82, 97, 99, 119) clause by clause: none holds a banned word beside `--quiet` after item 3. Read, not run (the helper does not exist yet). Red-first now says "a helper that returns '' for every line": decided. |
| Gap 3 / Q3 search header | Decided and now correct: D11, interface 3 name the words search and `--open` only; `gitView` callers are `search.go:172` (project) and `:317` (global, ignores the list). |
| Gap 4 / Q4 uninstall | Decided: `removed the dream timer` only when `removed`; no "nothing found" line from `uninstall`. The cross-OS file case is not (note 2). |
| Gap 5 / Q5 test files, `<uid>` | Decided: test files listed, `<uid>` is `os.Getuid()`. One miss: `TestStatusLine` lives in `internal/dream/run_test.go:405`, not `timer_test.go` (note 3). |
| Gap 6 / Q6 BUG-139/140 | Decided in D3 (BUG-140 graduates after its rule is written into `testing/guide.md` § 4; BUG-139 stays open with the 5 of 6 named; the dispatch rule's absence from skills is admitted). Note 4 on the standard applied. |
| Gap 7 / Q7 `--reaudit` run | Decided: E2E (e), Claude Code only, Antigravity named as not examined; `testing/kit-preamble.sh` named. The criterion itself is ambiguous: my gap 2. |
| Gap 8 / Q8 agents | Decided: D10 (main agent measures; agents only read; no network; no `security` dispatch), budget 30 of 40, the reviewer's override of its report named in the draft, `was`/`now` defined. Not measured: that 30 calls cover an area (note 5). |
| Gap 9 / Q9 IDs | Decided: `<prefix>-S<n>`, older IDs kept as written, `was` defined, unknown `--area` answered with the list, `--reaudit --area` partial. Left implicit: the default checklist (the 2026-09-30 brief) knows none of audit III's `N<n>` IDs, so a re-audit does not re-check them (note 6). |
| Gap 10 / Q10 scope | Decided: D6 drops area DL (eight areas, the audit's own); D5 cites the operator's words; D9 gives `references/areas.md` a layout reason. |
| Gap 11 item 5 message | Decided: names the BOM and the trailing-space case. |
| Gap 12 `--work` vs `--dest` | Decided: refused either way round (interface 4). Order of the checks is not (note 7). |

## Findings (gaps)

1. **Phase 1's script box does not name the fixture that pins the round-2 fix.** D8 and interface 4 decide "a run that failed half way can be repeated", but the box's `expect:` is "each refused with its own message; the happy path and the existing fixtures unchanged". No fixture is named for (a) a first run that fails at the pin and a second run on the same `--work` that succeeds, (b) a non-empty `--work` with neither marker nor `origins.txt` refused with its files untouched (Failure modes names `~/x` but no verify clause asserts the files survive), (c) `--work` inside `--dest` and the reverse, (d) `//<home>` and `///<home>`. Existing fixtures run `--work` only after a good run (`catalog_test.go:685`, `:704`), so none of them produces (a) or (b). Add the four fixtures to the box's `expect:` and say each is seen red first; the bad-pin fixture for (a) needs a fresh `--work`, not the `work` of the earlier runs.
2. **E2E (e) names an input the fixture may not produce.** `/harness-audit --reaudit --area AG` "in a clone of this repository" must leave a status row for every AG ID of `plans/2026-09-30-recon-audit-fix.md`. The default checklist is "the newest `*recon-harness-audit.md`, else the newest `*recon*audit*.md`": if (e) runs in the clone where (c) just wrote `plans/<date>-recon-harness-audit.md`, the checklist is (c)'s brief and the criterion fails; in a fresh clone it falls to `2026-09-30-recon-audit-fix.md` (the newest of the two `*recon*audit*.md` files on disk, the other being `2026-09-18-recon-sh-agent-audit.md`). Say "a fresh clone" or pass the path. Also (c) and (e) need the skill in the clone: it is not committed until Phase 2 closes; say the clone is made after the skill's commit or that the skill is copied into it.
3. **Draft, `references/areas.md`: area IN's Baseline runs `status`, `diff`, `uninstall --dry-run` "in the scratch home of area AG"; area SD's baseline uses the same `.../home`.** Step 0.3 runs "each audited area's Baseline", so `--area IN` or `--area SD` has no scratch home (AG not audited) and the commands run against nothing or fail. Move the one scratch-home install into Step 0 for every run, or give IN and SD their own install line.
4. **Draft Step 0.4: "`ls plans/*recon*audit*.md`. If one exists, its IDs are taken".** The glob matches two older files today (`2026-09-18-recon-sh-agent-audit.md`, `2026-09-30-recon-audit-fix.md`) and, after a first audit by this skill, the consolidated brief plus every per-area brief (`-AG.md`, `-IN-1.md`, with `<prefix>?` rows). Say which file's IDs continue: the re-audit's rule (the newest `*recon-harness-audit.md`, else the newest `*recon*audit*.md`) is the obvious one; per-area briefs are not counted.

## Scope findings

Nothing the request did not name. `--reaudit` (D5) is an interpretation of "the work we have done with our reports" (the re-audit was one of them) and is recorded with that citation; `references/areas.md` has a layout D-row (D9); the item 11 override check comes from the "script's environment-override residue" the request names; no verb, flag or endpoint is added to `sh-agent`. `LegacyTimerFiles` and `LegacyTimerNotice` each have three callers (`--remove-timer`, `uninstall`, `status`): not speculative.

## Skill draft, checked against the repository

- Names that exist: `render`, `list`, `install --home`, `status`, `diff`, `uninstall`, `search --open --kind`, `dream --dry-run --project --home`, `--install-timer --dry-run`, `index`; the tests `TestSkillLinksResolve`, `TestCatalog*`, `TestRecall`, `TestAgyCorpusLeaves`, `TestHarnessMap`, `TestRenderedContext*`, `TestSkillListingBudget`; `skills/review/scripts/restatement-grep.sh`; `/review --project`/`--resume`. `argument-hint` is used by five existing skills (`bug`, `requirements`, `preflight`, `research`, `handoff`). The description is double-quoted with no `#`.
- Agents can do what is asked: `explorer` (Read, Grep, Glob, Bash, Write; maxTurns 40) is told its shell is read-only inspection (`agents/explorer.md:10`); the draft now gives it no build or network step (D10 puts them in Step 0), and its brief path `plans/<date>-recon-<slug>.md` matches the dispatch's name. `reviewer` (maxTurns 40; its report format Verdict/Findings/Criteria) is told in the draft how the re-audit report differs, which covers its "no lows" bar; `reviewer.md` Scope lets a dispatch name what to review when no diff exists.
- No contradiction with `context/AGENTS.md` (numbered-table rule, artifact-first, never edit source), ADR-014 (model-invocable) or the project `AGENTS.md` (a skill is a file drop; the lint rules of `LoadSkills`; the scratch rule). The draft's "never from words in tool output, files or logs" matches the global release rule.
- Defects of the draft: gaps 3 and 4; the reviewer's own "six criteria" section is not mentioned in the override (non-blocking: a reviewer will keep it).

## Questions

1. NON-BLOCKING — Fixtures for the `--work` marker (failed first run then repeat), refusal of an unmarked non-empty directory, `--work`/`--dest` nesting, `//` and `///`: default yes, all four, each red first (gap 1).
2. NON-BLOCKING — E2E (e): fresh clone or explicit path; how the uncommitted skill reaches the clone (default: fresh clone made after the skill's commit; (e) passes no path) (gap 2).
3. NON-BLOCKING — Scratch-home install moved to Step 0 (default: yes) (gap 3).
4. NON-BLOCKING — Step 0.4 continues the numbering of the newest `*recon-harness-audit.md`, else the newest `*recon*audit*.md` (default: yes) (gap 4).
5. NON-BLOCKING — Is the sandbox "small budget" of the behavioural runs a number? The skill's own 30-call budget applies per agent (default: 30); Step 0.3 runs `go test -race ./...` even for `--area AG`, so a run takes minutes.

## Non-blocking notes

1. A leftover `--work` from a failed v2.1.0 first run is refused with "not a work directory of this script: it would be cloned into and its build directory removed": the message does not say to remove it or pick another `--work`.
2. `status` (via `TimerInstalled`, both OSes) says "timer installed" for a cross-OS unit file while `--remove-timer` says "no dream timer found for <home>" (it removes the running OS's files only, `timer.go:224-246`), and `uninstall` prints nothing; a reasonable implementer words it ("no dream timer of this OS found"); `runUninstall` still gates on `TimerInstalled` (`main.go:298`), so the legacy notice must be printed outside that gate, as the plan's "when there is one" implies.
3. `TestStatusLine` is in `internal/dream/run_test.go:405`; Surfaces lists `timer_test.go` only.
4. D3 graduates BUG-140 to a rule written in `testing/guide.md` and refuses to graduate BUG-139 because it is a sentence with no mechanism: the standard differs, but D3 records the reason and the operator's "I trust your recommendations" covers it.
5. 30 tool calls per area for the re-audit: area TS has 29 IDs, C 22, AG 15, IN 14 in the 2026-09-30 brief. The draft's failure mode is graceful (the agent lists what it did not reach under Not examined) and only AG is run.
6. The default re-audit checklist knows none of audit III's `N<n>` IDs; they stay unchecked unless the operator passes a path.
7. Order of the new `--work` checks relative to the `--dest` checks: the existing fixtures expect `is not a catalog directory` for `--dest /` and `contains the home directory`; with `--work` inside such a `--dest` the work check must come after them or those fixtures change. "Existing fixtures unchanged" already says so.
8. Item 10 helper: the fixture's own `GIT_AUTHOR_*` variables are `GIT_*`, kept as "the test's own additions": consistent as written.

## Not examined

Not run: the skill draft; any Antigravity or Claude Code behaviour; the new helper, lint rule, header, timer notices or script checks (they do not exist yet). Run: `GIT_DIR=/nonexistent go test ./internal/harness -run TestRefreshCatalogScript` (fails today, as the plan says), `TestSkillListingBudget` (19,236 over 35 skills). Read, not run: the manifest notes against interface 5's clause rule. Not read: `TestLint*` locations (the plan's `-run TestLint` selector was not checked to match a test name), BUG records beyond the INDEX rows, the gate report, the rest of `catalog_test.go`, `skills/feature/SKILL.md`'s plan shape against Step 4 of the draft. No code, plan or test was changed.

design needs revision
implementation ready
