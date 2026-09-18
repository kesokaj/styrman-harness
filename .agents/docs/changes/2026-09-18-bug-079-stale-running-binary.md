---
owner: main
date: 2026-09-18
---
# BUG-079 — `update` never replaces the file it runs from, and says so now     (2026-09-18)

## Summary

On devbox, `sh-agent update` to v1.2.2 succeeded and the next `sh-agent version` printed v1.2.1: the update had been run by `/usr/local/bin/sh-agent`, which `update` never replaces — it writes `~/.local/bin/sh-agent` and nothing else — and bash's command hash kept running the old file although PATH order was right. `update` now says, after it has written the installed binary, when the process that ran it is a different file that a shell reaches by name, and what to do about it. The guides put `hash -r` and `type sh-agent` before `which -a`. Unreleased.

## Type

bugfix

## What changed

- **`internal/harness/pathguard.go`**: `RunningElsewhere(pathEnv, binPath, exe)` — the PATH entry that reaches the running process's file by name when that file is not the installed one (`os.SameFile`; a link to the installed file is not reported; a file no PATH entry reaches is nobody's `sh-agent`; nothing is reported when the installed file does not exist). `NoteRunningElsewhere(home, out)` prints: `note: this update ran from <path>, not from the installed binary <installed>; update replaces only the installed one, so the file it ran from stays as it is. A shell that has already resolved `sh-agent` to that path keeps running it until `hash -r` (bash) or `rehash` (zsh); point it at the installed binary (`ln -sfn`) or remove it — when that path is root-owned, prefer removing it`.
- **`internal/harness/release.go`**: `UpdateFromRelease` prints the note only after a replacement and the new binary's install succeed; the up-to-date, not-a-release and older-release branches return before it (D2 of `plans/2026-09-18-sh-agent-works-everywhere.md`: on the older-release branch the running file is the newer one).
- **`cmd/sh-agent/main.go`**: the checkout `update` prints the same note after `Reinstall`; `install` never does; `--check` never does.
- **README § Install without cloning, `deploy/{local,remote}.md` § Verify, `testing/guide.md` § 5** (`documentation`): `hash -r`/`rehash` first, `sh-agent version` unqualified, `type sh-agent` when wrong, then `which -a`; five places that said `update` "replaces the binary" now say `~/.local/bin/sh-agent`; the by-hand install chains declared an accepted e2e-only check.
- **`.agents/docs/bugs/BUG-079-…`**, `INDEX.md`: `fixed`.

## What it does not do

It cannot see a shell's hash table. It says the fact the shell cannot: the file that ran the update is not the one just installed. The documented unqualified `sh-agent version` is still the check that catches a stale hash on its own.

## Tests

`TestRunningElsewhere` (8 cases over real files: the installed file itself, a link to it, a system copy reached by name, either PATH order, a file no entry reaches, an alias link to the system copy, a missing installed file, a missing exe), `TestUpdateFromReleaseReplacesAndInstalls` (the positive note beside the code, with the test process reachable by name), `TestUpToDateUpdateSaysNothingAboutTheRunningBinary` (`update` and `update --check` on the up-to-date path, installed file present, process reachable: silent), `TestCheckoutUpdateNamesTheBinaryItRanFrom`. Seven must-fail mutations in `.agents/scratch/bug079-mustfail/mutate.py`, each failing as a test failure.

## Definition of Done

`bug` skill: problem note; goldfish diagnosis (`plans/2026-09-18-recon-bug-079-diagnosis.md`, agreed: `release.go` writes `BinPath`, nothing compares it with `os.Executable()`); failing test first (`TestUpdateNamesTheBinaryItRanFrom`, failed for the message's absence); three review rounds at the cap:

- Round 1 (`bc745dc`): `verifier` PASS (229); `reviewer` APPROVE; `security` 1 medium — the note fired into an empty home naming a file that did not exist — + 5 low; docs lens APPROVE WITH NITS (five "replaces the binary" restatements; `guide.md` § 5). Fixed: the installed file must exist; `./sh-agent`; wording; the RTL characters an explorer had left in its brief.
- Round 2 (`0f2da25`): `verifier` PASS (230, the note observed live); `reviewer` APPROVE; `security` — the medium half closed: the note still fired on the older-release branch, calling a newer running binary stale. **Second round on one mechanism → the operator**: their D2, only after a binary was written.
- Round 3 (`a078848`): the note moved into `internal/harness`, printed by `UpdateFromRelease` after a replacement — the only place its positive case can be tested, since the release fixture compiles a test key into an unexported variable; `verifier` PASS (226, rename accounted; live note, `--check` silent); `reviewer` APPROVE; `security` no high, the medium closed on both paths; a pre-existing medium (`v1.02.2` passes `IsReleaseVersion`) and three lows go to Phase 2 of the plan as their own boxes, with the correction that `verifyRelease`'s marker check stops the rollback channel as described (`reviews/2026-09-18-security-c868ed4-a078848.md` § Outcome).
- `documentation` for the five restatements and `guide.md` § 5; `bugs` for the record. Every finding has an `## Outcome` line.

## Where it stops

Unreleased; ships with the plan's Phase 5 release. The audit items above are Phase 2's.
