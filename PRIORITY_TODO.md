# PRIORITY TODO — audit fix (operator, 2026-09-30)

This file goes ahead of everything else in this project (`AGENTS.md` § Rules for working here). It lists the phases of plan `.agents/docs/plans/2026-09-30-audit-fix.md`; the plan holds the tasks, the evidence (`.agents/docs/plans/2026-09-30-recon-audit-fix.md`) holds the findings by ID. A phase is ticked here only when its Checkpoint box in the plan is ticked. When every phase is ticked and the re-audit is published, this file is emptied to its title and the rule that names it is removed in the same commit.

Operator's words (2026-09-30): "I want everything fixed, improved, optimized, add new features from your audit report. Do your research first before execute. Split this up in phases and tasks and everything needs to be verified. When everything is done, re-run the audit."

## Phases

- [ ] Phase 0 — Research, plan, cold-reader check (research briefs under `plans/2026-09-30-research-audit-fix-*.md`)
- [ ] Phase 1 — Code fixes under the bug skill: reader leaf order (SD1), duplicate hits (SD2), proposal forgery (SD3), over-redaction (IN1), install orphans (IN2), dry-run guard (IN4), the small verb defects (IN11–14), dream lows (SD4–SD11)
- [ ] Phase 2 — Trail and spec truth: delivery claim, spec drifts, handoff correction, stopped-review markers, BUG-122/123 graduation with mechanisms, stale plan closed, project § Testing
- [ ] Phase 3 — Global context under budget and consistent; templates; REQUIREMENTS edits (operator approval) with ADR
- [ ] Phase 4 — Agents: lanes, tools, derive, description test, behavioural runs
- [ ] Phase 5 — Workflow skills: every WS finding; preflight and review fixes; `/review --resume` lens; size tiers
- [ ] Phase 6 — Technical skills: golang, frontend, typescript, terraform-gcp, gke, cloud-run-deploy, docker, gcloud-setup, json-api-style, bash, mermaid; strikes; VENDORED rows; link-check test
- [ ] Phase 7 — New skills and templates: requirements, spec, operate-gcp, data-gcp, api auth, testing strategy, explain-project; init-project templates
- [ ] Phase 8 — Installer simplifications (S1, S2, S5–S8, IN5–IN10)
- [ ] Phase 9 — Real-home install, sandbox runs, re-audit, handoff

## Operator decisions

Answered 2026-09-30 ("1 yes 2 yes 3 yes"): python and anti-ui-slop struck (D28); the four REQUIREMENTS edits approved, applied in Phase 3 with ADR-048 (D29); every verb stays (D30). No Open rows remain in the plan.
