---
owner: bugs
date: 2026-09-14
id: BUG-015
status: fixed
area: installer/binary
---
# BUG-015: install rebuilds from source even in a read-only, git-less copy of the tree

## Symptom
Since commit 82fd4a6 (2026-09-13, "installer rebuilds from source", the BUG-006
fix), `sh-agent install --repo <dir>` rebuilds the binary with `go build` into
`<dir>/.agents/scratch/sh-agent-build-<pid>/` whenever `go` is on PATH — with no
check that `<dir>` is this repo's own git checkout. Both downstream mirrors are
plain copies with no `.git`: styrman-crew's `image/harness/` (installed at
container boot by `image/entrypoint.sh:145-146` as the project user, image
root-owned and `chmod -R go-w`, Go 1.27.1 on PATH, under `timeout 20`) and
argolis-devbox's `deploy/harness/` (installed by
`deploy/scripts.d/032_setup-harness.sh:253` via `runuser`, PATH has no `go`).
Predicted from the code, not executed on either consumer (read-only in this
session; `.agents/scratch/mirror-facts.json` is a scratch note, not a cited
source): in crew the scratch `MkdirAll` fails EACCES, install exits 1, the
entrypoint withholds the harness stamp and logs "WARN: sh-agent install
reported errors or timed out — harness not deployed this boot" every boot; in
devbox install prints the misleading "warning: go is not on PATH; the running
binary was copied, not rebuilt — install Go and run `sh-agent update`" on a
tree where `update` is refused by design.

## Repro
Copy `context/`, `skills/`, `agents/`, `commands/`, `rules/` into a directory
with no `.git`, put `go` on PATH, run `sh-agent install --repo <dir>` — before
the fix it attempts a `go build` and either fails (unwritable tree) or
prints the toolchain warning on a tree that was never meant to build itself.
`TestPlainCopyNeverRebuilds` (cmd/sh-agent/main_test.go) reproduces this; it
failed with "go: go.mod file not found" before the fix.

## Root cause
`shouldRebuild(goOnPath, child)` considered only the toolchain and the child
marker, not whether the tree is a git checkout. The rebuild exists to keep
the deployed binary at the checkout's commit (BUG-006's fix), which is
meaningless for a plain copy whose image already built the running binary
from that same tree — the two downstream mirrors are exactly that shape, and
`.agents/docs/deploy/remote.md § Downstream mirrors` names both, but the
BUG-006 fix was not re-checked against that section before landing.

## Fix
`shouldRebuild(goOnPath, child, ownCheckout bool)`; new
`harness.IsCheckout(repo)` (`os.Lstat(<repo>/.git)`) gates the rebuild; a
plain copy installs the running binary as-is and prints "plain copy, not its
own git checkout: the running binary is installed as is" — the no-toolchain
warning now prints only for a checkout. `TestPlainCopyNeverRebuilds` and an
extended `TestShouldRebuild` (ownCheckout dimension) cover it. Describing
text moved with the behavior: AGENTS.md "Rules for working here", README.md
command table, `.agents/docs/deploy/remote.md § Downstream mirrors`.

## Prevention
A behavior that depends on the environment the binary runs in (toolchain on
PATH, writable tree, git present) must be designed against every deployment
shape the trail names — `.agents/docs/deploy/remote.md § Downstream mirrors`
lists both consumers and their call sites. A change to `install` is not done
until that section has been re-read and those call sites traced, not just
the sandbox scratch-home path this repo tests against itself.

## Related
BUG-006 (the rebuild this fix narrows), BUG-004 (a prior installer defect
that only showed up on a downstream consumer, not the sandbox). Recommend
the `install`/`status` review checklist (or `documentation`'s remote.md
maintenance step) add: "re-read Downstream mirrors and trace both call
sites" as a standing item, since this is the second installer defect
(BUG-004, BUG-015) invisible from the sandbox alone.
