---
owner: reviewer
date: 2026-09-16
---
# Review — bug-037-docs   (2026-09-16)

Range reviewed: working tree, `git diff HEAD -- internal/harness/adopt.go internal/harness/harness_test.go` (`--docs` lens per `skills/review/SKILL.md`)

## Verdict: APPROVE

The diff is a bug fix (BUG-037) that tightens `Adopt`'s path handling: it absolutizes `from`, validates the resolved name is a single path segment, and refuses a source/destination overlap instead of letting `os.RemoveAll` delete the checkout, `skills/`, or the source itself. No living doc asserts the old, buggy behaviour, and none of the four named mentions of `adopt` (README.md:29, README.md:80, `skills/distill/SKILL.md:137`, `skills/find-skills/SKILL.md:106`) describe `.`, `..`, or self-referential paths as inputs — all four describe the ordinary by-name/by-path adoption of a globally-installed skill, which this diff leaves working (`TestAdopt` unchanged and still green). No statement went stale and no removed rule needs a new home.

## Findings

None that meet the bar (a docs statement that no longer matches the code, or an undocumented changed interface). One observation, below severity for a finding:

- **[low / confidence: low]** `README.md:29` — the adopt row documents the happy path and the pre-existing "already exists; pass --force" refusal, but not the new overlap refusal (`"%s and %s overlap; adopt copies a skill directory from outside the repo's skills/"`). This is a safety refusal, not a new feature, and no doc contradicts it, so it does not meet the finding bar on its own; flagging only in case a future user hits it and looks for it in the daily-commands table.

## Docs cross-check (grep for old wording)

- `README.md:29` — "`sh-agent adopt <name|path>` ... `--force` replaces an existing skill of the same name." Still true; `--force` on a genuine name/dest mismatch is untouched by this diff.
- `README.md:80` — "A skill you already installed globally with `npx skills add -g` can be pulled in with `sh-agent adopt <name>`." Still true; unaffected (source is under `home`, dest under `repo/skills`, never overlapping in this flow).
- `skills/distill/SKILL.md:137` — adopts from a *different* project's `.agents/skills/<name>` into the harness repo; source and dest are in different repositories, so the new overlap check never fires here. Still accurate.
- `skills/find-skills/SKILL.md:106` — same pattern, global skill by name. Still accurate.
- `context/AGENTS.md:37` — "`sh-agent adopt` promotes one later" (a project-local skill into the harness). Unaffected; still accurate.
- `.agents/docs/testing/sh-agent/guide.md` and `.agents/docs/deploy/*.md` — no mention of `adopt` at all, before or after this diff; not a regression introduced by this scope.
- `.agents/docs/bugs/BUG-037-adopt-force-deletes-checkout.md` and `INDEX.md` — excluded from the grep per the lens (bug records), but noted for the caller: the bug's own "Fix" section names a further remediation step ("copy into a temp directory under `skills/` and rename it over `dest` only after the copy succeeds") that this diff does not implement — `os.RemoveAll(dest)` still runs before `copyTree`. `INDEX.md:43` still shows BUG-037 `open`, which is consistent with the diff being a partial fix, not a doc/behavior mismatch. This is a code-lens/completeness question, not a docs finding; flagging so it isn't lost.

## Criteria

1. Dead code — n/a, docs lens.
2. Flow tracing — n/a, docs lens.
3. Standards conformance — n/a, docs lens.
4. Directory layout — n/a, docs lens.
5. Simplicity (DIY-except-security) — n/a, docs lens.
6. Decoupling — n/a, docs lens.

## Not examined

- The code-correctness question of whether BUG-037's fix is complete (temp-dir-and-rename not implemented) — out of scope for the docs lens; the underlying commit is not closed (INDEX.md still lists it `open`), consistent with the docs.
- `.agents/docs/plans/2026-09-16-recon-diagnose-adopt-force-deletion.md` (untracked, in git status) — read for context on root cause, not treated as a living doc to reconcile against.
- Any downstream mirror repos (styrman-crew, argolis-devbox) — outside this repo's scope.

## Outcome

- **No finding.** The low note, that `README.md`'s `adopt` row does not mention the new overlap refusal, is `left`: the row states what the command does, nothing in it contradicts the code, and the README does not enumerate error paths for the other commands either.
