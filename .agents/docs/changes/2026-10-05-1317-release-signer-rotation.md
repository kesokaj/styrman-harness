---
owner: main
date: 2026-10-05
---
# Release signer rotated; the release key signs through the SSH agent (ADR-052)

## Summary

`a0a54b4`. The previous release key's private half was pasted into a chat session (exposed; no known misuse), and it doubled as dave's GitHub login key. New signer `SHA256:BwBzPnfX…` compiled into `internal/harness/release.go` and printed in `deploy/remote.md` § Releases and README § Install without cloning (with the v2.4.0 bridge note). The procedure no longer reads `~/.ssh/id_ed25519`: it checks the key is in the SSH agent (`ssh-add -L`), signs the tag with `user.signingkey='key::<public key>'` and `SHA256SUMS` with `ssh-keygen -Y sign -f <public key file>`. Project `AGENTS.md` § Release says the same. v2.4.0 is the bridge: tag by the new key, `SHA256SUMS.sig` by the retired one, so v2.3.2 places accept it.

## Review fixes (41af344, 48ae45c, 18b9018)

The code, docs and security reviews of a0a54b4 found the procedure ran the bridge with the wrong key. `41af344`: a v2.4.0-only bridge block in `remote.md` § Releases (two signatures over one `SHA256SUMS`: `SHA256SUMS.current.sig` by the current key, `SHA256SUMS.sig` by the retired key; step 3b with a v2.3.2 binary built from its tag and with the v2.4.0 binary both ways), the zip-place and binary-only edges in `delivery.md` § The model, ADR-052a, README's v2.4.0 path on `SHA256SUMS.current.sig`, release-key agent hygiene, and `TestEveryDocumentedReleaseKeyIsKnown`. `48ae45c`: `git -C <toplevel> archive` (a bare archive from the scratch subdirectory packs nothing), explicit retired-key load and delete, `remote.md:103` points at the current key, the test confines the retired key to the bridge block and reads every key type. `18b9018`: fingerprint check of the retired key file before load or delete, an agent-must-answer check, a tab before a key. Credential copies in scratch sandbox homes deleted by name (BUG-160, open for the kit fix).

## Verification

`TestReleaseSignerMatchesTheDeployGuide` passes, and fails with the old line restored in `release.go` (mutation run). Agent signing checked with a throwaway key on OpenSSH 10.0p2 / git 2.47.3: `ssh-keygen -Y sign -f <pub>` through the agent, then `-Y verify`, and the `ssh-add -L` check, all exit 0. Verifier PASS at 18b9018 (`testing/2026-10-05-1425-verify-v2.4.0-gate-r4.md`); main-agent rows and every mutation run in `testing/2026-10-05-verify-v2.4.0-gate-rows-r2.md`. Reviews: code r1 and r2, docs r1, security r1, r2 and r3 under `reviews/2026-10-05-*signer-rotation*` and `-b549d71-a0a54b4-docs`; the 18b9018 diff (4 lines of docs, 1 of test) read by the main agent.

## Not covered

The new key's private file at `~/.ssh/styrman-release_ed25519` has no passphrase yet (the operator sets it). The operator deleted the `dave` GitHub key (161370382) at ~13:27Z; the same key was re-added at 13:29:05Z (165439402), by nothing on this box; the operator is removing it and what re-adds it on dave.
