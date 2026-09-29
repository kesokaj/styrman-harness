---
owner: reviewer
date: 2026-09-29
---
# Review — search-index-code-r2   (2026-09-29)

Range reviewed: 565bb1e..aae16af

## Status: COMPLETE

## Verdict: APPROVE

The diff is the recorded security-round-2/round-3 fix set (D32–D35) plus its docs. Each fix is narrowly scoped, has a targeted test in `internal/search/security_test.go`, and the accompanying `README.md`/spec/`AGENTS.md`/`DECISIONS.md` text was updated in the same commits. `go build ./...` and `go test ./internal/search/...` pass. Nothing is broken; the two points below are low-severity and don't block.

## Findings

- **[low / confidence: medium]** `internal/search/search.go:322` (`dropReason`) — the "no longer a git checkout" branch is detected by `strings.Contains(err.Error(), "not a git repository")`, i.e. by matching git's English stderr text. `git` localizes this message under `gettext` catalogs on a non-C locale, so on such a machine a genuinely-gone checkout would never match, fall into the generic `err != nil` branch, and the registration would be kept forever with a recurring "note: skipped … this time" rather than ever being dropped. This doesn't leak anything (the safe/keep direction is the one taken on ambiguity, matching D35's intent), but it silently defeats the intended cleanup path on such machines. Low impact given this is a single-operator Linux/macOS CLI tool that in practice runs `git` under `LANG=C`/no gettext catalogs installed, so I'm not blocking on it — worth a one-line note if it's ever seen in practice.

- **[low / confidence: low]** `internal/search/search.go:602-603` — the `printable` doc comment says "Tabs and newlines are folded to spaces by the callers first," but most call sites (`sc.label`/`rawLabel`, the `--refs` note/hint lines, `dropReason`'s note lines, the rebuilt-DB note) pass raw, unfolded text straight into `printable` — only `cleanSnippet` and `markTarget` actually fold whitespace first. `printable` still neutralizes a raw `\n`/`\t` (both are `< 0x20`, mapped to `?`), so there's no functional/security gap — the comment just overstates what upstream folding happens before most calls. Not a correctness issue, so borderline as a "finding" at all; noting it since the project's own rule treats a comment that no longer matches the code it documents as reportable.

## Criteria

1. **Dead code** — pass. No orphaned functions; `docsRootOf` still has a live caller (`prepare`'s local-scope path), the old `hash-object`-based loop was fully replaced (not left dangling), and the removed `strings.ContainsAny(p, "\n")` filter was only needed for the newline-joined `git hash-object --stdin-paths` protocol that this diff deletes — its removal is correct, not a regression (paths are now read individually via `os.ReadFile`/`filepath.Join`, immune to embedded `\n`, and any raw `\n` that reaches an output line is still neutralized by `printable`).
2. **Flow tracing** — pass. `blobID`'s SHA-1 git-blob hash is computed directly from `os.ReadFile` bytes and threaded through `entry.oid` into `store.apply`/`diffCount` exactly as the old `hash-object` output was. The new `dropReason`/`register` gate on `.agents/docs` is applied consistently at both registration (`prepare`, gated on `os.Stat(.agents/docs)`) and reconciliation (`reconcileGlobal` → `dropReason`), and `docsRootOf` vs. the hardcoded `".agents/docs"` in `reconcileGlobal` don't diverge because the gate already guarantees `.agents/docs` exists whenever `docsRootOf` is called on a registered root. `printable` is applied at every output site that carries indexed/attacker-reachable content (label, snippet, refs rows, hints, drop/skip notes, rebuilt note); the `TestOutputStripsControlCharacters`/`TestHintAndNotesArePrintable`/`TestRebuiltNoteIsPrintable` tests exercise this concretely, each with a real failing condition (verified `go test ./internal/search/...` passes). `dsn()` correctly percent-encodes `?`/`#`/`%` in the path via `net/url`, verified against the driver's DSN-splitting behavior. Every new test can fail (each plants an attacker-shaped fixture — a `\x1b` byte, a secret-shaped name, a missing `.agents/docs` — and asserts the specific defect is absent), consistent with the review-standards bar.
3. **Standards conformance** — pass. `secretShaped`'s case-insensitive match and the `.git-credentials`/`credentials` additions are mirrored in `.gitignore` and in `README.md`'s `search` row, so the "every restatement of a changed rule" grep is satisfied for this vocabulary. `AGENTS.md`'s `govulncheck` line and `DECISIONS.md` ADR-045a were updated together in-range, consistent with the "describing text moves with behaviour" rule.
4. **Directory layout** — pass. `SecretShaped` lives in `internal/harness` (the package that already owns the credential-shape list for the installer) and is consumed by `internal/search` through an exported function, not a reach-in; tests stay beside their packages.
5. **Simplicity (DIY-except-security)** — pass. `crypto/sha1` is used only to reproduce git's own blob-ID scheme for content addressing (not a security primitive), which is stdlib-as-designed, not hand-rolled crypto. No new dependency was added for this fix set; `go.mod`/`go.sum` churn is `modernc.org/sqlite` moving from indirect to direct plus its own transitive build deps (`go mod tidy`), not a new functional dependency.
6. **Decoupling** — pass. `internal/search` depends on `internal/harness` only through the one exported `SecretShaped` function, sized to what it needs; no new cross-package reach-ins were introduced.

## Not examined

- The three review-report markdown files added in this range (`reviews/2026-09-29-review-search-index-code.md`, `reviews/2026-09-29-security-search-index.md`, `reviews/2026-09-29-review-8ee5d70-HEAD-docs.md`) — out of the dispatch's named scope (Go source + README/AGENTS/.gitignore/spec), and they are other reviewers' own artifacts, not source to critique.
- `internal/search/parse.go`, `store.go`'s `ensureSchema`/`applyLocked`, and the rest of `search.go` outside the diffed hunks — unchanged in this range, not re-reviewed.
- Behavioural verification beyond `go build ./...` and `go test ./internal/search/...` (both green); no manual CLI smoke test was run against a real `~/.config/sh-agent`, per the dispatch's code-only scope.
- Locale-dependent git stderr text (finding 1 above) was reasoned about, not empirically reproduced under a non-C locale.

## Outcome

The reviewer stopped at its turn limit after writing this report (Status COMPLETE).

- 1. `dropReason` matched git's English stderr — fixed at `internal/search/reconcile.go` (`git`): every git call runs with `LC_ALL=C` and an empty `LANGUAGE`; `TestDropReasonIgnoresLocale` (a fake `git` that answers in German unless `LC_ALL=C`: failed before the fix, passes after).
- 2. `printable`'s doc comment said callers fold whitespace — fixed at `internal/search/search.go` (`printable`): the comment says tabs and newlines become `?` except in snippets.
