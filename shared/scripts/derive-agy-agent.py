#!/usr/bin/env python3
"""Derive an Antigravity CLI agent .md from a Claude Code agent .md.

Claude Code is the source of truth (claude/agents/*.md). This script produces
the agy-cli flavor per the Antigravity custom-agent spec:
  - keeps portable frontmatter keys: name, description
  - adds  subagent: true / mainAgent: false  (invokable via invoke_subagent,
    never selectable as the primary agent)
  - maps Claude model aliases to agy model tiers (opus->pro, haiku->flash;
    sonnet and anything else -> omitted, i.e. agy's default "inherit")
  - maps the tools allowlist to agy tool names (Read->view_file, Grep->grep_search,
    Glob->find_by_name, Bash->run_command, Write->write_to_file,
    Edit->replace_file_content); an unmapped tool is an error, never a guess —
    a misspelled agy tool can hang the subagent
  - maps `skills:` to ~/.gemini/config/skills/<name> paths
  - drops Claude-only keys (memory, permissionMode, ...)
  - rewrites ~/.claude/skills/ paths to ~/.gemini/config/skills/

Usage: derive-agy-agent.py <claude-agent.md> <output.md>
"""

import re
import sys
from typing import Optional

MODEL_MAP = {"opus": "pro", "haiku": "flash"}
TOOL_MAP = {
    "Read": "view_file",
    "Grep": "grep_search",
    "Glob": "find_by_name",
    "Bash": "run_command",
    "Write": "write_to_file",
    "Edit": "replace_file_content",
}
AGY_SKILLS_DIR = "~/.gemini/config/skills/"


def main() -> None:
    if len(sys.argv) != 3:
        sys.exit(__doc__)
    src, dst = sys.argv[1], sys.argv[2]

    with open(src, encoding="utf-8") as f:
        text = f.read()

    m = re.match(r"^---\n(.*?)\n---\n", text, re.S)
    if not m:
        sys.exit(f"{src}: no YAML frontmatter block found")
    frontmatter, body = m.group(1), text[m.end():]

    def field(key: str) -> Optional[str]:
        fm = re.search(rf"^{key}:\s*(.+)$", frontmatter, re.M)
        return fm.group(1).strip() if fm else None

    name, description = field("name"), field("description")
    if not name or not description:
        sys.exit(f"{src}: frontmatter must contain single-line name and description")

    out = [f"name: {name}", f"description: {description}",
           "subagent: true", "mainAgent: false"]
    tier = MODEL_MAP.get(field("model") or "")
    if tier:
        out.append(f"model: {tier}")
    tools = field("tools")
    if tools:
        mapped = []
        for tool in [t.strip() for t in tools.split(",") if t.strip()]:
            if tool not in TOOL_MAP:
                sys.exit(f"{src}: tool {tool!r} has no agy mapping — add it to TOOL_MAP")
            mapped.append(TOOL_MAP[tool])
        out.append("tools: [" + ", ".join(mapped) + "]")
    skills = field("skills")
    if skills:
        paths = [AGY_SKILLS_DIR + s.strip() for s in skills.split(",") if s.strip()]
        out.append("skills: [" + ", ".join(paths) + "]")

    body = body.replace("~/.claude/skills/", "~/.gemini/config/skills/")

    with open(dst, "w", encoding="utf-8") as f:
        f.write("---\n" + "\n".join(out) + "\n---\n" + body)


if __name__ == "__main__":
    main()
