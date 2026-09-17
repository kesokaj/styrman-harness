---
owner: bugs
date: 2026-09-17
id: BUG-055
status: fixed
area: internal/harness (report.go push/checkout resolution)
---
# BUG-055: `sh-agent report` pushed to whatever remote its checkout resolved to, with no check that it was the feedback repository

## Symptom
Security review found nothing verified that the checkout `report` commits and pushes from is the feedback repository. `--feedback-repo` pointed at an unrelated project clone committed and pushed the report there; `GIT_DIR` set in the environment redirected the commit and push to another repository, both exiting 0.

## Repro
Probed against 015cf2d: `--feedback-repo <unrelated project clone>` exited 0, with `reports/probe/<ts>/feedback.jsonl` and `report.json` committed and pushed to that project's `origin main`. `GIT_DIR=<other repo>/.git` with a correct `--feedback-repo` also exited 0, committing and pushing into the other repository.

## Root cause
`internal/harness/report.go:363-386,562-570` (015cf2d) accepted any directory with a `.git` entry that wasn't itself a harness checkout as `--feedback-repo`, and the default `~/.config/sh-agent/feedback` got no check at all. `gitOut` also honoured `GIT_DIR`/`GIT_WORK_TREE` from the environment, so those redirected every git call regardless of the resolved path. Nothing compared the checkout's actual remote against the expected feedback URL before writing or pushing.

## Fix
`verifyFeedback` (`internal/harness/report.go`, commit 012b4ca) requires the checkout's own top level and `remote.origin.url` to exactly match the feedback URL, on both the explicit and default path; `SendReport` refuses `GIT_DIR`, `GIT_WORK_TREE`, `GIT_INDEX_FILE`, `GIT_OBJECT_DIRECTORY`, `GIT_COMMON_DIR` in the environment; the push is always the explicit `push -u origin HEAD`. `TestResolveFeedbackFailures` and `TestSendReportRefusesRedirectsAndLinks` cover both; must-fail seen.

## Prevention
A tool that commits and pushes to a repository resolved from a flag, default, or the environment must verify the resolved remote matches the expected target before writing anything — not just that a `.git` exists. Refuse (do not silently follow) redirect-shaped environment variables (`GIT_DIR` and siblings) in any tool that shells out to git.

## Related
BUG-053, BUG-054, BUG-056 (same review). "A tool that pushes verifies its remote before writing" is a candidate standing rule for `agents/security.md` or a `golang`/git-usage checklist — actionable (a specific remote-equality check) and generalizable to any future git-writing tool in this repo.
