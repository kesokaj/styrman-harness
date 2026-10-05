---
owner: main
date: 2026-10-05
---
# Verification — real-home install from the checkout, smoke on both CLIs (Linux box)   (2026-10-05)

Range reviewed: 0be3b97 installed (`sh-agent status`: from commit 0be3b97, vendor all, 1,790 files, drift none).

## Status: PASS

## Install

The operator chose option A in their own turn ("go option a"): this machine moves from a release place (v2.3.2, embedded) to a checkout place. `go run ./cmd/sh-agent install` refused first ("this home was installed from the release binary v2.3.2, not from a checkout; pass --repo …"); `go run ./cmd/sh-agent install --repo /home/ostling-simon/projects/styrman-harness` → "19 written, 1771 unchanged … 1790 files managed". Post-install: `hash -r`; `type sh-agent` → `~/.local/bin/sh-agent`; `sh-agent version` → `0be3b977b74f`; `status` → `repo … (main @ 0be3b97, clean)`, `binary 0be3b977b74f`, no `shadowed` line, `drift none`. `rules/critical-thinking.md` byte-identical to `~/.claude/rules/` and `~/.gemini/config/rules/`.

## Runs (real `$HOME`, cwd a fresh clone of `~/sandbox/ct/fixture` under `~/sandbox/smoke/`, outside every checkout)

| Check | Claude Code | Antigravity |
|---|---|---|
| Context, no tools: canary, the rule file, its line on corrections, the approval phrases | PASS — GANNET; `critical-thinking.md` from `~/.claude/rules/`; no apologies, re-check evidence; `go all approved`, `go approved 1-3, skip 5` | PASS — the same four answers; store: `<RULE[…/rules/AGENTS.md]` and `<RULE[…/rules/critical-thinking.md]`, canary line present, 0 `<truncated` |
| Skill auto-load (`main.go` review, no edits) | PASS — `golang`; tree unchanged | PASS — `golang` (with `review-standards`); tree unchanged |
| Agent dispatch | PASS — `verifier` on an uncommitted comment: log `testing/2026-10-05-verify-main-doc-comment.md`, `Status: FAIL` (no test covers `main`: the verifier's rule, expected) | PASS — `explorer` (`invoke_subagent` in the store): brief `plans/2026-10-05-recon-repo-map.md` |

Antigravity's tool checks ran with `--dangerously-skip-permissions` inside the throwaway clones (operator approval of 2026-10-05, test runs only); the real `~/.gemini/antigravity-cli/settings.json` does not exist. All replies used the new reply shape (explanation, recommendation, numbered table, `go approved` line).

## Not verified here

The Mac and the other places (each moves when the operator runs `sh-agent update` there; a release place gets the rule only from the next release). Whether the sessions' resume step registered the smoke clones in the real global search index (`~/.config/sh-agent`) was not checked.
