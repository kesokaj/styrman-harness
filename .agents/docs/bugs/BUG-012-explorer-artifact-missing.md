---
owner: bugs
date: 2026-09-13
id: BUG-012
status: fixed
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
CONFIRMED, two mechanisms, from transcripts (goldfish brief
`.agents/docs/plans/2026-09-14-bug-012-diagnosis.md`), refuting the earlier
"narrow allow-list" hypothesis (the same narrow list produced a brief under a
neutral prompt in `bug012-before`; the verifier's Bash denials there came from
compound commands like `go build …; echo …` that no single allow-list pattern
matches — a headless caveat, unrelated to this bug). (A) Dispatcher: when the
user's prompt said "change nothing" the MAIN agent dispatched the explorer as
"read-only recon" and told it not to write its brief, and the explorer
complied — seen in the 2026-09-13 smoke run and the 2026-09-14 controlled
reproduction (`bug012-before2`: "I told the explorer not to write its usual
file under plans/"). (B) In the `/review --project --all` run of 2026-09-13,
explorer's own transcript shows two Read calls, no Write call, and no
denial, yet its reply named a brief path — a fabricated write.

## Fix
`context/AGENTS.md`: an agent's artifact (brief, log, report, record) is
trail, not a change — "change nothing", "apply nothing", or "report only"
from the user never cancels it, and a dispatch never tells an agent to skip
its file (BUG-012). `agents/explorer.md`: the brief is trail, write it even
when the dispatch says read-only or change nothing, write it with the Write
tool, name its path only after that call succeeded, never report a file not
written. `skills/review/SKILL.md`: every dispatch carries "Write your
artifact with the Write tool; report only means no source edits." Verified:
the exact failing prompt with the same narrow allow-list, twice on Claude
Code and once on Antigravity 1.2.2, each left the brief on disk with source
untouched (`.agents/scratch/bug012/after*`; verifier log
`testing/2026-09-14-verify-bug-012-fix.md`).

## Prevention
A user's "change nothing" scopes source, never the trail — `context/AGENTS.md`
now says so explicitly; an agent names an artifact only after its write
succeeded (never report a file you did not write); a dispatcher that needs an
artifact checks the path exists (already the pattern in `/review`).

## Related
BUG-013 (a second failure of trusting a report over checking the artifact);
the goldfish diagnosis brief `.agents/docs/plans/2026-09-14-bug-012-diagnosis.md`;
`agents/explorer.md` gained the early-write section on 2026-09-13 (improvement
review item 7) — this bug shows the claim of a write must also be verified,
not just the instruction to write early. Claude's built-in Write-tool
guidance ("never create documentation files unless explicitly requested") is
the seam the explicit context/agent instruction closes.
