---
owner: main
date: 2026-10-09
---
# Metrics store — Definition of Done: review fixes, BUG-164/165, verification

- Reviews (dispatched over 2c0f215..166aa39):
  - `reviewer`, code + docs in one: `reviews/2026-10-08-review-2c0f215-166aa39-code-docs.md`, approve with nits.
  - `security`, three rounds, the third on the operator's "Fix both in this feature": `reviews/2026-10-08-security-2c0f215-166aa39.md`.

  Both records end with `## Outcome`.
- Fixed, each under a test that was red first:
  - `Verb` skips git's global options before the subcommand, and recognises `<<\EOF` here-docs.
  - Antigravity tool and generation rows update on re-ingest.
  - Stale `metrics-agy-*` and `dream-metrics-*` copies are swept by `ingest` and by `dream`.
  - `export` leaves out the `cmd_hash` salt.
  - `ingest` and `export` refuse flags they do not read.
  - `--project` matches the root or anything under it (a case-sensitive prefix), and events group by git toplevel.
  - The `signedRepo` fixture no longer flakes. A stray test output file was removed by name.
- **BUG-164 (escaped since v2.0.0, fixed).** `git log` and `git status` in project roots ran repository-defined programs: a signature program and a clean filter. `internal/trail.command` is now the only way git runs in a project root. It sets no lock, no fsmonitor, no signature check and `GIT_NO_LAZY_FETCH=1`, and it skips partial clones. `status` is replaced by `ls-tree` plus `hash-object --no-filters`, and `Toplevel` must be a parent of the directory. This covers the dreamer's read too (D39).
- **BUG-165 (open).** `sh-agent search` runs an unhardened `git status`. It is outside this range and is the next task.
- Docs:
  - `deploy/local.md`: the dreamer section names the nightly ingest and `metrics export` as the backup.
  - Spec: the git hardening, partial clones, `--project`, the export salt, and flag refusals.
  - `dream.md` and README: the raw-blob comparison.
  - The plan is closed with D37–D39, and ROADMAP Now is cleared.
- Verification: `verifier` wrote `testing/2026-10-09-verify-metrics-store.md`, **PASS**: gofmt and vet clean; `go test -race -v -count=1 -timeout 30m ./...` pass 610, fail 0, skip 3 (env-gated); § 3g tier pass 274.
- Not done, for the operator:
  - The first real ingest and the real-home smoke.
  - The REQUIREMENTS R3 wording ("`show --json`" → the numbers-only summary the code sends).
  - The Claude Code `retro` matrix control run at 2c0f215.
  - The interactive map page, which has lagged the README map since 2026-09-17.
  - `govulncheck`, which is not installed on this box.
- Agents dispatched: `reviewer` (1), `security` (3 rounds), `verifier` (1).
