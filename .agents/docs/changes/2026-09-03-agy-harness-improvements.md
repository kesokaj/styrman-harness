---
owner: changelog
date: 2026-09-03
scope: changeset:cd5213ec7b02
---

# Harness Audit and Standardization

Comprehensive audit and enhancement of `styrman-harness` across 5 evaluation phases and 6 approved improvement proposals.

## Changes

### 1. Trail Freshness Fix
- Fixed `internal/trail/check.go` to correctly identify when HEAD commits touch zero source files, preventing false `STALE` failures on doc-only commits.
- Added `TestCheckFreshnessDocOnlyCommit` in `internal/trail/check_test.go`.

### 2. Standards Conformance & Skill Optimization
- Consolidated `anti-ui-slop` playbooks into standard `references/` directory and updated links.
- Moved `codebase-design` markdown guides (`DESIGN-IT-TWICE.md`, `DEEPENING.md`) into `references/`.
- Pruned dead vendor artifacts (`agents/openai.yaml`, checksums, manifests).
- Added explicit negative trigger boundaries (`Do NOT use for...`) across all 28 skills in `shared/skills/`.
- Updated `grill-me` launcher for cross-platform execution on both Claude Code and Antigravity.

### 3. Native Git Pre-Commit Hook (Defense-in-Depth)
- Created `.githooks/pre-commit` executing `harness hook pretooluse` to block credential leaks across all IDEs and tools.
- Updated `setup.sh` to configure `git config core.hooksPath .githooks`.

### 4. Cold-Start Resume Protocol Streamlining
- Refactored Resume Protocol in `shared/AGENTS.md` (and mirrors `claude/CLAUDE.md`, `antigravity/GEMINI.md`) to use progressive disclosure, cutting cold-start context from ~16k to ~9.5k tokens.

### 5. Proportional Security Row
- Added `securityLines: 150` threshold to `internal/trail/owed.go`, evaluating security criteria in the `reviewer` pass for changes < 150 lines without dependency updates.
- Added `TouchesDependencyManifest()` in `internal/trail/source.go`.

### 6. Tuning & Derivation Documentation
- Documented `sonnet` model inheritance in `internal/install/derive.go`.
- Configured `"effortLevel": "medium"` in `claude/settings.json` for fast coding turns.
