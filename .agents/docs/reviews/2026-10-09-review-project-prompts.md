---
owner: reviewer
date: 2026-10-09
---
# Review — project-prompts   (2026-10-09)

Range reviewed: Scope: the whole repository at HEAD (c267990), no diff.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing in the surface I read is broken. There are no dated prompt idioms, every `sh-agent` verb and flag the files name exists in the binary, and the thresholds that several files restate agree with each other. The biggest problem is a conflict in reply shape. `skills/review/SKILL.md:116` (2026-09-13) tells the main thread to report a review "In chat, at most ten lines: the worst finding first, one verdict per lens…". `context/AGENTS.md:5`, introduced on 2026-09-17 and so newer, gives every reply that reports a different shape: an explanation line, a recommendation line, then ONE NUMBERED TABLE. Both files load in the same session. There are two smaller problems. First, all eight agent definitions cap their hand-off reply at a fixed number of lines (prompt-audit 1f). Second, the `init-project` description restates a Claude Code version claim that REQUIREMENTS R2 says has one home, and the skill-listing budget it rides in has only 160 bytes left. Coverage is partial: see Not checked.

### Step 0 assumptions
- Scope: the dispatch's list, which is project `AGENTS.md`, `context/AGENTS.md`, `context/dev-setup.md`, `agents/*.md`, `rules/*.md`, `commands/` (empty) and every in-scope `skills/<name>/`. Vendored files are excluded: `grilling`, `codebase-design`, `find-skills`, the `gke-basics` fork, `skills/google-catalog/catalog/`, and the vendored `skills/review-standards/references/prompt-audit/prompt-audit.md` (the `skills/VENDORED.md` row).
- Target models: Claude and Gemini, as the dispatch states. A finding that rests on one model's documented behaviour is `flag`. Group 2 stale-fact and conflict findings are actionable.

### Inventory (Step 1)
- Two context files, one rule, eight agents, no commands, and 32 in-scope skill directories, about 700 KB of markdown, most of it in `references/`. The largest are `frontend` (73 KB), `review-standards` (70 KB), `mermaid-syntax` (45 KB), `golang` (38 KB), `terraform` (32 KB) and `bash` (32 KB).
- No request-building code: the Go binary assembles no model prompt.

### Summary by group
- Group 1 (dated prompt text): one finding (1f: numeric reply ceilings in eight agents). There are zero hits for `think step by step`, `<scratchpad>`, "be thorough", "don't narrate", trait claims or `You are a helpful`. Apart from the vendored prompt-audit file, no file has more than two capitalised `MUST/NEVER/ALWAYS/CRITICAL/IMPORTANT`.
- Group 2 (brittle configuration files): two findings and three flags.
- Group 3 (descriptions): covered by Group 2's description finding. The routing text in the remaining descriptions is trigger text (keep-list 6).
- Group 4: not applicable, because the scope holds no request-building code. Sub-agent roster check: zero.

### Checks that passed (evidence)
- **Verbs and flags.** These are named across the scope, and each exists in `cmd/sh-agent/main.go:144-200`: install, update `--from-zip`/`--check`, status, list, render, adopt, search `--refs`/`--kind`/`--open`/`--global`, index, memory add|show, dream `--install-timer`/`--dry-run`, metrics and report. The `--kind` values the scope uses (bug, plan, spec) are in `internal/search/parse.go:250`. `--open --kind plan` and `--refs … --kind bug` are combinations the search code accepts (`internal/search/search.go:88-107`, and tests at `search_test.go:369,460`).
- **Size pins at HEAD.** Run with `go test ./internal/harness -run 'TestRenderedContext|TestSkillListingBudget' -v`, all PASS:
  - hand-written context: 14,936 of 15,000 bytes (64 bytes left; 436 over the 14,500 goal);
  - Antigravity block: 17,345 of 22,000;
  - skill listing: 19,340 of 19,500 (160 bytes left).
  
  None of the hunks below adds bytes to `context/` or `rules/`. Hunk 2 removes about 70 bytes from the listing.
- **Thresholds agree:**
  - reviewer at 50 lines, 5 files, or a bug fix: `context/AGENTS.md:16` and `agents/reviewer.md:3`;
  - documentation over 30 lines: `context/AGENTS.md:17` and `agents/documentation.md:3`;
  - two review rounds: `context/AGENTS.md:16` and `skills/review/SKILL.md:100`.
- **Tools agree with bodies.** `bugs`, `changelog` and `documentation` have no Bash, and their bodies say "you have no shell" (`agents/bugs.md:20`, `agents/changelog.md:77`, `agents/documentation.md:38,48`). `skills/bug/SKILL.md:55` and `skills/review/SKILL.md:96,102` send `sh-agent search` output to those agents in the dispatch for that reason.
- **Install paths agree.** The paths named in skills (`~/.gemini/config/rules/AGENTS.md`, `~/.gemini/config/skills/…`) match `internal/harness/plan.go:14-15,57-62`.
- **Roster check:** no redundant pair.
  - `reviewer` and `security` share tools but differ in checklist, model (sonnet/opus) and artifact name (`agents/security.md:30-75`).
  - `bugs`, `changelog` and `documentation` share tools but have disjoint lanes, which their descriptions name.
  - `explorer` and `dreamer` differ in task and in where they write.

## Findings

Ordered by confidence. Each finding gives Location, Evidence, Pattern, Why obsolete, Confidence and Action.

### F1 — /review's chat report shape conflicts with the context's reply shape
- **Location:** `skills/review/SKILL.md:116` vs `context/AGENTS.md:5`
- **Evidence:** review: "In chat, at most ten lines: the worst finding first, one verdict per lens, the record's path (or the single artifact's), and what `--fix` changed or rebutted." Context: "A reply that reports, proposes or asks opens with a one-line explanation and a one-line recommendation, then ONE NUMBERED TABLE, highest priority first — `#`, the problem in one line, your recommendation … — then how to proceed".
- **Pattern:** Group 2, instruction files that contradict each other (plus 1f, a numeric ceiling).
- **Why obsolete:** Both load in the main thread when `/review` runs, and they prescribe different shapes for the same report. The skill does not name the rule it overrides. `git blame` shows the review line is from 2026-09-13 (157934bd). `git log -S'ONE NUMBERED TABLE'` shows the context rule arrived on 2026-09-17 (77f47f7), and the current form of line 5 is from 2026-10-05 (ee28fe33). The context is newer, so the older skill line is the one to align.
- **Confidence:** medium. The two could be read as compatible (a ten-line reply holding a short table), but a review with more than about five findings cannot meet both.
- **Action:** rewrite the skill line to point at the context's shape, with no numeric cap (hunk 1). It does not touch `context/`, so no size pin moves.

### F2 — the `init-project` description restates the Claude Code version claim whose one home is REQUIREMENTS R2
- **Location:** `skills/init-project/SKILL.md:3`, the description
- **Evidence:** "(no `CLAUDE.md`: Claude Code 2.1.277+ reads `AGENTS.md` natively)". `REQUIREMENTS.md:11` says: "This sentence is the one home of that list; the context and the init-project skill point at it."
- **Pattern:** Group 2, volatile specifics (a version number in a description) and duplicated information across SKILL.md files.
- **Why obsolete:** The description dates from 2026-09-21 (b013b52e), before the one-home edit of 2026-10-01 (ddac2d89, ADR-048). That later commit updated the body (`:61`, which now points at R2) and left the description alone. The description drops R2's 2.1.281 qualifier for Bedrock, Vertex, Foundry and telemetry-disabled setups, and it is sent with every request on both CLIs. The listing it rides in is at 19,340 of 19,500 bytes, and this parenthetical routes nothing.
- **Confidence:** medium. In substance it agrees with R2's headline clause, so this is a restatement, not an opposite rule.
- **Action:** remove the parenthetical (hunk 2). `:61` already carries the claim and points at R2.

### F3 — every agent caps its hand-off reply at "at most ten lines"
- **Locations:** `agents/reviewer.md:101`, `agents/security.md:83`, `agents/verifier.md:84`, `agents/documentation.md:64`, `agents/dreamer.md:81`, `agents/explorer.md:77`, `agents/changelog.md:90`, `agents/bugs.md:92`
- **Evidence:** "Reply in **at most ten lines**: the verdict, the report's path, and up to three findings worst-first. Nothing else — the detail is already on disk, and repeating it into the caller's context is what fills a session's window…" (reviewer). The other seven are the same shape.
- **Pattern:** 1f, output-shaping choreography (numeric output ceiling).
- **Why obsolete:** The reference says a stated operational reason (here, the caller's context window) does not turn a numeric clamp into a keeper. The goal should be stated as audience and outcome framing, and each reply's content list kept as its format instruction. Provenance: `git log -S` traces the line to the rebuild commit e521826. No bug row or ADR names it (INDEX.md and DECISIONS.md grep: no hit).
- **Confidence:** medium. The cap bounds a pointer message, not the work product, so the harm the row names (starved reasoning) is weak here. The operator may reasonably decline.
- **Action:** rewrite: keep each content list, drop the number (hunk 3, shown for reviewer; the same one-word edit applies at the seven other lines). Agents are outside both size pins.

### Flags (low confidence, no hunk)
- **L1** `skills/review-standards/SKILL.md:3,8` vs `context/AGENTS.md:16`. The skill names its loaders as "the reviewer agent, or `/preflight` reading a diff below the reviewer threshold". The context tells the main thread of every task to "read the diff yourself" below the threshold, and `context/AGENTS.md:60` tells authors to apply its § 8. That means the description does not route the main thread's own small-diff review to the criteria. This is a gap, not opposite rules. A fix would add bytes to a listing with 160 left, so it should pay for itself by shortening the same description. Operator's call.
- **L2** `context/AGENTS.md:64` (and `skills/init-project/SKILL.md:61`) send the reader to "the harness's REQUIREMENTS R2". `internal/harness/plan.go:57-62` installs only `context/`, `skills/`, `agents/`, `commands/` and `rules/`, so a session in a user project on a release-binary or zip place has no `REQUIREMENTS.md` to open. The pointer resolves only where a checkout exists. This is a delivery question (AGENTS.md rule 6), not a prompt edit.
- **L3** `skills/feedback-log/SKILL.md`, Behavior Rules item 3: "When in doubt, log." This matches the 1a row "If in doubt, use [tool]", but it carries its own reason (the cost asymmetry, and the handoff may cite only logged entries), and a "When NOT to Log" list bounds it. Flag only.
- **L4** Document-length ceilings: `skills/handoff/SKILL.md:133` ("under 150 lines"), `skills/explain-project/SKILL.md:23` ("Under 200 lines") and `skills/distill/SKILL.md:90` ("under 200 lines"). These are 1f-shaped numeric caps, but on documents a fresh reader loads, not on reasoning. Flag only.
- **L5** `context/AGENTS.md`. The hand-written context is 64 bytes under its 15,000 ceiling and 436 bytes over its 14,500 goal (`internal/harness/size_test.go:59-60`). Any future finding that adds a context sentence must name the sentence it removes. This is not a defect today.

## Proposed diff

Proposed only: the F1 and F2 rows are Group 2 conflict or volatile-specifics edits, which the reference marks as for the operator to confirm. Nothing here has been applied.

```diff
--- a/skills/review/SKILL.md
+++ b/skills/review/SKILL.md
@@ -116 +116 @@
-In chat, at most ten lines: the worst finding first, one verdict per lens, the record's path (or the single artifact's), and what `--fix` changed or rebutted. A review is step 2 of the Definition of Done; its other steps follow as usual.
+In chat, in the global context's reply shape: the explanation line carries one verdict per lens and the record's path (or the single artifact's), the recommendation line says what to do next, and the numbered table lists the findings worst-first with what `--fix` changed or rebutted in its recommendation column. A review is step 2 of the Definition of Done; its other steps follow as usual.
```

```diff
--- a/skills/init-project/SKILL.md
+++ b/skills/init-project/SKILL.md
@@ -3 +3 @@
-description: "Scaffold a new project with the standard agent directory structure (.agents/ tree with docs, rules, skills, agents, workflows, commands), AGENTS.md, REQUIREMENTS.md, TODO.md, and ROADMAP.md (no `CLAUDE.md`: Claude Code 2.1.277+ reads `AGENTS.md` natively), or bring an older project's scaffold up to the current templates without touching what the project wrote. …
+description: "Scaffold a new project with the standard agent directory structure (.agents/ tree with docs, rules, skills, agents, workflows, commands), AGENTS.md, REQUIREMENTS.md, TODO.md, and ROADMAP.md, or bring an older project's scaffold up to the current templates without touching what the project wrote. …
```
(The rest of the description line is unchanged.)

```diff
--- a/agents/reviewer.md
+++ b/agents/reviewer.md
@@ -101 +101 @@
-Reply in **at most ten lines**: the verdict, the report's path, and up to three findings worst-first. Nothing else — the detail is already on disk, and repeating it into the caller's context is what fills a session's window with what it just wrote.
+Reply with only the verdict, the report's path, and up to three findings worst-first — the detail is already on disk, and repeating it into the caller's context is what fills a session's window with what it just wrote.
```
(The same edit applies at `agents/security.md:83`, `agents/verifier.md:84`, `agents/documentation.md:64`, `agents/dreamer.md:81`, `agents/explorer.md:77`, `agents/changelog.md:90` and `agents/bugs.md:92`: "Reply in **at most ten lines**:" becomes "Reply with only". Before applying, grep for tests and docs that quote "at most ten lines" (Step 6, out-of-band dependencies). I did not run that grep.)

## Criteria
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

## Not checked
- **Read in full:** `context/AGENTS.md`, `context/dev-setup.md`, `agents/security.md`, `skills/feedback-log/SKILL.md`, `skills/review/SKILL.md:96-116`, `skills/init-project/SKILL.md:61`, and the frontmatter of every agent and every in-scope skill.
- **Read only through the session's deployed copies, assumed identical to HEAD but not diffed:** project `AGENTS.md`, `agents/reviewer.md`, `rules/critical-thinking.md`.
- **Not read beyond grep hits** (signal greps for Group 1, verb and flag greps, threshold greps, path greps):
  - the bodies of `agents/bugs.md`, `changelog.md`, `documentation.md`, `dreamer.md`, `explorer.md` and `verifier.md`;
  - every in-scope SKILL.md body except `feedback-log`, and the parts of `review` and `init-project` above;
  - every `references/`, `assets/` and `scripts/` file under the in-scope skills, including the `frontend`, `review-standards` (go-review, typescript-review), `golang`, `terraform`, `bash`, `harness-audit` and `mermaid-syntax` references.
- **Description-versus-body routing** was checked only for `feedback-log`, `init-project` and `review-standards`, not for the other 29 descriptions.
- **Agent `maxTurns` against body instructions** was checked only for the no-shell agents' tool lists, not for turn budgets.
- **No grep** for out-of-band copies of the proposed-diff strings (Step 6 and 7).
