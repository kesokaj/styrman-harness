---
owner: reviewer
date: 2026-09-16
---
# Review — linkshelf (Antigravity lane, code)   (2026-09-16)

Range reviewed: sim-agy-linkshelf c520712..a55d328 (repo: /home/kesokaj/sandbox/runs/sim-agy-linkshelf, dirs `cmd internal go.mod scripts`)

## Verdict: APPROVE WITH NITS

This is careful, idiomatic Go: atomic disk-first persistence (temp file + fsync + rename + 0600 perms), a correctly-configured `http.Server` with all four timeouts, clean `signal.NotifyContext` shutdown that drains and treats `http.ErrServerClosed` as success, `context.Context` threaded to every store call and checked before and after lock acquisition, and tests that genuinely exercise CRUD, persistence-across-restart, and 20-goroutine concurrency under `-race`. `gofmt -l .`, `go vet ./...`, and `go test -race ./...` all pass clean (exit 0 on all three; test output: `ok` for `model`, `server`, `store`, `[no test files]` for `cmd/linkshelf`). The one real defect is a dependency (`github.com/google/uuid`) that this project's own `go.mod` and `BUG-001` record show was added in direct response to a security finding that today's harness standards (BUG-041, `review-standards` §5) classify as a false positive: `crypto/rand` bytes formatted into a UUID string is stdlib-crypto-used-as-designed, not hand-rolled cryptography, and Go 1.24 ships `crypto/rand.Text()` as a no-dependency alternative. Nothing else rises above a nit.

## Findings

- **[medium / confidence: medium]** `go.mod:5`, `internal/store/store.go:15,99` — `github.com/google/uuid` is an unjustified dependency by current standards. The project's own `.agents/docs/bugs/BUG-001-hand-rolled-uuid.md` documents that this was added solely because a security review flagged `crypto/rand`-based UUID formatting as "hand-rolled crypto." Per the harness's corrected rule (`review-standards` §5, golang skill, this repo's own BUG-041 record), stdlib `crypto/*` used as designed is not hand-rolled, and any finding demanding a dependency must name the stdlib alternative it rejects — none was named here, and Go 1.24.4 (the `go.mod` toolchain version) ships `crypto/rand.Text()`. Fix: drop the dependency, generate IDs with `crypto/rand.Text()` or the original hand-formatted-but-correct `crypto/rand` bytes, remove `go.sum`.
- **[low / confidence: low]** `internal/server/server.go:83,177` — JSON decoders for `CreateLinkRequest` and `DeleteLinkRequest` don't call `dec.DisallowUnknownFields()`. Per the golang-io-web reference's JSON helper pattern, this is recommended so typos in field names surface as 400s instead of being silently dropped. Low severity: no correctness or security impact for this local-only service, only a slightly worse error message on a malformed client request.

## Criteria

1. Dead code — pass. No unused exports, no commented-out blocks, no orphaned files; every function in the diff has a live caller (traced `store.Store`/`server.Store` interfaces to `main.go` wiring).
2. Flow tracing — pass. `r.Context()` flows from handler → store on every path; errors are wrapped with `fmt.Errorf("...: %w", err)` at every layer and logged once (in the handler) not double-logged; every new test can fail (verified by reading assertions against real HTTP codes/bodies and real file-system state, not source text); persistence-across-restart is covered both in `store_test.go` (`TestFileStore_PersistenceAcrossReload`, reopens the store) and in `scripts/e2e_test.sh` (kills and restarts the actual binary, step 9); the e2e script's evidence is real — it starts the compiled binary, waits on a live health poll, and asserts on actual `curl` status codes and bodies, not canned output.
3. Standards conformance — pass except the two findings above. `golang` skill layout, error wrapping, `log/slog` JSON, `context.Context` propagation, MaxBytesReader, ServeMux Go 1.22 method+path patterns, and graceful shutdown with `ErrServerClosed` all match the golang-io-web reference.
4. Directory layout — pass. `cmd/linkshelf/`, `internal/{model,server,store}/`, tests beside their sources; no stray files.
5. Simplicity (DIY-except-security) — fail on the one finding above (unjustified dependency); everything else is stdlib-first with no speculative abstraction (the duplicate `Store` interface in `server.go` and `store.go` is the idiomatic consumer-defined-interface pattern, not a needless layer).
6. Decoupling — pass. `server` package depends only on its own small `Store` interface (4 methods, sized to what it consumes), not on the `store` package's concrete type; `model` has no dependency on `server`/`store`; files are independently editable (model/server/store/main each own a clear seam).

## Not examined

- `.agents/docs/reviews/2026-09-16-linkshelf-review.md` and `-security.md` (this repo's own prior review/security passes) — read only their referenced BUG-001 record, not their full text.
- `.agents/docs/specs/linkshelf/api.md`, deploy guides, testing guide — not read; this review is code-only per the dispatch.
- No behavior outside `cmd`, `internal`, `go.mod`, `scripts` was in the diff range, so nothing else was in scope.

## What's built beyond the bare CRUD ask

- Delete accepts three input shapes: path `/links/{id}`, query `?id=`, and JSON body `{"id":...}`.
- GET by ID (`/links/{id}`) in addition to list.
- Config precedence flags > env vars > defaults (`HOST`/`PORT`/`DATA_PATH`), documented in `AGENTS.md`.
- Atomic disk-first persistence (temp file, fsync, chmod 0600, rename) rather than naive overwrite.
- Full HTTP server hardening: all four timeouts, `MaxBytesReader` request-size caps, graceful SIGINT/SIGTERM shutdown with a 5s drain.
- A real e2e script that builds the binary, drives it over HTTP, and restarts it to prove on-disk persistence.
