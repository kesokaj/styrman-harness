---
owner: main
date: 2026-10-08
---
# Metrics store — Phase 1: the store

- New `internal/metrics/store.go`:
  - `Path`, `Open` (the search index's open rules), forward-only `migrations` (schema 1 = the spec's tables), a per-store salt, refusal of a newer schema (`ErrNewerSchema`, store untouched), and `Store.Export` (JSONL, 0600, refuses an existing file).
  - Tests in `store_test.go`: mode and path, symlink refusal, salt stable across reopen, v1 rows surviving an appended migration, newer schema refused with bytes unchanged, export round trip and refusal.
- New `internal/search/export.go`:
  - `PrepareDBFiles`, `DSN`, `IsBusy`, `OpenCopy` and `KindOf`, exported for the metrics store.
  - `ProtoFields`, a schema-less protobuf decoder that keeps field paths. Tests in `protofields_test.go`: paths, truncation, the depth bound.
- Restatements of "the only package that imports the driver" were corrected in `AGENTS.md` (Layout, which gains an `internal/metrics` line, and Rules) and in the `internal/search` package comments.
- Plan rows D31 (`TrailDates` goes to a new `internal/trail` in Phase 4) and D32 (a symlinked `~/.config/sh-agent` directory is allowed, as in search; the spec is corrected).
- Verification: `go test -race -timeout 30m ./...` exit 0, and gofmt and vet are clean. Mutation control on a copy: with the migration loop disabled, 3 tests fail.
- Agents dispatched: none. The diff was read by the author; the feature's one review runs at the Definition of Done.
