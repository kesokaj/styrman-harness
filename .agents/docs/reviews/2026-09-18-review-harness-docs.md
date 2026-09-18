---
owner: reviewer
date: 2026-09-18
---
# Review — harness-docs (2026-09-18)

Range reviewed: `internal/harness` and `cmd/sh-agent`, read whole as they stand in the working tree (docs lens — not a diff), compared against `README.md`, `.agents/docs/specs/sh-agent/delivery.md`, `.agents/docs/testing/guide.md`, `.agents/docs/deploy/{local,remote}.md`, and the newest `.agents/docs/changes/` entry (`2026-09-18-install-names-every-change.md`).

## Verdict: APPROVE WITH NITS

The docs are unusually well kept — almost every command, flag, and refusal message I checked in `cmd/sh-agent/main.go` and `internal/harness/*.go` is quoted or paraphrased accurately in README § Daily commands, the deploy guides, and the delivery spec, including recent additions (the per-file `add`/`update`/`mode`/`remove` lines, the BUG-078/079 shadow and stale-binary notes, the release verification chain). Two gaps are real, both narrow: `install`'s `keep`/`kept` output for a foreign `~/.claude/CLAUDE.md` is undocumented, and README's "so `sh-agent update` works from anywhere" overstates what the checkout-discovery fallback actually does — it can refuse.

## Findings

- **[medium / confidence: high]** `internal/harness/apply.go:35` and `:143` (the `keep    ~/%s (not written by sh-agent; …)` line and the `%d kept` count in `Install`'s summary line) — no living doc names either. README.md:162 (install row of § Daily commands) documents the per-run output as `add` / `update` / `mode` / `remove` plus "counts the unchanged," and `.agents/docs/specs/sh-agent/delivery.md:43` repeats the same four-category description ("added, updated …, corrected the mode of, or removed, with the unchanged as a count"). Neither mentions the fifth category: when an unmanaged `~/.claude/CLAUDE.md` exists and differs from the harness stub, every install run prints a `keep` line for it and folds it into a separate `kept` count in the summary (`%d written, %d unchanged, %d removed, %d kept — %d files managed`). A reader following the docs and then running `sh-agent install` sees an unexplained line and count. Fix: add `keep` to the README/delivery.md description, or point at the existing "existing ~/.claude/CLAUDE.md … is left alone" paragraph (README.md:196) and note that install now also prints and counts it per run.

- **[medium / confidence: medium]** `cmd/sh-agent/main.go:546` (`knownCheckout`) — install and update refuse to build from a checkout that was only found by walking up from the working directory (or from the manifest's recorded path) when a *different* checkout is already the one on record; the two refusal messages ("this home was installed from the release binary …, not from a checkout …" and "checkout … is not the installed one (…); pass --repo … to build and install from it") appear nowhere in a living doc. README.md:171 and `.agents/docs/deploy/remote.md:77` both describe the discovery order ("finds the checkout from `--repo`, then `$SH_AGENT_REPO`, then the working directory or any of its parents, then the path recorded at the last install") and README specifically concludes "so `sh-agent update` works from anywhere" — which is true only when the discovered checkout matches the installed manifest's recorded repo, or `--repo`/`$SH_AGENT_REPO` is set explicitly. On a second-or-later install/update run from an unregistered checkout-shaped directory (an unpacked copy under review, a second clone) the command refuses instead of "working," and no doc says so. Fix: add one line to README § Daily commands (or `remote.md` § First install) naming the refusal and its two messages.

- **[low / confidence: low, unverified]** `cmd/sh-agent/main.go:301-303` — `report` on a release binary prints `report: --repo is not used here; this release binary reports the harness built into it` when `--repo` is passed; not documented in README's `report` row or `.agents/docs/deploy/local.md` § Report back. Narrow edge case (passing `--repo` to `report` specifically on a release-binary install); I did not check whether this is deliberately left out as too minor to document elsewhere.

## Criteria

Docs lens: criteria 1-6 are the language/code-quality lens and do not apply to this dispatch.

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- `internal/harness/adopt.go`, `frontmatter.go`, `derive.go`, `state.go`, `source.go`, `pathguard.go` (beyond `NoteRunningElsewhere`/`PathShadows`), `report.go`'s full redaction-shape list, and `release.go`'s edge-case messages (`"this binary reports version %q, which is not a release; leaving it alone"`, the stale `release-*` cleanup note) were not individually diffed against every doc line — spot-checked only via grep for their exact wording; I did not find contradictions but did not exhaustively re-read each doc paragraph against each message.
- `README.md` sections beyond § Daily commands, § Install, § Install without cloning, § The loop, and § What gets installed where were not read in full (e.g. § Adding things by hand's full table, anything past line ~220).
- `.agents/docs/deploy/remote.md` beyond line ~120 (the rest of § Releases, any sections after it) was not read.
- `.agents/docs/testing/guide.md` was read in full and cross-checked against `apply.go`'s dry-run/rebuild messages; no discrepancy found there.
- I did not run any command (per hard limits: no network, no real-`$HOME` install, no `sh-agent report`), so all comparisons are static text-vs-text, not observed output.
