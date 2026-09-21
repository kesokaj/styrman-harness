---
owner: main
date: 2026-09-21
---
# Preflight — v1.5.0   (2026-09-21)

Range reviewed: v1.4.0..c7e178e

VERDICT: SHIP

Nothing blocks: no verifier FAIL, BLOCKED or IN PROGRESS; no high-severity security finding; the security report has its `## Not examined` section and names no release-diff path there; reviewer APPROVE; no REQUIREMENTS violation (the § Known limits edit was operator-approved, ADR-035); no shadowing project agent (`.claude/agents`, `.agents/agents` absent).

## Findings, worst first, with the decision taken

| # | Source | Severity / confidence | Finding | Decision |
|---|---|---|---|---|
| 1 | security | medium / high | The sandbox log (`testing/2026-09-21-verify-sandbox-pruned-context.md:13`) prescribes copying live credential files into a scratch home with no `chmod`, lifetime or deletion named; mitigations verified (gitignored, untracked, not embeddable, deleted). | fixed now, docs only: guide § 4 states the rule (chmod 600, run lifetime, delete by name; prefer the real-home checks). The dated log stays as written. |
| 2 | reviewer | low / medium | `size_test.go` sums `len(ctx) + Σ len(r.Raw)` with no per-file wrapper allowance; exact while `rules/` is empty, ~300 bytes optimistic per rule file later. | accepted, named in the release entry; follow-up: a measured per-file allowance before the first `rules/` file (the plan already defers `rules/` behind a second-file probe). Owner: main. |
| 3 | security | low / high | `--dangerously-skip-permissions` became the fallback for headless Antigravity checks with live credentials; the six `permissions.allow` entries recorded on 2026-09-16 are gone from the real settings. | guide § 5 now says the allow list is the standing way and the flag is for throwaway clones only; restoring the operator's settings is the operator's (named to them). |
| 4 | security + reviewer | low / medium | The canary word is in tracked files (README, this repo's AGENTS.md, the guide), so a run started inside a harness checkout could answer it from the project context. | fixed now, docs only: guide § 5 requires a demo clone outside any harness checkout and, on Antigravity, the store check beside the answer. |
| 5 | docs lens | medium / high | `context/AGENTS.md` DoD step 2 dropped "a saved diff file only when the tree will stand still, and the reviewer pins the blobs it read". | rebutted: plan D16 — `agents/reviewer.md` never had a blob-pinning duty, so the allowance was unenforced; the rule is stricter now (a git range or a path list). Advisory in the gate. Owner: main, no action. |
| 6 | reviewer | note | `skills/cloud-run-deploy/SKILL.md:37` pre-flight bullet is unnumbered inside a numbered list. | accepted; a one-character follow-up, not worth reopening verification for. Owner: main. |

## Rows

- Reviewer (opus): `reviews/2026-09-21-review-v1.4.0-HEAD-code.md` — APPROVE WITH NITS; canary traced to all `RenderContext` callers and both vendor targets; both new tests can fail; the pruned rules' new homes hold. Earlier reports covering parts of the range: `reviews/2026-09-21-review-worktree-code.md`, `-docs.md`, `-docs-2.md`, `-docs-3.md`.
- Security: `reviews/2026-09-21-security-v1.4.0-HEAD.md` — FINDINGS, none high; govulncheck clean; go.mod/go.sum untouched; no hand-rolled crypto; GCP guidance strictly tightened.
- Verification: `testing/2026-09-21-verify-preflight-v1.5.0.md` — PASS (gofmt, vet, race suite, size test 20,205 / 22,000, govulncheck, list, scratch install drift none, canary last, vendors identical). Cited: `testing/2026-09-21-verify-real-home-smoke-5.md`, `-verify-sandbox-pruned-context.md`, `-context-consolidation-phase-2-behavioural.md`.
- Docs lens (Step 4): `reviews/2026-09-21-review-v1.4.0-HEAD-docs.md` — one advisory finding (row 5).
- Trail freshness (Step 3): newest `changes/` entry is the Phase 4 entry of this release's work — pass; no future `date:` — pass; `changes/` and `DECISIONS.md` append-only (the one modification, 5e08629, is ADR-035 appended: 10 insertions, 0 deletions) — pass; `ROADMAP.md`/`TODO.md` touched twice in the range — report-only.

## What was not checked

The Mac and the devbox (they update to this release afterwards); the full skill matrix; guide § 3a–3c (no `cmd/sh-agent` verb changed); the audit prompt's reference limits beyond the three vendor pages fetched on 2026-09-21; the per-file rules overhead (row 2) beyond the one probe.
