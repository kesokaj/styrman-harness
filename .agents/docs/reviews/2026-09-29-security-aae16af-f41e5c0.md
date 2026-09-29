---
owner: security
date: 2026-09-29
---
# Security review — aae16af-f41e5c0   (2026-09-29)

Range reviewed: aae16af..f41e5c0

## Status: COMPLETE

## Verdict: FINDINGS

The code changes in this range add no weakness of their own. The `LC_ALL=C` git environment, the `parseRefs` prefilter, the prepared inserts, the dropped `refs(oid)` index and the `printable` comment are all sound. The findings come from the range's instruction change. `context/AGENTS.md:9` (§ Resume, then checkpoint) and `skills/handoff/SKILL.md:154` now make every session in an existing project run `sh-agent search --open` before it does anything else. That automatically runs the project-index path (`prepare` → `guardProject` → `openStore`), which has two gaps that were already there before this range.

The worse gap is that `guardProject` refuses a symlinked `.agents` or `.agents/scratch` directory but never checks `.agents/scratch/search.db` itself. SQLite resolves a symlinked main-database path and writes through it. In a probe I confirmed that an untracked, gitignored `search.db` symlink made `sh-agent search --open` write an 86,016-byte SQLite database, holding the checkout's text, into the file the link named. That is the exact write-through the guard's own comment says it refuses.

The second gap is that the project index is created mode 0644. It copies the full text of untracked files regardless of their own permission bits.

## Findings

- **[medium / confidence: high]** `internal/search/search.go:128-145` (`guardProject`), reached from `prepare` at `search.go:158-167` and `store.go:144-147`. **Pre-existing; not introduced by the range. The range's new resume step now fires it at the start of every session.**
  - **The gap:** `guardProject` checks `.agents` and `.agents/scratch` with `os.Lstat` but never checks the database file. SQLite 3.51.2 (modernc v1.46.1) resolves symlinks in the main-database path (`_appendOnePathElement`, `SQLITE_MAX_SYMLINK`), so a symlinked `search.db` is followed and opened read-write with create.
  - **Probe** (scratch repository `.agents/scratch/sec-aae16af/probe`):
    - The setup was `.gitignore` = `.agents/scratch/*` plus an untracked `.agents/scratch/search.db` symlink pointing at an empty scratch file. `git check-ignore` passed.
    - `sh-agent search --open` then exited 0, and the target went from 0 bytes to 86,016 bytes starting `SQLite format 3`.
    - A git-delivered (tracked) symlink is refused: `check-ignore` exits 1 and the target stayed at 0 bytes. So the planted link has to arrive outside git, for example in an archive of a working tree that includes `.git` and the scratch directory, or through a patch applied with `git apply`.
  - **Effects:**
    - Verified: an empty target file is filled with a SQLite database.
    - Expected from the code but not probed: a missing target path is created. A foreign SQLite database gets `PRAGMA journal_mode(WAL)` (from the DSN) applied before `ensureSchema` rejects it as stale. A symlink to `~/.config/sh-agent/search.db` writes the whole checkout into the global index.
    - Expected to be safe: an existing non-SQLite file should hit NOTADB, after which `removeDB` unlinks only the link.
  - **Remediation (stdlib only):** in `guardProject`, `os.Lstat` `.agents/scratch/search.db` and refuse unless it is absent or a regular file. Optionally check `-wal`/`-shm` too; SQLite already opens those with `O_NOFOLLOW` (vendored `lib/sqlite_linux_amd64.go:23704`, `:24968`).
  - **Test:** plant an untracked symlinked `search.db` in an ignored scratch directory, pointing at a scratch file. Assert the command refuses and the target's size and bytes are unchanged. The current code fails this test.
- **[low / confidence: high]** `internal/search/store.go:142-147` (`tryOpen`). **Pre-existing; now reached at every session start.**
  - **The gap:** the project index is created with SQLite's default mode 0644. I observed 0644 under umask 0002, so the mode comes from SQLite, not from the umask. `MkdirAll(..., 0o700)` only sets the mode of directories it creates. A scaffolded `.agents/scratch/` that already exists (0775 in the probe) is left as it is.
  - **Why it matters:** the index holds the full text of every indexed file, including untracked files that are not secret-shaped by name. In the probe, a 0600 `local-settings.txt` canary came back from `sh-agent search canary` out of a 0644 `search.db`. On a machine where other local users can traverse the checkout, the index widens a file's owner-only permission to world-readable.
  - **Remediation (stdlib only):** create the database file 0600 before `sql.Open`, for example `os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600)` then close it, or `os.Chmod` it after creation. SQLite gives its `-wal`/`-shm` files the main file's mode, so they follow. Test: build an index in a fresh scratch checkout and assert `search.db` has mode 0600.

## Checked

- **Secrets:** grepped every added line of `git diff aae16af..f41e5c0` for token and key shapes (AWS, GitHub, Slack, Anthropic/OpenAI `sk-`, Google `AIza`, PEM headers, password/secret/api_key/token assignments, e-mail, session IDs). The only hits were the `idToken`/`pathToken` regex declarations. No secret appears in logs or error text. The new `git()` error text keeps git's stderr, and every note that prints it goes through `printable`.
- **Injection:**
  - The new prepared inserts bind every value as a parameter. `gc`'s table names are constants.
  - `parseRefs`' prefilter is a pure superset guard. It checks for `BUG-`/`ADR-` as literal substrings, which `idToken` requires. It checks for `/` or a `.<ext>` from the same `refExtensions` slice the regex is built from. Go initialises `refExtensions` before `pathToken` by dependency order. RE2 is linear, so there is no ReDoS. `hasRefExtension` is at most 15 linear scans of a line, and files are capped at 1 MiB. `TestRefsPrefilterIsEquivalent` passes.
  - `LC_ALL=C`/`LANGUAGE=` is appended after `os.Environ()`, and exec keeps the last duplicate. Every other variable is inherited exactly as before, when Env was nil. With `-z` output, path quoting is unaffected.
- **Git config exposure:** `git status` in a checkout still runs that checkout's own git config, for example `core.fsmonitor`. The resume step adds no exposure here, because sessions already run `git status`, and a git-delivered `.git/config` cannot be set by a remote.
- **Authz/authn:**
  - The range contains no auth.
  - `dropReason` now matches git's C-locale text under any locale. That can only unregister a cache entry, and the entry comes back on the next local search. A failed check still keeps the registration (`TestDropReasonIgnoresLocale` is a real must-fail).
  - `--open` in the resume step is project-scoped (`sc.projects = []proj{p}` unless `--global`), so no other project's trail enters a session. Record titles pass through `printable`. Its lines come from working-tree files the session already trusts, untracked ones included.
- **Crypto:** none added. The pre-existing `sha1` (git blob naming) and `sha256` (content keys for config files) in `reconcile.go` are content identifiers, i.e. `crypto/*` used as designed. Nothing is hand-rolled.
- **Dependencies:** `go.mod`/`go.sum` are unchanged in the range. `govulncheck ./...` (go1.27.1) found 0 vulnerabilities affecting the code. At module level it lists GO-2026-5024 (`golang.org/x/sys@v0.37.0`, windows `NewNTUnicodeString`), which the code does not call; it predates the range. `npm audit` does not apply: the range touches no JS files.
- **GCP surface:** nothing in the range touches IAM, service accounts, deploy flags or secrets.
- **Tests:** `go vet ./internal/search/` is clean. `go test -race -count=1 ./internal/search/` passes (6.4 s), run on the working tree at f41e5c0.
- **Dropped `refs(oid)` index** (not a security issue, for the code reviewer): `gc()`'s `DELETE FROM refs WHERE oid IN (...)` still filters on `refs.oid`, so each GC now scans the whole `refs` table. That affects performance only. `schemaVersion` stays "1", so older databases keep the index, which does no harm.

## Not examined

- `cmd/sh-agent` main: whether errors that carry git's stderr (the `gitView` failures in `prepare`) are printed to the terminal without passing through `printable`. This is pre-existing, and I did not read it.
- `search.go` outside lines 112-365, 596-625 and 840-877 (the `words`/`refs` query paths). `parse.go` outside lines 170-260 (`parseRecords`, `chunkText`, `splitLines`). `global_test.go`, `recall_test.go`.
- The trail-only files in the range (handoff, plans, cold reads, testing logs, the change entries): these were secret-grepped only, not read.
- Not probed:
  - a symlink pointing at a missing path (creation);
  - a foreign SQLite target (WAL conversion);
  - a symlink pointing at the global index;
  - macOS behaviour (the probe ran on Linux only).
- How Antigravity behaves when running the new resume step.
- The probe and its build are left in `.agents/scratch/sec-aae16af/` (gitignored). This includes a staged symlink in the probe repository's own index; the real repository was not touched.

## Outcome

Each fix has a test in `internal/search/security_test.go` that failed before it (run bare on f41e5c0: both FAIL) and a mutation that compiles and fails it (S1–S3).

- 1. [medium] `search.db` itself was not checked for a symlink — fixed at `internal/search/store.go` (`prepareFiles`, called by `tryOpen` for both indexes): the database and its `-wal` and `-shm` companions are refused unless absent or regular files, and a new database is created with `O_EXCL`; `TestRefuseSymlinkedDatabase` (each of the three names as a symlink, and the global search index: refused, the link's target unchanged). Spec § Failure modes says so.
- 2. [low] the index was created 0644 — fixed in the same function: created 0600, and an existing file with group or other bits is tightened on open; `TestDatabaseIsPrivate` (both indexes and their companions 0600; a 0644 index tightened). Spec § Where the indexes live says so.
