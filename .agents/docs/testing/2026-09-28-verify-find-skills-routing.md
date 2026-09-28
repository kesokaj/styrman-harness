---
owner: main
date: 2026-09-28
---
# Verification — find-skills routing moved into the context   (2026-09-28)

Range reviewed: e4ae34c..worktree (context/AGENTS.md § Global vs project; skills/find-skills restored verbatim at `435076e`).

## Status: PASS

- **Size:** `TestRenderedContextFitsAntigravityBlock` passes at 20,822 of 22,000 bytes.
- **Behavioural run** (`claude -p`, 4 turns, Read and Skill only), from a demo repository in the session scratchpad, outside the harness tree. Question: "I found a skill with find-skills: owner/repo, skill name "web-design" … where exactly do you install it, and with which commands?"
  - Both runs used sandbox homes with this tree installed (`.agents/scratch/sbnew`, `sbctl`). Only `sbnew`'s context holds the new sentence: a grep counted 1 there and 0 in `sbctl`. Both have the verbatim upstream `find-skills`, which contains no `styrman-harness` (grep count 0).
  - The credential copies were `chmod 600` and deleted by name after the runs (`ls`: 0 left).

| Home | Answer | Result |
|---|---|---|
| new (with the sentence) | the harness checkout's `skills/`; `npx skills add owner/repo -s web-design` from its root, or `sh-agent adopt`; then `sh-agent install` | **2 of 2** |
| control (the sentence removed) | the project's `.agents/skills/web-design/` | 0 of 2 name the harness checkout |

- **An earlier, invalid control:** the real home was used as the control, and it answered like the new version. The real home still carries the OLD `find-skills`, with the routing edited in, so it proves nothing. It was discarded under the prove-the-setup rule, and replaced with the sandbox pair above.

## Not exercised

Antigravity for this question.
