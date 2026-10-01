---
owner: security
date: 2026-10-01
---
# Security review — worktree-p6   (2026-10-01)

Range reviewed: working tree against HEAD 8be6c06 (Phase 6 of the audit-fix plan, uncommitted), rounds 1 to 3, the same files each time: `skills/terraform/**`, `skills/cloud-run-deploy/SKILL.md`, `skills/docker/**`, `skills/gcloud-setup/SKILL.md`, `skills/golang/references/service.md`, `skills/json-api-style/SKILL.md`, `skills/init-project/assets/agents-tree/docs/deploy/gcp.md`, `internal/harness/skill_links_test.go`; `git diff HEAD -- <files>` and each file whole. Round 3 read all 17 files in that scope end to end (2,024 lines, concatenated in `.agents/scratch/secrev-p6-r3/all.txt`). What a round did not reach is under that round's Not examined.

## Status: COMPLETE

## Verdict: FINDINGS

Round 3: one medium and four lows remain. Every other round-2 fix holds, including the `.dockerignore` allowlist, which was verified by build. The worst problem is in the `/tick` fix itself. The comment under the scheduler's `oidc_token` (`skills/terraform/references/resources-gcp.md:316-317`, echoed at `:229`) tells a public service's handler to call `idtoken.Validate` with the service URL as audience, and to do nothing else. `Validate` accepts any Google-signed ID token whose `aud` matches (`google.golang.org/api@v0.287.1/idtoken/validate.go:160`, with RS256 checked against Google's shared certs at `:185`). It never compares `email` or `iss`. Any service account in any project can mint such a token for any audience, so the check admits anyone with a Google Cloud account. The lows:
- `resources-gcp.md:237` still says the runtime role is DML-only, and `:3` still calls `db_url` only `sensitive`.
- Immutable tags now refuse the `:latest` alias that `cloud-run-deploy/SKILL.md:54` and the `gcp.md` template still offer.
- The new owner/DML recipe in `service.md:316` misses Cloud SQL's `cloudsqlsuperuser` membership for API-created users.
- The allowlist paragraph claims that a secret in any shape stays out. A `key.json` or `*.pem` inside `internal/` or `web/` still gets in, as the build showed.

Round 2 verdict, kept as written:
Round 2: all five round-1 mediums are fixed, and nothing above low remains. The `.dockerignore` probe was re-run with the new list taken from the file. `infra/terraform.tfvars`, `web/.env.local`, `keys/x.key` and `keys/x.pem` no longer enter the context, and neither do the plan file, `.terraform/terraform.tfstate` or root `.env`. The worst thing left is in that same fix. Moving from `*.env*` to `**/.env*` brought a root `app.env` back into the build context; the same build with the round-1 list keeps it out. A `key.json` in the project directory also still gets through, and that is exactly where `gcloud-setup`'s last-resort sample puts one. So do SSH keys without an extension, `web/.npmrc` and the tfvars JSON forms. The rest are lows:
- `idtoken.Validate` silently drops its audience check when the configured audience is empty.
- The migration-role split is asserted ("cannot drop tables") but never wired.
- The reference's live `allUsers` block would make "the standard app" public.
- The public-mechanism rule is missing from the deploy-time skill.
- A public service leaves the scheduler's `/tick` open to anyone.
- Several round-1 sub-points the Outcome marks fixed are not.

Round 1 verdict, kept as written:
No secret value, no hand-rolled crypto and no unsafe IAM default appears in any sample. The runtime SA is dedicated, `run.invoker` is granted per service, `public` defaults to false, secrets come from Secret Manager, and no SA keys are created. The worst problem is the `.dockerignore` list that every Go Dockerfile here copies: its patterns match only at the build-context root. So the secrets the terraform and frontend layouts keep below the root get through, together with any key outside the root: `infra/terraform.tfvars` (the value fed to `secret_data_wo`), the plan file (which the skill says "holds every value in clear"), `web/.env*`, and `keys/*.key|*.pem`. The Go stage's `COPY . .` then puts them into the build context and the build-stage layer, and `setup.md:24` claims this cannot happen. I checked this with a build, not by reading. Next come three gates that pass when they should not. `trivy config .` exits 0 on findings. `golang:<the go.mod go line>` pins the patch level the module states as its minimum, so this repo's own `go 1.24.4` would build on a toolchain with known stdlib advisories. The Terraform samples read secrets at `latest`, which contradicts the pin rule in `cloud-run-deploy`. The skill that should cover IAP JWT validation does not exist.

## Findings (round 1)

- **[medium / confidence: high]** `skills/docker/references/setup.md:10-24`: the patterns `.agents/scratch`, `*.env*`, `*.pem`, `*.key`, `*.tfvars` and `*.tfstate` match only at the context root. Docker matches `.dockerignore` with Go `filepath.Match` relative to the root, and only a `**/` prefix reaches subdirectories. That is unlike `.gitignore`. Verified: I copied the exact list (lines 11-21) into `.agents/scratch/secprobe-p6-dockerignore/` with dummy files and ran `FROM scratch` + `COPY . /ctx/` through `docker buildx build --output type=local`. Root `.env` was excluded. `infra/terraform.tfvars`, `infra/.agents/scratch/tfplan`, `web/.env.local`, `keys/sa.key` and `keys/tls.pem` all reached the output. These are exactly the paths `skills/terraform/SKILL.md:27,31,74` (run from `infra/`) and the frontend's `web/` produce. The distroless final stage copies only `/app`, which keeps the runtime image clean. The secrets still sit in the build context sent to the daemon or Colima VM and in the Go stage's cached layer. They also travel with any `--target build` image or `--cache-to … mode=max` export. Line 24 promises otherwise. Fix: an allowlist (`*`, then `!go.mod`, `!go.sum`, `!cmd/`, `!internal/`, `!web/`, then re-exclude `**/.env*`). Otherwise prefix every entry with `**/` and add `**/*.tfvars.json`, `**/*.tfstate*`, `**/.terraform`, `**/tfplan`. Either way, correct the claim on line 24.
- **[medium / confidence: high]** `skills/terraform/SKILL.md:70,76`, `skills/terraform/references/checks.md:3,11`: both say to run each check bare and read its exit status. But the command given is `trivy config .`, and trivy exits 0 whatever it finds unless `--exit-code 1` is passed. The table's exit-code column mentions the flag, but neither command does. The one misconfiguration scanner in the procedure ("public buckets, missing deletion protection, plaintext secrets") therefore always reads as clean. Fix: `trivy config --exit-code 1 .` in both places.
- **[medium / confidence: high]** `skills/docker/SKILL.md:60,75,115`: `FROM golang:<the go.mod go line>` takes the module's `go` line as the image tag. That line is a minimum language version, and since Go 1.21 it is usually three-part. This repo's own `go.mod` says `go 1.24.4`. Following the rule here gives `FROM golang:1.24.4`, a frozen patch that `AGENTS.md` says carries the `os`/`os/exec` stdlib advisories fixed in 1.24.6. The official image sets `GOTOOLCHAIN=local`, so nothing upgrades it. Every stdlib CVE fixed after the module's floor ends up compiled into the static binary. Fix: `golang:<major.minor>`, taken from the go line (`go 1.24.4` → `golang:1.24`, which tracks the latest patch). Never use a patch tag, and run `govulncheck -mode=binary` on the built binary as part of the smoke test.
- **[medium / confidence: high]** `skills/terraform/assets/example/main.tf:64-71,81-84` and `skills/terraform/references/resources-gcp.md:155-158`: the Cloud Run env var reads its secret at `version = "latest"`. This contradicts `skills/cloud-run-deploy/SKILL.md:101` ("Pin the version (`:3`), not `:latest`, for env vars") and the template's "pinned versions" row (`gcp.md:39`). With `latest`, a principal holding only `secretmanager.versions.add` decides what the next instance of a running revision reads. That needs no deploy and no `run.services.update`. Instances of one revision can also disagree, and a rollback does not roll the secret back. The example has a second defect: the service does not depend on `google_secret_manager_secret_version.app`, so the first apply can start a revision before any version exists. Fix: `version = google_secret_manager_secret_version.app.version`, which pins and orders in one reference. Do the same in resources-gcp.md.
- **[medium / confidence: high]** `skills/cloud-run-deploy/SKILL.md:96,100`: "The app still validates `x-goog-iap-jwt-assertion` (`api-auth-gcp` skill)" and "Connection code: `data-gcp` skill" point at skills that do not exist (`ls skills/`). Phase 7 is meant to write them. An agent with nothing to read will compose the ES256 signature, `aud`, `iss` and expiry checks itself. That is hand-rolled JWT verification, which the global context bans as high severity. This was also raised by the p6 code, docs and prompts reviews; it is recorded here for the security angle. Fix: name the vetted library in the cell, for example `google.golang.org/api/idtoken.Validate(ctx, token, audience)` with the audience Google documents for IAP on Cloud Run. The stdlib alternative, composing JWT/ES256 verification over `crypto/ecdsa`, is rejected because it is the banned hand-rolling. Otherwise ship Phase 6 and 7 together. The same row should also say where `roles/iap.httpsResourceAccessor` is granted: on the service's IAP resource, not the project, where it opens every IAP-protected resource.
- **[low / confidence: high]** `skills/terraform/references/resources-gcp.md:174-181`: the Direct VPC egress sample hardcodes `network = "default"`, `subnetwork = "default"`. The auto-mode default VPC carries `default-allow-internal` (every port from 10.128.0.0/9) and SSH/RDP open to 0.0.0.0/0, so the service joins a flat network with every VM in the project. Many organisations also delete that network (`compute.skipDefaultNetworkCreation`). Phase 6 deleted `security-compliance.md`, whose "DON'T: Use Default VPC" was the only caution, and the new sample does the opposite. The same sample also pairs public-IP Cloud SQL (`:113-115`) with `vpc_access`, against `SKILL.md:63` ("no `vpc_access` at all without a private target"). Fix: a dedicated `google_compute_network` (`auto_create_subnetworks = false`) and a `/26`+ `google_compute_subnetwork` passed in by variable. Set `network_tags` on the service for egress firewall rules, add `google_service_networking_connection` for private Cloud SQL, and comment out the block in the public-IP sample.
- **[low / confidence: medium]** `skills/cloud-run-deploy/SKILL.md:95,139-142`, `skills/terraform/SKILL.md:53`, `assets/example/main.tf:54-55`: the two skills make a service public by different mechanisms. Terraform uses `invoker_iam_disabled = true`; gcloud uses `--allow-unauthenticated`, which writes an `allUsers` binding. The IAM-gated path is "omit the flag", which leaves an existing `allUsers` binding or a disabled invoker check in place. Verify calls the unauthenticated 403 "proves nothing", yet for `Public: no` that 403 is the one probe that proves the gate. A service recorded `Public: no` can stay public across deploys with every Verify step green. Fix: when the Public row says no, pass `--no-allow-unauthenticated` explicitly (and `--invoker-iam-check`, the gcloud spelling `--[no-]invoker-iam-check`, to confirm against `gcloud run deploy --help`). Add to Verify a bare `curl -s -o /dev/null -w '%{http_code}'` that must print 403. Have `gcp.md` record which public mechanism the project uses and who holds `roles/run.invoker`.
- **[low / confidence: high]** `skills/terraform/references/resources-gcp.md:105-117`: the Cloud SQL sample enables public IP but sets no `ssl_mode`, so it inherits the API default that accepts unencrypted connections. Only the absence of `authorized_networks` keeps clients on the connector path, and nothing pins that if a later edit adds one. The sample also has no `backup_configuration` (Terraform's default is off), so a destructive incident or a bad migration cannot be recovered. The skill's own step 5 (`trivy config`) has Cloud SQL SSL and backup checks that would flag the sample (confidence medium on the exact check IDs). That teaches the agent to treat trivy output as noise. Fix: `ssl_mode = "ENCRYPTED_ONLY"` and `backup_configuration { enabled = true  point_in_time_recovery_enabled = true }`.
- **[low / confidence: high]** `skills/docker/references/setup.md:35-39,75-76`, `skills/docker/SKILL.md:50`, `skills/cloud-run-deploy/SKILL.md:58` (pre-existing lines, in the blast radius): compose publishes Postgres as `"${DB_PORT:-5432}:5432"` with superuser `app`/`app`, and the app and smoke tests use `-p 8080:8080`, all bound to every interface. On a Linux box Docker inserts its port rules ahead of UFW and firewalld, so the database is reachable from the LAN despite a host firewall. A Postgres superuser can also run `COPY … FROM PROGRAM` inside the container. Fix: bind to loopback with `"127.0.0.1:${DB_PORT:-5432}:5432"` and `-p 127.0.0.1:8080:8080`.
- **[low / confidence: medium]** `skills/docker/SKILL.md:60,67,75,87`, `references/setup.md:34`: base images are pinned by tag only (`gcr.io/distroless/static-debian12`, `node:24-alpine`, compose `postgres:16`). Rule 1 calls a tag a pin and never mentions a digest. Rule 5 forbids token build args but offers no sanctioned path for a build that needs one, such as private modules. Fix: state the trade-off (a tag gets patch rebuilds; a digest is reproducible but needs deliberate bumps). Pin the runtime stage by `@sha256:` with the tag as a comment, and name BuildKit `RUN --mount=type=secret` as the only build-time secret path.
- **[low / confidence: medium]** Tool installs without integrity checks:
  - `skills/terraform/SKILL.md:72-77` and `checks.md:5-11` require `terraform`, `tflint` and `trivy` but name no install source, version or verification. Trivy's own release channel was compromised in 2026 (incident details from memory, not re-verified).
  - `skills/golang/references/service.md:376` pipes `curl … golangci-lint.run/install.sh | sh`. The script fetched from a website is itself unverified.
  - `skills/gcloud-setup/SKILL.md:49` downloads an unversioned tarball and runs `install.sh` with no checksum.

  Fix: pinned versions verified against the publisher's signed checksums. That means HashiCorp `SHA256SUMS` plus its signature or the HashiCorp apt repo, release `checksums.txt` plus cosign or `gh attestation verify` for trivy, tflint and golangci-lint, and the signed apt repo (or Google's published SHA-256) for gcloud.
- **[low / confidence: high]** `skills/terraform/references/checks.md:29-36`, `resources-gcp.md:3`, `SKILL.md:27`: secrets at rest around Terraform.
  - `infra/.gitignore` misses `*.tfvars.json` and `*.auto.tfvars.json`, which Terraform loads automatically. The installer's secret-shaped list (`context/AGENTS.md:47`) misses them too.
  - `var.db_url` is only `sensitive`, not `ephemeral` like the example's `app_secret`, so its value is written into the saved plan file.
  - The main path keeps the secret in plaintext in `infra/terraform.tfvars`, the path finding 1 shows leaking into image builds.

  Fix: add the `.json` forms, mark `db_url` `ephemeral = true`, and make the documented alternative (`gcloud secrets versions add … --data-file=-`, or `TF_VAR_*` read with `read -s`) the default.
- **[low / confidence: medium]** `skills/terraform/SKILL.md:36-40`: the state bucket name `PROJECT-tfstate` is predictable in a global namespace. Nothing says that a create failing with "already exists" is a stop. A pre-registered bucket in someone else's project, made writable to all, would receive this project's state. Fix: treat a 409 on create as a stop, and before `init` confirm the bucket's `projectNumber` (`gcloud storage buckets describe`) matches the project, or add a random suffix recorded in `gcp.md`.
- **[low / confidence: medium]** `skills/terraform/references/resources-gcp.md:215-227`, `skills/golang/references/service.md:288,302-305,372`: the migration job runs as the runtime SA with the runtime DSN. The serving identity therefore owns the schema (DDL), so a SQL injection or a compromised instance can drop tables. The job sample also mounts neither `/cloudsql` nor the DSN secret, so it cannot connect as written. Fix: a separate `${app}-migrate` SA and DB role that owns the schema, with the runtime role DML-only.
- **[low / confidence: medium]** `skills/golang/references/service.md:155-280,80,250-273`: the logging contract names `httpRequest` as a liftable key. Its natural fill, `requestUrl`, carries the query string. The section neither points back to the never-log rule at `skills/golang/SKILL.md:34` nor names headers (`Authorization`, `Cookie`, `x-goog-iap-jwt-assertion`) and bodies. At the parsing boundary, `readJSON` echoes the decoder's error text to the client, including Go type names such as `greeting.name`. Client-supplied `traceparent` and `X-Cloud-Trace-Context` ids are length-checked but not hex-checked, so a caller can attach its log lines to any trace. The slog JSON handler escapes values, so there is no log forging. Fix: a one-line pointer to SKILL.md:34 that names headers and bodies; `requestUrl` without the query; a fixed client message plus `BadRequest.fieldViolations`; a 32/16 lowercase-hex check.
- **[low / confidence: high]** `skills/cloud-run-deploy/SKILL.md:29` (pre-existing): lists `.env` as a project-ID cross-check, against `skills/gcloud-setup/SKILL.md:128`, which says the ID is never discovered from `.env`. Reading `.env` puts secret values into the model transcript. Fix: drop `.env` from the sentence.
- **[low / confidence: high]** `internal/harness/skill_links_test.go:33,45,67-71`: no security impact. The test is read-only. `filepath.WalkDir` does not descend directory symlinks, but `os.ReadFile` follows a symlinked `*.md`. `os.Stat(filepath.Join(base, target))` resolves `..` and symlinks, so a link can resolve outside the skill and the repository (existence only). A link that escapes the skill passes in the checkout and dangles in an installed `~/.claude/skills/<name>/`. Fix (stdlib): require `filepath.IsLocal` of the target relative to the skill root (or `skills/` for cross-skill paths), or resolve through `os.Root`. The test also cannot see bare skill names (finding 5).

## Checked (round 1)

- **Secrets:** no secret value in any HCL, Dockerfile, env sample or log call. `APP_SECRET` and `DATABASE_URL` come from Secret Manager refs. The example uses a write-only `secret_data_wo` fed by a `sensitive` + `ephemeral` variable. `--set-env-vars` is forbidden for secrets (`cloud-run-deploy:102`, template `:47`), and the template stores names, never values (`:39`). The JSON conventions keep stack traces and secrets out of `message`, and a 5xx carries no cause (`json-api-conventions.md:114,117`). Gaps: findings 1, 4, 12, 16.
- **Injection:** the terraform example has no `templatefile`, `local-exec`, provisioner or `external` data source; variables reach only names and the image path. Shell samples interpolate operator-recorded defaults, double-quoted. SQL guidance uses `$1` parameters, no query builder. JSON input is bounded at 1 MiB with unknown fields refused.
- **Authz/authn:** the runtime SA is dedicated per service and never the default compute account. `run.invoker` is per service, never at project level. Scheduler uses OIDC with its own SA; the deployer holds `iam.serviceAccountUser` on the SA; the IAP service agent holds `run.invoker`. `public` defaults to false. SA keys are never created (terraform) and are a last resort on the operator's say-so (gcloud-setup, BUG-019 fix intact). No `roles/editor` anywhere. Gaps: findings 5, 7.
- **Crypto:** nothing composed in scope. `crypto/rand.Text` is recommended for tokens (`service.md:394`). The IAP JWT gap is finding 5.
- **Dependencies:** no `go.mod` change in range; the test is stdlib only. The dependencies the skills recommend to downstream projects (pgx/v5, cloudsqlconn, golang-migrate, OTel, profiler) are each justified in the text. `govulncheck ./...`: exit 0, "Your code is affected by 0 vulnerabilities", plus 1 vulnerability in a required module that is not called (output in `.agents/scratch/govulncheck-p6.txt`; `-show verbose` not run, so not identified by ID). `npm audit`: not applicable, `git ls-files` tracks no `package.json`.
- **GCP surface:** no `--allow-unauthenticated` by default. The state bucket has public-access prevention, uniform bucket-level access and versioning. Ingress is `all` with IAM for the gated default, and `internal` is documented. `deletion_protection` stays at its default. Gaps: findings 6, 7, 8, 13.

## Not examined (round 1)

- `skills/json-api-style/references/json-api-conventions.md` lines 1-84 (grepped for error, pagination and secret terms only; lines 85-149 read); `skills/golang/SKILL.md` (grepped for its logging rules only).
- Phase 6 files outside the caller's scope: `skills/frontend/**` (including how its node stage and `go:embed` meet this Dockerfile), `skills/golang/references/testing.md`, `skills/gke-basics`, `skills/bash`, `skills/mermaid-syntax`, `skills/review*`, `skills/VENDORED.md`, the plan and skill-matrix edits.
- The deleted terraform references: `security-compliance.md` read by headings only; `ci-cd-workflows.md`, `code-patterns.md`, `module-patterns.md`, `quick-reference.md` and `testing-frameworks.md` not read.
- No `terraform validate`, `tflint` or `trivy` run against `assets/example/` (none installed). No `gcloud` installed, so the flag claims (`--[no-]invoker-iam-check`, `--iap`, pinned `--set-secrets`) and the IAP role and audience strings were not checked against `--help`.
- The 2026 trivy release-channel compromise and trivy's Cloud SQL check IDs are cited from memory, not re-verified.
- `context/dev-setup.md`, the root-step gate the docker skill's line 109 points at: unchanged in this range and not reviewed.

## Outcome
- [medium] `.dockerignore` patterns match only at the root — fixed at `skills/docker/references/setup.md` (`**/` prefixes, `**/tfplan` added, the reason stated with the probe's date)
- [medium] `trivy config .` without `--exit-code 1` — fixed at `skills/terraform/SKILL.md:76` and `references/checks.md:11`
- [medium] `golang:<the go.mod go line>` pins a patch — fixed at `skills/docker/SKILL.md` (rule 1 and both samples: `golang:<major.minor>`; a digest pin named for reproducible builds)
- [medium] secrets read at `version = "latest"`, version race — fixed (pinned to the managed version; version resource in `depends_on`)
- [medium] IAP validation pointer to a missing skill — fixed in the cell (`idtoken.Validate` with the IAP audience plus an explicit `iss` check; never hand-rolled); the skill itself is Phase 7 (D38)
- [low] Direct VPC egress on the auto-mode `default` network — fixed at `references/resources-gcp.md` (own VPC + subnet blocks; the comment says why)
- [low] public by two mechanisms — fixed at `skills/terraform/SKILL.md` item 6 (Terraform owns the setting when the service is in HCL; `cloud-run-deploy` passes neither flag)
- [low] Cloud SQL without `ssl_mode` — fixed (`ENCRYPTED_ONLY`, no `authorized_networks`, commented)
- [low] compose and smoke ports on every interface — fixed (`127.0.0.1:` in `setup.md`, `docker/SKILL.md:50`, `cloud-run-deploy/SKILL.md:58`)
- [low] tag-only base image pins — fixed (rule 1 names the digest pin)
- [low] tool installs without integrity checks — left: operator — the gcloud apt/tarball lines are Google's documented install path; a checksum step would add a step the operator must take (the simpler option, recommended: keep Google's path)
- [low] secrets at rest around Terraform — fixed at `skills/terraform/SKILL.md` (tfvars never hold a secret value; `TF_VAR_app_secret` from Secret Manager or a prompt)
- [low] predictable state bucket name — fixed (an "already exists" create is a stop; suffixed name)
- [low] migration job as the runtime SA — fixed (own `migrate` SA with DDL; runtime role DML only, in `resources-gcp.md` and `service.md` item 3)
- [low] `httpRequest.requestUrl` carries the query string — fixed at `service.md` (never log the query, `Authorization` or `Cookie`; pointer to the SKILL.md rule)
- [low] `.env` as a project-ID cross-check — fixed at `skills/cloud-run-deploy/SKILL.md:29` (`deploy/gcp.md` instead)
- [low] link test follows symlinked `*.md` — left: no security impact (read-only test over a repository we control), as the finding itself says

## Round 2

### Findings (round 2)

- **[low / confidence: high]** `skills/docker/references/setup.md:17-28`, with `skills/gcloud-setup/SKILL.md:146,149`: the fix holds for every path round 1 measured, but it introduces a regression and leaves gaps.
  - Probe: `.agents/scratch/secprobe-p6-dockerignore-r2/`. The list was extracted from the file by script. The build was `FROM scratch` + `COPY . /ctx/` through `docker buildx build --output type=local`, exit 0. The same context with the round-1 list is in `ctx-r1list/`, for comparison.
  - Out with the new list: `infra/terraform.tfvars`, `web/.env.local`, `keys/x.key`, `keys/x.pem`, `infra/.agents/scratch/tfplan`, `infra/.terraform/terraform.tfstate`, root `.env`, `.agents/scratch/*` and `.agents/memory/*`.
  - Still in with the new list:
    - a root `app.env`, a **regression**: it is OUT with the round-1 `*.env*`;
    - `config/prod.env`, `creds/sa-key.json` and `keys/id_ed25519`;
    - `web/.npmrc`;
    - `infra/terraform.tfvars.json` and `infra/prod.auto.tfvars.json`;
    - `infra/terraform.tfstate.backup`.
  - Why these matter:
    - `gcloud-setup`'s last-resort sample uses `--key-file=key.json`, a file in the working directory. A root `key.json` is credential JSON, which is on the harness's own secret-shaped list, and it lands in the Go stage's `COPY . .` layer and in any cache export of it.
    - Docker rule 5 names `env_file` as the local secret path, and a `*.env` name gets in.
    - The node stage's `COPY web/ ./` copies `.npmrc`.
    - Terraform auto-loads both tfvars JSON forms.
    - Local state and its `.backup` exist while the example's backend block is commented out.
  - Line 25 ("a stray `COPY . .` cannot bake them into a layer") still claims more than the list delivers.
  - Fix: an allowlist, which is closed by default and shorter than any denylist: `*`, `!go.mod`, `!go.sum`, `!cmd/`, `!internal/`, `!web/`, then `web/node_modules`, `web/dist`, `web/.output`, `**/.env*`, `**/.npmrc`. If the denylist stays, add `**/*.env`, `**/*.tfvars.json`, `**/*.tfstate.*` and `**/.npmrc`, and reword line 25 to say what it covers. Either way, the `gcloud-setup` sample's key path should sit outside the repository (for example `~/.config/gcloud/<sa>.json`).
- **[low / confidence: high]** `skills/cloud-run-deploy/SKILL.md:96`: `idtoken.Validate` skips the audience check when the audience is the empty string.
  - Evidence, in the module cache: `google.golang.org/api@v0.287.1/idtoken/validate.go:109` says "If audience is not empty the audience claim of the Token is validated", and `:160` does it with `if audience != "" && payload.Audience != audience`.
  - The risk: an app reads the IAP audience from configuration and gets `""` from an unset env var. It then accepts any genuine IAP assertion minted for any other IAP resource, and the cell's `iss` check passes for those too. The app-side check is the layer that counts when IAP is off or bypassed.
  - The cell's own `iss` claim is right: `Payload.Issuer` is parsed and never compared.
  - Fix, one clause: "with a non-empty audience — `Validate` skips `aud` when it is empty, so startup fails on an empty configured audience".
  - Round 1's other half of this finding is still open: the cell does not say that `roles/iap.httpsResourceAccessor` is granted on the service, never the project. The template does now record who holds it.
- **[low / confidence: high]** `skills/golang/references/service.md:294,309-311,378`, `skills/terraform/references/resources-gcp.md:233-259`: "a compromised instance then cannot drop tables" is stated as a guarantee, but nothing makes it hold.
  - With the password DSN the sample shows (`user=app … password=…`), the Postgres role comes from the DSN, not from the service account. `roles/cloudsql.client` on the `migrate` SA only lets it open a connection.
  - Neither file creates a second role, gives the job a DSN secret and accessor binding of its own, or grants the runtime role DML only.
  - The job block mounts neither `/cloudsql` nor any secret.
  - The easiest completion is the job reading the runtime `DATABASE_URL`. Migrations then run as `app`, `app` must own the schema, and the runtime role has DDL.
  - Fix: say the wiring in one sentence and put the job's volume and secret in the HCL. One way is two roles: an owner role runs migrations, and `app` gets DML through `ALTER DEFAULT PRIVILEGES FOR ROLE <owner> … GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app`. Each role's DSN is a secret readable only by its own SA. The other way is IAM database authentication with each SA as its own database user (`cloudsqlconn.WithIAMAuthN`, already in the text).
- **[low / confidence: medium]** `skills/terraform/references/resources-gcp.md:211-226`: in § Public or IAM-gated, the live resource is the `allUsers` invoker binding. Its alternatives (`invoker_iam_disabled`, and the IAM-gated scheduler member) are comments. Building "the standard app" from the reference's blocks in order, which is where `SKILL.md` § The standard app points, produces a public service. That goes against the example's `public = false` and against `SKILL.md` item 6, which names only `invoker_iam_disabled` for public. Fix: comment out the `allUsers` block as its sibling is, or make the IAM-gated member the live block.
- **[low / confidence: medium]** `skills/cloud-run-deploy/SKILL.md:95,143`: round-1 finding 7 was fixed on the Terraform side only.
  - The deploy-time skill does not carry the rule "when the service is in HCL, neither flag". That rule lives only at `skills/terraform/SKILL.md:53`.
  - Row 95 says "Public services only", so an agent adds `--allow-unauthenticated`, an `allUsers` binding, to a Terraform-managed public service. Terraform's `_member` resources never remove that binding, so a later `public = false` leaves the service public.
  - For a gcloud-only service, omitting the flag keeps an existing `allUsers` binding.
  - Verify still says a bare curl's 403 "proves nothing", yet for `Public: no` it is the one probe that proves the gate is on.
  - Fix: one clause in row 95 pointing at terraform item 6. Reword line 143 to "a bare curl must 403 (the gate's proof); then the authenticated probe". This adds no step for the operator.
- **[low / confidence: high]** `skills/terraform/references/resources-gcp.md:274-288`, with `:211-214` and `assets/example/main.tf:55`: the scheduler sends an OIDC token, but Cloud Run checks it only while invoker IAM is enforced. With `invoker_iam_disabled = true` or an `allUsers` binding, which are the reference's own public forms, `/tick` is callable by anyone, and nothing says the app must then check the token itself. Fix: one line under the scheduler block. On a public service, either keep scheduled endpoints on a separate IAM-gated service, or have the handler verify the token with `idtoken.Validate` against the service URL and compare its `email` claim with the scheduler SA.
- **[low / confidence: medium]** `skills/terraform/references/resources-gcp.md:55-68`, `assets/example/main.tf:15-28`, `skills/cloud-run-deploy/SKILL.md:54,69-71`: the text calls the SHA tag immutable ("rollbacks and audits need immutable tags"), but Artifact Registry Docker tags are mutable unless the repository enforces it, and neither the HCL nor the `gcloud artifacts repositories create` line does. A holder of `artifactregistry.writer` can re-point `:<sha>`. Running revisions are pinned by digest, but a redeploy of that SHA pulls the replacement. Fix, as an operator choice: `docker_config { immutable_tags = true }` plus the create command's immutable-tags flag (confirm the spelling with `--help`), which rules out a moving `latest` alias. The other option is to drop "immutable" from the text.
- **[low / confidence: high]** Round-1 sub-points the Outcome marks fixed that are not:
  - `var.db_url` is still only `sensitive` (`resources-gcp.md:3,85`). It is not `ephemeral` like the example's `app_secret`, so its value is saved in the plan file, which `SKILL.md:31` says holds secrets in clear.
  - The `infra/.gitignore` in `checks.md:29-36` still lacks `*.tfvars.json` and `*.auto.tfvars.json`.
  - Cloud SQL still has no `backup_configuration` (`resources-gcp.md:105-119`), so a destructive migration or injection cannot be undone.
  - `service.md:86` still returns the decoder's `err.Error()` to the client.
  - `service.md:256-277` still checks trace ids by length only, so a client can choose any 32-character string, `/` included, for `projects/<p>/traces/<id>`.
  - Docker rule 5 (`docker/SKILL.md:91-92`) still names no sanctioned build-time secret path (`RUN --mount=type=secret`).
  - `service.md:163` names `Authorization` and `Cookie` but not `x-goog-iap-jwt-assertion`, which the IAP row has the app read. The assertion is a signed identity assertion carrying the user's email.

### Fix verification (round 2)

- `.dockerignore`: fixed for the measured paths, verified by build. The regression and gaps are finding 1.
- `trivy config --exit-code 1 .`: fixed at `skills/terraform/SKILL.md:76` and `references/checks.md:11`.
- `golang:<major.minor>`: fixed at `skills/docker/SKILL.md:68-69,86,128`, and the digest clause is present. Informational: rule 1 says "never `:latest`" yet names `gcr.io/distroless/static-debian12` with no tag, which is an implicit `:latest`. Distroless publishes no version tags, so its only pin is the digest.
- Secret versions pinned: fixed at `main.tf:69,81-85` and `resources-gcp.md:175`. The `.version` reference also orders the resources.
- IAP cell: fixed. It names `idtoken.Validate` and forbids a hand-rolled check, and the `iss` statement matches the library source. The empty-audience gap and the accessor scope are finding 2.
- Own VPC and subnet: fixed (`resources-gcp.md:140-152,193-199`). The VPC is custom-mode with implied deny-ingress, so the `vpc_access` block left live beside the public-IP SQL sample no longer has a security impact.
- `ssl_mode = "ENCRYPTED_ONLY"` with no `authorized_networks`: fixed (`resources-gcp.md:115-116`).
- Loopback ports: fixed (`setup.md:40,80`, `docker/SKILL.md:50`, `cloud-run-deploy/SKILL.md:58`).
- tfvars never hold a secret value: fixed (`terraform/SKILL.md:27`). The rest is in the last finding.
- State bucket "already exists" is a stop: fixed (`terraform/SKILL.md:41`). From memory, not verified: gcloud reports a name held by another project as "not available" (409), not "already exists". Keying the stop on any non-zero exit of the create is safer.
- Migration job under its own SA: the text and the SA are fixed. The enforcement is missing (finding 3).
- Never-log sentence: fixed (`service.md:163`, pointing at `golang/SKILL.md:34`). The IAP header gap is in the last finding.
- `ErrorInfo` in error bodies: leaks nothing. `reason` (`INTERNAL`, `MALFORMED_BODY`, `NAME_REQUIRED`) and `domain` (`greet.example.com`) are compile-time constants. The 500 path sends a fixed "internal error" and logs the cause server-side only (`service.md:61-71`). The only dynamic text reflected to the client is still `readJSON`'s (last finding).
- `.env` cross-check dropped: fixed (`cloud-run-deploy/SKILL.md:29-30`). A docs note with no security impact: the sentence now calls `deploy/gcp.md` "a cross-check, not a source" two lines after making it the source.

### Checked (round 2)

- **Secrets:** no secret value in any changed sample. Write-only versions, pinned version refs, the tfvars rule, and names-not-values in the template all hold. `git diff HEAD` removed lines grepped for security terms: every removed caution has a stronger replacement, except the `*.env*` → `**/.env*` change (finding 1). The deleted `security-compliance.md` was AWS-centric; its GCP-relevant DON'Ts (secrets in variables, default VPC, open firewall) now have counterparts.
- **Injection:** no new input surface. `readJSON` still bounds the body at 1 MiB and refuses unknown fields. SQL uses `$1` parameters. The reflected decoder text and the length-only trace check carry over (last finding).
- **Authz/authn:** findings 2-6. The runtime and migrate SAs and `run.invoker` per service are as round 1 found.
- **Crypto:** nothing composed. The IAP cell now names a vetted library (`google.golang.org/api/idtoken`) and forbids a hand-rolled JWT check. Consuming it is not hand-rolling, so no finding.
- **Dependencies:** `git diff HEAD --stat -- go.mod go.sum` is empty, and `skill_links_test.go` is stdlib only. `govulncheck ./...` on go1.27.1: exit 0, "Your code is affected by 0 vulnerabilities", plus 1 vulnerability in a required module that is not called (`.agents/scratch/govulncheck-p6-r2.txt`). That matches round 1 and the known GO-2026-5024 in `AGENTS.md`, but it was not re-identified with `-show verbose`. `npm audit`: not run, since no `package.json` is in scope; round 1 found none tracked.
- **GCP surface:** own VPC, `ssl_mode`, pinned secret versions and loopback ports verified. Findings 4-7 cover the rest.

### Not examined (round 2)

- `skills/json-api-style/references/json-api-conventions.md`: outside this round's scope list (round 1 read lines 85-149). `skills/golang/SKILL.md`: only line 34 grepped. `context/dev-setup.md`: not reviewed.
- The five other deleted terraform references (`ci-cd-workflows.md`, `code-patterns.md`, `module-patterns.md`, `quick-reference.md`, `testing-frameworks.md`): not read. `security-compliance.md`: only its headings and DO/DON'T lines, via `git show HEAD:`.
- No `terraform validate`, `tflint` or `trivy` run against `assets/example/` (`trivy` absent on this box). `gcloud` is not installed, so the immutable-tags flag, `--[no-]invoker-iam-check`, the IAP audience format `/projects/N/locations/R/services/S` and the issuer string were not checked against `--help` or Google's docs. The `idtoken` behaviour was checked against the module-cache source (v0.287.1).
- The gcloud 409 wording and the Cloud Run behaviour with invoker IAM disabled are stated from documentation memory, not re-verified.
- Phase 6 files outside the caller's scope (as in round 1): `skills/frontend/**`, `skills/golang/references/testing.md`, `skills/gke-basics`, `skills/bash`, `skills/mermaid-syntax`, `skills/review*`, `skills/VENDORED.md`, the plan and the skill-matrix edits.

## Outcome (round 2)
- [low] `.dockerignore` denylist regression and gaps — fixed at `skills/docker/references/setup.md`: the list is now the allowlist the finding proposed (`*`, `!go.mod`, `!go.sum`, `!cmd/`, `!internal/`, `!web/`, then `web/node_modules`, `web/dist`, `web/.output`, `**/.env*`, `**/*.env`, `**/.npmrc`), closed by default, and the paragraph says what it covers and why the denylist was dropped
- [low] `idtoken.Validate` skips `aud` when empty; accessor scope — fixed at `skills/cloud-run-deploy/SKILL.md:96` (non-empty audience, startup fails on an empty one; `roles/iap.httpsResourceAccessor` on the service, never the project)
- [low] DDL/DML claim not wired — fixed: `service.md` item 3 and the job block now say what IS wired (own SA; migrations run as the schema-owning `app` role over the same DSN) and give the role split as a project decision with its SQL; the false guarantee sentence is gone
- [low] live `allUsers` block — fixed at `resources-gcp.md` § Public or IAM-gated (the IAM-gated member is the live form, `allUsers` a comment with the Public: yes condition and the `/tick` consequence)
- [low] public-flag rule missing on the gcloud side; bare curl wording — fixed at `skills/cloud-run-deploy/SKILL.md:95,143`
- [low] `/tick` open on a public service — fixed (comment on the scheduler's `oidc_token`: the handler validates the token on a public service)
- [low] mutable tags — fixed (`docker_config { immutable_tags = true }` in both HCL sites; `--immutable-tags` on the `gcloud artifacts repositories create` line, flag present in gcloud 587 `--help`)
- [low] round-1 sub-points: `db_url` ephemeral (comment on the version resource names `sensitive` AND `ephemeral`); `infra/.gitignore` gains `*.tfvars.json`, `*.auto.tfvars.json`; Cloud SQL `backup_configuration` with PITR; `readJSON` logs the decoder error and returns a fixed message; trace ids checked as hex (`isHex`); docker rule 5 names `RUN --mount=type=secret`; `x-goog-iap-jwt-assertion` in the never-log list; distroless pinned by digest (rule 1 and the Dockerfile) — all fixed
- left: `gcloud-setup`'s last-resort `--key-file=key.json` path (the allowlist keeps it out of the context now; moving the sample path is a gcloud-setup change outside this phase's scope: Phase 9 docs lens)

## Round 3

### Findings (round 3), FAIL first

- **[medium / confidence: high] FAIL, the fix does not hold:** `skills/terraform/references/resources-gcp.md:316-317`, with `:228-229`.
  - The `/tick` remedy is "the handler validates it (idtoken.Validate with this audience) or anyone can call /tick". That authenticates any Google identity, not the scheduler.
  - `idtoken.Validate` checks the signature, the expiry, and `aud` when it is non-empty (`validate.go:160`, v0.287.1). RS256 tokens are checked against Google's shared certs (`:185`). `Issuer` is parsed (`:48`) and never compared, and `email` is not read at all.
  - Any service account in any project can mint a Google-signed ID token for an audience it chooses: the IAM Credentials `generateIdToken` call, or the metadata server's `identity?audience=`. On a public service that audience, the service URL, is public.
  - Result: round 2's open `/tick` stays open behind a check that looks like authentication. The precondition (Public: yes) is not the default, which is why this is medium rather than high.
  - Fix, one clause at both lines: after `Validate`, require `payload.Claims["email"]` to equal the scheduler SA's email and `email_verified` to be true. Google's push-authentication samples check this way (from memory). The alternative is to keep scheduled routes on a separate IAM-gated service. No new dependency.
- **[low / confidence: high] FAIL, restatements the fixes left:**
  - `resources-gcp.md:237`: "it owns the schema (DDL); the runtime SA's database role is DML only" heads the block whose own env comment (`:268-269`) and `service.md:316` now say migrations run as the schema-owning `app` role over the runtime DSN. A comment above a block is read as what the block delivers, so this is round 2's false guarantee, one copy left.
  - `resources-gcp.md:3`: "`var.db_url` is `sensitive`". Since `:88` now says `sensitive` AND `ephemeral`, and no `variable "db_url"` block exists in the skill, these two sentences are all an agent has to write it from.
  - Fix: `:237` says "its own identity, so IAM and logs tell the job from the service; the database role comes from the DSN (service.md § Database item 3)". `:3` says "`sensitive` and `ephemeral`".
- **[low / confidence: medium] Introduced by a fix:** `skills/cloud-run-deploy/SKILL.md:54` ("add :latest only as a convenience alias") and `skills/init-project/assets/agents-tree/docs/deploy/gcp.md:27` ("`latest` only as a convenience alias") conflict with the new `--immutable-tags` (`:71`) and `immutable_tags = true` (`resources-gcp.md:62`, `assets/example/main.tf:21`).
  - In a repository with immutable tags, a tag never moves to another digest, so the second release's `:latest` push is refused. Round 2 named this: "rules out a moving `latest` alias".
  - The failure lands mid-release, and the quickest way out (`--no-immutable-tags`, or dropping the HCL line) removes the control the fix added.
  - From Artifact Registry documentation memory; not tried against a registry (gcloud is absent).
  - Fix: drop the alias in both lines, for example "the SHA tag only: an immutable repository takes no moving alias". The same phrase in the global context (`context/dev-setup.md`, out of scope) needs the same edit.
- **[low / confidence: medium] Introduced by a fix:** `skills/golang/references/service.md:316`: the owner/runtime split recipe says "`ALTER DEFAULT PRIVILEGES FOR ROLE <owner> … GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app` so `app` gets DML only".
  - Cloud SQL makes a user created through its API (`google_sql_user`, `gcloud sql users create`, which is how this stack creates users) a member of `cloudsqlsuperuser`, with `CREATEROLE` and `CREATEDB`. That is from Cloud SQL documentation memory, not re-verified.
  - Default privileges add grants and take nothing away, so "DML only" does not follow.
  - It matters only once a project adopts the split, which the text calls a project decision.
  - Fix, one clause: `app` is created by the owner with SQL `CREATE ROLE app LOGIN`, or as an IAM database user, which starts with no privileges. It is never created through the API, and gets `REVOKE CREATE ON SCHEMA public FROM app`.
- **[low / confidence: high] Introduced by a fix, verified by build:** `skills/docker/references/setup.md:28-30` claims "a secret in any shape — `key.json`, `id_ed25519`, `*.tfvars.json`, `*.tfstate.backup` — stays out without being named".
  - That holds outside the re-admitted directories only. In the probe, `internal/config/sa-key.json` and `web/localhost-key.pem` got in. The latter is the mkcert file Vite's HTTPS dev server uses.
  - Inside `cmd/`, `internal/` and `web/`, only `.env*`, `*.env` and `.npmrc` are kept out.
  - Fix: reword, for example "outside the re-admitted directories; inside them only the shapes below are excluded, so no credential lives in `cmd/`, `internal/` or `web/`". Add no new patterns: growing a denylist inside the allowlist would repeat the round-2 mechanism.

### Fix verification (round 3)

- **`.dockerignore` allowlist (`setup.md:10-25`): PASS, verified by build.**
  - Probe: `.agents/scratch/secprobe-p6-dockerignore-r3/`. It copies the round-2 fixture and adds the build inputs plus six new shapes. The list was extracted from the file by `awk`. Every probed path was confirmed present in the fixture first.
  - Build: `FROM scratch` + `COPY . /ctx/` through `docker buildx build --output type=local`, exit 0.
  - Out: `app.env`, `config/prod.env`, `creds/sa-key.json`, `keys/id_ed25519`, `keys/x.key`, `keys/x.pem`, `web/.npmrc`, `web/.env.local`, `web/.env.production`, `infra/terraform.tfvars`, `infra/terraform.tfvars.json`, `infra/prod.auto.tfvars.json`, `infra/terraform.tfstate.backup`, root `key.json`, root `.env`, `cmd/app/.env`, `internal/config/dev.env`, `.git/`, `web/node_modules`, `web/dist`, `web/.output`.
  - In: `go.mod`, `go.sum`, `cmd/app/main.go`, `internal/x/x.go`, `internal/db/migrations/0001_init.up.sql`, `web/embed.go`, `web/package.json`, `web/package-lock.json`, `web/src/app.tsx`. Also in: `internal/config/sa-key.json` and `web/localhost-key.pem` (the last finding).
  - A root `migrations/` stays out until its `!` line is added, as the paragraph says.
- **`cloud-run-deploy/SKILL.md:96`: PASS.**
  - The non-empty-audience clause matches `validate.go:160`. The explicit `iss` check is needed, because `:48` parses the issuer and nothing compares it, and the cell names it. That check also rejects an RS256 token someone mints with the IAP audience string (its `iss` is `accounts.google.com`), so the cell is sound.
  - `roles/iap.httpsResourceAccessor` is granted on the service, never the project.
  - The audience format and issuer string were not checked against Google's docs.
- **`:95`: PASS.** Terraform owns the public setting, and the flag is never passed for a service in HCL.
- **`:143`: PASS as a comment.** No command in the block prints the status code for the IAM-gated case. `curl -s -o /dev/null -w '%{http_code}' "${URL}/healthz"` would make "MUST 403" checkable. Informational.
- **`:71` `--immutable-tags`: PASS.** The flag spelling was not checked here (gcloud absent). The alias conflict is the third finding.
- **`service.md` item 3: PASS.** It says what is wired and gives the split as a project decision. What is left is the second and fourth findings.
- **`readJSON` (`service.md:88-98`): PASS.** It returns a fixed client message and logs the decoder error at Debug, which is off under the Info handler `main` builds. Informational:
  - It logs through the default `slog` logger, not the one passed in.
  - The comment's "log it" (`:93`) sits beside `:172`'s "never log … a request body". The error carries at most a fragment of the client's own input, such as a number literal or an unknown field's name.
- **`isHex` (`service.md:269,285,291-301`): PASS for both header forms.** The `X-Cloud-Trace-Context` span id is re-formatted from a parsed `uint64`. A caller can still pick a well-formed trace id, which groups only its own lines.
- **`x-goog-iap-jwt-assertion` in the never-log list (`service.md:172`): PASS.**
- **`apiError.Field` → `BadRequest.fieldViolations` (`service.md:70-75`): PASS, leaks nothing.**
  - `field` and `description` are the handler's constants (`"name"`, `"name is required"`). `readJSON` never sets `Field`, so no decoder text and no client value reaches the detail.
  - `json.Encoder` escapes `<`, `>` and `&`.
  - A future handler that formats client input into `Msg` would reflect it in `message` and `description` alike. That is the client's own input, JSON-encoded.
- **Migrate job (`resources-gcp.md:237-289`): PASS for the wiring.**
  - It has its own SA with `cloudsql.client`, `migrate_reads_db_url` (accessor on that one secret), a pinned version ref, and the `/cloudsql` volume and mount.
  - The job now holds the runtime DSN, which `:268-269` states.
  - The stale header comment is the second finding.
- **§ Public or IAM-gated (`resources-gcp.md:216-232`): PASS.** No live block remains. `allUsers` is a comment conditioned on Public: yes, and the live invoker binding is the IAM-gated `scheduler_invokes`.
- **`oidc_token` comment (`resources-gcp.md:316-317`): FAIL.** See the first finding.
- **`immutable_tags = true`: PASS** in both HCL sites (`resources-gcp.md:62`, `assets/example/main.tf:21`).
- **Cloud SQL `backup_configuration` with PITR: PASS** (`resources-gcp.md:116-119`).
- **`db_url` `sensitive` AND `ephemeral`: PASS at `:88`.** `:3` is stale; see the second finding.
- **`checks.md` `.gitignore`: PASS** (`:34-35`).
  - `*.auto.tfvars.json` is already matched by `*.tfvars.json`, which is harmless.
  - `terraform/SKILL.md:31` still names only `*.tfvars` among the ignored files. Informational.
- **Distroless by digest: PASS** (`docker/SKILL.md:78,86`).
- **Rule 5 `RUN --mount=type=secret` (`docker/SKILL.md:91-94`): PASS.** Informational, functional, confidence medium:
  - The secret mounts at `/run/secrets/netrc` by default, but `go` reads `$NETRC` or `~/.netrc`, so the private fetch fails as written.
  - Fix: add `target=/root/.netrc`, or `NETRC=/run/secrets/netrc` for that `RUN`. That keeps an agent from reaching for the build arg the rule forbids.
- **Left by decision, not re-raised:**
  - `gcloud-setup`'s `--key-file=key.json`. A root `key.json` is now out of the build context, verified by build. One under `internal/` or `web/` gets in (the last finding).
  - The link test's symlink reads.
  - Tool installs. Under the same decision: gcloud's apt and tarball path, and golangci-lint's `curl … | sh` at `service.md:404`, which is that project's documented path.
- **Carried, no security impact:** `cloud-run-deploy/SKILL.md:29-30` calls `deploy/gcp.md` "a cross-check, not a source" one sentence after making it the source.
- **Removed lines (`git diff HEAD`, the eight modified files, grepped for security terms):** every removed caution has a stronger replacement.
  - The `.env` cross-check is gone.
  - `--set-secrets …:latest` is now pinned.
  - Docker rules 1, 4 and 5 and the setup claim were rewritten.
  - The AWS "store secrets in variables" and "open security groups" rows have GCP counterparts (`terraform/SKILL.md:61-62`, the own-VPC block).

### Checked (round 3)

- **Secrets:** all 17 files in scope read whole, and no secret value in any of them. Write-only versions, pinned version refs, the tfvars rule and names-not-values all hold. The build-context allowlist was verified by build, using dummy files only.
- **Injection:** no new input surface.
  - `readJSON` is bounded at 1 MiB, refuses unknown fields and returns a fixed message.
  - Trace ids are hex-checked before they are interpolated into a resource name.
  - SQL uses `$1` parameters.
  - The HCL has no `templatefile`, `local-exec`, provisioner or `external` data source.
- **Authz/authn:** the scheduler token is the first finding and the database role the fourth. The IAP cell is sound against the library source. The runtime, migrate and scheduler SAs are separate. `run.invoker` is granted per service, `public` defaults to false, and Terraform creates no SA keys.
- **Crypto:** nothing composed. `google.golang.org/api/idtoken` is a vetted library consumed as designed. The first finding is a missing claim comparison on its result, not hand-rolling.
- **Dependencies:** `git diff HEAD --stat -- go.mod go.sum` is empty, and `skill_links_test.go` is stdlib only.
  - `govulncheck -show verbose ./...` on go1.27.1: exit 0, "Your code is affected by 0 vulnerabilities". It also reports one module-level vulnerability that is not called, identified as GO-2026-5024: the known one that `AGENTS.md` names, through the SQLite driver's `golang.org/x/sys`. Output: `.agents/scratch/secrev-p6-r3/govulncheck.txt`.
  - `npm audit`: not applicable. There is no `package.json` in scope, and round 1 found none tracked.
- **GCP surface:** in place: immutable tags, backups with PITR, the own VPC, `ssl_mode`, the IAM-gated default and the IAP accessor scope. The public `/tick` is the first finding and the tag alias the third.

### Not examined (round 3)

- `gcloud`, `terraform`, `tflint` and `trivy` are absent on this box. So:
  - no `terraform validate` of `assets/example/`;
  - the `--immutable-tags` spelling, the IAP audience format `/projects/N/locations/R/services/S` and the issuer string were not checked;
  - Artifact Registry's refusal of a re-pushed `:latest` (third finding) and Cloud SQL's `cloudsqlsuperuser` default for API-created users (fourth finding) are stated from documentation memory.
- `context/dev-setup.md`, which also has "`latest` only as an alias" (third finding), and `context/AGENTS.md`'s secret-shaped list, which lacks `*.tfvars.json` (round 1): out of scope and not re-read.
- No real `go build` through the skill's Dockerfile with the allowlist. The probe copies the context and does not compile.
- The six deleted terraform references (`ci-cd-workflows.md`, `code-patterns.md`, `module-patterns.md`, `quick-reference.md`, `testing-frameworks.md`, and `security-compliance.md` beyond its headings): not read in any round. Their removed lines were not re-grepped this round; only the eight modified files' removed lines were.
- Phase 6 files outside the scope list: `skills/frontend/**`, `skills/golang/SKILL.md` beyond line 34, `skills/golang/references/testing.md`, `skills/json-api-style/references/json-api-conventions.md`, `skills/gke-basics`, `skills/bash`, `skills/mermaid-syntax`, `skills/review*`, `skills/VENDORED.md`, the plan and the skill-matrix edits.

## Outcome (round 3)
- [medium] `/tick` check incomplete — fixed at `skills/terraform/references/resources-gcp.md` (both comments: `idtoken.Validate` with the audience, then `email` == the scheduler SA and `email_verified` true; why the audience alone proves nothing) (review cap reached; fixed under the user's same-turn decision of 2026-10-01 — "Go with your recommendation, drop the latest alias" — and re-checked by Phase 9's lenses, not by a fourth round here)
- [low] "DML only" comment at `:237` — fixed (what is wired; the split is a project decision); `:3` prose names `sensitive = true` and `ephemeral = true`
- [low] `:latest` alias vs immutable tags — fixed per D40
- [low] owner/DML recipe and `cloudsqlsuperuser` — fixed at `service.md` item 3 (roles created in SQL, not through `google_sql_user`)
- [low] allowlist paragraph overclaims inside re-admitted directories — fixed (`**/*.pem`, `**/*.key`, `**/*-key.json`, `**/*credentials*.json` added to the tail; the paragraph says the guarantee holds outside the re-admitted directories and names the rule for inside them)
