---
owner: main
date: 2026-10-01
---
# Verification — Phase 6 skill-matrix rows on Claude Code (2026-10-01)

Range reviewed: working tree on 8be6c06 (Phase 6 edits uncommitted at run time, before the review fix round)

## Status: PASS

Procedure: `.agents/docs/testing/skill-matrix.md` § Procedure, run by `.agents/scratch/audit-fix/kit-phase6.sh <row>` (not committed, ADR-010), which sources `testing/kit-preamble.sh`, clones `~/sandbox/audit-fix/proj2` into `~/sandbox/audit-fix/runs/p6-<row>` outside every checkout, runs `claude -p` from the sandbox home `~/sandbox/audit-fix/sbNew` with the working tree installed (`sh-agent install --home …`, 27 skills, drift none), and runs the row's check cell in the clone with `OUT=<stdout file>`. The allowlist carries no `gcloud`, `docker`, `terraform` or `npx`: a cloud command is denied, which is the mock. Credentials copied in for the run and deleted by name afterwards. `xargs -P 3`. Skills loaded are read from the transcript's `Skill` tool uses; cloud commands from its `Bash` commands.

| row | pass | skills loaded | cloud commands attempted | secs |
|---|---|---|---|---|
| `typescript` | PASS | typescript | 0 | 12 |
| `bash` | PASS | bash | 0 | 27 |
| `frontend` | PASS | frontend | 0 | 35 |
| `golang` | PASS | golang, review-standards | 0 | 53 |
| `terraform` | PASS | terraform | 0 | 28 |
| `cloud-run` | PASS | preflight (a deploy request routes to the release gate, as the context says) | 2 (`gcloud config …`, denied: the mock) | 63 |
| `gcloud-auth` | PASS | gcloud-setup | 0 | 11 |
| `json-api` | PASS | json-api-style | 0 | 19 |
| `mermaid` | PASS | mermaid-syntax | 0 | 26 |
| `gke` | PASS | gke-basics | 0 | 98 |
| `docker` | PASS | docker, handoff | 0 | 151 |

11 of 11 rows PASS; every row loaded the skill under test by name.

## What was not exercised

The rows ran against the forks' tree, before the review fix round changed `docker` (three-stage Dockerfile, `**/` ignore list, loopback ports), `golang/references/service.md` (ErrorInfo, no BaseContext, logging rule), `terraform` (pinned secret versions, own VPC, migrate SA, `trivy --exit-code 1`), `cloud-run-deploy` (first-deploy sentence, IAP cell), `frontend` (default theme state, form error shape, two rules), `typescript` (Playwright `reuseExistingServer: false`) and `bash/SKILL.md` (routing lines): those edits are text the same rows load, and they are re-run in Phase 9 together with the Antigravity runs. The struck skills (`python`, `anti-ui-slop`) have no rows any more. Tool-backed checks (terraform validate, tflint, gcloud `--help`, docker build, npm build, tsc, shellcheck, Go sample compilation) are in the fork reports under `.agents/scratch/audit-fix/phase6-*.md` and re-run in part by the verifier logs `testing/2026-10-01-verify-audit-fix-phase-6.md` and `-r2.md`.
