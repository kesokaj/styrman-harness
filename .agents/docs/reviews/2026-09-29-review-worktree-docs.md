---
owner: reviewer
date: 2026-09-29
---
# Review — worktree-docs   (2026-09-29)

Range reviewed: working tree — changed files: .agents/docs/testing/skill-matrix.md, skills/VENDORED.md, skills/review-standards/SKILL.md, skills/review/SKILL.md; new/staged: skills/review-standards/references/prompt-audit/LICENSE.txt, skills/review-standards/references/prompt-audit/prompt-audit.md

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The change adds a fourth review lens (`--prompts`, vendoring Anthropic's `prompt-audit` reference) and is documented consistently everywhere the harness enumerates `/review`'s lenses: the skill's `description`, `argument-hint`, usage block, Step 1's flag list, the Step 3 lens table, and Step 5's per-lens fix behaviour in `skills/review/SKILL.md` all now include `prompts`/`--prompts`. `README.md`, `context/AGENTS.md`, `skills/preflight/SKILL.md` and `agents/*.md` do not enumerate the lens set at all (they reference `--docs` once, for a specific unrelated rule, or "lens" in the six-criteria sense), so none of them needed an update. The vendored file is byte-identical to upstream at the pinned commit (verified by re-fetching `8a1541c` and comparing sha256: `c33200da…1878f2` matches) and the `VENDORED.md` row's path, pin, and license are all correct. The one real gap is a paraphrase in `skills/review/SKILL.md`'s `--prompts` brief that silently drops a qualifier the vendored source states explicitly, which would make the reviewer skip the sub-agent-roster check in the single most common case this lens is used for.

## Findings

- **[medium / confidence: high]** `skills/review/SKILL.md:52` — the `--prompts` brief says "Its Group 4 is `not applicable` unless the scope holds request-building code," but the vendored source it summarizes (`skills/review-standards/references/prompt-audit/prompt-audit.md:173`) says Group 4 is not applicable when there's no request-building code "**except the sub-agent roster check, which applies to those agents' definition files**." Since `--prompts` scope is most often exactly `agents/*.md` / `.agents/agents/` with no request-building code (the harness's own primary use case per the new skill-matrix scenario, which scopes to `AGENTS.md` and `.agents/rules`), the brief's flattened wording tells the reviewer to skip the redundant-sub-agent-roster check in the case it was written for. Fix: restate the brief's Group 4 sentence to keep the roster-check exception, e.g. "Its Group 4 is `not applicable` unless the scope holds request-building code, except the sub-agent roster check, which still applies to any agent definition file in scope."

- **[low / confidence: low]** `.agents/docs/testing/skill-matrix.md:31` — the new `review-prompts` row grants `--add-dir <home>/.claude/skills/review-standards`, but no other `/review` scenario in this file (`review`, `review-shadow`, `review-shadow-name`, `review-project`, rows 26–30) grants that directory even though they equally dispatch the `reviewer` agent with the `review-standards` skill (which itself loads `references/go-review.md` / `references/typescript-review.md`). If skill references are normally reachable without `--add-dir`, this row's grant is either unnecessary scaffolding or evidence the other rows are under-permissioned for their own reference reads. I could not run the sandbox to settle which; flagging as unverified rather than asserting a defect.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens (the new reference lands under `skills/review-standards/references/prompt-audit/`, consistent with the existing `references/go-review.md` / `references/typescript-review.md` siblings and with `AGENTS.md`'s "Third-party skills are vendored verbatim... never edit them in place" rule — this is vendored verbatim as a reference, not a skill, and `VENDORED.md` records it)
5. Simplicity (DIY-except-security) — n/a — docs lens (stdlib/dependency question doesn't apply to a markdown reference; no new code dependency introduced)
6. Decoupling — n/a — docs lens

## `git grep` sweep for restatements of the changed lens set

The diff changes the enumerated lens list (`--code|--security|--simplify|--docs|--tests|--all` → adds `--prompts`) in four places, all in `skills/review/SKILL.md` itself (already updated together, verified above). Grepping the living docs (`README.md`, `AGENTS.md`, `context/`, `agents/`, every skill dir under `skills/` not in `skills/VENDORED.md`, `commands/`, `rules/`, `.agents/docs/specs/`, `.agents/docs/testing/guide.md`, `.agents/docs/deploy/`) for the old lens enumeration or for `--docs|--tests|--all` lens lists found no other restatement — no `dropped:` or `narrowed:` findings beyond the Group-4 paraphrase above, which is not a restatement of the *lens list* but a restatement of the *vendored file's own Group-4 rule*.

No rule, command, flag, or caveat was removed by this diff (it is purely additive except for the four places where the lens token list gained one more item in-line) — no `dropped:` findings.

## Not examined

- Whether `.agents/docs/testing/2026-09-29-verify-prompts-lens.md`'s in-flight verifier run (Status: IN PROGRESS at time of this review) or the concurrent behavioural probe in `.agents/scratch/prompts-probe/` actually confirm the `--prompts` lens works end to end — that is the verifier's job, not this docs review's, and its log was still IN PROGRESS when read.
- The `skill-matrix.md` `--add-dir` question above (low finding) — would need a sandbox run to confirm whether it's redundant or the other rows are missing it.
- `skills/review-standards/references/prompt-audit/LICENSE.txt` was not byte-verified against upstream (fetched the wrong URL, got a 404 from GitHub); its content is a standard Apache-2.0 license text with "Copyright 2026 Anthropic, PBC.", consistent with `VENDORED.md`'s Apache-2.0 claim, but the exact source URL was not located and confirmed.
- The plan file `.agents/docs/plans/2026-09-29-research-memory-and-prompt-audit.md` recommended vendoring under a new `skills/prompt-audit/references/` skill directory; the actual change instead nests the reference inside the existing `skills/review-standards/` skill. This is a deviation from that plan, but plans are excluded from the grep target the dispatch specifies (only living docs count), and no requirement pins the exact directory, so it is noted here for context rather than raised as a finding.

## Outcome
- `skills/review/SKILL.md:52` Group 4 paraphrase drops the roster-check exception — fixed: the lens row now reads "except its sub-agent roster check, which runs over every agent definition in the scope (two agents doing the same task with the same tools and near-duplicate prompts are one finding with a concrete roster edit)". Probe: a near-duplicate agent pair seeded in the fixture is found in 3 of 3 runs with the corrected brief; the one control run with the old wording found it too, so the contradiction between brief and reference is fixed while the predicted silent disablement did not reproduce (0 of 1) — `testing/2026-09-29-prompts-lens-behavioural.md` round 2.
- `skill-matrix.md:31` `--add-dir` unlike the other `/review` rows — rebutted: the other rows dispatch a reviewer whose `skills:` preload (`agents/reviewer.md:7-8`) carries `review-standards` into its context, while the `--prompts` brief makes it Read a file outside the project, which headless Claude Code refuses without `--add-dir` (`testing/2026-09-03-verify-sandbox.md` § Vendor behaviour: "reads outside the cwd need `--add-dir`"); the row's extra column now states this reason.
- Not examined, `LICENSE.txt` provenance — done by the main agent: fetched from `raw.githubusercontent.com/anthropics/skills/8a1541c4a3ffa5a20a5a91de0dcf3f0bab1d1ef4/skills/claude-api/LICENSE.txt` (11,345 bytes, Apache-2.0).
- Not examined, plan versus change (the research record recommended a `skills/prompt-audit/` skill) — left: the operator chose option 3, a `/review` lens, on 2026-09-29; the research record now says so under its recommendation.
