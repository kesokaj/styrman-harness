---
owner: main
date: 2026-09-03
---
# Sandbox verification of both CLIs, three harness bugs, deploy guides     (date: 2026-09-03)

## Summary

First live use of the rebuilt harness: both CLIs driven headless against a demo Go service from a throwaway home. Every run passed on the protocol level and surfaced three harness defects, all fixed.

- **BUG-003** — `context/dev-setup.md` still sent scratch to `/tmp`; the bullet is gone, the hard rule stands alone.
- **BUG-004** — Antigravity's YAML parser rejected `argument-hint: [… #issue]` and silently dropped the `bug` and `feature` skills. Six `argument-hint` values quoted. `internal/harness/frontmatter.go` gained `Lint` (an unquoted plain scalar with `#` is an error) and `entries` (the key walk shared with `Parse`); `LoadSkills`, `LoadAgents`, `LoadDocs` call it, so `sh-agent install|list` refuses the file and names the key.
- **BUG-005** — `agents/verifier.md` declared `isolation: worktree`; on Claude Code that is a fresh worktree at `HEAD`, blind to the uncommitted change the verifier exists to check. Key removed.
- **`init-project`** — notes on the permission a headless run needs to read the templates and on Claude Code's guard on `.claude/`; `.claude/worktrees/` added to the gitignore template.
- **Deploy guides** (operator request) — every project gets `docs/deploy/{local,remote,gcp}.md`: templates shipped by `init-project` (filled `local.md` at scaffold time; `remote.md` and `gcp.md` say *Not applicable yet* until a target exists), named in the context layout and Definition of Done step 3, `/preflight` deploys only by them, `cloud-run-deploy` reads and records the GCP defaults in `gcp.md` (the context's GCP rule now points there instead of `AGENTS.md`), the `documentation` agent owns them, `/handoff` lists them, and the project `AGENTS.md` template has a `## Deploy` section. The `documentation` agent's frontmatter rule was contradicting the context (which requires `owner:`/`date:` on every `docs/` artifact); it now matches. This repo's own guides: `local.md` (install into the real home), `remote.md` (the resync loop), `gcp.md` (not applicable, and why).

## Type

bug (BUG-003/004/005) and task (deploy guides)

## Files changed

- `internal/harness/frontmatter.go`, `internal/harness/inventory.go`, `internal/harness/harness_test.go` — lint and tests
- `skills/{audit,bug,distill,feature,preflight,research}/SKILL.md` — quoted `argument-hint`
- `agents/verifier.md` — no worktree isolation; `agents/documentation.md` — deploy guides, frontmatter rule
- `context/AGENTS.md`, `context/dev-setup.md` — deploy guides in the layout and DoD; GCP defaults recorded in `deploy/gcp.md`
- `skills/init-project/SKILL.md`, `skills/init-project/assets/AGENTS.md`, `skills/init-project/assets/agents-tree/docs/deploy/{local,remote,gcp}.md` — templates and steps
- `skills/preflight/SKILL.md`, `skills/cloud-run-deploy/SKILL.md`, `skills/handoff/SKILL.md` — deploy-guide references
- `AGENTS.md`, `.agents/docs/deploy/{local,remote,gcp}.md` — this repo's own guides
- `.agents/docs/bugs/BUG-003…`, `BUG-004…`, `BUG-005…`, `INDEX.md`; `.agents/docs/testing/2026-09-03-verify-sandbox.md`; `.agents/docs/plans/2026-09-03-rebuild.md` section 12

## Agents dispatched

`reviewer` (Sonnet, the Go diff only): one finding, a tab before `#` escaped the lint — fixed with a test. The markdown diff was read by the main thread. The sandbox runs dispatched `verifier`, `reviewer`, `security`, `documentation`, `bugs`, `changelog` inside the demo repos and are the subject of the test log).

## Verification

`gofmt`, `go vet`, `go test -race` clean; scratch-home install 260 files, no drift; both CLIs re-run after each fix. See `testing/2026-09-03-verify-sandbox.md`.

## Requirements

R2 (standards: both vendors must read every file), R3 (no scripts — the lint is in the binary), R6 (scratch home for every test).
