---
owner: main
date: 2026-09-29
---
# Search index, Phase 1: the project index   (2026-09-29-1839)

Plan: `plans/2026-09-29-search-index.md`. Spec: `specs/sh-agent/search.md`.

- **New verbs** `sh-agent search` (words, `--refs`, `--open`, `--kind`, `--limit`, `--global`) and `sh-agent index` (`--rebuild`, `--global`), dispatched before the harness source is resolved, so a release-binary or zip place can search.
- **New package** `internal/search`: a per-checkout SQLite index at `.agents/scratch/search.db` with FTS5 (porter), a metadata graph of the trail (bug and ADR status, plan boxes, dates) and references (IDs and paths per line). Every query first reconciles against `git ls-files` and `git status`; content is stored per blob ID, and every query joins through the current path-to-blob view.
- **Refusals:** outside a checkout, an unignored index path, a symlinked `.agents/scratch`, a write lock held past 10 s. An unreadable or older-schema database is rebuilt with a note.
- **Dependency:** `modernc.org/sqlite` v1.46.1, the newest release that keeps `go 1.24.4` (D21). `govulncheck` names one unreachable advisory in `golang.org/x/sys`, whose fix needs Go 1.25 (D25).
- **`report`** never sends `search.db` and names it in the skipped list; `ConfigDir`'s comment names it.
- **Performance fix during the phase:** word search ran the full-text match once per row when joined directly (0.6–2.5 s); a materialized match makes it 23–67 ms (D24).
- **Docs:** README rows for both verbs and the harness-state line; the spec.

Tests: `internal/search/search_test.go` (21 tests, among them branch switch, dirty and untracked files, staged rename with edit, binary/large/symlink/submodule skips, GC, FTS sync after GC, one blob at two paths, corrupt and old-schema rebuilds, chunking, records, references, ranking and caps, FTS operator characters, refs, open, three refusals, busy, two concurrent processes, argument errors); `cmd/sh-agent/search_test.go` (flag reordering, both verbs end to end, usage against the doc comment); `TestCollectNeverSendsSearchDB`. Written after the code: eleven mutations, each compiled; ten failed their named tests, one equivalent (D22).

Agents dispatched: `verifier` (`testing/2026-09-29-verify-search-index-phase-1.md`, PASS over c04f2d7..7c24503, Phase 0 included). The reviewer and security reviews run once over the whole change at Done.

Verification: `gofmt -l cmd internal` empty; `go vet ./...` exit 0; `go test -race ./...` exit 0; scratch-home install and `status` show the binary line without a warning and no drift.
