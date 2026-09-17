---
owner: bugs
date: 2026-09-17
id: BUG-060
status: fixed
area: internal/harness (report.go push target claim; docs/decision text)
---
# BUG-060: code comments and docs claimed a push "cannot be redirected" without naming the redirect mechanisms actually checked

## Symptom
`report.go`'s `FeedbackURL` comment, `commitAndPush`, `.agents/docs/deploy/local.md`, and `.agents/docs/specs/sh-agent/delivery.md` all stated the push "never goes to a remote" and "no pushurl or second url can redirect it." Security review (third round on this mechanism) found `url.<base>.insteadOf`/`pushInsteadOf` and a remote named after the URL string still redirect a URL-argument `git push`, in six probed variants (checkout-local config and via `GIT_CONFIG_*` env vars).

## Repro
`.agents/docs/reviews/2026-09-17-security-012b4ca-worktree.md` finding 1: with `url.<E>.pushInsteadOf=<O>`, `url.<E>.insteadOf=<O>`, `remote.<O>.url=<E>`, or `remote.<O>.pushurl=<E>` set (checkout-local or via `GIT_CONFIG_COUNT/KEY_0/VALUE_0`), `SendReport` returned `err=nil` while the feedback repository received 0 commits and the other repository received 1.

## Root cause
Two prior rounds narrowed the push-target hardening to the redirect shapes each round's finding named (`origin`'s own pushurl, a second `origin` url, an env-var echo of `origin`), and each time the surrounding prose was written as an absolute guarantee ("no pushurl or second url can redirect it") rather than scoped to the mechanisms actually checked. `url.*.insteadOf` and a remote named after the literal URL are separate git redirect mechanisms neither round tested.

## Fix
Not a code fix: per D14, the operator decided the push trusts the operator's own git configuration (this is a single-operator tool; `.git/config` never comes from the remote). `FeedbackURL`'s and `commitAndPush`'s comments, `local.md`, and `delivery.md` were narrowed to state precisely what holds (origin's own pushurl and extra urls cannot redirect) and what doesn't (URL rewrites and a remote named after the URL redirect it, as for any `git push`); D13 (3) corrected by D14 in `.agents/docs/plans/2026-09-17-delivery-rework.md`.

## Prevention
A claim that a destination "cannot be redirected" or "never goes to X" names the specific redirect mechanisms it was checked against, not an unscoped guarantee — the next round otherwise finds the mechanism nobody tried. When a review finding lands a second consecutive round on the same mechanism (here: the push target, three rounds running), stop patching and take it to the operator for what it must guarantee, per the project's own review-convergence rule.

## Related
BUG-026 (two consecutive rounds on one mechanism → stop and ask) is the harness rule this defect should have triggered earlier; this is the third occurrence of an absolute-sounding safety claim in prose outliving the mechanism it was verified against (see BUG-057, BUG-059 in this same session). Recommend the "name the mechanisms checked" clause be added explicitly to `context/AGENTS.md`'s "A behaviour change includes the text that DESCRIBES it" hard rule, or to `agents/security.md`'s remediation-writing guidance — actionable (grep for absolute words like "never"/"cannot" in a security-touching comment or doc and require a scope clause) and generalizable beyond this file.
