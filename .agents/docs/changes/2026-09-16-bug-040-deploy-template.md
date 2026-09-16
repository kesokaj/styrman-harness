---
owner: changelog
date: 2026-09-16
---
# BUG-040: init-project's deploy/local.md template shipped invented facts for a nonexistent service     (2026-09-16)

## Summary

Antigravity's `/init-project` committed a `deploy/local.md` with `PORT 8080` under Configure and `curl -i localhost:8080/links` under Verify for a service that did not yet exist, and with no Test section. The template's example values (`PORT`, `8080`, `healthz`) read as facts rather than placeholders because the template had no signal for what was fillable. Claude Code's run blanked the examples unprompted and added a Test section. Fixed: Build/Run/Test blocks in the template are now `<placeholder>` format; Configure row's instruction moved to an HTML comment above the table ("delete when no settings"); Verify carries no example command. SKILL.md Step 4 now states that a placeholder the answers do not fill is deleted, never replaced with an invented value.

## Type

bug fix

## Files changed

- `skills/init-project/SKILL.md` — Step 4 sentence clarified: placeholders the setup answers do not fill are deleted, not replaced with invented values.
- `skills/init-project/assets/agents-tree/docs/deploy/local.md` — Build/Run/Test blocks are `<placeholder>` format; Configure table has instruction in HTML comment and row is three placeholders only; Verify section has no example command.

## Agents dispatched

reviewer (`reviews/2026-09-16-review-bug-g2-deploy-template.md` APPROVE WITH NITS; one finding fixed: instruction moved from Configure cell to HTML comment above table), verifier (shared log: install succeeds with 260 files managed, drift none), bugs, changelog.

## Verification

Verifier log `.agents/docs/testing/2026-09-16-verify-linkshelf-sim-fixes.md`: install into scratch home succeeded (260 written, 0 unchanged, 0 removed, 0 kept). Template consistency with SKILL.md instruction confirmed (Markdown rendered, no broken syntax).

Owed but not done in this commit: per-CLI `/init-project` re-run at the fixed tree for both CLIs, recorded later in testing logs.

## Docs

`skills/init-project/SKILL.md` Step 4 updated; template `skills/init-project/assets/agents-tree/docs/deploy/local.md` updated to use placeholders consistently with the revised skill instruction.

## Requirements

n/a
