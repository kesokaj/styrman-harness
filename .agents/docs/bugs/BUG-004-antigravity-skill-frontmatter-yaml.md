---
owner: bugs
date: 2026-09-03
id: BUG-004
status: fixed
area: skills/frontmatter
---
# BUG-004: Antigravity silently dropped the `bug` and `feature` skills

## Symptom
Sandbox `cli.log`: `Failed to parse skill file …/skills/feature/SKILL.md: failed to parse frontmatter: yaml: line 2: did not find expected ',' or ']'` (same for `bug`). Neither skill was in the model's skill list; `/bug` and `/feature` did not exist on Antigravity. Claude Code listed both, so nothing looked wrong from there.

## Repro
`sh-agent install --home <sandbox>`, run `agy -p` once, `grep 'Failed to parse skill' ~/.gemini/antigravity-cli/log/cli-*.log`.

## Root cause
`argument-hint: [feature description or #issue]` is an unquoted YAML flow sequence and ` #` starts a comment, so a strict parser sees an unterminated list. Claude Code's lenient parser and the installer's own flattening parser both accepted it. Four more skills had unquoted `[…]` hints that happened to parse.

## Fix
Quoted the six `argument-hint` values. `internal/harness/frontmatter.go`: `Lint` rejects an unquoted plain scalar that starts with `#` or has `#` after a space or tab; `LoadSkills`, `LoadAgents`, and `LoadDocs` call it, so `sh-agent install|list` fails on the offending file and names the key. Tests `TestLintRejectsUnquotedComment`, `TestLoadSkillsRejectsUnquotedComment`. Goldfish dispatch skipped: the location was mechanically obvious from the log line.

## Prevention
The lenient parser is not the reference; the vendors' strict parsers are. After any frontmatter change, install into the sandbox home, run each CLI once, and grep its log for parse errors — `sh-agent status` says the files match, not that the vendor read them.

## Related
`.agents/docs/testing/2026-09-03-verify-sandbox.md`
