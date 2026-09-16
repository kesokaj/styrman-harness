---
owner: verifier
date: 2026-09-16
---
# Verification — attribution-write   (2026-09-16, third and final pass)

Range reviewed: working tree (uncommitted) vs HEAD (1254672) — `cmd/sh-agent/main.go`, `cmd/sh-agent/main_test.go` (security-review fix pass since the prior log: create path opens `O_CREATE|O_EXCL` and `Sync`s, insert path `Sync`s the temp file before rename, a non-object settings file falls through to `attributionHint`'s one-line note, CRLF files keep their line endings; test gained `"crlf file keeps its line endings"` and the `"not an object"`/`"null"` cases now assert the `note:` line) plus the still-uncommitted `REQUIREMENTS.md`, `context/AGENTS.md`, `ROADMAP.md`, `README.md`, `.agents/docs/architecture/DECISIONS.md`, `.agents/docs/deploy/local.md`, `.agents/docs/deploy/remote.md`.

## Status: PASS

Scope: `gofmt`, `go vet`, `go test -race -v -count=1 ./...`, `install`/`status`/reinstall against a fresh scratch home, `list` — complete for this change set.

## Commands
- `gofmt -l .` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -v -count=1 ./...` → exit 0; output saved to `.agents/scratch/verify-attribution-write/go-test-3.txt`; counts: pass 79, fail 0, skip 1 at 1254672+dirty for `go test -race -v -count=1 ./...`; vs previous: `.agents/docs/testing/2026-09-16-verify-retro-2026-09-16-applied.md`, pass 64→79, skip 1→1 — explained: the rise is this feature's full `TestAttributionNote`/`TestAttributionSet` addition plus, since the last `attribution-write` log (which had pass 78), one more subtest, `"crlf file keeps its line endings"` in `TestAttributionSet` (now 14 subtests, was 13), landed by the security-review fix pass; skip unchanged (`TestFindRepo`, pre-existing, unrelated to this change).
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-attribution-write-3` → exit 0. `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`; `note: …/.local/bin is not on PATH …`; `settings: …/.claude/settings.json created with "attribution": {"commit": "", "pr": ""} — Claude Code's commit and PR trailers are off`.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-attribution-write-3` → exit 0. `drift  none — 259 files match the repo`; no `note` line about attribution (key now present); rendered `.claude/AGENTS.md` is 93 lines (`wc -l`).
- Second `go run ./cmd/sh-agent install --home .agents/scratch/verify-attribution-write-3` (same scratch home, key already present) → exit 0. `0 written, 260 unchanged, 0 removed, 0 kept — 260 files managed`; PATH `note:` line still present (unrelated) but no `settings:` line — idempotent.
- `go run ./cmd/sh-agent list` → `### Skills (27)`, `### Subagents (7)`, matching the dispatch's expected counts.

## Coverage of changed code

`attributionSet` (`cmd/sh-agent/main.go`) branch-by-branch, against `TestAttributionSet` subtests in `cmd/sh-agent/main_test.go`:
- missing settings file → `"missing file is created with only the key"`
- insert into an object with other keys, preserving formatting/mode/trailing newline → `"other keys, formatting and mode stay"`
- CRLF file keeps its line endings (new this pass) → `"crlf file keeps its line endings"`
- empty object `{}` → `"empty object"`
- key already present, off/custom value → `"off value is left alone"`, `"custom value is left alone"`
- malformed JSON → `"malformed is left alone"` (now asserts the `note:` line falls through to `attributionHint`)
- valid JSON but not an object, e.g. `["x"]` → `"not an object is left alone"` (now asserts the `note:` line)
- bare `null` → `"null is left alone"` (now asserts the `note:` line)
- file over the 1 MiB read bound → `"a file over the read bound is left alone"`
- `os.Lstat` symlink guard → `"a symlinked settings.json is never written through"`
- `--dry-run` → `"dry run writes nothing"`
- `--vendor agy` → `"antigravity only writes nothing"`
- `status` after the file is later removed (`attributionHint`) → `"status alone keeps the note"`

`abandonTemp` (`cmd/sh-agent/main.go:616`, called from the four write-error paths inside the insert branch — `tmp.Write`, `tmp.Sync`, `os.Chmod`, `os.Rename`, per lines 591/594/600/604/608) — **still not covered**. No subtest forces `os.CreateTemp`, `Write`, `Sync`, `Chmod` or `Rename` to fail, so `abandonTemp` itself, including its own `os.Remove(tmpName)` failure branch, is exercised by no test in this change. Same gap named in both prior runs of this log; unchanged by this security-review fix pass, which added the `Sync` call but not a test that forces its failure.

`readSettings` — exercised indirectly by every `TestAttributionSet` subtest that supplies a settings file, and by `TestAttributionNote`'s malformed-JSON and directory-instead-of-file cases. The `os.Open` failure path on a file that stats as regular but fails to open (e.g. a permission-denied race between the `os.Lstat`/size check and the open) is **not covered**.

`attributionOff`/`attributionHint` — asserted by the `TestAttributionNote` cases; the fall-through from `attributionSet`'s not-an-object/null/malformed branches into `attributionHint`'s one-line note is now asserted directly in `TestAttributionSet` (see above), closing the coverage gap that would have existed had those branches silently diverged from `attributionHint`'s output format.

`REQUIREMENTS.md`, `context/AGENTS.md`, `ROADMAP.md`, `README.md`, `.agents/docs/deploy/{local,remote}.md`, `.agents/docs/architecture/DECISIONS.md` — prose only; re-read against the code's current behaviour and consistent with it; no test framework applies to prose, and this run did not re-dispatch `/review --docs` (out of scope for a verifier, and unchanged since the round-3 cold read already on disk).

## Not examined

- The `abandonTemp` write-error paths (`tmp.Write`, `tmp.Sync`, `os.Chmod`, `os.Rename`, and `abandonTemp`'s own `os.Remove` failure) — no test triggers them.
- `readSettings`'s `os.Open` failure branch on a regular file that fails to open after passing the Lstat/size checks.
- `.agents/docs/plans/2026-09-16-attribution-write.md` and the five `.agents/docs/reviews/2026-09-16-*attribution-write*.md` files — untracked, written by other agents, out of scope for this verification.
- No Antigravity-side smoke check (`agy` CLI start / context grep) was run for this change; the change is Claude-only (`attributionSet` no-ops for `--vendor agy`, confirmed by `"antigravity only writes nothing"`).

## `git status --short --untracked-files=all` check

Ran after the full command set. Lines outside `.agents/docs/`:
```
 M README.md
 M REQUIREMENTS.md
 M ROADMAP.md
 M cmd/sh-agent/main.go
 M cmd/sh-agent/main_test.go
 M context/AGENTS.md
```
These are the pre-existing tracked modifications that make up the change set itself (present in `git status` at the start of this session, per the dispatch's own snapshot) — not new artifacts produced by this run. No untracked file appears outside `.agents/docs/` (`git status --short --untracked-files=all | grep '^??' | grep -v '^?? \.agents/docs/'` → empty): this run's own commands (`go test`, `install`, `status`, `list`) left nothing behind outside scratch and the docs tree. The seven stray junk files noted in the two prior runs of this log (`["x"]\n`, `null\n`, etc.) are no longer present in `git status` — resolved since the last log.

## Environment

Linux 6.12.107+deb13-amd64, `go` toolchain from PATH (version not queried by this run), repo at commit 1254672 with the uncommitted changes listed in Range reviewed above (working tree dirty by design for this verification).
