---
owner: security
date: 2026-09-29
---
# Security review — search-index   (2026-09-29)

Range reviewed: c04f2d7..565bb1e (the committed range only; the working tree's uncommitted edits to `go.mod`, `go.sum`, `AGENTS.md` and `DECISIONS.md` were not reviewed)

## Status: COMPLETE

## Verdict: FINDINGS

The worst problem is in `internal/search/store.go:65`. `dsn()` pastes the index's filesystem path into a `file:` URI DSN without escaping it. modernc.org/sqlite splits a DSN at its first `?` and runs every `_pragma` query value as SQL through an exec that runs every statement in the string. So a checkout whose absolute path contains `?` (or `#` or `%`, which SQLite's URI parser also interprets) opens a different database file than the one `guardProject` checked. It also turns the rest of the directory name into connection parameters. A probe confirmed the misplacement. The global search index is the second, design-level problem. Any `sh-agent search` run registers that checkout silently and for good. That includes a third-party clone the explorer agent is told to search first. From then on, the checkout's whole `docs/` tree is served to every `--global` and `--open` query, and the retro skill treats those results as cross-project evidence for harness-scope rules. The remaining findings are low. SQL and FTS5 query construction are sound: every value is bound, and every FTS5 word is quoted.

## Findings

- **[medium / confidence: high (weakness), medium (exploitability)]** `internal/search/store.go:65` (`dsn`, used by `tryOpen` for both `projectDB(root)` and `globalDB(home)`). The path is inserted raw into `file:%s?_pragma=...`.
  - modernc.org/sqlite v1.46.1 `newConn` (conn.go) splits the DSN at the FIRST `?`. `applyQueryParams` (sqlite.go) runs each `_pragma` value as `"pragma " + v`. `stmt.exec` (stmt.go:177) loops over every statement in the string. SQLite, opened with `SQLITE_OPEN_URI`, also decodes `%XX` in the path and ends the path at `#`.
  - **Probe:** a scratch checkout at `.agents/scratch/secreview/probe/q?mark` built its index at `.agents/scratch/secreview/probe/q`, a SQLite file outside the checkout. The `journal_mode(WAL)` pragma was swallowed into a bogus parameter name. `guardProject`'s symlink check and gitignore check had approved a different path.
  - **In principle:** a directory name carrying `?_pragma=<percent-encoded SQL>` reaches SQL execution on the index connection. That includes `ATTACH`, which creates a SQLite file at any path the user can write. Candidate carriers: a submodule path chosen in an attacker's `.gitmodules`, with the agent's working directory inside the submodule (`--show-toplevel` is then the submodule's own directory), or the top-level folder of an extracted archive.
  - No payload was built.
  - **Remediation:** build the DSN from `url.URL{Scheme: "file", OmitHost: true, Path: path, RawQuery: <pragmas>}`, which percent-escapes `?`, `#` and `%` in the path (stdlib `net/url`, no dependency). Alternatively, refuse a path containing any of the three with a clear error.
  - **Test:** a checkout named `q?x` must get its database at `<root>/.agents/scratch/search.db` in WAL mode, and nothing may be created next to the checkout.

- **[medium / confidence: medium]** `internal/search/search.go:194` (`register` on every `prepare`) and `:247-321` (`reconcileGlobal`); spec `search.md:38`. Every checkout that runs `search` or `index` registers itself in `~/.config/sh-agent/search.db`, silently. Nothing checks that the checkout is the operator's own or harness-initialized: `docsRootOf` falls back to a plain `docs/`, which most third-party repositories have. The range's agent and skill text makes this automatic:
  - `agents/explorer.md`: the first command is `sh-agent search`.
  - context: "Run it before Grep".

  After one session in an untrusted clone, its `docs/**` is part of every later `--global` answer, from any directory:
  - `docs/plans/*.md` with `- [ ]` boxes and `docs/bugs/BUG-NNN*.md` with `status: open` are listed by `--open --global` as open records.
  - `skills/retro/SKILL.md` step 9 uses `--global` hits as "cross-project evidence" that makes a finding harness-scope. This is a persistent, cross-project indirect prompt-injection path into retro proposals for the harness's global rules. An operator review stands in front of any applied change.
  - Global hits are labelled `filepath.Base(root)`. On a name collision the FIRST-registered checkout keeps the bare name and later ones get `~2`. After `index --global --rebuild`, whichever checkout is searched first takes a trusted project's name.
  - Each `--global` query runs `git status` and `git hash-object` in every registered checkout. For a checkout whose `.git/config` came from the attacker (an archive with `.git/`), that re-runs `core.fsmonitor` and filter commands from unrelated directories. The same precondition already gives code execution on any git command, so this adds persistence, not new access.

  **Remediation (shortest first):**
  - register only a checkout whose docs root is `.agents/docs` (one condition before `register`);
  - print the full root for global hits whose base name is not unique;
  - optionally pass `-c core.fsmonitor=false` to the git calls made in registered checkouts.

- **[low / confidence: high]** `internal/search/reconcile.go:111` and `:127`. `gitView` excludes only paths containing `\n` before feeding `git hash-object --stdin-paths`. The probe confirmed that git also strips a trailing `\r` from each line and C-unquotes a line that starts with `"`: `x.md\r` and `"x.md"` both hashed as `x.md`, and `y.md\r` with no `y.md` failed with exit 128. Consequences:
  - the OID of one file is bound to another file's content;
  - content is hidden from the index, or substituted when the two files' order in the view lets the odd name be read first. Blobs are keyed by OID alone and shared across branches in a checkout's index and across projects in the global one;
  - `hash-object` fails outright when the unquoted or stripped name does not exist. `reconcileGlobal` returns on the first failing checkout (`search.go:305-307`), so one such file name under any registered checkout's docs root makes EVERY `--global` query fail until that checkout is fixed or the registrations are rebuilt.

  **Remediation:** exclude paths containing `\r` or starting with `"` (or any control character) together with `\n`. Better, key blob rows by a hash of the bytes `readIndexable` actually read, which also closes the race between `git status`/`hash-object` and the read. Make `reconcileGlobal` skip a checkout whose view fails, with a note, rather than fail every project.

- **[low / confidence: high]** `internal/search/search.go:552` (`cleanSnippet`), `:602` (`markTarget`), `:580`, `:670`, `:806`, `:301`. Snippets, record titles, status and date values, paths and dropped-root notes come from indexed file content and file names, and they are printed raw. `strings.Fields` removes only Unicode spaces, so ESC, BEL and C1 controls (CSI U+009B) pass through. A person running `sh-agent search` in a terminal over an untrusted checkout, or over a `--global` index that includes one, gets the file's terminal control sequences: rewritten or hidden lines, title changes, OSC 8 links, and OSC 52 clipboard writes where the terminal allows them. **Remediation:** map every `unicode.IsControl` rune except the `\x01`/`\x02` markers (which are replaced first) to a visible placeholder in each printed field (stdlib `strings.Map`).

- **[low / confidence: high]** `internal/search/store.go:137` (`tryOpen`: `os.MkdirAll(filepath.Dir(path), 0o755)`). Every per-checkout search calls `register`, which opens the global database. On a home where `~/.config/sh-agent` does not exist yet, the first search therefore creates it 0755, and SQLite creates the database 0644 (observed in the probe). That undoes the 0700 property `harness.ConfigDir` enforces (security review of the 2026-09-18 audit) until the next `install` runs `TightenConfigDir`. The file holds copies of `feedback.jsonl`, `retros.jsonl`, `proposals/*.md` and every registered trail. **Remediation:** create the global database's directory 0700 (`harness.ConfigDir`, or `os.MkdirAll(dir, 0o700)` for the global path), with a test that asserts the mode on a fresh home.

- **[low / confidence: medium]** `internal/search/reconcile.go:95` and `:21` (`--untracked-files=all`; `gcAge` of 14 days). Untracked files that are not gitignored are indexed whatever their name, including secret-shaped ones (`*.pem`, `*.key`, `*.tfvars`, `*.tfstate`, credential JSON not yet ignored). A keyword search then prints their lines as snippets into the agent's context, and the agents are now told to search first. Their content stays in `.agents/scratch/search.db` for 14 days after the file is deleted, and longer in free pages, since `secure_delete` is off. In the same way, `~/.config/sh-agent/search.db` keeps an UNREDACTED copy of `feedback.jsonl` lines for 14 days after the operator scrubs one. `report` itself redacts and never sends the database, which was verified. **Remediation:** skip the secret-shaped path list the git rule already names in `gitView`.

## Checked

- **Secrets:** no secret-shaped path and no credential pattern (private-key armor, AKIA, ghp_/github_pat_, AIza, xox, sk-, key/token/password literals) in the range's added lines. `report` sends only named files (`feedback.jsonl`, `retros.jsonl`, `state.json`, `proposals/*.md`); `search.db`, `-wal` and `-shm` are only listed as skipped (`report.go` at HEAD).
- **Injection:**
  - All SQL uses bound parameters. `sc.in()` builds only `?` placeholders, table names come from constant lists, and LIKE values are escaped with `ESCAPE '\'`.
  - Each FTS5 word is quoted with `"` doubled (`ftsQuery`), so no query syntax gets through.
  - The regexes are RE2, which rules out ReDoS. Files are capped at 1 MiB, and files with a NUL byte in their first 8 KiB are skipped.
  - `git` arguments are passed as argv with no shell. `-C <root>` consumes a leading `-`, and `exec.Command` refuses a binary found through a `.` in PATH.
  - DSN injection and the hash-object line parsing are the findings above.
- **Symlinks and paths:** `.agents` and `.agents/scratch` are refused when they are symlinks. A tracked `search.db` is refused by `git check-ignore`, since tracked paths are never ignored. SQLite opens the database, WAL and SHM files with `O_NOFOLLOW` (`lib/sqlite_linux_amd64.go:24968`, `:23704`), so a tracked `-wal`/`-shm` symlink fails closed. Symlinks, submodules and paths beyond a symlink are dropped from the view. An untracked, ignored `search.db` symlink can only arrive in an archive that carries `.git/`, and that archive already grants execution through `.git/config`.
- **Authz/authn:** none in scope. The command is local and runs as the invoking user.
- **Crypto:** `crypto/sha256` is used as a content key for the `~/.config` files (a checksum, consumed as designed). Nothing is hand-rolled.
- **Dependencies:**
  - `modernc.org/sqlite` v1.46.1 and ten transitive modules are justified by ADR-045 (operator-approved; R4 replaced; the `AGENTS.md` exception is at HEAD). The ADR rejects vendoring, cgo, LanceDB, a system `sqlite3` and "no index". A stdlib-only inverted index is not among the rejected options. The ADR says v1.60.1 was measured, but `go.mod` pins v1.46.1.
  - `go.mod` marks the directly imported `modernc.org/sqlite` as `// indirect`, so `go mod tidy` was not run.
  - `go mod verify`: all modules verified.
  - `govulncheck ./...` on the HEAD tree (go1.27.1): 0 reachable. One module-level advisory is not called: GO-2026-5024, in `golang.org/x/sys@v0.37.0`, Windows-only, fixed in v0.44.0, recorded as D25.
  - `npm audit` does not apply: the range has no npm changes.
- **GCP surface:** none in the range. No IAM, service accounts or deploy flags changed.
- **Probes:** everything is under `.agents/scratch/secreview/`: a HEAD tree copy made with `git archive HEAD`, a binary built from it, and two scratch repositories created with `git init` and addressed by absolute path. The search runs used `--home .agents/scratch/secreview/home` and never touched the real home.

## Not examined

- `internal/search/search_test.go`, `global_test.go`, `recall_test.go` and `cmd/sh-agent/search_test.go`: only their `Home`/`--home` isolation was read (each uses `t.TempDir()` or passes `--home`). The assertions, and whether tests cover the findings above, were not read.
- The Markdown trail (plans, specs other than `search.md`'s registration lines, testing logs, bug records, change entries), `README.md`, and the agent and skill diffs beyond the `sh-agent search` lines quoted above.
- The global-index trust finding comes from reading the code and the agent and skill text. No agent run was made to show an injected trail steering a retro.
- Case-insensitive filesystem behaviour (macOS) of `gitView`/`readIndexable` was not tested. The probe ran on Linux only.
- The rest of the transpiled SQLite/libc code in modernc beyond the DSN parsing and `O_NOFOLLOW` sites cited above.
- Commit `e848a39` (AGENTS.md § Testing, BUG-106/107 index rows): HEAD moved to it during this review. It is outside the reviewed range and was not read.

## Round 1

Each fix has a test that failed before it (run bare on the unfixed code, 2026-09-29: all six FAIL), in `internal/search/security_test.go`.

- 1. [medium] `dsn` puts the path into a `file:` URI unescaped — fixed at `internal/search/store.go:72`: the DSN is built with `net/url` (`url.URL{Scheme: "file", Path: …, RawQuery: …}`), so `?`, `#` and `%` stay in the path; `TestPathWithURISpecialCharacters` (a checkout named `q?mark#frag%20x`: the index lands inside it, nothing appears beside it).
- 2. [medium] every search registers the checkout for good — fixed at `internal/search/search.go:194`: only a checkout with an `.agents/docs` directory registers; README, `usage` and the spec say so; `TestRegistersOnlyHarnessProjects` (a clone with a plain `docs/` is not in `--global`, a harness one is). Registrations made before the fix stay until the checkout disappears or `index --global --rebuild`.
- 3. [low] `git hash-object --stdin-paths` misreads names with a trailing CR or a leading `"` — fixed at `internal/search/reconcile.go:138,157`: modified and untracked files get a git blob ID computed from their own bytes (`blobID`), and the `hash-object` call is gone; `TestHashesOddNamesItself` (both names found with their own content, `--global` still answers).
- 4. [low] control characters reach the terminal — fixed at `internal/search/search.go:564` (`printable`), applied in `cleanSnippet`, `markTarget`, every hit label, the `--refs` record line and every `--open` line; `TestOutputStripsControlCharacters` (words, `--refs`, `--open`: no ESC, no BEL).
- 5. [low] `~/.config/sh-agent` created 0755 — fixed at `internal/search/store.go:145`: directories the index creates are 0700; `TestGlobalDirIsPrivate`.
- 6. [low] untracked secret-shaped files indexed and kept for 14 days — fixed at `internal/search/reconcile.go:142`: names matching the installer's secret list (`harness.SecretShaped`, exported at `internal/harness/plan.go:212`) never enter the view; `TestSkipsSecretShapedFiles` (`.env` and `keys/server.pem` not found, `notes.md` found).
