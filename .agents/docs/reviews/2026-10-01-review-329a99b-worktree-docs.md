---
owner: reviewer
date: 2026-10-01
---
# Review — 329a99b worktree, docs lens   (2026-10-01)

Range reviewed: working tree against 329a99b (`git diff 329a99b` plus untracked files)

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Nine stale restatements of changed claims survive in living files (no high; 6 medium, 3 low). The worst: `skills/handoff/SKILL.md` still carries the old resume order and the old "both empty" rule, so a handoff's cold read runs a different protocol from the context; `AGENTS.md` rule 1 still says the dream timer "sends nothing off" the machine while the README map now draws and ADR-049 records its prompt to the model provider; the two size pins are named only once in three places. All claims in the lens list were checked with `restatement-grep.sh`; the ones that hold are listed under "Checked, no finding".

## Findings
- **[medium / high]** `skills/handoff/SKILL.md:168` — cold-read step: "AGENTS.md → TODO.md → this handoff → ROADMAP.md "Now" → ... → `bugs/INDEX.md` → `sh-agent search --open`". The context (context/AGENTS.md:9) now reads `AGENTS.md` → handoff → ROADMAP "Now" → plan → changes → `bugs/INDEX.md` → `TODO.md` → `search --open` → `memory show`. Change the chain to the new order (add `memory show`), so the check a handoff must pass is the protocol the next session runs.
- **[medium / high]** `skills/handoff/SKILL.md:66` — "both are empty (headers only) at every checkpoint, not just at Done". New rule: TODO empty; ROADMAP holds only "Now" naming the plan while one is open. The sentence's own tail ("while a plan is in progress "Now" names it") contradicts its head. Change to "`TODO.md` is empty at every checkpoint and `ROADMAP.md` holds only "Now" while a plan is open".
- **[low / high]** `skills/handoff/SKILL.md:36` — "the Resume Protocol gives the newest handoff the win on conflict" omits the new exception (an open plan that postdates the handoff). Add "except an open plan that postdates it".
- **[medium / high]** `agents/documentation.md:3,35,50` — description says "ROADMAP.md Now/Next/Later"; line 35 "Next / Later are yours"; line 50 "move items between Next and Later". `Later / Ideas` was removed from this repo's ROADMAP.md and the init-project template; the agent is told to maintain a section that no longer exists (and ideas now go to the handoff). Change all three to Now/Next and drop the Later move.
- **[medium / high]** `AGENTS.md` § How the harness works rule 1, last sentence — "The nightly dream timer ... moves no harness onto a machine and sends nothing off it". README.md:65,71 and ADR-049 now draw the dreamer's redacted prompt to the CLI's model provider (README.md:225 already says memory files are "sent nightly, redacted"). Change to "sends nothing to the feedback repository; the one thing that leaves the machine is the redacted prompt to the CLI's model provider (ADR-049)".
- **[medium / high]** `README.md:71` — "**Dream** (amber too, on each place)". The legend (README.md:65) now makes the dreamer olive, "learn", neither a way in nor back. Change "amber too" to "olive" (and keep the `sh-agent report` carry-back clause, which is the amber part).
- **[medium / high]** `skills/docker/SKILL.md:108` — "only when ... the project's `.claude/settings.json` records the standing permission". context/dev-setup.md:8 moved the record to `docs/deploy/local.md` § Operator steps (memory-loop.md:47 says the same). Change to "`docs/deploy/local.md` § Operator steps".
- **[low / high]** `AGENTS.md:24-25`, `README.md:274-276`, `.agents/docs/testing/guide.md:25-27` — each names only `TestRenderedContextFitsAntigravityBlock` (22,000 bytes); the new `TestRenderedContextUnderTarget` (15,000 bytes, goal 14,500, internal/harness/size_test.go) is named nowhere. Add one clause to each: "and TestRenderedContextUnderTarget pins it under 15,000 (goal 14,500)". The guide's command at :25 runs only the first test by `-run`; extend the regexp to `TestRenderedContext`. README.md:240 and the 22,000 figure are otherwise still true. Also guide.md:24 "must stay under 100 lines" was not re-measured.
- **[low / medium]** `skills/init-project/assets/agents-tree/docs/deploy/local.md:52` — the template `## Operator steps` comment says "standing procedures, never pending items"; the context now makes this section the home of the standing root-step permission record. Add "including the standing root-step permission the operator granted" so a scaffolded project has a place the rule can find.

## Checked, no finding
- Resume order: only skills/handoff restated it (above); README, repo AGENTS.md (its PRIORITY_TODO rule says "before AGENTS.md's own resume chain": still true), changelog agent, dream skill unaffected.
- Stub list: REQUIREMENTS.md:11 (R2) carries it (2.1.277, Bedrock, Vertex, Foundry, telemetry); init-project/SKILL.md:61 restates only the version, consistent.
- Deploy mechanics: `docker` states native-for-smoke and `linux/amd64` for Google (SKILL.md:37-41,109); `cloud-run-deploy` states local build first, Artifact Registry path, git-SHA tag, `latest` alias only (:18-61). dev-setup's pointer holds.
- `/preflight deploys only by the deploy guides` and the version-number rule now live in skills/preflight/SKILL.md:31; documentation.md:55 and deploy/remote.md:203 restate it for their own scope (remote.md's MINOR/PATCH list is the harness's own release rule, consistent).
- `memory-loop.md` § Known differences against R10: REQUIREMENTS.md R10 does not yet carry the path ("in the spec"); expected, the REQUIREMENTS edit awaits the operator. Not a finding.
- `delivery.md` § The model: dreamer colour and provider line are not restated there (no change needed; ADR-049 says untouched).
- Template `deploy/local.md` does have `## Operator steps` (line 50).

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens (restatement sweep is the whole review, above)
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
Code under internal/ and cmd/ (docs lens); `skills/init-project/assets/REQUIREMENTS.md` was read only as a diff for Scale/Latency moves, not against the global context; untracked files: the restatement script sees tracked files only, and `git status` showed none untracked at the start; the rendered context was not rebuilt to confirm byte counts; guide.md:24 (100-line claim) not re-measured; the interactive map page (outside the repo).

## Outcome

All nine findings fixed in the phase (handoff lines 36/66/168; project AGENTS.md rules 1 and 4 with ADR-049; README legend olive; documentation agent out of ROADMAP; docker skill points the standing permission at the deploy guide; the 15,000-byte pin named in AGENTS.md Setup commands, README and guide § 2 with both tests in the -run regexp; the template local.md names the standing permission). Round 2: `2026-10-01-review-329a99b-worktree-docs-r2.md`.
