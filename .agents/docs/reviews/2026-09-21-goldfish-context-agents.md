---
owner: reviewer
date: 2026-09-21
---
# Goldfish review — context/AGENTS.md   (2026-09-21)

Range reviewed: `context/AGENTS.md` alone (Part A), then compared against the repo (Part B), per the audit prompt's Step 5 / plan `2026-09-21-context-consolidation.md` Phase 3. G3 task: "close a bug record after a fix."

## Status: COMPLETE

## Verdict: SUFFICIENT? NO — MINIMAL? YES (every line earns its bytes; the gap is an omission, not slack)

The file is lean (16,569 bytes / ~4,142 tokens) and every sentence I read would cost a real mistake if deleted — no padding found. But it fails SUFFICIENCY on the assigned G3 task: closing a bug record is a two-file edit (`INDEX.md` status column and the record's own `status:` field) that context/AGENTS.md never describes — it documents only bug *creation* ("Dispatch `bugs`; it writes … with a prevention rule and updates `INDEX.md`") and rule *graduation* ("the commit that writes it has `bugs` move the matching index rows to the Graduated table"), and a goldfish reading only this file risks either leaving an open row stale or wrongly moving a plain fix into the Graduated table, conflating the two distinct mechanisms. The real answer lives only in `agents/bugs.md` ("When a fix lands later for an open bug, update its `status` in both files — that is the one allowed edit to a past record"), which context/AGENTS.md never points to.

## Part A — read only context/AGENTS.md and what it references by path

Referenced-by-path files read (Resume chain, the file's most concrete path references): `TODO.md`, `ROADMAP.md`, `REQUIREMENTS.md`, `.agents/docs/bugs/INDEX.md`, newest `.agents/docs/handoffs/2026-09-21-status-sha-question.md`, three newest `.agents/docs/changes/` entries, the active plan `.agents/docs/plans/2026-09-21-context-consolidation.md` (from ROADMAP "Now"). `context/dev-setup.md` was excluded: the task's own framing note (appended at install) is not a reference *within* the file's text.

### G1 — Comprehension

context/AGENTS.md is the global agent-context for this harness, deployed identically to Claude Code and the Antigravity CLI (overridden by a project's own `AGENTS.md`), so any agent working in any harness-managed project follows one discipline. It exists so progress survives across sessions and compaction: resume by reading, in fixed order, `AGENTS.md` → `TODO.md` → newest handoff → `ROADMAP.md` "Now" → the active plan → three newest `changes/` → bugs `INDEX.md`; state status and next step before editing; never guess intent when the trail doesn't say where work stopped. Work is classified (Trivial/Task/Feature/Bug) and the matching skill loaded. Every task runs a five-step Definition of Done: (1) `verifier` runs tests with the tree frozen during its dispatch; (2) `reviewer`/`security` review, findings fixed-or-rebutted in the report's Outcome, with escalation rules for repeated rounds; (3) `documentation` updates docs when behavior/interfaces/deploy changed; (4) `changelog` and `bugs` write records, the plan is ticked, `TODO.md`/`ROADMAP.md` stay empty at every checkpoint; (5) commit, push, `/handoff`. Subagents get an artifact, never the conversation, and any duty they must not skip is repeated in the dispatch text. Global (harness-wide) vs project (`.agents/{skills,agents,commands,rules}`) customizations must never share a name at the project level. A learning loop (`feedback-log`, `/retro`, `/distill`) turns corrections into graduated rules. Git is backup/versioning only — no CI, no hooks, strict staging/authorship/no-force-push rules. Hard rules cover: never hand-roll auth/crypto; dates from `date +%F`; verify writes with `git diff --numstat`/`git status`; `REQUIREMENTS.md` immutable without approval; scratch discipline (`.agents/scratch/`, never `/tmp`); several evidence rules (a claim states how it was verified; silence from instrumentation is not evidence; a must-fail mutation counts only if the code compiled and a named test failed); guard patterns for destructive/dry-run code; describing-text must track behaviour changes. Project layout conventions close the file.

Must do: read the resume chain in order; state status/next step before editing; classify work and load the matching skill; run the five DoD steps in order, dispatching the named subagents with a repeated critical-duty sentence; checkpoint at the fixed moments; number every failure via `bugs`; keep `TODO.md`/`ROADMAP.md` empty at checkpoints; commit/push per the git rules; obey every hard rule.

Must not: guess where work stopped; let "change nothing"/"report only" cancel a subagent's own artifact; let a subagent see the full conversation; reuse a harness skill/agent name at project level; commit secrets, force-push, rewrite history, or touch git config; hand-roll auth/crypto; call one green or root-only run "verified" for every path; let a dry-run mode perform its side effect; add a dependency the stdlib already covers.

### G2 — Critic (findings kept: 3, behaviour-changing only)

- **F1 [would cause a mistake].** The file documents bug *creation* ("Dispatch `bugs`; it writes … with a prevention rule and updates `INDEX.md`") and rule *graduation* ("the commit that writes it has `bugs` move the matching index rows to the Graduated table") but never states how an existing bug's `status` moves from `open` to `fixed`. A goldfish following only this text could (a) leave the INDEX row stuck at `open` forever after fixing it, or (b) misread "graduated" as the general closure step and wrongly move a plain fixed bug into the `## Graduated` table, corrupting the two-table schema (Graduated is reserved for a rule that became a permanent artifact, confirmed by `.agents/docs/bugs/INDEX.md`'s real Graduated table, which holds far fewer rows than "fixed" bugs in the main table).
- **F2 [would cause a mistake].** No schema is given for the individual `docs/bugs/BUG-NNN-*.md` record beyond "a prevention rule" — no mention that `status:` is frontmatter, that a past record has exactly one allowed later edit, or what sections exist. "Record the fix" (the G3 task's own wording) is unanswerable from this file alone.
- **F3 [would cause a mistake].** The file never says whether `.agents/docs/bugs/INDEX.md` is hand-edited or must always go through the `bugs` subagent. (The real answer, "Maintained by the `bugs` agent," is a comment inside `INDEX.md` itself, outside this file's reference set, and context/AGENTS.md points at INDEX.md only to be *read*, never as the sole writer.)

Kept only findings with a concrete failure mode; dropped speculative style notes.

### G3 — Readiness for "close a bug record after a fix"

1. **[BLOCKING]** How does a bug's `status` move from `open` to `fixed` in `INDEX.md` — hand-edited directly, or only via re-dispatching `bugs`?
2. **[BLOCKING]** What does the individual `BUG-NNN-*.md` record need beyond "a prevention rule" (a `status:` field, and is it the one allowed edit to a past record)?
3. **[BLOCKING]** Is `INDEX.md` ever edited directly by the main agent, or does every edit — including a status flip — go through `bugs`?
4. **[BLOCKING]** When does a fixed bug's prevention rule graduate to the `## Graduated` table versus stay a closed row in the main table — who judges "recurs," and is that a separate action from closing the bug?
5. **[NON-BLOCKING]** Does closing a bug get its own `changes/` entry, or does it fold into the fix's entry? (Reasonably inferred from DoD step 4's grouping and the repo's observed pattern, e.g. `changes/2026-09-16-bug-036-symlink-reads.md`.)
6. **[NON-BLOCKING]** Same commit as the fix, or a separate one? (Reasonably inferred: "commit … when a task or phase is done" implies one commit per completed task.)

## Part B — compared against the repo

G1 held up against `README.md`, the project `AGENTS.md`, `agents/bugs.md`, and `.agents/docs/` — no factual error found in G1 itself (the summary is accurate as far as it goes); the defect is a gap, not a wrong statement. The real closure procedure lives in `agents/bugs.md` §"The index" (verbatim): *"When a fix lands later for an open bug, update its `status` in both files — that is the one allowed edit to a past record."* — and its Records section confirms this is a `bugs`-subagent duty, not a freehand main-agent edit. That single sentence resolves Q1–Q3 outright and narrows Q4 (graduation is stated as a *separate*, later action, triggered only when a dispatch says a rule "graduated," never automatic on close). None of this is reachable from context/AGENTS.md or the files it references by path.

**Blocking G3 questions: 4 of 6.**

### Proposed fix (smallest byte cost, in the file a goldfish would look)

In `context/AGENTS.md`, the sentence ending "…Dispatch `bugs`; it writes `docs/bugs/BUG-NNN-*.md` with a prevention rule and updates `INDEX.md`." gains one clause:

> "…and updates `INDEX.md`; re-dispatch `bugs` to close a fixed one — flipping `status` in both the record and its index row is the only later edit to a past record."

~130 bytes. Answers Q1–Q3 (closure is a `bugs` re-dispatch, the field is `status`, both files move together) and leaves Q4 correctly separate (graduation already has its own sentence two lines down). Per the plan's own D10/Step-6 discipline, this addition should be paid for by trimming elsewhere and re-run against `TestRenderedContextFitsAntigravityBlock` before landing — Phase 1 left ~1,988 of the 22,000-byte budget unspent, so it likely fits without a further cut, but the size test is the authority, not this estimate.

## Report

- Tokens: 16,569 bytes ÷ 4 ≈ **4,142**.
- G1 errors: **0** (no factual error found; the gap is an omission, covered under G2/G3).
- G2 findings kept: **3** (F1–F3 above).
- G3 questions: **6 total, 4 blocking** (Q1–Q4).
- Verdict SUFFICIENT: **NO** — 4 blocking G3 questions and a G2 finding (F1) with a concrete corruption failure mode (misfiling a closed bug into the Graduated table).
- Verdict MINIMAL: **YES** — read the whole file; every sentence maps to a distinct, previously-recorded failure mode (many cite a BUG-NNN by name) or a distinct procedural step; no line was found that could be cut without losing exactly what F1–F3 show missing elsewhere — i.e., the file is already as tight as its coverage, it is simply missing the bug-closure clause, not carrying slack to trade for it.

## Not examined

- `docs/testing/guide.md`, `docs/deploy/{local,remote,gcp}.md`, `.agents/docs/retros/`, `~/.config/sh-agent/*` — mentioned in context/AGENTS.md but not load-bearing for the assigned G3 task (bug closure); not read.
- The other two Phase-3 goldfish files (`context/dev-setup.md`, this repo's own `AGENTS.md`) — out of this dispatch's scope by design (three separate reviewer dispatches per the plan).
- `skills/bug/SKILL.md` — deliberately not read for Part A since context/AGENTS.md names "the `bug` skill" without a path; read `agents/bugs.md` only for Part B comparison.
