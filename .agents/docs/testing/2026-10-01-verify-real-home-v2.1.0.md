---
owner: main
date: 2026-10-01
---
# Verification — real-home install of the v2.1.0 gated head and the guide § 5 smoke on both CLIs (2026-10-01)

Range reviewed: v2.0.0..159585a — one install into the real home, at the gated head 159585a, on the Linux box (`go run ./cmd/sh-agent install` per `deploy/local.md` § Build and install, as the deploy step of the release the operator asked for), then the six smoke checks of guide § 5.

## Status: PASS — 6 of 6 smoke checks on both CLIs at 159585a

Before the install the real home held the build of 8b6c6de: `sh-agent status` reported `drift 1518 missing, 16 differ, 0 stale` (the catalog's files and the skills changed since). Install: `1535 written, 249 unchanged, 0 removed, 0 kept — 1784 files managed`; `sh-agent status`: `repo has 35 skills, 8 agents, 0 commands, 0 rules`, `drift none — 1783 files match the repo`, `dream never run; timer not installed`; `sh-agent version` after `hash -r`, the bare name: `sh-agent 159585aad20d`. `diff -rq skills/google-catalog` against `~/.claude/skills/google-catalog` and `~/.gemini/config/skills/google-catalog`: empty.

| CLI | check | result | seconds |
|---|---|---|---|
| Claude Code | context | pass=1, denied 0 | 7 |
| Claude Code | agent | pass=1, denied 0 | 54 |
| Antigravity | context | pass=1, denied 0 | 8 |
| Claude Code | skill | pass=1, denied 0 | 70 |
| Antigravity | skill | pass=1, denied 0 | 46 |
| Antigravity | agent | pass=1, denied 0 | 134 |

Checks (the kit `kit-realhome.sh`, each in a fresh clone of the demo service, the CLI running from the real home): `context` — five facts from the global context with no tool (both CLIs answered 35 skills, the eight agents, and the canary word); `skill` — a Go review that names the `golang` skill and changes nothing; `agent` — a `verifier` dispatch whose log is on disk with a Status line. The kit's preamble refused nothing (0 "kit: FAIL" lines).

## Not checked

The Mac (no access this session). The nightly dream timer on this box: it is not installed, so the one-time remove and install the release notes name does not apply here. An interactive session on either CLI (does it prompt before reading under the installed skill directory): the operator's.
