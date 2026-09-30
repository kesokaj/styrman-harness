---
owner: research
date: 2026-09-30
---

# Research brief: replacing the vendored `terraform` skill with `terraform-gcp`

Every version and behaviour claim was checked on 2026-09-30 against the URL cited. Local findings come from `skills/terraform/` in this checkout.

## 0. The audit, confirmed locally

- `SKILL.md` plus six `references/*.md` (4,484 lines), **no `LICENSE`** though `SKILL.md` § License promises one (`VENDORED.md` already records the debt).
- `references/code-patterns.md` 249-258: `resource "aws_instance" "app"` sets both `for_each` and `count` — invalid HCL.
- `references/module-patterns.md` 876 gitignores `.terraform.lock.hcl`; `references/quick-reference.md` 449-450 says `git add` it (§ 4 settles it).
- Zero GCP mentions; 169 `aws_` hits; CI, pre-commit or MCP in `SKILL.md` and five of six references.

## 1. Versions (seen 2026-09-30)

| Thing | Current | Source |
|---|---|---|
| Terraform | **v1.16.4** (2026-09-23); v1.17.0-beta2 (2026-09-23) pre-release. 1.16.0 (2026-08-26) added `import` blocks inside modules, `lifecycle { destroy = false }`, a `store` block in `terraform_data` for ephemeral values. | https://github.com/hashicorp/terraform/releases ; https://raw.githubusercontent.com/hashicorp/terraform/v1.16.0/CHANGELOG.md |
| OpenTofu | **v1.13.0** (2026-09-29/30: blog says Sep 29, GitHub release stamped Sep 30 UTC); 1.12.6 is the previous line. | https://github.com/opentofu/opentofu/releases ; https://opentofu.org/blog/ |
| `hashicorp/google` | **v8.5.0** (2026-09-29); **8.0.0 shipped 2026-08-26**. 7.0.0 (prior major): `google_project_service.disable_on_destroy` no longer defaults to `true`; stricter import-ID validation; `enable_flow_logs` removed from subnetworks. 8.0.0, touching this stack: `google_secret_manager_secret_version.secret_data_wo_version` is now a **string** and **required** with `secret_data_wo`; backend service / global forwarding rule default `load_balancing_scheme` moved `EXTERNAL` → `EXTERNAL_MANAGED`; `google_cloud_run_v2_worker_pool` lost `custom_audiences`; IAP/ML Engine/Notebooks/BeyondCorp resources removed. The § 2 resources are unchanged in shape. Neither guide states a minimum Terraform version. | https://github.com/hashicorp/terraform-provider-google/releases ; https://raw.githubusercontent.com/hashicorp/terraform-provider-google/main/CHANGELOG.md ; `.../website/docs/guides/version_8_upgrade.html.markdown` ; `.../version_7_upgrade.html.markdown` |
| `google` vs `google-beta` | "All generally available (GA) products and features are available in both versions of the provider"; `google-beta` adds Preview features only. Every § 2 resource is in the GA `google` docs tree: **`google-beta` is not needed** for Cloud Run v2. | https://raw.githubusercontent.com/hashicorp/terraform-provider-google/main/website/docs/guides/provider_versions.html.markdown ; `.../website/docs/r/cloud_run_v2_service.html.markdown` |
| Cloud SQL Postgres | PostgreSQL **18 is the default** for new instances (enum `POSTGRES_18`; 17 is `POSTGRES_17`). | https://docs.cloud.google.com/sql/docs/postgres/db-versions ; https://docs.cloud.google.com/sql/docs/postgres/admin-api/rest/v1/SqlDatabaseVersion |

## 2. Canonical minimal HCL

Assembled from the provider's own examples (`.../website/docs/r/<name>.html.markdown` under the raw URL in § 1); `var.project`, `var.region`, `var.app`, `var.repo`, `var.git_sha` are the values `dev-setup.md` says to ask the operator for.

### versions.tf and the GCS backend

```hcl
terraform {
  required_version = "~> 1.16"
  required_providers {
    google = { source = "hashicorp/google", version = "~> 8.0" }
  }
  backend "gcs" {
    bucket = "PROJECT-tfstate"   # literal: backend blocks take no variables
    prefix = "APP"
  }
}

provider "google" {
  project = var.project
  region  = var.region
}
```

- Backend auth is ADC by default, else `impersonate_service_account`, `access_token`, `credentials`; optional `kms_encryption_key` (CMEK). https://developer.hashicorp.com/terraform/language/backend/gcs
- **Locking needs no table.** The backend writes `<prefix>/<workspace>.tflock` into the same bucket under a `DoesNotExist` precondition (`lockFile.If(storage.Conditions{DoesNotExist: true}).NewWriter`) and deletes it with `GenerationMatch`, so only the holder releases it. Always on. https://raw.githubusercontent.com/hashicorp/terraform/main/internal/backend/remote-state/gcs/client.go
- The bucket is made by hand, once (`gcloud storage buckets create gs://PROJECT-tfstate --location=REGION --uniform-bucket-level-access --public-access-prevention`, then `... update --versioning`), matching Google's guide (versioning, uniform access, `public_access_prevention = "enforced"`, `force_destroy = false`). https://docs.cloud.google.com/docs/terraform/resource-management/store-state

### Artifact Registry

```hcl
resource "google_artifact_registry_repository" "app" {
  location               = var.region
  repository_id          = var.repo
  format                 = "DOCKER"
  cleanup_policy_dry_run = false
  cleanup_policies {
    id     = "delete-untagged"
    action = "DELETE"
    condition { tag_state = "UNTAGGED" }
  }
}
```
Required: `format`, `repository_id` (`location` optional but always set here). `.../r/artifact_registry_repository.html.markdown`

### Runtime service account, secret, IAM

```hcl
resource "google_service_account" "run" {
  account_id   = "${var.app}-run"
  display_name = "${var.app} Cloud Run runtime"
}

resource "google_secret_manager_secret" "db_url" {
  secret_id = "${var.app}-database-url"
  replication { auto {} }
}

resource "google_secret_manager_secret_version" "db_url" {
  secret                 = google_secret_manager_secret.db_url.id
  secret_data_wo         = var.db_url   # write-only: never lands in state
  secret_data_wo_version = "1"          # string, required with _wo (provider 8.0)
}

resource "google_secret_manager_secret_iam_member" "run_reads_db_url" {
  project   = google_secret_manager_secret.db_url.project
  secret_id = google_secret_manager_secret.db_url.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = google_service_account.run.member   # "serviceAccount:<email>"
}
```
`secret_id` and `replication` are the only required secret args. `secret_data` (plaintext in state) and `secret_data_wo` are mutually exclusive; `deletion_policy` is `DELETE|DISABLE|ABANDON|PREVENT`. Or create only the secret in HCL and add versions with `gcloud secrets versions add`. `.../r/secret_manager_secret.html.markdown`, `.../r/secret_manager_secret_version.html.markdown`, `.../r/secret_manager_secret_iam.html.markdown`, `.../r/google_service_account.html.markdown` (the `member` attribute). Write-only arguments exist since Terraform 1.11 (the vendored skill's feature table; not re-verified today).

### Cloud Run service (secret env, Cloud SQL socket, Direct VPC egress)

```hcl
resource "google_cloud_run_v2_service" "app" {
  name     = var.app
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"   # or INGRESS_TRAFFIC_INTERNAL_ONLY / _INTERNAL_LOAD_BALANCER
  # deletion_protection defaults to true; set false only for throwaway environments

  template {
    service_account = google_service_account.run.email
    scaling { max_instance_count = 10 }

    containers {
      image = "${var.region}-docker.pkg.dev/${var.project}/${var.repo}/${var.app}:${var.git_sha}"
      env {
        name = "DATABASE_URL"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.db_url.secret_id
            version = "latest"
          }
        }
      }
      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance { instances = [google_sql_database_instance.db.connection_name] }
    }

    vpc_access {                       # only for private-IP targets
      network_interfaces {
        network    = "default"
        subnetwork = "default"
      }
      egress = "PRIVATE_RANGES_ONLY"   # or ALL_TRAFFIC
    }
  }
}
```
- `volumes { cloud_sql_instance }` + the `/cloudsql` mount is the HCL form of `gcloud run deploy --add-cloudsql-instances=CONNECTION_NAME`; the runtime SA needs `roles/cloudsql.client`. https://docs.cloud.google.com/run/docs/configuring/connect-cloudsql
- **Direct VPC egress is GA and current**; a connector is the exception for Cloud NAT cold starts (30 s+): "for better startup performance, we recommend using Serverless VPC Access connectors with Cloud NAT". No private-IP target, no `vpc_access`. https://docs.cloud.google.com/run/docs/configuring/vpc-direct-vpc
- `launch_stage` defaults to GA. `.../r/cloud_run_v2_service.html.markdown`

### Public vs IAM-gated invocation

```hcl
# Public — Google's recommended form (no allUsers binding; blocked where the
# org policy constraints/run.managed.requireInvokerIam is enforced):
#   invoker_iam_disabled = true      # top-level arg on google_cloud_run_v2_service

# Public — the IAM form (blocked where constraints/iam.allowedPolicyMemberDomains is enforced):
resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = google_cloud_run_v2_service.app.project
  location = google_cloud_run_v2_service.app.location
  name     = google_cloud_run_v2_service.app.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# IAM-gated: same resource, member = google_service_account.scheduler.member
```
`google_project_iam_member { role = "roles/run.invoker" }` grants every service in the project — too broad for one app; keep that resource for project-scoped roles. "The recommended way to make a public service is to disable the Cloud Run Invoker IAM check." https://docs.cloud.google.com/run/docs/authenticating/public ; `.../r/cloud_run_v2_service_iam.html.markdown` ; `.../r/google_project_iam.html.markdown` (never mix authoritative `_policy` with `_member`). Provider issue #19786 asking for `invoker_iam_disabled` is closed via PR #19833; the argument is in the current docs.

### Job and Scheduler

```hcl
resource "google_cloud_run_v2_job" "migrate" {
  name     = "${var.app}-migrate"
  location = var.region
  template {
    template {                        # yes, template.template
      service_account = google_service_account.run.email
      containers {
        image = "${var.region}-docker.pkg.dev/${var.project}/${var.repo}/${var.app}:${var.git_sha}"
        args  = ["migrate"]
      }
    }
  }
}

resource "google_cloud_scheduler_job" "tick" {   # google_service_account.scheduler: account_id = "${var.app}-scheduler"
  name             = "${var.app}-tick"
  region           = var.region
  schedule         = "*/15 * * * *"
  time_zone        = "Etc/UTC"
  attempt_deadline = "320s"
  http_target {
    http_method = "POST"
    uri         = "${google_cloud_run_v2_service.app.uri}/tick"
    oidc_token {
      service_account_email = google_service_account.scheduler.email
      audience              = google_cloud_run_v2_service.app.uri   # Google: "copy the full URL of your service"
    }
  }
}
```
The scheduler SA gets `roles/run.invoker` on the service (the IAM-gated member above). Job `deletion_protection` also defaults to true. https://docs.cloud.google.com/run/docs/triggering/using-scheduler ; `.../r/cloud_run_v2_job.html.markdown` ; `.../r/cloud_scheduler_job.html.markdown`

### Cloud SQL (Postgres)

```hcl
resource "google_sql_database_instance" "db" {
  name                = "${var.app}-db"
  database_version    = "POSTGRES_18"
  region              = var.region
  deletion_protection = true
  settings {
    tier    = "db-f1-micro"          # ENTERPRISE tiers; ENTERPRISE_PLUS needs performance-optimized tiers
    edition = "ENTERPRISE"           # PG 16+ defaults to ENTERPRISE_PLUS without this
    ip_configuration { ipv4_enabled = true }   # the /cloudsql socket works over public IP with no VPC
  }
}

resource "google_project_iam_member" "run_sql_client" {
  project = var.project
  role    = "roles/cloudsql.client"
  member  = google_service_account.run.member
}
```
Add `google_sql_database` for the database. Private IP instead: `ip_configuration { ipv4_enabled = false, private_network = <vpc self_link> }` plus the service's `vpc_access`. A `google_sql_user` password in HCL lands in state; prefer IAM database authentication. `.../r/sql_database_instance.html.markdown`

## 3. Auth for a solo operator with no CI

**Recommend `gcloud auth application-default login` (user credentials), with `--impersonate-service-account` as the optional step up. Not Workload Identity Federation.**

- Google: "When your code is running in a local development environment, such as a development workstation, the best option is to use the credentials associated with your user account"; "Service account keys are a security risk … choose a more secure alternative … whenever possible." https://docs.cloud.google.com/docs/authentication/set-up-adc-local-dev-environment
- WIF is for workloads outside Google Cloud ("AWS or Azure; on-premises Active Directory; deployment services, such as GitHub and GitLab"), never a human's workstation; with no CI there is nothing to federate. Workforce Identity Federation is the human variant, for an external IdP only. https://docs.cloud.google.com/iam/docs/workload-identity-federation
- Backend and provider both read ADC, so one login covers state and apply. Impersonating a dedicated `terraform` SA (`--impersonate-service-account`, or `impersonate_service_account` in the backend) narrows the operator's live roles and unifies the audit trail; worth it once more than one human applies, not before. Matches `gcloud-setup` and `dev-setup.md` § Auth as written.

## 4. `.terraform.lock.hcl`

Commit it. HashiCorp: "You should include this file in your version control repository so that you can discuss potential changes to your external dependencies via code review, just as you would discuss potential changes to your configuration itself." For the two machines in `dev-setup.md`: `terraform providers lock -platform=linux_amd64 -platform=darwin_arm64`, rerun after any provider bump; never commit `.terraform/`. https://developer.hashicorp.com/terraform/language/files/dependency-lock — so `module-patterns.md`'s `.gitignore` is the wrong half of the contradiction.

## 5. Testing without CI

| Tool | State (seen 2026-09-30) | Source |
|---|---|---|
| `terraform fmt -check`, `terraform validate` | built in; validate needs `init` (`-backend=false` works offline) | https://developer.hashicorp.com/terraform/cli/commands/plan |
| `terraform plan -detailed-exitcode -out=tfplan` | 0 = no diff, 1 = error, 2 = diff; the plan file holds secrets in clear, keep it in `.agents/scratch/` | same |
| `terraform test` (`*.tftest.hcl`) | since 1.6.0; `command = plan` or `apply`; `mock_provider` since 1.7.0; stable, though the docs never say "GA" | https://developer.hashicorp.com/terraform/language/tests |
| tflint v0.64.0 (2026-07-17) + `tflint-ruleset-google` **v0.40.0 (2026-09-23)**, "100+ rules" | `.tflint.hcl`: `plugin "google" { enabled = true; version = "0.40.0"; source = "github.com/terraform-linters/tflint-ruleset-google" }`, then `tflint --init` | https://github.com/terraform-linters/tflint/releases ; https://github.com/terraform-linters/tflint-ruleset-google (README, releases) |
| trivy v0.74.0 (2026-08-14), `trivy config .` | Go binary, tfsec's successor, also scans the image | https://github.com/aquasecurity/trivy/releases |
| checkov 3.3.20 (2026-09-27) | Python; overlaps trivy | https://github.com/bridgecrewio/checkov/releases |

**Recommend two add-ons over the built-in three (`fmt`, `validate`, `plan -detailed-exitcode`): tflint with the google ruleset, and `trivy config`.** Both are single Go binaries; tflint catches provider-specific mistakes before a plan touches the API; trivy covers policy with the binary that already scans the image. `terraform test` earns one paragraph for reusable modules only; an app root is tested by `plan`. Drop checkov.

## 6. Licence and alternatives

- Upstream `LICENSE`: **Apache License 2.0, "Copyright 2026 Anton Babenko"** (https://raw.githubusercontent.com/antonbabenko/terraform-skill/master/LICENSE). Apache-2.0 § 4 binds redistribution of the work *or a derivative*: licence text, attribution, modified files marked. A **clean rewrite that copies no upstream text owes nothing**; the current condensed `SKILL.md` does. Recommendation: write `terraform-gcp` from scratch (file layout, block order, `~>` pinning, committing the lock file are common conventions, not Babenko's text) and **delete `skills/terraform/`** with its `VENDORED.md` row, retiring the missing-`LICENSE` debt. If any sentence is carried over, add `LICENSE` beside it and a `VENDORED.md` row.
- **Google publishes no Terraform skill.** `google/skills` (Apache-2.0) has ~200 entries under `skills/cloud/`, none named terraform/iac/secret/artifact. `cloud-run-basics` (becc4b8, 2026-09-21) ships `references/iac-usage.md`: ~90 lines of `google_cloud_run_v2_service`, `_iam_member`, `_job`, `_worker_pool` with no service account, secrets or VPC — the shape of `gke-basics/references/iac-usage.md`. Not a substitute. https://api.github.com/repos/google/skills/contents/skills/cloud ; https://raw.githubusercontent.com/google/skills/main/skills/cloud/cloud-run-basics/references/iac-usage.md
- **`hashicorp/agent-skills`** (MPL-2.0) has 16 Terraform skills; the provider-neutral `terraform-style-guide` is ~350 lines, every example AWS, no CI/pre-commit/MCP. Vendoring it verbatim is legal (ship its LICENSE) but says in 350 AWS lines what 30 GCP lines can; do not vendor. https://github.com/hashicorp/agent-skills ; https://raw.githubusercontent.com/hashicorp/agent-skills/main/SKILLS.md

## Proposed skill outline: `terraform-gcp`

Clean-room, GCP only, no CI, no MCP, no pre-commit; every HCL sample comes from a provider doc page and cites it; no `VENDORED.md` row (nothing vendored).

**`SKILL.md` frontmatter**

```yaml
name: terraform-gcp
description: "Terraform (or OpenTofu) for Google Cloud: Cloud Run v2 with a dedicated service account, Artifact Registry, Secret Manager, Cloud SQL Postgres, Scheduler, Jobs, state in GCS, applied from the operator's machine. Use when writing, reviewing or planning HCL for GCP, choosing a resource or IAM binding, debugging state or a plan diff, or deciding what belongs in Terraform versus gcloud. Do NOT use for building or pushing the image (docker), a gcloud deploy or rollback (cloud-run-deploy), GKE clusters or Kubernetes objects (gke-basics), or gcloud login and ADC problems (gcloud-setup)."
```

**`SKILL.md` sections (target ~180 lines)**

1. *When Terraform, when gcloud* — `dev-setup.md`'s rule: Terraform when it is more than one service; image, new revision and rollback stay with `cloud-run-deploy`; `plan` before every apply; nothing goes to Google unasked.
2. *Layout* — one root per app under `infra/`: `versions.tf`, `main.tf`, `variables.tf`, `outputs.tf`, `terraform.tfvars` (gitignored, secret-shaped), `.terraform.lock.hcl` committed; the state bucket recorded in `docs/deploy/gcp.md`.
3. *Bootstrap once* — `google_project_service` for run, artifactregistry, secretmanager, sqladmin, cloudscheduler, iam; the state bucket by `gcloud` (chicken-and-egg); `init`; `providers lock` for both platforms.
4. *The standard app* — the § 2 snippets in order: SA → Artifact Registry → secret + IAM → Cloud SQL → service → invoker binding (public via `invoker_iam_disabled`, gated via `_iam_member`) → job → scheduler; each with the one-line "why" and the doc URL.
5. *Rules* — `~> 8.0` on `hashicorp/google`, no `google-beta` unless a Preview feature is named; `deletion_protection` at its default in anything shared; secrets never in `secret_data`, `--set-env-vars` or state; IAM by `_member`, never `_policy`; project-level bindings only for project-scoped roles (`cloudsql.client`); Direct VPC egress, connector only for the NAT cold-start case; image tag is the git SHA from `cloud-run-deploy`, passed as a variable; block order (`count`/`for_each` first, arguments, `labels`, `depends_on`, `lifecycle`); `for_each` for named sets, `count` for a boolean toggle, never both.
6. *Verify before apply* — `fmt -check`, `validate`, `plan -detailed-exitcode -out=.agents/scratch/tfplan` (2 = review the diff, 0 = nothing to do), `tflint --init && tflint`, `trivy config .`; `apply .agents/scratch/tfplan`; results into the project's `testing/` log per the Definition of Done.
7. *Auth* — § 3 in four lines, pointing at `gcloud-setup` for failures.
8. *Debugging state* — stuck lock (`<prefix>/default.tflock`; `force-unlock` only after confirming no run holds it), drift (`plan`, `import` block), `moved` for renames, `removed` / `lifecycle { destroy = false }` (1.16) to leave a resource in place.

**Reference files (three at most)**

- `references/resources-gcp.md` — the full copy of every § 2 snippet with the doc URL under each, plus the two private-IP variants (Cloud SQL `private_network`, service `vpc_access`) and a worker-pool note; ~250 lines. The single place HCL lives, so `SKILL.md` can stay short.
- `references/checks.md` — `.tflint.hcl` (ruleset google 0.40.0), the five commands with exit codes, the `.gitignore` (`.terraform/`, `*.tfstate*`, `*.tfvars`, `tfplan`; **not** the lock file), and one `terraform test` example with `command = plan` and `mock_provider "google"`; ~100 lines.
- `references/upgrades.md` — the provider-major checklist (last minor first, read the guide, clean `plan`) with the 7.0 and 8.0 items from § 1, and the Terraform minors the skill relies on (1.6 test, 1.7 mocks, 1.11 write-only, 1.16 `destroy = false`); ~60 lines.

**Then**: delete `skills/terraform/`, remove its `VENDORED.md` row, and repoint every description saying "(use terraform)" — `cloud-run-deploy`, `docker`, `frontend`, `gcloud-setup`, `python` — plus `dev-setup.md`'s "`terraform` skill" to `terraform-gcp` in the same commit (one vocabulary, one edit).
