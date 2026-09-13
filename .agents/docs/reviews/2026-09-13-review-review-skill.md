---
owner: reviewer
date: 2026-09-13
---
# Review — review-skill   (2026-09-13)

Range reviewed: 9ff33fe..working tree (uncommitted; `git diff 9ff33fe` + `git status --short`; new file `skills/review/SKILL.md` is untracked, blob not yet in git — reviewed from the working-tree file directly)

## Verdict: REQUEST CHANGES

The fold is otherwise clean and internally consistent (naming, dispatch table, artifact slugs, and the `/audit` repointing all check out against D1–D9), but the skill's own body text directly contradicts D8: it names Claude Code's built-in `/code-review` by name, which D8 explicitly forbids ("the skill never mentions a vendor's built-in"). A secondary, lower-severity gap: D7 claims `/review --project --all` makes "the same four dispatches the audit skill made," but the mechanics defined by D2/D3 (lenses combine, `--all` = every lens) actually produce six dispatches under that invocation (reviewer called three times — code, simplify, docs — plus security, verifier, explorer), not four.

## Findings

- **[medium / confidence: high]** `skills/review/SKILL.md:9` — "On Claude Code the built-in `/code-review` also exists; the harness path is this skill, so both CLIs behave alike." This names a vendor's built-in review command inside the skill body, directly contradicting D8's "the skill never mentions a vendor's built-in." (Note: this line is invisible to `git grep` over tracked files since `skills/review/` is untracked — confirm with `grep -n code-review skills/review/SKILL.md` before it is staged.) Fix: delete the sentence or replace it with vendor-neutral wording, e.g. "the harness path is this skill on both CLIs."

- **[low / confidence: medium]** `skills/review/SKILL.md` Step 3 table vs. plan D7 — D7 says `/review --project --all` makes "the same four dispatches the audit skill made" (explorer, reviewer, security, verifier — one reviewer call, per the deleted `skills/audit/SKILL.md`, whose single reviewer dispatch already covered all six `review-standards` criteria in one pass). But D2/D3, which the skill implements faithfully, define `--all` as every lens combining, and three lenses (`--code`, `--simplify`, `--docs`) all route to the `reviewer` agent (Step 3 table) — so `--project --all` actually issues six dispatches (three separate reviewer runs plus security, verifier, explorer), not four. This isn't a bug in the skill (it correctly implements D2/D3), but D7's description of equivalence with the old `/audit` is inaccurate and could surprise a user expecting audit-sized cost. Worth a one-line correction to D7 or a note in the skill that `--all` under `--project` is a superset of the old audit, not identical to it.

## Not examined

- Did not execute the skill end-to-end on either CLI (headless `/review` or `/review --project --all`); the plan's own e2e tiers (1)-(5) are unstarted work, not part of this diff's static content, and are out of scope for a text-only review.
- Did not review `.agents/docs/bugs/BUG-012-explorer-artifact-missing.md` or the `.agents/docs/reviews/2026-09-13-cold-read-*` files in depth — records under `.agents/docs/` were named out of scope by the dispatch.
- Did not check `agents/documentation.md`, `agents/changelog.md`, or `agents/bugs.md` beyond confirming `documentation` has edit tools for the `--fix --docs` path; they are not touched by this diff and are not named in the dispatch's contradiction list.
- Did not run `gofmt`/`go vet`/`go test` — no Go source is touched by this diff.

## Criteria

1. Dead code — pass. `skills/audit/SKILL.md` is fully removed via `git rm`; `git grep` (tracked files) and a direct file check confirm no leftover `/audit` references outside two named vendored exceptions; `skills/audit/` no longer exists on disk.
2. Flow tracing — pass except the D8 violation above. Cross-checked the skill's Step 3 dispatch table against `agents/reviewer.md`, `security.md`, `verifier.md`, `explorer.md`: artifact filename templates (`reviews/<date>-review-<slug>.md`, `security-<slug>.md`, `testing/<date>-verify-<slug>.md`, `plans/<date>-recon-<slug>.md`) all line up with the slugs the skill assigns (`<scope>-code`, `<scope>-simplify`, `<scope>-docs`, bare `<scope>` for security/verifier/explorer) with no filename collisions for any lens combination. `--fix` correctly relies on the main agent for edits since `reviewer`/`security` have no Edit tool (confirmed: both list only `Read, Grep, Glob, Bash, Write`). `documentation` (dispatched under `--fix --docs`) does have `Edit, Write`, so that step is executable as written.
3. Standards conformance — pass. Frontmatter values are quoted (no bare `#`, avoiding BUG-004); `argument-hint` matches D9 verbatim; skill count stays 27 per `sh-agent list` (confirmed by running it).
4. Directory layout — pass. New skill lands at `skills/review/SKILL.md`, the conventional location.
5. Simplicity — pass. No new dependency; the skill is a routing/text-composition document, matching the existing `preflight`/deleted-`audit` pattern in size and structure (64 lines, under the plan's 140-line ceiling).
6. Decoupling — pass. The skill computes scope text itself and hands it to unedited agents (D3's explicit design), so `reviewer`/`security`/`verifier`/`explorer` remain independently owned files; no agent file needed editing beyond the one-line description changes in `explorer.md` already covered under flow tracing.
