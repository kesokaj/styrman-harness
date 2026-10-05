---
owner: reviewer
date: 2026-10-05
---
# Review — kit-credential-failclosed   (2026-10-05)

Range reviewed: 218e301..0eb8320

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The rewrite does what it says for the cases it names, and I found nothing that ends the exit trap early or prints secret text. One claim is wider than the code. The new text in `guide.md` and the preamble comment say symlinks are followed, but a symlink to a directory is never walked, so a credential reachable only through one still gives "no credentials left in KIT_HOME" (reproduced below). Nothing in REQUIREMENTS.md or a recorded bug is violated. The scratch reproduction is in `.agents/scratch/rv`; the Go test itself was not run, only its seven new cases and the ELF and symlinked-KIT_HOME cases, replayed as bash against the old (218e301) and new preamble.

## Findings
- **[medium / confidence: high]** `.agents/docs/testing/kit-preamble.sh:98,115-123` and `.agents/docs/testing/guide.md:162` — the check can say "no credentials left" while one is present. The scan uses `find -H` with `-type f -o -type l`, so a symlink that points at a DIRECTORY is not entered. Repro: `ln -s ../realdir/x "$KIT_HOME/.ssh"`, with `id_rsa` holding a PEM header in `realdir/x`. New and old preamble both exit 0 with "no credentials left in KIT_HOME". The guide now says the check "follows symlinks", and the preamble comment says "Symlinks are followed", so a kit author will read that as covered. Fix, either of: (a) report every symlink whose target is a directory as `(directory link, not scanned)`, which is fail-closed and cannot walk `/`; or (b) narrow both claims to "a link to a file is read through; a link to a directory is not walked". Add a test case either way. The fingerprint's own comment already states why `-L` is not used.
- **[low / confidence: high]** `.agents/docs/testing/kit-preamble.sh:104` — `[ -e "$KIT_HOME" ] || return 0` answers "no credentials" when the path does not resolve from the trap's cwd. Repro: relative `KIT_HOME=rel-home`, the kit writes `.netrc`, then `cd /`. The result is exit 0 and "no credentials left in KIT_HOME". The old `[ -d ]` did the same, but the change's claim is "fails closed", and the preamble does not require an absolute `KIT_HOME` (line 30-31 only refuse the real home). A dangling symlink as `KIT_HOME` takes the same exit. Fix: refuse a non-absolute `KIT_HOME` at the top with exit 2, and let the check report a `KIT_HOME` that was set but is gone as "gone" rather than clean.
- **[low / confidence: medium]** `.agents/docs/testing/kit-preamble.sh:98,112` — the top-level `$KIT_HOME/bin` is pruned in full, not just the wrapper and `sh-agent.real`. A credential file placed there is invisible (repro: a PEM-header file at `$KIT_HOME/bin/leak` gives exit 0). The comment states it ("only the kit's own top-level bin/ is skipped"), so this is a documented bound. The two files the preamble creates can be skipped by name instead, which keeps the exemption to what the preamble wrote.
- **[low / confidence: medium]** `internal/harness/catalog_test.go:1597-1606` — no case covers the new `(could not be scanned)` branch (grep exit 2 on a readable file) or the `head -5` cap on find errors. Both were exercised by hand only (8 unreadable dirs print 5 lines; a 200 MB line scans in about 1 s). The new branch is untested code on the fail-closed path. This is a coverage gap, not a defect.
- **[low / confidence: low]** `.agents/docs/testing/kit-preamble.sh:99` — name list gaps: `oauth_creds.json` (Gemini CLI), `*.pem`, `.npmrc` are matched only by content. An empty or token-less fixture of those names passes. Content patterns catch real ones (`"refresh_token"`, PEM header), so this is only worth a line if the list is meant to be complete.

## What I verified (reproduced, not read)
- Exit trap under `set -euo pipefail`: an unreadable dir, an unreadable file, a dangling link, a fifo and 8 unreadable dirs all reach the verdict and print the FAIL block with exit 1. The failing `find | head` in `$()` does not abort the trap (command substitution does not inherit errexit in bash without `inherit_errexit`).
- Path handling: a name with a space and an embedded newline is found through `-print0` and printed as paths only. No case printed token text: grep stderr is dropped, and find's stderr carries paths only.
- A realistic home: `sh-agent install --home` into scratch (1790 files, including `.local/bin/sh-agent`) reports "no credentials left" with old and new preamble. No tracked file under `context/ skills/ agents/ rules/ commands/` matches `kit_cred_text`, so installed skills do not trigger false positives.
- Test cases against the old preamble (bash replay of each case): symlinked SSH key, key under nested bin, token in a binary file, `ya29.` token, `.netrc`, unreadable file and unreadable dir all give exit 0 on 218e301 and exit 1 on 0eb8320, so each fails without the fix. A symlinked `KIT_HOME` is also 0 then 1.
- ELF case: the old preamble also gives exit 0, because the file sits in a nested `bin/`, so it does not discriminate old from new. It is a pin for the new magic skip. Deleting the `case "$magic"` line gives exit 1, so it does fail without the skip. Mutation: the sed changed exactly the one line, and the mutated preamble ran (rc=1) rather than failing to build.
- The `grep -a` over a file with the token padded to 29 characters gives no match (correct, below the {30,} floor).
- The new tests' `chmod -R u+rwx` before `RemoveAll` is needed; the root skip is correct for the two unreadable cases.

## Criteria
1. Dead code — pass. Nothing left over; `kit_cred_paths` and all three variables are used.
2. Flow tracing — pass with the two gaps above (directory links; relative or vanished `KIT_HOME`). Trap flow, errexit, pipefail and SIGPIPE on `head -5` all behave. Not verified on bash 3.2 (see Not examined).
3. Standards conformance — fail on one restatement: `guide.md:162` and the comment at `kit-preamble.sh:98-100` claim symlinks are followed (first finding). The preamble's header comment (lines 21-27) still matches.
4. Directory layout — pass; both files are where the guide and the other kit tests expect them.
5. Simplicity — pass; stdlib and coreutils only, no new dependency, no crypto or auth code (it detects tokens, it does not handle them).
6. Decoupling — pass; one function, three regex variables, a single caller (`kit_on_exit`).

## Not examined
- macOS bash 3.2 and BSD find, grep, head and od were not run (Linux box). Read only: `find -H`, `-print0`, `read -d ''`, `[[ =~ $var ]]` with an unquoted variable, `grep -aqE --`, `head -c`, `od -An -tx1` are all available there to my knowledge, but unverified here.
- `go test ./internal/harness -run TestKitPreambleRefusesLeftoverCredentials` was not run; the cases were replayed as bash scripts, so the Go harness (`fixtureEnv`, the `Skip` on `/tmp` inside a checkout) is read, not executed.
- The changelog entry `2026-10-05-1759-kit-credential-check-fails-closed.md` and BUG-160's record were read only for context, not reviewed.
