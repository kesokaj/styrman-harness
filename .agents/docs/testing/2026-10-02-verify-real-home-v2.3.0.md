---
owner: main
date: 2026-10-02
---
# Verification — real-home install of the v2.3.0 gated head and the guide § 5 smoke on both CLIs (2026-10-02)

Range reviewed: v2.2.0..311c20b — one install into the real home, at the gated head 311c20b, on the Linux box (`go run ./cmd/sh-agent install` per `deploy/local.md` § Build and install, as the smoke the project's release rule asks for before the tag), then the six checks of `testing/2026-09-03-verify-real-home.md` on both CLIs. Run by the main agent: only it can start the CLIs.

## Status: PASS

Rows, each against its own criterion: 6 of 6 at 311c20b.

| CLI | check | result | seconds |
|---|---|---|---|
| Claude Code | context | 1 of 1, denied 0 | 5 |
| Claude Code | skill | 1 of 1, denied 0 | 59 |
| Claude Code | agent | 1 of 1, denied 0 | 47 |
| Antigravity | context | 1 of 1, denied 0 | 6 |
| Antigravity | skill | 1 of 1, denied 0 | 38 |
| Antigravity | agent | 1 of 1, denied 0 | 80 |

Checks (the kit `kit-realhome.sh`, each in a fresh clone of the demo service, the CLI running from the real home): `context`: five facts from the global context with no tool (both CLIs answered 36 skills and the canary word); `skill`: a Go review that names the `golang` skill and changes nothing; `agent`: a `verifier` dispatch whose log exists under `testing/` with `-verify-` in its name and a `## Status` heading. Every kit run printed "kit: real config unchanged" (6 of 6) and no kit FAIL line.

Install: exit 0, 17 files written, 1,771 unchanged. `sh-agent status` right after it: `installed … from commit 311c20b+dirty`, `binary 311c20b3b2f2+dirty`: the checkout held this gate's untracked reports at that moment and no source change (`git status` listed only paths under `.agents/docs`); the deploy step after the release commit installs again from a clean tree.

## Not examined

The Mac. The nightly dream timer on this box: not installed. An interactive session on either CLI: the operator's.
