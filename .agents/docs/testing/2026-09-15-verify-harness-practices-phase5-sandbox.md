---
owner: main
date: 2026-09-15
---
# Verification — Phase 5 checkpoint, sandbox start of each CLI   (2026-09-15)

Checkpoint of Phase 5 in `plans/2026-09-15-harness-practices-g01-g26.md`.

## Status: PASS

Both CLIs quoted Definition of Done step 2 with "or any fix under the `bug` skill". The newest Antigravity log counts 0 "Failed to parse" after I-20c changed the reviewer's frontmatter (BUG-004).

## Setup

- `go run ./cmd/sh-agent install --home /home/kesokaj/sandbox/home`: exit 0, `21 written, 239 unchanged, 0 removed, 0 kept — 260 files managed`.
- `go run ./cmd/sh-agent status --home /home/kesokaj/sandbox/home`: exit 0; `installed 2026-09-15T19:19:31+02:00 from commit b43dde1+dirty (vendor: all, 260 files, manifest ~/.config/sh-agent/state.json)`; `drift none — 259 files match the repo`. The tree was `b43dde1` with Phase 5 applied and not yet committed.
- From `/home/kesokaj/sandbox/demo-app`: `HOME=/home/kesokaj/sandbox/home claude -p "<prompt>" --permission-mode acceptEdits --output-format text`, exit 0.
- From `/home/kesokaj/sandbox/demo-app-agy`: `HOME=/home/kesokaj/sandbox/home agy -p "The project root is the current directory; ignore any other checkout. <prompt>" --mode accept-edits --print-timeout 9m`, exit 0.
- The prompt: "Do not use any tools. Quote verbatim step 2 of the Definition of Done in your global context, the sentence that says when the reviewer agent is dispatched."
- `grep -c "Failed to parse"` on the newest log, `cli-20260915_191932.log`, which this run created after `cli-20260915_191510.log`: `0`.

## Answers, verbatim

Claude Code:

```text
From the Definition of Done, step 2 (**Review**), the sentence naming when the reviewer is dispatched:

> **Review:** `reviewer` for 50+ lines, 5+ files, or any fix under the `bug` skill, given the diff only (a git range or a path list, never the intent; a saved diff file only when the tree will stand still, and the reviewer pins the blobs it read); below that read the diff yourself and say so in the changelog.
```

Antigravity CLI:

```text
"`reviewer` for 50+ lines, 5+ files, or any fix under the `bug` skill, given the diff only (a git range or a path list, never the intent; a saved diff file only when the tree will stand still, and the reviewer pins the blobs it read); below that read the diff yourself and say so in the changelog."
```

## Not examined

- Skill loads, agent dispatches and the matrix rows; Phase 6 runs them.
- One run per CLI; interactive starts; the Mac and the real home.

## Environment

Linux sandbox; Claude Code 2.1.272; Antigravity CLI 1.2.3. Raw outputs are in `.agents/scratch/phase5-sandbox/` (gitignored).
