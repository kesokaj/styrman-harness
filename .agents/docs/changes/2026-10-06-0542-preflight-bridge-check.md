---
owner: main
date: 2026-10-06
---
# /preflight waits for every place to cross a bridge release (BUG-161)

## Type

bugfix

## Summary

`skills/preflight/SKILL.md` Step 6: before the release commit, when the previous release was a bridge (one older installs need to move past a trust or format change, per its notes, release entry or an ADR), `/preflight` names it and asks the operator, in their own turn, whether every place has crossed; without a yes it stops before anything is committed or tagged, and an unattended run never tags past a bridge. Project `AGENTS.md` § Release points at it. Cause: v2.4.1 was tagged 58 minutes after the v2.4.0 bridge and stranded einherjar on v2.3.2 (BUG-161).

## Verification

Three lines; diff read by the main agent (under the reviewer threshold). `TestSkillListingBudget`, `TestLint*`, `TestRenderedContext*`, `TestHarnessMap` pass; scratch install (`--home .agents/scratch/home-b161`) deploys the new paragraph to both CLIs' `preflight` skill. Not run: a behavioural run of `/preflight` against a bridge; the next release owes the Antigravity run of this changed skill (AGENTS.md § Release).
