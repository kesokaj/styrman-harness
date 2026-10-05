---
owner: reviewer
date: 2026-10-05
---
# Review — signer-rotation-code-r2   (2026-10-05)

Range reviewed: a0a54b4..41af344

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The worst problem: bridge step 3b cannot run as written. `git archive v2.3.2` runs from `.agents/scratch/release-v2.4.0`, and from a subdirectory `git archive` archives only that subdirectory of the tree. v2.3.2's tree has no `.agents/scratch/release-v2.4.0`, so the archive is empty: 0 entries, exit 0. `tar` then succeeds with nothing to unpack, and `go build` fails with `stat …/v232src/cmd/sh-agent: directory not found` (rehearsed, below). The failure is closed, so nothing mis-signed ships. But by step 3b the v2.4.0 tag and `main` are already pushed, so `/preflight` stalls with a published tag that has no assets, and the operator has to improvise the fix in the middle of a release.

Second: `remote.md:103` (the signer paragraph, edited in this range) still tells a reader to "verify v2.4.0 against that key", meaning the retired, exposed key. ADR-052a, the install block at `remote.md:28` and README all say a human verifies `SHA256SUMS.current.sig` against the current key.

The rest of the bridge is correct. Bridge step 3, rehearsed with throwaway keys, signs, verifies, cross-refuses and removes the retired key from the agent exactly as written. The new test fails on a third `ssh-ed25519` key. Every finding in round 1 is fixed except the parts listed under Round 1 status below.

## Findings

- **[medium / confidence: high]** `.agents/docs/deploy/remote.md:197` (bridge 3b, `mkdir -p ../v232src && git archive v2.3.2 | tar -x -C ../v232src`). The command runs from `.agents/scratch/release-v2.4.0`, the cwd that step 2's `cd` set. `git archive` with no path archives only the cwd's subtree, so the result here is empty.
  - **Rehearsed:** from a subdirectory under `.agents/scratch/`, `git archive v2.3.2 | tar -t | wc -l` printed `0` and git exited 0. The build line as written then printed `stat …/v232src/cmd/sh-agent: directory not found` and exited 1.
  - **Fix:** `git -C "$(git rev-parse --show-toplevel)" archive v2.3.2 | tar -x -C ../v232src`. Rehearsed: the build exits 0, and `../sh-agent-v2.3.2 version` prints `sh-agent v2.3.2`.
  - **Why it matters:** the tag and the pushes come before the binary block, so a 3b that cannot run leaves a pushed tag with no assets ("a marker no place can update to", `delivery.md` § The model).
- **[medium / confidence: high]** `.agents/docs/deploy/remote.md:103`, last sentence. It reads: "its `SHA256SUMS.sig` is signed by the retired key, `ssh-ed25519 …BkX…`, so verify v2.4.0 against that key and every later release against the line here." That contradicts four other places:
  - ADR-052a: "`SHA256SUMS.current.sig` by the current key, the one every human check and README's v2.4.0 path verify, because the exposed key alone proves nothing now".
  - `remote.md:28`, the install block's v2.4.0 comment.
  - `README.md:103-104`.
  - `README.md:139-140`.

  A reader of § Releases who follows this sentence checks v2.4.0 against a key the security review calls proof of nothing. This range edited the paragraph (it added the `ssh-add -t 900` text) and left this clause standing.
  - **Fix:** "v2.4.0 also carries `SHA256SUMS.current.sig` by the current key; verify that file against the line here (ADR-052a). Its `SHA256SUMS.sig` by the retired key is for v2.3.2-and-older binaries only."
- **[low / confidence: high]** `internal/harness/release_test.go:478-503` `TestEveryDocumentedReleaseKeyIsKnown`. These gaps were all measured in a scratch copy of 41af344:
  - **Retired key accepted anywhere.** The test accepts the retired key at any position in either file. I swapped the retired key into the NORMAL step 3's `release-signer.pub` line (`remote.md:157`). Both this test and `TestReleaseSignerMatchesTheDeployGuide` still passed. That is round 1's high failure in reverse: a later release signed with the exposed key, with the tests green.
  - **Non-ed25519 keys missed.** The regexp sees only `ssh-ed25519 ` keys. An appended `sk-ssh-ed25519@openssh.com AAAA…` key, which is the hardware-key form `remote.md:103` recommends for the release key, passed. At a rotation to a hardware key, `strings.Cut(releaseSigner, " ssh-ed25519 ")` also turns the test into a Fatalf, so it does fail loudly there.
  - **Two files only.** It scans only `remote.md` and `README.md`. Today no other living file carries a key: a `git grep` found only `DECISIONS.md`, which is history, and `release.go`.

  ADR-052a's "refuses any key in `remote.md` or README that is neither" is therefore wider than the test, which refuses only ed25519 keys.
  - **Fix:** allow the retired key only between the bridge heading and the closing fence of its block, plus the `remote.md:103` sentence. Widen the regexp to `(?:sk-)?ssh-[a-z0-9-]+(?:@openssh\.com)? (AAAA[0-9A-Za-z+/=]+)`.
- **[low / confidence: medium]** `.agents/docs/deploy/remote.md:194` (`ssh-add -d retired-signer.pub && rm -f retired-signer.pub   # … delete its private file too`). Nothing in the block or its prose says where the retired private key is loaded from. Round 1 asked for this. The only deletion is a comment with no command. The load is fail-closed: line 189's `ssh-add -L … grep` exits 1 when the key is absent. The deletion is the step that `/preflight`, running commands, skips silently. ADR-052 Consequences names "this box's scratch copy" as a location.
  - **Fix:** name the file, and make the deletion a command (`rm -f <that path>`) followed by a `test ! -e` check.
- **[low / confidence: medium]** `.agents/docs/deploy/remote.md:201` (`! "./sh-agent-…" update --from-zip … > ../zipcheck-new.out 2>&1`). The negation passes on ANY non-zero exit (a flag error, an unwritable home, a missing file), not only on the signature refusal it claims to prove. Line 203-204 limits the damage, because the same binary must then print VERIFIED with the current signature. A refusal for an unrelated reason would still be read as "refuses the retired key".
  - **Fix:** keep the `!` line, then add `grep -q '<the refusal text release.go prints on a signature mismatch>' ../zipcheck-new.out`.
- **[low / confidence: high]** `.agents/docs/deploy/remote.md:177` (`# all six`) and `:216` ("lists the six assets"). For v2.4.0 the bridge prose adds a seventh asset. The upload change itself is stated in the bridge prose (`:181`, "step 4's upload list gains `SHA256SUMS.current.sig`"), but the loop at `:174` and these two counts are not marked. A literal run of step 4 uploads six assets and the listing check reads as satisfied, which leaves README's v2.4.0 path (`--pattern 'SHA256SUMS.current.sig'`) failing for every reader.
  - **Fix:** give the bridge block its own step-4 loop line, or add "(seven for v2.4.0)" to both counts.

## Round 1 status

From `reviews/2026-10-05-review-signer-rotation-code.md`:

1. **[high] No bridge in the procedure: fixed, with defects.** `remote.md:181-205` adds the bridge block. It carries:
   - a retired-key pre-check (`:189`)
   - the retired key's signature (`:190-191`)
   - a verify against the retired line (`:192-193`)
   - a v2.3.2 binary for 3b
   - a must-refuse check on the v2.4.0 binary

   I rehearsed step 3 with throwaway keys. It also adds a current-key signature, which goes past the original ask.

   Defects:
   - The 3b build cannot run as written (finding 1).
   - "Where the retired key is loaded from" is still open (finding 4).
2. **[medium] Zip places stranded: fixed.** The fix is in `delivery.md:32` (the zip edge and the binary-only edge, stated once, in rule 6's one home) and in ADR-052a's Context and Consequences. ADR-052a Consequences records the operator's version decision. The release-notes part could not be checked in this range: notes are built at release time.
3. **[medium] `remote.md:220` `id_ed25519` check: fixed.** It is now `remote.md:249`, "`ssh-add -l` must list the fingerprint …". The `:205` sub-point is also fixed: `remote.md:234` now marks the rehearsal as made "with the key then read from a file (the agent form above replaced it, ADR-052)".
4. **[low] ADR-051 consequence and handoff stale: open.** Nothing in the range says the v2.4.0 gate is re-run over a range ending at or after a0a54b4. The untracked `testing/2026-10-05-1317-verify-v2.4.0-gate.md` suggests a re-gate is under way, but I did not read it.
5. **[low] `remote.md:32` and README `:143` had no bridge pointer: fixed.** The pointers are at `remote.md:28` and `README.md:139-140`. But see finding 2: `remote.md:103` now disagrees with them.
6. **[low] Signer pin matched only lines starting with the email: fixed in part.** `TestEveryDocumentedReleaseKeyIsKnown` now matches every `ssh-ed25519` key in both files. Its gaps are in finding 3.

From `reviews/2026-10-05-review-b549d71-a0a54b4-docs.md`:

1. **[high] Bridge missing in steps 3/3b: fixed, with defects.** Same as code item 1. Step 3's verify line, the signing key and the 3b binary are each now named.
2. **[medium] `remote.md:220`: fixed.** Same as code item 3.
3. **[low] `remote.md:205` "private-key-derived fingerprint": fixed.** Marked historical (`remote.md:234`).
4. **[low] README "compiled into every sh-agent": fixed.** It now reads "from v2.4.0" (`README.md:100`).
5. **[low] README Mac block and `remote.md` install block: fixed.** The pointers are at `README.md:139-140` and `remote.md:28`. `remote.md:26`'s "§ Releases' allowed_signers line" is now right for the v2.4.0 path too, since that path verifies the current key.
6. **[low] Test coverage claim: fixed in part.** See finding 3.

## Agreement check: README, remote.md install block, delivery.md, ADR-052a

- **Agree.** These all say a human verifies v2.4.0 by `SHA256SUMS.current.sig` against the current key:
  - `README.md:103-104` (the gh chain)
  - `README.md:139-140` (the browser chain)
  - `remote.md:28` (the install block)
  - ADR-052a Decision
  - `delivery.md:32` ("verified by hand against the current key's `SHA256SUMS.current.sig` (README § Install without cloning)")
- **Disagrees.** `remote.md:103` says to verify v2.4.0 against the retired key (finding 2). ADR-052's own Consequences says the same ("uses the retired key's line"). As an append-only record, ADR-052a amends it generally ("amends ADR-052") without naming that sentence, which is acceptable.
- **Agree.** The binary-only edge reads the same in `delivery.md:32` and ADR-052a Context, and `README.md:167` agrees that `update` takes only the newest release.
- **Agree.** The bridge block's description of what 3b proves matches ADR-052a Decision line for line. The one exception is "removed from disk" (finding 4).

## Verified (how, and on which path)

- **Bridge step 3, every line as written.** Rehearsed with two throwaway ed25519 keys substituted for the real ones, in an `ssh-agent` whose socket was under `.agents/scratch/`, on this Linux box. Results:
  - Current-key sign plus `mv`: 0.
  - Verify of `.current.sig`: 0.
  - Retired key in agent: 0.
  - Retired sign: 0.
  - Retired verify: 0.
  - `ssh-add -d retired-signer.pub`: 0, "Identity removed". Afterwards the retired-key check exits 1 and the current-key check exits 0.
  - Must-fail cross-check: `.current.sig` against the retired allowed_signers exits 255.

  Scratch removed after the run.
- **Bridge step 3b build.** The as-written and corrected forms were run in `.agents/scratch/` (results in finding 1). Against the v2.3.2 tag:
  - The tag exists locally.
  - v2.3.2 has `--from-zip`, `--dry-run` and `--home` (`cmd/sh-agent/main.go:120,172,204`).
  - It prints the exact `VERIFIED against the signed SHA256SUMS of release ` string (`main.go:1161`).
  - Its `releaseSigner` is the retired key (`release.go:46`).

  So lines 199-200 grep the right text once the build works. I did not run 3b end to end, since a signed v2.4.0 zip needs a real key.
- **`TestEveryDocumentedReleaseKeyIsKnown`.** Run in a `git archive 41af344` copy under `.agents/scratch/`. Must-pass: `ok` with `TestReleaseSignerMatchesTheDeployGuide`. Must-fail: a third `ssh-ed25519 AAAA…` key appended to README fails at `release_test.go:497` with the expected message. The mutated code compiled and the named test failed. The two passing mutations are in finding 3. Each mutation was restored from a saved copy. `gofmt -l` is empty and `go vet ./internal/harness/` passes on the copy.

## Criteria

1. **Dead code: pass.** `retiredReleaseKey` is used by the new test. Nothing is orphaned. The bridge block calls itself history after v2.4.0, which is intended.
2. **Flow tracing: fail.** The bridge's 3b build cannot run from the release directory (finding 1). The must-refuse check is not tied to its cause (finding 5). The step-4 asset count is not adjusted (finding 6). Step 3's signing flow is correct (rehearsed).
3. **Standards conformance: fail.**
   - The restatement at `remote.md:103` contradicts ADR-052a (finding 2).
   - The ADR-052a coverage claim is wider than the test (finding 3).
   - The Go test is idiomatic, and gofmt and vet are clean.
   - Rule 6: `delivery.md` holds the edges once, and the guides point at it or give how-to only.
4. **Directory layout: pass.** No new files. Each edit is in the file that owns the claim.
5. **Simplicity (DIY-except-security): pass.** Signing and verification stay in `ssh-keygen`, `git` and the existing binary. Nothing cryptographic is hand-rolled. The second signature reuses the same tool path.
6. **Decoupling: pass.** The trust root is still one compiled-in var. The test reads docs through relative paths, like the existing pin.

## Not examined

- The full suite (`go test -race ./...`). Only the two signer tests ran, plus gofmt and vet on `internal/harness`.
- Bridge 3b end to end. It needs a zip signed by the real retired key, and I used neither real key.
- The refusal message the v2.4.0 binary prints on a signature mismatch. `release.go`/`main.go` at 41af344 were not read for it (finding 5's fix needs it).
- macOS behaviour of the bridge commands.
- The untracked records: changes `2026-10-05-1317-*`, the security review, both v2.4.0 gate logs, `BUG-160` and the modified `bugs/INDEX.md`. These appeared during this review and were not written by me.
- `skills/preflight/SKILL.md` beyond a grep for `SHA256SUMS`, `six`, `bridge` and `signer`, which found only its generic Step 5 artifacts line.
