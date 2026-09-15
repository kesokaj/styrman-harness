---
owner: bugs
date: 2026-09-15
id: BUG-020
status: open
area: skills/cloud-run-deploy
---
# BUG-020: cloud-run-deploy skill's own triggers bypass the /preflight gate

## Symptom
`skills/cloud-run-deploy/SKILL.md:6-9` description triggers on "deploy this service", "push to Cloud Run", "release", "push the image", "publish to Artifact Registry", "even a 'just redeploy it' request" — and never names `/preflight`. `context/AGENTS.md:45` says "Ship", "deploy", "release" or "cut vX.Y.Z" from the user's own turn means `/preflight`. Because this is an auto-loaded skill matching the same trigger words, an agent can follow it straight into a deploy without the release gate. Found by the 2026-09-15 docs-lens sweep, F20.

## Repro
Not reproduced as a failure; found by review — the skill's `description:` frontmatter was read against `context/AGENTS.md § Git`.

## Root cause
The skill predates the `/preflight` release-gate rule and restates its own entry triggers instead of deferring to the gate.

## Fix
Not fixed. This task only records the finding; no source changed.

## Prevention
A skill whose triggers overlap a gated command (here, the "Ship/deploy/release/cut vX.Y.Z" words that mean `/preflight`) names that gate as the entry point rather than restating the trigger words itself.

## Related
2026-09-15 docs-lens sweep, `.agents/docs/reviews/2026-09-15-review-sweep-merged.md` F20.
