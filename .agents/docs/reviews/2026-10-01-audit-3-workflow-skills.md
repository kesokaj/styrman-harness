---
owner: reviewer
date: 2026-10-01
---
# Review — audit-3-workflow-skills   (2026-10-01)

Range reviewed: scope at HEAD 2ba4a50, no diff: `skills/{feature,bug,review,review-standards,preflight,handoff,init-project,research,spec,requirements,testing-strategy,explain-project,retro,distill,dream,feedback-log,project-memory}/**` and `skills/VENDORED.md`'s notes on `grilling`, `codebase-design`, `find-skills`; checklist = recon sections WS and G of `.agents/docs/plans/2026-09-30-recon-audit-fix.md`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

No WS item is still broken as the recon described it. The worst open points are two medium-low procedural gaps: `distill` Step 1 dispatches the `dreamer` with a budget but without the "write the artifact early" clause, and the G4 / G10 coverage gaps have no plan row, decision or strike. Everything else is low. The criteria below are `n/a — audit`.

## Findings

### Status of every recon ID (verified in the files)

| ID | was | now | evidence |
|---|---|---|---|
| WS1 | med: review paragraphs 1,819 / 1,730 / 3,446 chars | fixed | `skills/review/SKILL.md`: Docs brief (l.60-68) and Prompts brief are numbered lists, Step 5 "One round" a numbered list; longest single line in the file is 1,387 chars (l.96) |
| WS2 | med: dated operator quotes in the handoff template | fixed | `handoff/SKILL.md` template (l.70-120) holds no quote; Writing Rules l.134-139 state the rules without dates or names |
| WS3 | med: init-project 5b paragraphs 1,248 / 2,456, `python3` | fixed | Step 5b is six numbered sub-steps (l.115-124); python3 named as a prerequisite with a `printf` fallback guarded by a `"`/`\` check (l.119); see new finding N3 for the parse check |
| WS4 | med: `--from-log` advertised, nothing reads the log | fixed | `distill/SKILL.md` l.40 reads `feedback.jsonl` (`gap` signal, category x3, `--global` = all projects); `argument-hint: "[--global] [--from-log]"` (l.4); `feedback-log:14` agrees; see N6 for the dangling "Look for:" |
| WS5 | med: dream "reads Claude Code's store" contradicts one store | fixed | `dream/SKILL.md` l.70 now: "the same store and the same path as on Claude Code: `<root>/.agents/memory`" |
| WS6 | low: Step 6.3 synthetic feedback entry | fixed | `distill/SKILL.md` l.158: "a distill never appends to `feedback.jsonl`" |
| WS7 | low: "loaded by /review" wrong | fixed | `review-standards/SKILL.md` l.3 and l.8: "the reviewer agent, or /preflight reading a diff below the reviewer threshold itself; `/review` never reviews by itself" |
| WS8 | low: scar tissue (styrman-harness, quotes, BUG ids) | fixed (residue: N7) | `grep -rn 'styrman-harness\|BUG-1[0-9][0-9]'` over the eight named files: 0 hits; only "(moved from the global context 2026-10-01)" provenance notes remain; harness-internal ids R10 / ADR-044 / D26 and one example quote remain, see N7 |
| WS9 | low: GCP shape prescribed for all three deploy guides | fixed | `preflight/SKILL.md` l.102: the GCP build/push/deploy sentence is scoped "for it that means"; `remote.md` and `local.md` follow their own steps |
| WS10 | low: research always writes integration-spec.md; `--kind plan` | fixed | `research/SKILL.md` l.22 (`--kind spec` first, `--kind plan` when empty), l.47 (`specs/<app>/<topic>.md` through the `spec` skill); stale restatement left in `agents/documentation.md:26`, see N8 |
| WS11 | low: VENDORED "still names skill-repair" stale | fixed | `skills/VENDORED.md` l.40-47: no "still names" claim (grep 0); the find-skills note says verbatim at 435076e |
| WS12 | low: rule 1 "never mention logging" contradicted | fixed | `feedback-log/SKILL.md` l.84: "the one thing said aloud is rule 7's memory refusal, in one line, and the handoff quotes the log's timestamps by design" |
| WS13 | low: four dispatches without budget / early-write clause | fixed (new gaps: N1, N2, N4) | `feature` l.89 (cold reader), l.101 (verifier, with the "after the FIRST command" text), `bug` l.55, `preflight` l.44 and l.57 all carry the budget; the early-write clause is in each |
| WS14 | med: documentation dispatched without refs output | fixed | `review` Step 4 (l.94) and Step 5.2 (l.105), `preflight` Step 4 (l.74), `feature` l.111, `bug` l.55: "carries `sh-agent search --refs <path>` output" |
| WS15 | med: preflight Step 1 skipped documentation and bugs | fixed | `preflight/SKILL.md` l.44: verifier, changelog, reviewer (conditional), documentation with refs output, bugs; l.45 names every undispatched row in the gate report |
| WS16 | med: reviewer row accepted a covering review by range alone | fixed | `preflight` l.59 (reviewer: `## Status: COMPLETE`), l.60 (security), l.61 (verifier PASS), l.88 (blocking list names IN PROGRESS reviewer and security reports) |
| WS17 | low: distill and retro share the dreamer artifact name, no reuse rule | partly | `distill` l.38 reuses `retros/<today>-dreamer.md`; `retro` l.47 has no reciprocal rule, so distill-then-retro overwrites the artifact distill cited (N5) |
| WS18 | low: 1,204-char handoff sentence with inline awk | fixed | `handoff/SKILL.md` l.158-169: six numbered questions, awk in its own fence; I ran the awk on `.agents/docs/plans/*.md` and it printed the one open plan, agreeing with `sh-agent search --open --kind plan` |
| WS19 | low: `--prompts` brief hardcodes two skill paths | fixed | `review/SKILL.md` Prompts brief item 2: "the first that exists; none found means say so in one line and stop" |
| WS20 | low: re-dispatched cold reader overwrite rule | fixed | `feature/SKILL.md` l.93: "a re-dispatched reviewer overwrites its artifact (the `review` skill's rule)" |
| WS21 | low: grilling's unnamed sub-agent | fixed | `skills/VENDORED.md` l.43-44 (explorer) and `review` Prompts brief item 4 (a vendored skill's unnamed sub-agent is covered by the VENDORED notes); the vendored file is unedited (D27 put the keep-list line in the brief, not the vendored prompt-audit file) |
| WS22 | low: feedback-log records project root | fixed | `feedback-log/SKILL.md` l.14 and the Log Format block (`"project": "/absolute/path..."`) |
| G1 | Requirements: nothing elicits or traces | fixed | `skills/requirements/SKILL.md` (elicit, tiers, draft in scratch, approve, ADR); `init-project/assets/REQUIREMENTS.md` has Must / Should / Could / Non-goals / Known limits, no SLA tables up front; init-project l.62 and l.144 route to `/requirements` |
| G2 | Spec: no skill or template | fixed | `skills/spec/SKILL.md` + `init-project/assets/agents-tree/docs/specs/TEMPLATE.md` (seven sections, same order as spec Step 2); `research` Step 4 writes through it |
| G3 | Design / ADR: no skill says what is structural | fixed | `feature/SKILL.md` l.76-85 "Structural decisions and ADRs" with the ADR shape; the ADR now belongs to the decision-holder (D4) |
| G4 | Plan / tasks: nothing sequences several plans, no Task breakdown help | open | `feature/SKILL.md` has no multi-plan or Task-breakdown text; the audit-fix plan has no D-row, strike or follow-up naming G4 (N9) |
| G5 | Test strategy: no skill, no guide template | fixed | `skills/testing-strategy/SKILL.md` (tiers, owed table, fixtures, guide shape); `assets/agents-tree/docs/testing/guide.md` has §1-§5 as the skill names them |
| G6 | Operate: nothing | fixed (skill exists; content not read, belongs to the technical-skills pass) | `skills/operate-gcp/` present; D31 maps G6 to it |
| G7 | Security / dependencies cadence, API auth | fixed (same limit) | `skills/operate-gcp/` (dependency cadence per D31) and `skills/api-auth-gcp/` present |
| G8 | Data: nothing | fixed (same limit) | `skills/data-gcp/` present |
| G9 | Onboarding / explain | fixed | `skills/explain-project/SKILL.md` writes `<docs-root>/onboarding.md`; `agents/documentation.md:3,29,58` carry it in its lane |
| G10 | Learn: over-weighted, automated half never ran on a real home | open | still five skills (feedback-log, retro, distill, dream, project-memory) plus two agents; no plan row, D-row or strike names G10; whether the nightly dream ever ran on a real home is not answerable from the skills (N9) |
| G11 | Size tiers | partly | `context/AGENTS.md:29` "Unscaffolded project" paragraph and spec's one-`specs/<app>/`-per-binary rule exist; plan l.19 lists a per-project size tier as "noted, not built"; nothing for multi-service beyond that |
| G12 | Trail cold reader | fixed | `review --resume` (Step 1, Resume brief, six questions) and `handoff` cold-reader check name each other as ways to run it |
| G13 | Templates | partly | README, spec TEMPLATE, tiered REQUIREMENTS, testing guide, gcp.md are templates under `init-project/assets/`; the plan, handoff, change, bug, review and retro shapes still live inline in `feature:29`, `handoff:70`, `retro:96` and the agents, which is what the gap named |

Count: 22 WS = 21 fixed, 1 partly (WS17), 0 open; 13 G = 9 fixed (G1, G2, G3, G5, G6, G7, G8, G9, G12), 2 partly (G11, G13), 2 open (G4, G10); 0 struck.

### New findings

- **N1** **[medium / confidence: high]** `skills/distill/SKILL.md:38` — the `dreamer` dispatch names "a 30-call budget" but not "write the artifact early, then update it"; `retro/SKILL.md:47` carries both, and the global rule says a duty a subagent must not skip is repeated in the dispatch text. Fix: add `with "write the artifact early, then update it"` after the budget.
- **N2** **[low / confidence: high]** `skills/feature/SKILL.md:113` — Step 6 dispatches `reviewer` and `security` ("`reviewer` (diff only) and `security` when due") with no budget or early-write clause, unlike Step 3 (l.89) and Step 4 (l.101) of the same file and `bug:55`. Fix: one sentence in Step 6 ("each dispatch to reviewer, security or verifier names a tool-call budget below its maxTurns and repeats write the artifact early").
- **N3** **[low / confidence: high]** `skills/init-project/SKILL.md:121` — the "file present without the key" branch validates the edit with `python3 -c 'import json...'` but the `printf` fallback exists only for the "file absent" branch (l.119); without `python3` the command fails, which the step reads as "a failed parse": it removes the inserted line and reports. Safe, but it blames the edit for a missing tool. Fix: state "no `python3`: skip the parse check and report the file's content for the user to read".
- **N4** **[low / confidence: medium]** `skills/preflight/SKILL.md:74` — Step 4 dispatches `reviewer` "with the docs brief exactly as `/review` Step 3 gives it" and names no budget or early-write clause; Step 2's clause (l.57) is scoped to Step 2's dispatches. Fix: add "with the Step 2 budget and early-write clause".
- **N5** **[low / confidence: medium]** `skills/retro/SKILL.md:47` vs `skills/distill/SKILL.md:38` — WS17 is closed one way only: distill reuses today's `retros/<date>-dreamer.md`, retro always dispatches and overwrites it, with a different `Promoted:` and `Since:`, under the record distill's Step 6 cites. Fix: one clause in retro Step 1 source 8: "when `retros/<today>-dreamer.md` exists, read it and dispatch again only when its `Since:` or `Promoted:` differs, writing `-2`".
- **N6** **[low / confidence: high]** `skills/distill/SKILL.md:38-42` — the `--from-log` paragraph (l.40) was inserted between "Look for:" and its table (`### Pattern Types`), so the sentence "Look for:" now introduces a paragraph about the log. Fix: end l.38 with a full stop and move l.40 after the Pattern Types table, or start l.40 with "Under `--from-log`" after the table.
- **N7** **[low / confidence: medium]** `skills/feedback-log/SKILL.md:24,29,30,94`, `skills/distill/SKILL.md:163`, `skills/dream/SKILL.md:22`, `skills/init-project/SKILL.md:113,118`, `skills/project-memory/SKILL.md:3,8` — harness-internal ids a project outside this repo cannot resolve (`R10`, `ADR-044`, "memory-loop D26", `specs/sh-agent/memory-loop.md`), and an operator-phrased example quote ("Is it not retro first, then handoff?"). Fix: drop the ids and keep the rule text; replace the quote with a neutral example.
- **N8** **[low / confidence: medium]** `agents/documentation.md:26` — still says "(an API's is usually integration-spec.md)", the shape `research` stopped writing (WS10). Fix: delete the parenthesis. (Outside this slice's files; recorded for the agents pass.)
- **N9** **[low / confidence: medium]** G4 / G10 with no record: the plan scope says "every finding by ID" and G1-G13 are in scope, yet neither G4 nor G10 has a D-row, follow-up or strike (`grep G4\|G10 audit-fix.md` finds only the scope line). Fix: one D-row each: struck or deferred with a reason.
- **N10** **[low / confidence: medium]** `skills/handoff/SKILL.md:80` — template line `Session ID: <conversation-id>` with "matches the feedback log's session field", but `feedback-log` records only the first 8 characters (l.48) and no step says where an agent reads the id. Fix: say "first 8 characters, as the feedback log does; `none` when the CLI exposes no id".
- **N11** **[low / confidence: low]** `skills/spec/SKILL.md:3`, `skills/requirements/SKILL.md:54,63` — "feature § Decisions" points at the plan's `## Decisions` (D-row shape); the ADR shape is under "Structural decisions and ADRs" (`feature:76`). Fix: name that heading.
- **N12** **[low / confidence: low]** `skills/feature/SKILL.md:52`, `skills/testing-strategy/SKILL.md:37,3` — testing-strategy says "load it before planning a feature's Verification section" and claims the plan "names the tiers from this table", but feature's Step 2 never points at the skill, and likewise `requirements` is not reachable from feature Step 1; a new skill is found only by description match. Fix: one clause each in feature Step 1 and Step 2.
- **N13** **[low / confidence: low]** `skills/handoff/SKILL.md:3-9`, `skills/feedback-log/SKILL.md:3-9` — descriptions are block scalars (`description: |`), not one line; they are under 1,024 characters (about 430 and 520) and `sh-agent list` accepts them, but every other skill uses a one-line quoted value and distill's template (l.67) teaches the block form. Fix: fold both to one quoted line, or accept and say so.

### Measurements

- Longest single line (paragraph as written) per SKILL.md, chars: feature 1,348 (l.91), bug 1,358 (l.55), review 1,387 (l.96), review-standards 855, preflight 974, handoff 659, init-project 1,277 (l.94), research 922, spec 779, requirements 731, testing-strategy 777, explain-project 710, retro 1,475 (l.143), distill 858, dream 704, feedback-log 801, project-memory 769. None exceeds 1,500; retro l.143 is 25 under.
- Blank-line-separated blocks (list items glued together) run longer: init-project 4,575 (l.115, Step 5b's six sub-steps, each its own list item), preflight 3,746 (l.97, Step 6's items), review-standards 3,598 (a checklist), retro 3,325, review 2,168 (the Docs brief list). These are lists of items each under 1,500; none is a prose paragraph.
- `go test ./internal/harness -run TestSkillLinksResolve -count=1`: ok. It checks markdown links only; the non-link section references I sampled in testing-strategy (`docker` skill § 4, `docker/references/setup.md` § Step 2, `typescript/references/testing.md` § E2E, `golang/references/testing.md` § Fakes) all resolve.
- `go run ./cmd/sh-agent list`: lists all 34 skills; every description in scope is under 1,024 characters (largest: review 857, spec 766, testing-strategy 764). Two are block scalars (N13).
- `grep -rn 'styrman-harness\|BUG-1[0-9][0-9]'` over the plan's eight files: 0. `go run ./cmd/sh-agent search --open --kind plan` lists one open plan (this audit's).

## Criteria
1. Dead code — n/a — audit
2. Flow tracing — n/a — audit
3. Standards conformance — n/a — audit
4. Directory layout — n/a — audit
5. Simplicity (DIY-except-security) — n/a — audit
6. Decoupling — n/a — audit

## Not examined
- The bodies of `skills/operate-gcp`, `skills/data-gcp`, `skills/api-auth-gcp` (G6-G8 marked fixed on existence and D31 only; the technical-skills pass owns them).
- `skills/review-standards/references/**` (go-review, typescript-review, vendored prompt-audit) beyond link resolution; `skills/review/scripts/restatement-grep.sh` was not run.
- The `changes/2026-10-01-*` entries were used only as pointers; no behavioural run of any skill on either CLI (matrix rows) was repeated, so "fixed" means present and consistent in the files, not exercised.
- Whether the nightly dream job has run on a real home (G10's second half).
- `grilling`, `codebase-design`, `find-skills` bodies; only their VENDORED notes were read.
- Template files under `init-project/assets/` other than REQUIREMENTS, specs/TEMPLATE and the guide's headings.

## Verdict

Counts over the 35 recon IDs: fixed 30, partly 3 (WS17, G11, G13), open 2 (G4, G10), struck 0. Highest open severity: medium (N1, distill's dreamer dispatch lacks the early-write clause); the two open IDs are low.

This slice is better than on 2026-09-30. All 22 WS items are closed or nearly so, verified in the files: the review brief and fix loop are numbered lists with no line over 1,387 characters (was 3,446), the handoff template carries no dated quotes, `--from-log` has a step that reads the log, dream no longer contradicts the one-store rule, preflight now runs documentation and bugs in Step 1 and refuses a covering review that is not `Status: COMPLETE`, and the BUG-id and project-name scars are gone (grep 0). Of the 13 gaps, nine now have a skill or a template (requirements, spec, testing-strategy, explain-project, operate-gcp, data-gcp, api-auth-gcp, the ADR shape in feature, `/review --resume`), links resolve, and every description is under 1,024 characters. What is left is edge work: a few dispatches without the early-write clause (N1, N2, N4), a retro/distill artifact-reuse rule that works in one direction (N5), harness ids still readable inside global skills (N7), and two gaps (G4, G10) nobody decided on. Nothing found makes a procedure unfollowable as written. "Fixed" here means present and consistent in the files, not re-run on either CLI.
