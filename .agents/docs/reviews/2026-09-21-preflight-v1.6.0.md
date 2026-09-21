---
owner: main
date: 2026-09-21
---
# Preflight — v1.6.0   (2026-09-21)

Range reviewed: v1.5.0..2facb44

VERDICT: SHIP

Nothing blocks: verifier PASS at f603a5a and again at 2facb44 (the post-lens commit holds one skill sentence and the ROADMAP line, read by the main agent); no high-severity security finding, the security report has its `## Not examined` section and names no release-diff path there; the source diff (7 lines in `context/AGENTS.md`, `skills/init-project/SKILL.md`, `REQUIREMENTS.md`, two deletions) is under the reviewer threshold and was read by the main agent; the REQUIREMENTS R2 edit was operator-approved as exact text (plan D4, ADR-036); no shadowing project agent.

## Findings, worst first, with the decision taken

| # | Source | Severity / confidence | Finding | Decision |
|---|---|---|---|---|
| 1 | security | medium / medium | Native `AGENTS.md` reading is wider than the stub: ancestor and nested `AGENTS.md` files (a vendored subtree's, a parent directory's) become live instructions; in this repo `context/AGENTS.md` and `skills/init-project/assets/AGENTS.md` are candidates when files beside them are read. | fixed now in 2facb44 (its own commit, so the release commit stays records-only), one sentence in `skills/init-project/SKILL.md` (say in the summary that such files will load; never keep a template `AGENTS.md` beside source); install, status and the size test re-run bare after it. A matching sentence in REQUIREMENTS R2 needs the operator's approval of exact text — handed over in the handoff. Owner: operator for R2. |
| 2 | security | low / high | Project instructions are no longer unconditional (first session after an upgrade, pre-2.1.277, Bedrock, Vertex, Foundry, telemetry off load none and say nothing); the global rules still load through `~/.claude/CLAUDE.md`. | accepted, already named in R2, the skill and ADR-036; R3 bars a hook. |
| 3 | docs lens (gate) | low / high | `ROADMAP.md` "Now" said "cold read pending" after the cold read had passed. | fixed now: the plan closes in this commit and "Now" is cleared. |

## Rows

- Reviewer: below threshold; the main agent read `git diff v1.5.0..HEAD` over the three source files and the two deletions. Docs lenses: `reviews/2026-09-21-review-v1.5.0-HEAD-docs.md` (v1.5.0..b013b52, one punctuation nit fixed) and `reviews/2026-09-21-review-v1.5.0-HEAD-docs-gate.md` (to HEAD, row 3).
- Security: `reviews/2026-09-21-security-v1.5.0-HEAD.md` — FINDINGS, none high; probe hygiene clean; govulncheck clean; no dependency change.
- Verification: `testing/2026-09-21-verify-preflight-v1.6.0.md` — PASS (gofmt, vet, race suite, size test 20,329 / 22,000, govulncheck, list, scratch install drift none, asset gone, global stub kept, root `CLAUDE.md` absent). Also `testing/2026-09-21-project-agents-md-native.md` (Phase 1) and `testing/2026-09-21-verify-agents-md-native.md` (the five probes).
- Trail freshness: newest `changes/` entry is this work's — pass; no future dates — pass; append-only (b013b52 appended ADR-036, 10/0) — pass; `ROADMAP.md` touched once in range — report-only.

## What was not checked

Nested and ancestor `AGENTS.md` loading was not measured (row 1 rests on the vendor docs); the devbox's Claude Code version; the Mac; the skill matrix beyond the `init-deploy` row's changed check; a first-session-after-upgrade run on a place other than this box.
