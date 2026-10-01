---
owner: reviewer
date: 2026-10-01
---
# Review — catalog-p1-p2-docs   (2026-10-01)

Range reviewed: 754a7ce..working tree (Phases 1 and 2, vendored entries excluded). Docs lens; every criterion below is `n/a — docs lens`.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

The worst problem is a rule that points at coverage that does not exist. `AGENTS.md:73` says the release rule for catalog entries "is met by the `google-catalog` router's matrix rows on both CLIs (`catalog-no-leak`, `catalog-route`, `catalog-precedence`)". No file under `.agents/docs/testing/` defines those rows (grep: zero hits outside plans and logs) and plan Phase 4, which writes them, is unticked. The same change narrows the NO-SHIP rule of v1.12.0 on the strength of rows nobody has written. Next: `delivery.md` and `remote.md` state the 19.9 MB, 13.9 MB and 7.1 MB sizes with no recorded measurement in the tracked trail. A few smaller drifts follow. The router skill's own text is accurate against the tree. All the tests I ran pass.

## Findings

- **[medium / confidence: high]** `AGENTS.md:73` — narrowed: the old rule (one Antigravity run of each changed skill and agent, NO-SHIP without it) → new clause, which waives per-entry runs and substitutes matrix rows `catalog-no-leak`, `catalog-route` and `catalog-precedence`. `git grep` over `.agents/docs/testing/skill-matrix.md`, `guide.md` and `sh-agent-matrix.md` finds none of the three; the plan puts them in Phase 4 (`plans/2026-10-01-google-skills-catalog.md:80`, unticked), and the security review (Phase 3) is also open. As written, a gate cannot check the substitute. Fix: add the three rows to `testing/skill-matrix.md` in the same commit as the clause, or word the clause "once Phase 4 lands" and keep the old rule for the catalog router until then.

- **[medium / confidence: medium]** `.agents/docs/specs/sh-agent/delivery.md:77` and `.agents/docs/deploy/remote.md:183` — the release-size figures (binary 19.9 MB against 13.9 MB, zip 7.1 MB, "about 6 MB more") appear in no verifier log or change entry in the tracked trail (`git grep '19\.9'` finds only these two lines). The plan's own box for "timing and size, before and after" (Phase 4, `plans/…catalog.md:83`) is unticked. The entry count I could check: 2,575 equals `git archive HEAD | tar -t | wc -l`, and 5,000 equals `maxZipEntries` in `internal/harness/zip.go:28`. Fix: record the measurement (command and date) in a testing log, or mark the sizes "measured on scratch, not yet logged". Also `delivery.md:77` lists the catalog under "Built, not yet released" while Phases 3 to 5 (security review, sandbox runs on both CLIs, Definition of Done) are open; say "vendored and tested in the suite; security review and sandbox runs pending".

- **[low / confidence: high]** `skills/VENDORED.md:41-42` — the `Notes:` line is left dangling. The new `## Vendored catalogs` section was inserted between `Notes:` and its bullet list, so the bullets (`grill-me`, `grilling`, `find-skills`, `codebase-design`, …, now `skills/VENDORED.md:59` onward) read as part of the catalogs section. Fix: move the catalogs section above `Notes:` or give the bullets their own `## Notes` heading.

- **[low / confidence: medium]** `.agents/docs/testing/guide.md`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/deploy/local.md` — none mentions `TestCatalog*`, `TestRouterDescription`, `TestRefreshCatalogScript`, `TestSkillListingBudget`, `scripts/refresh-catalog.sh`, or that the refresh needs the network and `git`, and `go` on PATH (it runs two full clones, 3m32s per BUG-146's record, and a `go test` writer). `AGENTS.md` § Testing sends skill changes to guide §2/§4, so a reader of the guide never learns that the catalog has its own tests, that `go test ./...` now embeds 778 catalog files, or that the install writes about 1,836 files instead of 266 (`testing/2026-10-01-verify-google-skills-catalog-phase-1.md:28`). Fix: one paragraph in the guide § for skills and one line in `local.md`.

- **[low / confidence: high]** `skills/google-catalog/SKILL.md:17` — "Provenance, pins and the refresh procedure: `catalog/MANIFEST.md`". `MANIFEST.md` holds provenance, pins, the Excluded list and Notes; it only names the script and gives no procedure. The procedure is in `AGENTS.md:62`, `skills/VENDORED.md:47-57` and the script header. Fix: point at `skills/VENDORED.md § Vendored catalogs` for the procedure.

- **[low / confidence: medium]** `skills/google-catalog/scripts/refresh-catalog.sh:2-20` — the header documents `--dest`, `--work` and `--exclude`, but not the two environment overrides `CATALOG_AGENTS_CLI_REPO` and `CATALOG_GOOGLE_SKILLS_REPO` (lines 52-53, used only by the test fixtures). `AGENTS.md:62` and `VENDORED.md` mention `--exclude` only, which is fine for the operator path. Fix: one comment line naming the two variables as test hooks.

- **[low / confidence: low]** `.agents/docs/specs/sh-agent/search.md:92` — "Binary size, linux/amd64 | 13.7 MB". A dated measurement, not a claim about now, but it will read as the size beside the 19.9 MB claim elsewhere. Fix: add the date or the commit to the cell. Unverified: whether the measurement predates the catalog.

## Restatement checks (no finding)

- Skill counts: every `34 skills` / `Skills (34)` hit in the living docs sits in a statement about the v2.0.0 tag (`delivery.md:77`, `remote.md:205`) or in dated testing logs. `README.md` carries no count. They are accurate as written and become stale at the next tag; `remote.md:205` already says "at v2.0.0".
- Third-party rule restatements: `AGENTS.md:62`, `README.md:269` (names three whole skills only), `skills/review/SKILL.md:63` and the script's header all carry the catalog clause where they restate the rule. `find-skills`, `gke-basics`, `data-gcp`, `operate-gcp`, `api-auth-gcp` and `gcloud-setup` contain no text pointing at upstream `google/skills` beyond the `gke-basics` VENDORED row, which stays accurate.
- `cloud-run-deploy` description: the old wording ("End-to-end release workflow", "Colima/Docker (ARM") survives only in dated reviews and logs. No living doc or test restates it.
- Tests: no test or comment still asserts the old behaviour. `skill_links_test.go:16` says "Skills vendored verbatim are skipped"; line 47 now also skips subtrees, which the comment at line 16 does not say (not a finding).
- Tests run: `TestCatalog*`, `TestRouterDescription`, `TestRestatementGrep*`, `TestNoTmpInHarnessProse`, `TestSkillLinksResolve` and `TestSkillListingBudget` all PASS (budget log: 19,236 of 19,500 bytes over 35 skills).

## Router skill against the tree

- Entry count 141 matches `MANIFEST.md` (7 + 1 + 131 + 2) and `VENDORED.md`. "about 36 KB" matches `INDEX.md` (36,418 bytes, 145 lines: 141 entries plus header lines).
- The seven harness skills named as winning (`cloud-run-deploy`, `data-gcp`, `terraform`, `gcloud-setup`, `gke-basics`, `operate-gcp`, `api-auth-gcp`) all exist under `skills/`.
- The four steps match the tree: INDEX line format (name, then description start), entry-relative `references/` paths, no executable bits (`MANIFEST.md`: "file modes only"), and the harness precedence. The one inaccuracy is the MANIFEST pointer above.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
- Vendored entries under `skills/google-catalog/catalog/<entry>/` (excluded by the brief) and the generated `CHECKSUMS`.
- `catalog_test.go` body beyond its header comments and the numbers it states; `refresh-catalog.sh` lines 60-196 (only the header and flags read, and the new-test comments).
- The 19.9 / 13.9 / 7.1 MB sizes were not re-measured. No release build was made.
- `skills/preflight/SKILL.md` and `context/AGENTS.md` were searched by the restatement script only, not read in full.
- The untracked `.agents/docs/reviews/2026-10-01-security-catalog-ops-iam-secops.md` is outside this lens.

## Outcome (main agent, 2026-10-01)

1. [medium] `AGENTS.md` § Release named matrix rows that existed nowhere: fixed — `testing/skill-matrix.md` holds `catalog-no-leak`, `catalog-route`, `catalog-precedence` (b0789f1) and `testing/guide.md` § 3f describes them.
2. [medium] sizes without a logged measurement: fixed — `testing/2026-10-01-verify-google-skills-catalog-phase-4.md` § Timing and size records the binary, zip, entry count, install and suite figures; `delivery.md` files the catalog under "Built, not yet released", which it is until the tag.
3. [low] `skills/VENDORED.md` section between "Notes:" and its bullets: fixed (bfd25d3).
- Lows: guide § 3f added; the router points at `skills/VENDORED.md` for the refresh procedure; the script header documents the two environment overrides; `deploy/local.md` is left without a refresh step (the procedure lives in VENDORED.md and guide § 3f; a standing operator step would restate it).
