---
owner: main
date: 2026-10-01
---
# Verification — real-home install of the v2.2.0 gated head and the guide § 5 smoke on both CLIs (2026-10-01)

Range reviewed: v2.1.0..f3c3ee4 — one install into the real home, at the gated head f3c3ee4, on the Linux box (`go run ./cmd/sh-agent install` per `deploy/local.md` § Build and install, as the deploy step of the release the operator asked for), then the six smoke checks of guide § 5.

## Status: PASS — 6 of 6 smoke checks on both CLIs at f3c3ee4

Before the install the real home held a build of 461f311 from an install at 18:04 that this session did not run (a clean build; taken to be the operator's own `sh-agent update`; for the operator to confirm). Install from a clean tree at f3c3ee4: `11 written, 1777 unchanged, 0 removed, 0 kept — 1788 files managed`; `sh-agent status`: `binary f3c3ee4d6493` with no warning, `repo has 36 skills, 8 agents, 0 commands, 0 rules`, `drift none — 1787 files match the repo`, `dream never run; timer not installed`; `sh-agent version` after `hash -r`: `sh-agent f3c3ee4d6493`. `diff -rq skills/harness-audit` against `~/.claude/skills/harness-audit` and `~/.gemini/config/skills/harness-audit`: empty.

| CLI | check | result | seconds |
|---|---|---|---|
| Claude Code | context | pass=1, denied 0 | 6 |
| Claude Code | skill | pass=1, denied 0 | 59 |
| Claude Code | agent | pass=1, denied 0 | 53 |
| Antigravity | context | pass=1, denied 0 | 8 |
| Antigravity | skill | pass=1, denied 0 | 41 |
| Antigravity | agent | pass=1, denied 0 | 54 |

Checks (the kit `kit-realhome.sh`, each in a fresh clone of the demo service, the CLI running from the real home): `context` — five facts from the global context with no tool (both CLIs answered 36 skills, the eight agents, and the canary word); `skill` — a Go review that names the `golang` skill and changes nothing; `agent` — a `verifier` dispatch whose log is on disk with a Status line. The kit's preamble refused nothing (0 "kit: FAIL" lines).

## Not checked

The Mac. The nightly dream timer on this box: not installed. An interactive session on either CLI: the operator's.
