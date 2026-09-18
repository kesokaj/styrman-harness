---
owner: explorer
date: 2026-09-18
---
# Recon — BUG-079 diagnosis (2026-09-18)

## Question
Why did `sh-agent update` (release path) report v1.2.2 installed, but `sh-agent version` in the same shell still printed v1.2.1, until `/usr/local/bin/sh-agent` was removed and the shell re-sourced?

## Summary
`UpdateFromRelease` (internal/harness/release.go:90-179) always writes the new binary to a FIXED path, `dest := filepath.Join(home, BinPath)` (release.go:162, `BinPath = ".local/bin/sh-agent"`, plan.go:16) — never to wherever the currently-running binary actually lives. It never calls `os.Executable()` or compares `dest` against the running process's own path. When the running binary is `/usr/local/bin/sh-agent` (repro: an older binary there, `~/.local/bin` earlier on PATH but bash's command hash already resolved `sh-agent` to `/usr/local/bin/sh-agent`), the update writes v1.2.2 to `~/.local/bin/sh-agent`, runs `<dest> install` (which correctly reports "260 unchanged" because it's self-consistent with `~/.local/bin`), and exits — leaving the actually-invoked `/usr/local/bin/sh-agent` at v1.2.1 completely untouched. The next `sh-agent version` in that shell re-uses bash's cached hash table entry for `/usr/local/bin/sh-agent`, so it prints v1.2.1 even though PATH order would prefer `~/.local/bin`. This is the same shadow-binary class of problem as BUG-078 (which produced `PathShadows` in pathguard.go), but the `update` code path never calls `PathShadows` to warn about it — only `install`'s `pathHint` (main.go:395) and `status` (main.go:648-653) do.

## Top candidate root causes (ranked)

1. **(most likely) `UpdateFromRelease` writes to a fixed `home`-derived path, never to the running binary's actual location, and does not detect/warn when they differ.**
   - Evidence: release.go:162 `dest := filepath.Join(home, BinPath)`; no `os.Executable()` call anywhere in release.go (confirmed via grep); the doc comment at release.go:159-161 explains the atomic-rename mechanics but never addresses *where* `dest` is relative to the invoking binary. `BinPath` is a package constant (plan.go:16), not derived from `os.Executable()`.
   - What would falsify it: if `dest` were computed from `os.Executable()` (or compared against it with a warning) this could not happen — check by re-reading release.go for any executable-path lookup; there is none currently.

2. **`update`'s output never runs `PathShadows` to tell the operator the write didn't land on the resolved binary.**
   - Evidence: `harness.PathShadows` (pathguard.go:28) exists specifically for this class of bug (its own doc comment cites BUG-078: "a downloaded, signature-verified, checksum-verified release binary was installed to /usr/local/bin and then ignored"). It is called from `install`'s `pathHint` (main.go:395) and from `status` (main.go:648-653), but the `update` case in main.go (lines 191-216, calling `harness.UpdateFromRelease`) never calls it. So even though the tooling to detect exactly this shadow exists, `update`'s own success output ("verified and installed ... 260 unchanged") gives no signal that the binary actually on PATH first (or hashed by the shell) was left behind.
   - Falsify: if `update` did call `PathShadows`/print a shadow warning after a successful release replace, the operator would have been told before running `version`. It currently does not — confirmed by grep showing only two call sites, both outside the `update` release branch.

3. **(contributing, not a code bug) bash's command hash table masks PATH order after the shell already resolved `sh-agent` once.** This is standard bash behavior (`hash -r` needed), not something sh-agent can fix directly, but it's what makes candidate 1's effect visible even when `~/.local/bin` precedes `/usr/local/bin` on PATH — the repro explicitly notes "or run `sh-agent update` in a shell that has already resolved `sh-agent` to `/usr/local/bin/sh-agent`". This amplifies but does not cause the core defect; candidate 1 is necessary and sufficient on its own (e.g., a fresh shell with `/usr/local/bin` genuinely first on PATH reproduces it without any hash involved).

## What `update` writes and what it would have needed to say
- Writes new binary to `filepath.Join(home, BinPath)` = `$HOME/.local/bin/sh-agent` (release.go:162), via `writeAtomic` (temp file + rename in the same dir).
- Nothing in `UpdateFromRelease` or the `update` case in main.go compares `dest` to `os.Executable()` of the running process, nor to any PATH-resolved `sh-agent`.
- For the operator to know before running `version`, the "verified and installed" message (release.go:169) or the subsequent install step would have needed to also report (via `PathShadows`, which already exists) that the binary actually invoked (`/usr/local/bin/sh-agent`, resolved by PATH or by the shell's hash) is NOT the file just replaced — i.e., the same warning `pathHint`/`status` already produce, but currently only reachable through `install`/`status`, not through `update`.

## Not examined
- Did not read `internal/harness/plan.go` in full (only grepped `BinPath` constant) or `cmd/sh-agent/main.go` lines 580-660 (`selfExecutable`, `binaryOlder`, status/pathHint helpers) in full detail — only grepped call sites.
- Did not inspect `release_test.go` or `main_test.go` for existing coverage of this exact shadow-on-update scenario (candidate for confirming whether this is already a known gap or genuinely untested).
- Did not confirm bash's exact hash-table semantics against the transcript beyond citing it as standard behavior.

diagnosis complete
