---
owner: main
date: 2026-09-29
---
# After v1.12.1: the missing Antigravity run of `/init-project`, the gate addendum     (date: 2026-09-29)

## Summary
The v1.12.1 gate report's verdict paragraph said the Release rule's Antigravity condition was met by a Claude Code run; its own "What was not checked" said the opposite. The run now exists: `/init-project`'s update path on Antigravity from the real home at v1.12.1, against the linked-store fixture — PASS, 1 of 1 (`testing/2026-09-29-verify-real-home-v1.12.1.md`): the Step 3/4 guards kept the committed link, Step 5b stopped at the trust check, nothing moved, no settings file. The gate report carries an `## Addendum` naming the overstatement (BUG-105); the matrix row records both CLIs. The release needed no change. Trail only: no source path touched.

## Type
docs

## Files changed
`.agents/docs/testing/2026-09-29-verify-real-home-v1.12.1.md` (new), `.agents/docs/reviews/2026-09-29-preflight-v1.12.1.md` (addendum), `.agents/docs/testing/skill-matrix.md` (one row), `.agents/docs/bugs/BUG-105-*.md` and the index.

## Agents dispatched
- `bugs` — BUG-105 (the verdict line written before "What was not checked" was read against the rule).
- No `verifier`, `reviewer`, `security` or `documentation`: no source, interface or documented behaviour changed.

## Verification
The log above; the fixture's facts asserted before the CLI started; stderr empty, no `invalid allow entry`, exit 0.
