---
owner: reviewer
date: 2026-10-01
---
# Review — c62e064-worktree-docs-r2   (2026-10-01)

Range reviewed: the working tree against commit c62e064 (`git diff c62e064` plus the untracked files of `git status --short --untracked-files=all`: handoff, round-1 review, `testing/kit-preamble.sh`, `skills/review/scripts/restatement-grep.sh`). Docs lens, round 2.

## Status: COMPLETE

## Verdict: REQUEST CHANGES

All six round-1 answers hold in the files as they stand (probed, below). The new pass finds one stale-pointer defect of the kind the lens exists for: six Graduated rows in `bugs/INDEX.md` still name the `AGENTS.md` § Testing paragraph and bullets that this change deleted, so the mechanism they record points at nothing. Two low findings follow. Nothing is wrong in the scripts themselves.

## Round-1 answers, checked
1. Preamble fingerprint: `ls -lR --time-style=full-iso | md5sum` is compared in the EXIT trap, which exits 1 with a FAIL line. Run by hand from the checkout: rc 0 unchanged; a bad `KIT_ASSERT_GREP` stops with rc 2 before anything runs. Holds.
2. `-i`/`--all` only before the first word, `--` ends options: `-i -- --append-system-prompt` exits 0 with hits; `-i --all -- kit-preamble guide` exits 0. Holds.
3. Living set: the exclusions are as stated; `guide.md`, `skill-matrix.md` and the preamble come back as hits. The two stale pointers now read `guide.md` § 4 § Behavioural runs (`guide.md:156`, `skill-matrix.md:14`). `tier-4` has no living hit left. Holds.
4. `memory-loop.md:12` names `metadata.node_type: memory`. Holds.
5. The review skill's `--docs` row names the script and how the skill directory is found (see finding 3 for where). Holds.
6. The preamble rebuilds `sh-agent.real` every run unless `KIT_BIN` is set. Holds (see finding 2 for its working directory).

## Findings
- **[medium / confidence: high]** `.agents/docs/bugs/INDEX.md:134,136,137,138,140,141` — the Graduated rows BUG-083, 087, 090, 091, 106, 107 name `AGENTS.md` § Testing (the "A behavioural or sandbox run proves its setup before its output counts" paragraph and its bullets, the behavioural-run sentence) as the graduated mechanism. That text was removed from `AGENTS.md:68` in this range; the procedure now lives in `guide.md` § 4 § Behavioural runs and the setup proof in `testing/kit-preamble.sh`. The same diff touches this file (BUG-123 moved to Graduated, pointing at the new home) but left these six rows. A graduation row whose mechanism no longer exists reads as a standing guard. Fix: repoint them to `guide.md` § 4 § Behavioural runs (083: the behavioural-run bullet; 087/090/091/106/107: `testing/kit-preamble.sh` plus that section). Found by `restatement-grep.sh -i -- "§ Testing"`.
- **[low / confidence: high]** `.agents/docs/testing/guide.md:167` — says the preamble lists the real `~/.config/sh-agent` again on exit "with a warning when it changed". `kit-preamble.sh` (`kit_on_exit`) does not warn: it prints a FAIL line and exits 1, so the run does not count. The description understates the gate; a reader would treat a changed directory as advisory. Fix: "and fails the kit (exit 1) when it changed".
- **[low / confidence: medium]** `skills/review/SKILL.md:51` — the sentence "Before reading by hand, run `bash <skill dir>/scripts/restatement-grep.sh …`" sits inside the `Under --project:` quotation, so a plain `--docs` brief, the common case and this very review, never carries it. The purpose stated ("so the first round reads every copy") applies to both scopes. Fix: move the sentence into the main `--docs` quotation. Related, low confidence: `kit-preamble.sh` builds `./cmd/sh-agent` relative to the working directory and the header's Inputs do not say the kit must run from the harness checkout; from a clone's directory it exits 2 (probed). Add one line to the header or build with `-C`.

## Criteria
1. Dead code — n/a — docs lens
2. Flow tracing — n/a — docs lens (the preamble and the script were run by hand: pass, fail and usage paths behave as documented)
3. Standards conformance — n/a — docs lens
4. Directory layout — n/a — docs lens
5. Simplicity (DIY-except-security) — n/a — docs lens
6. Decoupling — n/a — docs lens

## Not examined
Historical `plans/` hits for "AGENTS.md § Testing" (ticked boxes, left as records). The other modified plans, specs (`delivery.md`, `dream.md`, `search.md`) and handoff were not re-read this round beyond round 1; the Antigravity side of the preamble was not exercised (no agy run). The `README.md` and `deploy/` guides were not re-searched for the removed paragraph beyond the script's hits, which showed none.

## Outcome

Every finding fixed in the phase (see round 3, `2026-10-01-review-c62e064-worktree-docs-r3.md`, APPROVE with zero findings): the preamble fingerprints sizes and mtimes and fails the kit; the grep script takes dash-led words, keeps the guide, matrix and preamble as living files and excludes records and vendored skills; the two stale pointers and the six Graduated rows name guide § 4 § Behavioural runs and the preamble; `metadata.node_type`; the review skill names the script in the plain --docs brief; the preamble builds from its own checkout.
