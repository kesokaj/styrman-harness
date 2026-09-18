---
owner: bugs
date: 2026-09-18
id: BUG-078
status: open
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
Not fixed. Operator manually pointed `/opt/devbox/bin/sh-agent` at `/usr/local/bin/sh-agent`; re-install then correctly reported `release binary v1.2.1: installing the harness built into it`. No code or doc change has landed to prevent recurrence.

## Prevention
Candidate rule (not yet written into any file — proposing, not applying): after installing or replacing the `sh-agent` binary, run `which -a sh-agent` and compare every match against the binary just installed; a "verified install" claim is only true once `sh-agent version` (run through the shell, unqualified — not by absolute path) reports the expected version. Separately: a downloaded release artifact is run (`version`) and checksum-verified before being moved into a system path, and the download step uses `--clobber` so a stale file at the download path can't be silently reused.

## Related
Candidate for graduation into `.agents/docs/deploy/remote.md` § release/install procedure (or a `sh-agent install`/`status` PATH-shadow check) — not graduated by this record. Related to BUG-006, BUG-015, BUG-077 (binary/source-selection correctness on install/update).
