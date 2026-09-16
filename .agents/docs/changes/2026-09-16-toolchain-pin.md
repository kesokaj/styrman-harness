---
owner: changelog
date: 2026-09-16
---
# Toolchain pin to go1.25.8     (date: 2026-09-16)

## Summary
`go.mod` pins `toolchain go1.25.8` to clear the two reachable stdlib advisories three security reviews had reported; the `go` directive stays at 1.24.4 so nothing downstream is forced to move. Before the pin, `govulncheck ./...` reported 2 affecting vulnerabilities in `os` (GO-2026-4602) and `os/exec` (GO-2025-3956); after it, 0. A build with `GOTOOLCHAIN=local` on go1.24.4 still succeeds and tests still pass, so the pin is a soft upgrade, not a break.

## Type
bugfix

## Files changed
- `go.mod` — the `toolchain go1.25.8` line added
- `AGENTS.md` — setup commands updated: `gofmt -l cmd internal && go vet ./... && go test -race ./...` with a note that the pin uses `GOTOOLCHAIN=auto` by default and `GOTOOLCHAIN=local` still builds
- `README.md` — setup commands updated to match AGENTS.md

## Agents dispatched
changelog only; no source changed, so no verifier or reviewer dispatch.

## Verification
Measured by the main agent this session before and after the pin: with `GOTOOLCHAIN=auto`, `go version` reports go1.25.8, `gofmt -l cmd internal` and `go vet ./...` are clean, `go test -race -count=1 ./...` passes, and `govulncheck ./...` reports "No vulnerabilities found. Your code is affected by 0 vulnerabilities". With `GOTOOLCHAIN=local` on go1.24.4, `go version` reports go1.24.4, `go build ./...` succeeds and tests pass.

## Docs
AGENTS.md and README.md setup commands updated for the pin and its behavior under both toolchain modes. `.agents/docs/architecture/DECISIONS.md` (ADR-028) records the decision and consequences.

## Requirements
n/a — the pin is a security hardening of the build environment, no change to module exports or consumer requirements.
