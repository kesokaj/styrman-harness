---
name: grill-me
description: User-invoked command to launch a relentless grilling interview to sharpen a plan, decision, or design. Use when the user types /grill-me. Do NOT use for autonomous agent self-review or routine verification.
disable-model-invocation: true
---

Engage the grilling interview by loading and following the instructions in the `grilling` skill:
- On Claude Code: call the `Skill` tool with `"grilling"`.
- On Antigravity: read and apply the instructions from the `grilling` skill.

