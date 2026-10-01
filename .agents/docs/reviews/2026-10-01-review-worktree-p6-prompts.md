---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p6-prompts   (2026-10-01)

Range reviewed: working tree (uncommitted Phase 6 changes), prompts lens (`skills/review-standards/references/prompt-audit/prompt-audit.md`, Steps 0-6). Instruction files read whole: `skills/terraform/SKILL.md` + `references/{checks,resources-gcp,upgrades}.md` + `assets/example/main.tf` (grep and lines 60-75 only); `skills/golang/SKILL.md` + `references/{service,testing}.md`; `skills/json-api-style/SKILL.md` + reference; `skills/frontend/SKILL.md` + `references/{design,theme-system,accessibility,component-patterns,animation-recipes}.md`; `skills/typescript/SKILL.md` + `references/{tsconfig,testing}.md`; `skills/cloud-run-deploy/SKILL.md`; `skills/docker/SKILL.md` + `references/{health-checks,setup}.md`; `skills/gcloud-setup/SKILL.md`; `skills/gke-basics/SKILL.md` (Harness overrides paragraph only, its claims checked against the references); `skills/bash/SKILL.md` + `references/*.md` (changed hunks read in full, the rest by grep); `skills/mermaid-syntax/SKILL.md` + `references/*.md` (SKILL.md whole, the changed hunks of the references, one reference's head).

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Step 0 assumptions: scope as listed; target models Claude and Gemini (every model that reads a skill), so Group 1 findings with a model-specific reason are `flag`; Group 2 stale-fact and conflict findings are repository-grounded and actionable as written. Group 4: not applicable (no agent definitions, no request code). The worst problem: Phase 6 left five cross-skill conflicts where two instruction files give opposite instructions for the same thing, and a model following one produces something the other rejects. The Go sample's `/healthz` returns 204 and the docker skill's `-healthcheck` exits 1 on anything but 200; Terraform's example reads secrets at `version = "latest"` where cloud-run-deploy says never to; the Playwright config reuses whatever answers on :8080 where docker and review-standards say never to; the form sample parses a `{field: message}` error body that json-api-style and the Go adapter never produce; and `bash` still routes to Python, which the plan struck (D16, TS23). Group 1 (pressure language, scaffolds, fossils) is clean: no emphasis stacks, no "think step by step", no prefill, no anti-formatting rules. The prohibitions found (never `panic`, never `-it`, never a secret in `--set-env-vars`, never `:latest`) each trace to a rule the harness states with a reason, so they stay (keep-list 5).

Counts per group: Group 1 — 0; Group 2 — 9 findings (+ 8 flags); Group 3 — 0 (no tool descriptions; the skill `description` fields are routing text, assessed below); Group 4 — not applicable.

Description and routing checks (the dispatch's questions): every `description` is a trigger (when to use) plus a Do-NOT line; golang's and json-api-style's also list their contents (flag 6). "When NOT to use" lines exist and route to the right sibling in terraform, golang, json-api-style, docker, gcloud-setup, gke-basics, bash and mermaid; frontend, typescript and cloud-run-deploy carry the routing in the description only (adequate). Only `bash` still routes to `python` (finding 3); nothing routes to `anti-ui-slop`.

Conflicts the dispatch named, result: frontend build order and output path versus docker's build stage — agree on `vite build && tsc -b` and `dist/client`, but docker leaves the copy destination unstated (finding 7). docker base image versus golang's version facts — agree (`golang:<go.mod line>`, `wg.Go` 1.25, `cloudsqlconn` 1.26). cloud-run-deploy flags versus terraform — secret version conflict (finding 2); probe, ingress, VPC egress, concurrency and Cloud SQL agree. gcloud-setup install lines versus `context/dev-setup.md` — agree (the working-tree dev-setup points the standing-permission record at `docs/deploy/local.md` § Operator steps, as docker and gcloud-setup do; the copy in this session's loaded context says `.claude/settings.json` and is stale until the next install). json-api-style field naming versus golang tags — agree (lowerCamelCase); the error shape does not (finding 6). typescript tsconfig versus frontend — agree. bash § Installed commands versus the global context pointer — agree (`context/AGENTS.md:58`).

## Findings

Group 2, repository-grounded; ordered by confidence.

1. **[high] `skills/docker/references/health-checks.md:28` versus `skills/golang/references/service.md:103`**
   - Evidence: docker: `if err != nil || resp.StatusCode != http.StatusOK { os.Exit(1) }`; golang: `mux.Handle("GET /healthz", … w.WriteHeader(http.StatusNoContent))`.
   - Pattern: Group 2, instruction files that contradict each other.
   - Why: an agent that builds the service from the golang reference and adds the docker `-healthcheck` flag gets a container whose `HEALTHCHECK` and compose `healthcheck` report `unhealthy` forever (204 is not 200). Cloud Run's HTTP probe accepts 204, so the failure shows only locally.
   - Confidence: high (both texts quoted). Action: `rewrite` the docker check to any 2xx (hunk 1); the golang side is correct as written.

2. **[high] `skills/terraform/references/resources-gcp.md:157` and `skills/terraform/assets/example/main.tf:69` versus `skills/cloud-run-deploy/SKILL.md:101`**
   - Evidence: terraform `secret_key_ref { … version = "latest" }`; cloud-run-deploy: "Pin the version (`:3`), not `:latest`, for env vars: a rotated secret then ships with a deploy, not under a running revision." D23 records "`--set-secrets` with a pinned version" as the decision.
   - Pattern: Group 2, same setting with two defaults.
   - Why: the Terraform sample is the one an agent copies for the standing infrastructure; it reads `latest`, so a `secret_data_wo_version` bump rotates the secret under a running revision, the failure the deploy skill's rule exists to prevent.
   - Confidence: high. Action: `rewrite` to the pinned number the version resource exports (hunk 2). If the operator prefers `latest`, the direction flips and cloud-run-deploy changes instead, which D23 says it should not.

3. **[high] `skills/bash/SKILL.md:21` and `:27`**
   - Evidence: "Application logic that belongs in Go, Python, or TypeScript"; "Switch to Go or Python when:".
   - Pattern: Group 2, routing to a struck skill. Plan `2026-09-30-audit-fix.md` D16 (python struck), line 162 (TS4/TS23: "the `(use python)` routing phrases in other skills' Do-NOT lines deleted") and line 61 (`grep -rn python\|anti-ui-slop` over skills).
   - Why: `skills/python` is deleted in this tree; an agent told to switch a long script to Python has no skill behind it and the stack here is Go and TypeScript.
   - Confidence: high. Action: `rewrite` both to "Go or TypeScript" (hunk 3). `skills/init-project/SKILL.md:116` uses `python3` as a documented prerequisite (D13), which is a tool, not a routing, and stays.

4. **[high] `skills/typescript/references/testing.md:156`**
   - Evidence: `reuseExistingServer: !process.env['CI'],` on a `webServer` for `http://127.0.0.1:8080/healthz`; `skills/docker/SKILL.md:52` ("a reply from something else on :8080 is not evidence"); `skills/review-standards/SKILL.md` § 2 ("A self-hosting suite refuses to start when its port is already bound … probes the process it started rather than 'something answers'"); the harness runs no CI, so `CI` is never set and the option is always `true`.
   - Pattern: Group 2, instruction files that contradict each other.
   - Why: with the option true, Playwright runs the axe gate and the e2e tier against whatever answers `/healthz` on 8080 (a stale container, another project, `vite dev`), and the test reports green for a build it never ran. The file's own rule on the line below says the e2e tier is for the built Go binary.
   - Confidence: high. Action: `rewrite` to `reuseExistingServer: false` (hunk 4); Playwright then fails when the port is taken, which is the wanted behaviour.

5. **[high] `skills/frontend/references/component-patterns.md:192` and `:206` versus `skills/json-api-style/SKILL.md:31` and `skills/golang/references/service.md:63`**
   - Evidence: "the client shows what the API returns (`{ field: message }`)" and `errors: (await res.json()) as Record<string, string>   // Go returns { field: message }`; the Go adapter and json-api-style rule 7 emit `{"error": {"code", "status", "message", "details": [{"@type": "…BadRequest", "fieldViolations": [{"field", "description"}]}]}}`.
   - Pattern: Group 2, contradiction between sibling skills (frontend consumes what golang produces).
   - Why: against the Go service the same harness teaches, `errors` is `{error: {...}}`, so `state.errors["email"]` is `undefined` for every real validation failure and the form shows nothing while the request failed. The sample's type cast hides it from `tsc`.
   - Confidence: high. Action: `rewrite` the sample to read `fieldViolations` (hunk 5) and fix the sentence on line 192.

6. **[medium] `skills/golang/references/service.md:40-67` versus `skills/json-api-style/SKILL.md:31` and `references/json-api-conventions.md:115`**
   - Evidence: json-api-style: "an `ErrorInfo` (`reason` … `domain` … `metadata`) in every error response"; the sample `apiError` has `Status`, `Code`, `Msg` and `adapt` writes only `code`, `status`, `message`.
   - Pattern: Group 2, contradiction. Why: the sample presents itself as the service's complete error path and cites json-api-style, but a service built from it violates the rule the sibling states as mandatory.
   - Confidence: medium (the sample is a minimal skeleton; the rule says "every"). Action: `rewrite`: add a `Reason` field to `apiError` and emit the `ErrorInfo` detail from `adapt` (hunk 6), or relax rule 7 to "when the error has a stable reason"; the first keeps the newer, stricter text.

7. **[medium] `skills/frontend/SKILL.md:107-109` versus `skills/docker/SKILL.md:87-90`**
   - Evidence: frontend: "in Docker by a `node:24-alpine` stage whose output the Go stage copies (the `docker` skill has the Dockerfile)"; docker: "add a `node:24-alpine` build stage … and `COPY --from=` the output next to the binary".
   - Pattern: Group 2, dangling cross-reference and an ambiguous path. Why: docker holds three prose lines, no Dockerfile, so frontend's pointer resolves to nothing; "next to the binary" is wrong for `embed.go`, which needs `web/dist/client` to exist in the Go build stage before `go build` (and `.dockerignore` excludes `web/dist`, `references/setup.md:94`), so a model following the prose gets `pattern all:dist/client: no matching files`.
   - Confidence: medium (the placement could be inferred from frontend's layout). Action: `add` the stage to docker (hunk 7); frontend's pointer then holds.

8. **[medium] `skills/docker/references/setup.md:112` versus `skills/terraform/references/resources-gcp.md:107`**
   - Evidence: compose `image: postgres:16`; Terraform `database_version = "POSTGRES_18" # the default for new instances`.
   - Pattern: Group 2, the same setting (the database major version) with two values across the dev-to-prod path. Why: the local database is the integration tier's only Postgres; a migration or query that 16 accepts and 18 rejects (or the reverse) first fails in production.
   - Confidence: medium. Action: `rewrite` to `postgres:18` (hunk 8).

9. **[medium] `skills/frontend/SKILL.md:25` and `skills/frontend/references/theme-system.md:3` versus the same file's script and quick test**
   - Evidence: "Light mode is the default"; the script sets `data-theme` from `prefers-color-scheme` when nothing is stored, the toggle's initial state is `system`, and the quick test says "OS set to dark, `localStorage.theme` absent, refresh → dark".
   - Pattern: Group 2, internal contradiction. Why: an agent reading "light is the default" can hard-code `data-theme="light"` in the root or SSR shell, which makes the OS-following state unreachable; the shipped behaviour is "follow the OS, light when the OS gives no answer".
   - Confidence: medium. Action: `rewrite` (hunk 9).

### Flags (no edit proposed)

1. `skills/cloud-run-deploy/SKILL.md:96`, `:100` — `api-auth-gcp` and `data-gcp` do not exist yet. Provenance: plan D38 ("the pointers stay because Phase 7 follows in the same plan") and D24; keep-list 5 (recorded decision), and `review-worktree-p6-docs.md` already carries it. Risk to track: if Phase 7 does not ship in the same release, two dangling skill names ship; `TestSkillLinksResolve` checks files, not skill names.
2. `skills/frontend/references/animation-recipes.md:129-130` — `text-fg-inverse` on a `bg-overlay` wrapper: in dark theme `--fg-inverse` is `#16191d` over a 70 % black overlay, so the heading is near-invisible; no existing token is light in both themes (`--color-*: initial` removes `text-white`). The fix needs a new token (`--on-image`), which is a design choice for the operator.
3. `skills/terraform/references/checks.md:5-10` — the table runs `terraform init -backend=false && validate` (row 2) before `plan` (row 3); on a fresh clone `plan` then fails with "Backend initialization required" until a normal `terraform init` runs. The prose below the table names that init but not where it goes. Low: add the `init` as a row before the plan.
4. `skills/golang/references/service.md:16` — "Supported releases are the two newest majors plus the current one": the Go policy supports the current and the previous major; "plus the current one" restates one of the two. Low: wording.
5. Skills deployed into other projects cite repository paths that exist only in the harness checkout (`context/AGENTS.md`, `context/dev-setup.md`: docker/SKILL.md:55 and :109, gcloud-setup/SKILL.md:39, :92, :143, bash/references/testing-linting.md:73). On a deployed machine the content is in the model's context as `AGENTS.md`, so the pointer misleads rather than breaks. Low; a one-sentence convention ("`context/…` names the harness's source; installed, it is your global AGENTS.md") would cover them.
6. `description` fields of `golang` (names both references and their contents), `json-api-style` (lists six AIP numbers) and `mermaid-syntax` ("Use when:" followed by seven near-synonymous phrases) restate the body or enumerate triggers (Group 2, trigger-case enumeration; Group 3 trigger/behaviour split). All are under the 1,024-character limit and carry the routing keywords; trimming to intent categories is a preference, low.
7. History notes in instruction text: `skills/bash/SKILL.md:50`, `:54` ("(Moved from the global context 2026-10-01.)"), `skills/typescript/references/testing.md:3` and the last section ("see the harness's Phase 6 report for the exit status"). The plan (line 112) accepts the "(moved …)" notes as provenance; the Phase 6 report pointer will go stale. Low.
8. `skills/frontend/references/design.md:8` ("award-level composition … premium, current") is the vague kind of design direction; the file does name specific defaults to avoid (hero cards, stat strips, logo clouds, pill soup), which is the form that works (prompt-audit 1e exception). Optional: extend the named list from the styles the first result used. Low, nothing proposed.

Checked and clean (so the next run does not repeat it): `gke-basics` Harness overrides — every claim it makes about the references holds (`gke-golden-path.md:68` "don't ask users to paste project IDs", `gke-cluster-creation.md:132` "ALWAYS ask for `project_id`", the MCP-first preference in `cli-reference.md:8` and `mcp-usage.md:72`, the dangling `skills/gke/assets/…` path at `gke-security.md:117`, the `resources/recommended-ap.yaml` reference in the asset, the install roots `.claude/skills` and `.gemini/config/skills` in `plan.go:58`). Terraform versus cloud-run-deploy on probe path, ingress, public access, VPC egress and `/cloudsql`: consistent. mermaid: `%%{init}%%` appears only in the one sentence that tells a reader to recognise it; the changed examples all use frontmatter `config:`.

## Proposed diff

One hunk per finding, high or medium confidence; nothing applied.

```diff
# 1. skills/docker/references/health-checks.md
-    if err != nil || resp.StatusCode != http.StatusOK {
+    if err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
```

```diff
# 2a. skills/terraform/references/resources-gcp.md (Cloud Run v2 service, env block)
-            version = "latest"
+            version = google_secret_manager_secret_version.db_url.version # the pinned number: a rotation is a new revision
# 2b. skills/terraform/assets/example/main.tf:69
-            version = "latest"
+            version = google_secret_manager_secret_version.app.version
```

```diff
# 3. skills/bash/SKILL.md
-- Application logic that belongs in Go, Python, or TypeScript
+- Application logic that belongs in Go or TypeScript
-Shell is the right tool for glue scripts and automation. Switch to Go or Python when:
+Shell is the right tool for glue scripts and automation. Switch to Go when:
```

```diff
# 4. skills/typescript/references/testing.md:156
-    reuseExistingServer: !process.env['CI'],
+    reuseExistingServer: false, // a reply from something else on :8080 is not the build under test; Playwright fails when the port is taken
```

```diff
# 5a. skills/frontend/references/component-patterns.md:192 (Forms intro)
-the client shows what the API returns (`{ field: message }`).
+the client shows what the API returns: the `error.details[].fieldViolations[]` of the `json-api-style` error shape, mapped to field → message.
# 5b. skills/frontend/references/component-patterns.md:206
-  return { ok: false, errors: (await res.json()) as Record<string, string> }   // Go returns { field: message }
+  const body = (await res.json()) as { error?: { message?: string; details?: { fieldViolations?: { field: string; description: string }[] }[] } }
+  const errors: Record<string, string> = {}
+  for (const d of body.error?.details ?? []) for (const v of d.fieldViolations ?? []) errors[v.field] = v.description
+  return { ok: false, errors }   // json-api-style: BadRequest.fieldViolations
```

```diff
# 6. skills/golang/references/service.md (apiError and adapt)
 type apiError struct {
 	Status int    `json:"-"`
 	Code   string `json:"status"` // the google.rpc.Code name, see the json-api-style skill
 	Msg    string `json:"message"`
+	Reason string `json:"-"`      // UPPER_SNAKE_CASE, goes into the ErrorInfo detail
 }
@@ adapt, the INTERNAL fallback and the writeJSON call
-			ae = &apiError{Status: http.StatusInternalServerError, Code: "INTERNAL", Msg: "internal error"}
+			ae = &apiError{Status: http.StatusInternalServerError, Code: "INTERNAL", Msg: "internal error", Reason: "INTERNAL_ERROR"}
 		}
 		writeJSON(w, ae.Status, map[string]any{"error": map[string]any{
 			"code": ae.Status, "status": ae.Code, "message": ae.Msg,
+			"details": []any{map[string]any{
+				"@type": "type.googleapis.com/google.rpc.ErrorInfo", "reason": ae.Reason, "domain": "app.example.com",
+			}},
 		}})
(and a `Reason` on each `apiError` literal in `readJSON` and `hello`, e.g. "MALFORMED_JSON", "NAME_REQUIRED"; "domain" is the service's name, a constant in `main`)
```

```diff
# 7. skills/docker/SKILL.md, replacing the paragraph after the Key rules list
-If the Go binary also serves the `web/` frontend, add a `node:24-alpine` build stage
-(the only Node image here) that runs `npm ci && npm run build` in `web/` and
-`COPY --from=` the output next to the binary; the runtime image stays distroless.
-`web/dist/client` is built in the `node:24-alpine` stage, never committed.
+If the Go binary also serves the `web/` frontend (`embed.go` embeds `dist/client`), add a
+`node:24-alpine` stage (the only Node image here) and copy its output to `web/dist/client`
+in the Go stage after `COPY . .` and before `go build`; the runtime image stays distroless.
+`web/dist` is in `.dockerignore` and never committed, so this stage is the only source of it:
+
+```dockerfile
+FROM node:24-alpine AS web
+WORKDIR /web
+COPY web/package.json web/package-lock.json ./
+RUN npm ci
+COPY web/ .
+RUN npm run build
+# in the golang stage, after COPY . .
+COPY --from=web /web/dist/client ./web/dist/client
+```
```

```diff
# 8. skills/docker/references/setup.md:112
-    image: postgres:16
+    image: postgres:18   # the major version Cloud SQL runs here (terraform skill: POSTGRES_18)
```

```diff
# 9a. skills/frontend/SKILL.md:25
-- **Theme:** Every frontend ships with light/dark/system theme support. Light mode is the default. Tokens are
+- **Theme:** Every frontend ships with light/dark/system theme support. With no stored choice the page follows the OS and falls back to light. Tokens are
# 9b. skills/frontend/references/theme-system.md:3
-Every frontend ships with dark/light theme support. **Light mode is the default**; "system" (follow the OS) is the third state of the toggle.
+Every frontend ships with dark/light theme support. **With no stored choice the page follows the OS and falls back to light**; "system" (follow the OS) is the third state of the toggle and the initial one.
```

## Criteria
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

## Not examined
- `skills/mermaid-syntax/references/{class-state-complete,flowchart-complete,sequence-complete,other-diagrams}.md`: only the Phase 6 changed hunks (git diff), the `%%{init` and `packet-beta` greps and the first 30 lines of `flowchart-complete.md` were read; the unchanged syntax tables were not read line by line, and no Mermaid syntax claim was checked against the renderer.
- `skills/bash/references/{bash-fundamentals,makefile-patterns,testing-linting}.md`: Phase 6 hunks read in full and the rest searched (python, CI, brew, apt, github); not read end to end. One tension left unreported as a finding: `bash-fundamentals.md` shows `if [ -f "$file" ]` as the "Good" form while `SKILL.md` rule 3 says `[[ ]]`.
- `skills/terraform/assets/example/{main,outputs,variables,versions}.tf`: only `main.tf` lines 38 and 60-75 and the grep; not compared against `resources-gcp.md` block by block.
- No Go or HCL sample was compiled, planned or run (the plan's own verify row says the Go samples were compiled in scratch; not re-run here); no `gcloud`, `terraform`, `npm` or `docker` command was run. External facts (TypeScript 7.0, Vitest 5, Motion 13.4, Terraform 1.16, provider 8.x, Go 1.27 APIs, tflint ruleset 0.40.0, golangci-lint v2.14.0, Cloud Run flag behaviour) were taken as written, not verified; their verification dates in the files were not re-checked.
- Behavioural probes (prompt-audit Step 7) were not run: each finding is a repository contradiction, checked by reading both locations.
- `git blame` was not run on the changed lines (the files are uncommitted additions or rewrites in this working tree); provenance came from the plan's D-rows (D16, D18, D20-D26, D37, D38, TS4/TS23) and the earlier review `2026-10-01-review-worktree-p6-docs.md`. `INDEX.md` and `DECISIONS.md` were not opened, the plan D-rows being the recorded decisions for these files.

## Outcome
- F1 [high] docker `-healthcheck` wants 200, golang `/healthz` answers 204 — fixed at `skills/docker/references/health-checks.md` (any 2xx)
- F2 [high] terraform `version = "latest"` vs the pin rule; Playwright `reuseExistingServer` — fixed (both terraform sites pinned to the managed version; `reuseExistingServer: false` in `skills/typescript/references/testing.md`)
- F3 [high] form sample parses `{field: message}`; bash routes to Python — fixed at `skills/frontend/references/component-patterns.md` (an `ApiError` type and `fieldErrors()` over `fieldViolations`) and `skills/bash/SKILL.md:21,27`
- [medium] Go error body lacks `ErrorInfo` — fixed at `skills/golang/references/service.md` (`Reason` on `apiError`, `errorDomain`, the `details[]` entry in `adapt`; `go vet` clean)
- [medium] docker has no node stage; "next to the binary" breaks `go:embed` — fixed (three-stage Dockerfile)
- [medium] compose `postgres:16` vs Terraform `POSTGRES_18` — fixed at `skills/docker/references/setup.md` (`postgres:18`)
- [medium] "light is the default" vs the OS-following script — fixed at `skills/frontend/SKILL.md:25` and `references/theme-system.md:3` (system is the default on first load, resolving to light without an OS preference)
- flags (8): left — `api-auth-gcp`/`data-gcp` per D38; the rest are Phase 9 docs-lens items or judged readable

## Round 2

Range reviewed (round 2): working tree, prompts lens; the instruction files Phase 6 rewrote, read whole where they carry a round-1 fix and by targeted read plus repository-wide grep elsewhere: `skills/{terraform,golang,json-api-style,frontend,typescript,cloud-run-deploy,docker,gcloud-setup,bash,mermaid-syntax}/SKILL.md` and their `references/*.md`, `skills/gke-basics/SKILL.md` (Harness overrides paragraph only). Step 4 Group 2 re-run over the scope; Groups 1, 3 and 4 confirmed unchanged (below).

### Status: COMPLETE

### Verdict: REQUEST CHANGES

Round 1's seven fixes are each correct as written: F1 (docker self-check `StatusCode/100 != 2`, golang's 204 passes), F2 (both Terraform sites pin `google_secret_manager_secret_version.*.version`, `reuseExistingServer: false`), F3 (the form sample parses the AIP-193 body; bash says "Go or TypeScript"), the Go `ErrorInfo` detail, the three-stage Dockerfile (paths agree with `web/embed.go`'s `dist/client` and `web/package-lock.json`), `postgres:18`, and "system is the default" (the script, the toggle's initial state, the priority table and `typescript/references/testing.md`'s `theme: system` test all agree). The fixes leave two defects, one of them a new conflict: the Go sample now emits `ErrorInfo` but never the `BadRequest.fieldViolations` that json-api-style rule 7 requires for validation and that the corrected form sample reads, so against the Go service the harness teaches no field error ever reaches the form (finding 2); and the form's fallback `{ _: message }` is never rendered, so any error without a field violation shows nothing (finding 3). A third restatement of the Dockerfile tag rule still says the shorter, wrong wording (finding 1). Everything else checked agrees.

Step 0 assumptions as round 1 (targets Claude and Gemini; Group 2 findings are repository-grounded).

### Round 2 findings (Group 2, repository-grounded)

1. **[medium] `skills/golang/SKILL.md:56` and `skills/gke-basics/SKILL.md:23` versus `skills/docker/SKILL.md:68` and `:86`**
   - Evidence: golang: "the image's base is `golang:<the go.mod go line>` (docker skill)"; gke-basics: "use `golang:<the go.mod go line>` (the `docker` skill)"; docker rule 1: "`golang:<major.minor of the go.mod go line>` (`go 1.27.1` → `golang:1.27` … a three-part tag such as `golang:1.24.4` freezes a patch with known advisories)".
   - Pattern: Group 2, one rule restated in two wordings, the shorter one giving the wrong tag. The docs-lens review (`2026-10-01-review-worktree-p6-docs.md`, low) and the security review (`2026-10-01-security-worktree-p6.md`, medium) both name this; the fix reached docker only.
   - Why: go.mod's line is usually three-part since Go 1.21 (this repo's own is `go 1.24.4`); a model editing `gke-app-onboarding.md`'s Dockerfile "per the override" writes `golang:1.27.1`, the frozen patch docker forbids.
   - Confidence: high. Action: `rewrite` both (hunk 1).

2. **[medium] `skills/golang/references/service.md:21-33,101` versus `skills/json-api-style/SKILL.md:31` and `skills/frontend/references/component-patterns.md:192-206`**
   - Evidence: rule 7: "`BadRequest.fieldViolations[]` for validation"; the Go sample's validation failure is `&apiError{… Reason: "NAME_REQUIRED", Msg: "name is required"}` and `adapt` writes only the `ErrorInfo` entry; the form sample reads `fieldViolations` and nothing else (grep: `BadRequest` appears in no golang file).
   - Pattern: Group 2, producer and consumer sample disagree (round 1's finding 5 and 6, one hop further: the fixes made each side match rule 7 but not each other).
   - Why: a service built from `service.md` and a form built from `component-patterns.md` never connect: `errors` is always `{ _: message }`, and the form renders `errors["email"]` only. The sample is a skeleton, but `apiError` has no field to carry a violation, so a model extending it has to invent the shape.
   - Confidence: medium (rule 7 says "for validation", and the sample's only validation path is the one that omits it). Action: `rewrite` (hunk 2).

3. **[medium] `skills/frontend/references/component-patterns.md:200-206,219`**
   - Evidence: `fieldErrors` returns `{ _: body.error.message }` with the comment "no field violation: show the message once"; the markup reads `err("email")` and no other key; nothing renders `state.errors._`.
   - Pattern: Group 2, a sample whose stated behaviour its own markup does not implement. Why: a 401, 429, 500 or any 400 without a violation (everything the Go sample emits today, finding 2) leaves the form silently idle; the submit button re-enables and nothing is said. Confidence: high (both lines quoted). Action: `rewrite` (hunk 3).

4. **[low] `skills/frontend/references/theme-system.md:324` (and `:104`)**
   - Evidence: "Page loads in light mode on a light OS with nothing stored (default)". The rest of the test list matches the corrected claim (the cycle order equals the `next` map, the dark-OS bullet matches the script). The `(default)` beside "light mode" is the old claim's wording: with nothing stored the toggle reads "System" (the typescript test asserts that), the page merely resolves to light. `:104`'s "Light: VS Code Light Modern (default)" labels the `:root` palette and is correct.
   - Confidence: medium; Action: `rewrite` (hunk 4).

5. **[low] `skills/terraform/references/resources-gcp.md:175` and `:99`, `skills/terraform/SKILL.md:61`**
   - Evidence: the comment says "add the version resource to depends_on", but the attribute reference now orders it (`main.tf:83` keeps a `depends_on` with its own reason, harmless); and the alternative "create the secret in HCL and add versions with `gcloud secrets versions add`" leaves no resource whose `.version` can be pinned, while cloud-run-deploy's rule is "Pin the version (`:3`)".
   - Why: the comment is a leftover of the `latest` form; the alternative's reader has no stated pin. Cloud Run's pin rule and Terraform's now agree for the managed-version path; the gcloud path states nothing. Confidence: medium. Action: `rewrite` (hunk 5).

### Round 2 flags (no edit proposed)

- `skills/golang/references/testing.md:158` — the test's fake handler writes `{"error":{"code":400,"status":"INVALID_ARGUMENT","message":"name is required"}}` through `http.Error` (so `text/plain`, no `details`), beside `service.md`'s adapter that writes `application/json` with an `ErrorInfo`. It is a self-contained fake whose test asserts only the status; low.
- `skills/docker/references/health-checks.md:11-14` answers 200 "ok", `golang/references/service.md:109` answers 204. The 2xx self-check makes them agree; two samples of one endpoint is working redundancy (keep-list 8).
- Unchanged from round 1 and left per the Outcome: flags 1-8 (round 1).

### Checked and agreeing (round 2)

- docker Dockerfile versus frontend: `web/package.json` + `package-lock.json` copied then `npm ci`; `npm run build` = `vite build && tsc -b` → `/src/web/dist/client`; copied to `./web/dist/client` after `COPY . .` and before `go build`, where `web/embed.go` (`go:embed all:dist/client`, SKILL.md:37) reads it; `.dockerignore` carries `web/node_modules` and `web/dist` (`setup.md:14,16`), so the stage is the only source. `node:24-alpine` satisfies `engines node >=22.12`. frontend SKILL.md:107-109's pointer "the `docker` skill has the Dockerfile" now resolves.
- docker rule 1 versus its own `FROM` line (`golang:<major.minor of the go.mod go line>`), troubleshooting row (`golang:<major.minor>`) and `cloudsqlconn` 1.26 (golang `service.md:16,293`): agree.
- cloud-run-deploy `--set-secrets KEY=NAME:VERSION` pinned (`:101`, `:181`, `:192`) versus Terraform `version = …_secret_version.*.version`, `secret_data_wo_version` bump to rotate (`resources-gcp.md:98`): the same rule on the managed path. Not verified against the provider that `google_secret_manager_secret_version.version` is the plain integer `secret_key_ref.version` accepts.
- json-api-style rule 7 versus the Go adapter: `code` (HTTP status int), `status` (Code name), `message`, `details[0]` = `ErrorInfo` with `reason` (UPPER_SNAKE_CASE: `INTERNAL`, `MALFORMED_BODY`, `NAME_REQUIRED`) and `domain` (`errorDomain`); `metadata` omitted (listed in the rule as part of `ErrorInfo`, optional in AIP-193; the rule's wording "(… `metadata` map)" does not say optional; not raised). The `BadRequest` part is finding 2.
- No frontend reference other than the lines in finding 4 says light is the default (grepped `skills`, `context`, `README.md`); `reuseExistingServer` appears once, `false`; `postgres:` and `POSTGRES_` agree (18); no `python` routing remains in the scope (only a `fd "\.py$"` example).

### Group confirmations

Group 1 (pressure language, scaffolds, fossils): clean in the changed hunks (the fixes added no emphasis, no scaffolds). Group 3 (tool descriptions): not applicable. Group 4: not applicable.

### Proposed diff (round 2)

```diff
# 1a. skills/golang/SKILL.md:56
-the image's base is `golang:<the go.mod go line>` (docker skill).
+the image's base is `golang:<major.minor of the go.mod go line>` (docker skill).
# 1b. skills/gke-basics/SKILL.md:23
-use `golang:<the go.mod go line>` (the `docker` skill)
+use `golang:<major.minor of the go.mod go line>` (the `docker` skill)
```

```diff
# 2. skills/golang/references/service.md
 type apiError struct {
 	Status int    `json:"-"`
 	Code   string `json:"status"`
 	Reason string `json:"-"`
+	Field  string `json:"-"` // a validation failure names the request field: it becomes the BadRequest fieldViolation
 	Msg    string `json:"message"`
 }
@@ adapt, replace the details literal
-		writeJSON(w, ae.Status, map[string]any{"error": map[string]any{
-			"code": ae.Status, "status": ae.Code, "message": ae.Msg,
-			"details": []any{map[string]any{ // json-api-style rule 7: every error carries an ErrorInfo
-				"@type": "type.googleapis.com/google.rpc.ErrorInfo", "reason": ae.Reason, "domain": errorDomain,
-			}},
-		}})
+		details := []any{map[string]any{ // json-api-style rule 7: every error carries an ErrorInfo
+			"@type": "type.googleapis.com/google.rpc.ErrorInfo", "reason": ae.Reason, "domain": errorDomain,
+		}}
+		if ae.Field != "" { // ... and a validation error a BadRequest, at most once each
+			details = append(details, map[string]any{
+				"@type":           "type.googleapis.com/google.rpc.BadRequest",
+				"fieldViolations": []any{map[string]any{"field": ae.Field, "description": ae.Msg}},
+			})
+		}
+		writeJSON(w, ae.Status, map[string]any{"error": map[string]any{
+			"code": ae.Status, "status": ae.Code, "message": ae.Msg, "details": details,
+		}})
@@ hello
-		return &apiError{Status: http.StatusBadRequest, Code: "INVALID_ARGUMENT", Reason: "NAME_REQUIRED", Msg: "name is required"}
+		return &apiError{Status: http.StatusBadRequest, Code: "INVALID_ARGUMENT", Reason: "NAME_REQUIRED", Field: "name", Msg: "name is required"}
```

```diff
# 3. skills/frontend/references/component-patterns.md, in SignupForm, above the submit button
+      {!state.ok && state.errors._ && <p role="alert" className="text-sm text-fg">{state.errors._}</p>}
       <button type="submit" disabled={pending}
```

```diff
# 4. skills/frontend/references/theme-system.md:324
-- Page loads in light mode on a light OS with nothing stored (default)
+- Nothing stored, light OS: the page loads light and the toggle reads "System"
```

```diff
# 5a. skills/terraform/references/resources-gcp.md:175
-            version = google_secret_manager_secret_version.db_url.version # pinned to the managed version; add the version resource to depends_on
+            version = google_secret_manager_secret_version.db_url.version # pinned to the managed version; the reference orders it before the service
# 5b. skills/terraform/references/resources-gcp.md:99
-- The alternative that keeps values out of HCL entirely: create only the secret here and add versions with `gcloud secrets versions add NAME --data-file=-`.
+- The alternative that keeps values out of HCL entirely: create only the secret here and add versions with `gcloud secrets versions add NAME --data-file=-`; the service then pins the number that command prints (`version = "3"`), never `latest` (`cloud-run-deploy` pins too).
```

### Round 2 criteria
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

### Round 2 not examined
- No Go, HCL or TSX sample was compiled or run; hunk 2 and 3 are untested text. The provider's exported `version` attribute and Cloud Run's acceptance of it in `secret_key_ref` were taken from memory of the provider docs, not checked.
- `skills/mermaid-syntax`, `skills/bash` references, `skills/gcloud-setup` and `skills/cloud-run-deploy` beyond the secrets, probe and Dockerfile-adjacent lines: not re-read in round 2 (nothing in the round-1 fix list touched them except `bash/SKILL.md:21,27`, read).
- `skills/golang/references/service.md` and `skills/frontend/references/theme-system.md` were read where they carry the fixes, not end to end; `animation-recipes.md`, `design.md`, `accessibility.md` only by grep for the changed claims.
- Behavioural probes (prompt-audit Step 7): none; every finding is a repository contradiction checked by reading both locations.

## Outcome (round 2)
- [medium] `golang/SKILL.md:56`, `gke-basics/SKILL.md:23` still `golang:<the go.mod go line>` — fixed (`<major.minor of the go.mod go line>`)
- [medium] the Go sample emits no `BadRequest.fieldViolations` — fixed at `skills/golang/references/service.md` (`apiError.Field`; `adapt` appends the `BadRequest` detail when set; `NAME_REQUIRED` names `name`); block re-extracted, `go vet` clean
- [medium] the `_` fallback never rendered — fixed (see the code lens Outcome)
- [low] `theme-system.md:324` "(default)" — fixed ("system mode resolving to light")
- [low] `resources-gcp.md:175` comment — fixed (the reference orders the two resources)

## Round 3

Range reviewed (round 3): working tree, prompts lens; `skills/{terraform,golang,json-api-style,frontend,typescript,cloud-run-deploy,docker,gcloud-setup,bash,mermaid-syntax}/SKILL.md` and their `references/*.md`, `skills/gke-basics/SKILL.md` (Harness overrides paragraph only). Round-2 fixes verified by reading both ends of each; Step 4 Group 2 re-run over the scope.

### Status: COMPLETE

### Verdict: APPROVE WITH NITS

Every round-2 fix is correct as written, and the six cross-file questions the dispatch named all agree: docker's `FROM` for distroless (digest) matches rule 1 and says nothing `health-checks.md` contradicts; the allowlist `.dockerignore` admits what the three-stage Dockerfile copies (`go.mod`, `go.sum`, `cmd/`, `internal/`, `web/`; `!web/` re-admits `web/embed.go`, and the sample's migrations live in `internal/db/migrations`, covered by `!internal/`); cloud-run-deploy's "Terraform owns it" matches terraform item 6; json-api-style rule 7 matches `adapt`; `theme-system.md`'s quick test matches `frontend/SKILL.md:25` and `typescript/references/testing.md`. One conflict is new from a round-2 fix (finding 1, `--immutable-tags` against the `:latest` alias), one round-2 hunk was left unapplied (finding 2), and one claim added in round 2 is true only for adapter-made bodies (finding 3). None is high.

### Round 3 findings (Group 2, repository-grounded)

1. **[medium / confidence: medium] NEW from a round-2 fix: `skills/cloud-run-deploy/SKILL.md:54` versus `:71` (and `terraform/references/resources-gcp.md:62`, `assets/example/main.tf:21`)**
   - Evidence: `:54` "add :latest only as a convenience alias"; `:71` the repository is created with `--immutable-tags`; Terraform sets `immutable_tags = true`. The global context's rule is "`latest` only as an alias".
   - Why: on an immutable-tags repository a tag can be written once; the second `docker push …:latest` is refused, so an agent following `:54` and `:71` gets a failed push on the second release (or retags nothing and believes it did). I did not run it against Artifact Registry; the behaviour is from the documented meaning of immutable tags.
   - Action: `rewrite` `:54` to "`:latest` alias only on a repository without immutable tags; with them the SHA tag is the only tag" (hunk 1).

2. **[low / confidence: high] Left from round 2 (finding 5b, not applied): `skills/terraform/references/resources-gcp.md:99`**
   - Evidence: "create only the secret here and add versions with `gcloud secrets versions add …`" — the Cloud Run block at `:175` pins `google_secret_manager_secret_version.db_url.version`, a resource that alternative does not create, and states no pin; cloud-run-deploy pins `:3`.
   - Why: a reader taking the alternative has a service block that does not apply and no stated version rule. Action: `rewrite` (hunk 2). Round 2's Outcome fixed only the `:175` comment.

3. **[low / confidence: medium] Claim added by a round-2 fix: `skills/frontend/references/component-patterns.md:206`**
   - Evidence: the comment says the `_` line shows "a 401, 429, 500 or a bare 400"; the action does `(await res.json()) as ApiError` and `fieldErrors` reads `body.error.message`.
   - Why: the claim holds for bodies the Go adapter writes. A Go mux 404/405 ("404 page not found", `text/plain`), a Cloud Run 503 or an IAM 403 page is not JSON: `res.json()` rejects inside the action and the form errors out instead of showing the message. Action: narrow the comment to "an adapter-made error" or add `.catch(() => null)` and a fallback `{ _: res.statusText }` (hunk 3).

### Round 3 flags (no edit proposed)

- `skills/gke-basics/references/gke-app-onboarding.md:33` uses `FROM gcr.io/distroless/static:nonroot` (a floating tag); the Harness overrides paragraph (`SKILL.md:23`) corrects only the `golang:` line, so docker rule 1's digest pin and the vendored Dockerfile disagree. Vendored text, so the override is the place; low.
- `skills/docker/references/setup.md:27-33` — "a secret in any shape … stays out without being named" holds for the repository root; inside re-admitted `cmd/`, `internal/` and `web/` only the `.env`/`.npmrc` shapes (and `web/node_modules`, `web/dist`) are excluded, so a `key.json` under `internal/` enters the context. `:28` also suggests `migrations/` or `pkg/`, but the golang sample embeds `internal/db/migrations` (already admitted) and `pkg/` is for code another module imports. The frontend skill allows `ui/` for `web/`; the Dockerfile and `!web/` hardcode `web/` (the `:28` sentence covers the rename). Low, wording.
- `skills/json-api-style/SKILL.md:31` lists `metadata` in the mandatory `ErrorInfo`; `adapt` omits it (optional in AIP-193; round 2 left it). Unchanged.
- Round 1 flags 1-8 and round 2 flags: unchanged, left per the earlier Outcomes.

### Verified fixes (round 2)

- `golang/SKILL.md:56`, `gke-basics/SKILL.md:23`: both read `<major.minor of the go.mod go line>`, matching docker rule 1 and the Dockerfile `FROM`.
- `service.md`: `apiError.Field`; `adapt` appends a `BadRequest` detail with `fieldViolations[{field, description}]` once, `ErrorInfo` once; `hello` sets `Field: "name"`; the `[]map[string]string` literal is type-correct. The form's `fieldErrors` reads the same keys.
- `component-patterns.md`: `err("_")` renders the fallback; `export function SignupForm`. `theme-system.md:324` "system mode resolving to light" agrees with `frontend/SKILL.md:25` and the typescript test (`theme: system` toggle, light OS mocked).
- Terraform: `resources-gcp.md:175` comment reworded; the public/IAM-gated section's live content is by pointer (IAM-gated), `allUsers` a comment with its condition; `main.tf` uses `invoker_iam_disabled = var.public`, which matches item 6.
- docker: Dockerfile `FROM gcr.io/distroless/static-debian12@sha256:<digest>` equals rule 1 and the `distroless/static` mentions elsewhere; `health-checks.md` names no base image and uses `/app`, `COPY --from=build /app /app` agrees; `postgres:18` with `pgdata:/var/lib/postgresql` agrees with Terraform `POSTGRES_18`.
- cloud-run-deploy: `:95` (Terraform owns the setting, flag never passed) = terraform item 6; `:143` bare curl must 403 = the `:190` troubleshooting row; IAP cell (non-empty audience, accessor on the service, explicit issuer) is self-consistent.

### Group confirmations

Group 1: clean in the round-2 hunks (no emphasis, scaffolds or fossils added). Group 3: not applicable (no tool descriptions). Group 4: not applicable.

### Proposed diff (round 3)

```diff
# 1. skills/cloud-run-deploy/SKILL.md:54
-TAG="$(git rev-parse --short HEAD)"          # immutable tag; add :latest only as a convenience alias
+TAG="$(git rev-parse --short HEAD)"          # the only tag on an immutable-tags repository (:latest can be pushed once there, never again)
```

```diff
# 2. skills/terraform/references/resources-gcp.md:99
-- The alternative that keeps values out of HCL entirely: create only the secret here and add versions with `gcloud secrets versions add NAME --data-file=-`.
+- The alternative that keeps values out of HCL entirely: create only the secret here and add versions with `gcloud secrets versions add NAME --data-file=-`; the service block then has no version resource to reference, so it pins the number that command prints (`version = "3"`), never `latest`.
```

```diff
# 3. skills/frontend/references/component-patterns.md, in signup
-  return { ok: false, errors: fieldErrors((await res.json()) as ApiError) }   // the Go adapter's AIP-193 body
+  const body = (await res.json().catch(() => null)) as ApiError | null        // a proxy or mux error page is not JSON
+  return { ok: false, errors: body?.error ? fieldErrors(body) : { _: res.statusText || "Request failed" } }
```

### Round 3 criteria
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

### Round 3 not examined
- No Go, HCL, TSX or Dockerfile sample was compiled or built; the `.dockerignore` semantics (last matching rule wins, `!web/` re-admitting `web/embed.go`) are from Docker's documented behaviour, not a `docker build` run. Artifact Registry's refusal of a second push to an immutable tag was not exercised.
- Not re-read in round 3: `skills/mermaid-syntax`, `skills/bash`, `skills/gcloud-setup` (nothing in the round-2 fix list touched them), `golang/references/testing.md`, `frontend/references/{design,accessibility,animation-recipes}.md`. `terraform/assets/example/*.tf` only by grep for `version`, `invoker`, `immutable`.
- The provider's `google_secret_manager_secret_version.version` shape in `secret_key_ref.version` is still unverified (round 2).

## Outcome (round 3)
(to be filled by the author)

## Outcome (round 3)
- [medium, from a round-2 fix] `:latest` alias vs immutable tags — fixed per D40 (see the code lens Outcome) (review cap reached; fixed under the user's same-turn decision of 2026-10-01 — "Go with your recommendation, drop the latest alias" — and re-checked by Phase 9's lenses, not by a fourth round here)
- [low] `resources-gcp.md` gcloud-versions alternative states no pin — fixed (the service pins the version number the command prints)
- [low] the `_` comment overclaims on non-JSON replies — fixed (the guard and the comment)
- flags: the vendored gke-basics Dockerfile's floating distroless tag is upstream's (fork rule: SKILL.md only); the allowlist paragraph now says what holds inside re-admitted directories
