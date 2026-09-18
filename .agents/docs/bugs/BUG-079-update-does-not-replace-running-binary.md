---
owner: bugs
date: 2026-09-18
id: BUG-079
status: fixed
area: internal/harness/release.go (UpdateFromRelease), internal/harness/pathguard.go, shell hash cache
---
# BUG-079: `sh-agent update` reported success but left the running binary stale (bash hash cache)

## Symptom
On devbox, `sh-agent update` printed "verified and installed sh-agent-linux-amd64 v1.2.2" and "260 unchanged", but the next unqualified `sh-agent version` still printed v1.2.1. `which -a sh-agent` showed `~/.local/bin` first (correct order), giving no hint anything was wrong. Only after `rm -rf /usr/local/bin/sh-agent` (which broke the shell entirely: "No such file or directory") and `source .bashrc` did `sh-agent version` show v1.2.2.

## Repro
On a machine where `sh-agent` was earlier resolved by bash to `/usr/local/bin/sh-agent` (a stale v1.2.1 release binary) in the current shell session, running `~/.local/bin/sh-agent update` (found first on PATH) downloads and writes the new version only to `~/.local/bin/sh-agent`. The invoking shell's command hash table still points at `/usr/local/bin/sh-agent`; `sh-agent version` in that same session keeps running the old file. `which -a` reads PATH, not the hash table, so it shows no problem.

## Root cause
`UpdateFromRelease` in `internal/harness/release.go` writes the downloaded binary to `filepath.Join(home, BinPath)` and never checks whether that path is the file that is actually running (`os.Executable()`). This is a different mechanism from BUG-078: PATH order was correct here, so `PathShadows` found nothing to warn about. The actual cause is bash's command hash cache, which is invisible to any PATH-based check (`which -a`, `PathShadows`) and only cleared by `hash -r` or a new shell.
`Diagnosis missed:` BUG-078's fix modeled only PATH order; it did not model that a shell caches a resolved command path independent of PATH, so a correctly-ordered PATH can still run a stale file in an already-open shell. Two later review rounds also found the fix firing where it should not (no installed file; the older-release branch where the running file is the newer one) — the first fix asked "is the process the installed file?" and not "did this run install anything?"; the second round's brief line that would have found it: "which return paths of UpdateFromRelease install nothing?".

## Fix
Three commits on `main`: `bc745dc` — `internal/harness/pathguard.go` `RunningElsewhere(pathEnv, binPath, exe)` finds the PATH entry that reaches the running process's file by name when that file is not the installed `~/.local/bin/sh-agent` (`os.SameFile` identity; a link to the installed file is not reported; a file no PATH entry reaches is nobody's `sh-agent`); `update` prints a note on both the checkout and release paths; README and both deploy guides put `hash -r`/`rehash` and `type sh-agent` before `which -a`. `0f2da25` — review round 1: the installed file must exist or nothing is reported. `a078848` — review round 2 (operator's D2): the note is printed by `UpdateFromRelease` itself only after a replacement and the new binary's install succeed, and by the checkout `update` after `Reinstall`; an up-to-date, `--check`, not-a-release or older-release run prints nothing; message ends "when that path is root-owned, prefer removing it". Tests: `TestRunningElsewhere` (8 cases), `TestUpdateFromReleaseReplacesAndInstalls`, `TestUpToDateUpdateSaysNothingAboutTheRunningBinary`, `TestCheckoutUpdateNamesTheBinaryItRanFrom`; seven must-fail mutations. Verifier logs `testing/2026-09-18-verify-bug079-stale-running-binary{,-2,-3}.md`. `testing/guide.md` § 5 now carries the post-update check (`hash -r`/`type sh-agent` before `which -a`).

## Prevention
Before shipping a fix, `update`/`install` should compare the running binary's path (`os.Executable()`, resolved) via `os.SameFile` against the file it just wrote; when they differ, the output should say the running binary was not replaced and is stale, and name the remedy (`hash -r` in bash, `rehash` in zsh). Docs (README / `remote.md`) should add `type sh-agent` (or `hash -r` before `which -a`) to the post-install/update verification steps, since `which -a` cannot see a shell's hash cache. Now also in `testing/guide.md` § 5's post-update check (a doc, not yet a formal graduation target — see Related for the recommendation).

## Related
BUG-078 (PATH order shadow, same devbox, adjacent but distinct mechanism); BUG-076, BUG-077 (which binary is actually running). This is the third bug in this area (BUG-076/077/078/079) about the running binary vs. the installed one; recommend the prevention rule above be graduated alongside BUG-078's into `context/AGENTS.md` (hard rules) or a `.agents/rules/` file once an operator confirms the target.
