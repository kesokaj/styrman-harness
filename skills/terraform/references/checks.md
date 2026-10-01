# Checks — before every apply, with no CI

Five commands, run bare, each as its own call with its exit status read; the operator's shell is the pipeline.

| # | Command | Exit codes | Catches |
|---|---|---|---|
| 1 | `terraform fmt -check -recursive` | 0 clean, 3 files would change | formatting; run `terraform fmt -recursive` to fix |
| 2 | `terraform init -backend=false && terraform validate` | 0 valid, 1 error | syntax, unknown arguments, `for_each` with `count`, references to nothing; works offline once the provider is cached |
| 3 | `terraform plan -detailed-exitcode -out=.agents/scratch/tfplan` | 0 no diff, 1 error, 2 a diff | what the apply would do; the plan file holds every value in clear, so it lives in scratch |
| 4 | `tflint --init && tflint` | 0 clean, 2 issues, 1 error | provider-specific mistakes (invalid machine types, regions, enum values) before the API sees them |
| 5 | `trivy config --exit-code 1 .` | 0 clean, 1 findings (with `--exit-code 1`) | misconfiguration policy: public buckets, missing deletion protection, plaintext secrets |

Then `terraform apply .agents/scratch/tfplan` applies exactly the reviewed plan and nothing newer. `terraform init` (with the backend) is the one network step before a plan; `providers lock` after it when the lock file changed.

## `.tflint.hcl`

```hcl
plugin "google" {
  enabled = true
  version = "0.40.0"
  source  = "github.com/terraform-linters/tflint-ruleset-google"
}
```

`tflint --init` downloads the ruleset once (into `~/.tflint.d/plugins`); the version pin is a normal dependency bump to review. https://github.com/terraform-linters/tflint-ruleset-google

## `.gitignore` for `infra/`

```
.terraform/
*.tfstate
*.tfstate.*
*.tfvars
*.tfvars.json
*.auto.tfvars.json
tfplan
crash.log
```

Not ignored: `.terraform.lock.hcl`. It is committed, so a provider change is a reviewable diff, and it is locked for both machines of this stack with `terraform providers lock -platform=linux_amd64 -platform=darwin_arm64`, rerun after any provider bump. https://developer.hashicorp.com/terraform/language/files/dependency-lock

## `terraform test` for a reusable module

An app root is tested by its `plan`. A module reused by two roots earns one test file, run without credentials through a mocked provider (Terraform 1.7+):

```hcl
# tests/service.tftest.hcl
mock_provider "google" {}

variables {
  project = "example-project"
  region  = "europe-north1"
  app     = "demo"
  repo    = "apps"
  git_sha = "abc1234"
}

run "service_uses_runtime_sa" {
  command = plan

  assert {
    condition     = google_cloud_run_v2_service.app.template[0].service_account == google_service_account.run.email
    error_message = "the service must run as the dedicated runtime service account"
  }
}
```

`terraform test` from the module directory; `command = apply` only against a sandbox project the user named. https://developer.hashicorp.com/terraform/language/tests
