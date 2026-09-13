---
owner: reviewer
date: 2026-09-13
---
# Cold read (round 3) — `/review` plan   (2026-09-13)

Range reviewed: `.agents/docs/plans/2026-09-13-review-skill.md` (pre-implementation cold-reader check; no diff exists yet — `skills/review/` has not been written) plus the surfaces it names: `skills/audit/SKILL.md`, `skills/review-standards/SKILL.md`, `skills/preflight/SKILL.md`, `agents/explorer.md`, `agents/reviewer.md`, `agents/security.md`, `agents/verifier.md`, `agents/documentation.md`, `context/AGENTS.md`, `README.md`, root `AGENTS.md`, `REQUIREMENTS.md`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/architecture/DECISIONS.md`, `.agents/docs/bugs/INDEX.md`.

Six earlier cold reads of this same plan exist on disk (`2026-09-13-cold-read-review-skill-{1,3}.md`, `2026-09-13-review-cold-read-review-skill-2.md`, and round 2's three parallel reads `2026-09-13-cold-read-review-skill-r2-{1,2,3}.md`). I re-verified, against the live files rather than the prior reports' word, that this revision actually fixes what round 2 flagged as BLOCKING, then looked independently for anything new.

## Verdict: design ready / implementation ready

Every BLOCKING item raised across the three round-2 reads is now genuinely resolved in the plan text, and I confirmed each against the live tree rather than taking the plan's word for it:

- **`--fix` actor ambiguity** (r2-3's blocking gap): D5 now states outright "The MAIN agent applies every edit — `reviewer` and `security` have no edit tools and never touch source," which matches both agents' frontmatter (`tools: Read, Grep, Glob, Bash, Write` — no `Edit`) and their bodies ("You never modify source, config, or tests" / "you never edit source"). The `--docs` lens correctly routes to `documentation`, the one agent in this set that does carry an `Edit` tool (`agents/documentation.md` frontmatter). Fixed.
- **`--fix` on a single lens, round-history destination** (r2-2's blocking gap): D4 now lists "any run with `--fix`" as its own trigger for the merged record, independent of lens count, so the plan's own flagship example (`/review #42 --fix`, single default `--code` lens) gets a merged record to hold `Round N` sections. Fixed.
- **Verify-command false positives** (r2-1's medium finding): I ran the exact Phase-1 verify command, `git grep -nE '(^|[ \`(])/audit\b' -- . ':!.agents/docs' ':!skills/anti-ui-slop' ':!skills/gke-basics'`, against the current tree. It returns exactly the files/lines the plan's Surfaces list and D8 already name for editing (`agents/explorer.md:3`, `context/AGENTS.md:69`, `skills/audit/SKILL.md:3,14`, `skills/review-standards/SKILL.md:3,8`) — nothing from the two excluded vendored skills. After Phase 1's edits (delete `skills/audit/`, repoint the other three), this command will print nothing, exactly as the plan claims. Fixed.
- **Second "vendor's built-in" location in `skills/preflight/SKILL.md`** (r2-1 and r2-3 both flagged this): I ran `git grep -n "vendor's built-in"` and got exactly two files, three lines: `skills/audit/SKILL.md:3` (deleted with the directory) and `skills/preflight/SKILL.md:3` and `:50`. D8 explicitly names both the Step 2 body (`:50`) and the frontmatter description (`:3`). Fixed — nothing left unaccounted for.
- **D2/D3 untracked-file inconsistency** (r2-1's medium finding): current D2 now reads "none → the working tree, exactly as D3 defines it," inlining D3's literal text rather than restating a different rule. Fixed.
- **Skill-count arithmetic**: `ls -d skills/*/ | wc -l` = 27 today (including `audit`); delete one, add `review`, stays 27, matching `README.md:3` and `AGENTS.md:7`. Confirmed correct.
- **ADR-017 numbering**: the newest entry in `.agents/docs/architecture/DECISIONS.md` is ADR-016, so ADR-017 is free. No collision.

## Gaps

None found. Every interface, failure mode, and scope edge I checked resolves to a single reading, and the plan's own citations (D8's exact quoted replacement text, D3's literal dispatch strings, D4/D5's explicit actor and trigger rules) leave no fork between two careful implementers.

## Readiness — remaining questions (all NON-BLOCKING)

- D9 gives four token classes (flag, PR, range/ref, path) and two explicit errors (`--project` plus another scope selector; empty target) but not what happens to a token matching none of the four (a typo target, a nonexistent path/ref). A reasonable implementer errors naming the token, consistent with the spirit of the two errors already specified; nothing in the plan suggests silently treating it as a literal path.
- D9's precedence when a single token could satisfy two classes at once (an existing directory that shares its name with a valid git ref) is implied by the order the classes are listed in (flag → PR → range/ref → path) but not stated as an explicit rule. The collision is rare and a competent implementer follows the listed order.
- Same-day, non-`--fix` re-runs of `/review` against the same scope will overwrite the underlying per-lens agent artifacts (`agents/reviewer.md`, `agents/security.md`, `agents/verifier.md` templates are `$(date +%F)-<verb>-<slug>.md` with no same-day collision handling, unlike D4's explicit `-2`/`-3` rule for the merged record only). This is pre-existing behavior across every agent in the harness, not something this plan introduces, and D4's merged-record collision handling already covers the one durable artifact; a reasonable implementer accepts the overwrite.
- The `≤ 140 lines` budget on `skills/review/SKILL.md` (Phase 1, task 1) against the volume of D3's verbatim dispatch-brief text that D3 itself says must be "quoted in the skill" is a real fit risk. It is not a design fork (nothing here changes what gets built), so I have not counted it as a gap or question — only flagging it as an implementation risk worth watching in Phase 1's own verify step.

## Not examined

- Did not draft any part of `skills/review/SKILL.md` to test the line-budget risk noted above.
- Did not read `internal/harness/plan.go` / `derive.go` in depth; relied, as all six prior reads did, on `AGENTS.md`'s own rule that adding/deleting a skill is a file drop needing no installer change.
- Did not run any part of the Verification tiers (no code exists yet to install or smoke-test).
- Did not re-examine `commands/` or `.agents/rules/` beyond confirming both directories are empty.
- Did not check `agents/changelog.md` or `agents/bugs.md` bodies — not touched by this plan and not needed to assess it.

## Criteria (applied to the plan document, not to code — no diff exists)

1. Dead code — pass: `skills/audit/` is correctly slated for full deletion (`git rm -r`), not a comment-out; the verify grep confirms no reference survives outside the two vendored exceptions.
2. Flow tracing — pass: target → scope text → per-lens dispatch → artifact slug → merged-record trigger is consistent end to end, including the `--fix` actor question and the `--project`/`--all` dispatch count, both verified against the agents' actual tool lists.
3. Standards conformance — pass: plan format matches the `feature`-skill shape; frontmatter has `owner`/`date`.
4. Directory layout — pass: `skills/review/SKILL.md` lands where every other workflow skill lives.
5. Simplicity — pass: no new dependency; reuses the four existing subagents; the `--fix` loop is bounded and its edit responsibility is now assigned per lens without inventing new machinery.
6. Decoupling — pass: the router (`/review`) computes scope and hands agents literal text without editing any agent file; the one open seam from round 2 (the reviewer/security "never edit" rule vs. `--fix`) is now explicitly routed around rather than left implicit.
