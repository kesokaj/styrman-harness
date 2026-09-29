---
owner: reviewer
date: 2026-09-29
---
# Review — dream-docs   (2026-09-29)

Range reviewed: `git diff 80bf5b6` plus working-tree changes (uncommitted) — `skills/dream/SKILL.md` (new), `skills/retro/SKILL.md`, `context/AGENTS.md` (§ Learning loop, one sentence), `.agents/docs/testing/skill-matrix.md` (row `dream`), plus the modified plan and behavioural-test files shown in `git status`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The diff itself (the new `dream` skill, the `retro` Step 1/6 additions, the context sentence, the matrix row) is internally consistent and well cross-referenced. The problem is outside the diff: `dream` records now land in `<docs-root>/retros/` and in `retros.jsonl` as a third record kind, but four other living-doc restatements of "what's in `retros/`" / "what `retros.jsonl` holds" were not updated to include it — the same enumeration gap repeated in four places, which is exactly the failure mode the global context's describing-text rule calls out (a vocabulary gets a new member; every enumeration of it must move with it). None of these break anything functionally — `/dream` still writes where its own SKILL.md says — but a reader of `REQUIREMENTS.md`, `README.md`, `skills/handoff/SKILL.md`, or `skills/feedback-log/SKILL.md` alone would not learn that `dream` records exist.

## Findings

- **[medium / confidence: high]** `REQUIREMENTS.md:31` — "its retro and distill records in `.agents/docs/retros/`" still enumerates only two record kinds. `skills/dream/SKILL.md` Step 6 now writes a third: `<docs-root>/retros/<date>-dream.md`. This is the file the task brief names explicitly as the anchor for this check, and it is a project rule (REQUIREMENTS.md is the non-negotiables file) whose own text is now incomplete. Fix requires the operator's explicit approval to edit REQUIREMENTS.md (project AGENTS.md rule) — flagging it here rather than editing it.
- **[medium / confidence: high]** `README.md:210` — the "What gets installed where" table describes `retros.jsonl` as "(cross-project retro/distill index)". `skills/dream/SKILL.md` Step 6 appends `{"kind":"dream", ...}` lines to the same file (also from `/retro`'s Step 6 addition, when dream verdicts are applied through retro). The parenthetical undercounts the kinds the file now holds.
- **[medium / confidence: high]** `skills/handoff/SKILL.md:51` — the artifact table row "Retro and distill records (project scope) | `.agents/docs/retros/` |" omits the third kind of file now written to that directory by `/dream` (standalone dream records, per its Step 6). A next-session agent reading this table for what's in `retros/` would not expect a `-dream.md` file there.
- **[low / confidence: medium]** `skills/feedback-log/SKILL.md:91` — "`/distill` records there too" (after describing `/retro`'s record in `<docs-root>/retros/` and `retros.jsonl`) is the same enumeration, one clause short of naming `/dream`'s standalone record and its `retros.jsonl` line. Lower confidence/severity than the others because this file's own job (logging corrections) is unaffected; it's a secondary mention of the same fact already flagged above.

No other stale claim, broken interface reference, or undocumented command/flag was found. `skills/retro/SKILL.md` Step 1 item 7 and the `## Memory` section, `context/AGENTS.md`'s new sentence, and the `skill-matrix.md` `dream` row are all consistent with `skills/dream/SKILL.md`'s own text (store resolution, verdict vocabulary, `promote` flag semantics, `--dry-run` behaviour, record locations). The plan file's D7–D13 decisions and the Phase 1 behavioural-test log line up with what `skills/dream/SKILL.md` and `skills/retro/SKILL.md` actually say.

One claim I could not verify: `skills/dream/SKILL.md`'s Step 1/1 sourcing and the plan's D7 citation of `code.claude.com/docs/en/memory.md` § Storage location / § How it works (the settings-scope precedence, the 200-line/25KB load limit, the `type` frontmatter field) — no network access in this review, so the quoted external-doc text is taken on the plan's word, not independently checked.

## Criteria

1. Dead code — n/a, docs lens.
2. Flow tracing — n/a, docs lens (see Findings for the one flow that isn't fully traced through every restatement).
3. Standards conformance — fail on the global context's describing-text/vocabulary-enumeration rule, four instances, see Findings; the diff itself follows the `review` skill's docs-scope conventions (frontmatter, docs-root resolution, dated filenames).
4. Directory layout — n/a, docs lens; new file `skills/dream/SKILL.md` lands where every other skill lives.
5. Simplicity (DIY-except-security) — n/a, docs lens.
6. Decoupling — n/a, docs lens; the `dream` → `retro` handoff is a clean seam (dream never writes a rule; retro Step 6 is the one place applying promoted memories), described consistently in both skill files.

## Not examined

- Network verification of the `code.claude.com/docs/en/memory.md` citations (D7–D9 in the plan) — no internet access from this review; noted above as unverified.
- `.agents/docs/specs/sh-agent/delivery.md` — read via grep only (no hits for retro/distill/dream/memory); confirmed out of scope for this change (it documents the release/delivery model, not the learning loop).
- `.agents/docs/testing/guide.md` — checked via grep, no hits; this project has no single testing guide describing the learning loop, so no finding either way.
- The generated, installed copy of the context (`~/.claude/AGENTS.md`, `~/.claude/skills/`) — out of scope; `dream` has not yet been installed into the real home (the plan's Phase 2 real-home task is explicitly deferred to the closing handoff), so its absence from the installed skill list there is expected, not a docs defect in this repo.
- Did not re-verify the behavioural test logs' run-by-run claims (cost figures, transcript counts) against raw transcripts — took the `.agents/docs/testing/2026-09-29-dream-phase-1-behavioural.md` file at face value as a docs artifact, not re-run.

## Outcome
- `REQUIREMENTS.md:31` omits dream records — fixed: "its retro, distill and dream records", the exact edit put to the operator and approved on 2026-09-29 ("Approve the edit (Recommended)"); ADR-042 in `architecture/DECISIONS.md`.
- `README.md:210` `retros.jsonl` described as a retro/distill index — fixed by the `documentation` agent: "(cross-project retro, distill and dream index)"; its grep found one more restatement, `context/AGENTS.md` § Learning loop ("every retro and distill leaves a record"), fixed by the main agent to "every retro, distill and dream", size test re-run before the commit.
- `skills/handoff/SKILL.md:51` artifact table — fixed: "Retro, distill and dream records (project scope)".
- `skills/feedback-log/SKILL.md:91` — fixed: "`/distill` and `/dream` records there too".
