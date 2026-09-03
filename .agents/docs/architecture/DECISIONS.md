# Architecture Decisions

Append-only. To reverse a decision, append a new ADR and mark the old one `Superseded by ADR-NNN` — never rewrite it. The ADRs of the previous design (a hook-driven Go harness) were removed with it in commits 09fc5cc and 20a5abd; they remain readable in git history.

## ADR-001: The harness is context, skills, and subagents — no hooks, no gates     (2026-09-03) — Accepted

**Context:** The previous harness enforced its Definition of Done with lifecycle hooks, a commit gate, scope digests, and a dispatch ledger. It was reverted the same day it was finished: the machinery cost more attention than the rules it enforced were worth for a solo operator.

**Decision:** Behaviour lives in prose: the global context (`context/AGENTS.md`), the skills, and the subagent definitions. Nothing runs on tool calls or session events. The only executable is the installer.

**Consequences:** Rules are followed because they are read, not because they are enforced. A rule that is not followed is fixed by rewriting the rule (`/retro`), not by adding a gate. Both vendors are served by identical mechanisms, since neither needs hook support.

## ADR-002: One Go binary, `styrman`, with a git resync loop     (2026-09-03) — Accepted

**Context:** The operator works on several machines and wanted to push a harness change from one and pick it up on another without thinking. A bash script was considered first; it was replaced by Go so that `styrman update` (fetch, pull, rebuild, reinstall) and a manifest could be tested against a scratch home rather than by hand.

**Decision:** `cmd/styrman`, stdlib only. `install` writes an atomic copy of every managed file and records the manifest in `~/.config/styrman/state.json`; `update` fast-forwards the repo and re-runs `go run ./cmd/styrman install`, so a code change in the repo takes effect in the same update; the binary copies itself to `~/.local/bin/styrman`.

**Consequences:** Go is required on every machine (already true for the operator). Two copies of each skill exist on disk, one per vendor — a deliberate trade against symlinks, which neither CLI documents following across directories.

## ADR-003: Single source; Antigravity agents are derived at install     (2026-09-03) — Accepted

**Context:** The previous layout committed generated copies (`claude/CLAUDE.md`, `antigravity/GEMINI.md`, derived agents) that needed a drift check to stay honest.

**Decision:** `context/`, `agents/` (Claude Code format), and `skills/` are the only sources. `styrman install` assembles the context (AGENTS.md + dev-setup.md + generated inventory) and translates agent frontmatter for Antigravity in memory. Nothing generated is committed.

**Consequences:** No drift is possible between committed copies. The translation table (`toolMap`, `modelMap` in `derive.go`) is the one place vendor differences live; an unmapped tool fails the install instead of shipping a name Antigravity would hang on.

## ADR-004: Uninstall is manifest-only; hand-installed skills survive     (2026-09-03) — Accepted

**Context:** The previous installer mirrored whole directories and pruned anything the repo did not have, which wiped skills installed by hand with `npx skills add -g`. The operator wants adding skills by hand to stay easy.

**Decision:** `install` records every file it writes; `uninstall` and the resync step remove only paths from that manifest. A `~/.claude/CLAUDE.md` without the styrman marker is never overwritten. `styrman adopt` copies a hand-installed skill into the repo when it should become permanent.

**Consequences:** Directories under `~/.claude/skills` may hold both styrman's and the operator's skills; `status` reports only styrman's. A skill installed by hand into a vendor directory is not versioned and is not deployed to the other vendor until adopted.

## ADR-005: The context is named AGENTS.md everywhere     (2026-09-03) — Accepted

**Context:** The operator asked for the agents.md standard. Claude Code reads `CLAUDE.md`; Antigravity reads `AGENTS.md` or `GEMINI.md`.

**Decision:** The deployed file is `~/.claude/AGENTS.md` and `~/.gemini/config/rules/AGENTS.md`. `~/.claude/CLAUDE.md` and a project's `CLAUDE.md` are the single line `@AGENTS.md`. `/init-project` creates `.claude/{skills,agents,commands}` as symlinks into `.agents/`, matching what the skills.sh CLI does on its own.

**Consequences:** One canonical name at every level. Claude Code's import adds one hop; it is documented to load user-level imports without prompting.

## ADR-006: The binary is `sh-agent`; GCP defaults are asked for, never stored     (2026-09-03) — Accepted

**Context:** The operator asked for a less repo-specific binary name and, on reviewing `context/dev-setup.md`, rejected stored GCP defaults: a placeholder region or registry name is exactly the kind of value an agent will silently assume. They also asked for zero scripts in the repo.

**Decision:** The binary, its state directory (`~/.config/sh-agent`), its env var (`SH_AGENT_REPO`), and its file marker are named `sh-agent`; the repo keeps its name. The context says work is local by default and, when the user says GCP, the agent asks for project, region, Artifact Registry repository, service name, and visibility before doing anything, then records them in the project's `AGENTS.md`. The one shell script in the repo (`mermaid-syntax/scripts/validate-mermaid.sh`) was deleted along with the old hook scripts under `~/.agents/bin`.

**Consequences:** ADR-002's name is superseded by this one; its content stands. The first GCP action in a new project costs one question instead of one wrong assumption. Skills that need helper code ship it as documentation, not as scripts.

## ADR-007: Hooks stay out — researched and rejected     (2026-09-03) — Accepted

**Context:** After the rebuild the operator considered three hooks (a deploy gate for documentation freshness, a wrong-model warning, a loop protector). The options for both vendors were researched in `.agents/docs/plans/2026-09-03-research-hooks.md`: all three are buildable as `sh-agent hook` subcommands, but the Claude side requires a managed block in `~/.claude/settings.json`, and every hook adds a way for the harness to be wrong at run time.

**Decision:** No hooks. The operator's words: "Skip the hooks. I think this was my problem before. Keep the no hooks/scripts rule." R3 (no hooks, no gates, no CI) and the zero-scripts rule stand as written; the research is kept as the record of what was weighed.

**Consequences:** The Definition of Done, the model choice, and loop avoidance remain prose in the global context, enforced by reading. If a rule is not followed, the fix is `/retro` and a better rule, never a gate. Reopening this needs a new ADR that answers the research's cost section.

## ADR-008: Bugs are the project's memory; scratch lives in the project; the trail carries progress     (2026-09-03) — Accepted

**Context:** On reviewing the global context the operator asked for: a bug subagent that numbers every failure so the same mistake is not made twice; a project-local scratch space instead of `/tmp`; explicit checkpointing so resume has something to resume from; clear steps for plans, phases, long-running work, testing, verification, and e2e. Two sources were read in full for strategy — Böckeler's harness engineering article and the elephant/goldfish repository (see `.agents/docs/plans/2026-09-03-research-harness-strategies.md`).

**Decision:** Every failure gets a `BUG-NNN` record under `.agents/docs/bugs/` with a prevention rule, indexed in `INDEX.md`, written by the new `bugs` agent and read at resume; a rule that recurs graduates into a project rule or a test. All generated and temporary files go to `.agents/scratch/` (gitignored). Work is classified as trivial, task, feature, or bug; a feature has a plan file with phases and task checkboxes, a cold-reader check of the plan by a fresh `reviewer` before any code, one phase at a time, and each phase ends in a checkpoint (tick, changes entry, commit, push). A bug gets a goldfish diagnosis on symptom and repro only and a failing test before the fix. Testing has four named tiers (unit, integration, e2e, cloud smoke) and e2e is required before "done". The Definition of Done is an ordered sequence. The harness's own state lives in `~/.config/sh-agent/`. Adopted from the sources: guides and sensors as the vocabulary, the steering loop as the retro cadence, design-before-code with an artifact-only goldfish, fix-or-rebut with a round cap, the four user gates. Rejected: docs-in-chat, per-project rewritten command templates, and hooks (ADR-007).

**Consequences:** More is written to disk per task, and it is what makes the next session and the next agent cheaper. Bug records make repeated mistakes visible and give `/retro` concrete input. The project layout gains `docs/bugs/` and `scratch/`; `/init-project` creates both. Seven agents instead of six.

## ADR-009: The context stays short; procedures are skills; the reviewer is tiered     (2026-09-03) — Accepted

**Context:** The operator asked for an honest opinion of the rebuilt harness. The deployed context had grown to 200 lines (about 5k tokens every session). Claude Code's own guidance is that bloated instruction files make the model ignore the actual rules and that workflows relevant only sometimes belong in skills; Anthropic's context-engineering guidance treats context as a finite resource with diminishing returns. Every Opus review on a routine diff was also the largest recurring cost.

**Decision:** `context/AGENTS.md` is cut to under 70 lines with every rule kept. The feature and bug procedures become the `feature` and `bug` skills, model-invocable so they load themselves when the work matches, and callable as `/feature` and `/bug`. The Definition of Done keeps its five-step order in the context; the detail lives in the agents. The embedded inventory lists names only, since both CLIs inject descriptions themselves. A compaction instruction names what must survive. The `reviewer` runs on Sonnet and reports only correctness, security, requirement, and rule findings; `/preflight` dispatches it on the strongest model. Changelog entries name the agents dispatched so cost can be counted at the first retro. The remaining improvements (pruning skills, the Mac live check) are the operator's, and no further tuning happens before real use.

**Consequences:** About 3.5k fewer tokens per session before the first prompt. Procedures now compete for nothing until they are needed. The risk moves from "rule lost in a long file" to "skill not triggered"; the description of each procedural skill is therefore written as a trigger list, and the context keeps a one-line pointer per kind of work.

## ADR-010: Harness state stays outside the project and is approved, not pre-authorised; the skill matrix is a documented procedure, not a script     (2026-09-03) — Accepted

**Context:** The sandbox matrix showed that Claude Code denies headless access to `~/.config/sh-agent/` (feedback log, retro proposals, global handoffs) and prompts for it interactively. Two ways out were on the table: have `sh-agent install` add the directory to `permissions.additionalDirectories` in `~/.claude/settings.json`, or move the state into projects. The matrix itself is a bash driver, prompts, and checks under gitignored scratch, which exists on one machine only.

**Decision:** Neither escape hatch. `REQUIREMENTS.md` names editing vendor settings a non-goal and R7 fixes the state's location; the one prompt per session is accepted and documented (README, `deploy/local.md`), with the user-side setting named for those who want it gone. The matrix is recorded as documentation — `testing/skill-matrix.md` lists every scenario's prompt and check and the invocation flags — so any session on any machine can rebuild the kit in scratch; the scripts are not committed, keeping "the harness is context, skills, and subagents" true.

**Consequences:** One approval per Claude Code session that writes feedback or a proposal; nothing installed outside the customization directories. Re-running the matrix costs a few minutes of setup from the doc instead of a `git pull`; the trade is deliberate, and a Go `sh-agent sandbox` subcommand remains possible if the matrix becomes routine.
