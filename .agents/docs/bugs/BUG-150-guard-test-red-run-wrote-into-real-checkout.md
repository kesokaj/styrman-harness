---
owner: bugs
date: 2026-10-01
id: BUG-150
status: fixed
area: internal/harness/catalog_test.go (TestRefreshCatalogScript), skills/google-catalog/scripts/refresh-catalog.sh
---
# BUG-150: the red run of a guard test wrote into the real checkout

## Symptom
Writing a test for a guard that did not exist yet (the refresh script's `--work` must not be the checkout), the main agent ran the new test to see it fail. The case passed `--work <the real checkout>` with working fixture sources. Nothing refused it, so the script ran to the end: it cloned the two fixture repositories into the checkout's root (`agents-cli/`, `google-skills/`) and left `build/`, `origins.txt` and `writer.log` there. Five untracked paths; nothing tracked changed, nothing was committed. Found at once in the test's output ("got err <nil>") and by `git status`.

## Repro
`go test ./internal/harness -run TestRefreshCatalogScript` on the test as first written, before the guard was in the script.

## Root cause
The red-first step RUNS the unguarded path by design, and the fixture aimed that path at a real place with sources that work. The rule from the review one day earlier (R3-4: guard cases use sources that cannot be cloned, so a regressed guard stops at the clone) was applied to the `--dest` cases and not carried to the new `--work` cases.

## Fix
The five paths were checked against `git ls-files` (0 tracked), looked at (the clones' remotes pointed at the test's temporary directory; the files were stamped the minute of the run) and removed by name; the embedded directories were checked clean. The three forbidden-place cases for `--work` now run a COPY of the script inside a stand-in checkout under the test's temporary directory, with sources that cannot be cloned, and assert that nothing was written into the stand-in or the fixture home. The `/` case is skipped as root.

## Prevention
A test for a guard is written so that the guard's ABSENCE is harmless, because red-first runs it without the guard: the forbidden place is a stand-in under the test's temporary directory (copy the script or point it at a fake root) and the sources cannot run; never the real checkout, home or `/` with working inputs. After any red run of such a test, `git status --short --untracked-files=all` before the next edit.
The first half already existed (review R3-4), so it was not carried over: the rule was followed for `--dest` and missed for `--work`.

## Related
BUG-147 (the guard family; its doubly-backed guard); `reviews/2026-10-01-review-catalog-plan-code.md` finding R3-4 (the same lesson, one day old); BUG-094 (red-phase test with an outward side effect, same class).
