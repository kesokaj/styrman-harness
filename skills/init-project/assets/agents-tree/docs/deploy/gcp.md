---
owner: documentation
date: YYYY-MM-DD
target: gcp
---
# Deploy — GCP

Cloud Run first (`cloud-run-deploy` skill); GKE or Compute Engine only when Cloud Run does not fit. GCP is the only cloud. Nothing here runs unless the user says to deploy.

**Not applicable yet:** this project has no GCP target. Delete this line and fill the defaults when the user asks for one.

## Recorded defaults

Asked once, per the global context, never invented or copied from another project. `AGENTS.md` points here; the next session reads this table instead of asking again.

| Setting | Value |
|---|---|
| Project ID | |
| Region | |
| Artifact Registry repository | |
| Service name | |
| Public (unauthenticated) | yes / no |
| Service account | |

## Image

`${REGION}-docker.pkg.dev/${PROJECT}/${REPO}/${SERVICE}:${GIT_SHA}` — built with `docker buildx build --platform linux/amd64 --load`, tagged with the short SHA and nothing else: the repository enforces immutable tags (`--immutable-tags` / `immutable_tags = true`), so a moving `latest` alias cannot exist.

## Rollout and network

Decided once, each with its reason; the `cloud-run-deploy` skill reads the first five rows (its § What `docs/deploy/gcp.md` records) and the `terraform` skill the last; a row not yet decided stays empty, never guessed.

| Row | Value |
|---|---|
| Rollout | canary tag (`canary`), traffic slice (`10`), watch window, promote step (`--to-latest`) |
| Ingress / IAP / VPC | `--ingress` value; `--iap` yes/no and who holds `roles/iap.httpsResourceAccessor`; `--network`/`--subnet`/`--vpc-egress` or none |
| Cloud SQL | the `PROJECT:REGION:INSTANCE` for `--add-cloudsql-instances`, or none |
| Probes and scale | startup probe path; `--min-instances`/`--max-instances`/`--concurrency` when set |
| Secrets | `--set-secrets` names and pinned versions (names, never values) |
| Terraform state | the GCS bucket (`PROJECT-tfstate` or its suffixed name) and prefix, or none when nothing is in HCL (the `terraform` skill) |

## Procedure

Follow the `cloud-run-deploy` skill with the defaults and rows above: local test build and smoke run, amd64 build, push, `gcloud run deploy`, canary when the Rollout row says so, verify. The reason for every project-specific flag goes in the table above, not here.

## Secrets and configuration

<!-- Every secret via Secret Manager and `--set-secrets`, never `--set-env-vars`. List the names, not the values. -->

## Verify

<!-- `gcloud run services describe`, the health URL, the log query that shows the new revision serving. -->

## Rollback

```bash
gcloud run services update-traffic ${SERVICE} --region ${REGION} --to-revisions <previous>=100
```

## Infrastructure

<!-- When it is more than one service: the Terraform module path and how to plan/apply (`terraform` skill). -->
