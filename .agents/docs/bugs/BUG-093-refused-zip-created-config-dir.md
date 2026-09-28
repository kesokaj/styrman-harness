---
owner: bugs
date: 2026-09-28
id: BUG-093
status: fixed
area: cmd/sh-agent updateFromZip, internal/harness UnpackZip
---
# BUG-093: a refused `--from-zip` still created `~/.config/sh-agent`

## Symptom
E2E run 2 of `update --from-zip` against a zip with an escaping entry
(`…/x/../../escape`): the command correctly refused with exit 1 and the
escape message, but the target home held 2 entries afterward — an empty
`.config/sh-agent` — on a home that started empty.
(`.agents/docs/testing/2026-09-28-verify-update-from-zip-e2e.md`, run 2, first pass.)

## Repro
`go run ./cmd/sh-agent update --from-zip <zip-with-escaping-entry> --home <fresh dir>`,
then list `<fresh dir>` — before the fix it held an empty `.config/sh-agent`
tree; the command still reported exit 1.

## Root cause
`updateFromZip` called `harness.ConfigDir(home)` — which runs `os.MkdirAll`
on `~/.config/sh-agent` — before `harness.UnpackZip` validated the zip's
entry names. The refusal path (name escapes the unpack directory) fired
after the side effect, so a rejected zip still left a directory behind.
The unit tests only asserted no `zip-*` temp directory remained, so they
never caught the config dir being created; the plan's E2E criterion
("nothing is written under <scratch>") caught it because it counted every
entry under the home, not just the known temp artifact.

## Fix
`UnpackZip` creates its own parent directory, only after every entry-name
check has passed; `updateFromZip` no longer calls `ConfigDir` up front.
Regression test `TestUpdateFromZipRefusedZipWritesNothing`
(`cmd/sh-agent/zip_test.go`) failed before the fix ("a refused zip left 1
entries under the home") and passes after it. Uncommitted at time of
writing, part of the feature's single commit (plan
`.agents/docs/plans/2026-09-28-update-from-zip.md`).

## Prevention
A command that can refuse its input creates nothing before the refusal
point; its test for the refused case asserts the whole target is
unchanged (an empty home stays empty), not only that the obvious
temporary artifact is absent.

## Related
`.agents/docs/plans/2026-09-28-update-from-zip.md` Phase 3;
`.agents/docs/testing/2026-09-28-verify-update-from-zip-e2e.md`.
