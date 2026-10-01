---
owner: reviewer
date: 2026-10-01
---
# Cold read — gate lows and harness-audit plan, round 2   (2026-10-01)

Range reviewed: the plan file, the draft it names, and the surfaces it names at HEAD

## Status: COMPLETE

## Verdict: design needs revision; implementation not ready

One BLOCKING question and eleven non-blocking gaps. The worst: the `--work` guard (interface 4, item 7) refuses a work directory that holds no `origins.txt`, and the script writes `origins.txt` only AFTER both clones (`refresh-catalog.sh:143-149`: `fetch` x2, then `: > "$origins"`). A first run that fails in `fetch` (mistyped pin: "pin ... does not resolve", no network, a refused clone) leaves `agents-cli/` and `google-skills/` in the work directory and no `origins.txt`; the corrected re-run is then refused "not a work directory of this script". The plan asserts "the default work directory and any earlier one pass" and no fixture covers a failed first run (the bad-pin fixture, `catalog_test.go:779`, runs after a good run, so `origins.txt` exists). Phase 1 is otherwise implementable. The skill draft is coherent and every path, test, verb and agent it names exists; its listing arithmetic holds (236-byte description, 19,236 + 13 + 236 = 19,485 of 19,500, checked against `TestSkillListingBudget` output of 19,236 over 35 skills); the gaps in it are about agents' rules, not missing files.

## Round-1 items

| round 1 | status in the plan now |
|---|---|
| G1 (BLOCKING) ADR-014 / slash-only | Decided: D4 makes it model-invocable, cites ADR-014, budget arithmetic checked. The draft carries no `disable-model-invocation`. |
| G2 (BLOCKING) detection of a harness checkout | Decided: D7, tree shape at the git root; consistent with E2E (c) and the draft's Usage. |
| G3 (BLOCKING) skill procedure | Mostly decided: agents (explorer per area, security for IN/SD, reviewer for reaudit), parallel, severity scale, IDs, DL (D6), defaults of `--area` and `--reaudit`, `<n>`, one report per area, the tally file. Still open: a numeric tool-call budget; ID continuation versus audit III's `N<n>` IDs (gap 9); `S<n>` scope (gap 9). |
| G4 uninstall text, notice per OS, StatusLine, dry-run line | Decided except `uninstall` (gap 4): the plan says "lines unchanged", while item 1 still says `uninstall` does not say when it found nothing; the `TimerInstalled` (both OSes) versus `removed` (running OS) mismatch at `main.go:296-304` is not decided. Test files still missing from Surfaces. |
| G5 search header | Decided scope (project only, global ignores, words/--open). But the plan's claim about which commands show it is wrong (gap 3). |
| G6 `--work` | Partly: marker (`origins.txt`, D8) and `phys()` decided; marker timing is the BLOCKING gap (gap 1); `--work` inside `--dest` not addressed. |
| G7 notes helper | Decided: `quietNoteProblem`, clause split, banned words. Whole-word match is stated; case-insensitivity is not. New contradiction with item 6 (gap 2). |
| G8 notes wording | Decided by quoted text; the filestore text fails interface 5 (gap 2). |
| G9 `GIT_*` scope | Decided: one helper for every subprocess in the catalog tests. |
| G10 lint message/BOM | Decided enough: rule and message given; I checked the claim that skills and agents already refuse a `Split` failure (`inventory.go:185-187`, `:236-238`: true). `Split` itself refuses a BOM or a trailing-space fence (`frontmatter.go:24-26`), so the fixtures produce the condition. |
| G11 BUG-139/140/142 | Not decided (gap 6). |
| G12 docs clauses | Mostly: `testing/guide.md` § 2, script header, project `AGENTS.md`, README added. Still unnamed: where the `CATALOG_*_REPO` restriction is described outside the script (grep not run by the plan). |
| G13 verification wording | Partly: (c)'s "no tracked file changed outside `.agents/docs`" still names no command (`git status --porcelain` in the clone); NON-BLOCKING. |
| G14 Antigravity risk | Moot (model-invocable). |
| S1 | Moot. S2 (`--reaudit`): D5 cites the operator's words; those words support a repeatable audit, not a re-check mode; accepted as an interpretation, flagged. S3: `references/areas.md` has no D-row; low. |

## Findings (gaps)

1. **[BLOCKING] `--work` marker is written too late.** See Verdict. `refresh-catalog.sh:143-149`. Decide: write the marker (name and content) before the first `fetch`, and say what happens to a directory left by a failed run under v2.1.0 (no marker: refused, with the message naming the fix); add a fixture "first run fails at the pin, second run succeeds on the same `--work`" to Phase 1's script box.
2. **Interface 6's filestore note fails interface 5's test as quoted.** "Use the jump-host engine only after confirming OS Login, and never as it is: `jump_host_engine.py` passes `--quiet` … without `--quiet` …" is one clause (the split is `;` and ". ", not `:`) that names `--quiet` and holds `only` and `after`. The plan says the manifest passes "after item 3's two rewordings (checked when the test is written)" and that a failing note is reworded, so a default exists; the quoted text is nonetheless stated as the decided text and is wrong. Also: "red on a fixture line … before the change" cannot be shown while the helper does not exist (a build failure is no evidence, `review-standards` § 2): say the helper is first written with today's whole-line rule.
3. **Interface 3 misdescribes existing code.** "the words search, `--refs` and `--open` all show it, as does `sh-agent index`": `sc.status()` has two callers only, `search.go:643` (words) and `:877` (open). `--refs` calls no `status()`, and `Index` (`search.go:386-406`) prints a fixed key/value list with no header. Decide: which output carries the notice, and its shape for `--refs` and `index` (a new `left out` line?), or strike them from the claim; pinned index output and docs (`specs/sh-agent/search.md`) follow.
4. **Item 1, `uninstall`.** The table row says `uninstall` does not say when it found nothing; interface 1 leaves its lines unchanged. `runUninstall` gates on `TimerInstalled` (both OSes) and prints "removed the dream timer" after `RemoveTimer`, which removes the running OS's files only (`timer.go:224-246`): on a cross-OS file it claims a removal that did not happen. Decide: print on `removed`, and whether a "nothing found" line is wanted (the row says yes).
5. **Surfaces touched omits the test files** the signature change and the new behaviours edit (`internal/dream/timer_test.go`, `cmd/sh-agent/main_test.go`, `internal/harness/harness_test.go`, `internal/search/*_test.go`). `<uid>` in the `launchctl` text: whose (the running user's, as `RemoveTimer` uses `os.Getuid()`)? Default obvious; say it.
6. **D3 evidence not checked.** BUG-139: "closed as fixed by its two mechanisms (the lane sentence, and the dispatch rule that never names `ROADMAP.md` or `TODO.md`)" and "measured 5 of 6" stays: a rule that held 5 of 6 is not fixed; the dispatch rule exists only as the INDEX row's prevention text (`bugs/INDEX.md:98`), I found it in no skill, agent or context file (`agents/documentation.md:49,60` has the lane sentence only). BUG-140: "every kit separates the run from the check" is asserted; `testing/kit-preamble.sh` has nothing for it. Decide: closed with the residual named, or graduate the dispatch rule into a mechanism first.
7. **Phase 2 has no verification of `--reaudit`.** E2E (c) and (d) run the first-audit mode and the non-harness refusal only; the mode with the most new text (dispatch of `reviewer`, `<n>` counting, the tally) is never run. Add one run (`--reaudit --area AG`) or record it as not examined in the skill-matrix row. The sandbox runs are also subject to `testing/kit-preamble.sh` (project `AGENTS.md` § Testing); the plan does not name it.
8. **Skill draft against the agents it dispatches.**
   - `explorer` (tools Read, Grep, Glob, Bash, Write; maxTurns 40) is told by its own definition that Bash is "for read-only inspection … never for anything that mutates state" (`agents/explorer.md:10`); `areas.md` asks it for `go run … install --home .agents/scratch/<name>`, `go test -race ./...`, `govulncheck`, `git archive | tar -x`, a zip dry-run update, `SH_AGENT_RECALL=1 go test`. Scratch is allowed (`explorer.md` Records) but the contradiction stands; say the dispatch permits scratch writes.
   - Network: the draft's Step 1 forbids it "unless the area's block says to check a fact against a vendor's documentation"; DL's `gh release view` and TS's `git archive <pin>` diff against upstream need it, and TS says "where the area's dispatch allows the network" (circular). Default to state: no network, those measures are "not examined".
   - `reviewer` (`agents/reviewer.md`): its bar excludes "a wart … nobody meets yet" (the draft's `low`) and its report format is Verdict/Findings/Criteria; the draft's `### Status of every ID` shape must be said to replace it, and the dispatch must ask for lows.
   - No numeric tool-call budget; area IN is 42 files and about 17,000 lines of Go plus nine commands to run, against `maxTurns` 40.
   - "last security review … older than the last release": which files, which command (`sh-agent search --refs <path> --kind review`?) is undecided.
9. **IDs.** Audit III's new findings are `N<n>` (13 of them, `reviews/2026-10-01-audit-3-*.md`), not "new IDs continuing the area's numbering" as the draft's re-audit and Step 0.4 say; the default `--reaudit` checklist (`*recon*audit*.md`, today `plans/2026-09-30-recon-audit-fix.md`) knows none of those or any of audit III's statuses. `S<n>` "inside their area" collides across areas (`S1` in AG and IN) in the consolidated file; the old recon has no `S` rows. `was` in `| ID | was | now | evidence |` is not defined. `--area` with an unknown prefix, and `--reaudit --area`, are undefined.
10. **Scope: D6 (area DL) is "by: implementer"**; the request names "the work we have done with our reports", which had eight areas. It is defensible (three gate findings were in the release path) but per the rule a D-row must say who asked; the operator's confirmation is one line. `references/areas.md` has no D-row (low).
11. **Item 5 interface** is sound; low: the message says "does not end or parse" for a BOM-prefixed closed block, which is neither; name the BOM or trailing-space case in it.
12. **`--work` inside `--dest` (or the reverse)** is still unaddressed; the marker rule covers most of it. Low.

## Draft checks (against the repository)

- Paths, files, verbs: `render`, `list`, `install --home`, `status`, `diff`, `uninstall`, `search --open --kind`, `dream --dry-run --project`, `update`, `index` all exist (`--help`). Tests exist: `TestSkillLinksResolve`, `TestCatalog*`, `TestRecall`, `TestAgyCorpusLeaves`, `TestHarnessMap` (`cmd/sh-agent/main_test.go:1335`), `TestRenderedContext*`, `TestSkillListingBudget`. Files exist: `internal/harness/{release,zip,repo,report}.go`, `embed.go` (at the repository root), `deploy/{local,remote,gcp}.md`, `specs/sh-agent/{search,memory-loop,dream,delivery}.md`, `skills/init-project/assets/`, `skills/review/scripts/restatement-grep.sh` (its `-i` and `--` options match the draft's use), `.gemini/config/agents/<agent>/agent.md` (`plan.go:163`). `/review --project` and `--resume` exist (`skills/review/SKILL.md:19,27`).
- Skill coverage: the WS, TS and `google-catalog` lists hold exactly the 35 existing skills plus `harness-audit` (checked name by name).
- Agents: `explorer` has Read, Grep, Glob, Bash, Write, maxTurns 40 (`agents/explorer.md:4-7`): it has what the draft asks, subject to gap 8. `security` writes `reviews/<date>-security-<slug>.md`, matching the draft's name.
- Links: `TestSkillLinksResolve` regex only checks `references/…`-style backticks and markdown links; `references/areas.md` resolves; no draft path trips it.
- Description: 236 bytes (counted), name 13; 35 skills, 19,236 today.
- Contradictions with `context/AGENTS.md`, ADR-014 or project `AGENTS.md`: none found. Adding a skill stays a file drop. The numbered-table output matches the context's rule.

## Scope findings

- `--reaudit` (D5): a second mode; D-row cites words that do not name it; accepted as interpretation, say so to the operator in one line.
- Area DL (D6): see gap 10.
- `security` beside IN and SD: part of the same audit; no D-row; low.
- `references/areas.md`: no D-row; low.
- Nothing else extra: no endpoint, verb or flag on `sh-agent`.

## Questions

1. BLOCKING — Where and when is the `--work` marker written, so a failed first run does not poison the directory (gap 1)?
2. NON-BLOCKING — Final wording of the filestore note so it passes interface 5 (default: split at a `;`, drop `only after`) (gap 2).
3. NON-BLOCKING — Which of `--refs` and `index` carry the left-out notice, and in what shape (default: words and `--open` only, strike the other two claims) (gap 3).
4. NON-BLOCKING — `uninstall`: message from `removed`? A "nothing found" line? (default: print on `removed`; none for nothing) (gap 4).
5. NON-BLOCKING — Test files in Surfaces; `<uid>` (default: yes; `os.Getuid()`) (gap 5).
6. NON-BLOCKING — BUG-139 and BUG-140 closing state (default: close 140, keep 139 open with the 5 of 6 named) (gap 6).
7. NON-BLOCKING — A `--reaudit` run, or "not examined" in the matrix (default: one run on `--area AG`) (gap 7).
8. NON-BLOCKING — Explorer scratch writes; network default; reviewer's report shape and lows; numeric budget (default 30 of 40); security-age rule (gap 8).
9. NON-BLOCKING — ID rules vs `N<n>`, `S<n>`, `was`, unknown `--area` (gap 9).
10. NON-BLOCKING — Operator's one-line confirmation of DL and `--reaudit` (gap 10).

## Not examined

The skill draft was not run; no Antigravity or Claude Code behaviour of the skill. `go test` was run only for `TestSkillListingBudget` (cached). Not read: the full gate report, BUG records beyond the INDEX rows, `skills/feature/SKILL.md`'s plan shape against the draft's Step 4, README counts of skills (grep for the count wording found none), `google-catalog/SKILL.md` for the overrides, `TestHarnessMap`'s body. The "143 findings, 113 fixed" figures in the draft were not recounted. No code, plan or test was changed.

design needs revision
implementation not ready
