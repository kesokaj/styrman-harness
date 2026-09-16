---
owner: main
date: 2026-09-16
---
# Verification — the Go toolchain upgrade clears both advisories   (2026-09-16)

Range reviewed: no repository change. The tree is `bf48e64`, untouched; what changed is this box's Go.

## Status: PASS

Three security reviews across the day reported two standard-library advisories reachable from this code, and the v1.0.0 gate established that neither is exploitable here — `GO-2026-4602` is an escape from an `os.Root` this module never constructs, and `GO-2025-3956` needs `LookPath` to receive an empty or dotted path while the only call passes the literal `"go"`. They were left as hygiene rather than a defect. This closes them at the source, which is the toolchain, not the code.

## Why not the ways that do not work

`apt` has nothing: Debian trixie ships `golang-1.24-go 1.24.4-1` and no later `golang-1.*-go` package, while the fixes are in go1.24.6 and go1.25.8. Raising the `go` directive or adding a `toolchain` line to `go.mod` is refused by ADR-029 — measured earlier the same day to make the module unbuildable on any host that cannot fetch a toolchain, including through `go run`, before this repository's code executes. `/usr/lib/go-1.24` is root-owned and this project records no standing root permission, so replacing the system package was not this session's step to take.

## What was done, all user-level

```
go install golang.org/dl/go1.27.1@latest        → wrapper in ~/.go/bin
go1.27.1 download                               → SDK in ~/sdk/go1.27.1 (70.5 MB, checksum verified by the tool)
ln -sfn ~/sdk/go1.27.1/bin/go ~/.local/bin/go   → ~/.local/bin is position 1 on PATH
go install golang.org/x/vuln/cmd/govulncheck@latest
```

go1.27.1 is the newest patch available; `GOTOOLCHAIN=go1.27.1 go version` resolved it before anything was installed. govulncheck had to be rebuilt: the old binary carried go1.25 source-processing packages and refused the newer standard library with "unknown field rfd in struct literal". No root, no system file touched, and `/usr/bin/go` still reports go1.24.4 for anything that asks for it by absolute path.

## Measured after

| Check | Before | After |
|---|---|---|
| `go version` | go1.24.4 | **go1.27.1** |
| `govulncheck ./...` | `affected by 2 vulnerabilities` | **`No vulnerabilities found.`** |
| `gofmt -l cmd internal` | clean | clean — a newer gofmt reformatted nothing |
| `go vet ./...` | clean | clean |
| `go test -race -count=1 ./...` | pass | pass |
| Offline build (empty `GOMODCACHE`, `GOTOOLCHAIN=auto`, `GOPROXY=off`) | exit 0 | **exit 0** |
| Deployed binary's toolchain | go1.24.4 | **go1.27.1** (`go version -m ~/.local/bin/sh-agent`) |

The offline-build row is the one that mattered most: `go.mod` still says `go 1.24.4` with no `toolchain` line, so a consumer on an older Go is unaffected and a host with no network still builds. The upgrade is a property of this machine, not of the module.

`sh-agent install` was re-run afterwards, so the binary in `~/.local/bin` is built by go1.27.1; `status` reports `drift none`.

## What this does not change

The `v1.0.0` tag. It is unchanged and still correct: its gate report records the advisories as present-and-unreachable at the time it was cut, which is what was true. A future release cut from this box will simply not report them.

## Not examined

The Mac, which has its own Go and its own upgrade path (`brew upgrade go`). Whether the two advisories were ever reachable — the gate established they were not, and nothing here revisits that.

## Environment

Linux (Debian 13, amd64). `~/.local/bin/go` → `~/sdk/go1.27.1/bin/go`; undo by deleting that symlink. govulncheck v1.8.0, vulnerability database updated 2026-09-15.
