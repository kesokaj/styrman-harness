---
owner: changelog
date: 2026-09-16
---
# SSH-signed tag release procedure, ADR-026 (2026-09-16)

## Summary
Implemented ADR-026: releases are annotated `vX.Y.Z` tags on `main`, SSH-signed with the operator's ed25519 key (`~/.ssh/id_ed25519`, fingerprint `SHA256:5HpYQJjpG+yE1xXBTOIZPsTkbkgkVnV8z0eG8hBawyU`), cut by `/preflight vX.Y.Z` from the operator's own turn, pushed with a plain push (never moved or deleted). Signing uses per-command flags (`-c gpg.format=ssh -c user.signingkey=...`), never a config change. Before the tag: tree clean, release commit on gated head, release commit changed records only. After the tag: sits on HEAD, is annotated, last armor header is SSH SIGNATURE, `git verify-tag` passes against a scratch copy of the allowed-signers line read from crew's `/etc/styrman/allowed_signers`. Base ref `4aa19be` (commit both consumers carry); first tag `v1.0.0`; releases are MINOR (new/changed skill, agent, command, rule, context or installer behaviour), PATCH (fix to released behaviour only), or MAJOR (what deployments must react to). Tag is sole version tracking. Crew's channel (root timer polling this repo) consumes signed tags; `sh-agent update` on operator machines and argolis keeps following `main`; crew's mirror frozen at 4aa19be until channel ships. Rehearsed in `.agents/scratch/tagcheck/`: signing works headless; own key passes, foreign key "No principal matched", lightweight reads as `commit`, unsigned annotated "no signature found", unreadable file → operator's scratch copy used. Security review round 2: medium finding (tautological first target check) fixed by deletion; round 3 not run, operator waiting. Review round 2 REQUEST CHANGES findings fixed; documentation agent wrote sections twice, main agent applied command-form fixes to `remote.md` and `AGENTS.md` § Release; main agent fixed three changelog-agent mis-statements (which box runs git/OpenSSH, mirror sentence, base ref).

## Type
docs/decision

## Files changed
- `.agents/docs/architecture/DECISIONS.md` — ADR-026 (full decision: context, decision, consequences — ADR-023 sentence marked Superseded, ADR-010 marked Superseded in part)
- `.agents/docs/deploy/remote.md` — § Releases (signed tags, pre-push checks, verification, crew allowed-signers line); § Downstream mirrors (crew frozen at 4aa19be until channel ships, argolis refreshed as before); § Operator steps (only root step: `sudo -n cmp` comparing crew's file to scratch copy)
- `AGENTS.md` — § Release (versioning scheme: MINOR/PATCH/MAJOR; key files row; Release section is project-supplied like every skill)
- `README.md` — one sentence (releases are SSH-signed tags)
- `skills/preflight/SKILL.md` — Step 6 item 3 (follows a project's Release section; tags HEAD; confirms target)
- `agents/security.md` — unchanged from retro 2 (probe repository with `git init` and `git -C`)

## Agents dispatched
documentation ×2 (docs lens on ADR-026 and on deploy/remote.md, README updates), verifier ×2 (testing log for retro-2 task verifies signing, tag checks, and `/preflight vX.Y.Z` rehearsal must-pass and must-fail forms), reviewer ×2 (docs lens round 1 REQUEST CHANGES, round 2 APPROVE after fixes), security ×2 (round 1 one medium and five low; round 2 one medium rebutted, four low fixed, two left), changelog ×2 (ADR-026 embedded, this entry).

## Verification
Signing rehearsed in `.agents/scratch/tagcheck/` by main agent: per-command flags work headless with no passphrase; `git verify-tag` passes for operator's key, fails "No principal matched" for others. Tag checks rehearsed twice by verifier in round 2 of testing log (must-pass: lightweight shows type `commit`, unsigned shows "no signature found", four must-fail cases; final no-target form rehearsed by main agent must-pass and must-fail). Full test suite `go test -race -v -count=1 ./...` PASS (pass 79, fail 0, skip 1) in testing log. Logs: `.agents/docs/testing/2026-09-16-verify-retro-2-and-release-procedure.md` (round 2, rehearsing rewritten preflight block).

## Docs
`.agents/docs/architecture/DECISIONS.md` (ADR-026, marked prior ADRs Superseded), `.agents/docs/deploy/remote.md` (§ Releases, § Downstream mirrors, § Operator steps), `AGENTS.md` (§ Release), `README.md` (one sentence), `skills/preflight/SKILL.md` (Step 6 item 3). REQUIREMENTS.md unchanged (R3 and R4 hold; signed tag is versioning, not a new mechanism).

## Requirements
n/a — no new requirement; R3 (git is backup and versioning only) and R4 (update pulls from git) remain satisfied.
