---
owner: verifier
date: 2026-09-03
scope: changeset:f6f6c28eee97
verdict: PASS
---

# Verification — local tier   (2026-09-03)

Produced by `harness dispatch`, which ran the project's declared commands
directly. No model was involved: every command below exited zero.

Commands: `gofmt -l . | tee /dev/stderr | (! read)`, `go vet ./...`, `go test -race ./...`

## `gofmt -l . | tee /dev/stderr | (! read)` — PASS

```

```

## `go vet ./...` — PASS

```

```

## `go test -race ./...` — PASS

```
?   	github.com/kesokaj/styrman-harness/cmd/harness	[no test files]
ok  	github.com/kesokaj/styrman-harness/internal/dispatch	(cached)
ok  	github.com/kesokaj/styrman-harness/internal/gate	(cached)
ok  	github.com/kesokaj/styrman-harness/internal/hookio	(cached)
ok  	github.com/kesokaj/styrman-harness/internal/install	(cached)
ok  	github.com/kesokaj/styrman-harness/internal/trail	(cached)
```
