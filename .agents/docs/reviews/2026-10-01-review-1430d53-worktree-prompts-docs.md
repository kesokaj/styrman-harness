---
owner: reviewer
date: 2026-10-01
---
# Review — 1430d53-worktree-prompts-docs   (2026-10-01)

Range reviewed: the working tree against commit 1430d53 (`git diff 1430d53 -- agents skills AGENTS.md README.md .agents/docs/specs .agents/docs/testing/guide.md`); lenses `--prompts --docs`

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

No blocker. The shared "Write the artifact EARLY" core is byte-identical in all five agents (md5 of the paragraph equal), each followed by one trigger line; the verifier trigger keeps "after the FIRST verification command ... after EVERY command". The Antigravity model sentences agree with `derive.go` (`opus`->`pro`, `haiku`->`flash`, `sonnet` omitted = inherits; reviewer is sonnet and inherits, security is opus and is pinned to `pro`). The new delivery.md § Antigravity derivation matches `DeriveAgy` as it stands (name, quoted description, mapped tools, model map, `~/.claude/skills/` rewrite, `subagent: true`, `mainAgent: false`; `effort`, `maxTurns`, `skills:` dropped, D5b). Five findings, none above medium: one wrong path in the reviewer, one stale ADR lane in the changelog agent, two verifier rules dropped, two stale doc restatements.

## Findings
- **[medium / high]** `agents/reviewer.md` pre-flight step 2 — says "its `references/<language>.md`". The files are `references/go-review.md` and `references/typescript-review.md` (`ls skills/review-standards/references`), so an agent following the pattern opens `go.md`, which does not exist, and skips the language checklist. Change to `references/<language>-review.md`.
- **[low / high]** `agents/changelog.md` § Your lane (the tree block, line starting `└── architecture/` and `DECISIONS.md   # ADRs, append-only; mark old ones "Superseded by ADR-NNN", don't delete`) — still presents DECISIONS.md and the supersede procedure as the changelog agent's lane, against its own description, § What each artifact captures ("ADRs are not yours") and Records ("never an ADR"). Delete those two tree lines (and "ADRs" from the lane heading's implication); `skills/feature/SKILL.md:85` already says the decision-holder writes ADRs. Also the sentence "name the ADR a new decision supersedes" (§ The date, last paragraph) invites authoring an ADR claim: reword to "cite the ADR ids the dispatch's `--refs` output lists".
- **[low / medium]** `agents/verifier.md` — old rules narrowed or lost in the nine-step rewrite (`git show 1430d53:agents/verifier.md`):
  1. "When the result line is there, the next call is the Write that completes the log; a boot-log line the suite did not assert is `not captured`, not a reason to run again" and "never a second run to capture something the first run's log did not keep" — now only "never a second run" in step 5.
  2. "cite it from the log; later logs cite the guide and record outcomes only" and "either way, name them in the log" (the moved/left per-app guides) — gone; "several become one guide" also drops the old "and leave them".
  3. Inherited failures: "the reply says so, so the main agent numbers it" — gone (every failure gets a number is the main agent's duty; the verifier no longer prompts it).
  4. Suppression listing: old said where (source files / `*_test.go` / `*.test.ts`, `*.spec.ts`) and that the list goes "under Coverage of changed code"; step 7 lists them but the template no longer says where. Minor.
  The explained-cause examples (test removed with its code; go/git missing from PATH, no symlink support) were also dropped; judged not rule-bearing. Restore 1 and 3 at least, one clause each.
- **[low / high]** `.agents/docs/specs/sh-agent/memory-loop.md:49` (Known differences row "Subagent `effort` and `maxTurns`") — "Why it stays" cites `(D5a)`, which plan audit-fix D5b supersedes, and the row says nothing of the `skills:` key now dropped on probe evidence. Change `(D5a)` to `(D5b; delivery.md § Antigravity derivation)`; the sentence "the derived `agent.md` carries only keys the vendor documents" stays true.
- **[low / high]** `agents/documentation.md` Records "Never" line — the old "nor `REQUIREMENTS.md`" was removed from it; rule 6 still forbids editing it, so there is no gap in behaviour, but the Records line (what a hurried run reads) no longer names it. Add it back to the list or accept; not a rule loss. Also rule 8 states "`README.md` and `ROADMAP.md` at the root carry none" (frontmatter) — a claim about a file the agent does not own; drop "and `ROADMAP.md`".

## Prompts lens answers
1. Core byte-identical: yes (5 of 5 equal). One trigger line each: yes. Verifier keeps FIRST/EVERY: yes.
2. Verifier nine steps: see finding 3 (nothing else lost; step 4 keeps `-v -count=1`, redirect, counts, vs previous, explained; step 5 keeps foreground/background/kill; step 7 keeps `.only(`, reason rule, pre-existing?; step 8 keeps tree-moved rule; "never FAIL code you cannot see" kept in step 1).
3. changelog and bugs: both say "you have no shell" or take diff/`git log --stat` from the dispatch; neither runs a command (`git log --stat` appears only as dispatch content; `date +%F` / `sh-agent search` are dispatch inputs). changelog writes no ADR in Records/description/body, except the stale lane tree (finding 2). bugs clean.
4. documentation: no claim over ROADMAP/TODO/DECISIONS/godoc except the rule-8 mention above; "Never" line names `ROADMAP.md`, `TODO.md`, `architecture/DECISIONS.md` (DECISIONS is not at the root, but the line says "the three root files that are the main agent's" — wording only).
5. reviewer pre-flight names `~/.claude/skills/review-standards/SKILL.md` (derive rewrites it for Antigravity): yes, but the references pattern is wrong (finding 1).
6. Dated incidents/BUG ids in agent bodies: no date strings; `agents/dreamer.md:31` still carries "BUG-100" (the clause "it never replaces reading every row, BUG-100"), which is in 1430d53 too, not introduced by this range and not a "Related" convention; `agents/bugs.md:67` BUG-001 is the template row. Left as is (low, optional strike).
7. Model sentences: `skills/review/SKILL.md:44` and `skills/preflight/SKILL.md:52` agree with `derive.go` modelMap and the agents' `model:` keys.

## Docs lens (restatement-grep results)
- "inherit the global skills": living files hold none; the phrase survives only in `.agents/docs/plans/2026-09-29-audit-search-consumers.md:30`, `2026-09-30-audit-fix.md:202`, `2026-09-30-research-audit-fix-formats.md:68` (plans, history). Pass.
- "skills:": `delivery.md:87` matches `derive.go`; `AGENTS.md:58` (toolMap sentence) names only `tools:` and stays true; `harness_test.go:155` asserts `skills:` is not emitted. Only stale item: memory-loop.md:49 (finding 4). Plan `2026-09-30-audit-fix.md` D5/D5a still read as live decisions but D5b supersedes them in the same plan (plan records, not touched).
- "1,024" / "one line": `AGENTS.md:59` matches `inventory.go:252` (one line, two banned phrases), `descriptionLimit = 1024` (`frontmatter.go:105`), and `TestAgentDescriptionRules` / `TestLintDescriptionLimit` exist (`harness_test.go:2306`, `2328`). Pass.
- "git show" / "ADR" against changelog's old lane: README and specs name no changelog `git show` or ADR-authoring duty; `context/AGENTS.md:18` and `skills/feature/SKILL.md:85` already say the main agent writes ADRs. Only stale restatement is inside `agents/changelog.md` (finding 2). `bugs/INDEX.md` BUG-061 Graduated row already points to feature § Decisions.

## Criteria
1. Dead code — n/a (prompts and docs only).
2. Flow tracing — n/a for code; doc-to-code traced for delivery.md vs `derive.go` and AGENTS.md vs `inventory.go`/tests: agree.
3. Standards conformance — findings 1-5 above.
4. Directory layout — n/a (no new files in the range other than the untracked testing log, outside this range).
5. Simplicity (DIY-except-security) — n/a.
6. Decoupling — n/a.

## Not examined
Code under `internal/` (only read for the derive/model/description claims); `context/` beyond the changelog and size sentences; behavioural runs of the changed agent definitions (AGENTS.md § Testing requires `claude -p` runs; the untracked `.agents/docs/testing/2026-10-01-verify-agents-behavioural.md` was not read); `agents/security.md`, `agents/explorer.md`, `agents/dreamer.md` bodies beyond the shared core and the lines grepped; README's full text beyond the grepped lines; whether the 1,024 limit also applies to agent descriptions in code (test name says "skill or agent", not opened).

## Outcome

All five fixed in the phase: reviewer names `references/go-review.md` and `references/typescript-review.md`; changelog's lane no longer lists DECISIONS.md or the Superseded procedure; the verifier's three clauses are restored (name moved guides and cite the guide; `not captured` is not a reason to run again; the reply says so, so the main agent numbers it); memory-loop.md's row cites D5b and the dropped `skills:` key; documentation's Never line names REQUIREMENTS.md. The dreamer's BUG-100 parenthetical was removed too.
