---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p9-code   (2026-10-01)

Range reviewed: working tree — changed files per `git status --short --untracked-files=all` (Audit III fix round of Phase 9), read with `git diff HEAD -- <files>`; untracked files read whole (the requirements-edit proposal; the other untracked files are review/verify records, not read).

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

No defect in the code: every asked-about mechanism holds. The worst problem is a stale description: the 15,000-byte pin now measures only the hand-written part of the context, but three living files still say it pins the render. The rest are duplicated dispatch-budget clauses (reads as two clauses in one sentence in three skills) and one unverified Terraform drift risk.

Evidence run: `gofmt -l cmd internal` empty, `go vet ./...` clean, `go test -race` on cmd/sh-agent, internal/harness, internal/dream, internal/search all ok; `terraform validate` (tools in `.agents/scratch/tools`) Success on `skills/terraform/assets/example` and on all hcl blocks of `resources-gcp.md` assembled with the existing versions/variables fixture (it needed one fixture `variable "db_url"`, an unrelated pre-existing reference); `patch --dry-run -p1 REQUIREMENTS.md` of the proposal's diff rc=0.

## Findings
- **[medium / confidence: high]** `README.md:290`, `AGENTS.md:26`, `.agents/docs/testing/guide.md:27` — each still says `TestRenderedContextUnderTarget` pins the render ("pins it under 15,000", "the result under 15,000 bytes"). `internal/harness/size_test.go:73-81` now cuts at `\n## Installed` and measures only the hand-written part; the generated inventory counts only against the 22,000 pin. A describing-text restatement of a changed threshold left standing (global rule; review-standards § 3). Fix: say "the hand-written part (everything before the generated `## Installed` inventory)" in all three.
- **[low / confidence: high]** `skills/explain-project/SKILL.md:12`, `skills/distill/SKILL.md:38`, `skills/retro/SKILL.md:47` — the new budget clause is added while the old one stays in the same sentence: explain-project says the budget twice ("a tool-call budget below its `maxTurns`" before and after the quoted recon text); distill and retro each say "a tool-call budget below its `maxTurns`; write the artifact early" and then "a 30-call budget (with) write the artifact early". Not contradictory, but one sentence with two budgets. Fix: drop the new parenthetical and keep the existing clause (or the reverse) in each.
- **[low / confidence: medium]** `agents/changelog.md:3` — the description now says "never specs, README, ROADMAP or TODO" and, later in the same value, "never writes specs, README, or ROADMAP; those belong to documentation": the lane is stated twice in one description (the body's line 34 is the consistent single statement). Fix: drop the second sentence's list.
- **[low / confidence: low]** `skills/terraform/assets/example/main.tf:88`, `skills/terraform/references/resources-gcp.md:234` — `ignore_changes = [template[0].containers[0].image, traffic, client, client_version]` is valid (all four are attributes/blocks of `google_cloud_run_v2_service`; the example validates; creation still applies the configured image, ignore only suppresses later diffs; neither file sets `revision`, `traffic` or `labels`, so no name fights). Not verified (no GCP project, no apply): whether gcloud also writes template annotations/labels that the provider reports as drift; if a first `plan` after a gcloud deploy is not empty, add that attribute to the list.

## Asked-about mechanisms, answered
- size_test marker: `RenderContext` emits `…<main>\n<extras>\n` + `RenderInventoryCompact`, which starts `## Installed`, so `"\n## Installed"` is exact; no `## Installed` heading exists in `context/*.md` (grep). Hand-written growth past 15,000 still fails (it all precedes the marker); marker absent falls back to measuring the whole render (fails safe).
- dream: `finish(…, "N project(s) failed")` sets `LastRun.Error`; `StatusLine` (`internal/dream/dream.go:203`) prints "last run <date> failed: N project(s) failed". The new assertion `strings.Contains(st.LastRun.Error, "failed")` fails on the old code (empty string).
- memoryCmd: usage line to the writer + `exitCode(2)`; `run` unwraps `exitCode` with `errors.As` (main.go:80). Old test/doc callers of the old message: none (`grep 'usage: sh-agent memory'` finds only the new code and test). `TestMemoryUsageErrorExits2` fails on the old code (`errors.New` is not an `exitCode`).
- recall.tsv: BUG-139 and BUG-140 are the only `open` rows in `bugs/INDEX.md`; the audit-fix plan has unticked boxes; the search recall test passes. (`plans/2026-09-15-harness-practices-g01-g26.md` also has `- [ ]` lines, not in the expectation, and the suite is green, so the "open plan" query evidently does not count it; not investigated.)
- IAM `google_sql_user` blocks: references resolve (`google_service_account.run`/`.migrate`, `google_sql_database_instance.db`, `.member`), `for_each` over a map of strings is valid, `roles/cloudsql.instanceUser` is the IAM-login role; validated.
- service.md: `database/pgx/v5` with `mpg.WithInstance(sqlDB, &mpg.Config{})` matches the driver's API and `skills/data-gcp/SKILL.md:119`. Not compiled here (needs network); the author's block-compile log was not re-run.
- Docs: `agents/documentation.md` rule 3 matches its tools (`Read, Grep, Glob, Edit, Write`: no shell); `agents/changelog.md` names ROADMAP/TODO/DECISIONS as never its in the body and the tree line is gone; `PRIORITY_TODO.md:22` is true now (REQUIREMENTS.md unmodified, ADR-048 absent from DECISIONS.md); the proposal's apply claim holds (dry-run rc=0; `memory-loop.md:40` has the `Known differences (R10)` heading it points at).

## Criteria
1. Dead code — pass: nothing orphaned (`errors` still used in main.go; the changelog tree line removed, not commented out).
2. Flow tracing — pass: dream status path, memory exit path and the size marker traced end to end; the new tests would fail on the old code.
3. Standards conformance — fail (medium): three restatements of the 15,000 pin stale (finding 1); gofmt/vet clean.
4. Directory layout — pass.
5. Simplicity (DIY-except-security) — pass: no new dependency; IAM database auth, not passwords, in the HCL.
6. Decoupling — pass.

## Not examined
The untracked review and verify records under `.agents/docs/reviews/` and `.agents/docs/testing/`; `README.md`, `memory-loop.md`, `search.md` and `context/AGENTS.md` diffs beyond the lines grepped; the real-home install, sandbox CLI runs and re-audit of Phase 9; `terraform plan/apply` against GCP; the Go block compile of `service.md` (no network); the g01-g26 plan's open boxes.

## Outcome
- [medium] "rendered context under 15,000" in README, AGENTS.md and the guide while the pin measures the hand-written part — fixed (the five places name the hand-written context; the inventory counts only toward the 22,000)
- [low] the budget clause doubled in `explain-project:12`, `distill:38`, `retro:47`; the changelog description stated its lane twice — fixed (one clause per dispatch; one lane sentence)
- [low] gcloud-written annotations or labels may still show as drift under `ignore_changes` — left: no project to measure against; the list names what gcloud writes on every deploy (image, traffic, client fields); a label or annotation gcloud adds would be a `terraform plan` diff to read, which the skill's procedure already requires before every apply
