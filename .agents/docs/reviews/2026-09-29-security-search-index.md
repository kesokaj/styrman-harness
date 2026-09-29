---
owner: security
date: 2026-09-29
---
# Security review — search-index   (2026-09-29, round 3)

Range reviewed: c04f2d7..5fe69af. Round 3 covers the round-2 fixes in 5fe69af (254f3f5..5fe69af), every declaration they touched, and every reference to those declarations. Round 2 covered c04f2d7..254f3f5; round 1 covered c04f2d7..565bb1e.

## Status: COMPLETE

## Verdict: FINDINGS

Round-2 finding 4 is only half fixed, and a probe reproduced its worst case on HEAD.
- A clean tracked file still takes its ID from git's index and its content from a later read.
- Probe R edited a committed file while `search` waited for the write lock, then restored the file with `git checkout --`. From then on, `search` found the edited word in the clean file and none of its committed words.
- The Round 2 outcome and 5fe69af's message both call this finding fixed.

The half that was fixed introduced a memory regression. Every modified and untracked file's bytes are now held for the whole reconcile, including those whose blob the index already holds. Probe M measured a steady-state `search` at 187 MB resident over 180 MB of untracked files; the round-1 binary took 15 MB.

Three fixes are correct, and each has a test that fails on the round-1 code:
- the global view reads only `.agents/docs`;
- `printable` covers the hint and both notes;
- the secret list ignores case.

Two new gaps:
- The new `.agents/docs` check drops a registration on any stat error, which the project's hard rule forbids (probe S).
- README and the spec now claim `.env*`, which the list does not match (probe E).

Nothing here is above low. Secrets, injection, crypto, dependencies and GCP have nothing new.

## Findings

- **[low / confidence: high]** `internal/search/reconcile.go:103` (a tracked entry takes the ID from `git ls-files -s`) and `:354` (`readIndexable` with `e.data == nil`). This is the unfixed half of round-2 finding 4, now reproduced.
  - **How it happens:** the content is read in `applyLocked`, after `diffCount` and `BEGIN IMMEDIATE`, which waits up to the 10 s busy timeout, and after every blob indexed before it. Nothing compares that content with the ID.
  - **Probe R** (HEAD binary, scratch repository):
    - `b.md` was committed as "bravo committed text".
    - A second SQLite connection held `BEGIN IMMEDIATE` on the project's `search.db` for 4 s while `search bravo` ran. At 1.5 s the probe rewrote the file as "bravo poisoned zebra".
    - After the search finished, `git checkout -- b.md` made the file clean again: its bytes on disk were the committed ones.
    - `search zebra` then printed `b.md:1  doc … bravo poisoned »zebra«`. `search committed` found 0 hits. Both reported "reconciled 0 changed".
  - **How long it lasts:** the wrong bytes stay bound to the committed blob's ID until `--rebuild`. They drive snippets, `--refs`, and the record status that `--open` prints.
  - **Other consequences:**
    - In the global index the blob row is shared, so another project with the same committed `.agents/docs` blob is served these bytes.
    - A committed `.gitattributes` with `eol`, `working-tree-encoding` or `ident` makes the working-tree bytes differ from the blob by design.
    - `configView` has the same shape (the key from one read, the content from a second), but an append-only log never returns to an old key.
  - **Untrue descriptions:** these claims hold only for modified and untracked files:
    - the Round 2 outcome ("fixed");
    - 5fe69af's message ("hashed bytes are indexed bytes");
    - `entry`'s doc comment.
  - **Remediation (one mechanism for this finding and the next; it shortens the code):**
    - Remove `entry.data` and `readIndexable`'s `data` parameter.
    - In `applyLocked`, read an unknown blob's file once. Index it only when its blob ID, computed in the ID's own format, equals `e.oid`: git's SHA-1 blob format for a 40-hex ID, SHA-256 for a 64-hex ID, and the plain SHA-256 for `configView`'s `sha256:` key.
    - On a mismatch, insert nothing; the next query retries.
    - A file under an eol, encoding or ident conversion then stays unindexed. Say so in the spec.
    - Alternative for tracked entries: read the blob by ID with `git cat-file --batch`. Only IDs go on stdin, so round 1's path-quoting problem cannot arise. It needs `GIT_NO_LAZY_FETCH=1` so that a partial clone never fetches. It also gives `git()`'s unused `stdin` parameter a caller.
  - **Test:** probe R's arrangement. Hold the lock from a second connection, edit a clean tracked file during the wait, and restore it. Assert that the restored words are found and the edited ones are not. `TestModifiedFileReadOnce` then asserts that the indexed bytes hash to the ID, not a read count.

- **[low / confidence: high]** `internal/search/reconcile.go:145` (`entry{…, data: data}`). The round-2 fix introduced this regression.
  - `gitView` now keeps the bytes of every modified and untracked file until the reconcile ends. The Lstat gate at `:124-125` caps each file at 1 MiB, but nothing caps how many files there are.
  - That includes files whose blob the index already holds, where `applyLocked` never reads the bytes. Before the fix, each file was hashed and then released.
  - **Probe M:** a scratch checkout held 200 untracked files of 900 KB each. They were binary, so none were indexed. Maximum RSS of a steady-state `search memory`:

    | Binary | Run A | Run B |
    |---|---|---|
    | 254f3f5 | 15,160 KB | 15,960 KB |
    | HEAD | 187,428 KB | 187,268 KB |

  - **Who is exposed:**
    - Memory now grows with the size of the untracked, non-ignored tree, on every `search`.
    - The context tells every agent to run `search` first, and parallel subagents multiply the cost.
    - A data or build directory that is not ignored is enough to exhaust a laptop.
    - A clone's committed files are clean and carry no bytes, so this comes from the local tree, not from a repository an attacker controls.
  - **Remediation:** the previous finding's. Once the field is gone, the hashing read is one file at a time again.

- **[low / confidence: high]** `internal/search/search.go:305-314`. The new check drops a registration on any `os.Stat` error, not only "does not exist".
  - Causes include EACCES, EIO, ELOOP, and an unmounted or network volume.
  - The drop deletes the registration's rows from the global index. It prints "no longer has .agents/docs", which is then false.
  - It breaks the project's hard rule: "A guard in front of a destructive action keeps the data when its check ERRORS".
  - **Probe S:**
    - A checkout was registered, then `chmod 000 <root>/.agents`.
    - `--global wombat`, run from another checkout, printed `note: dropped …/r4 (no longer has .agents/docs)`.
    - With the mode restored, `--global wombat` still found 0 hits. Only a new search inside the checkout brought it back.
  - The older `checkoutRoot` branch at `:303` has the same shape. Any git failure drops the registration: `safe.directory` reporting "dubious ownership", `git` missing from PATH, or a corrupt index.
  - There is no attacker path. The cost is the availability of a registry that heals itself later.
  - **Remediation:**
    - Drop only when `errors.Is(err, fs.ErrNotExist)`, or when `err == nil && !fi.IsDir()`.
    - On any other error, skip the root for this query with a note, and keep the registration.
    - Treat the git branch the same way: drop only when the root is gone or git reports it is no longer a work tree.
  - **Test:** probe S's arrangement. The checkout answers again once the mode is restored, with no search inside it. Skip the test when running as root.

- **[low / confidence: high]** `README.md:185` and `.agents/docs/specs/sh-agent/search.md:35` now say the list covers `.env*`. The list at `internal/harness/plan.go:219` matches only `.env`, `.env.*`, `*.env` and `.envrc`.
  - **Probe E** (untracked files, HEAD binary):
    - Indexed and printed: `.env-prod`, `.env_local`, `.envprod`, `.pgpass`, `id_rsa_deploy`.
    - Skipped: `.ENV.LOCAL`, `Credentials`, `deploy.PEM`, `x/.Git-Credentials`. Case folding and the two new names work.
  - `plan.go:119` makes the same claim ("The names match the context's Git rule"). That rule stops a commit for `.env*`.
  - Round 2's option (b) was half taken, and no decision is recorded. The list still lacks `.pgpass`, `id_rsa*`, `id_ed25519*` and `id_ecdsa*`, and the Round 2 outcome has no `left: operator` line for them.
  - The rest of option (a), bumping `schemaVersion`, was not done: `store.go:28` is still `"1"`. So content indexed before the fix stays inside each existing `search.db` until gc (14 days), and in free pages after. That content includes `.ENV`, `.git-credentials` and `credentials` files. It is no longer served, because every answer joins `files` and the next reconcile deletes those rows.
  - **Remediation:**
    - Replace `.env`, `.env.*` and `.envrc` with `.env*`. The list gets shorter and matches the context's rule, README and the spec. Install and adopt then refuse `.env-*` names too, which the Git rule already stops.
    - Put the remaining names and the schema bump to the operator as one-line choices. Recommended: the bump, because it is a one-character change and the only way to purge content indexed before the fix.

- **[low / confidence: high]** `internal/search/store.go:128` (the rebuilt-index note). This is round-2 finding 3's third trigger, and it has no test.
  - The code is correct: it calls `printable`.
  - `TestHintAndNotesArePrintable` covers only the hint and the dropped note.
  - **Test:** a checkout whose path carries ESC, holding a `search.db` that contains garbage bytes. Assert that the rebuilt note is printed and carries no ESC.

## Checked

- **Global view (`.agents/docs` only):** correct. `docs := ".agents/docs"` is a constant.
  - `docsRootOf` keeps its one caller (`search.go:171`, the project scope), so the fix orphaned nothing.
  - `reg.docs` is still scanned at `:292` but never read, and `projects.docs_root` is written by `register` and read nowhere on the global path. This is dead data, and none of it is exported.
  - Registrations made before the round-1 fix are now dropped on the next `--global`.
  - A committed `.agents/docs` symlink passes the stat, but `gitView` lists nothing beyond it.
- **Printable output:** correct at `search.go:722`, `:314` and `store.go:128`.
  - `printable` maps `\n`, which is below 0x20, so round 2's newline-in-path case is closed.
  - The remaining raw `Fprint` sites print numbers, the user's own `--refs` argument (`:703`, `:706`, `:713`), or state values the code sets.
  - Error text on stderr is still raw, as round 2 noted.
- **Hashed bytes, modified and untracked:** correct.
  - When `data` is set, `readIndexable` skips its own Lstat. That is safe only because `gitView`'s Lstat regular-file and size gate (`:124-125`) runs first.
  - The Lstat-then-read race is the same one `readIndexable` had before.
  - The tracked half and the memory cost are findings 1 and 2.
- **Secret list:** `strings.ToLower` runs before `filepath.Match`, and every pattern is lower case.
  - Install (`plan.go:120`) and adopt (`adopt.go:161`) inherit the change.
  - `go run ./cmd/sh-agent install --home .agents/scratch/secreview3/insthome` exited 0 with 266 files, and no file was refused.
  - A case-insensitive `git ls-files` grep for the new names finds no tracked path, so the new bare `credentials` line in `.gitignore` hides nothing tracked.
- **Descriptions:**
  - `gitView`'s comment is now true.
  - `go doc` prints `SecretShaped` and `secretShaped` separately.
  - The spec's reconcile steps (lines 43-46) are accurate.
- **Declarations re-read:** `entry` (three constructions; one reader of `data`), `readFile` (three uses and the test), `readIndexable` (one caller), `secretShaped`/`SecretShaped`, `reconcileGlobal`, `register`, `printable`, `openStore`.
  - No export was left without a caller.
  - `git()`'s `stdin` parameter still has no caller that passes one (unexported).
- **Tests must fail on the old code:** HEAD's `security_test.go` was copied into a `git archive 254f3f5` tree, with a test-only `readFile` variable so it compiled. It compiled, and all four round-2 tests failed:
  - `TestGlobalDropsCheckoutWithoutAgentsDocs`, `TestHintAndNotesArePrintable` and `TestSkipsMoreSecretNames` failed with their own messages.
  - `TestModifiedFileReadOnce` counted 0 reads, because the variable was not wired into the old code. That failure is not evidence (see Not examined).
- **Secrets:** none in 5fe69af. The fixture `https://u:caracal@example.invalid` uses a reserved TLD.
- **Injection:** unchanged. Queries use bound parameters, and the hint's `LIKE` has an `ESCAPE`.
- **Authz/authn:** none in scope.
- **Crypto:** unchanged. `crypto/sha1` (the blob ID) and `crypto/sha256` (the config key) are content keys, used as designed.
- **Dependencies:**
  - 5fe69af changes none.
  - `govulncheck -show verbose ./...` (go1.27.1): 0 vulnerabilities affect the code.
  - The one advisory in a required module is GO-2026-5024, in `golang.org/x/sys@v0.37.0` (fixed in v0.44.0). It is not called; this is D25's advisory.
  - `npm audit` does not apply.
- **GCP surface:** none.
- **Build and tests:**
  - `gofmt -l cmd internal`: clean.
  - `go vet ./...`: clean.
  - `go test -race -count=1 ./...`: all four packages ok. This is the first full run in this review series.
- **Probes:** everything is under `.agents/scratch/secreview3/`.
  - The HEAD binary was built from the clean tree. The old binary came from `git archive 254f3f5`.
  - Each repository was made with `git init`, with its exit status and top level checked.
  - Every run passed `--home .agents/scratch/secreview3/home`, so the real home was never touched.
  - Afterwards, `git status` showed only this file.

## Not examined

- **macOS:** the case-insensitive filesystem was not tested. All probes ran on Linux.
- **Race width:** probe R held the lock from a second connection to open the window. The width of that window between two real agents was not measured.
- **`TestModifiedFileReadOnce` on the old code:** its failing state with the variable wired into the round-1 code, as the Round 2 outcome describes, was not reproduced.
- **Code not read:** `wanted`, `diffCount` and `insertContent` beyond their use of the view, and modernc internals.
- **Markdown not read:** everything in 5fe69af beyond `README.md:185`, `search.md:35-46` and `:73`, and plan D34. Commit `e848a39` was not read.
- **Antigravity:** no runs.
- **Stderr error text:** not re-probed; carried from round 2.
- **Round-2 finding 1:** registration, left by the operator as D34. The code still matches that decision (registration needs `.agents/docs` at `search.go:197`), so it is not raised again.

## Round 1

Each fix has a test that failed before it (run bare on the unfixed code, 2026-09-29: all six FAIL), in `internal/search/security_test.go`.

- 1. [medium] `dsn` puts the path into a `file:` URI unescaped — fixed at `internal/search/store.go:72`: the DSN is built with `net/url` (`url.URL{Scheme: "file", Path: …, RawQuery: …}`), so `?`, `#` and `%` stay in the path; `TestPathWithURISpecialCharacters` (a checkout named `q?mark#frag%20x`: the index lands inside it, nothing appears beside it).
- 2. [medium] every search registers the checkout for good — fixed at `internal/search/search.go:194`: only a checkout with an `.agents/docs` directory registers; README, `usage` and the spec say so; `TestRegistersOnlyHarnessProjects` (a clone with a plain `docs/` is not in `--global`, a harness one is). Registrations made before the fix stay until the checkout disappears or `index --global --rebuild`.
- 3. [low] `git hash-object --stdin-paths` misreads names with a trailing CR or a leading `"` — fixed at `internal/search/reconcile.go:138,157`: modified and untracked files get a git blob ID computed from their own bytes (`blobID`), and the `hash-object` call is gone; `TestHashesOddNamesItself` (both names found with their own content, `--global` still answers).
- 4. [low] control characters reach the terminal — fixed at `internal/search/search.go:564` (`printable`), applied in `cleanSnippet`, `markTarget`, every hit label, the `--refs` record line and every `--open` line; `TestOutputStripsControlCharacters` (words, `--refs`, `--open`: no ESC, no BEL).
- 5. [low] `~/.config/sh-agent` created 0755 — fixed at `internal/search/store.go:145`: directories the index creates are 0700; `TestGlobalDirIsPrivate`.
- 6. [low] untracked secret-shaped files indexed and kept for 14 days — fixed at `internal/search/reconcile.go:142`: names matching the installer's secret list (`harness.SecretShaped`, exported at `internal/harness/plan.go:212`) never enter the view; `TestSkipsSecretShapedFiles` (`.env` and `keys/server.pem` not found, `notes.md` found).

## Round 2

Each fix below has a test in `internal/search/security_test.go` that failed on the round-1 code (run bare 2026-09-29: all four FAIL, the read-once test after a behaviour-free `readFile` seam was added so it compiled).

- 1. [medium] a crafted clone that commits `.agents/docs` still registers — left: operator — the operator chose "Keep automatic (Recommended)" on 2026-09-29 over explicit registration; the residual (a clone the operator chooses to search in joins `--global` and `/retro` evidence until `index --global --rebuild`) is recorded in plan D34.
- 2. [low] `reconcileGlobal` fell back to a plain `docs/` — fixed at `internal/search/search.go` (`reconcileGlobal`): the global search index always reads `.agents/docs` and drops a registration that no longer has it, with `note: dropped <root> (no longer has .agents/docs)`; `TestGlobalDropsCheckoutWithoutAgentsDocs`.
- 3. [low] hint and note lines printed raw paths — fixed: `printable` on the `--refs` hint (`search.go`), the dropped-registration note (`search.go`) and the rebuilt-index note (`store.go`); the registrations-dropped note prints only a number; `TestHintAndNotesArePrintable` (a directory named with ESC, a vanished root named with ESC).
- 4. [low] content indexed from a second read — fixed at `internal/search/reconcile.go`: a modified or untracked entry keeps the bytes it was hashed from, and `readIndexable` uses them; `TestModifiedFileReadOnce` (one read).
- 5. [low] the secret-name list was case-sensitive and missed `.git-credentials` and `.aws/credentials` — fixed at `internal/harness/plan.go` (`secretShaped` lower-cases the name; `.git-credentials` and `credentials` added) and `.gitignore` kept in step; `TestSkipsMoreSecretNames`.
- 6. [low] stale descriptions — fixed: the `gitView` comment, `SecretShaped`'s doc comment moved off `secretShaped`'s, README's claim narrowed to "files with secret-shaped names … a secret inside an ordinarily named file is [indexed]"; the spec says the same.

## Round 3

The review's three-round cap. The operator chose, 2026-09-29, "Revert and document (Recommended)": undo the round-2 read-once change, record the race as a known limit, fix the check-error drop and the wording, stop reviewing.

- 1. [low] a clean tracked file can be indexed with content from an edit made and undone while `search` waits for the lock — left: operator — known limit in `specs/sh-agent/search.md` § Failure modes and plan D35; `index --rebuild` clears it.
- 2. [low] round 2's `entry.data` held every modified and untracked file in memory (15 MB to 187 MB) — fixed by reverting it (`internal/search/reconcile.go`: no `data` field, no `readFile` seam, `readIndexable` reads the file as before); `TestModifiedFileReadOnce` removed with it.
- 3. [low] a registration was dropped on any stat error — fixed at `internal/search/search.go` (`dropReason`): it goes only on positive evidence (root gone, git says not a checkout or another one, `.agents/docs` gone or not a directory); any other error keeps it and prints `note: skipped <root> this time (<error>)`; `TestRegistrationKeptWhenCheckErrors` (`.agents` mode 000: skipped, then served again after the mode is restored; mutation R1 fails it).
- 4. [low] README and spec claimed `.env*` — fixed: both list `.env`, `.env.*`, `*.env`, the list's own patterns.
- 5. [low] the rebuilt-index note had no test — fixed: `TestRebuiltNoteIsPrintable` (a checkout named with ESC and a corrupt index; mutation R2, the note without `printable`, fails it).

## Outcome

- Round 1, findings 1–6: fixed (commit 254f3f5), as the Round 1 lines say.
- Round 2, finding 1 (registration of a crafted clone): left: operator — "Keep automatic (Recommended)", plan D34.
- Round 2, findings 2, 3, 5, 6: fixed (commit 5fe69af).
- Round 2, finding 4 (second read): left: operator — reverted at round 3 with the round-3 finding 1 residual; plan D35.
- Round 3, finding 1: left: operator (plan D35); findings 2–5: fixed in the commit that adds this section.
