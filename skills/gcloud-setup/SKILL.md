---
name: gcloud-setup
description: |
  GCloud CLI authentication and project configuration. Use when encountering "ADC not found",
  "gcloud auth error", "No credentialed accounts", "Application Default Credentials", "project not set",
  "permission denied", a 403 from gcloud, bq, or a client library, or when setting up gcloud CLI
  for a new machine or project.
  Activate for any gcloud auth, gcloud config, or GCP credential question.
  Do NOT use for GKE cluster management (use gke-basics) or Terraform GCP resources (use terraform).
---

# GCloud CLI Setup & Authentication

Quick-reference for configuring `gcloud` CLI authentication and project settings (macOS or Linux).

**When NOT to use this skill:** GKE cluster management (use the `gke-basics` skill), Terraform GCP resources (use the `terraform` skill), or Cloud Run deployments (use the `cloud-run-deploy` skill).

## Quick Fixes

| Problem | Fix |
|---------|-----|
| ADC not found | `gcloud auth application-default login` |
| Project not set | `gcloud config set project PROJECT_ID` |
| Permission denied (403 with an account listed) | `gcloud projects get-iam-policy PROJECT_ID --flatten='bindings[].members' --filter='bindings.members:ACCOUNT' --format='value(bindings.role)'` shows the roles the account holds; compare with the role the call needs |
| Token expired | `gcloud auth application-default login` (re-run) |
| Quota exceeded | `gcloud compute project-info describe --project PROJECT_ID --format='yaml(quotas)'` for Compute quotas; other APIs: `gcloud services quota list --service=SERVICE.googleapis.com --consumer=projects/PROJECT_ID` (check `--help` for the exact verb on the installed version) |
| Wrong account | `gcloud auth login` (re-authenticate) |

## First-Time Setup

### 1. Install gcloud CLI

```bash
# macOS (Homebrew cask; `google-cloud-sdk` is the former token)
brew install --cask gcloud-cli
```

Linux (Debian/Ubuntu), the apt package — a root step under the standing-permission rule
(`context/dev-setup.md`), otherwise handed back with these exact commands:

```bash
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg
echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee /etc/apt/sources.list.d/google-cloud-sdk.list
sudo apt-get update && sudo apt-get install google-cloud-cli
```

The component manager is disabled in the apt package: extras are apt packages
(`google-cloud-cli-kubectl`, …), not `gcloud components install`. The no-root path is the
tarball: `curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-cli-linux-x86_64.tar.gz`,
check its SHA256 against the value printed on the install page (`sha256sum google-cloud-cli-linux-x86_64.tar.gz`) before `tar -xzf` and `./google-cloud-sdk/install.sh`, which keeps the component manager; the archive is fetched over TLS but carries no signature of its own.

### 2. Authenticate

```bash
# Login with your Google account (opens browser)
gcloud auth login

# Set up Application Default Credentials (required for SDKs and APIs)
gcloud auth application-default login
```

### 3. Set Default Project

```bash
# List available projects
gcloud projects list

# Set default project
gcloud config set project YOUR_PROJECT_ID

# Verify
gcloud config get project
```

### 4. Verify Setup

```bash
# Check authenticated accounts (* = active)
gcloud auth list

# Check current project
gcloud config get project

# Test API access
gcloud services list --enabled --limit=5
```

## Common Configurations

### Set Default Region/Zone

The region is the operator's to give, never assumed and never carried over from another project (`context/dev-setup.md`); a project records its own in `docs/deploy/gcp.md`.

```bash
gcloud config set compute/region "$REGION"      # e.g. europe-north1 — ask, do not guess
gcloud config set compute/zone "$REGION-a"
```

### Multiple Configurations (Named Profiles)

```bash
# Create a named config for a different project
gcloud config configurations create my-staging
gcloud config set project staging-project-id
gcloud config set compute/region europe-west1

# Switch between configs
gcloud config configurations activate default
gcloud config configurations activate my-staging

# List all configs
gcloud config configurations list
```

### Artifact Registry Authentication

```bash
# Configure Docker to push to Artifact Registry
gcloud auth configure-docker REGION-docker.pkg.dev
```

## Before asking the user to log in

1. Verify first: `gcloud auth list`. `No credentialed accounts.` means missing authentication; an account that IS listed plus a `403` means IAM permissions on that account, not authentication — investigate the roles instead of asking for a new login. The same applies to `bq`, client libraries and notebooks (`DefaultCredentialsError`).
2. Never fix credentials by code or alternative tools; the user runs `gcloud auth login` and `gcloud auth application-default login` in their own terminal.
3. After they confirm, `gcloud auth list` again, then continue the original task.

The project ID is never discovered from `gcloud config` or `.env`: the global context has the agent ask for the project's GCP defaults once and record them in the project's `docs/deploy/gcp.md`.

## Troubleshooting

### ADC Location

Application Default Credentials are stored at:
```
~/.config/gcloud/application_default_credentials.json
```

To reset: delete that file and re-run `gcloud auth application-default login`.

### Service account keys — last resort, not the path here

**The harness runs no CI (`context/AGENTS.md`), and a downloaded key file is a credential on disk that outlives the session.** Use `gcloud auth login` plus `gcloud auth application-default login` locally, and attached service-account identity on Cloud Run or GCE. Only when a platform offers neither, and the operator says so:
```bash
# Activate service account
gcloud auth activate-service-account --key-file=key.json

# Set ADC from service account
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json
```

> **Never commit service account keys to git.** Use Workload Identity for GKE, or Secret Manager for other compute.
