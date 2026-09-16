---
owner: bugs
date: 2026-09-15
id: BUG-029
status: fixed
area: testing (Phase 6 kit and staging scripts, `.agents/scratch/`)
---
# BUG-029: five checks and guards in scratch scripts matched the wrong state

## Symptom
Five separate defects in the 2026-09-15 harness practices work (Phases 6 and 7), each a check or
guard that fired (or stayed silent) on the wrong input:
(a) a control check exited 127 and was logged as a valid FAIL; (b) `pgrep -f 'phase6-kit/drive.sh'`
matched the grep's own shell; (c) a record-writer's needle `--vendor` matched two README lines and
wrote nothing; (d) an Outcome-section assert matched an inline substring; (e) a pre-commit guard on
`DOCS-REVIEW-VERDICT-MARKER`/`VERIFIER-STATUS-MARKER` blocked a commit over two review reports that
legitimately quote the placeholder names as findings.

## Repro
(a) base clone `.agents/scratch/basectl.nIaF`, control check run with a relative path, exit 127.
(b) `pgrep -f 'phase6-kit/drive.sh'` while the guard itself ran under that path.
(c) grep for `--vendor` against README.md, two unrelated matches.
(d) grep for `## Outcome` as a substring against a review report containing the string inline.
(e) `grep -rl 'DOCS-REVIEW-VERDICT-MARKER\|VERIFIER-STATUS-MARKER' .agents/docs` before commit
e31c3a2, matching `reviews/2026-09-15-review-c514936-preflight-docs.md` and
`…-preflight-advisory-docs.md`.

## Root cause
Each check or guard was exercised only in the state where it should fire, never in the state where
it must stay silent or pass, and a failing run's exit status/match reason was not read before being
trusted. (a) is a 127 (command not found from a relative path), not a check failure. (b) and (e) are
patterns that can match their own invocation or a legitimate quotation of the target string. (c) and
(d) are needles too generic for the corpus they searched.

## Fix
(a) absolute paths, log line corrected. (b) bracket pattern `'phase6-kit/driv[e]'`. (c) unique needle,
re-run. (d) line-anchored regex `^## Outcome`. (e) guard narrowed to `changes/` and `plans/`;
committed as e31c3a2.

## Prevention
Before a check or guard gates anything, run it once where it must pass and once where it must fail.
Read a failing control's exit status for its reason: 126, 127 and 2 (not executable, not found,
usage) are not a failing check. A pattern that searches process lists or files must not be able to
match its own invocation or a quotation of its target.

## Related
BUG-024 (a matrix check validated only against a must-fail state), BUG-010, BUG-013. This is the
third record with "validate a check in both states" as its prevention (BUG-024, BUG-013, here);
recommend it graduate into a `.agents/rules/` file or a testing-skill checklist item — actionable
and generalizable to any check/guard added in this repo.

Graduated 2026-09-16 to `context/AGENTS.md` § Hard rules, the "Silence from your own instrumentation" bullet.
