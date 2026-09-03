---
owner: main
date: 2026-09-03
---
# Verification — live sandbox, both CLIs headless   (2026-09-03)

Range reviewed: 87a862d..working tree (BUG-003, BUG-004, BUG-005 fixes; deploy guides)

## Status: PASS with three harness bugs found and fixed

## Setup

`HOME=/home/kesokaj/sandbox/home` logged in to Claude Code 2.1.259 and the Antigravity CLI 1.1.25; harness installed with `sh-agent install --home /home/kesokaj/sandbox/home` (254 → 260 files, no drift). Two copies of a three-file Go service: `/home/kesokaj/sandbox/demo-app` (Claude Code) and `demo-app-agy` (Antigravity). Headless invocations: `claude -p … --permission-mode acceptEdits --allowedTools …` and `agy -p … --mode accept-edits` with `permissions.allow` in `~/.gemini/antigravity-cli/settings.json`.

## Runs

| # | Run | Claude Code | Antigravity |
|---|---|---|---|
| 1 | Context questions, no tools (first heading, skill count, subagents, scratch rule, DoD order, which skill for a failing test / a new endpoint) | PASS — all five right; flagged the scratch contradiction → BUG-003 | PASS — all five right; `bug` and `feature` absent from its skill list → BUG-004 |
| 2 | `/init-project` with answers in the prompt | PASS — 21 files, `AGENTS.md`/`README.md` filled correctly, commit; `.claude/` symlinks refused by Claude Code's guard on that path, reported with the four lines to run | PASS — 25 files incl. symlinks, commit; it read the sibling checkout for the template text, so template filling is proven only on Claude Code |
| 3 | Task: `GET /version` through the Definition of Done | PASS — classified task; `verifier` blocked by worktree isolation → BUG-005, verified in main thread and again post-commit; `security` FINDINGS (1 medium, 4 low), two fixed, two rebutted in writing into `TODO.md`; `reviewer` (5 files) APPROVE WITH NITS; `documentation`, `bugs`, `changelog`; two commits | PASS — classified task; `verifier` PASS, `reviewer` APPROVE (143 lines), `documentation`, `changelog`; testing/reviews/changes records; one commit |
| 4 | `verifier` dispatched on an uncommitted edit after the BUG-005 fix and reinstall | PASS — ran in `/home/kesokaj/sandbox/demo-app` itself, `git status --short` showed the edit, verdict PASS, log written into the project's `testing/` | not applicable (Antigravity ignores `isolation`) |
| 5 | Project-layer discovery through the `.claude/*` symlinks (a throwaway `.agents/skills/demo-marker/SKILL.md`, no tools) | PASS — the skill was listed and its description quoted verbatim; symlinks stay, copies are not needed | not run — Antigravity reads `.agents/` natively |

Run 3 code re-verified by hand in both repos: `go vet ./... && go test -race ./...` → `ok` for `cmd/server`, `internal/greet`, `internal/version`.

## Harness bugs found

- **BUG-003** two scratch rules in the rendered context (`dev-setup.md` vs the hard rule). Fixed before run 2.
- **BUG-004** `argument-hint: [… #issue]` unparseable for Antigravity's YAML; `bug` and `feature` silently dropped. Six hints quoted; `sh-agent` now lints frontmatter (`TestLintRejectsUnquotedComment`, `TestLoadSkillsRejectsUnquotedComment`). Re-run of run 1 on Antigravity lists both skills; newest `cli-*.log` has zero `Failed to parse skill` lines.
- **BUG-005** `isolation: worktree` on `verifier` — Claude Code checks out `HEAD` in a throwaway worktree, so the verifier cannot see uncommitted work and its log is lost. Key removed.

## Vendor behaviour worth knowing (not harness bugs)

- Claude Code headless: reads outside the cwd need `--add-dir`; writes under the project's `.claude/` are refused by a guard (`mkdir` included, so copying instead of symlinking would not help; interactively it is one approval); subagents that use `Agent` need `"Agent"` in `--allowedTools`. `--dangerously-skip-permissions` is blocked by the auto-mode classifier on this box.
- Antigravity headless: one auto-denied tool ends the run with empty stdout and a stderr notice naming the rule; `--add-dir` makes the workspace ambiguous and the agent wanders (it ran `ps aux` to find its own args); the trusted workspace `/home/kesokaj/sandbox` let it read the sibling repo.
- Antigravity hides slash-only skills (`disable-model-invocation: true`) from the model's list, as intended; `/init-project` still expanded in print mode.

## Coverage of changed code

- `Lint`, `entries` — `TestLintRejectsUnquotedComment`, `TestLoadSkillsRejectsUnquotedComment`, existing `TestParseFrontmatter`
- `agents/verifier.md`, deploy templates, context edits — install into a scratch home (260 files, no drift), rendered context greps, run 4

## Commands

`gofmt -l .` → none; `go vet ./...` → clean; `go test -race ./...` → `ok`; `sh-agent install --home <scratch>` → `260 written`; `sh-agent status --home <scratch>` → `drift none — 259 files match the repo`.

## Not verified here

Interactive starts of both CLIs on the Mac and the real home (`/memory`, `/skills`, `/agents`); Antigravity's template filling in isolation; `/preflight` end to end (no deploy target in the demo).

## Environment

Linux (Debian 13, amd64), Go 1.24.4, Claude Code 2.1.259, Antigravity CLI 1.1.25 (Gemini 3.8 Flash).
