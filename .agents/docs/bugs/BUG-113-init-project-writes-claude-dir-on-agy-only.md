---
owner: bugs
date: 2026-09-30
id: BUG-113
status: fixed
area: skills/init-project, skills/dream
---
# BUG-113: /init-project writes .claude/ on an Antigravity-only machine

## Symptom
Operator, machine using only Antigravity, project uberwatcher, first `/init-project`: `.claude/` held skills/agents/commands/rules symlinks to `../.agents/*` and `settings.local.json`; `.gitignore` got `.claude/worktrees/` and `.claude/settings.local.json`. Operator: "It is empty right but it should not be there if we only use antigravity."

## Repro
Sandbox home installed with `sh-agent install --home <home> --vendor agy` (manifest `"vendor": "antigravity"`), bare Go repo, `agy -p '/init-project …' --mode accept-edits`. Pre-fix: `.claude/` with 4 links + settings.local.json, 2 `.claude` lines in `.gitignore` (fixture ~/sandbox/agyonly/proj, commit 7cf5460).

## Root cause
`skills/init-project/SKILL.md` Step 5 (links), Step 5b items 2-3 (settings.local.json, migration from ~/.claude/projects) and Step 6 (.gitignore lines) ran unconditionally; nothing asked which CLI the machine uses, although the install manifest already records it (`sh-agent install --vendor`). Diagnosis: goldfish dispatch skipped, location mechanically obvious.
Same class found by audit: `skills/dream/SKILL.md` Step 1 resolved the store only through `.claude/settings.local.json`, so the init fix alone would have made an agy-only project fall back to `~/.claude/projects/<slug>/memory` and report no memory; Step 2 read the global context only at `~/.claude/CLAUDE.md`.

## Fix
- `skills/init-project/SKILL.md` Step 1 "Claude Code or not": `grep -q '"vendor": "antigravity"' ~/.config/sh-agent/state.json`; on a match Steps 5, 5b items 2-3 and the two `.claude/` .gitignore lines are skipped, the `.agents/memory` store is still made; any other result keeps them. An update never removes an existing `.claude/`; Step 7 names it for the operator.
- `skills/dream/SKILL.md`: new Step 1 item 0 takes `<root>/.agents/memory` on a match (same grep); Step 2 also reads `~/.gemini/config/rules/AGENTS.md`.
- Audited and left (name both vendors or only read): preflight/review shadow checks, review `--prompts` path, dreamer, feedback-log, init-project AGENTS.md template, sh-agent (vendor-aware plan).
- Left for the operator, safe fallback: `context/dev-setup.md:8` and `skills/docker/SKILL.md:108` record the root-step permission only in `.claude/settings.json`, so an agy-only project always hands a root step back. Cosmetic: `skills/handoff/SKILL.md:78` shows `claude --resume`.
- Existing operator projects keep their `.claude/`; the operator removes it by hand.
- Test: `init_project_test.go` `TestInitProjectVendorCheckMatchesManifest` (both skills carry the grep; SaveState's manifest matches the needle for antigravity only; failed before with "skills/init-project/SKILL.md lacks the vendor check"). E2E rows `init-agy-only`, `dream-agy-only` in `.agents/docs/testing/skill-matrix.md`: fail pre-fix, pass post-fix (~/sandbox/agyonly/proj2, commit c8139bf; dream first line `store: …/proj2/.agents/memory (1 memory files; …: 0)`). Log: `.agents/docs/testing/2026-09-30-1721-verify-agy-only-init.md`.

## Prevention
A harness step that writes a vendor's own path (`.claude/`, `.gemini/`, a vendor settings file) runs only when that vendor is installed: read the install manifest's vendor, never assume both CLIs. A reader of that path has a fallback for the machine where it was never written.

## Related
BUG-102 (settings trust in the same two skills). Also names these files: BUG-025, BUG-027, BUG-040, BUG-043 (init-project); BUG-096 (dream). BUG-114 (budget overflow from this change's context clause). Review: `.agents/docs/reviews/2026-09-30-1721-review-agy-only-init.md`.
