## Environment & Stack

- **Local by default.** Build, run, and test on the machine in front of you, the way the project's `docs/deploy/local.md` says. Nothing goes to a cloud unless the user says so.
- **GCP is the only cloud.** When the user says to use GCP, deploy, or push an image, ask for the defaults before doing anything — project ID, region, Artifact Registry repository, service name, and whether the service is public — never assume or invent them, never fall back to another project's value, and record them in the project's `docs/deploy/gcp.md` so the next session does not ask again. Never suggest another cloud provider.
- **Backend: Go** (`golang` skill): `cmd/<app>/` + `internal/`, `log/slog` JSON logging.
- **Frontend: TanStack Start + TypeScript + Tailwind**, fully contained in one `web/` or `ui/` directory inside the Go project (`frontend`, `typescript` skills). Shipped JS runs on plain V8: no Node-only APIs outside the server entry.
- **Machines:** macOS (Apple Silicon) runs Docker via Colima; Linux boxes use plain Docker. Remote boxes or VMs only on explicit request — ask for access details, never assume them.
- **Root steps on the operator's own Linux box:** passwordless `sudo -n` is available. When the operator says to run a root step AND the project's `.claude/settings.json` records the standing permission, run it — do not hand it back — look at the target before deleting or overwriting it, and record the exact command plus its verification (`du` before an `rm -r`) in the change entry. Without that record, hand the step back with the exact command instead.
- **Test locally first:** the project's tests, then a native `docker build` + `docker run` smoke test (`docker` skill) before anything goes to Google.
- **Build for Google as amd64, always:** `docker buildx build --platform linux/amd64 --load …`; an arm64 image on Cloud Run fails with `exec format error`.
- **Images live in GCP Artifact Registry:** `${REGION}-docker.pkg.dev/${PROJECT}/${REPO}/${SERVICE}:${GIT_SHA}`, tagged with the git short SHA (`latest` only as an alias). GitHub Container Registry only on explicit request.
- **Compute:** Cloud Run first (`cloud-run-deploy` skill); GKE (`gke-basics`) or Compute Engine when Cloud Run does not fit. Infrastructure as code with Terraform (`terraform` skill) when it is more than one service.
- **Auth:** `gcloud auth login` and `gcloud auth application-default login` (`gcloud-setup` skill). Secrets reach Cloud Run through Secret Manager (`--set-secrets`), never `--set-env-vars`.
