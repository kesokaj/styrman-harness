---
owner: changelog
date: 2026-09-14
---
# Upstream items from styrman-crew v1.33.0     (2026-09-14)

## Summary

Four fixes for the items styrman-crew carried as "upstream" since its v1.31.0 mirror refresh, done here under this repo's Definition of Done during crew's v1.33.0 close-out: test isolation (`TestDryRunNeverBuilds` no longer assumes its own checkout), skill/agent name validation at parse boundary (`validName` function), symlink refusal in skill-file walks and copyTree, and LoadState guard consolidation (all manifest entries now validated through one checkpoint). Updated authorization language in context/dev-setup.md to name both operator consent and `.claude/settings.json` standing permission as required. Downstream styrman-crew v1.33.0 refreshes to this commit; argolis-devbox mirror update is a separate session.

## Type

task

## Files changed

- cmd/sh-agent/main_test.go — `TestDryRunNeverBuilds` uses new `checkoutShaped(t)` fixture (plainCopy + empty .git) instead of passing `../..` as --repo; isolates test from any tree shape
- internal/harness/inventory.go — NEW `validName(kind, file, name)` at parse boundary for skill and agent name:; rejects empty, `.`, `..`, containing `/` or `\`, failing `filepath.IsLocal`, or not `identLike` (first char [A-Za-z0-9], rest [A-Za-z0-9._-]); applied to explicit name and filename fallback
- internal/harness/state.go — `LoadState` drops manifest entries failing `filepath.IsLocal`; consolidates guard from apply.go resync loop (moved here, deleted from there)
- internal/harness/plan.go — skill-file walk refuses symlinks by name before reading
- internal/harness/adopt.go — `copyTree` refuses symlinks by name before reading (link's target was deployed under link's mode, `|0o600` on a link's 0777 → world-writable)
- internal/harness/pathguard_test.go — NEW file: TestSkillAndAgentNamesMustBeOneLocalSegment (ok-name, ../evil, a/b, /abs, .., a\b, "a b", a`b, café, .hidden, v1.2_x-y), TestResyncNeverFollowsANonLocalManifestEntry, TestLoadStateDropsNonLocalManifestEntries, TestSkillFilesAndAdoptCopiesRefuseSymlinks
- context/dev-setup.md:8 — "Root steps" sentence: operator consent AND `.claude/settings.json` standing permission both required; fallback updated to "hand the step back with the exact command"

## Agents dispatched

verifier, reviewer, security, bugs, changelog

## Verification

Verifier log `.agents/docs/testing/2026-09-14-upstream-items-verify.md`:
- `gofmt -l .` → exit 0, no output. PASS
- `go vet ./...` → exit 0, no output. PASS
- `go test -race ./... -count=1` → exit 0, all packages pass. PASS
- Sandbox `install --home .agents/scratch/x` → exit 0, status drift none. PASS
- Rendered `.claude/AGENTS.md` carries updated authorization sentence. PASS
- `sh-agent list --repo .` → Skills (27) Subagents (7), every shipped name passes charset rule. PASS
- Round 2 over the final tree (after the security fixes): PASS, appended to the same log as `## Round 2`.

Not exercised: sandbox CLI-start matrix.

## Docs

- context/dev-setup.md — changed (authorization requirement statement rendered into every consumer's AGENTS.md at install)
- No REQUIREMENTS.md change
- No README, specs, testing guide, or deploy guide change (internal mechanism only)

## Reviews

Reviewer `.agents/docs/reviews/2026-09-14-review-upstream-items.md` APPROVE WITH NITS — low: redundant trailing sentence (FIXED); informational: `validName` on ReadDir fallback unreachable (ACCEPTED).

Security `.agents/docs/reviews/2026-09-14-security-upstream-items.md` FINDINGS — high: symlink deployed under link's mode (FIXED: plan.go/adopt.go refusal + test); medium: IsLocal guard scope (FIXED: moved to LoadState + test); medium: name shape but not charset (FIXED: identLike + tests). Residuals recorded in outcome: `.git` file as worktree pointer, macOS case collisions, `.claude/settings.json` authorization record (pre-existing, unchanged).

## Downstream

styrman-crew v1.33.0 refreshes `image/harness/` to this commit; argolis-devbox mirror refresh is a separate session.

## Requirements

n/a (internal fixes; no new behavior exposed to operators).
