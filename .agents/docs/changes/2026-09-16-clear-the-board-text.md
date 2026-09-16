---
owner: changelog
date: 2026-09-16
---
# Instructions and rules: clear the board Phase 2 and 3 (2026-09-16)

## Summary
The instruction-text half of clearing every remaining open item from `.agents/docs/plans/2026-09-16-clear-the-board.md` Phase 2 and 3. Five bug records closed in text (BUG-019, BUG-020, BUG-021, BUG-022, BUG-031) and eight simulation follow-ups applied (G5, G10, G12, G13, G20; G1, G2, G7 verified already present). BUG-019: skills stopped instructing against the global hard rules (REGION required, registry derived, image tag the short SHA, amd64 ship target; daemon-start permission condition and reprobation; four dependencies with stdlib alternatives named). BUG-020: `cloud-run-deploy` now says ship/deploy/release in the user's turn means `/preflight`, lists itself as mechanics `/preflight` calls, refuses to ship on its own. BUG-021: `distill` defers to `feedback` as log owner, example carries project field, reason stated for cross-project filter exclusion. BUG-022: third-party registration scoped to public libraries, forbidden for private/internal, never agent-initiated. BUG-031: headless run ends only after every dispatched agent returns and its artifact is on disk, names any that has not. G5–G20: init-project says what filled sections mean and what empty becomes; review re-dispatch asks agent to re-read declarations a fix touched and references; runtime data comes from config with default outside repo or is gitignored; plain `go run` never dirties; documentation agent Setup section bound to contributor commands, tiers to testing guide; verifier `<app>` defined in resolution order and forbidden literal; both verifier and documentation restate it. REQUIREMENTS.md R4 amended on operator approval (ADR-027): "`update` refuses uncommitted changes on every path, because it deploys the working tree not the commit it fetched; `install` stays the deliberate way to deploy a dirty tree." (BUG-046, ADR-027).

## Type
docs/rules

## Files changed
- `skills/bash/references/makefile-patterns.md` — REGION required, registry derived, image tag the git short SHA, amd64 build target named as what ships to Google (BUG-019)
- `skills/gcloud-setup/SKILL.md` — region and zone no longer hardcoded; service-account key files reframed as last resort, not CI path (BUG-019)
- `skills/docker/SKILL.md` — daemon-start row states standing-permission condition, otherwise hands command back; smoke test captures container id and probes that container (BUG-019)
- `skills/bash/references/testing-linting.md`, `skills/frontend/references/accessibility.md` — harness runs no CI; local form named first (BUG-019)
- `skills/golang/references/golang-io-web.md` — Cobra+Viper noted with stdlib alternative (BUG-019)
- `skills/golang/references/golang-architecture.md` — `samber/lo` noted with stdlib alternative (BUG-019)
- `skills/golang/references/golang-core.md` — `samber/oops` and `goleak` noted with stdlib alternatives (BUG-019)
- `skills/golang/references/golang-testing.md` — `goleak` noted with stdlib alternative (BUG-019)
- `skills/cloud-run-deploy/SKILL.md` — description now says ship/deploy/release in user's turn means `/preflight`; skill is mechanics `/preflight` calls; Do NOT line forbids answering ship request on own (BUG-020)
- `skills/distill/SKILL.md` — defers to `skills/feedback/SKILL.md` as log owner; example carries project field; reason spelled: `/retro`'s cross-project filter excludes entries without project (BUG-021)
- `skills/golang/references/golang-ecosystem.md` — third-party registration (Context7, DeepWiki, OpenDeep, zRead) scoped to public library, forbidden for private/internal code, never agent-initiated (BUG-022)
- `context/AGENTS.md` — headless or print-mode run ends only after every dispatched agent returns and its artifact is on disk; names any that has not (BUG-031)
- `skills/init-project/SKILL.md` — says what filling a section means and what empty section becomes (G5)
- `skills/review/SKILL.md` — re-dispatch brief asks agent to re-read each declaration a fix touched and every reference to it (G10)
- `context/AGENTS.md` — runtime data path comes from configuration with default outside repo or is gitignored; plain `go run` never dirties the tree (G12)
- `agents/documentation.md`, project `AGENTS.md` template in init-project — Setup section bound to commands a contributor runs, tiers to testing guide (G13)
- `agents/verifier.md` — defines `<app>` in resolution order, forbids leaving it literal (G20)
- `agents/documentation.md` — restates `<app>` definition (G20)
- `.agents/docs/deploy/remote.md` — § The resync loop states the refusal, why it exists, the timer unit it makes safe, timer follows `main` rather than signed tag (BUG-046, ADR-027)
- `README.md` — command table and loop line agree with remote.md
- `REQUIREMENTS.md` R4 — amended on operator approval: "`update` refuses a checkout with uncommitted changes on every path, because it deploys the working tree rather than the commit it fetched; `install` stays the deliberate way to deploy a dirty tree." (ADR-027)
- `skills/preflight/SKILL.md` Step 4 — unchanged from retro 2 (verified present)

## Agents dispatched
verifier ×1, reviewer ×2 (code and docs lenses), security ×1, bugs ×1, changelog ×1

## Verification
`.agents/docs/testing/2026-09-16-verify-clear-the-board.md` — PASS. Every finding from `.agents/docs/reviews/2026-09-16-review-clear-the-board-code.md` (APPROVE WITH NITS, three findings), `2026-09-16-review-clear-the-board-docs.md` (REQUEST CHANGES, four findings), and `2026-09-16-security-clear-the-board.md` (two medium, one low) fixed, rebutted or left with an owner in each report's `## Outcome`. No pending findings block the entry.

## Docs
`.agents/docs/architecture/DECISIONS.md` (ADR-027, no changes to prior ADRs marked Superseded beyond prior entries), `.agents/docs/deploy/remote.md` (§ The resync loop, § Downstream mirrors unchanged), `AGENTS.md` (context unchanged except G12 runtime data rule), `README.md` (loop line), `REQUIREMENTS.md` (R4 amended on operator approval), nine skill references and two agent definitions (listed above). All new or amended text is committed with this entry.

## Requirements
REQUIREMENTS.md R4 amended on operator approval (ADR-027) — "`update` refuses a checkout with uncommitted changes on every path, because it deploys the working tree rather than the commit it fetched; `install` stays the deliberate way to deploy a dirty tree." No other requirements changed; all work satisfies existing REQUIREMENTS.md.
