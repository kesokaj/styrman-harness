---
owner: reviewer
date: 2026-10-05
---
# Review — kit-credential-failclosed round 2   (2026-10-05)

Range reviewed: 0eb8320..29cd652

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Nothing is broken. All three round 1 findings that could hide a credential (a link to a directory, a relative or vanished `KIT_HOME`, the pruned top-level `bin/`) are fixed, and I reproduced each fix in scratch. The exit trap still reaches its verdict under `set -euo pipefail`. No pipe feeds an early-exiting reader (BUG-155). The remaining nits are coverage gaps: the new absolute-`KIT_HOME` refusal and the dangling-link report have no test, and the changelog says "two test cases added" for four behaviour changes. Reproductions are in `.agents/scratch/rv2` (function replay, with a mutated copy) and in a scratchpad directory outside the checkout (the real preamble, because it refuses a clone inside a checkout).

## Round 1 findings, one by one

1. **Directory link reads clean (medium): FIXED.** Evidence:
   - `kit-preamble.sh:119` reports `(a link to a directory, not walked)`.
   - Replay of the exact test setup, `.ssh` -> `../kit-outside` holding `id_ed25519`: current function gives rc=1 and names `.../home/.ssh`; the same function with the line deleted gives rc=0 and "no credentials left". The new test therefore fails without the fix.
   - The guide (`guide.md:162`) and the preamble comment (lines 100-103) now say "link to a file read through, link to a directory reported, never walked", which matches the code.
2. **Relative or vanished `KIT_HOME` reads clean (low): FIXED for relative, PARTLY for vanished.**
   - The real preamble now exits 2 with "KIT_HOME must be an absolute path, got 'rel'" before anything else runs (reproduced, rc=2).
   - A dangling-link `KIT_HOME` now fails with `(a dangling link)` (reproduced, rc=1).
   - A `KIT_HOME` that was deleted outright still reports "no credentials left", rc=0 (reproduced). Nothing is left there to leak, so I accept it; round 1 asked for "gone", and this is the one place the fix is narrower than asked.
3. **Top-level `bin/` pruned in full (low): FIXED.** The prune is gone. The real preamble, with a real `go build` of `sh-agent.real` (ELF) plus the bash wrapper, reads clean (rc=0, no false positive). A PEM header in `$KIT_HOME/bin/stray` now gives rc=1 (reproduced; the new Go case passes). The exemption is now by ELF/Mach-O magic only, like everywhere else.
4. **No case for the `(could not be scanned)` branch (low): OPEN.** The diff adds none, and grep exit 2 on a readable file is hard to produce by hand. The `head -5` cap that round 1 also named is gone with the pipe, so only the branch is left. A coverage gap, not a defect; the changelog should not claim coverage of it.
5. **Name list gaps, `oauth_creds.json`, `*.pem`, `.npmrc` (low): OPEN, unchanged.** Real ones are still caught by content (`"refresh_token"`, PEM header). Only a token-less fixture of those names passes. Not blocking.

## Findings
- **[low / confidence: high]** `internal/harness/catalog_test.go:1597-1606` — no test covers the new `KIT_HOME must be an absolute path` exit 2 (`kit-preamble.sh:33`) or the `(a dangling link)` report (line 109). The Go test passes whether or not those two lines exist, so the changelog's "against 0eb8320 exactly those two fail" is true of the two added cases but silent on these two behaviours. Fix: one case that runs the preamble with `KIT_HOME=rel` and expects exit 2, and one with a dangling-link `KIT_HOME` expecting exit 1 and the FAIL line.
- **[low / confidence: medium]** `.agents/docs/testing/kit-preamble.sh:119` — every directory link under `KIT_HOME` is now a hard FAIL, with no allowance for one a CLI creates by design. I checked only the installer: a scratch `sh-agent install --home` holds no directory links and reads clean. A real Claude Code or Antigravity run was not examined. If one creates such a link, the cost is a failed run (fail closed, safe) that needs a named exemption, not a hole.
- **[low / confidence: high]** `.agents/docs/changes/2026-10-05-1759-kit-credential-check-fails-closed.md` (Round 2) — "the real sandbox homes still read clean" was not reproduced by me beyond one scratch install and one real-preamble run (both clean); the claim cites no log. Read the verifier log for the evidence.

## Preamble under `set -euo pipefail` (checked, not read)
- The function and trap were replayed verbatim with stubs for the fingerprint, under `set -euo pipefail`, and the real preamble was run for the verdicts below.
- Verdicts reached, in every case, with the right exit code: clean (rc=0), kit `exit 7` on a clean home (rc=7, propagated), kit `false` under errexit on a clean home (rc=1), dangling link plus `.netrc` plus a fifo (rc=1, both listed), unreadable directory (rc=1, "find could not read" with a path only).
- `unreadable=$(find ... 2>&1 >/dev/null) || true`: redirect order is right (stderr captured, stdout dropped), and the capture is whole, no pipe.
- Remaining pipes: `head -c 4 | od | tr` (head stops by byte count after od is already reading to EOF; `od` does not exit early) and `printf | sed`. Neither feeds an early-exiting reader. `grep -q` reads a file, not a pipe.
- The function's own return status cannot kill the trap: the `while` loop's status is its last body command, which is `continue`, `printf` or `case`, all 0; an empty loop is 0.

## Criteria
1. Dead code — pass. The prune and the `head -5` cap are removed with their uses; nothing left over.
2. Flow tracing — pass. Trap flow, exit-code propagation, errexit and pipefail traced and replayed; the two new test cases discriminate (dirlink: replay of the old behaviour gives rc=0; own-bin: the old `-prune` by construction).
3. Standards conformance — pass. Guide line and preamble comment match the code. The one restatement I grepped for (`symlink`, `top-level bin`) has no stale copy left. Test coverage is claimed beyond what is tested (findings 1 and 3).
4. Directory layout — pass.
5. Simplicity (DIY-except-security) — pass. Coreutils only; detects tokens, handles none.
6. Decoupling — pass. One function, one caller.

## Not examined
- A real Claude Code or Antigravity run through the preamble (whether either creates a directory link in its home).
- macOS bash 3.2 and BSD find/grep/od (Linux only here); `[[ =~ $var ]]`, `read -d ''` and `grep -aqE --` read, not run on it.
- `testing/2026-10-05-verify-kit-credential-check.md` and BUG-160's record were not read in detail; the verifier log was not re-checked against the commit.
- The full `go test ./...` was not run; only `TestKitPreambleRefusesLeftoverCredentials` (14 cases plus the clean, public-key, ELF and symlinked-`KIT_HOME` checks, all PASS).
