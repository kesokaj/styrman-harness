---
owner: main
date: 2026-10-09
id: BUG-168
status: fixed
area: cmd/sh-agent (vendor, update --from-zip, adopt, search and index), internal/search (validate)
---
# BUG-168: six usage errors exited 1 on stderr while the docs promised exit 2

## Symptom
README § Daily commands, `specs/sh-agent/delivery.md` and `specs/sh-agent/dream.md` say a usage error exits 2 with its reason on the verb's writer. These exited 1 with the message on stderr: `dream --vendor bogus` and `memory show --vendor bogus` (`unknown vendor`), `search` with no words, a bad `--limit` or `--kind` (`internal/search` `validate`), `adopt` with no name, `update --from-zip` with no file, and `--from-zip` given to a verb other than `update`. Found by the project review of 2026-10-09 (code lens § 7 vocabulary, low); shipped since each site was written.

## Repro
`TestUsageErrorsExit2` (`cmd/sh-agent/usage_test.go`): seven subtests, all red before the fix with `want exit 2, got <plain error>`.

## Root cause
The exit-2 convention was applied site by site: the flag parser, stray arguments, `memory`, `metrics` and dream's `ExitError` each print on the writer and return `exitCode(2)`, and no test covered the convention as a class, so each later site returned `errors.New(...)` and fell to `main`'s exit-1 path.

## Fix
- `cmd/sh-agent/main.go`: `(*ctx).vendor`, `memoryCmd`'s show branch, the three `--from-zip` refusals and `runAdopt` print the reason on the writer and return `exitCode(2)`; `runSearch` does the same for `search.IsUsage(err)`.
- `internal/search/search.go`: `validate` returns `usageError`; `IsUsage` carries it across the package boundary.
- Tests: `TestUsageErrorsExit2` (seven sites), `TestMainExitStatus` (main's 0, 1 and 2 through a helper process), `TestVersionVerb`; the four tests that pinned the plain error (`TestFromZip*`, `TestUpdateFromZipRefusals`, `TestSearchVerb`) now pin exit 2; the seventh site surfaced when the suite ran bare, by one of them.

## Prevention
A convention that spans verbs is pinned by one table test over every site, and a new usage refusal adds its row there. In `cmd/sh-agent`, a usage refusal is `fmt.Fprintln(c.out, ...)` then `return exitCode(2)`; an `errors.New` for a usage message is the review's finding.
