---
owner: bugs
date: 2026-09-14
id: BUG-016
status: fixed
area: cmd/sh-agent (tests)
---
# BUG-016: TestDryRunNeverBuilds asserted a checkout-only message from a test that never built a checkout

## Symptom
styrman-crew vendors this repo as a plain copy (no `.git`) under
`image/harness/`. Running the mirror's `go test ./...` there, `TestDryRunNeverBuilds`
FAILED: dry run printed the plain-copy line instead of `[dry-run] the binary
would be rebuilt from …`. styrman-crew recorded it as its own BUG-058
("upstream") on 2026-09-14 and its `VENDORED.md` said the mirror's tests are
"red by construction".

## Repro
Copy `context/ skills/ agents/ commands/ rules/` and `cmd/ internal/ go.mod`
to a directory with no `.git` and run
`go test ./cmd/sh-agent/ -run TestDryRunNeverBuilds` there.

## Root cause
`cmd/sh-agent/main_test.go`'s `TestDryRunNeverBuilds` passed `--repo` as
`filepath.Abs("../..")` — THIS checkout — and `IsCheckout`
(`internal/harness/repo.go`, an `Lstat` of `<repo>/.git`) is what selects the
rebuild-announcing dry-run message; in a plain copy `../..` is not a
checkout, so the test asserted a shape it never controlled. The sibling
`TestPlainCopyNeverRebuilds` already built its own fixture (`plainCopy(t)`);
the dry-run test did not.

## Fix
`cmd/sh-agent/main_test.go`: a `checkoutShaped(t)` fixture (line 236) —
`plainCopy(t)` plus an empty `.git` directory (`IsCheckout` needs only the
entry; `Inspect` tolerates one git cannot read) — and `TestDryRunNeverBuilds`
(line 199) now uses it instead of `../..`. Verified: `go test -race ./...`
green in the checkout; the fixture is what makes the result independent of
where the tree sits.

## Prevention
A test that asserts a shape-dependent path (checkout vs. plain copy) builds
that shape itself in a temp dir; `../..` in a test is the shape of wherever
the tree happens to sit, which a downstream mirror changes. Same family as
BUG-015's "designed against every deployment shape the trail names" —
`.agents/docs/deploy/remote.md § Downstream mirrors` — but this is a test
fixture, not the production code path, so recheck test fixtures too, not
just `install`/`status`, whenever a change touches `IsCheckout`.

## Related
BUG-015 (the plain-copy install path this test suite exists to cover);
styrman-crew's BUG-058 (the consumer's record of the same defect, filed as
"upstream" before this record existed).
