---
owner: reviewer
date: 2026-10-01
---
# Review — worktree-p5-code   (2026-10-01)

Range reviewed: working tree — changed files: context/AGENTS.md, agents/documentation.md, skills/VENDORED.md, skills/bug/SKILL.md, skills/distill/SKILL.md, skills/dream/SKILL.md, skills/feature/SKILL.md, skills/feedback-log/SKILL.md, skills/handoff/SKILL.md, skills/init-project/SKILL.md, skills/preflight/SKILL.md, skills/research/SKILL.md, skills/retro/SKILL.md, skills/review-standards/SKILL.md, skills/review/SKILL.md, skills/review/scripts/restatement-grep.sh, .agents/docs/bugs/INDEX.md, .agents/docs/plans/2026-09-30-audit-fix.md, .agents/docs/testing/skill-matrix.md; untracked, not ignored: internal/harness/restatement_script_test.go, .agents/docs/bugs/BUG-138-restatement-grep-silent-outside-harness.md, .agents/docs/testing/2026-10-01-verify-phase5-skill-matrix.md (diff read with `git diff HEAD -- <files>`; untracked files read whole)

## Status: COMPLETE

## Verdict: REQUEST CHANGES (round 1; round 2 below: APPROVE WITH NITS)

One test gap, no defect in the script itself. `TestRestatementGrepRunsOutsideHarness` covers only the hit arrangement (two fixtures, both expect exit 0 plus a hit line). The no-hit arrangement (exit 1, empty stdout) is not asserted, and the `--all` branch, which the same diff changed (`grep -nH`), is not run at all. BUG-138's own Prevention rule 1 and the dispatch both name "hit and no-hit, exit code and output shape". I ran the script by hand in a scratch non-harness repo (no-hit exit 1 with no output; `--all` hit exit 0 with `a.md:1:` prefix; `--all` miss exit 1; multi-word `-i` exit 0) and it behaves correctly, so the fix is right; the pin on it is incomplete. `go test ./internal/harness -count=1` passes (includes the size pins and the new test). The markdown procedures I traced resolve as written (see Criteria 2).

## Findings
- **[medium / confidence: high]** `internal/harness/restatement_script_test.go:60-84,86-112` — Only hit arrangements exist. Missing: (a) a no-hit run in the non-harness fixture asserting exit 1 and empty stdout (the exact "exit 1, silent" shape BUG-138 says is indistinguishable from failure; without the pair, a regression that makes every run exit 1 is caught only by the hit tests, and one that exits 0 or prints noise on no-hit is caught by nothing); (b) any `--all` run, though line 57 of `skills/review/scripts/restatement-grep.sh` (`grep -nH` feeding the `${line#*:*:}` strip) changed in this diff and is a different code path from the default branch; before `-H`, a single-file `--all` would have stripped the wrong prefix. Fix: add `run(dir, "--", "zzznomatch")` expecting `code == 1 && out == ""`, and `run(dir, "--all", "--", "name", "greets")` (hit, `README.md:1:`) plus one `--all` miss expecting exit 1, in the existing single-file fixture. Same arrangement the finding measured (project without VENDORED.md, one tracked file).
- **[low / confidence: medium]** `internal/harness/restatement_script_test.go:46` — `run` returns `stdout + stderr` joined, so the "output shape" assertion `strings.Contains(out, "README.md:1:")` cannot tell a stdout hit from an stderr message. The documented contract is `path:line: text` on stdout. Fix: return the two streams separately and assert on stdout; stderr empty on hit and no-hit.

## Criteria
1. Dead code — pass. Removed text in skills is replaced, not orphaned; `argument-hint --global` is defined (distill Step 3 table), `--from-log` is now implemented (distill:40). No unreferenced exports in the Go test.
2. Flow tracing — one finding (test gap above). Traced and fine: `/review --resume` -> `explorer` (tools include Bash, Write; writes `plans/<date>-recon-<slug>.md`, matching `plans/<date>-recon-project-resume.md`); handoff cold-reader names `/review --resume` and review's resume brief points at "handoff skill § Cold-reader check", which exists (handoff:156) and holds the awk block; preflight requires `## Status: COMPLETE` and both `reviewer` and `security` templates emit that line; `documentation` has no Bash (tools Read, Grep, Glob, Edit, Write) so the "carry `sh-agent search --refs` output" rule is consistent; feedback-log rule 1 cites "rule 7's memory refusal" and rule 7 does contain it; distill's "`retros.jsonl` line above" exists (Step 6.2); the review skill dir lookup order (`~/.claude`, `~/.gemini/config`, `.agents`) matches plan.go's install targets; init-project sub-steps 1-6 are in dependency order (root, trust and store, settings trust check, keys, migration, headless guard) and the Antigravity-only skip list (3-5) is consistent with sub-step 3's "skip 5". preflight Step 1 numbered list is in dependency order. Skill-matrix rows use the existing 5-column shape and `\|` escape. `restatement-grep.sh` verified by hand as above.
3. Standards conformance — pass. § 8: no guard-before-destructive or read-only-mode issue in the script; shell is `set -euo pipefail` with the empty-pipeline `|| true` now on the whole group (BUG-138's rule). Describing text: script header updated for `-H`; BUG-138 and INDEX row agree with the code; plan line matches the implemented `--resume` slug/path (`project-resume`). Context size pins pass. No stale restatement found for the lens list (`--all` "every lens but `--resume`" appears only in review/SKILL.md; README/specs have no copy).
4. Directory layout — pass. Test beside the other harness tests (`internal/harness/`), bug record under `bugs/`, verify log under `testing/`.
5. Simplicity (DIY-except-security) — pass. Stdlib only in the test; nothing auth/crypto related.
6. Decoupling — pass. The test reaches the script by relative path from `internal/harness` into `skills/`, the same coupling the existing embedded-plan test has; no cross-package reach-ins.

## Not examined
- Behavioural runs on either CLI (the verify log's 35-run matrix): I read its summary lines only, did not re-run any skill or agent.
- The must-fail mutation protocol for the new test (the BUG-138 record and verify log state it failed before each of the two script fixes; I did not reproduce the pre-fix failure, only confirmed the fixed script passes and behaves correctly by hand).
- Antigravity-specific behaviour of the changed skills (no run).
- `skills/dream`, `skills/research`, `skills/retro`, `skills/bug`, `skills/feature` one-line edits were read in the diff; their surrounding unchanged text was not re-read beyond the lines cited.

## Round 2

Range reviewed: working tree — changed files: context/AGENTS.md, agents/documentation.md, skills/VENDORED.md, skills/bug/SKILL.md, skills/distill/SKILL.md, skills/dream/SKILL.md, skills/feature/SKILL.md, skills/feedback-log/SKILL.md, skills/handoff/SKILL.md, skills/init-project/SKILL.md, skills/preflight/SKILL.md, skills/research/SKILL.md, skills/retro/SKILL.md, skills/review-standards/SKILL.md, skills/review/SKILL.md, skills/review/scripts/restatement-grep.sh, README.md, .agents/docs/bugs/INDEX.md, .agents/docs/plans/2026-09-30-audit-fix.md, .agents/docs/testing/skill-matrix.md; untracked, not ignored: internal/harness/restatement_script_test.go, .agents/docs/bugs/BUG-138-restatement-grep-silent-outside-harness.md, .agents/docs/testing/2026-10-01-verify-phase5-skill-matrix.md, .agents/docs/testing/2026-10-01-verify-audit-fix-phase-5.md (diff via `git diff HEAD -- <files>`; untracked read whole)

### Round-2 verdict: APPROVE WITH NITS

Both round-1 findings are fixed and verified. I ran the script by hand in a scratch repo (`.agents/scratch/r2/fx`, one tracked README.md with line 1 `GET /hello?name=<who> greets name.` and line 2 `second line, no match`) once per arrangement: `-i -- name` exit 0, stdout exactly `README.md:1:...`, stderr empty; `-- zzznomatch` exit 1, stdout and stderr empty; `--all -- hello greets` exit 0, the one prefixed line; `--all -- hello second` exit 1, empty. The `--all` miss does put the two words on different lines (hello line 1, second line 2), so it discriminates against an any-word implementation. `go test ./internal/harness -run 'TestRestatement|TestRenderedContext' -count=1` passes. One new low finding in a fix.

### Round-2 findings
- **[low / confidence: medium]** `agents/documentation.md:37` — the fix for "no shell" replaces `ls -d cmd/*/` with "a Glob for `cmd/*/` matching exactly one". Glob matches files, not directories; the pattern `cmd/*/` most likely returns nothing, so the agent falls through to the second and third rule and never finds `<app>` from `cmd/`. Fix: name a file-shaped pattern, e.g. `cmd/*/main.go` (one match), which is what the old `ls -d` meant. Not verified against the tool (no Glob call available to me).

### Round-1 ledger
- Test gap (medium) — fixed: `restatement_script_test.go:53-80` has the four arrangements, each proven by hand above; stdout and stderr are separate (`:37-51`), hit asserted on stdout exactly and stderr empty on hit and miss. The low finding (joined streams) — fixed.
- Harness-shaped test (`:82-107`) still reads `out` only; fine, it asserts only stdout.

### Fixes re-read, each declaration and its references
- `skills/review/SKILL.md:27` `--resume runs alone`: Step 1 is the only place the flag set is declared; `--all` is defined as "every lens but `--resume`"; handoff:158 and Step 4's one-lens case agree (one artifact = the record). No restatement of the old combination left (grepped README, specs, handoff).
- `skills/review/SKILL.md:83` three-path lookup: matches the install targets (`~/.claude/skills`, `~/.gemini/config/skills`, project `.agents/skills`); the awk program is under handoff § Cold-reader check as cited.
- `skills/review/SKILL.md:112` headless `--fix` sentence reads consistently with the cap paragraph and with bug:15 / feature:15, which now point at "the `review` skill Step 5".
- Size-pin command: `-run TestRenderedContext` in review:105, retro:143 and README:275 all match both pin names; the living plan line (audit-fix.md:176) names only `UnderTarget`, a deliberate single check, not a restatement. Remaining `...FitsAntigravityBlock` hits are dated records (testing/, reviews/, bugs) which are not living files.
- `skills/init-project/SKILL.md:112,118`: root via `--git-common-dir` equals `internal/memory/memory.go:62` (`Root`) and dream/SKILL.md:22; `top` is used only by sub-step 5's slug, and the Antigravity-only skip list (3 to 5) still holds. No orphaned reference to the old `--show-toplevel` root remains in this step.

### Criteria (round 2)
1. Dead code — pass; the fix introduced no unreferenced declaration (`run`, `git` helpers both used).
2. Flow tracing — pass apart from the Glob finding; the lookup, root and cap pointers resolve.
3. Standards conformance — pass; every restatement of the size-pin command and of the review cap now agrees.
4. Directory layout — pass.
5. Simplicity — pass; stdlib only.
6. Decoupling — pass.

### Round-2 not examined
- The `2026-10-01-verify-audit-fix-phase-5-r2.md` verify log (still IN PROGRESS while I ran; not read).
- Behavioural runs of the changed skills on either CLI.
- The Glob tool's behaviour on `cmd/*/` (reasoned, not run).

## Outcome
- round 1 [medium] test has no no-hit or `--all` arrangement — fixed at `internal/harness/restatement_script_test.go:47-74` (hit, miss, `--all` hit, `--all` miss; must-fail: `-H` dropped on the `--all` branch, `TestRestatementGrepRunsOutsideHarness` failed on the `--all hit` line, restored, `git diff --numstat` back to 8/4)
- round 1 [low] `run` merges stdout and stderr — fixed at `internal/harness/restatement_script_test.go:21-37` (separate streams; hits asserted on stdout, stderr asserted empty)
- round 2 [low] Glob for `cmd/*/` matches no directory — fixed at `agents/documentation.md:37` (`cmd/*/main.go`, its directory is `<app>`); behavioural run p5d3 after the edit: dispatched, Glob wording loaded
