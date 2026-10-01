---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p5-prompts   (2026-10-01)

Range reviewed: working tree — the changed instruction files, read whole as they are in the working tree: context/AGENTS.md, agents/documentation.md, skills/bug/SKILL.md, skills/distill/SKILL.md, skills/dream/SKILL.md, skills/feature/SKILL.md, skills/feedback-log/SKILL.md, skills/handoff/SKILL.md, skills/init-project/SKILL.md, skills/preflight/SKILL.md, skills/research/SKILL.md, skills/retro/SKILL.md, skills/review-standards/SKILL.md, skills/review/SKILL.md (diff against HEAD: `git diff HEAD -- <files>`; findings may land anywhere in these files)

## Status: COMPLETE

## Verdict: REQUEST CHANGES (four Group 2 findings, none a safety defect; all proposed only, nothing applied)

The worst problem: `skills/review/SKILL.md:105` and `skills/retro/SKILL.md:143` tell the main agent to run only `TestRenderedContextFitsAntigravityBlock` (22,000-byte pin) after a `context/` or `rules/` edit, while the project `AGENTS.md` Setup commands name two pins and `internal/harness/size_test.go` has a second, tighter one (`TestRenderedContextUnderTarget`, 15,000 bytes). An edit applied under `--fix --prompts` or a retro can pass the instructed check and fail `go test ./...`. The other three findings are an init-project root that disagrees with the Go store root, an instruction to a shell-less agent to run `ls`, and a dangling cross-reference to "the review loop cap". Every conflict you asked me to check came back consistent except these (list under Criteria).

## Step 0: assumptions, inventory, provenance
- Scope: the 14 files named above. Skipped: everything else (no user-level files read; `~/.claude/skills/review-standards/references/prompt-audit/prompt-audit.md` read as the method).
- Target models: Claude and Gemini (every model that reads the files). A finding whose only reason is one model's documented behaviour is `flag`; Group 2 stale-fact and conflict findings are actionable as written.
- Inventory: 1 global context, 1 agent definition (`documentation`), 12 skills (all `SKILL.md` only; no reference files in scope except the prompt-audit reference used as method). No tool descriptions, no request-building code, no few-shot blocks beyond format-pinning templates.
- Provenance: `bugs/INDEX.md`, `plans/2026-09-30-audit-fix.md`, `git blame` (the init-project 5b block and the "moved from" notes are uncommitted in the tree; HEAD had the `--show-toplevel` root too). Prohibitions in these files trace to bug rows or plan decisions and were kept (keep-list 5); the "moved from the global context 2026-10-01" notes are allowed by the plan (earlier reviews 329a99b-worktree-prompts and -r2 recorded that) and are a flag only.
- Volatile specifics checked by reading (no commands run): paths `skills/init-project/assets/{AGENTS,REQUIREMENTS,ROADMAP,TODO}.md` and `agents-tree/`, `skills/review/scripts/restatement-grep.sh` (flags `-i`, `--all`, `--` match the usage block), `skills/review-standards/references/prompt-audit/prompt-audit.md`, `architecture/DECISIONS.md` all exist; `sh-agent` flags `--open`, `--kind`, `--refs`, `--global` and `--vendor` are defined in `cmd/sh-agent/main.go:107-117`; the `"vendor": "antigravity"` probe in init-project and dream matches `state.go` (`MarshalIndent`, `VendorAntigravity`); both size-pin test names exist.

## Summary
Group 2 (brittle configuration files): 4 findings, all actionable. Group 1 (dated prompt text): 0 actionable; low-confidence flags only. Group 3 (tool descriptions): not applicable. Group 4: not applicable (no request-building code); sub-agent roster check over `agents/documentation.md`, the only agent in scope: 0 findings (no near-duplicate agent in scope).

## Findings

### F1 — one size pin run where the project has two
- **Location:** `skills/review/SKILL.md:105`; `skills/retro/SKILL.md:143`
- **Evidence:** review: "a `context/` or `rules/` edit runs `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1` bare before the commit"; retro: "run `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1` bare after the edit and before the commit". Project `AGENTS.md` Setup: "must also pass the two size pins in internal/harness/size_test.go … under 22,000 bytes … under 15,000 bytes". `size_test.go:26` and `:59` define both.
- **Pattern:** Group 2, instruction files that disagree on one point (a restatement narrower than its source).
- **Why obsolete:** the second pin (ceiling 15,000, rendered context measured at 14,882 bytes in the phase 3 commit message) is the binding one; the instructed `-run` filter cannot match it, so the check can be green while the suite is red.
- **Confidence:** high (the repository contradicts it).
- **Action:** rewrite (diff below).

### F2 — init-project resolves the store root differently from the code and the dream skill
- **Location:** `skills/init-project/SKILL.md:112` (sub-step 1); also feeds `:113`, `:116`, `:120`
- **Evidence:** "`root=$(git rev-parse --show-toplevel)`". `internal/memory/memory.go:60-72`: "Root is the main worktree's top level … so every worktree of a repository shares one store (D26)" via `--git-common-dir`; `skills/dream/SKILL.md:22`: "`root` is the main worktree's top level, the parent of `git rev-parse --path-format=absolute --git-common-dir`".
- **Pattern:** Group 2, two instruction files ruling differently on the same point; the code sides with dream.
- **Why obsolete:** run from a linked worktree, init-project creates `.agents/memory`, appends the gitignore lines and writes `autoMemoryDirectory` for the worktree's own path, while `sh-agent memory add`, `dream` and `nativeOn` (reads `<root>/.claude/settings.local.json`, `memory.go:230`) use the main worktree: two stores, settings the Go side never reads, and a migration that moves files into the wrong one.
- **Confidence:** medium. The code and the dream text are verified; that Claude Code's `~/.claude/projects/<slug>` name follows the directory the session ran in (so `slug` should keep `--show-toplevel`) I did not verify.
- **Action:** rewrite (diff below). Direction assumed: the newer, code-backed rule (D26) wins.

### F3 — an agent without a shell is told to run `ls`
- **Location:** `agents/documentation.md:37`
- **Evidence:** "the single `cmd/<app>` directory (`ls -d cmd/*/` printing exactly one)"; frontmatter `tools: Read, Grep, Glob, Edit, Write`; `skills/bug/SKILL.md:55`, `feature/SKILL.md:111`, `preflight/SKILL.md:44` all say the agent "has no shell".
- **Pattern:** Group 2, volatile specifics: a command the agent's tool list cannot run.
- **Why obsolete:** the instruction is unexecutable as written; Gemini, which hangs on an unrecognised tool, or Claude would improvise or guess `<app>`, and the rule's last clause ("never left literal in a path you write") is the one it protects.
- **Confidence:** high (tools list contradicts it).
- **Action:** rewrite (diff below).

### F4 — "the review loop's cap" is referenced where no loop is described
- **Location:** `skills/bug/SKILL.md:15` ("a review loop that hits its three-round cap (Step 5)"); `skills/feature/SKILL.md:15` ("a review loop that hits its cap")
- **Evidence:** bug Step 5 (`:55`) dispatches `reviewer` once and describes no rounds; feature Step 6 (`:113`) likewise. The cap is in `skills/review/SKILL.md:112` ("**The cap** is the third round") and the rounds rule in `review` Step 4 (`:94`), which `context/AGENTS.md:16` points at.
- **Pattern:** Group 2, a cross-reference that no longer resolves to what it names.
- **Why obsolete:** a fresh session following bug Step 5 finds no cap to hit; whether the gate waits depends on it knowing `review` owns the loop.
- **Confidence:** medium (the intended target is clear; the wording is the only defect).
- **Action:** rewrite (diff below).

### Flags (low confidence or decided elsewhere; no diff)
- **[flag, low]** `skills/review/SKILL.md:94`, `skills/handoff/SKILL.md:26`, `skills/preflight/SKILL.md:31`, `skills/bug/SKILL.md:21`, `skills/retro/SKILL.md:32`, `skills/review-standards/SKILL.md:43,76,84`, `skills/feature/SKILL.md:85`: "moved from the global context 2026-10-01" and "the ADR was the changelog agent's until 2026-10-01" are migration-relative phrasing (Group 1d, history narratives): a diff against a prior version no reader saw. The plan allows the note (earlier reviews 329a99b recorded it); left as a flag so the operator can strike the nine parentheticals once the move has settled.
- **[flag, low]** `agents/documentation.md:3` names `docs/onboarding.md` in the description but the lane tree (`:24-35`) and Records line do not list it; `:34` cites "(G13)", an id nothing in the repository resolves for the reader; `:59` calls `architecture/DECISIONS.md` a "root file".
- **[flag, low]** `context/AGENTS.md:64` "(the harness's REQUIREMENTS R2 lists them)": an id only this repository's checkout can resolve; earlier review 329a99b flagged the same.
- **[flag, low]** `skills/handoff/SKILL.md:36,154` and `agents/changelog.md:81` call the context section "the Resume Protocol"; its heading is "Resume, then checkpoint" (`context/AGENTS.md:7`).
- **[flag, low]** `skills/dream/SKILL.md:3,36-42` enumerates five verdicts; Step 4 and the failure table add `index` (rows for index lines, not memories). Vocabulary enumeration differs; both are readable.
- **[flag, low]** `skills/init-project/SKILL.md:121` (sub-step 6): "the store is still created" holds only if sub-step 2's `.gitignore` append and `mkdir` ran; the headless guard names the settings write only.
- **[flag, low]** `skills/feedback-log/SKILL.md:24,30` encode single sessions as rules (the "Is it not retro first, then handoff?" quote; "A missed filing"); recency-trap pattern, kept because they are signal definitions the retro matches on.

## Proposed diff

```diff
--- a/skills/review/SKILL.md          # F1
+++ b/skills/review/SKILL.md
-a `context/` or `rules/` edit runs `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1` bare before the commit.
+a `context/` or `rules/` edit runs `go test ./internal/harness -run 'TestRenderedContext(FitsAntigravityBlock|UnderTarget)' -count=1` bare before the commit.
--- a/skills/retro/SKILL.md           # F1 (same finding, second restatement)
+++ b/skills/retro/SKILL.md
-run `go test ./internal/harness -run TestRenderedContextFitsAntigravityBlock -count=1` bare after the edit and before the commit
+run `go test ./internal/harness -run 'TestRenderedContext(FitsAntigravityBlock|UnderTarget)' -count=1` bare after the edit and before the commit
```

```diff
--- a/skills/init-project/SKILL.md    # F2
+++ b/skills/init-project/SKILL.md
-1. **Root.** `root=$(git rev-parse --show-toplevel)`; when that fails or is empty, stop Step 5b and say so.
+1. **Root.** `root` is the main worktree's top level, the parent of `git rev-parse --path-format=absolute --git-common-dir`, so every worktree shares one store, as `sh-agent memory` and the `dream` skill resolve it; when that fails or is empty, stop Step 5b and say so. Sub-step 5's `slug` is computed from `git rev-parse --show-toplevel` (the directory Claude Code ran in), not from `root`.
```
(and in sub-step 5 `slug=$(printf '%s' "$root" | …)` becomes `slug=$(git rev-parse --show-toplevel | sed 's#[^A-Za-z0-9-]#-#g')`; verify the slug rule before applying — the one unverified premise of F2.)

```diff
--- a/agents/documentation.md         # F3
+++ b/agents/documentation.md
-`<app>` in `specs/<app>/` is the single `cmd/<app>` directory (`ls -d cmd/*/` printing exactly one), which a heading or title in `AGENTS.md` does not override,
+`<app>` in `specs/<app>/` is the single `cmd/<app>` directory (Glob `cmd/*/` matching exactly one), which a heading or title in `AGENTS.md` does not override,
```

```diff
--- a/skills/bug/SKILL.md             # F4
+++ b/skills/bug/SKILL.md
-a review loop that hits its three-round cap (Step 5)
+a review loop that hits its three-round cap (`review` Step 5, "The cap")
--- a/skills/feature/SKILL.md
+++ b/skills/feature/SKILL.md
-and a review loop that hits its cap.
+and a review loop that hits its cap (`review` Step 5, "The cap").
```

## Criteria
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

## Conflict checks you named (all read, results)
- review Step 1 (`--all` excludes `--resume`) vs description and usage block: consistent (description names `resume` as its own lens; usage says "never part of --all"; argument-hint lists both).
- review Step 3 table vs its three brief subsections: consistent (docs item 7 under `--project` exists; prompts and resume rows, slugs and artifact names match the closing paragraph).
- preflight Step 1 vs 2 vs 5: consistent; `## Status: COMPLETE` is required in Step 2 for a covering reviewer and security report and in Step 5's blocking list; a verifier log IN PROGRESS blocks in both.
- handoff § Cold-reader check vs review § Resume brief: same order, same six questions; the only wording differences are "Any/Every" and Q6's phrasing, none changes the question.
- init-project Step 5b sub-steps against each other: consistent except F2 (root) and the low flag on sub-step 6; sub-steps 3/4/5 chain (trust check, keys, migration only when the key is held) agrees, as do the Antigravity-only skips with Step 1.
- distill `--from-log` vs feedback-log: consistent (the log is the user's words; distill never appends; `retros.jsonl` is its index; `signal` values `gap` and the category field exist in the log format).
- context resume order vs handoff § Cold-reader check and review § Resume brief: identical sequence; the "open plan that postdates the handoff outranks it" clause appears in all three.

## Not examined
Nothing in scope was left unread. Not done: Step 7 of the audit (behavioural probes of a removal) is out of the report-only brief; no model was run on any rule. Not read: files outside the 14 (e.g. `agents/changelog.md`, other skills) except the lines cited above; `skills/review/scripts/restatement-grep.sh` and `internal/harness/restatement_script_test.go` (changed in the tree but not instruction prose; outside this lens's scope list). F2's slug premise is unverified.

## Outcome
- F1 [high] one size pin named where the project has two — fixed at `skills/review/SKILL.md:105`, `skills/retro/SKILL.md:143`, `README.md:275` (`-run TestRenderedContext`, both pins)
- F2 [medium] init-project Step 5b root differs from `internal/memory/memory.go:62` and the dream skill — fixed at `skills/init-project/SKILL.md:112,118` (root = main worktree via `--git-common-dir`; `top` from `--show-toplevel` serves only the migration slug, since Claude Code names `~/.claude/projects/<slug>` after the directory the session ran in — the premise the finding left unverified is accepted as the Claude Code convention the `init-memory` matrix row already measures)
- F3 [high] documentation agent told to run `ls` with no shell — fixed at `agents/documentation.md:37` (a Glob for `cmd/*/main.go`); behavioural runs p5d2 and p5d3 after the edit, both dispatched (p5d3's ROADMAP lane result is BUG-139, unrelated to this line)
- F4 [medium] "review loop cap" referenced where no loop is described — fixed at `skills/bug/SKILL.md:15`, `skills/feature/SKILL.md:15` (the `review` skill Step 5 named)
- flags (low, no diff): left — "moved from the global context 2026-10-01" notes stay until the Phase 9 re-audit decides on them as one class; `agents/documentation.md:3` `docs/onboarding.md`, "(G13)", "root file" and the "Resume Protocol" name in handoff/changelog are Phase 9 docs-lens items; dream's `index` verdict row and init-project sub-step 6's wording are left as readable; feedback-log's signal definitions are kept by the reviewer's own reasoning
