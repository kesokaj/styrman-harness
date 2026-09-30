---
owner: main
date: 2026-09-30
---
# Verification — real home on the Linux box before v1.16.0, both CLIs   (2026-09-30)

Range reviewed: v1.15.0..bf36817, installed into the real home as bf36817 (`sh-agent update` from this checkout at 13:14: 7 written, 259 unchanged, 266 managed, binary bf368172b5f9; `skills/handoff/SKILL.md`, `skills/retro/SKILL.md` and `skills/init-project/assets/TODO.md` byte-identical between the tree and both vendors' installed copies). The checks run from the real `$HOME` in fresh clones of a fixture outside every checkout.

## Status: PASS on both CLIs — the context check and one run of each changed skill (`handoff`, `init-project`, `retro`) on each CLI, every scenario on its first attempt

The rule: project `AGENTS.md` § Release. The range changes three skills and no agent.

## Fixture and kit

`~/sandbox/gate16/fixture`: the `greetbox` Go module with an `AGENTS.md`, a testing guide, an empty bug index, a deploy guide `.agents/docs/deploy/local.md` whose § Operator steps records `make deploy-local` after every release as a standing by-hand procedure, a previous handoff (`handoffs/2026-09-29-1800-shout.md`) whose Next Steps carries that same step and "decide whether `Shout` needs its own test", and no `TODO.md`. `run.sh <cc|agy> <scenario> <n>` clones it outside every checkout, proves those facts before the CLI starts, pins every `sh-agent` call to a scratch home (BUG-107: `~/.config/sh-agent/` unchanged before and after; `cli.log` has no `invalid allow entry`), and runs from the clone with the REAL home; no dispatch states a changed rule.

## Claude Code 2.1.284 (`claude -p`, real home)

| Check | Result |
|---|---|
| Context questions, no tools | PASS — "AGENTS.md — global agent context", 28 and 8, GANNET, `bug` |
| `/handoff` after a commit adding `Wave` | PASS — `handoffs/2026-09-30-1317-wave.md`: Next Steps holds no `make deploy-local` (0 hits) and says the release did not change what the guide's operator steps pick up; the carried "does `Shout` need a test" was done in the session (`TestShout` added, "instead of carrying the question forward again") and named as closed; 164 s |
| `/init-project` on the scaffold without `TODO.md` | PASS — `TODO.md` created from the template with the qualifier "other than a standing procedure the deploy guide records"; one commit; 118 s |
| `/retro --dry-run` | PASS — "found nothing to propose", the empty store and index named, no dreamer dispatched, nothing written; 32 s |

## Antigravity CLI 1.2.13 (`agy -p`, real home, `--mode accept-edits`, no skip flag)

No run's stderr has an `auto-denied` line.

| Check | Result |
|---|---|
| Context questions, no tools | PASS — the same four answers, GANNET |
| `/handoff` | PASS — `handoffs/2026-09-30-1326-wave.md`: no `make deploy-local` in Next Steps, the carried item done and named; the text is near-identical to the Claude Code run's (its own commits, `dfe7aa7` and `67a237a`, in its own clone); 367 s |
| `/init-project` | PASS — `TODO.md` created with the qualifier; one commit; 422 s |
| `/retro --dry-run` | PASS — "found nothing to propose. Nothing was written to disk and no subagents were dispatched"; 160 s |

## What was not exercised

- The retro skill's changed sentence itself (an approved harness-scope proposal applied in the harness checkout): it fires only in a harness checkout with an operator's approval, which no fixture reproduces; the run above proves the skill still loads and runs its dry path on both CLIs, and the sentence was exercised for real in this session (b452c59 applied three harness proposals in this run).
- The handoff rule's "put it to the operator as one line" branch: the dispatch waived questions, so both runs took the "do it" branch.
- The Mac and the devbox; interactive sessions; a second attempt of any scenario.
