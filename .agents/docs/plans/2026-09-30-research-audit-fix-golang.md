---
owner: research
date: 2026-09-30
---

# Research brief: fixing the `golang` skill (audit TS2, TS3, TS6–TS10, TS24)

Scope: `skills/golang/SKILL.md` (kept; its rules are the yardstick) and the seven `references/` files (5,863 lines, 279 KB). Dates seen are 2026-09-30 unless noted; local checks ran on go1.27.1 with `go.mod` at `go 1.24.4`.

## Recommendation (one line)

Delete the seven references and rewrite two short ones (`references/service.md`, `references/testing.md`) from the Go release notes and Google's docs; keep a `LICENSE` + `VENDORED.md` row only while any upstream sentence survives, and drop both once a `diff` against the upstream pin shows none. Reason: the files are 27 upstream skills flattened with their personas, multi-agent workflow and library defaults intact; 4 of 7 contradict `SKILL.md` § Dependencies on their first page.

## 1. Upstream: samber/cc-skills-golang

- Repository: <https://github.com/samber/cc-skills-golang> (MIT; description "A collection of Golang agentic skills that works"; default branch `main`; `pushed_at` 2026-09-07). `samber/cc-skills` (the repo the audit names) holds 24 non-Go skills and its README redirects Go readers here; the only Go-related skill left there is `promql-cli`, which our references cite as `samber/cc-skills@promql-cli` (4 places) — that pointer resolves to a skill this harness does not install.
- Layout: 46 directories `skills/golang-*` (each `SKILL.md` + `references/` + `evals/`), plus `rules/golang-always.mdc`. Our seven files flatten 27 of them under `## Sub-topic: golang-<name>` headings. **`golang-http` (the first sub-topic of `golang-io-web.md`) has never existed upstream** — `git log --all` over the clone shows no path matching `skills/*http*` — so that sub-topic's provenance is unknown; the trail (`changes/`, the 2026-09-15 review sweeps) does not say who wrote it.
- LICENSE (fetched raw): "MIT License / Copyright (c) 2026 Samuel Berthe / … The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software." So: a fork under the same name `golang` must carry the licence text as long as substantial portions of the upstream prose remain, however renamed. The current tree has neither a `LICENSE` nor a `VENDORED.md` row (audit TS3, confirmed).
- Most recent upstream commits (shallow clone, `git log`): HEAD `19a0626` 2026-09-07 "Revise Codex plugin installation instructions (#105)"; last commit touching `skills/`: `bac46b0` 2026-09-02 "docs: improve prose and sweep the skills (#102)"; `8f8e2fe` 2026-09-02 "feat(golang-modernize): add Go 1.27 support, fix sibling skill drift (#101)" touched `golang-concurrency`, `golang-testing`, `golang-dependency-management`, `golang-troubleshooting`. Per-directory last commits for the 27 flattened skills fall between 2026-08-19 (`golang-code-style`, `golang-cli`) and 2026-09-02.
- Drift against `bac46b0`: upstream already says `synctest.Test` (Go 1.25+) at `golang-testing/SKILL.md:217`, `wg.Go` 1.25+ at `golang-concurrency/SKILL.md:82`, and "do not create a new `tools.go`" at `golang-dependency-management/SKILL.md:112`; the v1 golangci-lint install path at `golang-lint/SKILL.md:39` is upstream's own bug.
- Verdict: a re-fetch fixes those three drifts and nothing else — the 29 personas, "5 parallel sub-agents" (×8), `ultrathink`, `Context7`, `.modernize` and the sqlx/Cobra/Viper/Prometheus defaults are upstream's design. Rewrite.

## 2. Go language and toolchain facts to pin

Dates (<https://go.dev/doc/devel/release>): go1.24.0 2025-02-11, go1.25.0 2025-08-12, go1.26.0 2026-02-10, go1.27.0 2026-08-19, go1.27.1 2026-09-01; "supported until there are two newer major releases" — 1.25–1.27 are supported, **1.24 is not**, and this repo's `go 1.24.4` is below the floor.

**Verified locally (lead with this):** with `go 1.24.4` on go1.27.1, `go build` passes but `go vet` and `go test` fail: `sync.(*WaitGroup).Go requires go1.25 or later (file is go1.24)`; with `go 1.25` they pass. go1.27 notes: "`go test` now invokes the `stdversion` vet check by default". A skill recommending 1.25+ APIs must say "raise the `go` line first".

| Feature | Version | Source line (release notes) |
|---|---|---|
| Swiss-table maps | 1.24 | go1.24: "a new builtin `map` implementation based on Swiss Tables" (`GOEXPERIMENT=noswissmap` to disable) — so `golang-core.md:230` "8-entry buckets and overflow chains" describes the pre-1.24 runtime |
| `go.mod` `tool` directive, `go get -tool` | 1.24 | go1.24: "removes the need for the previous workaround of adding tools as blank imports to a file conventionally named 'tools.go'" |
| `crypto/rand.Text`, `os.Root`, `weak`, `omitzero`, `runtime.AddCleanup`, generic type aliases | 1.24 | go1.24 |
| `testing.B.Loop`, `T.Context`/`B.Context` | 1.24 | go1.24: "`for b.Loop() { ... }` in place of … `b.N`"; "return a context that's canceled after the test completes and before test cleanup functions run" |
| `go vet` `tests` analyzer | 1.24 | go1.24 |
| `testing/synctest` experiment (`synctest.Run`, `GOEXPERIMENT=synctest`) | 1.24 only | go1.24 |
| `testing/synctest` GA as `synctest.Test(t, func(*testing.T))` | 1.25 | go1.25: "The `Test` function runs a test function in an isolated 'bubble'"; the 1.24 API "persists under `GOEXPERIMENT=synctest` but will be removed in Go 1.26". Verified: `go doc testing/synctest` on 1.27.1 lists only `Test`, `Wait`, `Sleep` |
| `synctest.Sleep` | 1.27 | go1.27 |
| `sync.WaitGroup.Go` | 1.25 | go1.25: "The new `WaitGroup.Go` method" — `golang-core.md:61` ("Go 1.24+") is wrong, `golang-architecture.md:750` ("1.25+") is right |
| Container-aware `GOMAXPROCS` | 1.25 | go1.25: "the runtime considers the CPU bandwidth limit of the cgroup"; `containermaxprocs=0`, `updatemaxprocs=0` — relevant to Cloud Run CPU limits |
| `go vet` `waitgroup`, `hostport` | 1.25 | go1.25 (1.27 renames `waitgroup` → `waitgroupgo` in `go fix`) |
| `encoding/json/v2` | experiment 1.25 (`GOEXPERIMENT=jsonv2`); 1.26 notes silent (still an experiment); **default in 1.27** | go1.27: new `encoding/json/v2` and `encoding/json/jsontext`; "The `encoding/json` package is now backed by the v2 implementation … the exact text of error messages may differ"; stricter: rejects invalid UTF-8 and duplicate object names; `GOEXPERIMENT=nojsonv2` opts out. `golang-architecture.md:764` "(Go 1.25+, experimental)" is stale |
| Green Tea GC | experiment 1.25; **default 1.26** | go1.26: "now enabled by default"; `GOEXPERIMENT=nogreenteagc` "expected removal in Go 1.27" |
| `go fix` = modernizers | 1.26 | go1.26: "completely revamped and is now the home of Go's _modernizers_", same analysis framework as `go vet`; 1.27 adds `atomictypes`, `embedlit`, `slicesbackward`, `unsafefuncs` |
| `new(expr)`, `T.ArtifactDir`, `slog.NewMultiHandler`, `ServeMux` slash redirect 307, PQ hybrid TLS default | 1.26 | go1.26 |
| Generic methods, `httptest.NewTestServer`, goroutine-leak profile GA, HTTP/1 body auto-drain on Close | 1.27 | go1.27 |

Not verified here: `golang-io-web.md:30` and `golang-architecture.md:727` date `httputil.ReverseProxy.Rewrite` to Go 1.26; it is older (1.20 from memory) — check go1.20 notes before keeping any date on it.

## 3. golangci-lint v2

- Install (<https://golangci-lint.run/docs/welcome/install/local/>): current **v2.14.0**; `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0`, with the warning "Using `go install`/`go get`, 'tools pattern', and `tool` command/directives installations aren't guaranteed to work. We recommend using binary installation": `curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.14.0`. `golang-ecosystem.md:235,250` use the v1 path (`…/golangci-lint/cmd/golangci-lint`), which installs v1.
- Defaults (<https://golangci-lint.run/docs/welcome/quick-start/>): "By default, the following linters are enabled: errcheck, govet, ineffassign, staticcheck, unused"; commands `golangci-lint run`, `golangci-lint fmt`, `golangci-lint migrate` (v1→v2 config), `golangci-lint config verify`.
- Config (<https://golangci-lint.run/docs/configuration/file/>): root `version: "2"`; `linters.default` ∈ `standard|all|none|fast`; `linters.enable|settings|exclusions`; `formatters.enable` ∈ `gci|gofmt|gofumpt|goimports|golines|swaggo`; `issues.new-from-rev`; `run.timeout` default `0` (`golang-ecosystem.md:412` "5m" is v1). `modernize` is a listed autofix linter (<https://golangci-lint.run/docs/linters/>).
- Minimal config for a stdlib-first service:

```yaml
version: "2"
linters:
  default: standard        # errcheck, govet, ineffassign, staticcheck, unused
  enable: [gosec, nolintlint, modernize]
formatters:
  enable: [gofmt, goimports]
```

## 4. Cloud Logging, Trace, Profiler, Error Reporting from Go on Cloud Run

Structured-log contract (<https://docs.cloud.google.com/logging/docs/structured-logging>, <https://docs.cloud.google.com/run/docs/logging>, LogEntry reference <https://docs.cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry>):
- Special JSON keys lifted out of `jsonPayload`: `severity`, `message`, `time`/`timestamp`, `logging.googleapis.com/trace`, `logging.googleapis.com/spanId`, `logging.googleapis.com/trace_sampled`, `logging.googleapis.com/sourceLocation` {`file`, `line`, `function`}, `logging.googleapis.com/labels`, `logging.googleapis.com/operation`, `logging.googleapis.com/insertId`, `httpRequest`.
- `severity` values: `DEFAULT DEBUG INFO NOTICE WARNING ERROR CRITICAL ALERT EMERGENCY`. slog prints `WARN`, so `LevelWarn` must be mapped to `WARNING`; `DEBUG/INFO/ERROR` match as-is.
- `trace` "must follow the format `projects/PROJECT_ID/traces/TRACE_ID`"; `spanId` "a 16-character, hexadecimal encoding of an 8-byte array". Cloud Run: "Container logs are not automatically correlated to request logs unless you use a Cloud Logging client library" — the trace field is what correlates them.
- Headers (<https://docs.cloud.google.com/trace/docs/trace-context>, seen 2026-09-25 UTC stamp): `traceparent: VERSION-TRACE_ID-PARENT_SPAN_ID-TRACE_FLAGS` (hex); legacy `X-Cloud-Trace-Context: TRACE_ID/SPAN_ID;o=OPTIONS` where "SPAN_ID is a 64-bit decimal representation"; "When possible, use the `traceparent` header". Cloud Run "automatically populates the W3C standard header `traceparent`" and samples "at a maximum rate of 0.1 requests per second for each instance" (<https://docs.cloud.google.com/run/docs/trace>, 2026-09-29).
- Official samples: `GoogleCloudPlatform/golang-samples/run/logging-manual/main.go` (2024-05-08) builds the trace as `fmt.Sprintf("projects/%s/traces/%s", projectID, traceParts[0])` from `X-Cloud-Trace-Context`; `opentelemetry/instrumentation/app/logger.go` (2025-01-13) is `slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replacer})` with `replacer` renaming `LevelKey→"severity"` (`LevelWarn→"WARNING"`), `TimeKey→"timestamp"`, `MessageKey→"message"`, wrapped by a handler whose `Handle` adds the three `logging.googleapis.com/*` keys from `trace.SpanContextFromContext(ctx)`. Caveat: it emits the raw `TraceID()` without the `projects/…/traces/` prefix the Cloud Run doc requires — prefix it. `AddSource: true` plus renaming `slog.SourceKey` to `logging.googleapis.com/sourceLocation` matches (slog's `Source` fields are `function`, `file`, `line`).
- Stdlib approach for the skill: one `ReplaceAttr` and one ~20-line wrapping handler that parses `traceparent` (fallback `X-Cloud-Trace-Context`) in middleware and adds trace/spanId/trace_sampled per record. No OTel needed for correlation alone.
- Error Reporting (<https://docs.cloud.google.com/error-reporting/docs/formatting-error-messages>, setup/go 2026-09-25): "automatically creates an error event when a log entry contains a stack trace and the severity level … is set to at least `ERROR`"; a message-only entry is forced in with `"@type": "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"`; evaluation order "`stack_trace`, then `exception`, and then `message`"; `serviceContext` {service, version} groups. So: log `ERROR` with `stack_trace` = `string(debug.Stack())` for panics/5xx and nothing else is needed.
- Cloud Trace: `github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace` is **v1.38.0 (2026-09-10) and deprecated — "will be archived after January 1st, 2027"** (<https://pkg.go.dev/github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace>). Replacement (<https://docs.cloud.google.com/stackdriver/docs/instrumentation/migrate-to-otlp-endpoints>, 2026-09-29): standard `otlptracegrpc` to `telemetry.googleapis.com:443` with `grpc.WithPerRPCCredentials(oauth.NewApplicationDefault(ctx))`, resource attribute `gcp.project_id=PROJECT_ID`, role `roles/telemetry.tracesWriter`; Google's first preference is still a collector when the environment allows one (<https://docs.cloud.google.com/trace/docs/setup/go-ot>, 2026-09-28: `autoexport.NewSpanExporter`, `autoprop.NewTextMapPropagator`). Skill position: OTel is a justified dependency only when spans are needed; log correlation alone is stdlib.
- Cloud Profiler (<https://docs.cloud.google.com/profiler/docs/profiling-go>, 2026-09-30): `cloud.google.com/go/profiler`, `profiler.Start(profiler.Config{Service, ServiceVersion})` "as early as possible", role `roles/cloudprofiler.agent`, "All officially maintained Go releases"; **no deprecation notice on the page** (a web search found only the Node.js repo archival of 2026-02-20). Optional; one import, one call.

## 5. Database on Cloud Run

- pgx v5.11.0 (2026-09-07, Go 1.25+; <https://pkg.go.dev/github.com/jackc/pgx/v5>): use the native interface when "the application only targets PostgreSQL" and "no other libraries that require `database/sql` are in use"; "`*pgx.Conn` … is not concurrency safe. Use package `…/pgxpool`". `database/sql` via `pgx/v5/stdlib` only when a library demands it (golang-migrate does).
- Cloud SQL (<https://docs.cloud.google.com/sql/docs/postgres/connect-run>): primary path is the Unix socket `/cloudsql/INSTANCE_CONNECTION_NAME` via `--add-cloudsql-instances` — "automatically encrypted without any additional configuration", no proxy, no module. The Go connector `cloud.google.com/go/cloudsqlconn` v1.25.3 (2026-09-28, **requires Go 1.26+**; <https://pkg.go.dev/cloud.google.com/go/cloudsqlconn>) adds IAM database auth (`WithIAMAuthN()`), `postgres/pgxv5.RegisterDriver(name, opts…)` for `database/sql`, `config.ConnConfig.DialFunc = d.Dial` for pgxpool, and `WithLazyRefresh()` for "environments where the CPU may be throttled" (Cloud Run without CPU-always-on). The Auth Proxy is the sidecar form of the same thing; not needed on Cloud Run.
- Migrations: golang-migrate v4.20.1 (2026-09-09; <https://pkg.go.dev/github.com/golang-migrate/migrate/v4>) with `source/iofs` over an `embed.FS` and `NewWithSourceInstance`; Atlas is Apache-2.0 with a Go SDK `ariga.io/atlas-go-sdk/atlasexec` that drives the CLI binary (<https://github.com/ariga/atlas>), so it is a second binary in the image.
- Solo-operator recommendation: pgx native + `pgxpool` over the `/cloudsql/…` socket (no extra module, one deploy flag); the connector only when IAM auth is wanted; golang-migrate embedded, as a `migrate` subcommand of the same binary run as a Cloud Run Job in `/preflight` before the new revision — not on every instance start.

## 6. Testing facts

- `t.Context()`/`b.Context()` 1.24; `b.Loop()` 1.24 (1.26: "no longer prevents inlining in the loop body"); `synctest.Test` 1.25, `synctest.Sleep` 1.27 (§ 2).
- `go test` flags (<https://pkg.go.dev/cmd/go#hdr-Testing_flags>): `-cover` "enable code coverage instrumentation"; `-covermode` "default is 'set' unless -race is enabled, in which case it is 'atomic'"; `-race` on linux/amd64, linux/arm64, darwin/amd64, darwin/arm64 (the two machines here); `-timeout` default `10m`; `-count`, `-shuffle`, `-run`.
- httptest: `NewRecorder` + `mux.ServeHTTP` for handlers, `NewServer` for clients, `NewTestServer` (1.27) for an in-memory network.
- Two references only: `references/testing.md` (~150 lines: table-driven, hand-written fakes, `t.Context`, `b.Loop`, `synctest.Test`, httptest, `-race`, coverage) and `references/service.md` (§ 7).

## 7. Proposed edits

| File (lines) | Verdict | Replacement |
|---|---|---|
| `golang-core.md` (662) | delete | 30 lines into `SKILL.md` § Concurrency: `wg.Go` **1.25**, `errgroup` only with a justification, "every goroutine has an owner and a stop". Drop personas (`:5,182,349,417`), sub-agent audits (`:81-90,387-395`), map internals (`:228-250`, pre-1.24), `samber/oops` at `:373,377`, `Context7` at `:321`, the "Deep Dive" pointers to files that do not exist. |
| `golang-architecture.md` (768) | delete | Style = `gofmt`, `go vet`, Google style guide link (already in `SKILL.md:62`). Modernize = one line: "`go fix ./...` (1.26+) and `go vet` are the modernizers; there is no `.modernize` file and no `/golang-modernize`". Keep the deprecated-packages table (`:713-727`) minus the unverified `Rewrite` row, in `service.md`. Drop "4+ args one per line even when the prompt asks otherwise" (`:15`) and `lo` at `:200`. |
| `golang-structure.md` (616) | delete | `SKILL.md:23,27,44` already state layout and hand-wired DI. Drop personas (`:5,370,473`), refactor sub-agents (`:478`), sqlx tag row (`:297`), Cobra+Viper (`:439`), 20 pointers to absent files. |
| `golang-testing.md` (532) | rewrite → `testing.md` | Replace `:141-186` (`synctest.Run` sample — **does not compile on 1.27.1: `undefined: synctest.Run`**, and on 1.24 `Run` took `func()`, not `func(*testing.T)`) with `synctest.Test`; drop "Experimental" note; drop persona/`ultrathink`/`gotests`/sub-agents (`:5-15`), `goleak` section (`:107-139`; stdlib: close what you start and assert), benchmark "Reference Files" to nothing (`:513-532`), `@promql-cli`. Keep `:340-382` (hand-written fake) verbatim — it is the one section that already matches `SKILL.md:44`. |
| `golang-io-web.md` (1068) | rewrite → `service.md` | Keep and correct the HTTP core (`:32-406`: routing 1.22, JSON helpers, timeouts, shutdown); add § Logging on Cloud Run (§ 4 contract, `ReplaceAttr` + trace handler, Error Reporting `stack_trace`), § Database (§ 5: pgx + pgxpool, `/cloudsql` socket, connector when IAM, golang-migrate embedded), § Lint (§ 3 config). Delete gRPC (`:463-655`), CLI/Cobra/Viper (`:881-1066`, 12 `assets/examples/*` pointers), "sqlx or pgx" MUST (`:668,674,695`), `alice` (`:28`), the `Rewrite (Go 1.26+)` row (`:30`). |
| `golang-ecosystem.md` (1088) | delete | Linter → § 3 block in `service.md` (fixes v1 install path `:235,250`, `tools.go` `:240-256`, `run.timeout` 5m `:412`). Observability (`:3-158`: Prometheus/Grafana/Datadog/Loki/Sentry/PostHog, five personas' worth) → § 4 in `service.md`. Security/troubleshooting/performance → already the `security` agent and `review-standards`; the ten "Reference Files" at `:572-596` point at nothing. |
| `golang-samber.md` (1129) | delete | Eight libraries, eight personas; `:722` and `:139,192,243` make `lo`/`oops` defaults against `SKILL.md:44,58`; `:627-641` does not compile (audit TS24), `:679-687` returns `nil` from a function with no result. `SKILL.md:58` already says the one sentence worth keeping. |
| `SKILL.md` | keep, edit | Reference table → two rows; add "raise `go` to 1.25+ before using `WaitGroup.Go`/`synctest.Test` — `go test` fails vet otherwise (verified 2026-09-30)"; JSON note: 1.27 `encoding/json` is v2-backed and rejects duplicate keys and invalid UTF-8. |
| new `skills/golang/LICENSE` | add | Upstream MIT text (Copyright (c) 2026 Samuel Berthe) — required while any substantial upstream portion remains (the `:340-382` fake sample, the HTTP core if kept). Remove in the commit whose `diff -r` against the pin shows no upstream prose left. |
| `skills/VENDORED.md` second table | add row | see below |

Proposed `VENDORED.md` row (second table, forks under upstream names):

`| golang | Samuel Berthe (samber) | github.com/samber/cc-skills-golang (`skills/golang-*`, 27 skills flattened into `references/`; the `golang-http` sub-topic has no upstream) | bac46b0 (2026-09-02; last commit touching `skills/`; HEAD 19a0626 2026-09-07) | `SKILL.md` local; `references/` rewritten 2026-09-30 to two files — surviving upstream text: <list after diff> | MIT |`

Per BUG-111 the pin is written only after `git archive bac46b0 skills | tar -x` into `.agents/scratch/` and `diff -r`; the flatten makes that diff per sub-topic, not per file.

## What was not checked

The `ReverseProxy.Rewrite` version; whether Cloud Run still injects `X-Cloud-Trace-Context` beside `traceparent`; whether Google's `logger.go` relies on a collector for the `projects/` prefix; Atlas commercial boundaries; `golang-http` provenance beyond the clone's history. Only this brief was written; the clone and compile fixtures are in the session scratchpad.
