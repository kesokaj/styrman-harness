---
owner: reviewer
date: 2026-10-01
---
# Review — 329a99b-worktree-prompts-r2   (2026-10-01)

Range reviewed: working tree against 329a99b (`git diff 329a99b -- context agents skills README.md ROADMAP.md AGENTS.md`), prompts lens (`/review --prompts`), round 2 of `2026-10-01-review-329a99b-worktree-prompts.md`

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Round 1's two high findings are closed: the instrumentation-silence clauses are in `review-standards` § 2 (line 43, all three clauses present), and "you read any agent-written entry against the diff" is back in context step 4. The ADR move is only half done: `agents/changelog.md` still tells the haiku agent to write ADRs in five places while its own line 72 says ADRs are not its job (medium). Smaller: a stale `documentation` description, two rules (compaction, headless "say so") that moved to `handoff` with no pointer from the context, and harness-specific ids left in global skills.

## Round 1 answers, checked in the files as they stand
1. Closed. `skills/review-standards/SKILL.md:43`: monitor/watcher/filter "can emit at all", "demonstrably hit", "needle matches its own negation", must-pass/must-fail, unusual condition. Context line 59 points at § 2.
2. Closed. `context/AGENTS.md:18` carries it; "agent-written entry or ADR" narrowed to entry because the main agent now writes ADRs (consistent).
3. Mostly closed. ADR shape + "structural" paragraph are in `skills/feature/SKILL.md:67-86` before Step 3. Changelog description and the template are changed, but the body still has ADR duties (finding 1). `documentation` no longer claims ROADMAP/godoc in its body (finding 2 on its description). Context says Now and Next, matching templates and ROADMAP.md. `handoff` 36/66/168 carry the new order, the open-plan exception, Now/Next. Project `AGENTS.md` rule 4 has the ADR-049 exception; rule 1 names the model-provider send. README legend names olive and its dashed line.
4. Closed. `TaskStop` is at `context/AGENTS.md:35`.
5. Mostly closed. `review-standards` § 7/§ 8 and the handoff addition carry no BUG ids; context line 64 says "the harness's REQUIREMENTS R2". Remaining ids: finding 4.
6. Plan pointers: partly closed (finding 5).

## Findings
- **[medium / confidence: high]** `agents/changelog.md:10,78,83,89` (and `:33`) — the ADR duty is struck from the description and the template (line 72: "ADRs are not yours") but the body still assigns it: line 10 "(except marking ADRs superseded)", line 78 "add an ADR if one is due", line 83 Writes "and an appended ADR in `DECISIONS.md`", line 89 Reply "(changelog path, ADR number)". An agent follows its Writes and Reply lines; two authors for one record. Fix: strike those clauses and the "ADR style" in line 12; keep line 33 only as the layout tree.
- **[low / confidence: high]** `agents/documentation.md:3` — description still says "never writes changelogs or ADRs (those belong to changelog)"; its own line 41 now says ADRs belong to the main agent. The description is what a dispatcher reads. Fix: "never writes changelogs (changelog), ADRs, TODO.md or ROADMAP.md (the main agent)".
- **[medium / confidence: medium]** `context/AGENTS.md:11,31` — two always-loaded rules left with no pointer. Compaction ("preserve the active plan path, modified files, test commands, open findings; re-read the plan and `git status`, restate status") and the headless "say so and name the agent that has not returned" now live only in `skills/handoff/SKILL.md:26`. A session being compacted, or a `claude -p` run, has not loaded `handoff`, so the rule is never read at the moment it applies. The render passes `TestRenderedContextUnderTarget` with headroom (11.7 KB source). Fix: one clause each in the context ("on compaction: `handoff` § Compaction"; "an unreturned agent is named in the output") or move them to where a headless run loads them.
- **[low / confidence: high]** Harness-specific ids in global skills: `skills/review-standards/SKILL.md` § 7 and § 8 intros ("plan audit-fix, C1"), and `skills/review/scripts/restatement-grep.sh` comment ("TestEmbeddedPlanEqualsCheckout, BUG-137"). No other project has that plan, test or bug. Round 1 allowed only the "(moved from the global context 2026-10-01)" provenance note. Fix: drop the plan id; say "a test in this repository" or move that comment out of the global script.
- **[low / confidence: high]** `.agents/docs/plans/2026-09-30-audit-fix.md` line 32 (second half: "review-standards § 4 is the one home of the code-shape rules"; "apply review-standards § 4") and line 105 ("when review-standards § 4 carries that rule") still say § 4; it landed in § 8 (§ 4 is Directory layout). Lines 29, 39, 111 are right. Fix the three pointers.
- **[low / confidence: medium]** `context/AGENTS.md:53` — the "a file tool turns an escape into its character; grep for the escape; build the needle from character codes" rule is now only in `skills/bash/SKILL.md:50`. The failure occurs when writing any text (Go, JSON, Markdown), and `bash` loads for shell work. One pointer word in line 53 fixes the reach.
- **[low / confidence: low]** `.agents/docs/architecture/DECISIONS.md` ADR-049 — records no operator approval; `AGENTS.md` rule 3 sends a map change to the operator first. I could not tell from the diff whether it was given.

## Nothing-lost diff of the old context (329a99b) against the new one plus landing skills
Landed and found: Records/ADR shape (feature), rounds + operator-choice + "arrangement the finding measured" (review:60), frozen tree + "cannot wait" (feature:101), "deploys only by the guides", live-question answer and version-number rule (preflight:31), `FIRST verification command` (verifier.md:18), bugs status-only edit and graduation (bugs.md), retro cadence, BEFORE-/handoff, retros.jsonl, promotion (retro:32), runtime-data rule (golang), Test locally / amd64 / Artifact Registry / GHCR (docker, cloud-run-deploy:19-61), installed command and escape rules (bash), ticked-line rule (handoff:66, TODO template), name-clash rule (init-project:91), handoff stamp (handoff:26).
Not landed, judged acceptable: "if the trail cannot say where work stopped, ask" (replaced by the unscaffolded-verb sentence; "never guess intent" stays), "run it before Grep" and `--global` for search (README and `search` cover it), "a name without the time is older than every timed name of its day" (init-project:24 covers the migration), "Running as root does not verify the non-root path" (the "other branch" sentence remains).
Landed but with weakened reach: findings 3 and 6.

## Criteria
1. Dead code — n/a — prompts lens; the leftover ADR clauses in `agents/changelog.md` are the nearest equivalent (finding 1).
2. Flow tracing — pointers from the context resolve: `review` Step 4, `review-standards` § 2, § 7, § 8, `feature` § Decisions, `bash` § Installed commands, `bug` Step 0, `retro` § When to Suggest, `preflight` all exist; the plan's § 4 pointers are stale.
3. Standards conformance — see findings 1, 2, 4.
4. Directory layout — n/a.
5. Simplicity — n/a.
6. Decoupling — ADR ownership is split between two agent files (findings 1, 2).

## Not examined
- Antigravity and Claude behavioural runs of the new context and agents (no CLI run in a prompts review).
- `agents/*.md` other than changelog, documentation, bugs, verifier: grep only.
- Whether `TestHarnessMap` pins the new README legend wording (the ADR says it still passes; not run).
- The remaining skill diffs (preflight, retro, init-project assets beyond ROADMAP/TODO/local.md) read in the diff only.
- Only `go test ./internal/harness -run RenderedContext` was run (both pins PASS); no other test.

## Outcome

Every finding fixed in the phase; see the next round (round 3, `2026-10-01-review-329a99b-worktree-prompts-r3.md`, APPROVE WITH NITS). The two lost rules landed in review-standards § 2 and context step 4; the escape-in-file-tool rule stays in the bash skill by design.
