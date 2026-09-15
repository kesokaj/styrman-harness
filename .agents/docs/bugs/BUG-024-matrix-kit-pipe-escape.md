---
owner: bugs
date: 2026-09-15
id: BUG-024
status: fixed
area: testing (skill-matrix kit)
---
# BUG-024: `explorer` and `init-deploy` rows failed the matrix although the run did what the row asks

## Symptom
Phase 6 matrix run of `plans/2026-09-15-harness-practices-g01-g26.md`: `explorer` row on Claude Code
failed although the explorer was dispatched, the output named `greet` and `version`, and no tracked
file changed. `init-deploy` failed on Claude Code too.

## Repro
Build a kit from `.agents/docs/testing/skill-matrix.md` by taking check cells verbatim. The explorer
check becomes `grep -qiE "greet\|version" "$OUT" && git diff --quiet` — in extended regex `\|` is a
literal pipe, so it never matches. Init-deploy's check passes `\|` as a literal argument to
`git log --oneline`, not a shell pipe to `grep`.

## Root cause
A Markdown table writes `|` inside a code span as `\|` (`skill-matrix.md:34`, `:22`); the kit-building
step kept the escape instead of stripping it before running the check as a real shell command. The
per-run must-fail control only requires the check to fail before the run, so it cannot catch a check
that fails on every input, including a passing one.

## Fix
Kit row data unescaped for preflight, distill, init-deploy and explorer rows before use.
`.agents/docs/testing/skill-matrix.md:15` (Procedure step 3) now states the kit unescapes `\|`. Finished
runs were re-checked with the corrected check; the Phase 6 e2e log records each result.

## Prevention
A kit built from a Markdown table unescapes `\|` in code spans before running it as shell. Before
trusting a check's FAIL, run the check once against a state where it must pass — a must-fail control
alone cannot catch a check that fails on everything.

## Related
`.agents/docs/testing/skill-matrix.md` Procedure step 3.
