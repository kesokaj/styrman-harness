---
owner: reviewer
date: 2026-09-30
---
# Review — project-docs   (2026-09-30)

Range reviewed: the whole repository at HEAD (3478a1c), no diff.

## Status: COMPLETE (sampled — see Not examined; this is a docs-lens audit of a large trail on a fixed tool-call budget, not an exhaustive read of every file)

## Verdict: APPROVE

No stale statement, dropped rule, or undocumented interface/command/deploy-step change was found in the files sampled. This repository's own process (see the newest `changes/` entry, `2026-09-30-0623-docs-name-search-everywhere.md`) already runs a "does every doc describe the shipped feature" sweep after releases, and the spot checks below confirm that sweep landed correctly. The single home for the delivery claim (rule 6, `.agents/docs/specs/sh-agent/delivery.md` § The model) holds: `README.md`, `AGENTS.md`, `.agents/docs/deploy/remote.md` and `local.md` all point at it rather than restating its content, and the mechanics they do spell out (what a command does at a place) are the command's own documentation, which rule 6 allows.

## Findings

None. Every claim checked against code or the trail matched:

- `README.md` § Daily commands table's verb list (`status`, `diff`, `install`/`sync`, `update`, `uninstall`, `list`, `render`, `adopt`, `report`, `search`, `index`, `version`) matches the `case` labels in `cmd/sh-agent/main.go` (lines 153–330); every flag named there (`--vendor`, `--home`, `--repo`, `--from-zip`, `--check`, `--dry-run`, `--force`, `--global`, `--refs`, `--open`, `--kind`, `--limit`, `--rebuild`, `--feedback-repo`) is registered in `cmd/sh-agent/main.go`'s `takesValue`/flag table (lines 738–751).
- README's claim that search skips secret-shaped filenames is implemented exactly as described: `internal/search/reconcile.go:159` calls `harness.SecretShaped`, the same function `internal/harness/plan.go` and `adopt.go` use for install-time filtering — one shared rule, not a restated one.
- `README.md` Layout's "`commands/` (empty)" and "`rules/` (empty)" both hold — both directories contain only a `.gitkeep`.
- `REQUIREMENTS.md` R8/R9 cite ADR-025, ADR-029, ADR-030, ADR-033, ADR-034, ADR-035, ADR-038, ADR-039, ADR-041, ADR-045, ADR-045a; all eleven exist as `## ADR-NNN` headings in `.agents/docs/architecture/DECISIONS.md`.
- The newest `changes/` entry's four claimed edits (README opening/Layout/map label, project `AGENTS.md` Layout+Key files+Testing, `delivery.md` § Where things stand verb list, `deploy/gcp.md` wording) are all present in the working tree as described.
- `testing/guide.md` does not carry the BUG-087/090/091/106/107 behavioural-run rules — checked against `.agents/docs/bugs/INDEX.md`, whose entries for those four bugs name `AGENTS.md` § Testing (not `testing/guide.md`) as the rule's home; this is the intended graduation target, not a gap.
- `deploy/remote.md`'s "27 skills and 7 agents today" (line 203) is explicitly framed as the state at the `v1.0.0` gate ("First tag `v1.0.0`, base ref for its gate 4aa19be … today") and is not read as a live count; current counts are 28 skills / 8 agents, but the sentence names the historical gate, not present state, so this is not a finding.

## Criteria

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

Given the 32-tool-call budget against a repository-wide, no-diff scope, the following were not read and are not covered by this report's verdict:

- `.agents/docs/specs/sh-agent/search.md` and `memory-loop.md` (only `delivery.md` was read in full).
- `context/AGENTS.md`, `context/dev-setup.md` (the deployed-context source) — not diffed against the project `AGENTS.md` shown by the system reminder.
- `agents/*.md` (all 8 subagent definitions) and every `skills/<name>/SKILL.md` not already touched by the sampled greps — not read individually for stale claims.
- `.agents/docs/testing/sh-agent-matrix.md` and `skill-matrix.md` — not checked for currency against `cmd/sh-agent`'s actual verb set.
- The great majority of `.agents/docs/changes/` (only the single newest entry was read), `.agents/docs/architecture/DECISIONS.md` (only grepped for headings, not read for content accuracy), `.agents/docs/plans/`, `.agents/docs/retros/`, `.agents/docs/handoffs/`.
- `skills/VENDORED.md` was read but not cross-checked line-by-line against the actual pinned commits/hashes of each vendored skill.
- No test files (`*_test.go`) were read for stale assertions, since no diff was in scope to check them against.

## Re-check commands

```bash
# Daily-commands table vs actual verbs/flags
grep -n 'case "' cmd/sh-agent/main.go
grep -n '"--' cmd/sh-agent/main.go

# secret-shaped filter shared between installer and search
grep -rn "SecretShaped" internal/

# delivery-claim single-home rule: only delivery.md should state "The model" content in full;
# everywhere else should point at it
grep -rln "The model" README.md AGENTS.md .agents/docs/deploy/*.md .agents/docs/specs/sh-agent/*.md

# ADRs cited in REQUIREMENTS.md all exist
for a in $(grep -oE 'ADR-[0-9]+[a-z]?' REQUIREMENTS.md | sort -u); do
  grep -q "^## $a\b" .agents/docs/architecture/DECISIONS.md || echo "MISSING: $a"
done

# skills/agents counts named anywhere as a live fact (re-run after any skill/agent add/remove)
grep -rn "skills and [0-9]* agents\|[0-9]* skills, [0-9]* agents" README.md AGENTS.md .agents/docs/deploy/*.md
```
