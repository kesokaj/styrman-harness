---
owner: explorer
date: 2026-09-29
---
# Recon — audit-search-consumers (2026-09-29)

## Status: COMPLETE

## Question
Audit every consumer in styrman-harness for whether the planned local SQLite index (per-project FTS5 + trail metadata graph; global feedback/retro index) would help, and how much: context/AGENTS.md procedures, non-vendored skills, agents/*.md, commands/, rules/, sh-agent verbs, testing/skill-matrix.md.

## Summary
The highest-value consumers are all "restatement/paraphrase grep" and "structured trail query" steps: `review-standards` skill's restatement check, `dream`'s graduation grep, `dreamer`'s rule-evidence grep, the global hard rule on reversed claims, and the resume/handoff awk-parsing of unticked plan boxes and INDEX.md rows — these replace fragile exact-grep or hand-rolled parsing with ranked/paraphrase-tolerant search or a direct structured answer. The clearest **non-consumers** are the verifier's own diff-scoped greps (suppressions, coverage), preflight's secret scan, and any `/review --docs` grep run "at the scope's head" of a historical range — all freshness-critical or need an exact literal match on a specific git revision the index (built from tracked files, presumably at HEAD) would not represent. `commands/` and `rules/` are both empty — nothing to audit there. Instruction placement: `context/AGENTS.md` has ~48 bytes of headroom under its pinned 22,000-byte budget, so a durable pointer to the index belongs in individual skill/agent files, not the global context; I directly observed (as this dispatch's own subagent) that Claude Code injects the full `~/.claude/CLAUDE.md`→`AGENTS.md` chain into a subagent's context, so a line already in `context/AGENTS.md` would reach Claude Code subagents for free — but with only ~48 bytes free that channel is effectively unusable, and Antigravity's equivalent inheritance for subagents is not stated anywhere in this repo (open question).

## Relevant files
- `context/AGENTS.md` — global procedures; Resume line 9 (unticked plan boxes, INDEX.md), Hard rules line 55 (restatement grep), line 51 (git-diff-based checks)
- `skills/review-standards/SKILL.md:48` — restatement-of-a-rule grep across context/, skills/, agents/, README.md
- `skills/review/SKILL.md:51` — `--docs` lens: `git grep` old wording "as the scope's head" across living docs
- `skills/dream/SKILL.md:32,41,48` — reads MEMORY.md + trail + context/skills/agents; graduation check is "one grep over key words, each hit read in full"
- `skills/retro/SKILL.md:36,43,68-69` — reads feedback.jsonl/retros.jsonl, exact-match controlled-vocabulary greps
- `skills/bug/SKILL.md:33` — reads `<docs-root>/bugs/INDEX.md` first
- `skills/feature/SKILL.md:19` — checks REQUIREMENTS.md/ROADMAP.md/bugs/INDEX.md
- `skills/handoff/SKILL.md:154` — Resume Protocol incl. an `awk` one-liner over every `plans/*.md` to find unticked boxes outside fenced blocks
- `skills/init-project/SKILL.md:24` — counts date-only handoff/change filenames via `ls`+`grep` (filename pattern, not content)
- `skills/preflight/SKILL.md:40` — secret scan on staged diff (must stay exact/live)
- `agents/bugs.md` — "if three records share a prevention rule, say so" (cross-record similarity, currently manual)
- `agents/dreamer.md:29,36` — reads every INDEX.md row + opens records; `git grep` key words over context/AGENTS.md/rules/skills before every rule row (restatement check, same shape as review-standards)
- `agents/reviewer.md`, `agents/verifier.md`, `agents/explorer.md`, `agents/changelog.md`, `agents/documentation.md`, `agents/security.md` — read in full; verifier is the clearest non-consumer (diff-scoped, freshness-critical)
- `internal/harness/report.go` — collects/pushes `feedback.jsonl`, `retros.jsonl`, `proposals/`; a write path, not a search consumer itself, but is the source the global index would ingest
- `internal/harness/derive.go` — no mention of context injection for Antigravity subagents; comment "markdown agents inherit the global skills anyway" is the only hint
- `REQUIREMENTS.md:49-51` — "Known limits": deployed context pinned under 22,000 bytes; `rules/` files share Antigravity's 24,000-byte block with it
- `.agents/docs/testing/skill-matrix.md` — behavioural-test procedure; no search-index consumer, but any new skill/agent behavior needs a matrix row + sandbox run before it can be claimed to work (project convention, not this audit's finding but relevant to landing the feature)
- `commands/`, `rules/` — both empty (`ls` returned nothing) — no consumers here

## Flows
- Restatement-check pattern (repeated 3×, independently, in `review-standards`, `dream`, `dreamer`): grep key words/subject nouns across a fixed small set of living files (`context/`, `skills/*/SKILL.md`, `agents/`, `README.md`, or the project's own `AGENTS.md`/rules) → read each hit in full → judge. This is exactly the "restatements of a claim" query the dispatch names.
- Resume/handoff pattern: read several named files in a fixed order (`AGENTS.md`→`TODO.md`→newest handoff→`ROADMAP.md`→active plan's unticked boxes→3 newest changes→`bugs/INDEX.md`) then run a hand-written `awk` script over `plans/*.md` to find unticked boxes ROADMAP does not already name. This `awk` is the kind of thing "unticked plan boxes" as a structured query directly replaces.
- Bug/ADR trail pattern (`dreamer`, `bug`, `feature`): read `INDEX.md` (all rows, newest last) then open matching records individually. A "every record that names file X" or "open bugs" graph query collapses the read-INDEX-then-open-records loop into one structured answer, though the records still need opening for their prose (root cause, prevention text).
- Verifier/preflight pattern: grep the **diff** (not the repo) for suppressions, secrets, added `.skip(`/`.only(`. This is scoped to uncommitted/staged changes on the live tree — an index over tracked files at some snapshot cannot see this and must never be substituted.

## Conventions
- Every subagent writes its artifact early and updates it (shared template across `explorer`, `verifier`, `reviewer`, `bugs`, `dreamer`) — any index-query instruction added to an agent file should sit in that agent's own numbered procedure, not as a late aside, to survive a cut-off run.
- Skills state "Docs root: `.agents/docs/` if it exists, else `docs/`" repeatedly — any index path/query wrapper needs to resolve the same way (index lives under `.agents/scratch/search.db`, independent of docs-root, but its queries target `<docs-root>`).
- The project explicitly documents a size budget as a *test-pinned* constraint (`REQUIREMENTS.md` "Known limits", `internal/harness/size_test.go`'s `TestRenderedContextFitsAntigravityBlock`) — any edit to `context/AGENTS.md` must be checked against that test, and per the dispatch's own note there is only ~48 bytes of slack, i.e. effectively none.

## Risks
- The `--docs` review lens and `dreamer`/`dream`'s restatement greps explicitly search "as the scope's head" or "before every rule row" — if the index is built once (or on-demand only from the current working tree) and a review scope is a past range `<a>..<b>`, using the index instead of `git grep` at `<b>` would silently search the wrong revision. Any instruction to prefer the index for these must carve out non-HEAD ranges.
- `retro`'s two exact-match greps (`rule_ref`, `target`) are deliberately literal matches against a controlled vocabulary — a ranked/fuzzy FTS hit could return a near-miss (different `rule_ref`) that reads as a false "this recurred elsewhere" positive; these should stay `grep -E` unless the index also supports exact-match mode.
- `verifier`'s own coverage/suppression checks and `preflight`'s secret scan are described in `AGENTS.md`'s Hard rules as never-piped, freshness-critical, live-tree operations — using a cached index here directly risks the "silence from your own instrumentation is not evidence" failure mode this project already guards against elsewhere.
- Whether Antigravity subagents receive `context/AGENTS.md` content at all is undetermined from this repo's code or docs (see Open questions) — an instruction placed only in `context/AGENTS.md` may not reach an Antigravity-dispatched agent, which argues for placing the instruction directly in each `agents/*.md` / `skills/*/SKILL.md` regardless of the context-injection answer.

## Open questions
- Does Antigravity inject the deployed global rules block (`~/.gemini/config/rules/AGENTS.md`) into a dispatched subagent's context, the way Claude Code demonstrably does (I received the full CLAUDE.md→AGENTS.md chain as this dispatch's own subagent)? Not stated anywhere in `internal/harness/derive.go`, `agents/*.md`, or `REQUIREMENTS.md`. `testing/skill-matrix.md`'s `ctx-deploy` row tests that the *main* agent can quote the context; no row tests a dispatched subagent's visibility into it on either CLI.
- Whether the planned index would be rebuilt automatically before each query (freshness) or on a stale cache — not specified in the dispatch; this determines whether the "verifier/preflight must not use it" risk above is absolute or just "must rebuild first."
- Prior research (`.agents/docs/plans/2026-09-29-research-sqlite-index.md`) already concluded CONDITIONAL/build-nothing based on measured grep latency (git grep 11ms on this repo) — this audit does not re-litigate that verdict, only maps consumers assuming the index is built.

## Not examined
- `internal/harness/plan.go` beyond the context-deployment grep (full install-plan logic not read)
- `skills/cloud-run-deploy`, `skills/docker`, `skills/gcloud-setup`, `skills/golang`, `skills/json-api-style`, `skills/mermaid-syntax`, `skills/python`, `skills/typescript` SKILL.md bodies were not read in full (grepped only for grep/read/scan keywords, none matched — low prior they are consumers, but not confirmed by full read)
- `.agents/docs/plans/2026-09-29-recon-project-index.md` and `2026-09-29-research-memory-and-prompt-audit.md` (sibling docs, not opened — may contain overlapping findings)
- Frequency/volume of past subagent trail-file reads (dispatch said "optional, only if cheap") — not measured
