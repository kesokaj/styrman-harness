---
owner: main
date: 2026-09-18
---
# A release binary installs only what it carries (ADR-034)     (2026-09-18)

## Summary

The operator installed v1.2.0 on a second machine, `argolis-devbox`, which still had an old unpacked bundle at `/opt/devbox/deploy/harness` recorded in its manifest. The release binary installed THAT tree's months-old content and reported `drift none` — measured against the old tree. It surfaced only because that machine's Claude Code, asked to quote the rule about a guard in front of a destructive action, answered with the `rm -r` rule rather than the fail-closed rule shipped in v1.2.0. The machine was stuck too: `install` kept choosing the copy, `update` refuses a plain copy, and no flag said "use what you carry".

A release binary now installs the harness built into it and never a tree on the machine — not the working directory, not the recorded path, not one named on the command line — and refuses `--repo`/`$SH_AGENT_REPO` rather than appearing to accept one. Development builds keep the old order exactly, because installing the working tree is how this harness is developed. The operator's words settled the scope after a narrower first fix: "The binary should only use the published git repo in github always and only" (ADR-034, BUG-077).

## Type

bugfix

## Files changed

- `internal/harness/repo.go` — `ResolveSource` splits by kind of build: a release binary returns its embedded files and refuses a named path; a development build keeps explicit → `$SH_AGENT_REPO` → working directory → recorded path → the not-found error
- `internal/harness/harness_test.go` — `TestResolveSource` now pins the new rule in five situations (recorded plain copy, recorded checkout, standing inside a checkout, `--repo`, `$SH_AGENT_REPO`) and keeps the development-build order; `gitInit` helper added
- `cmd/sh-agent/main_test.go` — `TestReleaseBinaryWithoutCheckout` asserts `--repo` is refused whether it points at a real tree or not
- `cmd/sh-agent/main.go` — the `--repo` usage text says a release binary refuses it
- `.agents/docs/architecture/DECISIONS.md` — ADR-034; `.agents/docs/bugs/BUG-077-*.md` and `INDEX.md`

## Agents dispatched

`bugs`, for BUG-077; its record was read against the diff before committing. The gate's reviewer and security dispatches cover this change as part of the v1.2.1 release range.

## Verification

`gofmt`, `go vet`, `go test -race -count=1 ./...` clean. Reproduced first, then fixed, then checked with a release-shaped binary (`-trimpath -ldflags "-X main.version=v1.2.1"`) in four situations: the devbox shape (a plain copy recorded in the manifest), standing inside a plain copy, standing inside the harness's own live checkout, and `--repo` at a real tree. The first three install the embedded files and report `source embedded v1.2.1`; the fourth is refused with the release named. The old behaviour is pinned as gone: the test that asserted a recorded tree wins for a release version now asserts the opposite, and the development-build cases still pass unchanged.

## Docs

`cmd/sh-agent/main.go`'s usage text. `README.md` § Install without cloning needed no change — it never told a reader to pass `--repo` — and `.agents/docs/deploy/local.md` describes `update`, not source selection. ADR-034 carries the reasoning.

## Requirements

R8 and R9 are unaffected in letter and strengthened in practice: "a release's `sh-agent` binary with the harness built in" now means only that, on every machine, whatever else is lying around. Plan D4's order stands for development builds; D14 stands and is narrowed from "last resort" to "only source" for a release binary.
