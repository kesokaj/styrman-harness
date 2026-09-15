---
owner: reviewer
date: 2026-09-15
---
# Review — coldread-harness-practices-r1-b   (2026-09-15)

Range reviewed: the plan file `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md` (untracked, read in full at HEAD `c192622`) and the surfaces it names, as they are at `c192622`. Cold-reader check, round 1, reviewer B. The reading started at I-11 to I-17 and Phases 3 and 4, then covered the rest.

## Verdict: REQUEST CHANGES

`design needs revision` / `implementation not ready`.

The worst problem: the Phase 3 and Phase 4 checkpoints expect the new agent and skill behaviour to be visible in this session's own dispatches, and it will not be. This repo has no `.claude/` directory (`ls -la .claude` finds nothing; only `.agents/skills -> ../skills` exists). So on Claude Code, `verifier` and `feature` load from the installed home, and installing into that home is explicitly out of scope (plan: Scope, Out). A checkpoint that expects "its log carries a counts line" can never be met as written.

Also broken:
- One verify check (5.3) cannot pass.
- The fail-closed release gate (I-19) keys on a report section that the reports it reads do not have.
- ADR-022 would land before ADR-021 in an append-only file.
- The G17 count rule and the plan's own Failure modes give opposite answers to the same case.

What holds: every Old anchor checked occurs exactly once (69 of them, including all of I-11 to I-17), `render | wc -l` prints 93 today, and the positive controls of 1.5, 2.3 and 2.5 do print lines at `c192622`.

## Findings (gaps)

- **[high / confidence: high]** Plan Phase 3 Checkpoint ("G12 is in the skill from Phase 2"), Phase 4 Checkpoint ("its log carries a counts line"), D6 — **edited files are not the loaded definitions.**
  - A session in this checkout on Claude Code loads `agents/verifier.md` and `skills/feature/SKILL.md` from the installed home. The repo has no `.claude/agents` or `.claude/skills` link.
  - I-15 and I-6 therefore do not reach the verifier dispatched at the Phase 4 checkpoint, or the feature skill driving Phase 3. The old verifier (`agents/verifier.md:71`) writes no `counts:` line, so the Phase 4 expect cannot be met. That leaves the next phase closed under the plan's own Step 4.
  - Decide one of these and write it on the checkpoint:
    - put the new rule text into the dispatch brief;
    - dispatch from the sandbox home after a scratch install;
    - move the expectation to Phase 6 (the `verifier-counts` row).
- **[medium / confidence: high]** `skills/preflight/SKILL.md:50` versus plan 5.3 — **5.3 cannot pass.**
  - 5.3 expects every line printed by `git grep -n "50+ lines or 5+ files" -- context skills agents` to name the bug case.
  - I-18 edits `preflight:40` but not `preflight:50` ("when the release diff is 50+ lines or 5+ files"). That line will print without the bug case.
  - Decide whether the release-diff reviewer threshold gains the bug case, or whether 5.3 excludes that line with a reason.
- **[medium / confidence: high]** I-19 against `skills/preflight/SKILL.md:75`, `skills/review/SKILL.md:56` and `agents/reviewer.md:23` — **the gate still fails open for most reports.**
  - A report counts as covering only if its "Not examined section" leaves the path out. Only agent artifacts have a section by that name (reviewer:23, security:21, verifier:21).
  - The reports the gate reads for coverage use other names:
    - `reviews/*-preflight-*` reports have "What was not checked" (preflight:75).
    - Merged `/review` records have "Not checked" (review:56).
  - A report without a section literally named "Not examined" names no gap, so it covers everything. That is fail-open for exactly those reports.
  - Fix: name all three headings in I-19, or rename the sections to one name in the same task.
- **[medium / confidence: high]** `.agents/docs/architecture/DECISIONS.md:157` (newest is ADR-020) versus 4.4 and 5.3 — **ADRs land out of order.**
  - ADR-022 is appended in Phase 4 and ADR-021 in Phase 5. The file is append-only (preflight Step 3 checks this), so it would read 020, 022, 021.
  - `changelog` numbers ADRs itself and may call the first one ADR-021, which would clash with I-18's text.
  - Fix: swap the numbers (G24 becomes ADR-021), or move ADR-021 earlier.
- **[medium / confidence: high]** I-15 "Test counts" versus Failure modes line "A G17 count comparison fails across machines: the FAIL line names the environment cause and stands" — **the two contradict.**
  - I-15 says a fall that the Environment section explains is not a FAIL line. The Failure modes line says an environment-caused fall is a FAIL that stands.
  - This repo has six environment-dependent skips (`internal/harness/harness_test.go:526,542,606,789`, `cmd/sh-agent/main_test.go:288`, `internal/harness/pathguard_test.go:120`), so it will happen between the Linux box and the Mac.
  - Pick one.
- **[medium / confidence: medium]** I-20 "A headless run with no user turn that reaches a gate … writes the question … into the handoff's Next Steps" versus the E2E edge case and the I-23 `bug-no-repro` check `grep -rqi "repro" .agents/docs/handoffs/` — **a correct agent can fail the edge-case check.**
  - The matrix drives `claude -p "<prompt>"` (`skill-matrix.md:13`), where the prompt is a user turn. I-20 also says "do not ask me anything" is not a waiver.
  - So an agent following the text can ask in its output and stop, with no handoff, and the check fails.
  - State whether a `-p` or `agy -p` run counts as "no user turn", or check `$OUT` for the question as well.
- **[medium / confidence: medium]** `agents/bugs.md:12,24,59,69,80` versus I-17 — **I-17 leaves describing text that contradicts it** (global hard rule: a behaviour change includes the text that describes it).
  - :69 would read "that is the one allowed edit to a past record" followed at once by a second allowed edit.
  - :80 still says "the only later edit to a past record is its `status`".
  - :59 still says "One table".
  - :12 and :24 take the next number "from `INDEX.md`". Once rows move to a Graduated table, the highest id can sit outside the main table, and the risk is a reused number.
  - Add these to I-17.
- **[low / confidence: high]** Plan Verification, golden path, versus the I-23 `feature` check — **the check is weaker than the criterion.** The criterion requires `expect:`, Checkpoint, a Cold-reader check box and `## Decisions`. The check greps only `Checkpoint` and `## Decisions`, so "yes" can be reported without the other two.
- **[low / confidence: high]** 3.1 and 3.2 — **two checks do not test what they claim.**
  - The first grep of 3.1, `## Round N`, already prints `1` at `c192622` (`skills/review/SKILL.md:60`), so it passes whether or not I-12 was applied. Grep a phrase that only I-12 adds, such as "no ledger line is still in flight".
  - 3.2 checks one of the four I-14 edits. `go-review.md`, `typescript-review.md` and `golang/SKILL.md` are not verified at all.
- **[low / confidence: medium]** I-1 against the 1.2 positive control — **the control cannot tell a working brief from a broken one.**
  - `git grep` does not follow symlinks, and `.agents/skills` here is a link to `skills/`. The I-1 brief searches `.agents/{rules,skills,…}/` and is dispatched at 1.2, before I-2 widens the living docs.
  - A brief that cannot see `skills/` at all would also produce a `dropped:` finding, so the control passes either way.
- **[low / confidence: medium]** I-22 and 5.1 — **the precedence test cannot support outcome A's wording.**
  - There is a control marker for skills but none for agents.
  - The test writes a Claude-format `.agents/agents/reviewer.md`, while Antigravity's own agent shape is `agents/<name>/agent.md` (`internal/harness/plan.go:60`).
  - "No tools allowed" is not given as a flag; on Claude Code, removing all tools can also hide the agent descriptions.
  - A CLI that never reads the project agent yields outcome A, whose New text says "because the vendors resolve a clash differently". The test would not have shown that.
- **[low / confidence: medium]** I-17 retro Step 6 — **the row move fires in only one path.**
  - The move is dispatched only under "Project scope, approved".
  - It never happens when:
    - a harness-scope proposal is applied later by `/retro --harness` (the bug's index lives in another project);
    - a rule graduates outside a retro (`context/AGENTS.md`, "A rule that recurs graduates…").
  - ADR-022's Decision says "when a bug's prevention rule graduates", which covers all of these.
- **[low / confidence: medium]** I-19 "A project's `AGENTS.md` may record other thresholds, and then those apply" — **replace or add is undecided.** If they replace the defaults, a project can switch off NO-SHIP on high-severity findings, which undoes "fails closed".

## Questions

BLOCKING:
1. Phases 3 and 4: how do the checkpoint dispatches get the new verifier and feature texts, when this repo has no `.claude/` link and real-home installs are out of scope?
2. 5.3: does `preflight:50` (the release-diff threshold) gain the bug case, or is it excluded from the check?
3. G17: is an environment-explained count fall a FAIL (Failure modes) or not (I-15)?
4. G21 and E2E: does a `claude -p` / `agy -p` run count as "no user turn"? Where must the repro question land for `bug-no-repro` to pass?

NON-BLOCKING:
5. ADR order: swap the numbers or reorder the phases. Default: number them in landing order.
6. I-10 else branch for `learning-loop:136`: the box state is unstated. Default: `[-]`, otherwise 2.4 fails.
7. I-14 `golang-ecosystem.md`: the Old and New code spans have mismatched delimiters (three backticks open, two close). The literal is clear from `:323`.
8. I-17 retro item 6: "on or after the newest retros record" when none exists. Default: all records.
9. I-15: count indented subtest `--- PASS:` lines? This repo has no `t.Run` calls, so it only matters downstream. Also, run the tests once with `-v` or twice? The template implies once.
10. I-22: the exact "no tools" flag on each CLI.
11. I-19: do project thresholds replace the defaults or add to them?

## Criteria (applied to a plan)

1. Dead code: pass. I-7 removes a sentence whose content the generated inventory already prints (`internal/harness/inventory.go:321`; render contains "can also be typed as" once).
2. Flow tracing: fail.
   - Loaded definitions are not the edited files (finding 1).
   - I-19 heading mismatch.
   - `preflight:50` is not traced by I-18.
   - `bugs.md` numbering is not traced by I-17.
3. Standards conformance: fail on the global describing-text rule (`agents/bugs.md`). The rest conforms: no Go code, `REQUIREMENTS.md` is untouched, and the `reviewer` description stays one quoted line.
4. Directory layout: pass. Every new file lands in `testing/`, `reviews/` or `plans/` under `.agents/docs`.
5. Simplicity (DIY-except-security): pass. Text edits only, no dependency, no security surface.
6. Decoupling: pass with one note. I-4 then I-20, and I-12 then I-20, edit the same sentences across phases. The plan orders them, and the anchors survive: I-12's New keeps the I-20 anchor.

## Not examined

- The full prose of I-6 blocks A and B, beyond the anchor, grep and ordering checks.
- `skills/init-project/assets/` beyond the template command lines.
- The research record beyond rows G17 to G24.
- The BUG-018 record body.
- The `security.md` and `explorer.md` bodies beyond the anchor line.
- Whether `skills/golang/references/*` has third-party provenance (only `skills/VENDORED.md` was checked).
- Whether the sandbox home and both CLI logins currently work.
- Matrix rows other than `feature`, `bug-no-repro` and `verifier-counts`.
- `handoffs/2026-09-14-session-close.md` beyond its Next Steps.
- Whether the I-8 grep and "no `changes/` entry" checks can run in a project without `<docs-root>/plans/`.
