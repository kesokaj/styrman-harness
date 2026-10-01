---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p6-code   (2026-10-01)

Range reviewed: working tree (uncommitted) — Phase 6 of the audit-fix plan: `skills/{terraform,golang,json-api-style,frontend,typescript,cloud-run-deploy,docker,gcloud-setup,bash,mermaid-syntax}/**`, `skills/gke-basics/SKILL.md`, `skills/VENDORED.md`, `AGENTS.md`, `skills/init-project/assets/agents-tree/docs/deploy/gcp.md`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/plans/2026-09-30-audit-fix.md`; deleted `skills/python/`, `skills/anti-ui-slop/`, the old golang and terraform references; untracked `internal/harness/skill_links_test.go`.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Nothing in the Go test is broken; the defects are in the skill text an agent will follow. Worst first: the docker skill tells an agent to copy the node stage's output "next to the binary" into the distroless runtime, while the frontend skill embeds `dist/client` with `go:embed` at `go build` time and says the docker skill has the Dockerfile (it has none), so the documented path produces a failing build or a binary with no UI. Two skills name skills that do not exist (`api-auth-gcp`, `data-gcp`). The Go service sample's `BaseContext` cancels in-flight requests on SIGTERM. The Terraform example can create the Cloud Run service before its secret has a version.

## Findings
- **[medium / high]** `skills/docker/SKILL.md:87-90` vs `skills/frontend/SKILL.md:106-110` and `skills/docker/references/setup.md:16` — contradiction on the frontend build. docker: node stage output is copied "next to the binary", runtime stays distroless. frontend: the Go stage copies it, `embed.go` embeds `dist/client`, "the docker skill has the Dockerfile". The docker skill has no such Dockerfile. `.dockerignore` lists `web/dist`, so `COPY . .` in the Go stage does not bring a host build either. Following docker as written, `go build` fails on the empty embed pattern (or the binary embeds nothing). Fix: put the real three-stage Dockerfile in docker (node stage first, `COPY --from=web /web/dist/client web/dist/client` into the golang stage BEFORE `go build`) and delete "next to the binary".
- **[medium / high]** `skills/cloud-run-deploy/SKILL.md:96,100` — points agents to `api-auth-gcp` and `data-gcp` skills; `ls skills/` has neither (grep of `context/`, `agents/`, `skills/` finds no other mention). The link test checks paths, not skill names, so it cannot catch this. An agent told to validate `x-goog-iap-jwt-assertion` with no skill to read will hand-roll JWT validation (forbidden by AGENTS.md hard rules). Fix: point `data-gcp` at `golang` § `references/service.md` Database; for IAP name a vetted library (for example `google.golang.org/api/idtoken`) or drop the sentence.
- **[medium / medium]** `skills/golang/references/service.md`, `run()` — `BaseContext: func(net.Listener) context.Context { return ctx }` where `ctx` is the SIGINT/SIGTERM `signal.NotifyContext`. On SIGTERM every in-flight request's `r.Context()` is cancelled immediately, before `srv.Shutdown` drains, so handlers using pgx or `*Context` calls abort. It defeats the graceful shutdown the sample and prose claim. Fix: remove `BaseContext`, or return `context.WithoutCancel(ctx)`.
- **[medium / medium]** `skills/terraform/assets/example/main.tf` (`google_cloud_run_v2_service.app`) and the Cloud Run block in `references/resources-gcp.md` — the service reads the secret with `version = "latest"` but `depends_on` omits `google_secret_manager_secret_version.app`; the secret itself is an implicit dependency, the version is not. On a first apply the service can be created before a version exists and Cloud Run rejects the revision. Fix: add the version resource to `depends_on`. Not verified by running Terraform (binary not installed here).
- **[medium / medium]** `skills/cloud-run-deploy/SKILL.md` § Canary rollout — "/preflight's deploy row runs this block; only a rollback skips it", but the first `gcloud run deploy` of a service that does not exist yet cannot use `--no-traffic` (gcloud refuses it when creating a service, from memory of the CLI, not checked on 587; the phase-6 fork report only verified the flag exists). No first-deploy case is stated anywhere in the skill. Fix: state "first deploy: no `--no-traffic`/`--tag`; canary starts from the second revision", or verify and drop the claim.
- **[low / medium]** `skills/cloud-run-deploy/SKILL.md` (`--set-secrets` row, "pin the version, not `:latest`") vs `skills/terraform/assets/example/main.tf` and `resources-gcp.md` (`version = "latest"`) — two skills disagree on one point; the terraform example is the one an agent copies. Fix: pick one and say why in both (for Terraform, `google_secret_manager_secret_version.app.version`, which also gives the dependency above).
- **[low / high]** `skills/golang/references/testing.md` benchmark/fuzz block — the file header says "Every block below is a complete file that compiles"; the `package greet` block uses `Hello`, which only the first block defines, so alone it fails `go vet` (`undefined: Hello`; I extracted the seven stdlib-only Go blocks and vetted each: six compile, this one does not). Fix: reword the claim or include `Hello`.
- **[low / medium]** `internal/harness/skill_links_test.go` via `vendoredSkills` (`harness_test.go:984`) — takes the first cell of every VENDORED.md table row, so the fork `gke-basics` (second table, "NOT verbatim", `SKILL.md` edited this phase) is skipped by the link check. Intentional-looking (its upstream dangling paths are named in the Harness overrides paragraph) but the test comment says "vendored verbatim"; its own edited `SKILL.md` text is unchecked.

## Criteria
1. Dead code — pass: python, anti-ui-slop, seven golang and six terraform references deleted; no orphaned references found, and the link test finds none dangling. `skills/python`/`anti-ui-slop` restatements in `context/`, `agents/`, `README.md`, `AGENTS.md`, `VENDORED.md` are gone (grep; the phase-6 verifier log shows the same).
2. Flow tracing — fail: three cross-file seams broken (docker/frontend build, cloud-run-deploy/missing skills, terraform/cloud-run-deploy secret version). `skill_links_test.go`: reasoned from the code — for a non-vendored skill a markdown link or backticked `references|assets|...` path outside fences that resolves against neither the file's dir, the skill root, nor `skills/` is appended to `missing` and fails with `t.Errorf`; ran it bare, `ok`. I did not run a must-fail mutation (would require editing a skill file; reasoned only). It checks paths only, not skill names.
3. Standards conformance — mixed: Go samples follow the golang skill (errors wrapped, `t.Context`, `wg.Go` marked 1.25, go line raised); the SC2155/SC2046 patterns are gone from the bash blocks (shellcheck on the extracted blocks, those codes: clean; the `for f in $(ls …)` left is labelled "Don't do this"); HCL has no `for_each`+`count` on one resource, the provider is pinned `~> 8.0`, and the secret version takes `var.app_secret` (sensitive, ephemeral), no literal. Docker sample is multi-stage, non-root (`USER nonroot`), `golang:<go.mod line>` wording matches gke-basics' override. Frontend build order `vite build && tsc -b` is consistent across frontend, typescript and docker (`npm run build`). Failures are the findings above. Unverifiable version claims (Terraform 1.16, provider 8.x, `POSTGRES_18`, golangci-lint v2.14.0, `httptest.NewTestServer` in 1.27, tflint ruleset 0.40.0, `secret_data_wo_version` "is a string") were not checked.
4. Directory layout — pass: new references under `skills/<name>/references/`, HCL under `skills/terraform/assets/example/`, test beside the harness tests in `internal/harness/`.
5. Simplicity (DIY-except-security) — pass with one exception: the IAP JWT pointer above invites hand-rolled JWT validation; no new Go dependency (the samples use pgx and golang-migrate, both already listed as accepted; `database/postgres` of golang-migrate pulls `lib/pq`, which the Dependencies list does not name).
6. Decoupling — pass: each skill is a self-contained directory; the cross-links are by skill name and the test pins paths. The skill-name pointer seam is where it broke.

## Not examined
- `terraform validate`/`fmt` on the HCL (no `terraform`/`tofu` binary here); HCL read only. Provider attribute names and version claims not checked against the registry.
- gcloud flag existence: not rerun; relied on the fork reports `.agents/scratch/audit-fix/phase6-*.md` only for the existence of flags, and read the diff for the rest.
- Go blocks using pgx and golang-migrate (`db.go`, `migrate.go`) were read, not compiled (modules not available offline).
- `skills/frontend/references/*` (accessibility, animation-recipes, component-patterns, theme-system; about 950 changed lines) and `skills/typescript/references/tsconfig.md` and most of `typescript/references/testing.md`: skimmed or grepped for the build-order and `dist/client` contradictions, not read through.
- `skills/json-api-style` references past the diff head, `skills/review-standards/SKILL.md`, `skills/review/SKILL.md`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/plans/2026-09-30-audit-fix.md`: not reviewed (outside the code lens or not reached).
- `skills/terraform/references/checks.md` and `upgrades.md` read; their `.tflint.hcl` plugin pin and Terraform minor claims not verified.

## Outcome
- [medium] docker vs frontend on `dist/client` — fixed at `skills/docker/SKILL.md` § Dockerfile (three stages; the Go stage copies `web/dist/client` before `go build`; the paragraph below it says why) — frontend's pointer now holds
- [medium] `api-auth-gcp` / `data-gcp` named before they exist — rebutted: D38 (Phase 7 is the next commit of the same plan); the two cells now say so and carry the interim answer (`idtoken.Validate` + issuer check; `service.md` § Cloud SQL), so no agent hand-rolls a JWT check meanwhile
- [medium] `BaseContext` cancels in-flight requests on SIGTERM — fixed at `skills/golang/references/service.md` `run()` (BaseContext removed; ctx only triggers `Shutdown`); block re-extracted and `go vet` clean in `.agents/scratch/go-samples/service_main`
- [medium] Cloud Run service lacks `depends_on` on the secret version — fixed at `skills/terraform/assets/example/main.tf` (version pinned to `google_secret_manager_secret_version.app.version`, resource in `depends_on`); `terraform validate` Success, `fmt -check` clean
- [medium] first-deploy `--no-traffic` claim — fixed at `skills/cloud-run-deploy/SKILL.md` § Canary rollout (the block runs from the second revision on; the first deploy is plain)
- [low] `--set-secrets` pin vs terraform `latest` — fixed (both terraform sites pinned to the managed version)
- [low] testing.md header claims every block compiles alone — fixed (the header now says beside the first block's `greet.go`)
- [low] link test skips the `gke-basics` fork entirely — fixed at `internal/harness/skill_links_test.go` (`forkSkills`: a fork's SKILL.md is checked, its references/assets are not); must-fail: a dangling link appended to `gke-basics/SKILL.md` failed the test, restored

## Round 2

Range reviewed: working tree, `git status --short --untracked-files=all` (Phase 6 of the audit-fix plan); diffs read with `git diff HEAD -- <files>`; `internal/harness/skill_links_test.go` read whole.

Status: COMPLETE

Verdict: APPROVE WITH NITS. All eight round-1 fixes hold; the one rebuttal (D38, Phase 7 skills) is accepted because the plan lists both skills as Phase 7 boxes (plan lines 171-172, D24) and the cells carry the interim answer. No new high finding. Two medium defects introduced or left in the fixed text: the Postgres 18 compose volume path, and the migrate job in the HCL reference that cannot reach the database.

### Verification of the fixes
- docker three-stage Dockerfile: node stage output `/src/web/dist/client` matches the `COPY --from=web` source; the copy lands after `COPY . .` and before `go build`; `web/dist` in `.dockerignore` is consistent. Holds.
- terraform: assembled every hcl block of `resources-gcp.md` (backend block stripped) with `var.project|region|app|repo|git_sha` and a sensitive `db_url`, in `.agents/scratch/rv-r2`, `terraform init` then `validate`: Success. Every reference resolves (`google_compute_network.app`, `google_compute_subnetwork.app`, `google_service_account.migrate`, `.scheduler`, `.db_url` version). The example root was validated by the main agent; I read it and found `version` pinned and the version in `depends_on`.
- `service.md` Go block: imports (context, encoding/json, errors, fmt, log/slog, net/http, os, os/signal, syscall, time) are all used; `BaseContext` gone; `Reason`, `errorDomain`, `ErrorInfo` all used by `adapt`. No exported leftover.
- Form sample (`component-patterns.md:196-221`): `ApiError`, `fieldErrors`, `res.json() as ApiError`, `?? []` chaining type-check by inspection (not compiled).
- `skill_links_test.go`: ran bare, ok; `gofmt -l` clean, `go vet` clean. `forkSkills` finds `gke-basics` by the header `Publisher / author`; the first table's header contains `| Source |` and resets the flag. `vendoredSkills` still has a caller. Did not rerun the must-fail mutation (relied on the main agent's report).
- cloud-run-deploy canary, `--iap` cell, `.env` drop, frontend theme wording (consistent across SKILL.md, theme-system.md:3 and the script at :182/:208), typescript `reuseExistingServer: false`, VENDORED/INDEX/gcp.md template/skill-matrix edits: read, consistent.

### Round 2 findings
- **[medium / medium]** `skills/docker/references/setup.md:38,46` — the compose file bumped to `postgres:18` but still mounts the named volume at `/var/lib/postgresql/data`. From memory of the docker-library change for 18 (not verified offline): the image's PGDATA moved to `/var/lib/postgresql/18/docker` and its VOLUME to `/var/lib/postgresql`, so the named volume stays empty and the data lives in an anonymous volume that `docker compose down` orphans. `SKILL.md`'s "Compose DB keeps old schema | `docker compose down -v`" row assumes the named volume holds the data. Fix: mount `pgdata:/var/lib/postgresql`, or pin `postgres:17` until verified; test by `up`, create a table, `down`, `up`, read it.
- **[medium / medium]** `skills/terraform/references/resources-gcp.md` § Job and Scheduler — the migrate job as written cannot migrate: no `cloud_sql_instance` volume, no `DATABASE_URL` secret env, and `google_service_account.migrate` is granted only `roles/cloudsql.client`, no `secretAccessor` on any secret (the DDL role's password is not modelled at all). `service.md` item 3 says the job owns the schema under its own SA. An agent copying the block gets a job that fails at `execute --wait`. Fix: add the volume, a secret env for the DDL user's URL with its own secret and accessor binding, or state in the block that those are omitted.
- **[low / medium]** `skills/terraform/references/resources-gcp.md` § APIs vs § Cloud Run v2 — the network and subnetwork blocks (`google_compute_*`) need `compute.googleapis.com`, which is not in the `apis` list (run, artifactregistry, secretmanager, sqladmin, cloudscheduler, iam); a first apply in a fresh project fails. Fix: add it with the "only with a private-IP target" note. Not verified by apply.
- **[low / high]** `skills/frontend/references/component-patterns.md:206,216` — `fieldErrors` returns `{ _: message }` when there are no field violations, commented "show the message once", but `SignupForm` only reads `err("email")`; a non-field error (409, 500, 429) renders nothing, and a non-JSON body (a 502 HTML page) makes `res.json()` throw inside the action. Fix: render `state.errors._` in a `role="alert"` paragraph and guard the parse.
- **[low / high]** `skills/golang/references/service.md` Database item 3 — the clause "(`gcloud run jobs execute --wait`) from `/preflight` before the new revision receives traffic" is now detached from "run as a Cloud Run Job" by the inserted own-SA sentence and reads as something a compromised instance cannot do; the sequencing rule (migrate before traffic, never on start) is lost. Reflow into two sentences.
- **[low / high]** `skills/cloud-run-deploy/SKILL.md:100` — points to `service.md` § Cloud SQL; the heading is "Database: Postgres on Cloud SQL" (`service.md:288`). The link test checks paths, not headings.
- **[low / low]** `skills/terraform/references/resources-gcp.md` — comment on the Cloud Run block says "add the version resource to depends_on" but the block has no `depends_on`; the `.version` reference already creates the edge, so the instruction is redundant, not wrong.

### Criteria (round 2)
1. Dead code — pass: no leftover declaration found from the fixes (`BaseContext` removed with no orphan; `vendoredSkills` still called).
2. Flow tracing — mostly pass: HCL resolves and validates; docker/frontend seam closed; failures are the migrate job wiring and the Postgres 18 volume.
3. Standards conformance — pass with the stale-description items above (`_` message, `§ Cloud SQL`).
4. Directory layout — pass.
5. Simplicity (DIY-except-security) — pass: IAP cell names `idtoken.Validate` plus an `iss` check, no hand-rolled JWT.
6. Decoupling — pass.

### Not examined (round 2)
- `terraform fmt` of the example root (main agent reports clean); the assembled reference blocks differ from fmt only in comment alignment, not checked as a requirement.
- Postgres 18 volume behaviour, Cloud SQL `ENCRYPTED_ONLY` with the socket, `idtoken.Validate` not checking `iss`, the `cloudsqlconn` Go floor: from memory, not run.
- TypeScript form sample not compiled; `skills/frontend/references/{accessibility,animation-recipes}.md`, `typescript/references/tsconfig.md`, bash references and mermaid references not reread this round; unchanged-since-round-1 files keep round 1's coverage.
- The `skill_links_test.go` must-fail mutation not rerun.

## Outcome (round 2)
- [medium] compose `postgres:18` with the old `/var/lib/postgresql/data` mount — fixed at `skills/docker/references/setup.md` (`pgdata:/var/lib/postgresql`; `docker image inspect postgres:18` on 2026-10-01: VOLUME `/var/lib/postgresql`, PGDATA `/var/lib/postgresql/18/docker`)
- [medium] migrate job has no secret, volume or accessor — fixed at `skills/terraform/references/resources-gcp.md` (the job reads the same `db_url` secret and mounts `/cloudsql`; `migrate_reads_db_url` IAM member added); assembled blocks `terraform validate` Success, `fmt -check` clean
- [low] `compute.googleapis.com` missing from the APIs — fixed (added to the `for_each` set)
- [low] the form's `_` message never rendered — fixed at `skills/frontend/references/component-patterns.md` (rendered after the fields); the block type-checks in the scratch app (`tsc -b` exit 0)
- [low] `service.md` item 3 garbled — fixed (rewritten: the SA split is wired, the role split is stated as a project decision with the `ALTER DEFAULT PRIVILEGES` shape)
- [low] `cloud-run-deploy:100` heading name — fixed (`§ Database: Postgres on Cloud SQL`)

## Round 3

Range reviewed: working tree, `git status --short --untracked-files=all` (Phase 6 of the audit-fix plan); diffs read with `git diff HEAD -- <files>`; `internal/harness/skill_links_test.go` read whole (unchanged since round 2).

Status: COMPLETE

Verdict: APPROVE WITH NITS. Every round-2 fix holds as stated, and the mechanism is settled: nothing in this round is a defect in the code samples, the HCL, or the Go. One medium defect introduced by a round-2 fix: `immutable_tags = true` (and `--immutable-tags`) on the Artifact Registry repository contradicts the documented `:latest` convenience alias, which a second push then fails on. The rest are low.

### Verification of the fixes
- docker `setup.md`: compose mounts `pgdata:/var/lib/postgresql` (matches the postgres:18 VOLUME per the recorded `docker image inspect`; not re-run). The `.dockerignore` allowlist: last-match-wins order is right (`*`, `!` re-admits, then the `web/` denials and the `**/` secret globs); every path the three-stage Dockerfile copies (`go.mod`, `go.sum`, `cmd/`, `internal/`, `web/`) is re-admitted and the `web/dist` denial agrees with the `COPY --from=web` into the Go stage. Holds.
- docker `SKILL.md`: rule 1 digest pin on the runtime `FROM` matches the rule text; `USER nonroot` exists in every distroless static variant; the three-part-tag example names no version. Rule 5 names `RUN --mount=type=secret` (see finding 3).
- terraform `resources-gcp.md`: re-assembled all eleven hcl blocks (backend stripped, `db_url` sensitive and ephemeral) in `.agents/scratch/rv-r3` with `.agents/scratch/tools/terraform`: `validate` Success, `fmt -check` clean. The migrate job reads `google_secret_manager_secret_version.db_url.version`, mounts `/cloudsql`, and `migrate_reads_db_url` grants accessor; `compute.googleapis.com` is in the set; Cloud SQL `backup_configuration` is accepted; the Public/IAM section is comments only and points at `scheduler_invokes`, which exists. The example root `fmt -check` clean; `docker_config` present in both. `checks.md` `.gitignore` has `*.tfvars.json` and `*.auto.tfvars.json`.
- golang `service.md`: re-extracted the main and gcplog blocks and ran `go vet`: clean. `apiError.Field` feeds a `BadRequest.fieldViolations` detail only when set; `readJSON` logs the decoder error at Debug and returns a fixed message; `isHex` checks length and alphabet on both trace paths (traceparent and `X-Cloud-Trace-Context`); never-log list names `x-goog-iap-jwt-assertion`. `golang/SKILL.md:56` and `gke-basics/SKILL.md:23` carry `<major.minor of the go.mod go line>`, consistent with docker rule 1.
- cloud-run-deploy: `--immutable-tags` exists (`gcloud artifacts repositories create --help` on the vendored SDK in `.agents/scratch/tools`); the IAP cell names a non-empty audience and the issuer check; `:100` heading is correct.
- frontend: the form block has `export function SignupForm` and renders the `_` message; the scratch app at `.agents/scratch/fe-app` was type-checked by the main agent (not re-run by me). `theme-system.md:324` reads consistently with the toggle text.
- `skill_links_test.go`: `go test ./internal/harness -run 'SkillLinks|Link'` ok; `gofmt -l internal/harness` clean.

### Round 3 findings
- **[medium / medium]** (new defect introduced by a round-2 fix) `skills/terraform/references/resources-gcp.md:62`, `skills/terraform/assets/example/main.tf:21`, `skills/cloud-run-deploy/SKILL.md:71` vs `:54` — the repository is now created with immutable tags, yet `cloud-run-deploy/SKILL.md:54` still says "add :latest only as a convenience alias" (and the global context's tagging rule says the same). Artifact Registry refuses to move or delete an immutable tag, so the second `docker push …:latest` fails, and so does re-pushing a git-SHA tag after a rebuild of the same commit (a rebuild is not byte-identical), which the deploy loop does after a failed canary. From memory of the Artifact Registry behaviour, not tried against a live repository. Fix: either drop the `:latest` alias clause at `cloud-run-deploy:54` and say a rebuild of the same commit needs a new tag (for example `${TAG}-2`), or leave tags mutable and say the SHA tag is a convention.
- **[low / medium]** `skills/docker/SKILL.md` rule 5 — `RUN --mount=type=secret,id=netrc` is named without `target=/root/.netrc`; the default mount is `/run/secrets/netrc`, which `go mod download` never reads, and the Dockerfile above carries no such mount. As written an agent following rule 5 gets a failing private-module fetch. Fix: write the full `RUN --mount=type=secret,id=netrc,target=/root/.netrc go mod download` line (the golang image's HOME is `/root`; verify).
- **[low / medium]** `skills/frontend/references/component-patterns.md:~205-217` — two remaining gaps in the form sample. (a) The "Forms" prose says React 19 "resets the fields on success"; React 19 also resets uncontrolled fields after an action that returns errors, so the user's typed email vanishes on a validation error unless `defaultValue` is fed back from state (from memory of facebook/react#29034, not run). (b) `await res.json()` still throws on a non-JSON body (a 502 HTML page), which round 2 asked to guard; the `_` rendering fixed only the first half. Fix: return the submitted values in state and use them as `defaultValue`, and `res.json().catch(() => null)` with a fallback message.
- **[low / low]** `skills/docker/SKILL.md` rule 1 — says "a minor written in this skill goes stale" and in the same sentence writes "`cloudsqlconn` needs 1.26 or later" and gives `go 1.27.1 → golang:1.27` as the example; harmless as an example, noted only because the claim and the sentence disagree.

### Criteria (round 3)
1. Dead code — pass: the allowlist `.dockerignore` and the `Public or IAM-gated` comment block leave no orphan; no unreferenced declaration in the vetted Go blocks.
2. Flow tracing — pass with one fail: HCL assembles and validates, the Go blocks vet, docker and frontend seams hold; the immutable-tags / `:latest` alias seam between terraform and cloud-run-deploy is the finding.
3. Standards conformance — pass with the stale-description items above (`:latest` alias, form reset prose).
4. Directory layout — pass.
5. Simplicity (DIY-except-security) — pass: IAP stays `idtoken.Validate` plus an issuer check; no hand-rolled JWT or crypto.
6. Decoupling — pass.

### Not examined (round 3)
- Artifact Registry immutable-tag behaviour, React 19 form reset, and BuildKit secret mount default target: from memory, not run.
- Postgres 18 volume path: relied on the recorded `docker image inspect`, not re-run (no daemon use this round).
- `fe-app` type-check not re-run by me; the Playwright and `web/` build were not run.
- Files unchanged since round 2 and not in the round-3 list (bash, mermaid, json-api-style, typescript references, `accessibility.md`, `animation-recipes.md`, the plan, skill-matrix, `review-standards`/`review` skills) keep rounds 1–2 coverage.
- `skill_links_test.go` must-fail mutation not rerun.

## Outcome (round 3)
- [medium, from a round-2 fix] immutable tags vs the `:latest` alias — fixed per D40 at `skills/cloud-run-deploy/SKILL.md:54`, `skills/init-project/assets/agents-tree/docs/deploy/gcp.md:27`, `skills/bash/references/makefile-patterns.md` (the `docker tag … :latest` line removed) (review cap reached; fixed under the user's same-turn decision of 2026-10-01 — "Go with your recommendation, drop the latest alias" — and re-checked by Phase 9's lenses, not by a fourth round here)
- [low] `RUN --mount=type=secret,id=netrc` without a target — fixed (`target=/root/.netrc`)
- [low] form resets fields on error too; `res.json()` unguarded — fixed at `skills/frontend/references/component-patterns.md` (prose says every action resets uncontrolled fields; `.catch(() => null)` with a status fallback); `tsc -b` exit 0 in the scratch app
