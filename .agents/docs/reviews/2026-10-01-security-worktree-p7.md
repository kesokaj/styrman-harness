---
owner: security
date: 2026-10-01
---
# Security review — worktree-p7   (2026-10-01)

Range reviewed: working tree against HEAD 677e8ee, round 3 (the last) — the same files as rounds 1–2 (`skills/api-auth-gcp/**`, `skills/data-gcp/**`, `skills/operate-gcp/**`, `skills/init-project/assets/agents-tree/docs/deploy/gcp.md`, `skills/init-project/assets/{README,REQUIREMENTS}.md`, `skills/requirements/SKILL.md`, `agents/documentation.md`, `agents/verifier.md`) plus `skills/golang/references/service.md` item 3 and `skills/terraform/references/resources-gcp.md` (§ Cloud SQL and § Job sections); `git diff HEAD -- <files>` and each file whole, every path read end to end or named under Not examined.

## Status: COMPLETE

## Verdict: FINDINGS

The round-2 medium is fixed in substance: the REVOKE is gone, and the test the caller asked for passes in full on PostgreSQL 18.6 when step 1 runs as a superuser. What is broken now is step 1 itself, on the platform the skill targets:
- **Step 1 needs a privilege Cloud SQL's `postgres` is not shown to have.** `GRANT "<migrate>" TO postgres` requires ADMIN OPTION on a role the platform created. With `postgres` modelled as Cloud SQL describes it (a non-superuser `cloudsqlsuperuser` member with CREATEROLE), it fails: `permission denied to grant role … Only roles with the ADMIN option … may grant this role`. The migration then cannot create a table. It fails closed, but whether Cloud SQL grants that ADMIN OPTION has not been observed.
- **Once step 1 does complete, the step-3 assertion cannot see the premise it backs.** With `public` owned by the migrate user, a runtime user that is a `cloudsqlsuperuser` member passes 000002. It then created a schema and a table, and could `SET ROLE` to a CREATEROLE role (measured). The assertion also runs once, not "ever".

One measured change fixes both: step 1 becomes `GRANT CREATE ON SCHEMA public TO <migrate>`.

The fixes also introduced two defects:
- A comment after a line-continuation backslash at `data-gcp/SKILL.md:97` breaks the job deploy command.
- "No password anywhere" is now false: the IAM path needs the `postgres` password.

Two older items stay open in part:
- `api-auth-gcp`'s "the one public shape".
- The `deploy/gcp.md` template's role-model row.

All other claimed fixes hold. govulncheck: no new advisory. No medium or high finding remains.

## Round 3

### The requested PostgreSQL 18 test, measured

A throwaway `postgres:18` container (18.6, `--network none`, trust auth) was stopped afterwards. The SQL was extracted verbatim from the current `data-gcp/SKILL.md` fences by awk; only `<app>` and `<project>` were substituted. Scripts and full output: `.agents/scratch/security-p7-r3/model{,2}.{sh,out}`.

Two LOGIN roles (`a-migrate@p.iam`, `a-run@p.iam`). Block 1 ran as the superuser; blocks 2–3 and `CREATE TABLE t(x int)` ran as the migrate role. The result was as the caller expected:

| Check | Expected | Measured |
|---|---|---|
| Block 1 (`GRANT … TO postgres; ALTER SCHEMA public OWNER TO …`) as a superuser | runs | exit 0 |
| 000001, 000002, `CREATE TABLE t` as migrate | run | exit 0 each |
| `has_table_privilege(run,'t','INSERT')` | true | `t` |
| `has_schema_privilege(run,'public','CREATE')` | false | `f` (`public` owner: the migrate role) |
| run: `CREATE TABLE u(x int)` | fails | `ERROR: permission denied for schema public` |
| run: `INSERT INTO t` / `DROP TABLE t` / `TRUNCATE t` | ok / fail / fail | exit 0 / `must be owner of table t` / `permission denied for table t` |
| 000002 after `GRANT CREATE ON SCHEMA public TO run` (must fail) | raises | `ERROR: runtime role has CREATE on schema public: the role split is broken` |

Three further arrangements:
- **Cloud SQL's `postgres` as a non-superuser member** (model B): block 1 fails at the GRANT.
- **A hypothesised platform grant** (model C: the IAM user granted TO `cloudsqlsuperuser` WITH ADMIN OPTION): block 1 completes.
- **The stray-grant case**: see R3-2.

`model.out` A18 and A20 are void. My own placeholder substitution corrupted `EXISTS`, so they are syntax errors and no evidence. They were re-run correctly as `model2.out` A5, A8, A10 and A12.

### Claimed fixes, verified (FAIL first)

1. **FAIL on stock PostgreSQL 18 when `postgres` is not a superuser; PASS as a superuser** — `data-gcp` § Roles step 1 (`:123-128`; the GRANT line was added during this round, see Not examined). This is R3-1 and R3-2. The REVOKE itself is gone (`:146-147`).
2. **PASS** — § Connect item 1 (`:60-67`):
   - `users create app` with no fallback and no lockout flag, its comment naming `cloudsqlsuperuser`;
   - `set-password --prompt-for-password`;
   - `secrets create`;
   - `read -rs DSN && printf … --data-file=-`;
   - `secretAccessor` on the one secret for `RUN_SA` and `MIGRATE_SA` (`:65-66`), not project-wide.
   - `DSN` is never unset, but it is not exported, so no child process sees it. Trivial; not a finding.
3. **PASS, with one defect the fix introduced** — § Connect item 2 (`:76-81`), § Migrations (`:104-121`), § Local (`:153-170`):
   - Both IAM users; Google's default quoted. I did not re-fetch it; see Not examined.
   - The IAM-path job uses `--set-env-vars` with a DSN that carries no secret, and wires the connector to golang-migrate with `stdlib.OpenDBFromPool`. This closes round 2's "show the IAM path's job dial".
   - § Local: `set-password postgres --prompt-for-password`, a proxy without `--auto-iam-authn` for that session, psql prompting, and a SELECT-only grant.
   - Checksum source and check named (`:168-170`).
   - The defect is `:97`, R3-4.
4. **PASS in part** — `api-auth-gcp`:
   - The issuer is in `:23`.
   - `:136-137`: a public service verifies a signature.
   - Firebase sensitive routes defined, and `email_verified` required for an email claim (`:79`). The sample stores the UID only.
   - `tick.New` refuses empty configuration (`go-samples.md:84-91`).
   - Still open: `:138-139`, "the one public shape" (R3-5).
5. **PASS** — `operate-gcp`:
   - Step 2 reads `changes/`, puts a security fix, SA change or rotation to the operator, and keeps a rotated version disabled (`:26-28`).
   - `--format` with four fields, and "never a raw entry" (`:60-62`).
   - Notes carry counts (`:82`).
   - Template `gcp.md:75` likewise.
   - Its sibling `gcp.md:54-58` was missed (R3-7).
6. **PASS** — the README template's secret row (`README.md:28`) and `requirements` gap 7 (`requirements/SKILL.md:38`). The role-model part of R2-5 is still open (R3-6).
7. **PASS** — `resources-gcp.md:140`: one IAM user per service account, grants per `data-gcp` § Roles. Its "no password anywhere" is R3-3.

### Findings (round 3)

- **[low / confidence: medium]** `skills/data-gcp/SKILL.md:124-127`, with `:158-163` and `skills/golang/references/service.md:316`: step 1 depends on a privilege `postgres` is not shown to have.
  - Confidence is high for the PostgreSQL 18 behaviour, which was measured. It is medium for Cloud SQL, whose ADMIN-OPTION grants were not observed.
  - **Why it fails.** On PostgreSQL 16 and later, `GRANT "<migrate>" TO postgres` needs ADMIN OPTION on the migrate role. Cloud SQL's `postgres` is a non-superuser `cloudsqlsuperuser` member with CREATEROLE, and the IAM user is created by the platform, not by `postgres`. Without that membership, `ALTER SCHEMA … OWNER TO` needs `SET ROLE` to the new owner, which is the skill's own comment.
  - **Model B, `model.out` B1–B4:**
    - B1: `ERROR: permission denied to grant role "b-migrate@p.iam" — Only roles with the ADMIN option on role … may grant this role`.
    - 000001 then gives `WARNING: no privileges were granted for "public"`.
    - The first `CREATE TABLE`: `permission denied for schema public`.
  - **Model C** (a hypothesised platform grant of the IAM user TO `cloudsqlsuperuser` WITH ADMIN OPTION): step 1 completes. The recipe therefore works only under a platform grant nobody has observed. § Local's `SET ROLE "<migrate>"` for a person's SELECT grant (`:161`) depends on the same membership.
  - It fails closed. The risk is the fix improvised under pressure, most likely a `cloudsqlsuperuser` grant.
  - **Remediation:**
    - Step 1 becomes `GRANT CREATE ON SCHEMA public TO "<app>-migrate@<project>.iam";`, run as `postgres`. Measured in `model2.out` B1–B6:
      - it works for a non-superuser member of the database owner, with no role membership;
      - DML works and the runtime user cannot create;
      - `public` stays with `pg_database_owner`, which also closes R3-2.
    - Drop `GRANT USAGE ON SCHEMA public` from 000001. PUBLIC already holds USAGE, and under this step 1 the line is a WARNING no-op.
    - The person's SELECT then needs another path, for example a NOLOGIN reader role that `postgres` creates (its creator holds ADMIN on it), given SELECT by default privileges in 000001. Not measured.
    - Alternatively, measure on a real POSTGRES_18 instance first: `SELECT member::regrole, admin_option FROM pg_auth_members WHERE roleid = '"<app>-migrate@<project>.iam"'::regrole`.
- **[low / confidence: high (measured on PostgreSQL 18)]** `skills/data-gcp/SKILL.md:137-144` (step 3), restated at `service.md:316`: once step 1 has moved `public` to the migrate user, the assertion does not test the premise it backs, that IAM users are not `cloudsqlsuperuser` members.
  - **Measured** (`model.out` A13–A17; `model2.out` A6–A7), after `GRANT cloudsqlsuperuser TO run`:
    - 000002 as written passes;
    - the runtime role has CREATE on the database (`run_create_db = t`);
    - it created `SCHEMA s2` and a table in it;
    - it can `SET ROLE cloudsqlsuperuser` (CREATEROLE).
  - **It checks once.** golang-migrate applies a version once, so 000002 checks the moment of the first deploy. `data-gcp:138` now says "at that point"; `service.md:316` still says "fails if it ever holds CREATE".
  - **Remediation:** either is enough, and (a) is shorter.
    - (a) R3-1's step 1. 000002 as written then raises on the membership (`model2.out` B8) and passes clean (B6).
    - (b) Widen the predicate to `has_schema_privilege(run,'public','CREATE') OR has_database_privilege(run, current_database(), 'CREATE') OR EXISTS (SELECT 1 FROM pg_roles r WHERE r.rolname = 'cloudsqlsuperuser' AND pg_has_role(run, r.oid, 'MEMBER'))`. Validated in `model2.out`: it passes clean (A5, A12), and raises with membership (A8) and with CREATE on the database (A10).
    - Either way, `service.md` says "when it first runs".
- **[low / confidence: high]** `skills/data-gcp/SKILL.md:114`, `skills/data-gcp/references/connect.md:108` and `skills/terraform/references/resources-gcp.md:140`: three places say there is no password on the IAM path: "no password exists anywhere", "No database password anywhere", and "with no password anywhere".
  - The IAM path now needs the built-in `postgres` password (`:159`), for step 1 and for every person grant. After step 1 it is a standing credential that, through the step-1 membership, holds the privileges of the owner of every table.
  - Remediation: "no application password; the `postgres` password is set at the prompt for the one-time owner step and person grants, and kept in the operator's password manager". Name it under the template's § Operator steps (R3-6).
- **[low (functional, introduced by a fix) / confidence: high (bash parse measured)]** `skills/data-gcp/SKILL.md:97`: `--args migrate --service-account "${MIGRATE_SA}" \   # the service binary's …` ends the command.
  - The backslash escapes a space rather than the newline, and `#` starts a comment.
  - The job deploy gets a stray `" "` argument and none of `--set-cloudsql-instances`, `--set-secrets`, `--max-retries 0` or `--task-timeout`. The next line runs as its own command. Measured with `echo` in `cont.sh`: `--second-line-flag: command not found`, exit 127.
  - It fails closed.
  - Remediation: move the comment above the command. A grep of the scoped skills, plus `cloud-run-deploy` and the template, for `\` followed by `#` finds only this line.
- **[low / confidence: high]** `skills/api-auth-gcp/SKILL.md:138-139`: "§ A scheduler route above is the one public shape here" survives. The Firebase row is a second public shape (round-1 low 2 / R2-2, the remaining part).
  - Remediation: "§ A scheduler route and § Firebase are the public shapes here; the scheduler route holds only with all four checks".
- **[low / confidence: medium]** R2-5, still open in part: no role-model row.
  - **Where it is missing:** `skills/init-project/assets/agents-tree/docs/deploy/gcp.md:33-40` has no row for the role model that `data-gcp:57,190-192` says is recorded there.
  - **§ Operator steps** (`:66`) names no database step: the `app` and `postgres` passwords, the DSN secret and its rotation, and step 1.
  - **The question is never asked.** `requirements/SKILL.md:38` asks about data and access, but not one role versus the owner/runtime split.
  - Remediation as in round 2: a `Database roles` row, and the database steps under § Operator steps.
- **[low / confidence: medium]** `skills/init-project/assets/agents-tree/docs/deploy/gcp.md:54-58`: § Rollback keeps the bare `update-traffic` command. `/preflight` deploys by this guide. The caveat (read `changes/` between GOOD and serving; a security fix, SA change or secret rotation goes to the operator) is only in § Operate's row (`:75`).
  - Remediation: one sentence in § Rollback pointing at that row.

Notes passed on for the code lens, not security:
- 000001 and 000002 name Cloud SQL IAM roles that the local compose database does not have, so `migrate up` against compose fails with `role … does not exist`.
- `service.md:366` uses golang-migrate's `database/postgres` driver package, while `data-gcp:90,118` names `database/pgx/v5`.

### Round-1 and round-2 lows, final status

| Item | Status | Line now |
|---|---|---|
| R1-1 Firebase sensitive routes, `email_verified` | fixed | `api-auth-gcp/SKILL.md:79`; the sample keeps the UID only (`go-samples.md:223`) |
| R1-2 Never list vs the Firebase row | still open in part | `:136-137` fixed; `:138-139` "the one public shape" (R3-5) |
| R1-3 tick 500 per request on empty config | fixed | `go-samples.md:84-91` |
| R1-4 rollback across a security fix, SA change or rotation | fixed, one stale copy | `operate-gcp/SKILL.md:26-28`, template `gcp.md:75`; `gcp.md:54-58` (R3-7) |
| R1-5 raw log entries in notes | fixed | `operate-gcp/SKILL.md:60-62,82` |
| R1-6 README secret rows | fixed | `init-project/assets/README.md:28` |
| R1-7 requirements data and access | fixed | `requirements/SKILL.md:38` |
| R1-8 Auth Proxy checksum | fixed | `data-gcp/SKILL.md:168-170` |
| R2-1 REVOKE | fixed (no REVOKE); its replacement is R3-1 and R3-2 | `data-gcp/SKILL.md:123-147` |
| R2-2 stale "three checks", issuer omissions | fixed, except "one public shape" | `:23`, `:139-140`, `go-samples.md:70-73`; `:138-139` (R3-5) |
| R2-3 item-1 operator block | fixed | `data-gcp/SKILL.md:60-67` |
| R2-4 `resources-gcp.md:140` third model | fixed (one IAM user per SA), apart from "no password anywhere" | `:140` (R3-3) |
| R2-5 role decision asked nowhere | still open in part | template `gcp.md:33-40,66`; `requirements/SKILL.md:38` (R3-6) |
| R2-6 "a GRANT by the schema owner" | fixed: SELECT only, as `postgres` via `SET ROLE` (which depends on R3-1) | `data-gcp/SKILL.md:161-163` |

None is "not applicable".

### Checked (round 3)

- **Secrets:**
  - A grep of all 14 scoped files for PEM private-key headers, Google API keys, GitHub, Stripe and Slack tokens, `"private_key"`, `PASSWORD '<literal>` and `--password=<value>` exited 1 (no hit). The needle was validated first on a planted PEM header in scratch (exit 0).
  - Path 1: the password goes only through prompts and stdin. Path 2: no application password; the `postgres` password goes through a prompt (R3-3).
- **Injection:**
  - The migration SQL names roles as quoted literals, with no interpolation. `has_schema_privilege` takes the name literally.
  - The Logging trace filter is unchanged and read-only.
- **Authn/authz:**
  - tick, in order: bearer present, `Validate`, the issuer in both spellings, `email` equal to the SA with `email_verified`, then `run`. It refuses empty configuration at construction.
  - The IAP middleware is unchanged (empty audience refused; `Issuer` compared).
  - Firebase stores the UID only.
- **Crypto:** nothing hand-rolled. Only claim comparisons after `idtoken`/Admin SDK verification.
- **Dependencies:**
  - `go.mod`/`go.sum` are unchanged (empty `git diff --stat`).
  - `govulncheck ./...` (go1.27.1) exited 0: 0 vulnerabilities affecting the code, and 1 in a required module that is not called. No new advisory. Output: `.agents/scratch/security-p7-r3/govulncheck.txt`.
  - `npm audit`: not applicable; no `package.json` is tracked.
- **GCP surface:**
  - `secretAccessor` is bound on the one secret, never the project.
  - `cloudsql.client`/`instanceUser` are project-scoped (the Cloud SQL norm).
  - The IAM-path job's env var carries no secret.
  - No `--allow-unauthenticated` change; no key files.
- **Agents:** `documentation.md` (+6/−5) and `verifier.md` (+2/−2) change lane text only. A grep of the diff for `name|description|tools|model|maxTurns|permissionMode|disallowedTools|color` exited 1. CLEAR.

### Not examined (round 3)

- **The file moved under the review.** `skills/data-gcp/SKILL.md` changed on disk twice during this round:
  - first the step-1 GRANT, the second `secretAccessor` binding, the IAM-path job paragraph and the § Local `postgres` login;
  - then the step-3 comment, which now says "at that point".
  - The findings are against the version hashed in `.agents/scratch/security-p7-r3/sums-1.txt`. Any later edit is unreviewed.
- **No real Cloud SQL instance.** Model B is my reading of Cloud SQL's documented `postgres` attributes, not an observation. Model C is a hypothesis. Not observed:
  - the ADMIN OPTION or memberships the platform gives API-created IAM users and `postgres` on POSTGRES_18;
  - who owns an API-created database.
- **Google's quoted default** ("granted no privileges on any databases, by default") was not re-fetched; no network was used.
- **gcloud was not run.** Not checked:
  - `--set-env-vars` parsing a value that holds `=` and spaces (the path-2 DSN);
  - `secrets versions add` on empty stdin.
- **From memory, not re-read:**
  - golang-migrate applying a version once, and sending a file as one query;
  - `stdlib.OpenDBFromPool`;
  - the `pgxmigrate.WithInstance` signature.
  - The Go samples were not compiled.
- **Read in part, as scoped:**
  - `service.md`: § Database only (`:310-401`);
  - `resources-gcp.md`: § Cloud SQL (`:105-143`) and § Job and Scheduler (`:236-328`) only.
- **No behavioural run** of an agent following these skills.
- **Outside the named scope:**
  - `context/AGENTS.md`, `README.md`;
  - `skills/{explain-project,spec,testing-strategy,cloud-run-deploy,feature,research,review,init-project}/SKILL.md`, the other new init-project assets, `skills/golang/references/testing.md`;
  - `.agents/docs/testing/skill-matrix.md`, `.agents/docs/bugs/INDEX.md`.
  - `cloud-run-deploy` and the template were grepped only for the backslash-comment shape.

## Round 2

Round-2 range: working tree against HEAD 677e8ee, round 2 — the same files as round 1 (`skills/api-auth-gcp/**`, `skills/data-gcp/**`, `skills/operate-gcp/**`, `skills/init-project/assets/agents-tree/docs/deploy/gcp.md`, `skills/init-project/assets/{README,REQUIREMENTS}.md`, `skills/requirements/SKILL.md`, `agents/documentation.md`, `agents/verifier.md`) plus `skills/golang/references/service.md` item 3 and `skills/terraform/references/resources-gcp.md` § Job (they restate the role recipe); `git diff HEAD -- <files>` and each file whole.

### Verdict (round 2): FINDINGS

The claimed fix for the worst round-1 finding does not work on the database version the skill pins.
- **What the recipe says.** `data-gcp` § Roles, restated in `golang/references/service.md` item 3, takes the runtime IAM user "down to DML" with `REVOKE cloudsqlsuperuser FROM "<app>-run@<project>.iam"`, run by the migrate user.
- **Why that fails.** Since PostgreSQL 16, only a role with ADMIN OPTION on `cloudsqlsuperuser` may revoke it, and only memberships that role granted itself.
- **What a local postgres:18.6 model showed.** The statement fails in one of two ways:
  - A plain member: an ERROR. The migration fails closed, and the obvious way to unblock it is to delete the line.
  - A member with ADMIN OPTION: the statement "succeeds" as a no-op, with a WARNING that pgx drops. The runtime user keeps CREATE on `public`, which is the DDL the split exists to remove.
- **It contradicts itself.** `:75` says a new IAM user has no privileges; `:103` says the API puts it in `cloudsqlsuperuser`. In the no-privilege reading, the migrate user cannot even create a table.

The other claimed fixes hold:
- The password path no longer puts the password in an argument list. `set-password --prompt-for-password` exists with echo disabled (gcloud 584 source).
- § Local uses the full email and no `PGPASSWORD`.
- The scheduler handler compares the issuer.
- The gcp.md template has § Operator steps.

The fixes left two stale restatements: "all three checks" in `api-auth-gcp`, and `resources-gcp.md:140`. All eight round-1 lows are still open; one sub-point is fixed. govulncheck reports no new advisory.

### The two questions

- **Is `REVOKE cloudsqlsuperuser FROM <role>` effective for an API-created user on Cloud SQL?** Not from another member, on PostgreSQL 16 and later (POSTGRES_18 here).
  - The premise "membership revocable by another member" is the pre-16 model, in which CREATEROLE sufficed.
  - PostgreSQL 16 added two things: revoking requires ADMIN OPTION on the role, and grantors are tracked. A plain REVOKE removes only the grants the current role made.
  - The model ran on local postgres:18.6 (`--network none`), with `cloudsqlsuperuser` NOLOGIN CREATEROLE CREATEDB owning database `app`, and both users made members by the superuser, as the platform does. Results:
    - A member with CREATEROLE but no ADMIN OPTION got: `ERROR: permission denied to revoke role "cloudsqlsuperuser" — Only roles with the ADMIN option on role "cloudsqlsuperuser" may revoke this role.`
    - A member with ADMIN OPTION got: `REVOKE ROLE` plus `WARNING: role "b_run" has not been granted membership in role "cloudsqlsuperuser" by role "b_migrate"`. The membership stayed in `pg_auth_members` afterwards.
    - `… GRANTED BY postgres` got: `ERROR: permission denied to revoke privileges granted by role "postgres"`.
  - Only the platform's own superuser could remove a membership it granted.
  - Whether Cloud SQL grants these memberships with ADMIN OPTION, under which grantor, or with a patched rule is not observed (no instance here).
- **Does revoking it from the runtime user leave the migrate user able to GRANT on objects it owns?** Yes.
  - The runtime user's membership does not touch the migrate user. Ownership carries the grant option.
  - Measured in all three arrangements where the migrate user could create tables: `CREATE TABLE`, `ALTER DEFAULT PRIVILEGES FOR ROLE <migrate> … TO <run>` and `GRANT SELECT ON <table> TO <run>` all succeeded, and `has_table_privilege(run, a_t, 'INSERT')` was true through the default.
  - The one statement that depends on membership is `GRANT USAGE ON SCHEMA public`. It needs the migrate user to own `public` (through `pg_database_owner`, that is, membership in the database owner).
  - Without that membership it is a silent `WARNING: no privileges were granted for "public"`. It is redundant anyway, since PUBLIC already holds USAGE.

### Claimed fixes, verified (FAIL first)

1. **FAIL** — `data-gcp` § Roles `:98-117`, restated at `golang/references/service.md:316`: the REVOKE does not take the runtime user down (finding R2-1).
2. **PASS, with defects** — `data-gcp` § Connect item 1 `:55-64`:
   - What holds:
     - One role, `app`, with "DDL at runtime" stated as the trade-off.
     - `gcloud sql users set-password … --prompt-for-password`. The 584 source, `flags.py:220-228`, says "character echo disabled".
     - `gcloud secrets versions add … --data-file=-` on stdin.
     - No `--password=` anywhere.
   - The defects are R2-3.
3. **PASS** — `resources-gcp.md:270-271` (§ Job) now matches.
   - Its stale sibling `:140` is R2-4.
   - `service.md` item 3 restates the failing recipe, which is R2-1.
4. **PASS** — `data-gcp` § Local `:124-131`:
   - The full email.
   - `--type=CLOUD_IAM_USER`.
   - `--auto-iam-authn`.
   - No `PGPASSWORD`.
   - The round-1 `${USER_EMAIL%%@*}` sub-point is fixed.
   - The new gap at `:128` is R2-6.
5. **PASS** — `api-auth-gcp` issuer check. `go-samples.md:101-105` compares `payload.Issuer` with `https://accounts.google.com` or `accounts.google.com`. It runs after `Validate` and before the `email`/`email_verified` check, and an IAP-signed ES256 token (`iss` `https://cloud.google.com/iap`) is now refused even before the audience mismatch. `SKILL.md:124-127` says four checks. Stale copies elsewhere are R2-2.
6. **PASS** — the gcp.md template has `## Operator steps` (`:64-66`) before `## Operate` (`:68`). Its content gap is R2-5.

### Findings

- **[medium / confidence: high for the PostgreSQL 18 behaviour (measured); medium for how Cloud SQL maps onto it]** `skills/data-gcp/SKILL.md:101-115` and `:75`, restated at `skills/golang/references/service.md:316`: the owner/runtime split recipe does not do what it claims.
  - **The REVOKE.** `REVOKE cloudsqlsuperuser FROM "<app>-run@<project>.iam"`, run by the migrate user, either ERRORs or is a silent no-op on PostgreSQL 16 and later (measured above).
    - In the error case, golang-migrate sends the file as one multi-statement query, which PostgreSQL runs as one implicit transaction. The migration fails closed. The natural fix under pressure is to delete the line, after which the skill still calls the runtime user "DML".
    - In the ADMIN-OPTION case, the version is recorded as applied. pgx discards the WARNING unless `OnNotice` is set.
    - In both cases, a runtime user that is still a member keeps `has_schema_privilege(run, 'public', 'CREATE') = t`. Membership in the database owner makes it owner of `public` (`pg_database_owner`), so it can create objects and drop any table in the schema. That is the DDL the split exists to remove, now hidden behind a claim that says otherwise. Path 1 at least states it.
  - **The starting state.** `:75` says "a new IAM user has NO privileges until GRANTed"; `:103` says the API "puts them in `cloudsqlsuperuser`".
    - If `:75` is right (my recollection of Cloud SQL's IAM-user docs; not verified here), the migrate user is not a member either.
    - Measured for that case: the REVOKE ERRORs, `GRANT USAGE ON SCHEMA public` is a silent no-op, and `CREATE TABLE` fails with `permission denied for schema public`. The recipe cannot build a schema without a bootstrap grant the skill does not give.
  - **The migrate job.** The job shape at `:90-94` is the password path's: a DSN secret over the `/cloudsql` socket. How the IAM path's job authenticates (the connector under golang-migrate's `WithInstance`) is not shown, so an agent improvises, most likely a password user.
  - **Remediation:**
    1. Replace the REVOKE with an assertion that fails the migration whenever the runtime user can do DDL. Example: `DO $$ BEGIN IF pg_has_role('<run>','cloudsqlsuperuser','MEMBER') OR has_schema_privilege('<run>','public','CREATE') THEN RAISE EXCEPTION 'runtime user can do DDL'; END IF; END $$;`. A platform default, a no-op REVOKE or a stray hand grant then fails closed.
    2. Measure once, on a real PG18 instance, what an API-created `CLOUD_IAM_SERVICE_ACCOUNT` user belongs to: `SELECT roleid::regrole, grantor::regrole, admin_option FROM pg_auth_members WHERE member = '<user>'::regrole`. Make `:75` and `:103` say the measured fact.
    3. If IAM users are not members, give the one-time operator bootstrap: as `postgres` through the Auth Proxy, for example `GRANT CREATE ON SCHEMA public TO "<migrate>"`. If they are members, say that a migration cannot complete the split on PostgreSQL 16 and later, and leave path 1's stated trade-off as the honest default.
    4. Show the IAM path's job dial, or name the gap.
    5. `service.md` item 3 changes in the same edit.
- **[low / confidence: high]** `skills/api-auth-gcp/SKILL.md:139`, `:23`, and `references/go-samples.md:70-73`: the issuer fix left stale restatements, and the round-1 Never-list contradiction is still there.
  - **`:139`.** "with all three checks", where `:124` now says four.
  - **The other two.** The Choose row at `:23` and the `tick` package doc at `go-samples.md:71-72` list audience, `email` and `email_verified`, without the issuer. A reader who copies the shorter list drops the new check.
  - **The Never list.** `:137-139` still makes § A scheduler route "the one shape" where a public service with an app-side check is acceptable, which forbids the Firebase row.
  - Remediation: "No `--allow-unauthenticated` on an API whose check is a plain, shared-secret or API-key header; a public service verifies a signature with `idtoken` or the Admin SDK (§ Firebase; § A scheduler route, all four checks)". Add the issuer to `:23` and to the package doc.
- **[low / confidence: medium]** `skills/data-gcp/SKILL.md:61-63`: the item-1 operator block has three defects.
  - **The create line.** `gcloud sql users create app … --password-policy-allowed-failed-attempts=5 2>/dev/null || gcloud sql users create app …` sets a lockout count without `--password-policy-enable-failed-attempts-check` (a separate flag; gcloud 584 `surface/sql/users/create.py:43-46`).
    - As far as I can show, it enforces nothing, and `2>/dev/null ||` hides whatever the API answers on Postgres.
    - If a lockout ever did apply, five bad logins by anyone with `cloudsql.client` would lock out the service's only database user.
  - **The DSN.** It is typed on stdin, which echoes, so the password is displayed in the terminal and its scrollback.
  - **The password.** It is chosen by a human and typed twice, with no source of randomness named.
  - **Missing steps.** `gcloud secrets create` and the two `secretAccessor` bindings are not shown, so an agent completes them unguided.
  - Remediation:
    - A plain `gcloud sql users create app --instance "${APP}-db"`.
    - The password from a password manager or generator, pasted at the prompt.
    - The DSN read with `read -rs DSN` and sent with `printf '%s' "$DSN" | gcloud secrets versions add … --data-file=-`. `printf` is a builtin, so it is not on argv. Then `unset DSN`.
- **[low / confidence: high]** `skills/terraform/references/resources-gcp.md:140`: a restatement of the role and credential model that the fix's grep missed.
  - It still says to prefer IAM database authentication with a `CLOUD_IAM_SERVICE_ACCOUNT` user "for the runtime SA".
  - Meanwhile, the § Job at `:270` runs migrations as the password user `app`. That is a third arrangement: an IAM runtime and a password owner. `data-gcp` neither describes it nor makes it work (its § Roles migration is written for a migrate IAM user, and `app` cannot REVOKE, per R2-1).
  - Remediation: one sentence pointing at `data-gcp` § Connect and § Roles for the role model, keeping only the state fact (a `google_sql_user` password lands in state).
- **[low / confidence: medium]** The database role decision is now "stated", but nothing asks for it or holds it. Places:
  - `skills/init-project/assets/agents-tree/docs/deploy/gcp.md:33-40,64-66`;
  - `skills/data-gcp/SKILL.md:55-57,155-157`.
  - `data-gcp` says the "DDL at runtime" trade-off and "the role model (one role, or owner + runtime)" are recorded in `deploy/gcp.md`. The template has no row for either.
  - Its § Operator steps names login, ADC and the monthly check, not the database password and DSN-secret steps (or their rotation), which item 1 makes operator steps.
  - The global context says the operator decides every architectural choice. A trade-off an agent writes into a file nobody reads is not a decision.
  - Remediation:
    - A `Database roles` row in § Rollout and network: "one role `app` (DDL at runtime; operator's choice, <date>) | owner + runtime".
    - The two database steps named in § Operator steps.
    - `data-gcp` puts path 1 to the operator as a one-line choice before using it for a service with public routes.
- **[low / confidence: medium]** `skills/data-gcp/SKILL.md:128`: the person reading production data needs "a GRANT by the schema owner", and the text does not say which.
  - The path of least resistance is a broad grant.
  - Remediation: "`GRANT SELECT ON ALL TABLES IN SCHEMA public` to the person's user, issued by the migrate user; never `cloudsqlsuperuser`, never write privileges for reading".

### Round-1 lows re-checked

| # | Round-1 low | Status | Line now | One-sentence fix |
|---|---|---|---|---|
| 1 | Firebase: revocation reserved for undefined "sensitive" routes; email trusted without `email_verified` | still open | `api-auth-gcp/SKILL.md:78-79`; `go-samples.md:201` | Define sensitive as every state-changing or admin route, and add "authorize on UID or custom claims; an email only with `email_verified` true". |
| 2 | Never list contradicts the Firebase row | still open, and now also stale ("three checks") | `api-auth-gcp/SKILL.md:137-139` | Rewrite as in R2-2: a public service verifies a signature (`idtoken` or the Admin SDK), never a plain header. |
| 3 | Scheduler `Handler` answers 500 per request on empty config instead of refusing to start | still open | `go-samples.md:85-90` | Return `(http.Handler, error)` from a constructor that refuses an empty audience or SA, as `iapauth.New` does. |
| 4 | Rollback target chosen by health alone; a rollback restores the old revision's SA and secret versions | still open | `operate-gcp/SKILL.md:26-27,35-39`; template `gcp.md:75` | Before step 2, read the `changes/` entries between GOOD and the serving revision, put a rollback across a security fix, SA change or secret rotation to the operator, and disable a rotated secret version. |
| 5 | Raw log entries (Cloud Run's `requestUrl` with query, `remoteIp`, `userAgent`) can reach the incident note and the committed bug record | still open | `operate-gcp/SKILL.md:56-58,76-79` | The note cites the query and the counts, never raw entries; an excerpt is printed with `--format` naming only severity, timestamp, status and message. |
| 6 | README Configuration table asks for every default, with no rule for a secret-bearing setting | still open | `init-project/assets/README.md:28-32` | Add "a secret's row says `Secret Manager: <name>` and has no default; a credential is never a code default". |
| 7 | `requirements` elicitation never asks about data and access | still open | `requirements/SKILL.md:32-37` | Add a seventh gap, "Data and access": data sensitivity, where secrets live, public or IAM-gated, who reads production data and through what, each with a recommended default. |
| 8 | Auth Proxy checksum: where the hash is published and how it is checked | still open, though the releases page is now named; the `${USER_EMAIL%%@*}` sub-point is **fixed** (`:129-130`) | `data-gcp/SKILL.md:133-135` | Name the hash's source (the SHA-256 table in that release's notes) and the check (`sha256sum -c -`, its exit status read bare). |

None is "not applicable". The Outcome's note that "those in the data-gcp role text are covered by the recipe" matches no low: none of the eight was in the role text.

### Checked

- **Secrets:**
  - A grep of every scoped file for password literals, `PASSWORD '…'`, Google API keys, PEM headers, GitHub, Stripe and Slack tokens, and `"private_key"` found no hits. The grep was validated first: it matched a positive planted in scratch (exit 0).
  - Placeholders only (`password=…`, `ops@example.com`). The Terraform DSN stays write-only (`secret_data_wo`, `resources-gcp.md:88`).
- **Injection:**
  - The migration SQL names roles as literal quoted identifiers, with no interpolation.
  - The `${TRACE_ID}` Logging filter is unchanged since round 1: read-only, within the operator's own access.
- **Authn/authz:**
  - The scheduler handler runs, in order: bearer present, `Validate` (audience), issuer in both spellings, `email` equal to the scheduler SA and `email_verified`, then `run`. It fails closed at each step.
  - The IAP middleware is unchanged (constructor refuses an empty audience; `Issuer` compared).
- **Crypto:** nothing hand-rolled. The issuer and email checks are claim comparisons after `idtoken.Validate`.
- **Dependencies:**
  - `git diff HEAD --stat -- go.mod go.sum` is empty.
  - `govulncheck ./...` (go1.27.1) exited 0: 0 vulnerabilities affecting the code, and 1 in a required module that is not called (GO-2026-5024, the known AGENTS.md entry). No new advisory. Output: `.agents/scratch/security-p7-r2/govulncheck.txt`.
  - `npm audit`: not applicable; no `package.json` is tracked.
- **GCP surface:**
  - IAM database users need `roles/cloudsql.instanceUser` and `roles/cloudsql.client`, which are project-scoped, the Cloud SQL norm.
  - No `--allow-unauthenticated` change; `--set-secrets` with pinned versions; no key files.
- **Agents:** `documentation.md` (the onboarding.md lane, "you never create it") and `verifier.md` (fills a guide that is still the template) change lane text only. No `tools:`, description, model or permission change. CLEAR.
- **PostgreSQL model:** a throwaway `postgres:18` container (18.6, `--network none`, trust auth), removed afterwards. It ran plain SQL and nothing else.

### Not examined

- **No real Cloud SQL instance.** These were not observed:
  - the memberships, grantors and ADMIN OPTION the platform gives API-created built-in and IAM users on PG18;
  - who owns an API-created database;
  - any Cloud SQL patch to role administration.
  - R2-1's Cloud SQL mapping rests on stock PostgreSQL 18.6 and on my reading of the docs.
- **gcloud.** It is not on `PATH`, and the SDK at `~/google-cloud-sdk` is 584.0.0, not the 587 the skills cite.
  - Only flag definitions were read (`flags.py`, `surface/sql/users/create.py`).
  - Not checked: API behaviour, such as whether a per-user lockout applies to Postgres; and `secrets versions add` refusing empty stdin.
- **Library behaviour from memory, not re-read from module source:**
  - golang-migrate's pgx driver sending the file as one query;
  - pgx dropping notices without `OnNotice`;
  - the `idtoken` certificate selection.
- **The Go samples were not compiled.**
- **Read only in part:**
  - `service.md`: its diff only (item 3); the rest of that file was not re-read.
  - `resources-gcp.md`: the grep hits plus lines 250-290 (§ Job). § Cloud SQL beyond line 140 was not re-read.
  - `agents/documentation.md` and `agents/verifier.md`: read whole, but long lines were cut at 700 characters (verifier `:14,21,23-26`). Those lines' changed hunks were read in full through the diff.
- **No behavioural run** of an agent following these skills.
- **Outside the named scope, not read this round:**
  - `context/AGENTS.md`, `README.md`;
  - `skills/{explain-project,spec,testing-strategy,cloud-run-deploy,feature,research,review}/SKILL.md`, `skills/init-project/SKILL.md`, the other new init-project assets, `skills/golang/references/testing.md`;
  - `.agents/docs/testing/skill-matrix.md`, `.agents/docs/bugs/INDEX.md`.

## Round 1

Round-1 range: working tree against HEAD 677e8ee — Phase 7 of the audit-fix plan, three skills that tell an agent how to do authentication, data access and operations on GCP, plus templates that record deploy and runbook values: `skills/api-auth-gcp/SKILL.md` + `references/go-samples.md`, `skills/data-gcp/SKILL.md` + `references/connect.md`, `skills/operate-gcp/SKILL.md` + `references/alerts.md`, `skills/init-project/assets/agents-tree/docs/deploy/gcp.md` (§ Operate), `skills/init-project/assets/README.md`, `skills/init-project/assets/REQUIREMENTS.md`, `skills/requirements/SKILL.md`, `agents/documentation.md`, `agents/verifier.md` (new files `git add -N`, read whole).

### Verdict (round 1): FINDINGS

The worst problem is the database privilege model in `data-gcp`. By default, the internet-facing service connects as the schema owner, through one DSN it shares with the migration job (user `app`, "which then owns the schema and has DDL at runtime"). The only ways this stack names for creating that user are `gcloud sql users create` and `google_sql_user`. The skill itself says users created that way land in `cloudsqlsuperuser`. So an SQL injection in the service gets DDL, `CREATEROLE` and `CREATEDB` by default.

The split that would fix this is left as "a project decision" that nothing asks about, and it now has two incompatible recipes:
- `data-gcp` makes the runtime role `app_rw` and grants default privileges for `app_owner`, a role nothing creates.
- `golang/references/service.md:316` keeps the API-created `app` as the runtime role and calls it "DML only".

Two medium secret-handling gaps sit in the same section:
- The `CREATE ROLE … PASSWORD '…'` statement runs "as the migration job's DSN", and the only automated place to run it is a committed, embedded migration file.
- The skill never shows how the `app` password and the DSN secret are created, so the obvious completion puts the password on argv and into the agent's transcript.

The authentication skill holds up:
- The IAP check refuses an empty audience and compares `iss`.
- The public `/tick` route checks `email` and `email_verified`.
- Roles are granted on the service.
- Nothing parses a token by hand.

Its findings are low. govulncheck reports no new advisory.

### Findings

- **[medium / confidence: high]** `skills/data-gcp/SKILL.md:49-53,89-104`, with `skills/golang/references/service.md:316` and `skills/terraform/references/resources-gcp.md:140`: the default runtime database role is the privileged owner, and the split has two incompatible recipes.
  - **The default.** One DSN, user `app`, serves both the migration job and the service (`:91`).
    - The skill never shows how `app` is created. The only creation paths this stack names are `gcloud sql users create` and `google_sql_user`, and the skill itself (`:93-94`) says those land in `cloudsqlsuperuser`.
    - So the public service's SQL-injection blast radius is: DDL (DROP), new login roles (persistence), and `CREATEDB`.
  - **The split.**
    - `data-gcp` creates the runtime role `app_rw` in SQL. It grants `ALTER DEFAULT PRIVILEGES FOR ROLE app_owner`, but no step creates `app_owner`, and the owner DSN user is `app` everywhere else (`:52`, `:91`). An agent following it either gets an error or grants defaults for a role that never creates the tables. That fails closed, which pushes the agent back to the one-DSN default.
    - `service.md:316` keeps `app` as the runtime role and calls it "DML only". That is p6 round 3's open low: default privileges add grants and remove none, so an API-created `app` stays in `cloudsqlsuperuser`.
    - The two skills cannot both be followed.
  - **The credential.** `terraform`'s resources-gcp.md:140 says to prefer IAM database authentication for the runtime SA. `data-gcp` makes the password socket "the default". The two skills give opposite defaults.
  - **No one asks.** Neither `requirements` nor `data-gcp` puts the split to the operator, so "a project decision" means the default always wins.
  - **Remediation:**
    - Keep one role model in `data-gcp`; `service.md` and resources-gcp.md point at it.
    - The owner is the migration DSN's user, created and named once.
    - The runtime is an IAM database user (the connector with `WithIAMAuthN`: no password and no `cloudsqlsuperuser`; privileges only through GRANTs a migration issues) or a role created in SQL with DML grants only. It never uses the owner DSN.
    - Make the split the default for any service with a public route. A single DSN is a one-line operator choice, recorded in `deploy/gcp.md`.
- **[medium / confidence: medium]** `skills/data-gcp/SKILL.md:96-98`: `CREATE ROLE app_rw LOGIN PASSWORD '…'` is to run "as the owner role (the migration job's DSN), once".
  - The only path the skill gives for running SQL as that DSN is a migration under `//go:embed migrations/*.sql` (`:74`). A migration file is committed to git and baked into every image, so the natural reading puts a plaintext password into source. That breaks the hard rule "Never commit or hardcode secrets".
  - Run by hand, a literal `PASSWORD '…'` is still kept in `pg_stat_statements` and in the server log whenever `log_statement` covers DDL. This is from Postgres documentation memory, not tried.
  - Remediation: say where it runs. The operator runs it once in psql through the Auth Proxy: `CREATE ROLE app_rw LOGIN;`, then `\password app_rw`. psql sends a SCRAM verifier, never the plaintext. Add "never in a migration file". The IAM runtime user from the first finding avoids the password entirely.
- **[medium / confidence: medium]** `skills/data-gcp/SKILL.md:49-53,84`: the skill never shows how the `app` password and the `${APP}-db-url` secret come into being.
  - It has no user-creation line, no `gcloud secrets create`/`versions add`, and no generation step.
  - The obvious completion for an agent is `gcloud sql users create app --password=<literal>` plus `echo "host=… password=<literal>" | gcloud secrets …`. That puts the password in argv, in shell history, and in the agent transcript, which the CLI persists on disk and the nightly dream job reads.
  - Remediation, with the value never printed:
    - Generate the password in a pipe that feeds Secret Manager directly (`--data-file=-`), and assemble the DSN secret in the same pipe.
    - The operator creates the user with `gcloud sql users create app --instance … --prompt-for-password`, or Terraform takes the write-only path the `terraform` skill already uses.
    - The agent prints these commands and never runs one whose argv carries the password.
  - Confidence is medium because `--prompt-for-password` was not checked: gcloud is absent here.
- **[low / confidence: medium]** `skills/api-auth-gcp/SKILL.md:77-79`, `references/go-samples.md:195-219`: the Firebase guidance stops at the token.
  - **Revocation.** `VerifyIDTokenAndCheckRevoked` is reserved for "sensitive routes only", and "sensitive" is never defined. An agent will default `checkRevoked=false`, so a disabled or revoked user keeps access, writes included, until the ID token expires (up to an hour).
  - **Email trust.** Nothing forbids authorizing on `token.Claims["email"]` without `email_verified`. With the email/password provider anyone can register an unverified `admin@<company>` address. The skill does make UID the user key, which is right. But the common admin allowlist keyed by email is not ruled out.
  - Remediation:
    - Define sensitive as any state-changing or admin route.
    - Add one line: "authorize on UID or custom claims; an email is trusted only with `email_verified` true".
  - Note: the project binding is implicit. `firebase.NewApp(ctx, nil)` takes the project from ADC. If that is the wrong project, the check fails closed, so this is not a weakness.
- **[low / confidence: high]** `skills/api-auth-gcp/SKILL.md:136-138`: the Never list contradicts the skill's own Firebase row.
  - The list says § A scheduler route is "the one shape" where a public service with an app-side check is acceptable. That forbids the Firebase row, which is a public service checking `VerifyIDToken`.
  - An agent reconciling the two loses the distinction the bullet exists to draw: a signature check versus a plain header check.
  - Remediation: "only when the check is a signature verification by `idtoken` or the Admin SDK (the Firebase row, § A scheduler route); never a plain, shared-secret or API-key header".
- **[low / confidence: high]** `skills/api-auth-gcp/references/go-samples.md:85-90`: the scheduler `Handler` checks for an empty audience or SA on every request and answers 500.
  - The skill's own rule (`SKILL.md:139`) is "refuse to start". The sample fails closed, but a misconfigured deploy passes its health checks and fails only at the first tick.
  - Remediation: a constructor that returns an error, like `iapauth.New`.
  - The missing `iss` check on this route is acceptable: an RS256 token verifies against Google's OAuth certs, and the email comparison binds the caller.
- **[low / confidence: medium]** `skills/operate-gcp/SKILL.md:21-39`, and the template `deploy/gcp.md` § Operate Rollback row: the rollback target is chosen by health alone.
  - On Cloud Run, the service account, secret versions and env vars belong to the revision. `update-traffic --to-revisions GOOD=100` therefore restores the old revision's identity and secrets.
  - So a release that narrowed the runtime SA, moved off a leaked secret version, or fixed an auth check is silently undone.
  - The command itself cannot widen invoker IAM or ingress (both service-level), and `--allow-unauthenticated` appears in no rollback. That part of the brief is CLEAR.
  - Remediation:
    - Before step 2, read the `changes/` entries between GOOD and the serving revision. When one is a security fix, an SA change or a secret rotation, put the rollback to the operator.
    - After a rotation, disable the replaced secret version so that a rollback fails closed rather than running on the old version.
- **[low / confidence: medium]** `skills/operate-gcp/SKILL.md:56-58,74-82`: no rule stops raw log output from being pasted into the incident note and from there into the committed bug record.
  - Cloud Run's own request log records `httpRequest.requestUrl` with the full query string, plus `remoteIp`, `referer` and `userAgent`, for every request, whatever the app logs.
  - The skill's `gcloud logging read` rows print whole entries. "Never log a query string" (`:56`) reads as if the app controlled this.
  - Remediation:
    - The note cites the query and the counts, never raw entries.
    - An excerpt is printed with `--format` naming only severity, timestamp, status and message.
    - Name the real control: no secret ever travels in a URL.
- **[low / confidence: medium]** `skills/init-project/assets/README.md:26-32`: the Configuration table asks for every flag and env var "with their defaults", and has no rule for a secret-bearing setting.
  - `init-project` Step 6 fills it from `os.Getenv` grep hits, and `explain-project` Step 2.2 asks for defaults the same way. An agent could therefore copy the compose DSN `postgres://app:app@…` (`connect.md:15`), or any credential default it finds, into a committed file.
  - The template itself has no secret-shaped placeholder. That part of the brief is CLEAR.
  - Remediation: one line, "a secret's row says `Secret Manager: <name>` and has no default; a credential is never a code default".
- **[low / confidence: medium]** `skills/requirements/SKILL.md:30-37`: the elicitation does not ask about data and access.
  - It asks who may use the system, but not:
    - what data the project holds and how sensitive it is;
    - where its secrets live;
    - whether the service is public;
    - who may read production data, and through what.
  - Those are exactly the decisions `api-auth-gcp` (public or IAM-gated) and `data-gcp` (the role split; the Auth Proxy "for reading production data") leave to the project.
  - Remediation: a seventh gap, "Data and access", each item with a recommended default.
  - The `REQUIREMENTS.md` template itself weakens no hard rule: Non-goals and Known limits are empty placeholders, and an item the operator approves is the operator's call.
- **[low / confidence: low]** `skills/data-gcp/SKILL.md:114-121`: "verify the published `sha256` beside it" does not say where the checksum is published.
  - From memory, the proxy binaries are served from `storage.googleapis.com/cloud-sql-connectors/cloud-sql-proxy/vX.Y.Z/`, and the SHA-256 table is in the GitHub release notes, not in a file beside each asset. An agent that finds no file may skip the check, or hash the download against itself.
  - Remediation: name the URL pattern, the source of the hash, and the command (`sha256sum -c -`, read bare).
  - Correctness only, fails closed: `user=${USER_EMAIL%%@*}` (`:116`) is wrong for a human IAM database user, whose name is the full email address.

### Checked

- **Authentication (api-auth-gcp, CLEAR apart from the lows above):**
  - IAP: `idtoken.Validate` with an audience refused when empty at construction (`go-samples.md:37-40`), plus an explicit `payload.Issuer` comparison (`:49`).
  - The never-log list carries `x-goog-iap-jwt-assertion` (`golang/references/service.md:172`, `operate-gcp/SKILL.md:57`, `api-auth-gcp/SKILL.md:63-64`). `X-Goog-Authenticated-User-Email` is explicitly not trusted.
  - Roles are granted on the service, never the project: `iap web add-iam-policy-binding --resource-type=cloud-run --service`, and `run.invoker` on the service for the IAP agent and for callers.
  - `/tick` on a public service: Validate, then `email` == the scheduler SA, then `email_verified`. This closes p6 round 3's medium.
  - Service-to-service: `idtoken.NewClient` with a non-empty audience, and the callee stays `--no-allow-unauthenticated`.
  - No service-account key files.
  - No hand-rolled JWT anywhere. A grep for `jwt.Parse`, base64 or token splitting finds only the Never bullet, and every sample was read line by line.
- **Data (data-gcp):**
  - DSN through `--set-secrets` with a pinned `${VERSION}`, never `--set-env-vars`; never logged (`connect.md:26-30`; pgx v5 redacts `password=` in parse errors, from memory).
  - `cloudsqlconn` shapes: `WithIAMAuthN` + `WithLazyRefresh`; the username is the SA email without `.gserviceaccount.com`.
  - Instance: `--ssl-mode=ENCRYPTED_ONLY`, no authorized networks, PITR and 7 retained backups, `--deletion-protection`.
  - The migrate job runs under its own SA with `--max-retries 0`.
  - `roles/cloudsql.client`/`instanceUser` are project-scoped. That is the Cloud SQL norm; an IAM condition on the instance would narrow it (no finding).
- **Operations (operate-gcp, alerts.md):**
  - The filters select on severity, status, trace and revision. Cloud Run logs no bodies and no request headers.
  - The 5xx policy needs more than five a minute sustained for five minutes, so it does not page on noise.
  - Rollback is `update-traffic` only, and `--remove-tags canary` closes the tagged URL.
  - The beta surfaces are named as beta: `gcloud beta run services logs tail`, `gcloud beta monitoring channels create`.
  - `policy.json` carries no secret.
  - The trace filter interpolates `${TRACE_ID}` inside double quotes. That is a read-only Logging query within the operator's own access, so no shell injection and no privilege gain.
- **Templates:**
  - `deploy/gcp.md` § Operate rows hold commands, names and ids, never values. § Secrets says "names, not the values", and Recorded defaults holds no secret.
  - The README template holds no secret-shaped placeholder.
  - The `REQUIREMENTS.md` template's Non-goals and Known limits are empty and weaken no hard rule.
- **Agents:** `documentation.md` (onboarding.md lane) and `verifier.md` (fill the template guide) edit lanes only. No `tools:`, description or permission change.
- **Secrets:** a grep of the in-scope files for key, token and password shapes finds only placeholders (`ada@example.com`, `ops@example.com`, `password=…`, `'…'`) and the local compose DSN `app:app`, a development default.
- **Dependencies:**
  - `go.mod`/`go.sum` are unchanged. The libraries the samples name for target projects are justified: idtoken, the Firebase Admin SDK, cloudsqlconn, pgx and golang-migrate. Stdlib has no OIDC/JWT verifier and no Postgres driver, and hand-rolling is banned.
  - `govulncheck ./...` (go1.27.1, exit 0) found 0 vulnerabilities affecting the code. One is in a required module and is not called: GO-2026-5024, `golang.org/x/sys@v0.37.0`, fixed in v0.44.0, the known entry in AGENTS.md. No new advisory. Output is in `.agents/scratch/security-p7/`.
  - `npm audit`: not applicable. No `package.json` is tracked.
- **GCP surface:**
  - `--allow-unauthenticated` appears only for the public-service case and the Never list.
  - Secrets go through Secret Manager with pinned versions.
  - No key files.
  - Every background shape uses a dedicated SA, never the compute default.

### Not examined

- `gcloud` is not installed on this box, so none of the gcloud claims were run:
  - `--iap` on `gcloud run deploy` as GA;
  - `gcloud iap web add-iam-policy-binding --resource-type=cloud-run`;
  - the IAP audience format `/projects/N/locations/R/services/S`;
  - `gcloud sql users create --prompt-for-password`;
  - `channels create` being beta-only;
  - `logs read --log-filter`/`--freshness`.
  - The skills say they were checked against gcloud 587, elsewhere.
- The Go samples were not compiled by me (the skill claims `go vet`).
- These library behaviours come from memory and from p6's reading of `validate.go`, not re-read from module source:
  - idtoken ES256 against the IAP certs, and `aud`/`iss`;
  - what Firebase `VerifyIDToken` checks;
  - pgx's DSN redaction.
- Partially read, outside the named scope: three new files were read with lines cut at 600 characters.
  - `skills/explain-project/SKILL.md` lines 3 and 12;
  - `skills/spec/SKILL.md` lines 3 and 29;
  - `skills/testing-strategy/SKILL.md` lines 3 and 52.
- Read as diffs only, outside the named scope:
  - `skills/init-project/SKILL.md` and `context/AGENTS.md`, with long lines cut at 900 characters;
  - `README.md`.
- Not read: the `.agents/docs/testing/skill-matrix.md` diff. The output was persisted and never read.
- Blast radius read only through grep or by section:
  - `skills/cloud-run-deploy/SKILL.md` § Rollback and flags;
  - `skills/golang/references/service.md` lines 172 and 316;
  - `skills/terraform/references/resources-gcp.md` § Cloud SQL (lines 105-142).
- No behavioural run of an agent following these skills.

### Outcome
- [medium] runtime role is the schema owner; two incompatible split recipes — fixed at `skills/data-gcp/SKILL.md` § Connect item 1 and § Roles (one recipe: the password path has one role `app` with the trade-off stated; the split is the IAM-authN path, two IAM users, the first migration revokes `cloudsqlsuperuser` from the runtime user and grants DML); `golang/references/service.md` item 3 and `terraform/references/resources-gcp.md:270` agree
- [medium] `CREATE ROLE … PASSWORD` in an embedded migration — fixed (the migration names roles only; no password exists on the IAM-authN path)
- [medium] password and DSN creation never shown — fixed (§ Connect item 1: `gcloud sql users set-password --prompt-for-password` and `gcloud secrets versions add --data-file=-` on stdin, both operator steps; never `--password=`, never argv or transcript; flags checked against gcloud 587)
- lows (8): see the round-2 report for what each lens re-checks; those in the data-gcp role text are covered by the recipe above; the rest are left for round 2 to name as fixed or open

## Outcome (round 2)
- [medium] `REVOKE cloudsqlsuperuser` does not take the runtime user down to DML; `:75` vs `:103` — fixed at `skills/data-gcp/SKILL.md` § Connect item 2 and § Roles: Google's documented default (IAM users have no privileges and are not `cloudsqlsuperuser` members; fetched 2026-10-01) is the one statement; `ALTER SCHEMA public OWNER TO <migrate user>` once as `postgres`; the first migration grants DML; `000002_assert_roles` raises if the runtime user has CREATE on `public` (fail closed, as you asked); no REVOKE anywhere; `service.md` item 3 and `resources-gcp.md:140,270` agree
- [low] stale "three checks" / issuer omissions — fixed (`api-auth-gcp/SKILL.md:23,139`, `go-samples.md` comment; the Never list now says a public service verifies a signature, never a header)
- [low] item-1 operator block (`2>/dev/null ||`, lockout flag, echoed DSN, no `secrets create`/binding) — fixed (`read -rs`, `secrets create`, `add-iam-policy-binding`, no fallback, no lockout flag)
- [low] `resources-gcp.md:140` third role model — fixed (one IAM user per SA, no password, grants per data-gcp § Roles)
- [low] the role decision asked nowhere — fixed: `requirements` gap 7 "Data and access" and the README template's secret-row rule
- [low] `:128` "a GRANT by the schema owner" — fixed (the migrate user, SELECT only)
- round-1 lows: 1 Firebase sensitive routes and `email_verified` — fixed (`api-auth-gcp/SKILL.md:78-79`); 2 Never list — fixed; 3 tick handler fails to start on empty config — fixed (`tick.New` returns an error; block vets); 4 rollback across a security fix/SA change/secret rotation — fixed (`operate-gcp` step 2 and the gcp.md Rollback row); 5 raw log entries in notes — fixed (`--format` naming four fields; counts not entries); 6 README secret rows — fixed; 7 requirements data-and-access gap — fixed; 8 Auth Proxy checksum source and check — fixed

## Outcome (round 3)
- [low] `GRANT <migrate> TO postgres` needs ADMIN OPTION Cloud SQL's `postgres` may not hold — fixed: step 1 is your measured fix, `GRANT CREATE ON SCHEMA public TO <migrate>`; no ownership transfer, `public` keeps its owner (review cap reached; fixed in the same round and re-checked by Phase 9's lenses)
- [low] the assertion passes once `public` belongs to the migrate user — fixed by the same change (`public` is not handed over); "ever" → "at that point" in `service.md` item 3
- [low, introduced by my round-3 edit] comment after a continuation backslash breaks the job command — fixed (`:97`); the verifier's final run executes the block with a shadowed `gcloud`
- [low, introduced by my round-3 edit] "no password exists anywhere" — fixed (the service and the job carry none; the built-in `postgres` user keeps one for the one-time step)
- [low] "the one public shape" vs the Firebase row — fixed (`api-auth-gcp/SKILL.md:138`: the one public shape without an end-user login)
- [low] gcp.md template: no role-model row, no database operator step, no rollback caveat in § Rollback — fixed (a `Database roles` row; § Operator steps names the one-time database steps; § Rollback carries the `changes/` caveat)
