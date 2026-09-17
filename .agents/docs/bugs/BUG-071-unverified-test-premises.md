---
owner: bugs
date: 2026-09-17
id: BUG-071
status: fixed
area: internal/harness/report_test.go (Phase 3 test premises)
---
# BUG-071: two Phase 3 tests set up conditions that did not hold, so they passed or failed for the wrong reason

## Symptom
Two test cases written during Phase 3 did not exercise what they claimed.
(a) A case meant to prove that a still-running report's clone is not swept
used `report-<os.Getpid()>` as the "live, unrelated" directory — exactly the
path the run under test clones into — so the assertion failed for a reason
unrelated to the sweep behaviour. (b) A case meant to prove the
unpushed-commit count fails safe made `.git/refs/remotes` unreadable, but git
read `.git/packed-refs` instead and answered normally, so the clone was
removed and the test failed; the fail-safe condition was only reproduced by
making `.git/packed-refs` itself unreadable (which makes git exit 128).

## Repro
`.agents/docs/reviews/2026-09-17-security-73f2c36-83f3cf6.md` § Outcome,
finding 1: "Measured while fixing: ... with `.git/packed-refs` unreadable git exits 128
(kept); with only `.git/refs/remotes` unreadable git still reads packed
refs and answers — the first shape I tried, which is why the test pins the
`packed-refs` one." `TestSendReportSweepsInterruptedClones` was fixed to use
pid 1 instead of the run's own pid.

## Root cause
A test premise ("this makes the refs unreadable", "this pid is alive and
unrelated to the run") was assumed rather than measured against the real
tool (git) and the real runtime (the test process's own pid) before being
built into the test.

## Fix
Not a source fix — a test fix, part of commit `dab2565`'s test additions.
(a) `TestSendReportSweepsInterruptedClones` uses a reaped pid and pid 1
instead of `os.Getpid()`. (b) `TestRemoveKeptClone`'s uncountable-clone
subtest makes `.git/packed-refs` unreadable instead of
`.git/refs/remotes`. Both were caught because the tests were run and read
before the fix, not after.

## Prevention
Before trusting a test that sets up an unusual condition, prove the
condition actually holds — run the command by hand in a scratch fixture and
read its output — because a premise that silently does not hold produces a
test that passes or fails for the wrong reason. This is BUG-024's rule
("validate a check against a must-pass state, not only a must-fail one"),
applied to test setup rather than to the assertion; it did not generalize on
its own and needed re-recording here for the setup case specifically.

## Related
BUG-024 (same rule, applied to the assertion side). Third occurrence of
"a premise or check was assumed rather than measured" in this project
(BUG-018/024 lineage) — if a fourth appears, recommend a project rule in
`AGENTS.md` or a `.agents/rules/` file requiring a scratch-measured probe
before any test that relies on an OS- or tool-level side condition (file
permissions, process liveness, upstream resolution); actionable (a review
checklist item: "does this test's setup print evidence it produced the
condition it claims?") and generalizable beyond this file. Changes:
`.agents/docs/changes/2026-09-17-release-binary-phase-3.md`.
