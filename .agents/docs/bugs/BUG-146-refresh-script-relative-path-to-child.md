---
owner: bugs
date: 2026-10-01
id: BUG-146
status: fixed
area: skills/google-catalog/scripts/refresh-catalog.sh
---
# BUG-146: refresh-catalog.sh with a relative --dest/--work cloned for 3m32s, then exited 1 "the catalog is empty"

## Symptom
Phase 1 verifier of plan `2026-10-01-google-skills-catalog` (`testing/2026-10-01-verify-google-skills-catalog-phase-1.md`, Status FAIL) ran `bash skills/google-catalog/scripts/refresh-catalog.sh --dest .agents/scratch/... --work .agents/scratch/... <pins>` from the repository root. After 3m32s of cloning it exited 1 with three "the catalog is empty" lines. With absolute paths the same command exited 0 in 2.3 s.

## Repro
Any relative `--dest` or `--work`, run from any directory.

## Root cause
The script hands its build directory to a child that runs elsewhere: `(cd "$repo" && CATALOG_DIR="$build" CATALOG_WRITE=1 go test ./internal/harness ...)`. `go test` runs each package's tests with the package directory as working directory, so a relative `$build` named a directory that does not exist there. The defaults were absolute (derived from the script's location) and `TestRefreshCatalogScript` passed absolute temp paths, so neither the author's run nor the test met the relative case. The rule existed ("a path handed to a child process started with another working directory is made absolute before the handoff", `review-standards` § 8, from BUG-080) and was not applied to a shell script.

## Fix
`abs()` in the script makes `--dest` and `--work` absolute against the caller's directory right after argument parsing. `TestRefreshCatalogScript` now runs the script from another directory with relative `--dest rel-dest --work rel-work` and asserts the catalog lands under the caller's directory (seen red: exit status 1; then green).

## Prevention
A script that takes a path flag and starts a child with another working directory (a `cd`, a `go test`, a `git -C`) makes the path absolute at parse time, and its test passes at least one RELATIVE path from a directory that is not the repository root.
The rule already existed (`review-standards` § 8): it was not followed.

## Related
BUG-080 (same class; this is the second occurrence). Graduation candidate: the `bash` skill's script checklist; actionable (a relative-path test catches it) and generalizable (any path-taking script).
