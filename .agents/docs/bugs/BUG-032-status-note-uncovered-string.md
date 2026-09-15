---
owner: bugs
date: 2026-09-15
id: BUG-032
status: fixed
area: cmd/sh-agent (tests)
---
# BUG-032: a changed status message shipped without a covering test, behind a coverage claim that named one test for two strings

## Symptom
The Definition of Done verifier over `c192622..a33978d` (`.agents/docs/testing/2026-09-15-verify-harness-practices-done.md`, Status FAIL) found no test for the `status` note printed at `cmd/sh-agent/main.go:309` when `~/.claude/CLAUDE.md` is hand-written. Commit c514936 reworded that note, together with the install keep-line at `internal/harness/apply.go:37`, so both tell the operator to move the file's content into a rule file under `rules/`, and its commit message says this landed "with a test". That test (`internal/harness/harness_test.go:257`, `TestInstallStatusResyncUninstall`) asserts only `apply.go`'s keep-line, not the `status` note.

## Repro
At a33978d: `grep -n 'rule file under rules' cmd/sh-agent/*_test.go` prints nothing, while `grep -n 'rule file under rules' cmd/sh-agent/main.go internal/harness/apply.go internal/harness/harness_test.go` prints the two production strings and one assertion covering only the second.

## Root cause
One fix edited the same guidance string in two packages (`cmd/sh-agent` and `internal/harness`). Coverage was checked for one of the two edited files and then stated for the change as a whole. Nobody mapped each changed string to its own assertion, so the reworded `status` note shipped untested.

## Fix
`TestStatusNamesForeignClaudeMD` added to `cmd/sh-agent/main_test.go`: installs into a temp home that already holds a hand-written `CLAUDE.md`, runs `status`, and asserts the note's text. Confirmed failable: it passed as written, failed (exit 1) with the note's string swapped out in `main.go`, and passed again after `main.go` was restored. `gofmt -l .` clean, `go vet ./...` exit 0, `go test -race ./...` exit 0.

## Prevention
When one change edits the same message or rule in more than one file, each edited file gets its own assertion, and the coverage claim names the covering test per file plus what is left uncovered — the global hard rule "Claims of test COVERAGE are the sharpest case: name what is not covered." A reviewer checks a per-file coverage map, not one failable test standing for the whole diff.

## Related
BUG-027 (describing text checked against old wording, not the rule's subject) and the global hard rule that a claim that something works states how it was verified and on which path. This is a third record sharing the "coverage claim must be checked per-file/per-path, not generalised from one passing check" theme (with BUG-013's exit-status-before-output and BUG-027's subject-not-wording grep); it is actionable (a per-file assertion map is a concrete check) and generalizable beyond this file. Recommend it graduate into the `bugs`/`reviewer` review checklist or `context/AGENTS.md`'s coverage-claim rule, next to the existing "name what is not covered" line.
