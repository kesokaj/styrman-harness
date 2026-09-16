---
owner: reviewer
date: 2026-09-16
---
# Review — clear-the-board-docs   (2026-09-16)

Range reviewed: `e08fa65..HEAD` (docs lens; also read `cmd/sh-agent/main.go`, `internal/harness/{adopt.go,apply.go}` and their tests as the mechanism the docs describe)

## Verdict: REQUEST CHANGES

The prose edits themselves are careful and internally consistent — I found no place that still states the old `gofmt -l .` or the old unconditional `update` behaviour among the living docs the dispatch named. But the range leaves two gaps of its own kind: one living doc the `gofmt -l .` claim explicitly covers was never touched (`.agents/docs/testing/sh-agent/guide.md`, the standing per-app guide `agents/verifier.md` itself says must be updated "in the same run that changes the steps"), and the newly-introduced BUG-048 — cited by name in four files and a test comment — has no bug record and no row in `bugs/INDEX.md` at all, open or graduated, unlike every other BUG-NNN this range touches. `REQUIREMENTS.md` R4 also goes unmentioned by the range and is left describing `update` as unconditional.

## Findings

- **[medium / confidence: high]** `.agents/docs/testing/sh-agent/guide.md:14` — still reads `gofmt -l .   # must print nothing`. This is not a dated log; its own header says "Update this file in the same run that changes the steps," and it is exactly the file the dispatch's `gofmt` claim asks to check ("Any other place that tells someone to run `gofmt -l .` for THIS repository is a finding"). Fix: `gofmt -l cmd internal` with the same BUG-048 caveat used in `AGENTS.md`/`README.md`/`deploy/local.md`.
- **[medium / confidence: high]** `.agents/docs/bugs/INDEX.md` — BUG-048 is named and relied on in `AGENTS.md:12`, `README.md:97`, `.agents/docs/deploy/local.md:10`, `agents/verifier.md:48`, and `cmd/sh-agent/main_test.go`'s plan-tick comment, but has no `docs/bugs/BUG-048-*.md` record and no row anywhere in `INDEX.md` (checked both the main table and Graduated). Every other bug this range closes (019, 020, 021, 022, 031, 046) at least has an existing row; BUG-048 has none. Per this project's own Definition of Done ("Every failure gets a number... Dispatch `bugs`; it writes `docs/bugs/BUG-NNN-*.md`... and updates `INDEX.md`"), this is incomplete, and the plan's Phase 1 checkpoint (`.agents/docs/plans/2026-09-16-clear-the-board.md:97`) is ticked `[x]` for work that includes BUG-048 without that record existing.
- **[medium / confidence: high]** `.agents/docs/plans/2026-09-16-clear-the-board.md:97` — Phase 1 Checkpoint is ticked `[x]` and claims "`changes/` entry, commit, push, 2026-09-16," but no `.agents/docs/changes/` entry describes this range: `git log --oneline e08fa65..HEAD` shows only two commits (6008f32, 3c46356), and `grep -rl "BUG-046\|clear-the-board\|BUG-048" .agents/docs/changes/` finds only `2026-09-16-retro-2-proposals-applied.md`, an earlier same-day entry that merely records BUG-046 as newly *opened*, not this range's fixes. The tick's evidence does not match what is on disk.
- **[low / confidence: medium]** `REQUIREMENTS.md:19` (R4) — "`update` pulls the repo from git and reinstalls" is now `narrowed:` by `cmd/sh-agent/main.go`'s new dirty-checkout refusal (added at `cmd/sh-agent/main.go:251`): R4 no longer states the condition under which `update` does *not* reinstall. `REQUIREMENTS.md` is immutable without the operator's explicit approval, so the range was right not to edit it unilaterally, but the gap the dispatch asked to check for is real and worth a line to the operator (propose the exact R4 edit, per the project's own trajectory-change rule) rather than leaving it silently stale.

No other restatement of the checked claims was found stale: `README.md`, `AGENTS.md`, `context/AGENTS.md`, `agents/documentation.md`/`agents/verifier.md`, `.agents/docs/deploy/remote.md` all agree with the new `update` refusal; `<app>` resolves identically in `agents/verifier.md:37` and `agents/documentation.md:39`; `context/AGENTS.md`'s two new rules (headless-run-waits-for-agents, runtime-data-path) have no conflicting restatement in `skills/handoff`, `skills/review`, `skills/bug`, `skills/feature`, `agents/*.md`, or the init-project templates — the "headless"/"print-mode" hits in `bug`/`feature`/`review` SKILL.md describe a different mechanism (waiting on a *user* reply at a gate, pre-existing text) and don't contradict the new dispatched-agent rule; `skills/distill/SKILL.md`'s rewritten example carries exactly the nine fields `skills/feedback/SKILL.md` lists (`ts, project, session, category, signal, original, corrected, user_quote, rule_ref`), and `skills/retro/SKILL.md:65-66`'s filter syntax matches what distill's prose describes; `skills/cloud-run-deploy/SKILL.md`'s new description is consistent with `skills/preflight/SKILL.md` (which already routes its own deploy step through this skill) and doesn't collide with `skills/docker/SKILL.md`'s Do NOT line, and it still uses a valid `description: |` block scalar the installer's frontmatter linter accepts.

## Criteria

1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens
3. Standards conformance — FAIL: two of the range's own claims (the `gofmt -l .` sweep, and "every failure gets a number") are not fully carried through — see findings above. Everything else checked (the `update`, `<app>`, context-rule, and distill/feedback/retro restatements) conforms.
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined

The Go source/test changes (`cmd/sh-agent/main.go`, `main_test.go`, `internal/harness/{adopt.go,apply.go,harness_test.go}`) beyond reading them for the doc claims they support — a code-lens review of that logic is out of scope here. The full `docs/testing/*-verify-*.md` dated-log corpus was skimmed for `gofmt -l .` hits but treated as historical/excluded per the dispatch's scope, not individually verified line by line. `skills/bash/references/{makefile-patterns.md,bash-fundamentals.md,testing-linting.md}`, `skills/gcloud-setup/SKILL.md`, `skills/frontend/references/accessibility.md`, and `skills/golang/references/*` prose edits (BUG-019/020/021/022 text) were read and found internally consistent with `context/AGENTS.md` but not cross-checked against every other skill in the repo beyond the named grep targets. Did not verify the git commit trailers or push state.

## Outcome

1. [medium] `.agents/docs/testing/sh-agent/guide.md:14` still said `gofmt -l .` — **fixed**: the guide now shows `gofmt -l cmd internal` with the reason, so all five sites of that command agree (`AGENTS.md`, `README.md`, `deploy/local.md`, the guide, and the caveat in `agents/verifier.md`). This was the one site the BUG-048 sweep missed, and it was the standing per-app guide rather than a dated log, which is exactly the kind of file the sweep had to reach.
2. [medium] BUG-048 named in five places with no record and no index row — **fixed**: the `bugs` agent wrote `BUG-048` and placed it in the Graduated table in this same commit.
3. [medium] the Phase 1 checkpoint was ticked claiming a `changes/` entry that did not exist — **fixed**: `changes/2026-09-16-clear-the-board-code.md` and `changes/2026-09-16-clear-the-board-text.md` cover the range, written in this commit. The tick was accurate about the verification and the commit and wrong about the entry; the plan's checkpoint line now matches what is on disk.
4. [low] `REQUIREMENTS.md` R4 narrowed by the new refusal — **fixed with the operator's approval**, asked and given in the same turn ("Amend R4, ADR, one commit"): R4 gains "`update` refuses a checkout with uncommitted changes on every path, because it deploys the working tree rather than the commit it fetched; `install` stays the deliberate way to deploy a dirty tree", recorded as ADR-027, committed together with the requirement edit as the trajectory rule prescribes.
