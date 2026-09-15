---
owner: reviewer
date: 2026-09-15
---
# Review — sweep-context-agents   (2026-09-15)

Range reviewed: the working tree at d1ff612 plus uncommitted edits (not a diff). Named files read whole: context/AGENTS.md, context/dev-setup.md, agents/bugs.md, agents/changelog.md, agents/documentation.md, agents/explorer.md, agents/reviewer.md, agents/security.md, agents/verifier.md. Lens: `--docs` with the `--project` classification.

## Verdict: REQUEST CHANGES

The worst problem is a split owner for ROADMAP.md, which leaves the documentation agent too narrow to follow the rule that says what it owns. The global context gives the main agent the whole of ROADMAP.md and requires it to be empty at Done (context/AGENTS.md:18), and it gives `documentation` "the project `AGENTS.md` setup and deploy sections" plus the deploy-path trigger (context/AGENTS.md:17). `agents/documentation.md` still claims ROADMAP Next and Later, and its description still claims "Now". It has no lane entry for the AGENTS.md sections: its Records rule (:52) limits its edits to its lane, so it cannot write those sections. It also has no deploy-path trigger: its stop rule (:12) tells it to refuse any dispatch that is not about behavior, interfaces or the plan. Five more statements have fallen behind the rule they restate:
- reviewer: blob pinning is missing.
- bugs: its dispatch triggers are narrower than the context's.
- changelog: "Resume Protocol" is a stale section name, and it says ROADMAP belongs to documentation.
- deploy guides: the gcp guides have no "Operator steps" section.
- reviewer: the bug-index path is hard-coded.

`explorer`, `security` and `verifier` are accurate. None of these files breaks a REQUIREMENTS.md item.

## Findings

- **[medium / confidence: high]** divergent — `agents/documentation.md:30` and `:43` (Next / Later are documentation's; "move items between Next and Later") and `agents/changelog.md:3` ("never writes specs, README, or ROADMAP; those belong to documentation") ↔ `context/AGENTS.md:18` ("you, the main agent, keep `ROADMAP.md` and `TODO.md` current, and at Done both are EMPTY … neither file is a parking lot for ideas"). The context's documentation ownership list at `context/AGENTS.md:17` does not include ROADMAP. Commit e5833a3 ("one TODO/ROADMAP owner and both empty at Done") added both texts, so the split was never resolved. What it causes: two writers for one file, and a documentation agent that keeps Next and Later filled when the context requires them empty. **Should say:** documentation.md:30/:43 drops its ROADMAP claim, or at most proposes the change in its reply; changelog.md:3 says "ROADMAP (the main agent's)".
- **[medium / confidence: high]** divergent (inside one file) — `agents/documentation.md:3`: the description says it maintains "ROADMAP.md Now/Next/Later". The same description, and `:30`/`:36`, say "Now" belongs to the main agent. This wording dates from e521826 and was never updated. The description is the text a dispatcher reads. **Should say:** no ROADMAP claim (see the finding above), and nothing about "Now".
- **[medium / confidence: high]** dropped: `context/AGENTS.md:17` ("the project `AGENTS.md` setup and deploy sections", added in d87df30). `agents/documentation.md:20-32` (the lane) has no entry for it, and `:52` (Records) restricts its edits to "the files in the lane above", so the agent's own rules forbid the edit the context gives it. **Should say:** add `AGENTS.md (project root): Setup commands and Deploy sections only` to the lane at :24-31, and cover it in the Records line at :52.
- **[medium / confidence: high]** narrowed: `context/AGENTS.md:17` ("when behavior, interfaces, the plan, or a deploy path changed", from 1aba06a) → `agents/documentation.md:12` ("You fire when behavior, interfaces, or the plan changed … if dispatched for one, say so and stop") and `:3` (the description has the same three triggers and does not list the deploy guides). A dispatch for a change to a deploy step only can be refused under :12. **Should say:** "behavior, interfaces, the plan, or a deploy path" in both places; the description also names the deploy guides.
- **[medium / confidence: high]** narrowed: `context/AGENTS.md:31` ("a test failing for a reason other than the change under test, a fix that took more than one attempt", from 87a862d) → `agents/bugs.md:3` ("Dispatch after a failing test that took more than one attempt", the 5063167 wording, never updated). Two context triggers became one narrower trigger. A test that fails for an unrelated reason on the first try, and a non-test fix that needed several attempts, both drop out of the description the dispatcher reads. `skills/bug/SKILL.md:3` already uses the new wording. **Should say:** "a test failing for a reason other than the change under test, a fix that took more than one attempt, a regression, …", in the same order as context/AGENTS.md:31.
- **[medium / confidence: medium]** divergent — `context/AGENTS.md:16` ("a saved diff file only when the tree will stand still, and the reviewer pins the blobs it read", added in 50ac771) ↔ `agents/reviewer.md:28`. The reviewer's Scope section says nothing about saved diff files or pinning blobs, and no other living doc does either (`git grep -i blob` finds nothing outside vendored or reference files). The context gives the reviewer a duty its definition does not carry. BUG-014's prevention rule (INDEX.md:20) says a global rule alone does not reliably reach a subagent. **Should say:** reviewer.md:28 adds "Given a saved diff file, record the blob id (`git hash-object` / `git ls-files -s`) of every file you read on the `Range reviewed:` line."
- **[low / confidence: high]** stale — `agents/changelog.md:84` "The Resume Protocol reads the three newest entries". The context heading was renamed from "Resume protocol — first action in any existing project" to "Resume, then checkpoint" (`context/AGENTS.md:7`) in 87a862d. The same stale name is at `skills/handoff/SKILL.md:150` and `:154` ("the Resume Protocol in the global context"). **Should say:** "Resume (context: *Resume, then checkpoint*) reads …" in all three places.
- **[low / confidence: high]** divergent — `agents/documentation.md:28` and `:48` ("Each deploy guide carries an `## Operator steps` section"; ADR-018 item 7) ↔ `skills/init-project/assets/agents-tree/docs/deploy/gcp.md`, whose headings are Recorded defaults, Image, Procedure, Secrets and configuration, Verify, Rollback and Infrastructure, with no Operator steps (the local and remote templates have one at :46 and :39). This repo's `.agents/docs/deploy/gcp.md` has none either. **Should say:** the gcp template gains `## Operator steps`, and this repo's "Not applicable" gcp guide gains the section with "none"; or documentation.md:48 names gcp as the exception.
- **[low / confidence: medium]** divergent — `agents/reviewer.md:32` reads `.agents/docs/bugs/INDEX.md` at a fixed path, but the same file resolves its docs root with a fallback to `docs/` (`:54-57`), as the `bugs` agent does (`agents/bugs.md:16`). In a project that uses `docs/`, the check for a recorded bug being reintroduced reads a file that does not exist. **Should say:** "`<docs-root>/bugs/INDEX.md`", with the docs root resolved first.
- **[low / confidence: medium]** divergent (outside the named files; `context/AGENTS.md:18` points at the file) — `skills/init-project/assets/ROADMAP.md:13-17` asks for "1-3 upcoming milestones" under Next and "Unscheduled directions worth keeping visible" under Later / Ideas. Its own comment at `:3` and `context/AGENTS.md:18` ("neither file is a parking lot for ideas"; ideas go in the handoff) say otherwise. **Should say:** the Next and Later comments point at plan files and the handoff instead of inviting entries.

### Unverified

- `context/dev-setup.md:8`: the new condition "AND that project's `.claude/settings.json` records the standing permission" (cceee54). No living doc or skill says how a project records that permission: init-project creates no settings.json, and this repo has no `.claude/`. I could not check what satisfies the condition.
- `context/AGENTS.md:15`: the verifier "runs … every tier the plan names". `agents/verifier.md:37` runs only what AGENTS.md and the testing guide name, `:51` says "Never deploy anything", and the cloud smoke tier runs inside `/preflight` (`skills/feature/SKILL.md:78`). I did not check whether dispatches pass the plan's tiers.
- `agents/changelog.md:61` (BUG-034 in styrman-crew): that repo was not read.
- Claude Code and Antigravity runtime claims (`context/AGENTS.md:37`, "/skills and /agents list both layers"; the Antigravity loading path at `:3`) were not exercised; the paths match `internal/harness/plan.go:58-94`.

### Observation (not a docs finding)

- The installed `~/.claude/AGENTS.md` is behind the repo: the rendered line 81 (the root-steps sentence) is the pre-cceee54 version, without the condition. The operator's real home has not been updated since cceee54. The standing operator step is `.agents/docs/deploy/local.md:43`.

## Classification (--project)

| File | Class | Proof |
|---|---|---|
| context/AGENTS.md | accurate | Deploy paths match `internal/harness/plan.go:88-94`. The retro cadence at :41 matches `cmd/sh-agent/main.go:267-271` and `internal/harness/state.go:81-90`. BUG-011 and BUG-012 at :58/:33 match `bugs/INDEX.md:17-18`. All ten skills named at :71 exist. The render is 93 lines, under the 100-line budget in `testing/sh-agent/guide.md:23`. It is the newer side of every divergence above. |
| context/dev-setup.md | accurate | `--set-secrets` at :13 matches `agents/security.md:35`; the amd64 build at :10 matches `agents/verifier.md:46`. The only open point at :8 is listed as unverified. |
| agents/bugs.md | outdated | `agents/bugs.md:3` against `context/AGENTS.md:31` (the trigger list). |
| agents/changelog.md | outdated | `agents/changelog.md:84` against `context/AGENTS.md:7` (Resume Protocol); `:3` against `context/AGENTS.md:18` (ROADMAP owner). The template at :39-68 matches the newest entry, `changes/2026-09-15-research-beads-elephant-goldfish.md:7-28`. |
| agents/documentation.md | outdated | `:3`, `:30`, `:43` against `context/AGENTS.md:18`; `:12`, `:20-32`, `:52` against `context/AGENTS.md:17`. |
| agents/explorer.md | accurate | ":3 scan phase of /review --project" matches `skills/review/SKILL.md:50`; the plans/ path at :40 matches `skills/review/SKILL.md:52`; BUG-012 behaviour at :65 matches `bugs/INDEX.md:18`. |
| agents/reviewer.md | outdated | `:28` against `context/AGENTS.md:16` (blob pinning); `:32` against its own `:54-57`. The "two questions" at :28 matches `skills/feature/SKILL.md:66`; "strongest model" at :12 matches `skills/preflight/SKILL.md:50`. |
| agents/security.md | accurate | The dispatch triggers at :3 match `context/AGENTS.md:16`; the ban at :33 matches `context/AGENTS.md:49`; the findings path at :46 matches `skills/review/SKILL.md:52`. |
| agents/verifier.md | accurate | The guide rule at :49 matches `testing/sh-agent/guide.md:7-9` and `agents/documentation.md:34`; BUG-014 at :90 matches `bugs/INDEX.md:20`; the tools at :4 are all in `internal/harness/derive.go:15-22`. |

No file is redundant.

## Re-check commands

```bash
git grep -n 'Now/Next/Later\|Next / Later are yours\|between Next and Later' -- agents skills context
git grep -n 'README, or ROADMAP; those belong to documentation' -- agents
git grep -n 'setup and deploy sections' -- context agents
git grep -n 'You fire when \*\*behavior, interfaces, or the plan\*\*' -- agents
git grep -n 'failing test that took more than one attempt' -- agents skills context
git grep -n -i 'pins the blobs\|blob' -- agents/reviewer.md context/AGENTS.md
git grep -n 'Resume Protocol' -- agents skills context README.md AGENTS.md
git grep -n 'Operator steps' -- skills/init-project/assets .agents/docs/deploy agents/documentation.md
git grep -n '\.agents/docs/bugs/INDEX.md' -- agents
git grep -n 'upcoming milestones\|worth keeping visible' -- skills/init-project/assets
```

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- The skills were not read whole. Only the passages the named files point at were read: `review` (Steps 3-5), `feature` (Steps 2-5), `bug` (headings and gates), `init-project` (tree and symlink steps, ROADMAP template, gcp template headings), `preflight:50`, `handoff:150-154`, `retro` (`--harness`), `research:45`.
- `commands/` and `rules/` are empty; the project has no `.agents/rules/`.
- `.agents/docs/specs/` is empty.
- `testing/skill-matrix.md` and the dated logs were not compared.
- The external styrman-crew bug ids were not checked.
