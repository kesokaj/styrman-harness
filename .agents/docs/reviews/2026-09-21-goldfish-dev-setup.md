---
owner: reviewer
date: 2026-09-21
---
# Goldfish test — context/dev-setup.md   (2026-09-21)

Range reviewed: `context/dev-setup.md` only (Part A), then compared against the repo (Part B). Audit-prompt Step 5, single-file goldfish test — not a diff review.

## Status: COMPLETE

File size: 2551 bytes → ~638 tokens (bytes ÷ 4).

## Part A — read cold (file + named skills only)

### G1 Comprehension

`context/dev-setup.md` is the "Environment & Stack" section appended to the global context at install (per `AGENTS.md`'s Environment & Stack heading and this repo's own `AGENTS.md` Layout table: "environment, stack, GCP defaults — appended to the context at install"). It tells an agent working for this operator where and how code runs: local-first by default, GCP as the only cloud and only on explicit instruction, the stack (Go backend in `cmd/<app>/`+`internal/`, TanStack Start/TypeScript/Tailwind frontend inside `web/` or `ui/`), which machines exist (macOS/Apple Silicon with Colima, or Linux with plain Docker), a root-step protocol for the operator's own Linux box, and the release pipeline order: local test → native docker build/run smoke test → `linux/amd64` build → push to Artifact Registry → Cloud Run deploy, with secrets only via Secret Manager. It names eight skills by backticked name (`golang`, `frontend`, `typescript`, `docker`, `cloud-run-deploy`, `gke-basics`, `terraform`, `gcloud-setup`) that hold the mechanics.

Must: ask for the five GCP defaults (project ID, region, Artifact Registry repo, service name, public/private) before any GCP action and record them in the project's `docs/deploy/gcp.md`; test locally and smoke-test in a native Docker build before anything goes to Google; always build `linux/amd64` for anything bound for Cloud Run; tag images by git SHA; put secrets through `--set-secrets`, never `--set-env-vars`; ask for remote-box/VM access details rather than assume them; on the operator's own Linux box, run a root step only when both the operator said so AND the project's `.claude/settings.json` records standing permission, inspecting the target first and recording the command plus verification.

Must not: never suggest a cloud provider other than GCP; never invent or assume the five GCP defaults or carry one over from another project; never fall back to a value from another project; never put a secret in `--set-env-vars`; never run a root step without both conditions being met (hand the exact command back instead); never let shipped frontend JS use Node-only APIs outside the server entry.

### G2 Critic (findings that would change agent behaviour)

1. **Docs-root ambiguity.** The file writes bare `docs/deploy/local.md` and `docs/deploy/gcp.md`, with no mention of the `.agents/docs/` vs `docs/` fallback that the `cloud-run-deploy` skill states explicitly ("read the project's `docs/deploy/gcp.md` (`.agents/docs/` or `docs/`)"). An agent reading only this file (before it loads the skill) could look in the wrong root on a project using `.agents/docs/`, or write the recorded defaults to a path the project doesn't use.
2. **Root-step protocol scoped only to "the operator's own Linux box."** No guidance for a root step needed on a *remote* Linux box or VM once access is granted (the file's separate remote-box bullet only covers "ask for access details," not what happens once you have them and a root step is needed there). An agent could either wrongly apply the standing-permission shortcut on a remote box or wrongly refuse a legitimate one — ambiguous either way.
3. **No Cloud Run service-account guidance.** Neither this file nor `cloud-run-deploy` (checked in Part A) says which identity a newly deployed Cloud Run service runs as; the deploy command shown has no `--service-account` flag, so the true default (the project's Compute Engine default service account, broad permissions) is unstated. This matters because the same file requires secrets via Secret Manager, which needs an IAM binding the deploy recipe never mentions granting.

### G3 Readiness — deploy a new Go service to Cloud Run for the first time, from a laptop

Questions needed for a correct first pass (using this file + its named skills only):
1. Is the laptop macOS or Linux? (Determines Colima vs plain-Docker path; dev-setup.md states both exist but not which applies to "this" laptop.)
2. Does `docs/deploy/gcp.md` already hold the five GCP defaults, or must I ask for all five now?
3. Is `gcloud` already installed and authenticated on this laptop (`gcloud auth list`)? (`gcloud-setup`'s install section only documents macOS/Homebrew — no Linux install path.)
4. Does the Artifact Registry repo already exist, or do I need to create it?
5. Does the Go service already have a Dockerfile satisfying the skill's rules (listens on `$PORT`, binds `0.0.0.0`, distroless multi-stage)?
6. Should the service be public (`--allow-unauthenticated`) or IAM-gated — operator's call, not assumed.
7. Does the service need secrets, and do the named Secret Manager secrets already exist?
8. What service account should the Cloud Run service run as, and does it have the IAM bindings the app needs (e.g. Secret Manager access)?
9. Is this the operator's real project checkout with a `.git` (needed for the SHA tag), and has `gcloud auth configure-docker ${REGION}-docker.pkg.dev` run once already for this machine/region?

## Part B — compared against the repo

Read: `README.md`, `context/AGENTS.md`, `.agents/docs/deploy/{local,remote,gcp}.md`, skill files above.

### G1 errors/omissions found
- No factual error in G1 against the repo. One omission: `context/AGENTS.md`'s own "Environment & Stack" heading is exactly this file's insertion point (confirmed by grepping — see below), so G1's characterization is accurate, not a gap to fix.

### G3 questions — BLOCKING / NON-BLOCKING
1. macOS vs Linux — **BLOCKING** (changes which commands run at all).
2. gcp.md populated? — **NON-BLOCKING**, already covered: `cloud-run-deploy` skill's own Pre-flight step 1 asks and writes it when missing — a fact the skill already holds, not a gap in this file.
3. gcloud installed/authenticated — **NON-BLOCKING** on macOS (skill documents it); **BLOCKING** on Linux (no install path exists anywhere in the read set — `gcloud-setup` SKILL.md has no Linux install section). This is a skill gap, not this file's, per the task's scoping rule, so no fix proposed here.
4. Artifact Registry repo exists — **NON-BLOCKING**, `cloud-run-deploy` gives the create command inline.
5. Dockerfile conformance — **NON-BLOCKING**, `docker` skill covers it fully.
6. Public vs IAM-gated — **NON-BLOCKING**, `cloud-run-deploy`'s flag table forces a deliberate choice.
7. Secrets exist in Secret Manager — **NON-BLOCKING**, mechanics covered by `--set-secrets`; existence-check is an operational step, not a knowledge gap.
8. Service-account/IAM for the deployed service — **BLOCKING**, and genuinely missing from both this file and every named skill (checked `cloud-run-deploy` in full).
9. Registry auth done once per machine/region — **NON-BLOCKING**, `cloud-run-deploy` Pre-flight step 5 covers it.

### Proposed fixes (fewest bytes, in the file a goldfish would check)
- For G2 finding 1 (docs-root ambiguity): change `context/dev-setup.md` line 3 and line 4's `docs/deploy/...` to `.agents/docs/deploy/... (or docs/)` — but this duplicates wording the `cloud-run-deploy` skill already carries verbatim; per the size-pinned budget (ADR-035) and the "a fact a skill already holds is not a gap in this file" rule, the cheaper fix is likely to leave dev-setup.md alone (the skill is loaded before any GCP write happens) — **flagged for the operator's call, not applied.**
- For G3 blocking Q1 (OS): no fix needed in this file — `context/AGENTS.md`'s Environment & Stack section (the operator-profile line, not dev-setup.md) already states "macOS (Apple Silicon) or Linux," and the file's own "Machines" bullet already branches on both; the true gap is that neither file tells the agent to *ask* which one this run is on. One line would close it economically:
  - `context/dev-setup.md` line 7, append: "; ask which one this machine is if not evident." (~13 tokens)
- For G3 blocking Q8 (service account / IAM): add to `skills/cloud-run-deploy/SKILL.md` Push & Deploy section (not this file — it's deploy mechanics, this file's job is only to point at the skill, which it already does): a line naming `--service-account` and the IAM bind step for Secret Manager access. Out of scope for `dev-setup.md` itself.

## Report

- **Tokens:** 2551 bytes ÷ 4 ≈ **638 tokens**.
- **G1 errors:** 0 factual errors; file's self-description confirmed accurate.
- **G2 findings kept:** 3 (docs-root ambiguity; root-step protocol silent on remote boxes; Cloud Run service-account/IAM silent).
- **G3 questions:** 9 total, **2 BLOCKING** (OS identity; Cloud Run service-account/IAM), 7 non-blocking (covered by named skills).
- **Verdict SUFFICIENT:** No — the two blocking gaps (which OS, and service-account/IAM for secrets access) are not answered anywhere in the file or its named skills, so a first-pass Cloud Run deploy from a laptop cannot be done correctly with only this file, even skills included.
- **Verdict MINIMAL:** No line in the file is dead weight — every bullet gates a distinct decision (cloud choice, GCP defaults, stack layout, machine type, root-step protocol, build/test order, architecture, registry, compute choice, auth) and removing any one would let a real mistake through (e.g. removing the `--set-secrets` line would let a secret land in `--set-env-vars`, which the file explicitly bans). No line is proposed for removal.

## Not examined
- `skills/golang/SKILL.md`, `skills/frontend/SKILL.md`, `skills/typescript/SKILL.md`, `skills/gke-basics/SKILL.md`, `skills/terraform/SKILL.md` bodies were not read in full (only confirmed to exist/sized); they are named by the file but not load-bearing for the Cloud-Run-deploy G3 scenario, so their omission does not change the counts above, but a goldfish test of a different task (e.g. GKE or Terraform) would need them read.
- `.agents/docs/deploy/local.md` and `.agents/docs/deploy/remote.md` were listed but not read line-by-line; only `gcp.md`'s presence and the `README.md`/`context/AGENTS.md` cross-reference were checked for Part B.
