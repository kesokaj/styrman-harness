---
owner: reviewer
date: 2026-09-30
---
# Review — worktree-retro-code   (2026-09-30)

Range reviewed: working tree — changed files: AGENTS.md agents/verifier.md skills/bug/SKILL.md skills/feature/SKILL.md skills/golang/SKILL.md skills/review-standards/SKILL.md; untracked, not ignored: none; read the diff with git diff HEAD -- <files>.

## Status: COMPLETE

## Verdict: APPROVE WITH NITS

Seven prose rule additions from `.agents/docs/retros/2026-09-30-retro.md` proposals 1–7. Each is grounded in a real bug record (BUG-109, BUG-110, BUG-111, BUG-112, and the two feature/verifier findings) and each names a check a reader can run. I found no rule that restates something already stated elsewhere in the harness, and no rule that adds a gate or step for the operator. The one real issue is a length contradiction inside `skills/feature/SKILL.md`: the new consumer-audit paragraph tells the scope restatement to carry "one line per consumer," but the same file's Step 1 opening sentence caps the restatement at "one or two sentences" — a feature with several consumers cannot satisfy both literally. Everything else is sound; two claims made in the AGENTS.md text (about `TestVendoredRowsCarryAPin` and the redaction needle shapes) were checked against the actual Go source and both hold.

## Findings

- **[low / confidence: medium]** `skills/feature/SKILL.md:19,25` — Step 1's opening rule fixes the scope restatement at "one or two sentences," but the new capability paragraph says the consumer list "goes in the scope restatement, one line per consumer." For a capability with three or four consumers this is no longer one or two sentences; the two rules are not reconciled (e.g. "…restatement, followed by one line per consumer" or an explicit carve-out the way the review-standards test-can-fail bullet marks its own exception). Low severity because both readings converge on the same practical outcome (list every consumer), but a strict reader hits a real conflict in adjacent sentences of the same step.

- **[low / confidence: low]** `skills/bug/SKILL.md:49` — the class-fix sentence sits inside "Step 4: Smallest fix," directly after "No adjacent refactors… no defensive code… no feature flags," and requires fixing "every site that has the defect in the same change" once the cause is a class. This is defensible (a class defect is the same bug at multiple sites, not a refactor) and is exactly what BUG-109's prevention rule asks for, but it does enlarge what "smallest fix" can mean without saying so explicitly next to the "no adjacent refactors" sentence it now sits beside. Not a correctness problem — flagging only because a fast reader could apply "smallest fix" to argue against fixing a second site the class-grep turns up, which is precisely the failure BUG-109 recorded.

No other contradictions found. The `agents/verifier.md` tree-under-test bullet narrows "the tree stays frozen" (global `context/AGENTS.md` § Definition of Done step 1) to source/tests/config rather than contradicting it — it resolves retro finding 5 (parallel reviewer/security dispatches writing under the docs root) without weakening the frozen-tree guarantee for source. The `skills/golang/SKILL.md` glob rule, the `skills/review-standards/SKILL.md` order-pin exception and consolidation-order bullet, and both `AGENTS.md` rules (vendor pin verification, redaction needle-by-kind) are each new ground with no prior restatement anywhere in `context/`, `skills/`, or `agents/` (checked by grep for their subject nouns: `Glob`, `consolidat`, `redact`/`needle`, vendored-pin language). Two of the AGENTS.md claims were checked against the live Go source rather than taken on faith: `TestVendoredRowsCarryAPin` (`internal/harness/harness_test.go:1829`) does refuse a missing repository or pin and does not check a pin's correctness, matching the rule's "not a wrong one" claim; `report.go`'s `prefixNeedles`/`keywordNeedles` (`internal/harness/report.go:36,52`) do split needles the way the new rule describes (vendor shape vs. distinctive prefix).

## Criteria

1. **Dead code** — pass. No functions, files, or exports touched; this is prose-only.
2. **Flow tracing** — pass with the one note above. Each new sentence's referenced check (`git archive | diff -r`, `git grep 'filepath.Glob'`, `git status --porcelain --untracked-files=all`, `git grep` for a class of call) is concrete and runnable; the two AGENTS.md claims about existing tests were spot-checked against the actual test code and hold.
3. **Standards conformance** — pass. No rule already stated elsewhere in `context/`, `skills/`, or `agents/` was left unchanged; grepped subject nouns (`Glob`, `redact`, `needle`, `consolidat`, `consumer`, `frozen`/`tree`) turned up no other restatement needing a matching edit.
4. **Directory layout** — n/a; no files added or moved.
5. **Simplicity (DIY-except-security)** — pass. No new dependency, no hand-rolled crypto/auth introduced or touched.
6. **Decoupling** — pass. Each rule lands in the file that owns its concern (Go pattern → `golang`; class-fix discipline → `bug`; test/consolidation-order discipline → `review-standards`; tree-freeze semantics → `verifier`; project-specific vendoring/redaction discipline → project `AGENTS.md`); no rule reaches across a skill boundary it doesn't own.

## Not examined

- The actual implementation these rules describe (`internal/harness/apply.go` `WriteAtomic`, `internal/harness/release.go`/`sweep.go`, `internal/harness/report.go` full needle list) was spot-checked only for the two claims above, not reviewed line by line — it is prior, already-shipped code (v1.13.1), not part of this diff.
- `skills/VENDORED.md` and `.agents/docs/bugs/INDEX.md` were read for context but are outside the six changed files and were not reviewed for their own correctness beyond confirming they don't contradict the new AGENTS.md text.
- No attempt was made to run `sh-agent search --refs` for `skills/golang/SKILL.md`, `skills/feature/SKILL.md`, or `skills/review-standards/SKILL.md` for `--kind bug` matches beyond what returned (searches for these three returned only pre-existing, unrelated bug references, already reflected above).

## Outcome
- 1 the consumer list vs the one-or-two-sentence restatement — fixed at `skills/feature/SKILL.md` Step 1: the consumers use the numbered-list form the same step already defines for a request over a list.
- 2 "smallest fix" vs "every site of the class" — fixed at `skills/bug/SKILL.md` Step 4: the sentence now says the smallest fix covers every site of the class and that this is not an adjacent refactor.
