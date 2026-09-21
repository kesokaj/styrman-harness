---
owner: reviewer
date: 2026-09-21
---
# Review — worktree-docs-3   (2026-09-21)

Range reviewed: working tree — changed files: context/AGENTS.md, context/dev-setup.md, skills/cloud-run-deploy/SKILL.md, AGENTS.md; read via `git diff HEAD -- context/AGENTS.md context/dev-setup.md skills/cloud-run-deploy/SKILL.md AGENTS.md`.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Three of the four named checks are clean: the `bugs` re-dispatch clause in `context/AGENTS.md` matches `agents/bugs.md`'s "one allowed edit" rule verbatim in substance; the new `--service-account` row and pre-flight bullet in `skills/cloud-run-deploy/SKILL.md` agree with `skills/gcloud-setup/SKILL.md` and `skills/gke-basics/SKILL.md` on Secret Manager/service-account practice (Cloud Run attached identity vs. GKE Workload Identity — different mechanisms, no contradiction); and the project `AGENTS.md` § Release parenthetical names real step numbers and titles from `skills/preflight/SKILL.md`. The one real defect is in `skills/cloud-run-deploy/SKILL.md` itself: the new "the service account it runs as" default added to `context/dev-setup.md` is not folded into this skill's own Pre-flight step 1, which still asks for and counts "all five" defaults, name-checking only the original five and never the service account — so an agent following this skill literally would not ask the operator to record the service account the deploy needs, even though later steps (line 37's bullet, and the `--service-account` table row) assume it is already recorded.

## Findings
- **[medium / confidence: high]** `skills/cloud-run-deploy/SKILL.md:24-28` — Pre-flight step 1 enumerates "project ID, region, Artifact Registry repo, service name, visibility" and says "ask the user for all five and write them there" when `docs/deploy/gcp.md` is missing or empty. This is now inconsistent with (a) `context/dev-setup.md`'s updated GCP-defaults list, which adds "the service account it runs as" as a default to ask for up front, and (b) the template `skills/init-project/assets/agents-tree/docs/deploy/gcp.md`, whose "Recorded defaults" table already has a "Service account" row. The new bullet at line 37 ("the operator names it with the other defaults") and the `--service-account` table row at line 84 both assume the service account was already gathered in step 1, but step 1's own text still says "five" and does not list it. Fix: add the service account to the step-1 list and change "all five" to "all six" (or equivalent), so an agent that only reads step 1 still asks for it.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — checked: the diff's rule/statement changes were grepped for other restatements in the living docs (`README.md`, `AGENTS.md`, `context/`, `agents/`, non-vendored `skills/`, `.agents/docs/specs`, `.agents/docs/testing/guide.md`, `.agents/docs/deploy/`). The bugs-status clause and the preflight-steps parenthetical have no stale counterpart; the GCP-defaults list has the one stale counterpart in Finding 1.
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- The three goldfish reports (`2026-09-21-goldfish-context-agents.md`, `-dev-setup.md`, `-project-agents.md`) were opened only to understand each edit's intent, per the dispatch; they were not themselves reviewed as deliverables (they are dated reviews/logs, not living docs, and the dispatch's grep scope excludes `reviews/`).
- `.agents/docs/specs/` and `.agents/docs/testing/guide.md` were grepped for the specific wording under test (service-account defaults, bugs-status rule, GCP defaults list, preflight steps) and returned no hits; they were not read end to end for unrelated staleness outside this diff's scope.
- `README.md` § How it works and the rest of its content were not read in full — only grepped for the specific terms this diff introduces or changes (Artifact Registry repository/service name list, service-account/secretAccessor wording, bugs-status wording).
- No check was made of `.agents/docs/deploy/local.md` or `.agents/docs/deploy/remote.md` (out of scope for this diff's four named checks; the diff touches no local/remote deploy claims).

## Outcome

1. fixed: `skills/cloud-run-deploy/SKILL.md` Pre-flight step 1 now lists the service account and says "all six" (same commit as the Phase 3 checkpoint).
