---
owner: bugs
date: 2026-09-18
id: BUG-079
status: open
area: internal/harness/release.go (UpdateFromRelease), shell hash cache
---
# BUG-079: `sh-agent update` reported success but left the running binary stale (bash hash cache)

## Symptom
On devbox, `sh-agent update` printed "verified and installed sh-agent-linux-amd64 v1.2.2" and "260 unchanged", but the next unqualified `sh-agent version` still printed v1.2.1. `which -a sh-agent` showed `~/.local/bin` first (correct order), giving no hint anything was wrong. Only after `rm -rf /usr/local/bin/sh-agent` (which broke the shell entirely: "No such file or directory") and `source .bashrc` did `sh-agent version` show v1.2.2.

## Repro
On a machine where `sh-agent` was earlier resolved by bash to `/usr/local/bin/sh-agent` (a stale v1.2.1 release binary) in the current shell session, running `~/.local/bin/sh-agent update` (found first on PATH) downloads and writes the new version only to `~/.local/bin/sh-agent`. The invoking shell's command hash table still points at `/usr/local/bin/sh-agent`; `sh-agent version` in that same session keeps running the old file. `which -a` reads PATH, not the hash table, so it shows no problem.

## Root cause
`UpdateFromRelease` in `internal/harness/release.go` writes the downloaded binary to `filepath.Join(home, BinPath)` and never checks whether that path is the file that is actually running (`os.Executable()`). This is a different mechanism from BUG-078: PATH order was correct here, so `PathShadows` found nothing to warn about. The actual cause is bash's command hash cache, which is invisible to any PATH-based check (`which -a`, `PathShadows`) and only cleared by `hash -r` or a new shell.
`Diagnosis missed:` BUG-078's fix modeled only PATH order; it did not model that a shell caches a resolved command path independent of PATH, so a correctly-ordered PATH can still run a stale file in an already-open shell.

## Fix
Not fixed.

## Prevention
Before shipping a fix, `update`/`install` should compare the running binary's path (`os.Executable()`, resolved) via `os.SameFile` against the file it just wrote; when they differ, the output should say the running binary was not replaced and is stale, and name the remedy (`hash -r` in bash, `rehash` in zsh). Docs (README / `remote.md`) should add `type sh-agent` (or `hash -r` before `which -a`) to the post-install/update verification steps, since `which -a` cannot see a shell's hash cache.

## Related
BUG-078 (PATH order shadow, same devbox, adjacent but distinct mechanism); BUG-076, BUG-077 (which binary is actually running). This is the third bug in this area (BUG-076/077/078/079) about the running binary vs. the installed one; recommend the prevention rule above be considered together with BUG-078's for graduation into `internal/harness/release.go` and `.agents/docs/deploy/remote.md` once fixed.
