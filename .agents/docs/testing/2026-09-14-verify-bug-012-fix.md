---
owner: verifier
date: 2026-09-14
---
# Verification — bug-012-fix   (2026-09-14)

Range reviewed: 157934b..working tree (uncommitted)

## Status: PASS

Scope: prose-only change (3 files) verified by repo checks + read (not re-run) of prior retest evidence under `.agents/scratch/bug012/` and sandbox run directories; complete for the scope given.

## Diff examined

`git diff 157934b -- context/AGENTS.md agents/explorer.md skills/review/SKILL.md` — matches the described shape exactly:

- `context/AGENTS.md` — one sentence appended to the skill/subagent paragraph: "An agent's artifact (brief, log, report, record) is trail, not a change: 'change nothing', 'apply nothing', or 'report only' from the user never cancels it — those words cover source, config, and tests, and a dispatch never tells an agent to skip its file (BUG-012)."
- `agents/explorer.md` — one paragraph added after the "Not examined" guidance: "The brief is trail, not a change to the project: write it even when the dispatch says read-only, change nothing, or report only — those words apply to source. Write it with the Write tool, and name its path in your reply only after that call has succeeded; never report a file you did not write."
- `skills/review/SKILL.md` — one clause added to the dispatch sentence in Step 3: ", and the sentence \"Write your artifact with the Write tool; report only means no source edits.\""

No other files in the diff. `.agents/docs/plans/2026-09-14-bug-012-diagnosis.md` is untracked (new), read and consistent with the fix (root-cause candidates 1/2 in the diagnosis — model skipping Write while fabricating a path — are exactly what the new "write even when told not to, name the path only after Write succeeds" instruction targets).

## Commands

- `git diff 157934b -- '*.go'` → no output — **no Go files changed in this range**, confirmed.
- `gofmt -l .` → no output (clean).
- `go vet ./...` → clean, no output.
- `go test -race ./...` → `ok github.com/kesokaj/styrman-harness/cmd/sh-agent 1.123s`, `ok github.com/kesokaj/styrman-harness/internal/harness 1.419s`.
- `go run ./cmd/sh-agent list` → `### Skills (27)` header, 27 rows listed, `### Subagents (7)` header, 7 rows listed. Matches expected 27/7.
- `go run ./cmd/sh-agent render | grep -c BUG-012` → `1` (the render output contains the sentence exactly once, from `context/AGENTS.md`).
- Scratch install: `go run ./cmd/sh-agent install --home /tmp/sh-agent-verify-bug012` → `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`. `go run ./cmd/sh-agent status --home /tmp/sh-agent-verify-bug012` → `drift  none — 259 files match the repo` (repo reported as `157934b, uncommitted changes`; binary `157934bdd7d2+dirty`). No drift.

## Retest evidence read (not re-run)

All files under `/home/kesokaj/workspace/styrman-harness/.agents/scratch/bug012/`:

- `before2.status`: `before2 exit 0`. `before2.out` tail: explorer's corrected findings are returned inline, ending "Nothing was changed: I told the explorer not to write its usual file under `plans/`, and `git status` shows no changes." Confirms the pre-fix repro: main agent instructed the explorer not to write, and it complied — no brief exists. Checked `/home/kesokaj/sandbox/runs/bug012-before2/.agents/docs/plans/` directly: empty (command returned nothing beyond the ls, i.e., no matching brief file).
- `after.status`: `after1 exit 0`, `after2 exit 0`, `after-agy exit 0`.
- `after1.out` tail: explorer's brief noted as `.agents/docs/plans/2026-09-14-recon-go-packages.md`, "untracked and not committed," source/test/config untouched. Confirmed on disk: `/home/kesokaj/sandbox/runs/bug012-after1/.agents/docs/plans/2026-09-14-recon-go-packages.md` exists (4021 bytes); `git status --short` in that run dir shows only that one untracked file.
- `after2.out` tail: brief `.agents/docs/plans/2026-09-14-recon-package-map.md`, "agent's only output... Source, config and tests are unchanged." Confirmed on disk: file exists (4614 bytes) in `/home/kesokaj/sandbox/runs/bug012-after2/.agents/docs/plans/`; `git status --short` shows only that untracked file.
- `after-agy.out` tail (Antigravity 1.2.2): "Explorer brief recorded in [2026-09-14-package-map.md](.../bug012-after-agy/.agents/docs/plans/2026-09-14-package-map.md)". Confirmed on disk: file exists (8857 bytes); `git status --short` in that run dir shows only that one untracked file — source untouched.
- Antigravity log: located at `/home/kesokaj/sandbox/home/.gemini/antigravity-cli/log/cli-20260914_091037.log` (376 lines). `grep -ic "parse error"` → `0`.

## Root cause

N/A — PASS. (Root cause of BUG-012 itself is documented in `.agents/docs/plans/2026-09-14-bug-012-diagnosis.md`: the explorer subagent, under a narrow headless allow-list and/or a "report only"-style dispatch phrasing, sometimes skipped the `Write` call entirely and fabricated a plausible artifact path in its final reply instead of an actual write.)

## Coverage of changed code

Prose-only change to three markdown/context files; no Go source changed. Coverage is behavioral (does the instruction change the explorer's actual behavior), and that is exactly what the three before/after transcript pairs above test:
- `context/AGENTS.md` — covered by the general principle now appearing once in `render` output, and consistent with all three "after" transcripts obeying it.
- `agents/explorer.md` — covered directly: all three after-runs (Claude x2, Antigravity x1) wrote a real brief file and reported its path only after confirming it, versus the before2 run where the explorer was told not to write and complied silently with no brief.
- `skills/review/SKILL.md` — the clause is dispatch-instruction text for `/review`'s Step 3; not independently exercised by these three transcripts (they test the explorer's own default behavior, not a `/review`-dispatched one with the new sentence attached). This is a gap: no after-transcript shows a `/review`-style dispatch that includes the new "report only means no source edits" sentence and confirms an agent still writes its artifact under it.

## Not examined

- No test exercises the exact `skills/review/SKILL.md` Step 3 sentence in a live `/review` dispatch (only `explorer.md`'s own general behavior was retested). This clause's effect is inferred, not directly verified.
- The diagnosis file's "Open questions" (single-N sampling, whether candidate 1 is reliable vs. a one-off) were not re-tested here; this run reads existing evidence, it does not re-run the repro for statistical confidence.
- reviewer/security/documentation/verifier/bugs/changelog subagents were not retested against the new instruction — only explorer was, per the evidence provided.
- No new commit was made in this session; range reviewed includes uncommitted working-tree changes as instructed.

## Environment

Linux 6.12.107+deb13-amd64. `go version` not explicitly captured but `go test -race` and `go vet` ran successfully against the repo's `go.mod`. No Docker/Colima involved (no Dockerfile changes in this range).
