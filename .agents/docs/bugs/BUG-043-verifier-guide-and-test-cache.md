---
owner: bugs
date: 2026-09-16
id: BUG-043
status: fixed
area: agents/verifier
---
# BUG-043: verifier skipped the first-use guide and let a cached Go test result hide an e2e regression

## Symptom
Two symptoms, both from the Claude Code lane of the linkshelf simulation (`.agents/docs/testing/2026-09-16-linkshelf-simulation.md`, gaps G16 and G17):

(a) Across twelve verifier dispatches on Claude Code, no `testing/linkshelf/guide.md` was ever written, although `context/AGENTS.md` step 3 says the verifier creates it on first use. Antigravity's third dispatch, on the same project, did write one.

(b) The project's own run had `go test -race ./...` print `ok linkshelf/cmd/linkshelf (cached)` while a `-count=1` rerun showed five of six e2e tests failing. The e2e tests build the binary at run time, which Go's test cache cannot see. The project recorded this itself as its own BUG-002.

## Repro
Sandbox `/home/kesokaj/sandbox/runs/sim-claude-linkshelf`: `.agents/docs/testing/` holds four verifier logs and no guide file; `.agents/docs/bugs/BUG-002-go-test-cache-hid-e2e-regression.md` records symptom (b) from inside that project.

## Root cause
The old guide sentence read "a procedure that no guide describes," which an `AGENTS.md` that lists the test commands satisfies on a plain reading — so the verifier judged the guide already described and never wrote one. Separately, the counts-run instruction named `-race` and `-v` but said nothing about the test cache, so a plain `go test -race ./...` was accepted as a real run even when Go served a cached pass for a binary-building e2e suite.

## Fix
`agents/verifier.md` (uncommitted at dispatch time, committed with this record): the guide is now written whenever `testing/<app>/guide.md` does not exist as a file — "an `AGENTS.md` that lists the commands does not replace it." The counts run always carries `-count=1` — "a verifier never takes a cached result … an `ok … (cached)` line is not a run."

## Prevention
An agent's "first time" duty is keyed to whether a specific file exists, never to a judgement call about whether another file "describes" the same thing. A verifier's counted test run is never a cached one — always force `-count=1` on the counts command, in this repo and in any project's own verification guide.

## Related
`.agents/docs/reviews/2026-09-16-review-bug-g16-verifier-cache-guide.md` (APPROVE WITH NITS; its one medium finding — a conditional `-count=1` trigger being undecidable — was fixed by deleting the condition). `.agents/docs/testing/2026-09-16-verify-bug-043-verifier.md` (verification, in progress) plus a live verifier dispatch in the sandbox at the fixed definition, owed and tracked in the simulation log's Re-runs section. Same simulation as BUG-039..042.

## Update 2026-09-16: first fix incomplete, recurrence in the same file

**Symptom (b2):** real-home smoke check 3 (`.agents/scratch/realhome/results/claude-3/`, log at `.agents/scratch/realhome/claude-demo/.agents/docs/testing/2026-09-16-verify-greet-comment.md`) showed Claude Code's verifier, dispatched at the 7a6ea90 definition, still writing no `testing/<app>/guide.md` on a project that already held one older verification log — the sandbox's `verifier-43` row at the same definition did write it, so the bug reproduced only when older logs existed.

**Root cause (b2):** the 7a6ea90 fix reworded one sentence in `agents/verifier.md` § What to run but left the same duty restated, in the old ambiguous form, in four other places: `agents/verifier.md` § Records ("when the procedure is new or its steps changed"), `agents/documentation.md` § Your lane and its tree comment, `skills/init-project/SKILL.md` Step 4, and `context/AGENTS.md` step 3 ("created by the verifier on first use"). "The procedure is new" and "first use" still read as false once an older log exists, so an agent consulting any of those four locations instead of the one edited sentence could still skip the guide.

`Diagnosis missed:` the first fix grepped only for the literal string `guide.md` inside the verifier's § What to run; it did not grep the duty's *subject* ("first use", "first time it runs", "procedure is new") across `agents/`, `skills/`, and `context/`. That grep — recommended by BUG-027's prevention rule and named again by the retro's Proposal 5 — would have found all four restatements before the fix was called done.

**Fix (b2):** all five locations now state the duty unconditionally on the file's absence: `agents/verifier.md` § Records → "whenever it is absent or its steps changed"; `agents/documentation.md` § Your lane and tree comment → "whenever it runs a project's procedure and finds the guide absent"; `skills/init-project/SKILL.md` Step 4 → "the first time it runs the project's procedure and finds it absent"; `context/AGENTS.md` step 3 → "written by the verifier whenever it is absent". Reviewed in `.agents/docs/reviews/2026-09-16-review-bug-043-guide-wording.md` (REQUEST CHANGES → all three findings fixed, see its Outcome); confirmed by `grep -rnE 'first use|first time it runs|procedure is new' context/ agents/ skills/*/SKILL.md README.md` returning empty.

**Verification (b2):** sandbox rows `verifier-43b` on both CLIs and a repeat of real-home check 3 on Claude Code at the new wording, tracked in `.agents/docs/testing/2026-09-16-linkshelf-simulation.md` § Re-runs and `.agents/docs/testing/2026-09-16-verify-real-home-3.md`.

**Prevention (b2):** when a duty's wording changes, grep the whole repo for the duty's *subject* (the condition it keys off — "first use", "first time", "procedure is new" — not the old sentence's literal text) across every `agents/`, `skills/*/SKILL.md`, and `context/` file, not just the one file being edited. Same rule as BUG-027; this is its second recurrence (see BUG-003/007/019/027) — a strong candidate to graduate into a standing reviewer check per the retro's Proposal 5.

## Related (update)
`.agents/docs/reviews/2026-09-16-review-bug-043-guide-wording.md` (REQUEST CHANGES, all findings fixed). BUG-027 (same prevention rule, first recurrence). Retro `.agents/docs/retros/2026-09-16-retro.md` Proposal 5 (a grep-by-subject check for restated duties).
