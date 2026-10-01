---
owner: main
date: 2026-10-01
---
# 2026-10-01-1520 — Google skills catalog, Phase 2: the listing budget

Plan `plans/2026-10-01-google-skills-catalog.md` Phase 2.

- `TestSkillListingBudget` (`internal/harness/size_test.go`): the names plus descriptions of the model-invocable skills are at most 19,500 bytes — Claude Code shows descriptions for about 19,800 and bare names beyond. Red at 19,585 with the router added; green at 19,236 after `skills/cloud-run-deploy/SKILL.md`'s description went from 967 to 618 bytes with its routing kept (a ship, deploy or release request means `/preflight`; when to use it directly; the four "Do NOT use" pointers).
- Docs: README § Layout, project `AGENTS.md` § Layout, § Rules (a catalog changes only through its refresh script; its diff gets a security review; entries never become global skills) and § Release (catalog entries are data: the router's matrix rows on both CLIs and a security review of a refresh's diff, not one run per entry — D11, D12), `specs/sh-agent/delivery.md` and `deploy/remote.md` (binary 13.9 → 20.0 MB for linux/amd64; zip 4.8 → 7.1 MB, 1,486 → 2,575 entries of the 5,000 `update --from-zip` accepts; measured with `go build -trimpath` and `git archive` on cf9c5e2), `testing/guide.md` § 3f.
- Docs lens over Phases 1 and 2 (`reviews/2026-10-01-review-catalog-p1-p2-docs.md`, REQUEST CHANGES): fixed here — `skills/VENDORED.md`'s catalog section moved above "Notes:", the router's pointer for the refresh procedure, the script header's environment overrides, the guide section; closing with Phase 4 — the matrix rows the release rule names.

Verification: `testing/2026-10-01-verify-google-skills-catalog-phase-2.md` PASS (356 top-level tests with -race, both size pins, the budget's must-fail control, scratch install 35 skills and no drift). The docs-lens fixes after that log are text in two skill files, a script comment and a guide section; the catalog, router and budget tests were re-run by the main agent after them (ok). Agents dispatched: verifier ×1, reviewer ×1 (docs lens); Phase 3's five security reviewers and Phase 4's sandbox rows are in flight.
