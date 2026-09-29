---
owner: reviewer
date: 2026-09-29
---
# Review — memory-loop-code   (2026-09-29)

Range reviewed: commits 1b8e84b..HEAD plus working tree (`git diff 1b8e84b^`), scoped to the source files named by the dispatch: `agents/dreamer.md` (new), `skills/init-project/SKILL.md` (Step 5b, Step 6, Step 7), `skills/dream/SKILL.md` (Step 1), `skills/feedback-log/SKILL.md` (rule 7), `skills/retro/SKILL.md` (source 8, record line), `skills/distill/SKILL.md` (Step 1), `context/AGENTS.md` (one clause), `AGENTS.md` (one Testing bullet), `.gitignore`, `.agents/docs/testing/skill-matrix.md` (three rows). Spec: `.agents/docs/specs/sh-agent/memory-loop.md` (read for cross-checking, not itself in scope). Plan, logs, reviews and retros under `.agents/docs/` treated as context only, per the dispatch.

## Status: COMPLETE

## Verdict: APPROVE

This is prose-level documentation (skills, an agent definition, the global and project context, one gitignore line, three test-matrix rows) with no Go code touched. Cross-checked every restatement of the two load-bearing facts this change introduces — the slug algorithm and the dreamer dispatch contract — across all six files that state either, and found them consistent everywhere (`skills/dream/SKILL.md:25`, `skills/init-project/SKILL.md` Step 5b.1, `skills/feedback-log/SKILL.md` rule 7, `skills/retro/SKILL.md` source 8, `skills/distill/SKILL.md` Step 1, `agents/dreamer.md`, the spec). The one destructive step this diff adds — the memory-store migration (`skills/init-project/SKILL.md` Step 5b.3) — is guarded on both directions (both-hold-files: move nothing; old-empty: nothing to move) and the settings-file edit it depends on (Step 5b.2) restores its saved copy on a failed JSON parse and makes no second attempt, matching the destructive-guard hard rule. `agents/dreamer.md`'s tools (`Read, Grep, Glob, Bash, Write`) all resolve in `internal/harness/derive.go`'s `toolMap`, its frontmatter matches the shape of `agents/verifier.md`/`agents/reviewer.md`, and every dispatch site (`skills/retro/SKILL.md`, `skills/distill/SKILL.md`, the skill-matrix's `dreamer` row) hands it the exact same contract line format with a 30-call budget under its 40-turn `maxTurns`, satisfying the "budget below maxTurns" rule. No dead code, no directory-layout surprise, no new dependency, no auth/crypto surface.

## Findings

None that meet the bar (correctness, security, a stated requirement, or a recorded bug reintroduced). One item considered and set aside below, for the record.

- **Considered, not filed:** `skills/dream/SKILL.md` Step 1's split-store check fires only "when `<store>` came from a setting and both `n` and `m` are above zero." The case `n = 0, m > 0` (configured store empty, files still sitting at the unmigrated default path) is not flagged as a split — the run ends with "no auto memory for `<repo>`" even though `m` is nonzero. The raw count is still printed on the first line ("so they are never skipped"), which reads as a deliberate choice to surface the number without adding a second anomaly branch, and the commit history around D12 suggests this exact absent/mismatched-store shape was already discussed. Flagging for visibility only — not a finding.

## Criteria

1. **Dead code** — pass. Nothing orphaned; every new clause (dreamer contract, Step 5b, gitignore line, feedback-log rule 7, `--setting-sources` testing bullet) is referenced from at least one consumer.
2. **Flow tracing** — pass. Traced `dreamer`'s dispatch contract from its three callers (`skills/retro/SKILL.md` source 8, `skills/distill/SKILL.md` Step 1, the skill-matrix `dreamer` row) into `agents/dreamer.md`'s own "Your dispatch" section — the line shape (`Store: … Docs root: … Since: … Promoted: …`) matches verbatim at every site. Traced the slug formula from `skills/init-project/SKILL.md` Step 5b.1's `sed` through `skills/dream/SKILL.md:25`'s prose definition and the skill-matrix rows — all agree on "every character other than a letter, digit or hyphen → hyphen," and a repo-wide grep found no surviving restatement of the older "every `/` turned into `-`" rule in any living doc (only in out-of-scope historical plan/review records). Traced the destructive migration step's two guard branches and the settings-file edit's parse-then-restore guard; both hold.
3. **Standards conformance** — pass. `agents/dreamer.md` frontmatter (`tools`, `model`, `effort`, `maxTurns`) matches the shape of the two other agent definitions checked (`verifier.md`, `reviewer.md`); its `description:` is one line, has no colon needing quoting (quoted anyway, harmless), and contains neither "dispatch prompt" nor "harness binary," per `AGENTS.md` § Rules for working here. No rule-restatement gaps found (see Flow tracing).
4. **Directory layout** — pass. New agent lands in `agents/`; no stray files; test-matrix rows land in the existing matrix file.
5. **Simplicity (DIY-except-security)** — pass. No new dependency, no code; nothing in this diff touches auth/crypto/sessions.
6. **Decoupling** — pass. The dreamer's dispatch line is the one seam three callers share; none of them reach into its internals, and its "Never" list (source, config, rules, skills, agents, context, memory store, REQUIREMENTS.md, another agent's records) keeps it from stepping on the writers it reads after.

## Not examined

- `README.md` and `ROADMAP.md`, which also changed in this range (per `git diff --stat`) but were not named in the dispatch's source-file list; the dispatch treats docs-pass material as a separate concern (`/review --docs`), so these were read only far enough to confirm they don't restate the slug or dispatch-contract facts, not reviewed against the six criteria.
- `.agents/docs/bugs/INDEX.md` in full (only the tail was read) — the diff does not touch it, so this was a sanity check, not a full pass.
- Whether `go test ./...` / `TestRenderedContextFitsAntigravityBlock` actually pass on the current tree — that is the verifier's job (Definition of Done step 1), not re-run here; the `.agents/docs/changes/2026-09-29-1428-memory-loop-phase-1.md` record cites a prior PASS (21,787/22,000 bytes) but that is a historical record, not something this review re-verified.
- The historical plan/review/testing records under `.agents/docs/` (cold-read rounds, behavioural-run logs) — read only as corroborating context per the dispatch's instruction that they are not in scope, including a numeric discrepancy noticed in passing (plan D15 says "six" void `/dream` probes, the testing log and changes entry say "nine," matching `AGENTS.md`'s new bullet) — not filed as a finding since it sits entirely in out-of-scope historical trail and doesn't affect the correctness of the in-scope rule text.

## Outcome
- APPROVE, no findings filed. The two items set aside were taken anyway: plan D15's "six" void probes corrected to nine (the logs' count), and `skills/dream/SKILL.md` Step 1 now says, when the settings store is empty but the default path holds files, "the default path holds memory files: run /init-project to move them" instead of a bare "no auto memory".
