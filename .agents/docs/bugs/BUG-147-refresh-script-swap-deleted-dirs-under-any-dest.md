---
owner: bugs
date: 2026-10-01
id: BUG-147
status: fixed
area: skills/google-catalog/scripts/refresh-catalog.sh
---
# BUG-147: refresh-catalog.sh's swap ran rm -rf on every directory under --dest, whatever --dest was

## Symptom
The swap step removed every directory under `--dest` before copying the new catalog in. With `--dest ""` the path became `/`-relative; with any directory that was not a catalog (a project root, the home directory) every subdirectory there was deleted.

## Repro
Reproduced red by the main agent in a temporary directory: `bash refresh-catalog.sh --dest "" --work <dir> <pin> <pin>` against a fixture directory holding `precious/` removed it. Found by the security review, finding S1 (`reviews/2026-10-01-security-catalog-solutions-router.md`), which had only read the code.

## Root cause
The script trusted `--dest` to be a catalog and deleted by directory listing, not by what the catalog itself named. Nothing checked the value or the place before a destructive `rm -rf`.

## Fix
`--dest` and `--work` must be non-empty; `/`, the home directory, the checkout and an ancestor of either are refused by name; a non-empty `--dest` without a `CHECKSUMS` file is refused ("not a catalog, nothing is replaced there"); the swap removes only directories named as first path components in the old `CHECKSUMS`. Tests: `TestRefreshCatalogScript` in `internal/harness/catalog_test.go`: the empty `--dest` and not-a-catalog fixtures assert `precious/` survives; six named-place cases assert each guard's own message (with the three guard cases removed all six fail). Fixed in the working tree, before any release carried the script.

## Prevention
A script that deletes under a path it was given deletes only what its own manifest names there, refuses an empty value and the places that can never be the target (`/`, home, the checkout, their ancestors), and its test asserts that a bystander directory survives. A guard backed by a second check gets a test on the guard's own message, or removing the guard fails nothing.

## Related
BUG-146, BUG-148 (same script, same plan `plans/2026-10-01-google-skills-catalog.md`); review `reviews/2026-10-01-security-catalog-solutions-router.md`.
