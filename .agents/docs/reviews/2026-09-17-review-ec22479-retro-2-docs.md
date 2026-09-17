---
owner: reviewer
date: 2026-09-17
---
# Review — retro-2 docs (context/AGENTS.md, AGENTS.md)   (2026-09-17)

Range reviewed: `git diff HEAD -- context/AGENTS.md AGENTS.md` (working tree, uncommitted)

## Verdict: APPROVE WITH NITS

Nothing in the four/five edited sentences is false or unactionable, and the cited bug IDs (BUG-057, 061, 062, 063, 064) all exist and support what is attributed to them at the sentence level closest to each citation. The real gaps are downstream: the DoD step-4 change (main agent may write the changes entry itself) was not propagated to three other checklists that still describe `changelog` as the one to dispatch with no self-write branch, and the new agents/*.md behavioural-testing rule is scoped narrowly enough to exclude exactly the kind of change (a `context/AGENTS.md` main-agent instruction) that BUG-064 — one of the bugs motivating this same diff — actually was, and whose own repro used a different, undocumented technique.

## Findings

- **[medium / confidence: high]** `skills/bug/SKILL.md:53`, `skills/feature/SKILL.md:98`, `skills/preflight/SKILL.md:40`, `skills/preflight/SKILL.md:83` — none of these restate the new DoD step 4 option ("or you write it yourself from the diff when you already hold the evidence"); each still names `changelog` as the step without qualification, as if dispatching it were the only route. `skills/feature/SKILL.md:88` already carries the self-write option ("`changelog` agent or yourself") so it is clearly a deliberate change, not an oversight elsewhere — per review-standards §3 ("every other restatement of a changed rule is changed with it"), these four are restatements left standing. Grep `changelog` across `skills/` to confirm before editing.

- **[medium / confidence: medium]** `AGENTS.md` § Testing, new sentence — "A change to `agents/*.md` also gets a behavioural run: `claude -p --agents <json …>` … the real home is never touched." This exactly matches how BUG-062 and BUG-063 were verified (`.agents/docs/testing/2026-09-17-bug-062-repro.md`, `-bug-063-adr-032-repro.md`, both using `--agents <json>` built from the working-tree file). But BUG-064 — cited in the very same diff's context/AGENTS.md sentence, and whose fix commit (`6b135da`) touched only `context/AGENTS.md`, `skills/review/SKILL.md`, `skills/preflight/SKILL.md`, never `agents/verifier.md` — was verified with a *different*, undocumented technique: `--append-system-prompt` on the main agent from the real `$HOME` (`.agents/docs/testing/2026-09-17-bug-064-repro.md`, Method line). The new Testing sentence covers only `agents/*.md` changes, so a future change to `context/AGENTS.md`'s main-agent behaviour (exactly BUG-064's class of change) has no standing rule requiring a behavioural run at all — only the generic install+status+parse-log check, which BUG-064 itself demonstrated would not have caught the intermittent-following defect. This is a real coverage gap in the doc the diff is editing, not a contradiction of fact.

- **[low / confidence: medium]** `context/AGENTS.md`, Subagent-duty sentence — "…because on Claude Code a rule inside an agent's definition was followed in some runs and not others while the dispatch text held every time it was tested (BUG-062, 063, 064)." BUG-062's actual fix (per its own Addendum) was ADR-032 deleting the `<app>`-resolution decision entirely, not moving the duty into dispatch text; the "dispatch text held every time it was tested" claim is demonstrated only for BUG-064 (`bug-064-repro.md`, 3/3 vs 2/3). Citing BUG-062 alongside 063/064 as evidence for the dispatch-text remedy overstates what that record shows; it is fair support only for the weaker, separate claim ("a rule inside an agent's definition was followed in some runs and not others"), which BUG-062's before/after and its addendum do support.

## Criteria

1. Dead code — n/a (prose-only diff).
2. Flow tracing — pass for the sentences examined individually; see findings above for the one place a claim in a sentence (Testing addition) does not trace forward to a consequence (context/AGENTS.md changes get no equivalent behavioural-run rule).
3. Standards conformance — fail on the "every other restatement changed with it" rule (finding 1); the bug-ID citations are otherwise accurate to their records.
4. Directory layout — n/a.
5. Simplicity — n/a, no code/dependency change.
6. Decoupling — n/a.

## Not examined

`agents/changelog.md` body beyond the grep already done (read in full; no restatement of "changelog always writes the entry" that needed changing there — it already allows the main agent to keep other files current and doesn't claim exclusivity over the changes entry). Vendored skills (out of scope per `skills/VENDORED.md`). The `chr(92)`/escape-quoting sentence in the Hard rules bullet (no bug ID cited, read and judged plausible/actionable on its face, not independently tested). Full read of `agents/documentation.md`, `agents/reviewer.md`, `agents/security.md`, `agents/verifier.md`, `agents/explorer.md` for other restatements of the DoD step-4 changelog rule or the subagent-duty rule beyond the ones grepped. `context/AGENTS.md` sections other than the four diff hunks and their immediate neighbors. `.agents/docs/testing/2026-09-17-verify-bug-064.md` (found but not opened).

## Outcome

- Finding 1 (the self-write allowance not in the skills that name `changelog`): fixed — `skills/bug` Step 5, `skills/feature` Step 6 and `skills/preflight` (the pending-work dispatch and the release entry) now say "`changelog` agent or yourself", matching `skills/feature` Step 4 and context § Definition of Done step 4.
- Finding 2 (the Testing addition covers agent definitions only, while BUG-064 was a context change): fixed — project `AGENTS.md` § Testing adds that a `context/` change meant to change the main agent's behaviour is run the same way through `claude -p --append-system-prompt`.
- Finding 3 (the dispatch claim cited all three bugs): fixed — the sentence now says the definition rules failed intermittently (BUG-062, 063) and that the dispatch sentence made the verifier write its guide in every run tested (BUG-064).
Re-installed into `.agents/scratch/retro2-home`: drift none; `go test -count=1 ./...` passes.
