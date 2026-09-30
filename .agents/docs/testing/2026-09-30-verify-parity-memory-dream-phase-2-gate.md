---
owner: verifier
date: 2026-09-30
---
# Verification — parity-memory-dream-phase-2-gate (2026-09-30)

Range reviewed: ae8f965..f554da8 plus uncommitted tree (plan, testing/guide.md, main's log 2026-09-30-verify-parity-memory-dream-phase-2.md). Guide: .agents/docs/testing/guide.md (exists).

## Status: PASS (re-run after the fixes; first run was FAIL)
Scope: Go checks, harness self-check, scratch-home installs (claude, agy, all) and evidence audit ran; no CLI started. FAIL is for one stale living-text claim and for evidence that does not support one sentence of main's log (below).

## Commands
- `gofmt -l cmd internal` → exit 0, prints nothing
- `go vet ./...` → exit 0
- `go test -race -count=1 -v ./...` → exit 0; ok for root, cmd/sh-agent, internal/harness, internal/memory, internal/search; counts: pass 406, fail 0, skip 2 at f554da8+dirty (dirty = docs only) for go test -race -count=1 -v ./...; vs previous: .agents/docs/testing/2026-09-30-verify-v1.16.2-gate.md, pass 375→406, skip 2→2 (skips: TestRecall needs SH_AGENT_RECALL=1; TestHelperIndex helper process)
  - TestRenderedContextFitsAntigravityBlock PASS: 21911 bytes (context 21911, rules 0), budget 22000 (89 bytes headroom)
  - TestInitProjectVendorCheckMatchesManifest PASS; TestEmbedHoldsOnlyTrackedFiles PASS
- Scratch-home installs under .agents/scratch/v/home-<vendor>, each `install --vendor V --home H` exit 0, `status` → "drift none":
  - claude: 134 files; .claude/skills/project-memory/SKILL.md present
  - agy: 133 files; .gemini/config/skills/project-memory/SKILL.md present
  - all: 267 files; both paths present
  (main's log says 135/134/268 files: one more each than I count; possibly a manifest/count difference at install time, not chased.)

## Root cause  (FAIL)
1. Stale living text: skills/dream/SKILL.md:50 says the index stays under 200 lines and 25 KB "because that is all Claude Code loads at start". After R10 Claude Code's native memory is off and `sh-agent memory show` does the loading (spec memory-loop.md:14). The claim that Claude Code loads the memory natively survives there.
2. Wording remnants (lower): context/AGENTS.md:42 "`/dream` tidies the project's auto memory"; skills/dream/SKILL.md:16 "consolidate the auto memory", :26 "no auto memory for <repo>"; skills/retro/SKILL.md:118 "no auto memory store". Not claims of native writing, but the old term.
3. Main's log, "Every stderr has 0 `auto-denied` lines" is not evidence: each ~/sandbox/p2/runs/*.err and dreamer/r*.err holds one line only (`rc=0 secs=N`), so CLI stderr was never captured there. The rc and secs numbers match the log exactly (0 everywhere; 14/12, 24/23, 9/11, 29/9 s etc.).
No remaining claim that /dream resolves autoMemoryDirectory: `git grep autoMemoryDirectory` in context/skills/agents/commands hits only skills/init-project/SKILL.md:115 (where /init-project writes it).

## Evidence audit (files under ~/sandbox/p2/runs)
Supported:
- claude-m1 deploy-on-fridays.md; agy-m1 friday-deploys.md; all-m1 deploy-fridays.md; all-m2 deploy-on-fridays.md: names and MEMORY.md index lines match the table; metadata lacks originSessionId in all four. Recall .out: agy answers "Friday" for agy-m1 and all-m1; cc/all-m2 session ids present; .out for remember of agy names the file.
- claude-n1 (control) indent-with-tabs.md carries originSessionId, title-case line "Indent with tabs", different file modes (664 vs 600): native writer. claude-n2 tabs-over-spaces.md has no originSessionId, slug-shape line, verb-style frontmatter.
- dreamer-r1..r5: retros/2026-09-30-dreamer.md each has Status COMPLETE, Scope line naming <clone>/.agents/memory with both promoted files, a skill candidate, open-migration.md left alone; working tree shows only untracked .agents/docs/retros/. Bug ids: r1, r2, r4, r5 cite BUG-001/002; r3 cites only "bugs/INDEX.md:5-6" and "two bug rows" (no ids), as main says.
Minor inconsistencies: Scope line counts "4 memory files" in r1/r4 but "3" in r2/r3/r5 although each store holds 3 memories plus MEMORY.md (r1, r2 checked); the dreamer counting differs run to run, unmentioned in main's log. Timestamps: 10 runs span 20:02:58-20:03:36, so the "clone m1/m2" runs ran interleaved; fine.
Not verifiable from disk: transcript claims (projects/*.jsonl, agy conversation stores), the real ~/.config/sh-agent hash, credentials deletion, 0 toPermission grep, the agents.json from working-tree dreamer.md.

## Coverage of changed code
No .go file changed in f554da8 beyond ae8f965's work (Phase 2 is text); Go tests above cover skill manifest, embed, size budget. No added //nolint, t.Skip, .only( found in the change set by inspection of the Go results (skip count unchanged at 2).

## Not examined
Transcript stores (sandbox projects/*.jsonl, agy conversations), hashes of the real home, the file-count difference 134 vs 135, `git diff` of the uncommitted plan/guide paragraph (only confirmed present, not read line by line), /review --docs items, other specs outside memory-loop.md.

## Environment
Linux 6.12 (Debian 13), Go toolchain as installed, no Docker needed. Tree at end matched tree at start (docs-only dirty paths; this log is new).

## Re-run after the fixes

Tree: HEAD f554da8 plus uncommitted fixes (context/AGENTS.md, skills/dream, skills/retro, plan, guide, the two logs). Source/tests/config did not change during this run (porcelain status identical at start and end). No CLI started; real home untouched.

Status: PASS. Scope: the three requested checks ran; not a full re-run of the scratch-home installs (see Not examined).

### (1) `git grep -n -i "auto memory\|Claude Code loads"` over context/ skills/ agents/ README.md .agents/docs/specs/
11 hits, all acceptable; none describes Claude Code writing or loading the store:
- specs/sh-agent/memory-loop.md:12 - "Claude Code's auto memory layout, unchanged": file layout. OK.
- specs/sh-agent/memory-loop.md:13 - native auto memory switched off per project via `autoMemoryEnabled: false`. OK.
- README.md:187 - native auto memory switched off by /init-project; every CLI writes via `sh-agent memory`. OK.
- README.md:216 - same, native switched off. OK.
- skills/feedback-log/SKILL.md:90 - native switched off, the verb is the one writer. OK.
- skills/init-project/SKILL.md:110 - item 2 switches native off; item 3 moves memory at Claude Code's old default path. OK.
- skills/init-project/SKILL.md:115 - the two settings keys; native stays on only where the workspace is untrusted, which `memory show` reports. Describes native memory being switched off. OK.
- skills/project-memory/SKILL.md:8 - native switched off by /init-project. OK.
- skills/preflight/SKILL.md:48 and skills/review/SKILL.md:42 - "Claude Code loads the project copy" of a shadowing review AGENT; unrelated to memory. OK.
- context/ and agents/: no hits (the stale claim from the first run is gone).

### (2) Main's log 2026-09-30-verify-parity-memory-dream-phase-2.md against the sandbox
- Stderr: ~/sandbox/p2/runs/*.err (10 files) and ~/sandbox/p2/dreamer/r1..r5.err each hold exactly one line, `rc=0 secs=N`; `grep -il auto-denied` over p2 and gate17 .err files finds nothing. Claim that the capture holds only the kit's rc line: CONFIRMED. The dreamer r*.err (34-36 s) are also rc lines only. Note main's log says "0 auto-denied lines" as a measured result; it is true, and the kit's capture form (`2> run.err`) is the same in gate17's ~/sandbox/gate17/sbrun.sh:13.
- gate17 claim ("the same capture recorded agy's auto-denied line"): the CURRENT gate17/runs/*.err are also rc-only (init-agy-only.err and dream-agy-only.err: 14 bytes). History exists in ~/sandbox/gate17/sb.log and sb2.log: `init-agy-only.err auto-denied: 1` and `dream-agy-only.err auto-denied: 1` (sbrun.sh:21 greps those files), and sb3.log shows 0 after the rerun overwrote them. So the claim is supported by the logs of earlier attempts, not by the files as they stand now. The actual stderr text cannot be checked now; only the counts.
- Scope-count variance: the retros' Scope lines say r1 4 files, r2 3, r3 3, r4 4, r5 3. Main's log says "4 in r1 and r4, 3 in r2, r3, r5": CONFIRMED. Stores of all five are identical: `MEMORY.md`, `open-migration.md`, `vet-before-tag-1.md`, `vet-before-tag-2.md` (3 memories plus the index; r3, r4, r5 listed explicitly). So 4 = index counted as a memory. CONFIRMED.
- Not re-checked here: the r3 id-needle failure and the matrix table (the artifacts were not re-read), nor E2E 1/2 tables.

### (3) Go checks (bare, separate calls)
- `go vet ./...` -> exit 0
- `go test -race -count=1 -v ./...` (output .agents/scratch/verify-p2b/go-test-1.txt) -> exit 0; ok for root, cmd/sh-agent, internal/harness, internal/memory, internal/search; counts: pass 406, fail 0, skip 2 at f554da8+dirty for go test -race -count=1 -v ./...; vs previous: this log's first run, pass 406->406, skip 2->2
- TestRenderedContextFitsAntigravityBlock PASS: 21906 bytes (context 21906, rules 0), budget 22000 (94 bytes headroom); first run saw 21911, the 5 bytes are the uncommitted context/AGENTS.md fix.

### Coverage of changed code (delta)
Changes since the first run are text only (context/AGENTS.md, skills/dream, skills/retro); covered by the size test, the harness/embed tests and install checks. No suppressions or skips added.

### Not examined (re-run)
- Scratch-home installs (claude/agy/all) and `gofmt -l` were not repeated; the embed/render tests that cover install content pass.
- The E2E and dreamer tables in main's log beyond the two claims above.
