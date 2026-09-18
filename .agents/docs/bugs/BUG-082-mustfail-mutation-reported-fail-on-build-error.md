---
owner: bugs
date: 2026-09-18
id: BUG-082
status: fixed
area: tooling (must-fail mutation harness, audit scripts)
---
# BUG-082: four must-fail mutations printed "fails as it must" when `go test` had failed at the build, so no test ever ran

## Symptom
Four must-fail mutations printed "fails as it must" while `go test` had actually failed at compilation ("declared and not used"), so no test ran at all. Found in `.agents/docs/changes/2026-09-18-sh-agent-audit.md` § Tests and `.agents/docs/reviews/2026-09-18-review-160aa16-3209cf3.md` § Outcome.

## Repro
A mutation written by `.agents/scratch/*-mustfail/mutate.py` removed a use of a variable, producing a Go build error; `go test` exited 1 from the compiler, and the mutation-testing harness treated any non-zero exit as the expected behavioural failure.

## Root cause
The harness keyed its must-fail verdict on `go test`'s exit status alone, not on whether the mutated code compiled and a specific named test then failed. A build failure and a test failure share the same exit code, so the check could not tell "the mechanism broke as expected" from "the mutation didn't compile".

## Fix
The harness now reads for a `FAIL: <test name>` line before accepting a must-fail verdict, and reports a build failure explicitly as `[build failed]` rather than as a pass of the mutation check; the four mutations were retargeted so they compile.

## Prevention
`context/AGENTS.md` § Hard rules, the mutation clause already states: a must-fail mutation "counts as a verdict only when the mutated code COMPILED and a named test failed — a build failure … is printed as such and is no evidence". This rule was added as part of the session that produced this bug and covers it going forward.

## Related
`BUG-068` (must-fail mutation hygiene, uniqueness and restore-by-copy); retro `2026-09-18-retro.md` finding 4, Proposal 4.
