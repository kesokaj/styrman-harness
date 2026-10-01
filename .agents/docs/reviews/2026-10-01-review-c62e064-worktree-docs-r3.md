---
owner: reviewer
date: 2026-10-01
---
# Review — c62e064-worktree-docs-r3   (2026-10-01)

Range reviewed: the working tree against commit c62e064 (`git diff c62e064` plus the untracked files of `git status --short --untracked-files=all`: handoff, round-1 and round-2 reviews, `testing/kit-preamble.sh`, `skills/review/scripts/restatement-grep.sh`). Docs lens, round 3.

## Status: COMPLETE

## Verdict: APPROVE

No defects in the files as they stand. The three round-2 answers hold; the final pass over `AGENTS.md` § Testing and guide § 4 § Behavioural runs found nothing dropped or narrowed.

## Round-2 answers, checked
1. Graduated rows: `bugs/INDEX.md:134,136,137,138,140,141` (BUG-083/087/090/091/106/107) and the new BUG-123 row (:147) point at `guide.md § 4 § Behavioural runs` and/or `testing/kit-preamble.sh`; no row names the deleted `AGENTS.md` paragraph. Holds.
2. `guide.md` § Behavioural runs says the preamble "FAILS the kit (exit 1) when it changed"; `kit_on_exit` does exit 1 with a FAIL line. Holds.
3. The restatement-grep sentence sits in the first quoted `--docs` brief (`skills/review/SKILL.md:51`, before `Under --project:`). `kit-preamble.sh` builds with `go build -C "$kit_repo"`, `kit_repo` derived from `BASH_SOURCE`, so the cwd no longer matters; the header's `KIT_BIN` line says so. Holds.

## Final pass: dropped / narrowed (AGENTS.md § Testing, old vs new)
- search thresholds (rebuild under 10 s, warm search under 100 ms): guide.md:120-121. Kept.
- install + status drift, one start per CLI, log grepped: guide § 2, § 3, § 4. Kept.
- behavioural-run conditions (stdin prompt, `--allowedTools`, subagent transcript, n of N, BUG-083, `--append-system-prompt`), real home untouched, scratch-HOME cannot authenticate, what a log did not exercise, skill-matrix: guide.md § Behavioural runs. Kept.
- real-home smoke incl. GANNET/BUG-085, each CLI: guide.md:175-176. Kept.
- agy settings bullet, must-fail control, `--setting-sources` bullet (D15), clone outside every checkout (BUG-106), `--home` wrapper and before/after listing (BUG-107): guide bullets plus the preamble. Kept; the dated "nine /dream probes" detail is a record, not a rule.
- `/review --docs` living-doc set (`context/`, `agents/`, non-vendored skills, searched at `skills/`): AGENTS.md:68 last sentence. Kept.
- Corpus measurement (971 / 0, 2026-09-30): AGENTS.md:68 and guide.md:135 agree.
- `tier-4`: no living hit remains (restatement-grep). Unticked boxes in `plans/2026-09-30-audit-fix.md:136,148` still say "AGENTS.md § Testing"; plans are outside the living set and that heading still points on to the guide, so not a finding.

## Findings
none

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
The preamble and the script were read, not re-run this round (round 2 ran both by hand; round 3 changed only the `-C` build, which I read). The other modified plans, specs, handoff, and the Antigravity side of the preamble were not re-read beyond rounds 1 and 2. No agy run.

## Outcome

No findings. Phase 2 of plan audit-fix closes on this round.
