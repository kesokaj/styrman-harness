---
owner: main
date: 2026-09-15
---
# Verification — harness practices G01–G26, end to end in the sandbox on both CLIs   (2026-09-15)

Range reviewed: the harness installed into the sandbox home `/home/kesokaj/sandbox/home`:
- `22c106c+dirty` (Phases 1–5) for 6.2, 6.3 and matrix passes 1 and 2. The dirty file was the matrix spec, which is not installed.
- `c192622`, the tree before this plan, for the comparison runs (B1).
- `fcf3b79` (Phases 1–7) for the same-conditions control (C1), the 7.3 init rows (I1) and the targeted rows (P7, T7).

## Status: PASS. Every row passes on both CLIs, directly or on its one re-run, except four Antigravity rows that D50 records as caused outside this plan. Every Phase 7 row at fcf3b79 passes on both CLIs.

Scope: tasks 6.2–6.4 of `plans/2026-09-15-harness-practices-g01-g26.md`, plus the Phase 7 runs that repeat matrix rows for skills changed after passes 1 and 2:
- `init-project`, changed in 641eccd, fe21671 and c514936;
- `preflight`, changed in 641eccd, fe21671 and e31c3a2;
- `bug`, `feature`, `handoff` and the `bugs` agent, changed in 641eccd.

The spec is `.agents/docs/testing/skill-matrix.md`. The kit is `.agents/scratch/phase6-kit/`, which is gitignored:
- `rows.json` is built from the spec, with each `\|` in a check unescaped (BUG-024).
- `run-one.sh` and `run-one-v2.sh` run one row in a fresh clone under `sandbox/runs/p6-<cli>-<row>-a<run>`. Each runs the row's check once on an empty output before the CLI starts, as a control that must fail, and again after the CLI exits.
- `drive.sh` and `drive-cfg.sh` drove passes 1 and 2, `compare.sh` drove B1, `head-control.sh` drove C1 and I1, `p7-preflight-claude.sh` drove P7, and `t7-targeted.sh` drove T7.
- Results are in `results/summary.tsv`.

## 6.2 Install and drift

- **19:24:30:** `go run ./cmd/sh-agent install --home /home/kesokaj/sandbox/home` exited 0, with 1 file written and 259 unchanged. `status` exited 0: installed from `22c106c+dirty`, `drift none — 259 files match the repo`.
- **22:07:53:** an install of c192622 was refused for a missing `--repo` (BUG-030).
- **22:09:30:** c192622 installed from its plain copy, rc=0, on the plain-copy path.
- **22:26:14:** this tree installed from fcf3b79, rc=0. `status` read `drift none — 259 files match the repo` and `binary fcf3b79a8c18`, with no warning.

## 6.3 One start of each CLI

Both CLIs started from the two demo repositories and exited 0. Each quoted two passages:
- Definition of Done step 2's "`reviewer` for 50+ lines, 5+ files, or any fix under the `bug` skill";
- the compaction sentence's "restate status and next step before the next edit".

`grep -c "Failed to parse"` printed `0` on the Antigravity log that start created, `cli-20260915_192430.log`.

## 6.4 Matrix rows

How the passes ran:
- **Pass 1** ran every row once, 19:48–21:07. The main lanes ran Claude Code at `-P 3` beside Antigravity at `-P 2`. Retro, distill and handoff write the shared proposal queue, so they ran one at a time with queue snapshots, finishing at 21:04:53.
- **Pass 2** re-ran every failed row once, 21:07:24–22:06:43. A pass on the re-run counts, and the table names it.
- **"re-check 1"** is BUG-024's corrected check, run on attempt 1's own output.
- **L1:** `documentation` on Claude Code also ran a long-limit diagnostic with a 3600 s cap.

Four rows failed twice, all on Antigravity, all at its 20-minute print timeout while waiting on subagents. D50 decides them. Each ran once with c192622's harness (B1). Each ran once more with this tree under the same load (C1): three Antigravity rows at once and nothing else, 22:26:44–22:42:06.

| row | CLI | attempt 1 | control 1 | re-check 1 | attempt 2 | diagnostic | c192622 (B1) | this tree, same load (C1) | result |
|---|---|---|---|---|---|---|---|---|---|
| `feature` | Claude Code | PASS rc=0 912s | FAIL | - | - | - | - | - | PASS |
| `feature` | Antigravity | PASS rc=0 1204s | FAIL | - | - | - | - | - | PASS |
| `bug` | Claude Code | PASS rc=0 837s | FAIL | - | - | - | - | - | PASS |
| `bug` | Antigravity | FAIL rc=0 1193s | FAIL | - | FAIL rc=0 1193s | - | PASS rc=0 323s | PASS rc=0 375s | outside this plan (D50): passed at c192622 and at this tree under the same load |
| `bug-no-repro` | Claude Code | PASS rc=0 127s | FAIL | - | - | - | - | - | PASS |
| `bug-no-repro` | Antigravity | PASS rc=0 1015s | FAIL | - | - | - | - | - | PASS |
| `review` | Claude Code | PASS rc=0 51s | FAIL | - | - | - | - | - | PASS |
| `review` | Antigravity | PASS rc=0 713s | FAIL | - | - | - | - | - | PASS |
| `review-project` | Claude Code | PASS rc=0 762s | FAIL | - | - | - | - | - | PASS |
| `review-project` | Antigravity | FAIL rc=0 1202s | FAIL | - | FAIL rc=0 1200s | - | PASS rc=0 230s | PASS rc=0 209s | outside this plan (D50): passed at c192622 and at this tree under the same load |
| `preflight` | Claude Code | FAIL rc=0 921s | FAIL | corrected: PASS | PASS rc=0 969s | - | - | - | PASS |
| `preflight` | Antigravity | FAIL rc=0 1203s | FAIL | - | FAIL rc=0 1203s | - | PASS rc=0 504s | PASS rc=0 340s | outside this plan (D50): passed at c192622 and at this tree under the same load |
| `retro` | Claude Code | PASS rc=0 197s | FAIL | - | - | - | - | - | PASS |
| `retro` | Antigravity | PASS rc=0 1000s | PASS | - | - | - | - | - | PASS |
| `distill` | Claude Code | PASS rc=0 670s | FAIL | - | - | - | - | - | PASS |
| `distill` | Antigravity | PASS rc=0 1202s | FAIL | - | - | - | - | - | PASS |
| `handoff` | Claude Code | PASS rc=0 223s | FAIL | - | - | - | - | - | PASS |
| `handoff` | Antigravity | FAIL rc=0 1197s | FAIL | - | PASS rc=0 1114s | - | - | - | PASS on re-run |
| `init-deploy` | Claude Code | FAIL rc=0 255s | FAIL | corrected: PASS | PASS rc=0 151s | - | - | - | PASS |
| `init-deploy` | Antigravity | PASS rc=0 1100s | FAIL | - | - | - | - | - | PASS |
| `explorer` | Claude Code | FAIL rc=0 61s | FAIL | corrected: PASS | PASS rc=0 71s | - | - | - | PASS |
| `explorer` | Antigravity | FAIL rc=0 1202s | FAIL | - | PASS rc=0 486s | - | - | - | PASS on re-run |
| `review-standards` | Claude Code | PASS rc=0 157s | FAIL | - | - | - | - | - | PASS |
| `review-standards` | Antigravity | PASS rc=0 466s | FAIL | - | - | - | - | - | PASS |
| `golang` | Claude Code | PASS rc=0 93s | FAIL | - | - | - | - | - | PASS |
| `golang` | Antigravity | PASS rc=0 245s | FAIL | - | - | - | - | - | PASS |
| `documentation` | Claude Code | FAIL rc=124 1501s | FAIL | - | PASS rc=0 1253s | FAIL rc=0 1072s | - | - | PASS on re-run |
| `documentation` | Antigravity | FAIL rc=0 1198s | FAIL | - | FAIL rc=0 1202s | - | PASS rc=0 774s | PASS rc=0 713s | outside this plan (D50): passed at c192622 and at this tree under the same load |
| `verifier-counts` | Claude Code | PASS rc=0 68s | FAIL | - | - | - | - | - | PASS |
| `verifier-counts` | Antigravity | FAIL rc=1 290s | FAIL | - | PASS rc=0 671s | - | - | - | PASS on re-run |

### Classification (D50)

`bug`, `review-project`, `preflight` and `documentation` on Antigravity are **caused outside this plan**. Both matrix attempts ran with five CLI processes at once and hit the 20-minute print timeout. In a quieter run each passed, on both harnesses:

| row | c192622 (B1) | this tree, same load (C1) |
|---|---|---|
| `bug` | 323 s | 375 s |
| `review-project` | 230 s | 209 s |
| `preflight` | 504 s | 340 s |
| `documentation` | 774 s | 713 s |

The failure followed the load or the service at that hour, not the harness. These rows are not counted as passes.

Each C1 pass was read, not only its check:
- **`bug`:** a failing-first test, the fix, BUG-001, a per-bug reviewer report, a verifier log and a commit.
- **`review-project`:** its merged record with six lenses.
- **`preflight`:** `2026-09-15-preflight-e2cad5c.md` with `VERDICT: SHIP`, plus a docs-lens review.
- **`documentation`:** updated `README.md` and `deploy/local.md`; wrote a changes entry, a review, a verifier log, a bug record and a handoff; committed.

Attempts that failed once, with their cause:
- **`documentation` on Claude Code:**
  - a1 hit the kit's 1500 s timeout mid-Definition of Done (exit 124). It passed on its re-run (1253 s).
  - The L1 diagnostic ended its turn while a background `security` agent was still running (BUG-031).
- **`verifier-counts` on Antigravity:** a1 exited 1 on `UNAVAILABLE (code 503): No capacity available`.
- **`handoff` and `explorer` on Antigravity:** a1 hit the print timeout. Both passed on their re-runs.
- **`explorer`, `init-deploy` and `preflight` on Claude Code:** a1 failed only on the kit's escaped pipe (BUG-024). Each passes on the corrected check and on its re-run.
- **`retro` on Antigravity:** a1's control passed, because Claude Code's retro had already written the queue file. The queue snapshots attribute its pass: the file grew from 2718 to 5366 bytes with a new section from its own run.

## Phase 7 rows at fcf3b79

| row | CLI | run | control | check | rc | time |
|---|---|---|---|---|---|---|
| `init-update` | Claude Code | I1 | FAIL | PASS | 0 | 121s |
| `init-update` | Antigravity | I1 | FAIL | PASS | 0 | 107s |
| `init-deploy` | Claude Code | I1 | FAIL | PASS | 0 | 147s |
| `init-deploy` | Antigravity | I1 | FAIL | PASS | 0 | 115s |
| `preflight` | Claude Code | P7 | FAIL | PASS | 0 | 975s |
| `preflight` | Antigravity | C1 | FAIL | PASS | 0 | 340s |
| `bug` | Claude Code | T7 | FAIL | PASS | 0 | 981s |
| `bug` | Antigravity | C1 | FAIL | PASS | 0 | 375s |
| `bug-no-repro` | Claude Code | T7 | FAIL | PASS | 0 | 144s |
| `bug-no-repro` | Antigravity | T7 | FAIL | PASS | 0 | 94s |
| `feature` | Claude Code | T7 | FAIL | PASS | 0 | 494s |
| `feature` | Antigravity | T7 | FAIL | PASS | 0 | 250s |
| `handoff` | Claude Code | T7 | FAIL | PASS | 0 | 300s |
| `handoff` | Antigravity | T7 | FAIL | PASS | 0 | 90s |
| `documentation` | Antigravity | C1 | FAIL | PASS | 0 | 713s |

- **`init-update` (I1):**
  - Both CLIs made one commit, `init: update scaffold to current templates`, and deleted, modified or renamed no file.
  - Claude Code added `.agents/docs/retros/.gitkeep`. Antigravity added that file and one `.gitignore` line.
  - Both kept the existing `.claude` links.
  - c514936 changed `init-project`'s description, so the newest Antigravity log (`cli-20260915_224206.log`) was grepped: `Failed to parse` appears 0 times.
- **`init-deploy` (I1)**, from a clone reset to 8cbb128:
  - One commit, `init: scaffold agent tree and project docs`, with only added files: 25 on Claude Code, 29 on Antigravity. The four `.claude` links are the difference.
  - Claude Code refused to create `.claude/` in the headless run. The run said so and printed the command for the operator. The row's check does not cover the links.
- **`preflight` on Antigravity (C1)** is the only Antigravity run of option C's Step 4:
  - It recorded three docs gaps with severity, confidence and an owner.
  - It shipped over the high-severity gap, as ADR-024 intends.
  - It filed none of them under a handoff's Next Steps. That is in the plan's Follow-ups, as the fifth finding on BUG-026's mechanism, for the operator.
  - On Claude Code, P7 (975 s) wrote `2026-09-15-preflight-9ae4eff.md` with `VERDICT: SHIP` and a `Docs freshness (Step 4, advisory)` section ("No docs gap blocks this release"). It filed its docs gaps, with their owner, under the Next Steps of a handoff it wrote, so option C's filing clause held on Claude Code. The kit's tool allowlist refused `govulncheck`, and the run said so.
- **T7:** the rows for 641eccd's edits, at most three CLI processes at once, 22:58:46–23:25:46. **`bug`** on Claude Code (981 s): an `explorer` diagnosis, the fix, BUG-002 and BUG-003, `reviewer` and `security` reports, a verifier log, `bugs` and `changelog` dispatched, and `TODO.md` emptied into its handoff; committed. **`bug-no-repro`** (Claude Code 144 s, Antigravity 94 s): each stopped at the missing-repro gate, changed no source, and committed only a handoff whose Next Steps put the gate question first. **`feature`** (Claude Code 494 s, Antigravity 250 s): a plan in the new template (expect clauses, Checkpoint boxes, a ticked Cold-reader box citing its review, a Decisions section) and no source changes. **`handoff`** (Claude Code 300 s, Antigravity 90 s): a handoff with every template section, and a commit. Claude Code moved its demo's two `TODO.md` items into the handoff's Next Steps. The Antigravity demo's `TODO.md` had no items to move, so that path ran on Claude Code only. Neither run changed the proposal queue.

## Evidence columns and their limits

- **Claude Code:** skills and agents come from the Skill and Agent tool calls in the run's own transcript. Denials come from error tool results that name a permission.
- **Antigravity:** skills and commands come from the step payloads of every conversation database that mentions the run's directory. A database can hold more than one run's steps, so a row's skills column and its `git push --dry-run` hits can belong to another run.
- **No push reached either sandbox origin:** `demo-app` and `demo-app-agy` have no commit or reflog entry after 2026-09-14. No Claude Code run executed a forbidden command.

### Evidence per run

| row | CLI | attempt | skills loaded | agents or subagent calls | executed forbidden commands | real denials or commands run |
|---|---|---|---|---|---|---|
| `bug` | agy | 1 | bug,feature,golang,handoff,preflight,review,review-standards | invoke_subagent x15 | git push --dry-run;git push --dry-run;git push --dry-run | commands 522 |
| `bug` | agy | 2 | bug,golang,handoff,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 534 |
| `bug` | agy | B1 | bug,golang,handoff,preflight,review,review-standards | invoke_subagent x18 | git push --dry-run;git push --dry-run;git push --dry-run | commands 501 |
| `bug` | agy | C1 | bug,feature,golang,handoff,review,review-standards | invoke_subagent x12 | - | commands 255 |
| `bug` | claude | 1 | handoff | explorer,verifier,reviewer,security,bugs,changelog | - | denials 10 |
| `bug` | claude | T7 | handoff | explorer,verifier,reviewer,security,bugs,changelog | - | denials 4 |
| `bug-no-repro` | agy | 1 | bug,golang,handoff,preflight,review,review-standards | invoke_subagent x12 | git push --dry-run;git push --dry-run;git push --dry-run | commands 501 |
| `bug-no-repro` | agy | T7 | bug,handoff | none | - | commands 54 |
| `bug-no-repro` | claude | 1 | handoff | - | - | denials 0 |
| `bug-no-repro` | claude | T7 | handoff | - | - | denials 0 |
| `distill` | agy | 1 | distill,feature,golang,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 483 |
| `distill` | claude | 1 | handoff | reviewer,changelog,bugs | - | denials 1 |
| `documentation` | agy | 1 | golang,handoff,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 465 |
| `documentation` | agy | 2 | feature,golang,handoff,preflight,review,review-standards | invoke_subagent x27 | git push --dry-run;git push --dry-run;git push --dry-run | commands 651 |
| `documentation` | agy | B1 | feature,golang,handoff,review,review-standards | invoke_subagent x21 | - | commands 225 |
| `documentation` | agy | C1 | feature,handoff,review,review-standards | invoke_subagent x21 | - | commands 219 |
| `documentation` | claude | 1 | golang | verifier,reviewer,security,bugs | - | denials 11 |
| `documentation` | claude | 2 | golang,handoff,review | verifier,reviewer,security,reviewer,documentation,verifier,bugs,changelog | - | denials 15 |
| `documentation` | claude | L1 | golang,review | verifier,reviewer,security,reviewer | - | denials 18 |
| `explorer` | agy | 1 | golang,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 441 |
| `explorer` | agy | 2 | golang,preflight,review,review-standards | invoke_subagent x12 | git push --dry-run;git push --dry-run;git push --dry-run | commands 429 |
| `explorer` | claude | 1 | - | explorer | - | denials 0 |
| `explorer` | claude | 2 | - | explorer | - | denials 0 |
| `feature` | agy | 1 | feature,golang,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 510 |
| `feature` | agy | T7 | feature,review-standards | invoke_subagent x3 | - | commands 51 |
| `feature` | claude | 1 | golang,json-api-style | reviewer,reviewer | - | denials 1 |
| `feature` | claude | T7 | handoff,json-api-style | reviewer | - | denials 2 |
| `golang` | agy | 1 | golang,preflight,review,review-standards | invoke_subagent x9 | git push --dry-run;git push --dry-run;git push --dry-run | commands 492 |
| `golang` | claude | 1 | golang,review-standards | - | - | denials 1 |
| `handoff` | agy | 1 | feature,golang,handoff,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 555 |
| `handoff` | agy | 2 | bug,feature,golang,handoff,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 564 |
| `handoff` | agy | T7 | handoff | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 45 |
| `handoff` | claude | 1 | - | changelog | - | denials 0 |
| `handoff` | claude | T7 | - | changelog | - | denials 0 |
| `init-deploy` | agy | 1 | golang,init-project,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 477 |
| `init-deploy` | agy | I1 | init-project | none | - | commands 69 |
| `init-deploy` | claude | 1 | - | changelog | - | denials 2 |
| `init-deploy` | claude | 2 | - | - | - | denials 1 |
| `init-deploy` | claude | I1 | - | - | - | denials 0 |
| `init-update` | agy | I1 | init-project | none | - | commands 69 |
| `init-update` | claude | I1 | - | - | - | denials 0 |
| `preflight` | agy | 1 | golang,preflight,review,review-standards | invoke_subagent x12 | git push --dry-run;git push --dry-run;git push --dry-run | commands 660 |
| `preflight` | agy | 2 | feature,golang,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 753 |
| `preflight` | agy | B1 | golang,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 414 |
| `preflight` | agy | C1 | feature,preflight,review,review-standards | invoke_subagent x6 | git push --dry-run origin main 2\u003e&1 || true;git push --dry-run origin main 2\u003e&1 || true;git push --dry-run origin main 2\u003e&1 || true | commands 228 |
| `preflight` | claude | 1 | review | reviewer,security,verifier,reviewer,documentation,bugs | - | denials 18 |
| `preflight` | claude | 2 | - | reviewer,security,verifier,reviewer,documentation,changelog | - | denials 7 |
| `preflight` | claude | P7 | review | reviewer,security,verifier,reviewer,changelog | git -C /home/kesokaj/sandbox/runs/p6-claude-preflight-aP7 check-ignore -v -n --no-index .claude/sett | denials 7 |
| `retro` | agy | 1 | bug,distill,golang,handoff,preflight,retro,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 840 |
| `retro` | claude | 1 | - | - | - | denials 1 |
| `review` | agy | 1 | feature,golang,preflight,review,review-standards | invoke_subagent x12 | git push --dry-run;git push --dry-run;git push --dry-run | commands 468 |
| `review` | claude | 1 | - | reviewer | - | denials 0 |
| `review-project` | agy | 1 | golang,preflight,review,review-standards | invoke_subagent x12 | git push --dry-run;git push --dry-run;git push --dry-run | commands 721 |
| `review-project` | agy | 2 | feature,golang,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 675 |
| `review-project` | agy | B1 | golang,preflight,review,review-standards | invoke_subagent x12 | git push --dry-run;git push --dry-run;git push --dry-run | commands 594 |
| `review-project` | agy | C1 | codebase-design,feature,review,review-standards | invoke_subagent x3 | - | commands 237 |
| `review-project` | claude | 1 | - | reviewer,security,reviewer,reviewer,verifier,explorer | - | denials 9 |
| `review-standards` | agy | 1 | codebase-design,golang,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 414 |
| `review-standards` | claude | 1 | golang,review-standards | - | - | denials 4 |
| `verifier-counts` | agy | 1 | feature,golang,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 462 |
| `verifier-counts` | agy | 2 | golang,preflight,review,review-standards | none | git push --dry-run;git push --dry-run;git push --dry-run | commands 561 |
| `verifier-counts` | claude | 1 | - | verifier | - | denials 0 |

### Proposal queue snapshots (retro, distill, handoff)

- `agy-distill-a1`: before `2026-09-15-retro.md 5366 ffc634cbce4f`, after `2026-09-15-retro.md 5366 ffc634cbce4f`
- `agy-handoff-a1`: before `2026-09-15-retro.md 5366 ffc634cbce4f`, after `2026-09-15-retro.md 5366 ffc634cbce4f`
- `agy-handoff-a2`: before `2026-09-15-retro.md 5366 ffc634cbce4f`, after `2026-09-15-retro.md 5366 ffc634cbce4f`
- `agy-handoff-aT7`: before `2026-09-15-retro.md 5366 ffc634cbce4f`, after `2026-09-15-retro.md 5366 ffc634cbce4f`
- `agy-retro-a1`: before `2026-09-15-retro.md 2718 8e8aecec8708`, after `2026-09-15-retro.md 5366 ffc634cbce4f`
- `claude-distill-a1`: before `2026-09-15-retro.md 5366 ffc634cbce4f`, after `2026-09-15-retro.md 5366 ffc634cbce4f`
- `claude-handoff-a1`: before `2026-09-15-retro.md 5366 ffc634cbce4f`, after `2026-09-15-retro.md 5366 ffc634cbce4f`
- `claude-handoff-aT7`: before `2026-09-15-retro.md 5366 ffc634cbce4f`, after `2026-09-15-retro.md 5366 ffc634cbce4f`
- `claude-retro-a1`: before `empty`, after `2026-09-15-retro.md 2718 8e8aecec8708`

## Edit to row map

The map covers the plan's I-edits as they stood at 22c106c. An edit marked "C1 run only" passed on Antigravity only in the same-conditions control. The Phase 7 rows above exercise later edits to the same files.

| edit | rows that exercise it | status |
|---|---|---|
| I-1a | review-project (CC), preflight (CC), review-project (AG, C1 run only), preflight (AG, C1 run only) | exercised |
| I-2a | bug-no-repro (CC), bug-no-repro (AG) | exercised |
| I-3a | - | not examined: no matrix row exercises it |
| I-4a | - | not examined: harness repo file, not installed |
| I-5a | - | not examined: harness repo file, not installed |
| I-5b | - | not examined: harness repo file, not installed |
| I-5c | - | not examined: harness repo file, not installed |
| I-5d | - | not examined: harness repo file, not installed |
| I-5e | - | not examined: harness repo file, not installed |
| I-7a | feature (CC), feature (AG) | exercised |
| I-7b | feature (CC), feature (AG) | exercised |
| I-7c | feature (CC), feature (AG) | exercised |
| I-7d | - | not examined: feature row stops after the cold read |
| I-7e | - | not examined: feature row stops after the cold read |
| I-8a | - | not examined: no matrix row exercises it |
| I-8b | - | not examined: no matrix row exercises it |
| I-8c | 6.3 CLI start | exercised |
| I-8d | - | not examined: no matrix row exercises it |
| I-8e | - | not examined: no matrix row exercises it |
| I-9a | handoff (CC), handoff (AG), bug-no-repro (CC), bug-no-repro (AG) | exercised |
| I-9b | handoff (CC), handoff (AG) | exercised |
| I-9c | handoff (CC), handoff (AG) | exercised |
| I-9d | handoff (CC), handoff (AG) | exercised |
| I-10a | bug (CC), bug (AG, C1 run only) | exercised |
| I-10b | documentation (CC), documentation (AG, C1 run only) | exercised |
| I-11a | - | not examined: harness repo trail |
| I-11b | - | not examined: harness repo trail |
| I-11c | - | not examined: harness repo trail |
| I-11d | - | not examined: harness repo trail |
| I-11e | - | not examined: harness repo trail |
| I-11f | - | not examined: harness repo trail |
| I-11g | - | not examined: harness repo trail |
| I-12a | - | not examined: no row runs /review --fix |
| I-12b | - | not examined: no row runs /review --fix |
| I-12c | - | not examined: no row runs /review --fix |
| I-12d | - | not examined: no row runs /review --fix |
| I-12e | - | not examined: no matrix row exercises it |
| I-12f | - | not examined: no row runs /review --fix |
| I-13a | - | not examined: no row runs /review --fix |
| I-13b | preflight (CC), preflight (AG, C1 run only) | exercised |
| I-14a | review-standards (CC), review-standards (AG), review (CC), review (AG), review-project (CC), bug (CC), review-project (AG, C1 run only), bug (AG, C1 run only) | exercised |
| I-14b | review-standards (CC), review-standards (AG), review (CC), review (AG), review-project (CC), bug (CC), review-project (AG, C1 run only), bug (AG, C1 run only) | exercised |
| I-14c | review-standards (CC), review-standards (AG), review (CC), review (AG), review-project (CC), bug (CC), review-project (AG, C1 run only), bug (AG, C1 run only) | exercised |
| I-15a | review-standards (CC), review-standards (AG), review (CC), review (AG), bug (CC), bug (AG, C1 run only) | exercised |
| I-15b | - | not examined: no matrix row exercises it |
| I-15c | golang (CC), golang (AG) | exercised |
| I-15d | - | not examined: no matrix row exercises it |
| I-15e | - | not examined: no matrix row exercises it |
| I-15f | - | not examined: no matrix row exercises it |
| I-16a | verifier-counts (CC), verifier-counts (AG), bug (CC), documentation (CC), bug (AG, C1 run only), documentation (AG, C1 run only) | exercised |
| I-16b | verifier-counts (CC), verifier-counts (AG), bug (CC), documentation (CC), bug (AG, C1 run only), documentation (AG, C1 run only) | exercised |
| I-16c | verifier-counts (CC), verifier-counts (AG) | exercised |
| I-16d | verifier-counts (CC), verifier-counts (AG), bug (CC), documentation (CC), bug (AG, C1 run only), documentation (AG, C1 run only) | exercised |
| I-17a | review (CC), review (AG), review-project (CC), bug (CC), preflight (CC), documentation (CC), review-project (AG, C1 run only), bug (AG, C1 run only), preflight (AG, C1 run only), documentation (AG, C1 run only) | exercised |
| I-17b | verifier-counts (CC), verifier-counts (AG), bug (CC), documentation (CC), bug (AG, C1 run only), documentation (AG, C1 run only) | exercised |
| I-17c | review-project (CC), documentation (CC), review-project (AG, C1 run only), documentation (AG, C1 run only) | exercised |
| I-17d | explorer (CC), explorer (AG), bug (CC), bug (AG, C1 run only) | exercised |
| I-17e | review (CC), review (AG), review-project (CC), bug (CC), preflight (CC), documentation (CC), review-project (AG, C1 run only), bug (AG, C1 run only), preflight (AG, C1 run only), documentation (AG, C1 run only) | exercised |
| I-17f | review-project (CC), documentation (CC), review-project (AG, C1 run only), documentation (AG, C1 run only) | exercised |
| I-18a | bug (CC), documentation (CC), preflight (CC), bug (AG, C1 run only), documentation (AG, C1 run only), preflight (AG, C1 run only) | exercised |
| I-18b | bug (CC), bug (AG, C1 run only) | exercised |
| I-18c | bug (CC), bug (AG, C1 run only) | exercised |
| I-18d | bug (CC), bug (AG, C1 run only) | exercised |
| I-18e | - | not examined: no matrix row exercises it |
| I-18f | bug (CC), bug (AG, C1 run only) | exercised |
| I-18g | - | not examined: no matrix row exercises it |
| I-18h | bug (CC), bug (AG, C1 run only) | exercised |
| I-18i | retro (CC), retro (AG) | exercised |
| I-18j | retro (CC), retro (AG) | exercised |
| I-18k | - | not examined: no matrix row exercises it |
| I-18l | - | not examined: no matrix row exercises it |
| I-18m | - | not examined: no matrix row exercises it |
| I-18n | init-deploy (CC), init-deploy (AG) | exercised |
| I-19a | 5.1 precedence test | exercised |
| I-19b | retro (CC), retro (AG) | exercised |
| I-19c | retro (CC), retro (AG) | exercised |
| I-19d | distill (CC), distill (AG) | exercised |
| I-19e | - | not examined: no matrix row exercises it |
| I-20a | bug (CC), bug (AG, C1 run only), 6.3 CLI start | exercised |
| I-20b | bug (CC), bug (AG, C1 run only) | exercised |
| I-20c | bug (CC), bug (AG, C1 run only), 5.7 install | exercised |
| I-20d | - | not examined: no matrix row exercises it |
| I-20e | bug (CC), bug (AG, C1 run only) | exercised |
| I-21d | preflight (CC), preflight (AG, C1 run only) | exercised |
| I-21e | preflight (CC), preflight (AG, C1 run only) | exercised |
| I-21f | - | not examined: preflight row cannot ship (no remote, NO-SHIP expected) |
| I-21g | - | not examined: preflight row cannot ship |
| I-22a | feature (CC), feature (AG) | exercised |
| I-22b | bug-no-repro (CC), bug-no-repro (AG) | exercised |
| I-22c | bug-no-repro (CC), bug-no-repro (AG) | exercised |
| I-22e | - | not examined: no row runs /review --fix |
| I-22f | bug-no-repro (CC), bug-no-repro (AG) | exercised |
| I-23a | - | not examined: the init-deploy prompt supplies the commands |
| I-23b | - | not examined: no existing CLAUDE.md at base 8cbb128; changed again by D42 |
| I-24a | feature (CC), feature (AG) | exercised |
| I-24b | bug-no-repro (CC), bug-no-repro (AG) | exercised |
| I-24c | verifier-counts (CC), verifier-counts (AG) | exercised |
| I-24d | feature (CC), feature (AG) | exercised |
| I-25a | documentation (CC), documentation (AG, C1 run only) | exercised |
| I-25b | documentation (CC), documentation (AG, C1 run only) | exercised |
| I-25c | documentation (CC), documentation (AG, C1 run only) | exercised |

## Not examined

- **`/review --fix` end to end:** no row passes `--fix`.
- **A release that ships through a remote:** the demos have no remote, so preflight's Step 6 tag and push did not run.
- **The operator's real home:** 6.2–6.4 ran in the sandbox home only. `testing/2026-09-15-verify-real-home.md` covers the real home.
- **The Mac.**
- **Antigravity's interactive prompts:** every run was headless.
- **Claude Code's `documentation` row at fcf3b79.** Its C1 counterpart covers the `bugs` agent on Antigravity, and Claude Code's `bug` row in T7 covers it there.
- **The I-edits below**, each with its reason.

Edits not examined end to end: I-3a (no matrix row exercises it); I-4a (harness repo file, not installed); I-5a (harness repo file, not installed); I-5b (harness repo file, not installed); I-5c (harness repo file, not installed); I-5d (harness repo file, not installed); I-5e (harness repo file, not installed); I-7d (feature row stops after the cold read); I-7e (feature row stops after the cold read); I-8a (no matrix row exercises it); I-8b (no matrix row exercises it); I-8d (no matrix row exercises it); I-8e (no matrix row exercises it); I-11a (harness repo trail); I-11b (harness repo trail); I-11c (harness repo trail); I-11d (harness repo trail); I-11e (harness repo trail); I-11f (harness repo trail); I-11g (harness repo trail); I-12a (no row runs /review --fix); I-12b (no row runs /review --fix); I-12c (no row runs /review --fix); I-12d (no row runs /review --fix); I-12e (no matrix row exercises it); I-12f (no row runs /review --fix); I-13a (no row runs /review --fix); I-15b (no matrix row exercises it); I-15d (no matrix row exercises it); I-15e (no matrix row exercises it); I-15f (no matrix row exercises it); I-18e (no matrix row exercises it); I-18g (no matrix row exercises it); I-18k (no matrix row exercises it); I-18l (no matrix row exercises it); I-18m (no matrix row exercises it); I-19e (no matrix row exercises it); I-20d (no matrix row exercises it); I-21f (preflight row cannot ship (no remote, NO-SHIP expected)); I-21g (preflight row cannot ship); I-22e (no row runs /review --fix); I-23a (the init-deploy prompt supplies the commands); I-23b (no existing CLAUDE.md at base 8cbb128; changed again by D42)

## Environment

- **OS:** Linux 6.12.107+deb13-amd64 x86_64, Debian GNU/Linux 13.
- **Tools:** Claude Code 2.1.272, Antigravity CLI 1.2.3, go1.24.4 linux/amd64.
- **Sandbox home:** `/home/kesokaj/sandbox/home`.
- **Demo repositories:** `/home/kesokaj/sandbox/demo-app` for Claude Code, `demo-app-agy` for Antigravity.
- **Run clones:** under `/home/kesokaj/sandbox/runs/`.
