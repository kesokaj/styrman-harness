---
owner: main
date: 2026-09-15
---
# Plan — harness practices from beads, elephant-goldfish and goldfish, G01–G26   (2026-09-15)

Start prompt for a fresh session in this checkout: `/feature Implement .agents/docs/plans/2026-09-15-harness-practices-g01-g26.md as written; the scope is confirmed, resume at the first unticked box.`

## Why

The operator asked for three repositories to be researched for ways to fight drift, keep more tasks and phases on disk, run more and smaller verifications, and plan better. The research record (`.agents/docs/plans/2026-09-15-research-beads-elephant-goldfish-goldfish.md`) verified 73 candidates into 40 groups. The operator chose all 26 add-now and add-next groups on 2026-09-15: "go with the first bundle g01-g26, start /feature", then "go all, use your default thresholds for G20".

## Scope

In: groups G01–G26 of the research record, as text edits to `context/AGENTS.md`, skills, agents, this repo's `AGENTS.md`, `README.md`, the sh-agent testing guide, the skill matrix, two stale plans, `DECISIONS.md` (ADR-021, ADR-022) and BUG-018's status. Group ids are the research record's.
Out: groups G27–G40; any Go code; `REQUIREMENTS.md`; installing into a real home (the operator's step, R6); refreshing the downstream mirrors (their own projects).
Follow-ups noted, not built: none yet (task 1.6 and Step 4 file discoveries here, each ending `found in: <task>`).

## Surfaces touched

- Global context: `context/AGENTS.md` lines 11, 16, 26, 31, 37, 71.
- Skills: `skills/feature/SKILL.md` (gates, Steps 1–4, 6); `skills/bug/SKILL.md` (new gates section, :29, :37, :45, :49); `skills/review/SKILL.md` (:48, :56, :60); `skills/preflight/SKILL.md` (:40, :50, :52, :65, :77, :81); `skills/handoff/SKILL.md` (:32, :33, :154); `skills/retro/SKILL.md` (:43, :56, :58, :75, :129, :131); `skills/distill/SKILL.md` (:98); `skills/init-project/SKILL.md` (:43, :44); `skills/review-standards/SKILL.md` (§1–§3) with `references/go-review.md` and `references/typescript-review.md` (§ Tests); `skills/golang/SKILL.md` (:23); `skills/golang/references/golang-ecosystem.md` (:323). None is vendored: `skills/VENDORED.md` lists grilling, codebase-design, find-skills, anti-ui-slop, gke-basics and terraform.
- Agents, bodies only except the `reviewer` description: `agents/verifier.md`, `agents/reviewer.md`, `agents/security.md`, `agents/explorer.md`, `agents/bugs.md` (:12, :24, template, :59, :69, :75, :80), `agents/changelog.md` (:94).
- This repo: `AGENTS.md` (:7, :46), `README.md` (:3, :86, :87), `.agents/docs/testing/sh-agent/guide.md` (:22), `.agents/docs/testing/skill-matrix.md`, `.agents/docs/plans/2026-09-03-rebuild.md` (:58, :97, :98, :99, :101, :113), `.agents/docs/plans/2026-09-13-learning-loop-project-level.md` (:136), `.agents/docs/architecture/DECISIONS.md`, BUG-018's record and `INDEX.md` row (status only).
- Scratch only, never committed: the docs-lens control repository `.agents/scratch/docs-lens-control/` (task 1.2) and the sandbox project `runs/g26-precedence` (task 5.1).
- Line numbers are as of commit `c192622`; edits are located by their old text.

## How edits are written, applied and dispatched

- Interfaces are listed in application order, phase by phase. Each edit is a heading `#### I-<n><letter> — <file>` followed by fenced blocks: `old` is text that must occur exactly once in the file at the moment the edit is applied, and `new` replaces it verbatim (an insert repeats the old text inside `new`). A range edit has `from` (the first line replaced) and `to` (the first line kept, or `<EOF>`), and its `new` is followed by one blank line before `to`. `new-A` and `new-B` are alternatives a task chooses between. Fenced blocks with any other label are reference text, not edits.
- Before any file of a task is written, every `old` of the task is confirmed to occur exactly once in the current file (BUG-010). A rehearsal parses this plan, applies every edit in order to copies of the files at `c192622`, and asserts each `old` occurs once at its moment (D15).
- In-session dispatches: this checkout has no `.claude/` directory, so a Claude Code session here loads skills and agents from the installed home, which does not carry this plan's edits, and R6 rules out installing there. Every dispatch whose behaviour this plan changes therefore carries the changed text in its brief: a `--docs` review gets the I-1a brief verbatim (plus I-4a's paths in this repo), and from the Phase 4 checkpoint on the verifier gets the repo's `agents/verifier.md` § What to run pasted with "apply these rules as written; they are not installed in this session". Phase 6 exercises the installed new definitions in the sandbox (D9).

## Interfaces — the exact edits

### I-1 — the `--docs` brief hunts old wording, dropped and narrowed rules (G06)

#### I-1a — skills/review/SKILL.md
```old
report only." |
```
```new
report only. For every rule or statement the scope changes or removes, `git grep` its old wording in the living docs as they are at the scope's head (the working tree, or `<b>` of a range; a symlinked directory is searched at its target) — `README.md`, `AGENTS.md`, `.agents/{rules,skills,agents,commands}/`, `<docs-root>/specs/`, `<docs-root>/testing/<app>/guide.md`, `<docs-root>/deploy/`, plus any paths the project's `AGENTS.md` adds — and never in `changes/`, `DECISIONS.md`, `handoffs/`, `reviews/`, `retros/`, bug records or dated logs: a hit that still says the old thing is a finding naming both file:line; a removed rule, command, flag or caveat with no home by meaning in those files is a finding `dropped: <old file:line>`; a rule whose new home leaves out a condition or qualifier the old text had is a finding `narrowed: <old file:line> → <new file:line>`; a claim you could not check is listed as unverified." Under `--project`: "classify each living doc as accurate, outdated or redundant with the file:line that proves it, and end with the `git grep` commands that re-check every stale reference." |
```

### I-2 — restore the repro clause (G06; closes BUG-018)

#### I-2a — skills/bug/SKILL.md
```old
**ask the user for one**
```
```new
**ask the user for one**, whatever autonomous instruction is in effect
```

### I-3 — a retro diff lists each removed rule's new home (G06)

#### I-3a — skills/retro/SKILL.md
```old
the **why** — what went wrong and how the change prevents it.
```
```new
the **why** — what went wrong and how the change prevents it. A diff that deletes or rewrites a rule lists under it each removed rule's new home (file:line) or `dropped: <why>`.
```

### I-4 — this repo's docs lens covers context, agents and skills (G06; `documentation` agent)

#### I-4a — AGENTS.md
```old
run from that home on each CLI before the install is called done.
```
```new
run from that home on each CLI before the install is called done. A change to `context/`, `agents/` or `skills/` also gets `/review --docs` over its range; here the living docs also include `context/`, `agents/`, and every skill directory under `skills/` not listed in `skills/VENDORED.md`, searched at `skills/` because `.agents/skills` is a symlink.
```

### I-5 — hand-typed counts become pointers (G08; `documentation` agent)

#### I-5a — README.md
```old
27 skills in the [Agent Skills](https://agentskills.io) `SKILL.md` standard, seven subagents,
```
```new
skills in the [Agent Skills](https://agentskills.io) `SKILL.md` standard, subagents,
```

#### I-5b — README.md
```old
seven subagents in Claude Code format
```
```new
subagents in Claude Code format
```

#### I-5c — README.md
```old
27 skills, each a directory
```
```new
skills, each a directory
```

#### I-5d — AGENTS.md
```old
27 skills (Agent Skills `SKILL.md` standard), seven subagents,
```
```new
skills (Agent Skills `SKILL.md` standard), subagents,
```

#### I-5e — .agents/docs/testing/sh-agent/guide.md
```old
# sanity: 27 skills, 7 agents parse; counts change only when a skill/agent is added or removed
```
```new
# sanity: every skill and agent parses; its section headers print the counts
```

### I-6 — BUG-018 status (G06)

Dispatch `bugs` to set `status: fixed` in `BUG-018-bug-skill-missing-autonomy-clause.md` and its `INDEX.md` row, the one edit `agents/bugs.md` allows on a past record. The Phase 1 changes entry names the fix (I-2a), and the Phase 5 changes entry says that I-22 replaces the clause with a gates section (D16).

### I-7 — the feature skill: Steps 1–4 and 6 (G02, G03, G04, G05, G09, G12, G13, G14)

#### I-7a — skills/feature/SKILL.md
```from
## Step 1: Scope
```
```to
## Step 5: E2E before "done"
```
````new
## Step 1: Scope

Restate the request in one or two sentences: "I read this as: X. Confirm or correct." If `$ARGUMENTS` is `#<n>` or an issue URL, read it first (`gh issue view <n> --json title,body,labels,comments`). Check `REQUIREMENTS.md`, `ROADMAP.md`, and `docs/bugs/INDEX.md` for anything that constrains or has bitten this area before; if `ROADMAP.md` "Now" already names a plan with unticked boxes, say so in the restatement. Misreads here cost the most later.

When the request ranges over a LIST — a backlog, `TODO.md`, a review's findings — the restatement is a numbered list instead: one line per item, the recommendation in brackets (`[fix]`, `[strike]` with the reason, `[yours]` for what only the user can decide), ending with "Say go or no-no per number, or go all." The answers ARE the scope; strikes the user approves are applied before the plan is written, with the reason on each line. A single request that leaves gaps in who, what, constraints or failure modes gets the same numbered form, one gap per line with your recommended answer; a gap the user does not settle now becomes an `Open:` row in the plan with your default.

When the confirmed scope touches an area this session has not read, or an estimated 10+ files, dispatch `explorer` with the scope restatement only, asking for the closest existing surfaces with file:line, the patterns to mirror, and what the request leaves open; `## Surfaces touched` cites its brief (`<docs-root>/plans/<date>-recon-<slug>.md`).

## Step 2: Plan file — `<docs-root>/plans/<date>-<slug>.md`

```markdown
---
owner: main
date: YYYY-MM-DD
---
# Plan — <feature>   (<date>)

## Why
<user problem; cite the issue, conversation, or REQUIREMENTS item>

## Scope
In: … / Out: … / Follow-ups noted, not built: …

## Surfaces touched
<files, packages, routes, tables, external services; the recon brief when Step 1 dispatched one>

## Interfaces
<function signatures, request/response shapes, schema changes, config keys — concrete enough that two implementers converge>

## Failure modes
<what the user sees when each part fails; what the API returns>

## Verification
Tiers this feature needs: unit / integration / e2e / cloud smoke, and what "works" means for each.
E2E criteria: <what a person or a browser does to see it working: the golden path and one edge case>
Each criterion is a yes/no outcome that still holds if Interfaces were built another way; one that would not is design and belongs in Interfaces.

## Phases
- [ ] Cold-reader check — design ready and implementation ready (<reviews/ paths>), or the user's same-turn override quoted with its date
### Phase 1 — <name> (dependency order: data model and migrations → services → API → UI → e2e)
- [ ] <task> — verify: <tier, command, who: AI | user | both> — expect: <what it prints or returns when done; an empty result names its positive control>
- [ ] …
- [ ] Checkpoint — expect: <observable end state of the phase> — evidence: <verify log path, or command and exit status>, <date>
### Phase 2 — …

## Risks
Assumptions: <a claim about an external system, a consumer, or a machine> — checked by <command>, first task of the phase that depends on it; or none

## Decisions
D<n> — <date> — <decision> — by: <the user's words quoted | cold read round N | implementer>
Open: <question verbatim> — default: <what the implementer does>
Rejected: <alternative> — <why>
```

Point `ROADMAP.md` "Now" at the plan. For a small feature the plan can be short; the sections stay. A task that adds or changes a command, flag, endpoint, config key or deploy step ends `— docs: <file>` and is ticked only when that file says so; other tasks carry no docs clause. A box that will not be done ends as `- [-] <task> — dropped: <why> | superseded by <plan § task> | moved to <handoff path> Next Steps (<date>)`; `[-]` is closed, never remaining work.

## Step 3: Cold-reader check — no code until it passes

Dispatch the `reviewer` agent with **only the plan file**, nothing from the conversation, and this brief:

> You are a fresh reader with no prior context. Read this plan and the code surfaces it names. Answer two things. (1) Gaps: what must be decided before implementation starts that the plan does not decide — interfaces, failure modes, scope edges, misunderstandings of existing code (cite file:line), a Verification or E2E criterion that names an implementation choice, an assumption asserted but not checked, a task changing a command, flag, endpoint, config key or deploy step with no docs clause. (2) Readiness: could a first-pass implementer ship this without asking the author a question? List every question you would have to ask, and mark each BLOCKING (two careful implementers would build different things) or NON-BLOCKING (the plan states a default, or a reasonable implementer decides it alone). End with `design ready` only if you have zero gaps, and `implementation ready` only if you have zero BLOCKING questions; otherwise `design needs revision` / `implementation not ready`.

Address every gap and blocking question in the plan, or rebut it in writing citing `REQUIREMENTS.md`, `AGENTS.md`, or the user's words; write each rebuttal and each non-blocking question as a row under `## Decisions`. Re-run until both tokens are "ready". After three rounds without convergence the feature is under-specified: stop and ask the user which requirement to clarify or cut — and say in the question whether the residue is a real ambiguity or one narrow item you would resolve the same way, with your recommendation first. Tick the Cold-reader check box only on both tokens or a quoted same-turn override; residue the user accepts at the cap is an `Open:` row. Until then, no source edits; reading, `git log`, and running existing tests are fine.

## Step 4: Phases, one at a time

For each phase, implement its tasks in order. For a task whose verify clause names a test, write that test first and run it bare: it must fail for the task's reason, because a test that passes before the code exists is not the test (review-standards § 2); in a compiled language first add the signature from Interfaces with a zero-value body, since a compile error is not the red. Tasks verified by a walkthrough or a container run are unchanged and still need covering tests at Done. After each task, run its named verification and tick its box when the result meets its expect.

Between checkpoints the plan holds the state. At a stopping point inside a phase — a blocker, a long run, the context may compact, the session may end — append to the box you are on `— in progress: <where it stopped>; blocker: none | <what>`, and delete the note when you tick the box. Work found mid-phase is filed when found: a blocker is a new box above the task it blocks, ending `— found in: <task>`; something broken goes through the `bug` skill; deferrable work is one line under Scope's Follow-ups ending `found in: <task>`; a question only the user can decide is an `Open:` row under `## Decisions`, named in the Step 6 report. A decision the plan did not make — an interface detail, a scope edge, a workaround — gets a dated D-row before work continues; a structural one also gets an ADR.

When the phase's tasks are ticked, run the project's full test command **as its own tool call, bare, and read its exit status** (never `test | tail && git commit` — BUG-033/041 in styrman-crew both shipped a red test that way). In a plan with two or more phases, a phase that changed source then gets a `verifier` dispatch over `<commit before the phase>..HEAD` (or the working tree), slug `<plan-slug>-phase-<N>`; its log is committed with the phase. Copy the evidence onto the phase's Checkpoint box from the log or the call, never from memory, and tick it only when the evidence meets its expect — a missed expectation leaves it unticked and the next phase closed. Then write the `changes/` entry for the phase (`changelog` agent or yourself), and commit and push as separate calls. That is the checkpoint. A new session resumes at the first unticked box.

Per-layer verification, typical Go + TanStack Start project: migration applied and schema inspected → package tests with `-race` → API hit with `curl` → UI driven in a real browser (golden path, one edge case) → container built and run locally (`docker build`, `docker run`, `curl /healthz`).
````

#### I-7b — skills/feature/SKILL.md
```from
## Step 6: Definition of Done
```
```to
<EOF>
```
```new
## Step 6: Definition of Done

Steps 1–5 of the global Definition of Done: `verifier` log, `reviewer` (diff only) and `security` when due, `documentation` when behavior or interfaces changed, `changelog` naming the agents dispatched, `bugs` for anything that went wrong on the way; before the commit, close the plan — every `[ ]` becomes `[x]` or `[-]`, and a `who: user` box not done moves to the handoff's Next Steps — then commit and push. Then a report in chat: feature in one line, files by phase, tests added, cold-reader result and how each gap was resolved, review outcome with rebuttals verbatim, e2e result, follow-ups, and every `Open:` row.
```

### I-8 — the global context: checkpoint box, compaction, overgrown tasks (G01, G03, G19)

#### I-8a — context/AGENTS.md
```old
(tick its tasks, `changes/` entry, commit, push)
```
```new
(tick its Checkpoint box, `changes/` entry, commit, push)
```

#### I-8b — context/AGENTS.md
```old
Progress that lives only in chat is lost.
```
```new
Progress that lives only in chat, a compaction summary, or the CLI's own memory is lost: write it to the trail file that owns it.
```

#### I-8c — context/AGENTS.md
```old
the test commands, and open findings.
```
```new
the test commands, and open findings; after a compaction, re-read that plan and `git status` and restate status and next step before the next edit.
```

#### I-8d — context/AGENTS.md
```old
three-line plan in chat, do it, Definition of Done.
```
```new
three-line plan in chat, do it, Definition of Done; a Task that meets a design question or outgrows one sitting becomes a Feature for the remainder, under the feature skill's gates.
```

#### I-8e — context/AGENTS.md
```old
Workflow skills, also typed as `/name`: `/init-project` `/handoff` `/review` `/research` `/preflight` `/retro` `/distill` `/grilling` `/feature` `/bug`. The harness is
```
```new
The harness is
```

### I-9 — the handoff reads the plan's state and checks for drift (G03, G04, G05)

#### I-9a — skills/handoff/SKILL.md
```old
3. **What's still open** — Unfinished work, known issues, blocked items
```
```new
3. **What's still open** — the active plan's first unticked box with its in-progress note, then its Follow-ups and `Open:` rows, then unfinished work, known issues, blocked items
```

#### I-9b — skills/handoff/SKILL.md
```old
4. **Key decisions** — Architecture choices, trade-offs made, rejected alternatives
```
```new
4. **Key decisions** — the active plan's `## Decisions` rows first, then architecture choices, trade-offs made, rejected alternatives from the transcript
```

#### I-9c — skills/handoff/SKILL.md
```old
Ask for four things:
```
```new
Ask for six things:
```

#### I-9d — skills/handoff/SKILL.md
```old
any contradiction between handoff, ROADMAP, TODO, and plan,
```
```new
any contradiction between handoff, ROADMAP, TODO, and plan, any plan with `[ ]` boxes that `ROADMAP.md` "Now" does not name (`grep -l '^- \[ \]' <docs-root>/plans/*.md`), any commit since the previous handoff that changed paths outside `<docs-root>` with no `changes/` entry covering it,
```

### I-10 — noticed work goes to the plan or the handoff, not TODO.md (G03)

#### I-10a — skills/bug/SKILL.md
```old
No adjacent refactors, no defensive code
```
```new
No adjacent refactors (note what you saw under the active plan's Follow-ups, or the handoff's Next Steps when no plan is active), no defensive code
```

#### I-10b — agents/changelog.md
```old
any deferred work the main agent should put in `TODO.md`
```
```new
any deferred work the main agent should file under the active plan's Follow-ups or the handoff's Next Steps
```

### I-11 — close the stale plan boxes (G05)

#### I-11a — .agents/docs/plans/2026-09-03-rebuild.md
```old
- [ ] **(you)** On the Mac: check that Antigravity CLI lists the skills (`/skills`) and the seven agents (`/agents`) and that `~/.gemini/config/rules/AGENTS.md` is picked up (ask agy what its rules say) — all three verified on the Linux box 2026-09-03; Mac only
```
```new
- [x] **(you)** On the Mac: check that Antigravity CLI lists the skills (`/skills`) and the seven agents (`/agents`) and that `~/.gemini/config/rules/AGENTS.md` is picked up (ask agy what its rules say) — all three verified on the Linux box 2026-09-03; Mac only — verified 2026-09-04 (`testing/2026-09-04-verify-mac-real-home.md`)
```

#### I-11b — .agents/docs/plans/2026-09-03-rebuild.md
```old
- [ ] **(you)** Verified on the Mac: `go run ./cmd/sh-agent install` the first time (clone `~/workspace/styrman-harness` if absent), `sh-agent update` after that, then `sh-agent status` (done on the Linux box 2026-09-03; Mac only — see section 12)
```
```new
- [-] **(you)** Verified on the Mac: `go run ./cmd/sh-agent install` the first time (clone `~/workspace/styrman-harness` if absent), `sh-agent update` after that, then `sh-agent status` (done on the Linux box 2026-09-03; Mac only — see section 12) — moved to `plans/2026-09-15-harness-practices-g01-g26.md` Done list (2026-09-15)
```

#### I-11c — .agents/docs/plans/2026-09-03-rebuild.md
```old
- [ ] **(you)** `sh-agent update` on the other machines (the Mac — see section 12)
```
```new
- [-] **(you)** `sh-agent update` on the other machines (the Mac — see section 12) — moved to `plans/2026-09-15-harness-practices-g01-g26.md` Done list (2026-09-15)
```

#### I-11d — .agents/docs/plans/2026-09-03-rebuild.md
```old
- [ ] **(you)** On the Mac: check Claude Code `/memory` shows `~/.claude/CLAUDE.md` importing `AGENTS.md`, and `/skills` lists the 28
```
```new
- [-] **(you)** On the Mac: check Claude Code `/memory` shows `~/.claude/CLAUDE.md` importing `AGENTS.md`, and `/skills` lists the 28 — moved to `plans/2026-09-15-harness-practices-g01-g26.md` Done list (2026-09-15)
```

#### I-11e — .agents/docs/plans/2026-09-03-rebuild.md
```old
- [ ] **(you)** Run `sh-agent update --check` from another machine after a push to confirm the resync loop (the Mac — see section 12)
```
```new
- [-] **(you)** Run `sh-agent update --check` from another machine after a push to confirm the resync loop (the Mac — see section 12) — moved to `plans/2026-09-15-harness-practices-g01-g26.md` Done list (2026-09-15)
```

#### I-11f — .agents/docs/plans/2026-09-03-rebuild.md
```old
- [ ] **(you, later)** Live check on the Mac: Claude Code `/memory` and `/skills`; Antigravity `/skills`, `/agents`, and a question that only the global context can answer (headless equivalents pass on the Linux box; Mac only — see section 12)
```
```new
- [-] **(you, later)** Live check on the Mac: Claude Code `/memory` and `/skills`; Antigravity `/skills`, `/agents`, and a question that only the global context can answer (headless equivalents pass on the Linux box; Mac only — see section 12) — moved to `plans/2026-09-15-harness-practices-g01-g26.md` Done list (2026-09-15)
```

#### I-11g — .agents/docs/plans/2026-09-13-learning-loop-project-level.md
```old
- [ ] Operator: on this box run `sh-agent update` (or `go run ./cmd/sh-agent install`) once, since the installed 2026-09-03 binary cannot rebuild itself; on the Mac `sh-agent update` — who: user
```
```new
- [-] Operator: on this box run `sh-agent update` (or `go run ./cmd/sh-agent install`) once, since the installed 2026-09-03 binary cannot rebuild itself; on the Mac `sh-agent update` — who: user — moved to `plans/2026-09-15-harness-practices-g01-g26.md` Done list (2026-09-15)
```

### I-12 — review rounds keep a ledger and hand over the fix list (G10)

#### I-12a — skills/review/SKILL.md
```old
A round: apply each finding or rebut it in writing in the record; run the project's checks bare and read the status; re-dispatch the same lens agents over the same scope with the same slugs (they overwrite their artifacts; the record gains a `## Round N` section and, at the end, its `## Outcome`).
```
```new
A round: apply each finding or rebut it in the record; run the project's checks bare and read the status; write the record's `## Round N` section, one line per finding so far — the finding quoted, then `fixed at <file:line>`, `rebutted: <reason>` (Step 4) or `left: <owner>`; re-dispatch the same lens agents over the same scope with the same slugs (they overwrite their artifacts), each brief followed by `---`, that section's fixed lines only (`<file:line> — what changed`), and "Verify each fix is correct and complete; look for defects the fixes introduced." Rebuttal reasons stay in the record, out of the dispatch. The record ends with its `## Outcome`.
```

#### I-12b — skills/review/SKILL.md
```old
Stop when a round reports nothing, after the third round
```
```new
Stop when a round reports nothing and no ledger line is still in flight (a `--docs` list with `documentation`, a `--tests` gap not yet written), after the third round
```

### I-13 — a rebuttal needs a checkable reason (G11)

#### I-13a — skills/review/SKILL.md
```old
, `rebutted` (why), or `left` (who owns it)
```
```new
, `rebutted` (one checkable reason: not a defect because X, with X at file:line; a trade-off recorded at file:line, in an ADR, or in an AGENTS.md section; or the user's words, quoted — never an intent the next reader cannot verify, which is fixed instead, and never a command that exited non-zero), or `left` (who owns it; a finding outside the scope is `left`, not `rebutted`)
```

#### I-13b — skills/preflight/SKILL.md
```old
— `fixed`, `rebutted`, or `left`,
```
```new
— `fixed`, `rebutted` (a reason as `/review` Step 4 defines it), or `left`,
```

### I-14 — the reviewer checklist gains the missing hunt items (G15, G16)

#### I-14a — skills/review-standards/SKILL.md
```old
- [ ] Verdict is always "delete"
```
```new
- [ ] Leftover debug output and temporary scaffolding added for the change (prints, `console.log`, debug logs).
- [ ] Verdict is always "delete"
```

#### I-14b — skills/review-standards/SKILL.md
```old
it reports coverage that does not exist.
```
```new
it reports coverage that does not exist.
- [ ] A write that can fail halfway leaves no partial state; a check-then-act on a file or row is not a TOCTOU race; a fallback does not mask the error it replaces; empty, NaN, Unicode and timezone inputs are handled where the input allows them.
- [ ] A source change whose generated counterpart (route tree, protobuf, `go generate` output, embedded build) is absent from the diff.
- [ ] An assertion changed in an existing test traces to a behaviour change in the source diff; a suppression directive, skip, `.only` or lint-config exclusion the diff adds names its rule or condition and a reason other than making the check pass — otherwise it is a finding.
```

#### I-14c — skills/review-standards/SKILL.md
```old
are verified as run, not re-litigated line-by-line.
```
```new
are verified as run, not re-litigated line-by-line.
- [ ] Comments, doc-comments, test names and spec rows that name behaviour the diff changed still describe it — the global context's describing-text rule; a stale one is a finding, not a comment request.
```

### I-15 — suppressions in the language files (G16)

#### I-15a — skills/review-standards/references/go-review.md
```old
no sleeps as synchronization — use channels or `testing/synctest`-style fakes.
```
```new
no sleeps as synchronization — use channels or `testing/synctest`-style fakes.
- [ ] `t.Skip` only behind a named prerequisite (`testing.Short()`, a missing binary) with a message; `//nolint:<linter> // <reason>` and `//lint:ignore <check> <reason>`, never bare.
```

#### I-15b — skills/review-standards/references/typescript-review.md
```old
cases for exported generics.
```
```new
cases for exported generics.
- [ ] No `it.only` or `describe.only`; `.skip` and `eslint-disable` carry a reason (`-- <reason>`).
```

#### I-15c — skills/golang/SKILL.md
```old
Code must pass `go vet` and `staticcheck`.
```
```new
Code must pass `go vet` and `staticcheck`. A vet, lint or test failure is fixed in the source; a `//nolint` or a skip is for a demonstrated false positive and names the linter and the reason.
```

#### I-15d — skills/golang/references/golang-ecosystem.md
```old
`unused`, `misspell`)
```
```new
`unused`, `misspell`, `nolintlint`)
```

### I-16 — the verifier runs every command, counts tests, names inherited failures and suppressions (G07, G16, G17)

#### I-16a — agents/verifier.md
```old
When `AGENTS.md` names nothing, auto-detect:
```
```new
- **Every command on its own:** run each command as its own call and read its own exit status and output — a chain in `AGENTS.md` (`a && b && c`) is split into its commands, and `gofmt -l` fails by printing, not by exit status; record every command's outcome even after an earlier one failed.
- **Test counts:** the project's full Go test command (never a `-run` narrowing) runs once, with `-v` added and its output saved to `.agents/scratch/verify-<slug>/go-test.txt`; read its exit status, then count the `--- PASS:`, `--- FAIL:` and `--- SKIP:` lines, indented subtest lines included, and write `counts: pass N, fail N, skip N at <git rev-parse --short HEAD>[+dirty]`. Find the newest earlier log in `testing/` with a counts line for the same command and write `vs previous: <path>, pass a→b, skip c→d`. A fall in pass or a rise in skip is a FAIL line unless `git diff <that sha>` explains it (a test removed with its code, a skip added with its reason) or the Environment section does (go or git missing from PATH, no symlink support, another working directory); then the line says `explained: <cause>`. With no earlier counts write `vs previous: none`. A frontend runner's counts are copied from its summary, else `counts: not printed`; a project with neither writes `counts: not applicable`.

When `AGENTS.md` names nothing, auto-detect:
```

#### I-16b — agents/verifier.md
```old
A changed file with no covering test is a FAIL line naming the file — reported, never patched.
```
```new
A changed file with no covering test is a FAIL line naming the file — reported, never patched.
- **Inherited failures:** a failing test whose test file and the code it exercises are untouched by the change set is marked `pre-existing?` with that evidence; Status stays FAIL, and the reply says so, so the main agent numbers it.
- **Suppressions added:** from the change set's diff, list under Coverage of changed code, as file:line, every added `//nolint`, `//lint:ignore`, `@ts-ignore`, `@ts-expect-error` or `eslint-disable` line in source files, every added `t.Skip`, `t.SkipNow` or `//go:build ignore` in `*_test.go`, and every added `.skip(` or `.only(` in `*.test.ts` or `*.spec.ts`; an added `.only(`, or a skip or ignore with no reason, is a FAIL line; a bare `//nolint` is a FAIL only when the project's commands run golangci-lint.
```

#### I-16c — agents/verifier.md
```old
- `go test ./... -race` → exact outcome (counts, failures verbatim)
```
```new
- `go test -race -v ./...` → exit 0; counts: pass N, fail 0, skip N at <sha>; vs previous: <path> (failures verbatim)
```

### I-17 — a cut tool result counts as not examined (G18)

#### I-17a — agents/reviewer.md
```old
naming what you did not reach
```
```new
naming what you did not reach (a tool result that came back cut — a preview, a saved-output path, a line limit — counts as not reached until you have read the rest)
```

#### I-17b — agents/verifier.md
```old
naming what you did not reach
```
```new
naming what you did not reach (a tool result that came back cut — a preview, a saved-output path, a line limit — counts as not reached until you have read the rest)
```

#### I-17c — agents/security.md
```old
naming what you did not reach
```
```new
naming what you did not reach (a tool result that came back cut — a preview, a saved-output path, a line limit — counts as not reached until you have read the rest)
```

#### I-17d — agents/explorer.md
```old
naming what you did not reach
```
```new
naming what you did not reach (a tool result that came back cut — a preview, a saved-output path, a line limit — counts as not reached until you have read the rest)
```

### I-18 — diagnosis misses, graduation tests, graduated index rows (G23, G24)

#### I-18a — agents/bugs.md
```old
The next number comes from `INDEX.md`;
```
```new
The next number comes from `INDEX.md`, both its tables;
```

#### I-18b — agents/bugs.md
```old
Find the next number from `INDEX.md` (or the highest `BUG-*.md` filename)
```
```new
Find the next number from `INDEX.md` — the highest id in either table — or the highest `BUG-*.md` filename
```

#### I-18c — agents/bugs.md
```old
Why it happened — the mechanism, not the symptom again.
```
```new
Why it happened — the mechanism, not the symptom again.
`Diagnosis missed:` <only when the dispatch names one: the cause the goldfish diagnosis did not rank, and the line of its brief that would have found it>
```

#### I-18d — agents/bugs.md
```old
One table, newest last, one row per bug:
```
```new
One main table, newest last, one row per bug, and below it a `## Graduated` table once a rule graduates (see below):
```

#### I-18e — agents/bugs.md
```old
that is the one allowed edit to a past record.
```
```new
that is the one allowed edit to a past record. When a dispatch says a bug's prevention rule graduated into an `AGENTS.md`, a `.agents/rules/` file or a test, move that bug's row from the main table to the `## Graduated` table at the bottom of `INDEX.md` — `| BUG-NNN | area | → <home file> |` — and leave the record file untouched.
```

#### I-18f — agents/bugs.md
```old
or a test). You do not write that rule;
```
```new
or a test), and say whether the rule is actionable (a specific check would catch it) and generalizable (it applies beyond this file or session) — information for whoever approves it, never a reason to withhold the recommendation. You do not write that rule;
```

#### I-18g — agents/bugs.md
```old
the only later edit to a past record is its `status`, in both files.
```
```new
the only later edit to a past record is its `status`, in both files; moving a row to the Graduated table edits the index, not a record.
```

#### I-18h — skills/bug/SKILL.md
```old
that is where the harness learns.
```
```new
that is where the harness learns. When Step 2 ran and the reviewer, the verifier or the fix surfaced a cause the diagnosis did not rank, pass `bugs` a `Diagnosis missed:` line — what it was, and which line of the Step 2 brief would have found it.
```

#### I-18i — skills/retro/SKILL.md
```old
whose prevention rule is not yet written anywhere
```
```new
whose prevention rule is not yet written anywhere, and `Diagnosis missed:` lines in bug records dated on or after the newest `<docs-root>/retros/` record (all records when there is none)
```

#### I-18j — skills/retro/SKILL.md
```old
Show a numbered table — title, scope, target, one-line why — and ask which to apply.
```
```new
Show a numbered table — title, scope, target, one-line why, and for a bug-derived proposal that targets an `AGENTS.md` whether its rule is actionable and generalizable (information, never a filter) — and ask which to apply.
```

#### I-18k — skills/retro/SKILL.md
```old
status `applied` in the record and the index.
```
```new
status `applied` in the record and the index; when the applied rule is a bug's prevention rule, dispatch `bugs` in the same task to move that bug's row to the index's Graduated table.
```

#### I-18l — context/AGENTS.md
```old
A rule that recurs graduates into the project's `AGENTS.md`, a `.agents/rules/` file, or a test.
```
```new
A rule that recurs graduates into the project's `AGENTS.md`, a `.agents/rules/` file, or a test, and the commit that writes it has `bugs` move the matching index rows to the Graduated table.
```

ADR-021, appended to `.agents/docs/architecture/DECISIONS.md` by `changelog`:
```text
## ADR-021: Bug index rows move to a Graduated table when their rule graduates     (2026-09-15) — Accepted

**Context:** The 2026-09-13 improvement review refused capping the bug index as a speculative benefit. Resume reads `INDEX.md` whole in every session, and the 2026-09-15 research counted a 60-row index in styrman-crew (not re-checked). The operator approved group G24 on 2026-09-15 ("go all").

**Decision:** No size or age cap. When a bug's prevention rule graduates into an `AGENTS.md`, a `.agents/rules/` file or a test, the commit that writes it has `bugs` move that bug's row to a `## Graduated` table at the bottom of `INDEX.md` naming the rule's new home; record files are untouched, and new ids continue from the highest id in either table. This reopens the 2026-09-13 refusal for this narrower rule.

**Consequences:** The main table shrinks only as rules graduate. `context/AGENTS.md`, `agents/bugs.md` and `skills/retro/SKILL.md` say so together.
```

### I-19 — a project skill or agent never reuses a harness name (G26)

The test, in the sandbox of `.agents/docs/testing/2026-09-03-verify-sandbox.md`, with this tree installed into the sandbox home. A throwaway git repository `runs/g26-precedence` holds these files, with the frontmatter shown and no body beyond one line:
```text
.agents/skills/review/SKILL.md                name: review             description: "MARKER-G26-PROJECT-SKILL project copy of review"
.agents/skills/g26-control/SKILL.md           name: g26-control        description: "MARKER-G26-CONTROL-SKILL"
.agents/agents/reviewer.md                    name: reviewer           description: "MARKER-G26-PROJECT-AGENT project copy of reviewer"   (Claude Code shape)
.agents/agents/reviewer/agent.md              name: reviewer           description: "MARKER-G26-PROJECT-AGENT project copy of reviewer"   (Antigravity shape, no tools key)
.agents/agents/g26-control-agent.md           name: g26-control-agent  description: "MARKER-G26-CONTROL-AGENT"   (Claude Code shape)
.agents/agents/g26-control-agent/agent.md     name: g26-control-agent  description: "MARKER-G26-CONTROL-AGENT"   (Antigravity shape, no tools key)
.claude/skills -> ../.agents/skills
.claude/agents -> ../.agents/agents
```
Each CLI gets, with the sandbox driver's usual flags: "Do not use any tools. Quote verbatim the description you have for the skill named review, the agent named reviewer, the skill named g26-control and the agent named g26-control-agent." The answers go verbatim into `testing/<date>-verify-g26-precedence.md`, with one line per CLI and layer: control quoted or not, project marker quoted or not.

Outcome B holds only when both CLIs quote both control markers and both project markers. Every other result, including an unquoted control, is outcome A (D12).

#### I-19a — context/AGENTS.md
```old
Both load; the project layer has priority on both vendors.
```
```new-A
Both load, but do not rely on which copy wins a name clash: in a project, a skill or agent never reuses a harness name — extend one with a differently named skill or a `.agents/rules/` file.
```
```new-B
Both load; the project layer has priority on both vendors. Still, in a project a skill or agent never reuses a harness name — extend one with a differently named skill or a `.agents/rules/` file, so harness fixes keep arriving.
```

#### I-19b — skills/retro/SKILL.md
```old
| Skill knowledge or workflow | `.agents/skills/<name>/SKILL.md` |
```
```new
| Skill knowledge or workflow | `.agents/skills/<name>/SKILL.md` (never a harness skill's name; see `sh-agent list`) |
```

#### I-19c — skills/retro/SKILL.md
```old
| Subagent behavior | `.agents/agents/<name>.md` |
```
```new
| Subagent behavior | `.agents/agents/<name>.md` (never a harness agent's name) |
```

#### I-19d — skills/distill/SKILL.md
```old
| Pattern is unique to this project | `.agents/skills/<name>/SKILL.md` |
```
```new
| Pattern is unique to this project | `.agents/skills/<name>/SKILL.md` (never a harness skill's name; see `sh-agent list`) |
```

### I-20 — every bug fix gets the reviewer, whatever its size (G22; one commit with ADR-022)

#### I-20a — context/AGENTS.md
```old
2. **Review:** `reviewer` for 50+ lines or 5+ files,
```
```new
2. **Review:** `reviewer` for 50+ lines, 5+ files, or any fix under the `bug` skill,
```

#### I-20b — skills/bug/SKILL.md
```old
, `reviewer` for 50+ lines or 5+ files with the diff only,
```
```new
, `reviewer` with the diff only, whatever its size (Step 0 already exempted the trivial),
```

#### I-20c — agents/reviewer.md
```old
when it reaches 50 lines or 5 files,
```
```new
when it reaches 50 lines or 5 files or fixes a bug,
```

#### I-20d — skills/preflight/SKILL.md
```old
when that work alone is 50+ lines or 5+ files.
```
```new
when that work alone is 50+ lines or 5+ files or is a bug fix.
```

ADR-022, appended to `.agents/docs/architecture/DECISIONS.md` by `changelog`:
```text
## ADR-022: Every bug fix gets the reviewer, whatever its size     (2026-09-15) — Accepted

**Context:** ADR-009 tiered review cost, and the 50-line / 5-file threshold was kept in the 2026-09-03 research. The bug skill says small diffs hide bugs well, yet the threshold exempted small fixes from review. The 2026-09-15 research (group G22, from elephant-goldfish) recommended reviewing every non-trivial bug fix; the operator approved it on 2026-09-15 ("go all").

**Decision:** Work under the `bug` skill gets a `reviewer` dispatch over its diff whatever its size; the bug skill's Step 0 still exempts the trivial. `context/AGENTS.md` step 2, `skills/bug/SKILL.md`, `agents/reviewer.md` and `skills/preflight/SKILL.md` Step 1 change together. This amends ADR-009's threshold for bug work only.

**Consequences:** One Sonnet reviewer dispatch per non-trivial bug fix. Features and tasks keep the 50-line / 5-file threshold, and so does the release gate's review of the release diff (`skills/preflight/SKILL.md` Step 2), which reads a sub-threshold release diff itself.
```

### I-21 — the release gate fails closed and defines blocking (G20)

#### I-21a — skills/preflight/SKILL.md
```old
decide whether it subsumes the commits in question;
```
```new
decide whether it subsumes the commits in question — a report covers a commit only if its `Range reviewed:` line reaches it and its not-examined section (`Not examined`, `Not checked`, or `What was not checked`) names none of that commit's paths;
```

#### I-21b — skills/preflight/SKILL.md
```old
whose `Range reviewed:` reaches HEAD satisfies this row; cite it.
```
```new
whose `Range reviewed:` reaches HEAD and whose not-examined section names no path in the release diff satisfies this row; cite it.
```

#### I-21c — skills/preflight/SKILL.md
```old
Otherwise dispatch the `verifier` agent now over the release range.
```
```new
Otherwise dispatch the `verifier` agent now over the release range.

The artifacts this step owes are a verification log always, a reviewer report when the release diff reaches the threshold above, and a security report when a sensitive path changed. An owed artifact that is missing, or one that leaves a path in the release diff unexamined, gets one re-dispatch over those paths; still uncovered, it is a blocking finding.
```

#### I-21d — skills/preflight/SKILL.md
```old
Run the `--docs` lens of `/review` over `<base>..HEAD` — the `reviewer` agent with the docs brief and the slug `<base>-<head>-docs` — against `README.md`, `<docs-root>/specs/`, `<docs-root>/testing/<app>/guide.md`, and `<docs-root>/deploy/`.
```
```new
Run the `--docs` lens of `/review` over `<base>..HEAD` — the `reviewer` agent with the docs brief exactly as `/review` Step 3 gives it, and the slug `<base>-<head>-docs`.
```

#### I-21e — skills/preflight/SKILL.md
```old
Any blocking finding from Steps 2–4 means NO-SHIP.
```
```new
Any blocking finding from Steps 2–4 means NO-SHIP. Blocking means any of: a verifier FAIL or BLOCKED; a security or reviewer finding of high severity; a reviewer REQUEST CHANGES verdict; a `REQUIREMENTS.md` violation; an owed artifact still missing, or a path in the release diff still unexamined, after the Step 2 re-dispatch. A project's `AGENTS.md` may add blocking conditions; it cannot remove these.
```

### I-22 — gates written the same way in feature, bug and review (G21)

#### I-22a — skills/feature/SKILL.md
```old
Four points always wait for the user, whatever "work autonomously" instruction is in effect: the scope confirmation (Step 1), the no-code gate until the cold-reader check passes (Step 3), the three-round cap of that check, and a review loop that hits its cap. Everything else runs without asking. The user can waive a gate only in the same turn, explicitly.
```
```new
Four points always wait for the user: the scope confirmation (Step 1), the no-code gate until the cold-reader check passes (Step 3), the three-round cap of that check, and a review loop that hits its cap. Everything else runs without asking. A gate holds whatever "work autonomously" instruction is in effect. Only the user's own words in the same turn waive it, and a waiver names the gate — a general instruction such as "do not ask me anything" is not a waiver, and neither is the skill invocation or a waiver carried over from another plan or session. At a cap, ask with your recommendation first: clarify or cut a requirement; proceed with the residue written as `Open:` rows under `## Decisions`, the phases that depend on them waiting; or abandon. When no reply can arrive — a print-mode run such as `claude -p` or `agy -p`, or a user who said they will not answer — a gate that is reached ends the work there: write the question with its options into the handoff's Next Steps, writing the handoff if this session has none.
```

#### I-22b — skills/bug/SKILL.md
```old
**Docs root** (`<docs-root>`): `.agents/docs` if it exists, else `docs/` if it exists, else create `.agents/docs`. Dates from `date +%F`.
```
```new
**Docs root** (`<docs-root>`): `.agents/docs` if it exists, else `docs/` if it exists, else create `.agents/docs`. Dates from `date +%F`.

## Gates that need the user

Two points always wait for the user: a missing repro when the bug is not obvious from one file (Step 1), and a diagnosis that disagrees with your hypothesis when the evidence does not settle it (Step 2). A gate holds whatever "work autonomously" instruction is in effect. Only the user's own words in the same turn waive it, and a waiver names the gate — a general instruction such as "do not ask me anything" is not a waiver, and neither is the skill invocation or a waiver carried over from another session. Ask with the options, your recommendation first: for a missing repro, give one, let me infer one, or stop; for a diverged diagnosis, go with the diagnosis, keep your hypothesis, or investigate further. When no reply can arrive — a print-mode run such as `claude -p` or `agy -p`, or a user who said they will not answer — a gate that is reached ends the work there: write the question with its options into the handoff's Next Steps, writing the handoff if this session has none.
```

#### I-22c — skills/bug/SKILL.md
```old
**ask the user for one**, whatever autonomous instruction is in effect (steps, URL, log line, failing test). Never invent a repro; a wrong-shaped guess wastes the diagnosis.
```
```new
**ask the user for one** (steps, URL, log line, failing test; a gate, above). Never invent a repro on your own: one is inferred only when the user picks that option at the gate, and it is printed before Step 2 runs, because a wrong-shaped guess wastes the diagnosis.
```

#### I-22d — skills/bug/SKILL.md
```old
Disagreement: read its evidence and re-investigate before choosing;
```
```new
Disagreement: read its evidence and re-investigate before choosing, and when the evidence does not settle it, ask (a gate, above);
```

#### I-22e — skills/review/SKILL.md
```old
A headless run with `--fix` applies; without it, every run only reports.
```
```new
The three-round cap is a gate: it holds whatever "work autonomously" instruction is in effect, and only the user's own words in the same turn, naming it, waive it. At the cap, ask with your recommendation first: accept the remaining findings as `left` with an owner in the Outcome section, run one more round, or abandon the fix. A headless run with `--fix` applies up to the cap; when no reply can arrive there — a print-mode run, or a user who said they will not answer — it stops and writes the question into its record and into the handoff's Next Steps, writing the handoff if the session has none. Without `--fix`, every run only reports.
```

### I-23 — `/init-project` derives commands and folds an existing CLAUDE.md (G25)

#### I-23a — skills/init-project/SKILL.md
```old
ask the user for a one-paragraph project overview, the build/run/test commands, and the deploy target (none, a remote box, or GCP) and fill the sections;
```
```new
ask the user for a one-paragraph project overview and the deploy target (none, a remote box, or GCP); when `go.mod`, a `web/` or `ui/` `package.json`, or a Dockerfile exists, derive the build, run and test commands from them, print them and ask "confirm or correct", otherwise ask for them; delete every template command line the user did not confirm; fill the sections;
```

#### I-23b — skills/init-project/SKILL.md
```old
- `CLAUDE.md` — stays the single line `@AGENTS.md`, so Claude Code reads the agents.md file.
```
```new
- `CLAUDE.md` — the single line `@AGENTS.md`, so Claude Code reads the agents.md file. If it already existed with content, never overwrite it: ask whether to fold that content into `AGENTS.md`; on yes, fold it and make `CLAUDE.md` the single line; on no, keep the content and add `@AGENTS.md` as its first line. Say which in the summary.
```

### I-24 — the skill matrix checks what the E2E criteria claim (G02, G07, G17, G21)

#### I-24a — .agents/docs/testing/skill-matrix.md
```old
&& ! ls internal/metrics
```
```new
&& grep -q "expect:" .agents/docs/plans/*metric* && grep -q "\- \[.\] Checkpoint" .agents/docs/plans/*metric* && grep -q "\- \[.\] Cold-reader check" .agents/docs/plans/*metric* && grep -q "## Decisions" .agents/docs/plans/*metric* && ! ls internal/metrics
```

#### I-24b — .agents/docs/testing/skill-matrix.md
```old
| `review` | demo HEAD, one uncommitted edit |
```
```new
| `bug-no-repro` | demo HEAD |  | /bug The service sometimes feels slow. Do not ask me anything. There is no remote: commit, skip the push. Finish with what you did. | `git diff --quiet origin/HEAD -- cmd internal && [ -z "$(git ls-files --others --exclude-standard -- cmd internal)" ] && grep -qi "repro" "$OUT" && grep -rqi "repro" .agents/docs/handoffs/` |
| `review` | demo HEAD, one uncommitted edit |
```

#### I-24c — .agents/docs/testing/skill-matrix.md
```old
| `ctx-deploy` | demo HEAD |
```
```new
| `verifier-counts` | demo HEAD + an uncommitted `internal/greet/extra_test.go` holding one passing test |  | Dispatch the verifier agent on the uncommitted change and relay its status line. Change nothing else. | `grep -q "counts: pass" .agents/docs/testing/*verify*.md && grep -q "vs previous:" .agents/docs/testing/*verify*.md` |
| `ctx-deploy` | demo HEAD |
```

## Failure modes

- An `old` text does not occur exactly once when its task starts: nothing of that task is written; re-read the file, fix the edit here with a D-row, re-run the rehearsal, then write (BUG-010).
- The installer refuses a file (BUG-004): only `agents/reviewer.md` changes frontmatter, and its description stays one double-quoted line (`AGENTS.md` "Rules for working here").
- `sh-agent render | wc -l` shows other than 93: an edit added or removed a newline; fix it before the checkpoint.
- A CLI log shows a parse error for a skill or agent: the phase stays open until fixed.
- The docs-lens control (1.2) misses a planted finding or flags the negative: the brief does not work; rework I-1a and re-run 1.2 before anything else.
- The one-time sweep (1.6) returns many findings: a finding on a surface this plan edits is fixed in the phase that edits it; any other is a Follow-ups line, or a bug record when it is a real defect; the plan gains no groups.
- A G26 control is not quoted on a CLI: outcome A (D12); the log says which layer did not load.
- A matrix row fails because the demo or a CLI changed rather than the harness: the e2e log says so; it is not recorded as a pass.
- A count comparison falls: a FAIL line unless the diff or the environment explains it, in which case the line names the cause (I-16a, D17).
- A headless run outlives the tool timeout: a detached driver with its own status file, whose exit status is read before any output is cited (BUG-013).

## Verification

Tiers:
- Unit: `go test -race ./...` exits 0 at every checkpoint; prose tests such as `TestNoTmpInHarnessProse` stay green. No Go code changes.
- Integration: the rehearsal passes; `go run ./cmd/sh-agent install --home .agents/scratch/<slug>-home` exits 0; `status --home` says drift none; `list` parses; `render | wc -l` prints 93; each task's check meets its expect.
- E2E: Phase 6 in the sandbox on both CLIs.
- Cloud smoke: none.

E2E criteria:
- Golden path: in the sandbox demo, `/feature` writes a plan with `expect:` on its task lines, a Checkpoint box per phase, a Cold-reader check box and a `## Decisions` section, on both CLIs — yes when the `feature` row's check (I-24a) exits 0 on both.
- Edge case: `/bug` with no repro and "do not ask me anything" changes nothing under `cmd/` or `internal/`, tracked or untracked, and writes the repro question into a handoff, on both CLIs — yes when the `bug-no-repro` row's check (I-24b) exits 0 on both.

Each criterion is a yes/no outcome that still holds if the interface wording were different.

## Phases

- [ ] Cold-reader check — design ready and implementation ready (<reviews/ paths>), or the user's same-turn override quoted with its date

### Phase 1 — the docs lens and the lost rule (G06, G08)
- [ ] 1.1 I-1a — verify: integration, `grep -c 'narrowed: <old file:line>' skills/review/SKILL.md`, who: AI — expect: `1` (positive control: `git show c192622:skills/review/SKILL.md | grep -c 'narrowed:'` prints `0`)
- [ ] 1.2 Docs-lens control, after 1.1: build `.agents/scratch/docs-lens-control/` as below and dispatch `reviewer` there with the I-1a brief verbatim and the scope `Scope: commits base..change (git log base..change); the diff is git diff base...change.` — verify: integration, the report the reviewer writes inside that repository, who: AI — expect: a `dropped:` finding for "Never deploy on Fridays."; a finding naming `README.md` for "from the `main` branch"; a `narrowed:` finding for the linter rule; no finding for "Tag releases with the git short SHA."
- [ ] 1.3 I-2a, then I-6 — verify: integration, `grep -c "whatever autonomous instruction is in effect" skills/bug/SKILL.md` and `grep '^| BUG-018' .agents/docs/bugs/INDEX.md`, who: AI — expect: `1`, and the row's status cell reads `fixed`
- [ ] 1.4 I-3a — verify: integration, `grep -c 'dropped: <why>' skills/retro/SKILL.md`, who: AI — expect: `1`
- [ ] 1.5 I-4a and I-5a–e through `documentation` (D7) — verify: integration, `git grep -nE '27 skills|seven subagents|7 agents parse' -- README.md AGENTS.md .agents/docs/testing/sh-agent/guide.md` and `grep -c 'searched at `skills/` because' AGENTS.md`, who: AI — expect: the first prints nothing (positive control: the same pattern over `c192622` prints lines) and the second prints `1`
- [ ] 1.6 One-time sweep: dispatch `reviewer` with scope `Scope: the whole repository at HEAD (<sha>), no diff.`, the I-1a brief including its `--project` clause, and I-4a's extra paths — verify: integration, the report under `reviews/`, who: AI — expect: a classification of every living doc and closing re-check greps; each finding fixed, filed under Follow-ups, or recorded as a bug per Failure modes
- [ ] Checkpoint — expect: `go test -race ./...` exit 0; scratch install exit 0 and drift none; `render | wc -l` prints 93; 1.2's expectations met; a `verifier` log over the phase range with Status PASS — evidence: <commands, exit status, report and log paths>, <date>

The control repository of 1.2, created with `git init` under `.agents/scratch/docs-lens-control/` (reference text, not edits):
```text
commit base
  AGENTS.md
    # AGENTS.md — control
    ## Rules
    - Never deploy on Fridays.
    - Deploy only from the `main` branch.
    - Always run the linter before every commit, even for docs-only changes.
    - Tag releases with the git short SHA.
  README.md
    # control
    Deploys go out from the `main` branch.
commit change (tag it `change`; tag the first `base`)
  AGENTS.md
    # AGENTS.md — control
    ## Rules
    - Deploy only from a release tag.
    - See `.agents/skills/release/SKILL.md` for release steps.
  .agents/skills/release/SKILL.md
    ---
    name: release
    description: "Release steps for the control repository."
    ---
    - Always run the linter before every commit.
    - Tag releases with the git short SHA.
  README.md unchanged
```

### Phase 2 — plans, phases and memory (G01, G02, G03, G04, G05, G09, G12, G13, G14, G19)
- [ ] 2.1 I-7a and I-7b — verify: integration, for each pattern in "— expect:", "- [ ] Checkpoint —", "## Decisions", "- [ ] Cold-reader check —", "— in progress:", "recon-<slug>", "phase-<N>", "close the plan": `grep -cF -- "<pattern>" skills/feature/SKILL.md`, plus `grep -cF "then tick the boxes in the plan" skills/feature/SKILL.md`, who: AI — expect: every pattern prints at least `1`, and the last prints `0` (positive control: over `c192622` every pattern prints `0` and the last prints `1`)
- [ ] 2.2 I-8a–e — verify: integration, `go run ./cmd/sh-agent render` into `.agents/scratch/render.txt`, then `grep -cF "restate status and next step before the next edit"`, `grep -cF "Workflow skills, also typed"`, `grep -cF "can also be typed as"` on it and `wc -l` on it, who: AI — expect: `1`, `0`, `1`, `93`
- [ ] 2.3 I-9a–d and I-10a–b — verify: integration, `grep -cF "Ask for six things" skills/handoff/SKILL.md` and `git grep -nF 'should put in `TODO.md`' -- agents skills context`, who: AI — expect: `1`, and nothing printed (positive control: over `c192622` it prints `agents/changelog.md:94`)
- [ ] 2.4 I-11a–g — verify: integration, `grep -l '^- \[ \]' .agents/docs/plans/*.md`, who: AI — expect: only this plan's path
- [ ] 2.5 Describing text — verify: integration, `git grep -nE "tick its tasks|then tick the boxes|Workflow skills, also typed|should put in .TODO" -- context skills agents README.md AGENTS.md`, who: AI — expect: nothing printed (positive control: over `c192622` at least three lines)
- [ ] Checkpoint — expect: `go test -race ./...` exit 0; drift none; `render | wc -l` prints 93; a `verifier` log over the phase range with Status PASS — evidence: <commands, exit status, log path>, <date>

### Phase 3 — review loop and review standards (G10, G11, G15, G16 skill half)
- [ ] 3.1 I-12a–b, then I-13a–b — verify: integration, `grep -cF "that section's fixed lines only" skills/review/SKILL.md`, `grep -cF "never a command that exited non-zero" skills/review/SKILL.md`, `grep -cF 'a reason as `/review` Step 4 defines it' skills/preflight/SKILL.md`, who: AI — expect: `1`, `1`, `1` (positive control: each prints `0` over `c192622`)
- [ ] 3.2 I-14a–c and I-15a–d — verify: integration, `grep -cF` for "Leftover debug output", "TOCTOU race", "generated counterpart", "assertion changed in an existing test" and "describing-text rule" in `skills/review-standards/SKILL.md`; "lint:ignore" in `go-review.md`; "describe.only" in `typescript-review.md`; "demonstrated false positive" in `skills/golang/SKILL.md`; and `sed -n 323p skills/golang/references/golang-ecosystem.md` piped to `grep -cF nolintlint`, who: AI — expect: every grep prints `1`
- [ ] Checkpoint — expect: tests exit 0; drift none; a `verifier` log over the phase range with Status PASS — evidence: <log path>, <date>

### Phase 4 — agents and records (G07, G16 verifier half, G17, G18, G23, G24)
- [ ] 4.1 I-16a–c — verify: integration, `grep -cF` for "Every command on its own", "Test counts:", "Inherited failures:" and "Suppressions added:" in `agents/verifier.md`, who: AI — expect: each prints `1`
- [ ] 4.2 I-17a–d — verify: integration, `grep -cF "counts as not reached until you have read the rest" agents/reviewer.md agents/verifier.md agents/security.md agents/explorer.md`, who: AI — expect: `1` for each file
- [ ] 4.3 I-18a–l — verify: integration, `grep -cF "Diagnosis missed:"` in `agents/bugs.md`, `skills/bug/SKILL.md` and `skills/retro/SKILL.md`; `grep -cF "Graduated"` in `agents/bugs.md`, `skills/retro/SKILL.md` and `context/AGENTS.md`, who: AI — expect: at least `1` in each
- [ ] 4.4 ADR-021 through `changelog`, with I-18's text — verify: integration, `grep -c '^## ADR-021: Bug index rows move' .agents/docs/architecture/DECISIONS.md`, who: AI — expect: `1`
- [ ] 4.5 The derived Antigravity agents carry I-16 and I-17 — verify: integration, scratch install, then `grep -cF "counts as not reached"` in `<scratch>/.gemini/config/agents/{reviewer,verifier,security,explorer}/agent.md`, who: AI — expect: `1` for each file
- [ ] Checkpoint — expect: tests exit 0; drift none; a `verifier` log over the phase range with Status PASS, dispatched with the repo's § What to run pasted (D9), carrying a `counts:` line and a `vs previous:` line — evidence: <log path>, <date>

### Phase 5 — gates, release gate, init and names (G20, G21, G22, G25, G26)
- [ ] 5.1 The I-19 precedence test — verify: e2e, `testing/<date>-verify-g26-precedence.md`, who: AI — expect: one line per CLI and layer; the outcome, A or B, is a D-row
- [ ] 5.2 I-19a with that outcome's `new`, then I-19b–d — verify: integration, the render file's `grep -cF "never reuses a harness name"`, `grep -cF "never a harness" skills/retro/SKILL.md skills/distill/SKILL.md`, who: AI — expect: `1`; `2` and `1`
- [ ] 5.3 I-20a–d with ADR-022, all in the phase commit — verify: integration, `git grep -nF "50+ lines or 5+ files" -- context skills agents` and `grep -c '^## ADR-022: Every bug fix gets the reviewer' .agents/docs/architecture/DECISIONS.md`, who: AI — expect: exactly two lines, both in `skills/preflight/SKILL.md`, the Step 1 line ending "or is a bug fix." and the Step 2 release-diff line (D8); `1`
- [ ] 5.4 I-21a–e — verify: integration, `grep -cF` for "Blocking means any of", "names no path in the release diff", "The artifacts this step owes" and "exactly as `/review` Step 3 gives it" in `skills/preflight/SKILL.md`, who: AI — expect: each prints `1`
- [ ] 5.5 I-22a–e — verify: integration, `grep -c "^## Gates that need the user" skills/feature/SKILL.md skills/bug/SKILL.md`, `grep -cF "The three-round cap is a gate" skills/review/SKILL.md`, `grep -cF "whatever autonomous instruction is in effect" skills/bug/SKILL.md`, `grep -cF "Never invent a repro on your own" skills/bug/SKILL.md`, who: AI — expect: `1` each; `1`; `0`; `1`
- [ ] 5.6 I-23a–b — verify: integration, `grep -cF` for "derive the build, run and test commands" and "add `@AGENTS.md` as its first line" in `skills/init-project/SKILL.md`, who: AI — expect: each prints `1`
- [ ] 5.7 The reviewer description reaches both vendors — verify: integration, scratch install exit 0, then `grep -cF "or fixes a bug" <scratch>/.claude/agents/reviewer.md <scratch>/.gemini/config/agents/reviewer/agent.md`, who: AI — expect: `1` for each file
- [ ] Checkpoint — expect: tests exit 0; drift none; `render | wc -l` prints 93; a `verifier` log over the phase range with Status PASS (D9) — evidence: <log path>, <date>

### Phase 6 — end to end in the sandbox, both CLIs
- [ ] 6.1 I-24a–c — verify: integration, `grep -cF` for "`bug-no-repro`" and "`verifier-counts`" in `.agents/docs/testing/skill-matrix.md`, who: AI — expect: each prints `1`
- [ ] 6.2 Install this tree into the sandbox home and check drift — verify: e2e, `sh-agent install --home <sandbox home>` then `sh-agent status --home <sandbox home>`, who: AI — expect: exit 0 and drift none
- [ ] 6.3 One start of each CLI with a context question: quote Definition of Done step 2's reviewer trigger and the compaction sentence — verify: e2e, the two outputs and the newest Antigravity `cli-*.log`, who: AI — expect: both quote "or any fix under the `bug` skill" and "restate status and next step"; `grep -c "Failed to parse"` on the log prints `0`
- [ ] 6.4 Matrix rows on both CLIs, detached with a status file, after `git -C <a fresh clone> rev-parse --verify origin/HEAD` succeeds: `feature`, `bug`, `bug-no-repro`, `review`, `review-project`, `preflight`, `retro`, `distill`, `handoff`, `init-deploy`, `explorer`, `review-standards`, `golang`, `verifier-counts` — verify: e2e, each row's check in its clone, who: AI — expect: every check exits 0 on both CLIs; a failure gets a bug record and a fix, or a `[-]` box with the reason and an `Open:` row
- [ ] Checkpoint — expect: 6.2 to 6.4 pass, recorded in `testing/<date>-verify-harness-practices-e2e.md` — evidence: <log path>, <date>

### Done — the Definition of Done
- [ ] `verifier` over `c192622..HEAD`, every tier above, dispatched with the repo's § What to run pasted (D9) — who: AI
- [ ] `reviewer` over the diff by commit range, and a `--docs` review over the same range with the I-1a brief and I-4a's paths — who: AI
- [ ] `documentation` for anything the reviews name in README, the testing guide or the deploy guides — who: AI
- [ ] `changelog` naming the agents dispatched; `bugs` for anything that went wrong — who: AI
- [ ] Close the plan; `ROADMAP.md` "Now" and `TODO.md` empty — who: AI
- [ ] Commit, push, `/handoff` with its cold read. Its Next Steps name, for the operator: on this box, `sh-agent update` (or `go run ./cmd/sh-agent install`) and `sh-agent status`; on the Mac, `sh-agent update --check`, `sh-agent update`, `sh-agent status`, Claude Code `/memory` and `/skills`, Antigravity `/skills` and `/agents`; on both, the three smoke checks of `.agents/docs/testing/2026-09-03-verify-real-home.md` on each CLI; then the downstream mirror refresh (`deploy/remote.md` § Downstream mirrors) — who: AI

## Risks

Assumptions:
- Appending text to an existing context line adds no rendered line — checked by 2.2.
- `skills/golang` and `skills/review-standards` are not vendored — checked on 2026-09-15 against `skills/VENDORED.md`.
- Both CLIs load project skills and agents from `.agents/` and the `.claude/` symlinks, in the shapes I-19 writes — checked by 5.1's control markers.
- A fresh clone of the sandbox demo has `origin/HEAD` — checked in 6.4 before the rows run.
- The sandbox home and both CLIs' logins still work — checked by 6.2 before 6.4.

Other risks: `skills/feature/SKILL.md` grows from 82 to about 100 lines and Step 4 to four paragraphs. The Phase 6 matrix takes hours on Claude Code. The one-time sweep can surface more than this plan fixes.

## Decisions

D1 — 2026-09-15 — scope is G01–G26 of the research record — by: "go with the first bundle g01-g26, start /feature"; "go all, use your default thresholds for G20"
D2 — 2026-09-15 — G20's blocking thresholds: a verifier FAIL or BLOCKED; a security or reviewer finding of high severity; a reviewer REQUEST CHANGES; a `REQUIREMENTS.md` violation; a release path no report examined — by: "use your default thresholds for G20"
D3 — 2026-09-15 — G24 lands with ADR-021 (Phase 4), reopening the 2026-09-13 refusal to cap the bug index; G22 lands with ADR-022 (Phase 5), amending ADR-009's threshold for bug work; the numbers follow landing order — by: "go all"; cold read round 1
D4 — 2026-09-15 — the first sentence of context line 71 goes, to offset I-8b and I-8c's words; its fact lives in the generated inventory (`internal/harness/inventory.go:321`) — by: implementer, from the research set judge's G01
D5 — 2026-09-15 — the feature skill's Steps 1–4 and 6 are specified as range replacements, because seven groups edit Step 4 — by: implementer
D6 — 2026-09-15 — this plan already uses the new template shape (expect clauses, Checkpoint boxes, a Decisions section), so it exercises what it builds — by: implementer
D7 — 2026-09-15 — `README.md`, this repo's `AGENTS.md` and the testing guide are edited by `documentation` with the exact texts of I-4a and I-5a–e; if it declines a section outside its lane, the main agent applies the same text and the changes entry says so — by: implementer; cold read round 1
D8 — 2026-09-15 — the release gate's release-diff reviewer threshold (`skills/preflight/SKILL.md` Step 2) stays; a bug fix is reviewed in its own task (I-20), and the gate reads a sub-threshold release diff itself — by: cold read round 1
D9 — 2026-09-15 — in-session dispatches whose behaviour this plan changes carry the changed text in their briefs, because this checkout's session loads the installed definitions; Phase 6 exercises the installed new definitions in the sandbox — by: cold read round 1
D10 — 2026-09-15 — the docs-lens control runs in a scratch repository with three planted findings and one negative, away from this repo's bug index — by: cold read round 1
D11 — 2026-09-15 — `I-1a` greps at the scope's head, follows a symlinked directory to its target, and reports `narrowed:` for a rule that kept its home but lost a qualifier (BUG-018's mechanism) — by: cold read round 1
D12 — 2026-09-15 — the G26 test carries skill and agent controls in both vendor shapes; only a full B result keeps line 37's claim, and every other result writes the cause-neutral outcome A — by: cold read round 1
D13 — 2026-09-15 — "no reply can arrive" (a print-mode run, or a user who said they will not answer) is when a reached gate ends the work and writes a handoff; the `bug-no-repro` check also catches untracked files — by: cold read round 1
D14 — 2026-09-15 — the release gate names the artifacts it owes, accepts the three not-examined section names the agents and records use, and lets a project add blocking conditions but never remove them — by: cold read round 1
D15 — 2026-09-15 — every edit is a fenced old/new block in application order, and a rehearsal applies them all to copies of the files at `c192622` before implementation. Run on 2026-09-15 for this revision: 82 edits parsed, 0 failures, 26 files, for both I-19a outcomes; against the rehearsed end state every task expectation above holds, `sh-agent render` prints 93 lines, `sh-agent list` parses 27 skills and 7 agents, and no line containing `/tmp` was added — by: cold read round 1
D16 — 2026-09-15 — BUG-018 gets a status change only (`agents/bugs.md` allows no other edit); the Phase 1 and Phase 5 changes entries carry its fix history — by: cold read round 1
D17 — 2026-09-15 — a count fall that the diff or the environment explains is not a FAIL; its line names the cause — by: cold read round 1
D18 — 2026-09-15 — every phase checkpoint, Phases 1 and 2 included, gets a verifier log, because context, skill and agent prose is this repo's source — by: cold read round 1
D19 — 2026-09-15 — the stale operator boxes point at this plan's Done list, whose last row names each of their checks — by: cold read round 1
D20 — 2026-09-15 — the bug gate's "let me infer one" option and the "never invent a repro" sentence are reconciled in I-22c — by: cold read round 1
D21 — 2026-09-15 — `agents/bugs.md` describing text (next number, one table, the one allowed edit) and `context/AGENTS.md`'s graduation sentence change with G24, so the row move has a trigger outside `/retro` — by: cold read round 1
Open: findings of the one-time sweep (1.6) outside this plan's surfaces — default: Follow-ups lines, or a bug record for a real defect; no new groups
Open: G26 outcome A or B — default: decided by 5.1; any unquoted control means A
Rejected: several smaller plans — the order constraints (G02 before G03, G06 before G21, G22's four lines in one commit) belong in one place
Rejected: a per-phase `documentation` dispatch for docs clauses — the main agent ticks docs clauses, and `documentation` runs at Done
