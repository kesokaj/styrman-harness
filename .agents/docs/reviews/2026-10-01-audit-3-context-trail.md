---
owner: reviewer
date: 2026-10-01
---
# Review — audit-3-context-trail   (2026-10-01)

Range reviewed: whole repository at HEAD 2ba4a50 (no diff), slice C + TR of `.agents/docs/plans/2026-09-30-recon-audit-fix.md`: `context/AGENTS.md`, `context/dev-setup.md`, `README.md`, `AGENTS.md`, `REQUIREMENTS.md`, `ROADMAP.md`, `TODO.md`, `PRIORITY_TODO.md`, `.agents/docs/specs/sh-agent/*.md` (memory-loop.md § Known differences read; delivery.md spot-checked), `.agents/docs/testing/guide.md` (spot-checked), `.agents/docs/deploy/*.md` (local.md § Operator steps, remote.md spot-checked), `.agents/docs/architecture/DECISIONS.md` (ADR-047a, ADR-049, ADR-048 absence), `.agents/docs/bugs/INDEX.md`, newest three `changes/` (phase 3 entry read in full; phase 7 and 8 not), newest handoff (2026-10-01-0601), `skills/init-project/assets/*`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS (audit pass; no code)

Worst problem: `PRIORITY_TODO.md:22` tells the next reader that the four REQUIREMENTS edits were "applied in Phase 3 with ADR-048". They were not: `REQUIREMENTS.md` is untouched, no ADR-048 exists, and the plan's two boxes (`plans/2026-09-30-audit-fix.md:116-117`) are still open awaiting the operator's yes to the exact text. The plan, the Phase 3 changes entry and the Phase 3 line of `PRIORITY_TODO.md` are honest about the wait; the file that goes first on resume is not. The text awaiting approval lives only in gitignored scratch. The context itself is now consistent and 7,031 bytes under the Antigravity block, but only 31 bytes under the harness's own 15,000-byte pin.

## Findings

### Recon ID table

| ID | was | now | evidence |
|---|---|---|---|
| C1 | high | fixed | `go test ./internal/harness -run TestRenderedContext -count=1`: 14,969 B vs 22,000 and 15,000 pins, both PASS; whole rules moved to review-standards §7/§8, bash, golang, preflight, bug, retro (headroom: see N3) |
| C2 | high | fixed | `context/AGENTS.md:18` — TODO empty; ROADMAP holds "Now" (active plan while open) and "Next" (≤3) |
| C3 | high | fixed | `skills/init-project/assets/{TODO,ROADMAP}.md` carry title/note and Now/Next only; repo `ROADMAP.md` has no Later section; grep for "Later / Ideas" over context, skills, agents, README, ROADMAP, TODO: none |
| C4 | med | fixed (context side) | `context/AGENTS.md:64` points at "REQUIREMENTS R2" for the stub list; R2 still lists Vertex/Foundry, so both agree; the D29 R2 reword awaits approval (see C12) |
| C5 | med | fixed | `context/AGENTS.md:3` defines `docs/` = `.agents/docs/` in the preamble; notation is `docs/…` throughout |
| C6 | med | fixed | `context/AGENTS.md:9` TODO.md read seventh; order identical to `skills/review/SKILL.md:79` resume brief |
| C7 | med | fixed | `context/dev-setup.md:8` gates on `docs/deploy/local.md` § Operator steps (exists, local.md:112); no `settings.json` in dev-setup |
| C8 | med | fixed | `skills/init-project/assets/agents-tree/docs/bugs/INDEX.md:8-12` has `## Graduated` with the move rule |
| C9 | med | fixed | `context/AGENTS.md` § Hard rules is one short bullet; the escape/installed-command/`rm -r` lessons moved to `bash` (§ Installed commands, skills/bash/SKILL.md:52) |
| C10 | med | fixed | `context/AGENTS.md` § Hard rules: one bullet each, pointing at `review-standards` § 2 and § 8 |
| C11 | med | fixed | DoD step 2 (`context/AGENTS.md:16`) now ≈150 words; rounds protocol is `review` Step 4 (skills/review/SKILL.md:94), checked to exist |
| C12 | med | open | `REQUIREMENTS.md` R3, R4, R8, R9 provenance parentheticals unchanged (R9 still ≈190 words); awaits operator text approval (D29a) |
| C13 | low | fixed | `context/AGENTS.md` § Git has no "bypass hooks" |
| C14 | low | fixed | docs-only ≤10 lines now only in the Trivial bullet (`context/AGENTS.md:25`) |
| C15 | low | partly | `memory-loop.md:40-48` § Known differences exists with 5 rows (D17: `/goal` stays); R10 text still says "in the spec" with no pointer, while `memory-loop.md:42` says "R10 points here" (N5); `TaskStop` (context:35) is a Claude-only name not listed there (N6) |
| C16 | low | fixed | `context/AGENTS.md:9` "A verb that refuses because the project is not scaffolded is reported, then `/init-project` offered" |
| C17 | low | fixed | project `AGENTS.md` § Testing is 145 words, names § 3e, `internal/memory`, `internal/dream` and the corpus measurement |
| C18 | low | fixed | `skills/init-project/assets/REQUIREMENTS.md:7-30` leads with Must/Should/Could; no SLA tables first |
| C19 | low | fixed | `skills/init-project/assets/AGENTS.md:42` "verification tiers the plan names" |
| C20 | low | fixed | `ROADMAP.md:3` and `assets/ROADMAP.md:3` header comments describe the two sections |
| C21 | low | fixed | `README.md:65` olive "learn"; dashed olive line to the model provider; ADR-049 (DECISIONS.md tail); `TestHarnessMap` PASS |
| C22 | low | open | `REQUIREMENTS.md` R7 still "auto memory store"; awaits approval |
| C-scar | low | partly | context incidents all moved (tree-freeze, token filename, `hash -r`, BUG-085 text, `exec format error`); remaining: `REQUIREMENTS.md` Known limits BUG-085 measurement text (pending the same approval) |
| TR1 | high | partly | mechanisms now exist: `skills/review/scripts/restatement-grep.sh` (run by review Step 3.3), `testing/kit-preamble.sh`; BUG-122/123/086 in the Graduated table (INDEX.md:101-141). Not yet proven: BUG-138 (script silent outside the harness) fixed, BUG-139 and BUG-140 are open (INDEX.md:98-99) and are new recurrences of the "dispatch names the file" and "run proves its own result" classes |
| TR2 | med | fixed | `handoffs/2026-10-01-0601-audit-fix-phases-0-1.md` Environment State: installed from 032ece1, log scope named; the 2206 handoff is left as written (append-only), the newest wins |
| TR3 | med | partly | `context/AGENTS.md` § Hard rules now says time stamps come from the shell "in the same call that writes them"; the misordered 2145 entry cannot be fixed (append-only) and stays |
| TR4 | med | fixed | one real IN PROGRESS review remains, `reviews/2026-09-30-review-a275b96-head-code.md:9`, with an appended `## Stopped … void` section (line 20ff); `grep -E '^(## )?Status:? *IN PROGRESS'` over reviews and testing lists only it. D12 |
| TR5 | low | fixed | `plans/2026-09-15-harness-practices-g01-g26.md`: the four remaining `- [ ]` (lines 180-184) are the template fence; the only other plan with open boxes is the active audit-fix plan (7) |
| TR6 | low | open | `~/.config/sh-agent/proposals/` has only the six `*-retro.md`; no `proposals/dream/`; D10 keeps the timer opt-in, so the dreamer is still unrun on this box (the `sh-agent-update.timer` is the operator's update timer, not the dreamer's) |
| TR7 | info | open | `.agents/docs` is now 11 MB (7.72 MB on 09-30); reviews directory 386 files; grows with every audit |

### New findings

- **[medium / high]** `PRIORITY_TODO.md:22` — says the four REQUIREMENTS edits were "approved, applied in Phase 3 with ADR-048 (D29)". Reality: not applied, no ADR-048 (grep of DECISIONS.md), plan lines 116-117 unticked, `changes/2026-10-01-0632-audit-fix-phase-3.md:29` "REQUIREMENTS.md untouched". The same file ticks Phase 3 `[x]` (line 15) while its text says the edit "waits". Fix: reword line 22 to "approved in intent (D29); exact text awaits the operator's yes (D29a); not applied". The trail is honest in the plan, the changes entry and PRIORITY_TODO line 15, not in line 22 — a reader who stops at the first-read file concludes the requirements are done.
- **[medium / high]** `.agents/scratch/audit-fix/requirements.diff` — the exact text the operator must approve exists only in a gitignored directory (it does still apply cleanly to REQUIREMENTS.md by `patch --dry-run`). No handoff after Phase 3 exists to carry it as an operator Next Step, and the newest handoff's Next Step 2 says "when Phase 3 shows it". A scratch cleanup, or a fresh session, loses the proposal. Fix: copy the diff into the plan (or `plans/…-requirements-diff.md`) and state in a handoff that the operator's answer is still owed.
- **[medium / high]** `internal/harness/size_test.go:51` (`TestRenderedContextUnderTarget`) — render is 14,969 B against a hard 15,000 pin: 31 bytes. The generated `## Installed` block (`inventory.go:358`) is inside the measured text and lists every skill (≈15 B each), so two more skills in `skills/` fail `go test`, while R5 and `AGENTS.md` say adding a skill is "a file drop — never a code change". This recreates C1 (94 B headroom then, 31 B now) and the 14,500 goal is missed by 469 B (reported, not failing). Fix: pin the hand-written context separately from the generated Installed block, or raise the pin to the 22,000-byte Antigravity budget and keep 15,000 as a reported goal.
- **[low / high]** `handoffs/2026-10-01-0601-audit-fix-phases-0-1.md` — newest handoff says "Phase 2 in progress", Next Step 1 "Finish Phase 2"; eight phases later `PRIORITY_TODO.md` has Phase 9 as the only open phase. Resume still works only through the context's "except an open plan that postdates it" clause (`context/AGENTS.md:9`); `/review --resume` question 3 would report the contradiction. Fix: the Phase 9 handoff (planned) replaces it; an interim handoff is optional.
- **[low / high]** `ROADMAP.md:5` — "Last reviewed: 2026-09-30" after eight phases landed on 2026-10-01.
- **[low / medium]** `.agents/docs/specs/sh-agent/memory-loop.md:42` says "REQUIREMENTS.md R10 points here"; R10 (REQUIREMENTS.md:50) has no pointer until the pending edit lands. Fix: reword to "R10 will point here" or let the approved edit land first.
- **[low / high]** `context/AGENTS.md:35` — "stopped (`TaskStop`)" is a Claude Code tool name in a context R10 says must run the same on both CLIs; not listed in § Known differences. Fix: "stopped with the CLI's own stop" or add a Known differences row.
- **[info]** REQUIREMENTS vs code, sampled: R4 dirty-tree refusal (`TestUpdateNeverInstallsDirtyTree`, `TestInstallMarksDirtyTree`), R4 attribution key (settings.go), R9 zip UNVERIFIED (release.go:434), R9 "temporary directory" still holds (zip.go:96 `MkdirTemp` beside the zip), R3/R4 dream clauses (ADR-047a). No code contradicts REQUIREMENTS.md in what I read. Not run: the behavioural tests themselves beyond the three named below.

## Criteria
1. Dead code — n/a — audit
2. Flow tracing — n/a — audit
3. Standards conformance — n/a — audit
4. Directory layout — n/a — audit
5. Simplicity (DIY-except-security) — n/a — audit
6. Decoupling — n/a — audit

## Verdict (counts)

30 IDs: fixed 22, partly 4 (C15, C-scar, TR1, TR3), open 4 (C12, C22, TR6, TR7), struck 0. (C4 counted as fixed on the context side; its R2 reword shares the C12/C22 approval.) Highest open severity: medium (C12, the REQUIREMENTS provenance bulk, plus new findings N1-N3); no high is open. Resume order versus the `handoff` and `review --resume` briefs: identical. Templates versus the skills that fill them: `testing/guide.md` template is detected by the verifier (agents/verifier.md:21), TODO/ROADMAP/INDEX templates match the context's Records rule; no mismatch found. The Phase 3 REQUIREMENTS gap is honest in the plan and the Phase 3 changes entry, dishonest in `PRIORITY_TODO.md:22`.

Better than 2026-09-30: yes, clearly. The two always-loaded high findings are gone with evidence: the context renders at 14,969 B (was 21,906 of 22,000), no rule contradicts another (ROADMAP/TODO, resume order, docs-root notation), the templates no longer breed parking lots, the stopped review is void-marked, and the recurring bug classes have mechanisms on disk (a script and a kit preamble) instead of longer paragraphs. What is not better: the immutable REQUIREMENTS file is exactly as heavy and as stale as before because the operator's text approval is outstanding (and the text is stored only in scratch), the context's headroom is back to a few dozen bytes against its own pin, and the trail's front-door files (PRIORITY_TODO line 22, the 0601 handoff, ROADMAP date) lag the work.

## Not examined

- `specs/sh-agent/delivery.md`, `dream.md`, `search.md` beyond grep hits; `testing/guide.md` beyond §2/§4 header spot-check; `deploy/{gcp,remote}.md` beyond one grep.
- `changes/2026-10-01-0941-*` and `-1027-*` (phases 7, 8), read only via the plan and git log subject lines; the phase 3 entry was read in full.
- The behavioural claims in the phase logs (35 runs, sandbox matrix) were not re-run; only `TestRenderedContext*`, `TestHarnessMap`, `TestSkillLinksResolve`, `TestEmbeddedPlanEqualsCheckout` were executed (all PASS).
- Whether the hourly `sh-agent-update.timer` has since moved the real home past 032ece1 (the handoff's "installed from 032ece1" is time-dependent); `sh-agent status` was not run.
- README map prose beyond the legend lines and Dream bullet; TR7's trail-to-source ratio not recomputed.
