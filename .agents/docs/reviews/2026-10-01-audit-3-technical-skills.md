---
owner: reviewer
date: 2026-10-01
---
# Review — audit-3-technical-skills   (2026-10-01)

Range reviewed: no diff — `skills/{golang,typescript,frontend,bash,docker,json-api-style,mermaid-syntax,gcloud-setup,cloud-run-deploy,terraform,gke-basics,operate-gcp,data-gcp,api-auth-gcp}/**`, `skills/VENDORED.md` at HEAD 2ba4a50

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

No TS finding is open and nothing in the slice is broken at the level of "cannot compile or run". Two cross-skill contradictions remain, both medium: the terraform skill keeps the Cloud Run image in HCL while cloud-run-deploy deploys, canaries and rolls back the same service with gcloud (the next `terraform apply` undoes all of it); and golang's migration sample uses the lib/pq-backed golang-migrate driver while data-gcp prescribes the pgx/v5 driver.

## Findings

Audit I checklist (every TS ID, verified in the files at HEAD):

| ID | was | now | evidence |
|---|---|---|---|
| TS1 | high: terraform 100% AWS | fixed | `grep -rciE 'aws_\|\baws\b\|amazon\|s3\|dynamodb' skills/terraform` = 0 in every file (was 175); 0 CI/pre-commit/MCP demands (SKILL.md:8 only negates them); no LICENSE promise; example copied to scratch: `fmt -check` 0, `init -backend=false` + `validate` = "Success! The configuration is valid" (Terraform 1.16.4); checks.md ships the five bare commands; `for_each`/`count` rule stated correctly (SKILL.md "never both") |
| TS2 | high: ≈30 dangling golang pointers | fixed | golang is now SKILL.md + 2 references; `TestSkillLinksResolve` PASS (-v, ran); grep for "See <Title>", `assets/examples`, `templates/`, "Makefile template" = 0 |
| TS3 | high: flattened samber skills, no row/LICENSE | fixed | no persona blocks, Context7, "parallel sub-agents", `golang-samber.md`; golang is the harness's own clean-room text, so no VENDORED row or LICENSE is owed |
| TS4 | high: LICENSE label mismatch | struck (D28) | `skills/anti-ui-slop` absent; VENDORED.md:46 records the strike; 0 `anti-ui-slop`/`uizze` hits in scope |
| TS5 | high: Tailwind v3 / Vite SPA in references | fixed | grep `tailwind.config\|content:\s*\[\|main.tsx\|index.html` in frontend: only "no tailwind.config.*" prose and the SPA shell note; theme-system.md uses `@import "tailwindcss"` + `@theme inline` |
| TS6 | med: testify/uuid/sqlx/Cobra/mockgen against SKILL.md | fixed | grep `testify\|is.NotNil\|sqlx\|uuid.NewString\|cobra\|viper\|mockgen` in golang and bash = 0 (golangci-lint hits are lint, not the banned list) |
| TS7 | med: synctest.Run experimental | fixed | testing.md:`synctest.Test` GA 1.25, "the 1.24 Run experiment is gone"; block compiles and passes `go test -race` on go1.27.1; agrees with service.md version table |
| TS8 | med: tools.go, golangci v1 keys | fixed | `tools.go` appears only in service.md:424 as the "instead of" row; `.golangci.yml` is `version: "2"` with `default: standard` |
| TS9 | med: no GCP observability | fixed | service.md § Logging: Cloud Logging JSON contract, trace correlation, Error Reporting, Trace, Profiler; gcplog block vets clean |
| TS10 | med: off-harness workflow | fixed | grep `persona\|.modernize\|golang-modernize\|dependabot\|perf(scope)` = 0 |
| TS11 | med: terraform references | fixed | covered by TS1; references are checks.md, resources-gcp.md, upgrades.md only |
| TS12 | med: gke-basics MCP>gcloud, ask/never-ask, 2 dangling paths | partly | `diff -r` against upstream 6f0b877 (cloned to scratch): only SKILL.md differs, references/assets are verbatim, so the contradictions and the two dangling paths still exist in them; SKILL.md "Harness overrides" names and resolves each (gcloud over MCP, project asked once, both paths named). Deliberate (verbatim fork rule); `TestSkillLinksResolve` skips the skill, so nothing pins it |
| TS13 | med: anti-ui-slop vs frontend | struck (D28) | skill gone; frontend description now covers "ALL frontend work" with no competitor |
| TS14 | med: react-hook-form/zod/createServerFn | fixed | frontend SKILL.md:113 "never react-hook-form"; component-patterns.md § Forms native `<form action>`; `z.email()` with the deprecation noted; "No createServerFn: SPA mode" |
| TS15 | med: dark ring 1.8:1, hard-coded zinc/amber, inline-style toggle | fixed | recomputed with the WCAG formula: ring 6.31/5.94/5.68 light, 7.46/6.77/6.34 dark, error-on-tertiary 4.69, success 8.75 — all match theme-system.md; `--ring: var(--accent)` both themes; zinc appears only in the sentence saying the palette is removed; remaining `style={{}}` are motion values (animation-recipes, component-patterns:23), not theme styling |
| TS16 | med: cloud-run-deploy gaps | fixed | canary block, flags table, token verify present; every flag grepped in `gcloud run deploy --help` (SDK 587): `--iap`, `--startup-probe`, `--liveness-probe`, `--no-traffic`, `--tag`, `--cpu-boost`, `--vpc-egress`, `--ingress`, `--add-cloudsql-instances`, `--set-secrets` present, `--platform` absent as the skill says; `--to-tags/--remove-tags/--to-latest`, `--log-filter/--freshness`, proxy `--tag` present; health-checks.md:40-41 names the same probe flag |
| TS17 | low: golang:1.24/1.22 pins | fixed | docker SKILL.md pins by the go.mod line (`golang:<major.minor>`), distroless by digest; the gke-basics onboarding pin is covered by its override paragraph |
| TS18 | low: brew cask, no Linux path, Console answers | fixed | `brew install --cask gcloud-cli`; apt and tarball Linux paths; no "Cloud Console" answer left |
| TS19 | low: 2012 JSON style guide | fixed | json-api-style on AIP-140/142/126/158/193; error body `{"error":{code,status,message,details[ErrorInfo,BadRequest]}}` identical in golang service.md (adapter), json-api-conventions.md, frontend component-patterns.md and api-auth; tags camelCase in both skills |
| TS20 | low: bash SC2155/SC2046 etc. | fixed | shellcheck over the 47 bash blocks of bash-fundamentals and testing-linting: 0 SC2155, 0 SC2046 (the SC1072/SC2154/SC2034 hits are fragments with free variables); `IMAGE_TAG` defined once (makefile-patterns.md:128); `$(error)` only inside a target-time `$(if)`; no `/var/log` |
| TS21 | low: CONTEXT.md | fixed | VENDORED.md Notes carry the "read it as the project's AGENTS.md" line; codebase-design untouched |
| TS22 | low: anti-ui-slop NOTICE | struck (D28) | skill gone |
| TS23 | low: python | struck (D28) | `skills/python` absent; only `python3` (interpreter) and the kubernetes-python client row in gke-basics remain |
| TS24 | med: wg.Go / map internals / sql.Open / samber samples | fixed | wg.Go is 1.25 everywhere (SKILL.md, service.md table, testing.md); no samber; `sql.Open` error checked in migrate.go; all 10 complete Go files vet clean, 5 test blocks vet and `go test -race` pass (see Criteria 2) |
| TS25 | med: dist/client committed vs Node stage | fixed | frontend SKILL.md and docker SKILL.md agree: `dist/client` gitignored, built by the `node:24-alpine` stage, copied before `go build`, `web/dist` in `.dockerignore` |
| TS26 | med: .dockerignore twice | fixed | one list in docker/references/setup.md § Step 1; SKILL.md rule 4 points to it |
| TS27 | low: gke samples discard errors, deletion_protection=false | partly | `client-library-usage.md` and `iac-usage.md` unchanged (verbatim by rule); the override paragraph says "illustrative, not copy-paste" and names both defects |
| TS28 | low: mermaid %%{init}, packet-beta, provenance | fixed | SKILL.md:167 states `%%{init}` is the older accepted form and that frontmatter `config:` is used throughout; other-diagrams.md:405 uses `packet` (11.4) and says why; VENDORED note records "harness's own skill, no upstream identified" |
| TS29 | low: typescript thin | fixed | SKILL.md + `references/tsconfig.md` + `references/testing.md`; plain-V8 rule and tsconfig baseline present |

New findings:

- **[medium / confidence: medium]** `skills/terraform/SKILL.md:15` (also `references/resources-gcp.md:176,268`, `assets/example/main.tf:66`) vs `skills/cloud-run-deploy/SKILL.md` Push & Deploy / Canary / Rollback — the skill says "the image, the new revision and the rollback stay with cloud-run-deploy" while the HCL pins `image = "...:${var.git_sha}"` and owns the service; there is no `lifecycle { ignore_changes = [...] }` anywhere. A `gcloud run deploy` (new image, `--no-traffic --tag canary`, `update-traffic`, a rollback `--to-revisions`) changes the same fields Terraform manages, so the next `terraform plan` shows a diff and an apply resets the image and the traffic split to what HCL says, silently undoing a rollback. Fix: pick one owner of the image/traffic fields — either state that a service in HCL is deployed by `terraform apply -var git_sha=…` and cloud-run-deploy's canary is not used for it, or add `ignore_changes` on `template[0].containers[0].image`, `client`, `client_version` and the traffic block to the example and say the deploy skill owns them.
- **[medium / confidence: high]** `skills/golang/references/service.md:366` vs `skills/data-gcp/SKILL.md:90,119` — the compiling migrate sample imports `github.com/golang-migrate/migrate/v4/database/postgres` (its `go mod tidy` in scratch pulled `github.com/lib/pq v1.10.9`, indirect), while data-gcp prescribes `database/pgx/v5` and golang's own Dependencies list names pgx and golang-migrate only. Two skills give two drivers and the golang one adds an unlisted module. Fix: change the sample to `pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"` with `pgxmigrate.WithInstance(sqlDB, &pgxmigrate.Config{})` and the database name that driver registers, then re-run tidy and vet.
- **[low / confidence: medium]** `skills/terraform/references/resources-gcp.md:140` — the IAM-authentication bullet tells the reader to create one `google_sql_user` per SA but shows no HCL for it and no `roles/cloudsql.instanceUser` binding, which data-gcp/SKILL.md (path 2) says each SA needs besides `cloudsql.client`; only `cloudsql.client` is in the HCL (lines 132-136, 246-249). A Terraform-managed instance following only this file hits an IAM login failure. Fix: add the `google_sql_user` and `google_project_iam_member` instanceUser blocks, or point to data-gcp for the roles in that bullet.
- **[low / confidence: low]** `skills/docker/SKILL.md:3`, and `frontend`, `cloud-run-deploy`, `gcloud-setup`, `typescript`, `bash` — descriptions are `|` block scalars (several lines), not one line; each is under 1,024 characters and `TestLintDescriptionLimit` passes, and AGENTS.md's one-line rule names agents, so this is only a risk note for Antigravity's frontmatter reader (BUG-004 class). No change unless a sandbox start shows one dropped.

## Criteria
1. Dead code — n/a — audit
2. Flow tracing — n/a — audit (evidence gathered for the findings: `TestSkillLinksResolve` and the description tests PASS; service.md ×4 files, api-auth-gcp go-samples ×4, data-gcp connect.md ×2 compiled with `go mod tidy && go vet` on go1.27.1; testing.md ×5 blocks vet and `go test -race` pass; `httptest.NewTestServer`, `synctest.Test`, `slog.NewMultiHandler` exist in `go doc`; terraform example validates; `gcloud` command groups of operate-gcp, data-gcp, api-auth-gcp, cloud-run-deploy resolve on SDK 587)
3. Standards conformance — n/a — audit
4. Directory layout — n/a — audit
5. Simplicity (DIY-except-security) — n/a — audit (no hand-rolled auth: IAP, Firebase and OIDC samples call `idtoken` and the Firebase Admin SDK)
6. Decoupling — n/a — audit

## Verdict (counts)

Fixed 23, partly 2 (TS12, TS27: verbatim upstream references stay as they are, named by the override paragraph), open 0, struck 4 (TS4, TS13, TS22, TS23 under D28). Highest open severity: medium (the two new cross-skill findings above; none of the 29 original IDs is open).

Better than on 2026-09-30, clearly. The two numbers the first audit measured: terraform AWS identifiers 175 → 0 (and `terraform validate` now succeeds on a shipped example, where the old skill had invalid HCL), and golang dangling pointers ≈30 → 0 with a test (`TestSkillLinksResolve`) that keeps it so for every non-vendored skill. Every Go block in the slice that claims to compile does (15 blocks compiled, none failed; the 3 fragments in health-checks.md and json-api-style were not), the contrast table in frontend matches a recomputation to the second decimal, the cloud-run-deploy flags exist in gcloud 587 and `--platform` is absent as stated, and gke-basics differs from its pin in SKILL.md only. What remains is seam work between the new skills (terraform vs the gcloud deploy path, golang vs data-gcp migrate driver), not content decay.

## Not examined
- `golang-migrate` and the `service.md` migrate sample were compiled but not run against a database; the data-gcp SQL roles recipe (GRANT, ALTER DEFAULT PRIVILEGES, assertion migration) was read, not executed on Postgres 18; the `connect.md` IAM dialer was vetted, not run.
- Not executed: the frontend TanStack Start build (version pins in SKILL.md are claims from the skill's own measurement), the typescript skill's tsconfig/vitest references, the docker three-stage Dockerfile (no `docker build`), `trivy` and `tflint` against the example (`terraform validate` and `fmt` only).
- Not read line by line: `skills/gke-basics/references/*` beyond the diff against the pin (verbatim), `mermaid-syntax` examples and class/flowchart/sequence references, `bash/references/makefile-patterns.md` beyond TS20's lines, `frontend/references/{design,animation-recipes}.md`, `typescript/references/*`, `json-api-style` conventions beyond the error section, api-auth-gcp prose after § Firebase, `operate-gcp/SKILL.md` beyond its commands, the middle of `data-gcp/references/connect.md`.
- Shellcheck covered the bash skill's two references only, not the shell blocks of the cloud skills; `terraform` HCL blocks in `resources-gcp.md` other than the shipped example were not extracted and validated.
- No Antigravity run (the audit is read-only); the dispatch's "one line" description check was done by reading frontmatter, not through an Antigravity start.
