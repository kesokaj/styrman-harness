---
owner: verifier
date: 2026-10-01
---
# Verification — audit-fix-phase-7   (2026-10-01)

Range reviewed: working tree on 677e8ee (Phase 7)

## Status: PASS

Scope: complete for the seven steps the dispatch named; guide `.agents/docs/testing/guide.md` exists (cited, not rewritten). No CLI behavioural runs (skill-matrix rows not run here).

## Commands
- `gofmt -l cmd internal` -> exit 0, empty
- `go vet ./...` -> exit 0
- `go test -race -v -count=1 ./...` -> exit 0; every package ok (root, cmd/sh-agent, internal/dream, harness, memory, search); counts: pass 492, fail 0, skip 3 at 677e8ee+dirty for go test -race -v -count=1 ./...; vs previous: none with this exact command (phase 6 log cited 491 non-verbose runs; no fall). Output: .agents/scratch/p7-gotest.txt. Named tests PASS: TestSkillLinksResolve, TestEmbedHoldsOnlyTrackedFiles, TestEmbeddedPlanEqualsCheckout, TestAgentDescriptionRules, TestLintDescriptionLimit, TestRenderedContextUnderTarget. (The 3 skips were not itemised.)
- `go run ./cmd/sh-agent list` -> exit 0; `Skills (34)`, `Subagents (8)`
- `go run ./cmd/sh-agent render | wc -c` -> 14925 (< 15,000)
- `install --home .agents/scratch/verify-p7` -> exit 0, 266 written; `status --home ...` -> exit 0: repo has 34 skills, 8 agents; `drift none — 265 files match the repo`
- `ls .../.claude/skills | wc -l` -> 34
- templates deployed (specs/TEMPLATE.md and assets/README.md under init-project) -> exit 0
- Go samples: 6 of 6 vetted (api-auth-gcp 4, data-gcp 2, operate-gcp 0 blocks); each own module, `go mod tidy` + `go vet` exit 0; all 6 were complete files with `package`, none wrapped. Dir: .agents/scratch/verify-p7-go/
- gcloud flags (gcloud --help of the matching command; beta ones under `gcloud beta`): 29 of 30 found, e.g. run services logs read --log-filter/--freshness/--limit; beta run services logs tail --region; update-traffic --to-revisions; sql instances create --enable-point-in-time-recovery/--availability-type/--no-assign-ip/--database-flags/--edition/--deletion-protection/--backup-start-time/--database-version; run jobs deploy --args; scheduler jobs create http --oauth-service-account-email/--http-method; tasks create-http-task --oidc-token-audience/--oidc-service-account-email; pubsub subscriptions create --ack-deadline/--push-auth-service-account; iap web add-iam-policy-binding --resource-type; beta monitoring channels create --type; run deploy --iap/--add-cloudsql-instances/--cpu-boost/--no-cpu-throttling/--no-allow-unauthenticated/--min-instances/--concurrency. NOT FOUND: literal `CLOUD_IAM_SERVICE_ACCOUNT` in `sql users create` help (the help describes `--type` and links to a list of user types; the value is not enumerated in the text, so unconfirmed, not shown wrong).
- Greps: TEMPLATE.md in research/SKILL.md = 1 (>=1); testing/guide.md lines in init-project/SKILL.md = 3 (>=2); onboarding.md in agents/documentation.md = 3 (>=3); 'untouched template' in agents/verifier.md = 1; 'added 2026-10-01' in skill-matrix.md = 10; '^## Operate' in deploy/gcp.md template = 1; '## Skills' in README.md = 1
- `find skills -type f -perm /111` -> empty; `grep -L '^description:' skills/*/SKILL.md` -> empty

## Coverage of changed code
No Go source changed in the range; the skills/agents/templates are pinned by the named tests above (link resolution, embed, plan equality, description rules, size pins). Suppressions added: none checked in Go files (none in the change set). Behavioural coverage of the seven new skills is the skill-matrix rows, not run here.

## What was not exercised (copied from the fork reports)
- spec-req: no CLI behavioural run of either skill; Antigravity deployment of assets/ under the skill not checked by the fork (this run's drift check covered the Claude Code scratch home only); README skills table: fork found none (README now has `## Skills`, count 1).
- testing-explain: the two skills on a CLI; the README template's placeholder fill on a real init run; verifier agent's "when absent" vs template placeholders (this range's verifier.md contains 'untouched template', 1 hit; not read for sufficiency).
- gcp: no command ran against a GCP project; alert policy JSON not validated by the Monitoring API; scheduler sample's `email_verified` bool claim not exercised against a live token; cloudsqlconn sample compiles but never dialed an instance; cost figures from the brief.

## Not examined
The 3 skipped tests; the seven skill-matrix rows on any CLI; Antigravity install; sandbox starts; the `CLOUD_IAM_SERVICE_ACCOUNT` value; the remaining flags in the first 40 listed that were not spot-checked (e.g. --auto-iam-authn, --network, --instance, --iap on iap-less commands); govulncheck.

## Environment
Linux 6.12 (Debian 13), go from PATH, gcloud SDK in .agents/scratch/tools, network available. Start vs end `git status`: only this log (a docs artifact) was added; no source change.
