---
owner: main
date: 2026-09-28
---
# Verification — `update --from-zip` E2E with a real GitHub zip   (2026-09-28)

Range reviewed: 1cb2608..worktree. Plan: `plans/2026-09-28-update-from-zip.md`, Phase 3, E2E criteria 1–3.

## Status: PASS (after one fix, BUG-093)

## Setup, proven before the runs

- **The zip:** `gh api repos/kesokaj/styrman-harness/zipball/main` into `.agents/scratch/e2e-zip/main.zip`. It has 936 entries under ONE top folder, `kesokaj-styrman-harness-1cb26085…`, and `context/AGENTS.md` sits directly under that folder. It also holds this repository's tracked symlink `.agents/skills` as a symlink entry, which is the D6 case. All of this was checked with Python's `zipfile` before any run.
- **The binary:** `go build -ldflags "-X main.version=v0.0.1"`, which reports `sh-agent v0.0.1`, a release version. Its own source is therefore the embedded files, so a context that matches the zip can only have come from `--from-zip`.

## Runs

| # | Criterion | Result |
|---|---|---|
| 1 | Golden path: `update --from-zip main.zip --home home` | exit 0. It printed `zip … (sha256 0d66ecfc…): installing its harness with the running binary v0.0.1; the binary itself is not replaced`, installed 258 files, and left no `zip-*` directory under `home/.config/sh-agent/`. `status`: `source … this home was installed from a zip …`, `installed … from zip sha256:0d66ecfc8bda`, `drift not compared — … 258 of 258 manifest files present`, exit 0. The installed `~/.claude/AGENTS.md` is byte-identical (`cmp`) to the same zip's context rendered separately (`go run ./cmd/sh-agent render --repo <unzipped>`), and to the Antigravity copy. It ends with GANNET. |
| 2 | A zip with an `…/x/../../escape` entry (the fixture was checked to hold it) | exit 1, with `entry "…/x/../../escape" would escape the unpack directory: a .. element`. No `escape` file anywhere. **First run: the home held 2 entries**, an empty `.config/sh-agent`, because `ConfigDir` created it before the zip was checked (BUG-093). After the fix: **0 entries**. |
| 3 | Dry run | On a fresh home: exit 0, with a `[dry-run] zip …` line and **0 entries** in the home. On the populated home from run 1: a hash of every file and a list of every path, taken before and after, are **identical**, and no `zip-*` directory was left. |

## Fix made during the E2E

`UnpackZip` now creates the parent directory itself, after every name check has passed, and `updateFromZip` no longer calls `ConfigDir` first. The new test `TestUpdateFromZipRefusedZipWritesNothing` failed before the fix ("a refused zip left 1 entries under the home") and passes after it. The full race suite is green.

## Not exercised

- A zip downloaded through the browser's "Download ZIP" button. It has the same one-top-folder layout, with `styrman-harness-main/` as the folder name, and the unit test `TestUnpackZipGitHubTopFolder` covers that name.
- Antigravity reading the installed files. The install writes the same files as `install`, and those files are compared byte for byte above.
- macOS.

## Rerun after the review fixes (D11)

The binary was rebuilt from the working tree after the review round. The fixes it carries: `UnpackZip` validates before writing and returns the sum of the bytes it read; the `--from-zip` refusals run before any command; `update --dry-run` is read-only; `diff` has a zip branch; stale `zip-*` directories are swept.

| # | Result |
|---|---|
| 1 | exit 0. The first line carries `sha256 0d66ecfc…, UNVERIFIED — a zip carries no signature; …`. `status` shows `from zip sha256:0d66ecfc8bda` and `258 of 258 manifest files present`. The installed context is still byte-identical to the zip's render, and no `zip-*` directory is left. `diff`: `not compared — this home was installed from a zip that is not kept`. |
| 2 | Slip zip: exit 1, with **0** entries in a fresh home. |
| 3 | Dry run: a fresh home keeps **0** entries, and the populated home's hash list is identical before and after. |

## Rerun after round 2 (D13)

The binary was rebuilt again. E2E 1: exit 0, and the context is still byte-identical to the zip's render. E2E 2: exit 1, 0 entries. E2E 3: exit 0, 0 entries.

## Rerun after round 3

Final code (the sweep runs after a successful install; flag-like values are refused): E2E 1 exit 0, context identical; E2E 2 exit 1, 0 entries; E2E 3 exit 0, 0 entries.
