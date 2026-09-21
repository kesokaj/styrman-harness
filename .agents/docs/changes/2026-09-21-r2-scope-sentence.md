---
owner: main
date: 2026-09-21
---
# REQUIREMENTS R2: native `AGENTS.md` reading's wider scope, on record; v1.6.0 smoke on the real home     (2026-09-21)

## Summary

The operator approved the sentence the v1.6.0 gate's security finding asked for: R2 now says native reading also loads an `AGENTS.md` in any directory above the project and, on read, in any subdirectory, and that a template `AGENTS.md` never sits beside source (ADR-037, same commit). No release: `REQUIREMENTS.md` is not part of what a tag delivers. The three real-home smoke checks passed on both CLIs on v1.6.0 (`testing/2026-09-21-verify-real-home-smoke-6.md`); the Mac updated to v1.6.0 by the operator.

## Type

docs (requirements, ADR, testing log)

## Verification

Local: gofmt, vet, `go test -race -count=1 ./...` exit 0, scratch install drift none. Real home: the smoke log above.

## Review

Docs only, eleven lines, read by the main agent.

## Agents dispatched

`changelog` (ADR-037).
