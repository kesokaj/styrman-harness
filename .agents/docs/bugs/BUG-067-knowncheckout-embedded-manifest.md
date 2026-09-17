---
owner: bugs
date: 2026-09-17
id: BUG-067
status: fixed
area: cmd/sh-agent/main.go (knownCheckout)
---
# BUG-067: `knownCheckout` resolved an `embedded:vX.Y.Z` manifest value as a filesystem path

## Symptom
An install manifest written by an embedded install records `"repo": "embedded:vX.Y.Z"` — a release name, not a path. `knownCheckout` handed that value to `samePath`, which resolves it with `filepath.EvalSymlinks` against the process's working directory. `ResolveSource` had already been hardened to skip the `embedded:` prefix and was tested for it (`internal/harness/repo.go:42`, `TestResolveSource`); `knownCheckout` was a second reader of the same field and was not.

## Repro
Not independently reproduced by a failing exploit: no bypass was demonstrated (the resolved relative value never equals the absolute checkout path in the cases checked), so the measured impact was a misleading message, not a wrong action. Found by code reading during the security review (`.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md`, finding 2).

## Root cause
When the `Repo` field gained a value of a new KIND (a release name where paths were previously stored), only the reader that prompted the change (`ResolveSource`) was guarded. `knownCheckout` was not re-examined even though it reads the same field.

## Fix
`knownCheckout` now refuses on `EmbeddedPrefix` before comparing, naming the release. `TestKnownCheckoutRefusesEmbeddedManifest` covers it and fails with the guard removed (commit 73f2c36).

## Prevention
When one field gains a value of a new kind, grep every reader of that field in the same change — guarding only the reader that prompted the change is half a fix. Same shape as BUG-036.

## Related
BUG-036 (grep every reader of the same file class). BUG-066 (same review, finding 1). Review: `.agents/docs/reviews/2026-09-17-security-f7ee817-21ba8dc.md`.
