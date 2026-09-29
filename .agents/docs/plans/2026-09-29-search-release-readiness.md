---
owner: main
date: 2026-09-29
---
# Plan — `sh-agent search`: release readiness   (2026-09-29)

## Why

The operator, 2026-09-29, after the search index plan (`plans/2026-09-29-search-index.md`) closed: "I want everything to be done before we cut a new version. Nothing waiting, nothing in todo roadmap or waiting." The open items are that handoff's Next Steps (`handoffs/2026-09-29-2102-search-index-done.md`): the main agent does not use `sh-agent search` in real work (0 of 28 sessions, D36 of the search plan), the Antigravity runs the release gate needs are incomplete, and the benefit on a project with a large trail was never measured over several sessions.

## Scope

In:
- A clean A/B of the whole harness per arm on Claude Code: `--setting-sources project,local` (measured 2026-09-29: user-level skills and context not loaded, project skills loaded), each arm's installed `.claude/{skills,agents,commands}` copied into the run's project, its rendered context appended once. This removes the first comparison's two caveats (v1.12.1 skills in both arms; the context loaded twice).
- Main-agent adoption: measured with the new skills in place; if the main agent still does not search, one change to the context or a skill (D2), measured again.
- A multi-session test on a large trail: three sessions of ordinary work in a clone of this repository at v1.12.1, per arm.
- Every Antigravity run `AGENTS.md` § Release requires: one completed run of each changed skill and agent.
- Reviews and verification of anything this plan changes; docs; the real-home smoke after the operator's install on this box.
- Trail closed: `TODO.md` and `ROADMAP.md` empty, the handoff with nothing queued, so `/preflight v1.13.0` is the only next step.

Out: the release itself (`/preflight` from the operator's turn); the known limits D25 (`x/sys` advisory, needs Go 1.25) and D35 (the edit-during-lock race), both documented decisions; installs on the other machines (after the release).

Follow-ups noted, not built: none yet.

## Surfaces touched

- `.agents/scratch/eff-kit/` (kit, gitignored, ADR-010), new files beside the Phase 5 kit, which stays as it is: `clean-session.sh` (one clean-arm session), `lane-clean.sh <arm> <app|trail> <n>` (a fresh project, then sessions s1–s3; results in `results-clean/<arm>-<app|trail>-<n>-s<k>.{json,err}`, projects in `~/sandbox/search/eff/cwork/`), `metrics-clean.py` (per lane type, session and arm; the main agent's and the subagents' `sh-agent search` calls in separate columns; prints the D3 count), `score-trail.py` (s1 drift fixed, s2 answers out of 5, s3 `go test ./cmd/sh-agent -run '^TestUsageMatchesDocComment$' -v` shows the test ran and passed, pushed), prompts `LT1`–`LT3` in `prompts.sh`. Each arm's Claude Code install: `~/sandbox/search/eff/chome-<arm>/` (`sh-agent install --vendor claude` of v1.12.1 and of this tree).
- `context/AGENTS.md` or one skill, only if D2 changes wording.
- `.agents/docs/testing/` logs; `testing/skill-matrix.md` if a scenario changes.

## Interfaces

Clean-arm Claude session: `cd <project>; PATH=<arm wrapper dir>:$PATH claude -p --setting-sources project,local --append-system-prompt "$(cat ctx-<arm>.md)" --permission-mode acceptEdits --allowedTools <the Phase 5 list> < prompt`, where `<project>/.claude/` holds a copy of `<arm home>/.claude/{skills,agents,commands}` from `sh-agent install --vendor claude --home <arm home>` of that arm's tree, excluded from git by `<project>/.git/info/exclude`. No `--agents`, no `--add-dir`.

Large-trail lane (per arm, one clone of this repository at v1.12.1 with its `.agents/docs` trail, the drift seed of the search plan's Phase 5, and a bare origin): s1 "resume from the trail, then fix the seeded drift: `sh-agent update --check` exits 2 when behind — make every living file agree, as a Task under the Definition of Done"; s2 "resume from the trail, then answer the five trail questions of the search plan's Phase 5 in ANSWERS.md with file:line, commit"; s3 "resume from the trail, then, as a Task under the Definition of Done, add a Go test `TestUsageMatchesDocComment` to `cmd/sh-agent` that fails when a verb in the usage text is missing from `main.go`'s doc comment, or the other way round"; the exact prompts are `LT1`–`LT3`. Each prompt carries its own "Do not ask me anything" and its push instruction; s2 is a question session that commits and pushes.

## Failure modes

- A run auto-denied on Antigravity: the command is added by name and the run repeated (project `AGENTS.md` § Testing).
- `agy -p` returns before a dispatched subagent: the run is repeated up to three times; still void after three, the gate is told in the log, and the operator decides (Open row).

## Verification

E2E criteria:
- The main agent of the search arm runs `sh-agent search` in at least 4 of 6 large-trail sessions (Phase 1, or Phase 2 after the change), or the operator has accepted the measured number in their own turn.
- The clean A/B's per-session table exists for Claude Code with both arms delivering; the main agent's `sh-agent search` count is reported for each arm.
- Every changed skill and agent has one completed Antigravity run on disk (artifact or answer) with no auto-denied line.
- After the operator's install on this box, the three real-home smoke checks pass on both CLIs, and `sh-agent search --refs ADR-040` answers from the installed binary.
- `TODO.md` and `ROADMAP.md` hold headers only; the newest handoff's Next Steps name only `/preflight v1.13.0`.

## Phases

- [x] Cold-reader check — round 1 `reviews/2026-09-29-coldread-search-release-readiness-r1.md` (one blocking gap, D3), round 2 `-r2.md` (two, D4), round 3 `-r3.md` (one narrow blocking item, one non-blocking) — cap reached; the operator, 2026-09-29: "Proceed with my fix (Recommended)", resolved as D5

### Phase 1 — Clean A/B on Claude Code, with the new skills
- [x] Kit: the files named in Surfaces touched; a dry run of one `trail` project per arm (setup only, no session) shows its `.claude/skills` holds that arm's skills (the search arm's `feature` skill contains `sh-agent search --refs`, the baseline's does not) and the drift seed in place — verify: AI — expect: `grep -c` 2 and 0, the seed line found — done 2026-09-29: `grep -c` 2 and 0, the seed line found, `.claude/` excluded, tree clean, in both arms' dry-run projects
- [ ] Sample apps, both arms, 2 runs each, 3 sessions (`lane-clean.sh <arm> linkshelf|tasklog <n>`) — verify: AI, `outcomes.py` pointed at `cwork/`, `metrics-clean.py` — expect: the app rows of the table, main-agent and subagent searches in separate columns
- [x] Large-trail lane, both arms, 2 runs each (`lane-clean.sh <arm> trail <n>`) — verify: AI, `score-trail.py`, `metrics-clean.py` — expect: per session the trail rows of the table and the three scores — done 2026-09-29: 12 of 12 sessions delivered; every lane: drift fixed, 5 of 5 answers, the new test passes
- [x] Adoption decision (D3): count the large-trail sessions of the search arm in which the MAIN agent (not a subagent) ran `sh-agent search` — verify: AI, the D3 line `metrics-clean.py` prints (main-agent calls only, from the main transcript; subagent transcripts are counted in their own column) — expect: a D-row stating `n of 6` and `adopted` (n ≥ 4) or `not adopted` (n ≤ 3); `not adopted` opens Phase 2, `adopted` marks Phase 2's boxes `[-]` with that D-row as the reason — done 2026-09-29: 3 of 6, `not adopted` (D9); Phase 2 runs
- [ ] Checkpoint — expect: `testing/<date>-search-adoption.md` committed — evidence: <path>, <date>

### Phase 2 — Main-agent adoption (runs when D3 says `not adopted`)
- [x] D2 decided with the operator: the change and its exact text — who: both — done 2026-09-29: pre-decided, D2
- [x] The change, `TestRenderedContextFitsAntigravityBlock` PASS if context — verify: unit — expect: PASS — done 2026-09-29: PASS at 21,996 of 22,000 bytes, GANNET last; the handoff skill's restatement of the order updated (D9)
- [x] Re-run of the probe (`search-first`, 3 with, 3 control) and the large-trail lane, search arm, 2 runs — verify: AI — expect: the main agent runs `sh-agent search` in at least 4 of the 6 large-trail sessions; fewer is an `Open:` row put to the operator with the numbers, never a silent pass — done 2026-09-29: large-trail 6 of 6 sessions (bar: 4), quality unchanged; the probe on the 65-file fixture 0 of 3 with the change and 0 of 3 control — the step is followed on a real-sized trail and skipped on a tiny one (D9)
- [ ] Checkpoint — expect: log and change committed; `/review --docs` over the change — evidence: <path>, <date>

### Phase 2b — Optimizations (the operator: "run optimizations")
- [x] Measure, then keep only what measures better — verify: AI, `index --rebuild` three times each, and an equivalence run of `parseRefs` over every tracked file — expect: faster build, smaller index, identical references — done 2026-09-29: build 2.2 s → 1.7 s and index 24.1 MB → 21.7 MB on this repository; the prefiltered `parseRefs` returned the same 40,842 references on 1,077 files; kept: insert statements prepared once per reconcile, `parseRefs` skipping lines that cannot match, the `refs(oid)` index dropped (it served only garbage collection); rejected: merging the FTS index after a build (the file grew to 25.8 MB from free pages); `TestRefsPrefilterIsEquivalent` added, mutations O1–O3 each fail it; the extension list is one source for the regex and the prefilter (D8)

### Phase 3 — Antigravity runs the release gate needs
- [ ] One completed run of each changed agent without one yet (`dreamer`; `security` on the final wording) and each changed skill (`review-standards`, `dream`, `handoff`, `bug`, `feature`, `retro`, `research`, `init-project`), in a sandbox home installed from this tree — verify: AI — expect: artifact or answer on disk, no auto-denied line, up to three attempts each
- [ ] Checkpoint — expect: `testing/<date>-search-agy-gate.md` committed — evidence: <path>, <date>

### Phase 4 — Real home and close
- [ ] `sh-agent update` on this box (D7, on the operator's word), then `hash -r` — verify: AI, `sh-agent status` — expect: binary line without a warning, drift none — first run 2026-09-29 21:32 at 2716d08: 33 written, `binary 2716d0858bc1` without a warning, GANNET last, `sh-agent search --refs ADR-040` answered; repeated after the last change of this plan
- [ ] Real-home smoke on both CLIs (`testing/2026-09-03-verify-real-home.md`) plus `sh-agent search --refs ADR-040` — verify: AI — expect: PASS, GANNET
- [ ] Definition of Done: verifier if source changed, reviews for what changed, changes entries, plan closed, `TODO.md`/`ROADMAP.md` empty, handoff naming only `/preflight v1.13.0` — verify: AI — expect: all on disk and pushed

## Risks

Assumptions: `--setting-sources project,local` excludes the real home's skills and context and loads project skills — checked 2026-09-29 by a probe (a project skill listed, no harness skill, no GANNET). Headless Claude Code reads `.claude/` in the project (only writes are refused) — checked by the same probe.

## Decisions

D1 — 2026-09-29 — The clean A/B replaces the first comparison's setup for everything in this plan; the first comparison stays as recorded — by: implementer, from the probe.
D4 — 2026-09-29 — The kit is new files beside the Phase 5 kit (named in Surfaces touched), because the Phase 5 `metrics.py` sums main and subagent searches and runs the larger task as one session; the matrix task adds a Zip place column, since the table has no row shape for a place — by: cold read round 2 (`reviews/2026-09-29-coldread-search-release-readiness-r2.md`).
D3 — 2026-09-29 — The adoption bar: the main agent of the search arm runs `sh-agent search` in at least 4 of the 6 large-trail sessions (2 runs × 3 sessions); that task is where the trail is large enough for search to matter; sample-app sessions are reported but do not decide — by: cold read round 1 (`reviews/2026-09-29-coldread-search-release-readiness-r1.md`, blocking question 1).
D2 — 2026-09-29 — If D3 says `not adopted`: `sh-agent search --open` joins the context's resume reading order, paid for by an equal cut named in a D-row, then measured again; if it still misses the bar, the number is reported and the release goes ahead — by: the operator, "Resume step (Recommended)".
D5 — 2026-09-29 — Round 3 residue: the third trail session's task is the usage/doc-comment test (the zip-place column documented a place without using `--from-zip`); the prompts carry their own waivers — by: the operator, "Proceed with my fix (Recommended)".
D6 — 2026-09-29 — The release is v1.13.0, a MINOR: new verbs, nothing an existing install must react to — by: the operator, "v1.13.0 (Recommended)", and "Cut new version update docs. This is dod". Scope's Out line on the release is superseded: `/preflight v1.13.0` runs at the end of this plan.
D10 — 2026-09-29 — Reviews of this plan's source changes (aae16af..f41e5c0): code APPROVE (`reviews/2026-09-29-review-aae16af-f41e5c0-code.md`); security two findings, both fixed with tests that failed first (`reviews/2026-09-29-security-aae16af-f41e5c0.md` § Outcome): the index files are refused unless regular files, and are mode 0600 — by: security review.
D9 — 2026-09-29 — Adoption: with the first wording and the new skills the main agent ran `sh-agent search` in 3 of 6 large-trail sessions (`not adopted`). Per D2, `→ `sh-agent search --open` (what is open)` joins the context's resume reading order after `bugs/INDEX.md`, paid for by cutting "`--open` lists open bugs and plans; " from the Search first entry (net +8 bytes, 21,996 of 22,000); `skills/handoff/SKILL.md`'s restatement of the order gains the step. Measured again in a separate arm: 6 of 6 large-trail sessions, same quality; the trail-question session read 12k characters against the baseline's 27k and cost 0.58 against 0.66 USD; the whole lane cost about 15% more than the baseline in both search arms, mostly in s1, where the search arms' agents checked the new bug number against the other branch and renumbered it — by: measured by implementer; change per the operator's D2.
D8 — 2026-09-29 — Optimizations kept only where measured (Phase 2b); warm search was already 25–40 ms and is unchanged; the schema version stays 1 because no release carries the old index — a development cache keeps its extra index until `index --rebuild` — by: the operator, "run optimizations"; measured by implementer.
D7 — 2026-09-29 — The operator asked for the local harness to be updated ("Don't forget to update local harness"): `sh-agent update` runs on this box from this session, once the tree is clean — by: the operator.
