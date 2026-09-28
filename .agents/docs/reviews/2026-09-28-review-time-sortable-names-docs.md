---
owner: reviewer
date: 2026-09-28
---
# Review — time-sortable-names-docs   (2026-09-28)

Range reviewed: `17fa89b..e5a49d4`, limited to `-- context agents skills`

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Docs-lens only. The diff itself (context/AGENTS.md, agents/changelog.md, skills/handoff/SKILL.md, skills/init-project/SKILL.md, skills/preflight/SKILL.md line 86) is internally consistent and lands the three named rule changes correctly. But two restatements of the rules the scope changed were left standing elsewhere in the living docs, both still telling an agent to do the old thing: `skills/handoff/SKILL.md:62` still instructs "Run `date +%F`... use that date everywhere below" though Step 5 of the same file now requires `date +%F-%H%M` for the handoff filename, and `skills/preflight/SKILL.md:40` still tells the caller to dispatch the `changelog` agent "with today's date" though the agent's own contract (agents/changelog.md) and the other preflight call site (line 86, inside this diff) now require the `date +%F-%H%M` stamp.

## Findings

- **[medium / confidence: high]** `skills/handoff/SKILL.md:62` — "Run `date +%F` and use that date everywhere below — never a date from memory." This line was not touched by the diff, but Step 5 of the very same file (line 145, changed in this diff) now saves the handoff as `$DOCS_ROOT/handoffs/$(date +%F-%H%M)-<topic>.md`. Following Step 3's literal instruction produces a date-only filename, which `init-project`'s new counter (skills/init-project/SKILL.md:24) will flag as a date-only name on the very next audit. Fix: change line 62 to capture the stamp (`date +%F-%H%M`) and use its first ten characters where a bare date is needed (frontmatter, heading) and the whole stamp for the filename — mirroring what `agents/changelog.md:22` now says.

- **[medium / confidence: high]** `skills/preflight/SKILL.md:40` — "dispatch the ... `changelog` agent (with today's date; or write the entry yourself, ...)" — still the pre-diff wording. The diff updated the *other* changelog dispatch in the same file (line 86: "Dispatch the `changelog` agent with the stamp from `date +%F-%H%M`"), and `agents/changelog.md`'s own contract now expects a stamp, not a bare date (a bare date now produces an explicitly flagged date-only filename per agents/changelog.md:22). Line 40 is the non-release dispatch path (used on every preflight run with dirty source paths, not just at tag time), so it is the more common of the two call sites in this file and now disagrees with its sibling at line 86. Fix: reword line 40 to match line 86's "with the stamp from `date +%F-%H%M`".

## Criteria

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- `README.md`, `<docs-root>/specs/`, `<docs-root>/testing/guide.md`, and `<docs-root>/deploy/*.md` were grepped for the old naming patterns (`date +%F` alone, `YYYY-MM-DD-<slug>`, `changes/`, `handoffs/`) and returned no stale hits related to this scope; not read line-by-line beyond that, since the scope is limited to `context agents skills` and none of those trees appeared in the diff.
- `agents/bugs.md`, `agents/reviewer.md`, `agents/security.md`, `agents/verifier.md`, `agents/explorer.md` each still use bare `date +%F` for their own artifact naming (`reviews/`, `testing/`, `plans/recon-*`) — these are outside the renamed scope (only `handoffs/` and `changes/` were renamed) and are correct as written; not findings.
- The newest `changes/` entry (`.agents/docs/changes/2026-09-28-1428-time-sortable-names-phase-1.md`) was checked only for filename shape (confirms the new `YYYY-MM-DD-HHMM-<slug>.md` convention is already in live use) — its body content was not read in full.
- `ROADMAP.md` and the uncommitted plan file (`.agents/docs/plans/2026-09-28-time-sortable-trail-names.md`, shown modified in git status) are outside the `context agents skills` path filter and were not reviewed.
- Vendored skills (per `skills/VENDORED.md`: grilling, codebase-design, find-skills, anti-ui-slop, gke-basics, terraform) were excluded from the grep as instructed.

## Outcome

- Finding 1 (`skills/handoff/SKILL.md:62`): fixed. Step 3.1 now names `date +%F-%H%M` for the file name (Step 5) beside `date +%F` for dates. Plan D5.
- Finding 2 (`skills/preflight/SKILL.md:40`): fixed. The dirty-tree changelog dispatch passes the stamp from `date +%F-%H%M`, like `:86`. Plan D5.
- Re-check: `git grep -nE "today's date" -- skills agents context` leaves only dispatches to agents that write neither handoffs nor change entries (`agents/bugs.md:20`, `skills/preflight/SKILL.md:48`, `skills/review/SKILL.md:43`). These are out of scope and unchanged.
- The behavioural evidence in `testing/2026-09-28-verify-time-sortable-names-phase-2.md` predates both edits. The handoff probe wrote timed names 3 of 3 with the Step 3.1 text as it was, so the edit removes a contradiction rather than changing measured behaviour. The re-verification after the fixes is the verifier run at Done.
