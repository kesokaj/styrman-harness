---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p9-docs   (2026-10-01)

Range reviewed: working tree, changed files per `git status --short --untracked-files=all` (the Audit III fix round of Phase 9), `git diff HEAD -- <files>`; docs lens over the living docs.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

Worst first: the new Terraform `ignore_changes` on the Cloud Run image contradicts the HCL beside it and `terraform/SKILL.md`, which still says the image tag reaches Terraform as a variable (an apply with a new `git_sha` now does nothing for the service, while the migrate job still follows it). The size-pin edit in `size_test.go` changed what `TestRenderedContextUnderTarget` measures, and four living files plus the test's own constant comment still say "rendered context under 15,000". Everything else checked is consistent: the REQUIREMENTS-edit state, the recall figure, IAM users, ROADMAP date, PRIORITY_TODO Phase 9.

## Findings
- **[medium / confidence: high]** `skills/terraform/SKILL.md:15` vs `skills/terraform/assets/example/main.tf:66,86-88` and `skills/terraform/references/resources-gcp.md:199,230-234` — SKILL.md says "the image tag is the git short SHA … passed into Terraform as a variable" and line 16 says `/preflight` is the deploy path with this skill as "the HCL under it". The diff adds `ignore_changes = [template[0].containers[0].image, traffic, client, client_version]` on the service, while the same file still writes `image = ".../${var.git_sha}"`. After the first create, a changed `var.git_sha` is ignored on the service: the variable is dead for the service, and a reader following line 15 expects an apply to roll the image. The migrate job (`resources-gcp.md:296`) has no `ignore_changes`, so job and service diverge. The ownership split ("cloud-run-deploy owns the revisions") is stated only in the HCL comments: not in `terraform/SKILL.md`, not in `cloud-run-deploy/SKILL.md` (its only Terraform sentence is the `--allow-unauthenticated` row, line 95), and the `gcp.md` template (`skills/init-project/assets/agents-tree/docs/deploy/gcp.md:31,40,65`) does not mention it. Fix: rewrite SKILL.md:15 to the split (first create via var, later images by `cloud-run-deploy`, ignore_changes named), say the same in one sentence in `cloud-run-deploy`, and either ignore the job image too or say why not.
- **[medium / confidence: high]** `AGENTS.md:26`, `README.md:290`, `.agents/docs/testing/guide.md:27`, `skills/retro/SKILL.md:143`, `internal/harness/size_test.go:52-57` — all still describe the 15,000-byte pin as bounding the rendered context ("pins it under 15,000", "the harness's own 15,000-byte ceiling"). The diff makes the test measure only the hand-written part, cut at `\n## Installed`. I ran it: hand-written 14,107 bytes, render 14,992 (8 bytes under 15,000 today, but no longer enforced). The retro skill tells an operator that a green run means the render is under 15,000; it no longer does. Fix: each copy says "hand-written part (everything before `## Installed`)", and `size_test.go`'s constant comment gets the same.
- **[low / confidence: high]** `.agents/docs/specs/sh-agent/search.md:87` vs `:11` — the performance table still reads "Full build, 1,080 files | 1.7 s" while line 11 now says "2.1 s for 1,351 files on 2026-10-01". Two file counts and two times in one spec. Re-date the row or drop one figure.
- **[low / confidence: medium]** `.agents/docs/specs/sh-agent/search.md:89` — the 2026-10-01 re-measure cites `testing/2026-09-29-search-recall.md`, which holds the old 0.653/0.379 run, so no testing record carries 0.633/0.296. I re-ran `SH_AGENT_RECALL=1 go test ./internal/search -run TestRecall`: index mean 0.633 over 41 confirmed. The grep arm (0.296) I did not see in the filtered output. Name the 2026-10-01 verifier log as the source if it holds the numbers.

## Restatement results
- Recall figures: `0.653`/`0.379` appear only in `plans/2026-09-29-search-index.md` (dated plan, stays) and in `search.md:89` as the history. No other living file quotes them. OK.
- Rebuild time: `under two seconds` has no remaining hit. `1,080 files` survives only at `search.md:87` (finding 3).
- REQUIREMENTS-edit state: the proposal file, `PRIORITY_TODO.md:22`, `memory-loop.md:42` and plan D29/D29a (`plans/2026-09-30-audit-fix.md:116` unticked, `:229`, `:231`) agree: approved in intent, exact text awaiting the operator, `REQUIREMENTS.md` untouched (`git diff HEAD --stat -- REQUIREMENTS.md` empty). "applied in Phase 3" has no hit. `ADR-048` is only in the plan (as the future record) and in PRIORITY_TODO ("no ADR-048 exists yet"). The proposal's diff applies to the current `REQUIREMENTS.md` (`patch --dry-run`, no reject). OK.
- Changelog lanes: `agents/changelog.md` description and body agree that ROADMAP, TODO and DECISIONS are not its; `documentation.md` item 4 agrees. The description states the "never specs, README, ROADMAP" lane twice (once with TODO, once without); not a contradiction.
- `TaskStop`: only `context/AGENTS.md:35` uses the word ("or the CLI's own stop"); no copy in `skills/review` or the agents to diverge. OK.
- IAM users in HCL vs `data-gcp/SKILL.md:74-81`: `trimsuffix(email, ".gserviceaccount.com")` yields `${APP}-run@${PROJECT}.iam` / `${APP}-migrate@…`, the same shape; `roles/cloudsql.instanceUser` for both via `iam_db_login`, `roles/cloudsql.client` for both (`resources-gcp.md:133`, `:274`). OK.
- `ROADMAP.md` "Last reviewed: 2026-10-01" matches today. `PRIORITY_TODO.md` still has Phase 9 unticked. OK.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined / unverified
- `skills/golang/references/service.md:366`: the claim that `database/postgres` pulls `lib/pq` and that `database/pgx/v5` is the right migrate import path; not checked against golang-migrate.
- The Terraform HCL was not run through `terraform validate`/`tflint` (not installed here); `ignore_changes` syntax judged by reading.
- `grep arm 0.296` and the "2.1 s / 1,351 files" build time not re-measured.
- Go code diffs (`cmd/sh-agent/main.go`, `memory_test.go`, `internal/dream/*`) and the `reviews/audit-3-*` and `testing/2026-10-01-verify-real-home.md` files: not read (outside this docs scope).
- Minor nits left out: `explain-project/SKILL.md` Step 1 now names the explorer budget twice; the README `--from-zip` sentence repeats the dry-run cleanup it already says later; a doubled blank line in `main.tf` before `lifecycle`.

## Outcome
- [medium] the Terraform/gcloud ownership split lived only in HCL comments; `var.git_sha` inert after the first apply; the job diverged — fixed at `skills/terraform/SKILL.md:15` (Terraform sets the image once, then ignores it; shape versus revisions), `skills/cloud-run-deploy/SKILL.md` § Canary (one sentence), the gcp.md template's Terraform row, and `resources-gcp.md`'s job block (`ignore_changes` on its image); assembled blocks validate
- [medium] "rendered context under 15,000" restated in five places while the pin measures the hand-written part — fixed at `AGENTS.md` § Setup, `README.md`, `testing/guide.md` § 2, `skills/retro/SKILL.md:143`, `size_test.go` comment
- [low] `search.md:87` old build numbers; the recall citation — fixed (1,351 files / 2.1 s with the 2026-09-29 figures kept beside; the re-measure cites the Phase 9 verify log; the grep arm 0.296 was measured in the same run)
