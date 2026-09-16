---
owner: reviewer
date: 2026-09-16
---
# Review — 4aa19be-HEAD-docs (2026-09-16)

Range reviewed: `4aa19be..HEAD` — lens `--docs`, Step 4 of the v1.0.0 release gate.

## Verdict: APPROVE WITH NITS

Nothing in the diff leaves a living doc contradicting the code or another living doc. The five things the dispatch named as must-agree (update's dirty-checkout refusal, the toolchain pin, the `gofmt -l cmd internal` gate, the signed-tag release procedure, the command table) are consistent everywhere I checked. The one real gap: `install`'s two new user-visible refusals — a secret-shaped file inside a skill directory, and two skills sharing one frontmatter `name` — are implemented and tested but not described anywhere a cold reader of the docs would look (README command table, AGENTS.md "Rules for working here", REQUIREMENTS.md R5). A user who hits either sees only the Go error string from `install`.

## Findings

- **[low / confidence: high]** `README.md` command table (`sh-agent install` / `sh-agent adopt` rows), `AGENTS.md` "Rules for working here" — `internal/harness/plan.go`'s `secretShaped` refusal (install aborts when a skill directory contains a `.env`/`*.pem`/`*.key`/etc.-shaped file) and `internal/harness/inventory.go`'s `LoadSkills` duplicate-name refusal are new, user-reachable failure modes with no mention in any living doc. `AGENTS.md`'s "Rules for working here" bullet list already documents two sibling installer refusals in the same shape ("the installer refuses one without [a description]... It also refuses any frontmatter value that a strict YAML parser reads as a comment") — these two belong beside them. Advisory only, per the dispatch; does not block this release.
- **[low / confidence: medium]** `agents/verifier.md:48` still opens with `gofmt -l .` (must print nothing) as the Go auto-detect default, with an exception clause for scratch-containing projects added in this diff. Not a contradiction here — verifier Step 2 reads the project's own `AGENTS.md` first, which names `gofmt -l cmd internal` explicitly, and the exception clause says the project command wins — but a cold reader of `verifier.md` alone (its default path, without reaching the exception) could still be misled about what "the" gate command is for a Go project in general. No fix needed for this repository's own docs; flagging only because the dispatch's `gofmt -l .` check named this file's text as a hit.

## Criteria

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — the specific must-agree checks below, all confirmed as of the working tree:
   - `update` dirty-checkout refusal (ADR-027): `REQUIREMENTS.md:19`, `README.md:25,40`, `AGENTS.md` (Release section, unaffected — update isn't there), `.agents/docs/deploy/remote.md:27-33,45,99`, `.agents/docs/deploy/local.md:31,43` all state the refusal consistently; no doc still describes unconditional pull-and-reinstall.
   - `go.mod` toolchain pin (ADR-028): `AGENTS.md:17`, `README.md` (Working on the harness section) both say `toolchain go1.25.8` clears the two stdlib advisories while the `go` directive stays 1.24.4; matches ADR-028's Decision exactly. `.agents/docs/testing/sh-agent/guide.md` and the deploy guides say nothing that contradicts it.
   - `gofmt -l cmd internal` (BUG-048): present and consistent in `AGENTS.md:12`, `README.md` (Working on the harness), `.agents/docs/deploy/local.md:18`, `.agents/docs/testing/sh-agent/guide.md:14`. No living doc for *this repository* still says `gofmt -l .` — the one remaining `gofmt -l .` hit, `agents/verifier.md:48`, is the harness-wide default with an explicit override clause (see Finding 2 above), not a per-repo instruction.
   - Signed-tag release procedure: `AGENTS.md` § Release, `.agents/docs/deploy/remote.md` § Releases, `.agents/docs/architecture/DECISIONS.md` ADR-026, and `skills/preflight/SKILL.md` Step 6 item 3 all describe the same one procedure (tag HEAD with no target, verify parent-is-gated-head and records-only diff before tagging, `verify-tag` against the allowed-signers line before push) — `preflight`'s Step 6 defers by name to "the project's AGENTS.md (a Release section) or the deploy guide it points to," which is exactly what this repo supplies.
   - `install`'s secret-shaped-file and duplicate-skill-name refusals: implemented (`plan.go`, `inventory.go`) and tested (`TestPlanRefusesSecretShapedFiles`, duplicate-name test per BUG-050/BUG-049 rows in `bugs/INDEX.md`) but undocumented in any living doc a user would read before hitting them (Finding 1).
   - `README.md` command table vs. binary behaviour: checked every row (`status`, `diff`, `install`/`sync`, `update`, `uninstall`, `list`, `render`, `adopt`, `version`) against the diffed Go source; all match except the two omissions in Finding 1.
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

- `context/AGENTS.md`, `agents/documentation.md`, `agents/security.md` diffs — read the diff stat but did not re-grep their subject nouns beyond the security-trigger ("crypto") wording change already covered by BUG-023's closure; low risk, small diffs.
- `skills/bash/references/*`, `skills/docker/SKILL.md`, `skills/gcloud-setup/SKILL.md`, `skills/golang/references/*`, `skills/cloud-run-deploy/SKILL.md`, `skills/distill/SKILL.md`, `skills/frontend/references/accessibility.md`, `skills/retro/SKILL.md`, `skills/review/SKILL.md` — all show small diffs (2-15 lines) in the stat; not individually re-read against their own subject nouns for stale restatements elsewhere in the skill tree.
- `.agents/docs/gcp.md` — unchanged in this range; not re-checked against the new toolchain/release-procedure claims since nothing in the diff touches GCP.
- `cmd/sh-agent/main.go` / `main_test.go`, `internal/harness/apply.go`, `internal/harness/repo.go`, `internal/harness/harness_test.go` — read only the parts relevant to the five must-agree checks (update refusal, secret-shaped refusal, duplicate-name refusal); not reviewed line-by-line for code-quality findings, which is out of scope for the `--docs` lens.
