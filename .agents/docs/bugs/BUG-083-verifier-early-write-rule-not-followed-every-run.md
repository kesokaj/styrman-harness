---
owner: bugs
date: 2026-09-18
id: BUG-083
status: open
area: agents/verifier
---
# BUG-083: rewritten "write the artifact EARLY" rule in verifier.md followed in only 1 of 3 runs

## Symptom
`agents/verifier.md` § Write the artifact EARLY was rewritten (retro proposal 6, `retros/2026-09-18-retro.md`) to say "Write your log file after the FIRST command you run … with a Status line of IN PROGRESS — and update it after EVERY command". Behavioural run per project `AGENTS.md` § Testing: `claude -p --agents <json from the working-tree agents/verifier.md>` in a scratch clone `.agents/scratch/verifier-run`, three runs, dispatch text NOT repeating the early-write sentence (definition alone under test). 1 of 3 runs wrote the log with IN PROGRESS after the first command and updated it after the second; 2 of 3 ran every command, then wrote once at the end with no IN PROGRESS. Order read from subagent transcripts under `~/.claude/projects/…verifier-run/*/subagents/`. A first attempt without `--allowedTools` was blocked at the permission gate and is not evidence.

## Repro
1. `claude -p --agents <json built from working-tree agents/verifier.md> --allowedTools ...` in `.agents/scratch/verifier-run`, three runs.
2. Read each run's subagent jsonl for the order of `Bash` and `Write` tool_use blocks.
3. Result: 1/3 IN PROGRESS-then-update; 2/3 single end-of-run write.

## Root cause
A rule inside an agent's definition is followed in some runs and not others on Claude Code — the same mechanism as BUG-062, BUG-063 and BUG-064 (dispatch text made the verifier write its guide every run).

## Fix
Not fixed in the definition — the reworded rule is kept as applied. The enforcement is the existing context rule: `context/AGENTS.md` § Work already says a dispatch to `verifier` repeats "write the artifact early"; that dispatch text must now be phrased to say "write the log after the FIRST command, IN PROGRESS, and update it after every command."

## Prevention
A change to an agent definition is verified by runs whose dispatch does NOT repeat the changed rule, and the count is reported as n of N, never "works". The sentence that must hold every run goes in the dispatch text (BUG-064) — the definition alone is not reliable enforcement on Claude Code.
`Diagnosis missed:` n/a.

## Related
BUG-062, BUG-063, BUG-064 — same mechanism (definition-only rules are unreliable on Claude Code; the enforcement belongs in dispatch text). Third+ occurrence of this exact prevention rule; already graduated once into `context/AGENTS.md` § Work per BUG-064 — this record shows the graduated rule was not applied when the dispatch text for this change should have repeated it.
