---
name: cloud-run-deploy
description: |
  End-to-end release workflow for Google Cloud Run: build container images locally with
  Colima/Docker (ARM for local testing, AMD64 for deploy), push to GCP Artifact Registry,
  deploy to Cloud Run, verify, and roll back. This is the HOW of a Cloud Run deploy, not
  the decision to do one: "ship", "deploy", "release" or "cut vX.Y.Z" in the user's own turn
  means `/preflight`, which runs its gate and then calls this skill's steps at its deploy row
  — so enter through `/preflight` and never run these steps in its place. Use directly for the
  mechanics: an image build or Artifact Registry push, a rollback, a `gcloud run` flag, a
  failed revision, or a deploy `/preflight` has already gated.
  Do NOT use to answer a ship/deploy/release request on its own — that is `/preflight`.
  Do NOT use for Kubernetes/GKE workloads (use gke-basics), local-only container work
  (use docker), gcloud auth problems (use gcloud-setup), or Terraform-managed
  infrastructure (use terraform).
---

Release pipeline: **local build (Colima) → Artifact Registry → Cloud Run**. Build
locally first, always. Artifact Registry is the primary registry; GitHub Container
Registry only when explicitly requested.

## Pre-flight

1. **Defaults:** read the project's `docs/deploy/gcp.md` (`.agents/docs/` or `docs/`):
   project ID, region, Artifact Registry repo, service name, the service account it
   runs as, visibility. If it is missing or empty, ask the user for all six and write
   them there before anything else. Never guess a project; `gcloud config get project` and `.env` are
   cross-checks, not sources.
2. **Auth:** if any gcloud/docker-push command fails with a credential error, stop and
   apply the `gcloud-setup` skill.
3. **Colima (macOS only):** `colima status || colima start`; on Linux the native daemon is
   already running.
4. **Region:** confirm once and reuse everywhere (`REGION`, e.g. `europe-north1`);
   Artifact Registry repo and Cloud Run service should live in the same region.
5. **Registry auth (once per machine/region):**
   `gcloud auth configure-docker ${REGION}-docker.pkg.dev`
6. **Service account:** the service account the service will run as exists (`gcloud iam service-accounts describe NAME@PROJECT.iam.gserviceaccount.com`) and holds `roles/secretmanager.secretAccessor` on every secret the deploy names — the operator names it with the other defaults (`docs/deploy/gcp.md`); never fall back to the default compute account.

## Build — dual architecture

Cloud Run executes `linux/amd64`. The ARM build exists for fast local testing on
Apple Silicon; the AMD64 build is what ships. Verify both, deploy the AMD64 one.

```bash
IMAGE="${REGION}-docker.pkg.dev/${PROJECT}/${REPO}/${SERVICE}"
TAG="$(git rev-parse --short HEAD)"          # immutable tag; add :latest only as a convenience alias

# Local test build (native ARM on Apple Silicon — fast)
docker build -t "${SERVICE}:local" .
docker run --rm -p 8080:8080 "${SERVICE}:local"   # smoke-test before anything ships

# Deploy build (AMD64 under emulation on Colima)
docker buildx build --platform linux/amd64 -t "${IMAGE}:${TAG}" --load .
```

Notes:

- `buildx --platform linux/amd64` on Apple Silicon runs under QEMU emulation —
  slower but correct. If buildx complains about the builder, create one once:
  `docker buildx create --use`.
- Tag with the git SHA, not just `latest` — rollbacks and audits need immutable tags.
- If the Artifact Registry repo doesn't exist yet:
  `gcloud artifacts repositories create ${REPO} --repository-format=docker --location=${REGION}`

## Push & Deploy

```bash
docker push "${IMAGE}:${TAG}"

gcloud run deploy "${SERVICE}" \
  --image "${IMAGE}:${TAG}" \
  --region "${REGION}" \
  --platform managed \
  --port 8080
```

Common flags — apply deliberately, not by default:

| Flag | Use |
|---|---|
| `--allow-unauthenticated` | Public services only. Omit for internal services (IAM-gated is the safe default). |
| `--set-env-vars KEY=VALUE` | Plain config. Never put secrets here. |
| `--set-secrets KEY=SECRET_NAME:latest` | Secrets via Secret Manager — the only acceptable way to pass credentials. |
| `--service-account NAME@PROJECT.iam.gserviceaccount.com` | A dedicated service account per service, never the project's default compute account; grant it `roles/secretmanager.secretAccessor` on each secret `--set-secrets` names and nothing it does not use. |
| `--min-instances` / `--max-instances` | Cold-start vs cost trade-off; confirm with the user before setting min > 0 (it bills continuously). |
| `--memory` / `--cpu` | Only when the defaults (512Mi / 1 CPU) are proven insufficient. |

The service must listen on `$PORT` (Cloud Run injects it, default 8080) and bind
`0.0.0.0`, not localhost.

## Verify

```bash
URL=$(gcloud run services describe "${SERVICE}" --region "${REGION}" --format 'value(status.url)')
curl -fsS "${URL}/healthz" || curl -fsSI "${URL}"
gcloud run services logs read "${SERVICE}" --region "${REGION}" --limit 50
```

A deploy is not done until the new revision serves traffic and the logs are clean.
Report the actual revision name and URL — never claim success without the curl.

## Rollback

```bash
gcloud run revisions list --service "${SERVICE}" --region "${REGION}"
gcloud run services update-traffic "${SERVICE}" --region "${REGION}" \
  --to-revisions "${GOOD_REVISION}=100"
```

Because images are tagged by git SHA, the previous revision's image is always
recoverable; `update-traffic` is instant and does not rebuild anything.

## Troubleshooting

| Symptom | Likely cause → fix |
|---|---|
| `exec format error` in Cloud Run logs | ARM image was pushed — rebuild with `--platform linux/amd64` |
| Container fails to start / `PORT` timeout | App not listening on `$PORT` on `0.0.0.0`, or crashes on boot — check `gcloud run services logs read` |
| `denied: Permission ... artifactregistry.repositories.uploadArtifacts` | Missing `roles/artifactregistry.writer` or registry auth not configured — rerun `gcloud auth configure-docker`, then check IAM |
| `403` calling the service | Service is IAM-gated (no `--allow-unauthenticated`) — invoke with an identity token: `curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" $URL` |
| Push is very slow from Mac | QEMU emulation makes AMD64 layers slow to build; the push itself is network-bound — build once, tag once, reuse layers |
