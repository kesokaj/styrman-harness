---
owner: bugs
date: 2026-09-15
id: BUG-030
status: fixed
area: testing (Phase 6 kit)
---
# BUG-030: Phase 6.4 comparison's first launch died at its base install

## Symptom
`.agents/scratch/phase6-kit/compare.sh` exited 1 on its base-tree install:
`sh-agent: checkout /home/kesokaj/workspace/styrman-harness/.agents/scratch/base-c192622 is not the
installed one (/home/kesokaj/workspace/styrman-harness); pass --repo
/home/kesokaj/workspace/styrman-harness/.agents/scratch/base-c192622 to build and install from it`.

## Repro
Script installs c192622's harness (a plain copy at `.agents/scratch/base-c192622`) into the sandbox
home to re-run rows that failed twice: `go run ./cmd/sh-agent install --home <sandbox>` run from that
copy, without `--repo`.

## Root cause
The installer's checkout guard correctly refused to build from a tree it did not recognize as the
installed one. The script never passed `--repo` for the second tree, and it was launched detached
(long-running, backgrounded) without its first CLI step ever having been run by hand first, so the
missing flag was not caught before launch.

## Fix
Added `--repo` to every install and status call, for both trees; relaunched 22:09:30. Base install
then took the plain-copy path with rc=0. First attempt's files kept as `results/attempt1-*`; exit
trap had already reinstalled this tree (`done 2026-09-15T22:07:53+02:00`).

## Prevention
Run a detached script's setup steps once in the foreground before launching it, or give it a dry-run
mode. A script that installs a second tree into a home passes `--repo` on every install and status
call for that tree.

## Related
BUG-006, BUG-015 (installer's checkout/build-source guards); BUG-013 (detached long-running process
needs its own status file, run here as `results/attempt1-*`).
