---
owner: bugs
date: 2026-09-17
id: BUG-051
status: fixed
area: internal/harness (report_test.go)
---
# BUG-051: `git init --bare` failed with "GIT_WORK_TREE not allowed without GIT_DIR" inside an isolated test

## Symptom
`TestSendReport` failed on its first run at `git init --bare -q -b main`:
`fatal: GIT_WORK_TREE (or --work-tree=<directory>) not allowed without specifying GIT_DIR`.

## Repro
A test helper (`gitEnv` in `internal/harness/report_test.go`) isolated git from
the operator's config with `t.Setenv("GIT_DIR", "")` and
`t.Setenv("GIT_WORK_TREE", "")`, then ran `git init --bare` in a subprocess.

## Root cause
Git treats an environment variable set to the empty string as SET, not unset.
Blanking `GIT_WORK_TREE` this way left it "set to empty", which git rejects
without a paired `GIT_DIR`. `t.Setenv(k, "")` only registers the value for
`t.Cleanup` to restore; it does not remove the variable from the process
environment during the test.

## Fix
`internal/harness/report_test.go` `gitEnv`: for `GIT_DIR`, `GIT_WORK_TREE`,
`GIT_INDEX_FILE`, call `t.Setenv(k, "")` (to register the restore) followed by
`os.Unsetenv(k)`. Confirmed the pre-existing use at
`internal/harness/harness_test.go:650-653` (`TestBehindAndPull` end-of-test
cleanup) already follows this same pattern and does not carry the defect.

## Prevention
Before isolating git (or any tool with "empty means set" semantics) from the
environment in a test, unset variables with `os.Unsetenv` after
`t.Setenv(k, "")`, not `t.Setenv(k, "")` alone; a blank string is not an unset
variable.

## Related
BUG-052 (same test, same first run).
