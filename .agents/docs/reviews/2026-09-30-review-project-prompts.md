---
owner: reviewer
date: 2026-09-30
---
# Review — project-prompts   (2026-09-30)

Range reviewed: the whole repository at HEAD (3478a1c)

## Status: COMPLETE

Lens: `review-standards/references/prompt-audit/prompt-audit.md` (Anthropic prompt-audit reference, Apache-2.0, vendored verbatim per `skills/VENDORED.md`), executed as a `--project`-scope audit over this harness checkout's own prompt surface.

## Step 0 — Scope and target model

**Scope**, from the dispatch: `AGENTS.md` (root), `context/AGENTS.md`, `context/dev-setup.md`, `agents/*.md` (8 files), `skills/*/SKILL.md` less the six vendored skills (`grilling`, `codebase-design`, `find-skills`, `anti-ui-slop`, `gke-basics`, `terraform`, per `skills/VENDORED.md`), `commands/*.md`, `rules/*.md`. No `CLAUDE.md` or `GEMINI.md` exists at the root (by design — root `AGENTS.md` says "no CLAUDE.md, no GEMINI.md"; R2). `commands/` and `rules/` hold only `.gitkeep` — no prompt content there to audit. `.claude/settings.local.json` exists but is a settings file (Step 1 explicitly excludes settings/credential files from reading) — not opened. `.agents/rules/` does not exist in this repo.

**Target models**: this harness explicitly deploys one context to two vendors — Claude Code and the Antigravity CLI (Gemini) — so every instruction file in scope is read by both. No finding in this report rests on one model's documented behaviour alone; per the dispatch, any such finding would be downgraded to `flag`, but none arose.

## Step 1 — Inventory

- `AGENTS.md` (root, 90 lines) — project rules, source of the global context.
- `context/AGENTS.md` (67 lines) — the global context, deployed verbatim to both vendors.
- `context/dev-setup.md` (13 lines) — appended to the context at install.
- `agents/bugs.md`, `agents/changelog.md`, `agents/documentation.md`, `agents/dreamer.md`, `agents/explorer.md`, `agents/reviewer.md`, `agents/security.md`, `agents/verifier.md` — 8 subagent definitions (76–108 lines each), Claude Code format.
- 22 non-vendored `skills/<name>/SKILL.md`: `bash`, `bug`, `cloud-run-deploy`, `distill`, `docker`, `dream`, `feature`, `feedback-log`, `frontend`, `gcloud-setup`, `golang`, `handoff`, `init-project`, `json-api-style`, `mermaid-syntax`, `preflight`, `python`, `research`, `retro`, `review`, `review-standards`, `typescript` (28–281 lines each).
- `commands/`, `rules/` — empty (`.gitkeep` only).
- No request-building code, no tool-schema definitions, no few-shot blocks, no application system-prompt assembly code anywhere in scope: this is a pure agent-configuration surface (Group 4 mostly not applicable, per the dispatch).

## Step 2 — Provenance

Checked `.agents/docs/bugs/INDEX.md` (through BUG-089) and `.agents/docs/architecture/DECISIONS.md` (ADR ids cross-referenced below) for every prohibition and emphatic line that looked like a candidate. Every prohibition/emphasis instance found in scope either (a) carries its own stated reason in the same sentence, or (b) traces to a specific bug row or ADR already cited inline in the text (e.g. BUG-004 frontmatter quoting, BUG-006/015 install-vs-copy, BUG-041 crypto/rand vs. hand-rolled, BUG-085 the GANNET canary and the 24,000-byte cap, ADR-029/030/033/035/038/045/045a — all verified to exist in `DECISIONS.md`). None of the historical dates or bug/ADR citations read as decorative "history narrative" (Group 2 row) — each is load-bearing: it is either the reason a still-active constraint exists, or (in `REQUIREMENTS.md`'s case) the required amendment trail for a file that is immutable without the operator's approval.

## Step 3/4 — Classification and anti-pattern scan

Ran the Group 1 signal greps (§1a–1f of the lens: `MUST|NEVER|ALWAYS|CRITICAL|IMPORTANT` density, `think step by step`/`<scratchpad>`/prefill scaffolds, hedge words, word/turn-count ceilings, grader vocabulary, date-conditionals, anti-formatting rules, trait claims, identity stubs, `Remember,|Again,|As stated above`, retired-model names, reminder cadences) across `context/`, `agents/`, and all 22 non-vendored `SKILL.md` files, then read the full text of `context/AGENTS.md`, `context/dev-setup.md`, root `AGENTS.md`, all 8 `agents/*.md`, and `skills/{bug,feature,preflight}/SKILL.md` in full (the three most likely to carry patch-accretion cruft, per their bug history of repeated gate revisions — BUG-018, BUG-026, BUG-042).

Result: no Group 1 hits. The one non-trivial grep hit (`MUST`/`NEVER` in `mermaid-syntax/SKILL.md:112,175`) is exact-syntax contract text for a fragile operation (Mermaid frontmatter and namespace rules) — keep-list items 3 and 4, not cruft. The `v11.5.0` version note in the same file documents a specific historical Mermaid bug range and its fix; it is tool-contract knowledge, not a stale prompt scaffold, and nothing in the repository contradicts it. No `<scratchpad>`, no assistant-prefill pattern, no forced-tool-choice pattern, no numeric word/turn caps, no grader vocabulary, no anti-formatting rules, no identity-stub role lines, no retired-model names, no reminder-cadence text anywhere in scope.

Every prohibition in `context/AGENTS.md` § Hard rules and root `AGENTS.md` carries its own reason in the same sentence (e.g. "Never hand-roll auth, crypto… — vetted libraries only. Calling the standard library's `crypto/*` as designed… is not hand-rolling") — these read as reasoned, current constraints (keep-list item 5), not inherited pressure language from older models. The project's own stated design ("No hooks, gates, or scripts stand behind these rules; they hold because they are read") explains the density of `NEVER`/`ALWAYS` without it being dated-model over-triggering language: each instance is a distinct, reasoned constraint rather than a repeated intensifier.

**Group 2 — stale facts / cross-file contradiction.** Checked every file path and test name that `AGENTS.md` (root) and `context/AGENTS.md` name against the actual tree: `internal/harness/plan.go`, `internal/harness/derive.go`, `REQUIREMENTS.md`, `.agents/docs/plans/2026-09-03-rebuild.md`, `.agents/docs/deploy/remote.md`, `.agents/docs/specs/sh-agent/delivery.md`, `internal/harness/report.go`, `internal/search/search.go`, `.agents/docs/testing/guide.md`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/testing/2026-09-03-verify-sandbox.md`, `.agents/docs/testing/2026-09-03-verify-real-home.md`, `.agents/docs/retros/2026-09-18-retro-2.md`, `.agents/docs/specs/sh-agent/search.md`, `.agents/docs/plans/2026-09-29-dreamer-memory-loop.md`, `REQUIREMENTS.md § Known limits` — all present. Named test functions (`TestHarnessMap`, `TestRenderedContextFitsAntigravityBlock`, `TestSendReport`, `TestDryRunNeverBuilds`, `TestReadsNeverLeaveCheckout`) all resolve to real files under `internal/`/`cmd/`. All 8 ADR ids cited in the two AGENTS.md files exist in `DECISIONS.md`. No cross-file contradiction found between `context/AGENTS.md`, root `AGENTS.md`, and the 8 agent definitions: the repeated "docs root resolution" and "Records/Scratch/Never" boilerplate across all 8 agent files is identical working redundancy (keep-list item 8), and the ownership split between `changelog` (changes/, DECISIONS.md) and `documentation` (README, specs, testing/deploy guides, ROADMAP Next/Later) is stated consistently in both files with no overlap claimed twice.

**Group 3 — tool descriptions.** Not applicable in the conventional sense: this project defines no custom API tool schemas. The closest analogue, the `tools:` frontmatter lists on the 8 agents, name only Claude Code's built-in tools (`Read`, `Grep`, `Glob`, `Edit`, `Write`, `Bash`) — no worked examples, no scolding cross-references, no fake dialogue embedded in any agent's tool list or body.

**Group 4 — request config / sub-agent roster.** No request-building code in scope (not applicable). Roster check (the one Group-4 item that always applies): compared all 8 agents pairwise for near-duplicate task+toolset pairs. `bugs`, `changelog`, and `documentation` share the identical tool list (`Read, Grep, Glob, Edit, Write`) but write to disjoint artifact classes (bug records, changelog/ADRs, and how-it-works docs respectively) with no overlapping ownership claimed — not a duplicate roster entry. `explorer` and `dreamer` share `Read, Grep, Glob, Bash, Write` but operate over disjoint inputs (codebase recon vs. memory-store/trail consolidation) and produce disjoint artifacts (recon brief vs. retro candidate list) — not a duplicate. `reviewer` and `security` share `Read, Grep, Glob, Bash, Write` but have explicitly distinct checklists (six general-correctness criteria vs. six security-specific checks), different trigger conditions (size threshold vs. "touches auth/crypto/secrets/IAM/parsing/dependencies, any size"), and different models/effort (`sonnet`/high vs. `opus`/max) — a deliberate two-lens split named in the Definition of Done, not redundancy. No roster finding.

## Findings

None at high or medium confidence. The surface is clean: every prohibition and emphasis instance found carries its own reason or traces to a bug/ADR record that still applies, no dated scaffolds or scoring/grader vocabulary appear anywhere in scope, every path and test name the instruction files cite resolves in the current tree, and the sub-agent roster has no redundant pair.

Two low-confidence, non-actionable observations (left out of the proposed diff per the lens's own threshold — flagged, not fixed):

- **[flag / confidence: low]** `skills/mermaid-syntax/SKILL.md:159` — "Fix: Upgrade to Mermaid v11.5.0+." This is idiom-dated only in the trivial sense that any pinned library version eventually ages; nothing in the repository contradicts it (it documents a specific historical bug range, not the current recommended version), so it is not a Group-2 stale-fact finding under the lens's own bar (a claim the repository contradicts). No action proposed.
- **[flag / confidence: low]** `feedback-log/SKILL.md:56` and `distill/SKILL.md:156` both reference a 2026-09-13 format cutover for `feedback.jsonl`'s `rule_ref` field ("Entries written before 2026-09-13 use free text and are simply never matched"). This is migration-relative phrasing (Group 1d signal) but is functionally necessary here — it tells the agent how to treat old log rows it may still read — so it is not proposed for removal; noted only because the lens's signal list would flag the phrasing pattern.

## Proposed diff

None. No finding reached the high/medium-confidence, actionable bar (Step 6 of the lens). An audit that finds nothing changes nothing.

## Criteria

1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

## Not examined

- Full line-by-line reads of 19 of the 22 non-vendored `SKILL.md` files (`bash`, `cloud-run-deploy`, `distill`, `docker`, `dream`, `feedback-log`, `frontend`, `gcloud-setup`, `golang`, `handoff`, `init-project`, `json-api-style`, `mermaid-syntax`, `python`, `research`, `retro`, `review`, `review-standards`, `typescript`): these were covered by the Group 1 signal-grep sweep (§1a–1f patterns) across their full text plus targeted spot-checks of the hits it returned, not by a full sentence-by-sentence read against Step 3's "could the model already know this" test. `bug`, `feature`, and `preflight` were read in full.
- `review-standards/SKILL.md` was not re-read from this repo's copy on disk; its content is known verbatim because it is the file that invoked this audit (this repo is the source `sh-agent` deploys from), so the deployed and checked-out copies are the same file by construction.
- External verification of the Mermaid-version claim (`v11.5.0`) against the currently shipping Mermaid release — out of the repository's own record, so left as a low-confidence flag rather than checked.
- `.claude/settings.local.json` — a settings/credential-adjacent file the lens's Step 1 says not to read; not opened, consistent with the lens.
- No PR or diff range exists for this dispatch (whole-repo `--project` scope was named explicitly), so no "what changed since the last audit" comparison was possible; this is a from-scratch pass, not a delta from a prior prompt-audit report (none exists in `.agents/docs/reviews/` to diff against).
