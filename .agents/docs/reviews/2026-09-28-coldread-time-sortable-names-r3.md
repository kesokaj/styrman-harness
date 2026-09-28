---
owner: reviewer
date: 2026-09-28
---
# Cold-read review (round 3) — time-sortable trail names   (2026-09-28)

Range reviewed: `.agents/docs/plans/2026-09-28-time-sortable-trail-names.md` (plan document, cold-reader check — no code diff exists yet). Earlier rounds: `.agents/docs/reviews/2026-09-28-coldread-time-sortable-names-r1.md`, `-r2.md`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Both rounds' blocking findings are closed in this text. Round 1's gap (changelog agent has no Bash) is resolved by D3: the dispatcher passes the `date +%F-%H%M` stamp, the agent's `tools:` list is untouched, and I independently verified `agents/changelog.md:4` still reads `Read, Grep, Glob, Edit, Write`. Round 2's blocking gap (did the Phase 2 handoff-skill probe rename the copied file's frontmatter `name:` to avoid colliding with the installed `handoff` skill) is now answered explicitly: "the copy's frontmatter line `name: handoff` is changed to `name: handoff-probe`." Round 2's medium finding (`/init-project`'s directory-missing vs. directory-unlistable ambiguity) is now worded to require a two-step check ("A directory that does not exist counts zero; for one that exists, run `ls <dir>`…"), which is enough for a reasonable implementer to add a `[ -d ]` guard without asking — no longer blocking. I re-verified every file:line citation in `## Surfaces touched` and re-ran the describing-text sweep myself (`git grep -nE '(changes|handoffs)/[^ \`]*(date|YYYY|<)' -- context agents skills commands rules README.md AGENTS.md`); it returns exactly the four hits the plan names. One small, non-blocking citation error survives from round 2 unfixed, and one pre-existing inconsistency in the surrounding text the plan chooses not to touch is worth a decision.

## Findings

- **[low / confidence: high]** `agents/changelog.md` — `## Surfaces touched` still cites `:23, the date rules` for the line that becomes "Use the stamp given in the dispatch prompt…", but that text is rule 1 at line 22 (verified by direct read); line 23 is rule 2 ("If absent, run `date +%F` yourself."), which `## Interfaces` says stays unchanged. Round 2 flagged the same citation and called it non-blocking because the quoted replacement text in `## Interfaces` (lines ~82–84) unambiguously identifies rule 1 by content, and Phase 1's verify clause (`grep -nE 'HHMM|%H%M|<stamp>'`) is line-agnostic. Still true here — not blocking, but it was not fixed between rounds and is worth a one-line correction before Phase 1 so an implementer scanning by line number doesn't edit rule 2 by mistake.
- **[low / confidence: medium]** `agents/changelog.md:23` itself ("If absent, run `date +%F` yourself.") is an instruction the agent cannot execute — its `tools:` frontmatter (`:4`) has no Bash, and the plan's own D3 reasoning states this explicitly for rule 1's replacement. The plan leaves rule 2 "as they are, for the date only" without noting that it is already unreachable in production (a HEAD-level defect, not introduced here). This is out of the plan's diff and not blocking implementation of this feature, but since the plan is already editing the surrounding rule block and is aware of the Bash constraint, it's worth a decision: fix the wording here (e.g. "derive it from filesystem context" per rule 3's own fallback) or file a BUG record so it isn't silently carried forward again.

## Criteria

1. Dead code — n/a, no code in this plan; the one candidate (agents/changelog.md's unreachable "run date yourself" fallback, finding 2 above) is pre-existing and outside the diff, not introduced by this plan.
2. Flow tracing — passes at every hop checked: context sentence → `context/AGENTS.md:65` (verified); handoff skill `:43,:62,:145` (verified, `:62` correctly untouched); changelog agent dispatcher → agent rule chain (`preflight:86` → `agents/changelog.md:3` description → rule 1, content-verified despite the `:23` citation slip); `/init-project` `:22,:106` (verified, now resolves round 2's ambiguity). The five-line grep (`HHMM|%H%M|<stamp>`) independently covers every changelog-agent edit including the lane-tree comment at `:30`, which the broader describing-text sweep's regex cannot reach (a space after `changes/` breaks that pattern) — redundant coverage, not a gap.
3. Standards conformance — the describing-text sweep is exact (re-run, 4/4 hits match); no other restatement of the naming pattern found outside the plan's own list.
4. Directory layout — n/a, no new files or directories.
5. Simplicity — pass; no new dependency; the `/init-project` grep pipeline is the minimum needed and its directory-existence ambiguity from round 2 is now resolved by wording alone.
6. Decoupling — pass; the five touched surfaces (context, handoff skill, changelog agent, init-project, preflight) are independently editable and the sweep step is the cross-file drift check.

## Answers to the three questions

**(1) Gaps.** None blocking. Two small, non-blocking items: the `agents/changelog.md:23` citation (should read `:22`, unfixed since round 2) and the unreachable "run `date +%F` yourself" fallback the plan knowingly leaves in place while editing the same block (worth a BUG record, not a plan blocker). Everything else — file:line citations, the sweep's 4-hit count, the two previously-blocking gaps from rounds 1 and 2 — checks out against the live files.

**(2) Scope.** Nothing beyond `## Why`'s approved scope items 1–8 and the D-rows' documented responses to rounds 1–2 (D3 for the Bash question). No new endpoint, shape, or unused configuration.

**(3) Readiness.** Zero BLOCKING questions.
- NON-BLOCKING: fix `agents/changelog.md`'s Surfaces citation from `:23` to `:22` (content in `## Interfaces` already disambiguates; a careful implementer follows the quoted text, not the number).
- NON-BLOCKING: decide whether to correct or bug-record the dead "run `date +%F` yourself" fallback at `agents/changelog.md:23`; leaving it as-is changes nothing this plan's tests check.
- NON-BLOCKING: `/init-project`'s exact `[ -d ]` wording, left to the implementer per round 2 (now resolvable from the text alone).

## Not examined

- No git diff exists yet — this is a plan, not a change; nothing beyond the plan text and the current-state files it names.
- `.agents/scratch/agy-probe/` and `.agents/scratch/trail-search/` (the measured probe data backing `## Why` and the rejected `INDEX.md` alternative) — not opened; their headline numbers are taken from the plan's own text, as in rounds 1 and 2.
- Claude Code's native project-skill slash-command resolution mechanics (directory name vs. frontmatter `name:`) — not verified against product docs or a live run; round 2 flagged this as a question and the current plan text now answers it (rename the frontmatter) without my needing to resolve the underlying mechanic myself, since the plan's chosen answer (rename both) is safe regardless of which key the loader actually uses.
- `internal/harness/size_test.go` was read but not executed in this round; round 2 already ran it and recorded 20,329 of 22,000 B, which I did not re-run.

design ready
implementation ready
