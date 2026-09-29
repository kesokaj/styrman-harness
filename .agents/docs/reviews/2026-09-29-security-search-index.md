---
owner: security
date: 2026-09-29
---
# Security review — search-index   (2026-09-29, round 2)

Range reviewed: c04f2d7..254f3f5. Round 2 covers the six round-1 fixes in 254f3f5, every declaration they touched, and every reference to those declarations. Round 1 covered c04f2d7..565bb1e.

## Status: COMPLETE

## Verdict: FINDINGS

Three fixes are correct and complete: fix 1 (the DSN), fix 3 for the odd file names, and fix 5 (0700). Probes confirmed each.

Fix 2 closes only accidental registration. Its gate is an `.agents/docs` directory, which any repository can commit. So a crafted clone still reaches every `--global` answer and `/retro`'s cross-project evidence, exactly as round 1 described. The gate is also checked only at registration. `reconcileGlobal` re-derives each registered root's docs root, and that falls back to a plain `docs/`. Probe: a checkout that had `.agents/docs` when first searched served its committed `docs/evil.md` as `trail` once `.agents/docs` was removed. Registrations made before the fix keep answering the same way.

Fixes 4 and 6 have gaps. The `--refs` hint line prints indexed paths raw. Probe: a committed directory name carrying ESC reached the output unfiltered. The secret list is case-sensitive and misses common credential files. Probe: `.ENV`, `deploy.PEM`, `.git-credentials` and `.aws/credentials` were all indexed and printed.

The fixes also left three descriptions untrue: a stale `gitView` comment, a doc comment merged into the wrong function in `plan.go`, and claims in the spec and README.

## Findings

- **[medium / confidence: medium]** `internal/search/search.go:194` (the registration gate). This is round 1's finding 2, narrowed but not closed for a targeted repository.
  - The gate is `os.Stat(<root>/.agents/docs)`. That directory layout is public and trivially committed.
  - A clone that ships `.agents/docs/` registers on the first `sh-agent search` inside it. That search is explorer's first command, and the context tells agents to search before Grep.
  - From then on the clone answers every `--global` and `--open --global` query from any directory, and feeds `/retro` step 9's cross-project evidence for harness-scope rules.
  - `README.md:185` calls these checkouts "harness-scaffolded", but nothing checks that.
  - **Remediation (a design question for the operator, since registration is plan and spec behaviour):**
    - (a) Register only on an explicit `sh-agent index --global` run inside the checkout, never on `search`. Recommended: the trust decision becomes the operator's, the implicit call leaves `prepare`, and `/init-project` can run it once.
    - (b) Keep implicit registration and record the residual in `search.md`, with the operator's review of retro proposals as the control.

- **[low / confidence: high]** `internal/search/search.go:309-312` (`reconcileGlobal`: `docs := docsRootOf(r.root)`). The gate is enforced only at registration.
  - Every `--global` query re-derives each registered root's docs root. `docsRootOf` falls back to a plain `docs/`, and `.agents/docs` is never re-checked.
  - **Probe B:** a checkout registered with `.agents/docs/plans/p.md`, then `git rm -r .agents/docs` plus a committed `docs/evil.md`. `--global pelican`, run from another checkout, printed `pb:docs/evil.md:1  trail`.
  - The same code keeps every pre-fix registration of a plain-`docs/` clone answering. A branch switch or pull that drops `.agents/docs` turns an ordinary `docs/` into trail.
  - So `search.md:38`'s claim does not hold for the running code: "a third-party clone with a plain `docs/` never registers, so it cannot answer `--global` queries or feed `/retro`".
  - **Remediation (it makes the code shorter):** in `reconcileGlobal`, index `.agents/docs/` only. Drop a registered root that lacks that directory exactly like a vanished one, with the existing note. This also retires pre-fix registrations without a rebuild.
  - **Test:** register a checkout, remove its `.agents/docs`, add a committed `docs/x.md`, and assert that `--global` does not find it and that the note is printed.

- **[low / confidence: high]** `internal/search/search.go:714` (`hint: did you mean --refs %s`, with paths from `sameName` at `:785`), `:306` (`note: dropped %s`) and `internal/search/store.go:128` (`note: index at %s … rebuilt`). Fix 4's `printable` is not applied to these three print sites.
  - The hint prints indexed paths, so a committed directory name is enough to reach it.
  - **Probe A:** a checkout with `sub<ESC>[7mREV<ESC>[0m/guide.md` committed. `search --refs guide.md` printed the raw ESC sequence in its hint line. The word search over the same file printed `sub?[7mREV?[0m`, so the check can fire, and it fired only on the hint.
  - **Probe E:** a registered checkout whose directory name carries ESC was moved away. The next `--global` query printed the raw sequence in the dropped note.
  - The same fix removed the `\n` exclusion from `gitView`. A path can now carry a newline, and the hint would print it as a separate forged line into the agent's context. This follows from the same raw `Fprintf`; the newline case was not probed.
  - So `search.md:23`'s claim does not hold while the hint is raw: "a file cannot drive the terminal".
  - **Remediation:** wrap the three `Fprintf` calls in `printable`. Extend `TestOutputStripsControlCharacters` with a `--refs <bare name>` case whose directory name carries ESC and a newline.
  - Error text on stderr (`cmd/sh-agent/main.go:76`, and paths inside `fmt.Errorf`) is the last raw channel. It is lower priority.

- **[low / confidence: medium]** `internal/search/reconcile.go:126` versus `:342`. Fix 3 is correct for the odd names (probe and test). The race named in round 1's remediation is still open.
  - The blob ID comes from `gitView`'s read, or from git's index for a clean tracked file. The content stored under that ID comes from a second read, in `applyLocked` → `readIndexable`, after `BEGIN IMMEDIATE`, which can wait out the busy timeout.
  - A write between the two reads stores the new bytes under the old ID. Two cases make this likely: an agent editing while another searches, and the first full index of a checkout, which takes seconds.
  - `blobID` equals git's blob ID. For a clean tracked file, the poisoned ID is therefore the committed blob's own. Once the file returns to that content (checkout, stash, branch switch), the index serves the other bytes for it for as long as the blob stays referenced. Only `--rebuild` cures it.
  - In the global index, blob rows are shared across projects under that ID. One checkout's working tree therefore also decides another project's identical blob. Example: a committed `.gitattributes` with `working-tree-encoding=UTF-16` makes a copied file read as binary, which marks the shared blob skipped.
  - Not reproduced: there was no hook to hold the window open.
  - **Remediation:** index the bytes that were hashed.
    - For modified and untracked entries, carry the bytes from `gitView` into `applyLocked` instead of re-reading.
    - For a tracked entry with a SHA-1-shaped ID, compare `blobID(bytes)` with the ID. On a mismatch, insert nothing, so the next query retries.
    - Files under an eol or encoding conversion then stay unindexed; say so in the spec.

- **[low / confidence: high]** `internal/harness/plan.go:214-224`, reached from `internal/search/reconcile.go:142`. Fix 6 reuses the installer's list, which matches case-sensitively and by base name.
  - **Probe D** (untracked files in a scratch checkout): `.ENV`, `deploy.PEM`, `.git-credentials`, `.aws/credentials` and `id_rsa_deploy` were all indexed, and their lines were printed as snippets. `.env` was skipped.
  - On macOS's default case-insensitive volume, `.ENV` is the file a dotenv loader opens as `.env`.
  - `README.md:185` says secret-shaped files "are never indexed".
  - `schemaVersion` is still `"1"`. Secret-shaped content indexed before the fix therefore stays in each existing `.agents/scratch/search.db` until `gc` (14 days), and in free pages after.
  - **Remediation (operator's one-line choice, because extending the list also extends install's and adopt's refusals):**
    - (a) Recommended, the simpler option: case-fold the match (`strings.ToLower`), reword README to "names on the installer's secret list", and bump `schemaVersion` so existing indexes rebuild once.
    - (b) Also add `.git-credentials`, a bare `credentials`, `.pgpass`, `id_rsa*`, `id_ed25519*` and `id_ecdsa*` to the shared list.

- **[low / confidence: high]** Descriptions the fixes left untrue.
  - `internal/search/reconcile.go:72-75`: `gitView`'s doc comment still says modified and untracked blob IDs come "from one `git hash-object --stdin-paths`". `git()`'s `stdin` parameter (`:34-38`) no longer has a caller that passes one.
  - `internal/harness/plan.go:206-212`: the new `SecretShaped` comment was inserted directly under `secretShaped`'s comment, with no blank line between them. `go doc ./internal/harness SecretShaped` prints both paragraphs, starting "secretShaped reports whether a filename…", and `secretShaped` is left undocumented. Fix: move that paragraph back above `func secretShaped`.
  - `search.md:23`, `search.md:38` and `README.md:185` make the claims the probes above disproved. They change together with those fixes.

## Checked

- **Fix 1 (DSN):** correct and complete.
  - `url.URL` escapes `?`, `#` and `%` in the path. It leaves `&`, `=` and `;` literal, and SQLite gives those no meaning in the path. The driver's first `?` is now the real query separator.
  - **Probe C:** a checkout named `q?a=b&c;d#e%3Ff` got its index inside itself, in WAL mode (header bytes 18–19 = `2 2`, which the unit test does not assert). Nothing appeared beside it.
  - `dsn` has one caller. `search_test.go:78` still opens `"file:"+projectDB(root)` raw, but only under `t.TempDir()`.
  - `remote.md` names no Windows release target. The `file://C:/…` form `url.URL` would produce there therefore does not arise.
- **Fix 2:** the gate works at registration as coded (`TestRegistersOnlyHarnessProjects` passes). `os.Stat` follows a committed `.agents/docs` symlink, but `gitView` lists nothing beyond a symlink, so no outside content enters. The residuals are the first two findings.
- **Fix 3:**
  - `hash-object` and the `\n` exclusion are gone.
  - `blobID` is git's blob format over stdlib `crypto/sha1`: a content key used as designed, not hand-rolled. It is plain SHA-1 where git uses collision-detecting SHA-1DC. That matters only for a chosen-prefix pair whose two halves the attacker authors. Informational.
  - Round 1's "skip a failing checkout" was not done. `reconcileGlobal` still fails every `--global` query when one registered checkout's `git ls-files` or `git status` fails, for example on a corrupt `.git/index`. A clone's content can no longer trigger it.
- **Fix 4:**
  - `printable` covers C0, DEL and C1.
  - Invalid UTF-8 (a raw 0x9b CSI byte) comes out of `strings.Map` as U+FFFD.
  - The markers are replaced before `printable` runs, which is the right order.
  - Bidi format characters (U+202E and similar) are not control characters. They pass through and are not a finding.
  - The gaps are listed in the findings above.
- **Fix 5:** on a fresh probe home, `~/.config` and `~/.config/sh-agent` were created 0700, and `search.db` is 0644 inside them. A 0755 directory left by the pre-fix binary stays 0755 until the next `install` runs `TightenConfigDir`, as README says.
- **Fix 6:**
  - Tracked and untracked names are skipped alike.
  - `--global` trail views are covered too, because they also go through `gitView`.
  - A skipped modified or untracked file is still read once, in memory, to hash it.
- **Declarations re-read:** `dsn`, `register` and the gate, `gitView`, `blobID`, `git`, `label`/`rawLabel`, `printable`, `tryOpen`, `SecretShaped`/`secretShaped`.
  - `SecretShaped` has one caller, in search.
  - `secretShaped` keeps its three: `plan.go:120`, `adopt.go:161` and `zip_test.go`.
  - `rawLabel` has one caller, `label`.
  - No export was left without a caller.
- **Secrets:** none in 254f3f5's added lines.
- **Injection:**
  - SQL and FTS5 are unchanged since round 1: bound parameters, quoted FTS5 words.
  - The DSN is fixed (above).
  - No path reaches git as an argument or on stdin any more.
- **Authz/authn:** none in scope.
- **Crypto:** `crypto/sha1` (the blob ID) and `crypto/sha256` (the config-file key) are content keys. Nothing is hand-rolled.
- **Dependencies:**
  - `go.mod` now requires `modernc.org/sqlite v1.46.1` directly, which fixes round 1's tidy nit. That version bundles SQLite 3.51.2 (2026-01-09).
  - `go mod verify`: all modules verified.
  - `govulncheck ./...` (go1.27.1): 0 vulnerabilities in the code, 0 in imported packages, and 1 in a required module that is not called. It is taken to be round 1's Windows-only `golang.org/x/sys` advisory (D25); `-show verbose` was not run.
  - `npm audit` does not apply.
- **GCP surface:** none.
- **Tests:** the six round-1 tests pass (`go test -count=1 -v -run '<the six>' ./internal/search/`).
- **Probes:** everything is under `.agents/scratch/secreview2/`.
  - A binary was built from the clean HEAD tree.
  - The scratch repositories were made with `git init`, with its exit status checked, and addressed with `git -C`.
  - Every run passed `--home .agents/scratch/secreview2/home`, so the real home was never touched.

## Not examined

- `go test -race ./...` was not run in full. Only the six security tests ran.
- Three results were reasoned from the code, not reproduced: the read-versus-read race, the `working-tree-encoding` cross-project skip, and the newline-in-path hint line. The ESC hint and the dropped note were probed.
- macOS (a case-insensitive filesystem) was not tested. All probes ran on Linux.
- `govulncheck -show verbose` was not run, so the one module-level advisory was not named directly.
- The Markdown in 254f3f5 beyond the `search.md` and README lines cited: ADR-045a, the code and docs reviews, and the plan. Commit `e848a39` (AGENTS.md § Testing, bug-index rows) was not read.
- Items left from round 1: the assertions of every test other than the six, and modernc internals beyond DSN parsing and `O_NOFOLLOW`.

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
