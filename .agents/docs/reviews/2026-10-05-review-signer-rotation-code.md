---
owner: reviewer
date: 2026-10-05
---
# Review — signer-rotation-code   (2026-10-05)

Range reviewed: b549d71..a0a54b4

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The worst problem: the release procedure that `/preflight` runs line by line has no bridge path, even though ADR-052 says v2.4.0 is a bridge release. Run as written for v2.4.0, step 3 signs `SHA256SUMS` with the NEW key. Step 3 then verifies it against the new key, and step 3b checks it with the v2.4.0 binary, which trusts the new key. Every check passes, and the result is a release that no v2.3.2 binary-only place can update to, now or later. The second problem: zip places are stranded after the rotation, and nothing says so. Third, the Operator steps in `remote.md` still tell a new machine to check `~/.ssh/id_ed25519`. The new commands themselves are correct. I ran each one with throwaway keys and it works, and each check fails when it should. The new signer line and its fingerprint are the same everywhere they appear.

## Findings

- **[high / confidence: high]** `.agents/docs/deploy/remote.md:106-123,155-163` (§ Releases procedure, steps 0–3b), against `.agents/docs/architecture/DECISIONS.md` ADR-052 bullet 3 and `remote.md:102`. The prose says v2.4.0's `SHA256SUMS.sig` is signed by the retired key. The procedure has no step that does that:
  - The pre-check (`ssh-add -L … grep -qxF '<new key>'`) only checks for the new key.
  - Step 3 writes the new key's public half to `release-signer.pub` and signs with it.
  - Step 3's verify uses the new allowed-signers line.
  - Step 3b runs `./sh-agent-$(go env GOOS)-$(go env GOARCH)`, which is built from the v2.4.0 tag and so carries the new key.

  For v2.4.0, every line exits 0 and certifies a release whose signature every v2.3.2 binary refuses (`release.go` verifies `SHA256SUMS.sig` only against its compiled-in signer). Binary-only places are stranded and need a manual reinstall. `/preflight` follows this block literally ("each command its own call"), so the prose note does not protect the run. ADR-052's "a v2.3.2 binary built from the tag" is also ambiguous: a binary built from the v2.4.0 tag is not a v2.3.2 binary.

  Fix: add an explicit v2.4.0-only bridge block:
  1. Pre-check that the retired key is in the agent.
  2. Sign with the retired key's public half.
  3. Verify against `ostling.simon@gmail.com ssh-ed25519 …BkX…`.
  4. For 3b, use a binary built from `git archive v2.3.2` into `.agents/scratch/` and require the `VERIFIED … release v2.4.0` line.
  5. Add a must-fail check that the v2.4.0 binary refuses the bridge signature.

  Also say where the retired private key is loaded from for this one signing.
- **[medium / confidence: high]** `.agents/docs/architecture/DECISIONS.md` ADR-052 Consequences and `.agents/docs/specs/sh-agent/delivery.md:77` (the v2.4.0 sentence). Zip places are left stranded and nothing records it. Per delivery.md § The model (line 32), `update --from-zip` "installs the unpacked tree with the RUNNING binary" and never replaces it, and a release zip whose signed pair does not verify is refused. A zip place on a v2.3.2 binary therefore accepts the v2.4.0 zip (retired key) but keeps a binary that trusts only the retired key. It then refuses every later release zip, which is signed by the new key alone.

  ADR-052 only says "every place should run `sh-agent update` soon", which does not move a zip place's binary. The delivery spec's v2.4.0 sentence (the one home for what a tag delivers, AGENTS.md rule 6) does not mention the rotation at all.

  Fix:
  - Append a note (ADRs are append-only) and add a delivery.md sentence: a zip place must download the v2.4.0 binary by hand.
  - Add the same to the v2.4.0 release notes.
  - Have the operator decide the version number against the rule in `remote.md:207`: "MAJOR for what an existing deployment must react to".
- **[medium / confidence: high]** `.agents/docs/deploy/remote.md:220` (§ Operator steps). It still says "On a new machine, `ssh-keygen -y -f ~/.ssh/id_ed25519 | ssh-keygen -lf -` must print the fingerprint in § Releases before that machine cuts a release". That contradicts `remote.md:102` ("It is never `~/.ssh/id_ed25519`") and the new `ssh-add -L` pre-check. On a correctly set-up machine the check fails, because `id_ed25519` is a login key. Satisfying it means putting the release key at the default path, which is what caused ADR-052 in the first place. This is a restatement the rotation missed.

  Fix: replace it with the `ssh-add -L | cut -d' ' -f1,2 | grep -qxF '<signer>'` check. Related: `remote.md:205` still lists "the private-key-derived fingerprint" among the "stricter forms above", and that check is no longer in the block. Drop the clause or mark it historical.
- **[low / confidence: high]** `.agents/docs/architecture/DECISIONS.md` ADR-051 Consequences ("the v2.4.0 tag still goes on `caa37d6`") and `.agents/docs/handoffs/2026-10-05-1248-v2.4.0-tag-still-pending.md`. Both are false after a0a54b4. `caa37d6` compiles in the exposed key, and the source change in `release.go` after the gated head `b549d71` means the procedure's `HEAD^ == <gated head>` and records-only `git diff --quiet` checks now need a new gate. The procedure guards against tagging the wrong commit, but ADR-052 does not say that it supersedes that part of ADR-051.

  Fix: one sentence in an appended note saying the v2.4.0 gate is re-run over a range that ends at or after a0a54b4.
- **[low / confidence: high]** `.agents/docs/deploy/remote.md:32` (the release-binary download chain) and `README.md:143` (the Mac/browser chain). Both carry only the new signer line and no v2.4.0 note. The bridge comment exists only above README's first chain (`README.md:103-104`). With `V=v2.4.0` both chains fail closed at `ssh-keygen -Y verify`, so this is safe but confusing.

  Fix: point both at the bridge note, or print the retired line beside them.
- **[low / confidence: medium]** `internal/harness/release_test.go:449` `TestReleaseSignerMatchesTheDeployGuide` pins only `remote.md` lines that start with `ostling.simon@gmail.com ssh-`. These copies of the key are not pinned:
  - the three new signing lines (`ssh-add -L` check `remote.md:106`, `key::` in `:117`, `release-signer.pub` in `:156`)
  - README's two lines (`README.md:109,143`)

  They are all consistent today. I checked every occurrence: 7 in remote.md, 2 in README, 1 in release.go, 1 in DECISIONS.md, all byte-identical. At the next rotation they can drift unseen. The test's own comment says "two copies of a key are two keys the day one of them is rotated".

  Fix: match the base64 key body anywhere in remote.md and README, not only at line start.

## Verified (how, and on which path)

- **Fingerprints.** `printf '<key>' | ssh-keygen -lf -` gives `SHA256:BwBzPnfXKuKrWP6p10IZgcDG8uwSRXJOFv0Izl+79yQ` for the new key, which matches remote.md:102 and ADR-052. The retired key gives `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`, which matches the old guide.
- **Signing commands, rehearsed** with two throwaway ed25519 keys and an `ssh-agent` whose socket was under `.agents/scratch/rv-signer` (since deleted), on OpenSSH 10.0p2 and git 2.47.3 on this Linux box:
  - `ssh-add -L | cut -d' ' -f1,2 | grep -qxF '<pub>'`: exit 0 with the key loaded; exit 1 with an empty agent, with only a different key loaded, and with no agent (`SSH_AUTH_SOCK` unset).
  - `git -c gpg.format=ssh -c user.signingkey='key::<pub>' tag -s`: exit 0 with the key in the agent; exit 128 ("Couldn't find key in agent?") with the key absent, and no tag created. `verify-tag` against an allowed-signers line built the procedure's way says `Good "git" signature for ostling.simon@gmail.com`.
  - `ssh-keygen -Y sign -f <pubfile> -n sh-agent-release SHA256SUMS`: exit 0 and writes `.sig` with the key in the agent; exit 255 ("No private key found") and no `.sig` with the key absent.
- **Go change.** `gofmt -l internal/harness/release.go` is empty, `go vet ./internal/harness/` passes, and `go test -count=1 -run 'TestReleaseSignerMatchesTheDeployGuide|TestChecksumFor' ./internal/harness/` is ok. I did not re-run the must-fail mutation that the change entry claims. The test's logic (exact line equality) fails on the old line by inspection.
- **Old-key search.** `git grep` at a0a54b4 for the old key, the old fingerprint and `id_ed25519` found, outside the history records:
  - the bridge note (`remote.md:102`, intended)
  - ADR-026/052 text (append-only history)
  - `remote.md:220` (finding above)
  - unrelated secret-name lists (`plan.go:222`, `harness_test.go`, `.gitignore`, `skills/docker`)

  There are no stale copies in `context/`, `skills/preflight`, or `agents/`.

## Criteria

1. **Dead code: pass.** The diff only swaps a constant's value and adds a comment. Nothing is orphaned.
2. **Flow tracing: fail.** The new signer flows correctly from `releaseSigner` to `verifyRelease` to the guide pin. The release flow does not: the bridge step ADR-052 requires is in no executable step, and the 3b self-check would pass a mis-signed bridge (finding 1). The zip place's binary is never replaced, so it loses the update path (finding 2).
3. **Standards conformance: fail.** The rotation missed restatements at `remote.md:220` and `:205` (finding 3), and the ADR-051 consequence is stale (finding 4). The `release.go` comment is accurate. Go style is clean.
4. **Directory layout: pass.** No new files. The edits sit in the files that own each claim.
5. **Simplicity (DIY-except-security): pass.** Signing and verification stay in `ssh-keygen` and `git` as designed, and nothing is hand-rolled. Signing through the agent by public key removes the private-key path, which is the simpler and safer shape.
6. **Decoupling: pass.** The trust root is still one compiled-in variable, and the test swap seam is unchanged.

## Not examined

- The records in the range, read only through grep:
  - `handoffs/2026-10-05-1242-*`, `handoffs/2026-10-05-1248-*`
  - `reviews/2026-10-05-preflight-v2.4.0.md`, `reviews/2026-10-05-review-v2.3.2-0c0846c-{code,docs}.md`
  - `testing/2026-10-05-verify-v2.3.3-gate{,-rows}.md`
  - `changes/2026-10-05-1238-*`, `changes/2026-10-05-1248-*`
- The untracked `changes/2026-10-05-1317-release-signer-rotation.md`, read for context only (outside the range).
- `reviews/2026-10-05-security-v2.4.0-signer-rotation.md`, not read.
- macOS: OpenSSH and git behaviour for the `key::` signing and the agent-held `-Y sign` was not rehearsed there.
- The real release key and the retired key: neither was used. The rehearsal used throwaway keys only.
- The full test suite: not run. Only the signer pin and checksum tests were run.
- Whether ADR-051's REQUIREMENTS edit also needed a ROADMAP.md change in the same commit: not assessed.
