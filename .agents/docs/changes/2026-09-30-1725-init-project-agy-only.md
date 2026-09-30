---
owner: main
date: 2026-09-30
---
# /init-project writes nothing under .claude/ on an Antigravity-only machine; /dream finds the store there     (2026-09-30)

## Summary
The operator, on a machine that uses only Antigravity, found `.claude/{skills,agents,commands,rules}` links and `.claude/settings.local.json` in a project after its first `/init-project` (BUG-113). Steps 5, 5b (items 2–3) and 6 of the skill ran unconditionally. Step 1 now reads the install manifest (`grep -q '"vendor": "antigravity"' ~/.config/sh-agent/state.json`); on a match every `.claude/` write is skipped and the `.agents/memory` store is still made; any other result keeps the old behavior; an update never removes an existing `.claude/` and names it for the operator. The operator chose the manifest as the signal and a by-hand cleanup for existing projects. The operator's audit request of the same class (a Claude-only path assumed without checking the vendor) found `/dream` Step 1 resolving the store only through `.claude/settings.local.json`, which the fix alone would have broken on such a machine: new item 0 takes `<root>/.agents/memory` there, and its global-context read names the Antigravity path too. Every other site names both vendors or only reads; the root-step permission record in `.claude/settings.json` (`context/dev-setup.md:8`, `skills/docker/SKILL.md:108`) is left as a safe fallback (an agy-only project hands a root step back). The context line for this pushed the Antigravity rules block over budget (BUG-114, found by review); shortened to 21998 of 22000 bytes.

## Type
bug (BUG-113, BUG-114)

## Files changed
- `skills/init-project/SKILL.md` — Step 1 "Claude Code or not"; Steps 5, 5b, 6, 7 and the intro follow it
- `skills/dream/SKILL.md` — Step 1 item 0; Step 2 global context on Antigravity
- `context/AGENTS.md:38` — `.claude/*` linked by `/init-project` (Claude Code only); "in a project, " cut for the budget
- `init_project_test.go` — `TestInitProjectVendorCheckMatchesManifest` (new)
- `.agents/docs/specs/sh-agent/memory-loop.md`, `README.md`, `.agents/docs/deploy/local.md` — the vendor decides whether projects get `.claude/`
- `.agents/docs/testing/skill-matrix.md` — rows `init-agy-only`, `dream-agy-only`
- `.agents/docs/bugs/BUG-113-*.md`, `BUG-114-*.md`, `INDEX.md`

## Agents dispatched
reviewer (`reviews/2026-09-30-1721-review-agy-only-init.md`, REQUEST CHANGES → finding 1 fixed, Outcome written), verifier (`testing/2026-09-30-1721-verify-agy-only-init.md`, PASS; a first run was stopped when the context line had to change), documentation (README, `deploy/local.md`; the README paragraph rewritten by the main agent after the agent's edit left a sentence broken), bugs (BUG-113, BUG-114). No security review: no auth, secrets, parsing or dependency touched.

## Verification
`go test -race -count=1 ./...` 375 pass, 0 fail. Antigravity 1.2.13 in a sandbox home installed `--vendor agy`: pre-fix `/init-project` made `.claude/` (the operator's symptom); post-fix `/init-project` made none and `/dream` reported `store: …/.agents/memory (1 memory files; …)`. Not exercised: the Claude Code path of both skills (the release gate's real-home runs).

## Docs
README § memory store, `deploy/local.md` § Configure, `specs/sh-agent/memory-loop.md`.

## Requirements
n/a
