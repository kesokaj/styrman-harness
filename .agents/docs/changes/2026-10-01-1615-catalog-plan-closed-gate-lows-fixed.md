---
owner: main
date: 2026-10-01
---
# 2026-10-01-1615 — Google skills catalog plan closed; the lows of its last review round and of the v2.1.0 gate's first pass fixed

Plan `plans/2026-10-01-google-skills-catalog.md` is closed (no open box; `ROADMAP.md` Now is empty). This entry covers what changed after 7dd0645, from four reports over v2.0.0..7dd0645: code review round 3, the docs lens, and the two security slices of the gate.

- **Dream timer (BUG-149)**: the "real home" that keeps the bare unit name is read from the account database (`os/user`), not from `HOME`. A sandbox run sets `HOME` to its own directory, so under the old lookup a sandbox `--remove-timer` or `uninstall` would have used the real timer's name. `TestUnitForIgnoresHomeOverride`. Limit: a build without cgo on macOS can only read `HOME` (`specs/sh-agent/dream.md`). `renderSystemd` escapes `$` in `ExecStart` (`TestRenderSystemdEscapesDollar`).
- **Refresh script**: exclusions are one name per line and compared whole (excluding `delta epsilon` no longer excludes `delta` and `epsilon`); the `--exclude` hint for an odd upstream name is shell-quoted; `--dest`, the checkout and the home directory are resolved through `..` and symlinks before the guard compares them; the header names the destination guards (D28).
- **Search**: `TestNoIndexMarkerDeletedCountsNoMore` covers the marker removed from the work tree.
- **Router** (`skills/google-catalog/SKILL.md` step 4): anything that sends the project's code, data, logs, prompts or credentials off the machine is shown to the operator first, like a change.
- **Tests only**: the six destination-guard cases name sources that cannot be cloned; a bystander directory survives the swap; the symlinked source set is refused; two more spellings of the checkout are refused.
- **Docs**: README (the timer row: per-home name, `--home`, the one-time remove and install after v2.0.0; "top-level frontmatter line"), `specs/sh-agent/memory-loop.md` (a symlinked `MEMORY.md` is refused by `memory add` and `show` only), `testing/guide.md` (the listing budget counts the skills a model is shown; rows that load third-party text run with read-only grants), project `AGENTS.md` (the same budget wording; a description can fail `go test`), `skills/VENDORED.md` (the script's guards; `npx skills@<version>`), `ROADMAP.md`. `bugs/BUG-145` had lost its escape sequences; restored.

Left with an owner (handoff): the lint skipped for a command or rule whose frontmatter does not parse (a new refusal); a notice naming `.sh-agent-noindex` directories in `status`; `--work` unchecked; hostnames in `deploy/remote.md` if the repository is ever published.

Verification before the commit (main agent): `gofmt -l cmd internal` empty, `go vet ./...`, `go test -race -count=1 ./...` all six packages ok; seen red first: the home override, the dollar sign, the hint quoting, the two-word exclusion, the two path spellings; checked with the fix removed on saved copies: the deleted marker, the bystander directory, the symlinked source set. The gate's verifier runs over this commit next. Not covered: a real `systemctl` or `launchctl` call for BUG-149 (the test stops at the unit name); macOS.

Agents dispatched: reviewer ×2 (code round 3, docs lens), security ×2 (code and text, trail), verifier ×1 (PASS over v2.0.0..7dd0645), documentation ×1, bugs ×1 (BUG-149).
