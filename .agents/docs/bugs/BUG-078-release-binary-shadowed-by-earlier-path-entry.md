---
owner: bugs
date: 2026-09-18
id: BUG-078
status: fixed
area: installer/binary (PATH resolution, release install procedure)
---
# BUG-078: verified v1.2.1 release binary installed to /usr/local/bin, but an earlier PATH entry kept running a stale (devel) build

## Symptom
On devbox: operator verified `SHA256SUMS.sig`, checksum, and `./sh-agent-linux-amd64 version` printed `sh-agent v1.2.1`, then `sudo mv -f` to `/usr/local/bin/sh-agent`, `hash -r`. `sh-agent version` then printed `sh-agent (devel)`. `sh-agent install --vendor all` / `status` reported:
```
plain copy, not its own git checkout: the running binary is installed as is
repo       /opt/devbox/deploy/harness (plain copy, not its own git checkout; update is refused)
binary     (devel)
```
`claude -p` quoted stale (pre-v1.2.0) context content, proving the wrong tree was installed.

## Repro
Not reproduced in this session (measured live on devbox by the operator). `which -a sh-agent` after the move showed three binaries on PATH:
```
/home/admin/.local/bin/sh-agent   3784864 bytes, 10:21  (earlier PATH entry, wrote the bad install)
/opt/devbox/bin/sh-agent          (devel)               (earlier PATH entry, took over after #1 was deleted)
/usr/local/bin/sh-agent           7052333 bytes, v1.2.1  (the verified release binary, PATH position 11)
```

## Root cause
The install procedure verified the *downloaded file* (signature, checksum, `./sh-agent-linux-amd64 version`) but never verified the *binary the shell resolves* after the move. `/usr/local/bin` was not first on `PATH`; two other `sh-agent` binaries — a stale `~/.local/bin/sh-agent` and a `(devel)` build at `/opt/devbox/bin/sh-agent` from the provisioning tree — shadowed it, so unqualified `sh-agent` commands kept running the old binary and installed the stale plain copy at `/opt/devbox/deploy/harness` that it resolved from disk — a development build carries no embedded harness — instead of the files built into v1.2.1.
`Diagnosis missed:` the first instruction block treated "the file was moved into place" as "the binary is installed" (no post-move verification at all, and no `--clobber` on the download, so a pre-existing stale file at the download path was moved into the system path). The second attempt fixed the download-time check but still treated "`~/.local/bin` is where install writes" as the only candidate shadow, so deleting that one file was believed sufficient — `which -a` was not run until `/opt/devbox/bin/sh-agent` was also found to be shadowing.

## Fix
`52574a0` adds `internal/harness/pathguard.go` `PathShadows`: `install` and `status` warn when an earlier PATH entry holds a different file than the one just installed. `8ce58a0` (review round 1) made identity `os.SameFile`, so a symlink or hard link to the installed binary is not a shadow; `b029ea1` (round 2) deleted the lexical directory comparison that had decided "on PATH" beside it, after two review rounds each found a directory form on which the two disagreed; `28e4747` (round 3) made `install --dry-run` print no PATH claim. README § Install without cloning and `remote.md` § First install now verify the downloaded file before the move and the shell's resolution (`sh-agent version` unqualified, `which -a`) after it. Each round's findings and their disposition are in `.agents/docs/reviews/2026-09-18-{review,security}-4cdc0b5-*.md` `## Outcome` sections.

This fixes only the ordering PATH-shadow: a system or Homebrew directory listed before `~/.local/bin`. It does NOT catch the ordering that actually happened on devbox — the wrong (stale) binary ran the install itself and wrote its output first on its own PATH (`~/.local/bin`), so from that binary's own point of view nothing was shadowed; the code-side guard cannot see a binary it isn't running as. That case is caught only by the doc-side verify step (unqualified `sh-agent version`, `which -a` after the move), not by code.
`Diagnosis missed:` the first fix (before review) decided "on PATH" by comparing directory names alongside a separate file-identity check; two review rounds each surfaced a directory form (`.local/bin//`, a symlinked directory) on which the two checks disagreed. Rather than patch a third form, the operator was asked what the directory-comparison mechanism was for and it was deleted — round 2 made identity-reached-by-the-walk the sole definition of "on PATH."

## Prevention
Candidate rule (not yet written into any file — proposing, not applying): after installing or replacing the `sh-agent` binary, run `which -a sh-agent` and compare every match against the binary just installed; a "verified install" claim is only true once `sh-agent version` (run through the shell, unqualified — not by absolute path) reports the expected version. Separately: a downloaded release artifact is run (`version`) and checksum-verified before being moved into a system path, and the download step uses `--clobber` so a stale file at the download path can't be silently reused.

## Related
Graduated: the prevention rule now lives in `context/AGENTS.md` § Hard rules, "A claim that something works states HOW it was verified" bullet, sentence "An installed COMMAND is verified by the shell that will run it — the bare name, in that shell, after `hash -r` (bash) or `rehash` (zsh), then `type <name>` when it disagrees — never by the file that was moved into place" (BUG-078, BUG-079 named there). Related to BUG-006, BUG-015, BUG-077 (binary/source-selection correctness on install/update).
