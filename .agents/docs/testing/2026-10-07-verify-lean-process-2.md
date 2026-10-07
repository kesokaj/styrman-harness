---
owner: verifier
date: 2026-10-07
---
# Verification — lean-process-2   (2026-10-07)

Range reviewed: f82457b..working tree (uncommitted; no Go source changed)

## Status: PASS

Scope: step 1 tier and the named installer/kit checks ran; run.sh NOT run (by dispatch); guide.md exists (modified in this change set); not complete as to behavioural runs (guide §2/§4) which the dispatch did not name

## Commands
- 1. pwd /home/ostling-simon/projects/styrman-harness; HEAD f82457b; status snapshot .agents/scratch/verify-status1.txt (18 paths before this log)
- 4. `gofmt -l cmd internal` -> exit 0, printed nothing
- 4. `go vet ./...` -> exit 0
- 5. `go test -race -timeout 30m -v -count=1 ./...` -> exit 0, all 6 packages ok (root 9.1s, cmd/sh-agent 698s, dream 20s, harness 408s, memory 12s, search 82s); counts: pass 559, fail 0, skip 3 at f82457b+dirty for go test -race -timeout 30m -v -count=1 ./...; vs previous: .agents/docs/testing/2026-10-07-verify-lean-process.md, pass 559->559, skip 3->3. Skips: TestAgyCorpusLeaves, TestRecall, TestHelperIndex (same count as before). Output .agents/scratch/verify-lean-process-2/go-test-1.txt
- 4. `go test ./internal/harness -run TestRenderedContext -v -count=1` -> exit 0; Antigravity block 17131 bytes (budget 22000); hand-written context 14722 bytes (target 15000); both PASS
- 4. `go run ./cmd/sh-agent install --home .agents/scratch/home-verify-lean2` -> exit 0 (1790 written). Checks: `.claude/agents/security.md` line 6 `effort: high` yes; `.claude/AGENTS.md` contains "last, after step 2's fixes" (1 match) yes; `.claude/skills/review/SKILL.md` contains "select:SendMessage" (1 match) yes
- 4. `bash -n .agents/docs/testing/ab-kit/run.sh` -> exit 0
- 4. `python3 -I -m py_compile` accept.py and metrics.py -> exit 0 both. Note: -I ignores PYTHONPYCACHEPREFIX, so this wrote ab-kit/__pycache__/*.pyc into the tree; I deleted those two files and the empty dir by name, then recompiled with sys.pycache_prefix=.agents/scratch/pyc (exit 0, tree clean of pycache).
- 8. Tree after: `git status` differs from step 1 only by this log; no source, test or config path moved.

## Coverage of changed code
- No Go source changed. Markdown (context/, agents/, skills/) is covered by the harness tests in `go test ./...`: lint/frontmatter rules, TestSkillListingBudget, the two size pins, agent description rules. Verified by the install checks above for the three named strings.
- NOT covered by any automated test: the wording/behaviour of the process changes in skills/{bug,feature,preflight,review}, agents/security.md and context/AGENTS.md (behavioural runs per guide §4 not run; not named by dispatch); ab-kit/run.sh (not run, token cost; only `bash -n`), accept.py and metrics.py (syntax-compiled only, no tests); prompt.md and prefix-full.md (data).
- Suppressions: no diff scan done for added nolint/skip/.only beyond the fact that no Go or TS source changed (not examined in the markdown/kit text).

## Not examined
run.sh execution and the main agent's control runs (cited in .agents/docs/changes/2026-10-07-1129-lean-process-2.md, not re-run); behavioural runs per guide §4/§5; Antigravity sandbox start; govulncheck; the content of the guide.md edit as a command change (no verification command in it was run beyond AGENTS.md's Setup commands, which are the pre-change ones as AGENTS.md is unmodified); suppression scan of the diff.

## Environment
Linux 6.18.28; Go as on PATH (`go version` not recorded); Python 3.13; no Docker used.
