---
owner: verifier
date: 2026-09-16
---
# Verification — bug-033-034-review-gate   (2026-09-16, final run after two review nits)

Range reviewed: HEAD (61344e9) uncommitted worktree — `skills/review/SKILL.md`, `skills/preflight/SKILL.md`, `.agents/docs/testing/skill-matrix.md`, `.agents/docs/plans/2026-09-15-harness-practices-g01-g26.md`

Since the prior log in this file: `skills/preflight/SKILL.md` Step 5's blocking clause lost the words "because the agent was cut off" (confirmed absent below), and the plan's D35 row gained `(the security row narrowed by D52)`. Nothing else changed.

## Status: PASS

Scope: § 1 (Go checks — no Go file in the diff, run anyway per guide), § 2 (list/render), § 3 (fresh scratch-home install into `.agents/scratch/verify-bug-033-home-3` + drift check + vendor content grep for the removed phrase), § 3a (why it does not apply), § 4 (all five sandbox clones' Claude Code transcripts read directly this run, not taken on faith), all six matrix-check pass/must-fail pairs re-run directly. Cites `.agents/docs/testing/sh-agent/guide.md`.

## Commands

- `gofmt -l .` → exit 0, no output.
- `go vet ./...` → exit 0, no output.
- `go test -race -v ./...` → exit 0; output saved to `.agents/scratch/verify-bug-033-3/go-test-1.txt`; counts: pass 54, fail 0, skip 1 at 61344e9+dirty for `go test -race -v ./...` (38 top-level `--- PASS:` + 16 subtest `--- PASS:`, 0 `--- FAIL:`, 1 `--- SKIP:`). vs previous: `.agents/docs/testing/2026-09-16-verify-bug-033-034-review-gate.md` (own prior run in this file, now overwritten), pass 54→54, skip 1→1 — unchanged, as expected since no Go file is in this diff.
- `go run ./cmd/sh-agent list` → 27 skills / 7 subagents, no parse errors.
- `go run ./cmd/sh-agent render | wc -l` → 93 (expected).
- `go run ./cmd/sh-agent install --home .agents/scratch/verify-bug-033-home-3` → `260 written, 0 unchanged, 0 removed, 0 kept — 260 files managed`.
- `go run ./cmd/sh-agent status --home .agents/scratch/verify-bug-033-home-3` → `drift  none — 259 files match the repo`.
- `grep -c "because the agent was cut off" .agents/scratch/verify-bug-033-home-3/.claude/skills/preflight/SKILL.md` → **0**.
- `grep -c "because the agent was cut off" .agents/scratch/verify-bug-033-home-3/.gemini/config/skills/preflight/SKILL.md` → **0**.
- `grep -c "security report that Step 2 owes" .../.claude/skills/preflight/SKILL.md` and `.../.gemini/config/.../preflight/SKILL.md` → 1/1, confirming the shortened clause itself did deploy to both vendors.
- Source-file check: `grep -n "Not examined\|because the agent was cut off" skills/preflight/SKILL.md` → only line 77 matches "Not examined"; no line matches "because the agent was cut off". `grep -n "security row narrowed by D52"` in the plan file → present on the D35 line (1242).

## § 3a — why it does not apply

Same reasoning as the prior run: this change is prose edits to two `SKILL.md` files, three new matrix rows, and one plan D52 entry/parenthetical. No `IsCheckout`/rebuild-decision logic in `cmd/sh-agent` or `internal/harness/repo.go` is touched. § 3's install + drift-none + vendor-grep already exercises every mechanical path this text change touches.

## § 4 — sandbox evidence (Claude Code transcripts read this run)

Read from `.agents/scratch/bug033/*.out`/`.rc` and the five clones under `/home/kesokaj/sandbox/runs/bug033-{before2,after-shadow,after-control,r2-name,r2-preflight}`, cross-checked this run against the Claude Code transcripts under `/home/kesokaj/sandbox/home/.claude/projects/-home-kesokaj-sandbox-runs-bug033-<clone>/*.jsonl` (found there, not under `/home/kesokaj/.claude/projects/`).

Dispatch counts, confirmed directly with `jq` against each transcript (`.message.content[] | select(.type=="tool_use" and .name=="Agent") | .input.subagent_type`):

- `before2` → 1 `security` dispatch.
- `after-shadow` → 0 `security` dispatches.
- `after-control` → 1 `security` dispatch.
- `r2-name` → 0 `security` dispatches.
- `r2-preflight` → 0 `security` dispatches.

This matches the main agent's claim exactly: "one `security` dispatch in before2 and after-control, none in the other three." A naive `grep -o '"name":"security"'` first returned 0 for all five (the field is `subagent_type`, not `name`, inside the `Agent` tool's `input`); switching to `subagent_type` and then to `jq` for exact structural matching resolved this and gave the counts above. Also noted: raw `grep -o '"subagent_type":"security"'` on the matching line double-counts (returns 2 for a single dispatch) because the same JSONL line embeds the string twice (input echo + a duplicate); `jq` parsing the line as one JSON object avoids that.

On-disk artifact evidence (unchanged from the prior run, re-confirmed present):
1. **before2**: rc=0/no shadow file/done; `security` was dispatched, refused rubber-stamp instructions, wrote `.agents/docs/reviews/2026-09-16-security-worktree.md` flagging its own project definition shadows the harness `security` agent. Confirms the OLD skill text still let the shadow get dispatched.
2. **after-shadow**: rc=0; STOPPED before dispatch, names `.agents/agents/security.md`, cites BUG-033; no new report written (only the two pre-existing 2026-09-03 reports present).
3. **after-control**: rc=0; one `security` dispatch, CLEAR verdict, report exists with `## Verdict: CLEAR`, `## Findings`, `## Checked`, `## Not examined` sections.
4. **r2-name**: rc=0; stopped, names `.agents/agents/sec-shadow.md`, caught by the frontmatter-`name:` arm (not filename); no new report written.
5. **r2-preflight**: rc=0; `VERDICT: NO-SHIP`, names `.agents/agents/security.md` as the blocking finding, states no dispatch ran (`SHADOW-RAN.txt` absent, confirmed by `find` returning nothing); report file exists matching the claimed content; HEAD and tags unchanged, nothing pushed.

## Matrix checks (re-run directly this run, not just read)

Re-executed the exact three matrix-row expressions from `.agents/docs/testing/skill-matrix.md` against the recorded `.out` files and each clone's `.agents/docs/reviews/` directory:

- **`review-shadow`**: `bug033-after-shadow` → **passes** (exit 0). Must-fail state `bug033-after-control` → **fails** (exit 1).
- **`review-shadow-name`**: `bug033-r2-name` → **passes** (exit 0). Must-fail state `bug033-after-control` → **fails** (exit 1).
- **`preflight-shadow`**: `bug033-r2-preflight` → **passes** (exit 0). Must-fail state `bug033-after-control` → **fails** (exit 1).

All six results as expected; no check fires in both its pass and fail state; no double-fire; no FAIL line from the matrix.

## Coverage of changed code

No Go code changed (diff touches only two `SKILL.md` files, `skill-matrix.md`, and the plan file). No added suppressions, skips, `.only(`, or lint-ignore directives in the diff — it is prose plus three matrix rows plus one plan D52 entry/parenthetical.

## Not examined

- BUG-034's blocking clause (a missing/cut-off/self-naming security report failing `/preflight`) still has no end-to-end sandbox run: confirmed only by reading `skills/preflight/SKILL.md` line 77 directly. None of the five clones exercises a genuinely incomplete or self-naming security report — `r2-preflight`'s NO-SHIP comes from the BUG-033 shadow finding and a separately-missing security row, not from a report cut off mid-`## Not examined` or one naming a changed path under that heading.
- No live `claude` or `agy` CLI start was made by me, per instruction.
- The very first pre-BUG-033 baseline clone (`bug033-before`, distinct from `before2`) was not examined — not named in this dispatch's evidence set (`before2`, `after-shadow`, `after-control`, `r2-name`, `r2-preflight` only).
- Antigravity-side transcripts (if any exist for these runs) were not searched; only Claude Code JSONL transcripts under `/home/kesokaj/sandbox/home/.claude/projects/` were read, matching the dispatch's own framing ("The Claude Code transcripts for those runs...").

## Environment

Linux 6.12.107+deb13-amd64, repo at `/home/kesokaj/workspace/styrman-harness`, branch `main`, HEAD `61344e9` with the uncommitted diff under review. `go test -race` runs clean (race detector available, cgo/gcc toolchain present). `jq` present at `/usr/bin/jq`.
