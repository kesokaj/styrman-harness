---
owner: main
date: 2026-09-15
---
# Research — beads, elephant-goldfish and goldfish against this harness   (2026-09-15)

## Question

The operator asked for three repositories to be researched in full and weighed against this harness: `gastownhall/beads`, `vshvedov/elephant-goldfish` and `lukacf/goldfish`. Can their skills and ways of working merge into what the harness already has? The interests named were drift, more tasks and phases, more and smaller verifications, and planning. The deliverable is a table of add or not, each row with a short explanation, a recommendation, and what the change would do.

## Method

- Shallow clones in `.agents/scratch/research/` (gitignored): beads `f56632a` (2026-09-12), elephant-goldfish `b8ebb3d` (2026-05-12), goldfish `5dbe820` (2026-07-08).
- Six read-only readers: two for beads, one for elephant-goldfish, two for goldfish, one mapping the harness. One proposer per repo then wrote twenty add, adapt or skip candidates each, sixty in all.
- Three independent skeptics checked every candidate, each prompted to refute. The harness lens asked whether the harness already does it. The source lens asked whether the repo really does it and what it needs. The fit lens checked R1 to R7, the non-goals and the ADR-009 context budget. A completeness critic, a second round for uncovered mechanisms and one judge of the whole set followed.
- The critic named eight mechanisms no candidate covered. A second round wrote thirteen more candidates, eight from beads and five from goldfish, and the same three lenses checked them. A set judge then merged all seventy-three into forty groups, measured the cumulative cost per file, listed conflicts and ordered a first bundle.
- The first run's verification stage failed on a model usage limit. A resume re-ran finished agents, so it was stopped, and a new run read the cached candidates from files; the failure is recorded as BUG-017. Readers and first-round proposers ran on Fable 5.1; everything after ran on Opus 5.
- The main agent spot-checked the footprint facts below against the clones, and the harness findings under "Drift found in the harness itself" against git history.

## The three repositories

**beads (`bd`), MIT.** A Go CLI issue tracker for coding agents. Every unit of work is an issue in a Dolt SQL database under `.beads/`; typed dependency edges form a graph; `bd ready` returns the claimable frontier; agents claim, close, and file `discovered-from` work. Agent behaviour is injected by `bd prime` at session start and after compaction, through a `SessionStart` hook that `bd setup claude` writes into `.claude/settings.json` (`cmd/bd/setup/claude.go:63,279-287`). Embedded Dolt needs a CGO build (`beads_nocgo.go:45`), Dolt is the only backend (`docs/cli-reference/init.md:13`), and `bd init` manages five git hooks (`cmd/bd/hooks.go:24`). Its agent skill (`plugins/beads/skills/beads/`) carries practices that work without the tool: the two-week test (`SKILL.md:23`), overwrite-only handoff notes (`resources/WORKFLOWS.md:368`), and ready fronts instead of phases (`resources/WORKFLOWS.md:147`).

**elephant-goldfish, MIT.** Prose only: 31 markdown and yaml files, no script and no hook (`git ls-files`). Five commands per vendor implement the elephant and goldfish pattern: brainstorm, prd, new-feature, fix-bug and precommit-review. The elephant is the full-context session; a goldfish is a fresh subagent given only an artifact. Its core was researched and largely adopted on 2026-09-03 (`.agents/docs/plans/2026-09-03-research-harness-strategies.md`). This second read looked for what the first one missed. The repo states its own cost: "The price is the higher token consumption and longer irrational work. Be warned." (`README.md:17`).

**goldfish (`goldfish-ml`), AGPL-3.0.** Not an elephant and goldfish kit: a Python MCP server for agentic machine-learning research. It exposes 43 MCP tools (`@mcp.tool` in `src/goldfish/server_tools/`) and keeps SQLite inside a hidden sibling git repository, run by a daemon with Docker or GCP execution. Runs are contracts: `results_spec` before, `finalize_run` after, and a gate that blocks new runs while a finished run is unfinalized (`src/goldfish/experiment_model/records.py:1778`). A regenerated `STATE.md` is "Claude's primary context recovery mechanism after compaction" (`src/goldfish/state/state_md.py:1-4`). Its AI pre-run review approves on timeout or error (`src/goldfish/pre_run_review.py:340,351`). None of the machinery transfers. Its process ideas do, mostly from its own `CLAUDE.md` (test-first, "No exceptions", `CLAUDE.md:25,50`) and the historical plans in `docs/archive/`.

Licensing: every adopted idea is restated in the harness's own words. MIT text copied verbatim would need its notice; AGPL text is not copied.

## The harness today, where the candidates land

- The rendered global context is 93 of its 100-line budget (`.agents/docs/retros/2026-09-14-retro.md:44`), and ADR-009 keeps procedures in skills. A candidate that can live in a skill or an agent file does not go in the context.
- Checkpoints happen per plan phase (`context/AGENTS.md:11`); a task line carries a verify clause (`skills/feature/SKILL.md:53`); the verifier runs once per task at Done (`context/AGENTS.md:15,21`).
- The cold-reader check is one reviewer dispatch with two questions (`skills/feature/SKILL.md:64-68`); `/review --fix` re-dispatches up to three rounds (`skills/review/SKILL.md:60`).
- Recurring failures the candidates target: describing text that outlived a moved rule (BUG-003, BUG-007), the Definition of Done skipped mid-session (BUG-008), verification cited without reading its status (BUG-011, BUG-013), and environment assumptions found late (BUG-015, BUG-016).

## Recommendation table

Forty groups cover all seventy-three candidates. **Add now** is a file-drop edit with verified value and no rule conflict. **Add next** is worth doing but costs a dispatch, grows a file every feature fills, or needs an operator decision. **Later** names its trigger. **Skip** conflicts with a requirement or is already present. Each row states the change as corrected by the verification lenses. Candidate ids point at `.agents/scratch/research/` (scratch, not committed).

### Add now

| # | Change | From | What the change would do | Where | Candidates |
|---|---|---|---|---|---|
| G01 | The trail outranks the compaction summary and the CLI's own memory | beads | After a compaction the agent re-reads the active plan and `git status` and restates status and next step before editing. A fact the next session needs goes to the trail file that owns it, never only to the CLI's per-machine memory. Both edits land on line 11; the first sentence of line 71 is cut to offset the words, so no rendered line is added. | `context/AGENTS.md:11`, `:71` | BD-01, BD-14 |
| G02 | Each task and phase states its expected result and closes on evidence | goldfish, beads | A task's verify clause says what the command prints when the task is done; an empty result names its positive control. Each phase ends with a Checkpoint box holding the expected end state and the evidence (log path, or command and exit status). A miss leaves it unticked, so the next phase cannot start. Verification criteria are yes/no outcomes that survive a different implementation. | `skills/feature/SKILL.md` | GF-01, GF-02, BD-08 |
| G03 | Mid-phase progress and discovered work live in the plan | beads | A task box is ticked when its verification meets its expectation, so a resume lands on the task actually in progress. At a stopping point the open box carries where work stopped and any blocker. Work found mid-phase is filed when found: a blocker as a box above the task it blocks, deferrable work as a Follow-ups line naming where it was found. Handoff, changelog and the bug skill point at the plan instead of TODO.md or the transcript. | `skills/feature`, `skills/bug`, `skills/handoff`, `agents/changelog.md`, `context/AGENTS.md:11` | BD-02, BD-04 |
| G04 | One Decisions section and a cold-reader box in the plan | goldfish, elephant-goldfish | The template gains `## Decisions`: decisions as they are taken, open questions with the implementer's default, cold-read rebuttals and rejected alternatives. The first box under Phases is the cold-reader check, so the no-code gate survives a resume until both tokens are ready or a same-turn override is quoted on it. Step 1 asks one numbered round when a request leaves gaps. Four existing plans already carry a Decisions section by hand. | `skills/feature/SKILL.md`, `skills/handoff/SKILL.md` | GF-05, GF-06, EG-06, EG-10 |
| G05 | Plans close explicitly; the handoff cold read catches trail drift | beads, goldfish | A box that will not be done becomes `[-]` with dropped, superseded or moved-to-handoff and a date; closing a plan turns every `[ ]` into `[x]` or `[-]`. The handoff's cold reader names any plan with open boxes that ROADMAP does not name, and any source commit without a `changes/` entry. The seven open operator boxes in two finished plans are closed once. | `skills/feature`, `skills/handoff`, two plan files | BD-06, BD-R2-03, GF-14 |
| G06 | The docs lens hunts old wording and dropped rules | beads, goldfish | For each rule the scope changes or removes, `/review --docs` greps its old wording in living docs and reports each stale hit with both file:line, plus a `dropped:` row for a rule that left with no new home. Under `--project` it lists stale references and ends with the greps that re-check them. The same change restores the clause the 87a862d trim dropped from the missing-repro gate. | `skills/review`, `AGENTS.md`, `skills/retro`, `skills/bug` | BD-09, GF-10, BD-R2-02, BD-13 |
| G07 | The verifier runs every command on its own | elephant-goldfish | Each command runs as its own call with its own exit status; a chained setup line is split, and `gofmt -l` is judged by its output. Every outcome is logged even after an earlier failure. A failing test the change never touched is marked `pre-existing?`; Status stays FAIL and the main agent records a bug. | `agents/verifier.md` | EG-04 |
| G08 | Hand-typed inventory counts become pointers | beads | README, the repo AGENTS.md and the testing guide stop typing skill and agent counts; `sh-agent list` and the generated inventory carry them. A trivial docs-only edit. | `README.md`, `AGENTS.md`, testing guide | BD-R2-05 |
| G09 | Test first on feature tasks | goldfish | A task whose verify clause names a test gets that test written and seen failing for the task's reason before the code; in a compiled language a zero-value signature lands first, because a compile error is not the red. | `skills/feature/SKILL.md` | GF-04 |
| G10 | Review rounds keep a ledger and hand over the fix list | elephant-goldfish | Before each `/review --fix` re-dispatch the record gains `## Round N`, one line per finding: fixed, rebutted or left. The re-dispatched reviewer receives only the fixed lines and "verify each fix; look for defects the fixes introduced". Rebuttal reasons stay out of the dispatch, so no author intent reaches the reviewer. | `skills/review/SKILL.md` | EG-01, EG-02 |
| G11 | A rebuttal needs a checkable reason | elephant-goldfish, goldfish | An Outcome `rebutted` line cites file:line, an ADR or an AGENTS.md section, or the user's quoted words; intent is not a reason. An out-of-scope finding becomes `left` with an owner. A command that exited non-zero is never rebutted. | `skills/review`, `skills/preflight` | EG-05, GF-07 |

### Add next

| # | Change | From | What the change would do, and why it waits | Candidates |
|---|---|---|---|---|
| G12 | A verifier log per phase | goldfish | In a plan with two or more phases, each phase that changed source gets its own verifier run as its Checkpoint evidence. One Sonnet dispatch per phase; land after G02 shows the Checkpoint box is filled honestly. | GF-03 |
| G13 | Recon before the plan | elephant-goldfish | For an unread area or about ten files, `explorer` maps the closest surfaces first and the plan cites its brief. One dispatch per larger feature. | EG-09 |
| G14 | Environment assumptions and required docs in the plan | goldfish, beads | Risks lists each assumption about an external system with the command that checks it in Phase 1. A task changing a command, flag, endpoint, config key or deploy step names its doc file and stays unticked until that file is updated. Grows the template right after G02 and G04. | GF-08, BD-07 |
| G15 | Reviewer hunt items | elephant-goldfish | Review standards gain leftover debug output, masked fallbacks and partial writes, missing regenerated files, and stale comments and test names. Read on every review; no incident recorded yet. | EG-14 |
| G16 | Silenced checks become visible | elephant-goldfish, goldfish | A changed assertion with no source change, and a new skip, `.only`, `//nolint` or `eslint-disable` without a reason, become review findings; the verifier lists suppressions the diff adds. No incident recorded yet. | EG-03, GF-R2-02, GF-R2-03, GF-R2-05 |
| G17 | Test counts compared with the previous log | goldfish | Verifier logs record pass, fail and skip counts at a commit and flag a fall the diff does not explain. False alarms across machines are likely until real logs carry counts; land after G07. | GF-R2-01 |
| G18 | A cut tool result counts as not examined | beads | Reviewer, security, verifier and explorer treat a truncated output as unread until they read the rest. The one recorded near-miss would not have been caught by it. | BD-R2-01 |
| G19 | A Task that outgrows one sitting becomes a Feature | beads | The remainder moves into a plan under the feature gates. About eighteen words in every session, and the feature skill's description already triggers on design questions. | BD-03 |
| G20 | The release gate fails closed | goldfish | A report covers a path only if its Not examined section leaves it out; a gap gets one re-dispatch, then NO-SHIP. "Blocking" is defined from verdict and severity fields, with thresholds the operator picks. | GF-13, GF-17 |
| G21 | Same gate wording in feature, bug and review | elephant-goldfish | Each skill names its gates; a waiver is the user's words in the same turn; a cap offers clarify or cut, proceed with residue, or abandon; a headless run stops and writes the question into the handoff. Effort M plus matrix re-runs on both CLIs. | EG-07 |
| G22 | Every bug fix gets the reviewer | elephant-goldfish | Amends the 50-line threshold for bug-skill work: four files and an ADR in one commit, one Sonnet dispatch per bug. | EG-19 |
| G23 | Diagnosis misses and graduation tests | elephant-goldfish, goldfish | Bug records name causes the goldfish diagnosis missed; a proposed prevention rule states whether it is actionable and general. Rarely fires: one of sixteen records ran a diagnosis. | EG-13, GF-09 |
| G24 | Bug index rows shrink when their rule graduates | beads | A graduated rule's row moves to a short table naming its new home. Reopens the approved 2026-09-13 refusal to cap the index. | BD-11 |
| G25 | `/init-project` derives commands and folds CLAUDE.md | elephant-goldfish, beads | Derived build, run and test commands are confirmed before they replace template defaults; an existing CLAUDE.md is folded into AGENTS.md instead of skipped. Fires only at init. | EG-18, BD-20 |
| G26 | A project skill never reuses a harness name | beads | Replaces the line-37 claim that the project layer has priority on both vendors, which Claude Code's docs report as untrue for skills. Test how each CLI resolves a same-named project skill in the sandbox first. | BD-R2-06 |

### Later

| # | Change | From | Trigger | Candidates |
|---|---|---|---|---|
| G27 | A comprehension paraphrase in the cold read | elephant-goldfish | A plan misread that the Gaps question missed. | EG-08 |
| G28 | No runnable destructive one-liner in a handoff's Next Steps | beads | A destructive step executed from the trail. Only ever in the handoff Writing Rules: in the context it would contradict `context/dev-setup.md:8`. | BD-10 |
| G29 | A cross-vendor second opinion on reviews | goldfish | A review missing what the other vendor's model catches. | GF-19 |
| G30 | An upgrade test from an older manifest shape | beads | A binary change breaking a home installed by an earlier commit. The plain-copy leg is already covered. | BD-R2-04 |
| G31 | The verifier names stale generated output | elephant-goldfish | A consumer shipping stale generated files, and a decision on who runs the generator. | EG-15 |
| G32 | Plan template extras: surface obligations, a Flow section, an ordering check | elephant-goldfish, beads | Trail evidence for any of them; the template already takes about sixteen lines from this set. | EG-16, EG-17, BD-05 |

### Skip

| # | Mechanism | Why not | Candidates |
|---|---|---|---|
| G33 | Adopting `bd` or goldfish as tools | `bd` needs a Dolt database per project, git hooks and a settings.json hook, and its `bd prime` forbids markdown plans. Goldfish is an AGPL Python MCP server with SQLite and a daemon. Both conflict with R3, R5, R7, ADR-007 and the settings non-goal. | BD-16, GF-11 |
| G34 | Anything that needs a hook, CI job, daemon or clock | R3, ADR-007 and the hooks and settings non-goals; the prose halves are in G01 and G05. | BD-15, GF-12, GF-20, BD-R2-07 |
| G35 | Already present | The close protocol is Definition of Done step 5; bug records plus `changes/` entries already carry what goldfish's changelog does. | BD-18, GF-18 |
| G36 | Concurrency primitives and execution hints | No concurrent writers to the trail; model hints cannot be honoured on Antigravity. | BD-17, BD-19 |
| G37 | Contradicts a recorded decision | Stopping before commit (the commit is the checkpoint, ADR-008), per-project template rewriting (ADR-008, R4), uncapped unanimous panels (the three-round cap). | EG-12, EG-20, GF-15 |
| G38 | Ceremony with no recorded need | A brainstorm lens panel costs three to six dispatches per use; mandatory reason strings cannot be enforced as prose and would leak intent to reviewers. | EG-11, GF-16 |
| G39 | Declared baselines and trend alerts | Test counts are exact, so the previous log is the right reference (G17). | GF-R2-04 |
| G40 | Version stamps on skills or handoffs | Skills and binary ship from one commit; `sh-agent status` and each mirror's COMMIT file already show staleness. | BD-12, BD-R2-08 |

## Suggested first bundle, in order

1. G06 first, because it reviews the edits that follow and restores the lost repro clause.
2. G01, before G03, since both edit `context/AGENTS.md:11`.
3. G02, the Checkpoint box that G03 and G05 build on.
4. G03 and G04 in the same pass as G02, so feature Step 4 is rewritten once.
5. G05, after G02 defines the box states.
6. G07, the prerequisite for G17.
7. G08, which removes text and can ride along in any commit.
8. Then G09, and G10 with G11 as one review-skill change.

## Drift found in the harness itself

- The 2026-09-03 context trim (87a862d) removed "never skipped on the strength of an 'autonomous' instruction" from the missing-repro gate. `skills/bug/SKILL.md:29` still asks for a repro but no longer says an autonomy instruction cannot skip it. Verified with `git show`.
- Seven operator boxes remain unticked in two finished plans: six in `2026-09-03-rebuild.md`, one in `2026-09-13-learning-loop-project-level.md`, while ROADMAP "Now" is empty. Verified with `grep`.
- `context/AGENTS.md:37` says the project layer has priority on both vendors. A verification lens reports that Claude Code's documentation resolves skills personal-over-project; this was not tested.
- The repo's own setup line (`AGENTS.md:12`) chains `gofmt -l`, `go vet` and `go test` with `&&`. `gofmt -l` exits 0 while listing unformatted files, and a vet failure hides the test result.
- Claude Code's per-machine memory for this repo holds a status fact about a downstream mirror that lives nowhere in the trail.
- Skill and agent counts are typed by hand in README, the repo AGENTS.md and the testing guide, while the deployed inventory is generated.

## Cost of the set

- `context/AGENTS.md` stays at 93 rendered lines: every edit lands on an existing line. Words are the real cost. G01's words are offset by the line-71 cut; G19, G22 and G26 would add about seventy words and stay add next.
- `skills/feature/SKILL.md` grows from 82 to about 98 lines, and Step 4 passes 300 words; split it into short paragraphs when G02, G03, G04 and G09 land.
- `skills/review/SKILL.md` stays near 64 lines, but three long lines grow by about 180 words together.
- `skills/review-standards/SKILL.md` is preloaded into every reviewer dispatch, so G15 and G16 stay compact.
- No group needs code, a hook, a settings edit, a dependency or a new agent. An agent body change is verified with one sandbox dispatch; G22's description change needs one too.

## Conflicts to respect when implementing

- G02 lands before G03, and the line-11 parenthesis becomes "tick its Checkpoint box", or "tick its tasks" contradicts per-task ticking.
- G10 pastes only the fixed lines into a re-dispatch, never rebuttal reasons, to keep the reviewer diff-only.
- G11 limits "never rebutted" to a command that exited non-zero; a coverage finding stays rebuttable, as `/review --tests` allows.
- Every widening of the cold-reader brief folds into its Gaps question, so "two questions" stays true in `skills/feature/SKILL.md:66` and `agents/reviewer.md:28`.
- G21 states that a headless stop at a cap wins over "a headless run with --fix applies".
- G22's four line edits and its ADR land in one commit, or the four lines contradict each other.
- G24 reopens an operator-approved decision; G26 waits for a sandbox test of vendor precedence.

## Corrections from verification

The lenses changed these facts; the tables use the corrected versions.

- **goldfish history, not practice.** The per-phase gate checklists, "Done when" boxes, the five-agent Ralph loop, the documentation audit, the contract tests and the failure-pattern graduation criteria live only in `docs/archive/`. The shipped pattern store has no inclusion criteria, and nothing reads approved patterns.
- **goldfish memory.** `log_thought` writes a database row at once; the line reaches `STATE.md` only at the next regeneration, and `## Active Goal` has no caller that sets it. A reason is optional on `run()`. `vs_previous` subtracts raw values without checking that the metric matches.
- **goldfish suppressions.** No linter configuration enforces "never suppress"; the rule is text, and its forty suppressions carry codes but no reasons.
- **beads enforces less than its docs suggest.** Quality-gate labels are a documentation example that `bd close` never checks. The `bd swarm validate` direction check is a title-keyword heuristic. Compaction has one tier and targets database size. Version awareness is a one-shot `.beads/.local_version` check. `bd defer` is a snooze, not a terminal state; only closed is terminal. The upgrade gate runs six scenarios against one previous version in CI, not the documented matrix.
- **beads truncation notice.** It prefixes `bd prime`'s three built-in outputs only; extending it to diffs and logs is the harness's own generalization.
- **beads contradicts itself on markdown.** The plugin skill's two-week test allows a TodoWrite list for short work; every `bd prime` output forbids TodoWrite, TaskCreate and markdown files for task tracking.
- **elephant-goldfish scripts.** The lint baseline never checks out main, `git stash create` ignores `--include-untracked` so `git clean -fd` can delete untracked files, and it writes to `/tmp`. Its codegen step runs after the tests in every adapter, and no stale-output check exists. The headless outcome in G21 has no source precedent.
- **The harness has more than the proposals credited.** The cold-reader check already runs three readers per round in practice. The 2026-09-13 improvement review already refused capping the bug index. ROADMAP "Now" is not a phase marker. No explicit one-active-plan rule exists. Appending to a context line adds words, not rendered lines.

## Not verified

- Nothing here is implemented or tested; every "what the change would do" is a prediction from the text of the change.
- Claude Code's skill precedence behind G26 comes from a documentation summary, not a live test.
- The 60-row styrman-crew bug index behind G24 comes from a verifier's count and was not re-checked.
- Dave Rensin's article behind elephant-goldfish was not read; the repo's paraphrase was taken as stated.
- Not read in full: beads `docs/reference/`, `docs/architecture/` and most of the generated CLI reference; goldfish `cloud/`, `infra/`, `io/`, `metrics/`, `datasets/` and its web dashboard.
