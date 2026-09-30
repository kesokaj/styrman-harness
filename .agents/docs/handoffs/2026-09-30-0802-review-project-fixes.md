---
owner: main
date: 2026-09-30
---
# Handoff — 2026-09-30 — project review fixed; v1.13.1 is the operator's call

Session ID: 014yCPpqSoS7aMhqehoUno6S        (same machine: `claude --resume 014yCPpqSoS7aMhqehoUno6S`)
Project: /home/kesokaj/workspace/styrman-harness
Branch: main

## Context

The operator ran `/review --project --all` at 3478a1c (v1.13.0 plus one docs commit), then said "fix all issues in your recommended way" and asked whether a new version is needed. The fixes are committed and pushed as e59638c; nothing is queued. Two findings are recorded `left: operator` with a recommended answer each, and the release is a PATCH the operator cuts.

## Completed

- The review run (commit 12ef66b): seven artifacts and `.agents/docs/reviews/2026-09-30-review-project-merged.md`. No FAIL; one medium (simplify), five low (security).
- The fixes (commit e59638c), three review rounds, all in the merged record's Round 1–3 and Outcome sections:
  - `internal/harness/sweep.go` `SweepStaleDirs` replaces the glob sweep in `release.go` and `sweepStaleZipDirs` in `cmd/sh-agent/main.go` (BUG-109).
  - `internal/harness/report.go`: sixteen redaction needles added, no leading `\b`; `pushError` scrubs URL credentials (BUG-110).
  - `skills/VENDORED.md`: `gke-basics` (google/skills 6f0b877) and `terraform` (antonbabenko/terraform-skill v1.6.0, 2271bc4) recorded as forks with their local changes; `AGENTS.md` vendoring rule names them; `TestVendoredRowsCarryAPin` (BUG-111).
  - `WriteAtomic` exported and shared by the settings writer; `abandon` names a temp file it could not remove; close-before-rename restored after round 2 swapped it, pinned by `TestWriteAtomicClosesBeforeRename` (BUG-112).
  - ADR-046: the checkout place's unsigned `main` pull is an accepted risk.
  - Docs: `.agents/docs/deploy/local.md:66` and `.agents/docs/specs/sh-agent/delivery.md:56` describe the redaction kinds.
- Checks: `gofmt`, `go vet`, `go test -race ./...` (371 pass, 0 fail, 2 skip), scratch install with no drift, `govulncheck` clean; verifier logs for the project and for the fixes (three rounds, mutation protocol on every fix).

## Open / In Progress

- Nothing in progress. No plan is active; `ROADMAP.md` and `TODO.md` are empty.
- Requirements edits awaiting approval: none.

## Key Decisions

- The two unpinned skills were NOT re-vendored verbatim: both had been rewritten locally, and a verbatim copy would change what agents read (an Antigravity run per changed skill before any release). Recorded as pinned forks instead; the rename question is the operator's (below).
- Finding 3 (checkout trusts `main`'s tip) accepted as designed (R9), ADR-046, rather than adding a signature check, which would be a path change under AGENTS.md rule 3.
- The redaction deny-list was extended three times in three rounds; the security agent and the main agent agree a fourth batch is the wrong shape, so the remainder goes to the operator as one choice.
- Round 3 was the cap. A round-1 fix caused the round-2 defect (WriteAtomic's ordering); its purpose was not in doubt, so the original order was restored and BUG-112 records the two attempts.

## Corrections

- None this session (the feedback log's newest entry is 2026-09-29T21:02:42+02:00, from the previous session).

## Agent notes (unverified)

- `commands/` and `rules/` hold only `.gitkeep` since fc89509 while the project AGENTS.md layout table describes them as holding deployed prompts and rules (explorer). Both are wired into `LoadDocs`; whether that is intended scaffolding is a one-line check for the operator.
- Two read-only globs built from the home path remain (`internal/harness/state.go` `LastRetro`, `internal/search/reconcile.go` `configView`); the security agent confirmed neither can do more than read.
- The main-agent probe of the push-error shape used `/tmp/x` and `/tmp/y` twice before switching to scratch; no file of the project's was involved and nothing was committed from there.

## Environment State

- Branch: main, clean, pushed (HEAD e59638c).
- Running processes: none. Every agent dispatched this session has returned; the round-3 verifier was stopped after its log was on disk.
- Scratch: `.agents/scratch/vendor/` holds the two upstream clones used for the provenance diffs (gitignored, safe to delete).

## Artifacts to Read

- `.agents/docs/reviews/2026-09-30-review-project-merged.md` — the review, the three rounds, the Outcome with the two operator choices.
- `.agents/docs/changes/2026-09-30-0758-review-project-fixes.md` — the change entry, agents dispatched, PATCH candidate.
- `.agents/docs/reviews/2026-09-30-security-worktree.md` — round-3 security check: the remaining redaction shapes, with the per-kind rules it recommends.
- `.agents/docs/testing/2026-09-30-verify-review-fixes.md` — verifier log, rounds 1–3.
- `.agents/docs/bugs/INDEX.md` — BUG-109 to BUG-112.
- `.agents/docs/architecture/DECISIONS.md` § ADR-046.

## Next Steps

1. Operator: cut v1.13.1 with `/preflight v1.13.1` from your own turn. The change alters released `update` and `report` behaviour by fixing defects only, so PATCH. No skill or agent definition changed, so the gate's Antigravity-run rule does not add runs; the real-home smoke on both CLIs still applies.
2. Operator choice, redaction remainder (merged record § Outcome, finding 4): (a) one rule per kind — `(?i)hooks\.slack\.com/`, `://[^/\s:@]{16,}@`, an `sk-None-` prefix — recommended, a Task with a failing test each; or (b) accept and leave the record as is.
3. Operator choice, fork names (finding 5): keep `gke-basics` and `terraform` under their upstream names (recommended; the manifest now says a re-fetch never overwrites them) or rename them to satisfy the "differently named skill" rule, which touches the context's Installed list, three skill descriptions and the test matrices.
4. Operator: say whether empty `commands/` and `rules/` are intended; if not, the AGENTS.md layout line is a trivial edit.
