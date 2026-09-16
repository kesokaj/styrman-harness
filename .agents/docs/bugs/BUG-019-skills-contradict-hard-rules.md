---
owner: bugs
date: 2026-09-15
id: BUG-019
status: fixed
area: skills (hard rules)
---
# BUG-019: skill texts and references instruct against the global hard rules

## Symptom
Several skill files carry guidance that contradicts `context/AGENTS.md` hard rules:
- `skills/bash/references/makefile-patterns.md`: `docker-build` (:133) runs a native `$(DOCKER) build`, `REGISTRY` (:129) is fixed to `us-docker.pkg.dev`, and both deploy targets hardcode `--region=us-central1` (:193,196). `context/dev-setup.md` says to ask for the region and never assume it (:4), to build for Google as amd64 always (:10), and to tag images `${REGION}-docker.pkg.dev/${PROJECT}/${REPO}/${SERVICE}:${GIT_SHA}` (:11).
- `skills/docker/SKILL.md:105`: `sudo -n systemctl start docker` as a fix. `context/dev-setup.md:8` runs a root step only when the operator says so and the project's `.claude/settings.json` records the standing permission; otherwise the exact command goes back to the operator.
- `skills/docker/SKILL.md:50-51`: `curl -fsS http://localhost:8080/healthz` against whatever answers on the port, where `context/AGENTS.md:57` says a green run counts only when it probed the process it started.
- `skills/bash/references/testing-linting.md:74,285`: GitHub Actions CI sections; `skills/gcloud-setup/SKILL.md:125-130`: service-account key files "for CI/CD"; `skills/frontend/references/accessibility.md:239`: axe-core "in CI". `context/AGENTS.md:45` says no CI, ever.
- `skills/golang/references/golang-io-web.md:893`: Cobra+Viper as the default CLI stack; `golang-architecture.md:200`: samber/lo; `golang-core.md:373`: samber/oops; `golang-testing.md:27`: goleak as a SHOULD — `context/AGENTS.md:49` says stdlib first, every dependency justified.
- `skills/golang/references/golang-ecosystem.md:792,795,799`: benchmark baselines written to `/tmp/report-*.txt` — `context/AGENTS.md:53` says never `/tmp`.

Found by the 2026-09-15 docs-lens sweep, `.agents/docs/reviews/2026-09-15-review-sweep-merged.md`, findings F23, F31, F33, F36, F44, F45, F46, F48, F61.

## Repro
Not reproduced as a failure; found by review (grep sweep against `context/AGENTS.md` hard rules), not by running anything.

## Root cause
Imported and older skill guidance was never cross-checked against the global hard rules when those rules were written or later tightened. `TestNoTmpInHarnessProse` (the existing `/tmp` prose check) scans only `SKILL.md` files, not their `references/` subdirectories, so drift in reference docs goes undetected.

## Fix
Fixed 2026-09-16 (Phases 2-3, clear the board): the contradicting lines in
`skills/*/references/` were corrected in text, and `TestNoTmpInHarnessProse`
now walks every non-vendored skill directory's `references/`, `commands/`
and `rules/`, not just `SKILL.md` (terraform's Terraform-native `/tmp` usage
exempted).

## Prevention
A change to a global hard rule greps `skills/*/SKILL.md` **and**
`skills/*/references/` for guidance that contradicts it. Graduated:
`TestNoTmpInHarnessProse` — its prevention rule is now a test, not prose.

## Related
BUG-003 and BUG-007 have the same shape: a rule moved or tightened and its old wording survived elsewhere. BUG-014 is a neighbour: a global rule did not reach an agent definition. With this record the shape recurs three times, which makes it retro input for a standing project rule or a broader prose test in `internal/harness`. Source: `.agents/docs/reviews/2026-09-15-review-sweep-merged.md`.
