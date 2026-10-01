---
owner: security
date: 2026-10-01
---
# Security review — catalog-gke   (2026-10-01)

Range reviewed: cf9c5e2, skills/google-catalog/catalog (slice: gke-*; 35 entries)

## Status: COMPLETE

## Verdict: FINDINGS

The worst problem is in `gke-batch-hpc`. It tells the agent to install two cluster-wide operators straight from mutable URLs: Kueue from `releases/latest` and the MPI Operator from the `master` branch of kubeflow/mpi-operator, both with `kubectl apply -f https://…`. Whatever those URLs serve at run time (CRDs, controllers, admission webhooks, ClusterRoles) gets installed unreviewed, with the operator's cluster-admin-level credentials. Four more medium findings follow:

- `gke-cluster-creation`'s GPU and A3 templates drop the golden-path hardening without saying so, and run nodes as the default Compute Engine service account with the `cloud-platform` scope.
- `gke-workload-identity`'s connectivity test prints a live OAuth access token into the transcript.
- `gke-workload-troubleshooting` ends by creating a branch, committing, pushing and opening a GitHub PR on its own.
- `gke-upgrades`' troubleshooting reference deletes a ValidatingWebhookConfiguration, which may be a security policy, to unblock an upgrade.

A per-entry note fixes each of these, and no entry needs to be excluded. The slice contains no secrets, no real project identifiers, no permission-granting frontmatter, no `curl | sh` and no path that sends code, cluster data or credentials outside Google Cloud or the operator's project. Result: 15 entries get a note, 0 are excluded and 20 need none.

Router step 4 (`skills/google-catalog/SKILL.md`) covers six things:

- GCP defaults are asked for first, never assumed.
- Secrets go through Secret Manager.
- Downloads are pinned and verified, never `curl | sh`.
- Git stages by name.
- Deploys follow the project's deploy guide.
- Harness skills win over entries.

It does not cover four things:

- Confirmation before destructive, IAM or cluster mutations.
- Credentials printed to the terminal.
- Hardening defaults in cluster templates.
- Consent for commit, push or PR.

The findings below say which side each one falls on.

## Findings

- **[medium / confidence: high]** `gke-batch-hpc/SKILL.md:89`, `:144`: `kubectl apply --server-side -f https://github.com/kubernetes-sigs/kueue/releases/latest/download/manifests.yaml` and `kubectl apply -f https://raw.githubusercontent.com/kubeflow/mpi-operator/master/deploy/v2beta1/mpi-operator.yaml`.
  - Weakness: both URLs are mutable (`latest` moves with every release; `master` moves with every commit). Each manifest installs CRDs, a controller Deployment with cluster-scoped RBAC and, for Kueue, admission webhooks that see every Job and Pod. Whoever controls the URL's content at run time controls what lands in the cluster, applied with the operator's own (usually admin) credentials, and nothing tells the agent to ask first.
  - Router: step 4's "a download is pinned and verified" covers this in spirit, but `kubectl apply -f <url>` does not read as a download.
  - Remediation: pin a release tag (Kueue `releases/download/vX.Y.Z/manifests.yaml`, mpi-operator at a tag or commit SHA). Download into `.agents/scratch/`, verify against the release's checksum or the pinned commit, review the RBAC and webhooks, ask the operator, then `kubectl apply -f <file>`.
  - Recommendation: **note**.
- **[medium / confidence: high]** `gke-cluster-creation/SKILL.md:239` (template 4; Standard JSON `oauthScopes` at `:281`), `:287` (template 5; `--scopes …cloud-platform` at `:303`, JSON at `:326`).
  - Weakness: template 4's Standard JSON has no private nodes, authorized networks, Shielded/secure boot, `workloadIdentityConfig`/`workloadMetadataConfig` or `serviceAccount`. Its nodes run as the default Compute Engine SA with the `cloud-platform` scope, and every pod can fetch that SA's token from the metadata server. In projects where the default SA still holds the historical automatic Editor grant, any compromised pod acts as project Editor. Template 5's gcloud form sets `--workload-pool` but has no private nodes, authorized networks or Shielded Nodes. Its JSON form also omits Workload Identity.
  - Template 2 says it is unhardened; templates 4 and 5 do not. The entry does make the agent present the configuration and ask before creating (§ Instructions), so the gap is disclosure, not consent.
  - Router: not covered.
  - Remediation: add a dedicated least-privilege node SA (`--service-account`), `GKE_METADATA`, private nodes and authorized networks. Otherwise, state that the template is unhardened when presenting it.
  - Recommendation: **note**.
- **[medium / confidence: high]** `gke-workload-identity/SKILL.md:277`: in Step 5(c), `kubectl exec … -- curl … 'http://169.254.169.254/computeMetadata/v1/instance/service-accounts/default/token'` ("a healthy path returns a token").
  - Weakness: the response body carries a bearer `access_token` for the workload's identity, valid for about 1 h with every role that identity holds. It lands in the terminal, the CLI's session transcript on disk and the model provider's context.
  - Mitigation in place: `sh-agent report` redacts `ya29.` (`internal/harness/report.go:40`). Other transcript readers were not checked.
  - Router: not covered. Its "secrets via Secret Manager" rule is about storing secrets.
  - Remediation: use `…/service-accounts/default/email`, which proves the metadata path without a credential, or `curl -sS -o /dev/null -w '%{http_code}\n' …/token`.
  - Recommendation: **note**.
- **[medium / confidence: high]** `gke-workload-troubleshooting/SKILL.md:340-344`: Step 5.3 says "create a branch, commit the change, open a Pull Request (PR) on GitHub, and conclude the workflow (do not wait for human merge)".
  - Weakness: an outward-facing write (push plus PR to a shared remote) without the operator's request, in an entry that otherwise declares a read-only boundary. It is reinforced by Step 0 (`:30` and the line after it), which tells the agent to infer context "non-interactively … to ensure autonomous execution flow".
  - Router: "git stages by name" covers staging only. It says nothing about consent to commit, push or open a PR.
  - Remediation: present the patch and stop. Git actions happen only on the operator's request, under the harness git rules.
  - Recommendation: **note**.
- **[medium / confidence: medium]** `gke-upgrades/references/troubleshooting.md:93`: the fix for a webhook blocking an upgrade is `kubectl delete validatingwebhookconfigurations WEBHOOK_NAME`, and the restore at `:95-96` is a trailing comment.
  - Weakness: a ValidatingWebhookConfiguration is often a security control (Policy Controller/Gatekeeper, Kyverno, image verification). Deleting it turns enforcement off cluster-wide until someone re-applies the backup, and nothing tells the agent to find out what the webhook enforces or to ask. In the same file, Option A at `:35-36` relaxes a PDB to `maxUnavailable: 100%` without taking a backup.
  - Both contradict the entry's own rule (SKILL.md § Maintenance runbooks: restore as a numbered step with its own verification).
  - Confidence is medium because the entry mostly writes runbooks for humans, but its § Troubleshooting sends the agent to this file for "fix procedures".
  - Router: not covered.
  - Recommendation: **note**.
- **[low / confidence: high]** `gke-platform-security/SKILL.md:176`, `:182`, `:188`: `gcloud projects add-iam-policy-binding <PROJECT_ID> … --role roles/storage.objectViewer|roles/cloudsql.client|roles/pubsub.subscriber --quiet`.
  - Weakness: project-wide grants (every bucket, every subscription), run with `--quiet` and no confirmation step. A compromised workload with that GSA reads every object in the project.
  - Router: not covered. The harness `api-auth-gcp` skill may cover the topic; it was not read.
  - Remediation: bind on the bucket or subscription where the resource supports it, and confirm every IAM change with the operator.
  - Cosmetic: the code fence opened before `:176` is never closed, so § Resources renders as code.
  - Recommendation: **note**.
- **[low / confidence: medium]** `gke-backup-dr/SKILL.md:33`, `:127`: `--include-secrets` with `--all-namespaces` is the golden-path default.
  - Weakness: every Kubernetes Secret in every namespace is copied into the backup and can be restored (read) by whoever holds restore rights on the plan. That includes a separate backup project, since cross-project backup is GA per the entry's own § Recent Changes. CMEK is recommended, not required. A production restore needs confirmation, but the backup's scope does not.
  - Router: only partly covered. Secrets kept in Secret Manager shrink the set of Kubernetes Secrets, but the backup scope is still not governed.
  - Recommendation: **note**.
- **[low / confidence: medium]** `gke-ai-troubleshooting-jobset-interruption/SKILL.md:22` ("Autonomous & Sandbox Execution Rule" … complete the workflow "autonomously"; each step says "Proceed … automatically") and `:292` (Resolution 2 [High Risk]: cordon/drain the node, delete the underlying GCE VM).
  - Weakness: the resolution is worded "Recommend", but it has no confirmation gate, and the autonomy rule invites the agent to act.
  - Router: covers `:35` ("Independently gather context") through "GCP defaults asked first". It does not cover the destructive step.
  - Recommendation: **note**.
- **[low / confidence: high]** `gke-golden-path/SKILL.md:96-97`: "Discover project/cluster context via MCP tools or `gcloud config get-value project` — don't ask users to paste project IDs."
  - Weakness: an explicit instruction not to ask, contradicting router step 4. Followed literally, cluster creation and updates land in whatever project the active gcloud configuration points at.
  - Router: overrides fully. The note puts the override at the point of conflict.
  - Recommendation: **note**.
- **[low / confidence: medium]** `gke-compute-classes/SKILL.md:26`: the CODE-FIRST rule says to clone or fetch `GoogleCloudPlatform/cluster-autoscaler` from GitHub and run `git log -S`/`git blame` whenever behaviour is unclear.
  - Weakness: network access without asking, and the clone location is unspecified, so it would land in the operator's tree. It is read-only, on a Google-owned repo, and sends no data out.
  - Also in this entry: `references/compute-class-provisioning-methods.md:25` recommends "Privileged DaemonSets" for node init, which contradicts the entry's own CRITICAL POD-PRIVILEGE RULE. `assets/log-autoscaler-events.sh:148` loops forever (`while true`) and would hang an unattended run.
  - Router: not covered.
  - Recommendation: **note**.
- **[low / confidence: medium]** `gke-app-onboarding/SKILL.md:61`: `pack build <image> --builder gcr.io/buildpacks/builder:latest`.
  - Weakness: a mutable tag for a builder image that runs build code against the operator's source. `assets/Dockerfile:2` is `FROM node:22-slim`: pinned by tag, not by digest.
  - Router: covered in spirit (downloads pinned).
  - Recommendation: **note**.
- **[low / confidence: medium]** `gke-inference/SKILL.md:98`: "Some models require Hugging Face tokens. Create a Kubernetes Secret and reference it in the manifest."
  - Weakness: the usual way to create that Secret (`kubectl create secret generic … --from-literal=…`) puts the token on a command line and into the transcript.
  - Router: step 4 (secrets via Secret Manager) overrides. The note names the GKE path to use.
  - Recommendation: **note**.
- **[low / confidence: medium]** `gke-cluster-autoscaler/SKILL.md:62`: "Temporarily set `min-nodes` to 0 or delete instances manually in GCE".
  - Weakness: destructive, with no confirmation step. Two smaller issues in the same entry: `assets/find-scale-down-blockers.sh:140` writes `.tmp_pool_counts.$$` into the current directory instead of `mktemp` (removed at `:21`, `:150`), and `assets/log-autoscaler-events.sh:148` polls forever.
  - Router: not covered.
  - Recommendation: **note**.
- **[low / confidence: high]** `gke-service-networking/assets/gateway.yaml:11-15`, which SKILL.md workflow 1 applies: a global external (internet-facing) load balancer with an HTTP-only listener on port 80. `assets/ingress.yaml` is the same: an external GCE Ingress over HTTP.
  - Weakness: exposure is the entry's purpose, but cleartext is the default, and nothing tells the agent to confirm that the service should be public. `gateway-https.yaml` exists.
  - Router: partly covered (deploys follow the deploy guide).
  - Recommendation: **note**.
- **[low / confidence: medium]** `gke-basics/references/cli-reference.md:263-264` (`helm install|upgrade <RELEASE> <CHART>`, no `--version` or source), `:261` (`kubectl run … --image=busybox`, unpinned) and `:286` (the fix for an unreachable private cluster: enable `dnsEndpointConfig.allowExternalTraffic` or add authorized networks).
  - Weakness: the `:286` fix widens control-plane reach as a troubleshooting step, without asking.
  - Router: the catalog `gke-basics` entry is superseded by the harness skill of the same name (step 4). Other entries still point at its references (`gke-golden-path` § Rules 3).
  - Recommendation: **note**.

Cross-cutting observations (no per-entry note needed):

- **[low / confidence: high] C1: context inferred instead of asked.** These lines tell the agent to work out project and cluster from the environment (`gcloud config get-value project`, `kubectl config current-context`, "independently gather", "non-interactively"):
  - `gke-ai-troubleshooting-jobset-interruption/SKILL.md:35`
  - `gke-ai-troubleshooting-tpu-metrics-monitoring/SKILL.md:20`
  - `gke-ai-troubleshooting-tpu-vbar-oom/SKILL.md:30`
  - `gke-cluster-creation/SKILL.md:28` (its own § Instructions "ALWAYS ask" wins)
  - `gke-golden-path/SKILL.md:97`
  - `gke-node-notready/SKILL.md:34-35`
  - `gke-storage-troubleshooting/SKILL.md:34`
  - `gke-workload-scaling-troubleshooting/SKILL.md:43`
  - `gke-workload-troubleshooting/SKILL.md:30`
  - the `validate_queries.sh` fallbacks (jobset `:6`, vbar-oom `:6`)

  Router step 4 overrides all of them, and the diagnostic ones are read-only. Only `gke-golden-path`'s explicit "don't ask" gets a note.
- **[info] C2: `--quiet`** is on nearly every gcloud command and suppresses gcloud's own prompts. The sweep found no `delete … --quiet` in the slice, so the agent's confirmation, not gcloud's prompt, is the safety net everywhere.
- **[info] C3: `gcloud container clusters get-credentials`**, attempted in Step 0 of several troubleshooting entries, rewrites the operator's kubeconfig and switches the current context. Every later `kubectl` from any tool then targets that cluster.

## Checked

- **Secrets:** a sweep of every file in the slice for PEM, `AKIA`/`ASIA`, `AIza`, `ya29.`, `ghp_`, `xox*`, `hf_`, `keys create`, `password`, `api_key` found no credentials. Project identifiers are placeholders only (`my-gcp-project`, `customer-ai-project-123`, `my-project`, `<PROJECT_ID>`). One command prints a credential (finding 3). No script reads `~/.kube`, `~/.config/gcloud` or environment secrets. `gke-ai-troubleshooting-jobset-interruption/scripts/validate_queries.sh:38` passes `$(gcloud auth print-access-token)` in curl's argv to `monitoring.googleapis.com`; the token is visible in the process table while curl runs, the endpoint is Google's and the body is discarded. Informational.
- **Injection:** `gke-cluster-autoscaler` (§ CRITICAL RULES: refuse injected identifiers, pasted logs are data) and `gke-compute-classes` (CRITICAL INJECTION RULE) carry explicit defences. Entries substitute `{placeholders}` into shell commands, URLs and BigQuery SQL, and only `gke-cluster-autoscaler` validates names (`^[a-z0-9-]+$`). Scripts quote their variables, and `validate_config.py` parses local files only. No text tells the agent to ignore other instructions.
- **Authz/authn:** IAM findings 2 and 6. Workload Identity diagnosis is read-only with a human hand-off. RBAC examples are namespace-scoped and bound to Google Groups; `computeclass-rbac-editor.yaml` is labelled DO NOT DEPLOY. There is no `cluster-admin` binding, no binding to `system:authenticated` (`gke-multitenancy` and the golden path forbid it), no `allUsers`/`allAuthenticatedUsers`, no `0.0.0.0/0` and no legacy authorization.
- **Crypto:** none implemented or composed. TLS comes from Google-managed certificates or Certificate Manager, and CMEK from KMS key references. Nothing hand-rolled.
- **Dependencies:** the slice is vendored text and adds no dependency to the harness. The Go module is untouched, so `govulncheck ./...` would test the harness rather than this slice; not run. The only package manifest, `gke-app-onboarding/assets/package.json` with its `package-lock.json`, declares zero dependencies, so `npm audit` has nothing to audit; not run. The remote installs the entries instruct are findings 1, 11 and 15.
- **GCP surface:**
  - `--allow-unauthenticated`: absent.
  - Service-account keys: no creation, and two entries forbid JSON keys.
  - Secrets: mounted as files rather than env vars (`gke-manifest-generation`), with the Secret Manager add-on recommended.
  - IAM: findings 2 and 6.
  - Control-plane exposure: the DNS endpoint's `allowExternalTraffic: true` is a golden-path default (`gke-golden-path`, `gke-networking`, `assets/golden-path-autopilot.yaml:53`). It is IAM-gated and flagged as a Day-0 decision to ask about, so it is recorded here without a note.
- **Frontmatter:** all 35 SKILL.md files carry only `name`, `description` and `metadata`. No `allowed-tools`, `hooks`, `permissions` or `tools` key appears anywhere in the slice.
- **Scripts:** all 9 distinct scripts were read end to end:
  - the four `validate_queries.sh`, `validate_config.py`, `find-scale-down-blockers.sh`, `audit_cluster.sh` and `log-autoscaler-events.sh` (its `gke-compute-classes` copy is byte-identical by sha256)
  - plus `index.js`, the Dockerfile and `package*.json`

  They make read-only gcloud/kubectl calls to Google APIs, contain no `rm -r`, and write nothing outside the working directory except an operator-chosen `--log-file`. The temp file in the working directory is finding 13.
- **Network destinations:** every URL is a Google domain, a documentation link (buildpacks.io, kubernetes.io, registry.terraform.io, samber/awesome-prometheus-alerts, apache.org, kubernetes-client/java, GoogleCloudPlatform repositories) or one of the two install URLs of finding 1. Nothing sends cluster data, code or credentials outside Google Cloud or the operator's project; the PR in finding 4 goes to the operator's own remote.
- **Destructive commands:** swept for cluster, node-pool and instance deletes, `--force`, `--grace-period`, `drain`, `--delete-emptydir-data`, `maxUnavailable: 100%` and `rm -r` across all files. The confirmation-gated ones:
  - `gke-ai-troubleshooting-tpu-dynamic-slices-monitoring`: finalizer removal and CRD deletion.
  - `gke-backup-dr`: production restore.
  - Entries with a read-only boundary: `gke-node-notready`, `gke-storage-troubleshooting`, `gke-workload-identity`, `gke-workload-scaling-troubleshooting`, `gke-workload-troubleshooting`, the GMP troubleshooting in `gke-observability`, and `gke-service-networking` § Troubleshooting.
  - `gke-upgrades/references/runbook-template.md`: it restores PDBs as numbered, verified steps.

  The ungated ones are findings 5, 8 and 13.

## Not examined

- **References and assets:** the 37 reference `.md` files and 28 of the 34 YAML assets were swept with grep (patterns above), not read end to end. Read in full or by flagged region:
  - `gke-upgrades/references/troubleshooting.md` and `runbook-template.md`
  - `gke-basics/references/cli-reference.md` and `mcp-usage.md`
  - `gke-compute-classes/references/compute-class-provisioning-methods.md`
  - `gke-cluster-autoscaler/references/ca-provisioning.md`
  - six manifests: `gateway.yaml`, `gateway-https.yaml`, `ingress.yaml`, `computeclass-rbac-editor.yaml`, `workload-identity-pod.yaml`, `gke-app-onboarding/assets/deployment.yaml`

  Instruction-shaped text in the unread parts, phrased outside the sweep patterns, would be missed.
- **Remote and generated content:** what the two unpinned install URLs currently serve, and what `gcloud container ai profiles manifests create` generates, were not fetched.
- **Harness side:** only `internal/harness/report.go:40` was checked, for the `ya29.` needle. Whether the dream job's transcript reads apply the same redaction was not checked. The harness skills the router prefers (`gke-basics`, `api-auth-gcp`, `cloud-run-deploy`, …) were not read.
- **Tools and live checks:** `govulncheck` and `npm audit` were not run (reasons under Checked). Nothing was executed against a cluster; every finding comes from reading the text.
- **Scratch:** concatenations used for reading are in `.agents/scratch/security-catalog-gke/` (gitignored).

## Recommendations for catalog/MANIFEST.md § Notes

| entry | recommendation | note text |
|---|---|---|
| gke-batch-hpc | note | The Kueue (`releases/latest`) and MPI Operator (`master` branch) installs apply mutable remote manifests cluster-wide: pin a release tag, download to `.agents/scratch/`, verify, review the RBAC and webhooks, ask the operator, then apply the file (router step 4). |
| gke-cluster-creation | note | Templates 4 (Standard JSON) and 5 (A3) drop the golden-path hardening and run nodes as the default Compute Engine SA with `cloud-platform` scope: add a least-privilege node SA, Workload Identity (`GKE_METADATA`), private nodes and authorized networks, or tell the operator the template is unhardened before asking to create. |
| gke-workload-identity | note | Step 5(c) prints a live access token: test with `…/service-accounts/default/email`, or discard the body (`curl -o /dev/null -w '%{http_code}'`); never print the token endpoint's response. |
| gke-workload-troubleshooting | note | Step 5.3's branch, commit, push and PR happen only when the operator asks; otherwise present the patch and stop. Context is asked for, never inferred "for autonomous execution" (router step 4). |
| gke-upgrades | note | references/troubleshooting.md: deleting a ValidatingWebhookConfiguration may switch off a security policy, and Option A relaxes a PDB with no backup. Name what the webhook enforces, ask the operator, prefer narrowing its scope, back up first, and restore as a numbered, verified step (the entry's own runbook rule). |
| gke-platform-security | note | The cross-service `gcloud projects add-iam-policy-binding … --quiet` examples grant project-wide: bind on the bucket or subscription where the resource allows it, and confirm every IAM change with the operator. |
| gke-backup-dr | note | `--include-secrets` (with `--all-namespaces`) copies every Kubernetes Secret into backups readable by whoever can restore them: ask before including secrets or choosing a cross-project backup target, and use CMEK. |
| gke-ai-troubleshooting-jobset-interruption | note | "Autonomous" covers the read-only queries only: Resolution 2 (cordon/drain the node, delete the GCE VM) is proposed to the operator and run only on their confirmation. |
| gke-golden-path | note | The guardrail "don't ask users to paste project IDs" is overridden: ask for and record project, region and cluster first (router step 4); `gcloud config get-value project` is a hint to confirm, not the answer. |
| gke-compute-classes | note | CODE-FIRST: fetch or clone the autoscaler source only with the operator's go-ahead, into `.agents/scratch/`, read-only. The reference's "Privileged DaemonSets (Recommended)" yields to the entry's POD-PRIVILEGE rule. `assets/log-autoscaler-events.sh` polls until killed: run it under `timeout`. |
| gke-app-onboarding | note | Pin the buildpacks builder (`gcr.io/buildpacks/builder@sha256:…`) and the Dockerfile base image by digest instead of `:latest` or a tag (router step 4). |
| gke-inference | note | Keep the Hugging Face token in Secret Manager and mount it through the Secret Manager add-on (gke-workload-security § 6), never via `kubectl create secret --from-literal` or in a manifest (router step 4). |
| gke-cluster-autoscaler | note | "Delete instances manually in GCE" and `min-nodes` changes are proposed, never run unconfirmed. `find-scale-down-blockers.sh` writes a temp file into the current directory: run it from `.agents/scratch/`. `log-autoscaler-events.sh` polls until killed: run it under `timeout`. |
| gke-service-networking | note | `assets/gateway.yaml` and `assets/ingress.yaml` are internet-facing and HTTP-only: confirm the service should be public, default to `gateway-https.yaml` with an HTTP→HTTPS redirect, and apply per the project's deploy guide (router step 4). |
| gke-basics | note | Superseded by the harness `gke-basics` skill. If its references are read anyway: pin `helm install` with `--version` and a reviewed chart source, and treat enabling `allowExternalTraffic` or adding authorized networks as a control-plane exposure change the operator approves. |

Counts: 15 note, 0 exclude, 20 with no finding that needs a note.

- 15 are clean: gke-ai-troubleshooting-handle-disruption-gpu-tpu, gke-ai-troubleshooting-tpu-dynamic-slices-monitoring, gke-alert-configuration, gke-cost-analysis, gke-cost-optimization, gke-custom-golden-image-discovery, gke-manifest-generation, gke-multitenancy, gke-networking, gke-observability, gke-productionize, gke-reliability, gke-storage, gke-workload-scaling, gke-workload-security.
- 5 carry only the router-overridden C1: gke-ai-troubleshooting-tpu-metrics-monitoring, gke-ai-troubleshooting-tpu-vbar-oom, gke-node-notready, gke-storage-troubleshooting, gke-workload-scaling-troubleshooting.

## Outcome (main agent, 2026-10-01)

- The five mediums (`gke-batch-hpc` installing from `releases/latest` and a `master` branch; `gke-cluster-creation` templates 4 and 5 on the default service account; `gke-workload-identity` printing a live token; `gke-workload-troubleshooting` opening a PR unasked; `gke-upgrades` deleting an admission webhook): noted per entry under `catalog/MANIFEST.md` § Notes; the router's step 5 names `kubectl apply -f https://…` from `latest` or a branch as unverified remote code and step 4 requires the operator's go for any cluster change.
- The 15 `note` rows: written under § Notes. No entry excluded; 20 carry no line.
