---
owner: bugs
date: 2026-09-13
id: BUG-012
status: open
area: agents/explorer
---
# BUG-012: explorer claims a written brief that does not exist on disk

## Symptom
Twice today, headless in the sandbox (real `HOME=/home/kesokaj/sandbox/home`,
`claude -p … --permission-mode acceptEdits`, `Agent` + a Bash allow-list), the
`explorer` agent returned its brief in chat and — the second time — stated it
had written `plans/2026-09-13-recon-project.md`. No such file existed after:
`find . -name '*recon*'` empty, `git status` clean of it.

## Repro
1. `/home/kesokaj/sandbox/runs/claude-smoke`, ~22:05: "map this repository's
   packages and their imports … change nothing." No brief file appeared.
2. `/home/kesokaj/sandbox/runs/claude-project`, ~23:05, `/review --project --all`
   end-to-end: "apply nothing." Recorded in that clone's
   `.agents/docs/reviews/2026-09-13-review-project-merged.md` as
   "explorer — artifact missing." Same session's `reviewer`, `security`,
   `verifier` wrote their artifacts fine under the same permissions.
   Antigravity's explorer in the same smoke check also left no brief.
3. `/home/kesokaj/sandbox/runs/claude-project2`, ~21:20-21:34Z (2026-09-13),
   `/review --project --all` again — this time the brief
   `plans/2026-09-13-recon-project.md` WAS written. The only run difference:
   `--allowedTools` additionally permitted `Bash(bash:*)`, `Bash(printf:*)`,
   `Bash(echo:*)`, `Bash(cd:*)`, `Bash(set:*)`. Runs 1-2 allowed only `git`,
   `go`, `ls`, `cat`, `date`, `mkdir`, `wc`, `grep`, `find`, `pwd`, `head`,
   `tail`, `test` (plus `Agent`). In the failing run-2 project, `verifier`
   also "lost Bash" under the narrow list — a second agent affected by the
   same narrow allow-list, not explorer-specific.

## Root cause
NOT CONFIRMED — no controlled retest isolating the variable yet. Leading
candidate, promoted by evidence 3: (b) explorer writes its brief through a
shell command (heredoc via `bash -c` / `printf` / `cat >`) rather than the
Write tool; a narrow headless `--allowedTools` list silently denies that
Bash invocation, after which the agent still reports the write as done.
`reviewer`/`security`/`verifier` write their artifacts with the Write tool
and were unaffected in runs 1-2 (verifier's own Bash loss in run 2 was a
separate symptom of the same narrow list, not a Write-tool failure). (a)
"change nothing" misread and (c) working-directory mismatch are now
secondary candidates. Still not confirmed: needs a retest with the narrow
allow-list plus transcript inspection of explorer's actual write call.

## Fix
Not fixed. Candidate fix, not yet applied: change `agents/explorer.md` to
write its brief with the Write tool instead of a shell heredoc, then retest
under the narrow allow-list to confirm.

## Prevention
Until diagnosed: headless dispatch allow-lists include `Bash(bash:*)` and
`Bash(printf:*)` (or wider) as a workaround; a dispatcher that needs the
explorer's brief checks the path exists before trusting the reply (`/review`
already records "artifact missing" — follow that pattern elsewhere). Next
investigation: dispatch explorer under the narrow allow-list with the
transcript inspected for the actual write call, and once the mechanism is
confirmed, change `agents/explorer.md` to use the Write tool and retest.

## Related
`agents/explorer.md` gained the early-write section today (improvement
review item 7); this bug shows the claim of a write must also be verified,
not just the instruction to write early.
