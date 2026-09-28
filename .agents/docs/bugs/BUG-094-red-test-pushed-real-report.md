---
owner: bugs
date: 2026-09-28
id: BUG-094
status: fixed
area: tests / cmd/sh-agent report path
---
# BUG-094: A red test run pushed a real report to the feedback repository

## Symptom
While fixing a review finding ("`report --from-zip x` still sends a report"), `TestFromZipRefusedOnOtherCommands` (`cmd/sh-agent/zip_test.go`) was written before the refusal existed and run to confirm it was red. With no refusal in place, `run([]string{"report", "--from-zip", "x.zip", "--home", <temp>})` executed the real `sh-agent report`: it cloned `styrman-harness-feedback` over the operator's SSH credential and pushed commit `c5d47b2` ("report: dave 2026-09-28-172114") at 2026-09-28T15:21:17Z.

## Repro
Run the test with the pre-fix `main.go` (no `--from-zip` refusal ahead of the command switch) and a real, unrestricted `PATH`/`HOME` — the command runs for real and pushes.

## Root cause
A red-phase test for a missing refusal or guard executes the real command it means to refuse, if nothing cuts that command's outward reach first. The test isolated only `--home` (scratch), not `PATH` or credential reachability; `~/.ssh` and `gh` auth were still reachable from the real home, so the real `report` ran to completion and pushed.

## Fix
The refusal now runs before any command (`cmd/sh-agent/main.go`, the `zipSet` block ahead of the first `switch`). The test now sets `PATH` to an empty temp directory so git/gh are unreachable even if the refusal regresses. A mutation run (refusal disabled, `PATH` empty) pushed nothing, confirmed via `gh api .../commits`. Every later run used `-timeout`.

Remediation of the already-pushed record (`reports/dave/2026-09-28-172114/report.json` in `styrman-harness-feedback`, commit `c5d47b2`) was operator-approved ("All approved", 2026-09-28). Content read back: metrics only (host "dave", linux/amd64, `harness_commit` 1cb2608, drift 257, 27 skills, 7 agents, `go1.27.1`), no files, no redactions, no secrets, no corrections.

The pushed report was removed from `styrman-harness-feedback` by an ordinary commit, `1703cea` ("Remove the report a styrman-harness test pushed by mistake (BUG-094)…"), no force-push.

## Prevention
A test for a command with an outward-facing side effect (report, download, push, deploy) cuts that reach first — empty `PATH`, no reachable credentials — before the command can run, and above all when the test is written ahead of the refusal or guard it checks: the red run executes the real thing. `--home` pointing at scratch does not isolate credentials that live elsewhere (`~/.ssh`, `gh` auth).

## Related
Shares its mechanism with BUG-071 ("a test that sets up an unusual condition proves the condition holds") and the Definition of Done step 1 rule that a verifier dispatch runs on a frozen tree — both are cases of a test/run not confirming its own isolation before trusting the result. Candidate home: `context/AGENTS.md` § Hard rules, near "A test that sets up an unusual condition proves the condition HOLDS before it is trusted" — this is the same principle applied to outward-facing side effects specifically, and could be folded into that same bullet rather than added as a new one.
