# Upgrades — the provider major checklist and the Terraform minors this skill relies on

## Bumping `hashicorp/google` to a new major

1. Move to the last minor of the current major first and get a clean `plan` (exit 0): a diff there is drift, not the upgrade.
2. Read the provider's upgrade guide for the new major (`website/docs/guides/version_N_upgrade.html.markdown` in https://github.com/hashicorp/terraform-provider-google) and list every item that names a resource in the root.
3. Change the pin (`~> N.0`), `terraform init -upgrade`, `terraform providers lock -platform=linux_amd64 -platform=darwin_arm64`, commit the lock file.
4. `terraform plan -detailed-exitcode`: exit 0 is the goal; a 2 is read line by line, and a replacement (`-/+`) of a stateful resource stops the upgrade until it is understood.
5. Apply in its own step, after the plan was read, never together with a resource change.

### Items for this stack

| Major | Change | Effect here |
|---|---|---|
| 7.0 | `google_project_service.disable_on_destroy` no longer defaults to `true` | set it explicitly (`false`) so a destroy never disables a shared API |
| 7.0 | import-ID validation is stricter; `enable_flow_logs` removed from subnetworks | `import` blocks use the full resource ID form from the doc page |
| 8.0 | `google_secret_manager_secret_version.secret_data_wo_version` is a **string** and **required** with `secret_data_wo` | every write-only version carries `secret_data_wo_version = "1"` |
| 8.0 | backend service and global forwarding rule default `load_balancing_scheme` moved `EXTERNAL` → `EXTERNAL_MANAGED` | only when a load balancer is in the root: pin the scheme explicitly before the bump |
| 8.0 | `google_cloud_run_v2_worker_pool` lost `custom_audiences`; IAP, ML Engine, Notebooks and BeyondCorp resources removed | not used by the standard app |

## Terraform minors the skill depends on

| Version | Feature | Used for |
|---|---|---|
| 1.5 | `import` blocks | adopting a resource created outside Terraform, plan-reviewable |
| 1.6 | `terraform test` (`*.tftest.hcl`) | module tests |
| 1.7 | `mock_provider` | tests without credentials |
| 1.11 | write-only arguments | `secret_data_wo` |
| 1.16 | `import` inside modules; `lifecycle { destroy = false }`; `store` in `terraform_data` for ephemeral values | leaving a resource in place while removing it from HCL |

OpenTofu 1.13 tracks the same language features; the `required_version` constraint is written for whichever binary the project uses, and the two are not mixed on one state.
