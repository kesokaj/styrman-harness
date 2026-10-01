---
owner: reviewer
date: 2026-09-30
---
# Review — 4764c5f-worktree-code-r2   (2026-09-30)

Range reviewed: working tree against commit 4764c5f (`git diff 4764c5f`; 16 files under cmd/sh-agent, internal/dream, internal/harness, internal/search). Round 2; round 1 is `2026-09-30-review-4764c5f-worktree-code.md`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

The four round-1 answers hold (see "Round 1 answers"). The worst remaining problem is a gap in the IN2 fix: a failed stale-file removal in the resync loop returns without saving the manifest, so the files this run already wrote stay unlisted and uninstall leaves them. Two smaller defects in the IN3 (symlink) change: Status and Diff treat a file under a linked directory as a normal file while Install keeps it, and a stranded doc-comment. Nothing is broken on the main paths. Not run by me: go test, vet, gofmt (read-only review, no re-run; the verify log records them).

## Findings
- **[medium / confidence: high]** `internal/harness/apply.go:~141` — the resync loop's `return fmt.Errorf("removing stale %s: %w", ...)` does not go through `saveOnError`. By then every planned file is already written (and on disk); the manifest on disk is still the OLD one, so new files written in this run are unlisted and `Uninstall` leaves them: the exact IN2 failure, on a path the fix did not cover (only the `WriteAtomic` error calls `saveOnError`). Fix: `return saveOnError(fmt.Errorf(...))`; add a test where a stale file's directory is unwritable (as in TestInstallMidFailureKeepsManifestComplete).
- **[low / confidence: high]** `internal/harness/apply.go:220,254` (Status, Diff) — they test only `Lstat(dest)` for a symlink, not `underLink(home, f.Rel)`. For a file under an operator's linked directory (the TestUninstallAndLinkedDirs arrangement: `~/.claude/skills` linked) Install prints `keep ... (a symlink ...)` and never touches it, while Status reads through the link and reports `differs`/`missing` (drift) and Diff prints a diff, forever: a drift the tool refuses to fix. Fix: use the same predicate (`symlink || underLink`) in Status and Diff; TestInstallDoesNotFollowLinkedDest asserts `foreign` only for direct links.
- **[low / confidence: high]** `internal/harness/apply.go:~405-425` — `underLink` was inserted between `pruneEmptyDirs`'s doc-comment and its `func`, so the "pruneEmptyDirs removes now-empty directories ..." text now heads `underLink`'s doc and `pruneEmptyDirs` has none: a description that names the wrong function. Fix: move `underLink` above that comment.
- **[low / confidence: low]** `internal/harness/apply.go:64` — when `~/.claude`, `~/.claude/skills` or `~/.local/bin` (the planned `BinPath`) is a link (common with dotfile managers), Install writes nothing, prints one keep line per file and exits 0; for `.local/bin/sh-agent` the binary is then never replaced by `update`. This is the stated intent of the security fix, and `status` shows the binary line, so it is not a defect in the code as written; flagged only because it is a behaviour change for such machines and is an operator decision (no test for the `.local/bin` case).

## Round 1 answers
1. Deploy guide `update --check` sentence: still stale at `.agents/docs/deploy/remote.md:75` ("works on a dirty checkout"); deferred to the documentation pass as stated, a docs file outside this diff, not counted here. It must be fixed in that pass (and the line-92 timer advice read against it).
2. Redaction value rule: accepted. Quoted values whole, per-key length rule, short shapes pinned as redacted and `token=abcdefg` / `secret = my passphrase` pinned as accepted under-redaction in `TestReportRedactShapes`; reading the `accept` func and the FindAll loop, no input I traced (`token: secret=abcdefgh`, `secret: x token=abcdef12`, unterminated quote) escapes a value the rule should catch.
3. Nested-first parse: rebutted acceptably. `TestAgyCorpusLeaves` runs on the corpus (env-gated measurement, openCopy, the reader's own query); the before count (971, D30) and after count (971, 0 keys) are in the verify log. Wire types of ordinary prose make a false nested parse improbable.
4. Symlink in `saveOnError`, resync loop and `Uninstall`: fixed (`apply.go:47`, `:134`, `:327`), with `TestUninstallAndLinkedDirs` and `TestResyncSkipsLink`. Gap found in Status/Diff above.

## Criteria
1. Dead code — pass: `copyDB` still used by openCopy; nothing orphaned.
2. Flow tracing — fail (minor): resync-remove error path skips the manifest save; Status/Diff disagree with Install on linked directories. Other hops traced OK: `answer()` (sole caller realRun), `removeCopies(conf, now)` callers, `find()` chunk-id dedupe, `update()` ordering (dirty check above both read-only returns, tests changed with it), openCopy callers.
3. Standards conformance — fail (minor): misplaced doc-comment; remote.md:75 restatement still stale (deferred, see above).
4. Directory layout — pass: no new source files; tests beside code.
5. Simplicity (DIY-except-security) — pass: stdlib only; redaction is not auth/crypto.
6. Decoupling — pass: `copyFile` var is a seam inside one package; no new cross-package reach-ins.

## Not examined
- Bodies of run_test.go, dream_test.go (beyond TestAgyCorpusLeaves), collect_test.go, search_test.go, main_test.go, memory_test.go.
- Did not re-run gofmt, vet, `go test -race`, govulncheck, or any behavioural/sandbox run (AGENTS.md § Testing); rely on the verify log for those.
- `trailDates` for a project that is a subdirectory of a repo or has non-ASCII paths (git quotes them; falls back to mtime, unchanged from the old behaviour); no timeout on its git calls.
- The repo's living-doc restatements of the dream render indent and fence escape in `specs/sh-agent/dream.md`.

## Outcome

Answered in the next round (see `2026-09-30-review-4764c5f-worktree-code.md`); the fixes and their tests are in the phase commit. Every finding of this round was fixed, except the deploy-guide sentence, which the documentation pass of the phase fixes, and the nested-first prose-loss concern, rebutted by the corpus measurement (971 leaves and 0 argument keys before and after, `testing/2026-09-30-verify-audit-fix-phase-1.md`).
