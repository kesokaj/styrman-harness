---
owner: reviewer
date: 2026-09-28
---
# Review — time-sortable-names-code (2026-09-28)

Range reviewed: `17fa89b..e5a49d4`, limited to `-- context agents skills` (5 files, 114 diff lines: `agents/changelog.md`, `context/AGENTS.md`, `skills/handoff/SKILL.md`, `skills/init-project/SKILL.md`, `skills/preflight/SKILL.md`)

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The change (switching `handoffs/` and `changes/` filenames from `YYYY-MM-DD-<slug>.md` to `YYYY-MM-DD-HHMM-<slug>.md`) is small, self-consistent across the files it touches, and verified: `TestRenderedContextFitsAntigravityBlock` still passes (20498/22000 bytes), a scratch `sh-agent install`/`status` cycle shows 0 drift, and the new date-only-name counting regex in `init-project/SKILL.md` was hand-tested and counts correctly. One real cross-reference was missed: `skills/handoff/SKILL.md` Step 3 still tells the agent to capture a single `date +%F` value and "use that date everywhere below," but Step 5 (changed in this diff) now needs a `date +%F-%H%M` stamp — the instruction wasn't updated to match its own consequence.

## Findings

- **[low / confidence: medium]** `skills/handoff/SKILL.md:62` vs `:145` — Step 3.1 says "Run `date +%F` and use that date everywhere below — never a date from memory," but Step 5's save command (changed in this diff, e5a49d4) now computes `$(date +%F-%H%M)`, a different value than what Step 3.1 told the agent to capture and reuse. The frontmatter (`:74`) and heading (`:76`) templates still use plain `YYYY-MM-DD`, so as written there are now two different "the date" values in the same document with no instruction telling the agent which one to capture when, or to derive the filename stamp from the same captured moment rather than a fresh `date` call that could tick over to a new day relative to the frontmatter/heading date in a long-running session. This is the same class of gap the project's own describing-text rule (`context/AGENTS.md` § Hard rules, and review-standards §3) is meant to catch: a rule changed in one place (Step 5) whose sibling restatement (Step 3.1) wasn't updated with it. Fix: either have Step 3.1 capture the full `date +%F-%H%M` stamp up front and Step 5 reuse it, or explicitly say Step 5 recomputes a fresh timestamp for the filename only.

- **[none — verified fine]** `skills/init-project/SKILL.md:24` — the new date-only-name counting pipeline (`grep -E '^[0-9]{4}-[0-9]{2}-[0-9]{2}-' | grep -vcE '^[0-9]{4}-[0-9]{2}-[0-9]{2}-[0-9]{4}-'`) was hand-tested against a mix of date-only and timed filenames (including a slug that itself starts with 4 digits) and counts correctly (2 of 4 files, the two true date-only names). No finding.

## Criteria

1. **Dead code** — pass. No functions/files orphaned; this is a docs/skill-prose change only.
2. **Flow tracing** — one gap found (the finding above): the handoff skill's own internal cross-reference between "capture the date once" (Step 3) and "the filename now needs a timestamp" (Step 5) wasn't reconciled. Everywhere else the stamp threads correctly: `agents/changelog.md` description → dispatch instruction in `skills/preflight/SKILL.md:86` both say "the stamp from `date +%F-%H%M`"; `agents/changelog.md`'s own fallback logic (no Bash tool available to the agent, confirmed via its unchanged `tools:` line) was correctly adjusted to drop the old "run `date +%F` yourself" fallback (which was never actually executable by this agent) in favor of filesystem-context inference or `DATE-UNKNOWN`.
3. **Standards conformance** — pass except the Finding above. Grepped `YYYY-MM-DD` and `date +%F` across `context/`, `agents/`, `skills/`, `commands/`, `rules/`, `README.md`: all other restatements (reviews/, plans/, testing/, retros/ naming; `date +%F` used for other artifact types) are legitimately out of this change's scope (the commit history / plan names only `handoffs/` and `changes/` as the phase-1/phase-2 target), and none of them claim the old `changes/`/`handoffs/` format.
4. **Directory layout** — n/a, no new files; edits land in the same files that already owned this convention.
5. **Simplicity / DIY-except-security** — pass. Pure shell `date`/`grep`, no new dependency, no crypto/auth surface touched.
6. **Decoupling** — pass. Each file edited independently states the convention for its own agent/skill; no shared helper needed, consistent with this repo's "file drop, never code" rule for skill/agent changes.

## Not examined

- `.agents/docs/plans/2026-09-28-time-sortable-trail-names.md` and `ROADMAP.md` — modified in this range per `git status` but excluded by the dispatch's path filter (`-- context agents skills`); not reviewed here.
- `.agents/docs/testing/2026-09-28-verify-time-sortable-names-phase-2.md` — untracked, also outside the path filter; not reviewed.
- Behavioural verification (an actual dispatch of `changelog` or `/handoff` exercising the new stamp end-to-end) was not run by this reviewer — only the static file-drop consistency and the `install`/`status` drift check were exercised, per the tool-call budget for this dispatch.

## Outcome

- Finding 1 (`skills/handoff/SKILL.md:62` vs `:145`): fixed. Step 3.1 now reads "Run `date +%F` and use that date everywhere below, and `date +%F-%H%M` for the handoff's file name (Step 5)". It was fixed together with the docs-lens finding 1 (`reviews/2026-09-28-review-time-sortable-names-docs.md`), plan D5.
