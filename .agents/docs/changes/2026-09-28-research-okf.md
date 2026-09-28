---
owner: main
date: 2026-09-28
---
# Research: OKF (Open Knowledge Format) against the harness — not now     (date: 2026-09-28)

## Summary
The operator asked whether the harness implements OKF, Google Cloud's markdown-plus-frontmatter knowledge format, and whether it should. It does not. The research record weighs seven approaches, from watching through a vendored skill, a harness skill, a conformant trail and an `sh-agent` validator to a project pilot and a formal rejection. It recommends A: watch, with named triggers, the realistic one being a project that has to import, publish or maintain a bundle.

The operator approved A. They answered the one open scoping question, whether the three unsurveyed install places do BigQuery work: "no BigQuery work on other machines". On the operator's request, one of the record's found-along-the-way items is fixed: frontmatter on the audit prompt.

## Type
task

## Files changed
- `.agents/docs/plans/2026-09-28-research-okf.md` — new research record: answer, what OKF is, approaches A–G, comparison and rule-fit tables, recommendation and triggers, the operator's answer, found-along-the-way table, corrections from verification, not verified.
- `.agents/docs/plans/2026-09-21-context-audit-prompt.md` — adds the `owner: main` / `date: 2026-09-21` frontmatter that `context/AGENTS.md` requires of every docs artifact (found-along-the-way #2). The immutable handoff `2026-09-04-argolis-devbox-instructions.md`, which also lacks frontmatter, is left alone.
- `.agents/docs/testing/2026-09-28-verify-research-okf.md` — verifier log, PASS.

## Agents dispatched
- **Workflow run `wf_a9b1f1c8-9a1`:** 8 agents, 0 errors.
  - 4 read-only readers: spec, ecosystem, harness map, consumer projects.
  - 1 drafter.
  - 3 adversarial verification lenses: external facts, harness claims and compliance, devil's advocate. All three returned "holds with corrections", none changed the recommendation, and their corrections are applied in the record.
- **`verifier`:** one run, PASS.

## Verification
- **The main agent spot-checked against primary sources:**
  - the conformance and trust wording, against the spec clone at `ad30107` (`SPEC.md:3,57-63,403-410,736-764`);
  - zero OKF hits in the saved Claude Code (2.1.283), agy (1.2.12), Gemini CLI and Codex changelogs;
  - the fabricioctelles skill's `scripts/validate.sh`, against ADR-006, ADR-007 and `embed.go:9`;
  - ADR-005's wording;
  - `deploy/remote.md:56`;
  - the frontmatter keys on our own skills (8 carry `argument-hint`).
- **The `verifier` ran on the tree holding only the new record** (`testing/2026-09-28-verify-research-okf.md`):
  - gofmt, vet and `go test -race`: green;
  - rendered context: 20,329 of 22,000 bytes;
  - scratch-home install and status: drift none;
  - `list`: 27 skills, 7 agents;
  - the record: text, not binary, with no control bytes.
- **Written after the verifier returned, and outside its run:** the #2 frontmatter fix (4 lines, docs only) and the one table cell that marks it fixed. The operator's answer in §6 was in the record before the verifier ran.
- **No `reviewer`:** there is no source diff. The main agent read the record itself.

## Docs
No docs changed. This is a research record and one frontmatter fix; no interface, command or requirement was amended.

## Requirements
n/a
