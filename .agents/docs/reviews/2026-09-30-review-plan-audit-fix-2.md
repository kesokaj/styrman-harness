---
owner: reviewer
date: 2026-09-30
---
# Review — plan audit-fix, cold read round 2   (2026-09-30)

Range reviewed: `.agents/docs/plans/2026-09-30-audit-fix.md` (cold read, round 2, no diff), checked against the code it names. Round 1's 14 questions were checked against the plan as it stands.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Round 1's blockers are closed: C1 now lists a second cut list, `/goal` is consistent, IN2 states the manifest as (old still on disk) union (written), `--kind` is one kind per call, the python/anti-ui-slop grep is anchored and matches only the two directories being deleted, the tools are fetched into scratch, the `skills:` probe is scheduled before emission, and the Lint threshold is one number. What remains: (a) one real contradiction on REQUIREMENTS.md, where D29 says "approved in intent" while the global rule needs approval of the exact edit and the plan applies it in the same commit as it shows it; (b) the verbatim early-write block would overwrite the verifier's own duty; (c) two fixtures (SD7, IN3) that do not produce the condition they claim; (d) stale restatements of decisions the plan itself superseded (D2, TS1, Failure modes, Risks). Everything else is NON-BLOCKING.

## (1) Gaps

- **[high / medium]** REQUIREMENTS edits: Phase 3 line 116 applies the four edits and writes ADR-048 in one commit and "quotes the exact diff in the checkpoint report"; D29 says "approved in intent; the exact text is shown in the Phase 3 checkpoint report". Global rule: "propose the exact edit, approval, ADR, update". Approval of an intent is not approval of text that does not exist yet. Failure modes (line 62) and Risks (line 189) still say the default is "proposed, not applied; Open row stays". Three statements, two behaviours. Decide: show the exact diff to the operator BEFORE the commit (stop at that task), or record that "2 yes" covers the wording the implementer writes.
- **[medium / high]** Early-write block (line 45, AG10): "the 13-line block in `agents/reviewer.md` is the canonical text, copied verbatim into security, verifier, explorer and dreamer (their variants deleted)". The reviewer block says "Write your report file ... the range line". `agents/verifier.md:12-20` carries the verifier's own duty: write the log after the FIRST verification command and update it after EVERY command; the global AGENTS.md repeats that as the verifier-specific dispatch text. A verbatim reviewer copy deletes it. The verify clause ("5 identical blocks, diff pairwise empty") would pass while the verifier rule is gone. Fix: keep one shared core plus a one-line per-agent trigger ("verifier: after the first command and after every command"), and change the verify to pin both.
- **[medium / medium]** SD7 fixture (line 92): "a fixture that truncates the -wal after the db is copied" does not produce the condition the retry is keyed on. `isUnreadable` (`internal/search/store.go:94`) is SQLITE_CORRUPT / NOTADB on OPEN; SQLite treats a short or invalid -wal as empty (frame checksums end the log) and opens the older content without error, so the retry never fires and the test cannot fail. `copyDB` (`internal/search/dream.go:242`) is one synchronous loop over three files, so "after the db is copied" needs a test seam the plan does not name, and the one-second wait needs an injectable sleep or the test takes a second per run. State: where the retry wraps (copy plus open, in the three callers at dream.go:60/117/203, not inside copyDB), the seam, and a fixture that demonstrably yields CORRUPT (e.g. a db copy from mid-checkpoint, or a db file whose page 1 is torn), proved to hold before the test is trusted.
- **[medium / medium]** IN3 (line 87): `apply.go:55-57` reads the destination with `os.ReadFile` (follows a link) and calls `os.Stat` and `os.Chmod` (both follow). The defect only appears when the linked target's bytes equal the planned ones AND its mode differs from `f.Mode`; a fixture with equal mode leaves "link untouched" true without the fix. The plan also says "a `keep` line", but the existing `keep` text (apply.go:38) claims "not written by sh-agent; move only rules meant for every machine into rules/", which is wrong for a link. Decide the line text and whether a linked dest enters the manifest (if it does, `uninstall` removes the operator's link at apply.go:290).
- **[medium / medium]** Stale restatements of decisions the plan itself supersedes: D2 (line 194) still reads "target <= 14,000 ... a warning test at 15,000" while Interfaces and Phase 3 say 14,500 and a FAILING test above 15,000 (D2a supersedes only the 18,000); Phase 6 TS1 (line 152) says "LICENSE added; VENDORED row updated" against D8a and line 50 ("no licence owed, VENDORED row removed"); D16 still says "Open below" though D28 decided it; D5's Budget prose line against D5a. An implementer following the task text builds the superseded design.
- **[medium / medium]** C1 arithmetic still does not reach the target by the plan's own numbers. The second list's parentheses sum to 2,950, not the stated ≈ 3,150; first list 4,600; total 7,550. Needed 21,906 - 14,500 = 7,406, before adding back the G11 paragraph (≈ 250), seven skill names (≈ 150) and the resume-line clause (≈ 100): net ≈ 14,850. The Phase 3 and Phase 7 verify "<= 14,500" will fail while the 15,000 pin passes. Fix: say which it is (ceiling 15,000 with 14,500 a goal), or name one more whole rule to leave.
- **[medium / medium]** AGENTS.md § Testing requires a behavioural run for any change to `agents/*.md`. Phase 4 runs changelog and verifier only (5 each); the plan also rewrites reviewer, security, explorer, dreamer, bugs and documentation (the early-write block, lane edits, records lines, ADR removal) with no run and no stated exemption. Likewise C1 removes whole rules from the always-loaded context but only the G11 paragraph gets an `--append-system-prompt` run; nothing checks the agent still follows a moved rule (for example the budget-and-artifact clause). Name the exemption and its reason, or add the runs.
- **[medium / low]** IN4 (line 88): `update` returns on `checkOnly` both when up to date (main.go:696) and when behind (main.go:701), before the dirty check (:711). Moving the refusal above both changes the documented `update --check` contract (README:162 and :183: "only reports, exit 2 when behind") for a dirty tree: it now exits 1. The plan names `--dry-run` and the code's parameter is `checkOnly`; it does not say whether the refusal applies to `--check` too, and gives no docs clause for README:183 or `specs/sh-agent/*.md`.
- **[low / medium]** Kit and guide wiring: `TestRenderedContextUnderTarget` is a new test, but guide § 2 (`.agents/docs/testing/guide.md:20-28`) and AGENTS.md Setup commands name only the 22,000 pin; no clause updates them. `kit-preamble.sh` inputs `KIT_ASSERT_GREP` / `KIT_ASSERT_COMMIT` have no format (one line each? several?), and the preamble is required in Phase 4 but the Phase 3 behavioural run does not source it or say the scratch project lies outside every checkout (BUG-106; the plan's Failure modes line 60 says it does).
- **[low / medium]** Phase 5 verify names skill-matrix rows `review-docs` and `review-resume`; `testing/skill-matrix.md` has neither, and no task says to write them or on what fixture. Same for rows of the seven new skills (Phase 7 "row PASS"). AGENTS.md § Testing says matrix rows are added; say who writes them and the fixture (demo HEAD, as the `distill`/`handoff` rows use).
- **[low / medium]** `explain-project` writes `docs/onboarding.md`: no surface in the plan or spec owns the file (the documentation agent's lane lists README, specs, guide, deploy guides), and `grep -rn onboarding` over specs, README and init-project finds nothing. Add to the lane or say it is the skill's own artifact.
- **[low / medium]** Interfaces line 39 still carries thinking-out-loud: "`skills/project-memory`? no: to the search spec pointer". The decision is recoverable (search spec pointer) but the text should read as one.
- **[low / low]** `restatement-grep.sh`: the plan leaves the exec bit and the interpreter unspecified. `plan.go:138` preserves the mode `&0o755|0o600`, so a non-executable git mode installs as non-executable; the skill should call it as `sh <path>` or the file be committed 0755. Also `git ls-files` lists tracked files only, so a file created in the same session is invisible until staged; say so in the script's header.
- **[low / low]** TR1 "expect: 46 rows": the Graduated table today has 46 lines starting `|` including header and separator (44 data rows); after two moves the count is 46 data rows or 48 lines. State which.
- **[low / low]** `AG6` verifier `wc -w < 1,400` from 2,053 (a 32% cut) while adding nine numbered steps and the 13-line block; achievable but the word count is a proxy; the behavioural run (step marker) is the real check.
- **[low / low]** Risks line 189 asserts the 971-leaf corpus exists "as DB copies under `~/.gemini/antigravity-cli/conversations`"; I confirmed 97 entries exist there now, not 971 leaves (leaves are counted by the measurement, not files), so this is consistent, but the Phase 1 verify should say "N recorded", as SD9 does.

## (2) Scope — built beyond the request in `## Why`

The request: fix every audit finding, add features from the audit, research first, phases, verify, PRIORITY_TODO, re-audit. Each item below is either covered by an audit ID or has a D-row; the ones with weak cover:

- `TestRenderedContextUnderTarget` (second pin): D2 "implementer". Acceptable as the verification of C1, but it duplicates the `render | wc -c` check; say it stays as the regression pin. Low.
- README template in init-project assets (Phase 7 line 172): no ID or D-row. Name the G-item or strike.
- `explain-project` (D31) and the `--resume` lens (D14): D-rows by "implementer" cite audit gaps G9 and G12; covered by "add new features from your audit report". No finding.
- `lock_other.go` stub and a flock around `Install` and `UpdateFromRelease` (IN8): covered by the audit ID IN8.
- Nothing else found beyond the request: one Lint threshold, one `skills:` mechanism after the probe, no Budget prose line.

## (3) Readiness — questions I would have to ask

1. BLOCKING. REQUIREMENTS edits: show the exact diff to the operator and stop before committing, or commit on "2 yes"? (Failure modes / Risks say "not applied"; Phase 3 says applied.)
2. BLOCKING. Early-write block: keep the verifier's "first command / every command" rule as a per-agent line, or accept the reviewer wording verbatim?
3. BLOCKING. SD7: what fixture produces SQLITE_CORRUPT or NOTADB on the copy, where does the retry wrap, and what is the injectable sleep? (Two implementers would build different retry seams and the plan's fixture passes without the fix.)
4. BLOCKING. IN3: what is the linked-destination fixture's mode, what does the line say, and does the link enter the manifest?
5. NON-BLOCKING. C1: is 14,500 a hard verify or a goal under the 15,000 pin? (Default: hard pin 15,000; report the rest.)
6. NON-BLOCKING. Which of D2 (14,000), TS1 (LICENSE added), D16 (Open), D5 (Budget line) is stale? (Default: the later D-row wins, per the supersedes lines.)
7. NON-BLOCKING. IN4: does the dirty refusal apply to `update --check` (exit 1) or only to `--dry-run`?
8. NON-BLOCKING. Behavioural runs for the other agents and the moved context rules: exempt or run?
9. NON-BLOCKING. Skill-matrix rows and fixtures for `review-docs`, `review-resume` and the seven new skills.
10. NON-BLOCKING. Script exec bit, `kit-preamble.sh` input format, where `docs/onboarding.md` is owned, what "46 rows" counts.

## Round 1 questions: status

1 C1 arithmetic: answered (second list), residual gap above. 2 `/goal`: answered (kept, C15). 3 IN2: answered. 4 `--kind`: answered. 5 grep: answered (anchored; matches only the dirs being deleted). 6 AG12b: answered (one threshold 1,024; nothing refused today; vendored fork rule). Longest first-line description measured here is 775 (`skills/review`), the plan's 963 is cloud-run-deploy's multi-line block; `Lint` must join continuation lines, which the plan does not say (NON-BLOCKING, a reasonable implementer does it). 7 `skills:` probe: answered (Phase 4, n/3, before emission). 8 tools: answered. 9 threshold: answered. 10 SD7/IN8/SD9: SD7 answered in text but the fixture fails (above); IN8 answered; SD9 answered. 11 content: answered (AG9 sentence, early-write block, step marker); grilling note answered (D27). 12 links test / restatement script: answered, exec-bit gap above. 13 docs clauses and homes: mostly answered; IN4 and guide § 2 gaps above. 14 re-audit publish: answered (reports are the record, fallback stated).

## Criteria (the plan as a design document)

1. Dead code — pass: strikes python and anti-ui-slop, dead pointers, samber file; `TestSkillLinksResolve` walks non-vendored skills only.
2. Flow tracing — fail: SD7 fixture and IN3 fixture cannot fail as written; IN4 changes the `--check` contract without tracing README:162/183.
3. Standards conformance — fail: stale restatements (D2, TS1, D16, Failure modes/Risks vs D29); REQUIREMENTS approval path; behavioural runs for changed agents; docs clauses for guide § 2 and IN4.
4. Directory layout — pass: `skills/review/scripts/`, `.agents/docs/testing/kit-preamble.sh`, `assets/agents-tree/docs/specs/TEMPLATE.md`; `docs/onboarding.md` ownership unstated.
5. Simplicity (DIY-except-security) — pass: stdlib flock, no new dependency; `api-auth-gcp` names vetted libraries and the `iss` check, no hand-rolled crypto.
6. Decoupling — pass: phases touch separate files; Phase 3 and Phase 4 both edit `agents/verifier.md`, `bugs.md`, `documentation.md` in sequence (stated by order).

## Not examined

The recon brief and the five research briefs were read only where needed for the cited facts (the plan's technical claims about Terraform, Cloud Run flags, Tailwind, Vitest and Go release versions were not verified). The "measured 963" cloud-run-deploy description length and the claim that no vendored description exceeds 1,024 were not reproduced (my measurement reads the first line only). The second-list byte estimates were summed, not re-measured line by line. `skills/preflight` restatement of the Ship/deploy paragraph was confirmed present only at `SKILL.md:29`; the version-number rule's target was not read. Phase 6 and Phase 7 technical content (D20-D26) was not assessed.

design needs revision
implementation not ready
