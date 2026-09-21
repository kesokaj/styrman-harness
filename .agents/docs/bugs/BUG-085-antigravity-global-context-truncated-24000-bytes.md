---
owner: bugs
date: 2026-09-21
id: BUG-085
status: fixed
area: context / installer (Antigravity global rules)
---
# BUG-085: Antigravity loads only the first 24,000 bytes of the deployed global context

## Symptom
Antigravity CLI (1.2.7) truncates `~/.gemini/config/rules/AGENTS.md` (27,114 bytes, identical content to `~/.claude/AGENTS.md`) at byte 24,000, mid-sentence ("Once given, record them in"), wrapping it as `<RULE[...AGENTS.md]>...<truncated 3114 bytes>` inside `<user_rules>`. The missing tail drops the rest of the GCP-defaults bullet, the Backend/Frontend/Machines/root-steps/test-locally/amd64/Artifact-Registry/Cloud-Run/gcloud bullets of `context/dev-setup.md`, and the whole generated `## Installed` inventory, on every Antigravity session. Claude Code loads the same file in full.

## Repro
`agy -p 'Do not use any tools. Answer only from the global rules already loaded in your context. What is the LAST section heading of the global context, and how many skills does it say are installed?' --print-timeout 3m` → answers `## Environment & Stack` and UNKNOWN, and self-reports truncated by 3,114 bytes. `claude -p` with the same prompt answers `## Installed` and "Skills (27)". Confirmed in two conversation DBs: `~/.gemini/antigravity-cli/conversations/48ae34c2-...db` (2026-09-21, truncated 3,114 of 27,114) and `...3c034e5c-...db` (2026-09-18, truncated 2,553 of 26,553, v1.3.0-era). `strings <db> | grep -n -B1 -A1 '<truncated'`.

## Root cause
`context/AGENTS.md` grew from 7,349 chars (ADR-009 trim, 87a862d) to 23,194 chars (42fad8f) as nineteen BUG prevention rules graduated into it; no test pins the deployed context's rendered size against any vendor cap (`TestHarnessMap` only pins verb names). The Antigravity vendor docs state a 12,000-character rules-file limit; the observed cut is 24,000 bytes — whether the cap is per-file or per `<user_rules>` block is unknown. The crossing commit is 5bc56fb (2026-09-17): rendered context reached 24,948 bytes, above the observed cut; every release since v1.2.0 (v1.2.0-v1.4.0) ships a context Antigravity truncates. The existing real-home smoke procedure (`.agents/docs/testing/2026-09-03-verify-real-home.md`) asks for the skill count and pruned names, which sit in the truncated tail, but the CLI answers those from its own internal skill listing rather than from the file, so the check passed on a truncated file; the one post-crossing Antigravity context check (d278f5a) asked about a Hard-rules sentence that loads (bytes 13,875-22,609), so it never probed the tail.
`Diagnosis missed:` no diagnosis dispatch preceded this — found directly by an agent-context audit, not a goldfish brief.

Probe (2026-09-21): a second rules file `~/.gemini/config/rules/aa-canary-bug085.md` (45 bytes, one canary sentence) was placed beside `AGENTS.md` (27,114 bytes) in the real home; `agy -p` was asked for the canary word and for AGENTS.md's last heading and truncation. The canary was UNKNOWN and the string does not occur anywhere in the session's conversation DB — no `<RULE[...aa-canary-bug085.md]>` block was built at all. AGENTS.md's own truncation grew from 3,114 to 3,456 bytes lost (342 bytes more), i.e. it was allowed only 23,658 bytes with the second file present. The cap is therefore on the whole global-rules block, not per file: a second file costs the first one room and, once the block fills, is dropped silently with no marker of its own. Splitting the context into several `rules/` files cannot fix this bug; pruning the deployed context under the cap is the only fix. Whether a second file loads when the block is under the cap was not tested.

## Fix
Pruned `context/AGENTS.md` and `context/dev-setup.md` under the cap in commit 1df57f3 (Phase 1 of `.agents/docs/plans/2026-09-21-context-consolidation.md`; 27,163 → 20,012 bytes, 44 operator-approved actions), then 20,205 bytes after Phase 3 (807de3f). Pinned by `TestRenderedContextFitsAntigravityBlock` (`internal/harness/size_test.go`, budget 22,000 bytes; must-fail mutation verified in 94a6596). A canary sentence sits as the last line of the deployed context (770e76a), asked for by the real-home smoke check. The real home on this box was updated to 94a6596 on 2026-09-21; Antigravity's own conversation store then held one `<RULE[…/AGENTS.md]>` block with the canary line and zero `<truncated` markers (`.agents/docs/testing/2026-09-21-verify-real-home-smoke-5.md`).

Still open elsewhere: the Mac and the devbox get the fix at their next `sh-agent update` to a release carrying it (v1.5.0 pending).

The proposed two-file probe to learn whether the cap is per-file was run 2026-09-21: the cap is per whole rules block, not per file, and a second file is silently dropped once the block is full, so splitting context across multiple files is not a viable fix.

## Prevention
A deployed file that a vendor caps is pinned by a test that renders it from the real repo against the observed cap with a margin. Every smoke check asks one question answerable ONLY from the file's last section — never a question (like a skill count) the CLI can answer from its own internal listing, because that gives a false pass on a truncated file.

## Related
BUG-004 (Antigravity silently drops content it cannot parse — same vendor, same "silent partial load" shape). `context/AGENTS.md` § Hard rules, "silence from your own instrumentation is not evidence" (a check that cannot fail is not a check — the smoke check here could not fail because it never read the truncated tail). ADR-009 ("the context stays short"). `.agents/docs/plans/2026-09-21-context-audit-assessment.md` § Defect.
