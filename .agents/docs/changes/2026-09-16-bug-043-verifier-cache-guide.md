---
owner: changelog
date: 2026-09-16
---
# BUG-043: verifier skipped the first-use guide and let a cached Go test result hide an e2e regression     (2026-09-16)

## Summary

The verifier agent wrote no testing guide despite twelve dispatches on Claude Code, because its condition for writing checked whether an `AGENTS.md` "describes" the test procedure — a subjective trigger that an `AGENTS.md` listing the commands satisfied. Separately, the verifier's counts run accepted a cached `ok … (cached)` result as a real run, hiding five failing e2e tests until `-count=1` forced a fresh result. Fixed in `agents/verifier.md`: the guide is now written whenever `testing/<app>/guide.md` does not exist as a file (not when another file "describes" it), and the counts run always carries `-count=1` (a verifier never takes a cached result).

## Type

bug fix

## Files changed

- `agents/verifier.md` — line 54: guide-writing condition changed from "a procedure that no guide describes" to "whenever `testing/<app>/guide.md` does not exist … an `AGENTS.md` that lists the commands does not replace it." Line 43: counts run always carries `-count=1` because "a verifier never takes a cached result … an `ok … (cached)` line is not a run."

## Agents dispatched

reviewer (`reviews/2026-09-16-review-bug-g16-verifier-cache-guide.md` APPROVE WITH NITS; medium finding on an undecidable conditional trigger fixed by deletion of the condition, as recorded in Outcome), verifier (`.agents/docs/testing/2026-09-16-verify-bug-043-verifier.md`, being written now), bugs (`.agents/docs/bugs/BUG-043-verifier-guide-and-test-cache.md`, already written), changelog.

## Verification

Verifier log `.agents/docs/testing/2026-09-16-verify-bug-043-verifier.md` (being written). Bug record `.agents/docs/bugs/BUG-043-verifier-guide-and-test-cache.md` names the sandbox repro path (no guide after twelve dispatches on Claude Code, `ok … (cached)` hid five failing e2e tests) and owed: a live verifier dispatch at the fixed definition, recorded in the simulation log's § Re-runs section.

## Docs

`agents/verifier.md` (the agent definition itself) changed; no other docs changed — the `context/AGENTS.md` step 3 reference to "the verifier creates [the guide] on first use" is consistent with the fix and required no update.

## Requirements

n/a
