---
owner: bugs
date: 2026-09-16
id: BUG-048
status: fixed
area: docs/setup commands
---
# BUG-048: the documented gate `gofmt -l .` fails on junk left in `.agents/scratch`

## Symptom
`gofmt -l . && go vet ./... && go test -race ./...` (the documented setup
command) failed with `.agents/scratch/tagcheck/src.go:1:1: expected
'package', found x` — a two-byte junk file left over from a release-procedure
rehearsal in gitignored scratch, unrelated to the change under test.

## Repro
Leave any non-gofmt-clean or invalid `.go` file under `.agents/scratch/` and
run `gofmt -l .` from the repo root; it fails even though `go build ./...`,
`go vet ./...` and `go test ./...` pass, because the Go toolchain skips
dot-directories and `gofmt` does not.

## Root cause
`gofmt -l .` walks every directory including dot-directories; `go build`,
`go vet`, and `go test` do not walk `.agents/scratch`. The global scratch
rule puts every generated file in `.agents/scratch/`, so any Go project under
this harness that runs a bare `gofmt -l .` has a standing trap: a scratch
artifact fails a gate that has nothing to do with the change under test.

## Fix
Documented command changed to `gofmt -l cmd internal` in `AGENTS.md:12`,
`README.md`, `.agents/docs/deploy/local.md`, and
`.agents/docs/testing/sh-agent/guide.md`; the junk file was deleted after
`git ls-files .agents/scratch/tagcheck` confirmed nothing tracked there.

## Prevention
A verification command that walks the whole tree is scoped to the source
directories on any project that keeps generated files under
`.agents/scratch/`; `agents/verifier.md` § What to run now names this caveat
for its Go auto-detect so the rule reaches every project, not just this one.

## Related
BUG-019/003/007 are the same shape (a rule's old wording or a tool's blind
spot survives a change); this is a fourth instance and the natural home
(`agents/verifier.md`) is where it graduated.
