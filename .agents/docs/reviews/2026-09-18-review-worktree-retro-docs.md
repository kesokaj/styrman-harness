---
owner: reviewer
date: 2026-09-18
---
# Review — worktree-retro-docs (2026-09-18)

Range reviewed: working tree diff via `git diff HEAD -- context/AGENTS.md skills/feedback-log/SKILL.md agents/verifier.md skills/feature/SKILL.md`

## Verdict: REQUEST CHANGES

Docs-lens review (no source, config, or tests in scope). Nothing is factually wrong in the new text, but two restatements of the "write the artifact early" rule were left on the old wording after `agents/verifier.md` replaced it, the new feedback-log triggers 6/7 have no matching `signal` vocabulary value, and two of the four cited events in `context/AGENTS.md`'s Hard rules have no BUG record backing them despite the project's own "every failure gets a number" rule.

## Findings

- **[medium / confidence: high]** `agents/reviewer.md:20` and `agents/security.md:18` — both still read "Write your report file as soon as you have one finding worth recording, then UPDATE it as you go...", the exact phrasing `agents/verifier.md` just replaced with the per-command IN-PROGRESS-log rule. Criterion 3 ("every other restatement of a changed rule is changed with it") applies: these are the same "write the artifact early" rule for sibling read-only agents, now inconsistent with the verifier's stricter version. `agents/explorer.md:17` carries the same old wording too but explorer writes a brief, not a pass/fail log, so it is the weaker case of the three — still worth a decision either way, not a silent skip. `restated: agents/reviewer.md:20, agents/security.md:18 (also agents/explorer.md:17) still say the old thing; agents/verifier.md:18-19 has the new thing.`

- **[medium / confidence: medium]** `skills/feedback-log/SKILL.md:28-29` (new triggers 6 "Redirection" and 7 "A missed filing") vs. the Signal Types table (`skills/feedback-log/SKILL.md:73-78`, values `correction|override|restatement|format|gap`) — neither new trigger maps cleanly to an existing signal value. Trigger 6 ("the user changes what is being built... even when no word of it is 'no'") is closest to `override` but `override` is documented narrowly as "User changed agent's choice to something else" (tool/pattern), not scope/direction. Trigger 7 ("the user asks for something a rule says should already exist... the ask IS the correction") is closest to `gap` ("Agent lacked knowledge it should have had") but the trigger's own text says the agent *did* know the rule and just didn't file it — that's not a knowledge gap. Since `rule_ref`/`signal` values are the controlled vocabulary `/retro` matches on across projects, an entry logged under triggers 6/7 will be coerced into a signal that doesn't describe it, or left ambiguous. Either add `redirection` and a value for trigger 7 to the table, or the SKILL.md states explicitly which existing value each maps to.

- **[low / confidence: medium]** `context/AGENTS.md` Hard rules — of the four sentences this scope changes, two cite BUG ids that exist (`BUG-078`, `BUG-079` on the "verified by the shell" sentence; `BUG-080` on the "read-only mode" sentence — both confirmed present in `.agents/docs/bugs/INDEX.md`). The other two cite dated events with no BUG id at all: "(three rounds on one mechanism in styrman-harness, 2026-09-18, each passing a fresh-fixture test)" on the Review bullet, and "(four mutations read as 'fails as it must' on a build error, styrman-harness 2026-09-18)" on the mutation-testing bullet. `git grep` of "fresh-fixture", "fresh fixture", "build error", "mutations", "declared and not used" against `.agents/docs/bugs/INDEX.md` returns nothing — no BUG-NNN record documents either event, unlike the other two sentences in the same diff that do cite one. The project's own rule ("Every failure gets a number... a fix that took more than one attempt, a review finding that was a real defect... Dispatch `bugs`") applies to both: three rounds on one mechanism and four mutations misread as passing are exactly the failure classes that rule names.

## Not a finding (checked, clean)

- `agents/verifier.md` § Write the artifact EARLY: the "roughly two thirds of your budget" sentence survives unchanged in the diff (still present, untouched by the replacement).
- `context/AGENTS.md`'s `## Installed` block and skill/agent counts: unaffected by this scope — no skill, agent, command, or rule was added or removed, only prose edited in place. No update needed, correctly left alone.
- README.md § Install without cloning, `.agents/docs/deploy/local.md` § Verify, `.agents/docs/deploy/remote.md` (lines ~46-52) all already state the same "verify from the shell that will run it" procedure (`hash -r`/`rehash`, `type sh-agent`, `which -a`) that the new `context/AGENTS.md` sentence describes, each citing BUG-078/BUG-079. Consistent.

## Criteria

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — FAIL: rule restated in `agents/reviewer.md` and `agents/security.md` not updated alongside `agents/verifier.md`'s changed wording (see finding above); feedback-log's new triggers not reflected in its own signal vocabulary table.
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- `skills/feature/SKILL.md` diff (added clause about a test fixture in a verify clause not producing its claimed condition): checked it has no stale restatement elsewhere via targeted grep of "verify clause" and "test fixture" wording, but did not exhaustively grep every skill under `skills/` for a paraphrase of this specific clause — only the four named-in-brief targets (context/AGENTS.md, feedback-log, verifier.md, feature/SKILL.md sentences) got full grep treatment.
- Did not check `.agents/docs/specs/` or the newest `changes/` entry against this diff in detail beyond the bug-index and deploy/README checks above — no changes/ entry appears to reference this exact diff yet (working tree, uncommitted).
- Did not run any test suite (docs-only scope; none applicable).

## Outcome

- Finding 1 (`agents/reviewer.md:20`, `agents/security.md:18`, `agents/explorer.md:17` keep "one finding worth recording"): left — owner: the next `/retro`. Proposal 6 was the verifier's: its passing commands read as "no finding yet", which is why its rule now names the first command. For a reviewer, security or explorer the finding IS the unit, and whether they want the per-step form is a separate proposal; noted in the handoff.
- Finding 2 (triggers 6 and 7 without a signal value): fixed at `skills/feedback-log/SKILL.md:28-29` (each trigger names its signal) and `:79` (`redirection` row).
- Finding 3 (two sentences citing no BUG id): fixed — the events are recorded as BUG-081 and BUG-082 (`bugs/INDEX.md` Graduated table) and cited at `context/AGENTS.md:26` and `:56`.
