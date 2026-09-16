---
owner: reviewer
date: 2026-09-16
---
# Review — retro-2026-09-16-applied   (2026-09-16)

Range reviewed: working tree, six named paths — `agents/changelog.md`, `agents/reviewer.md`,
`context/AGENTS.md`, `skills/review-standards/SKILL.md`, `skills/golang/SKILL.md`,
`skills/init-project/SKILL.md` (the seventh proposal, `.agents/docs/bugs/INDEX.md`, is out of scope
per the dispatch). Blob shas as reviewed:

| File | blob sha |
|---|---|
| `agents/changelog.md` | `15d569ba9bcedd9190c139df9f5bd3cc090cf9a8` |
| `agents/reviewer.md` | `2e032a23aac97ba4f92a52f912f861da918f908b` |
| `context/AGENTS.md` | `b574997bd857a1a7ec0673af8a3729a06c7ce6db` |
| `skills/review-standards/SKILL.md` | `ee41f4305954ac7cfc8d7228c0e100c119b0c259` |
| `skills/golang/SKILL.md` | `98c6c9e45cdc879a50250bdd4fa646331df69127` |
| `skills/init-project/SKILL.md` | `6899f707c5d7480d87e5799bb439b16b29b64fd2` |

Cross-checked against `~/.config/sh-agent/proposals/2026-09-16-retro.md` (Proposals 1–8) and
`.agents/docs/retros/2026-09-16-retro.md`. Note: the proposals file was rewritten mid-review to mark
all eight rows "applied … the harness commit carrying [this review file]" before this review or any
commit existed. I did not act on that text — it names nothing for me to change and asks nothing of
me — but it is worth the operator's attention: a queue file describing its own future commit before
the commit exists is the kind of claim the global context's "a claim that something works states HOW
it was verified" rule exists to catch, and it sits outside a file this review touches.

## Verdict: APPROVE WITH NITS

Six of the six diffs match their proposals verbatim, land in the section each proposal named, and
parse cleanly (`go run ./cmd/sh-agent list`: 27 skills, 7 agents, no error). The one real problem is
that Proposal 2's own subject — "TODO.md and ROADMAP.md are EMPTY, and when" — has a restatement in
`skills/handoff/SKILL.md` that this proposal set did not touch, so the trail now contradicts itself
on the exact question Proposal 5 (grep every restatement) was written to stop recurring. Everything
else checked out clean.

## Findings

- **[medium / confidence: high]** `skills/handoff/SKILL.md` Step 3 item 3 — still reads "Bring
  `ROADMAP.md` and `TODO.md` current yourself (Definition of Done step 4): **when the task is done**
  both are empty; while a plan is in progress … `TODO.md` holds only items the plan will consume."
  This directly contradicts the new `context/AGENTS.md:18`: "both are EMPTY (headers only) **at every
  checkpoint, not just at Done**." A reader of the handoff skill alone (which is exactly when this
  fires — at a handoff, itself a checkpoint) is told emptiness is a Done-only obligation while a plan
  is running; a reader of the context is told otherwise. This is the precise class of bug Proposal 5
  (`review-standards § 3`'s new grep-every-restatement item, and BUG-003/007/019/027/042) targets,
  and it is not one of the six files in scope, so the proposal set that fixed the anchor left its
  clearest restatement unfixed. Fix: update `skills/handoff/SKILL.md` Step 3 item 3 to match — "both
  are EMPTY at every checkpoint, not just at Done" — in the same change that lands this queue.

## Checks that came back clean (no finding)

- **Step-4 wording vs. the Resume/checkpoint paragraph** (`context/AGENTS.md:11`) — the checkpoint
  paragraph names what triggers a checkpoint but makes no claim about TODO/ROADMAP emptiness, so
  there is nothing for the new step-4 sentence to contradict there.
- **Step-2 new clause vs. the review-cap sentence** (`context/AGENTS.md:16`) — "three rounds without
  convergence means stop and ask the user" is unconditional; the new "fix the rule, then re-dispatch"
  clause is a continuation of the existing SECOND-round sentence (per Proposal 3's own diff, appended
  to that sentence) and does not claim an exception to the three-round cap. A round-3 re-dispatch that
  still contradicts still hits "stop and ask the user." No contradiction found.
- **New Hard-rules sentence vs. BUG-029's row and the verifier's counts rule** — the appended sentence
  ("a check is validated once in the must-pass state and once in the must-fail state…") restates
  BUG-029's prevention rule consistently; the verifier's separate `-count=1` / pass-fail-skip counting
  rule (`agents/verifier.md:43`) is about test-run counts, an unrelated mechanism — no clash.
  BUG-029's own INDEX row is untouched by this diff (out of scope, file 7).
- **review-standards § 3's new item vs. its existing describing-text item** — not redundant: the
  existing item covers comments/doc-comments/test names/spec rows describing behaviour *the diff
  itself* changed; the new item covers rules/bans/thresholds/triggers restated *elsewhere in the
  repository* (`context/`, `skills/*/SKILL.md`, `agents/`, `README.md`), with named grep roots the
  older item lacks. Different scope, both earn their place.
- **Stale-text grep** (`at Done`, `EMPTY`, `headers only`, `finished lines`, `two questions`,
  `reference library`, `simple_expansion`) across `skills/`, `context/`, `agents/`, `README.md`,
  `skills/handoff/SKILL.md`, `skills/feature/SKILL.md` — only the one hit above. `two questions` hits
  are all in `.agents/docs/` history (already fixed per BUG-042's record, `agents/reviewer.md` no
  longer says "two questions" in the live file). `reference library` in `skills/golang/SKILL.md`'s
  description is consistent with the new reference-reading line, not stale. `simple_expansion` now
  appears in exactly one live skill file (the one this diff added it to) plus historical logs.
- **Frontmatter/parsing** — `go run ./cmd/sh-agent list` reports 27 skills, 7 agents, no parse error.
- **Placement** — golang's new bullet lands under "Structure & Style" (as Proposal 6 specifies);
  `references/golang-io-web.md` and `references/golang-testing.md` both exist. init-project's new
  clause lands in Step 2's headless sentence (as Proposal 7 specifies). changelog's new commits rule
  lands under `## Rules` (as Proposal 1 specifies). All match their proposal diffs verbatim.

## Criteria

1. Dead code — pass; no removed/orphaned code, these are prose-only additions.
2. Flow tracing — one gap: the new EMPTY-at-every-checkpoint rule is not traced into
   `skills/handoff/SKILL.md`'s own restatement of the same obligation (see Finding above).
3. Standards conformance — pass otherwise; each edit matches its proposal's target section and
   wording exactly.
4. Directory layout — n/a (no new files).
5. Simplicity — pass; all six edits are prose additions, no new dependency or abstraction.
6. Decoupling — pass; each edit is scoped to the file/section its proposal named, no cross-file
   collision among the six diffs themselves.

## Not examined

- `.agents/docs/bugs/INDEX.md` (the seventh proposal) — explicitly out of scope per the dispatch.
- Every other project's harness mirror (`crew`, `argolis-devbox`) — not part of this diff.
- Full-text read of `skills/feature/SKILL.md` beyond the grepped terms — read the phase/checkpoint
  sections relevant to the cross-checks above, not the whole file end to end.

## Outcome

- Finding 1 (medium) — fixed: `skills/handoff/SKILL.md` Step 3 item 3 now reads "both are empty (headers only) at every checkpoint, not just at Done — a ticked line is deleted in the commit that ticks it", matching `context/AGENTS.md` step 4. Found by the reviewer, not by Proposal 2's own grep — Proposal 5's rule caught its first case in the commit that introduces it.
- Note on the queue file — accepted: the eight `**Status:** applied` lines were written after the edits but before this review and the commit existed; the commit that carries this report makes them true, and the next retro writes the status after the review, not before.
- Docs lens (`reviews/2026-09-16-review-worktree-docs-retro.md`, REQUEST CHANGES) — its one finding fixed: the contradiction clause added to `skills/review/SKILL.md` Step 5.
- Round 1 closed 2026-09-16; no second round owed (two restatements, no mechanism finding).
