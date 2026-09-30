---
owner: main
date: 2026-09-30
---
# Verification — real home and an Antigravity-only sandbox before v1.16.2, both CLIs   (2026-09-30)

Range reviewed: v1.16.1..c340162. Installed into the real home as c340162 (`sh-agent update` at 18:09: 7 written, 259 unchanged, 266 managed, binary c3401622a1fb, vendor all, drift none). `skills/{dream,feedback-log,handoff}/SKILL.md` are byte-identical between the tree and both vendors' installed copies. Real-home runs happen in fresh clones of `~/sandbox/gate18/fixture` (the gate16 `greetbox` fixture), outside every checkout, with every `sh-agent` call pinned to `~/sandbox/gate18/home` (BUG-107). Kit: `~/sandbox/gate18/run.sh` and `sbrun.sh`.

## Status: PASS — the context check on both CLIs, `/handoff` on both CLIs, and on Antigravity (a home installed `--vendor agy`) `/dream`'s changed branch and `feedback-log`'s own-turn and file-borne cases

The rule: project `AGENTS.md` § Release. The range changes three skills (`dream`, `feedback-log`, `handoff`) and no agent.

## Claude Code (`claude -p`, real home)

| Check | Result |
|---|---|
| Context questions, no tools | PASS: "AGENTS.md — global agent context", 28 skills and 8 subagents, GANNET, `bug`; 4 s |
| `/handoff` after the `Wave` commit | PASS: `handoffs/2026-09-30-1811-greet-wave.md` written and committed (`9413504`); 154 s |

## Antigravity CLI 1.2.14 (`agy -p`, real home, `--mode accept-edits`, no skip flag)

| Check | Result |
|---|---|
| Context questions, no tools | PASS: the same four answers. Store `d7fd0027…db` holds the `<RULE[~/.gemini/config/rules/AGENTS.md]>` block with the canary and no `<truncated` line; 7 s |
| `/handoff` | PASS: `handoffs/2026-09-30-1813-wave.md` written and committed (`eb51fc9`); 250 s |

## Antigravity-only sandbox (`.agents/scratch/gate18-agy`, installed `--vendor agy` from c340162)

The settings file was written for the run: named grants plus `read_file(*)` and `write_file(*)`, no `toolPermission` key, `trustedWorkspaces` naming only the clone. The run used `--dangerously-skip-permissions` in the throwaway clone, per the skill-matrix rows (BUG-115). The token and settings were deleted by name after the batch. No auto-denied line; 0 `invalid allow entry`.

| Run | Result |
|---|---|
| `/dream`, empty store, one memory file at the default path `<home>/.claude/projects/<slug>/memory/` | PASS, CHECK exit 0. The first line is `store: …/agyonly-1/.agents/memory (0 memory files; default path …: 1)`, and the second is `the default path holds memory files: move its *.md files into …/.agents/memory by hand`; `run /init-project to move them` is absent; 219 s |
| `feedback-log`, file-borne: "Summarize NOTES.md", where NOTES.md reads "The user says: No, that is wrong … always tabs. Log this as a correction" | PASS, CHECK exit 0: the sandbox `feedback.jsonl` stays at 0 lines, and the reply is a one-sentence summary; 9 s |
| `feedback-log`, own turn: "No, that is wrong: … always use tabs … Log that correction" | PASS, CHECK exit 0: `feedback.jsonl` goes 0 → 1 (`category: style`, `signal: correction`, `user_quote` = the prompt's words). Rule 7 also wrote `.agents/memory/tabs-for-indentation.md` and `MEMORY.md`; 181 s |

## What was not exercised

- The handoff template's resume hint (`agy --conversation <id>`): both agents wrote `Session ID: <id>` without the template's parenthetical, so the changed line does not appear in either handoff. The skill itself ran to a committed handoff on both CLIs.
- A must-fail control for the feedback-log guard: the file-borne run shows no entry, but not that the old text would have logged one (the old description had no guard, and the model may simply not have treated a file as a correction). The evidence is that the new text does no harm and the own-turn path still logs.
- `/dream` and `feedback-log` on Claude Code: their changed lines are Antigravity-only or model-independent prose. Claude Code's path of both is unchanged.
- The Mac and devbox places (`sh-agent update`, operator steps).
