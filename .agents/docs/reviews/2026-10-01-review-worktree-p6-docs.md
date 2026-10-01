---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p6-docs   (2026-10-01)

Range reviewed: working tree (git status --short --untracked-files=all), Phase 6 of the audit-fix plan: `skills/{terraform,golang,json-api-style,frontend,typescript,cloud-run-deploy,docker,gcloud-setup,bash,mermaid-syntax}/**`, `skills/gke-basics/SKILL.md`, `skills/VENDORED.md`, `AGENTS.md`, `skills/init-project/assets/agents-tree/docs/deploy/gcp.md`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/plans/2026-09-30-audit-fix.md`; deleted `skills/python/`, `skills/anti-ui-slop/`, seven `skills/golang/references/golang-*.md`, six `skills/terraform/references/*.md`; untracked `internal/harness/skill_links_test.go` (read whole). Lens: docs (`--docs`), report only.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

No stale reference to a struck skill or a deleted file survives in the living docs: the restatement grep for `anti-ui-slop`, `python`, `samber`, `golang-*`, `terraform-skill`, `Anton Babenko`, the old flags and the old Tailwind/Zod/mermaid wording hit only dated plans, bug records and the new VENDORED note. README, context, specs, the testing guide and `deploy/*.md` carry no skill count or skill list that went stale. The skill-matrix rows for `python` and `anti-ui-slop` are gone and every remaining technical skill still has a row. The worst problems are two statements that point at something that does not exist yet (skills Phase 7 will write; a gcp.md slot the template lacks). `TestSkillLinksResolve`, `TestVendored*` and `TestNoTmp*` pass (`go test ./internal/harness -run 'TestSkillLinks|TestVendored|TestNoTmp'` ok).

## Findings
- **[medium / confidence: high]** `skills/cloud-run-deploy/SKILL.md:96` and `:100` — point at the `api-auth-gcp` and `data-gcp` skills ("`api-auth-gcp` skill", "Connection code: `data-gcp` skill"); neither exists under `skills/` (Phase 7 writes them; plan D38 records the pointers as deliberate). Until Phase 7 lands in the same release, an agent told to read them finds no such skill, and `TestSkillLinksResolve` cannot catch it because it checks file links, not skill names. Fix: either ship Phase 6 and 7 together, or reword the two cells to say the skill is not yet written (or drop the pointers) until Phase 7 commits.
- **[low / confidence: high]** `skills/terraform/SKILL.md:17` vs `skills/init-project/assets/agents-tree/docs/deploy/gcp.md:16-23,59-61` — the skill says "the state bucket's name is recorded there too" (in `deploy/gcp.md`), but the template has no row for it: § Recorded defaults lists six settings, the new § Rollout and network table has no state-bucket row, and § Infrastructure's comment names only "the Terraform module path and how to plan/apply". A project scaffolded from the template has no place to record it; the claim has no home. Fix: add the state bucket (and its location) to the `## Infrastructure` comment or as a row.
- **[low / confidence: high]** `skills/VENDORED.md` (paragraph under "Apache-2.0 skills — forks") — "Two skills that used to sit here are gone: `terraform` … ; a third, `gcloud-auth-verification` …" names two skills, not three; the "Two" and the "a third" do not agree. Fix: "Two skills that used to sit here are gone: `terraform` …; and `gcloud-auth-verification` …" (or drop the count).
- **[low / confidence: high]** `.agents/docs/bugs/INDEX.md:108` (Graduated table, BUG-022) — "Graduated to `skills/golang/references/golang-ecosystem.md` (registration scoped to public libraries)": the file is deleted by this change and the new `golang` skill and its two references (`service.md`, `testing.md`) carry no registration rule (`grep -i regist` finds only the `init()` driver-registration line and a pgx driver comment). The rule's subject (the old reference recommending registration) went with the file, so the rule is moot, but the index row now names a path that does not exist. `dropped: skills/golang/references/golang-ecosystem.md` (BUG-022 guard). Fix: repoint the row (the bugs agent owns it) to the golang skill or mark the row as retired with the file.
- **[low / confidence: medium]** `skills/frontend/SKILL.md` — `dropped: skills/anti-ui-slop/references/new-work.md:13-16` ("Build the complete state": loading, empty, error, success, disabled or recovery states required by the task) and `anti-ui-slop/SKILL.md` "Finish" (render and inspect once; fix clipping, overlap, inert interactions). `frontend` carries the restraint rules (Apps & Dashboards, lines 119-135) and a Loading States section in `references/component-patterns.md:248`, but no rule requiring the state set, and no render-and-inspect step. Partly covered by `skills/feature/SKILL.md:107` ("walk the golden path and one edge case (empty state, error response …)"). Whether the strike (operator D28) was meant to drop these is the operator's call; named so it is a decision, not an accident.
- **[low / confidence: high]** `.agents/docs/testing/skill-matrix.md:63` — the `gcloud-auth` row's check `grep -qiE "gcloud-auth-verification\|gcloud-setup"` still accepts the name of a skill removed on 2026-09-13 (VENDORED.md says it was merged into `gcloud-setup`). Not changed by Phase 6, but the row was touched in the same file and the stale alternative can pass a run that loaded nothing current. Fix: check for `gcloud-setup` only.

## Checked and clean (no finding)
- `README.md`: no skill count or list names `python`, `anti-ui-slop` or a fork pair; line 256 ("Record third-party provenance in `skills/VENDORED.md`") still true (`mermaid-syntax` carries no row, and VENDORED.md says why).
- `context/dev-setup.md:10` (`terraform`, `gke-basics`, `cloud-run-deploy`) and `context/AGENTS.md:55` ("the `golang` skill has the runtime-data rule": present at `skills/golang/SKILL.md:45-46`) resolve. `agents/reviewer.md:35`, `skills/review-standards/SKILL.md:47,84` point at `golang`, `typescript`, `frontend`, which exist; no pointer to a deleted `golang-*.md`, `python` or `anti-ui-slop` anywhere in `context/`, `agents/`, `skills/*/SKILL.md`, `README.md`, `AGENTS.md`, specs, guide or deploy docs.
- Routing phrases `(use terraform)` in `bash`, `cloud-run-deploy`, `docker`, `frontend`, `gcloud-setup` still name a skill that exists; no `(use python)` remains.
- `AGENTS.md:62` ("the one fork … `gke-basics`; `terraform` is the harness's own") matches `VENDORED.md` (one row in the second table, no terraform row, no LICENSE under `skills/terraform`).
- `skills-lock.json` does not exist, so no lock entry for a struck skill.
- `skills/init-project/assets/agents-tree/docs/deploy/gcp.md` § Rollout and network matches `skills/cloud-run-deploy/SKILL.md` § What `docs/deploy/gcp.md` records row for row (five rows).
- `docs/deploy/{local,remote,gcp}.md` of this repo: no mention of the changed skills.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- Technical correctness of skill content (Go samples compile, `terraform validate`, flag names, `vite build` order D37): plan Phase 6 verify commands; not rerun here (docs lens).
- `.agents/docs/plans/2026-09-30-audit-fix.md`: Phase 6 boxes (lines 154-163) were still unticked except the tools row at review time; ticking them and the `changes/` entry for Phase 6 is not yet on disk (newest `changes/` entry is `2026-10-01-0737-audit-fix-phase-5.md`). Trail bookkeeping, not a doc contradiction; left to the commit.
- `skills/*/references/**` line-by-line prose for internal cross-reference accuracy (only `TestSkillLinksResolve` and the grep sweeps above covered links).
- Records (`changes/`, `DECISIONS.md`, `handoffs/`, `reviews/`, `retros/`, bug records, dated logs) were not searched by design; `bugs/INDEX.md` was read only for the BUG-022 row.
- `skills/review/SKILL.md` and `skills/review-standards/SKILL.md` one-line edits are in the working tree but outside the dispatched scope; seen in the diff, no problem noticed.
- Antigravity behaviour of any changed skill: not run.

## Outcome
- [medium] Phase 7 skill names in cloud-run-deploy — rebutted: D38; cells now name the phase and the interim answer (see the code lens Outcome)
- [low] state bucket has no slot in the gcp.md template — fixed at `skills/init-project/assets/agents-tree/docs/deploy/gcp.md` (row "Terraform state" in § Rollout and network)
- [low] VENDORED "a third" names two — fixed at `skills/VENDORED.md` (paragraph reworded)
- [low] BUG-022 Graduated row points at the deleted `golang-ecosystem.md` — fixed: the rule's new home is `skills/VENDORED.md`'s opening paragraph; the INDEX row repointed by `bugs` (the one edit a past record allows)
- [low] anti-ui-slop's "complete state set" and "render and inspect once" rules dropped — fixed at `skills/frontend/SKILL.md` § Apps & Dashboards (two bullets)
- [low] matrix row `gcloud-auth` accepts the removed `gcloud-auth-verification` — fixed at `.agents/docs/testing/skill-matrix.md:63`

## Round 2

Range reviewed (round 2): working tree (`git status --short --untracked-files=all`), Phase 6 of the audit-fix plan, read with `git diff HEAD -- <files>` and by hand; scope as in round 1 plus the other lenses' fixes. Lens: docs, report only.

Verdict (round 2): APPROVE WITH NITS. All five round-1 fixes are correct and complete (checked below). Two restatements the other lenses' fixes left behind remain, both low.

### Findings
- **[low / confidence: high]** `skills/golang/SKILL.md:56` and `skills/gke-basics/SKILL.md:23` — both say the image base is `golang:<the go.mod go line>`. Read literally, `go 1.27.1` gives the three-part `golang:1.27.1`, which `skills/docker/SKILL.md:86` (rule 1) now forbids ("a three-part tag such as `golang:1.24.4` freezes a patch with known advisories") and `:68-69` writes as `golang:<major.minor of the go.mod go line>`. Three files restate one rule in two wordings, and the shorter one gives the wrong tag. Fix: use `golang:<major.minor of the go.mod go line>` in both (the `restatement-grep` on `go.mod go line` found exactly these three sites plus docker's own `:128`).
- **[low / confidence: medium]** `skills/init-project/assets/agents-tree/docs/deploy/gcp.md:25` vs `skills/cloud-run-deploy/SKILL.md:170-181` — the template's § Rollout and network intro says every row is "read by the `cloud-run-deploy` skill (its § What `docs/deploy/gcp.md` records)", and round 1 confirmed five rows matched; the new sixth row "Terraform state" (`gcp.md:40`) is read by the `terraform` skill and is absent from cloud-run-deploy's table (five rows, "these rows"). The intro now overstates. Fix: reword the intro ("…each read by `cloud-run-deploy`, except Terraform state, which the `terraform` skill reads") or add the row to cloud-run-deploy's table.

### Round-1 fixes verified
- Terraform state row: present (`gcp.md:40`); content agrees with `skills/terraform/SKILL.md:17,36,41` (`PROJECT-tfstate`, suffixed name on collision, prefix). Only the intro nit above.
- `VENDORED.md`: "Two skills … `terraform` …; and `gcloud-auth-verification` …" agrees; the opening paragraph carries BUG-022's rule; `bugs/INDEX.md:108` points at that paragraph. Correct.
- `frontend/SKILL.md:121-122`: both rules present, worded as rules with the evidence named. Correct.
- `skill-matrix.md:63`: accepts `gcloud-setup` only. `restatement-grep` for `gcloud-auth-verification` leaves only `VENDORED.md:27` (the removal note) and the PRIORITY_TODO line. Correct.
- `cloud-run-deploy` `:96,:100`: each names "audit-fix Phase 7"; `:100` carries an interim answer (`service.md`, whose heading is "Database: Postgres on Cloud SQL", `service.md:288`); `:96` (`--iap`) names the idtoken rule inline, so no dead end. Rebuttal D38 accepted.

### Other lenses' changes, restatement sweep
- Docker three stages and `golang:` tag: docker SKILL consistent (`:68,69,86,128`); the golang/gke-basics wording is the finding above.
- `.dockerignore`: one list at `docker/references/setup.md:9-23` (`**/` prefixes); SKILL rule 4 points at it; no other living doc quotes an old list (README, context, preflight list secret-shaped names for commits, a different rule). Clean.
- Terraform secret versions: `terraform/assets/example/main.tf:69` and `references/resources-gcp.md:175` pin to the managed version; `cloud-run-deploy:101` says pin, not `:latest`; no living doc says `latest` for a secret. Clean.
- Frontend default theme: `SKILL.md:25`, `theme-system.md:3,181-182,209,224-229` all say system-on-first-load, light when the OS states no preference; `:104` ("Light … (default)") labels the palette, `:324` ("light OS, nothing stored (default)") is consistent. Clean.
- Compose `postgres:18` (`setup.md:38`) agrees with `terraform/SKILL.md:51` and `resources-gcp.md:107` (`POSTGRES_18`). Clean.
- healthz any 2xx: `health-checks.md:28` (comment says golang's handler answers 204; `golang/references/service.md:109` does). The docker skill's own example handler (`health-checks.md:11-16`) answers 200; both are accepted by the check. Clean.

### Not examined (round 2)
- Technical correctness of changed skill content (Terraform `validate`, Dockerfile build, cloudsqlconn Go 1.26 floor): not rerun, docs lens.
- `skills/*/references/**` prose beyond the swept restatements.
- Records (`changes/`, bug records, plans); Antigravity behaviour.


## Outcome (round 2)
- [low] `golang/SKILL.md:56`, `gke-basics/SKILL.md:23` base-image wording — fixed (`<major.minor of the go.mod go line>`)
- [low] template intro says every row is read by `cloud-run-deploy` — fixed at `skills/init-project/assets/agents-tree/docs/deploy/gcp.md` (five rows for `cloud-run-deploy`, the Terraform state row for `terraform`)

## Round 3

Range reviewed (round 3): working tree (`git status --short --untracked-files=all`), Phase 6 of the audit-fix plan, read with `git diff HEAD -- <files>` and by hand, plus `bash skills/review/scripts/restatement-grep.sh -i -- <key words>` for each round-2 change (`dockerignore`, `distroless`, `postgres:`/`var/lib/postgresql`, `immutable`, `db_url`/`migrate`, `fieldViolation`/`ErrorInfo`/`field:`, `allow-unauthenticated`/`iap`, `go.mod go line`, `mutable`). Lens: docs, report only.

### Status (round 3): COMPLETE

### Verdict (round 3): APPROVE WITH NITS

Both round-2 fixes are correct and complete, and no round-2 fix from the other lenses left a stale restatement of the allowlist, the digest pin, compose, the migrate job's secret, `fieldViolations` or the IAP cell. One low nit: a comment in cloud-run-deploy still states the SHA tag is immutable without the condition the same file adds 17 lines later.

### Findings
- **[low / confidence: medium]** `skills/cloud-run-deploy/SKILL.md:54` — the build snippet's comment says `# immutable tag`, a bare claim. Artifact Registry tags are mutable unless the repository was created with `--immutable-tags` (`:71`, `resources-gcp.md:62`, which both say "the SHA tag is an audit trail only with this"); `:71` covers only the "repository doesn't exist yet" path, so a repository that already exists and was created without the flag still gets this comment read as a guarantee. This is a stale restatement left by the immutable-tags fix (the other sites, `:69`, `:71`, `main.tf:21`, `resources-gcp.md:62`, are conditional or accurate). Fix: `# SHA tag; immutable only when the repository has immutable tags (below)`.

### Round-2 fixes verified
- `skills/golang/SKILL.md:56` and `skills/gke-basics/SKILL.md:23`: both `golang:<major.minor of the go.mod go line>`; agree with `skills/docker/SKILL.md:68,69,86`. The `restatement-grep` for `go.mod go line` finds no other wording. Correct.
- `skills/init-project/assets/agents-tree/docs/deploy/gcp.md:25`: first five rows for `cloud-run-deploy`, the Terraform state row for `terraform`; `cloud-run-deploy/SKILL.md:170-181` still has the five rows, row for row. Correct.

### Round-2 restatement sweep (other lenses' fixes)
- `.dockerignore` allowlist: one list at `skills/docker/references/setup.md:8-33`; `docker/SKILL.md:90` (rule 4) points at it, `:105` says `web/dist/` is in it (it is, `setup.md:18`), rule 5 (`:91-95`) is about layers and secrets, not the ignore list. `AGENTS.md`, `context/*.md`, README quote no ignore list; `context/AGENTS.md:47`'s secret-shaped list is about commits (not a restatement). No denylist survives. Clean.
- Distroless digest: `docker/SKILL.md:78,86` pin by `@sha256:<digest>`; `:87-88,95,104,128`, `health-checks.md:3,21,76` use `distroless/static` generically, not as a pin. Clean.
- Compose: `postgres:18` at `/var/lib/postgresql` (`setup.md:43,51`) agrees with `terraform/SKILL.md:51` and `resources-gcp.md:110` (`POSTGRES_18`); no other living doc names a Postgres major or volume path. Clean.
- Immutable tags: see the finding; `terraform/assets/example/main.tf:21` and `resources-gcp.md:62` carry `immutable_tags = true` consistently; no other living doc calls the SHA tag immutable.
- Migrate job reads `db_url`: `resources-gcp.md:249-272` grants the migrate SA the secret and mounts `db_url`; `golang/references/service.md:316` (owner DSN a second secret "readable only by the job's SA") agrees; no living doc says the job needs no secret (grep for `no secret`, `needs no`, `secret-free` found only unrelated hits). `terraform/SKILL.md:54` names the job without the secret, which is not a "needs no secret" claim. Clean.
- Form and `fieldViolations`: `frontend/references/component-patterns.md:192-216` reads `error.details[].fieldViolations`; `golang/references/service.md:43-44,67-73,110` emits `ErrorInfo` plus `BadRequest.fieldViolations`; `json-api-style/SKILL.md:31` (rule 7) and `json-api-conventions.md:96-125` agree. No living doc says errors are `{field: message}` (the only such phrase is in `plans/2026-09-30-research-audit-fix-frontend.md:88`, a dated research record). Clean.
- cloud-run-deploy IAP and `--allow-unauthenticated` rows (`:95-96`): `terraform/SKILL.md:53` agrees (Terraform owns the setting, neither flag passed); `agents/security.md:29` only lists the flag as a surface; `:178` and the template's `gcp.md:36` carry the `--iap` row. `api-auth-gcp` and `data-gcp` pointers name Phase 7 with an interim answer (accepted in round 1). Clean.
- Skill matrix, testing guide, README, `context/dev-setup.md`, `skills/preflight/SKILL.md`: no row or sentence restates any of the above.

### Not examined (round 3)
- Technical correctness of the changed skill content (Dockerfile build, `terraform validate`, the digest or `idtoken` claims); docs lens.
- `skills/*/references/**` prose beyond the swept restatements; Antigravity behaviour.
- Records (`changes/`, bug records, plans, other reviews and verify logs).

## Outcome (round 3)
- [low] `cloud-run-deploy/SKILL.md:54` bare "immutable tag" comment — fixed per D40 (the comment names the repository setting that makes it so and drops the alias)
