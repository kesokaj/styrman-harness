---
owner: reviewer
date: 2026-09-13
---
# Cold read — plan `.agents/docs/plans/2026-09-13-review-skill.md`   (2026-09-13)

Range reviewed: plan file + named surfaces (skills/review-standards, skills/preflight, skills/audit, agents/explorer.md, agents/reviewer.md, agents/security.md, agents/verifier.md, context/AGENTS.md, README.md, AGENTS.md). No diff — this is a cold read of a plan document, not a code diff.

## Verdict: design needs revision / implementation not ready

The plan (D1–D9) is unusually thorough and its cross-references to the surfaces it names check out — every `/audit` mention outside `.agents/docs` (agents/explorer.md:3, context/AGENTS.md:69, skills/review-standards/SKILL.md:3 and :8, skills/audit/SKILL.md itself) is accounted for in the Surfaces-touched list and the Phase 1 grep check, and the "27 skills" count (README.md:3, README.md:87, AGENTS.md:7) is consistent with one skill removed and one added. The worst problem is that the plan defines a user-facing range syntax (`<a>..<b>`, two-dot) but never pins down whether the agents it dispatches actually diff two-dot or three-dot — and the agents' own stated default is three-dot, which computes a materially different diff once branches have diverged. A second real gap: the plan defines an error case for `--project` combined with another scope token but leaves every other cross-category combination (a path plus a range, a path plus a PR, two PRs) undefined.

## Gaps — must be decided before implementation

1. **Two-dot vs three-dot diff semantics for a range/ref target.** D2 defines the syntax as `<a>..<b>` (two-dot) and "a single `<ref>` (meaning `<ref>..HEAD`)" — two-dot again. D3's dispatch brief for that case is the bare text `Scope: <a>..<b>`, with no git command spelled out (contrast the working-tree and PR briefs, which both spell out the exact command: `git diff HEAD -- <files>`, `gh pr diff <n>`). But `agents/reviewer.md:28` ("Establish the diff with `git diff`, `git diff --staged`, or `git diff <base>...HEAD`") and `agents/security.md:26` ("Establish the diff (`git diff`, `git diff --staged`, `git diff <base>...HEAD`)") both state the agents' own default as **three-dot**. `git diff a..b` and `git diff a...b` differ whenever the two refs have diverged — the common case for reviewing a feature branch against a main that has moved on — so which one runs changes what gets reviewed. Cite: plan D2, D3 (range clause); `agents/reviewer.md:28`; `agents/security.md:26`.

2. **Cross-category target combinations are undefined.** D9 defines exactly one combination error: "`--project` together with paths, a range, or a PR is an error ('one scope')." It says nothing about a path token given together with a range token, a path with a PR, two PR numbers, etc. (e.g. `/review internal/api main..HEAD --code`). Two implementers could reasonably choose to error the same as the `--project` case, silently take the first/last classified token, or attempt to combine them (diff of a range restricted to paths) — these produce different behavior. Cite: plan D9.

3. **`Range reviewed:` line for `--project` runs.** Every callee agent's report template (reviewer, security, verifier) has a mandatory `Range reviewed: <base>..<head>` line, but `--project` has no base/head — D3's own text ("the whole repository at HEAD, no diff") is a reasonable value to drop in, but the plan never says to use it there; a minor documentation gap, not blocking.

## Readiness — questions a first-pass implementer would have to ask

**BLOCKING**
- Q1 (ties to gap 1): For a range or single-ref target, does the skill instruct the dispatched agents to run `git diff <a>..<b>` (matching the literal two-dot syntax in D2) or `git diff <a>...<b>` (matching the agents' own stated default convention)? Two careful implementers land on different diffs for the same command.
- Q2 (ties to gap 2): What happens when a user supplies tokens from two different scope categories that D9 doesn't cover (path + range, path + PR, two PRs, etc.)? Error like the `--project` case, or some other resolution? Two careful implementers build different parsers here.

**NON-BLOCKING** (the plan states a workable default, or a reasonable implementer decides alone)
- Q3: Exact `Range reviewed:` text for `--project` — D3's own scope text ("the whole repository at HEAD, no diff") is the obvious value to reuse.
- Q4: For the `--docs` lens under `--fix`, is the docs-lens reviewer re-dispatched in the same round after `documentation` applies its edits, to check convergence, same as the code/security/simplify lenses? D5's general "apply the findings … then re-dispatch the same lens agents" reads as yes.
- Q5: Are `--fix` rounds joint across every active lens (one shared round counter, all lenses re-dispatched together each round) or independent per lens with separate counters? D5's phrasing ("re-dispatch **the same lens agents**", plural) most naturally reads as one joint loop; worth stating explicitly in the skill text but a competent implementer would land here anyway.
- Q6: Precedence when a single token both exists as a path and validates via `git rev-parse --verify` (e.g. a directory sharing a branch's name). D9's rule list order (flag → PR → ref/range → path, in that order) reads naturally as first-match-wins, resolving this without another decision.
- Q7: Exact check for "a PR of another branch without checking it out" (Failure modes section) — comparing `git rev-parse HEAD` against `gh pr view <n> --json headRefOid` is the obvious mechanism; not spelled out but not contentious.

## Not examined

Did not read the prior review artifacts already on disk (`2026-09-13-cold-read-review-skill-*`, `2026-09-13-review-cold-read-review-skill-2.md`) — instructed to read fresh, with no prior context, so this report may repeat or diverge from earlier rounds' findings by design. Did not run any commands beyond static grep/read (no e2e, no sandbox install, no attempt to actually build `skills/review/SKILL.md` since it does not exist yet — Phase 1 is unstarted). Did not check every reference in `.agents/docs/testing/skill-matrix.md` beyond the `audit` row and the new `review` row the plan names. Did not review `context/dev-setup.md` beyond a grep for "audit" (no hits).
