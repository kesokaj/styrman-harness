---
owner: reviewer
date: 2026-10-05
---
# Review — worktree-docs   (2026-10-05)

Range reviewed: working tree vs HEAD (ad2130b): .agents/docs/architecture/DECISIONS.md, AGENTS.md, README.md, context/AGENTS.md, rules/critical-thinking.md (docs lens)

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

No living doc states the old facts about `rules/` (grep for "empty", "none", old Layout wording finds only the two lines this scope already changed) and the numbers in ADR-050 reproduce (`go test ./internal/harness -run TestRenderedContext`: hand-written 14,700, Antigravity block 16,815 = 15,653 + 1,162). What remains is restatement drift of the new reply/approval protocol in skills, plus an undocumented new rule-file key.

## Findings
- **[medium / confidence: high]** `skills/feature/SKILL.md:21` and `skills/dream/SKILL.md:56` vs `context/AGENTS.md:5` — the context now defines the approval vocabulary as `go all approved` / `go approved 1-3, skip 5`; the feature skill still tells the agent to end a list with "Say go or no-no per number, or go all." and dream says "go per number, or go all". Three different answer phrasings for the same numbered-table gate; the context says only the context's approval counts ("Nothing proposed runs until they approve it ... and only that"), so an agent following a skill offers a phrase the context's protocol does not name. Fix: have the two skills point at the context's phrasing, or state in the context that "go" forms are equivalent.
- **[low / confidence: high]** `skills/requirements/SKILL.md:30` → `context/AGENTS.md:5` — "Ask in the form the global context prescribes ... never a long preamble, never options without a recommendation": the context no longer says either clause (the struck sentence was "always briefly, never options without a recommendation"), and the form now opens with a one-line explanation and one-line recommendation before the table, which this skill's paraphrase omits. `skills/spec/SKILL.md:43` and `skills/harness-audit/SKILL.md:67` restate the table columns only (not wrong, but omit the opening lines and the approval step). Fix: requirements:30 drops the dead quote and points at the context.
- **[low / confidence: medium]** `context/AGENTS.md:5` vs `context/AGENTS.md:56` (Hard rule "Lead with what is broken") — two openers for a reply that reports: "opens with a one-line explanation and a one-line recommendation" vs "Lead with what is broken". Same file, no precedence stated. The reviewer/security agent reply formats (`agents/reviewer.md` "Output") follow their dispatch so they are not affected, but the main agent has two rules for the first line. Fix: one clause saying the explanation line is what is broken when something is.
- **[low / confidence: medium]** `README.md:240` — the rules row names Claude's `paths:` and Antigravity's `trigger: glob` + `globs:` but not `trigger: always_on`, the key the first shipped rule (`rules/critical-thinking.md:2`) uses to be always-on on Antigravity. A reader adding the second rule has no documented key for always-on. Fix: add `trigger: always_on` to that row.
- **[medium / confidence: high]** `.agents/docs/architecture/DECISIONS.md` ADR-050 Consequences (ADR is a record, not edited here; reported because the pointed-to log disagrees) — it says the behavioural runs "are in `testing/2026-10-05-verify-critical-thinking.md`", but that log's `## Status: FAIL` (Antigravity row P2: the rule stopped a direct instruction that the control carried out 3 of 3; P3 table-and-recommendation form 0 of 1; P4 no valid run; the log also says the anti-sycophancy rule shows no difference from the control on Claude Code). The context text "A direct instruction is its own approval" is therefore not yet true on Antigravity. The new log is untracked; the change must not ship as "behaviour verified" until the operator decides the P2 fix. (Not a stale doc line, a claim whose evidence is red.)

## Update after the first write
- The tree changed while I reviewed: `skills/dream/SKILL.md`, `skills/feature/SKILL.md` and `skills/requirements/SKILL.md` now carry edits I did not make. They address findings 1 and 3 (the skills now cite the context's `go all approved` / `go approved 1-3, skip 5` and its Operator-paragraph form). Read only as a diff, not re-reviewed in full. One defect in them: `skills/feature/SKILL.md:21` now reads `...skip 5") The answers ARE the scope` with the sentence's full stop lost (low, cosmetic).
- `go test ./...` (run once, tail read): `internal/harness`, `internal/dream`, `internal/memory`, `internal/search` ok; `cmd/sh-agent` FAILED with a 600.012s timeout panic. I did not read the cause or re-run it; a slow or hung test unrelated to this scope is possible, but it is unverified and the log cannot be called green.

## Checks done with no finding
- `AGENTS.md:43` (Layout) and `README.md:278` now say rules/ has one file; no other living file says rules/ is empty or "by design" empty (`git grep` over README, AGENTS, context, skills, agents, specs, guide, deploy, REQUIREMENTS, internal, cmd).
- Size claims: `.agents/docs/testing/guide.md:27-31`, `AGENTS.md:24-26`, `README.md:288-289` restate the 15,000 / 22,000 pins and goal 14,500; unchanged and still true. ADR-050 figures reproduced.
- "disagree before complying" was moved from context to `rules/critical-thinking.md:14`; no other living copy.
- `README.md` map and delivery spec: no path moved (ADR-050 says so; install/update/report/release untouched), rule 3 not triggered.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- Cause of the `cmd/sh-agent` 600s test timeout (see update); not re-run.
- `.agents/docs/specs/` read only by grep for the changed claims, not whole; `deploy/*.md` likewise.
- Antigravity behaviour claims (rule loads as always-on at `~/.gemini/config/rules/`) taken from the verify log, not re-run.
- The rest of the verify log (below row P4) and `changes/` newest entry (none yet describes this scope).

## Outcome

- [medium] feature/dream approval phrasing — fixed: `skills/feature/SKILL.md:21` and `skills/dream/SKILL.md:56` now name the global context's approval line ("go all approved", "go approved 1-3, skip 5").
- [low] requirements:30 dead quote — fixed: `skills/requirements/SKILL.md:30` points at the Operator paragraph's form (one-line explanation and recommendation, then the table) and drops the quote.
- [low] two openers (context:5 vs context:56) — rebutted: not a conflict; the explanation line of `context/AGENTS.md:5` is where "what is broken" goes, and `:56` names what leads, not the line's position; no reply can satisfy one and break the other.
- [low] README:240 lacks `trigger: always_on` — fixed by `documentation`: `README.md:240` now names it.
- [medium] ADR-050 cites a FAIL log as evidence — fixed: ADR-050 Consequences states the log is `Status: FAIL` and why (Antigravity P2; no measured difference for the anti-sycophancy rule on Claude Code); the P2 fix is left: operator (handoff Next Steps).
- [nit, round-1 update] `skills/feature/SKILL.md:21` lost its full stop — fixed at `skills/feature/SKILL.md:21`.
