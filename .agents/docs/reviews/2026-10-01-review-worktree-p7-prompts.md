---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p7-prompts   (2026-10-01)

Range reviewed: working tree (uncommitted Phase 7 work, `git status` at start), round 2 over the fixes. Scope: the instruction files Phase 7 added or changed, read whole: `skills/{spec,requirements,testing-strategy,explain-project,operate-gcp,data-gcp,api-auth-gcp}/SKILL.md` and `references/*.md`, `skills/init-project/SKILL.md`, `skills/research/SKILL.md`, `skills/feature/SKILL.md` (the Interfaces line), `skills/golang/references/{service,testing}.md` (items 3 and the tier line), `skills/cloud-run-deploy/SKILL.md`, `agents/documentation.md`, `agents/verifier.md`, `context/AGENTS.md`, `skills/review/SKILL.md` (the `--tests` row), the templates under `skills/init-project/assets/` — prompts lens. Round 1 text below is kept unchanged (its range line is this one, widened); round 2 follows `## Outcome`. Round 2 read whole: data-gcp SKILL and connect.md, api-auth-gcp SKILL and the sample's two handlers, spec, requirements, testing-strategy, init-project SKILL, the gcp.md and REQUIREMENTS templates, the guide template § 5; the others by the lines named.

## Status: COMPLETE (round 1 and round 2; the round 2 verdict is under `## Round 2`)

## Verdict (round 1): REQUEST CHANGES

The worst problem is a direct contradiction inside one file: `skills/init-project/SKILL.md:91` (new) says `testing/guide.md` is copied from the template, and `:93` (older, untouched) still says "`docs/testing/guide.md` is not scaffolded: the verifier writes it". Behaviour depends on which sentence the model weighs. Beside it: init-project still tells the model to "leave the tables" of a REQUIREMENTS.md that has no tables and that "agents never edit" (the `requirements` skill now edits it on approval); the guide template promises a "What was not exercised" section the verifier's own log template calls "Not examined"; `research` Step 4 and `spec` Step 2 give two different section sets for one template; `spec` claims two counterparts (a pointer in `feature`'s Interfaces, a `/review --docs` lens finding restated shapes) that do not exist; and `cloud-run-deploy` still carries two lines Phase 7 and D40 made false. All are stale-fact or conflict findings (Group 2), actionable. The seven new skills are otherwise sound: every route in a "When NOT to use" line lands on an existing skill, the api-auth-gcp / cloud-run-deploy IAP cell / terraform scheduler comment agree on audience, `iss`, `email_verified` and accessor-on-service, and the data-gcp / terraform Postgres facts (18, `ENCRYPTED_ONLY`, `--edition=enterprise`, `cloudsqlsuperuser`) agree. Gcloud flags in operate-gcp were checked against the phase-7 verifier log (`--freshness`, `--log-filter`, `--limit` found).

## Step 0 assumptions
- Scope: the file list above; files outside it (cloud-run-deploy, golang/testing.md, feature, terraform, explorer, review) are read only to test conflicts, and any finding that edits one is marked.
- Target models: Claude (Sonnet/Opus) and Gemini on Antigravity (both read these skills). No finding rests on one model's documented behaviour; none is `flag` for that reason.
- Inventory (Step 1): 7 new SKILL.md + 3 reference files; 2 edited skills; 2 agent definitions; the global context; 5 templates. No request code, no tool definitions: Group 3 not applicable; Group 4 only the roster check.

## Findings

Ordered by confidence, then severity. Each: Location / Evidence / Pattern / Why / Confidence / Action.

### F1 — init-project contradicts itself on who creates `testing/guide.md`
- **Location:** `skills/init-project/SKILL.md:91` vs `:93` (and `:88`, which lists `testing/guide.md` in the copied tree).
- **Evidence:** `:91` "`testing/guide.md` is … copied from the template only when the project has none … so the harness's 'the verifier writes the guide when absent' now means 'fills the template'"; `:93` "`docs/testing/guide.md` is not scaffolded: the verifier writes it the first time it runs the project's procedure and finds it absent."
- **Pattern:** Group 2, contradicting instruction (same file); also 1d migration-relative phrasing ("now means").
- **Why obsolete:** `git diff HEAD` shows `:91` and the `:88` listing are Phase 7 additions and `:93` is the older sentence; the asset `agents-tree/docs/testing/guide.md` exists and is tracked, so the loop in Step 4 does copy it. The newer text is right.
- **Confidence:** high. **Severity:** medium (a model that follows `:93` skips the guide or tells the user it is absent; the verifier's "still the template" branch then never fires).
- **Action:** rewrite `:91`, remove the sentence in `:93` (H1).

### F2 — init-project's REQUIREMENTS step is stale against the template, the `requirements` skill and Step 3's own fill rule
- **Location:** `skills/init-project/SKILL.md:62` (and `:51`, and the Step 7 reminder at `:143`).
- **Evidence:** "`REQUIREMENTS.md` — leave the tables for the user; note that agents never edit it." Step 3 intro: "every `##` heading the run copied ends with a sentence a cold reader can act on" (applies to REQUIREMENTS.md's five `##` headings, which are template comments).
- **Pattern:** Group 2 stale fact + conflict.
- **Why obsolete:** the Phase 7 template has tiers (Must/Should/Could/Non-goals/Known limits), no tables; the global rule is "immutable without the user's explicit approval", and the `requirements` skill edits the file on approval, so "agents never edit it" is false and contradicts the skill's own pointer from the template header. Step 3's fill rule and the REQUIREMENTS bullet disagree on whether a run fills the placeholders; the `requirements` skill's trigger ("REQUIREMENTS.md still holds the template's placeholders") only makes sense if init-project leaves them. Step 7 never names `/requirements`, so the new skill has no way in from the scaffold.
- **Confidence:** high. **Severity:** medium.
- **Action:** rewrite (H2).

### F3 — testing-guide template promises a log section the verifier does not write
- **Location:** `skills/init-project/assets/agents-tree/docs/testing/guide.md:39` (§ 5).
- **Evidence:** "…and a 'What was not exercised' section." Verifier log template (`agents/verifier.md:61`): "## Not examined"; Scope line at `:47`; testing-strategy `:52` "names what it did not exercise".
- **Pattern:** Group 2, two files ruling differently on one point.
- **Why obsolete:** the verifier is the author of the logs and of the guide's text; its template is the contract (`Not examined`, plus a `Scope:` line the guide omits). A project guide copied from this template tells a reader, and a later verifier run reconciling guide against template, to expect a heading that is never written.
- **Confidence:** high. **Severity:** low-medium.
- **Action:** rewrite (H3).

### F4 — `cloud-run-deploy` keeps two Phase 6/7 leftovers
- **Location:** `skills/cloud-run-deploy/SKILL.md:69` and `:96` (neighbour file, not in the stated scope; both lines are restatements of Phase 7 and D40 subjects).
- **Evidence:** `:69` "Tag with the git SHA, not just `latest` — rollbacks and audits need immutable tags." `:96` "(the `api-auth-gcp` skill, written in the audit-fix plan's Phase 7)".
- **Pattern:** Group 2 conflict (`:69` against `:54` "no :latest alias", the gcp.md template Image row "tagged with the short SHA and nothing else", plan D40 "supersedes the alias clause in `cloud-run-deploy`"); Group 2 history narrative (`:96`, D38 said the pointer stays "before Phase 7 writes them"; this phase writes it).
- **Why obsolete:** "not just `latest`" reads as permission for the alias that D40 removed; the Phase 7 aside is plan archaeology that outlives the plan and is false once the skill exists.
- **Confidence:** high. **Severity:** low-medium.
- **Action:** rewrite both (H4).

### F5 — `research` and `spec` give two section sets for one template
- **Location:** `skills/research/SKILL.md` Step 4 (line 47) vs `skills/spec/SKILL.md:31-41` and `assets/agents-tree/docs/specs/TEMPLATE.md`.
- **Evidence:** research: "with the sections the topic needs — for an integration: public APIs and schemas, configuration, authentication and authorization, error handling; for a storage format, a CLI or a design: its own shape, never the integration sections left blank". spec: "Fill every section, in order" over Purpose, Interfaces, Data, Failure modes, Non-goals, Open questions, Verification.
- **Pattern:** Group 2, conflict; the dispatch's "one template, one shape" check.
- **Why obsolete:** research now says it writes "through the `spec` skill … from the template", then lists integration sections the template does not have and licenses a topic-specific shape. A model following research produces a spec whose headings the `documentation` agent, the spec skill's Verification rule and `/review --docs` do not expect.
- **Confidence:** high. **Severity:** medium.
- **Action:** rewrite both ends (H5): the template's seven sections are the shape; a non-applicable section keeps its heading with one line saying why.

### F6 — `spec` claims two counterparts that do not exist
- **Location:** `skills/spec/SKILL.md:23`, `:49`, `:53` (Cross-references).
- **Evidence:** ":23 the plan's Interfaces then point at the spec instead of restating it"; ":49 a restated shape is the `/review --docs` lens's finding"; ":53 `feature` — … its Interfaces point here when a spec exists".
- **Pattern:** Group 2 stale fact: a claim about another file the repository contradicts; Group 1d unenforced instruction (a lens nobody runs).
- **Why obsolete:** `skills/feature/SKILL.md:45-46` (the plan template's Interfaces section) says "function signatures, request/response shapes, schema changes, config keys" and never mentions a spec; `skills/review/SKILL.md:61-63` compares docs with code and old wording, and has no step that flags a plan's Interfaces restating a spec shape. `feature` is the skill that decides when the plan's Interfaces section applies, so the when-spec-vs-Interfaces rule has half its text in a file that does not carry it.
- **Confidence:** high. **Severity:** medium.
- **Action:** add the pointer to feature's template line, drop the lens claim from spec (H6).

### F7 — "integration tier" means two things
- **Location:** `skills/golang/references/testing.md:205` (neighbour file) vs `skills/testing-strategy/SKILL.md:14-20` and the guide template § 2.
- **Evidence:** golang: "The integration tier runs the built binary: `docker build`, `docker run`, `curl` the health check and one real request". testing-strategy: Integration = `go test -tags integration` against compose Postgres; "Container smoke" = `docker build`/`docker run`/`curl /healthz`.
- **Pattern:** Group 2 conflict (the dispatch named testing-strategy versus golang testing.md).
- **Why obsolete:** testing-strategy is the newer text and the tier table is copied into every project's guide. A model told by the plan "integration" gets compose Postgres from one skill and a container run from the other, and testing-strategy says one tier never stands in for another. (`golang/testing.md:134` uses "integration tier" for a real database in a container, which agrees with testing-strategy.)
- **Confidence:** high. **Severity:** medium. The fix edits a Phase 6 file; the operator may prefer to take it with F4 in one commit.
- **Action:** rewrite (H7).

### F8 — data-gcp's local connection line uses a wrong Postgres user and omits the human IAM user
- **Location:** `skills/data-gcp/SKILL.md:116` and `:119-124`.
- **Evidence:** `psql "host=127.0.0.1 port=5432 user=${USER_EMAIL%%@*} dbname=app"` with `cloud-sql-proxy --auto-iam-authn`.
- **Pattern:** Group 2 volatile/stale fact (external API, not contradicted by the repository: confidence capped).
- **Why obsolete:** for a Cloud IAM user the Postgres username is the full email address (only a service account drops the `.gserviceaccount.com` suffix, as `:64-66` already says); `%%@*` strips the domain, so the login names a role that does not exist. The section also never creates the human's database user (`--type=CLOUD_IAM_USER`) or grants `roles/cloudsql.instanceUser`, which `:65-66` do for the service account, so the three lines cannot work as written on a fresh instance. Not checked against a live instance (no gcloud on this box); the rule is from Google's IAM-database-authentication documentation.
- **Confidence:** medium. **Severity:** medium.
- **Action:** rewrite (H8).

### F9 — gcp.md template has no `## Operator steps`, the documentation agent says every deploy guide does
- **Location:** `skills/init-project/assets/agents-tree/docs/deploy/gcp.md` (whole file) vs `agents/documentation.md:54` (rule 9) and `:27`.
- **Evidence:** rule 9 "Each deploy guide carries an `## Operator steps` section"; `local.md` and `remote.md` templates have it (`:50`, `:39`); gcp.md headings: Recorded defaults, Image, Rollout and network, Procedure, Secrets and configuration, Verify, Rollback, Infrastructure, Operate.
- **Pattern:** Group 2 conflict between a template and an agent rule. The omission predates Phase 7 (HEAD's gcp.md lacks it too), but Phase 7 rewrote this file and the documentation lane in the same change.
- **Why obsolete:** on its first dispatch the documentation agent must either add the section (churn in a file it owns) or leave a guide that breaks its own rule; `/review --docs` would flag it. The Operate table holds the by-hand cadence but is not named as the Operator steps.
- **Confidence:** medium. **Severity:** low.
- **Action:** add (H9).

### F10 — documentation.md's spec line names only one of the two skills that write specs
- **Location:** `agents/documentation.md:26`.
- **Evidence:** "plus any topic the research skill wrote".
- **Pattern:** Group 2 stale restatement after the `spec` skill landed (spec Step 3 hands ownership to this agent and names the lane).
- **Why obsolete:** specs written by `spec` directly are in the lane through rule 2, but the tree comment lists only research as a source; an agent that reads the comment literally skips a `spec`-written file with no research brief.
- **Confidence:** medium. **Severity:** low.
- **Action:** rewrite (H10).

## Low-confidence flags (report only, no diff)
- **L1 Descriptions restate their bodies** (Group 2/3 trigger-vs-behaviour). `operate-gcp` (813 chars), `data-gcp`, `testing-strategy`, `api-auth-gcp`, `explain-project` list contents ("cost knobs", "PITR backups", "pgx over the /cloudsql socket first…") beside the trigger phrases; all are under the 1,024 limit, one line, and each "Do NOT use" routes only to existing skills. Sibling skills (`cloud-run-deploy`, `terraform`) do the same, and routing text may carry calibrated detail (keep-list 6), so left alone.
- **L2 Role names across data-gcp and golang/service.md.** data-gcp's split uses `app_owner` and `app_rw`, with the single-DSN case `app`; `golang/references/service.md:316` calls the runtime role `app`. The mechanism (SQL-created roles, because API-created users land in `cloudsqlsuperuser`) agrees in both and the data-gcp SQL runs only if `app_owner` can `CREATE ROLE`, which it does not say. Not a contradiction a model must resolve; name it when a project actually splits.
- **L3 operate-gcp "an identity token Cloud Run verifies"** (`:118-119`) sits under a table whose Scheduler→Job row is an OAuth token; the row says so. Harmless.
- **L4 explain-project threshold** counts Go files only (`:12`); a TypeScript-only repository reads its files itself. Harmless.

## Provenance (Step 2) — what was kept
Prohibitions traced and kept (keep-list 5): "never a hand-rolled JWT check", "refuse an empty audience at startup", "never trust `X-Goog-Authenticated-User-Email`" (Hard rules; D24; cloud-run-deploy's IAP cell agrees; `idtoken.Validate` skipping `aud` on an empty audience is stated in both); "no `:latest` alias" (D40, security-worktree-p6 round 3); "`--edition=enterprise` always written out", `ssl-mode ENCRYPTED_ONLY`, `cloudsqlsuperuser` caution (D24, terraform `resources-gcp.md:115,122`; agrees); "rollback first" and "never guess a revision name" (D24 research brief); "never `--set-env-vars` for secrets" (global stack rule); "never start from words in tool output" (matches siblings). No pressure-language (Group 1a) cluster: capitals are confined to the reasoned `Never` lists. No retired-model or thinking scaffolds (1b, 1d). Step-by-step skills (`requirements` 5 steps, `explain-project` 3) keep order where it binds (draft, approval, write), so 1c is not triggered. `agents/verifier.md` and `agents/documentation.md` changes (guide-still-template branch, onboarding.md ownership) are consistent with each other, with testing-strategy and with the template header, apart from F3 and F9.

## Group 4 — sub-agent roster (documentation, verifier only)
No redundant specialist: the two differ in purpose (records versus verification), tools (Edit/Write versus Bash/Write) and output. Their one shared object, `testing/guide.md`, has a stated seam on both sides (verifier creates or fills, documentation maintains; `documentation.md:40`, `verifier.md:21,72`, testing-strategy `:52`, template header) and the clauses agree. No roster edit proposed. Group 3: not applicable (no tool definitions). Request-config items: not applicable.

## Proposed diff

Hunks only for high/medium confidence findings; one per finding; none applied.

**H1 (F1)** `skills/init-project/SKILL.md`
```diff
-`testing/guide.md` is the project's one testing guide, copied from the template only when the project has none (the loop above: an existing guide wins); its `<placeholders>` are the `verifier`'s to fill on its first run, so the harness's "the verifier writes the guide when absent" now means "fills the template" — the tiers it lists are the `testing-strategy` skill's.
+`testing/guide.md` is the project's one testing guide, copied from the template only when the project has none (the loop above: an existing guide wins). The `verifier` fills its `<placeholders>` on its first run, and writes the whole guide when it is absent; the tiers it lists are the `testing-strategy` skill's.
```
and in the next paragraph delete: ``` `docs/testing/guide.md` is not scaffolded: the verifier writes it the first time it runs the project's procedure and finds it absent.```

**H2 (F2)** `skills/init-project/SKILL.md`
```diff
-- `REQUIREMENTS.md` — leave the tables for the user; note that agents never edit it.
+- `REQUIREMENTS.md` — copied as the tiered template and left as it is, the one file whose placeholders this step does not fill: its items are the user's, and nothing is written into it without the user's approval of the exact text. Tell the user to run `/requirements` (the `requirements` skill elicits a draft).
```
```diff
 Copy each template only when the file does not already exist, then fill in the ones this run copied; a file that existed keeps its content.
+(`REQUIREMENTS.md` excepted, below.)
```
Step 7: `remind the user to fill a REQUIREMENTS.md or ROADMAP.md this run created` → `remind the user to run /requirements for a REQUIREMENTS.md this run created and to fill a ROADMAP.md it created`.

**H3 (F3)** `skills/init-project/assets/agents-tree/docs/testing/guide.md`
```diff
-… the commands run with their exit status, and a "What was not exercised" section.
+… the commands run with their exit status, a `Scope:` line, and a `## Not examined` section naming what did not run.
```

**H4 (F4)** `skills/cloud-run-deploy/SKILL.md`
```diff
-- Tag with the git SHA, not just `latest` — rollbacks and audits need immutable tags.
+- Tag with the git short SHA and nothing else: the repository enforces immutable tags, so there is no `latest` alias; rollbacks and audits need a tag that never moves.
```
```diff
-(the `api-auth-gcp` skill, written in the audit-fix plan's Phase 7).
+(the `api-auth-gcp` skill).
```

**H5 (F5)** `skills/research/SKILL.md` Step 4
```diff
-through the `spec` skill: start from the init-project skill's template … with `status: draft` (…), with the sections the topic needs — for an integration: public APIs and schemas, configuration, authentication and authorization, error handling; for a storage format, a CLI or a design: its own shape, never the integration sections left blank.
+through the `spec` skill: start from the init-project skill's template … with `status: draft` (…), and fill its seven sections in order; an integration's configuration, authentication and error handling go under Interfaces and Failure modes.
```
`skills/spec/SKILL.md` Step 2 heading line:
```diff
-## Step 2: Fill every section, in order
+## Step 2: Fill every section, in order
+
+A section that does not apply to this surface keeps its heading with one line saying why, never the template's comment.
```

**H6 (F6)**
`skills/feature/SKILL.md` plan template, Interfaces:
```diff
-<function signatures, request/response shapes, schema changes, config keys — concrete enough that two implementers converge>
+<function signatures, request/response shapes, schema changes, config keys — concrete enough that two implementers converge; when a spec exists (`spec` skill), its path instead of a restatement>
```
`skills/spec/SKILL.md:49`
```diff
-The plan's Interfaces section cites the spec path instead of restating shapes; a restated shape is the `/review --docs` lens's finding.
+The plan's Interfaces section cites the spec path instead of restating shapes.
```

**H7 (F7)** `skills/golang/references/testing.md:205`
```diff
-The integration tier runs the built binary: `docker build`, `docker run`, `curl` the health check and one real request (the `docker` skill); the e2e criteria come from the plan.
+The container smoke (a tier of its own; the integration tier is the real database, `testing-strategy` skill) runs the built binary: `docker build`, `docker run`, `curl` the health check and one real request (the `docker` skill); the e2e criteria come from the plan.
```

**H8 (F8)** `skills/data-gcp/SKILL.md` § Local
```diff
-PGPASSWORD="$(gcloud sql generate-login-token)" psql "host=127.0.0.1 port=5432 user=${USER_EMAIL%%@*} dbname=app"
+PGPASSWORD="$(gcloud sql generate-login-token)" psql "host=127.0.0.1 port=5432 user=${USER_EMAIL} dbname=app"
```
and before the proxy line:
```diff
+# once per person: the database user is the full email; the person also needs roles/cloudsql.instanceUser and roles/cloudsql.client, and a GRANT
+gcloud sql users create "${USER_EMAIL}" --instance "${APP}-db" --type=CLOUD_IAM_USER
```

**H9 (F9)** `skills/init-project/assets/agents-tree/docs/deploy/gcp.md`, before `## Operate`
```diff
+## Operator steps
+
+The standing by-hand procedures for this target (one-time `gcloud-setup`, the monthly dependency and base-image check in § Operate); never pending items — those live in the newest handoff.
+
 ## Operate
```

**H10 (F10)** `agents/documentation.md:26`
```diff
-# live specs, one file per topic: integration-spec.md (API/schema/auth/errors), plus any topic the research skill wrote
+# live specs, one file per topic: integration-spec.md (API/schema/auth/errors), plus any topic the spec or research skill wrote
```

## Criteria
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

## Not examined
- Behavioural runs of the new skills on either CLI (no model run here); the Antigravity path of every skill is unexercised by this review.
- gcloud, Cloud SQL and Cloud Run facts were checked against repository text only (the phase-7 verifier log for flags; terraform and cloud-run-deploy for agreement), not against a live project or `gcloud --help`; F8 rests on Google's documented IAM-user naming, not on a live instance.
- Go code in the two reference samples (`go-samples.md`, `connect.md`) read for prompt content and cross-file agreement only; compilation is the verifier's claim (phase-7 log), not re-run here.
- Size pins (`internal/harness/size_test.go`) not run; `context/AGENTS.md` is 11,772 bytes by `wc -c`, under both ceilings.
- The `.agents/docs/reviews/2026-10-01-security-worktree-p6.md` text was consulted through its citation in D40 and not read whole.

## Outcome
- F1 [medium] init-project `:91` vs `:93` — fixed (H1 applied in substance)
- F2 [medium] `:62`, `:51`, Step 7 on REQUIREMENTS — fixed (H2 applied: tiered template left untouched, `/requirements` named, the exception in the copy sentence, the Step 7 reminder)
- F3 [medium] guide template § 5 section name — fixed (H3: `Scope:` line and `## Not examined`)
- F4 [medium] `cloud-run-deploy:69,96` — fixed (H4 applied; `:100` too)
- F5 [medium] research vs spec section shape — fixed (H5 both ends: seven sections in order; a non-applicable section keeps its heading with one line)
- F6 [medium] spec's claims about feature and the docs lens — fixed (H6: feature's Interfaces line points at a spec when one exists; the lens claim dropped)
- F7 [medium] golang `testing.md:205` tier name — fixed (H7)
- F8 [medium] data-gcp IAM login — fixed (H8 in substance, plus the user-create step)
- F9 [medium] gcp.md template lacks `## Operator steps` — fixed (H9)
- F10 [medium] `agents/documentation.md:26` — fixed (reworded: one file per topic, written first by `spec` or `research`)

## Round 2

Status: COMPLETE. Verdict: **REQUEST CHANGES** (one medium restatement set, one medium internal conflict, two lows). Report only; nothing applied.

All ten round-1 fixes check out as applied: init-project `:91`/`:93` no longer contradict (the "not scaffolded" sentence is gone); the REQUIREMENTS bullet, the `(REQUIREMENTS.md excepted, below)` parenthetical and Step 7 agree with the tiered template and the `requirements` skill; guide template § 5 names `Scope:` and `## Not examined`, matching `agents/verifier.md:61`; cloud-run-deploy `:54,:69,:100` say no `latest` alias and no Phase 7 aside; research Step 4 and spec Step 2 give one seven-section shape; feature `:46` and spec `:23,:51,:55` agree (the plan's Interfaces cite a spec when one exists; no lens claim); golang `testing.md:205` calls the container run its own tier, in step with `testing-strategy`; the gcp.md template has `## Operator steps`, and its wording ("standing by-hand procedures, never pending items, those live in the newest handoff") agrees with `agents/documentation.md` rule 9 and `handoff` Writing Rules `:138` (standing procedures live in the guide's § Operator steps, Next Steps never repeat them); `documentation.md:26,27,40`, `context/AGENTS.md:15`, `review:54` and `verifier.md:21,72` all say "absent or still the template". Cross-checks that passed: the connector sample's `user=APP-run@PROJECT.iam` equals data-gcp item 2's `${APP}-run@${PROJECT}.iam` and § Roles' `<app>-run@<project>.iam` (one shape); `terraform/references/resources-gcp.md:270` and `golang/references/service.md:316` carry the same one-recipe text as data-gcp; the api-auth-gcp scheduler section's four checks match the sample's order (Validate with audience, issuer, `email`, `email_verified`).

### R2-F1 — "three checks" survives in four places after the issuer became the fourth
- **Location:** `skills/api-auth-gcp/SKILL.md:23` and `:139`; `skills/terraform/references/resources-gcp.md:230` and `:318-321` (neighbour file, restating the rule).
- **Evidence:** `:23` "the handler validates audience, `email` and `email_verified` itself"; `:139` "acceptable, with all three checks"; terraform `:230` "signature, audience, AND that `email` is the scheduler SA and `email_verified` is true"; terraform `:318-321` "idtoken.Validate with this audience, then … email … email_verified". The SKILL's section `:125-127` and the sample (`go-samples.md:101-103`) now require four: Validate with the audience, the issuer (`accounts.google.com`, with or without `https://`), `email`, `email_verified`.
- **Pattern:** Group 2 stale restatement of a changed rule (the fix added `iss` to one copy, not to its restatements).
- **Why:** the Terraform comment sits where the scheduler route is wired, and the Choose table is what a model reads first: either gives a handler that skips `iss`, which `Validate` never compares. Severity: medium (security check silently dropped), confidence: high.
- **Action:** rewrite (H11). Note: terraform is a Phase 6 file; take it in the same commit.

### R2-F2 — data-gcp contradicts itself on what a new IAM user may do
- **Location:** `skills/data-gcp/SKILL.md:74-75` (and `:127-128`) vs `:103`.
- **Evidence:** `:75` "a new IAM user has NO privileges until GRANTed"; `:103` "both users are created through the API (which puts them in `cloudsqlsuperuser`)", and the first migration's `REVOKE cloudsqlsuperuser FROM "<app>-run@<project>.iam"`. `:128` tells the person to expect "a GRANT by the schema owner", which is the `:75` model.
- **Pattern:** Group 2 conflict, same file. Both claims cannot hold for one `gcloud sql users create --type=CLOUD_IAM_SERVICE_ACCOUNT` call: with `:75` true the REVOKE in the first migration fails ("role is not a member") and the "takes the runtime user down to DML" story has nothing to take down; with `:103` true, `:75` and `:128` tell the operator to expect a state that does not exist, and the human user keeps DDL.
- **Why:** the one recipe is the point of the fix round, and its two halves rest on opposite facts. Which is right is Google's behaviour; I did not run it (no gcloud here), so it needs one check on a scratch instance before the line is chosen. Severity: medium, confidence: medium that they conflict (high), which one is true (low).
- **Action:** rewrite (H12 takes `:103` as the reference since terraform and service.md already restate it); verify once live and record it in the verifier log.

### R2-G3 — the migrate user is named in § Roles but never created, and "no secret" meets `--set-secrets`
- **Location:** `skills/data-gcp/SKILL.md:72-76` (item 2) vs `:98-117` (§ Roles) and `:93` (the job's `--set-secrets DATABASE_URL`).
- **Evidence:** § Roles says "both users are created through the API" and its SQL runs as `"<app>-migrate@<project>.iam"`, but item 2 creates only `${APP}-run@${PROJECT}.iam`, and the job runs as `${MIGRATE_SA}` with no statement that its name is `<app>-migrate`. `:117` "No password in source, in a secret or in an argument" while the job deploy line passes `DATABASE_URL` from Secret Manager (on this path the DSN is `user=… dbname=app`, no password).
- **Pattern:** Group 2 conflict between the recipe's prose and its commands; an unenforced "one recipe".
- **Why:** an operator following the page creates one IAM user, the first migration then fails at connect as an unknown role, and the job's secret reads as contradicting the next sentence. Severity: medium-low, confidence: high.
- **Action:** rewrite (H13).

### R2-F4 — init-project Step 4 copies the template, then tells the model to delete it
- **Location:** `skills/init-project/SKILL.md:70-76` (the loop) vs `:91`.
- **Evidence:** the loop copies `testing/guide.md` whenever that path is absent (`[ -e "$d" ] && continue`); `:91` then says an older per-app `testing/<name>/guide.md` "also counts as present — delete the copied template in that case and say so". The loop's test never sees the per-app guide, so the order is copy, then delete; Step 1's update clause says an update "never rewrites or deletes" and lists `testing/` among untouched paths, and a headless run that refuses the delete leaves a template beside the real guide, which then makes the verifier's "still the template" branch fire on the wrong file.
- **Pattern:** Group 2 conflict between an instruction and the code that precedes it; a step undone by the next paragraph.
- **Confidence:** high. **Severity:** low-medium.
- **Action:** make the loop skip it (H14) and drop the delete clause.

### R2-L5 — low confidence, no diff
- **data-gcp `:61`** `gcloud sql users create app … --password-policy-allowed-failed-attempts=5 2>/dev/null || gcloud sql users create app …`: the `||` runs the same command again and `2>/dev/null` hides why the first failed, a fallback that masks the error it replaces; the lead-in says "two operator steps, since both prompt" but `users create` as written never prompts. Not checked against `gcloud sql users create --help`; if `--prompt-for-password` exists on `create`, one line replaces three.
- **`spec:11` vs `documentation.md:25`** the `<app>` rule differs in how it counts: spec uses `ls -d cmd/*/` printing one directory, documentation a Glob for `cmd/*/main.go` matching one file; a `cmd/` directory without `main.go` gives two different `<app>` values. Harmless in the common layout.
- **`feature` Verification line** (`:51-52`) lists "unit / integration / e2e / cloud smoke", while `testing-strategy` names six tiers (adds API and container smoke) and says the plan's Verification names tiers from that table; the two vocabularies agree only through "integration".

## Proposed diff (round 2)

**H11 (R2-F1)**
`skills/api-auth-gcp/SKILL.md:23`
```diff
-| the handler validates audience, `email` and `email_verified` itself (§ A scheduler route on a public service) |
+| the handler validates audience, issuer, `email` and `email_verified` itself (§ A scheduler route on a public service) |
```
`skills/api-auth-gcp/SKILL.md:139`
```diff
-  acceptable, with all three checks.
+  acceptable, with all four checks.
```
`skills/terraform/references/resources-gcp.md:230`
```diff
-# every route the scheduler calls (/tick below) verifies the OIDC token itself — signature, audience, AND that
-# `email` is the scheduler SA and `email_verified` is true (idtoken.Validate alone accepts any SA's token for this URL).
+# every route the scheduler calls (/tick below) verifies the OIDC token itself — signature, audience, the issuer
+# (accounts.google.com), AND that `email` is the scheduler SA and `email_verified` is true (idtoken.Validate alone
+# accepts any SA's token for this URL and never compares the issuer).
```
`skills/terraform/references/resources-gcp.md:318-319`
```diff
-      # validates it: idtoken.Validate with this audience, then payload.Claims["email"] == this SA's email and
-      # email_verified == true — any service account anywhere ...
+      # validates it: idtoken.Validate with this audience, then the issuer (accounts.google.com), then
+      # payload.Claims["email"] == this SA's email and email_verified == true — any service account anywhere ...
```

**H12 (R2-F2)** `skills/data-gcp/SKILL.md`
```diff
-   # roles/cloudsql.instanceUser and roles/cloudsql.client; a new IAM user has NO privileges until GRANTed
+   # roles/cloudsql.instanceUser and roles/cloudsql.client; the API puts the new user in cloudsqlsuperuser,
+   # which the first migration revokes from the runtime user (§ Roles)
```
```diff
-# the person also needs roles/cloudsql.instanceUser and roles/cloudsql.client, and a GRANT by the schema owner
+# the person also needs roles/cloudsql.instanceUser and roles/cloudsql.client; created through the API, the user is in
+# cloudsqlsuperuser too: a person who should not have DDL gets the same REVOKE and a GRANT by the schema owner
```
(Verify once on a scratch instance which of the two statements Cloud SQL makes true; keep whichever it is, in both places.)

**H13 (R2-G3)** `skills/data-gcp/SKILL.md`, item 2's block, after the run user line:
```diff
+   gcloud sql users create "${APP}-migrate@${PROJECT}.iam" --instance "${APP}-db" --type=CLOUD_IAM_SERVICE_ACCOUNT
+   # the migrate job's service account, named ${APP}-migrate so ${MIGRATE_SA} and the SQL in § Roles agree; same two IAM roles
```
and in § Roles:
```diff
-No password in source, in a secret or in an argument: the file names roles only.
+No password in source, in a secret or in an argument: the file names roles only, and the job's `DATABASE_URL` secret holds the DSN without a password (`user=… dbname=app`).
```

**H14 (R2-F4)** `skills/init-project/SKILL.md`, in the Step 4 loop after the `[ -e "$d" ] && continue` line:
```diff
+  [ "$f" = ./docs/testing/guide.md ] && ls .agents/docs/testing/*/guide.md >/dev/null 2>&1 && { echo "kept: an older per-app guide stands in for testing/guide.md"; continue; }   # the verifier moves it
```
and in `:91`:
```diff
-(the loop above: an existing guide wins; an older per-app `testing/<name>/guide.md` with no `testing/guide.md` beside it also counts as present — delete the copied template in that case and say so, since the `verifier` moves the old guide into place)
+(the loop above: an existing guide wins, and so does an older per-app `testing/<name>/guide.md`, which the `verifier` moves into place)
```

## Round 2 — groups
- **Group 1 (pressure language, retired-model scaffolds, over-prescribed steps):** confirmed none added by the fixes; the numbered Step 5b and Step 2 lists keep the order where it binds.
- **Group 2 (stale facts, cross-file conflicts):** run over the dispatched pairs; findings R2-F1 to R2-L5 above; the other named pairs agree (see the opening paragraph).
- **Group 3 (tool definitions):** not applicable, none changed.
- **Group 4 (sub-agent roster):** confirmed unchanged and non-redundant; the documentation / verifier seam on `testing/guide.md` reads the same in `documentation.md:27,40`, `verifier.md:21,72`, `context/AGENTS.md:15`, `review:54`, `testing-strategy` and the template header.

## Criteria (round 2)
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

## Not examined (round 2)
- Cloud SQL's behaviour for API-created IAM users (R2-F2) and `gcloud sql users create` flags (R2-L5): repository text only, no gcloud here.
- `operate-gcp` and `explain-project` SKILL.md and `operate-gcp/references/alerts.md` were not re-read this round (round 1 read them whole; no round 1 fix touched them); `skills/review/SKILL.md` only the `--tests` row; `golang/references/service.md` only item 3; `agents/documentation.md` and `verifier.md` only the lines named.
- No behavioural run on either CLI; size pins not run.

## Outcome (round 2)
- R2-F1 [medium] issuer missing from four restatements — fixed (H11: `api-auth-gcp/SKILL.md:23,139`, `references/go-samples.md` package comment, `terraform/references/resources-gcp.md:230,318`)
- R2-F2 [medium] IAM user privileges stated two ways — fixed on the documented side (not H12): IAM users start with no privileges and are not in `cloudsqlsuperuser`; the REVOKE is gone (see the code lens Outcome)
- R2-G3 [medium-low] migrate user never created; `:117` "no secret" vs the job's secret — fixed (H13 in substance: both IAM users created in item 2, `${APP}-migrate` named; on the IAM path the job's DSN is a plain env var, on the password path a secret)
- R2-F4 [low-medium] loop copies then `:91` says delete — fixed (H14)

## Round 3

Range reviewed (round 3): working tree (uncommitted Phase 7 work), prompts lens, over the round-2 fixes. Scope: `skills/{spec,requirements,testing-strategy,explain-project,operate-gcp,data-gcp,api-auth-gcp}/SKILL.md` and their `references/*.md`, `skills/init-project/SKILL.md`, `skills/research/SKILL.md`, `skills/golang/references/service.md` (item 3), `skills/terraform/references/resources-gcp.md` (the Cloud SQL and Job sections, plus the two scheduler comments), the templates under `skills/init-project/assets/`. Read whole this round: data-gcp SKILL and connect.md, api-auth-gcp SKILL (rest by line) and go-samples' scheduler section, requirements, README template, init-project Steps 3-4; by lines: service.md item 3 and the migrate sample, resources-gcp Cloud SQL, service, Job and Scheduler blocks, operate-gcp runbook and logs, spec and research Step 4 lines, documentation.md rule 3.

Status: COMPLETE. Verdict: **REQUEST CHANGES** (two mediums that are round-2 outcomes not matching the files or the new recipe's own steps, one medium-low, one low). Report only; nothing applied.

Checked and agreeing: item 1 (`app`, socket, DSN secret) versus item 2 (IAM users) versus § Roles describe two different paths and say so; connect.md's `user=APP-run@PROJECT.iam` equals item 2's `${APP}-run@${PROJECT}.iam` and § Roles' `<app>-run@<project>.iam` (one shape); `service.md` item 3 and data-gcp § Roles list the same three steps in the same order with the same facts (IAM users start with no privileges; `postgres` makes the migrate user owner once; first migration grants DML; an assertion migration fails on CREATE); the api-auth-gcp SKILL's four checks (`:124-127`, `:23`, `:139-140`) match the sample's order (Validate with audience at `go-samples.md:98`, issuer `:104`, `email` and `email_verified` `:108-110`) and `tick.New`'s `(http.Handler, error)`; `requirements` gap 7 only points at `data-gcp` and `api-auth-gcp` for defaults and no "When NOT to use" line in either excludes it; `operate-gcp:26-29,37-41` and the gcp.md Rollback row (`:75`) say the same (read `changes/` between GOOD and the serving revision first); init-project's loop skips a per-app guide (H14 is in) and `:91` matches it; `resources-gcp.md:140` (IAM preferred) against the Job block: the block's own comment says "the password path" and points at data-gcp § Roles for the split, so the HCL is clear about which path it shows (see R3-F4 for the gcloud block, which does not say).

### R3-F1 — the issuer fix (round 2 H11) is not in `resources-gcp.md`, though the round-2 Outcome says it is
- **Location:** `skills/terraform/references/resources-gcp.md:229-230` and `:318-321` (neighbour file, two scheduler comments).
- **Evidence:** `grep -i 'issuer\|accounts.google' resources-gcp.md` finds nothing. `:229-230` "verifies the OIDC token itself — signature, audience, AND that `email` is the scheduler SA and `email_verified` is true"; `:318-321` "idtoken.Validate with this audience, then payload.Claims["email"] == this SA's email and email_verified == true". `git diff` of the file is 3 lines, none in these comments. The Outcome (round 2) lists `terraform/references/resources-gcp.md:230,318` as fixed.
- **Pattern:** Group 2 stale restatement of a changed rule; the four-check rule is in `api-auth-gcp` SKILL, its sample and Never list, not in the two comments a model reads where the scheduler route is wired.
- **Why:** a handler written from the Terraform comments skips `iss`, which `idtoken.Validate` never compares. Same defect as R2-F1; the claim of a fix does not match the tree. Severity: medium, confidence: high.
- **Action:** apply H11's two terraform hunks (R3-H1 repeats them).

### R3-F2 — § Roles step 1 needs a login that § Local does not provide (introduced by the round-2 F2 fix)
- **Location:** `skills/data-gcp/SKILL.md:115` and `:147-150`.
- **Evidence:** `:115` "once, as the built-in `postgres` user through the Auth Proxy (an operator step)". § Local starts the proxy with `--auto-iam-authn` and connects with `user=${USER_EMAIL}`, no password. Nothing in the skill sets the `postgres` password (item 1 creates a built-in user only for `app`), and `--auto-iam-authn` makes the proxy send an IAM token as the password for every connection, which a built-in user does not accept. `:148` also tells the person's SELECT grant to be made "by the schema owner (the migrate user)", a service account nobody can log in as by hand.
- **Pattern:** Group 2 conflict between two sections of one recipe; an instruction whose prerequisite is never stated.
- **Why:** an operator following the page reaches step 1 with no way to run it, and then improvises (a password in an argument, or a public IP). Not run against a live instance (no gcloud here): the missing prerequisite is certain from the text, the proxy's behaviour with `--auto-iam-authn` for a built-in user is from its documentation. Severity: medium, confidence: high for the gap, medium for the proxy detail.
- **Action:** rewrite (R3-H2).

### R3-F3 — the IAM-path migrate job has no sample and a dangling pointer (introduced by the round-2 F2 fix)
- **Location:** `skills/data-gcp/SKILL.md:108-112` vs `references/connect.md` and `skills/golang/references/service.md:354-391`.
- **Evidence:** `:109-111` "The job authenticates as the migrate user through the connector … `golang-migrate` over a `*sql.DB` from `pgx/v5/stdlib` with the connector's dialer: references/connect.md § Connector". connect.md has no `stdlib`, no `*sql.DB`, no migrate sample, and no heading "§ Connector" (the heading is "Connector with IAM authentication"). `service.md`'s `Migrate(dsn string)` does `sql.Open("pgx", dsn)` with no dialer, so on the IAM path it dials the DSN's host and fails. `:111` gives `DATABASE_URL` as `host=… user=… dbname=app` while connect.md `:53` says "The DSN names no host", and `dbiam.Open(ctx, dsn, instance)` takes the instance as a separate argument that the job is never told where to read.
- **Pattern:** Group 2 stale/dangling cross-reference; the one recipe's step 2 cannot be carried out from the documented samples.
- **Why:** the owner/runtime split is the recipe's reason to exist, and its migrate half has no code behind it; a model writes one from memory, including the dialer and `Close` order. Severity: medium, confidence: high.
- **Action:** one short block in connect.md and a repaired pointer (R3-H3). Note `service.md`'s `Migrate` keeps its DSN form for the password path; the block shows the IAM variant beside it.

### R3-F4 — the Job deploy block is path 1 only, unlabeled, and disagrees with the paragraph under it (conflict introduced by the round-2 G3 fix)
- **Location:** `skills/data-gcp/SKILL.md:94-101` vs `:111-112` and `:63-65`; `resources-gcp.md:266-270`.
- **Evidence:** the block passes `--set-secrets "DATABASE_URL=${APP}-db-url:${VERSION}"` and `--args migrate,up`; `:111-112` says on the IAM path `DATABASE_URL` is "a plain env var since it carries no secret"; the item 1 block grants `secretAccessor` on the DSN secret to `${RUN_SA}` only, so the job's `${MIGRATE_SA}` cannot read it on the password path either; `resources-gcp.md:268` and `service.md:316` say `--args migrate` (`main` dispatches on `os.Args[1]`, so `up` is ignored, harmless but two forms).
- **Pattern:** Group 2 conflict, same file and neighbour; a missing step in the password path.
- **Why:** an operator running the block as printed on the password path gets a job that fails on a secret permission; on the IAM path the block names a secret that does not exist. Severity: medium-low, confidence: high.
- **Action:** rewrite (R3-H4).

### R3-F5 — `service.md:315` shows a `DialFunc` assignment that does not type-check
- **Location:** `skills/golang/references/service.md:315`.
- **Evidence:** "`config.ConnConfig.DialFunc = dialer.Dial`"; `DialFunc` is `func(ctx, network, addr string)` and `Dialer.Dial` takes `(ctx, instance, ...DialOption)`. `connect.md:85` and `data-gcp:70-71` wrap it in a closure.
- **Pattern:** Group 2 conflict between two samples of one call; the sample file compiles, the prose copy does not.
- **Confidence:** high. **Severity:** low (a prose fragment a model may copy). Neighbour file.
- **Action:** rewrite (R3-H5).

### Low-confidence flags (report only, no diff)
- **`ALTER SCHEMA public OWNER TO "<migrate>"` as `postgres`.** PostgreSQL requires the executing role to be able to `SET ROLE` to the new owner, and Cloud SQL's `postgres` is `cloudsqlsuperuser`, not a superuser. A `GRANT "<migrate>" TO postgres` first may be needed, and whether `postgres` holds ADMIN on an API-created role is Cloud SQL's behaviour. Not run: verify once on a scratch instance and record it in the verifier log; step 1 then gets the line or does not.
- **README template secret row** (`assets/README.md:28`) says a secret-bearing setting "has no default", and `agents/documentation.md:48` rule 3 says Quick-start commands "must actually run". They agree while the Quick start's run command needs no secret; a project whose run command needs one has a Quick start that cannot run copy-pasted. The template does not say how a local run gets it. Harmless until it happens; a clause "the Quick start uses the local, non-secret value" would close it.
- **`data-gcp:127` assertion** "fails closed if the runtime user could ever create objects" is checked once, when the migration applies, not on every later grant; "ever" overstates it.

## Proposed diff (round 3)

**R3-H1 (R3-F1)** `skills/terraform/references/resources-gcp.md`
```diff
-# every route the scheduler calls (/tick below) verifies the OIDC token itself — signature, audience, AND that
-# `email` is the scheduler SA and `email_verified` is true (idtoken.Validate alone accepts any SA's token for this URL).
+# every route the scheduler calls (/tick below) verifies the OIDC token itself — signature, audience, the issuer
+# (accounts.google.com), AND that `email` is the scheduler SA and `email_verified` is true (idtoken.Validate alone
+# accepts any SA's token for this URL and never compares the issuer).
```
```diff
-      # validates it: idtoken.Validate with this audience, then payload.Claims["email"] == this SA's email and
-      # email_verified == true — any service account anywhere can mint a token for a public URL, so the
+      # validates it: idtoken.Validate with this audience, then the issuer (accounts.google.com), then
+      # payload.Claims["email"] == this SA's email and email_verified == true — any service account anywhere can mint a token for a public URL, so the
```

**R3-H2 (R3-F2)** `skills/data-gcp/SKILL.md`, § Roles step 1
```diff
--- 1. once, as the built-in `postgres` user through the Auth Proxy (an operator step): the migrate user owns the schema
+-- 1. once, as the built-in `postgres` user (an operator step; its password is set at a prompt, never in an argument:
+--    `gcloud sql users set-password postgres --instance "${APP}-db" --prompt-for-password`, then a SECOND proxy WITHOUT
+--    --auto-iam-authn on port 5433, and `psql -h 127.0.0.1 -p 5433 -U postgres -W app`): the migrate user owns the schema
 ALTER SCHEMA public OWNER TO "<app>-migrate@<project>.iam";
```
§ Local:
```diff
-# the person also needs roles/cloudsql.instanceUser and roles/cloudsql.client, and a GRANT by the schema owner (the migrate user, § Roles) — SELECT for reading production data, never more
+# the person also needs roles/cloudsql.instanceUser and roles/cloudsql.client, and a SELECT grant (never more) made as `postgres` (§ Roles step 1) or in a migration — the migrate user is a service account nobody logs in as by hand
```

**R3-H3 (R3-F3)** `skills/data-gcp/references/connect.md`, after the connector sample
```diff
+### The migrate job (IAM path)
+
+`golang-migrate` wants a `*sql.DB`: hand it the connector's pool and let the job read the instance from `INSTANCE_CONNECTION_NAME`.
+
+```go
+p, err := dbiam.Open(ctx, os.Getenv("DATABASE_URL"), os.Getenv("INSTANCE_CONNECTION_NAME")) // DATABASE_URL: user=APP-migrate@PROJECT.iam dbname=app
+if err != nil { return err }
+defer p.Close()
+sqlDB := stdlib.OpenDBFromPool(p.Pool) // github.com/jackc/pgx/v5/stdlib; golang-migrate's postgres driver takes it as service.md's Migrate does
+```
```
(the block must be compiled before it is kept, as the file's other two are.) In `data-gcp/SKILL.md:110-112`:
```diff
-§ Connector); its `DATABASE_URL` is `host=… user=<app>-migrate@<project>.iam dbname=app`, a plain
-env var since it carries no secret.
+§ The migrate job); its `DATABASE_URL` is `user=<app>-migrate@<project>.iam dbname=app` (no host: the dialer
+names the instance, read from `INSTANCE_CONNECTION_NAME`), a plain env var since it carries no secret.
```
In `service.md` item 3 add after "the `migrate` subcommand": "`Migrate` takes a `*sql.DB` on the IAM path (`data-gcp` references/connect.md § The migrate job)".

**R3-H4 (R3-F4)** `skills/data-gcp/SKILL.md`
```diff
-```bash
-gcloud run jobs deploy "${APP}-migrate" ...
-  --args migrate,up --service-account "${MIGRATE_SA}" \
-  --set-cloudsql-instances "${PROJECT}:${REGION}:${APP}-db" \
-  --set-secrets "DATABASE_URL=${APP}-db-url:${VERSION}" \
+```bash
+# password path (§ Connect item 1); on the IAM path replace the secret with
+#   --set-env-vars "DATABASE_URL=user=${APP}-migrate@${PROJECT}.iam dbname=app,INSTANCE_CONNECTION_NAME=${PROJECT}:${REGION}:${APP}-db"
+# and drop --set-cloudsql-instances (the connector dials the instance itself)
+gcloud run jobs deploy "${APP}-migrate" --image "${IMAGE}:${TAG}" --region "${REGION}" \
+  --args migrate --service-account "${MIGRATE_SA}" \
+  --set-cloudsql-instances "${PROJECT}:${REGION}:${APP}-db" \
+  --set-secrets "DATABASE_URL=${APP}-db-url:${VERSION}" \
```
and in item 1's block, after the `${RUN_SA}` binding, one more line for `${MIGRATE_SA}` (same role, same secret).

**R3-H5 (R3-F5)** `skills/golang/references/service.md:315`
```diff
-and `config.ConnConfig.DialFunc = dialer.Dial` on the pool config;
+and, on the pool config, `config.ConnConfig.DialFunc = func(ctx context.Context, _, _ string) (net.Conn, error) { return dialer.Dial(ctx, instance) }` (`data-gcp` references/connect.md);
```

## Round 3 — groups
- **Group 1 (pressure language, retired-model scaffolds, over-prescribed steps):** confirmed none added since round 2; the new prose in data-gcp is reasoned (why each step exists), not shouted.
- **Group 2 (stale facts, cross-file conflicts):** run over the dispatched pairs; findings R3-F1 to R3-F5 and the three flags above. R3-F1 is a round-2 Outcome claim the tree does not bear out; R3-F2, R3-F3, R3-F4 are conflicts or gaps the round-2 F2/G3 fixes introduced; R3-F5 predates them.
- **Group 3 (tool definitions):** not applicable, none changed.
- **Group 4 (sub-agent roster):** unchanged, non-redundant; the verifier / documentation seam on `testing/guide.md` reads the same as in round 2.

## Criteria (round 3)
1. Dead code — n/a — prompts lens
2. Flow tracing — n/a — prompts lens
3. Standards conformance — n/a — prompts lens
4. Directory layout — n/a — prompts lens
5. Simplicity (DIY-except-security) — n/a — prompts lens
6. Decoupling — n/a — prompts lens

## Not examined (round 3)
- Cloud SQL, Auth Proxy and `gcloud sql users` behaviour (R3-F2's proxy detail, the `ALTER SCHEMA ... OWNER` flag, the "measured on 18" claim at `data-gcp:135-136`): repository text only, no gcloud or instance here.
- `explain-project`, `testing-strategy` and `operate-gcp/references/alerts.md` not re-read (no round-2 fix touched them); `spec` and `research` only by the lines named; `skills/init-project` templates other than README, guide and gcp.md Rollback row not re-read.
- The Go blocks were read for prompt content only, not compiled; no behavioural run on either CLI.

## Outcome (round 3)
- R3-F1 [medium] terraform comments omit the issuer — fixed at `resources-gcp.md` (both comments; H11's terraform half had not been applied in round 2 — the main agent's miss, caught by the round-3 verifier's grep) (review cap reached; fixed in the same round and re-checked by Phase 9's lenses, not by a fourth round here)
- R3-F2 [medium, from the round-2 fix] `postgres` password never set; the person's grant assigned to a login nobody uses — fixed (§ Local: `set-password postgres --prompt-for-password`, a proxy session without `--auto-iam-authn`, `SET ROLE` to the migrate user for the grant)
- R3-F3 [medium] the migrate job's connector path has no sample — fixed (see the code lens Outcome)
- R3-F4 [medium-low] the job block, the migrate SA's accessor binding, `--args migrate,up` — fixed (IAM-path paragraph; the second accessor binding; `--args migrate`)
- R3-F5 [low] `DialFunc = dialer.Dial` does not type-check — fixed at `service.md` item 2 (the wrapper named; `connect.md` has it)
- flags: `GRANT <migrate> TO postgres` — fixed; assertion wording — fixed ("at that point"); the README secret rule stays
